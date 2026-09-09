# Changelog

All notable changes to the USC (Universal Skill Compiler) project will be documented in this file.

## [v0.1.3] - 2026-09-10

### 🐛 修复与加固 (Bug Fixes & Hardening)
- **标准 CLI `--help` / `-h` / `help` 与 `--version` 标志支持**：
  - 修复执行 `usc --help` 被误识别为 `Unknown command` 的问题，现可正确输出完整使用指引并以退出状态码 0 正常结束。
  - 新增 `usc --version` / `-v` 原生版本输出。
- **自动化脚本退出状态码（Exit Code）规范化**：
  - 修复未知命令、非法参数或编译/安装流程失败时仍返回成功状态码 0 的严重隐患。
  - 当收到未知指令或流程异常时严格返回退出码 1，确保 CI/CD 自动化流水线与自主 Agent 脚本能够准确判别执行成败。
- **全阶段异常熔断拦截**：
  - 在 `RunIngest`、`RunDecontaminate`、`RunMinimize`、`RunReAudit`、`RunSandbox` 及 `RunAttest` 遇到严重违规（如越权或 Attestation 异常）时，立即终端熔断并返回非零状态码，禁止静默忽略错误。

---

## [v0.1.2] - 2026-09-10

### ✨ 更新功能 (New Features)
- **AI Agent 一键自动化安装脚本 (`install.ps1` / `install.sh`)**：
  - 支持 Windows PowerShell 与 Linux/macOS 纯单行全自动安装。
  - 自动侦测并编译 Go 源码、注入系统 PATH 环境变量、并将 USC 原生能力注册为 Agent Skill。
- **直接支持远程 URL 与 ClawHub 协议 (`usc-core/fetcher`)**：
  - 推出统一 Fetcher 模块，自动解析 ClawHub 技能 URL（`https://clawhub.ai/...`、`@owner/skill`）及 GitHub 仓库。
  - 自动提取技能意图、API 协议元数据、目的地规范及权限契约。
- **零外部依赖洁净室运行器 (Zero-Dependency Clean-Room Runner)**：
  - 重建代码默认采用 Python 标准库与无外挂依赖设计，杜绝在宿主 Agent 环境下因缺少第三方库而报错。
- **全流程一键编译与自动加载 (`usc install <URL>`)**：
  - 打通“远程拉取 -> 洁净室编译去污 -> INV-7 权限单调收缩 -> 机器证明包签发 -> 自动部署至活跃 Agent”端到端闭环。
- **Antigravity CLI (agy) 路径精准适配**：
  - 技能安装首选对齐 `~/.gemini/config/skills/`，实现 Antigravity Agent 即时载入。

### 🐛 修复与加固 (Bug Fixes & Hardening)
- 修复部署目录误包含 `.usc` 后缀的问题，规范统一为标准语义技能名称。
- 修复在包含多个 Agent 的系统环境中，自动感知识别宿主 Agent 环境变量（`ANTIGRAVITY_AGENT` 等）的优先级机制。

---

## [v0.1.1] - 2026-09-10

### ✨ 更新功能 (New Features)
- **多 Agent 原生生态适配与新手一键载入 (`usc-core/adapter`)**：
  - 全面支持 **OpenClaw**、**Hermes Agent**、**Claude Code**、**AGY CLI (Antigravity CLI)**、**Codex (OpenAI Assistants)** 五大异构 Agent 运行时。
  - **新手一键自动侦测与加载 (`usc install`)**：面向电脑新手用户，无需手动编辑配置或查找路径，自动嗅探系统已安装的 Agent Runtime，一键将通过零信任审核的 Skill 部署到位。
  - **原生格式导出 (`usc export`)**：支持将编译产物独立导出为各生态的原生规范目录结构（OpenClaw `skill.yaml`、Hermes `manifest.json`、Claude Code `tool.json`、AGY CLI `SKILL.md`、Codex `function.json`）。
  - **环境状态自检 (`usc targets`)**：清晰可视化列出 5 大平台在当前计算机上的安装识别状态与技能存储路径。

---

## [v0.1.0] - 2026-09-10

