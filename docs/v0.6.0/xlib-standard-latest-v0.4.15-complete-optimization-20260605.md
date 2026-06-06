# xlib-standard 最新版本 v0.4.15 深度分析与完整优化方案

> 生成日期：2026-06-05
> 分析基线：`ZoneCNH/xlib-standard` 最新发布 `v0.4.15`，commit `c0fc3813e156cf35a37ddd0033432a78943bb32b`
> 分析方式：基于 GitHub 当前仓库文件、最新 release/tag、README、Makefile、goalcli、Harness、Debt、Workflow、Rules、Evidence 文件进行事实核验。未在本地执行完整 CI，因此本文对“可运行通过”的判断只引用仓库已有命令声明与文件事实，不伪造本地运行结果。

---

## 0. 一句话结论

`xlib-standard` 已经不再是普通 Go 模板仓库，而是一个正在成型的 **标准生产系统**：它同时承担 Standard Source、Go Reference Template、Generator、Harness、Evidence Runtime、Goal Runtime、Debt Control Plane、Rules Registry 与 Downstream Adoption 的职责。

但最新 `v0.4.15` 仍没有彻底达到“自我进化、禁止漂移、证据事实、Harness 自证、Goal 化交付、全债务闭环”的目标。主要原因不是缺少文档或命令，而是 **事实源仍分散、验证深度不一致、债务检测仍偏浅、Evidence 仍偏 manifest 化、Goal/Rule/Traceability/Release 之间仍有语义漂移**。

最终建议：

```text
不要继续堆脚本、文档和 checklist。
应该把当前仓库升级为 Standard Production Kernel：

Canonical Facts
  → Standard Graph
  → Goal Graph
  → Debt Graph
  → Harness Proof Graph
  → Evidence Ledger
  → Release Policy Decision
  → Downstream Conformance Proof
  → Retrospective Rule Patch
  → 下一轮自我进化
```

---

## 1. 最新版本事实核验

### 1.1 最新版本

当前 GitHub Release 最新版本是：

```text
version: v0.4.15
commit:  c0fc3813e156cf35a37ddd0033432a78943bb32b
release: 2026-06-05 10:24 UTC 左右
主题: Preserve leader-owned cleanup rule
```

`v0.4.15` 相比 `v0.4.13` 只领先 2 个 commit，但变更面覆盖治理、adoption、branch cleanup、rules deep analysis、render tests、Makefile、Harness 与文档。

### 1.2 仓库声明的五类职责

README 当前明确声明仓库承担五类职责：

```text
Standard Source
Go Reference Template
Generator
Harness
Evidence Runtime
```

它还声明：标准文本、模板、generator、Harness gate 和 Evidence runtime 必须一起通过 release gate 验证。

这说明仓库真实身份已经不是“模板库”，而是 **标准生产与验证运行时**。

### 1.3 最新增量的价值

`v0.4.15` 最新增量主要增强了三件事：

1. **adoption-check**：新增 downstream governance pack / ruleset / Makefile / Harness / registry 的采纳检查。
2. **branch-governance**：把破坏性分支操作明确为 leader-owned，worker 只能审计、分类、准备 backup evidence。
3. **rules deep analysis**：新增 `.agent/rules/` 深度分析，承认规则系统中存在 exit code 冲突、generated artifact 元信息漂移、机器渲染文件过大、规则索引重复、已知 P0 gap 等问题。

这些是正向改进，但它们也暴露了更深的问题：系统正在快速扩张，若没有统一事实内核，很快会从“强治理”滑向“治理系统自身失控”。

---

## 2. 当前成熟度评分

> 这是基于仓库文件事实的结构性评分，不等于 `goalcli score` 的发布门禁结果。

| 维度 | 当前估计 | 说明 |
|---|---:|---|
| 标准身份清晰度 | 9.0 | README 与标准目录已经明确五类职责和非目标。 |
| Harness 覆盖面 | 8.8 | Makefile 与 `.agent/harness/harness.yaml` 覆盖大量 required/final/extended gate。 |
| Evidence 链路 | 8.2 | manifest、checksum、workflow artifact、debt evidence 已成型，但还不是事件账本/证明图。 |
| Goal Runtime | 8.0 | Goal v3.1 与 G12-G16 已接入，但仍有 dry-run/planned/implemented 语义混用风险。 |
| 禁止漂移能力 | 6.6 | 最新版本已经出现 `v0.4.15` 与 `v0.4.13` 常量/命令示例漂移。 |
| 债务治理深度 | 6.5 | 债务 gate 存在，但结构债、实现债、测试债、领域债检测仍明显浅。 |
| Traceability 完整性 | 6.8 | `traceability-check` 已有实现，但规则文档仍称其为 GAP，且实现只做路径存在性级别。 |
| Downstream Proof | 6.2 | adoption-check 已接入，但 downstream adoption 仍主要是 gap/status 证明，不是 D7 adoption proof。 |
| 自我进化闭环 | 6.8 | retrospective/self-improving 有入口，但还未形成 escape → invariant → fixture → detector → policy 的强闭环。 |

综合判断：

```text
当前系统成熟度：约 8.3 / 10
距离 9.8 release-grade 标准生产系统的主要差距：
1. 单一事实源缺失。
2. Gate proof depth 未类型化。
3. Debt detector 语义深度不足。
4. Evidence 仍偏 artifact manifest，不是 append-only proof ledger。
5. Traceability / Rules / Implementation 状态存在事实冲突。
6. Downstream adoption 尚未达到 proof-based conformance。
```

---

## 3. 底层本质

### 3.1 表层问题

你提出的问题包括：

- 结构债：分层违规 import、L2 互相耦合、循环依赖、上帝模块。
- 实现债：重复代码、过时模式、补丁热点。
- 测试债：缺失测试、脆弱测试、金字塔倒置。
- 文档债：ADR 缺失、文档与代码不一致。
- 依赖债：过期、废弃、有 CVE 的第三方库。
- 领域债：模型与业务语言不一致。
- 标准系统债：Harness、自我验证、证明、Evidence、事实、禁止漂移、Goal 化、自我进化。

### 3.2 底层问题

这些不是六七类互不相关的债务，而是同一个底层问题：

```text
仓库缺少一个能把“事实、规则、目标、代码、测试、文档、依赖、证据、发布、下游采用”统一建模的生产内核。
```

