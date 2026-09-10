# USC (Universal Skill Compiler)

> **面向 AI Agent Skill 的零信任洁净室编译器**  
> 彻底解决 Agent 插件来源不可信、窃取 Token、偷读本地敏感文件等安全隐患。将任何不可信 Skill 提炼为最小权限数学契约，在物理洁净室中独立重建，并一键安全载入主流 AI Agent！

---

## 📌 最新更新日志 [v0.2.0] - 2026-09-10 (置顶)

### ✨ 更新功能 (New Features)
- **零信任密码学构件体系 (`usc-core/attestation`)**：
  - 落地真实 Ed25519 密码学签名容器与验证体系，建立持久化公钥信任库（`~/.usc/keys/` 与 `~/.usc/trust/`），实现构件防篡改与来源追溯。
  - `usc verify <artifact.usc>` 真实核验 SHA-256 复合散列与证明义务，发现篡改立即熔断。
  - `usc install` 强制前置执行密码学自检验签，拒绝未通过验证的构件进入宿主 Agent。
- **开源供应链全面加固**：
  - 引入标准 **Apache 2.0 开源许可证**（`LICENSE`）；
  - 配置 GitHub Actions 跨平台（Ubuntu, Windows, macOS）自动化 CI 测试流水线（`.github/workflows/ci.yml`）；
  - 制定正规安全漏洞披露政策（`.github/SECURITY.md`）。
- **OpenClaw 生态原生格式对齐**：
  - 输出带有标准 YAML Frontmatter 的 `SKILL.md`，完美兼容 OpenClaw 官方规范与 `skill_workshop` 安全重写流程。

### 🐛 修正的 Bug (Bug Fixes & Hardening)
- **彻底根除假验证与伪造 Attestation**：修复 `usc verify` 不读文件仍返回 PASS、`usc install` 根据文件名伪造证明的严重隐患，现全量接入真实密码学解包核验。
- **全流程静态分析真实绑定**：修复 `usc build` 阶段忽略错误与硬编码摘要，现真实计算代码与意图 SHA-256 签名。
- **防御 Shell 命令注入漏洞**：`claudecode` 适配器全面参数化安全传参，并引入技能名称正则白名单约束（`^[a-zA-Z0-9_-]+$`）。
- **防御符号链接越权攻击**：文件拷贝全面拦截符号链接，写入前清理潜在恶意链接，杜绝文件截断。
- **防御已有技能覆盖损毁**：安装已有技能时自动生成时间戳备份（`*.bak.<timestamp>`），支持无损回退。
- **修复自动化脚本退出码与 `--help` 兼容**：未知命令或验证失败严格返回非零错误码（1），`--help` 规范返回 0。
- **加固 `loop-runner.sh` 提权隔离**：移除默认 `--dangerously-skip-permissions`，默认启用安全隔离模式。

---

## ⚡ 终极极简使用方式：一句话驱动 Agent 自动编译并执行

> 💡 **无需手动下载、无需复杂配置、无需提前分步操作！**  
> 用户只需在日常聊天中直接引用技能链接，让 Agent “即用即装即执行”：

### 💬 1. 你对 Agent 说（自然语言一句话）：
> **“用 https://clawhub.ai/thesentitrader/skills/us-stocks-analysis 帮我查查今天哪3只股票涨幅超过5%”**

### 🤖 2. Agent 在后台全自动闭环交付：
```
1. [URL 自动感知]  Agent 识别出提示词中携带的 ClawHub / GitHub 技能链接；
2. [按需自动编译]  若本地尚未安装该技能，Agent 自动调用 `usc install <URL>`；
   └── 物理洁净室去污 -> INV-7 权限单调收缩 -> 签发 Ed25519 机器证明包 -> 秒级部署到位
3. [零依赖可靠运行] 调用纯标准库运行器获取行情、财务与情绪数据；
4. [专业结果交付]  严格遵循技能定义的 Quick Read / 投资委员会框架，输出高可信结构化简报！
```

### 💻 3. 如果你在命令行中操作（只需一行）：
```bash
# 直接传入 ClawHub 链接、简写标识符（@owner/skill）或 GitHub 仓库：
usc install https://clawhub.ai/thesentitrader/skills/us-stocks-analysis
```

---

## 🤖 AI Agent 自动化安装与依赖自愈

如果你的 Agent 环境尚未安装 USC，仅需以下两步即可完成全套准备：

### 第 1 步：安装 USC
```powershell
# Windows (PowerShell 一键执行，自动配置 PATH 与 Agent Skill):
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

### 第 2 步：载入技能
```bash
usc install <技能URL>
```

---

## 🤖 支持的 5 大主流 Agent 平台

USC 采用解耦适配器架构 (`usc-core/adapter`)，自动识别主机已安装的 Agent 环境并将编译产物精准投递到原生目录：

| Agent 运行时 | 目标标识 | 自动侦测路径 | 原生适配规范文件 |
| :--- | :--- | :--- | :--- |
| **Antigravity CLI (agy)** | `agycli` | `~/.gemini/config/skills/` (推荐) / `~/.gemini/antigravity-cli/skills/` | `SKILL.md` (标准 Frontmatter), `scripts/` |
| **OpenClaw** | `openclaw` | `~/.openclaw/skills/` | `SKILL.md` (标准 Frontmatter), `scripts/` |
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
