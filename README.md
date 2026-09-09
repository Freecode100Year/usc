# USC (Universal Skill Compiler)

## 📌 最新更新日志 (Changelog)

### [v0.1.0] - 2026-09-10

#### ✨ 更新功能 (New Features)
- **文档与使用指南全面完善**：自述文件中完整补充 USC CLI 10 个子命令的详细参数、交互示例、工作流说明及机器证明包（Proof Bundle）检验指南。
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

# 📖 项目概述与核心信条

> **Trust the contract, not the code.**（信任契约，而不是代码。）  
> **Transfer capability, not trust.**（传递能力，不传递信任。）  
> **Probabilistic understanding, deterministic enforcement.**（概率式理解，确定性执行。）  
> **Enforce what you declare. Observe what you enforce. Prove what you observe.**（声明什么就强制什么；强制什么就观测什么；观测什么就证明什么。）

USC（Universal Skill Compiler）是一个**针对不可信意图、未知实现和异构 Agent Runtime 的零信任洁净室重建编译器**。它将不可信 Agent Skill 分解为意图、行为与权限，将其压缩为最小可审查信任契约，再针对不同 Agent Runtime 独立重建，并对编译过程、生成产物和运行时行为提供机器可验证的证明。

---

# 🚀 快速上手与使用方法 (Usage Guide)

### 1. 编译安装

在具备 Go 1.22+ 环境的终端中克隆并编译：

```bash
# 克隆仓库
git clone https://github.com/Freecode100Year/usc.git
cd usc

# 编译生成 CLI 命令行工具
go build -o bin/usc.exe ./cmd/usc

# (Linux / macOS)
# go build -o bin/usc ./cmd/usc
```

编译完成后即可通过 `bin/usc.exe`（或将 `bin` 加入 `PATH`）使用命令行。

---

### 2. 核心 CLI 命令详解

USC 提供了 10 个标准化的零信任编译器与运行时管理命令：

```text
Core Commands:
  analyze <source>               分析不可信 Skill 意图、观测真实行为并计算发散度
  extract <source>               无毒提取 Canonical Intent IR，隔离原始提示词
  rebuild <blueprint> [--target] 在物理洁净室中独立重建目标运行时代码
  build <source> [--target]      执行端到端 7 阶段零信任编译并输出机器证明包
  verify <artifact.usc>          验证编译产物完整性与 Ed25519 签名
  verify-proof <proof-dir>       独立核验机器证明包（Proof Bundle）5 大证明义务
  run <artifact.usc>             在受限沙箱与能力代理（Broker）中执行 Skill
  trace -f <skill-id>            实时追踪运行时能力中介与鉴权请求流
  replay <trace.usctrace>        执行确定性决策重放（Deterministic Decision Replay）
  top                            显示系统零信任安全状态与实时监控仪表盘
```

---

#### 🔍 命令 1：分析不可信 Skill (`usc analyze`)
用于审查任意来源（ClawHub、GitHub、ZIP、本地脚本等）的不可信代码。通过 AST 分析其声明意图（$I_{decl}$）与观测到的实际行为（$B_{obs}$），检测越权残留（$B_{extra}$）并计算风险加权发散度分数（$D_N$）。

```bash
usc analyze ./examples/weather-skill
```
**输出示例**：
```text
[+] Analyzing untrusted source: ./examples/weather-skill
------------------------------------------------------------
Status:             PASS
Divergence Score:   0.00
Observed Caps:      1
Extra Behaviors:    0
```

---

#### 🧬 命令 2：提取标准意图 IR (`usc extract`)
将不可信 Skill 中的有效目标转换为平台中立的 `Canonical Intent IR`。彻底过滤可能夹带间接提示词注入（Prompt Injection）的原始文本与注释，实行物理级污染隔离。

```bash
usc extract ./examples/weather-skill
```
**输出示例**：
```text
[+] Extracting Canonical Intent IR for: ./examples/weather-skill
Schema:      usc.intent.v0.1
Intent ID:   intent_weather-skill
Status:      TAINT_QUARANTINED
Canonical IR extracted without raw prompt leakage.
```

---

#### 🏗️ 命令 3：洁净室独立重建 (`usc rebuild`)
读取冻结后的 Minimal Capability Blueprint 与 Canonical Intent IR，在无源码挂载、关闭外网的物理洁净室中为目标平台（如 `hermes`、`openclaw`、`langgraph`）重新生成功能代码。

```bash
usc rebuild blueprint.json --target hermes
```
**输出示例**：
```text
[+] Rebuilding candidate in Clean-Room for target: hermes
Physical Barrier:  ACTIVE (Separate Process, Net Denied)
Raw Source Leaks:  0 bytes
Codegen Status:    GENERATED_UNTRUSTED -> RE_AUDITED (PASS)
```

---

#### ⚙️ 命令 4：端到端零信任编译 (`usc build`)
串联执行完整的 7 阶段流水线：
`INGEST (10%)` $\rightarrow$ `DECONTAMINATE (25%)` $\rightarrow$ `MINIMIZE (15%)` $\rightarrow$ `CLEAN_REBUILD (20%)` $\rightarrow$ `RE_AUDIT (10%)` $\rightarrow$ `SANDBOX (15%)` $\rightarrow$ `ATTEST (5%)`。

