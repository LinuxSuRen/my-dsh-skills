---
name: onvif-rtsp-dev
description: ONVIF 与 RTSP/RTP 音视频流开发方法论:WS-Discovery 发现、SOAP 端点与多 Profile、PTZ 变焦与能力区分、对讲回传 backchannel、RTSP TCP interleaved 服务端、SDP 多轨协商、H.264/AAC/G.711 打包与时间戳规则、ffmpeg 跨平台设备采集与屏幕推流、私有协议适配、端口漂移与嵌入架构,及乱码/饿死/乱序/灰色画面等实测陷阱。跨平台(Go/Java/ArkTS 通用)。
whenToUse: 开发或调试 ONVIF 设备端/客户端(摄像头、NVR、网关、发现工具)与 RTSP 推拉流(服务端、拉流客户端、对讲回传),实现 ffmpeg 跨平台设备采集或屏幕推流,适配私有云台协议,排查设备发现不到、Profile/取流地址错误、PTZ 不工作、拉流花屏/变速/丢包告警、屏幕采集灰色、设备枚举为空、SOAP 应答乱码等问题时使用。
---

# ONVIF 与 RTSP/RTP 音视频流开发方法论

## 1. ONVIF 发现层(WS-Discovery)

- 组播组 `239.255.255.250:3702`,客户端发 Probe,设备回 ProbeMatches(**单播回源**,带 RelatesTo 对应 Probe 的 MessageID)
- 设备上线主动发 Hello、下线发 Bye(客户端可零 Probe 感知设备)
- 应答必须带 `AppSequence(InstanceId/MessageNumber)`,Scopes 空格分隔:类型 `onvif://www.onvif.org/type/NetworkVideoTransmitter`、名称/硬件带宿主机型便于多设备区分
- **Android 等系统需持有组播锁**(WifiManager.MulticastLock)才能收到 Probe,WiFi 驱动默认过滤组播
- 调试:python socket 发自定义 Probe 即可;应答 XML 必须按 **UTF-8** 编码(Latin-1 截断会让中文设备名变乱码,`Content-Length` 也按 UTF-8 字节数)

## 2. ONVIF SOAP 服务端

- 单一 device_service 端点即可承载全部操作(GetProfiles/GetStreamUri/GetCapabilities/GetServices/GetDeviceInformation/GetSystemDateAndTime/PTZ 系列按 Body 动作名分发),客户端兼容任何路径
- 按动作名匹配 Body 用 **tag 级正则**(任意命名空间前缀)防误撞——`Stop` 这类短词干会撞中其他动作
- 属性值引号单双都兼容;数值解析拿不到就回退默认,**不要静默吞掉**
- GetCapabilities/GetServices 决定客户端如何发现 PTZ/Media XAddr;GetDeviceInformation 返回真实机型/序列号(权限受限时生成本机持久 UUID 兜底)

## 3. ONVIF Profile 与取流

- 一个 Profile = 一路输出:VideoSourceConfiguration + VideoEncoderConfiguration + (可选)AudioEncoderConfiguration + (可选)PTZConfiguration
- **客户端以 Profile 内 PTZConfiguration 元素的存在判定 PTZ 可用**——漏了它,再强的 PTZ 实现客户端也不显示控制
- token 命名:`profile_N`、`ptz_N` 保持同序号;GetStreamUri 带 ProfileToken 回来时按尾号路由到对应码流/相机
- GetStreamUri 返回 `rtsp://ip:port/camN`;GetSnapshotUri 返回 HTTP 单帧地址(不支持 MJPEG-RTSP 的客户端降级轮询)
- Profile 声明的分辨率/帧率/码率与实际输出一致(动态生成,别硬编码)

## 4. ONVIF PTZ 与能力区分