只要这些对象分散在 README、Makefile、Go 常量、shell 脚本、workflow YAML、`.agent`、docs、release manifest、registry 中，系统就一定会继续漂移。

### 3.3 最不可再拆解的真理

1. **事实只能有一个权威源。** 其他地方只能是投影或生成产物。
2. **通过不是事实，证据才是事实。** `passed` 必须绑定命令、输入、输出、exit code、环境、digest、时间、runner、policy version。
3. **存在不等于实现。** file exists、command registered、dry-run ready、implemented、executed、evidence verified、release usable 是不同状态。
4. **检查器必须自证。** 一个 gate 没有 negative fixture，就不能证明它真的会在坏情况失败。
5. **债务不是文本标签，而是 broken invariant。** 每个 debt finding 必须绑定 invariant、location、evidence、severity、owner、expiry、remediation、goal。
6. **标准系统必须能生产自身标准。** 否则每次人工修文档、脚本、workflow 都会制造新漂移。
7. **自我进化不能靠总结。** 每次逃逸问题都必须升级为 rule、detector、fixture、policy 或 prompt patch。

---

## 4. 最新版本发现的关键结构性问题

## 4.1 P0/P1：最新版本事实漂移

这是当前最新版本最值得优先处理的问题。

仓库最新 release 是 `v0.4.15`，但多个当前权威/半权威位置仍保留 `v0.4.13`：

```text
cmd/goalcli/governance.go
  projectReleaseVersion = "v0.4.13"

internal/tools/releasemanifest/main.go
  defaultReleaseVersion = "v0.4.13"

README.md
  release-preflight VERSION=v0.4.13

.agent/harness/harness.yaml
  release_preflight command 使用 VERSION=v0.4.13

docs/project-structural-analysis-20260605.md
  仍以 v0.4.13 对齐结果为主叙事
```

### 影响

这直接违反“禁止漂移”和“事实唯一化”。

如果 `v0.4.15` 是真实 release，而 release manifest / goalcli version / harness preflight 仍指向 `v0.4.13`，那么系统会出现三类风险：

1. release evidence 可能生成旧版本字段。
2. 文档示例可能指导用户验证旧版本。
3. Agent/CI/审计读取不同文件时得到不同当前版本。

### 修复原则

不要手工把所有 `v0.4.13` 改成 `v0.4.15` 后结束。那只是下一轮漂移的开始。

正确修法是新增单一事实源：

```yaml
# .xlib/facts/xlib.yaml
schema_version: xlib-facts/v1
module: github.com/ZoneCNH/xlib-standard
current_release:
  version: v0.4.15
  commit: c0fc3813e156cf35a37ddd0033432a78943bb32b
  released_at: 2026-06-05T10:24:00Z
runtime:
  goal_runtime_version: v3.1
  governance_runtime_version: v2.9.3
tools:
  go: "1.23.0"
  golangci_lint: "v2.1.6"
  govulncheck: "v1.1.4"
```

然后由它投影生成：

```text
cmd/goalcli/version_gen.go
internal/tools/releasemanifest/version_gen.go
.agent/harness/harness.generated.yaml fragment
README release command snippet
docs/release.md release command snippet
.github/workflows env/tool versions
```

新增 gate：

```bash
GOWORK=off go run ./cmd/goalcli fact audit --strict
```

release gate 必须阻断任何当前事实漂移。

---

## 4.2 P1：GoalCLI 仍是治理上帝入口

`cmd/goalcli/main.go` 通过一个巨大 `switch` 分发大量命令：version、doctor、guards、context、schema、debt、adoption、goal-runtime、planned commands、external scripts、evidence、release、security、traceability、rules 等。

虽然仓库已经有 `debt.go`、`traceability.go`、`adoption_check.go`、`goalruntime.go` 等拆分文件，但 `main.go` 仍承担：

```text
command registry
policy routing
external script orchestration
security dispatch
score dispatch
planned command dispatch
usage authority
```

### 结构风险

1. 新命令继续向一个 switch 聚集。
2. 命令注册、usage、Makefile、registry、implementation-status 仍需多处同步。
3. planned/dry-run/implemented/real-executed 命令在同一个入口中暴露，用户容易把“可调用”误认为“可证明”。
4. CLI 既是控制面，又是若干 domain checker 的执行面，职责边界不够硬。

### 优化目标

把 `goalcli` 从“上帝模块”改成 thin shell：

```text
cmd/goalcli
  main.go                    # 只负责 Cobra/flag/dispatch shell
internal/controlplane
  commandregistry            # 命令注册与状态
  policy                     # release/debt/goal policy decision
internal/factkernel
internal/debtkernel
internal/tracekernel
internal/harnesskernel
internal/evidencekernel
internal/downstreamkernel
```

命令由 registry 注册，不再靠 main.go 手工 switch 与 usage 文本同步。

---

## 4.3 P1：Traceability 事实冲突

当前仓库存在一个典型“事实源冲突”：

- `.agent/rules/iron-rules.md` 仍声明 `traceability-check` 是 GAP，尚未实现。
- `cmd/goalcli/traceability.go` 已经实现了 `traceability-check`。
- `.agent/registries/command-implementation-status.yaml` 把 `traceability-check` 放入 implemented / release usable 的命令组。
- `docs/reports/rules-deep-analysis-20260605.md` 仍把 Traceability Gate 未落地视为最高风险债务之一。

### 真实判断

`traceability-check` 已经有实现，但它当前实现的证明深度有限：

```text
已做：
- 解析 markdown traceability matrix。
- 找 REQ 行。
- 检查主要产物列里的路径是否存在。
- 有 gap 时返回 9。

未充分做：
- 不验证 Goal → Req → AC → Task → Issue → Commit → PR → Evidence → Release 全链路图。
- 不验证 Evidence digest。
- 不验证 issue/PR/commit 与 REQ 的双向关系。
- 不验证每个 P0 规则是否有 active detector。
- 不验证 evidence 事件是否可 replay。
```

因此它不是“未实现”，也不是“完整实现”。它应标记为：

```text
implementation_status: partial_implemented
proof_depth: D2/D3
release_usable_for: path existence traceability only
not_release_usable_for: full lifecycle proof graph
```

### 修复原则

给所有 gate 加 proof depth，而不是用 implemented / not implemented 二元状态。

---

## 4.4 P1：Debt Control Plane 仍偏浅

