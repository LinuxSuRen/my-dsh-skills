---
name: open-source-contribution
description: 开源贡献规范：Issue 与 Pull Request 的提交要求、Review 礼仪与流程、good-first-issue 创建模板、开源文档写作规范、仓库门面（About 与 README）完善、CLI 工具收录（hd-home）、开源仓库自动化（GitHub Actions）。
whenToUse: 代表用户参与开源项目时使用——提交 issue、创建/更新 Pull Request、执行或响应 code review、创建 good-first-issue、编写开源相关文档、创建或完善开源仓库门面、判断开源 CLI 工具是否收录进 hd-home、为开源仓库配置自动化；不限于特定项目。
---

# 开源贡献规范

提炼自 [LinuxSuRen/open-source-best-practice](https://github.com/LinuxSuRen/open-source-best-practice)（OSBP），并与其保持同步演进。总原则：**公开、透明**——公开的不仅是结果，更是过程；review 不是审核，任何人都可以是 reviewer。

## 1. Issue

提交前**先搜索现有 issue 列表**，确认没有重复后再新建。常见误区：只有标题没有内容、只给现象不给上下文、只给截图不给文字（不利于检索）。

写作要求：

- 标题简洁、规范，用标签或标题前缀分类，例如 `Bug: xxx`、`Proposal: xxx`、`Question: xxx`；
- 提供完整上下文：环境信息、版本、操作步骤、预期与实际结果、错误/异常的**关键文字**（截图不能替代文字）；
- 与 UI 相关的问题附截图；
- 语言遵循对应社区期望的规定（看已有 issue 的惯例）。

## 2. Pull Request

核心纪律：**一个 PR 只包含一类修改**；单次改动再多，也要新开一个分支（禁止直接在 master/main 上提交变更）。

创建前：

- 修复的问题已有对应 issue 时，先确认没有人提交过 PR，并在 issue 下留言说明修复计划；
- 预计变更较大时，先创建 issue 描述想法，视难易程度、争议性等征得反馈后再动手；
- 首次向某项目提交 PR 前，浏览已合并 PR 的评论列表与格式，向社区惯例看齐。

提交时：

- 提交前充分自测；未准备好接受 review 时，标题加 `WIP: ` 前缀；
- 关注 commit 记录（"卷面分"），只保留希望合并进上游的记录，其余提前 squash；思考能否把"大型" PR 拆成多个以减轻 reviewer 压力；
- PR 描述给出尽可能多的详情：相关 issue、解决的问题、方便 review 的上下文、自测过程；可能引起争议的部分给出解释；
- 认为合并后可修复对应 issue 时，使用 `close #123` 等[关联表达式](https://docs.github.com/en/issues/tracking-your-work-with-issues/using-issues/linking-a-pull-request-to-an-issue)；
- 涉及 UI 的改动，给出修改前后的效果截图。

更新时：

- 避免 `--force` 强制推送——reviewers 将无法轻松看到最新修改的部分；
- 避免同一主题的 PR 反复关闭、新建，避免同一 PR 中频繁提交；
- PR 超过一周没有得到 review，可 cc 相关 team；没有 team 可 cc 时，找最近合并过类似 PR 的人，说明缘由并对打扰表示歉意；
- 禁止通过即时聊天工具催促特定的人 review。

## 3. Review

生命周期：PR 不适用于有紧急合并需求的场景；预期在 2～7 天内完成 review 并合并。

作为作者：

- 提交后先自行检查一遍，发现问题就标记"进行中"（`WIP: `）；
- 对不确定的部分，主动以评论形式写出自己的观点；
- 没有人有义务 review 你的 PR（包括维护者）；对每位 reviewer 表示感谢；最佳 review 周期内不要催促；
- 确实需要请求帮助时给出说明，优先 @ team，其次才 @ 个人。

作为 reviewer：

- 观点明确，避免模棱两可的评论——作者要据此决定是否修改；
- 不确定的问题这样表述："我感觉这里可能有问题"，给出建议做法与理由，并**注明该评论不阻碍 PR 合并**；
- 确定有问题时，给出能证明观点的信息或数据，有权威资料的一并给链接（例如官方社区文档）；
- 优先请求一个 team review；也可以请最后修改该文件的贡献者、或对应 issue 的创建者帮忙；
- review 完成但还需人工验证时，用评论命令（如 `/hold`）阻止过早自动合并，避免人为干预自动化过程；合并前尽可能多地运行自动化测试。

## 4. 创建 good-first-issue

判断标准：对"初次接触该项目的人"友好。具体要求：

- 有清晰的技术栈要求——"新人"只表明初次接触，与技术水平无关，应列出完成该 issue 所需技能；
- 有清晰的上下文描述，不含糊背景；
- 没有明显或潜在的时间约束，不与 milestone/版本计划挂钩；
- （可选）有助于贡献者了解项目结构与贡献流程。

模板：

```markdown
## Background

## Technical requirement

## Expect

## Potential TODO list
```

## 5. 开源文档写作

面向开源实践类文档的写作规范（写作风格可参考 [document-style-guide](https://github.com/ruanyf/document-style-guide)）：

- 正确使用中英文标点（中文语境用中文标点：`。` 而非 `.`）；
- 专有名词书写规范（英文注意词性、大小写）：`Spring Cloud OpenFeign` 而非 `spring cloud openfeign`；
- 避免缩写，首次出现时说明：`Spring Cloud Alibaba（后文简称 SCA）`；
- Markdown 代码块标注对应语言；
- 专有名词首次出现时给出链接；
- 重点内容加粗或着重强调；
- 明确受众群体，站在受众角度确保参考价值；使用合适的人称代词；
- 完稿后自己通读一遍再发布。

## 6. 仓库门面（About 与 README）

开源仓库应在第一时间让陌生开发者看懂项目是做什么的：

- **GitHub About**：一句话准确描述项目是什么、解决什么问题，不堆砌关键词；用小写 kebab-case 的 topics 增加可发现性；有官网或文档时填写 homepage。
- **README**：面向「第一次听说本项目的人」写作，至少覆盖：项目是什么、解决什么问题；快速开始（安装 + 最小可用示例）；使用与配置方式；如何贡献（链接贡献指南）；License。badge 等装饰不替代上述内容。
- 新建开源仓库或完善仓库门面时，逐项检查以上内容是否**简洁、准确、与项目现状一致**。

## 7. CLI 工具收录（hd-home）

- **触发条件**（需同时满足）：当前项目是开源项目；属于命令行工具（CLI）；GitHub Releases 上有对应的二进制资产。
- 满足时**提示研发**是否把项目收录到 [linuxsuren/hd-home](https://github.com/linuxsuren/hd-home)（`hd` 下载器的工具注册表，收录后可通过 `hd get <name>` 安装）；**研发同意后才行动**，不得擅自提交。
- 收录方式：向 hd-home 提 PR 新增 `config/<org>/<repo>.yml`，字段参照其 [CONTRIBUTION.md](https://github.com/linuxsuren/hd-home/blob/master/CONTRIBUTION.md)：
  - `filename` 用 `{{.Name}}`、`{{.OS}}`、`{{.Arch}}`、`{{.Version}}` 模板匹配 release 资产名；
  - 按需提供 `binary`（压缩包内二进制名）、`targetBinary`、`tar`、`formatOverrides`、`replacements`、`categories`、`versionCmd` 等；
  - **禁止直接编辑 hd-home 的 README.md**——它由 `README.tpl` 经 yaml-readme 自动生成。
- 提交前先确认研发的分支策略；PR 描述给出工具名、仓库链接与安装验证方式（如 `hd get <name>`）。

## 8. 开源仓库自动化（GitHub Actions）

GitHub 开源项目**推荐全部配备**以下自动化；agent 检查开源仓库时发现缺失，应提示研发补充：

- **CI**：push/PR 触发构建 + 测试（按技术栈选 `go test`、`npm test` 等）；
- **Lint**：静态检查（golangci-lint、eslint 等），可并入 CI job；
- **Release**：打 tag 触发多平台构建并发布二进制或制品；
- **CodeQL**：安全漏洞扫描（`.github/workflows/codeql.yml`）；
- **Dependabot**：依赖自动更新（`.github/dependabot.yml`，覆盖包管理与 github-actions 生态）。

新增或修改 workflow 走 PR；遵循最小权限原则（默认 `permissions: contents: read`，发布 job 才给 `contents: write`）。
