# USC (Universal Skill Compiler)

> **面向 AI Agent Skill 的零信任洁净室编译器**  
> 彻底解决 Agent 插件来源不可信、窃取 Token、偷读本地文件等安全风险。将任何不可信 Skill 提炼为最小权限契约，在物理洁净室中独立重建，并一键载入你的主流 Agent！

---

## ⚡ 30 秒极简上手（新手向）

无论是从 GitHub、ClawHub 下载的脚本，还是社区分享的压缩包，新手只需两步即可安全载入：

```bash
# 1. 下载或编译 USC
go build -o bin/usc.exe ./cmd/usc

# 2. 检查你电脑安装了哪些 Agent (支持 OpenClaw / Hermes / Claude Code / AGY CLI / Codex)
usc targets

# 3. 一键安全编译并自动装入你的 Agent 技能库
usc build ./untrusted-skill
usc install dist/untrusted-skill.usc
```

> 💡 **自动装载**：`usc install` 会自动识别你电脑已安装的 Agent 软件，将其转换为原生插件格式并放入对应文件夹，**无需任何手动配置**。

---

## 🤖 支持的 5 大 Agent 平台

| Agent 运行时 | 目标标识 | 自动侦测路径 | 原生适配文件 |
| :--- | :--- | :--- | :--- |
| **OpenClaw** | `openclaw` | `~/.openclaw/skills/` | `skill.yaml`, `runner.py` |
| **Hermes Agent** | `hermes` | `~/.hermes/skills/` | `manifest.json`, `index.js` |
| **Claude Code** | `claudecode` | `~/.claude/skills/` | `tool.json`, `execute.sh` |
| **AGY CLI** | `agycli` | `~/.gemini/antigravity-cli/skills/` | `SKILL.md` (标准 Frontmatter) |
| **OpenAI Codex** | `codex` | `~/.codex/tools/` | `function.json`, `index.js` |

> 也可以手动指定目标：`usc install dist/skill.usc --target claudecode`  
> 或者导出为离线包供分享：`usc export dist/skill.usc --target openclaw --out ./my-skills`

---

## 🛠️ 常用命令速查

| 命令 | 用途 | 适用场景 |
| :--- | :--- | :--- |
| `usc targets` | 扫描电脑已安装的 Agent 运行环境 | **新手第一步**，查看可用 Agent |
| `usc install <artifact>` | **一键安装**到本地已激活的 Agent | **日常使用**，全自动配置 |
| `usc build <source>` | 端到端 7 阶段洁净室全量安全编译 | 将不可信源码转换为受证明的 `.usc` |
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
   自动剥离多余越权行为（如天气插件偷读 `~/.ssh` 或上传遥测），攻击面削减率（ASR）高达 **90%+**。
4. **机器可验证密码学证明包 (Proof Bundle)**：
   每次编译均输出带有 Ed25519 签名的完整证明包（包含哈希审计链、SPDX 2.3 SBOM、洁净室测谎记录），结果无法伪造。

---

## 📂 深入与进阶文档

- 🏛️ **白皮书完整规范 (RFC)**：[docs/spec/USC-ZERO-TRUST-v0.1.md](docs/spec/USC-ZERO-TRUST-v0.1.md)
- 💡 **技术演进与架构优化建议书**：[RECOMMENDATIONS.md](RECOMMENDATIONS.md)
- 📋 **历史版本更新日志**：[CHANGELOG.md](CHANGELOG.md)
