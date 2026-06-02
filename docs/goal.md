# xlib-standard 10分标准升级 Goal 方案 v1.0

> 文件类型：完整 Goal Runtime 可执行方案  
> 目标仓库：`https://github.com/ZoneCNH/xlib-standard`  
> 关联基础库：`https://github.com/ZoneCNH/kernel`、`https://github.com/ZoneCNH/xlib-standard`  
> 关联应用：`x.go`  
> 当前评分基线：7.1 / 10  
> 目标评分：10 / 10  
> 执行标准：Goal Runtime Prompt v3.1  
> 默认执行模式：Full  
> 当前日期：2026-06-02  
> 时区：Asia/Tokyo  
> 完成声明格式：`DONE with evidence:`

---

## 0. 执行总提示

你是负责 `x.go` 基础库体系标准化、工程化和复利资产化的架构执行 Agent。

你的任务不是写一份普通方案，而是把当前 `xlib-standard` 从“可用的 Go 基础库模板仓库”升级为“10 分标准的基础库标准源 + 模板工厂 + Harness/Evidence 运行时”。

你必须按照 Goal Runtime Prompt v3.1 执行完整闭环：

```text
Goal
  → Context Recovery
  → Spec
  → Design
  → Plan
  → Tasks
  → Execution
  → Verification
  → Evidence
  → Review
  → Release
  → Retrospective
  → Self-improving
```

不允许只输出计划。  
不允许只修改 README。  
不允许没有 Evidence 就声称完成。  
不允许继续保留 `baselib-template` / `foundationx` 旧命名作为主口径。  
不允许让 `xlib-standard`、`kernel`、`x.go` 三者边界继续模糊。

---

## 1. 问题的底层本质

当前 `xlib-standard` 的核心问题不是“缺少文件”，而是 **身份、边界、执行证据和下游继承关系没有完全一致**。

现状可以概括为：

```text
仓库名：xlib-standard
go.mod：github.com/ZoneCNH/baselib-template
README 主标题：baselib-template
生成示例：foundationx
最新目标命名：xlib-standard + kernel
实际能力：标准文档 + Go 模板 + generator + Harness + Evidence
边界文档：又说 xlib-standard 不应包含 generator/runtime 实现
```

这导致系统存在结构性矛盾：

```text
xlib-standard 到底是标准源？
还是模板实现仓库？
baselib-template 是否废弃？
foundationx 是否已被 kernel 替代？
下游库到底继承哪个标准？
Agent 执行时以哪个仓库、哪个 module path、哪个命名为事实？
```

因此，10 分升级的第一性目标不是“补几个脚本”，而是：

> 建立一个唯一、稳定、可执行、可验证、可继承、可复利演进的基础库标准工厂。

---

## 2. 不可再拆解的基本真理

### 2.1 标准必须有唯一事实源

基础库体系中，所有下游库必须能回答：

```text
我的标准从哪里继承？
我的模板从哪里生成？
我的 gate 从哪里来？
我的 Evidence 如何证明？
我和 x.go 的边界在哪里？
我是否允许依赖 kernel？
```

如果标准源不唯一，下游库会分叉。

### 2.2 模板必须可运行，不只是可阅读

10 分标准不是 README 写得完整，而是：

```text
make ci 能跑
make release-check 能跑
生成库能跑
contracts 能校验
boundary 能阻断错误依赖
Evidence 能证明完成状态
```

### 2.3 Evidence 是完成声明的一部分

没有 Evidence，不能声明 DONE。

必须使用：

```text
DONE with evidence:
- scope:
- gates:
- artifacts:
- known gaps:
```

### 2.4 基础库不得反向依赖业务层

`xlib-standard`、`kernel`、`postgresx`、`redisx`、`kafkax`、`taosx`、`ossx`、`clickhousex` 都不得依赖：

```text
github.com/bytechainx/x.go
github.com/ZoneCNH/x.go
x.go/internal/*
MacroRegime
MarketRegime
TradingSignal
Kline
OrderBook
Position
RiskGate
BTCUSDT
ETHUSDT
```

### 2.5 kernel 是 L0 基础能力，不是业务框架

`kernel` 的职责是提供：

```text
context
error
config interface
logging interface
metrics interface
lifecycle
clock
retry
backoff
health
trace
version
```

它不应该包含：

```text
PostgreSQL runtime
Redis runtime
Kafka runtime
TDengine runtime
OSS runtime
x.go 业务模型
交易模型
宏观状态模型
```

### 2.6 Gate 必须可机器验证

人审可以提高质量，但不能代替机器 gate。

10 分标准必须有：

```text
Semantic Gate
Executable Gate
Hybrid Gate
Evidence Gate
Release Gate
Retrospective Gate
```

---

## 3. 被误认为真理的常见假设

| 假设 | 为什么是错的 | 正确裁决 |
|---|---|---|
| 仓库名改成 `xlib-standard` 就已经完成标准化 | module path、文档、脚本、CI、生成示例仍可能保留旧名 | 必须全链路命名统一 |
| README 说明标准即可 | 下游生成、CI、contracts、Evidence 才是真正标准 | 文档必须被 gate 约束 |
| `baselib-template` 可以长期并存 | 会形成双标准、双入口、双继承关系 | 要么迁移，要么废弃并写入 ADR |
| `foundationx` 只是旧名，不影响执行 | integration、docs、生成示例会继续污染下游 | 必须统一为 `kernel` |
| Shell/Python gate 足够 | 可用，但不利于 x.go Go 化治理和长期维护 | 核心 gate 应 Go 化 |
| Evidence manifest 有了就足够 | 还需要 workflow artifact、tag、commit、tree、source digest、contract fingerprint 的强绑定 | Evidence 需要可复验 |
| xlib-standard 可以同时什么都做 | 可以，但必须在边界文档中明确裁决 | 需要 ADR 锁定角色 |

---

## 4. 可以被打破的限制

### 4.1 可以打破“标准仓库不能包含模板实现”的限制

如果当前仓库已经包含 Go 模板、generator、Harness、Evidence，那么可以正式裁决：

