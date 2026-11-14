---
name: async-job-sse-upload-dev
description: Web 长任务与大文件传输开发方法论(Go 后端+浏览器前端):202+job_id 异步任务框架(类型级并发上限、取消与失败区分、panic 隔离、过期清理)、SSE 进度流实现(X-Accel-Buffering、变化才推、REST 兜底)与前端双清理语义、统一分片上传引擎(WriteAt 直写峰值磁盘 1×、乱序断点幂等、LimitReader 防粘连)与客户端(重试先对账、流式 SHA-256、服务端口径进度)、GB 级下载禁 fetch→blob、nginx 代理缓冲与容器文件替换陷阱。
whenToUse: 开发备份/导入/生成/导出等耗时操作的 API 与前端进度展示、大文件(百 MB~GB 级)上传下载、仪表盘实时统计推送(SSE 选型与实现),或排查 SSE 无事件/200 后 body 永久 pending、上传重复入库或内存暴涨、分片内容错位、断点续传失效、dev 代理截断长连接、bind mount EBUSY 等问题时使用。
---

# Web 长任务与大文件传输方法论

## 1. 长任务 API 标准形态

- POST 业务接口→**202+job_id**(不在 HTTP 请求里跑耗时操作);GET /jobs/{id} 一次性查询兜底(SSE 断了还能查);DELETE /jobs/{id} 取消
- per-type 并发上限(超限 409+明确文案"已有 N 个 X 任务在运行");进度回调 update(Progress{Stage,Current,Total})自动算百分比封顶 100
- ctx 取消→**cancelled 与 failed 区分**(用户取消≠执行失败);safeCall recover panic 转"任务内部错误"(一个任务 panic 不带崩进程);终态过期清理(如 10min)

## 2. SSE 实现要点(服务端)

- 三件套响应头:`Content-Type: text/event-stream`+`Cache-Control: no-cache`+**`X-Accel-Buffering: no`**(经 nginx 必须,否则响应被代理缓冲,事件到不了客户端)
- 连接即推当前快照(防首帧空白);内部轮询(如 500ms)状态/百分比**变化才写帧**,终态推完即关(连接不挂死);客户端断开经 `Request.Context().Done()` 退出
- **为何 SSE 而非轮询/WebSocket**:进度是服务端→客户端单向流——轮询有延迟且空转打 DB;WS 双向过重、鉴权易漏(反例:裸挂 WS 的接口后来专门收敛统一鉴权)、经 nginx 需额外 upgrade 配置;SSE 纯 HTTP,复用现有 auth 中间件与权限码,EventSource 自带重连
- 仪表盘常驻 SSE 数据面:fields 参数按需返回块(默认只推轻量统计);序列化对比(如 5s)变化才推;数据源用内存缓存(零 DB 查询);**同一数据 REST 与 SSE 必须同源**,否则前端两套口径

## 3. 前端 SSE 消费

- **两种清理语义不可混用**:常驻流 useEffect return es.close()+onerror 空处理(交给浏览器自动重连);一次性任务收到终态即 close+resolve、onerror close+reject——混用导致连接泄漏或 Promise 永不落定
- 用命名事件 addEventListener("job");JSON.parse 失败静默忽略,单条坏消息不能杀整条流
- dev 态 Next.js 代理会缓冲/截断长连接:SSE 用绝对后端地址直连绕过代理,直连不可得时靠 EventSource 自动重连兜底

## 4. 分片上传引擎(服务端)

- init 预分配目标文件(Truncate 稀疏占位)→PUT raw body 分片按偏移 **WriteAt 直写**——无"分片落盘+最终拼接"步骤,峰值磁盘占用 1×
- 乱序/断点:received_chunks 落库,Status 接口供客户端查缺口;幂等重传(重写覆盖+不重复标记)
- **分片边界精确**:LimitReader(expected+1) 读满期望后多余字节视为错误(防两个分片粘在一个请求里静默错位);末片按剩余字节算期望大小
- complete 校验大小/分片齐全/可选整文件 sha256(流式,内存恒定);pending(如 2h)/终态(如 24h)过期清理;消费完成自动清理暂存与会话
- purpose 场景化(ValidateInit/OnComplete 钩子)让多场景(瓦片/恢复包)复用同一引擎,替代各写一套上传;暂存目录独立于业务目录(业务目录清理白名单会误删暂存文件)
- 大包双通道消费:multipart 直传(小包)+upload_id 引用分片成品(大包,突破网关单请求体积限制;分片上限(如 32MiB)低于网关限制即无需改网关)

