# goalkit evidence ledger

`.agent/evidence/ledger.jsonl` is the source evidence ledger for goalkit v0.1.0.
Generated evidence packs may be written under `release/evidence/goalkit/`, but generated
release artifacts are not the source of truth and must not be treated as canonical ledger input.

The PR-4 G12-G16 commands are xlibgate/Harness-backed, non-blocking checks. They report
`mva_status: not-complete` until a later PR activates blocking Harness policy with fresh evidence.
