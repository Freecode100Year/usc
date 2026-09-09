# USC Zero-Trust Architecture Specification v0.1

**Universal Skill Compiler**

> **Trust the contract, not the code.**
> **Transfer capability, not trust.**
> **Probabilistic understanding, deterministic enforcement.**
> **Enforce what you declare. Observe what you enforce. Prove what you observe.**

---

# 0. Executive Definition

USC 的本质不是 Agent Skill 的代码转换器，而是：

> **针对不可信意图、未知实现和异构 Agent Runtime 的零信任洁净室重建编译器。**

英文定义：

> **USC is a zero-trust clean-room compiler that decomposes untrusted Agent Skills into intent, behavior and authority, compresses them into a minimal reviewable trust contract, independently rebuilds the functionality for heterogeneous Agent runtimes, and emits machine-verifiable attestations for both compilation and execution.**

中文定义：

> **USC 是一个零信任洁净室编译器：它将不可信 Agent Skill 分解为意图、行为与权限，将其压缩为最小可审查信任契约，再针对不同 Agent Runtime 独立重建，并对编译过程、生成产物和运行时行为提供机器可验证的证明。**

传统转换器解决：

```text
Source A
   │
Convert
   │
Target B
```

USC 解决：

```text
Untrusted Implementation
        │
Intent Extraction
        │
Intent Decontamination
        │
Necessary Behavior
        │
Capability Minimization
        │
Minimal Trust Contract
────────────────────────
   TRUST BARRIER 
────────────────────────
Independent Rebuild
        │
Independent Re-Audit
        │
Runtime Validation
        │
Cryptographic Attestation
```

因此：

> **跨平台只是 USC 编译流水线的降阶产物。**

USC 真正的核心技术护城河是：

1. Zero-Trust Invariants
2. Intent Decontamination
3. Capability Lattice
4. Monotonic Privilege Reduction
5. Capability Minimization
6. Secret Capability Broker
7. Physical Clean-Room Barrier
8. Independent Re-Audit
9. Runtime Mediation
10. Machine-Verifiable Proof Chain
11. White-Box Observability
12. Deterministic Decision Replay

---

# 1. Core Creed

USC 遵循四条根本信条：

> **Trust no Skill.**
> 不信任任何 Skill。

> **Grant the minimum.**
> 只授予最小权限。

> **Rebuild from intent.**
> 从意图重新构建。

> **Transfer capability, not trust.**
> 传递能力，不传递信任。

进一步形成工程总原则：

> **Probabilistic systems may propose. Deterministic mathematics decides.**

即：

> **概率系统负责提出候选，确定性数学负责裁决。**

LLM、启发式规则、静态模型可以帮助理解：

```text
What does this Skill probably want to do?
```

但最终安全裁决必须由：

```text
Capability Lattice
Set Algebra
Policy
Proof Obligations
Invariant Verification
Deterministic Guards
```

完成。

---

# 2. Trust State Machine

USC 中不存在从源码可信直接跳到产物安全的概念。

所有状态严格单向演进：

```text
[ UNTRUSTED ]
      │
      │ Isolated Parsing
      │ Taint Tagging
      ▼
[ DERIVED ]
      │
      │ Intent Decontamination
      │ Necessity Analysis
      ▼
[ CONSTRAINED ]
      │
      │ Minimal Capability Blueprint Frozen
      │
══════╪════════════════════════════════════
      │ PHYSICAL TRUST BARRIER
══════╪════════════════════════════════════
      │
      │ Clean-Room Codegen
      ▼
[ GENERATED_UNTRUSTED ]
      │
      │ Independent Re-Audit
      ▼
[ RE-AUDITED ]
      │
      │ Sandbox / Guard Validation
      ▼
[ VALIDATED ]
      │
      │ Proof Obligation Closure
      │ Cryptographic Signing
      ▼
[ ATTESTED ]
```

状态定义：

## UNTRUSTED

包括：

```text
ClawHub Skill
GitHub Repository
ZIP
SKILL.md
Prompt
Shell Script
Python Script
Dependency Manifest
Binary
```

所有源输入统一标记：

```text
TAINTED_UNTRUSTED
```

来源不会提升信任。