```text
xlib-standard = Standard Source + Go Reference Template + Harness/Evidence Runtime
```

这是推荐路径。

### 4.2 可以打破“每个基础库手工复制模板”的限制

通过 generator 统一渲染：

```text
xlib-standard
  → kernel
  → configx
  → observex
  → postgresx
  → redisx
  → kafkax
  → taosx
  → ossx
  → clickhousex
```

### 4.3 可以打破“gate 散落在 shell/python”的限制

升级为统一 Go CLI：

```text
cmd/xlibgate
  docs-check
  boundary
  secrets
  contracts
  render
  render-check
  evidence
  release-check
  score
```

### 4.4 可以打破“标准不可度量”的限制

引入 10 分评分器：

```text
xlibgate score
```

输出：

```json
{
  "score": 9.6,
  "failed_dimensions": [],
  "warnings": [],
  "evidence": "release/manifest/latest.json"
}
```

---

## 5. 从零设计的新方案

### 5.1 目标架构

```text
                    ┌──────────────────────────────┐
                    │        xlib-standard          │
                    │ Standard + Template + Gates   │
                    │ Evidence + Goal Runtime       │
                    └──────────────┬───────────────┘
                                   │
             ┌─────────────────────┼─────────────────────┐
             │                     │                     │
             ▼                     ▼                     ▼
        ┌─────────┐          ┌──────────┐          ┌──────────┐
        │ kernel  │          │ configx  │          │ observex │
        │   L0    │          │   L0/L1  │          │   L0/L1  │
        └────┬────┘          └────┬─────┘          └────┬─────┘
             │                    │                     │
             ├────────────┬───────┴───────────┬─────────┘
             ▼            ▼                   ▼
       ┌──────────┐  ┌──────────┐       ┌──────────┐
       │ postgresx│  │  redisx  │       │  kafkax  │
       └────┬─────┘  └────┬─────┘       └────┬─────┘
            │             │                  │
            ├─────────────┼──────────────────┤
            ▼             ▼                  ▼
       ┌──────────┐  ┌──────────┐       ┌────────────┐
       │  taosx   │  │  ossx    │       │clickhousex │
       └────┬─────┘  └────┬─────┘       └────┬───────┘
            │             │                  │
            └─────────────┴──────────┬───────┘
                                      ▼
                                   ┌──────┐
                                   │ x.go │
                                   └──────┘
```

### 5.2 推荐角色裁决

本方案采用方案 B：

```text
xlib-standard 不是纯文档仓库。
xlib-standard 是基础库标准源 + Go 标准模板实现 + Harness/Evidence runtime。
```

因此需要修改边界文档：

```text
旧口径：
xlib-standard 禁止 generator、模板脚本和 runtime 实现。

新口径：
xlib-standard 允许包含 Go Reference Template、generator、xlibgate、contracts、Harness、Evidence runtime；
但禁止包含 profile runtime、业务模型、真实生产连接、x.go 依赖。
```

### 5.3 标准分层

```text
Standard Layer:
  xlib-standard

L0 Kernel Layer:
  kernel

L0/L1 Common Layer:
  configx
  observex
  testkitx

L1 Infrastructure Adapter Layer:
  postgresx
  redisx
  kafkax
  taosx
  ossx
  clickhousex

L2 Technical Composition Layer:
  storagex
  cachex
  eventx
  datax

Business/Application Layer:
  x.go
```

### 5.4 命名统一表

| 旧名 | 新名 | 裁决 |
|---|---|---|
| `baselib-template` | `xlib-standard` | 旧名废弃，仅在 migration ADR 中保留 |
| `github.com/ZoneCNH/baselib-template` | `github.com/ZoneCNH/xlib-standard` | 必须修改 go.mod 与所有 imports |
| `foundationx` | `kernel` | 旧名废弃 |
| `github.com/ZoneCNH/foundationx` | `github.com/ZoneCNH/kernel` | integration 示例改为 kernel |
| `pkg/templatex` | `pkg/templatex` 或 `pkg/xlibtemplate` | 可保留作为 reference template package，但文档必须解释 |
| `templatex_` metrics prefix | `xlib_template_` 或生成库 prefix | 建议生成库自动替换 |

### 5.5 生成链路

```bash
go run ./cmd/xlibgate render \
  --module-name kernel \
  --module-path github.com/ZoneCNH/kernel \
  --package-name kernel \
  --out ../kernel
```

下游生成后必须通过：

```bash
cd ../kernel
GOWORK=off make release-check
```

---

## 6. Goal Runtime v3.1 对象模型

### 6.1 Goal

```text
GOAL-20260602-001
将 ZoneCNH/xlib-standard 升级为 10 分标准基础库工厂，完成仓库角色裁决、命名统一、kernel 对齐、Go 化 gate、Evidence 强化、下游生成兼容、Goal Runtime v3.1 工程化和评分器落地，使其成为 x.go 基础库体系的唯一标准源与模板执行源。
```

### 6.2 Spec

```text
SPEC-xlib-standard-v2.0
```

### 6.3 Design

```text
DESIGN-xlib-standard-v2.0
```

### 6.4 Plan

```text
PLAN-GOAL-20260602-001-v1.0
```

### 6.5 状态机

```text
INIT
  → CONTEXT_READY
  → GOAL_READY
  → SPEC_READY
  → DESIGN_READY
  → PLAN_READY
  → TASKS_READY
  → EXECUTING
  → VERIFYING
  → REVIEWING
  → RELEASING
  → RETROSPECTING
  → DONE
```

异常状态：

```text
BLOCKED
FAILED
NEEDS_RESEARCH
NEEDS_DECISION
NEEDS_REPLAN
NEEDS_ROLLBACK
NEEDS_HUMAN_APPROVAL
INCONSISTENT_STATE
```

---

## 7. Requirements

### REQ-001：仓库角色裁决

`xlib-standard` 必须明确自身角色：

```text
Standard Source + Go Reference Template + Generator + Harness + Evidence Runtime
```

不得继续出现自相矛盾的边界描述。

验收标准：

