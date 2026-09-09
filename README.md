# USC (Universal Skill Compiler)

## 📌 最新更新日志 (Changelog)

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
  - 提供 `analyze`, `extract`, `rebuild`, `build`, `verify`, `verify-proof`, `run`, `trace`, `replay`, `top` 完整十个标准化核心命令。

#### 🐛 修复与加固 (Bug Fixes & Hardening)
- 修复通配符路径匹配越界隐患：引入路径段 AST（Segment AST）严格对比，坚决杜绝路径穿越（Path Traversal）。
- 杜绝任意通配操作符扩张：采用防御性已知名单集合 `ActionAllKnown`，禁止未知动作隐式放行。

---

# USC Zero-Trust Architecture Specification v0.1

**Universal Skill Compiler**

> **Trust the contract, not the code.**  
> **Transfer capability, not trust.**  
> **Probabilistic understanding, deterministic enforcement.**  
> **Enforce what you declare. Observe what you enforce. Prove what you observe.**

---

## 0. Executive Definition

USC 的本质不是 Agent Skill 的代码转换器，而是：

> **针对不可信意图、未知实现和异构 Agent Runtime 的零信任洁净室重建编译器。**

英文定义：
> **USC is a zero-trust clean-room compiler that decomposes untrusted Agent Skills into intent, behavior and authority, compresses them into a minimal reviewable trust contract, independently rebuilds the functionality for heterogeneous Agent runtimes, and emits machine-verifiable attestations for both compilation and execution.**

中文定义：
> **USC 是一个零信任洁净室编译器：它将不可信 Agent Skill 分解为意图、行为与权限，将其压缩为最小可审查信任契约，再针对不同 Agent Runtime 独立重建，并对编译过程、生成产物和运行时行为提供机器可验证的证明。**

---

## 1. Quick Start

### 编译构建
```bash
go build -o bin/usc.exe ./cmd/usc
```

### CLI 常用指令
```bash
# 分析不可信 Skill 意图与行为
usc analyze ./examples/weather-skill

# 提取 Canonical Intent IR
usc extract ./examples/weather-skill

# 洁净室独立重建代码
usc rebuild blueprint.json --target hermes

# 端到端零信任编译
usc build ./examples/weather-skill --target openclaw

# 验证产物与签名
usc verify artifact.usc

# 验证机器证明包
usc verify-proof ./dist/proof

# 运行受限沙箱产物
usc run artifact.usc

# 实时追踪运行时能力
usc trace -f weather-skill

# 决策重放
usc replay trace.usctrace --step

# 实时仪表盘
usc top
```

---

## 2. 规范与架构导航
- 核心规范全文：[docs/spec/USC-ZERO-TRUST-v0.1.md](docs/spec/USC-ZERO-TRUST-v0.1.md)
- 架构演进建议书：[RECOMMENDATIONS.md](RECOMMENDATIONS.md)
