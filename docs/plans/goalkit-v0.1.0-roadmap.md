# goalkit v0.1.0 roadmap

Status: execution roadmap. Normative authority lives in `docs/standard/goalkit-runtime.md` and `docs/adr/ADR-20260603-001-goalkit-xlibgate-runtime.md`.

## Authority map

- Root proposal: `.worktree/goalkit-v0.1.0-plan.md`.
- Context brief: `.omx/context/goalkit-v0.1.0-team-20260603T005302Z.md`.
- Runtime standard: `docs/standard/goalkit-runtime.md`.
- Runtime ADR: `docs/adr/ADR-20260603-001-goalkit-xlibgate-runtime.md`.
- Source evidence ledger: `.agent/evidence/ledger.jsonl`.

## PR alias table

| Root plan area | Roadmap alias | Current v0.1.0 state |
| --- | --- | --- |
| PR-0/1/2 authority and registry setup | PR-1 authority split | documented, not full MVA completion |
| PR-3/4/5 Harness and command wiring | PR-4 command-backed slice | implemented as non-blocking G12-G16 equivalents |
| PR-6/8 evidence and final checks | PR-5+ blocking activation | pending fresh Harness evidence |
| PR-9+ CLI/product expansion | v0.2.0+ | out of v0.1.0 scope |

## MVA rule

The current slice exposes `goal-acceptance`, `goal-delivery`, `goal-handover`, `goal-downstream-adoption`, `goal-certify`, and `goal-runtime-final`. These commands are evidence-producing checks, but they do not make the MVA complete until Harness policy marks the gates required and fresh ledger evidence is captured.
