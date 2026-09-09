# USC (Universal Skill Compiler) 零信任架构演进与优化建议书

> **文档性质**：独立技术改进与演进建议书（RFC Proposals）  
> **核心原则**：核心代码实现严格遵照《USC Zero-Trust Architecture Specification v0.1》蓝图，不随意发挥主观想象。本建议书将设计中具有长远工程价值、前沿技术升级空间及工程落地细节的具体改进点独立梳理归档，作为后续 v0.2+ 版本迭代和 RFC 讨论的官方参考。

---

## 1. 权限求解器（Capability Minimization Solver）升级建议

### 现状与限制（v0.1）
在 v0.1 规范中，求解器抽象接口 `Solver` 采用启发式/贪心算法（`GreedySolver`）结合硬编码规则约束。贪心策略虽能以 $O(N)$ 复杂度快速计算最小能力集合，但在复杂的跨组件依赖树或存在多候选执行路径（Pareto Frontier）时，容易陷入局部最优解，导致能力缩减未达全局最紧界。

### 建议演进（v0.2+）
1. **引入 Z3 SMT (Satisfiability Modulo Theories) 约束求解器**：
   - 将 Capability Lattice、路径前缀包含关系、动作掩码约束及企业策略（Enterprise Policy Profile）转化为一阶逻辑谓词逻辑公式与 Horn 子句。
   - 利用 SMT-LIB2 标准定义能力约束集合，证明 $Satisfy(I, W, T, C)$ 的可行域解集。
2. **多目标 ILP (Integer Linear Programming) 优化**：
   - 目标函数：$\min \sum w_i R_i$，其中决策变量为原子能力授予布尔指示变量 $x_c \in \{0, 1\}$。
   - 求解网络面、文件系统面、凭据暴露面与系统调用的多目标 Pareto 前沿解，使用户可在“极度安全（Strict）”与“高兼容性（Compatible）”之间选择最优投影。

---

## 2. 操作系统级沙箱与洁净室（Physical Clean-Room）跨平台方案

### 现状与限制（v0.1）
规范第 11-12 节基于 Linux 内核特性（`/proc/self/mountinfo`, `/proc/self/fd`, namespaces, seccomp-bpf, process capabilities）定义了物理洁净室测谎（Clean-Room Proof）。然而在 Windows 宿主环境下，缺乏原生 Linux 命名空间机制。

### 建议演进
1. **Windows 原生沙箱隔离层**：
   - 使用 **Windows AppContainer** 隔离进程执行空间，剥夺 `SeDebugPrivilege` 等危险 Token 特权。
   - 结合 **Windows Job Objects** 限制 CPU、内存、进程生成树（禁止派生子进程 `JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE`）。
   - 文件系统重定向：利用 Windows Filter Driver 或挂载独立的虚拟只读 VHDX / 内存盘。
2. **Linux 环境生产化落地**：
   - 集成 **Landlock LSM**（Linux 5.13+）在无 root 权限下实现非特权规则沙盒限制文件读写。
   - 采用 **eBPF (TC / Sockops / LSM Probes)** 替代传统 iptables 拦截外部流量，实现纳秒级强制中介与 destination-bound 校验。
3. **轻量 WASM/WASI 沙箱作为跨平台运行时标准**：
   - 为 Target Runtime 增加 WASI (WebAssembly System Interface) 编译目标，在用户态实现零系统调用逃逸的内存安全隔离。

---

## 3. 对抗性去污（Decontamination）与 Prompt 注入防御升级

### 现状与限制（v0.1）
在提取 Canonical Intent IR 时，若原不可信 Skill 在 `SKILL.md`、注释或元数据中潜藏复杂的间接提示词注入（Indirect Prompt Injection），仅依赖 AST 静态解析无法完全抵御隐蔽语义攻击。

### 建议演进
1. **双重盲审语义防火墙 (Dual-Blind Semantic Boundary)**：
   - 引入专用语义脱敏隔离器，将不可信自然语言描述通过严格受限的模式语法（Context-Free Grammar / PEG）格式化为非自然语言的抽象数据对象。
2. **结构化语义投影（Structured Schema Projection）**：
   - 严禁 LLM 直接输出自然语言作为下游输入，仅允许输出受 JSON Schema 强类型约束的 AST 节点。
3. **差分投毒检测 (Differential Poisoning Audit)**：
   - 针对源文件运行无上下文静态切片器，对比提取前后的 AST 依赖图。一旦发现文本内容与代码逻辑存在偏离（如文档声明为“天气查询”，但 AST 引用了系统命令），直接触发 $D_N > \tau_{critical}$ 的硬阻断。

---

## 4. 凭据代理（Secret Capability Broker）硬件绑定与零知识证明

### 现状与限制（v0.1）
v0.1 的 Credential Broker 在宿主用户态内通过内存检查拦截凭据直接暴露，将明文注入到 TLS 报文中。但如果宿主进程受到高权限恶意代码内存嗅探，仍有泄露风险。

### 建议演进
1. **TPM 2.0 / 硬件加密机 (HSM) 绑定**：
   - 将 Capability Handle 的私钥与本地硬件 TPM（Trusted Platform Module）或硬件安全区域（Enclave, 如 Intel SGX / AMD SEV）绑定。
   - 签署 Attestation 时使用 TPM PCR 寄存器测量编译器状态，确保证明具有防篡改硬件背书。
2. **TLS 终止与 mTLS 硬件解耦**：
   - Broker 作为独立前置 Sidecar Proxy 运行，Skill 生成的代码仅通过 Unix Domain Socket 发起明文无凭据请求，由 Proxy 完成鉴权注入与双向 TLS，彻底消除不可信代码访问任何密钥载荷的可能性。

---

## 5. 供应链安全标准（SLSA Level 4 / in-toto / Sigstore）深度集成

### 建议演进
1. **兼容 SLSA v1.0 (Supply-chain Levels for Software Artifacts)**：
   - 将 USC Machine Proof Bundle 映射为 SLSA Provenance v1.0 规范，输出标准 `provenance.slsa.json`。
2. **Sigstore Cosign 开放密钥基础设施**：
   - 支持通过 OIDC 无私钥签名（Keyless Signing）将 Attestation 证书记录写入 Rekor 透明度日志系统，供任何第三方 Agent 基础设施通过公共账本实时审计。
3. **SBOM 标准化升级**：
   - 统一输出符合 CycloneDX 1.5 和 SPDX 3.0 的软件物料清单，涵盖所有编译时内嵌依赖与 Polyfill 散列。

---

## 6. 确定性构建与重放性能优化（Performance & Determinism）

### 建议演进
1. **无锁环形缓冲区（Lock-Free Ring Buffer）**：
   - Flight Recorder 在高吞吐运行时环境下，采用原子 CAS 无锁环形缓冲区记录事件，将追踪对 Skill 运行时性能损耗降至 1% 以下。
2. **RFC 8785 JCS 向量化加速**：
   - 在审计账本高频计算中，为 JCS 规范化 JSON 与 SHA-256 采用 SIMD/AVX2 指令集优化，降低哈希链追加延迟。
3. **分布式决策重放检查点 (Decision Snapshots)**：
   - 在长时运行 Skill 追踪中，引入每 1000 步的 State Snapshot 机制，支持任意时刻的快进重放（Fast-forward Replay），避免从 Genesis 开始全程回放。
