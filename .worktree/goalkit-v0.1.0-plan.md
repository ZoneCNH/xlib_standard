# goalkit v0.1.0 × xlib-standard 完整可执行 Goal 方案

> 目标仓库：`git@github.com:ZoneCNH/xlib-standard.git`  
> 起始版本：`goalkit v0.1.0`  
> 对齐标准：Goal Runtime Prompt v3.1 / Goal Runtime v3.1.1  
> 文档性质：最终完整可执行 Goal 文档  
> 适用对象：Agent、Codex、Maintainer、Reviewer、Release Manager、Downstream Maintainer  
> 状态：设计封顶，进入执行  
> 生成时间：2026-06-03

---

## 0. 最高结论

`xlib-standard` 的最终目标不是成为一个普通模板仓库，而是成为：

```text
基础库标准源
+ Go Reference Template
+ Generator
+ Harness
+ Evidence Runtime
+ Goal Runtime
+ 标准工厂
```

但前面完整讨论也说明：Goal Runtime v3.1.1 已经足够完整，继续横向扩张会导致结构债。

因此从现在开始，统一以 `goalkit v0.1.0` 作为起始版本，将体系收敛为：

```text
Goal Kernel
+ Harness Runtime
+ Evidence Ledger
+ Completion Extension
+ Release Extension
+ Ecosystem Extension
+ Governance Extension
+ Automation Surface
```

最终一句话：

> goalkit v0.1.0 的目标，是把 Goal Runtime 从“全能治理平台”收敛为“Goal 小内核 + Harness 控制面 + Evidence Ledger + 可插拔扩展”的可执行系统。

---

## 1. North Star Goal

```text
GOAL-20260603-XLIB-GOALKIT-001

把 xlib-standard 从“基础库模板仓库”升级为
“goalkit v0.1.0 驱动的基础库标准工厂”。
```

最终必须达成：

```text
1. Goal 可被结构化表达。
2. Harness 可决定 Mode / Gate / Blocking / Evidence。
3. xlibgate 可执行 Harness 裁决。
4. Evidence Ledger 可记录事实链。
5. Full Mode 可完成 Acceptance / Delivery / Handover / Certificate。
6. Release 只在 release_impacting=true 时触发。
7. Downstream 只在 downstream_impacting=true 时触发。
8. Automation 分阶段启用，不能绕过治理。
9. 最终完成只能是 DONE with evidence。
```

---

## 2. 核心结构性修正

历史方案中曾经出现过以下风险：

```text
1. Goal 过度膨胀成全能治理平台。
2. Harness 被隐藏在 Makefile / xlibgate / Gate 分类下面。
3. goalkit 与 xlibgate 的边界不够清晰。
4. Evidence 仍可能停留在 Markdown 报告层。
5. 自动化 Issue / PR / Commit / Release 的发布风险过高。
```

最终修正：

```text
1. Goal Core 小型化。
2. Harness Runtime 一级化。
3. goalkit v0.1.0 不提供独立 CLI。
4. xlibgate 是唯一机器执行器。
5. Evidence Ledger 使用 JSONL。
6. A-Z 改成 Capability Catalog，而不是必经流程。
7. PR-28 自动化必须分阶段启用。
```

---

## 3. goalkit v0.1.0 定位

`goalkit v0.1.0` 不是独立 CLI。

在 v0.1.0 中：

```text
goalkit = Runtime 规范 + 对象模型 + 版本号 + .agent 文件 + schemas + templates
xlibgate = 唯一机器执行器
Makefile = 人类 / CI 命令入口
CI = 远程执行环境
Evidence Ledger = 事实源
```

边界表：

| 层 | 定位 | v0.1.0 是否实现 |
|---|---|---|
| goalkit | Runtime 规范 + 对象模型 + 版本号 | 是 |
| xlibgate | 唯一执行器 | 是 |
| Makefile | 人类/CI 入口 | 是 |
| goalkit CLI | 独立 CLI | 否，推迟到 v0.2.0 或以后 |

