---
name: android-camera-dev
description: 安卓相机采集与推流开发方法论：Camera2 多摄并发限制与降级、JPEG 直出/YUV 回退、DRI 重启标记归一化、JPEG_ORIENTATION 方向跟随、相机被驱逐自愈、HAL 卡死防护、资源战 dumpsys 分析、拍照产物元数据无损嵌入（EXIF/XMP）。真机实测（小米电视/工控板/OPPO/联想 Android 16 四类设备）。
whenToUse: 开发安卓摄像头采集、MJPEG/RTSP/ONVIF 推流类应用（Camera2 API），排查绿屏/花屏、画面方向不对、多摄开不了、相机被系统或特权应用抢占、HAL 卡死 ANR、息屏断流等问题时使用。
---

# 安卓相机采集与推流开发方法论

## 1. 采集链路选择

- **JPEG 直出优先**（ImageReader + ImageFormat.JPEG，帧带 EXIF 可判来源）；设备不支持时回退 YUV_420_888 软压缩（YuvImage→NV21，注意 rowStride/pixelStride 语义因机型而异）
- **带 DRI 重启标记的 JPEG 帧必须重编码归一化**（标记扫描：SOS 前出现 0xFFDD 即中招）——RTP-JPEG 客户端普遍不支持重启标记，直发必绿屏；BitmapFactory→compress(Q85) 一并统一 Huffman 表、去 EXIF，代价约每帧 15-30ms
- 归一化判断要逐帧做（同设备不同分辨率模式可能不同），无 DRI 直通零开销

## 2. 多摄与并发

- `getConcurrentStreamingCameraIds()`（API 30+）未声明组合的 HAL，框架只允许一 app 同时开一路——第二路 openCamera 报 `Too many cameras already open`（REJECT）
- **降级策略**：全部相机按最大输出面积降序抢开，先开的高分辨率相机占住 ISP（即"默认分辨率最高"）；开启失败的路静默跳过，全败才算失败
- 用户选择：SharedPreferences 持久化相机编号（-1=自动），界面切换即管线重启
- ONVIF/状态页只广播**真实在推流**的相机，相机关闭即反注册；Profile token 与相机编号严格对应（profile_{N+1}↔N），切换不影响其他相机编号

## 3. 方向跟随

- 加速度传感器判定摆放：横放 0°/180°、竖放 90°/270°、**平放保持上次**（阈值：|a|>4 且主轴超副轴 2 倍）
- 用 `CaptureRequest.JPEG_ORIENTATION` 让 **HAL 直接转正像素**，快照/MJPEG/RTSP 下游零感知（不依赖 EXIF 方向标签——RTSP 生态普遍无视它）
- 公式：后摄 `(sensorOrientation - displayRotation + 360) % 360`；**前摄镜像相反** `(sensorOrientation + displayRotation) % 360`；前后摄 sensorOrientation 常见 90/270，必须 per-camera 计算
- 映射语义以后摄 orientation=90 真机校准（参考 ohos-ipcam-streamer OrientationSensor）

## 4. 稳定性（真机血泪）

- **openCamera 的 connectDevice binder 可能被卡死的 HAL 永久阻塞**（OPPO 实测主线程 ANR）→ 相机启动必须放工作线程
- 相机被高优先级客户端 **EVICT**（看家应用/shell 特权抢占）：onDisconnected/onError 即时自愈重开并重新入册；帧停滞 >6s 看门狗兜底；连续失败升级管线重启
- **管线重启走服务内存活方式**（startService 带 extra，onStartCommand 重建内部资源）——不能 stopService+延迟 startService，应用跌入 cached 态后定制 ROM 会静默丢弃后台启动
- 相机 HAL 整体卡死后（所有 app 都连不上相机）唯一恢复手段是**重启设备**
- 息屏/后台强停对抗：前台服务+常驻通知是底线；ROM 级冻结（OPPO Hans）代码无解需用户加白

## 5. 诊断命令

```bash
# 相机资源战事件流：CONNECT/DISCONNECT/EVICT/REJECT 语义判抢占与并发限制
adb shell dumpsys media.camera | grep -E "CONNECT|DISCONNECT|EVICT|REJECT|concurrent"
# 主线程 ANR 栈（卡 connectDevice = HAL 卡死实锤）
adb shell dumpsys dropbox --print data_app_anr
```

- **JPEG 结构解剖是花屏/绿屏定位第一步**（python 遍历 marker）：SOF 尺寸与采样因子、有无 DRI、DHT 表数——对比好帧/坏帧的段差异直接锁定根因
- 像素级验证：解码后采样统计绿像素占比（G>90 且 G>1.5×max(R,B)），比肉眼客观

## 6. 拍照产物元数据无损嵌入

- JPEG 走 EXIF APP1（ImageDescription）+ XMP APP1 双通道，已存在的段一律不动；两槽都被相机占用时降级写 COM 段；PNG 用 eXIf chunk + IHDR 后的 tEXt
- **元数据缺失绝不能让拍照失败**：任何不支持/畸形输入都返回原图不报错——元数据是增强项不是依赖项，采集链路对它零容忍度反过来（元数据写入失败 ≠ 拍照失败）