当前债务系统是非常好的入口，但还不是完整债务治理内核。

### 当前能力

`Makefile` 已经有：

```text
debt
architecture
domain
docs-drift
dependency-debt
security-debt
testing-debt
implementation-debt
downstream-debt
debt-evidence
debt-trend
debt-patch-suggest
debt-lifecycle-check
```

`internal/debtcheck` 已经输出 debt report、score、digest、sections、findings。

### 核心不足

| 债务类型 | 当前检测方式 | 问题 |
|---|---|---|
| 结构债 | 只扫描 legacy `x.go` import | 不检测 layer graph、L2↔L2、循环依赖、internal/public 边界、god module。 |
| 实现债 | marker 文本 `xlib-implementation-debt` | 不检测重复代码、复杂度、补丁热点、过时模式。 |
| 测试债 | marker 文本 `xlib-testing-debt` | 不检测覆盖率、风险覆盖、脆弱测试、mutation、flaky、金字塔倒置。 |
| 文档债 | marker + docs-check grep | 不检测 ADR 缺失、API 文档漂移、版本事实漂移、历史快照混读。 |
| 依赖债 | `@latest` 和 `curl | bash` | 不检测 CVE evidence、新依赖 owner、license、deprecated、SBOM、runtime purpose expiry。 |
| 领域债 | marker 文本 `xlib-domain-forbidden` | 不检测 ubiquitous language、bounded context、旧术语、DDD 模型偏移。 |
| Downstream 债 | required token / status file | 更像 registry completeness，不是实际 downstream conformance proof。 |

### 更深的问题

当前 `Finding` schema 只有：

```go
ID
Severity
Path
Message
```

缺少治理必要字段：

```text
invariant_id
proof_depth
detector_id
detector_version
confidence
owner
goal_id
adr
remediation
expiry
waiver
first_seen
last_seen
hash
```

此外，当前 scoring 允许 P1/P2 findings 在 enforce 模式下仍返回 passed，只要 score 不低于阈值。这对“观察型债务”可以接受，但对 release-blocking P1 不够精确。应该由 rule registry 的 `release_blocking` 和 policy 决定，而不是只由 P0/P1/P2 计分间接决定。

---

## 4.5 P1：docs-check 仍是文本 needle 检查

`scripts/check_docs.sh` 的主要机制是：

```bash
require_text file needle
```

它检查大量文件是否包含固定字符串，例如 Makefile 片段、workflow 字符串、文档入口、工具版本、目标库名、Docker 文本等。

### 价值

它能防止明显遗漏，是很好的早期防线。

### 问题

它不能理解语义：

1. 字符串存在不代表结构正确。
2. 旧版本字符串可能同时存在，导致新旧事实并存。
3. 文档变得越来越像“为了过 grep 而写”。
4. 变更成本高，误报/漏报都会增长。
5. 它无法判断当前事实与历史快照的上下文差异。

### 正确方向

保留 `docs-check` 作为 D1/D2 低成本防线，但新增结构化检查：

```bash
goalcli docs graph-check
```

它应解析：

```text
Markdown frontmatter
ADR registry
fact projection markers
link graph
release version facts
current/historical snapshot metadata
command references
schema references
```

---

## 4.6 P1：Harness 覆盖很广，但缺少 proof depth 类型系统

`.agent/harness/harness.yaml` 当前列出大量 required/final/extended gates，也显式记录 duplicate command aliases。但 gate 的语义深度不清：

```text
file exists
dry-run verify
external script executes
unit test executes
release evidence verifies
traceability graph verifies
downstream adoption proves
```

这些都被放进 Harness 里，但它们证明的东西完全不同。

### 解决方式：Proof Depth

定义 D0-D7：

| Depth | 含义 | 例子 | 可否 release blocking |
|---|---|---|---|
| D0 | presence | 文件存在 | 否，只能辅助 |
| D1 | text needle | grep 文本 | 否，只能辅助 |
| D2 | syntax/schema | YAML/JSON/Markdown 结构校验 | 可辅助 |
| D3 | static graph | import graph、trace graph、doc graph | 可阻断结构类问题 |
| D4 | execution | 命令真实执行并返回 exit code | 可 release blocking |
| D5 | negative fixture / mutation | 证明坏样本会失败 | 标准系统关键 gate 必须达到 |
| D6 | replay / attestation | 证据可重放、digest 可验证 | release final 必须达到 |
| D7 | downstream conformance | 下游真实采用并通过 | adoption claim 必须达到 |

release-blocking gate 最低 D4；标准/Harness/Evidence/Traceability gate 应逐步达到 D5/D6；downstream adoption 只能由 D7 声明。

---

## 4.7 P1：Evidence 仍偏 manifest，不是完整 Evidence Ledger

当前 manifest 已记录很多字段，也有 checksum 和 workflow artifact。这是很好的基础。

但完整标准生产系统需要的是：

```text
Evidence Event Ledger
  + Proof Graph
  + Manifest Projection
```

而不是把 `release/manifest/latest.json` 当成唯一中心。

### 当前风险

1. `latest.json` 是生成产物，语义上容易被误认为当前事实。
2. 单个 manifest 不表达每个 gate 的输入/输出因果图。
3. checksum 证明完整性，不证明语义有效性。
4. CI workflow 上传 debt evidence，release workflow 主要上传 manifest，artifact 集合存在不对称。
5. `CHECK_STATUS=passed` 这类变量仍需要更强 policy 绑定，避免“传参即通过”的错觉。

### 正确结构

```text
release/evidence/events/*.jsonl       # append-only event ledger
release/proof/latest.proof.json       # proof graph projection
release/manifest/latest.json          # release summary projection
release/attestations/*.json           # optional attestation projection
```

manifest 应该是 ledger 的投影，而不是唯一事实源。

---

## 4.8 P1：Security 证据 freshness 仍需强制可观测

当前安全策略已经更合理：CI/Release 默认不跑外部漏洞库，Security workflow 每周强制 `govulncheck`。

但新风险是：

```text
默认 gate 通过 ≠ 漏洞扫描证据新鲜。
```

如果 Security workflow 长期失败、被禁用、权限异常或 `vars.XLIB_ENABLE_VULNCHECK` 为空，日常 release-check 仍可能通过。

### 建议

新增：

```bash
goalcli security freshness-check --max-age 168h --require-scheduled-success
```

