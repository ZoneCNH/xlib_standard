# goalkit v0.1.0 — 综合设计提案（已拆分）

> 仓库：`github.com/ZoneCNH/xlib-standard`
> 版本：goalkit v0.1.0 / Goal Runtime v3.1.1
> 状态：综合设计提案；不得作为执行终稿。权威规范、ADR、路线图与迁移索引已拆分到 `split/docs/`。
> 生成时间：2026-06-03
> 来源：综合 `.worktree/goal/` 38 个文件深度分析后重构；逐文件 checksum、重复组与处置证据见 `split/docs/plans/goalkit-v0.1.0-migration-index.md`。

---

## 0. 本文件的目的

本文件是 `.worktree/goal/` 目录中来源材料的**综合设计提案**，不是权威替代文件，也不作为执行入口。

权威层级与职责分离如下：

- `split/docs/standard/goalkit-runtime.md`：规范性标准与运行时契约。
- `split/docs/adr/ADR-20260603-001-goalkit-xlibgate-runtime.md`：架构决策与被拒绝方案。
- `split/docs/plans/goalkit-v0.1.0-roadmap.md`：PR 顺序、依赖、回滚与验收命令。
- `split/docs/plans/goalkit-v0.1.0-migration-index.md`：`goal/` 来源文件到拆分工件的逐文件迁移证据。

执行者必须以拆分后的 `split/docs/` 工件为审计入口；本文件仅保留为历史上下文与综合提案。

它解决了原始文件群的以下问题：

| 问题            | 原状                                       | 本文修正                 |
| --------------- | ------------------------------------------ | ------------------------ |
| 完全重复        | goalkit.md 与 goalkit_v0_1_0...md checksum 相同 | 在迁移索引中逐项列出并保留审计证据 |
| 局部重复        | 自动化内容在 3 个文件中重复                | 合并为第 10 节           |
| 架构漂移        | 4 个架构文档迭代未收敛                     | 合并为第 2-5 节          |
| 脱离仓库        | 不引用已有 .agent/、xlibgate               | 基于仓库实际状态         |
| PR 过多         | 28 个 PR，后 20 个极薄                     | 收敛为 12 个 PR          |
| DONE 膨胀       | Full Mode DONE 公式 37 项                  | 按 Mode 分级，最多 11 项 |
| 与 debt.md 冲突 | goalkit 与 debt 方案重叠                   | 明确分工边界             |

---

## 1. 最高结论

```
Goal Runtime v3.1.1 的目标能力边界已经收敛为提案。
下一阶段不是把本文件当作终稿执行，而是按拆分工件逐步落地、验证并审计。
```

goalkit v0.1.0 = 把 Goal Runtime 从设计文档产品化为可执行系统：

```
Goal Kernel        — 定义"做什么"
Harness Runtime    — 定义"怎么裁判"
Evidence Ledger    — 定义"怎么证明"
Extensions         — 按需启用重型流程
Automation Surface — 安全自动化边界
```

---

## 2. 架构：Goal Kernel + Harness Runtime + Extensions

### 2.1 最终分层

```
Goal Kernel           （必选，所有 Mode）
  ├── Goal
  ├── Spec
  ├── Design
  ├── Plan
  ├── Task
  ├── Test
  ├── Evidence
  └── Review

Harness Runtime       （必选，所有 Mode）
  ├── Mode Router
  ├── Gate Registry           ← .agent/harness.yaml
  ├── Command Registry        ← .agent/command-registry.yaml
  ├── Blocking Policy
  ├── Evidence Policy
  ├── Artifact Policy
  └── Failure Budget

Completion Extension  （Full Mode 或 release_verify）
  ├── Acceptance
  ├── Delivery
  ├── Handover
  └── Completion Certificate

Release Extension     （release_impacting=true）
  ├── Version Plan
  ├── CHANGELOG
  ├── Release Manifest
  ├── Release Preflight
  └── Rollback

Ecosystem Extension   （downstream_impacting=true）
  ├── Downstream Impact Matrix
  ├── Adoption Decision
  └── x.go Consumer Boundary

Governance Extension  （Runtime/Harness/Policy 变更）
  ├── Runtime-as-Code
  ├── Policy-as-Code
  ├── Trust Root
  ├── Drift Check
  └── Budget / Anti-bloat

Automation Surface    （最后启用，分阶段）
  ├── Issue Sync
  ├── PR Sync
  ├── Auto Commit Guard
  ├── Version Plan
  ├── Release Draft
  └── Release Publish
```

#### Downstream 库清单（Ecosystem Extension 检查范围）

| 层级 | 库名                                                        | 说明               |
| ---- | ----------------------------------------------------------- | ------------------ |
| L0   | xlib-standard（kernel）                                     | 基础标准库本体     |
| L1   | configx, observex, testkitx                                 | 核心基础设施扩展   |
| L2   | postgresx, redisx, kafkax, taosx, ossx, clickhousex         | 数据/存储/消息扩展 |

`xlibgate downstream` 和 `make goal-downstream-adoption` 的检查范围覆盖 L1 + L2 全部库。L0 是自身，不需要 downstream 检查。

### 2.2 核心边界表

| 组件            | 职责                                | 不负责                   |
| --------------- | ----------------------------------- | ------------------------ |
| Goal Kernel     | 目标、规格、设计、计划、任务、证据  | 发布、认证、自动化       |
| Harness Runtime | Mode 路由、Gate 编排、Evidence 强制 | 执行检查命令本身         |
| xlibgate        | 执行 Gate 检查命令                  | 定义规则                 |
| Makefile        | 人类/CI 入口                        | 复杂决策逻辑             |
| Evidence Ledger | 机器可验证的事实源                  | 人类可读报告（那是视图） |
| Policy-as-Code  | 定义规则                            | 执行调度                 |

