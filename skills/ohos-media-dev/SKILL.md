---
name: ohos-media-dev
description: 鸿蒙音视频应用开发方法论:Camera Kit 多摄并发(并发类型打开)、HCODEC 多实例与 buffer 输入、码率 int64 键、NV12 32 对齐、OHAudio 采集/播放、长时任务与熄屏、后台相机限制、HAL 超时与自愈,均为真机实测陷阱。
whenToUse: 开发 HarmonyOS/OpenHarmony 上的摄像头采集、硬件编解码、麦克风/扬声器、RTSP 推流类应用,或排查编码码率失控、画面色度损坏、编码器饿死、熄屏断流、后台冻结、相机打不开等问题时使用。
---

# 鸿蒙(HarmonyOS)音视频应用开发方法论

## 1. 相机(Camera Kit)

- **多摄并发**:每颗相机独立 CameraInput+Session;必须查询 `getCameraConcurrentInfos` 并对**所有相机**(含第一颗)用其返回的并发类型 `input.open(type)` 打开——首颗若用普通 `open()` 独占,后续相机报 7400102 Device conflict
- 并发不足的设备单路失败自动跳过,全部失败才算失败(降级而非报错)
- preview profile 按"最接近目标"挑选;实际尺寸可能≠请求值,C 侧缓冲按**实际 profile 尺寸**分配(用两阶段:先建 surface,session 配置前 setFormat(实际宽高))
- 采集 buffer 平面布局实测差异:planeCount=3 且 plane[1].offset > plane[2].offset 时是 **NV21**(拷贝需交换 UV);rowStride/columnStride 语义以真机为准,首个 buffer dump 一次确认
- 会话打开偶发 7400103(Session not config),重试一次即过;**HAL 调用全部加超时保护**(stop/release/close 个别永不回调,不设超时会卡死状态机;深度休眠下给到 8s)

## 2. 视频硬编(HCODEC)

- **多实例可用**(双路 720p 实测),但**不支持进程内二次配置同一实例**——重配置走"销毁→重建"或退出应用重启;方向/分辨率切换用"保存配置→terminateSelf→重开自动恢复"模式
- **码率是 int64 键**:必须 `OH_AVFormat_SetLongValue(OH_MD_KEY_BITRATE, ...)`,SetIntValue 写入编码器读不到,实测码率失控到 14Mbps;`BITRATE_MODE_CBR` 配套
- **输出宽 32 字节对齐**:720 宽(非 32 对齐)硬编实测底部大面积色度损坏(绿块);旋转后裁剪到对齐(704x1280)即恢复
- **buffer 输入模式陷阱**:Start 后**不可清空 need-input 队列**——codec 启动即投递缓冲,清空会"泄漏"已投递缓冲(codec 视为未归还不再投递),该路编码器永久饿死(push-enter 有 have=0 特征);重建路径的队列清理放在销毁之后
- 编码帧回调经 threadsafe function 送 ArkTS,帧携带相机编号

## 3. 音频(OHAudio / OH_AudioEncoder)

- 采集:OH_AudioStreamBuilder(capturer)回调 OnReadData 是**实时线程**,只做拷贝入队;**AAC 编码器在该机型 CreateByMime("audio/aac") 可能失败**,回退 `CreateByName("avenc_aac")`
- 音频编码器仅提供 legacy AVMemory 回调(SetCallback + onNeedInputData/onNeedOutputData),写入用 `OH_AVMemory_GetAddr`,提交用新式 `PushInputData(index+attr)`——两者混用是官方示例做法;编码器输出裸 AAC 帧(无 ADTS)
- G.711 旁路:采集 48kHz(48000/8000=6 整数倍)→ 混单声道 → **FIR 低通(33 阶 sinc,截止 3.4kHz)再 6:1 抽取**;裸均值抽取无抗混叠,高频折叠回语音频段听感沙哑
- 播放(对讲):renderer 用**回调拉取式**(SetRendererWriteDataCallback),环形缓冲存收到的 PCMA、回调里 A-law 解码填充,无数据填静音;满则丢最旧保实时

## 4. 后台与熄屏

- 长时任务必须**双管齐下**:`backgroundTaskManager.startBackgroundRunning(DATA_TRANSFER)` + module.json5 里 EntryAbility 声明 `"backgroundModes": ["dataTransfer"]`——缺声明注册失败(9800005 bgMode invalid),**熄屏即被冻结**(freezer cgroup 可查)
- **后台相机是系统隐私限制**:应用退后台后 `CreateCameraInput` 返回 null/会话被回收,长时任务也无法豁免——不做强行支持,动态适配:后台暂停重试(避免相机设备半开楔死 7400201,该状态仅进程退出可解除),回前台自动重建管线恢复
- 帧停滞看门狗分级自愈:单次停滞重开会话 → 连续多次全量重建管线(采集面+编码器+会话,与停止/启动同路径可重复执行)

## 5. 工程杂项

- 权限:相机/麦克风都是 user_grant,UIAbility 里 `requestPermissionsFromUser` 一次申请;麦克风不可用时功能降级而非失败
- hilog **不支持 %f**,浮点日志乘 100 取整用 %d,否则输出乱码
- napi 垫片模式:C 层管采集/旋转/编码/音频全链,ArkTS 只做服务与 UI;按相机多实例(实例数组),全局单例改多实例时注意 listener context 指向各自实例
- 系统返回键:`onBackPress` 返回 true 拦截 + `moveAbilityToBackground` 进后台;退出用显式关闭按钮(停止推流+长时任务+terminateSelf)

## 6. 实测陷阱速查

| 现象 | 根因 |
|---|---|
| 码率失控 14Mbps | BITRATE 用 SetIntValue(int64 键读不到)|
| 竖屏底部绿块 | 输出宽 720 非 32 对齐,硬编色度损坏 |
| 单路编码器零输出 | Start 后清了 need-input 队列,codec 缓冲泄漏饿死 |
| 第二颗相机打不开(7400102) | 首颗未用并发类型打开,独占设备 |
| 熄屏流断/进程冻结 | module.json5 缺 backgroundModes 声明(9800005)|
| 后台回来相机全失效 | 系统回收后台相机会话 → 看门狗重建 |
| 相机反复 7400201 | 后台期间反复重开致设备半开,仅进程退出解除 |
| AAC 编码器创建失败 | CreateByMime 不支持,回退 CreateByName(avenc_aac) |
| G.711 人声沙哑 | 抽取无抗混叠滤波 |