```text
AC-REQ-001-001: docs/adr/ADR-20260602-001-xlib-standard-role.md 存在。
AC-REQ-001-002: README.md 明确 xlib-standard 的五类职责。
AC-REQ-001-003: docs/standard/module-boundary.md 不再声明 xlib-standard 禁止自身 generator/Harness/Evidence 实现。
AC-REQ-001-004: docs/standard/repository-roles.md 不再把 baselib-template 作为主仓库角色。
```

### REQ-002：命名统一

必须完成以下替换：

```text
baselib-template → xlib-standard
github.com/ZoneCNH/baselib-template → github.com/ZoneCNH/xlib-standard
foundationx → kernel
github.com/ZoneCNH/foundationx → github.com/ZoneCNH/kernel
```

保留旧名的唯一允许位置：

```text
docs/adr/*migration*
CHANGELOG migration section
docs/migration/*
```

验收标准：

```text
AC-REQ-002-001: go.mod module 为 github.com/ZoneCNH/xlib-standard。
AC-REQ-002-002: 非 migration 文档中不再出现 baselib-template 主口径。
AC-REQ-002-003: integration 默认渲染 kernel。
AC-REQ-002-004: contracts tests import 新 module path。
AC-REQ-002-005: xlibgate stale-name gate 可阻断旧名回归。
```

### REQ-003：kernel 下游兼容

`xlib-standard` 必须能生成并校验 `kernel` 形态。

验收标准：

```text
AC-REQ-003-001: scripts/run_integration.sh 或 cmd/xlibgate integration 包含 kernel case。
AC-REQ-003-002: kernel case 使用 module path github.com/ZoneCNH/kernel。
AC-REQ-003-003: 渲染后的 kernel 通过 go test ./...。
AC-REQ-003-004: 渲染后的 kernel 通过 make contracts。
AC-REQ-003-005: 渲染后的 kernel 通过 make boundary。
AC-REQ-003-006: 渲染后的 kernel 能生成 release manifest。
```

### REQ-004：核心 Gate Go 化

必须新增统一 CLI：

```text
cmd/xlibgate
```

子命令：

```text
docs-check
boundary
secrets
contracts
render
render-check
integration
evidence
release-check
score
```

验收标准：

```text
AC-REQ-004-001: cmd/xlibgate/main.go 存在。
AC-REQ-004-002: Makefile 核心 gate 优先调用 go run ./cmd/xlibgate。
AC-REQ-004-003: scripts/*.sh 可作为兼容包装，但不能是唯一实现。
AC-REQ-004-004: CI 运行 go-based gate。
```

### REQ-005：Goal Runtime v3.1 工程化

`.agent/` 必须从轻量说明升级为 v3.1 可执行运行时。

必须包含：

```text
.agent/goal-runtime.md
.agent/object-model.md
.agent/state-machine.md
.agent/traceability-matrix.md
.agent/harness.yaml
.agent/evidence-protocol.md
.agent/review-template.md
.agent/release-template.md
.agent/retrospective-template.md
.agent/risk-register.md
.agent/decision-log.md
.agent/rollback-protocol.md
.agent/prompt-patches.md
.agent/harness-patches.md
.agent/rule-patches.md
```

验收标准：

```text
AC-REQ-005-001: .agent/state-machine.md 包含 v3.1 完整状态机。
AC-REQ-005-002: .agent/object-model.md 包含 Goal/Spec/Requirement/AC/Design/ADR/Plan/Task/Test/Evidence/Risk/Decision/Review/Release/Retrospective/Patch。
AC-REQ-005-003: .agent/traceability-matrix.md 至少覆盖本 Goal 的所有 REQ。
AC-REQ-005-004: docs-check 校验 .agent 关键文件存在。
```

### REQ-006：Evidence 强化

release manifest 必须记录：

```text
module
version
commit
tree_sha
source_digest
tracked_file_count
go_version
generated_at
generated_by
tree_state
checks
contracts
dependencies
tools
artifacts
workflow_run_id
artifact_name
artifact_url
score
known_risks
breaking_changes
```

验收标准：

```text
AC-REQ-006-001: release/manifest/template.json 包含 workflow_run_id、artifact_url、score。
AC-REQ-006-002: releasemanifest tool 生成并校验新增字段。
AC-REQ-006-003: CI summary 输出 manifest SHA256、artifact name、workflow URL。
AC-REQ-006-004: release-final-check 要求 tree_state=clean。
```

### REQ-007：10 分评分器

必须新增评分器：

```bash
go run ./cmd/xlibgate score
```

评分维度：

```text
repository_identity
naming_consistency
standard_boundary
go_module_integrity
api_template
contracts
tests
harness
ci
release_evidence
security
goal_runtime
kernel_compatibility
downstream_generation
documentation
retrospective
```

验收标准：

```text
AC-REQ-007-001: score 输出 JSON。
AC-REQ-007-002: score < 9.5 时 release-final-check 失败。
AC-REQ-007-003: score 结果写入 release manifest。
AC-REQ-007-004: docs/scorecard.md 说明评分规则。
```

### REQ-008：安全 Gate 升级

保留现有 grep gate，同时新增更强扫描模式：

```text
gitleaks 或等价 Go 实现扫描
private key pattern
AWS key
GitHub token
connection string
.env 泄漏
/home/k8s/secrets/env/* 内容泄漏
```

验收标准：

```text
AC-REQ-008-001: make security 包含 govulncheck。
AC-REQ-008-002: make security 包含 secret scan。
AC-REQ-008-003: secret scan 可识别常见 token/private key/connection string。
AC-REQ-008-004: 文档声明不得把 /home/k8s/secrets/env/* 写入源码、README、测试日志、Release Manifest、PR 描述。
```

### REQ-009：下游基础库生成矩阵

integration 必须至少覆盖：

```text
kernel
configx
redisx
postgresx
kafkax
taosx
ossx
clickhousex
```

MVA 阶段可先覆盖：

```text
kernel
corekit
```

Full 阶段覆盖全矩阵。

验收标准：

