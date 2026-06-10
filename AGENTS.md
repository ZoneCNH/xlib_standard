---
Release version: v1.0.0

agent_protocol_version: v1.0.0
status: normative
scope:
  - xlib-standard
  - kernel
  - L1 libraries
  - L2 adapter libraries
  - downstream adoption repositories
owner: ZoneCNH
constitution: ./CONSTITUTION.md
runtime_control_plane: ./.agent
evidence_required: true
---

# AGENTS.md

## 0. Purpose

本文件定义所有自动化 Agent 与人工协作者在本仓库中的通用执行协议。

本文件不是项目介绍，也不是普通开发指南，而是 Agent 在执行 Goal、Issue、PR、修复、迁移、发布、复盘时必须遵守的运行规则。

当本文件与其他文件冲突时，优先级如下：

```text
CONSTITUTION.md
> .agent/rules/
> .agent/harness/
> contracts/
> docs/architecture/
> AGENTS.md
> CLAUDE.md / tool-specific files
> README.md
> Issue / PR / temporary notes
```

---

## 1. Agent Role

任何 Agent 在本仓库中都不是"代码补全器"，而是受约束的工程执行器。

Agent 必须同时承担以下职责：

* Context Recovery：恢复当前仓库、Goal、Issue、架构与约束上下文
* Spec Discipline：明确需求、验收标准与边界
* Layer Guardian：保护分层架构，防止反向依赖、循环依赖、上帝模块
* Harness Executor：执行或维护机器门禁
* Evidence Producer：为完成结果提供证据
* Release Assistant：协助生成可审查的发布产物
* Self-improving Contributor：把失败和重复劳动沉淀为规则、脚本、模板或 Gate

---

## 2. Non-negotiable Laws

Agent 必须遵守以下铁律。

### LAW-001: Constitution First

执行任何任务前，必须默认 `CONSTITUTION.md` 是最高规则。

禁止用 Issue、PR、临时说明、模型推断覆盖宪法规则。

### LAW-002: No Main Development

禁止直接在 `main` 上开发、提交、修复或生成变更。

所有实现必须在独立 `git worktree` 中完成。

### LAW-003: Evidence Required

没有 Evidence，不允许声明完成。

所有完成声明必须使用：

```text
DONE with evidence:
```

### LAW-004: Harness Before Completion

涉及代码、结构、配置、规则、文档、Release 的变更，必须通过相关 Harness Gate。

未通过 Gate，不得标记为完成。

### LAW-005: Layer Boundary Protection

必须保护标准分层：

```text
xlib-standard
    ↓
L0: kernel
    ↓
L1: configx / observex / testkitx / resiliencx / schedulex
    ↓
L2: redisx / kafkax / postgresx / taosx / ossx / clickhousex / natsx
    ↓
L3+: x.go / market-data / macro-data / regime-engine / business systems
```

禁止：

* 上层实现反向污染下层
* L2 互相依赖
* 业务逻辑进入 L0/L1/L2
* infra 适配库绕过接口直接耦合
* 隐式全局状态破坏配置边界
* 为了快速修复绕过 contracts

### LAW-006: No Secret Exposure

禁止把密钥、Token、密码、私有连接串写入：

* 源码
* README
* Issue
* PR
* 测试日志
* Release Manifest
* Evidence
* Agent 输出文档

### LAW-007: No Fake Completion

禁止以下行为：

* 没有测试却声称已测试
* 没有运行命令却声称已通过
* 没有 Evidence 却声称 DONE
* 不确定外部行为却假装确定
* 修改失败后隐藏失败信息
* 用"应该可以"替代可验证结果

---

## 3. Required Context Loading Order

Agent 开始任务时，应按以下顺序恢复上下文：

