# Rollback Protocol

1. Identify failing REQ, file set, command and owner.
2. Revert only the smallest scoped change or apply a forward fix if safer.
3. Re-run the proof command that failed plus `git diff --check`.
4. Record rollback decision in `.agent/decision-log.md` or task result.
5. Never roll back another worker's code, Makefile, CI, manifest or gate implementation without leader approval.
