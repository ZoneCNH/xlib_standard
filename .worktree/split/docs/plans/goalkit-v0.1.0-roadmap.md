# goalkit v0.1.0 — PR 路线图

> 本文件是执行计划，非权威规范。标准正文见 `split/docs/standard/goalkit-runtime.md`；ADR 见 `split/docs/adr/ADR-20260603-001-goalkit-xlibgate-runtime.md`；来源迁移证据见 `split/docs/plans/goalkit-v0.1.0-migration-index.md`。`split/` 为当前评审暂存路径；最终落地路径以 ADR 指定的 `docs/` 与 `.agent/` 位置为准。

## 0. 执行约束

- 本路线图不得声明 G12-G16 已实现或已通过；实现状态以 `.agent/command-implementation-status.yaml` 和实际 gate evidence 为准。
- 示例 `GOAL_ID` 必须使用标准格式。测试 fixture 使用 `GOAL-20260603-XLIB-RUNTIME-001`，不得使用 `test` 代替。
- downstream 目标使用 `goal-downstream-adoption`，不得覆盖或误写为已有 `downstream-adoption` 目标。

## 1. Current-State Matrix

| 能力                                       | specified | registered | dry_run_ready | implemented | verified |
| ------------------------------------------ | --------- | ---------- | ------------- | ----------- | -------- |
| Goal Kernel 对象模型                       | ✓         | ✓          | —             | —           | —        |
| Mode Router                                | ✓         | ✓          | ✓             | —           | —        |
| G12 acceptance                             | ✓         | —          | —             | —           | —        |
| G13 delivery                               | ✓         | —          | —             | —           | —        |
| G14 handover                               | ✓         | —          | —             | —           | —        |
| G15 downstream adoption (gap-only dry run) | ✓         | ✓          | ✓             | —           | —        |
| G16 certify                                | ✓         | —          | —             | —           | —        |
| Evidence Ledger                            | ✓         | ✓          | ✓             | partial     | —        |
| Automation Surface                         | ✓         | —          | —             | —           | —        |
| debt 命令                                  | ✓         | ✓          | ✓             | ✓           | ✓        |
| governance 命令                            | ✓         | ✓          | ✓             | ✓           | ✓        |
| score/version/doctor                       | ✓         | ✓          | ✓             | ✓           | ✓        |

## 2. PR Dependency Graph

```
Phase 1: Core MVA (7 天)
┌─────┐  ┌─────┐  ┌─────┐
│PR-1 │  │PR-2 │  │PR-3 │  ← 可并行
└──┬──┘  └──┬──┘  └──┬──┘
   │        │        │
   └────────┼────────┘
            │
         ┌──▼──┐
         │PR-4 │  ← 依赖 PR-1~3
         └──┬──┘
            │
         ┌──▼──┐
         │PR-5 │  ← 依赖 PR-4
         └─────┘

Phase 2: 可信治理 (30 天)
PR-6, PR-7, PR-8 ← 依赖 PR-5

Phase 3: 生态与协作 (60 天)
PR-9, PR-10 ← 依赖 PR-8

Phase 4: 成熟化 (90 天)
PR-11, PR-12 ← 依赖 PR-10
```

## 3. PR 详细规格

### PR-1: Templates + Docs

- **前置条件**：无
- **范围**：`.agent/acceptance,delivery,handover,downstream,certification,signoff/` + `docs/standard/`
- **产出**：静态模板与文档
- **不得声明**：任何命令可执行、任何 gate 已通过
- **验收命令**：`GOWORK=off make docs-check && GOWORK=off make governance-check`
- **回滚**：`git revert`，模板无副作用
- **参考执行包**：xlib_standard_pr1_goal_runtime_v3_1_1_templates_docs_execution_pack.md

### PR-2: Schemas + Output Contract

- **前置条件**：无（可与 PR-1 并行）
- **范围**：`.agent/schemas/` gate_result.schema.json, evidence.schema.json
- **产出**：JSON Schema 事实源
- **不得声明**：命令可执行
- **验收命令**：`GOWORK=off make docs-check && GOWORK=off make governance-check`
- **回滚**：`git revert`
- **参考执行包**：xlib_standard_pr2_schemas_output_contract_execution_pack.md

