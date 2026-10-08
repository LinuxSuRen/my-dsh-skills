---
name: mqtt-device-gateway-dev
description: MQTT 设备控制面开发方法论（Go/paho + mosquitto）:共享连接 OnConnected 多播防订阅静默丢失、重连窗口订阅恢复、client-id 唯一、RPC over MQTT 协议设计(request_id 幂等/deadline/控制指令不 retain)、retained 配置分发与零长度清理标记、dynsec 动态凭证(%u ACL 隔离/删 client 即吊销/幂等收敛/首启引导死锁)、TCP 探测连接诊断分类、业务键去重的合法重发白名单。源自机器人平台 broker 侧+设备侧双端实战。
whenToUse: 开发或排查 Go 服务与边缘设备经 MQTT 互联的控制面(设备网关、RPC 下发、状态上报、配置分发、mosquitto dynsec 凭证管理),尤其共享 paho 连接多模块、弱网重连后消息丢失或 RPC 恒超时、broker ACL 静默拒绝、凭证吊销/互踢、同键消息被去重丢弃等问题时使用。
---

# MQTT 设备控制面开发方法论

## 1. paho 客户端生存参数

- AutoReconnect+ConnectRetry+ResumeSubs;`CleanSession(true)`+`ResumeSubs(false)`(broker 侧订阅不比进程活得久,避免消息积压在无 handler 的会话);`SetOrderMatters(false)` 解除回调串行(慢 handler 不阻塞入站投递)
- SUBACK 等待加超时上限(如 30s),防"连接成功但永不 ACK 订阅"的 goroutine 泄漏
- **重连窗口期 Subscribe 会被拒**("not currently connected..."/"reconnecting state and cleansession is true"):自维护订阅表,被拒也照样记住订阅,OnConnect 回调统一重放
- **client-id 必须全局唯一**:同一 broker 上两条连接(如事件总线+心跳)共用 client-id 会互相顶替,broker 反复踢旧连接(session takeover)
- `ssl://` broker 无 CA 配置时启动即报错,不允许静默降级明文
- 连接失败错误要带定位信息:broker 地址、client_id、超时值、常见原因清单(不可达/DNS/TLS/认证被拒)——ConnectRetry 静默重试期间看不到真实原因,超时才暴露

## 2. 共享连接的 OnConnected 必须多播

- 多模块共享一条 paho 连接时,**单槽 OnConnected 会被后注册者覆盖先注册者**——broker 重启/踢线重连后部分订阅静默丢失:表现为无关模块超时(如 RPC 请求可达、响应无人接收恒超时,极难联想到订阅丢失)
- 对策:回调改多播列表(只增不减),通知时拷贝快照遍历;补"重连后全部重订阅"回归测试

## 3. RPC over MQTT 协议

- topic 设计:`{prefix}/{id}/rpc/request`(下行,QoS1,不 retain)+`{prefix}/{id}/rpc/response/{request_id}`(上行);平台侧通配订阅 `{prefix}/+/rpc/response/+` 一次订阅路由全部在途请求(避免每次调用 Subscribe/Unsubscribe),pending map+request_id 关联响应
- **控制指令绝不 retain**:机器人离线期间的请求由调用方按 deadline 超时,语义等价调用失败;retain 会让指令残留在 broker、设备上线后补发
- request_id 幂等:QoS1 重放同一请求只执行一次(去重表按 request_id)
- deadline_unix/timeout_ms 双字段+最大上限 cap,过期请求直接回 TIMEOUT 不执行;响应发布用 `context.WithoutCancel`+短超时(不被调用方取消打断)
- 未知 method 回 METHOD_UNSUPPORTED 而非超时——版本漂移可观测,调用方能区分"设备旧版"与"设备失联"
- 调用方超时取 ctx deadline 与配置较小者;超时错误带 topic 与排查提示(确认设备标识一致且设备在线)

## 4. topic 与 retained 语义

- 配置/服务发现用 retained(如 platform/info 广播平台地址):设备连上即取、零配置;等待超时回退环境变量
- per-entity topic(如 task/status/{runCode})+通配符 "+" 订阅:retained 消息不互相覆盖,可按实体单独清理
- 新旧 payload 格式共存时**换新 topic**(如 status→status_full),避免订阅方按新格式解析旧消息失败
- **零长度 payload 是 retained 清理标记**,broker 会投给订阅者:handler 须静默跳过,否则解码告警噪音
- 双端 topic 契约变更需两侧同批发版;订阅 topic 与对端发布 topic 不同步=静默收不到(无报错)

