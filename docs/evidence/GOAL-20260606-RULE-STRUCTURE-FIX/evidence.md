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
- `XLIB_ENABLE_VULNCHECK=1 XLIB_FORCE_VULNCHECK=1 GOWORK=off make security`
- `GOWORK=off make integration`
- `GOWORK=off make p2-runtime-check`
- `GOWORK=off make release-check`
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
`XLIB_ENABLE_VULNCHECK=1` is set. A separate forced vulnerability run passed
with `XLIB_ENABLE_VULNCHECK=1 XLIB_FORCE_VULNCHECK=1 GOWORK=off make security`;
the scan reported no reachable vulnerabilities and the secret scan passed.

Release and runtime notes:

- `GOWORK=off make integration` rendered and verified `kernel`, `configx`, and
  `redisx` templates.
- `GOWORK=off make p2-runtime-check` passed runtime install, upgrade, release,
  evidence replay, conformance, pack, downstream dry-run, ownership, and
  execution-context checks.
- `GOWORK=off make release-check` passed the full release harness, including
  formatting, vet, tests, race tests, governance, rules, integration, dependency
  governance, release evidence, and checksum checks.

Push and release execution notes:

- `git push -u origin agent/team-rule-structure-fix` pushed the branch and set
  its upstream to `origin/agent/team-rule-structure-fix`.
- `git ls-remote --tags origin refs/tags/v0.5.0` confirmed remote tag
  `v0.5.0` exists at `80435cfb7df48de784e462a71346f706f147c53b`.
- `gh release view v0.5.0 --repo ZoneCNH/xlib-standard` confirmed GitHub
  Release `v0.5.0` is published at
  `https://github.com/ZoneCNH/xlib-standard/releases/tag/v0.5.0`.
- `XLIB_CONTEXT=release_verify GOWORK=off make release-preflight VERSION=v0.5.0`
  stopped at the required release branch gate:
  `ERROR: release preflight must run on main; current branch is agent/team-rule-structure-fix`.
  No duplicate release or tag was created from the repair branch.

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

- No proof-based external downstream adoption repository run is claimed; local
  `adoption-check`, `integration`, `downstream-baseline --dry-run --verify`,
  and `downstream-adoption --dry-run --verify` passed.
- Existing traceability lifecycle graph gaps remain reported as partial.
