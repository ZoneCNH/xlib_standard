# Harness 规则

> 源自 Goal 完整规则 v1.0 §11

## RULE-HARNESS-001：Harness Gate 是阶段裁判

> **SSOT**: Gate 矩阵的权威定义在 [`core-rules.md` §131](./core-rules.md#131-lite--standard--full-gate-矩阵)（Lite/Standard/Full 三级）。执行顺序见 [`core-rules.md` §137](./core-rules.md#137-执行顺序规则)。本文件仅引用，避免重复定义导致漂移。

Gate 矩阵概要（详细定义见 core-rules.md §131）：

| Gate | Lite | Standard | Full |
|------|------|----------|------|
| schema-check | 必须 | 必须 | 必须 |
| worktree-check | 必须 | 必须 | 必须 |
| evidence-check | 必须 | 必须 | 必须 |
| traceability-check | 可选 | 必须 | 必须 |
| design-check | 可选 | 推荐 | 必须 |
| risk-check | 可选 | 必须 | 必须 |
| pr-check | 推荐 | 必须 | 必须 |
| release-check | 可选 | 推荐 | 必须 |
| retro-check | 推荐 | 必须 | 必须 |
| adoption-check | 不需要 | 可选 | 必须 |

## RULE-HARNESS-002：Gate 必须有标准结构

```yaml
gate_id:
name:
type: semantic | executable | hybrid
severity: P0 | P1 | P2 | P3
target:
rules:
commands:
pass_condition:
fail_condition:
evidence_required:
```

## RULE-HARNESS-003：P0 Gate 失败必须阻断

P0 Gate 失败时：

```text
禁止 commit
禁止 PR ready
禁止 merge
禁止 release
禁止 DONE
```

## RULE-HARNESS-004：Harness 结果必须归档

必须生成：

```text
reports/context-check.json
reports/spec-check.json
reports/design-check.json
reports/task-check.json
reports/worktree-check.txt
reports/evidence-check.json
reports/release-check.json
```
