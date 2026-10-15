---
name: installer-target-execution
description: 安装程序在目标机上的执行规范：执行通道抽象、提权矩阵、凭据安全、跨平台适配、传输与幂等。
whenToUse: 实现安装器/部署器的目标机执行层（本机或 SSH 远程执行命令、传文件、提权、处理凭据）时使用；包体设计见 installer-bundle-design，装前校验见 installer-preflight。
---

# 目标机执行规范

> 沉淀自 deploy-tool 的 Runner/SSH/sudo/跨平台实现与踩坑（darwin/Windows 本机安装、sudo -S 与 stdin 冲突、docker 组授权等）。

## 执行通道抽象

- 定义统一的执行接口（Run/MkdirAll/UploadFile/WriteFile/Exists/FileMD5/Platform），本机与 SSH 实现同构，**安装流程不感知底层**；
- 接口带 `Platform()`（目标机 OS）：df 语法、时间校准、docker 安装方式、默认目录等一切平台差异经它分支，不要散落 if GOOS；
- 路径统一 `/` 分隔（工具端可能在 Windows，`\` 会让远端 bash/tar 炸）。

## 提权矩阵（按序探测，一次确定全局前缀）

1. root（`id -u` == 0）→ 无需前缀；
2. `sudo -n` 免密；
3. `sudo -S` 密码（stdin 传递，**绝不落命令行参数**）；
4. 无 sudo 但 docker 可直连（docker 组用户）→ 免提权 + 提示用家目录。

- 提权前缀是**管道前缀**（`printf 'pw\n' | sudo -S -p ''`）：与"也要读 stdin 的命令"（如 `docker login --password-stdin`）组合时会互踩——按提权形态分别构造（sudo -S 形态喂两行：sudo 密码 + 业务密码）；
- macOS/Windows 无 sudo 概念：docker 由 Docker Desktop 提供、当前用户直连；不可用时**不自动安装**，报错并给官方下载链接；
- Linux 装完 docker 后顺手**授予登录用户 docker 组**（`usermod -aG`，组不存在先建；`-aG` 只追加）：组成员按 /etc/group 数据库判断（与会话有效组无关），提示重新登录生效；失败仅告警。

## 凭据安全

- 所有口令（仓库/SSH/sudo/证书口令）经 stdin 或 SSH 加密通道传递，不进命令行（ps 可见）、不进浏览器持久存储；
- 拼接 shell 前对嵌入值做白名单校验（如用户名 `^[a-z_][a-z0-9._-]{0,31}$`）+ 单引号转义（shQuote）双保险；
- 避免嵌套引号：需要组合命令时宁可分步执行，不用 `sh -c '...$var...'` 包 shQuote 的值。

## 文件传输与幂等

- 上传优先 scp（BatchMode 防密码场景挂起），失败切内置 SFTP；目标机未启用 sftp 子系统时可经提权自动启用（改 sshd_config + reload，不断连）后重建连接；
- **MD5 增量上传**：传前比对远端同名文件，一致跳过；传后再校验，不一致报错——大文件断点免重传 + 防静默损坏；
- **重装幂等**：解压前 `rm -rf` 旧目录（root 属主内容必须经 sudo 清）；失败时保留已上传包，重跑自动复用。

## 跨平台适配清单（Linux/macOS/Windows 本机安装）

| 差异点 | 处理 |
|---|---|
| 默认安装目录 | Linux `/opt/...`；macOS/Windows 家目录（/opt 普通用户不可写且无需提权） |
| 磁盘检测 | GNU `df -B1 --output=avail` / BSD `df -k` 换算；目录不存在时**向上取最近存在的祖先** |
| 时间校准 | `date -s` 是 GNU 语法，仅 Linux；其他平台跳过并说明 |
| bash 依赖 | Windows 需 Git for Windows；启动时 `exec.LookPath("bash")` 检查，缺失给指引 |
| docker 引擎安装 | 仅 Linux（离线物料/get.docker.com）；macOS/Windows 一律手动 Docker Desktop + 官方链接 |
