---
name: web-launch-dx
description: Web 应用启动体验规范：后端监听端口支持漂移（被占时向后探测，默认开启、可 strict 关闭），漂移后的实际地址传导到日志、对外宣告 URL 与前端 dev 代理；服务就绪后自动打开实际地址（默认开启、可关闭）。
whenToUse: 开发或改造任何 web 应用或带 HTTP 服务的工具时——新增服务监听、配置端口、处理端口占用冲突、决定启动后是否自动打开浏览器、接前端 dev 代理时。
---

# Web 应用启动体验：端口漂移与自动打开

同时开发多个 web 应用，8080/3000/5173 这类常用端口必然相撞。两条硬规则：**端口漂移**与**启动后自动打开**——都默认开启，都提供开关。

## 1. 端口漂移：默认开，可关

**行为定义**：

- 配置端口被占用时向后逐个探测（`port+1`、`port+2`…），直到可监听；设探测上限（如 +20），上限内全部占用才报错退出。
- 仅「地址被占」（EADDRINUSE）触发漂移；权限不足、地址不可用等其他错误直接报错，不盲目换口。
- 跨平台判定：unix 用 `EADDRINUSE`；Windows 用 `WSAEADDRINUSE`（数值 10048）判定，常量仅存在于 Windows syscall 包，跨平台编译用数值。

**开关**：默认开启。提供 strict 关闭开关——env（如 `PORT_DRIFT=off`）或配置项（如 `server.port_drift: false`）；关闭后端口被占即启动失败，报错按 `user-facing-errors` 给普通人能懂的提示 + 占用端口排查命令（如 `lsof -i :8080`）。

**范围**：HTTP 是主场景；RTSP/UDP 等附属监听同款漂移（UDP 全被占可降级，如 tcp-only），上限内全占用才失败。

参考实现（[local-onvif-adapter](https://github.com/linuxsuren/local-onvif-adapter) 的 `listenWithDrift`）：

```go
ln, port, drifted, err := listenWithDrift(cfg.HTTPAddr, maxPortDrift, logger)
if drifted {
    logger.Warn("http port drifted", "from", cfg.HTTPAddr, "to", ln.Addr().String())
}
```

## 2. 漂移结果必须传导（否则等于没漂）

1. **日志显式宣告**：漂移必须 Warn 打出 from→to；最终监听地址无论是否漂移都打印。
2. **先定端口，再算地址**：对外宣告的 URL（回调、注册端点、流地址、WebSocket）一律用实际端口计算——漂移解析必须发生在任何地址计算之前。
3. **前端 dev 代理不能写死后端端口**：分离模式（vite dev + 后端）下，proxy target 硬编码 `127.0.0.1:8080` 会因后端漂移而断链（反例：open-cal-hub 的 vite.config）；用 env 注入实际端口（启动脚本写 `.env.development`、vite 读 `VITE_API_PORT` 组装 proxy target），或 dev 模式下显式关闭后端漂移并写进 README。
4. **前端 dev server 自身**：Vite 默认即漂移（`strictPort: false`），不要随意改 `true`。

## 3. 启动后自动打开：默认开，可关

- 服务**就绪后**打开浏览器到实际地址（漂移后的端口）；就绪判定用「端口可连」或首个 200，避免打开即连接被拒。
- 开关：默认开启；env（如 `OPEN_BROWSER=off`）或配置项关闭；服务器/容器/headless 场景显式关闭，或检测到无显示环境时自动跳过。
- 跨平台打开命令：darwin `open <url>`、linux `xdg-open <url>`、windows `rundll32 url.dll,FileProtocolHandler <url>`；打开失败只 Warn 不阻断启动。

## 4. 落地清单

1. 启动入口：`listenWithDrift` + strict 开关 + 漂移 Warn 日志
2. 地址计算全部后置于实际端口确定之后
3. 自动打开：就绪后触发 + 开关 + 跨平台命令 + 失败不阻断
4. env/配置命名随项目惯例，README 写明默认值与关闭方法
5. dev 代理目标改为可注入，杜绝写死端口

## 衔接

- `user-facing-errors`：端口占用/全占用失败的报错给稳定提示 + 排查命令，技术根因进日志
- `frontend-dev-workflow`：vite 侧保持默认漂移；env 注入实践落在前端项目结构里
- `admin-credential-bootstrap`：首启打印的访问地址同样必须是漂移后的实际地址