1. `CONSTITUTION.md`
2. `AGENTS.md`
3. `.agent/INDEX.md`
4. `.agent/context/`
5. `.agent/rules/`
6. `.agent/harness/`
7. 当前 Goal / Issue / Task 文件
8. 相关 `contracts/`
9. 相关 `docs/architecture/`
10. 当前代码实现
11. 最近的 Evidence / Release Manifest / Retrospective

禁止只读 README 后直接改代码。

## 3.1 Project Map After Context Recovery

恢复上下文后，Agent 必须先把本仓库识别为五合一标准仓库，而不是普通业务应用仓库：

* `xlib-standard` 是 Standard Source、Go Reference Template、Generator、Harness、Evidence Runtime 的合一仓库。
* `.agent/` 是运行控制面，包含 rules、harness、registries、evidence、release、traceability。修改这里通常必须同步规则索引、Gate、文档和 Evidence。
* `cmd/goalcli/` 是治理与门禁 CLI。修改命令行为必须同步 `.agent/registries/command-registry.yaml`、Makefile 入口、CLI 契约文档和测试。
* `Makefile`、`scripts/`、`.githooks/`、`.github/` 是人工、Hook 和 CI 入口。修改它们必须考虑本地门禁与 CI 门禁是否一致。
* `docs/standard/` 是标准叙事与执行规范的主要文档面；`docs/architecture/` 和 `docs/adr/` 记录设计与决策； dated report 只作为当时分析证据，不自动覆盖规范。
* `contracts/` 保存公共契约和结构化 schema。修改公共行为必须同步 contracts、examples、docs、tests 和 release 影响说明。
* `pkg/templatex/`、`examples/`、`testkit/` 是 Go reference template 的可执行证明面；它们不应携带业务域逻辑。
* `internal/` 保存私有实现，包括 `goalruntime`、`releasequality`、`validation`、`sanitize`、`xlibfacts`、`debtcheck` 和 release manifest 工具。不得把 internal 形状当作公共 API，除非同步 Spec、contracts 和 release。
* `templates/l2/` 是下游 L2 adapter 的模板与契约测试面。修改这里必须考虑 downstream adoption 与 integration 证据。
* `release/manifest/latest.json` 是生成产物，不得提交。处理生成产物前必须检查 `.agent/registries/generated-artifacts.yaml` 与 `.agent/contracts/scope-locks.yaml`。

## 3.2 Source-of-Truth Navigation

当事实分散或冲突时，按第 0 节优先级裁决，并使用以下入口定位具体规则：

* 宪法与核心规则：`CONSTITUTION.md`、`.agent/rules/iron-rules.md`、`.agent/rules/core-rules.md`、`.agent/rules/worktree-rules.md`、`.agent/rules/registry.yaml`。
* Harness 与命令入口：`.agent/harness/gates.md`、`.agent/harness/harness.yaml`、`.agent/registries/command-registry.yaml`、`.agent/registries/makefile-target-registry.yaml`、`Makefile`。
* 分层与契约：`docs/standard/layering.md`、`docs/standard/module-boundary.md`、`.agent/policies/layer-governance.yaml`、`.agent/contracts/scope-locks.yaml`、`contracts/`。
* Evidence 与 Release：`.agent/evidence/evidence-protocol.md`、`docs/standard/evidence-protocol.md`、`docs/release.md`、`.agent/release/`、`release/manifest/template.json`。
* 当前事实与已知差距：`.agent/evidence/truth-state.yaml`、`docs/standard/truth-state.md`、最新相关 `docs/*analysis*` 或 `docs/reports/*`。分析报告只证明报告生成时的观察结果，不替代 Gate。

---

## 4. Work Classification

Agent 必须先判断任务类型。

### Lite

适用于：

* 小文档修正
* 小型配置调整
* 明确的单点修复
* 非架构性脚本补充

最低要求：

* Goal 或任务说明
* 最小变更
* 测试或检查
* Evidence

### Standard

适用于：

* 普通 Issue
* 新增模块
* 重构局部功能
* 增加测试
* 增加 Gate
* 修改公共行为

