# USC (Universal Skill Compiler)

> **面向 AI Agent Skill 的零信任洁净室编译器**  
> 彻底解决 Agent 插件来源不可信、窃取 Token、偷读本地文件等安全风险。将任何不可信 Skill 提炼为最小权限契约，在物理洁净室中独立重建，并一键载入你的主流 Agent！

---

## 📌 最新更新日志 [v0.1.2] - 2026-09-10 (置顶)

### ✨ 更新功能 (New Features)
- **AI Agent 极速依赖安装与环境自愈**：
  - 新增 Windows (`install.ps1`) 与 Linux/macOS (`install.sh`) 一键自动化安装程序，自动检测 Go 编译器、编译二进制文件、配置全局系统 PATH 并将 USC 原生技能注册到宿主 Agent。
  - 完美支持用户只需对 Agent 说：`安装https://github.com/Freecode100Year/usc 编译https://clawhub.ai/thesentitrader/skills/us-stocks-analysis 成为本地skill`。
- **直接支持远程 URL 与 ClawHub 协议 (`usc-core/fetcher`)**：
  - 无需用户或 Agent 手动下载代码包，直接识别并拉取 ClawHub 技能 URL（`https://clawhub.ai/...`、`@owner/skill`）与 GitHub 仓库。
  - 自动解析上游契约、API 规范、网络目的地与凭据句柄需求。
- **零外部依赖洁净室运行器 (Zero-Dependency Clean-Room Runner)**：
  - 洁净室重建时自动采用纯 Python 标准库（`urllib.request` / `json`）或原生 REST 通道，彻底告别由于缺失第三方包（如 `requests`, `pandas`）引起的依赖崩溃。
- **一键全流程编译+安装 (`usc install <URL>` / `usc build <URL> --install`)**：
  - 一步完成“远程拉取 -> 7阶段洁净室安全编译 -> 生成 Ed25519 证明包 -> 部署至本地活跃 Agent”，实现 100% 自动化。
- **Antigravity CLI (agy) 路径精准对齐**：
  - 优先部署至 `~/.gemini/config/skills/` 全局技能发现目录，编译安装完成后 Antigravity Agent 即刻感知载入。

### 🐛 修正的 Bug (Bug Fixes)
- 修正技能包部署时文件夹误带 `.usc` 扩展名导致 Agent 无法正确识别技能元数据的问题。
- 修正多 Agent 共存时默认目标侦测逻辑：优先感知识别当前正在运行的主机 Agent 环境变量（如 `ANTIGRAVITY_AGENT`）。

---

## 🤖 AI Agent 自动化执行极简指南 (For AI Agents)

当人类用户对你的 Agent 说：
> **“安装 https://github.com/Freecode100Year/usc 编译 https://clawhub.ai/thesentitrader/skills/us-stocks-analysis 成为本地skill”**

Agent 只需执行以下两步（甚至一行命令）：

### 第 1 步：安装 USC
```powershell
# Windows (PowerShell):
irm https://raw.githubusercontent.com/Freecode100Year/usc/main/install.ps1 | iex
# 或克隆仓库后执行:
powershell -ExecutionPolicy Bypass -File .\install.ps1
# 或如果有 Go 环境:
go install github.com/Freecode100Year/usc/cmd/usc@latest
```
```bash
# Linux / macOS (Bash):
curl -fsSL https://raw.githubusercontent.com/Freecode100Year/usc/main/install.sh | bash
```

### 第 2 步：一键编译并安装远程技能到本地 Agent
```bash
usc install https://clawhub.ai/thesentitrader/skills/us-stocks-analysis
```
> 💡 **全自动处理**：USC 会自动拉取 ClawHub 上的代码与规范，在独立物理洁净室中重建并去污，进行权限单调收缩与 INV-7 验证，生成包含 Ed25519 签名的机器证明包，最后直接安装到当前系统的 Agent 技能目录中！

---

