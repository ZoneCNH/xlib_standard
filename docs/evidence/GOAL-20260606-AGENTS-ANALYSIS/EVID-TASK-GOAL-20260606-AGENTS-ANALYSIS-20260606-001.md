# EVID-TASK-GOAL-20260606-AGENTS-ANALYSIS-20260606-001

## Scope

- Goal: GOAL-20260606-AGENTS-ANALYSIS
- Task: Deeply analyze current project governance and update `AGENTS.md`.
- Worktree: `/home/xlib-standard/.worktree/workspaces/update-agents-analysis`
- Branch: `codex/update-agents-analysis`
- Base commit: `6adcdff 合并规则结构治理修复`
- Classification: Lite documentation / governance navigation update.
- Branch note: User did not provide an existing durable Goal or Task ID; descriptive Codex branch was used and this Evidence records that exception.

## Context Recovery

Reviewed authoritative and project-specific surfaces before editing:

- `CONSTITUTION.md`
- `AGENTS.md`
- `.agent/INDEX.md`
- `.agent/rules/iron-rules.md`
- `.agent/rules/core-rules.md`
- `.agent/rules/worktree-rules.md`
- `.agent/harness/gates.md`
- `Makefile`
- `.agent/registries/command-registry.yaml`
- `.agent/evidence/evidence-protocol.md`
- `docs/standard/README.md`
- `docs/standard/xlib-standard.md`
- `docs/standard/harness-gates.md`
- `docs/standard/evidence-protocol.md`
- `docs/standard/security-and-secret-policy.md`
- `docs/standard/layering.md`
- `docs/architecture/README.md`
- `.agent/policies/layer-governance.yaml`
- `.agent/contracts/scope-locks.yaml`
- `README.md`
- `docs/project-structural-analysis-20260605.md`

## Findings

- The repository is a combined Standard Source, Go Reference Template, Generator, Harness and Evidence Runtime, not a normal business application.
- Existing `AGENTS.md` already defined the broad protocol, but did not give future agents a compact project map or source-of-truth navigation path.
- Worktree requirements were present as a law, but the concrete branch/worktree preflight was not explicit in `AGENTS.md`.
- Gate selection was listed as a generic sequence, but did not map common change classes to the relevant gates.
- Evidence rules existed, but did not distinguish local task evidence from release-final evidence and generated manifest boundaries.
- Known proof boundaries, especially traceability depth and local-only downstream proof, were documented elsewhere but absent from `AGENTS.md`.

## Changes

- Updated `AGENTS.md` with:
  - Project map after context recovery.
  - Source-of-truth navigation for rules, gates, contracts, evidence and known gaps.
  - Worktree and branch preflight protocol.
  - Generated artifact boundary reminder.
  - Gate selection by change class.
  - Evidence granularity rules.
  - AutoResearch and current proof-boundary notes.

## Commands

- PASS: `GOWORK=off make docs-check`
  - Output summary: `go run ./cmd/goalcli docs-check`; `docs-check passed`.
- PASS: `GOWORK=off make rules-verify`
  - Output summary: `rules total: 419`; `rules active: 388`; `goalcli subcommands available: 124`; `makefile targets available: 118`; `all active rules have valid enforced_by commands`.
- PASS: `GOWORK=off go test ./...`
  - Output summary: all packages passed, including `cmd/goalcli`, `contracts`, `examples/*`, `internal/*`, `pkg/templatex`, `scripts`, `templates/l2/test/contract`, and `testkit`.
- PASS: `git diff --check`
  - Output summary: no whitespace errors.
- PASS: `git status --short`
  - Output summary at final verification time: `M AGENTS.md`; `?? docs/evidence/GOAL-20260606-AGENTS-ANALYSIS/`.

## Acceptance Criteria

- `AGENTS.md` gives agents a current, project-specific navigation protocol.
- The update preserves Constitution-first, worktree-only and evidence-required laws.
- The update does not change code, public APIs, contracts, generated artifacts or release files.
- Relevant documentation gates pass, or any gap is recorded honestly.

## Risks and Gaps

- This is a local documentation/governance update, not release evidence.
- Release-final gates are not required for this Lite update and are not claimed.
- Downstream adoption is not exercised or claimed.
- Security vulnerability scanning is not claimed unless the corresponding gate is run.
- `GOWORK=off make release-check`, `GOWORK=off make release-final-check`, `GOWORK=off make security`, `GOWORK=off make dependency-check`, `GOWORK=off make integration`, and downstream smoke tests were not run because this change is limited to `AGENTS.md` and local Evidence.
