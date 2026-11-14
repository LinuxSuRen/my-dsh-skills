---
name: github-ci-release
description: GitHub 项目自动化门面：CI 与 Release 必须有对应 GitHub Actions（构建/测试/lint + tag 触发发布）；Go 项目提供二进制与容器双轨分发（goreleaser 多平台 + 多架构镜像推 ghcr），同版本同源。
whenToUse: 新建 GitHub 仓库、补齐 CI/release workflow、配置分支保护的合并门槛、为 Go 项目做发布（goreleaser、Docker 镜像、ghcr）、或检查项目自动化完整性时。
---

# GitHub 项目 CI 与 Release 规范

核心原则：**CI 与 Release 是仓库门面的一部分，不是可选项**。访客判断项目可信度的第一眼就是：有没有绿的 CI 徽章、有没有像样的 Release。新仓库开工第一件事就是放 CI，首个版本发布前 release workflow 必须就位——不要打了 tag 才补。

## 1. CI（每个 GitHub 项目必备）

- `.github/workflows/ci.yml`，最小覆盖：**构建 + 测试 + lint**，按技术栈选择（`go build ./... && go test ./...` + golangci-lint；`tsc --noEmit` + vitest；mvn/gradle 等）。
- 触发：`push`（master/main）+ `pull_request`；用 `concurrency` 取消同分支排队，避免任务堆积。
- 依赖缓存加速：`actions/setup-go` 自带模块缓存、`setup-node` 的 `cache`；无缓存的长任务再考虑手动缓存。
- **CI 与分支保护联动**：master 保护规则要求 CI 通过才可合并（arch-dev-workflow 的分支保护条目），CI 不作为门槛等于没有 CI。
- 用 dependabot 保持 actions 版本更新（本仓库实践：checkout/setup-go/golangci-lint-action/gh-release 定期升版）。
- 不在本机跑的重型验证（docker build、E2E）尽量下沉 CI（arch-dev-workflow 验证策略）。

## 2. Release（tag 触发自动化）

- `.github/workflows/release.yml`：打 `v*` tag 触发；构建产物 + 自动 release notes + 附加 checksums。
- 版本号 `vX.Y.Z` 语义化；Release 发布后不可变（arch-dev-workflow：防篡改已发布产物，修问题发新版本而非覆盖）。
- release notes 用 GitHub 自动生成（`generate_release_notes: true`）起步，重要版本手写要点。

## 3. Go 项目分发双轨（二进制 + 容器缺一不可）

### 二进制
- **goreleaser** 一站式：多平台交叉编译 + 压缩包 + checksums + 自动上传 GitHub Release。
- 平台矩阵：linux/darwin/windows × amd64/arm64（arm64 时代默认双架构，按用户群裁剪）。
- 轻量补充：文档给出 `go install github.com/<owner>/<repo>@latest`；CLI 工具满足收录条件时进 hd-home（见 open-source-contribution）。

### 容器
- **多阶段 Dockerfile**：builder（golang）→ 运行时层（distroless/alpine/scratch），最终镜像只留二进制与必需文件。
- **多架构**：docker buildx，linux/amd64 + arm64；推 **ghcr.io**——GitHub 原生 registry，`GITHUB_TOKEN` + `permissions: packages: write` 即可，无需额外凭证。
- tag 三层：`latest`、`vX.Y.Z`、`vX.Y`；镜像名 = 仓库名、命名空间 = 原始 owner（与 project-branding 一致：分发物在原始项目名下，不在 fork 名下）。
- 静态编译（CGO_ENABLED=0）可让 scratch 镜像与交叉编译同时简化。

### 双轨同源
- 二进制与镜像**同一 tag、同一 commit** 构建触发（版本单源，见 installer-bundle-design）；禁止 tag 与镜像各自漂移。
- README 安装章节写全三条路径：`go install` / Release 二进制 / 容器镜像，各附一行命令。

## 4. 检查清单（新仓库/盘点旧仓库）

1. CI workflow 存在且 PR 上真实跑（绿徽章）？
2. 分支保护把 CI 设为合并门槛？
3. release workflow 在首个 tag 前就位？
4. Go 项目：goreleaser 配置 + 多阶段多架构镜像推 ghcr？
5. 二进制与镜像同版本同源、checksums 齐全？
6. README 安装章节三路齐全（go install/二进制/容器）？
7. dependabat 覆盖 actions 版本？
