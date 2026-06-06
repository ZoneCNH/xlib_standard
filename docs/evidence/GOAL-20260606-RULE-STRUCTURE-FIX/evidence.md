# Evidence: GOAL-20260606-RULE-STRUCTURE-FIX

## Goal

Deep-analyze and repair the current project rule structure using the OMX team
runtime, with governance gates proving that the rule, harness, documentation,
and command registry surfaces remain consistent.

## Team Runtime

- Team: `fix-xlib-standard-rul-205cfc11`
- Worktree: `/home/xlib-standard/.worktree/workspaces/team-rule-structure-fix`
- Branch: `agent/team-rule-structure-fix`
- Final status command: `omx team status fix-xlib-standard-rul-205cfc11 --json`
- Final status summary: `phase=complete`, `tasks.completed=3`,
  `tasks.pending=0`, `tasks.in_progress=0`, `tasks.failed=0`,
  `workers.dead=0`, `workers.non_reporting=0`

Worker mailbox reconciliation:

- `worker-1`: task 1 completed and verified; commits integrated.
- `worker-2`: task 2 completed at commit
  `7cc7c7dbe395977edf6dc6027c2f2fe10d2ce36c`; duplicate later
  cherry-pick failed because the task was already integrated; worker
  reconciliation confirmed no amend or new assignment was needed.
- `worker-3`: task 3 completed; evidence commit `5d0f704`; verification
  passed.

## Leader Reconciliation

The team patches repaired stale rule and authority references. The leader added
the missing human startup entrypoints required by `AGENTS.md` and aligned the
machine registry with those entrypoints:

- Added `.agent/INDEX.md` as the human context-recovery map.
- Added `.agent/context/README.md` as the context directory contract.
- Added `docs/architecture/README.md` as the architecture entrypoint.
- Updated `.agent/index.yaml` authority order and file registry entries.
- Updated runtime ownership for `docs/architecture/`.
- Updated goalcli required index paths and classification tests.
- Updated `CLAUDE.md` context loading order to include `.agent/INDEX.md`.

## Verification Commands

Passed:

- `GOWORK=off go test ./cmd/goalcli -run TestCommandRegistryRequiresFullCommandSurface -count=1`
- `GOWORK=off go run ./cmd/goalcli docs-check`
- `GOWORK=off go run ./cmd/goalcli rules-consistency-check`
- `GOWORK=off go run ./cmd/goalcli rules-verify`
- `GOWORK=off go run ./cmd/goalcli command-registry`
- `GOWORK=off go test ./... -count=1`
- `GOWORK=off go vet ./...`
- `GOWORK=off make docs-check`
- `GOWORK=off make rules-consistency-check`
- `GOWORK=off make evidence-check`
- `GOWORK=off make governance-check`
- `GOWORK=off make lint`
- `GOWORK=off make rules-verify`
- `git diff --check`

Known failed-then-fixed gate:

- Initial `GOWORK=off go run ./cmd/goalcli command-registry` failed because
  `.agent/index.yaml` did not list
  `.agent/evidence/team-rule-structure-fix-worker-3.md`.
- Added that evidence entry to `.agent/index.yaml`.
- Reran `GOWORK=off go run ./cmd/goalcli command-registry`; it passed.

## Governance Notes

`GOWORK=off make governance-check` passed. During that run,
`traceability-check` reported the current contract status as passing while
listing `REQ-001` through `REQ-012` with partial lifecycle graph gaps. Those
gaps are pre-existing governance debt surfaced by the gate output, not a blocker
under the current passing contract.

Security note: the governance run reported govulncheck as suspended unless
`XLIB_ENABLE_VULNCHECK=1` is set; the secret scan passed.

## Changed Files

- `.agent/INDEX.md`
- `.agent/context/README.md`
- `.agent/index.yaml`
- `.agent/policies/runtime-file-ownership.yaml`
- `CLAUDE.md`
- `cmd/goalcli/governance.go`
- `cmd/goalcli/main_test.go`
- `docs/architecture/README.md`
- `docs/evidence/GOAL-20260606-RULE-STRUCTURE-FIX/evidence.md`

## Residual Risk

- Downstream adoption repositories were not run.
- Template release and e2e checks were not run.
- Govulncheck was not run without the suspended-gate environment override.
- Existing traceability lifecycle graph gaps remain reported as partial.