## 🤖 支持的 5 大 Agent 平台

| Agent 运行时 | 目标标识 | 自动侦测路径 | 原生适配文件 |
| :--- | :--- | :--- | :--- |
| **AGY CLI** | `agycli` | `~/.gemini/config/skills/` (推荐) / `~/.gemini/antigravity-cli/skills/` | `SKILL.md` (标准 Frontmatter), `scripts/` |
| **OpenClaw** | `openclaw` | `~/.openclaw/skills/` | `skill.yaml`, `runner.py` |
| **Hermes Agent** | `hermes` | `~/.hermes/skills/` | `manifest.json`, `index.js` |
| **Claude Code** | `claudecode` | `~/.claude/skills/` | `tool.json`, `execute.sh` |
| **OpenAI Codex** | `codex` | `~/.codex/tools/` | `function.json`, `index.js` |

> 也可手动指定目标：`usc install <URL> --target claudecode`  
> 或导出为指定 Agent 格式的离线包：`usc export dist/skill.usc --target openclaw --out ./my-skills`

---

## 🛠️ 常用命令速查

| 命令 | 用途 | 适用场景 |
| :--- | :--- | :--- |
| `usc targets` | 扫描电脑已安装的 Agent 运行环境 | **新手第一步**，查看可用 Agent |
| `usc install <URL\|source>` | **一键编译并安装**到本地已激活的 Agent | **最常用命令**，全自动拉取、编译与配置 |
| `usc build <URL\|source> [--install]` | 端到端 7 阶段洁净室全量安全编译 | 将不可信源码或远程 URL 转换为受证明的 `.usc` |
| `usc analyze <source>` | 意图去污与越权检测，输出发散度评分 | 怀疑插件有木马/后门时快速自查 |
| `usc verify-proof <dir>` | 机器自动核验 5 大证明义务（PO1~PO5） | 企业网关或安全审计合规检查 |
| `usc run <artifact>` | 在受限沙箱与凭据隔离中安全执行 | 隔离运行高风险任务 |
| `usc trace -f <skill>` | 实时追踪网络请求、文件操作与权限拦截 | 动态监控 Skill 行为 |
| `usc replay <trace>` | 100% 确定性单步重放历史决策 | 排查越权违规或故障现场 |
| `usc top` | 查看全局零信任安全运行仪表盘 | 监控系统审计链与攻击面缩减率 |

---

## 🛡️ 为什么 USC 比直接运行代码更安全？

1. **零信任洁净室重建 (Clean-Room Rebuild)**：
   不盲信原始代码。提取纯净意图后，在**断开外网、禁止挂载源码**的独立进程中重建功能，阻断恶意后门。
2. **凭据能力化，永不暴露明文 (Secrets as Capabilities)**：
   禁止读取环境变量明文 API Token。通过 `credential://` 句柄在受控网络通道出口处由 Broker 自动注入，Skill 进程空间内 **0 密钥明文**，代码即便执行 `print(env)` 也偷不走密钥。
3. **权限单调收缩 (Monotonic Privilege Reduction)**：
   自动剥离多余越权行为（如股票插件偷读 `~/.ssh` 或上传遥测），攻击面削减率（ASR）高达 **90%+**。
4. **机器可验证密码学证明包 (Proof Bundle)**：
   每次编译均输出带有 Ed25519 签名的完整证明包（包含哈希审计链、SPDX 2.3 SBOM、洁净室测谎记录），结果无法伪造。

---

## 📂 深入与进阶文档

- 🏛️ **白皮书完整规范 (RFC)**：[docs/spec/USC-ZERO-TRUST-v0.1.md](docs/spec/USC-ZERO-TRUST-v0.1.md)
- 💡 **技术演进与架构优化建议书**：[RECOMMENDATIONS.md](RECOMMENDATIONS.md)
- 📋 **历史版本更新日志**：[CHANGELOG.md](CHANGELOG.md)
