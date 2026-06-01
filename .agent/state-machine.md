# State Machine

```text
intake -> scope_lock -> plan -> implement -> verify -> review -> release -> retrospective -> complete
                         |          |          |          |             |
                         v          v          v          v             v
                      blocked <--- fix <--- changes_requested <--- rollback
```

## States

- `intake`: goal/context/task loaded; owner identified.
- `scope_lock`: worker scope and forbidden files recorded.
- `plan`: required artifacts, ACs, risks and verification commands mapped.
- `implement`: scoped files edited.
- `verify`: tests, docs-check, boundary/contracts/integration/release/score checks run or gaps recorded.
- `review`: reviewer validates evidence and scope compliance.
- `release`: manifest, checksum, version and final gate recorded.
- `retrospective`: defects feed prompt/harness/rule patches.
- `complete`: DONE with evidence and no open blocker.
- `blocked`: owner/action recorded; no silent partial completion.
- `rollback`: revert or mitigation path executed from rollback protocol.

## Transition rules

- `implement` cannot start before scope is locked.
- `complete` requires traceability matrix closure for all REQs.
- `release` requires `GOWORK=off make release-final-check` and score gate, unless another worker owns the still-missing executable gate and this slice records the gap.