### PR-3: Runtime Index + ADR

- **前置条件**：无（可与 PR-1/2 并行）
- **范围**：`.agent/` runtime index + `docs/adr/`
- **产出**：兼容性声明、版本索引
- **不得声明**：命令可执行
- **验收命令**：`GOWORK=off make docs-check && GOWORK=off make governance-check`
- **回滚**：`git revert`
- **参考执行包**：xlib_standard_pr3_runtime_index_compatibility_adr_execution_pack.md

### PR-4: Harness + xlibgate 实现

- **前置条件**：PR-1, PR-2, PR-3 全部合并
- **范围**：Makefile goal-\* targets + `cmd/xlibgate/` acceptance,delivery,handover,downstream,certify + `internal/goalruntime/` + `testdata/`
- **产出**：可执行命令 + fixtures + tests
- **不得声明**：G12-G16 为 blocking required gate
- **Fixture ID**：`GOAL-20260603-XLIB-RUNTIME-001`（不得使用 `test` 作为验收 fixture）
- **验收命令**：`GOAL_ID=GOAL-20260603-XLIB-RUNTIME-001 GOWORK=off make goal-acceptance && GOAL_ID=GOAL-20260603-XLIB-RUNTIME-001 GOWORK=off make goal-runtime-final`
- **回滚**：回退 Makefile goal-\* targets、`.agent/command-registry.yaml` 的 G12-G16 注册、`cmd/xlibgate/` 新子命令、`internal/goalruntime/`、fixtures 与相关测试；不得触碰已存在的 debt/governance 命令。
- **参考执行包**：xlib_standard_pr4_makefile_harness_command_registry_execution_pack.md + xlib_standard_pr5_xlibgate_commands_fixtures_tests_execution_pack.md

### PR-5: Generated Artifact Policy + Blocking

- **前置条件**：PR-4 合并
- **范围**：`.agent/evidence-artifact-policy.yaml` 更新 + `harness.yaml` blocking activation
- **产出**：Full Mode G12-G16 阻断生效
- **不得声明**：Automation Surface 可用
- **验收命令**：`GOAL_ID=GOAL-20260603-XLIB-RUNTIME-001 GOAL_RUNTIME_MODE=FULL GOWORK=off make goal-runtime-final`
- **回滚**：回退 `.agent/harness.yaml` 的 blocking activation 与 `.agent/evidence-artifact-policy.yaml`；保留 PR-4 中已经验证的非 blocking 命令实现。
- **参考执行包**：xlib_standard_pr6 + pr8 execution packs

### PR-6: Freeze + Drift Control

- **前置条件**：PR-5 合并
- **范围**：freeze/unfreeze 机制 + drift detection
- **产出**：配置冻结与漂移检测能力
- **验收命令**：`GOWORK=off make governance-check`
- **回滚**：`git revert`

### PR-7: Runtime-as-Code + Policy-as-Code

- **前置条件**：PR-5 合并
- **范围**：runtime 配置代码化 + policy 定义
- **产出**：声明式 runtime 与 policy 配置
- **验收命令**：`GOWORK=off make governance-check`
- **回滚**：`git revert`

### PR-8: Test Harness + Trust Root

- **前置条件**：PR-5 合并
- **范围**：测试 harness 框架 + trust root 验证
- **产出**：可测试、可信任的 gate 执行环境
- **验收命令**：`GOWORK=off make test && GOWORK=off make governance-check`
- **回滚**：`git revert`

### PR-9: Downstream Adoption Orchestrator

- **前置条件**：PR-8 合并
- **范围**：downstream 采纳编排器
- **产出**：自动化 downstream 依赖追踪与通知
- **验收命令**：`GOAL_ID=GOAL-20260603-XLIB-RUNTIME-001 GOWORK=off make goal-downstream-adoption && GOWORK=off make governance-check`
- **回滚**：`git revert`，不得覆盖或删除已有 `downstream-adoption` target；若新增 target 与既有 target 冲突，优先回退 PR-9。