### 2.3 与已有仓库结构的映射

仓库已有 `.agent/` 目录包含 30+ 文件。goalkit v0.1.0 **不重建**，而是基于已有结构扩展：

| 已有文件                               | goalkit 中的角色                |
| -------------------------------------- | ------------------------------- |
| `.agent/harness.yaml`                  | Harness Runtime 配置主文件      |
| `.agent/command-registry.yaml`         | Command Registry 事实源         |
| `.agent/gates.md`                      | Gate 文档（保留，补充 G12-G16） |
| `.agent/evidence-protocol.md`          | Evidence Ledger 规范基础        |
| `.agent/evidence-artifact-policy.yaml` | Artifact Policy 事实源          |
| `cmd/xlibgate/main.go`                 | CLI 执行器入口                  |
| `.agent/command-implementation-status.yaml` | 命令实现状态生命周期事实源      |
| `cmd/xlibgate/governance.go`           | 已有治理命令实现模式            |

新增文件（PR 中创建）：

```
.agent/acceptance/           ← PR-1
.agent/delivery/             ← PR-1
.agent/handover/             ← PR-1
.agent/downstream/           ← PR-1（不覆盖已有 downstream-*）
.agent/certification/        ← PR-1
.agent/schemas/              ← PR-2
internal/goalruntime/        ← PR-4
testdata/goal-runtime-v3-1-1/ ← PR-4
```

---

## 3. Goal 对象模型

```yaml
goal:
  id: GOAL-YYYYMMDD-NNN
  title: ""
  intent: ""
  mode: LITE | STANDARD | FULL
  change_level: L0 | L1 | L2 | L3 | L4 | L5
  release_impacting: false
  downstream_impacting: false

spec:
  requirements: []
  acceptance_criteria: []
  out_of_scope: []

execution:
  design: ""
  plan: ""
  tasks: []
  tests: []

evidence:
  ledger: []
  required: []

extensions:
  completion: { enabled: false }
  release: { enabled: false }
  downstream: { enabled: false }
  governance: { enabled: false }
  automation: { enabled: false }
```

关键开关：`mode`、`release_impacting`、`downstream_impacting` 决定启用哪些 Extension。

---

## 4. Mode 路由与 DONE 公式

### 4.1 Lite Mode

适用：文档修正、注释修复、链接修复。

DONE 公式（4 项）：

```
DONE = Task Complete
     + docs-check PASS
     + Evidence Summary
     + Review
```

### 4.2 Standard Mode

适用：模板补充、docs/standard 新文档、.agent 模板调整。

DONE 公式（6 项）：

```
DONE = Spec PASS
     + Task Complete
     + Test PASS
     + Evidence
     + Review
     + Optional Acceptance
```

### 4.3 Full Mode

适用：Makefile、harness.yaml、xlibgate、release manifest、CI、public API、security。

DONE 公式（11 项）：

```
DONE = Spec PASS
     + Design PASS
     + Plan PASS
     + Task PASS
     + Test PASS
     + Evidence PASS
     + Review PASS
     + Acceptance PASS
     + Delivery PASS
     + Handover PASS
     + Completion Certificate
```

### 4.4 Release DONE（叠加）

仅当 `release_impacting=true` 时在 Full DONE 基础上追加：

```
     + Version Plan
     + Release Manifest
     + Preflight PASS
     + Rollback Plan
```

### 4.5 Mode 路由规则

```
变更文件包含            → Mode
─────────────────────────────────
docs/ 纯文档修正        → LITE
.agent/ 模板、docs/standard → STANDARD
Makefile / harness.yaml → FULL
cmd/xlibgate/           → FULL
release/manifest/       → FULL
.github/workflows/      → FULL
public API 变更         → FULL
```

冲突解决：当变更跨多个 Mode 区间时，取最高 Mode。例如同时修改 `docs/` 和 `Makefile`，取 FULL。

---

## 5. Gate 收敛模型

不再使用无限编号。收敛为五类：

### Intent Gates（防目标不清）

```
Context Gate / Goal Gate / Spec Gate / Design Gate / Plan Gate
```

### Execution Gates（防没执行）

```
Task Gate / Implementation Gate / Test Gate / Evidence Gate / Review Gate
```

### Completion Gates — G12-G16（Full Mode，本次新增）

```
G12 Acceptance Gate     → make goal-acceptance
G13 Delivery Gate       → make goal-delivery
G14 Handover Gate       → make goal-handover
G15 Downstream Gate     → make goal-downstream-adoption
G16 Completion Gate     → make goal-certify
```

### Release Gates（release_impacting=true）

```
Version Gate / Manifest Gate / Preflight Gate / Trust Root Gate / Rollback Gate
```

### Evolution Gates（长期治理）

```
Retrospective Gate / Policy Gate / Runtime Drift Gate / Budget Gate
```

### Gate 失败处理策略

Gate 返回非零 exit code 时的处理决策树：

| Exit Code | 名称                | 处理方式                             |
| --------- | ------------------- | ------------------------------------ |
| 1         | FAIL                | 修复问题后重跑 Gate                  |
| 2         | BLOCKED             | 解除外部依赖后重跑                   |
| 3         | INVALID_INPUT       | 修正参数后重跑                       |
| 4         | INCONSISTENT_STATE  | 修复 Registry/Artifact 状态后重跑    |
| 5-9       | MISSING_*           | 补齐缺失物（Evidence/Artifact/Manifest/Certificate）后重跑 |