```text
AC-REQ-009-001: docs/downstream-matrix.md 存在。
AC-REQ-009-002: integration 至少覆盖 kernel。
AC-REQ-009-003: Full mode 覆盖所有目标库。
AC-REQ-009-004: 每个库都有预期 module path、package name、layer、allowed deps、forbidden deps。
```

### REQ-010：x.go 集成边界

必须明确：

```text
x.go 可以消费基础库。
基础库不得依赖 x.go。
x.go 的 Market Data / Macro Data / Regime Engine 不下沉到基础库。
x.go secrets 只能由调用方显式传入。
```

验收标准：

```text
AC-REQ-010-001: docs/xgo-integration-boundary.md 存在。
AC-REQ-010-002: boundary gate 禁止 x.go import。
AC-REQ-010-003: boundary gate 禁止业务词污染。
AC-REQ-010-004: docs 明确 /home/k8s/secrets/env/* 只作为调用方路径约束，不由基础库默认读取。
```

---

## 8. Design

### 8.1 目录结构目标

```text
xlib-standard/
  README.md
  AGENTS.md
  CHANGELOG.md
  go.mod
  Makefile
  .golangci.yml

  .github/
    workflows/
      ci.yml
      release.yml
      extended.yml

  .agent/
    goal-runtime.md
    object-model.md
    state-machine.md
    traceability-matrix.md
    harness.yaml
    evidence-protocol.md
    review-template.md
    release-template.md
    retrospective-template.md
    risk-register.md
    decision-log.md
    rollback-protocol.md
    prompt-patches.md
    harness-patches.md
    rule-patches.md

  cmd/
    xlibgate/
      main.go
      docscheck/
      boundary/
      secrets/
      contracts/
      render/
      rendercheck/
      integration/
      evidence/
      score/

  pkg/
    templatex/
      config.go
      client.go
      errors.go
      health.go
      metrics.go
      options.go
      version.go
      doc.go

  internal/
    sanitize/
    validation/
    runtime/
    tools/
      releasemanifest/

  contracts/
    config.schema.json
    error.schema.json
    health.schema.json
    metrics.md
    manifest.schema.json

  examples/
    basic/
    config/
    health/

  testkit/

  docs/
    standard/
    adr/
    migration/
    scorecard.md
    downstream-matrix.md
    xgo-integration-boundary.md
    generation.md
    design.md
    spec.md
    testing.md
    release.md

  release/
    manifest/
      template.json

  scripts/
    compatibility wrappers only
```

### 8.2 xlibgate 子命令设计

```text
xlibgate docs-check
  校验 README/docs/.agent 关键文件、链接、占位符、旧名污染。

xlibgate boundary
  校验禁止 x.go 依赖、禁止业务词、禁止 internal 反向依赖 public package。

xlibgate secrets
  校验 token、secret、password、private key、connection string、/home/k8s/secrets/env 泄漏。

xlibgate contracts
  校验 config/error/health/metrics/manifest contracts 与代码常量一致。

xlibgate render
  生成下游库。

xlibgate render-check
  校验生成库无旧名残留、module path 正确、package 正确。

xlibgate integration
  执行下游渲染 smoke。

xlibgate evidence
  生成 release manifest 与 checksum。

xlibgate release-check
  校验 manifest、checksum、clean tree、score threshold。

xlibgate score
  输出 10 分评分 JSON。
```

### 8.3 Evidence 设计

Manifest schema：

```json
{
  "module": "github.com/ZoneCNH/xlib-standard",
  "version": "v0.2.0",
  "commit": "...",
  "tree_sha": "...",
  "source_digest": "sha256:...",
  "tracked_file_count": 0,
  "go_version": "go1.23.x",
  "generated_at": "2026-06-02T00:00:00Z",
  "generated_by": "cmd/xlibgate evidence",
  "tree_state": "clean",
  "score": {
    "value": 10.0,
    "threshold": 9.5,
    "dimensions": {}
  },
  "checks": {},
  "contracts": [],
  "dependencies": [],
  "tools": {},
  "artifacts": [
    "release/manifest/latest.json",
    "release/manifest/latest.json.sha256"
  ],
  "ci": {
    "workflow_run_id": "",
    "workflow_url": "",
    "artifact_name": "",
    "artifact_url": ""
  },
  "notes": {
    "breaking_changes": "none",
    "known_risks": []
  }
}
```

### 8.4 Scorecard 设计

总分 10 分：

| 维度 | 权重 |
|---|---:|
| repository_identity | 0.8 |
| naming_consistency | 0.8 |
| standard_boundary | 0.8 |
| go_module_integrity | 0.6 |
| api_template | 0.7 |
| contracts | 0.7 |
| tests | 0.7 |
| harness | 0.8 |
| ci | 0.6 |
| release_evidence | 0.8 |
| security | 0.7 |
| goal_runtime | 0.8 |
| kernel_compatibility | 0.7 |
| downstream_generation | 0.6 |
| documentation | 0.4 |
| retrospective | 0.3 |

Release threshold：

```text
release-check: score >= 9.0
release-final-check: score >= 9.5
10分发布目标: score == 10.0 或所有 P0/P1 维度满分且总分 >= 9.8
```

---

## 9. ADR

### ADR-20260602-001：xlib-standard 仓库角色裁决

```text
Status: Accepted

Decision:
  xlib-standard 是基础库标准源 + Go reference template + generator + Harness/Evidence runtime。
  不再维护 baselib-template 作为独立主口径。

Rejected:
  1. 拆分为 xlib-standard + xlib-template 两个仓库。
  2. 继续保留 baselib-template 作为事实模板源。

Consequences:
  1. 修改 go.mod 到 github.com/ZoneCNH/xlib-standard。
  2. 所有旧命名迁移到 migration 文档。
  3. module-boundary 更新为允许 reference implementation。
```

### ADR-20260602-002：foundationx 统一迁移为 kernel

```text
Status: Accepted

Decision:
  foundationx 改名为 kernel。
  kernel 是 L0 基础能力库，默认路径 github.com/ZoneCNH/kernel。

Rejected:
  1. 同时支持 foundationx 和 kernel 作为平级主名。
  2. 继续在 integration 中使用 foundationx 示例。

Consequences:
  1. run_integration 默认渲染 kernel。
  2. docs/repository-roles 使用 kernel。
  3. generation examples 使用 kernel。
```

