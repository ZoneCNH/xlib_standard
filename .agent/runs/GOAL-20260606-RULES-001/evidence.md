# Evidence: GOAL-20260606-RULES-001

Task: TASK-GOAL-20260606-RULES-001-001
Evidence ID: EVID-TASK-GOAL-20260606-RULES-001-001-20260606-001
Run timestamp: 2026-06-06T05:24:37+08:00
Worktree: /home/xlib-standard/.worktree/workspaces/rules-governance-drift-fix
Branch: fix/rules-governance-drift-20260606
Status: passed

## Objective

Deep-analyze `.agent/rules`, use OMX team execution to identify governance drift, and repair stale rule/runtime documentation so the rule text matches current harness commands and gate locations.

## Team Result

Team: `in-home-xlib-standard-131e4dc5`

Before shutdown, `omx team status in-home-xlib-standard-131e4dc5` reported:

```text
team=in-home-xlib-standard-131e4dc5 phase=complete
tasks: total=2 pending=0 blocked=0 in_progress=0 completed=2 failed=0
workers: total=2 dead=0 non_reporting=0
```

Leader reconciliation:

- Worker 1 found stale links and stale governance narrative around traceability, self-improving checks, and P0 coverage; those findings were applied.
- Worker 2 confirmed the live enforcement surfaces in `Makefile`, `cmd/goalcli`, `scripts/extract_rules.py`, and `scripts/verify_rules.py`; no code change was required for those surfaces.
- `omx team shutdown in-home-xlib-standard-131e4dc5` completed with no worker diffs to merge.
- A post-shutdown status check returned `No team state found for in-home-xlib-standard-131e4dc5`, confirming the team state was removed after graceful shutdown.

## Changed Files

- `.agent/rules/README.md`
- `.agent/rules/iron-rules.md`
- `.agent/runtime/standard/goal-runtime-canonical.md`
- `.agent/runs/GOAL-20260606-RULES-001/evidence.md`

## Changes

- Fixed stale documentation links from rule/runtime docs to their actual repository paths.
- Replaced stale `trace-check` / `.agent/traceability-matrix.yaml` references with the current `goalcli traceability-check`, `make traceability-check`, and `.agent/traceability/traceability-matrix.md` contract.
- Replaced stale self-improving gap wording with the current `goalcli retro-check`, `goalcli self-improving-check`, and `make retro-check` enforcement contract.
- Updated P0 coverage language so it reflects current active gate coverage instead of an obsolete "known P0 gap" section.
- Clarified that rule coverage counts are a snapshot and that live gates are the source of truth.

## Verification Commands

All commands below were run from `/home/xlib-standard/.worktree/workspaces/rules-governance-drift-fix`.

```bash
omx team status in-home-xlib-standard-131e4dc5
git diff --check
GOWORK=off make rules-verify
GOWORK=off make rules-consistency-check
GOWORK=off make traceability-check
GOWORK=off make self-improving-check
GOWORK=off go test ./cmd/goalcli -run 'Test.*(Rules|CommandImplementationStatus|CommandRegistry|Traceability|SelfImproving)' -count=1
GOWORK=off make docs-check
GOWORK=off make governance-check
GOWORK=off make worktree-check
GOWORK=off make evidence-check
rg -n "trace-check|traceability-matrix\\.yaml|已知 P0 Gap|Known P0 Gap|尚未实现|待实现|P0 gap|P0 Gap" .agent/rules .agent/runtime/standard/goal-runtime-canonical.md
```

## Verification Results

- `omx team status in-home-xlib-standard-131e4dc5`: after shutdown, returned `No team state found for in-home-xlib-standard-131e4dc5`.
- `git diff --check`: PASS.
- `GOWORK=off make rules-verify`: PASS; 419 total rules, 363 active rules, all active rules have valid `enforced_by` commands.
- `GOWORK=off make rules-consistency-check`: PASS; canonical, iron rules, and registry references are consistent.
- `GOWORK=off make traceability-check`: PASS.
- `GOWORK=off make self-improving-check`: PASS; retrospective present, 3 patch registries present, retro gate present, 0 patch entries.
- `GOWORK=off go test ./cmd/goalcli -run 'Test.*(Rules|CommandImplementationStatus|CommandRegistry|Traceability|SelfImproving)' -count=1`: PASS.
- `GOWORK=off make docs-check`: PASS.
- `GOWORK=off make governance-check`: PASS.
- `GOWORK=off make worktree-check`: PASS; worktree path and branch are valid for local write.
- `GOWORK=off make evidence-check`: PASS; registry contract satisfied.
- Stale-term scan for removed gap/link wording: PASS by absence of matches.

## Risks And Boundaries

- `make governance-check` reports `govulncheck suspended; set XLIB_ENABLE_VULNCHECK=1 to run vulnerability scan`; vulnerability scanning was not claimed.
- No registry or generated rule artifact was regenerated in this task; the fix is limited to stale governance/rule documentation.
- Worker 2 identified a remaining improvement opportunity: add a freshness gate that asserts `python3 scripts/extract_rules.py && python3 scripts/render_domain_rules.py` leaves a clean diff and that `.agent/registries/generated-artifacts.yaml generated_by` matches rule-doc generation commands.

## Follow-up

- Add the generated-rule freshness gate described above if this repository wants machine enforcement for drift between rule docs, generated registries, and generated-artifact metadata.