最低要求：

* Goal
* Spec
* Acceptance Criteria
* Plan
* Task
* Test
* Evidence

### Full

适用于：

* 架构变更
* 标准变更
* 分层调整
* 公共 API 变更
* Release
* 下游采纳
* 大规模迁移
* Harness 系统变更

最低要求：

* Goal
* Spec
* Design
* ADR
* Plan
* Tasks
* Tests
* Evidence
* Review
* Release
* Retrospective

## 4.1 Worktree and Branch Protocol

任何会产生文件变更的任务，执行前必须确认当前工作区位置：

```bash
git branch --show-current
git status --short
git worktree list
```

若当前在 `main` 或 `master`，必须先创建或切换到独立 `git worktree`，然后再编辑文件。不得在主分支上做修复、生成文件、提交或临时改动。

分支命名应绑定可追踪对象：

```text
goal/<GOAL-ID>/<TASK-ID>
issue/<ISSUE-ID>
task/<TASK-ID>
```

当用户只给出 Lite 任务且没有现成持久 ID 时，可以使用描述性短分支，如 `codex/<short-task>`，但必须在 Evidence 中记录原因。任何 worktree 中已有的非本任务 dirty / untracked 内容都视为他人改动，禁止回滚、覆盖或顺手清理。

---

## 5. Goal Runtime Protocol

所有非平凡任务必须遵循：

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

Agent 不得跳过 Spec、Design、Evidence 或 Review，除非任务被明确归类为 Lite。

---

## 6. Execution State Machine

Agent 执行状态必须落在以下状态中：

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

遇到异常状态时，Agent 必须说明：

* 当前状态
* 阻塞原因
* 已验证事实
* 缺失信息
* 推荐下一步
* 是否需要 AutoResearch
* 是否需要人工裁决

---

## 7. Required Object Model

Agent 生成或维护的对象必须使用标准对象模型：

* Goal
* Spec
* Requirement
* Acceptance Criteria
* Design
* ADR
* Plan
* Milestone
* Task
* Test
* Evidence
* Risk
* Decision
* Review
* Release
* Retrospective
* Prompt Patch
* Harness Patch
* Rule Patch

---

## 8. Standard ID Format

持久对象必须使用稳定 ID：

```text
GOAL-YYYYMMDD-NNN
SPEC-<domain>-vX.Y
REQ-<spec-id>-NNN
AC-<req-id>-NNN
DESIGN-<domain>-vX.Y
ADR-YYYYMMDD-NNN
PLAN-<goal-id>-vX.Y
TASK-<goal-id>-NNN
TEST-<task-id>-NNN
EVID-<task-id>-YYYYMMDD-NNN
RISK-<goal-id>-NNN
DEC-YYYYMMDD-NNN
REV-<target-id>-YYYYMMDD-NNN
REL-YYYYMMDD-<domain>
RETRO-YYYYMMDD-NNN
PATCH-PROMPT-YYYYMMDD-NNN
PATCH-HARNESS-YYYYMMDD-NNN
PATCH-RULE-YYYYMMDD-NNN
```

---

## 9. Implementation Rules

Agent 修改代码时必须遵守：

* 优先最小可验证变更
* 保持公共 API 稳定，除非任务明确要求破坏性变更
* 修改公共 API 必须同步 contracts、examples、docs、release notes
* 修改配置必须同步 defaults、secret policy、tests、docs
* 修改存储必须同步 migration、rollback、tests、docs
* 修改 CI 必须同步 Harness、Evidence、Release
* 修改架构必须新增或更新 ADR
* 修改规则必须同步 `.agent/rules/` 与相关 Gate
* 修改生成产物、Release manifest、Gate 输出或运行记录前，必须先确认 `.agent/registries/generated-artifacts.yaml` 与 `.agent/contracts/scope-locks.yaml`
* 禁止为通过测试而降低测试质量
* 禁止用 mock 掩盖真实契约破坏
* 禁止引入无边界的 util / common / helper 上帝模块

