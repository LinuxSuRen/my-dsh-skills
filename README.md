# my-dsh-skills

[DeepSeek Harness (DSH)](https://www.npmjs.com/package/@deepseek-ai/dsh) 个人技能库——AI 编码助手按需加载的可复用任务指令集。

## 技能列表

| 技能 | 说明 |
|---|---|
| [`arch-dev-workflow`](arch-dev-workflow/SKILL.md) | 软硬件架构师通用工程方法论：方案先行与最小实现、Git 纪律、验证策略、接口与标识符设计原则、跨端契约与展示分离。适用于任何技术栈的研发任务。 |
| [`frontend-dev-workflow`](frontend-dev-workflow/SKILL.md) | Rick 的前端研发品味：Vue 3 + Vite 默认、按需上 TS/React、Element Plus 与手写设计令牌双轨、极简依赖、fetch 封装 API 层、自清理组件模式、中文文案规范。提炼自 linuxsuren 开源前端项目（open-pdf、onvif-ai、smart-chat 等）。 |
| [`git-identity`](git-identity/SKILL.md) | 提交代码必须使用 Rick 本人的 Git 身份（user.name=Rick, email=linuxsuren@users.noreply.github.com），禁止使用 agent/占位身份提交。 |
| [`branching-workflow`](branching-workflow/SKILL.md) | 每次提交/推送前必须先询问研发：新建分支还是当前分支、走 PR 还是直推；不得擅自决定分支策略。agent 只负责 review PR 并加 comment（建议定期/自动化 review），合并由人工执行或收到明确指令。 |
| [`open-source-contribution`](open-source-contribution/SKILL.md) | 开源贡献规范：Issue/PR 提交要求、Review 礼仪、good-first-issue 模板、开源文档写作、仓库门面（About/README）完善、CLI 工具收录（hd-home）、开源仓库自动化（GitHub Actions）。提炼自 [open-source-best-practice](https://github.com/LinuxSuRen/open-source-best-practice)。 |
| [`android-device-debugging`](android-device-debugging/SKILL.md) | 安卓真机联调方法论：adb 无线调试配对、各厂商 ROM 陷阱（OPPO 冻结器与授权限制、华为 HDC、定制板 app-idle 强停）、相机资源战日志分析、ANR 定位、USB 隧道隔离网络故障、Android CI 构建静默失败根因清单。提炼自 tv-uvc-streamer 多设备实战（小米电视/工控板/OPPO/联想 Android 16）。 |
| [`onvif-rtsp-dev`](onvif-rtsp-dev/SKILL.md) | ONVIF 与 RTSP/RTP 音视频流开发方法论：WS-Discovery、SOAP 多 Profile、PTZ、对讲回传、TCP interleaved、SDP 多轨协商、H.264/AAC/G.711/JPEG(RFC 2435) 打包、ffmpeg 跨平台采集、子进程托管与录像分段、私有协议适配、ONVIF 设备网关代理、按需推流+MediaMTX/WHEP 分发与 Web 播放端（WHEP/MSE）实测陷阱。跨平台（Go/Java/ArkTS/TS 通用）。 |
| [`ohos-media-dev`](ohos-media-dev/SKILL.md) | 鸿蒙音视频应用开发方法论：Camera Kit 多摄并发、HCODEC 多实例与码率 int64 键、NV12 32 对齐、OHAudio 采集/播放、长时任务与熄屏、后台相机限制、HAL 超时自愈。 |
| [`android-camera-dev`](android-camera-dev/SKILL.md) | 安卓相机采集与推流开发方法论：Camera2 多摄并发降级、JPEG 直出/YUV 回退、DRI 归一化、JPEG_ORIENTATION 方向跟随、驱逐自愈、HAL 卡死防护、资源战 dumpsys 分析、拍照产物元数据无损嵌入。真机实测四类设备。 |
| [`mqtt-device-gateway-dev`](mqtt-device-gateway-dev/SKILL.md) | MQTT 设备控制面开发方法论（Go/paho + mosquitto）：共享连接 OnConnected 多播防订阅静默丢失、重连订阅恢复、RPC over MQTT 协议设计、retained 语义、dynsec 动态凭证（%u ACL 隔离/即时吊销/幂等收敛）、连接诊断分类、业务键去重白名单。源自设备控制面双端实战。 |
| [`edge-offline-pipeline-dev`](edge-offline-pipeline-dev/SKILL.md) | 边缘设备弱网数据管道开发方法论（Go）：目录即队列的崩溃安全磁盘重试队列（rename 原子状态机/文件名编码重试时间/边写边算 sha256/双重限额）、录像分段录制→remux→幂等上传、弱网事件上报不阻塞任务、离线 SQLite 任务队列、告警边沿触发。 |
| [`async-job-sse-upload-dev`](async-job-sse-upload-dev/SKILL.md) | Web 长任务与大文件传输方法论（Go 后端+浏览器前端）：202+job_id 异步任务框架、SSE 进度流（X-Accel-Buffering/双清理语义）、统一分片上传引擎（WriteAt 直写峰值磁盘 1×）与客户端（重试先对账/流式 SHA-256）、GB 级下载禁 fetch→blob、nginx/BFF 代理与容器文件替换陷阱。 |
| [`skill-refresh`](skill-refresh/SKILL.md) | 技能库保鲜：按周期（state.md 配置）调研工具/语言最佳实践，生成带来源的更新提案，研发确认后才修改与提交。 |
| [`subagent-parallel-dev`](subagent-parallel-dev/SKILL.md) | 多子代理并行开发大型跨仓库任务：契约先行、早期契约审查、独立复测纪律（实跑复核+stash/纯净worktree对照证实既有失败）、真实服务E2E兜底、有进展才提交。提炼自 6 仓库 7 子代理实战。 |
| [`redis-stream-rpc`](redis-stream-rpc/SKILL.md) | Redis Streams 版 RPC 总线：注册/心跳/命令/回执流布局、entry 双编码、幂等键、PEL 回收(XAUTOCLAIM)、权重优先窗口仲裁、令牌随流下发与恢复重派。与 mqtt-device-gateway-dev 互补。 |
| [`skill-distillation`](skill-distillation/SKILL.md) | 任务后经验沉淀：复杂/高难度/耗时长的任务完成后总结可复用经验，区分通用技能与项目级技能，经研发确认后写入对应仓库。 |
| [`robot-cameras`](robot-cameras/SKILL.md) | 机器人（机器狗等）与网络摄像机（IPC）RTSP 取流速查与探测：智身科技 M1、宇泛灵猫 Cyvet、海康（实测通道ID规则）/大华/宇视/水星/TP-LINK/雄迈品牌格式、宇树 Go 系列机载网络；含新型号 ffprobe/ffplay 探测接入流程。 |
| [`installer-bundle-design`](installer-bundle-design/SKILL.md) | 安装包设计方法论：清单驱动（manifest 自描述、字段只增不改）、单一产物与编排收敛/按 profile 裁剪、白名单自动收集编排引用、在线/离线双模同源版本语义、对外端口全参数化。适用于任何打包→目标机安装类交付工具。 |
| [`installer-target-execution`](installer-target-execution/SKILL.md) | 目标机执行层方法论：执行通道抽象（Run/Upload/Platform）、提权矩阵（root/sudo -n/sudo -S/docker 组，stdin 凭据安全与双行喂入）、SFTP/md5 增量传输、重装幂等、docker 组授权语义、Linux/macOS/Windows 三平台适配清单。 |
| [`installer-preflight`](installer-preflight/SKILL.md) | 安装前置校验与故障预防：分层校验（建连前→上传前→拉镜像前）、架构/端口预检、错误定位到元凶并带下一步动作、Docker 高频陷阱表（containerd save 多架构、inspect 架构误报、bind 源变目录、半参数化端口、拼接式地址）、证书固定名输入、shim 冒烟验证策略。 |
| [`ui-responsive-list`](ui-responsive-list/SKILL.md) | 列表页 UI 自适应规范：筛选区「放得下才显示 + 折叠展开多行」（ResizeObserver 实测宽度、纯函数折叠计算、隐藏控件保持挂载）与表格/卡片双视图按视口定默认（窄屏卡片、宽屏表格，用户偏好优先）。提炼自多套管理端 Web 列表页改造实战。 |
| [`user-facing-errors`](user-facing-errors/SKILL.md) | 用户可见错误规范：最终用户只看到稳定错误码+普通人能读懂的提示，技术根因只进日志，绝不进对用户的响应。 |
| [`patent-disclosure`](patent-disclosure/SKILL.md) | 中国专利交底书：发明/实用新型/外观设计的专利点挖掘、轻量查新与成文，首篇定稿后可做保护型 1+N 专利布局（围栏）。含填表、线稿、CAD STEP 解析、公式 OMML、mermaid 出图与 Word 定稿工具链。 |
| [`patent-application`](patent-application/SKILL.md) | 专利申请文件撰写：把已有交底改写为权利要求书、说明书、摘要与黑白附图四件套；含材料门禁、标号核对、支持性审查、权要审计工具。须显式点名并指定交底目录。 |
| [`patent-docket`](patent-docket/SKILL.md) | 专利案卷会稿：单 agent 协调者，按发明人/工程师给的材料一趟串起交底书与申请文件；立案、分诊、派工、记轮次，申请问题清单上的缺口最多改三轮。须显式点名。 |
| [`patent-search`](patent-search/SKILL.md) | 中国专利公布公告著录检索：发明人/申请人/分类号/名称/摘要高级查询（epub.cnipa.gov.cn），可从单图或权利要求抽关键字生成检索式；结果落 Markdown 报告并给机读 JSON。 |
| [`patent-reader`](patent-reader/SKILL.md) | 专利通俗解读：公开号/PDF/全文读成通俗笔记、Canvas 图谱与术语网，可写入 Obsidian 库；独权拆成稳定特征行（claim_features.json）作为对照表输入。 |
| [`patent-chart`](patent-chart/SKILL.md) | 权利要求对照表（claim chart）：独权拆成技术特征逐格比对对比文件/产品/标准，附证据与强弱；用于无效对照、FTO 初筛、侵权/EoU、SEP。只出底稿，不出法律意见。须显式点名。 |
| [`patent-map`](patent-map/SKILL.md) | 专利地图：基于已解读入库的 Obsidian 案例库摊开五种图（语义地形、申请人重叠、文内引证等），bge-small-zh 向量本地渲染。须显式点名且需已有解读入库。 |
| [`patent-oa`](patent-oa/SKILL.md) | 审查答复（OA）辅助：审查意见问答与意见陈述草稿，案例库支撑检索，库薄时引导案例入库；用户确认采纳后才出 Word 陈述。须显式点名。 |
| [`patent-exam-policy`](patent-exam-policy/SKILL.md) | 专利政策简报：对照国知局近期审查口径，说明对交底写法/本稿的影响并提示申请书式；技能进化仅为旁路，默认只改交底包。须显式点名。 |

以上 9 个 `patent-*` 技能源自 [handsomestWei/patent-disclosure-skill](https://github.com/handsomestWei/patent-disclosure-skill)（MIT，© handsomestWei），经适配融入本仓库：原 `skills/<name>/` 布局扁平化为顶层目录，包间路径引用与依赖清单同步改写，各目录内保留原作者 [LICENSE](patent-disclosure/LICENSE)。交底主链路需要 Python 3.9+，在 `patent-disclosure/` 下执行 `pip install -r requirements.txt`（python-docx、latex2mathml、playwright 等）；其余可选依赖见各包 `tools/README.md`，完整使用文档见上游 [INSTALL.md](https://github.com/handsomestWei/patent-disclosure-skill/blob/main/INSTALL.md)。

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