# USC (Universal Skill Compiler)

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

# 🚀 编译与完整命令手册 (CLI Full Reference)

### 编译安装
```bash
go build -o bin/usc.exe ./cmd/usc
```

### 指令全集

| 指令 | 作用说明 | 典型场景 |
| :--- | :--- | :--- |
| `usc targets` | 探测已支持的 Agent 环境 | **新手必用**，查看当前电脑具备哪些可用的 Agent 运行时 |
| `usc install <artifact>` | 一键安装到目标 Agent | **新手最爱**，免配置自动部署到本地 Agent 技能库 |
| `usc export <artifact>` | 导出为目标 Agent 原生技能 | 生成对应 Agent 原生目录结构供离线导入 |
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
- 📋 **历史更新日志**：[CHANGELOG.md](CHANGELOG.md)