### ADR-20260602-003：核心 Gate Go 化

```text
Status: Accepted

Decision:
  新增 cmd/xlibgate，核心 gate 使用 Go 实现。
  scripts/*.sh 仅保留兼容包装。

Rejected:
  1. 继续让 bash/python 承担主 gate。
  2. 每个 gate 分散实现。

Consequences:
  1. Makefile 调用 xlibgate。
  2. CI 以 xlibgate 为标准入口。
  3. 更容易被 x.go 复用。
```

---

## 10. Traceability Matrix

| Requirement | Acceptance Criteria | Design Section | Task | Test | Evidence |
|---|---|---|---|---|---|
| REQ-001 | AC-REQ-001-001~004 | 8.1, 9 | TASK-001~004 | TEST-001 | EVID-001 |
| REQ-002 | AC-REQ-002-001~005 | 5.4, 8.1 | TASK-005~010 | TEST-002 | EVID-002 |
| REQ-003 | AC-REQ-003-001~006 | 5.5, 8.2 | TASK-011~014 | TEST-003 | EVID-003 |
| REQ-004 | AC-REQ-004-001~004 | 8.2 | TASK-015~022 | TEST-004 | EVID-004 |
| REQ-005 | AC-REQ-005-001~004 | 6, 8.1 | TASK-023~030 | TEST-005 | EVID-005 |
| REQ-006 | AC-REQ-006-001~004 | 8.3 | TASK-031~036 | TEST-006 | EVID-006 |
| REQ-007 | AC-REQ-007-001~004 | 8.4 | TASK-037~041 | TEST-007 | EVID-007 |
| REQ-008 | AC-REQ-008-001~004 | 8.2 | TASK-042~045 | TEST-008 | EVID-008 |
| REQ-009 | AC-REQ-009-001~004 | 5.1, 8.2 | TASK-046~050 | TEST-009 | EVID-009 |
| REQ-010 | AC-REQ-010-001~004 | 5.3, 8.2 | TASK-051~054 | TEST-010 | EVID-010 |

---

## 11. Task Breakdown

### Milestone M1：身份裁决与命名统一

```text
TASK-001: 新增 ADR-20260602-001-xlib-standard-role.md
TASK-002: 新增 ADR-20260602-002-kernel-rename.md
TASK-003: 修改 docs/standard/module-boundary.md
TASK-004: 修改 docs/standard/repository-roles.md
TASK-005: 修改 go.mod module path
TASK-006: 替换 imports 中 github.com/ZoneCNH/baselib-template
TASK-007: 替换 README 主标题与标准说明
TASK-008: 替换 AGENTS.md 中旧名
TASK-009: 替换 docs/generation.md 生成示例
TASK-010: 新增 docs/migration/baselib-template-to-xlib-standard.md
```

### Milestone M2：kernel 兼容

```text
TASK-011: 修改 integration 渲染 case 为 kernel
TASK-012: 修改 generation 示例为 kernel
TASK-013: 修改 downstream matrix，加入 kernel
TASK-014: 添加 kernel 渲染后的 stale-name 检查
```

### Milestone M3：xlibgate Go 化

```text
TASK-015: 创建 cmd/xlibgate/main.go
TASK-016: 实现 docs-check 子命令
TASK-017: 实现 boundary 子命令
TASK-018: 实现 secrets 子命令
TASK-019: 实现 contracts 子命令
TASK-020: 实现 render 子命令
TASK-021: 实现 render-check 子命令
TASK-022: 实现 integration 子命令
```

### Milestone M4：Goal Runtime v3.1 工程化

```text
TASK-023: 更新 .agent/goal-runtime.md
TASK-024: 新增 .agent/object-model.md
TASK-025: 新增 .agent/state-machine.md
TASK-026: 新增 .agent/traceability-matrix.md
TASK-027: 新增 .agent/risk-register.md
TASK-028: 新增 .agent/decision-log.md
TASK-029: 新增 .agent/rollback-protocol.md
TASK-030: 新增 patch 模板文件
```

### Milestone M5：Evidence 与评分器

```text
TASK-031: 扩展 release manifest schema
TASK-032: 更新 releasemanifest 生成器
TASK-033: 更新 checksum 校验
TASK-034: 更新 CI artifact summary
TASK-035: 增加 workflow metadata 注入
TASK-036: release-final-check 强制 clean tree + score
TASK-037: 实现 xlibgate score
TASK-038: 新增 docs/scorecard.md
TASK-039: score 写入 manifest
TASK-040: score 小于阈值时 gate 失败
TASK-041: score tests 覆盖维度
```

### Milestone M6：安全与下游矩阵

```text
TASK-042: 强化 secret scanner
TASK-043: 增加 connection string 检测
TASK-044: 增加 /home/k8s/secrets/env 泄漏检测
TASK-045: 更新 security docs
TASK-046: 新增 docs/downstream-matrix.md
TASK-047: integration 增加 configx/redisx/postgresx/kafkax/taosx/ossx/clickhousex case
TASK-048: 每个下游 case 执行 go test ./...
TASK-049: 每个下游 case 执行 make contracts/boundary/evidence
TASK-050: downstream matrix 写入 release manifest
```

### Milestone M7：x.go 边界与最终发布

```text
TASK-051: 新增 docs/xgo-integration-boundary.md
TASK-052: boundary gate 增强 x.go 业务词禁止项
TASK-053: docs 明确 secrets 显式传入原则
TASK-054: 生成最终 release evidence
TASK-055: 完成 review
TASK-056: 完成 retrospective
TASK-057: 输出 DONE with evidence
```

---

## 12. Harness Gates

### Required Gate

```bash
GOWORK=off make fmt
GOWORK=off make vet
GOWORK=off make lint
GOWORK=off make test
GOWORK=off make race
GOWORK=off make boundary
GOWORK=off make security
GOWORK=off make contracts
GOWORK=off make integration
GOWORK=off make docs-check
CHECK_STATUS=passed GOWORK=off make evidence
RELEASE_EVIDENCE_REQUIRE_PASSED=1 GOWORK=off make release-evidence-check
```