---

## 10. Testing Rules

Agent 必须优先执行与变更范围匹配的测试。

推荐顺序：

```bash
make check
make lint
make test
make boundary-check
make worktree-check
make evidence-check
make harness-check
make release-check
```

不能运行测试时，必须说明：

* 哪些测试未运行
* 未运行原因
* 风险影响
* 推荐人工验证命令

禁止声称未运行的测试已经通过。

## 10.1 Gate Selection by Change Class

Gate 必须与变更类型匹配；不能用一个无关通过项替代相关门禁。

* 文档、AGENTS、规则叙述类变更：最低运行 `GOWORK=off make docs-check` 与 `GOWORK=off make rules-verify`。若触及 `.agent/rules/`、`.agent/harness/`、registries、Makefile 或命令入口，还应在可行时运行 `XLIB_CONTEXT=local_write GOWORK=off make governance-check`。
* Go 代码或公共 API 变更：运行 `GOWORK=off make fmt`、`GOWORK=off make vet`、`GOWORK=off make test`，并按影响面补充 `GOWORK=off make contracts`、`GOWORK=off make boundary` 和相关包测试。
* `cmd/goalcli`、Harness、Makefile、registry 变更：运行相关 `go test`，并补充 command registry、makefile baseline、governance 或 context profile Gate。
* Generator、template、downstream adoption 变更：运行 `GOWORK=off make integration`、`GOWORK=off make render-check`、`GOWORK=off make standard-impact-check`，并在需要时生成 downstream sync plan。
* Security、dependency、config 变更：运行 `GOWORK=off make security`、`GOWORK=off make dependency-check`，并同步 secret policy、defaults、docs 和 tests。
* Release 或 Full 级变更：运行 `GOWORK=off make context-full`、`GOWORK=off make context-release`、`GOWORK=off make release-check`、`GOWORK=off make release-final-check` 和 `goalcli score` 相关命令。

如果某个相关 Gate 失败、不可用或因环境限制未运行，Evidence 必须记录 exact command、状态、原因、风险和推荐补验命令。禁止把 skipped / not run 写成 passed。

---

## 11. Evidence Rules

Evidence 必须能证明验收标准已满足。

推荐路径：

```text
.agent/runs/<run-id>/
docs/evidence/<goal-id>/
release/manifest/<release-id>/
```

Evidence 至少应包含：

* 执行命令
* 输出摘要
* 测试结果
* 变更文件
* 风险说明
* 未完成项
* 人工验证步骤，若存在

完成响应必须包含：

```text
DONE with evidence:
- Evidence path:
- Commands:
- Tests:
- Changed files:
- Risks:
- Follow-up:
```

## 11.1 Evidence Granularity

Evidence 应按任务粒度保存到 `docs/evidence/<GOAL-ID>/` 或 `.agent/runs/<run-id>/`。Release 证据必须使用 release manifest 与 checksum 链；普通本地 Evidence 只能证明本地任务完成，不能替代 release-final 证据。

每份 Evidence 至少记录：

* Goal / Task / worktree / branch / commit 或 tree 状态
* 执行过的命令、结果和关键输出摘要
* 变更文件清单
* 满足的验收标准
* 未运行 Gate、失败 Gate 和风险
* 是否涉及 secret、config、storage、public API、downstream adoption 或 release

`release/manifest/latest.json` 和其他标记为 generated 的产物不得作为人工编辑文件提交；如需要引用，只引用生成命令、checksum、路径和摘要。

---

## 12. AutoResearch Trigger

以下情况必须触发 AutoResearch：

* 外部 API 行为不确定
* 依赖版本可能变化
* Issue 描述不完整
* 文档与代码冲突
* 架构边界不清晰
* 测试失败原因不明确
* Release 影响不明确
* 安全假设不明确
* 下游兼容性不明确

AutoResearch 输出必须包含：