并在 release-final-check 中要求：

```text
security.last_successful_vulncheck_at <= 168h
security.workflow_conclusion == success
security.govulncheck_version == canonical facts version
```

---

## 4.9 P1：Downstream adoption 仍不是 proof-based conformance

`adoption-check` 是 `v0.4.15` 的重要进步，它会检查 downstream governance pack、workflow、ruleset、harness、registry、Makefile 等。

但这仍主要是：

```text
结构存在性 + 文本契约 + ruleset 参数
```

不是：

```text
真实 downstream repo 渲染 → 构建 → 测试 → debt evidence → release evidence → adoption proof
```

当前 `.agent/registries/command-implementation-status.yaml` 也诚实地区分了 downstream-baseline/downstream-adoption 仍是 gap report only，不得声明 proof-based adoption。

### 需要补齐

```bash
goalcli downstream conformance --repo ../kernel --profile release --write-proof
goalcli downstream conformance --repo ../configx --profile release --write-proof
goalcli downstream conformance --matrix .xlib/downstream/matrix.yaml
```

只有 D7 proof 才能写：

```yaml
adoption_status: adopted
proof_based_adoption: true
```

---

## 5. 被误认为真理的常见假设

| 常见假设 | 为什么不可靠 | 应替换为 |
|---|---|---|
| 有 README 就有标准 | README 是投影，不是事实源 | Canonical facts + generated docs |
| 有 Makefile target 就有 gate | target 可能只是 dry-run 或外部脚本 | command status + proof depth |
| 有 CI passed 就可以 release | CI 可能没有 fresh security/downstream/proof evidence | release policy decision |
| 有 checksum 就有证据 | checksum 只证明文件没变 | subject digest + command + input/output + policy |
| 有 debt score 就无债 | 检测器浅时 score 会虚高 | invariant coverage + detector depth |
| 有 traceability-check 就有全链路追踪 | 当前只证明路径存在性层面 | lifecycle proof graph |
| adoption-check 通过就已采用 | 只证明 pack/规则存在 | downstream D7 conformance proof |
| 自我进化是写 retrospective | retrospective 不等于系统进化 | escape → rule → fixture → detector → policy |

---

## 6. 从零设计的新方案：Standard Production Kernel

### 6.1 总体架构

```text
┌────────────────────────────────────────────────────────────────┐
│                   Standard Production Kernel                    │
├────────────────────────────────────────────────────────────────┤
│ 1. Canonical Fact Kernel                                        │
│    version/tools/gates/downstream/rules/status/policies          │
├────────────────────────────────────────────────────────────────┤
│ 2. Standard Graph                                                │
│    standards/specs/contracts/public API/modules/layers            │
├────────────────────────────────────────────────────────────────┤
│ 3. Goal Graph                                                    │
│    Goal → Req → AC → Task → Issue → Commit → PR                  │
├────────────────────────────────────────────────────────────────┤
│ 4. Debt Graph                                                    │
│    BrokenInvariant → Finding → Evidence → Remediation            │
├────────────────────────────────────────────────────────────────┤
│ 5. Harness Proof Runtime                                         │
│    Gate → ProofDepth → Fixture → Command → Result                │
├────────────────────────────────────────────────────────────────┤
│ 6. Evidence Ledger                                               │
│    append-only event chain + digest + projection                  │
├────────────────────────────────────────────────────────────────┤
│ 7. Release Policy Decider                                        │
│    no drift + no blocking debt + proof graph complete             │
├────────────────────────────────────────────────────────────────┤
│ 8. Downstream Conformance Lab                                    │
│    render/build/test/debt/release/adoption proof                  │
├────────────────────────────────────────────────────────────────┤
│ 9. Self-improving Loop                                           │
│    escape → hypothesis → patch → fixture → detector → policy      │
└────────────────────────────────────────────────────────────────┘
```

### 6.2 推荐目录结构

```text
.xlib/
  facts/
    xlib.yaml
    tools.yaml
    gates.yaml
    downstream.yaml
    release.yaml
    projections.yaml
  architecture/
    layers.yaml
    import-rules.yaml
    package-boundaries.yaml
  debt/
    rules.yaml
    rule-registry.yaml
    exceptions.yaml
    finding.schema.json
    fixtures/
      architecture/
        pass/
        fail-l2-cycle/
        fail-layer-import/
      dependency/
      docs/
      domain/
      implementation/
      testing/
  goals/
    goal.schema.json
    active/
    completed/
  proof/
    proof.schema.json
    gate-depth.yaml
  evidence/
    event.schema.json
    ledger-policy.yaml
  downstream/
    matrix.yaml
    conformance-profiles.yaml
    adoption-status.yaml
  policies/
    release-policy.yaml
    waiver-policy.yaml

internal/
  factkernel/
  projection/
  architecturegraph/
  debtkernel/
  tracekernel/
  harnesskernel/
  evidencekernel/
  releasepolicy/
  downstreamkernel/

cmd/goalcli/
  main.go                 # thin shell only
```

---

## 7. Canonical Facts：彻底解决漂移

### 7.1 `xlib.yaml`

```yaml
schema_version: xlib-facts/v1
module: github.com/ZoneCNH/xlib-standard
repository: https://github.com/ZoneCNH/xlib-standard
current_release:
  version: v0.4.15
  commit: c0fc3813e156cf35a37ddd0033432a78943bb32b
  tag: v0.4.15
runtime:
  goal_runtime: "3.1"
  governance_runtime: "2.9.3"
tools:
  go: "1.23.0"
  golangci_lint: "v2.1.6"
  govulncheck: "v1.1.4"
security:
  vulncheck_interval_hours: 168
  default_enable_vulncheck: false
release:
  min_score: 9.8
  require_gowork_off: true
  require_clean_tree: true
  manifest: release/manifest/latest.json
  checksum: release/manifest/latest.json.sha256
downstream:
  default_target: kernel
  libraries:
    - kernel
    - configx
    - observex
    - testkitx
    - postgresx
    - redisx
    - kafkax
    - natsx
    - taosx
    - ossx
    - clickhousex
```

### 7.2 `projections.yaml`