原则：Gate 失败**不允许绕过**，只允许修复后重试或降级 Mode（降级需记录 Evidence）。

---

## 6. Evidence Ledger

Evidence 不是 Markdown 文件堆，而是结构化事实记录。

```yaml
evidence:
  id: EVID-TASK-001-20260603-001
  goal_id: GOAL-20260603-001
  task_id: TASK-001
  gate_id: G12_ACCEPTANCE
  commit: abc123
  command: "make goal-acceptance GOAL_ID=..."
  exit_code: 0
  output_artifact: ".agent/acceptance/acceptance_report.md"
  checksum: sha256:...
  actor: agent | human
  generated_at: "2026-06-03T12:00:00Z"
  supports:
    - AC-001
    - TEST-001
```

规则：

- Markdown report 是**视图**
- Evidence Ledger 是**事实源**
- 生成的 report 文件**不提交**源码历史

---

## 7. 与 debt.md 的分工

`.worktree/debt.md`（3766 行）和 goalkit 方案有重叠。必须明确分工：

| 维度          | goalkit v0.1.0                                   | debt.md                                    |
| ------------- | ------------------------------------------------ | ------------------------------------------ |
| 定位          | Goal Runtime 框架与工具包                        | 债务治理的自动化落地                       |
| 覆盖          | Goal Kernel、Harness、Gate、Evidence、Extensions | debt-evidence、debt gates、downstream scan |
| xlibgate 命令 | acceptance/delivery/handover/downstream/certify  | debt（已实现）、governance（已实现）       |
| Makefile 目标 | goal-\* 系列                                     | debt-evidence、debt、score 等（已存在）    |
| 优先级        | 新增 G12-G16 framework                           | 利用已有 framework 做 debt 治理            |
| 执行顺序      | 先建 Goal Runtime 框架                           | 在框架上运行 debt goals                    |

原则：**goalkit 建框架，debt 用框架**。两者不应互相阻塞。

---

## 8. PR 执行路线（从 28 个收敛为 12 个）

原始方案有 28 个 PR，其中 PR-9~PR-27 每个仅 230 行、内容极薄。重新合并为 12 个可执行 PR：

### Phase 1：Core MVA（7 天）

| PR   | 名称                                 | 范围                                                                                                                         | 产出                          |
| ---- | ------------------------------------ | ---------------------------------------------------------------------------------------------------------------------------- | ----------------------------- |
| PR-1 | Templates + Docs                     | .agent/acceptance,delivery,handover,downstream,certification,signoff/ + docs/standard/                                       | 静态模板与文档                |
| PR-2 | Schemas + Output Contract            | .agent/schemas/ gate_result.schema.json, evidence.schema.json                                                                | JSON Schema 事实源            |
| PR-3 | Runtime Index + ADR                  | .agent/ runtime index + docs/adr/                                                                                            | 兼容性声明、版本索引          |
| PR-4 | Harness + xlibgate 实现              | Makefile goal-\* targets + cmd/xlibgate/ acceptance,delivery,handover,downstream,certify + internal/goalruntime/ + testdata/ | 可执行命令 + fixtures + tests |
| PR-5 | Generated Artifact Policy + Blocking | .agent/evidence-artifact-policy.yaml 更新 + harness.yaml blocking activation                                                 | Full Mode G12-G16 阻断生效    |

MVA 验收命令：

```bash
GOAL_ID=GOAL-20260603-XLIB-RUNTIME-001 \
GOAL_RUNTIME_MODE=FULL \
GOWORK=off make goal-runtime-final
```

成功标准：

```
acceptance PASS + delivery PASS + handover PASS
+ downstream decision recorded
+ completion certificate generated (not committed)
+ DONE with evidence
```

### Phase 2：可信治理（30 天）

| PR   | 名称                             | 范围                                                  |
| ---- | -------------------------------- | ----------------------------------------------------- |
| PR-6 | Freeze + Drift Control           | runtime freeze certificate + drift-check              |
| PR-7 | Runtime-as-Code + Policy-as-Code | runtime.v3.1.1.yaml + policy rules YAML               |
| PR-8 | Test Harness + Trust Root        | golden/negative fixtures for Gates + trust root check |

### Phase 3：生态与协作（60 天）

| PR    | 名称                             | 范围                                                   |
| ----- | -------------------------------- | ------------------------------------------------------ |
| PR-9  | Downstream Adoption Orchestrator | downstream sync matrix automation + adoption dashboard |
| PR-10 | Runtime Observability + Budget   | SLO metrics + anti-bloat check + budget gate           |

### Phase 4：成熟化（90 天）

| PR    | 名称                          | 范围                                                                                            |
| ----- | ----------------------------- | ----------------------------------------------------------------------------------------------- |
| PR-11 | DX + Conformance + Publishing | onboarding guide + conformance benchmark + standard publishing                                  |
| PR-12 | Issue/PR/Release Automation   | issue-sync, pr-sync, auto-commit, version-plan, release-draft, release-publish（分 6 个 Stage） |

### 原始 PR 合并映射

| 新 PR | 合并自原始 PR                                       |
| ----- | --------------------------------------------------- |
| PR-1  | 原 PR-1                                             |
| PR-2  | 原 PR-2                                             |
| PR-3  | 原 PR-3                                             |
| PR-4  | 原 PR-4 + PR-5（Makefile + xlibgate 一起做更高效）  |
| PR-5  | 原 PR-6 + PR-8（artifact policy + blocking 一起做） |
| PR-6  | 原 PR-9                                             |
| PR-7  | 原 PR-13 + PR-18                                    |
| PR-8  | 原 PR-19 + PR-20                                    |
| PR-9  | 原 PR-14 + PR-15                                    |
| PR-10 | 原 PR-21 + PR-22                                    |
| PR-11 | 原 PR-23 + PR-24 + PR-25 + PR-26 + PR-27            |
| PR-12 | 原 PR-28（保持不变，但最后做）                      |

