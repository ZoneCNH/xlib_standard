# ADR-20260603-001: goalkit v0.1.0 通过 xlibgate 和 Harness Runtime 执行

Status: Accepted

## Decision

goalkit v0.1.0 不作为 standalone CLI 交付。可执行面是调用 `xlibgate` 的 Makefile targets，Harness Runtime 继续作为 policy/control plane。

source evidence ledger 是 `.agent/evidence/ledger.jsonl`。`release/evidence/goalkit/` 下的 generated packs 只是 derived artifacts。

## Rationale

该决策让 MVA 与现有 governance gate architecture 保持一致，避免第二套 command authority，并确保 G12-G16 evidence 可以由 xlib-standard 现有 runtime 审计。

## Consequences

- G12-G16 equivalents 是 command-backed，并且在 goalkit v0.1.0 MVA evidence scope 内是 blocking。
- 只有 fresh source-ledger evidence 证明 full MVA 后，reports 才能暴露 `mva_status: complete`。
- v0.1.0 future work 不得引入 mandatory external `goalkit` CLI。

## Rejected

拒绝 standalone `goalkit` CLI，因为它会绕过 v0.1.0 的 `xlibgate` executor 和 Harness Runtime control-plane contract。