硬规则：

```text
v0.1.0 不实现 goalkit CLI。
所有执行能力归 xlibgate。
```

---

## 4. 最终架构

```text
Goal Kernel
  ↓
Harness Runtime
  ↓
xlibgate Executors
  ↓
Gate Results
  ↓
Evidence Ledger
  ↓
Completion / Release / Ecosystem / Governance / Automation Extensions
  ↓
DONE with evidence
```

完整分层：

```text
Goal Kernel
+ Harness Runtime
+ Evidence Ledger
+ Completion Extension
+ Release Extension
+ Ecosystem Extension
+ Governance Extension
+ Automation Surface
```

---

## 5. Goal Kernel

Goal Kernel 必须保持小型。

只包含：

```text
Goal
Spec
Design
Plan
Task
Test
Evidence
Review
```

只回答：

```text
要做什么？
为什么做？
怎么做？
谁来做？
如何证明完成？
```

不直接包含：

```text
Release
Publishing
Conformance
Ecosystem
Automation
Observability
```

---

## 6. Harness Runtime

Harness Runtime 是 goalkit v0.1.0 的一级控制面。

Harness 不是：

```text
不是 Makefile
不是 xlibgate
不是 CI
不是 Policy
不是 Evidence
```

Harness 是：

```text
机器裁判
+ 执行路由器
+ Gate 编排器
+ Evidence 强制器
+ Blocking 策略控制器
```

Harness 决定：

```text
当前 Goal 应该走 Lite / Standard / Full？
哪些 Gate 必须跑？
哪些 Gate 是 optional？
哪些失败必须阻断？
哪些报告必须生成？
哪些生成物不得提交？
哪些 Evidence 才能证明 DONE？
```

关系固定：

```text
Policy 定义规则
Runtime-as-Code 定义结构
Harness 决定如何执行
Makefile 暴露入口
xlibgate 执行检查
CI 自动运行
Evidence 记录结果
```

---

## 7. Harness Source of Truth

Harness 必须是数据驱动控制面。

### 7.1 最小 Source of Truth

```text
.agent/harness.yaml
.agent/registry/runtime.yaml
.agent/registry/commands.yaml
.agent/evidence/ledger.jsonl
```

### 7.2 增强 Source of Truth

```text
.agent/registry/gates.yaml
.agent/registry/makefile_baseline.yaml
.agent/harness/mode_routing.yaml
.agent/harness/impact_rules.yaml
.agent/harness/blocking_policy.yaml
.agent/harness/evidence_policy.yaml
.agent/harness/artifact_policy.yaml
.agent/harness/non_acceptance_policy.yaml
.agent/harness/failure_budget.yaml
```

执行关系：

```text
xlibgate 读取 Harness 配置
xlibgate 执行 Gate
Makefile 调用 xlibgate
CI 调用 Makefile
Evidence Ledger 记录执行事实
```

---

## 8. Evidence Ledger

Evidence 不应只是 Markdown 报告。

主文件：

```text
.agent/evidence/ledger.jsonl
```

每行一个 JSON 对象：

```json
{
  "id": "EVID-TASK-001-20260603-001",
  "goal_id": "GOAL-20260603-XLIB-GOALKIT-001",
  "task_id": "TASK-001",
  "gate_id": "G12_ACCEPTANCE",
  "commit": "abc123",
  "command": "GOWORK=off make goal-acceptance GOAL_ID=...",
  "exit_code": 0,
  "output_artifact": ".agent/acceptance/acceptance_report.md",
  "checksum": "sha256:...",
  "actor": "agent",
  "generated_at": "2026-06-03T00:00:00Z",
  "supports": ["AC-001", "TEST-001"]
}
```

规则：

```text
1. Append-only。
2. 不重写旧证据。
3. 修正使用 revision。
4. artifact 必须有 checksum。
5. CI 中必须记录 commit。
6. Markdown report 是人类视图，不是事实源。
```