跳过的原始 PR：

- 原 PR-7（CI artifacts scorecard）：与已有 `make score` 重叠
- 原 PR-10（pack generator）：过早优化
- 原 PR-11（runtime pack CI）：过早优化
- 原 PR-12（runtime lifecycle）：与 PR-6 合并
- 原 PR-16（agent team orchestration）：非 goalkit 核心
- 原 PR-17（self-improving）：非 goalkit 核心

### 执行包引用表

本节仅保留历史归并关系。执行者不得使用缩写文件名或本表作为来源处置依据；逐文件精确路径、SHA-256、重复组、迁移目标与处置状态以 `split/docs/plans/goalkit-v0.1.0-migration-index.md` 为准。

| 新 PR | 原始执行包范围 |
| ----- | -------------- |
| PR-1  | PR-1 templates/docs 执行包 |
| PR-2  | PR-2 schemas/output-contract 执行包 |
| PR-3  | PR-3 runtime-index/compatibility/ADR 执行包 |
| PR-4  | PR-4 Makefile/harness/command-registry 与 PR-5 xlibgate/fixtures/tests 执行包 |
| PR-5  | PR-6 generated-artifact-policy 与 PR-8 blocking-activation 执行包 |
| PR-6  | PR-9 freeze/drift-control 执行包 |
| PR-7  | PR-13 runtime-as-code 与 PR-18 policy-as-code 执行包 |
| PR-8  | PR-12 operations lifecycle 与 PR-19 harness/golden-fixtures 执行包 |
| PR-9  | PR-14 downstream orchestrator 与 PR-15 ecosystem/dashboard 执行包 |
| PR-10 | PR-20/21/22 trust-root、observability、budget 执行包 |
| PR-11 | PR-16/17/23/24/25/26/27 orchestration、DX、versioning、conformance、publishing、constitution 执行包 |
| PR-12 | xlib_standard_issue_pr_commit_release_automation_patch.md                         |

> 文件名中 `...` 表示省略的中间路径段。执行时在 `.worktree/goal/` 下按前缀匹配即可。

---

## 9. xlibgate 命令契约

### 9.1 新增命令

```
xlibgate acceptance --goal <id> [--mode LITE|STANDARD|FULL] [--format text|json] [--dry-run]
xlibgate delivery   --goal <id> [--mode ...] [--format ...]
xlibgate handover   --goal <id> [--mode ...] [--format ...]
xlibgate downstream --goal <id> [--mode ...] [--format ...]
xlibgate certify    --goal <id> [--mode ...] [--format ...]
```

### 9.2 统一参数

| 参数           | 必填 | 默认     | 说明           |
| -------------- | ---- | -------- | -------------- |
| --goal         | 是   | -        | Goal ID        |
| --mode         | 否   | STANDARD | 执行模式       |
| --format       | 否   | text     | 输出格式       |
| --artifact-dir | 否   | .agent   | 生成物根目录   |
| --dry-run      | 否   | false    | 不写生成物     |
| --verify       | 否   | false    | 只验证，不生成 |
| --strict       | 否   | false    | 严格模式       |

### 9.3 Exit Code 契约

| Code | 名称                | 含义                      |
| ---- | ------------------- | ------------------------- |
| 0    | PASS                | Gate 通过                 |
| 1    | FAIL                | 通用失败                  |
| 2    | BLOCKED             | 外部依赖阻塞              |
| 3    | INVALID_INPUT       | 参数缺失                  |
| 4    | INCONSISTENT_STATE  | 状态冲突                  |
| 5    | MISSING_EVIDENCE    | 缺 Evidence               |
| 6    | NON_ACCEPTANCE      | 触发不接受条件            |
| 7    | MISSING_ARTIFACT    | 缺交付物                  |
| 8    | MISSING_MANIFEST    | 缺 Release Manifest       |
| 9    | MISSING_CERTIFICATE | 缺 Completion Certificate |

### 9.4 Gate Result Envelope（JSON）

```json
{
  "schema_version": "3.1.1",
  "goal_id": "",
  "gate_id": "G12_ACCEPTANCE",
  "command": "acceptance",
  "status": "pass | pass_with_risk | fail | blocked | skipped",
  "mode": "FULL",
  "report_path": ".agent/acceptance/acceptance_report.md",
  "generated_artifacts": [],
  "evidence": [],
  "risks": [],
  "errors": [],
  "warnings": [],
  "duration_ms": 0,
  "generated_at": ""
}
```

### 9.5 Go 实现结构

基于仓库已有 `cmd/xlibgate/` 风格：

```
cmd/xlibgate/
  main.go              ← 已有，注册新命令
  governance.go        ← 已有，治理命令模式参考
  acceptance.go        ← 新增
  delivery.go          ← 新增
  handover.go          ← 新增
  downstream.go        ← 新增
  certify.go           ← 新增

internal/goalruntime/
  result.go            ← GateResult, GateStatus, ExitCode
  options.go           ← 统一参数解析
  report.go            ← Report 生成
  acceptance.go        ← 业务逻辑
  delivery.go
  handover.go
  downstream.go
  certify.go
  *_test.go

testdata/goal-runtime-v3-1-1/
  acceptance/pass/
  acceptance/missing-evidence/
  delivery/pass/
  delivery/missing-artifact/
  handover/pass/
  downstream/pass/
  downstream/xgo-boundary-fail/
  certify/pass/
  certify/missing-acceptance/
```