### Extended Gate

```bash
GOWORK=off make property
GOWORK=off make golden
FUZZ_SMOKE_TIME=10s GOWORK=off make fuzz-smoke
GOWORK=off make ci-extended
GOWORK=off make release-check-extended
```

### Final Gate

```bash
GOWORK=off make release-final-check
GOWORK=off make release-preflight VERSION=v0.2.0
go run ./cmd/xlibgate score --min 9.8
```

### Semantic Gate

```text
- README 与 docs 角色定义一致。
- xlib-standard 与 kernel 命名一致。
- 旧名只存在于 migration ADR。
- Goal Runtime v3.1 对象完整。
- x.go 边界明确。
```

### Executable Gate

```text
- go test ./...
- go test -race ./...
- xlibgate integration
- xlibgate evidence
- xlibgate score
```

### Hybrid Gate

```text
- docs-check 校验文档结构与占位符。
- contracts 校验 schema 与代码常量。
- boundary 校验依赖与业务词。
- release-check 校验证据与当前源码一致。
```

---

## 13. Evidence Protocol

### 13.1 Task Evidence

每个 Task 完成必须记录：

```text
EVID-<task-id>-<date>-NNN
- changed files:
- commands:
- result:
- artifact:
- known gaps:
```

### 13.2 Goal Evidence

Goal 完成必须记录：

```text
DONE with evidence:
- scope: GOAL-20260602-001
- gates:
  - GOWORK=off make release-final-check: passed
  - GOWORK=off make release-preflight VERSION=v0.2.0: passed
  - go run ./cmd/xlibgate score --min 9.8: passed
- artifacts:
  - release/manifest/latest.json
  - release/manifest/latest.json.sha256
  - docs/scorecard.md
  - docs/downstream-matrix.md
  - .agent/traceability-matrix.md
- known gaps:
  - none
```

### 13.3 不允许的 Evidence

```text
- “我看了，应该没问题”
- “README 已更新”
- “CI 应该会过”
- “未运行但预计通过”
- “跳过安全检查”
- “本地环境没有工具，所以视为通过”
```

---

## 14. Definition of Done

### 14.1 Task DoD

```text
- 对应文件已修改。
- 对应测试已新增或更新。
- 相关 gate 已通过。
- Evidence 已记录。
```

### 14.2 Issue DoD

```text
- Issue 下所有 Task 完成。
- Traceability Matrix 已更新。
- Review 通过。
- 无 P0/P1 known gaps。
```

### 14.3 Goal DoD

```text
- 所有 REQ 完成。
- release-final-check 通过。
- score >= 9.8。
- kernel integration 通过。
- Evidence manifest 生成。
- Retrospective 输出 Prompt/Harness/Rule Patch。
```

### 14.4 Release DoD

```text
- tag version 与 manifest version 一致。
- commit/tree/source digest 一致。
- contracts fingerprint 一致。
- dependency inventory 一致。
- CI artifact 上传成功。
- release-final-check 通过。
```

### 14.5 Retrospective DoD

```text
- 记录本次失败点。
- 记录新增 gate 建议。
- 记录 Prompt Patch。
- 记录 Harness Patch。
- 记录 Rule Patch。
- 生成下一轮 Issue candidates。
```

---

## 15. Risk Register

| Risk ID | 风险 | 等级 | 触发条件 | 缓解方案 |
|---|---|---:|---|---|
| RISK-001 | 改 module path 导致 imports 大面积失败 | P0 | `go test ./...` 失败 | 全局替换 + gofmt + go mod tidy |
| RISK-002 | 旧名清理过度，migration 文档也被误删 | P1 | docs/migration 缺上下文 | stale-name gate 支持 allowlist |
| RISK-003 | xlibgate 一次性 Go 化过大 | P1 | 超过 1 天无法完成 | 先包装现有脚本，再逐步迁移 |
| RISK-004 | kernel 生成后不符合真实 kernel 需求 | P1 | kernel 后续实现冲突 | kernel 专用 profile ADR |
| RISK-005 | score 变成形式主义 | P1 | 只检查文件存在 | score 必须绑定 executable evidence |
| RISK-006 | security scanner 误报过多 | P2 | CI 阻塞 | allowlist 必须显式、可审计 |
| RISK-007 | 下游矩阵过大拖慢 CI | P2 | integration 时间过长 | MVA 覆盖 kernel/corekit，Full nightly 覆盖全矩阵 |
| RISK-008 | Goal Runtime 文档膨胀但不可执行 | P1 | `.agent` 只写概念 | docs-check 校验对象与 traceability |

---

## 16. Rollback Protocol

### 16.1 命名迁移回滚

如果 module path 迁移失败：

```bash
git revert <migration_commit>
```

然后进入：

```text
NEEDS_REPLAN
```

重新拆分：

```text
Phase 1: docs rename
Phase 2: go.mod/import rename
Phase 3: integration rename
```

### 16.2 xlibgate 回滚

如果 Go 化 gate 引发大面积失败：

```text
保留 scripts/*.sh 为主入口
xlibgate 作为 experimental
Makefile 回滚到 scripts
新增 Issue：逐个迁移 gate
```

### 16.3 Evidence schema 回滚

如果 manifest schema 导致 release blocked：

```text
保留旧 schema v1
新增 manifest schema_version
v2 字段先 optional
下一版本再强制 required
```

---

## 17. Human Approval Gates

以下变更必须人工确认：

```text
- 是否最终废弃 baselib-template 名称。
- 是否允许 xlib-standard 同时包含标准与实现。
- 是否将 foundationx 全面迁移为 kernel。
- 是否把 score >= 9.8 作为 release-final-check 强制条件。
- 是否在 CI 中执行完整 downstream matrix。
```

默认本方案裁决：