Git 策略：

```text
.agent/evidence/ledger.jsonl 可 git tracked。
.agent/evidence/reports/*.md 可按需 tracked。
.agent/evidence/artifacts/ 默认 ignored。
```

---

## 9. Goal 对象模型

```yaml
goal:
  id: GOAL-YYYYMMDD-NNN
  title:
  intent:
  mode: LITE | STANDARD | FULL
  change_level: L0 | L1 | L2 | L3 | L4 | L5
  release_impacting: true | false
  downstream_impacting: true | false
  automation_allowed: true | false

spec:
  requirements:
  acceptance_criteria:
  out_of_scope:

execution:
  design:
  plan:
  tasks:
  tests:

evidence:
  ledger:
  required:
  generated:

extensions:
  completion:
    enabled: true | false
  release:
    enabled: true | false
  downstream:
    enabled: true | false
  governance:
    enabled: true | false
  automation:
    enabled: true | false
```

关键开关：

```text
release_impacting
downstream_impacting
automation_allowed
completion.enabled
release.enabled
downstream.enabled
```

---

## 10. Mode 路由与 DONE 公式

### 10.1 Lite Mode

适合：

```text
文档错别字
链接修复
注释修复
非标准源说明补充
```

需要：

```text
docs-check
Evidence summary
Review
```

DONE：

```text
DONE
= Task Complete
+ docs-check PASS
+ Evidence
+ Review
```

---

### 10.2 Standard Mode

适合：

```text
模板补充
docs/standard 新文档
.agent 模板调整
非阻断规则补充
```

需要：

```text
docs-check
governance-check
Spec / Plan / Task
Evidence
Review
```

DONE：

```text
DONE
= Spec PASS
+ Task Complete
+ Test PASS
+ Evidence
+ Review
+ Optional Acceptance
```

---

### 10.3 Full Mode

适合：

```text
Makefile
.agent/harness.yaml
xlibgate
release manifest
CI
generator
public API
security
downstream
automation
```

必须：

```text
Acceptance
Delivery
Handover
Downstream Decision
Completion Certificate
Release Manifest if release-impacting
Trust Root if release-impacting
```

DONE：

```text
DONE
= Spec PASS
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

---

### 10.4 Release DONE

只在 `release_impacting=true` 启用：

```text
DONE
= Full DONE
+ Version Plan
+ Release Manifest
+ Preflight PASS
+ Rollback Plan
+ Release Evidence
```

---

## 11. Impact 判定机制

`release_impacting` / `downstream_impacting` 不能只靠作者手填。

建议文件：

```text
.agent/harness/impact_rules.yaml
```

示例：

```yaml
release_impacting:
  triggers:
    - paths: ["VERSION", "CHANGELOG.md", "release/**"]
    - paths: ["Makefile", ".agent/harness.yaml", ".github/workflows/**"]
    - paths: ["cmd/xlibgate/**", "internal/**"]
    - paths: [".agent/schemas/**", "schemas/**"]
    - paths: ["templates/**", "generator/**"]
  default: false
  override_requires_evidence: true

downstream_impacting:
  triggers:
    - paths: ["templates/**", "generator/**", "api/**", "schemas/**", ".agent/schemas/**"]
    - paths: ["docs/standard/downstream-runtime-sync-policy.md"]
  default: false
  override_requires_evidence: true

security_impacting:
  triggers:
    - paths: [".github/workflows/**", "scripts/security*", ".agent/security/**"]
  default: false
  override_requires_evidence: true
