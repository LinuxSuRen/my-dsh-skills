---
name: askills-sync
description: 仓库根存在 askills 二进制时，在技能变更落库后（提交/合并/pull）或安装落后时，重建二进制并调用 askills install 把最新技能同步到各 AI 编码工具目录；覆盖式幂等，重跑即升级。
whenToUse: 在 askills 源仓库（my-dsh-skills）完成技能新增/修改并落库之后，或发现某工具加载的技能落后于仓库内容时——先确认仓库根存在可执行的 ./askills（.gitignore 忽略的构建产物），存在才走本技能。
---

# askills 技能同步规范

仓库根存在 `./askills` 二进制时，在合适的时机调用它把最新技能同步安装到各 AI 编码工具（dsh、codex、opencode 等 agents 兼容工具）。

## 前提判定

- 检查仓库根 `./askills` 存在且可执行：`test -x ./askills`。它已被 `.gitignore` 忽略，是本地构建产物，不是提交内容。
- 不存在时本技能不适用：向研发提示可 `go build -o askills ./cmd/askills` 构建或 `go run ./cmd/askills ...` 单次运行，由研发决定是否构建与安装，不擅自写入工具目录。

## 关键语义（调用前必须知道）

1. **技能嵌在二进制里**：`//go:embed skills` 目录通配，安装内容 = 构建那一刻的快照。技能变更后**必须先重建二进制**再 install，否则装出的是旧技能。
2. **覆盖式幂等**：每个技能目录先删后拷，重跑即升级。默认目标 `~/.agents/skills`（[Agent Skills](https://agentskills.io/) 共享标准目录，dsh/codex/opencode 都会发现）；`--tool dsh|opencode|codex|all` 换目标，`--dir` 指定任意目录，`--skill a,b` 只装子集。
3. **软链会被替换为拷贝**：安装目录里若是指向仓库的软链，install 会删除软链换成拷贝副本——此后以安装副本为准，仓库再改动必须重新 install 才生效；软链的「改仓库即生效」特性随之消失。

## 合适的时机

| 时机 | 动作 |
|---|---|
| 技能变更落库后：本仓库 commit / PR 合并 / pull 更新了 `skills/` | 重建 `./askills` → `./askills install`（按需 `--tool`） |
| 会话中发现工具加载的技能与仓库不一致（安装落后） | 同上——幂等操作，无需犹豫 |
| 没有技能变更 | 不重装；嵌内内容不会自己变 |

## 标准操作

```bash
# 在仓库根执行
go build -o askills ./cmd/askills   # 重建：嵌入最新 skills/
./askills list                      # 核对二进制内技能清单
./askills install                   # 全部 → ~/.agents/skills（默认，多工具共享）
./askills install --tool all        # agents/dsh/opencode/codex 四目标全装
./askills install --skill git-identity,branching-workflow  # 只装子集
```

安装后验证：目标目录出现对应技能目录与 `SKILL.md`；新开会话（或技能目录刷新）即加载新内容。

## 边界与衔接

- 本操作只写用户级工具目录，不改仓库内容、不构成提交；`./askills` 始终保持未跟踪。
- 与 `skill-refresh`、`skill-distillation` 衔接：两者产出技能变更并提交后，收尾一步即走本技能同步安装，形成「改仓库 → 同步工具」闭环。
- 与 `branching-workflow` 衔接：重建的二进制若要提交入库（正常不提交），先问分支策略。