### ✨ 更新功能 (New Features)
- **零信任洁净室编译器架构落地**：完整实现 USC Zero-Trust Architecture Specification v0.1 核心体系。
- **能力偏序半格与代数计算 (`usc-core/capability`)**：
  - 实现结构化能力类型系统（网络、文件系统、凭据、进程/命令执行、系统调用等），防御性拒绝野生通配符。
  - 实现严格偏序公理验证（自反性、反对称性、传递性判定）。
  - 实现最大公共下界（Meet，$\sqcap$）与最小共同上界（Join，$\sqcup$）能力收敛计算。
  - 实现多维风险成本向量计算：$Cost(C) = (R_n, R_f, R_s, R_w, R_e, R_p, R_{persist})$ 以及能力削减率（CCR）与攻击面削减率（ASR）。
- **意图去污与必要行为分解 (`usc-core/intent`)**：
  - 建立三元分析模型：声明意图 $I_{decl}$、观测行为 $B_{obs}$ 与必要行为 $B_{nec}$。
  - 实现非必要越权行为提取 $B_{extra} = B_{obs} \setminus B_{nec}$ 与缺失行为检测 $B_{missing} = B_{nec} \setminus B_{obs}$。
  - 风险加权差异量化 $D_N = \sum w(c)$，超阈值自动执行 `HARD_BLOCK` 熔断拦截或 `STRIP` 剥离。
- **能力最小化求解器 (`usc-core/minimizer`)**：
  - 定义解耦求解器抽象接口 `Solver`。
  - 实现贪心求解器 `GreedySolver` 与约束求解器 `ConstraintSolver`，冻结 Minimal Capability Blueprint。
- **零信任不变量与证明义务校验 (`usc-core/verifier`)**：
  - 严密守护 INV-1 至 INV-9 九大零信任不变量。
  - 实现 INV-7 权限单调收缩验证器（$C_{out} \sqsubseteq C_{in}$，越权即 `HARD_BLOCK`）。
  - 实现 INV-9 洁净室物理隔离度量验证器（独立进程、无未授权挂载、零源码碎片泄漏）。
  - 自动化证明义务推演：$PO_1 \land PO_2 \land PO_3 \land PO_4 \land PO_5 \Rightarrow \text{ATTESTED}$。
- **无损审计总线与哈希账本 (`usc-core/audit`)**：
  - 严格双总线架构：无损审计总线（Lossless Audit Bus）与尽力而为 UI 总线（Best-Effort UI Bus）。
  - RFC 8785 JCS (Canonical JSON) 规范化确定性序列化。
  - 增量密码学哈希链：$H_n = \text{SHA256}(H_{n-1} \parallel \text{Canonical}(Event_n))$。
- **凭据与运行时能力中介 Broker (`usc-core/broker`)**：
  - 彻底切断明文环境变量，实行 `credential://` 能力句柄隔离。
  - 严格限制凭据仅在受限传输通道目的地直接注入，不可信代码地址空间永无 Secret 明文。
  - 进程边界与文件沙箱隔离守护。
- **白盒观测与决策重放 (`usc-core/replay`)**：
  - 环形缓冲飞行记录器（Flight Recorder）生成 `.usctrace`。
  - 确定性决策单步重放工具（Deterministic Decision Replay）。
- **机器可验证密码学证明包 (`usc-core/attestation`)**：
  - 生成并校验 Ed25519 签名的 `attestation.json` 与完备 Proof Bundle。
- **完备 CLI 工具链 (`cmd/usc`)**：
  - 提供 `analyze`, `extract`, `rebuild`, `build`, `verify`, `verify-proof`, `run`, `trace`, `replay`, `top`, `targets`, `export`, `install` 完整命令行。

### 🐛 修复与加固 (Bug Fixes & Hardening)
- 修复通配符路径匹配越界隐患：引入路径段 AST（Segment AST）严格对比，坚决杜绝路径穿越（Path Traversal）。
- 杜绝任意通配操作符扩张：采用防御性已知名单集合 `ActionAllKnown`，禁止未知动作隐式放行。
