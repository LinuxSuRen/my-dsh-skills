---
name: robot-cameras
description: 机器人（机器狗等）与网络摄像机（IPC）RTSP 取流地址速查与探测方法：智身科技 M1、宇泛灵猫 Cyvet、海康/大华/宇视等品牌格式、宇树 Go 系列机载网络等。
whenToUse: 为机器人/机器狗设备做视频取流、录像、推流开发或联调时使用；接入新型号摄像头时按探测流程验证并补充本技能。
---

# 机器人摄像头取流

## 已验证型号

### 智身科技 M1（机器狗）

前后各一个本体摄像头，RTSP 端口 8554，路径按方位命名，IP 固定 `192.168.168.168`：

- 前：`rtsp://192.168.168.168:8554/front`
- 后：`rtsp://192.168.168.168:8554/back`

### 宇泛灵猫 Cyvet

按通道号取流，`stream=0` 为主码流，设备 IP 视部署而定：

- `rtsp://<设备IP>:554/live?channel=<通道号>&stream=0`

## 网络摄像机（IPC）品牌 RTSP 格式

### 海康威视 Hikvision（实测）

- 新平台 IPC/NVR：`rtsp://<user>:<pass>@<ip>:554/Streaming/Channels/<通道ID>`，**通道ID = 通道号×10 + 码流号**（101=通道1主码流、102=通道1子码流、201=通道2主码流）；
- 经典旧格式 IPC：`rtsp://<ip>:554/h264/ch1/main/av_stream`；
- 实测：`/Streaming/Channels/102`（通道1子码流）用于录像场景。

### 大华 Dahua（公开资料，未实测）

- `rtsp://<user>:<pass>@<ip>:554/cam/realmonitor?channel=<通道号>&subtype=<0主|1子>`

### 其他常见品牌（公开资料，未实测）

- 宇视 Uniview：`rtsp://<ip>:554/media/video1`（主）、`/media/video2`（子）；
- 水星 Mercury / TP-LINK：`rtsp://<user>:<pass>@<ip>:554/stream1`（主）、`/stream2`（子）；
- 雄迈 Xiongmai（XM）：`rtsp://<ip>:554/user=admin&password=&channel=1&stream=0.sdp?real_stream`；部分设备需在网络设置中关闭再重开 RTSP 才生效；
- 通用 ONVIF 设备：优先走 WS-Discovery + GetStreamUri 动态取流，按 通道/码流 参数适配多厂商，避免硬编码品牌 URL。

## 网络资料（未实测，以官方文档为准）

### 宇树（Unitree）Go1 / Go2

- 机载网络固定为 `192.168.123.x`：机载计算单元 `.161`，上位机常用 `.162`；
- 官方默认以 WebRTC 提供低延迟视频流；RTSP 地址随固件与版本有差异，**使用前以官方文档或实测为准，不凭记忆填地址**。

## 探测方法（新型号接入流程）

1. 确认设备 IP 与本机同网段（`ping` / `arp`）；
2. 端口未知时优先尝试常见 RTSP 端口：554、8554；
3. 用 ffprobe 验证流可解析：`ffprobe -v error -show_streams "rtsp://<ip>:<port>/<path>"`；
4. 用 ffplay 人工确认画面与延迟：`ffplay -fflags nobuffer "rtsp://..."`；
5. 验证通过后，把型号、URL 模板与验证命令补充进本技能，走技能修订流程（提案→研发确认→PR）。

与 `onvif-rtsp-dev`、`android-camera-dev` 技能联动：ONVIF 设备先走发现流程；安卓平台摄像头走 adb 相关技能。
