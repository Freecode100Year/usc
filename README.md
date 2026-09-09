# USC (Universal Skill Compiler)

## 📌 最新更新日志 (Changelog)

### [v0.1.1] - 2026-09-10

#### ✨ 更新功能 (New Features)
- **多 Agent 原生生态适配与新手一键载入 (`usc-core/adapter`)**：
  - 全面支持 **OpenClaw**、**Hermes Agent**、**Claude Code**、**AGY CLI (Antigravity CLI)**、**Codex (OpenAI Assistants)** 五大异构 Agent 运行时。
  - **新手一键自动侦测与加载 (`usc install`)**：面向电脑新手用户，无需手动编辑配置或查找路径，自动嗅探系统已安装的 Agent Runtime，一键将通过零信任审核的 Skill 部署到位。
  - **原生格式导出 (`usc export`)**：支持将编译产物独立导出为各生态的原生规范目录结构（OpenClaw `skill.yaml`、Hermes `manifest.json`、Claude Code `tool.json`、AGY CLI `SKILL.md`、Codex `function.json`）。
  - **环境状态自检 (`usc targets`)**：清晰可视化列出 5 大平台在当前计算机上的安装识别状态与技能存储路径。

---

### [v0.1.0] - 2026-09-10

#### ✨ 更新功能 (New Features)
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

#### 🐛 修复与加固 (Bug Fixes & Hardening)
- 修复通配符路径匹配越界隐患：引入路径段 AST（Segment AST）严格对比，坚决杜绝路径穿越（Path Traversal）。
- 杜绝任意通配操作符扩张：采用防御性已知名单集合 `ActionAllKnown`，禁止未知动作隐式放行。

---

# 📖 项目概述与核心信条

> **Trust the contract, not the code.**（信任契约，而不是代码。）  
> **Transfer capability, not trust.**（传递能力，不传递信任。）  
> **Probabilistic understanding, deterministic enforcement.**（概率式理解，确定性执行。）  
> **Enforce what you declare. Observe what you enforce. Prove what you observe.**（声明什么就强制什么；强制什么就观测什么；观测什么就证明什么。）

USC（Universal Skill Compiler）是一个**针对不可信意图、未知实现和异构 Agent Runtime 的零信任洁净室重建编译器**。它将不可信 Agent Skill 分解为意图、行为与权限，将其压缩为最小可审查信任契约，再针对不同 Agent Runtime 独立重建，并对编译过程、生成产物和运行时行为提供机器可验证的证明。

---

# 🌟 新手小白快速上手指南 (Beginner's Guide)

大部分用户可能并非底层安全专家或系统极客，USC 专门针对电脑新手提供了**全自动侦测与一键加载**方案，只需三步即可在主流 AI Agent 中安全使用任何 Skill：

```text
[不可信 Skill 源码] 
       │ 1. 一键编译并完成零信任洁净室审查
       ▼
  usc build ./my-skill
       │ 2. 自动生成 .usc 制品与机器证明包
       ▼
  usc install dist/my-skill.usc
       │ 3. 自动嗅探你的电脑安装了哪个 Agent，直接部署到位！
       ▼
[OpenClaw / Hermes / Claude Code / AGY CLI / Codex 立即可用！]
```

---

## 🎯 电脑新手常用操作（2 条命令搞定）

### 第一步：检查你的电脑已安装哪些 Agent
运行：
```bash
usc targets
```
**终端将自动扫描并显示**：
```text
================================================================
          SUPPORTED AGENT RUNTIMES & LOCAL DETECTION            
================================================================
 • openclaw     : OpenClaw AI Runtime        [DETECTED: READY]
   Path: C:\Users\yourname\.openclaw\skills
 • hermes       : Hermes Autonomous Agent    [NOT DETECTED]
   Path: C:\Users\yourname\.hermes\skills
 • claudecode   : Claude Code CLI            [DETECTED: READY]
   Path: C:\Users\yourname\.claude\skills
 • agycli       : Antigravity CLI (agy)      [DETECTED: READY]
   Path: C:\Users\yourname\.gemini\antigravity-cli\skills
 • codex        : OpenAI Codex / Assistants  [DETECTED: READY]
   Path: C:\Users\yourname\.codex\tools
================================================================
Tip for beginners: Run 'usc install <artifact.usc>' to auto-load!
```

---

### 第二步：一键自动安装到 Agent 中
编译后，直接运行：
```bash
usc install dist/weather-skill.usc
```
USC 会自动检测本地已激活的 Agent 运行时，将格式转换为原生规范并自动放入对应的技能文件夹中，无需任何手动复杂配置！

---

# 🤖 支持的 5 大主流 Agent 运行时详情