- Zoom 属 PTZ 服务:ContinuousMove(速度 ±1)、AbsoluteMove(位置)、RelativeMove、Stop、GetStatus(当前位置)、GetConfigurationOptions(能力)
- **能力按轴分开上报**:GetConfigurationOptions 的 Spaces 里 PanTilt 与 Zoom 独立——无云台设备只报 ZoomSpaces,客户端据此隐藏方向键
- 位置空间 [0,1] 归一,设备内部映射物理量(如 zoomRatio min–max 线性);**自定义空间合法**:另报自有 URI 空间携带真实范围(如倍率 0.55–100),客户端即可显示 "2.5x" 而非百分比
- ContinuousMove 用定时器步进模拟(如 200ms 一拍全范围 4%×速度),Stop/超时双保险
- 注意命令 token 元素名:`GetConfigurationOptions` 用 **ConfigurationToken**,其余多用 ProfileToken——解析两者都兼容,否则多设备能力不同时会答错

## 5. RTSP 服务端骨架(TCP interleaved)

- 方法机:OPTIONS(Public 列表)→ DESCRIBE(SDP)→ SETUP(Transport 应答)→ PLAY → TEARDOWN
- **CSeq 必须原样回显**,不匹配 ffmpeg 类客户端直接断连
- 只支持 TCP 时 UDP SETUP 回 **461 Unsupported Transport**(客户端自动改 TCP,比不回应友好)
- PLAY 应答带 RTP-Info;新客户端从**关键帧**起播(waitKey 丢 P 帧直到 IDR),否则花屏
- URL 即路由:`/camN` 对应第 N 路源,DESCRIBE 与 SETUP 都从中解析相机号

## 6. SDP 多轨协商

- 每媒体一条 `m=` 段,`a=control:trackID=N` 作 SETUP 句柄——**客户端 SETUP 哪条轨就收哪条**(音/视频按需协商)
- H.264:`a=rtpmap:96 H264/90000` + fmtp 带 `packetization-mode=1;sprop-parameter-sets=<b64SPS>,<b64PPS>`(从码流拆 SPS/PPS 动态生成)
- AAC(RFC 3640):`MPEG4-GENERIC/<rate>/<ch>` + fmtp `mode=AAC-hbr;sizelength=13;indexlength=3;indexdeltalength=3;config=<ASC>`
  - ASC 两字节:OT(5bit)|freqIdx(4bit)|ch(4bit);AAC-LC 44100 立体声=0x1210、48000=0x1190
- G.711:静态类型 PCMA=8 / PCMU=0,`PCMA/8000/1`
- 对讲回传轨:`a=sendonly` 标记(gortsplib 系客户端以此识别 backchannel),G.711 PCMA/8000 是事实标准

## 7. RTP 打包

- **H.264(RFC 6184)**:单 NAL ≤MTU 直封装;超长分片 FU-A,marker 只在帧末片;去起始码
- **JPEG(RFC 2435)**:主头 8B = type-specific(1)|fragment-offset(3)|type(1)|quality(1)|width/8(1)|height/8(1);熵数据按字节分片,marker 仅末片
  - type 0=4:2:2、1=4:2:0,由 **SOF 首分量采样因子**决定(0x21→0,0x22→1);宽高同样从 SOF 动态解析,**禁止写死分辨率**
  - quality=255 时首片附量化表头(MBZ 2B+表长 2B)+64 字节裸表序列;提取 DQT 时去掉每表 1 字节表 ID
  - **含 DRI 重启标记的帧不可直接打包**:type 64+ 附 Restart Marker Header 的规范做法 ffmpeg rtpdec_jpeg 不支持(Unimplemented,直接丢包)——部分相机 HAL 默认输出带 DRI 的 JPEG(OPPO 实测),采集层需重编码归一化(BitmapFactory→compress 顺带统一 Huffman、去 EXIF)
- **AAC(RFC 3640)**:载荷 = AU-headers-length(2B,bit 计)=0x0010 + AU header(13bit 尺寸<<3 + 3bit index)+ AAC 裸帧;每 AU 固定 1024 样本做时间戳步进
- **G.711**:1 字节=1 样本;20ms 一包(8k 下 160 字节)
- 时间戳:视频 90kHz、音频按采样率时钟;**整条流(含参数集)严格单调递增**(同毫秒多包强制 +1),否则客户端报 non-monotonic DTS;同一接入单元(SPS/PPS/IDR)共享同一时间戳

