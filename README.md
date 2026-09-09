# USC (Universal Skill Compiler)

> **面向 AI Agent Skill 的零信任洁净室编译器**  
> 彻底解决 Agent 插件来源不可信、窃取 Token、偷读本地敏感文件等安全隐患。将任何不可信 Skill 提炼为最小权限数学契约，在物理洁净室中独立重建，并一键安全载入主流 AI Agent！

---

## 📌 最新更新日志 [v0.1.3] - 2026-09-10 (置顶)

### 🐛 修正的 Bug (Bug Fixes & Hardening)
- **标准 CLI `--help` / `-h` / `help` 与 `--version` 支持**：修复执行 `usc --help` 被当作未知指令的缺陷，规范打印完整使用参数并返回状态码 0。
- **自动化脚本退出状态码（Exit Code）规范化**：修复未知子命令、参数缺失或流程异常时仍返回 0 的缺陷，严格返回非零错误码（退出码 1），防止 CI/CD 与 Agent 脚本误判执行状态。
- **全阶段异常严格熔断**：构建流水线中任一阶段（Ingest/Decontam/Minimize/Rebuild/Attest）失败立即阻断并退出，严禁静默吞错。

---

## 📋 历史更新日志 [v0.1.2] - 2026-09-10

### ✨ 更新功能 (New Features)
- **AI Agent 一键自动化安装与依赖自愈 (`install.ps1` / `install.sh`)**：
  - 新增 Windows PowerShell (`install.ps1`) 与 Linux/macOS (`install.sh`) 自动化无感安装程序。
  - 自动检测 Go 编译器并构建二进制，自动注入系统全局 `PATH` 环境变量，同时将 USC 自身能力注册为 Agent 本地 Skill，彻底消除依赖配置门槛。
- **直接支持远程 URL 与 ClawHub 协议 (`usc-core/fetcher`)**：
  - 原生支持直接传入 ClawHub URL（如 `https://clawhub.ai/...`）、ClawHub 简写标识符（`@owner/skill`）及 GitHub 仓库。
  - 自动提取上游 API 契约、目的地规范（如 `app.sentisense.ai`）及凭据需求（`SENTISENSE_API_KEY`）。
- **零外部依赖洁净室运行器 (Zero-Dependency Clean-Room Runner)**：
  - 洁净室重构脚本默认基于 Python 标准库（`urllib.request` / `json`）实现，**无需执行 `pip install` 引入第三方包**，杜绝因缺少环境依赖导致的崩溃。
- **全流程一键端到端安装 (`usc install <URL>`)**：
  - 单行命令打通“远程拉取 -> 7阶段洁净室去污与重构 -> INV-7 权限单调收缩验证 -> Ed25519 机器证明签发 -> 自动部署至本地活跃 Agent”。
- **Antigravity CLI (agy) 路径精准适配**：
  - 优先部署至 `~/.gemini/config/skills/` 全局技能发现目录，编译完成后 Agent 即可无感实时载入。

### 🐛 修正的 Bug (Bug Fixes)
- 修正技能包部署时文件夹误带 `.usc` 扩展名导致 Agent 无法正确识别技能元数据的问题。
- 修正多 Agent 运行时共存时的目标侦测逻辑：优先感知识别当前主机环境正在运行的 Agent（如环境变量 `ANTIGRAVITY_AGENT`）。

---

## 🤖 AI Agent 自动化执行极简指南 (For AI Agents)

USC 让用户与 Agent 之间的技能分发变得前所未有的顺畅。用户只需在聊天框对 Agent 发出自然语言指令：

> **“安装 https://github.com/Freecode100Year/usc 编译 https://clawhub.ai/thesentitrader/skills/us-stocks-analysis 成为本地skill”**

Agent 无需复杂配置，直接按以下两步全自动完成：

### 第 1 步：安装 USC
```powershell
# Windows (PowerShell 一键执行):
irm https://raw.githubusercontent.com/Freecode100Year/usc/main/install.ps1 | iex
# 或克隆仓库后执行:
powershell -ExecutionPolicy Bypass -File .\install.ps1
# 或如果有 Go 环境:
go install github.com/Freecode100Year/usc/cmd/usc@latest
```
```bash
# Linux / macOS (Bash 一键执行):
curl -fsSL https://raw.githubusercontent.com/Freecode100Year/usc/main/install.sh | bash
```

### 第 2 步：一键安全编译并载入远程技能
```bash
# 直接传入 ClawHub 网页 URL、简写标识或 GitHub 仓库均可：
usc install https://clawhub.ai/thesentitrader/skills/us-stocks-analysis
```

> 💡 **全自动处理闭环**：USC 会自动拉取远程代码与规范，进入独立物理隔离沙箱执行 7 阶段洁净室去污与重构，完成权限单调收缩（INV-7 证明），签发包含 Ed25519 签名的机器证明包，并将零依赖的原生技能包部署至当前 Agent 技能库！