## 5. 分片上传客户端

- 并发窗口(如 3,夹紧 1-6)+单片重试(如 3 次):**重试前先查 status 对账**——响应丢失≠分片丢失,received_chunks 含该片则跳过;会话已定稿立即停止重试(幂等/进度/断点续传一石三鸟)
- 会话过期(404)自动重新 initiate,经 onUploadId 回调让调用方持久化新 id(刷新页面可续传)
- **大文件 sha256 流式计算**:Web Crypto digest 需完整 buffer,GB 级不可行——自实现流式 SHA-256(file.stream() 逐块喂入,内存恒定=单片)
- 进度以**服务端 returned received_chunks 口径**回调(断点/重试后仍准确);file.slice() 零拷贝切片,禁止一次性读入整个文件(曾致内存/CPU 陡增)
- 大小阈值分流:>50MB 自动走分片+sha256,小包保持 multipart 直传,双通道并存

## 6. 浏览器大文件通道决策

- **GB 级下载禁 fetch→response.blob()**(整个文件攒进 JS 内存,传输完成才失败):用 window.location.assign 同源直连,HttpOnly cookie 自动携带,交给浏览器下载管理器(边下边写盘、可恢复、零页面内存)
- 通道决策树:上传走分片引擎、下载走原生导航、任务进度走 SSE、双向实时走 WS+轮询降级(WS 断线时轮询兜底如 30s,恢复正常即停);推送去重集合封顶(如 200)防长会话内存增长
- BFF 代理坑(Next.js rewrites 为例):特定前缀规则必须排在通配规则之前(否则 404);dev 代理有 body 上限(大文件上传需调);CSP connect-src/media-src/worker-src 需 `blob:`;HTTPS 页面禁嵌 http iframe(混合内容拦截)

## 7. 基础设施陷阱(nginx/容器)

- **nginx proxy_buffering**:超内存缓冲(约 32KB)的响应写 proxy_temp 临时文件,磁盘紧张写不进则 nginx 停止回读上游——表现为**大文件下载 200 后 body 永久 pending、后端无错误日志**;对策 X-Accel-Buffering: no(下载与 SSE 两处都要)
- 磁盘空间:写前 statfs 预检(余量+目录实际大小,恢复类按包×3 预留),写入期 ENOSPC 统一转"可用 X MB,需 Y MB"可操作文案;statfs 不可用跳过不阻断由 ENOSPC 兜底;Windows 平台 API 差异拆平台文件
- **bind mount 挂载点 rename 返回 EBUSY**:绝不 rename 挂载点根目录,只操作其内部条目;staging 建在目标目录内(同文件系统,条目级 rename 零拷贝);旧内容移入隐藏暂存(. 开头)→新内容移入→成功删暂存、失败回滚;隐藏目录打包时跳过(防污染下一份备份)

## 8. 实测陷阱速查

| 现象 | 根因 |
|---|---|
| SSE 连上无事件 | nginx/代理缓冲,缺 X-Accel-Buffering: no;dev 代理截断长连接 |
| 大文件下载 200 后 body 永久 pending | proxy_buffering 写 temp 失败,nginx 停止回读上游 |
| 前端 Promise 永不落定/连接泄漏 | 常驻流与一次性任务的 SSE 清理语义混用 |
| 分片内容静默错位 | 两个分片粘在一个请求,缺 LimitReader(expected+1) 边界校验 |
| 上传重复入库 | 无幂等键;重试前未查 status 对账 received_chunks |
| 大文件上传前端内存爆 | Web Crypto digest 需完整 buffer/一次性读入文件,应流式 |
| 断点续传后进度条跳变 | 进度按本地已发数而非服务端 received_chunks 口径 |
| GB 级下载传输完才报错 | fetch→blob 把整个文件攒进 JS 内存,应原生导航下载 |
| 恢复替换目录报 EBUSY | rename 了 bind mount 挂载点,应只动内部条目+staging 回滚 |
| SSE 与 REST 数据口径不一致 | 同一数据两处实现未同源(合并进 SSE 后 REST 须同步补) |
| 任务 panic 整个服务崩 | job 执行未 recover 隔离 |