---

## DERIVED

从不可信实现提取出的事实：

```text
Declared Intent
Observed Behavior
AST
Call Graph
Data Flow
Dependency Facts
Candidate Capability Graph
```

这些依然可能受源内容投毒影响。

DERIVED 不是授权状态。

---

## CONSTRAINED

通过：

```text
Intent Decontamination
Functional Necessity Analysis
Capability Minimization
Policy Resolution
Secret Flow Constraints
```

生成：

```text
Canonical Intent IR
Minimal Capability Blueprint
Approved Policy Contract
```

并被冻结为构建基线。

---

## GENERATED_UNTRUSTED

Clean-Room Generator 根据 CONSTRAINED 输入重新生成代码。

虽然不继承原始代码，生成物仍可能存在：

```text
LLM hallucination
Codegen bug
Adapter bug
Privilege widening
Prompt injection reintroduction
Incorrect dependency
```

所以必须标记：

```text
GENERATED_UNTRUSTED
```

---

## RE-AUDITED

由与 Codegen 独立的安全 Pass 对生成物重新进行：

```text
AST analysis
Capability extraction
Secret flow analysis
Prompt analysis
Policy comparison
INV-7 assertion
Dependency analysis
```

---

## VALIDATED

产物经过受限 Sandbox：

```text
Smoke Validation
Policy Denial Tests
Mock External Services
Expected Capability Tests
Exit Code Checks
```

v0.1 中：

```text
SMOKE_VALIDATED
```

不应被表述为完整业务正确性证明。

---

## ATTESTED

仅代表：

```text
Artifact identity verified
Compiler identity verified
Blueprint verified
Policy verified
Proof obligations satisfied
Audit ledger integrity verified
SBOM generated
Signature valid
```

ATTESTED 不代表绝对安全。

ATTESTED 也不必然等于 DEPLOYABLE。

部署资格应单独表示：

```text
Artifact State:
ATTESTED

Deployment Eligibility:
TARGET_POLICY_APPROVED
```

---

# 3. The Nine Zero-Trust Invariants

## INV-1 — Origin Neutrality

> 外部代码绝不因来源而获信。

形式：

$$Trust(Source)=UNTRUSTED$$

无论来源：

```text
ClawHub
Official GitHub
Vendor example
Community repo
Local ZIP
```

均执行：

```text
TAINTED_UNTRUSTED
```

---

# INV-2 — Zero Grant Baseline

> 索取不代表赋予。

初始能力授予：

$$C_0=\varnothing$$

原 Skill 声明：

```text
net.egress:*
fs.read:*
shell.exec:*
```

只能作为：

```text
Intent clue
Observed declaration
Analysis evidence
```

不能成为授权依据。

---

# INV-3 — Authority Severance

> 意图可以继承，权限绝不继承。

禁止：

```text
Original Permission ──> Target Permission
```

允许：

```text
Original Source ──> Intent Analysis ──> Necessary Operation ──> Minimal Capability
```

因此必须切断：

$$Authority_{source} \nRightarrow Authority_{target}$$

---

# INV-4 — Candidate Distrust

> 洁净室生成的候选实现，在独立重审之前仍然不可信。

$$Generated \neq Trusted$$

生成状态必须：

```text
GENERATED_UNTRUSTED
```

---

# INV-5 — Secrets as Capabilities

> Secret 是能力，不是环境变量。

禁止：

```text
process.env.GITHUB_TOKEN
os.getenv("TOKEN")
export SECRET=...
```

生成代码只能获得：

```text
credential://github/pr_reader
```

凭据流：

```text
Secret Vault
    │
Capability Handle
    │
USC Credential Broker
    │
Authorized Transform
    │
Authorized Network Primitive
    │
Authorized Destination
```

必须满足：

$$Secret \rightarrow Reader \rightarrow Transform \rightarrow Primitive \rightarrow Destination$$

并且：

> **Secret material MUST NOT enter the untrusted execution address space.**

---

# INV-6 — Strict Provenance

> 每一次转换都必须有因果证据。

任何转换节点必须回答：

```text
Where did this fact come from?
Why was this transformation performed?
Which rule allowed it?
Which intent step requires it?
What confidence level supports it?
```