```text
已废弃 baselib-template 主口径。
xlib-standard = 标准 + 模板 + gate + evidence。
foundationx → kernel。
release-final-check 要求 score >= 9.5；10 分目标要求 >= 9.8。
CI required 覆盖 kernel；nightly/full 覆盖全 downstream matrix。
```

---

## 18. Failure Budget

```text
P0 failure budget: 0
P1 failure budget: 0 for release-final-check
P2 failure budget: <= 3 warnings allowed
P3 docs polish: 可进入 retrospective backlog
```

P0 包括：

```text
go test failure
race failure
x.go dependency leak
secret leak
release manifest mismatch
module path inconsistency
kernel integration failure
score < 9.5
```

---

## 19. 最小可行行动 MVA

MVA 目标：用最小改动把分数从 7.1 拉到 8.5+。

### MVA 必做

```text
1. 新增 ADR：xlib-standard 角色裁决。
2. go.mod 改为 github.com/ZoneCNH/xlib-standard。
3. 全局替换 baselib-template 主口径。
4. 全局替换 foundationx 主口径为 kernel。
5. integration 默认渲染 kernel。
6. docs/module-boundary 解除自相矛盾。
7. contracts tests import 新 module path。
8. 新增 stale-name gate。
9. 运行 GOWORK=off make release-check。
10. 输出 DONE with evidence。
```

### MVA 不做

```text
- 不一次性完成所有 xlibgate Go 化。
- 不一次性覆盖所有下游基础库。
- 不一次性重写全部 .agent 文档。
```

---

## 20. 1 天行动计划

目标：修复 P0 身份矛盾。

```text
Day 1 / Block 1:
  - 新增 ADR-20260602-001-xlib-standard-role.md
  - 新增 ADR-20260602-002-kernel-rename.md

Day 1 / Block 2:
  - 修改 go.mod
  - 替换 imports
  - go mod tidy

Day 1 / Block 3:
  - 修改 README/AGENTS/docs/standard
  - 替换 generation 示例为 kernel

Day 1 / Block 4:
  - 修改 run_integration.sh kernel case
  - 修改 check_rendered_template allowlist
  - 新增 stale-name 检查

Day 1 / Verification:
  - GOWORK=off go test ./...
  - GOWORK=off make contracts
  - GOWORK=off make boundary
  - GOWORK=off make integration
  - GOWORK=off make release-check
```

交付物：

```text
- P0 命名统一完成
- kernel 渲染通过
- release manifest 生成
- 当前评分预计：8.5 / 10
```

---

## 21. 7 天行动计划

目标：把仓库升级为 9+ 标准工厂。

### Day 2：xlibgate 框架

```text
- 创建 cmd/xlibgate
- 实现 CLI skeleton
- docs-check 迁移为 Go
- boundary 迁移为 Go
```

### Day 3：Contracts / Secrets / Render

```text
- contracts 子命令
- secrets 子命令
- render 子命令
- render-check 子命令
```

### Day 4：Evidence v2

```text
- manifest schema v2
- workflow metadata
- artifact metadata
- score placeholder
```

### Day 5：Goal Runtime v3.1

```text
- .agent/object-model.md
- .agent/state-machine.md
- .agent/traceability-matrix.md
- .agent/risk-register.md
- .agent/rollback-protocol.md
```

### Day 6：Scorecard

```text
- xlibgate score
- docs/scorecard.md
- release-check score threshold
```

### Day 7：Review + Release Candidate

```text
- release-check-extended
- release-final-check
- retrospective
- prompt/harness/rule patches
```

交付物：

```text
- 当前评分预计：9.2 ~ 9.5 / 10
```

---

## 22. 30 天行动计划

目标：达到 10 分标准，并形成可复利基础库工厂。

### Week 1：标准源自洽

```text
- 角色裁决
- 命名统一
- kernel 兼容
- xlibgate 初版
- Evidence v2
```

### Week 2：下游矩阵

```text
- configx profile
- observex profile
- redisx profile
- postgresx profile
- kafkax profile
- taosx profile
- ossx profile
- clickhousex profile
```

### Week 3：质量与安全

```text
- gitleaks / enhanced secret scan
- manifest schema validation
- downstream compatibility scoring
- benchmark smoke
- fuzz smoke 扩展
- property tests 扩展
```

### Week 4：复利工程

```text
- Prompt Patch
- Harness Patch
- Rule Patch
- New Issue Candidates
- x.go integration examples
- kernel bootstrap validation
- release v0.2.0 / v1.0.0-rc
```

最终交付物：

```text
- xlib-standard 10分版本
- kernel 可生成/可校验版本
- downstream matrix
- xlibgate CLI
- release manifest v2
- Goal Runtime v3.1 工程化
- score >= 9.8
```

---

## 23. 衡量指标

### 23.1 工程指标

```text
go_test_pass_rate = 100%
race_pass_rate = 100%
contract_gate_pass_rate = 100%
boundary_gate_pass_rate = 100%
security_gate_pass_rate = 100%
release_check_pass_rate = 100%
kernel_render_success = true
score >= 9.8
```

### 23.2 一致性指标

```text
old_name_leak_count = 0
module_path_mismatch_count = 0
foundationx_main_usage_count = 0
xgo_dependency_count = 0
business_term_leak_count = 0
```

### 23.3 Evidence 指标

```text
manifest_exists = true
manifest_checksum_valid = true
source_digest_match = true
contract_fingerprint_match = true
dependency_inventory_match = true
workflow_artifact_uploaded = true
tree_state = clean
```

### 23.4 复利指标

```text
new_library_bootstrap_time <= 10 minutes
new_library_required_gate_pass_time <= 30 minutes
downstream_profile_reuse_rate >= 80%
manual_copy_paste_steps <= 2
```

---

## 24. 迭代优化机制

### 24.1 每次失败都要产生补丁

```text
Test failure → Harness Patch
Docs ambiguity → Prompt Patch
Boundary leak → Rule Patch
Release failure → CI Gate Suggestion
Downstream incompatibility → Template Patch
```

### 24.2 Retrospective 输出格式

```text
RETRO-YYYYMMDD-NNN
- What failed:
- Why it failed:
- Detection gap:
- Prevention patch:
- Prompt Patch:
- Harness Patch:
- Rule Patch:
- New Issue Candidates:
```