```yaml
schema_version: xlib-projections/v1
projections:
  - id: goalcli-version
    source: .xlib/facts/xlib.yaml#/current_release/version
    target: cmd/goalcli/version_gen.go
    mode: generated
  - id: manifest-default-version
    source: .xlib/facts/xlib.yaml#/current_release/version
    target: internal/tools/releasemanifest/version_gen.go
    mode: generated
  - id: readme-release-preflight
    source: .xlib/facts/xlib.yaml#/current_release/version
    target: README.md
    marker: "<!-- xlib:release-preflight -->"
    mode: managed-block
  - id: harness-release-preflight
    source: .xlib/facts/xlib.yaml#/current_release/version
    target: .agent/harness/harness.yaml
    mode: structured-yaml
  - id: workflow-tool-versions
    source: .xlib/facts/xlib.yaml#/tools
    target: .github/workflows/*.yml
    mode: structured-yaml
```

### 7.3 新命令

```bash
goalcli fact validate
goalcli fact render --check
goalcli fact render --write
goalcli fact audit --strict
goalcli fact explain version
```

### 7.4 release policy

```text
任何当前事实漂移直接阻断 release。
历史文档允许旧版本，但必须标记：snapshot_date、applies_to_version、historical: true。
```

---

## 8. Debt Kernel：同时解决六类债务

### 8.1 统一 Finding Schema

```json
{
  "schema_version": "debt-finding/v2",
  "id": "debt.architecture.layer-violation",
  "section": "architecture",
  "invariant_id": "ARCH-LAYER-001",
  "severity": "P0",
  "release_blocking": true,
  "proof_depth": "D3",
  "detector": {
    "id": "architecturegraph.import-layer-check",
    "version": "v1",
    "confidence": 0.98
  },
  "location": {
    "path": "pkg/foo/foo.go",
    "line": 12,
    "symbol": "pkg/foo"
  },
  "evidence": {
    "kind": "import_edge",
    "from": "github.com/ZoneCNH/xlib-standard/pkg/foo",
    "to": "github.com/ZoneCNH/xlib-standard/internal/bar",
    "digest": "sha256:..."
  },
  "goal_id": "GOAL-20260605-DEBT-ARCH-001",
  "adr": "docs/adr/ADR-20260605-architecture-layers.md",
  "owner": "standard-maintainer",
  "remediation": "Move dependency behind L1 interface or invert dependency through contract.",
  "first_seen": "2026-06-05T00:00:00Z",
  "expiry": null,
  "waiver": null,
  "status": "open"
}
```

### 8.2 结构债解决方案

新增：

```yaml
# .xlib/architecture/layers.yaml
schema_version: architecture-layers/v1
layers:
  L0:
    packages:
      - contracts/...
      - internal/validation/...
  L1:
    packages:
      - pkg/templatex
      - testkit/...
  L2:
    packages:
      - internal/tools/...
      - internal/debtcheck/...
      - internal/releasequality/...
  L3:
    packages: []
forbidden_edges:
  - from: L0
    to: L1
  - from: L0
    to: L2
  - from: L1
    to: L2
  - from: L2
    to: L2
    unless_declared: true
  - from: any
    to: github.com/ZoneCNH/x.go
```

新增 detectors：

```bash
goalcli architecture graph --json
 goalcli architecture check --layers .xlib/architecture/layers.yaml
 goalcli architecture cycles
 goalcli architecture god-modules --max-fanout 12 --max-fanin 20
 goalcli architecture boundary --public-api pkg/templatex
```

必须检测：

- import layer violation。
- L2 ↔ L2 未声明耦合。
- cycle / SCC。
- public package 反向依赖 internal tooling。
- `cmd/goalcli` dispatch 过大。
- package fan-in/fan-out 异常。
- internal tools 与 release runtime 的边界污染。

### 8.3 实现债解决方案

新增 detectors：

```bash
goalcli implementation duplicate --threshold 30-lines
 goalcli implementation complexity --max-cyclomatic 12 --max-cognitive 18
 goalcli implementation hotspot --since 90d --top 20
 goalcli implementation obsolete-patterns
 goalcli implementation god-function --max-lines 120
```

重点修复对象：

1. `cmd/goalcli/main.go` 的巨大 dispatch。
2. `scripts/check_docs.sh` 的海量 needle。
3. external shell script 与 Go command 的重复职责。
4. release/version/tool facts 多处 hardcode。

### 8.4 测试债解决方案

新增：

```bash
goalcli testing coverage --risk-weighted
 goalcli testing pyramid --unit-min 70 --integration-max 25 --e2e-max 10
 goalcli testing mutation --profile smoke
 goalcli testing flaky --history release/evidence/test-history.jsonl
 goalcli testing negative-fixtures --all
```

测试不再只看 `go test ./...`，而要看：

```text
risk-weighted coverage
mutation score
negative fixture pass/fail
flaky rate
gate escape count
public API coverage
schema coverage
release manifest fixture isolation
```

### 8.5 文档债解决方案

新增：

```bash
goalcli docs graph-check
 goalcli docs adr-check
 goalcli docs fact-drift
 goalcli docs snapshot-index
 goalcli docs public-api-sync
```

规则：

1. 当前事实只来自 `.xlib/facts`。
2. README、docs、workflow 命令片段由 managed block 生成。
3. 历史报告必须带：

```yaml
historical_snapshot: true
snapshot_date: 2026-06-02
applies_to_version: v0.3.7
superseded_by: docs/project-structural-analysis-20260605.md
```

4. 任何 public API / release gate / evidence schema / security policy / downstream contract 改动必须有 ADR 或 decision log。

### 8.6 依赖债解决方案

新增：

```bash
goalcli dependency sbom --format cyclonedx
 goalcli dependency vuln --require-fresh 168h
 goalcli dependency purpose-check
 goalcli dependency license-check
 goalcli dependency deprecated-check
 goalcli dependency owner-check
```

每个 direct dependency 必须有：

```yaml
module: golang.org/x/tools
purpose: static analysis and repository tooling
owner: standard-maintainer
runtime: false
adr: null
license_policy: allowed
last_reviewed: 2026-06-05
review_interval_days: 90
```

### 8.7 领域债解决方案

新增：

```yaml
# .xlib/domain/language.yaml
schema_version: domain-language/v1
bounded_contexts:
  standard_source:
    allowed_terms:
      - Standard Source
      - Harness
      - Evidence Runtime
      - Goal Runtime
    forbidden_terms:
      - baselib-template
      - foundationx
  downstream_library:
    allowed_terms:
      - kernel
      - configx
      - redisx
    forbidden_terms:
      - x.go business model
```