形式：

$$Transformation \Rightarrow Provenance\neq\varnothing$$

对于低置信度行为：

```text
Low confidence + privilege expansion ──> HARD_BLOCK
Low confidence + uncertain functionality ──> REVIEW / DEGRADE
Low confidence + safe privilege contraction ──> ALLOW with provenance
```

---

# INV-7 — Monotonic Privilege Reduction

> 自动变换过程中，权限只能保持不变或收缩。

定义偏序：

$$C_a\sqsubseteq C_b$$

表示：

> $C_a$ 的授权范围不超过 $C_b$。

每个自动 Pass：

$$C_{i+1}\sqsubseteq C_i$$

生成代码必须：

$$C_{gen}\sqsubseteq C_{blueprint}$$

运行时必须：

$$C_{runtime}\sqsubseteq C_{generated}$$

因此：

$$C_{runtime}\sqsubseteq C_{blueprint}$$

任何违反：

$$C_{out}\not\sqsubseteq C_{in}$$

立即：

```text
HARD_BLOCK
```

---

# INV-8 — Deterministic Confinement

> 运行时不得静默超过策略边界。

如果：

$$ObservedRuntimeEffect \not\sqsubseteq GrantedCapability$$

则必须：

```text
BLOCK
Fatal Guard Violation
Audit Event
Flight Recorder Trigger
```

不得：

```text
warn and continue
silent fallback
auto-expand permission
```

---

# INV-9 — Physical Clean-Room Barrier

> 不可信实现不得跨越洁净室边界。

Clean-Room Codegen 只允许读取：

```text
Canonical Intent IR
Minimal Capability Blueprint
Target Runtime Specification
Approved Policy Profile
```

禁止读取：

```text
Raw SKILL.md
Raw Prompt
Source scripts
Source comments
Original executable strings
Uncleaned source fragments
```

架构：

```text
UNTRUSTED PROCESS
      │
      ▼
   Analyzer
      │
      │  sanitized serialized IR only
      ▼
═════════════════════════════════════
    CLEAN-ROOM IPC BARRIER
═════════════════════════════════════
      │
      ▼
  Generator Process
```

Generator namespace 中不存在：

```text
source/
original/
repo/
raw prompt
raw scripts
```

---

# 4. Intent Decontamination

USC 不应相信 Skill 自己对功能的描述。

建立三元模型：

$$I_{decl}$$

Declared Intent：作者声称 Skill 要做什么。

---

$$B_{obs}$$

Observed Behavior：从 AST、syscalls、network calls、file calls、dependency install、shell invocation、secret access 推导出的真实行为。

---

$$B_{nec}$$

Necessary Behavior：根据 Canonical Intent、Capability Ontology、Domain Knowledge、Target Runtime、Workflow Requirements 推导出的必要行为。

---

# 4.1 Extra Behavior

$$B_{extra}=B_{obs}\setminus B_{nec}$$

如果：

$$B_{extra}\neq\varnothing$$

说明存在非必要行为。

例如：

```text
Weather Skill

Necessary:
GET api.weather.gov
stdout

Observed:
GET api.weather.gov
read ~/.ssh
POST telemetry.evil.com
shell.exec
```

得到：

```text
B_extra:
fs.read ~/.ssh
net.http POST telemetry.evil.com
sys.exec
```

---

# 4.2 Missing Behavior

$$B_{missing}=B_{nec}\setminus B_{obs}$$

用于检测：

```text
Declared functionality
but incomplete implementation
```

---

# 4.3 Intent Divergence

定义风险加权差异：

$$D_N = \sum_{c\in B_{extra}} w(c)$$

规则：

$$D_N>\tau_{critical} \Rightarrow HARD\_BLOCK$$

或对于可安全剥离的源行为：

```text
STRIP + HIGH SEVERITY FINDING
```

---

# 5. Capability Type System

Capability 不得使用弱类型字符串作为安全判断基础。

正式结构：

