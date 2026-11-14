---
name: project-branding
description: 项目门面规范：所有项目必须有 logo（单源文件、README 与 Web 端 favicon/导航统一引用）；开源项目的链接一律指向原始上游仓库而非 fork 地址（remote 是 fork 时取 parent）。
whenToUse: 新建项目、补齐仓库门面、做 Web 端 header/favicon/登录页、写 README 头部与项目地址、或发现项目缺 logo、链接指到 fork 地址时。
---

# 项目门面：logo 与原始项目链接

## 1. 所有项目必须有 logo

- 无例外：CLI 工具、Web 服务、库都要有 logo；"小工具不用 logo"不成立。
- **单源文件**：仓库只放一份 logo 源文件（推荐 `docs/logo.svg`；必须位图则 `docs/logo.png` ≥512px），其余场合一律引用不复制。
- 引用位置：
  - README 顶部居中，置于标题上方：`<p align="center"><img src="docs/logo.svg" width="120"></p>`
  - Web 项目：favicon + 登录页/导航栏 logo 构建期引用同一源文件（复制进静态目录是构建产物，不算第二份源）
  - CLI：README 必引；`--version` 的 ASCII 标识可选
- SVG 优先（无损缩放）；位图从 SVG 导出，不在位图上二次编辑。

## 2. 开源项目链接原始仓库（不是 fork 地址）

- README、Web 页面 footer、仓库 About 侧栏、安装命令（`go install upstream/...@vX`、`docker pull upstream/...`）中的仓库链接**一律指向原始上游仓库**。
- fork 只是同步与贡献通道，不是门面：用户应被引导到上游看文档、提 issue/PR。
- **判定方法**——本仓库 remote origin 若是 fork，取 parent：

  ```bash
  gh repo view --json parent --jq '.parent.owner.login + "/" + .parent.name'
  # 输出 null 则本仓库就是原始仓库，直接用 origin
  ```

- 为他人项目贡献时的 fork 开发仓库：门面保持上游原样不改（贡献礼仪见 open-source-contribution）；只有自己名下的原始仓库才做门面定制。
- 上游已归档/失联的 fork 承继场景：显式注明来源（"fork of \<upstream\>，上游已停止维护"），不冒充原始项目。

## 3. 检查清单

1. 仓库有唯一 logo 源文件，README 顶部已引用？
2. Web 端 favicon/导航/登录页与 README 用的是同一份源？
3. `gh repo view --json parent` 确认链接指向的是原始仓库？
4. 安装命令里的模块路径/镜像名是上游的？
5. 页面 footer 的项目地址可直接点开且落在原始仓库？
