# goalkit v0.1.0 Agent Team Context

Generated: 2026-06-03T00:53:02Z
Workspace: /home/xlib-standard
Goal: 使用 agent teams 执行 `.worktree/goalkit-v0.1.0-plan.md`

## Target Result

Use OMX team workers to execute the goalkit v0.1.0 plan through a safe MVA path:

1. Resolve the execution boundary between `.worktree/goalkit-v0.1.0-plan.md` and the split normative artifacts.
2. Land the first executable goalkit MVA slice without falsely declaring unfinished capabilities complete.
3. Produce evidence for every completion claim, ending only when commands/tests support the claim.

## Current Evidence

- `.worktree/goalkit-v0.1.0-plan.md` is currently marked as a complete executable Goal plan.
- `.worktree/split/docs/plans/goalkit-v0.1.0-roadmap.md` says it is an execution plan and points to split standard/ADR artifacts as authority.
- `.worktree/split/docs/adr/ADR-20260603-001-goalkit-xlibgate-runtime.md` previously downgraded the original plan to a non-final proposal.
- `.worktree/split/docs/standard/goalkit-runtime.md` contains goalkit runtime standard material, but `docs/standard/goalkit-runtime.md` does not yet exist.
- `.agent/harness.yaml` is still schema version 3.1 and does not yet contain a goalkit v0.1.0 harness control-plane layout.
- `.agent/command-registry.yaml` is still schema version 2.9.3 and currently registers `goal-runtime` / `acceptance-matrix`, but not G12-G16 commands.
- Existing OMX team state is inactive and belongs to `docs-sync-recheck-20260603`; do not reuse it as current execution evidence.

## Conflicts To Resolve First

- The modified root plan calls `.agent/evidence/ledger.jsonl` the Evidence Ledger path; split runtime notes mention `release/evidence/goalkit/`. Choose one by ADR/standard update, or document compatibility if both are needed.
- The root plan uses PR-0..PR-28; split roadmap uses PR-1..PR-12. Preserve one canonical delivery map or explicitly map aliases.
- The root plan requires MVA PR-0, PR-1, PR-2, PR-3, PR-4, PR-5, PR-6, PR-8; split roadmap says Core MVA is PR-1~5. Do not mark MVA complete until this is reconciled.
- The root plan asks for G12-G16 executable gates; current command registry and Makefile do not yet expose those names.

## Initial Team Lanes

Lane A - authority and docs:
- Move or reconcile split normative artifacts into the repository authority path if appropriate.
- Ensure root plan, ADR, standard, roadmap, and migration index do not contradict implementation status.
- Keep all repository docs in Chinese except code identifiers, commands, fixed protocol names, and commit titles.

Lane B - runtime and command implementation:
- Map existing `cmd/xlibgate`, `Makefile`, `.agent/harness.yaml`, and `.agent/command-registry.yaml`.
- Implement the smallest safe MVA slice needed for G12-G16 or equivalent goalkit commands.
- Avoid new dependencies unless explicitly required.

Lane C - tests and evidence:
- Add or update focused regression tests/fixtures before behavior claims.
- Verify with targeted `go test`, then relevant `GOWORK=off make ...` gates.
- Record remaining non-acceptance risks if full verification is not yet possible.

## Hard Constraints

- Do not commit generated artifacts such as `release/manifest/latest.json`.
- Do not treat `completion_certificate.md` or generated Evidence artifacts as source files unless the policy explicitly permits it.
- Do not implement goalkit v0.1.0 as a mandatory external CLI; `xlibgate` remains the executor surface.
- Do not bypass Harness decisions in `xlibgate`.
- Do not claim `DONE with evidence:` unless evidence is fresh and command-backed.
- Preserve unrelated user/team work in the dirty tree.

## Suggested Stop Condition

Stop the team run when:

1. A reconciled goalkit MVA execution boundary is committed to local files.
2. The first MVA implementation slice is present or a specific blocker is proven.
3. Verification commands have been run, with outputs summarized.
4. Remaining phases are listed as explicit follow-up tasks, not implied as completed.