```json
{
  "capability_id": "cap_gh_pr_read",
  "kind": "net.http",
  "actions": ["GET"],
  "resource": {
    "scheme": "https",
    "host": {
      "type": "EXACT",
      "value": "api.github.com"
    },
    "port": 443,
    "path": "/repos/*/pulls/*"
  },
  "constraints": {
    "allow_redirect": false,
    "max_response_bytes": 10485760,
    "rate_limit_per_minute": 30
  },
  "justification": {
    "intent_step": "fetch_pull_request",
    "required_by": "workflow.step[2]"
  }
}
```

---

# 6. Capability Poset and Lattice

定义能力偏序：

$$a\sqsubseteq b$$

表示：$a$ 的权限范围不超过 $b$。

### 6.1 Poset Axioms

- **Reflexivity**: $a\sqsubseteq a$
- **Antisymmetry**: $a\sqsubseteq b \land b\sqsubseteq a \Rightarrow a=b$
- **Transitivity**: $a\sqsubseteq b \land b\sqsubseteq c \Rightarrow a\sqsubseteq c$

### 6.2 Meet and Join

$$C_{effective} = C_{blueprint} \sqcap C_{enterprise} \sqcap C_{runtime}$$

---

# 7. Authority Confinement Theorem

$$\boxed{C_{runtime} \sqsubseteq C_{blueprint}}$$

---

# 8. Capability Minimization

$$C_{min} = \arg\min_{C\in F} Cost(C)$$

$$Cost(C) = (R_n, R_f, R_s, R_w, R_e, R_p, R_{persist})$$

---

# 9. Information Flow Security

$$Public \sqsubseteq Internal \sqsubseteq Sensitive \sqsubseteq Secret$$

---

# 10. Secret Capability Broker

运行时不得把 Secret 交给 Skill。

---

# 11. Clean-Room Physical Isolation

Generator 必须：独立进程、无网络、无源码挂载。

---

# 12. Clean-Room Proof

实际测量：mountinfo, fd, network namespace, seccomp, process capabilities, environment, digest scan。

---

# 13. Compiler Pipeline

```text
UNTRUSTED INPUT
      │
1. INGEST
      │
2. DECONTAMINATE
      │
3. MINIMIZE
      │
CONSTRAINED BLUEPRINT
── TRUST BARRIER ──
4. CLEAN REBUILD
      │
GENERATED_UNTRUSTED
      │
5. RE-AUDIT
      │
6. SANDBOX
      │
7. ATTEST
```

---

# 14. Seven Observable Compiler Stages

INGEST (10%), DECONTAMINATE (25%), MINIMIZE (15%), CLEAN_REBUILD (20%), RE_AUDIT (10%), SANDBOX (15%), ATTEST (5%).

---

# 15. Proof Obligations

$$\boxed{ATTESTED \iff \bigwedge_i ProofObligation_i}$$

---

# 16. Diagnostics

PANIC, HARD_BLOCK, WARN, DROP, PASS。

---

# 17. Lossless Audit Ledger & 18. Audit Hash Chain

$$H_n = SHA256(H_{n-1} \parallel Canonical(Event_n))$$
Canonical JSON 采用 RFC 8785 JCS。

---

# 19. Observable Proof Principle & 20. Trust Compression

$$CCR = 1 - \frac{|C_{granted}|}{|C_{observed}|}, \quad ASR = 1 - \frac{Cost(C_{granted})}{Cost(C_{observed})}$$

---

# 21. White-Box Observability & 22. Runtime Enforcement Plane

Mediation Before Observation, Evidence Before Visualization, Metadata Before Content, Replay Decisions Before Replaying Reality.

---

# 23. Runtime Capability Trace & 24. IPC Transcript & 25. Flight Recorder & 26. Decision Replay

支持 `.usctrace` 环形缓冲记录与 `usc replay trace.usctrace --step` 决策重放。

---

# 27. Machine Proof Bundle & 28. Attestation

包含 `attestation.json`, `audit-chain.json`, `provenance.json`, `cleanroom-proof.json`, `policy.json`, `blueprint.json`, `sbom.spdx.json`。

---

# 29-31. Foundations & Verification

数学五层骨架（偏序格、优化、程序逻辑、信息流、追踪语义）。

---

# 32. Core Go Architecture

核心 Go 架构映射至 `usc-core/` 与 `cmd/usc/`。
