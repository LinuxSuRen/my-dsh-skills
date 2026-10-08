# my-dsh-skills

[DeepSeek Harness (DSH)](https://www.npmjs.com/package/@deepseek-ai/dsh) 个人技能库——AI 编码助手按需加载的可复用任务指令集。

## 技能列表

| 技能 | 说明 |
|---|---|
| [`arch-dev-workflow`](arch-dev-workflow/SKILL.md) | 软硬件架构师通用工程方法论：方案先行与最小实现、Git 纪律、验证策略、接口与标识符设计原则、跨端契约与展示分离。适用于任何技术栈的研发任务。 |
| [`frontend-dev-workflow`](frontend-dev-workflow/SKILL.md) | Rick 的前端研发品味：Vue 3 + Vite 默认、按需上 TS/React、Element Plus 与手写设计令牌双轨、极简依赖、fetch 封装 API 层、自清理组件模式、中文文案规范。提炼自 linuxsuren 开源前端项目（open-pdf、onvif-ai、smart-chat 等）。 |
| [`git-identity`](git-identity/SKILL.md) | 提交代码必须使用 Rick 本人的 Git 身份（user.name=Rick, email=linuxsuren@users.noreply.github.com），禁止使用 agent/占位身份提交。 |
| [`branching-workflow`](branching-workflow/SKILL.md) | 每次提交/推送前必须先询问研发：新建分支还是当前分支、走 PR 还是直推；不得擅自决定分支策略。 |
| [`open-source-contribution`](open-source-contribution/SKILL.md) | 开源贡献规范：Issue/PR 提交要求、Review 礼仪、good-first-issue 模板、开源文档写作、仓库门面（About/README）完善、CLI 工具收录（hd-home）、开源仓库自动化（GitHub Actions）。提炼自 [open-source-best-practice](https://github.com/LinuxSuRen/open-source-best-practice)。 |
| [`android-device-debugging`](android-device-debugging/SKILL.md) | 安卓真机联调方法论：adb 无线调试配对、各厂商 ROM 陷阱（OPPO 冻结器与授权限制、华为 HDC、定制板 app-idle 强停）、相机资源战日志分析、ANR 定位、USB 隧道隔离网络故障、Android CI 构建静默失败根因清单。提炼自 tv-uvc-streamer 多设备实战（小米电视/工控板/OPPO/联想 Android 16）。 |
| [`onvif-rtsp-dev`](onvif-rtsp-dev/SKILL.md) | ONVIF 与 RTSP/RTP 音视频流开发方法论：WS-Discovery、SOAP 多 Profile、PTZ、对讲回传、TCP interleaved、SDP 多轨协商、H.264/AAC/G.711/JPEG(RFC 2435) 打包、ffmpeg 跨平台采集、子进程托管与录像分段、私有协议适配、ONVIF 设备网关代理、按需推流+MediaMTX/WHEP 分发与 Web 播放端（WHEP/MSE）实测陷阱。跨平台（Go/Java/ArkTS/TS 通用）。 |
| [`ohos-media-dev`](ohos-media-dev/SKILL.md) | 鸿蒙音视频应用开发方法论：Camera Kit 多摄并发、HCODEC 多实例与码率 int64 键、NV12 32 对齐、OHAudio 采集/播放、长时任务与熄屏、后台相机限制、HAL 超时自愈。 |
| [`android-camera-dev`](android-camera-dev/SKILL.md) | 安卓相机采集与推流开发方法论：Camera2 多摄并发降级、JPEG 直出/YUV 回退、DRI 归一化、JPEG_ORIENTATION 方向跟随、驱逐自愈、HAL 卡死防护、资源战 dumpsys 分析、拍照产物元数据无损嵌入。真机实测四类设备。 |
| [`mqtt-device-gateway-dev`](mqtt-device-gateway-dev/SKILL.md) | MQTT 设备控制面开发方法论（Go/paho + mosquitto）：共享连接 OnConnected 多播防订阅静默丢失、重连订阅恢复、RPC over MQTT 协议设计、retained 语义、dynsec 动态凭证（%u ACL 隔离/即时吊销/幂等收敛）、连接诊断分类、业务键去重白名单。源自机器人平台双端实战。 |
| [`edge-offline-pipeline-dev`](edge-offline-pipeline-dev/SKILL.md) | 边缘设备弱网数据管道开发方法论（Go）：目录即队列的崩溃安全磁盘重试队列（rename 原子状态机/文件名编码重试时间/边写边算 sha256/双重限额）、录像分段录制→remux→幂等上传、弱网事件上报不阻塞任务、离线 SQLite 任务队列、告警边沿触发。 |
| [`async-job-sse-upload-dev`](async-job-sse-upload-dev/SKILL.md) | Web 长任务与大文件传输方法论（Go 后端+浏览器前端）：202+job_id 异步任务框架、SSE 进度流（X-Accel-Buffering/双清理语义）、统一分片上传引擎（WriteAt 直写峰值磁盘 1×）与客户端（重试先对账/流式 SHA-256）、GB 级下载禁 fetch→blob、nginx/BFF 代理与容器文件替换陷阱。 |
| [`skill-refresh`](skill-refresh/SKILL.md) | 技能库保鲜：按周期（state.md 配置）调研工具/语言最佳实践，生成带来源的更新提案，研发确认后才修改与提交。 |
| [`skill-distillation`](skill-distillation/SKILL.md) | 任务后经验沉淀：复杂/高难度/耗时长的任务完成后总结可复用经验，区分通用技能与项目级技能，经研发确认后写入对应仓库。 |

## 通过 CLI 安装（askills）

本仓库提供 Go 编写的 CLI `askills`，把全部技能嵌入单个二进制，一条命令安装到各 AI 编码工具：

```bash
go run ./cmd/askills list               # 查看内置技能
go run ./cmd/askills install            # 全部技能 → ~/.agents/skills（dsh、codex、opencode 等共享）
go run ./cmd/askills install --tool dsh # → ~/.dsh/skills（--tool 支持 dsh / opencode / codex / all）
go run ./cmd/askills uninstall          # 移除已安装的技能
```

默认目标 `~/.agents/skills` 是 [Agent Skills](https://agentskills.io/) 开放标准的共享目录，dsh、codex、opencode 等 agents 兼容工具都会发现它；也可用 `--dir` 指定任意目录。安装是覆盖式、幂等的，重跑即为升级。

## 如何使用

每个技能是一个目录，内含带 YAML frontmatter 的 `SKILL.md`（`name` + `description` 为必填，`whenToUse` 可选）。DSH 的 skill 文件系统提供商会扫描各根目录并监听变化，添加或修改后无需重启。

### 方式一：customSkillDirs（推荐）

clone 本仓库后，在 DSH profile 的 `~/.dsh/profiles/web/cordis.patch.yml` 追加：

```yaml
- id: skill-filesystem
  name: '@deepseek-ai/dsh-skill-filesystem'
  config:
    customSkillDirs:
      - /path/to/my-dsh-skills
```

之后在任意会话中说“按 arch-dev-workflow 的规范做”，或让 AI 在匹配任务时自动调用。仓库更新后 `git pull` 即可同步最新技能。

### 方式二：软链到用户技能目录

不想改配置时，直接把技能目录软链进 DSH 用户级技能根：

```bash
git clone https://github.com/linuxsuren/my-dsh-skills ~/Workspace/github/linuxsuren/my-dsh-skills
ln -s ~/Workspace/github/linuxsuren/my-dsh-skills/arch-dev-workflow ~/.dsh/skills/arch-dev-workflow
```

### 方式三：复制进项目

把某个技能目录复制到项目的 `.dsh/skills/` 或 `.agents/skills/` 下，随仓库提交，团队成员打开会话即自动获得（`~/.agents/skills` 亦可被其他兼容 agents 约定的 AI 工具共享）。

## 如何贡献

1. 新建 `<skill-name>/SKILL.md`（kebab-case 命名），写入 frontmatter 与指令正文：
   ```yaml
   ---
   name: my-skill
   description: 一句话说明技能做什么。
   whenToUse: 什么场景下应使用本技能。
   ---
   ```
2. 正文写清楚可执行的操作规则，避免与具体仓库耦合——项目专属规范应放进各项目自己的 `.dsh/skills/`。
3. 提 PR，更新本 README 的技能列表。

## License

[MIT](LICENSE)