```bash
usc build ./examples/weather-skill --target hermes
```
**输出示例**：
```text
[+] Building weather-skill for target hermes...
 [Stage 1/7] INGEST        (10%) ... PASS
 [Stage 2/7] DECONTAMINATE (25%) ... PASS
 [Stage 3/7] MINIMIZE      (15%) ... PASS
 [Stage 4/7] CLEAN_REBUILD (20%) ... PASS
 [Stage 5/7] RE_AUDIT      (10%) ... PASS
 [Stage 6/7] SANDBOX       (15%) ... PASS
 [Stage 7/7] ATTEST        ( 5%) ... PASS

[✓] Build Complete: dist/weather-skill.usc
    Machine Proof Bundle: dist/proof
    Capability Count Reduction (CCR): 71.4%
    Attack Surface Reduction   (ASR): 92.5%
    Artifact Status: ATTESTED
```

---

#### 🔐 命令 5：制品签名验证 (`usc verify`)
独立验证 `.usc` 产物的文件散列、编译器标识与 Ed25519 签名有效性，核对目标企业策略匹配度。

```bash
usc verify dist/weather-skill.usc
```
**输出示例**：
```text
[+] Verifying artifact: dist/weather-skill.usc
Artifact Digest:    sha256:d82e11a94f...
Attestation Status: VALID_ED25519_SIGNATURE
Policy Match:       TARGET_POLICY_APPROVED
Verdict:            PASS
```

---

#### 📜 命令 6：机器证明包全链验证 (`usc verify-proof`)
不信任任何编译报告或仪表盘，直接对 `dist/proof/` 目录中的全部机器证据进行严格独立计算与数学判定，验证 5 大证明义务（$PO_1 \land PO_2 \land PO_3 \land PO_4 \land PO_5$）：

```bash
usc verify-proof dist/proof
```
**输出示例**：
```text
[+] Verifying Machine Proof Bundle in: dist/proof
  [✓] attestation.json:      Valid schema & signature
  [✓] audit-chain.json:      Lossless hash chain verified (H0 -> Hn)
  [✓] cleanroom-proof.json:  0 raw source leaks, net denied verified
  [✓] blueprint.json:        Minimal lattice bound verified
  [✓] sbom.spdx.json:        SPDX 2.3 SBOM consistent

All 5 Proof Obligations satisfied: PO1 ∧ PO2 ∧ PO3 ∧ PO4 ∧ PO5 = true
Final State: ATTESTED
```

---

#### 🛡️ 命令 7：受限沙箱运行 (`usc run`)
在 USC 运行时强制中介平面（Enforcement Plane）中启动 Skill。禁止读取环境变量明文密钥，所有对外交互必须经过中介 Broker。

```bash
usc run dist/weather-skill.usc
```
**输出示例**：
```text
[+] Launching artifact inside guarded USC Runtime: dist/weather-skill.usc
Runtime Mediation Plane: ACTIVE
Secret Capability Broker: credential:// handles mapped
Egress Network Guard:    DESTINATION_CHECK_ENFORCED
Execution Confinement:   SECCOMP_SANDBOX_ACTIVE
[Runtime Output] Hello from zero-trust rebuilt Agent Skill!
```

---

#### 📡 命令 8：实时能力追踪 (`usc trace`)
实时监控正在执行的 Skill 发起的每一项权限请求、中介检查状态与拦截结果：

```bash
usc trace -f weather-skill
```
**输出示例**：
```text
[+] Streaming live runtime capability trace for: weather-skill
17:02:01.104 HANDLE_RESOLVE credential://github/pr_reader
17:02:01.105 HOST_CHECK     api.github.com PASS
17:02:01.105 METHOD_CHECK   GET PASS
17:02:01.106 PATH_CHECK     /repos/foo/bar/pulls/42 PASS
17:02:01.120 TLS_CONNECT    api.github.com:443
17:02:01.240 RESPONSE       200 / 14.2KB
```

---

#### ⏪ 命令 9：确定性决策重放 (`usc replay`)
基于飞行记录器（Flight Recorder）生成的 `.usctrace` 追踪轨迹与当时冻结的策略基线，逐步回放安全裁决，确保重放决策 100% 比特级确定：

```bash
usc replay trace.usctrace --step
```
**输出示例**：
```text
[+] Starting Deterministic Decision Replay for: trace.usctrace
Step 01: Expected=ALLOW Replayed=ALLOW [MATCH]
Step 02: Expected=DENY Replayed=DENY [MATCH]
Deterministic Replay Integrity: true (All steps bit-exact)
```

---

#### 📊 命令 10：实时安全监控仪表盘 (`usc top`)
查看活跃运行的 Skill 状态、洁净室运行情况、无损审计链长度、平均 ASR 与 CCR 削减比率：

```bash
usc top
```
**输出示例**：
```text
================================================================
             USC ZERO-TRUST SECURITY DASHBOARD                  
================================================================
 Active Skills:      1 running / 0 blocked
 Clean-Room Status:  ONLINE (Isolated)
 Audit Chain Length: 142 events (Hash Chain: OK)
 Average ASR:        92.5%
 Average CCR:        71.4%
 Flight Recorder:    RingBuffer active (0 violations)
================================================================
```

---

### 3. 机器证明包目录结构 (Machine Proof Bundle)

每次成功执行 `usc build` 后，系统将在 `dist/` 目录下生成完整的机器可验证证明集合：

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