### 9.6 Makefile 目标

```makefile
GOAL_RUNTIME_MODE ?= STANDARD

.PHONY: require-goal-id
require-goal-id:
	@if [ -z "$(GOAL_ID)" ]; then echo "GOAL_ID is required"; exit 3; fi

.PHONY: goal-acceptance
goal-acceptance: require-goal-id
	$(XLIBGATE) acceptance --goal $(GOAL_ID) --mode $(GOAL_RUNTIME_MODE)

.PHONY: goal-delivery
goal-delivery: require-goal-id
	$(XLIBGATE) delivery --goal $(GOAL_ID) --mode $(GOAL_RUNTIME_MODE)

.PHONY: goal-handover
goal-handover: require-goal-id
	$(XLIBGATE) handover --goal $(GOAL_ID) --mode $(GOAL_RUNTIME_MODE)

.PHONY: goal-downstream-adoption
goal-downstream-adoption: require-goal-id
	$(XLIBGATE) downstream --goal $(GOAL_ID) --mode $(GOAL_RUNTIME_MODE)

.PHONY: goal-certify
goal-certify: require-goal-id
	$(XLIBGATE) certify --goal $(GOAL_ID) --mode $(GOAL_RUNTIME_MODE)

.PHONY: goal-runtime-final
goal-runtime-final: require-goal-id
	GOWORK=off $(MAKE) docs-check
	GOWORK=off $(MAKE) governance-check
	GOWORK=off $(MAKE) evidence-check
	GOWORK=off $(MAKE) goal-acceptance GOAL_ID=$(GOAL_ID)
	GOWORK=off $(MAKE) goal-delivery GOAL_ID=$(GOAL_ID)
	GOWORK=off $(MAKE) goal-downstream-adoption GOAL_ID=$(GOAL_ID)
	GOWORK=off $(MAKE) goal-handover GOAL_ID=$(GOAL_ID)
	GOWORK=off $(MAKE) goal-certify GOAL_ID=$(GOAL_ID)
```

关键约束：`goal-downstream-adoption` 不得覆盖已有 `downstream-adoption`。

---

## 10. Automation Surface（PR-12，最后启用）

### 10.1 分阶段启用

| Stage | 能力                                      | 风险 |
| ----- | ----------------------------------------- | ---- |
| 1     | issue-sync, issue-plan, pr-body, pr-sync  | 低   |
| 2     | version-plan, release-plan, release-draft | 低   |
| 3     | auto-commit（仅 bot branch）              | 中   |
| 4     | release-publish（需 human approval）      | 高   |

### 10.2 绝对禁止

```
直接 push main
绕过 branch protection
绕过 Human Approval
Gate 未通过就合并
缺 Evidence 就关闭 Issue
无 Release Manifest 就发布
无 Completion Certificate 就声明 DONE
自动提交 generated artifacts
自动暴露 secrets
```

### 10.3 Issue ↔ Goal 映射

每个可执行 Issue 必须映射到一个 Goal 或 Task。

#### Issue Body 结构化模板

Issue body 必须包含 Goal Runtime 元数据块：

```markdown
## Goal Runtime

Goal ID:
Mode: Lite / Standard / Full
Change Level: L0 / L1 / L2 / L3 / L4 / L5
Runtime Version: v3.1.1

## Scope

In Scope:
-

Out of Scope:
-

## Acceptance

- [ ]

## Evidence Required

- [ ] Gate results
- [ ] Diff summary
- [ ] Release manifest if needed
```

#### Issue Labels 分类体系

```
goal:runtime     goal:lite     goal:standard     goal:full

state:triaged          state:goal-ready       state:spec-ready
state:design-ready     state:plan-ready       state:tasks-ready
state:executing        state:verifying        state:reviewing
state:accepted         state:delivered        state:released
state:done             state:blocked          state:needs-research
state:needs-decision   state:needs-replan     state:needs-rollback

risk:l0   risk:l1   risk:l2   risk:l3   risk:l4   risk:l5

area:agent        area:harness      area:xlibgate     area:release
area:evidence     area:docs         area:downstream   area:runtime-pack
area:policy       area:trust-root
```

#### Issue 状态同步表

| Goal Runtime State | Issue Label          |
| ------------------ | -------------------- |
| INIT               | state:triaged        |
| GOAL_READY         | state:goal-ready     |
| SPEC_READY         | state:spec-ready     |
| DESIGN_READY       | state:design-ready   |
| PLAN_READY         | state:plan-ready     |
| TASKS_READY        | state:tasks-ready    |
| EXECUTING          | state:executing      |
| VERIFYING          | state:verifying      |
| REVIEWING          | state:reviewing      |
| ACCEPTING          | state:accepted       |
| DELIVERING         | state:delivered      |
| RELEASING          | state:released       |
| DONE               | state:done           |
| BLOCKED            | state:blocked        |
| NEEDS_RESEARCH     | state:needs-research |
| NEEDS_DECISION     | state:needs-decision |
| NEEDS_REPLAN       | state:needs-replan   |
| NEEDS_ROLLBACK     | state:needs-rollback |

#### Issue 关闭条件

Issue 只有满足以下全部条件才可关闭：

```
PR merged
Evidence exists
Completion Certificate exists（if Full Mode）
Release Manifest exists（if release_impacting）
Final statement contains DONE with evidence
```

### 10.4 PR 自动化

#### PR Body 完整模板