新增命令：

```bash
goalcli domain language-check
 goalcli domain bounded-context-check
 goalcli domain banned-term-check
 goalcli domain migration-term-check
```

---

## 9. Harness 自证系统

### 9.1 Gate 必须有 negative fixture

每个 release-blocking gate 都必须证明：

```text
好样本通过。
坏样本失败。
失败原因可解释。
exit code 符合标准。
Evidence 记录输入/输出 digest。
```

### 9.2 Fixture 目录

```text
.xlib/harness/fixtures/
  architecture/
    pass-clean/
    fail-l2-cycle/
    fail-layer-import/
    fail-legacy-x-import/
  docs/
    pass-current-facts/
    fail-version-drift/
    fail-missing-adr/
  dependency/
    pass-approved-purpose/
    fail-latest/
    fail-missing-purpose/
  evidence/
    pass-ledger-chain/
    fail-broken-checksum/
    fail-missing-command/
  traceability/
    pass-full-chain/
    fail-missing-pr/
    fail-missing-evidence-digest/
  downstream/
    pass-adopted-fixture/
    fail-false-adoption-claim/
```

### 9.3 新命令

```bash
goalcli harness test --all
 goalcli harness test --gate architecture
 goalcli harness mutation --profile release
 goalcli harness proof --write release/proof/harness-proof.json
```

### 9.4 release policy

```text
任何 D4+ release-blocking gate 如果没有 D5 negative fixture，不能升级为 release-final gate。
```

---

## 10. Proof Graph：从 passed 升级为可计算证明

### 10.1 证明图节点

```text
Fact
Rule
Invariant
Goal
Requirement
AcceptanceCriteria
Task
Issue
Commit
PullRequest
Gate
CommandRun
EvidenceEvent
Artifact
Release
DownstreamRepo
```

### 10.2 证明图边

```text
derived_from
implements
verifies
blocks
produces
hashes
requires
supersedes
adopts
waives
expires
```

### 10.3 schema 示例

```json
{
  "schema_version": "proof-graph/v1",
  "subject": {
    "repo": "ZoneCNH/xlib-standard",
    "commit": "c0fc3813e156cf35a37ddd0033432a78943bb32b",
    "version": "v0.4.15"
  },
  "nodes": [
    {
      "id": "gate:fact-audit",
      "type": "Gate",
      "proof_depth": "D5",
      "release_blocking": true
    },
    {
      "id": "event:fact-audit:20260605",
      "type": "EvidenceEvent",
      "sha256": "..."
    }
  ],
  "edges": [
    {
      "from": "gate:fact-audit",
      "to": "event:fact-audit:20260605",
      "type": "produces"
    }
  ],
  "decision": {
    "release_allowed": false,
    "reasons": [
      "current release v0.4.15 conflicts with v0.4.13 projection in goalcli/releasemanifest/harness"
    ]
  }
}
```

---

## 11. Evidence Ledger：证据事实系统

### 11.1 event schema

```json
{
  "schema_version": "evidence-event/v1",
  "event_id": "evt_20260605_fact_audit_001",
  "previous_event_sha256": "...",
  "timestamp": "2026-06-05T12:00:00Z",
  "repo": "ZoneCNH/xlib-standard",
  "commit": "c0fc3813e156cf35a37ddd0033432a78943bb32b",
  "context": "release_verify",
  "command": "GOWORK=off go run ./cmd/goalcli fact audit --strict",
  "exit_code": 1,
  "inputs": [
    {"path": ".xlib/facts/xlib.yaml", "sha256": "..."},
    {"path": "cmd/goalcli/governance.go", "sha256": "..."}
  ],
  "outputs": [
    {"path": "release/proof/fact-audit.json", "sha256": "..."}
  ],
  "policy": {
    "id": "release-policy/v1",
    "sha256": "..."
  },
  "result": {
    "status": "failed",
    "findings": ["fact.version.projectReleaseVersion.v0.4.13"]
  }
}
```

### 11.2 原则

```text
release/manifest/latest.json 不是事实源。
它只是 evidence ledger + proof graph 的摘要投影。
```

---

## 12. Release Policy：禁止漂移与禁止伪完成

### 12.1 release 条件

release final 必须同时满足：

```text
1. current version fact 与 tag 一致。
2. 所有 current projection 无漂移。
3. 所有 P0 rule active，并有 implemented detector。
4. 所有 release-blocking gate proof_depth >= D4。
5. 所有 release-blocking gate 有 negative fixture，proof_depth >= D5。
6. debt report 无 P0，无 release-blocking P1/P2，无 expired waiver。
7. traceability graph 完整。
8. evidence ledger hash chain 完整。
9. release manifest 是 ledger/proof 的投影，不是手工产物。
10. security vuln evidence fresh。
11. downstream adoption claim 必须有 D7 proof；否则只能写 not_adopted/gap_report。
```

### 12.2 policy 伪代码

```text
release_allowed =
  facts.current_release.version == git.tag
  && fact_audit.status == passed
  && all(projections.status == in_sync)
  && all(blocking_gates.depth >= D4)
  && all(blocking_gates.negative_fixture == passed)
  && debt.p0 == 0
  && debt.release_blocking == 0
  && traceability.full_lifecycle == passed
  && evidence.ledger_chain == valid
  && security.vulncheck_age_hours <= 168
  && no_false_downstream_adoption_claim
```

---

## 13. Goal 化生产流

### 13.1 Goal 不再是 TODO

Goal 应成为一等生产对象：

```yaml
schema_version: goal/v2
id: GOAL-20260605-FACT-KERNEL-001
title: Build canonical fact kernel and strict drift audit
intent: Eliminate v0.4.15/v0.4.13 release fact drift
invariants:
  - FACT-RELEASE-001: current release version has exactly one source
acceptance_criteria:
  - fact audit fails on stale v0.4.13 projection
  - fact audit passes after generated projection update
artifacts:
  - .xlib/facts/xlib.yaml
  - internal/factkernel/...
  - cmd/goalcli/fact.go
negative_fixtures:
  - .xlib/harness/fixtures/docs/fail-version-drift
release_impact:
  required: true
rollback:
  command: git revert <commit>
evidence_schema:
  - evidence-event/v1
```

### 13.2 Goal lifecycle

