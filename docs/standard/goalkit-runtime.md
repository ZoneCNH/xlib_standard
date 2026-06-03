# goalkit v0.1.0 runtime standard

Status: normative for the v0.1.0 MVA slice.

## Authority

- `xlibgate` is the only machine executor for goalkit v0.1.0 commands.
- Harness Runtime is the control plane for routing, policy, and evidence interpretation.
- goalkit v0.1.0 is **not** an independent external CLI.
- `.agent/evidence/ledger.jsonl` is the source evidence ledger.
- `release/evidence/goalkit/` is a generated evidence pack location and is not the source ledger.

## Command surface

The PR-4 command-backed slice exposes G12-G16 equivalents through Makefile targets that delegate to `xlibgate`:

| Gate | Command / target | Blocking |
| --- | --- | --- |
| G12 acceptance | `goal-acceptance` | no |
| G13 delivery | `goal-delivery` | no |
| G14 handover | `goal-handover` | no |
| G15 downstream adoption | `goal-downstream-adoption` | no |
| G16 certify | `goal-certify` | no |
| G12-G16 final report | `goal-runtime-final` | no |

Each command requires `GOAL_ID` or `--goal-id` and reports `mva_status: not-complete`. The target result proves only the current command-backed Harness slice, not full MVA completion.

## Completion rule

Do not claim the goalkit v0.1.0 MVA is complete until Harness policy activates the required gates, fresh command evidence is recorded, and the root plan / roadmap aliases are reconciled in the evidence ledger.