```markdown
# PR

## Goal

Goal ID:
Issue:
Mode:
Change Level:
Runtime Version:

## Spec / Design / Plan

Spec:
Design:
Plan:

## Tasks

- TASK-

## Evidence

Evidence ID:
Evidence Path:
Gate Results:

## Acceptance / Delivery

Acceptance:
Delivery:
Handover:
Completion Certificate:

## Release

Version Plan:
CHANGELOG:
Release Manifest:
Rollback:

## Checklist

- [ ] Issue linked
- [ ] Goal linked
- [ ] Gates passed
- [ ] Evidence attached
- [ ] Docs updated
- [ ] CHANGELOG updated if needed
- [ ] Generated artifacts not committed
- [ ] Secrets not leaked
```

#### PR 自动更新时机表

| 触发事件                         | 自动动作                                              |
| -------------------------------- | ----------------------------------------------------- |
| Issue triaged                    | 写入 Goal ID 和 Mode                                  |
| Branch created                   | PR draft body 写入 Goal metadata                      |
| Gate completed                   | 更新 Gate Results                                     |
| Evidence generated               | 更新 Evidence 链接                                    |
| Release manifest generated       | 更新 Release Manifest artifact 链接                   |
| Completion Certificate generated | 更新 Completion Certificate 链接                      |
| All gates pass                   | 从 Draft 转为 Ready for review                        |
| PR merged                        | 同步 Issue state:done，前提是 DONE with evidence 成立 |

#### PR 策略

```
默认创建 Draft PR。
Gate 全部通过后转为 Ready for review。
只允许 Squash merge。
禁止直接 merge commit。
禁止绕过 review。
禁止绕过 required checks。
```

### 10.5 Commit 自动化

#### Commit Message 规范

格式：

```
<type>(<scope>): <summary>

Goal: GOAL-YYYYMMDD-NNN
Issue: #<number>
Task: TASK-<goal-id>-NNN
Evidence: EVID-<task-id>-YYYYMMDD-NNN
```

示例：

```
feat(agent): add goal acceptance templates

Goal: GOAL-20260603-002
Issue: #123
Task: TASK-GOAL-20260603-002-001
Evidence: EVID-TASK-GOAL-20260603-002-001-20260603-001
```

#### Auto-commit 分支限制

Bot 只允许提交到以下分支：

```
feat/goal-<goal-id>-<slug>
fix/goal-<goal-id>-<slug>
docs/goal-<goal-id>-<slug>
```

禁止提交到：

```
main
release/*
protected branches
```

#### Auto-commit 前置检查

自动提交前必须通过以下 6 项检查：

```
1. 工作区只包含当前 Task scope 内文件
2. 禁止提交 generated artifacts
3. secrets scan 通过
4. format / docs-check / minimal gate 通过
5. diff summary 已生成
6. commit message 符合规范
```

#### 禁止自动提交的文件

```
release/manifest/latest.json
release/manifest/latest.json.sha256
.agent/acceptance/acceptance_report.md
.agent/delivery/artifact_inventory.yaml
.agent/delivery/delivery_report.md
.agent/handover/handover_report.md
.agent/handover/next_agent_context.md
.agent/downstream/downstream_sync_matrix.md
.agent/downstream/downstream_adoption_report.md
.agent/certification/completion_certificate.md
release/runtime-pack/*.md
release/runtime-pack/*.json
release/runtime-pack/*.sha256
.agent/**/generated/*
```

### 10.6 自动化 xlibgate 命令

以下命令仅在 PR-12 实现：

| 命令                       | 说明                                        | 输出文件                                                             |
| -------------------------- | ------------------------------------------- | -------------------------------------------------------------------- |
| `xlibgate issue-sync`      | 读取 GitHub Issue，生成/更新 Goal 和 Task   | `.agent/registry/goals.yaml`、`.agent/registry/tasks.yaml`           |
| `xlibgate issue-plan`      | 从 Issue 生成执行计划和 PR 拆分建议         | `.agent/issues/issue-<number>-plan.md`                               |
| `xlibgate issue-state`     | 同步 Goal State 到 Issue labels / comment   | —（直接更新 GitHub Issue）                                           |
| `xlibgate pr-body`         | 根据 Goal / Evidence / Release 生成 PR body | `.agent/pr/pr_body.md`                                               |
| `xlibgate pr-sync`         | 同步 PR body、labels、linked issues         | —（直接更新 GitHub PR）                                              |
| `xlibgate auto-commit`     | 在 bot branch 内执行安全自动提交            | —（直接提交到 bot branch）                                           |
| `xlibgate version-plan`    | 根据变更级别生成版本变更计划                | `.agent/release/version_plan.md`、`.agent/release/version_plan.json` |
| `xlibgate release-plan`    | 生成 release plan                           | `.agent/release/release_plan.md`                                     |
| `xlibgate release-draft`   | 生成 GitHub Release draft 内容              | `.agent/release/release_draft.md`                                    |
| `xlibgate release-publish` | 全部 release gates 通过后发布版本           | —（创建 Git tag + GitHub Release）                                   |

命令参数示例：

```
xlibgate issue-sync --issue <number>
xlibgate issue-plan --issue <number>
xlibgate issue-state --issue <number> --state <STATE>
xlibgate pr-body --goal <id> --issue <number>
xlibgate pr-sync --pr <number> --goal <id>
xlibgate auto-commit --goal <id> --task <id>
xlibgate version-plan --goal <id>
xlibgate release-plan --goal <id>
xlibgate release-draft --goal <id> --version <version>
xlibgate release-publish --goal <id> --version <version>
```

### 10.7 Version 规划规则

#### 版本 bump 判定表

| bump  | 触发条件                                                                     |
| ----- | ---------------------------------------------------------------------------- |
| patch | 文档、模板、低风险修复                                                       |
| minor | 新增 Runtime capability、xlibgate command、schema 可选字段、非破坏性能力扩展 |
| major | 破坏性接口、Manifest 主结构变更、Harness schema 不兼容、Generator 不兼容     |