```text
Draft
 → Context Restored
 → Invariants Defined
 → Acceptance Proof Defined
 → Patch Proposed
 → Harness Executed
 → Evidence Written
 → Policy Decided
 → Released / Blocked
 → Retrospective
 → Rule/Detector/Fixture Upgraded
```

---

## 14. AutoResearch 与自我进化

### 14.1 AutoResearch 的边界

AutoResearch 只能做：

```text
发现变更
提出假设
生成 patch proposal
找相关规则/ADR/历史证据
生成候选 detector
生成候选 fixture
```

不能做：

```text
自行宣布 done
自行把 planned 升级为 implemented
自行把 artifact_exists 升级为 release_usable
自行声明 downstream adopted
```

### 14.2 自我进化闭环

```text
Gate Escape
 → Root Cause
 → Broken Invariant
 → New Rule
 → New Detector
 → Negative Fixture
 → Harness Proof
 → Policy Update
 → Evidence Event
 → Retrospective
 → Next Goal
```

### 14.3 复利机制

每次失败都不是单次修复，而是系统能力增长：

| 失败类型 | 复利产物 |
|---|---|
| 版本漂移 | fact projection detector |
| 文档与代码不一致 | docs graph detector |
| L2 cycle | architecture graph fixture |
| false adoption claim | downstream conformance policy |
| security workflow 失败 | freshness dashboard |
| traceability 漏洞 | proof graph edge detector |
| repeated patch hotspot | implementation hotspot detector |

---

## 15. 最小可行行动 MVA

目标：用最少改动，把系统从“强治理但会漂移”推进到“当前事实不可漂移”。

### MVA-1：建立当前版本事实源

新增：

```text
.xlib/facts/xlib.yaml
```

记录：

```text
current_release.version = v0.4.15
commit = c0fc381...
go = 1.23.0
golangci-lint = v2.1.6
govulncheck = v1.1.4
release_min_score = 9.8
```

### MVA-2：新增 fact audit

新增：

```bash
goalcli fact audit --strict
```

首批检查：

```text
cmd/goalcli/governance.go 不得 hardcode 旧版本
internal/tools/releasemanifest/main.go 不得 hardcode 旧版本
README release-preflight 不得指向旧版本
.agent/harness/harness.yaml 不得指向旧版本
current docs 不得把历史版本当当前事实
```

### MVA-3：修复 v0.4.13 漂移

用 generated file 或 managed block 替代手写常量。

### MVA-4：给 gate 增加 proof_depth 字段

先只改 Harness schema，不要求一次全部升级。

```yaml
- id: docs_check
  command: GOWORK=off make docs-check
  proof_depth: D1
  release_blocking: true
  target_depth: D4
```

### MVA-5：升级 debt finding schema

保持兼容旧 schema，但新增 optional fields：

```text
invariant_id
release_blocking
proof_depth
owner
expiry
remediation
detector
```

### MVA-6：加第一个 negative fixture

优先选择版本漂移，因为当前已经真实发生：

```text
.xlib/harness/fixtures/fact/fail-version-drift/
```

### MVA-7：把 traceability 状态改成 partial

同步：

```text
iron-rules.md
command-implementation-status.yaml
docs/reports/rules-deep-analysis-20260605.md
```

状态应从 gap/implemented 二选一改为：

```text
partial_implemented
proof_depth: D3
full_lifecycle_graph: gap
```

### MVA-8：release gate 加 fact audit

```makefile
release-check: require-gowork-off fact-audit ci integration ...
```

---

## 16. 1 天行动计划

当天目标：阻断当前已观察到的事实漂移。

1. 新增 `.xlib/facts/xlib.yaml`。
2. 新增 `internal/factkernel` 最小实现。
3. 新增 `goalcli fact audit --strict`。
4. 检测并修复所有当前 `v0.4.13` 投影漂移。
5. 修改 README 与 Harness 的 release-preflight 命令为 fact-generated block。
6. 给 `docs/project-structural-analysis-20260605.md` 加 historical/current status 注记，避免把 v0.4.13 对齐结果误读为 v0.4.15 当前事实。
7. 在 Makefile 增加：

```makefile
fact-audit:
	$(GOALCLI) fact audit --strict
```

8. 将 `fact-audit` 加入 `release-check` 与 `governance-check`。

验收：

```bash
GOWORK=off make fact-audit
GOWORK=off make docs-check
XLIB_CONTEXT=release_verify GOWORK=off make release-check
```

---

## 17. 7 天行动计划

7 天目标：从“漂移修复”升级为“可证明不漂移”。

1. 完成 `projections.yaml`。
2. 把 version/tool/gate/downstream facts 迁入 `.xlib/facts`。
3. 给 Harness gate 增加字段：

```text
proof_depth
release_blocking
authority
negative_fixture_required
target_depth
```

4. 为以下 gate 建 negative fixture：

```text
fact-audit
docs-check
architecture
dependency-debt
security
traceability-check
adoption-check
```

5. debt finding schema 升级到 v2。
6. architecture detector 从 legacy import 升级到 import graph。
7. docs-check 保留，但新增 `docs graph-check`。
8. Traceability 状态统一为 partial，并新增 full graph 设计。
9. Security workflow 输出 freshness artifact。
10. Release manifest 引用 debt evidence 与 security freshness evidence。

验收：

```bash
GOWORK=off make context-release
GOWORK=off make harness-fixtures
GOWORK=off make fact-audit
GOWORK=off make debt-evidence
```

---

## 18. 30 天行动计划

30 天目标：完整 Debt Control Plane + Proof Graph + Evidence Ledger 初版。

### 第 1 周：Fact Kernel / Drift Lock

- 所有当前事实收敛到 `.xlib/facts`。
- 所有手写版本号/工具版本改为 generated 或 audited projection。
- release gate 阻断 drift。

### 第 2 周：Debt Kernel v2

- architecture graph detector。
- implementation complexity/duplication/hotspot detector。
- testing risk coverage / pyramid detector。
- dependency SBOM/purpose/freshness detector。
- domain language detector。
- docs ADR/fact graph detector。

### 第 3 周：Proof Graph / Evidence Ledger

- evidence event schema。
- command run event writer。
- proof graph projection。
- manifest 从 proof/ledger 投影生成。
- release policy decider。

### 第 4 周：Downstream Conformance Lab

- kernel conformance fixture。
- configx/redisx representative conformance。
- adoption proof schema。
- false adoption claim gate。
- downstream proof dashboard。

