---
name: frontend-dev-workflow
description: Rick 的前端研发品味：Vue 3 + Vite 默认、按需上 TS/React、Element Plus 与手写设计令牌双轨、极简依赖、fetch 封装 API 层、自清理组件模式、中文文案规范。提炼自 linuxsuren 开源项目。
whenToUse: 任何前端开发任务——新建 web 子项目、写页面/组件、封装 API、选型 UI 方案、写样式或用户文案时，先按本技能对齐技术选型与代码风格，再动手。
---

# 前端研发工作流（Rick 的前端品味）

提炼自 [open-pdf](https://github.com/linuxsuren/open-pdf)、[open-cloud-web](https://github.com/linuxsuren/open-cloud-web)、[onvif-ai](https://github.com/linuxsuren/onvif-ai)、[onvif-local](https://github.com/linuxsuren/onvif-local)、[local-gb28181-adapter](https://github.com/linuxsuren/local-gb28181-adapter)、[smart-chat](https://github.com/linuxsuren/smart-chat) 等开源项目。核心态度：**克制、按需、最小依赖，把复杂度留给真正需要的地方。**

## 1. 技术选型（默认值，偏离需先确认）

| 维度 | 默认选择 | 例外条件 |
|---|---|---|
| 框架 | Vue 3 | 宿主/嵌入环境是 React（如 DSH 插件）才用 React 18 |
| 构建 | Vite，`package.json` 里只有 dev/build/preview | 不引脚手架全家桶 |
| 语言 | 简单工具/CRUD 管理台用 JS | 实时媒体、复杂状态机、长生命周期项目上 TS（`vue-tsc -b`） |
| UI 组件 | CRUD 管理台用 Element Plus（+ @element-plus/icons-vue） | 产品化/品牌化 UI 手写 CSS + 设计令牌，不用 Tailwind |
| 状态管理 | 无库：Vue 模块级 `reactive` store 或组件内 `ref`；React 手写 store + `useSyncExternalStore` | — |
| 路由 | 无 vue-router：`view = ref('xxx')` + `v-show` 标签页切换，视图拆到 `views/*.vue` | 真有多路由需求再议 |
| Lint/格式化 | 不配 ESLint/Prettier，风格靠本规范约束 | — |

依赖按需引入，禁止"以防万一"加依赖。典型白名单：`clsx`、`marked`、`dompurify`、`jmuxer` 这类单一职责小库。

## 2. 项目结构

Vue 子项目（`web/`）：

```
web/src/
  main.js          # createApp 引导：OAuth 回调处理 → boot 拉取 me → mount
  App.vue          # 页面编排：标签页导航 + 视图切换 + 对话框调度
  api.js           # 唯一的 HTTP 封装（见 §3）
  components/      # 可复用组件（CameraTable、SystemCard…）
  views/           # 页面级视图（多标签页项目）
  composables/     # useXxx() 有状态逻辑（TS 项目）
  style.css        # 全局设计令牌 + 基础样式
```

React 子项目：`src/{main.tsx, App.tsx, components/, hooks/, lib/{api.ts, store.ts}, styles/{tokens.css, *.module.css}}`。

## 3. API 层：单文件 fetch 封装

所有 HTTP 走一个 `api.js`/`api.ts` 模块，用 `fetch` 不用 axios：

1. **超时**：`AbortController` + `setTimeout`（15s），超时抛 `ApiError('请求超时，请检查网络后重试', 'timeout')`。
2. **网络错误**：包装为 `ApiError('无法连接后端服务，请确认服务已启动', 'network')`，禁止把原始异常抛给 UI。
3. **错误 envelope**：非 2xx 时解析响应体 `{error: {message, code}}`，兜底文案 `请求失败（HTTP xxx）`；`ApiError` 携带 `code` 供 UI 分支。
4. **响应 envelope**：`{data}` 自动解包返回；空响应返回 null；非 JSON 容错。
5. **按资源导出具名函数**：`getCameras()`、`createCamera(camera)`、`updateCamera(id, camera)`、`deleteCamera(id)`；路径参数一律 `encodeURIComponent`。
6. **鉴权**：token 存 localStorage；OAuth 回调 `/#token=` 在 main 引导期截获、写入并 `history.replaceState` 清理 URL。
7. **SSE/WebSocket**：同样收敛到 api/composables 模块，自动重连（约 2s），组件卸载自动断开。

## 4. 组件模式

Vue（一律 `<script setup>`，TS 时加 `lang="ts"`）：

- `ref` 优先于 `reactive`；派生值用 `computed`（如权限过滤标签页）。
- 数据流 props down / emit up；对话框/抽屉用 `v-model:visible` + 当前编辑对象 prop；新建与编辑共用同一对话框（`editingXxx` 为 null 表示新建）。
- 加载：`onMounted` 里 `Promise.all` 并行拉取，`loading` ref 驱动表格状态；列表兜底 `Array.isArray(x) ? x : []`。
- 轮询：`setInterval` 存模块级变量，`onUnmounted` 清理；**所有定时器、WebSocket、事件监听必须清理**。
- CRUD 编排集中在 App.vue：`openCreate`/`openEdit`/`handleRemove`，成功后 `ElMessage.success` + `loadAll()` 刷新。
- composable 内部自注册清理（`useWebSocket` 里 `onUnmounted(disconnect)`），调用方不必关心释放。

React：函数组件 + hooks；稳定引用用 `useCallback`/`useRef`；store 订阅用 `useSyncExternalStore(subscribe, getState)`，样式用 CSS Modules。

## 5. 样式与设计令牌

全局 `style.css` 在 `:root` 定义令牌，组件只消费变量不硬编码：

- 颜色：`--color-bg-base/panel/elevated`、`--color-border-subtle`、`--color-text-{dim,secondary,primary,bright}`、`--color-accent-*`，语义别名 `--color-success/warning/danger`。
- 间距：`--space-1..16`（4px 基准）；字体：`--font-display`/`--font-mono`。
- 中文优先字体栈：`'PingFang SC', 'Microsoft YaHei', 'Segoe UI', system-ui, sans-serif`；数据密集界面用 mono 强调数值。

命名 BEM（`app-header__brand`）；Vue 用 scoped style，React 用 `*.module.css`。

两种页面骨架按场景选：

- **管理台**：`.page` max-width 1200px 居中、24px/16px padding、浅灰背景（#f5f7fa）、卡片分块（`.block` 上间距 16px）、正文 14px。
- **控制台/产品**：100vh `app-shell`（header + grid 主区 + 固定宽侧栏），深色数据密集风格，令牌驱动。

视觉态度（对齐 [Anthropic frontend-design](https://github.com/anthropics/skills/tree/main/skills/frontend-design)）：克制、有意图；拒绝模板化默认样式——不用全大写标签、单词强调、渐变装饰、千篇一律卡片阵列；结构即信息；动效只用于回应操作（打开、展开、确认），不做入场动画堆砌。

## 6. 文案规范（与 arch-dev-workflow「展示与数据分离」一致）

- 用户可见文案由前端定义，后端只返回稳定 code/枚举；前端维护映射表，未知 code 回退兜底文案。
- 错误提示友好包装：说清发生了什么 + 怎么办（"请求超时，请检查网络后重试"）；禁止透传 `rpc error`、`connection refused` 等底层技术错误。
- 中文文案为主，专有名词（服务名、命令、协议名）保留英文；同一概念全文同一称呼；禁止后置括号补充解释。
- 操作动词全流程一致：保存就是"保存"，不要一处"提交"一处"发布"。

## 7. 验证

- TS 项目：`vue-tsc -b && vite build`（Vue）/ `tsc -b && vite build`（React）；JS 项目：`vite build`。
- 关键路径端到端用 Playwright（`e2e/*.spec.ts`）。
- 遵循 arch-dev-workflow：不在本机跑重型构建/部署验证，交给 CI。

## 8. 决策速查

| 场景 | 选择 |
|---|---|
| 新工具的配套管理界面 | Vue 3 + JS + Vite + Element Plus，管理台骨架 |
| 品牌化产品 UI / 演示大屏 | Vue 3 + TS + 手写令牌 CSS，app-shell 骨架 |
| 嵌入 React 宿主的页面 | React 18 + TS + CSS Modules + tokens.css |
| 需要全局状态 | 模块级 reactive store / useSyncExternalStore，不引库 |
| 需要多个页面 | ref 标签页 + views/，不引 vue-router |