#### Version Plan JSON 结构示例

```json
{
  "goal_id": "GOAL-20260603-002",
  "recommended_bump": "minor",
  "current_version": "v1.2.3",
  "next_version": "v1.3.0",
  "reason": [
    "adds new goal runtime capability",
    "adds new xlibgate commands",
    "adds new schemas"
  ],
  "requires_human_approval": true
}
```

#### Tag 规则

格式：`vMAJOR.MINOR.PATCH`

Tag message 必须包含：

```
Goal ID
Release ID
Manifest SHA256
Completion Certificate
Evidence Artifact
```

### 10.8 发布硬前置

发布前必须通过以下 12 项检查：

```
1.  workflow_dispatch 触发（禁止 PR event 直接触发）
2.  environment approval 通过
3.  version-plan PASS
4.  release-plan PASS
5.  CHANGELOG 已更新
6.  release-final-check PASS
7.  release-preflight PASS
8.  Release Manifest 存在且 checksum 匹配
9.  Completion Certificate 存在（Full Mode）
10. Trust Root PASS
11. Evidence Attestation PASS
12. Git tag 不存在冲突
```

### 10.9 GitHub Actions Workflow 设计

新增两个 workflow 文件：

```
.github/workflows/goal-runtime-sync.yml
.github/workflows/goal-runtime-release.yml
```

#### goal-runtime-sync.yml

**触发器 1：Issue opened / labeled / edited**

```yaml
on:
  issues:
    types: [opened, labeled, edited]
```

动作：issue-sync → issue-plan → comment plan。

**触发器 2：PR opened / synchronize / reopened**

```yaml
on:
  pull_request:
    types: [opened, synchronize, reopened, ready_for_review]
```

动作：pr-sync → docs-check → governance-check → evidence-check。

**触发器 3：PR closed + merged**

```yaml
on:
  pull_request:
    types: [closed]
```

条件：`merged == true`。动作：pr-sync → issue-state done（if evidence complete）→ version-plan（if release-impacting）。

#### goal-runtime-release.yml

仅 `workflow_dispatch` 触发：

```yaml
workflow_dispatch:
  inputs:
    goal_id:
      description: "Goal ID"
      required: true
    version:
      description: "Version to release"
      required: true
```

动作序列：

```
release-plan
→ release-final-check
→ release-preflight
→ release-draft
→ human approval（environment gate）
→ release-publish
```

#### 发布防护规则

```
1. release-publish 只能通过 workflow_dispatch 触发
2. release-publish 必须要求 environment approval
3. release-publish 不允许在 pull_request event 直接运行
4. tag 已存在时必须失败
5. manifest checksum 不匹配必须失败
6. Completion Certificate 缺失必须失败
7. GitHub token 权限最小化
```

### 10.10 Registry 扩展

PR-12 新增 5 个 registry 文件：

#### .agent/registry/issues.yaml

```yaml
issues:
  - issue: 123
    goal_id: GOAL-20260603-002
    state: executing
    mode: FULL
    linked_prs: []
```

#### .agent/registry/prs.yaml

```yaml
prs:
  - pr: 456
    goal_id: GOAL-20260603-002
    state: reviewing
    evidence: []
    release_impacting: true
```

#### .agent/registry/commits.yaml

```yaml
commits:
  - sha: abc123
    goal_id: GOAL-20260603-002
    task_id: TASK-GOAL-20260603-002-001
    branch: feat/goal-GOAL-20260603-002-acceptance
    auto: true
```

#### .agent/registry/versions.yaml

```yaml
versions:
  current: v1.2.3
  planned:
    - goal_id: GOAL-20260603-002
      next: v1.3.0
      bump: minor
      status: planned
```

#### .agent/registry/releases.yaml

```yaml
releases:
  - version: v1.3.0
    goal_id: GOAL-20260603-002
    tag: v1.3.0
    manifest_sha256: sha256:...
    certificate: true
    published_at: ""
```

### 10.11 自动化 Makefile Targets

```makefile
.PHONY: issue-sync
issue-sync:
	$(XLIBGATE) issue-sync --issue $(ISSUE)

.PHONY: issue-plan
issue-plan:
	$(XLIBGATE) issue-plan --issue $(ISSUE)

.PHONY: issue-state
issue-state:
	$(XLIBGATE) issue-state --issue $(ISSUE) --state $(STATE)

.PHONY: pr-body
pr-body:
	$(XLIBGATE) pr-body --goal $(GOAL_ID) --issue $(ISSUE)

.PHONY: pr-sync
pr-sync:
	$(XLIBGATE) pr-sync --pr $(PR) --goal $(GOAL_ID)

.PHONY: auto-commit
auto-commit:
	$(XLIBGATE) auto-commit --goal $(GOAL_ID) --task $(TASK_ID)

.PHONY: version-plan
version-plan:
	$(XLIBGATE) version-plan --goal $(GOAL_ID)

.PHONY: release-plan
release-plan:
	$(XLIBGATE) release-plan --goal $(GOAL_ID)

.PHONY: release-draft
release-draft:
	$(XLIBGATE) release-draft --goal $(GOAL_ID) --version $(VERSION)

.PHONY: release-publish
release-publish:
	$(XLIBGATE) release-publish --goal $(GOAL_ID) --version $(VERSION)
```

### 10.12 PR-12 验收清单

