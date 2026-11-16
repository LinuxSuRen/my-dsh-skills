---
name: branching-workflow
description: 每次提交/推送代码前，必须先询问 Rick 是新建分支还是用当前分支、走 PR 还是直推；不得擅自决定分支策略。AI agent 只负责 review PR 并加 comment，合并由人工执行或收到明确指令。UI 相关改动的 PR 必须附受影响页面的修改后截图。push 前必须检查远程默认分支是否前进，变了先合并进当前分支、解决冲突再推。
whenToUse: 任何即将执行 git commit / git push / 创建 PR / review PR / 合并 PR 的时机——提交前先问分支策略；合并前先确认授权；UI 改动的 PR 先备好页面截图；push 前查远程默认分支是否前进。
---

# Git 分支策略确认规范

所有代码改动在**提交之前**，必须先向 Rick（研发本人）确认分支策略，禁止擅自决定。

## 执行要求

1. **提交前必问**（二选一，让 Rick 明确选择）：
   - 新建分支（给出建议的分支名，如 `feat/xxx`、`fix/xxx`、`docs/xxx`）+ 是否创建 PR
   - 继续用当前分支（默认直推还是也要 PR，同样要确认）

   推荐话术：
   > 这个改动准备提交了：新建分支 `feat/xxx` 走 PR，还是直接提交到当前分支 `master`？

2. **禁止的擅自行为**：
   - 未经询问就 `git checkout -b` 新分支
   - 未经询问就直推 `master`/主干分支
   - 未经确认就 `gh pr create` 或合并 PR
   - 强推（`--force`）历史重写类操作——即使 Rick 要求过重写，执行前也要再次复述影响范围确认

2a. **push 前默认分支同步检查**：每次 `git push` 前先 `git fetch origin`，确认远程默认分支（一般是 `origin/master`；不确定时 `git symbolic-ref refs/remotes/origin/HEAD`）是否前进：
   - 前进了（`git rev-list --count HEAD..origin/master` > 0）：先 `git merge origin/master` 进当前分支，解决冲突并跑过项目验证命令（构建/测试）后再 push。冲突解决要理解两侧意图后取舍，禁止盲选一边。
   - 没前进：直接 push。
   - 例外——若前进的提交里含当前分支所叠加 PR 的 squash 版（master 出现单提交带 `(#N)` 后缀），直接 merge 会带入重复提交，按 3a 换基底重建后再 push。
   - 直推默认分支同样适用：push 前先 fetch，远程已前进就先合并（能快进则快进）再 push；禁止用 `--force` 盖掉他人提交。
   - PR 分支保持与默认分支同步，还能减少 review 冲突噪音与 GitHub mergeable 误报。

3. **PR 流程**（Rick 选择走 PR 时）：建分支 → 按规范提交 → push → `gh pr create`（标题/描述完整；UI 改动必须附页面截图，见下文专节）→ 等 Rick review 合并，**不要自行 merge**。

## PR Review 与合并权限

- **职责边界：agent 只 review + comment，不合并**。对仓库的 OPEN PR 主动给出审查意见——approve / request-changes（他人 PR）或 comment（本人账号 PR，GitHub 不允许自审生效），意见必须落成 PR 上的 comment，不只在会话里说。
- **合并权限属于研发**：默认人工合并；agent 仅在收到**包含"merge/合并"的明确指令**时才可执行合并。CI 变绿、review 通过、任务延续等都不构成合并授权。
- **review 检查项**：CI 状态、README 技能表、frontmatter 规范（name 与目录一致）、项目名泄漏、与既有技能重复度。embed 嵌入走 `//go:embed skills` 目录通配，新增技能目录无需登记清单。
- **结果呈现**：review 完成后，在会话中以**列表/表格**汇总每个 PR 的简要信息——编号、标题/分支、review 结论、遗留项、CI 状态——并为每个 PR 附上可点击的链接（`https://github.com/<org>/<repo>/pull/<N>`），方便研发直接打开处理。
- **定期/自动化 review**：会话开始时检查仓库有无待处理的 OPEN PR，有则主动发起/刷新 review；并向研发建议建立定期或自动化 review 机制（如定时唤起的 review 工作流），由研发决定是否启用。

3a. **叠加分支遇上 squash 合并**：上游 PR 若以 squash 方式合入（master 出现单个 `(\\#N)` 后缀提交、无原始 commit），叠加其上的分支必须重建——旧 commit SHA 全部不在 master 历史里，直接推送会显示 conflict/重复提交：
   - 换基底：`git rebase --onto master <旧基底> <分支>`（旧基底 = 原 PR 的最后一个原声 commit）
   - 顺带去重（同一修复已进上游）的 commit 用 cherry-pick 挑拣重建，丢弃已合并项
   - 重建后 `push -f` 更新 PR，mergeable 复核

4. **例外**：连续多个小改动属于同一任务且 Rick 已明确过策略（例如"接下来的改动都走 PR"），沿用已确认的策略即可，无需每次重复询问；但任务切换或涉及主干直推时必须重新确认。

5. 与 `git-identity` 技能配合：提交身份永远是 Rick 本人（`Rick <linuxsuren@users.noreply.github.com>`）。

## UI 改动的 PR 必须带页面截图

- **范围**：提交内容涉及 UI 部分——页面、组件、样式、布局、交互与用户可见文案——走 PR 时，PR 描述里必须放**受影响页面的修改后截图**；纯后端/构建/文档改动不适用。
- **覆盖度**：PR 动了几个页面/视图就给几张；响应式或按视口分支的改动（对齐 `ui-responsive-list`）补齐各断点视口；弹窗、错误态等关键交互状态能用截图表达的也带上。
- **怎么截**：本地把应用跑起来截图（Playwright 临时脚本即可，可复用或沉淀到项目 e2e 设施）；图片提交进 PR 分支并在描述中以相对路径引用，或直接粘贴进 PR 描述。
- **真实性**：截图必须来自本分支构建的真实渲染，禁止拿旧图、设计稿或其他分支产物充数；每张注明对应页面与视口。
- **评审对齐**：UI 改动 PR 缺截图视为描述不完整——review 时 request-changes（本人账号 PR 则 comment 指出），补齐后再过。