| 平台名称 | 标识符 (`--target`) | 自动侦测路径 | 生成的原生适配物 |
| :--- | :--- | :--- | :--- |
| **OpenClaw** | `openclaw` | `~/.openclaw/skills` | `skill.yaml`, `runner.py`（沙箱隔离驱动） |
| **Hermes Agent** | `hermes` | `~/.hermes/skills` | `manifest.json`, `index.js`（受限入口） |
| **Claude Code** | `claudecode` | `~/.claude/skills` | `tool.json`（Tool Use Schema）, `execute.sh` |
| **AGY CLI** | `agycli` | `~/.gemini/antigravity-cli/skills` | `SKILL.md`（标准 YAML Frontmatter 与元数据） |
| **OpenAI Codex** | `codex` | `~/.codex/tools` | `function.json`（Function Calling 契约）, `index.js` |

如果你希望明确指定安装到某一个 Agent，只需带上 `--target` 参数：
```bash
# 安装到 Claude Code
usc install dist/weather-skill.usc --target claudecode

# 安装到 OpenClaw
usc install dist/weather-skill.usc --target openclaw

# 安装到 Google Antigravity CLI (agy)
usc install dist/weather-skill.usc --target agycli

# 安装到 Hermes Agent
usc install dist/weather-skill.usc --target hermes

# 安装到 OpenAI Codex / Assistants
usc install dist/weather-skill.usc --target codex
```

---

### 📦 手动导出原生技能包（离线分发或二次分享）

如果你想将生成的技能打包发给其他人，使用 `export` 命令：
```bash
usc export dist/weather-skill.usc --target openclaw --out ./my-exported-skills
```
将在指定目录下生成可以直接复制使用的目标平台技能包。

---

# 🛠️ 进阶命令手册 (CLI Full Reference)

USC 提供完整的 13 个核心命令行指令：

| 指令 | 作用说明 | 典型场景 |
| :--- | :--- | :--- |
| `usc analyze <source>` | 意图去污与发散度分析 | 快速审查 GitHub / ClawHub 下载的不可信插件是否有后门 |
| `usc extract <source>` | 提取平台中立标准意图 IR | 剥离所有提示词注入攻击，生成纯净契约 |
| `usc rebuild <blueprint>` | 物理洁净室代码重建 | 隔绝网络和原始文件，由洁净室重新生成代码 |
| `usc build <source>` | 端到端 7 阶段全量编译 | 生成带有机器证明包的 `.usc` 二进制制品 |
| `usc verify <artifact>` | 验证制品签名与哈希 | 上线部署前验证制品完整性与合规性 |
| `usc verify-proof <proof-dir>`| 机器验证证明包全链 | 机器自动核验 5 大证明义务（$PO_1 \sim PO_5$） |
| `usc run <artifact>` | 受限沙箱与代理中介运行 | 安全执行 Skill，切断明文环境变量泄露风险 |
| `usc trace -f <skill-id>` | 实时运行态能力追踪 | 动态观察网络请求、文件访问与权限拦截 |
| `usc replay <trace>` | 确定性决策重放 | 发生违规或故障时，100% 逐步重放安全判定过程 |
| `usc top` | 全局安全仪表盘 | 查看当前系统运行状态、审计链与攻击面缩减率 |
| `usc targets` | 探测已支持的 Agent 环境 | 新手查看当前电脑具备哪些可用的 Agent 运行时 |
| `usc export <artifact>` | 导出为目标 Agent 原生技能 | 生成对应 Agent 原生目录结构供离线导入 |
| `usc install <artifact>` | 一键安装到目标 Agent | **新手最爱**，免配置自动部署到 Agent 技能库 |

---

### 机器证明包目录结构 (Machine Proof Bundle)

每次成功执行 `usc build` 后，系统将在 `dist/proof/` 目录下生成完整的机器可验证证明集合：

```text
dist/
├── weather-skill.usc         # 洁净室独立重建的沙箱二进制制品
└── proof/
    ├── attestation.json      # Ed25519 签名的机器证明元数据文件
    ├── audit-chain.json      # RFC 8785 JCS 规范化的无损哈希链全量账本
    ├── blueprint.json        # 冻结的最小能力蓝图（Minimal Capability Blueprint）
    ├── cleanroom-proof.json  # 物理洁净室测谎数据（0 源码泄漏、外网断开）
    ├── provenance.json       # 严格因果事实与转换溯源链
    ├── policy.json           # 编译采用的企业安全约束策略
    └── sbom.spdx.json        # SPDX 2.3 标准软件物料清单
```

---

# 📚 文档与技术规范导航
- 🏛️ **白皮书完整规范**：[docs/spec/USC-ZERO-TRUST-v0.1.md](docs/spec/USC-ZERO-TRUST-v0.1.md)
- 💡 **架构演进与优化建议书**：[RECOMMENDATIONS.md](RECOMMENDATIONS.md)