## 5. mosquitto dynsec 设备凭证

- `$CONTROL/dynamic-security/v1` JSON 协议,correlationData 关联请求/应答(未回带时按"唯一在途请求"兜底)
- 设备隔离用 pattern ACL `{prefix}/%u/#`(%u=用户名=设备标识)一设备一角色;外部第三方设备角色仅允许发布自己的 `external/{type}/%u/#`
- **删除/禁用 client 即时断开在线连接=即时吊销**;一设备一令牌(唯一索引),否则第二张令牌覆盖前者密码
- **幂等收敛式管理**:ensureRole 只补缺失 ACL(acltype+topic 匹配)不删既有条目——手工授权不受影响,平台新增必要 ACL 能收敛到存量部署;ensureClient 缺失建、存在则收敛密码(支持轮换)+补挂角色;启动后台重试 EnsureDefaults 直到成功,不阻塞服务启动

## 6. mosquitto 部署硬坑

- **dynsec 首启引导死锁**:插件在 dynamic-security.json 缺失时自动生成随机密码 admin(伴随 .pw 文件),应用拿约定凭证被 not authorised 拒且无法预知随机值——entrypoint 检测自动生成态( json 缺失或 .pw 存在)后用约定凭证重新引导;json 0644(mosquitto 降权后可读,内容仅哈希)
- 认证用户的全局 topic 规则必须用 **pattern** 关键字:写在 user 声明之前的 topic 规则只对匿名客户端生效,allow_anonymous false 下等于没写
- mTLS 模式下**所有 topic 操作都查 ACL**:服务发现类 retained topic(如 platform/info)漏配被静默拒绝,平台地址发不出去

## 7. 连接诊断分类

- 连接超时先做 TCP DialTimeout 探测:TCP 通=TLS/mTLS 握手或认证失败(报错带"去下载证书"类修复指引);TCP 不通=网络/防火墙
- 报错信息按上述分类组织而非裸透传,现场第一眼可定位

## 8. 业务键去重的边界

- 按业务键(如 taskID)去重防 QoS1 重放时,**必须识别"合法的同键重发"**(重试/恢复/续跑):如充电暂停后恢复会复用同一 taskID 重派发,被当重复丢弃=任务永久卡死
- 对策:合法重发带显式标记(如 is_charge_resume)绕过业务级去重;去重分事件级(eventID)与业务级(实体ID)两层,业务级放行白名单

## 9. 配置热更新

- 订阅 refresh topic 通知重拉:去抖窗口(如 1s)合并多次通知;重拉异步执行不阻塞 MQTT 回调 goroutine;重拉前后快照对比打印新增/移除/修改(凭据只记有无)
- 部分更新保护:upsert 事件不带 devices 时保留现有列表,不覆盖清空

## 10. 实测陷阱速查

| 现象 | 根因 |
|---|---|
| 重连后 RPC 请求可达、响应恒超时 | OnConnected 单槽覆盖,部分订阅静默丢失 |
| 重连窗口期 Subscribe 报 not currently connected | paho 重连态拒订阅;自维护订阅表 OnConnect 重放 |
| 两条连接反复互踢 | client-id 相同,broker session takeover |
| 设备上线后收到旧指令 | 控制指令用了 retained,broker 残留补发 |
| 恢复执行的任务被静默丢弃、永久卡死 | 重派发复用同 taskID 撞上业务级去重 |
| 订阅方持续 decode 告警噪音 | 零长度 retained 清理标记未静默跳过 |
| 应用连 broker 一直 not authorised | dynsec 首启自动生成随机密码,需 entrypoint 重新引导 |
| 认证用户收不到全局 topic 消息 | ACL 写在 user 声明前只对匿名生效,应使用 pattern |
| mTLS 下 retained 平台地址发不出 | 所有 topic 都查 ACL,服务发现 topic 漏配静默拒绝 |
| topic 改名后对端缓存收不到数据 | 双端契约不同步,需同批发版 |
| 手工加的 ACL 被平台管理覆盖丢失 | 角色管理非幂等收敛(删了重建),应只补缺失不删既有 |