## 8. TCP interleaved 双向

- 下行:`$` + channel(1B) + length(2B) + RTP,与 RTSP 文本同连接;**发送必须串行化**(promise 链/队列),并发写致 TCP 流内乱序,客户端报巨额丢包
- 上行(对讲):控制连接混入二进制 `$` 帧——接收解析先按字节识别 `$` 头完整取出,再回落文本 RTSP 解析
- 应答与 RTP 同 socket 时互斥(同一发送链);关闭连接:应答入队 flush 后再关,同步 close 会截断响应

## 9. 多客户端与多路

- 每客户端独立 session/seq/ssrc/时间戳基准;分发按 session 订阅的轨/相机过滤
- 关键帧前补发 SPS/PPS(编码器 IDR 常不带),SDP sprop 仅兜底
- 多路:一路一条 SDP/URL,与 ONVIF Profile 序号对齐

## 10. 实测陷阱速查

| 现象 | 根因 |
|---|---|
| 客户端发现不到设备 | 组播锁未持有 / ProbeMatches 没回源 / AppSequence 缺失 |
| 中文设备名乱码 | SOAP 按 Latin-1 编码截断 |
| PTZ 控制不显示 | Profile 缺 PTZConfiguration 元素 |
| 无云台设备仍显示方向键 | 客户端未按 Spaces 区分 PanTilt/Zoom |
| 变焦 200 但无效果 | token 元素名/引号解析失败被静默跳过 |
| ffmpeg "65533 packets lost" | 并发写 socket 致 RTP 乱序 → 发送串行化 |
| 拉流花屏/开头黑屏 | 未 waitKey 从 IDR 起播 |
| non-monotonic DTS 告警 | 同毫秒多包时间戳重复 |
| 音频变速变调 | 采样率写死,未按 SDP/ASC 动态 |
| AAC 解码噪声 | AU-headers-length 缺失或 ASC 错 |
| OpenCV 拉不到流 | 默认 UDP,应回 461 引导 TCP |
| 应答空响应 | 响应未 flush 就 close |
| macOS 屏幕采集灰色/挂起 | 未签名 ffmpeg 无权限,avfoundation 静默拒绝(改用 screencapture 管道) |
| Windows 枚举为空(FFmpeg 7.1+) | dshow 输出格式变更,旧解析器靠分节标题失明(需双格式兼容) |
| 快照失败(设备被占) | 取流进程独占设备,快照需从 RTSP 流拉帧而非直连 |
| 私有协议命令无效 | 缺校验和/速度编码错误(带符号补码 vs 无符号)/角度编码位数错误 |
| 管道推流 pipe EOF | StdoutPipe 与 cmd.Run 竞态(改用 os.Pipe + defer 清理) |
| 屏幕推流帧率极低 | Retina 全分辨率 PNG 3MB/帧管道拥塞(需 scale 缩小) |
| 端口漂移后客户端连不上 | 宣告地址未跟随实际端口(GetStreamUri/UI/Discovery 全部用实际端口) |
| Go 路由注册 panic | ServeMux 方法限定模式与全方法模式路径重叠(GET / vs /onvif/) |

## 11. ffmpeg 跨平台设备采集