```

判定流程：

```text
1. 作者在 Goal 中声明。
2. Harness 根据 changed files 自动推断。
3. 声明与推断不一致 → Gate FAIL 或要求 evidence override。
```

---

## 12. Extensions

### 12.1 Completion Extension

启用条件：

```text
GOAL_RUNTIME_MODE=FULL
XLIB_CONTEXT=release_verify
GOAL_RUNTIME_ACCEPTANCE=required
```

包含：

```text
Acceptance
Delivery
Handover
Completion Certificate
```

---

### 12.2 Release Extension

启用条件：

```text
release_impacting=true
```

包含：

```text
Version Plan
CHANGELOG
Release Manifest
Release Preflight
Tag
GitHub Release
Rollback
```

---

### 12.3 Ecosystem Extension

启用条件：

```text
downstream_impacting=true
```

包含：

```text
Downstream Impact
Adoption Matrix
Kernel Verification
x.go Consumer Boundary
Adoption Certificate
```

---

### 12.4 Governance Extension

启用条件：

```text
Runtime / Harness / Policy / Release / Security 变更
```

包含：

```text
Runtime-as-Code
Policy-as-Code
Trust Root
Runtime Test Harness
Drift Check
Budget
Conformance
Observability
```

---

### 12.5 Automation Surface

默认允许：

```text
issue-sync
issue-plan
pr-body
pr-sync
version-plan
release-draft
```

谨慎启用：

```text
auto-commit
release-publish
issue-close
```

绝对禁止：

```text
直接 push main
绕过 Review
绕过 Gate
无 Evidence 关闭 Issue
无 approval 发布 Release
```

---

## 13. A-Z 转为 Capability Catalog

A-Z 不是必经流程，而是能力库。

| 层 | Capability | 启用条件 |
|---|---|---|
| A | Completion | Full Mode |
| B | Schema | Runtime / CLI 变更 |
| C | Execution Contract | Harness / xlibgate 变更 |
| D | Rollout | Gate 变更 |
| E | Agent Prompt | Agent 执行 |
| F | Runbook | 人工落地 |
| G | SSOT / Drift | Runtime 文件变更 |
| H | Freeze | Runtime 封版 |
| I/J | Pack / Pack CI | Release / Runtime 交付 |
| K | Lifecycle | Runtime 运维 |
| L | Runtime-as-Code | 多文件规则漂移 |
| M/N | Downstream / Ecosystem | 下游影响 |
| O | Agent Teams | 多 Agent 并行 |
| P | Self-improving | 失败复盘 |
| Q | Policy-as-Code | 规则漂移 |
| R | Test Harness | Gate 可信度 |
| S | Trust Root | Release / 安全 |
| T | Observability | 运行期监控 |
| U | Budget | 防过重 |
| V | DX | Onboarding |
| W | Versioning | 版本演进 |
| X | Conformance | 认证 |
| Y | Publishing | 标准分发 |
| Z | Constitution | 封顶与停止 |

---

## 14. 执行路线

### Phase 0：边界冻结

```text
PR-0: ADR boundary freeze
```

必须包含：

```text
ADR-001 goalkit v0.1.0 不提供 CLI
ADR-002 Harness 是数据驱动控制面
ADR-003 Evidence Ledger 使用 JSONL
ADR-004 规范先于实现
ADR-005 Harness 不可绕过
```

---

### Phase 1：只读规范层

```text
PR-1: goalkit v0.1.0 skeleton
PR-2: schemas + fixtures
PR-3: runtime index + compatibility
```

只做：

```text
规范
Schema
Fixture
Runtime Index
```

不做：

```text
Makefile
xlibgate
CI
Release
```

---

### Phase 2：MVA 可运行核心

```text
PR-4: Harness Runtime Registration
PR-5: xlibgate harness commands
PR-6: generated artifact policy
PR-8: Full Mode blocking activation
```

MVA 验收：

```bash
GOAL_ID=GOAL-20260603-XLIB-GOALKIT-001 \
GOAL_RUNTIME_MODE=FULL \
GOWORK=off make goal-runtime-final
```

---

### Phase 3：可信治理

```text
PR-9 freeze + drift
PR-10 pack generator
PR-13 runtime-as-code
PR-18 policy-as-code
PR-19 test harness
PR-20 trust root
```

---

### Phase 4：生态和协作

```text
PR-14 downstream adoption
PR-15 ecosystem dashboard
PR-16 agent team orchestration
PR-17 self-improving
```

---

### Phase 5：成熟化

```text
PR-21 observability
PR-22 budget
PR-23 DX
PR-24 versioning
PR-25 conformance
PR-26 publishing
PR-27 constitution
```

---

### Phase 6：自动化

```text
PR-28 automation
```

最后启用，且必须分阶段。

---

## 15. Harness 验收清单

```text
- [ ] Harness 被定义为一级控制面
- [ ] .agent/harness.yaml 存在
- [ ] .agent/registry/runtime.yaml 存在
- [ ] .agent/registry/commands.yaml 存在
- [ ] .agent/evidence/ledger.jsonl 路径已定义
- [ ] Mode Router 支持 LITE / STANDARD / FULL
- [ ] impact_rules.yaml 定义 release_impacting
- [ ] impact_rules.yaml 定义 downstream_impacting
- [ ] impact_rules.yaml 定义 security_impacting
- [ ] Gate Registry 存在或等价声明存在
- [ ] Command Registry 存在
- [ ] Makefile Baseline 存在或等价检查存在
- [ ] Blocking Policy 存在
- [ ] Evidence Policy 存在
- [ ] Artifact Policy 存在
- [ ] Non-Acceptance Policy 存在
- [ ] xlibgate commands 不得绕过 Harness
- [ ] xlibgate harness self-check 可运行
```

---

## 16. 自动化验收清单

### Issue Sync

```text
- [ ] 同一个 Issue 不重复创建 Goal
- [ ] Issue 缺 Goal Runtime block 时失败
- [ ] Issue comment 使用 marker 更新，不重复刷屏
- [ ] Issue label 与 Goal State 一致
- [ ] Issue edited 后可重新 plan
- [ ] Issue close 必须依赖 DONE with evidence
```

### PR Sync

```text
- [ ] PR body 只更新 marker 区域
- [ ] Human Notes 不被覆盖
- [ ] Evidence 链接自动更新
- [ ] Issue link 自动更新
- [ ] PR labels 与 Mode / Area / Risk 一致
- [ ] Draft PR 不会被错误标记 ready
```

### Auto Commit

```text
- [ ] 当前分支必须是 agent/goal-*
- [ ] 禁止 main
- [ ] 禁止 release/*
- [ ] 禁止 hotfix/*
- [ ] scope guard 生效
- [ ] generated artifact guard 生效
- [ ] secret guard 生效
- [ ] minimal gate 通过后才 commit
- [ ] commit message 包含 Goal / Issue / Task / Evidence
```

### Release Publish

```text
- [ ] 只能 workflow_dispatch
- [ ] 必须 environment approval
- [ ] release-final-check PASS
- [ ] release-preflight PASS
- [ ] tag 不存在
- [ ] release manifest exists
- [ ] release manifest checksum valid
- [ ] completion certificate exists if Full Mode
- [ ] trust root PASS
- [ ] evidence attestation PASS
- [ ] release notes 不含 secrets
```

---

## 17. 自动化交付清单

### xlibgate 命令

```text
cmd/xlibgate issue-sync
cmd/xlibgate issue-plan
cmd/xlibgate issue-state
cmd/xlibgate pr-body
cmd/xlibgate pr-sync
cmd/xlibgate auto-commit
cmd/xlibgate version-plan
cmd/xlibgate release-plan
cmd/xlibgate release-draft
cmd/xlibgate release-publish
```

### GitHub Workflows

```text
.github/workflows/goal-runtime-sync.yml
.github/workflows/goal-runtime-release.yml
```

### Registries

```text
.agent/registry/issues.yaml
.agent/registry/prs.yaml
.agent/registry/commits.yaml
.agent/registry/versions.yaml
.agent/registry/releases.yaml
.agent/registry/automation.yaml
```

### Evidence

```text
.agent/evidence/automation/issue_sync.md
.agent/evidence/automation/pr_sync.md
.agent/evidence/automation/auto_commit.md
.agent/evidence/automation/version_plan.md
.agent/evidence/automation/release_plan.md
.agent/evidence/automation/release_draft.md
.agent/evidence/automation/release_publish.md
```

### Fixtures

```text
testdata/automation/issue-sync/pass
testdata/automation/issue-sync/missing-goal-block
testdata/automation/pr-sync/preserve-human-notes
testdata/automation/auto-commit/main-branch-denied
testdata/automation/auto-commit/generated-artifact-denied
testdata/automation/auto-commit/secret-denied
testdata/automation/version-plan/patch
testdata/automation/version-plan/minor
testdata/automation/version-plan/major
testdata/automation/release-publish/missing-approval
testdata/automation/release-publish/tag-exists
testdata/automation/release-publish/missing-manifest
```

---

## 18. Non-Acceptance Criteria

出现以下任一项，不得验收：

```text
- goalkit v0.1.0 被实现为必须有 CLI
- Harness 被硬编码在 Go 代码中，没有 YAML 来源
- xlibgate 绕过 Harness
- Evidence 只有 Markdown，没有 ledger.jsonl
- generated artifacts 被提交
- completion_certificate.md 被提交
- release/manifest/latest.json 被提交
- release-publish 可无 approval 运行
- auto-commit 可提交 main
- x.go business model 进入 xlib-standard
- release_impacting / downstream_impacting 只靠作者手填，无 Harness 校验
- PR-28 自动化默认启用 release-publish
```

---

## 19. MVA 成功定义

MVA 需要：

```text
PR-0
PR-1
PR-2
PR-3
PR-4
PR-5
PR-6
PR-8
```

MVA 成功定义：

```text
Goal Kernel 可表达目标
Harness Runtime 可路由 Mode
G12-G16 可执行
Evidence Ledger 可记录
Full Mode 可阻断
Completion Certificate 可生成
DONE with evidence 可成立
```

MVA 命令：

```bash
GOAL_ID=GOAL-20260603-XLIB-GOALKIT-001 \
GOAL_RUNTIME_MODE=FULL \
GOWORK=off make goal-runtime-final
```

---

## 20. 1 天 / 7 天 / 30 天计划

### 1 天

```text
PR-0
PR-1
```

产出：

```text
ADR 边界冻结
goalkit v0.1.0 骨架
Harness 最小锚点
Evidence Ledger 物理路径
```

---

### 7 天

```text
PR-0 ~ PR-8
```

目标：

```text
goalkit v0.1.0 MVA 可运行
```

---

### 30 天

```text
PR-9
PR-10
PR-13
PR-18
PR-19
PR-20
```

目标：

```text
可信治理、防漂移、测试、Trust Root
```

暂缓：

```text
Publishing
Badge
Full ecosystem certification
release-publish automation
```

---

## 21. 最终 DONE

Harness 完成：

```text
Harness Runtime PASS with evidence:
```

自动化完成：

```text
Issue / PR / Commit / Release automation PASS with evidence:
```

最终完成：

```text
DONE with evidence:
```

---

## 22. 最终结论

本方案已经把所有历史文档重新收敛为一个可执行 Goal：

```text
goalkit v0.1.0
= Goal Kernel
+ Harness Runtime
+ Evidence Ledger
+ Runtime Extensions
+ Safe Automation
```

最重要的结构性判断：

```text
Harness 是控制面。
Evidence Ledger 是事实源。
xlibgate 是执行器。
Makefile 是入口。
Automation 不能绕过 Harness。
```

最终一句话：

> 从 PR-0 开始，先冻结边界，再建立 Harness 和 Evidence Ledger，最后才逐步启用 Full Mode、Release 和自动化。
