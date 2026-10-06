# my-dsh-skills

[DeepSeek Harness (DSH)](https://www.npmjs.com/package/@deepseek-ai/dsh) 个人技能库——AI 编码助手按需加载的可复用任务指令集。

## 技能列表

| 技能 | 说明 |
|---|---|
| [`arch-dev-workflow`](arch-dev-workflow/SKILL.md) | 软硬件架构师通用工程方法论：方案先行与最小实现、Git 纪律、验证策略、接口与标识符设计原则、跨端契约与展示分离。适用于任何技术栈的研发任务。 |
| [`git-identity`](git-identity/SKILL.md) | 提交代码必须使用 Rick 本人的 Git 身份（user.name=Rick, email=linuxsuren@users.noreply.github.com），禁止使用 agent/占位身份提交。 |
| [`branching-workflow`](branching-workflow/SKILL.md) | 每次提交/推送前必须先询问研发：新建分支还是当前分支、走 PR 还是直推；不得擅自决定分支策略。 |
| [`open-source-contribution`](open-source-contribution/SKILL.md) | 开源贡献规范：Issue/PR 提交要求、Review 礼仪、good-first-issue 模板、开源文档写作、仓库门面（About/README）完善、CLI 工具收录（hd-home）。提炼自 [open-source-best-practice](https://github.com/LinuxSuRen/open-source-best-practice)。 |
| [`skill-refresh`](skill-refresh/SKILL.md) | 技能库保鲜：按周期（state.md 配置）调研工具/语言最佳实践，生成带来源的更新提案，研发确认后才修改与提交。 |

## 如何使用

每个技能是一个目录，内含带 YAML frontmatter 的 `SKILL.md`（`name` + `description` 为必填，`whenToUse` 可选）。DSH 的 skill 文件系统提供商会扫描各根目录并监听变化，添加或修改后无需重启。

### 方式一：customSkillDirs（推荐）

clone 本仓库后，在 DSH profile 的 `~/.dsh/profiles/web/cordis.patch.yml` 追加：

```yaml
- id: skill-filesystem
  name: '@deepseek-ai/dsh-skill-filesystem'
  config:
    customSkillDirs:
      - /path/to/my-dsh-skills
```

之后在任意会话中说“按 arch-dev-workflow 的规范做”，或让 AI 在匹配任务时自动调用。仓库更新后 `git pull` 即可同步最新技能。

### 方式二：软链到用户技能目录

不想改配置时，直接把技能目录软链进 DSH 用户级技能根：

```bash
git clone https://github.com/linuxsuren/my-dsh-skills ~/Workspace/github/linuxsuren/my-dsh-skills
ln -s ~/Workspace/github/linuxsuren/my-dsh-skills/arch-dev-workflow ~/.dsh/skills/arch-dev-workflow
```

### 方式三：复制进项目

把某个技能目录复制到项目的 `.dsh/skills/` 或 `.agents/skills/` 下，随仓库提交，团队成员打开会话即自动获得（`~/.agents/skills` 亦可被其他兼容 agents 约定的 AI 工具共享）。

## 如何贡献

1. 新建 `<skill-name>/SKILL.md`（kebab-case 命名），写入 frontmatter 与指令正文：
   ```yaml
   ---
   name: my-skill
   description: 一句话说明技能做什么。
   whenToUse: 什么场景下应使用本技能。
   ---
   ```
2. 正文写清楚可执行的操作规则，避免与具体仓库耦合——项目专属规范应放进各项目自己的 `.dsh/skills/`。
3. 提 PR，更新本 README 的技能列表。

## License

[MIT](LICENSE)