### PR-10: Runtime Observability + Budget

- **前置条件**：PR-8 合并
- **范围**：runtime 可观测性 + 预算控制
- **产出**：gate 执行指标采集与预算约束
- **验收命令**：`GOWORK=off make governance-check`
- **回滚**：`git revert`

### PR-11: DX + Conformance + Publishing

- **前置条件**：PR-10 合并
- **范围**：开发者体验优化 + 合规检查 + 发布流程
- **产出**：一致的 DX 与合规验证
- **验收命令**：`GOWORK=off make docs-check && GOWORK=off make governance-check`
- **回滚**：`git revert`

### PR-12: Issue/PR/Release Automation

- **前置条件**：PR-10 合并
- **范围**：Issue、PR、Release 全流程自动化（分 6 Stage 交付）
- **产出**：端到端自动化管线
- **验收命令**：`GOWORK=off make governance-check`
- **回滚**：按 Stage 逐步回退

## 4. Rollback Strategy

| 失败场景                               | 回滚行动                                                                                                                   |
| -------------------------------------- | -------------------------------------------------------------------------------------------------------------------------- |
| PR-4 合并后 G12-G16 命令失败           | PR-5 不得合并；回退 PR-4 的 Makefile targets、command registry、xlibgate 子命令、fixtures 与 tests                         |
| PR-5 blocking 导致已有 gate 误阻断     | 回退 harness.yaml 变更，保留 PR-4 命令                                                                                     |
| PR-9 downstream target 覆盖已有 target | 立即回退 PR-9 Makefile 变更，恢复原 `downstream-adoption`；`goal-downstream-adoption` 必须保持新增且无覆盖                 |
| generated artifact 误提交              | 以后续 revert/清理提交移除误提交产物，并更新 `.gitignore`/artifact policy；除非 release manager 明确批准，不改写已发布历史 |

## 5. goalkit 与 debt.md 分工

| 维度          | goalkit v0.1.0                                  | debt.md                              |
| ------------- | ----------------------------------------------- | ------------------------------------ |
| 定位          | Goal Runtime 框架与工具包                       | 债务治理的自动化落地                 |
| xlibgate 命令 | acceptance/delivery/handover/downstream/certify | debt（已实现）、governance（已实现） |
| 优先级        | 先建 Goal Runtime 框架                          | 在框架上运行 debt goals              |
| 执行顺序      | 先建框架                                        | 用框架                               |

原则：goalkit 建框架，debt 用框架。两者不应互相阻塞。

## 6. Timeline

| 阶段     | 时间  | PR       | 产出                       |
| -------- | ----- | -------- | -------------------------- |
| Core MVA | 7 天  | PR-1~5   | G12-G16 可执行 + blocking  |
| 可信治理 | 30 天 | PR-6~8   | 防漂移、可测试、可信任     |
| 生态协作 | 60 天 | PR-9~10  | downstream + observability |
| 成熟化   | 90 天 | PR-11~12 | DX + automation            |

## 7. Metrics

| 指标                  | source                   | command                | threshold | owner              |
| --------------------- | ------------------------ | ---------------------- | --------- | ------------------ |
| Goal 完成率           | evidence ledger          | evidence-check         | >80%      | goal owner         |
| 平均 PR 数/Goal       | git log + issue registry | issue-registry         | ≤3        | team lead          |
| 平均 Gate 耗时        | gate result JSON         | xlibgate --format json | <30s      | harness maintainer |
| Gate 失败率           | gate result JSON         | aggregate gate results | <20%      | harness maintainer |
| Evidence 缺失率       | evidence-check           | evidence-check         | 0%        | goal owner         |
| Release rollback 次数 | release manifest         | release history        | 0         | release manager    |
| Full Mode 使用比例    | goal registry            | goal-runtime           | <30%      | team lead          |
| 小改动误判 Full Mode  | mode routing log         | mode-router audit      | 0         | harness maintainer |
