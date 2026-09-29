---
name: git-identity
description: 提交代码必须使用 Rick 本人的 Git 身份（user.name=Rick, email=linuxsuren@users.noreply.github.com），禁止使用 agent/占位身份提交。
whenToUse: 任何需要 git commit 的场合——写代码、修 bug、生成文档后提交，一律先确认提交身份。
---

# Git 提交身份规范

所有 `git commit` 必须使用我本人的身份，不得使用任何 agent/占位身份（如 `agent@dsh`、`integrator@dsh`、`Copilot` 等）：

- `user.name` = `Rick`
- `user.email` = `linuxsuren@users.noreply.github.com`

## 执行要求

1. **提交前检查**：`git config user.name && git config user.email`，若不是上述身份则先修正（优先改仓库局部配置，不污染全局）：
   ```bash
   git config user.name "Rick"
   git config user.email "linuxsuren@users.noreply.github.com"
   ```
2. **内联身份提交**（双保险，即使局部配置被改掉也生效）：
   ```bash
   git -c user.name="Rick" -c user.email="linuxsuren@users.noreply.github.com" commit -m "..."
   ```
3. **禁止**在 commit 命令里出现任何 `user.email=agent/...`、`user.name=integrator` 之类的占位身份。
4. 提交信息语言与格式遵循项目惯例（本仓库：中文 + Conventional Commits）。
5. 若发现历史提交混入了占位身份，提醒我是否需要 `git rebase` 重写作者并强推（此操作有风险，必须先征得我同意）。