验收：

```bash
XLIB_CONTEXT=release_verify GOWORK=off make release-final-check
GOWORK=off goalcli proof verify --release v0.4.x
GOWORK=off goalcli downstream conformance --matrix .xlib/downstream/matrix.yaml
```

---

## 19. 90 天路线图

90 天目标：把 `xlib-standard` 升级为可复用的标准生产平台。

### 0-30 天：内核稳定

- fact kernel。
- drift audit。
- debt v2。
- proof depth。
- fixture harness。
- ledger v1。

### 31-60 天：下游证明

- kernel、configx、redisx、observex representative proof。
- downstream release/adoption status 从 status file 升级为 proof file。
- conformance profile 支持 L0/L1/L2。

### 61-90 天：自我进化

- escape analyzer。
- retrospective patch generator。
- detector suggestion generator。
- rule conflict detector。
- generated-artifacts rebuild verifier。
- automatic fact projection PR。

---

## 20. 指标体系

### 20.1 禁止漂移指标

```text
fact_drift_count = 0
current_version_projection_drift = 0
historical_doc_without_snapshot_metadata = 0
tool_version_drift = 0
gate_command_drift = 0
```

### 20.2 Harness 指标

```text
release_blocking_gate_count
gate_with_proof_depth_D4_plus_ratio
gate_with_negative_fixture_ratio
gate_escape_count
fixture_mutation_score
```

### 20.3 Debt 指标

```text
p0_debt_count = 0
release_blocking_debt_count = 0
waiver_expired_count = 0
architecture_cycle_count = 0
l2_to_l2_undeclared_edges = 0
duplication_score
complexity_hotspot_score
risk_weighted_coverage
```

### 20.4 Evidence 指标

```text
evidence_event_chain_valid = true
manifest_is_projection = true
command_run_events_coverage
artifact_digest_coverage
security_evidence_age_hours <= 168
release_policy_decision_reproducible = true
```

### 20.5 Downstream 指标

```text
downstream_targets_total
downstream_D7_adopted_count
downstream_gap_count
false_adoption_claim_count = 0
conformance_replay_success_rate
```

### 20.6 自我进化指标

```text
escape_to_rule_patch_rate
recurring_issue_rate
mean_time_to_detector
mean_time_to_negative_fixture
retrospective_action_completion_rate
```

---

## 21. 第一批建议 Patch 列表

### Patch 1：fix-current-release-facts

目标：消除 `v0.4.15` / `v0.4.13` 漂移。

涉及：

```text
.xlib/facts/xlib.yaml
cmd/goalcli/version_gen.go
internal/tools/releasemanifest/version_gen.go
README.md managed block
.agent/harness/harness.yaml
Makefile fact-audit target
```

### Patch 2：gate-proof-depth-schema

目标：给 Harness gate 增加 proof depth。

涉及：

```text
.agent/harness/harness.yaml
contracts/harness.schema.json
cmd/goalcli/harness_check.go
```

### Patch 3：debt-finding-v2

目标：升级 Finding schema。

涉及：

```text
internal/debtcheck
.agent/policies/debt/rules.yaml
.agent/registries/debt/rule-registry.yaml
contracts/debt-finding.schema.json
```

### Patch 4：architecture-graph-check

目标：从 legacy import 扫描升级到 layer/cycle detector。

涉及：

```text
.xlib/architecture/layers.yaml
internal/architecturegraph
cmd/goalcli/architecture.go
.xlib/harness/fixtures/architecture
```

### Patch 5：traceability-state-reconcile

目标：消除 traceability gap/implemented 事实冲突。

涉及：

```text
.agent/rules/iron-rules.md
.agent/registries/command-implementation-status.yaml
docs/reports/rules-deep-analysis-20260605.md
cmd/goalcli/traceability.go
```

状态改为：

```text
partial_implemented
proof_depth: D3
full_lifecycle_graph: gap
```

### Patch 6：evidence-ledger-v1

目标：manifest 从事实源降级为 projection。

涉及：

```text
release/evidence/events/
internal/evidencekernel
internal/tools/releasemanifest
scripts/check_release_evidence.sh
```

---

## 22. 最终推荐路径

### 不推荐

```text
继续往 README、Makefile、workflow、.agent 和 Go 常量中手工同步事实。
```

这会继续产生版本漂移、命令漂移、规则漂移、证据漂移。

### 推荐

```text
第一阶段：不拆仓，先单仓强分层。
第二阶段：建立 Canonical Fact Kernel。
第三阶段：把 Harness gate 类型化为 proof depth。
第四阶段：把 Debt Control Plane 从 marker 扫描升级为 graph/invariant 检测。
第五阶段：把 Evidence 从 manifest 升级为 append-only ledger + proof graph。
第六阶段：把 Goal Runtime 和 AutoResearch 接入 escape → detector → fixture → policy 的自我进化闭环。
第七阶段：用 Downstream Conformance Lab 证明 adoption，而不是声明 adoption。
```

### 最高优先级行动

立即做这三件事：

```text
1. 新增 .xlib/facts/xlib.yaml，修复 v0.4.15/v0.4.13 漂移。
2. 新增 goalcli fact audit --strict，并加入 release-check。
3. 把 traceability-check 状态从“gap / implemented”冲突改为“partial implemented + target full graph”。
```

这三件事完成后，系统才有资格继续谈“彻底禁止漂移”和“自我进化”。

---

## 23. 附录：本次分析用到的关键事实文件

```text
README.md
Makefile
cmd/goalcli/main.go
cmd/goalcli/governance.go
cmd/goalcli/adoption_check.go
cmd/goalcli/debt.go
cmd/goalcli/traceability.go
internal/debtcheck/debtcheck.go
internal/tools/releasemanifest/main.go
.agent/harness/harness.yaml
.agent/evidence/truth-state.yaml
.agent/registries/command-implementation-status.yaml
.agent/policies/debt/rules.yaml
.agent/registries/debt/rule-registry.yaml
.github/workflows/ci.yml
.github/workflows/release.yml
.github/workflows/security.yml
.github/rulesets/protect-main.json
scripts/check_docs.sh
docs/project-structural-analysis-20260605.md
docs/reports/rules-deep-analysis-20260605.md
.agent/rules/iron-rules.md
go.mod
.tool-versions
.golangci.yml
```