```
- [ ] issue-sync 可解析 Issue 并生成 Goal Registry
- [ ] issue-plan 可生成执行计划
- [ ] issue-state 可同步 Goal State 到 Issue labels
- [ ] pr-body 可生成符合模板的 PR body
- [ ] pr-sync 可同步 PR metadata、labels、linked issues
- [ ] auto-commit 有 scope / secret / generated artifact guard
- [ ] auto-commit 仅限 bot branch，禁止 main/release/*
- [ ] version-plan 可推荐版本并生成 JSON
- [ ] release-plan 可生成发布计划
- [ ] release-draft 可生成 GitHub Release 草稿
- [ ] release-publish 必须通过 human approval
- [ ] GitHub Actions 权限最小化
- [ ] Issue 关闭依赖 DONE with evidence
- [ ] docs-check PASS && governance-check PASS
```

---

## 11. Runtime Constitution

goalkit v0.1.0 必须遵守的不可违反原则：

```
C1.  Goal Kernel 必须保持小型。
C2.  Harness Runtime 是执行控制面，不是工具集。
C3.  Release/Downstream/Automation/Publishing 都是 Extension。
C4.  Extension 只能按 Mode/Impact 启用。
C5.  Evidence 必须 Ledger 化。
C6.  Gate 必须有成本预算。
C7.  A-Z 是 Capability Catalog，不是必经流程。
C8.  PR-12 自动化默认低风险启用。
C9.  Full Mode 不得滥用。
C10. 新 Gate 必须先证明已有 Gate 无法覆盖。
C11. Harness 定义裁判方式，xlibgate 执行裁判。
C12. Makefile 只是入口，不是规则源。
C13. Generated Artifact = 由 xlibgate/make 命令生成的报告、证书、清单，
     不得提交源码历史，由 .gitignore 覆盖。
C14. 每个 Extension 的生成物路径必须在 .gitignore 中声明。
```

---

## 12. 时间计划

### 1 天

```
PR-1：templates + docs
产出：.agent/acceptance,delivery,handover,downstream,certification,signoff/ + docs/standard/
验收：make docs-check && make governance-check
不碰：Makefile、xlibgate、CI
```

### 3 天

```
PR-2：schemas + output contract
PR-3：runtime index + ADR
产出：.agent/schemas/ JSON Schema + runtime version index + ADR
```

### 7 天

```
PR-4：Harness + xlibgate 实现
PR-5：artifact policy + blocking activation
产出：G12-G16 可执行 + Full Mode blocking 生效
验收：GOWORK=off make goal-runtime-final GOAL_ID=GOAL-20260603-XLIB-RUNTIME-001
```

### 30 天

```
PR-6 ~ PR-8：freeze/drift + runtime-as-code + test harness
产出：防漂移、可测试、可信任
```

### 60-90 天

```
PR-9 ~ PR-12：生态、监控、成熟化、自动化
产出：downstream orchestration + observability + DX + automation
```

---

## 13. MVA 成功定义

goalkit v0.1.0 的最小可验收状态（PR-1 ~ PR-5 完成后）：

```
Goal Kernel 可表达目标            ← Goal 对象模型可用
Harness Runtime 可路由 Mode       ← mode routing 生效
G12-G16 可执行                    ← xlibgate 5 个命令可运行
Evidence Ledger 可记录            ← gate result 可输出
Full Mode 可阻断                  ← blocking policy 生效
Completion Certificate 可生成     ← certify 命令可运行
DONE with evidence 可成立         ← goal-runtime-final 通过
```

一条命令验证 MVA：

```bash
GOAL_ID=GOAL-20260603-XLIB-RUNTIME-001 \
GOAL_RUNTIME_MODE=FULL \
GOWORK=off make goal-runtime-final
```

---

## 14. 衡量指标

### 执行指标

```
Goal 完成率
平均 PR 数 / Goal
平均 Gate 耗时
Gate 失败率
Evidence 缺失率
```

### 质量指标

```
Release rollback 次数
Manifest mismatch 次数
Generated artifact 误提交次数
x.go boundary violation 次数
```

### 成本指标

```
Full Mode 使用比例（目标 <30%）
小改动被误判 Full Mode 次数
Gate 平均耗时（目标 <30s）
```

---

## 15. 原始文件处置建议

`.worktree/goal/` 原始 38 个文件不得由本提案直接删除或归档。处置规则如下：

- 逐文件路径、SHA-256、重复组、迁移目标与处置说明以 `split/docs/plans/goalkit-v0.1.0-migration-index.md` 为唯一审计索引。
- 重复文件只能标记为重复候选并保留审计证据；不得因“本文件已替代”而直接删除。
- 归档动作只能在迁移索引审核通过后执行；归档前不得改变 `goal/` 来源文件内容。
- 如本提案与标准、ADR、roadmap 或迁移索引冲突，以拆分后的 `split/docs/` 工件为准。

归档建议：只有在 `split/docs/plans/goalkit-v0.1.0-migration-index.md` 审核通过后，才可移入 `goal/archive/`；未审核前不得删除来源文件。

---

## 16. 最终结论

goalkit v0.1.0 的正确路径：

```
不是继续扩展 Runtime 层。
不是新增 PR 序号。
不是更多的 Gate 编号。
不是更长的 DONE 公式。

而是：
  Goal Kernel 极小化
  Harness Runtime 一级控制面
  Evidence Ledger 化
  Extensions 按需启用
  Automation 分阶段
  28 PR → 12 PR
  37 项 DONE → 最多 11 项
  38 个来源文件 → 标准 / ADR / roadmap / migration index 拆分承接
```

一句话：

> goalkit v0.1.0 的终极目标，是让每一个 Goal 都能用最轻的 Mode 完成、用最少的 Gate 证明、用最小的 Evidence 闭环。