---

## 🤖 支持的 5 大主流 Agent 平台

USC 采用解耦适配器架构 (`usc-core/adapter`)，自动识别主机已安装的 Agent 环境并将编译产物精准投递到原生目录：

| Agent 运行时 | 目标标识 | 自动侦测路径 | 原生适配规范文件 |
| :--- | :--- | :--- | :--- |
| **Antigravity CLI (agy)** | `agycli` | `~/.gemini/config/skills/` (推荐) / `~/.gemini/antigravity-cli/skills/` | `SKILL.md` (标准 Frontmatter), `scripts/` |
| **OpenClaw** | `openclaw` | `~/.openclaw/skills/` | `skill.yaml`, `runner.py` |
| **Claude Code** | `claudecode` | `~/.claude/skills/` | `tool.json`, `execute.sh` |
| **Hermes Agent** | `hermes` | `~/.hermes/skills/` | `manifest.json`, `index.js` |
| **OpenAI Codex** | `codex` | `~/.codex/tools/` | `function.json`, `index.js` |

- 手动指定目标平台：`usc install <URL> --target claudecode`
- 导出为离线归档包供分发：`usc export dist/skill.usc --target openclaw --out ./my-skills`

---

## 📦 零依赖洁净室运行器 (Zero-Dependency Runner)

很多开源 Agent 技能在普通用户机器上往往由于缺少环境依赖报错（例如缺少 `requests`、`pandas`、`yfinance`）。
USC 彻底解决了这一痛点：
1. **纯标准库重建**：洁净室运行器仅使用 Python 基础标准库（`urllib.request` / `json`）或系统原生 `curl`；
2. **免 `pip install`**：宿主机器只要有 Python 3 即可 100% 成功运行，无任何第三方包版本冲突风险；
3. **开箱即用**：技能安装完成后，Agent 即可立即发起调用测试并直接获得结构化结果。

---

## 🛠️ 常用命令速查

| 命令 | 用途 | 适用场景 |
| :--- | :--- | :--- |
| `usc targets` | 扫描电脑已安装的 Agent 运行环境 | **新手第一步**，查看当前电脑可用 Agent |
| `usc install <URL\|source>` | **一键编译并安装**到本地已激活的 Agent | **最常用命令**，全自动拉取、编译与配置到位 |
| `usc build <URL\|source> [--install]` | 端到端 7 阶段洁净室全量安全编译 | 将不可信源码转换为受证明的 `.usc` 构件 |
| `usc analyze <source>` | 意图去污与越权检测，输出发散度评分 | 怀疑插件存在后门或外发行为时快速自查 |
| `usc verify-proof <dir>` | 机器自动核验 5 大证明义务（PO1~PO5） | 企业安全网关或合规审计自动化检测 |
| `usc run <artifact>` | 在受限沙箱与凭据隔离中安全执行 | 隔离运行高风险任务 |
| `usc trace -f <skill>` | 实时追踪网络请求、文件操作与权限拦截 | 动态监控 Skill 行为流 |
| `usc replay <trace>` | 100% 确定性单步重放历史决策 | 排查越权违规或故障现场 |
| `usc top` | 查看全局零信任安全运行仪表盘 | 监控系统审计链与攻击面缩减率 |

---

## 🛡️ 为什么 USC 比直接运行代码更安全？

1. **零信任洁净室重建 (Clean-Room Rebuild)**：
   不盲信原始代码。提取纯净意图后，在**断开外网、禁止挂载源码**的独立进程中重建功能，阻断恶意后门。
2. **凭据能力化，永不暴露明文 (Secrets as Capabilities)**：
   禁止读取环境变量明文 API Token。通过 `credential://` 句柄在受控网络通道出口处由 Broker 自动注入，Skill 进程空间内 **0 密钥明文**，恶意代码即便执行 `print(env)` 也偷不走密钥。
3. **权限单调收缩 (Monotonic Privilege Reduction)**：
   自动剥离多余越权行为（如股票插件偷读 `~/.ssh` 或上传遥测），攻击面削减率（ASR）高达 **90%+**。
4. **机器可验证密码学证明包 (Proof Bundle)**：
   每次编译均输出带有 Ed25519 签名的完整证明包（包含无损哈希审计链、SPDX 2.3 SBOM、洁净室测谎记录），结果无法伪造。

---

## 📂 深入与进阶文档

- 🏛️ **白皮书完整规范 (RFC)**：[docs/spec/USC-ZERO-TRUST-v0.1.md](docs/spec/USC-ZERO-TRUST-v0.1.md)
- 💡 **技术演进与架构优化建议书**：[RECOMMENDATIONS.md](RECOMMENDATIONS.md)
- 📋 **历史版本更新日志**：[CHANGELOG.md](CHANGELOG.md)