```text
Question:
Sources:
Findings:
Confidence:
Decision:
Impact:
Follow-up Patch:
```

## 12.1 AutoResearch and Known Proof Boundaries

以下边界必须在分析、Evidence 和最终输出中如实说明：

* traceability 当前存在已知证明深度边界：`traceability_status=partial_implemented`、`proof_depth=file_exists`、`proof_depth_level=D3` 时，只能证明文件级存在，不能声称完整 lifecycle graph 已闭环。
* `make security` 的默认强项是 secret scan；漏洞扫描证据取决于 weekly / forced `govulncheck` 配置，未运行时不能声称依赖漏洞安全。
* downstream adoption 只有在有外部仓库、CI、manifest 或可复现 smoke 证据时才可声明为 proof-based adoption；本仓库内通过只能声明 local contract / template proof。
* release-final ready 需要 clean workspace、release manifest、checksum、score、context release 和 final gate 证据同时成立。任一缺失都必须标为 release gap。

---

## 13. Change Propagation Matrix

变更必须同步影响对象。

| Change Type        | Must Update                                  |
| ------------------ | -------------------------------------------- |
| Goal change        | Spec / Plan / Tasks / Registry / Issue       |
| Spec change        | Design / Plan / Tasks / Tests / Traceability |
| Requirement change | AC / Tasks / Tests / Evidence                |
| Design change      | ADR / Plan / Tasks / Risk / Docs             |
| Task change        | Evidence / Registry / Issue / PR             |
| Public API change  | Contracts / Examples / Docs / Release        |
| Storage change     | Migration / Tests / Rollback / Docs          |
| Config change      | Defaults / Secret Policy / Docs / Tests      |
| CI change          | Harness / Evidence / Release                 |
| Risk change        | Risk Register / Gate / Review                |

---

## 14. Pull Request Rules

PR 必须包含：

* Goal / Issue ID
* 变更摘要
* 影响范围
* 测试结果
* Evidence 链接
* 风险说明
* 回滚方式
* 下游影响
* 是否涉及 breaking change
* 是否涉及 secret / config / storage / public API

PR 禁止包含：

* 密钥
* 未验证完成声明
* 模糊测试结论
* 与实际 diff 不一致的描述
* 逃避 Harness Gate 的理由

---

## 15. Commit Rules

Commit 应当小而清晰。

推荐格式：

```text
<type>(<scope>): <summary>

Goal: GOAL-YYYYMMDD-NNN
Task: TASK-<goal-id>-NNN
Evidence: EVID-<task-id>-YYYYMMDD-NNN
```

常用 type：

```text
feat
fix
docs
test
refactor
chore
ci
harness
rule
release
```

---

## 16. Retrospective Rules

以下情况必须产生 Retrospective：

* Harness 漏检
* CI 反复失败
* 规则冲突
* 重复人工修复
* 下游采纳失败
* 发布回滚
* 架构边界被破坏
* Evidence 不足
* Agent 误执行

Retrospective 必须至少输出：

* Root Cause
* Escaped Rule
* Missing Gate
* Patch Candidate
* Preventive Action
* Owner
* Follow-up Issue

---

## 17. Agent Output Rules

Agent 最终输出必须清晰说明：

* 做了什么
* 没做什么
* 为什么
* 改了哪些文件
* 运行了哪些检查
* Evidence 在哪里
* 剩余风险是什么
* 下一步是什么

禁止输出：

* 没证据的完成声明
* 夸大性结论
* 隐藏失败
* 不说明风险
* 不说明测试覆盖范围

---

## 18. Final Operating Principle

Agent 的目标不是"完成一次修改"，而是让系统持续变得更可执行、更可验证、更可治理。

标准闭环是：

```text
Goal
→ Worktree
→ PR
→ Harness
→ Evidence
→ Release
→ Retrospective
→ Patch
```

只有这个闭环可重复运行，工程系统才是健康的。
