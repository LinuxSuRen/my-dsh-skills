---
name: android-device-debugging
description: 安卓真机联调方法论：adb 无线调试配对、各厂商 ROM 陷阱（OPPO 冻结器与授权限制、华为 HDC、定制板 app-idle 强停）、相机资源战日志分析、ANR 定位、USB 隧道隔离网络故障、Android CI 构建静默失败根因清单。
whenToUse: 在安卓真机/平板/工控板/电视上部署调试应用、排查设备连不上或装不上、分析系统级强停/冻结/相机抢占、配置或修复 Android CI 构建时使用。
---

# 安卓真机联调方法论

## 1. 设备接入

### USB 直连排查顺序
1. `system_profiler SPUSBDataType`（macOS）确认物理枚举——连 USB 总线都看不到就是线/口/模式问题，与 adb 无关
2. `adb devices` 看 unauthorized（手机上点授权弹窗）还是 offline
3. 充电线陷阱：大量 USB-C 线是纯充电线；扩展坞的 Billboard Device 是"协商失败"特征，数据不过
4. 设备重启后 OPPO 系会把 USB 模式重置为「仅充电」——通知栏切「传输文件」

### 无线调试（Android 11+，配对码流程）
```bash
adb mdns services          # 自动发现，无需手抄 IP（含 _adb-tls-pairing._tcp）
adb pair <ip>:<配对端口> <6位码>
adb connect <ip>:<连接端口>  # 配对端口与连接端口是两个不同端口
```
- adb 31.0.2 起支持 pair 与 mdns；配对码有时效，出现即立刻配
- 重启 adb server（kill-server）会丢掉已连的 TCP transport，需要重连
- 5555 端口（经典 adb tcpip）与无线调试端口并存时，无线调试挂了先试 5555

## 2. 设备识别陷阱

| 现象 | 结论 | 对策 |
|---|---|---|
| USB 显示「HDC Device」（PID 0x12d1） | 华为系用 HDC 协议，非 adb | 换设备或用华为 hdc 工具 |
| HarmonyOS NEXT（纯血鸿蒙 5.x） | 不支持安卓 APK，hdc 只装 HAP | 提前问清系统版本，别浪费时间 |
| OPPO force-stop 后 `am start` 报 `Cannot make calls to a recycled instance` | ColorOS 内部 bug | 用 `monkey -p <pkg> -c android.intent.category.LAUNCHER 1` 启动 |
| OPPO `pm grant` 报无 GRANT_RUNTIME_PERMISSIONS / `cmd appops` 报无 MANAGE_APP_OPS_MODES | ROM 限制 shell 权限 | 界面内弹窗授权，或换思路 |

## 3. ROM 后台强杀与对抗

### 症状定位命令
```bash
logcat -d | grep -iE "Stopping service|app idle"     # app idle 强停（实测 60 秒即杀）
logcat -d | grep -iE "freeze|hans"                    # OPPO Hans 冻结器
logcat -d | grep "startForeground.*not allowed"       # 前台服务被 bg restriction 拒绝
```

### 对抗组合拳（按层次）
1. 前台服务 + 常驻通知（标准）
2. `cmd appops set <pkg> RUN_IN_BACKGROUND allow`、`RUN_ANY_IN_BACKGROUND allow`（部分 ROM 拒绝 shell 执行）
3. `cmd deviceidle whitelist +<pkg>`（doze 白名单）
4. `am set-standby-bucket <pkg> active`
5. 息屏冻结（OPPO 实测 `reason: traffic`）：用户侧关 WiFi 省电/加电池白名单，代码层无解

### 服务重启的隐藏坑
- `stopService() + 延迟 startService()`：应用跌入 cached 态后，定制 ROM 会**静默丢弃**后台 startService（服务记录归零、无异常日志）——需要重启采集时应改用**服务内存活管线重启**（服务不拆，onStartCommand 带 extra 重建内部资源）
- targetSdk < 26 豁免后台 startService 限制；但 Android 16 把最低可安装 targetSdk 提到 24（24 仍低于全部行为断点，零行为损失）

## 4. 相机资源战（dumpsys media.camera）

```bash
adb shell dumpsys media.camera | grep -E "CONNECT|DISCONNECT|EVICT|REJECT|concurrent"
```

| 事件 | 含义 |
|---|---|
| `REJECT ... Too many cameras already open` | HAL 不支持并发相机（无 concurrent 组合声明），或应用同时开两路 |
| `EVICT ... held by <pkg> (score N)` | 被更高优先级客户端驱逐（score 越高优先级越高，shell/系统应用无敌） |
| `CONNECT/DISCONNECT` 高频循环 | 系统看家/监控类应用周期性抢相机 |

- 相机被 EVICT 后应用收到 onDisconnected/onError——**采集框架必须自愈重开**，否则僵尸
- ANR 排查：`dumpsys dropbox --print data_app_anr` 拿主线程栈；主线程卡在 `ICameraService.connectDevice` = 相机 HAL 卡死，相机启动逻辑必须放工作线程

## 5. 网络故障隔离

- **USB 隧道法**：`adb forward tcp:18090 tcp:8090` 后 curl localhost——绕开 WiFi 验证应用层。一次实测：WiFi 路径快照卡死，USB 隧道 40MB/s，实锤 ROM WiFi 省电黑洞
- 应用层缓解：`WifiLock` 用 `WIFI_MODE_FULL_LOW_LATENCY`（API 29+）——`WIFI_MODE_FULL` 在 Android 10+ 已被系统无视
- TCP 卡在初始拥塞窗口（~5KB 后不动）= 大包黑洞特征；DF ping 全挡不一定是 MTU 问题（ROM 可能防火墙 ICMP）

## 6. Android CI 构建静默失败清单

每个都是"本地绿 CI 红"且无有效报错，按踩坑顺序：

1. **`gradle.properties` 里提交了代理配置**（`systemProp.http.proxyHost=127.0.0.1`）：CI 上 Gradle 把依赖请求发往 runner 不存在的代理，连接拒绝被静默记为「仓库无此构件」——本地能过是因为代理活着。本地代理放 `~/.gradle/gradle.properties`，永不入库
2. **vendored 代码被上游 .gitignore 排除**：第三方源码里生成的头文件（如 libusb 的 `version_describe.h`）在本地是历史构建残留，CI 检出即编译失败——`git add -f` 强制入库静态内容，构建才自洽
3. **AGP 构建期静默自动安装 NDK 失败**（`InstallFailedException`，600MB 无输出）：流水线里用 sdkmanager 显式预装 + `build.gradle` 钉 `ndkVersion`
4. `android-actions/setup-android@v3` 默认装的 `tools` 包已从 SDK 仓库下架 → 直接用 runner 预装 SDK，自己接受 license
5. 「`Resolved plugin` 成功后 0.1 秒报 not found」且零网络请求：先查配置层（代理/镜像/仓库过滤），别急着换 JDK、换 Gradle、换 action——这些都不是元凶但每个都像

## 7. 快速核对

- `getprop ro.product.model / ro.build.version.release / ro.product.cpu.abi`：机型三件套
- `dumpsys window | grep mCurrentFocus`：当前焦点界面（验证弹框/界面状态）
- `uiautomator dump` + bounds 解析：无源码 UI 自动化点击
- 像素级图像验证（绿屏/花屏）：PIL 采样统计绿像素占比，比肉眼可靠