- **设备枚举按平台分发**:Linux 扫 `/dev/video*`+sysfs(名称/连接方式/过滤 UVC 元数据节点 index!=0);macOS/Windows 解析 `ffmpeg -list_devices`
- **FFmpeg 7.1+ 统一设备列表**:dshow 不再打印分节标题,改为每行末尾 `(video)`/`(video, audio)`/`(audio)` 媒体类型标记——解析器必须**双格式兼容**(旧格式靠分节,新格式靠行尾括号组);正则限定小写媒体类型词,设备名含括号不会误判
- **macOS avfoundation**:新版在设备名后追加 `[uid:...] [serial:...]`,需剥离;`Capture screen N` 是屏幕设备,非摄像头
- **设备独占访问**:Windows(dshow)/Linux(v4l2) 采集设备通常只能被一个进程打开——取流进程持有设备时,**快照必须从 RTSP 流拉帧**(不能二次直连设备);macOS 允许多进程并发访问摄像头
- **私有协议适配(以 Skydroid TP 为例)**:命令末尾追加 2 位十六进制累加校验(整条含 `#` 的 ASCII 求和低 8 位);速度为带符号补码 ±100(0x64=+100 右/上,0x9C=-100 左/下);角度为 16 位补码 `%04X`(值=角度×100,范围 ±90°);UDP 通道无 ACK 需 3 次重发;姿态回报帧 `#tpUG2rGAC<yaw><pitch><roll>`(各 4 位 hex,÷100)

## 12. 屏幕采集推流

- **macOS 权限陷阱**:未签名 ffmpeg 的 avfoundation 屏幕采集被 macOS **静默拒绝**——不弹权限框,返回灰色帧或挂起(screencapture 系统命令可正常采集验证)
- **macOS 替代方案**:`screencapture -x -t png /tmp/frame.png` 循环 + `cat` 到管道 + ffmpeg `-f image2pipe -i pipe:` 读取,完全绕开签名权限
- Retina 全分辨率 PNG 每帧 ~3MB,管道吞吐不足——输出端加 `-vf scale=1280:-2` 降分辨率保帧率;macOS 屏幕源**不加 `-nostdin`**(需从 stdin 管道读帧)
- Go 管道连接用 `os.Pipe()` 而非 `StdoutPipe`(后者与 `cmd.Run` 的 Wait 竞态导致 pipe 过早关闭 → EOF)
- Linux 屏幕采集用 `-f x11grab -i :0.0`(一屏一 X display);Windows 用 `-f gdigrab -i desktop`

## 13. 嵌入式 RTSP 服务与端口漂移

- **内嵌 RTSP 服务器**(gortsplib 等):Publisher(ANNOUNCE/RECORD) → ServerStream → 多 Reader(DESCRIBE/PLAY) 分发;读者在发布者未就绪时**等待**(默认 5s,而非立刻 404);发布者断开自动释放路径,新发布者可踢旧接管
- **端口漂移**:HTTP/RTSP 端口被占时自动向后尝试(EADDRINUSE → next port,最多 20 个);**对外宣告的地址必须用实际端口**(GetStreamUri/GetSnapshotUri/WS-Discovery XAddr/UI 展示全部跟随),漂移结果不写配置,下次启动仍从首选端口探测
- UDP 传输端口(RTP 8000/RTCP 8001)被占时降级为仅 TCP 模式,不影响拉流
- Go 1.22+ ServeMux 模式冲突:`GET /`(方法限定+更宽路径)与 `/onvif/`(全方法+更窄路径)不能共存——根路由去掉方法前缀

## 14. 快照与 PTZ 架构

- **快照优先走流**:启用中的摄像头从 RTSP 服务拉帧(避免二次独占打开设备);失败回退直连;禁用摄像头直接直连;短 TTL 缓存(3s)防进程风暴
- **PTZ mock 架构**:无硬件时用虚拟状态机(绝对命令记录 + 连续移动按速度虚拟积分,满速 90°/s yaw/45°/s pitch);有姿态回报(Skydroid GAC)则用真实角度替代估算
- **焦距(Focus)属 Imaging 服务而非 PTZ**:GetOptions 声明 AutoFocusMode/DefaultSpeed/NearLimit/FarLimit;GetMoveOptions 返回 Absolute/Relative/Continuous 各自范围;Move 操作支持三种模式;PTZ GetServiceCapabilities 用 `Zoom="true"` 属性声明
- 预置位实现:自定义位存本地角度回放;内建位(回中/垂直向下/Follow/Lock/FPV 模式)直接发硬件命令
