# 状态机

> **SSOT**: 治理状态机的权威定义在 [`CONSTITUTION.md` §9](../../CONSTITUTION.md)（13 正常 + 8 异常状态）。本文件是执行面简化版（10 正常 + 2 异常），以下提供显式映射。

## 治理状态 → 执行状态映射

| 治理状态（CONSTITUTION §9） | 执行状态（本文件） | 说明 |
|---------------------------|-----------------|------|
| INIT | intake | 初始加载 |
| CONTEXT_READY | intake | 上下文已恢复 |
| GOAL_READY | intake | Goal 已定义 |
| SPEC_READY | scope_lock | Spec 已锁定 |
| DESIGN_READY | plan | 设计已完成 |
| PLAN_READY | plan | 计划已制定 |
| TASKS_READY | plan | 任务已拆解 |
| EXECUTING | implement | 执行中 |
| VERIFYING | verify | 验证中 |
| REVIEWING | review | 审查中 |
| RELEASING | release | 发布中 |
| RETROSPECTING | retrospective | 复盘中 |
| DONE | complete | 完成 |
| BLOCKED | blocked | 阻塞 |
| FAILED | blocked | 失败 |
| NEEDS_RESEARCH | blocked | 需要研究 |
| NEEDS_DECISION | blocked | 需要决策 |
| NEEDS_REPLAN | blocked | 需要重新计划 |
| NEEDS_ROLLBACK | rollback | 需要回滚 |
| NEEDS_HUMAN_APPROVAL | blocked | 需要人工批准 |
| INCONSISTENT_STATE | blocked | 状态不一致 |

## 执行状态机

```text
intake -> scope_lock -> plan -> implement -> verify -> review -> release -> retrospective -> complete
                         |          |          |          |             |
                         v          v          v          v             v
                      blocked <--- fix <--- changes_requested <--- rollback
```

## 状态

- `intake`: goal/context/task 已加载；owner 已识别。
- `scope_lock`: worker scope 和 forbidden files 已记录。
- `plan`: required artifacts、AC、risk 和 verification commands 已映射。
- `implement`: 已编辑限定 scope 内文件。
- `verify`: 已运行 tests、docs-check、boundary/contracts/integration/release/score checks，或已记录 gaps。
- `review`: reviewer 验证 Evidence 和 scope compliance。
- `release`: manifest、checksum、version 和 final gate 已记录。
- `retrospective`: defects 回流到 prompt/harness/rule patches。
- `complete`: DONE with evidence，且没有 open blocker。
- `blocked`: owner/action 已记录；不得静默部分完成。
- `rollback`: 按 rollback protocol 执行 revert 或 mitigation path。

## 转换规则

- `implement` 不能在 scope lock 之前开始。
- `complete` 要求所有 REQ 在 traceability matrix 中关闭。
- `release` 要求 `GOWORK=off make release-final-check` 和 score gate；如果仍缺失的 executable gate 由其他 worker 负责，本 slice 必须记录 gap。
