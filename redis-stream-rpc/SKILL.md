---
name: redis-stream-rpc
description: Redis Streams 版 RPC 总线模式：注册/心跳/命令/回执流布局、entry 按字段区分双编码、幂等键、consumer group 命名、PEL 回收（XAUTOCLAIM）、窗口仲裁与令牌随流下发。与 mqtt-device-gateway-dev（MQTT 侧）互补，提炼自机器人平台算法服务总线实战。
whenToUse: 平台需要与一组可插拔服务（算法引擎/插件/worker）经消息总线解耦通信、要求服务自注册与审批纳管、且已部署 Redis 时。
---

# Redis Streams 服务总线 RPC 模式

## 1. 流布局（四类流 + 结果总线）

```
{domain}.service.register     服务→平台  注册（低频、保留全史审计、XGroupCreateMkStream 起始 0）
{domain}.service.heartbeat    服务→平台  心跳（高频、MaxLen ~10000 截断、平铺字段）
{domain}.service.{sid}.command 平台→服务 命令（entry 双编码见 §2）
{domain}.service.{sid}.reply   服务→平台 同步回执（pending map + request_id + deadline）
既有结果总线                    服务→平台 业务结果（沿用平台既有流，消费组零改动）
```

- 消费组命名：命令流 `{sid}-command-group`（consumer=instance_id，同组多实例天然竞争消费=水平扩容）；平台回执流一个组多流订阅。
- 服务上线/恢复：平台幂等重发 `auth.grant`（令牌+平台地址）+ `models.sync`；**注册即触发**，不要依赖 offline→online 翻转（重启快于心跳阈值时不会翻转——实战踩坑）。

## 2. entry 按字段区分双编码（同流共存）

- JSON 命令：`payload`(JSON 字符串) + `command` + `request_id` + `deadline_unix`
- 二进制信封（如 protobuf）：`event`(bytes) + `event_id/task_id/type` 排查字段
- 心跳是平铺 map，注册/结果只有一个 `payload` 字段——三类解析路径分开写。

## 3. 可靠性语义

- **幂等**：request_id 稳定可重试（业务键派生，如 `det-{runCode}-{materialId}`）；同 request_id+service_id 去重。
- **deadline**：信封必带 `deadline_unix`，过期回 `EXPIRED` 不执行（人工兜底类服务可豁免）。
- **PEL 回收（必做）**：消费循环周期性 `XAUTOCLAIM <stream> <group> <consumer> 60000 0-0` 认领超时未 ACK 的滞留消息，走与正常消息相同处理——否则"消费到一半崩溃"的消息永久滞留。
- 服务不可用：平台先查审批态+在线态**快速失败**（不等超时）；派发失败的业务子集发失败事件闭环，不让上游挂起。
- **恢复重派**：SERVICE_OFFLINE 失败记录进内存重派队列（requestId+原始载荷+TTL），服务恢复回调时重算目标 fan-out，幂等键保证不重复。

## 4. 多服务仲裁（可插拔策略）

- 同一业务键允许 fan-out 多服务、采纳一路：`DispatchStrategy`（选目标）与 `ArbitrationStrategy`（采纳）**两个独立接口** + 按名注册表 + 配置选择。
- 默认权重优先语义：最高权重回 completed 短路定稿；窗口内（可配，如 5s）只有低权重则到期降级；全失败/超时合成失败事件闭环；**迟到高权重不覆盖已采纳**（防报告/告警抖动）仅存参考；仲裁器在事件入口与既有合并链路之间，对下游透明。
- 仲裁表内存 + deadline 即可（平台重启由业务超时兜底），`Sweep`/`Submit` 导出可同步驱动便于表格驱动测试。

## 5. 平台侧数据模型

- 服务表（审批态×在线态两维字段分离、权重为平台侧管理属性）；"做什么"是平台实体（任务表），"怎么做"由服务注册自声明（引擎能力：模态/输入形式/参数 schema）——服务无业务状态，重装后 sync 拉齐自愈。
- 二进制不走消息：文件经平台暂存/资产库，命令只带 URL 引用，方向恒为"服务按令牌主动访问平台"。