### 24.3 自动研究 AutoResearch 触发条件

```text
- Go / golangci-lint / govulncheck 行为不确定。
- GitHub Actions artifact API 行为变化。
- gitleaks 配置不确定。
- 下游生成库失败原因不明确。
- kernel 与 xlib-standard 设计冲突。
- x.go 边界要求变化。
```

### 24.4 Self-improving 机制

每次 release 后更新：

```text
.agent/prompt-patches.md
.agent/harness-patches.md
.agent/rule-patches.md
docs/scorecard.md
docs/downstream-matrix.md
CHANGELOG.md
```

---

## 25. 最终推荐路径

推荐采用：

```text
路径 B：xlib-standard = 标准源 + Go 模板实现 + Generator + Harness + Evidence Runtime
```

不推荐拆分为 `xlib-standard` + `xlib-template` 两个仓库，因为当前仓库已经有完整模板实现资产，拆分会增加迁移成本和同步成本。

最终路径：

```text
Phase 1:
  先把 xlib-standard 身份统一，清除 baselib-template/foundationx 主口径。

Phase 2:
  让 kernel 成为默认下游生成样板。

Phase 3:
  把核心 gate Go 化为 xlibgate。

Phase 4:
  强化 Evidence manifest 与 scorecard。

Phase 5:
  扩展 downstream matrix，形成基础库工厂。

Phase 6:
  用 retrospective 让标准自我增强。
```

---

## 26. 可执行 Agent Prompt

将以下 Prompt 交给执行 Agent：

```text
你是 xlib-standard 10分标准升级执行 Agent。

目标：
将 https://github.com/ZoneCNH/xlib-standard 从当前 7.1/10 升级到 10/10 标准基础库工厂。

执行标准：
Goal Runtime Prompt v3.1，Full 模式。

必须完成：
1. 裁决 xlib-standard 角色为 Standard Source + Go Reference Template + Generator + Harness/Evidence Runtime。
2. 将 baselib-template 主口径统一迁移为 xlib-standard。
3. 将 foundationx 主口径统一迁移为 kernel。
4. 修改 go.mod module path 为 github.com/ZoneCNH/xlib-standard。
5. 修复所有 imports、docs、scripts、contracts、CI、.agent 中的旧命名。
6. integration 默认渲染 github.com/ZoneCNH/kernel。
7. 新增或更新 ADR、migration docs、module boundary、repository roles。
8. 新增 xlibgate CLI，至少提供 docs-check/boundary/secrets/contracts/render-check/evidence/score 的可执行入口；允许第一阶段包装旧脚本，但必须形成 Go 化迁移计划。
9. 升级 .agent 为 Goal Runtime v3.1 对象模型、状态机、Traceability、Risk、Decision、Rollback、Retrospective。
10. 升级 release manifest，记录 score、workflow artifact、checksum、contract fingerprint、dependency inventory。
11. 新增 scorecard，release-final-check 要求 score >= 9.5；10分目标要求 >= 9.8。
12. 确保基础库不得依赖 x.go，不得泄漏 /home/k8s/secrets/env/*，不得包含业务模型。

验证命令：
- GOWORK=off go test ./...
- GOWORK=off go test -race ./...
- GOWORK=off make contracts
- GOWORK=off make boundary
- GOWORK=off make security
- GOWORK=off make integration
- GOWORK=off make docs-check
- CHECK_STATUS=passed GOWORK=off make evidence
- RELEASE_EVIDENCE_REQUIRE_PASSED=1 GOWORK=off make release-evidence-check
- GOWORK=off make release-final-check
- go run ./cmd/xlibgate score --min 9.8

完成声明必须使用：
DONE with evidence:
- scope: GOAL-20260602-001
- gates:
- artifacts:
- known gaps:
```

---

## 27. 10 分完成检查表

```text
[ ] go.mod module path = github.com/ZoneCNH/xlib-standard
[ ] README 主标题 = xlib-standard
[ ] baselib-template 只出现在 migration/ADR 历史说明中
[ ] foundationx 只出现在 migration/ADR 历史说明中
[ ] kernel 是默认生成示例
[ ] docs/standard/module-boundary 与实际仓库角色一致
[ ] docs/standard/repository-roles 与 kernel 命名一致
[ ] cmd/xlibgate 存在
[ ] Makefile 调用 xlibgate 或兼容包装
[ ] contracts 与代码常量一致
[ ] release manifest v2 可生成
[ ] release manifest checksum 可验证
[ ] CI 上传 release evidence artifact
[ ] .agent 包含 v3.1 对象模型与状态机
[ ] traceability matrix 覆盖所有 REQ
[ ] risk register 存在
[ ] rollback protocol 存在
[ ] scorecard 存在
[ ] score >= 9.8
[ ] kernel integration pass
[ ] x.go dependency count = 0
[ ] business term leak count = 0
[ ] secret leak count = 0
[ ] DONE with evidence 输出
```

---

## 28. 最终结论

当前 `xlib-standard` 不是失败项目，而是一个已经具备 70% 工程基础的标准模板工厂雏形。

它要达到 10 分，关键不是继续堆文档，而是完成以下五件事：

```text
1. 角色唯一化。
2. 命名彻底统一。
3. kernel 下游样板化。
4. Gate Go 化与 Evidence 强化。
5. Goal Runtime v3.1 真正工程化。
```

最终推荐路径：

```text
先用 1 天完成 P0 角色/命名/kernel 修复，把评分拉到 8.5；
再用 7 天完成 xlibgate + Evidence v2 + Goal Runtime 工程化，把评分拉到 9.5；
最后用 30 天扩展 downstream matrix + security + scorecard + retrospective，冲 10 分。
```

完成后的 `xlib-standard` 应该成为：

```text
x.go 基础库体系的唯一标准事实源；
kernel 和所有 L1/L2 基础库的生成源；
Agent Teams 的可执行 Goal Runtime；
CI/Harness/Evidence 的统一裁决器；
长期自我改进的复利工程资产。
```
