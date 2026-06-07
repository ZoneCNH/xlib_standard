# EVID-TASK-GOAL-20260607-001-004-20260607-001

## Goal

- Goal: GOAL-20260607-001
- Task: TASK-GOAL-20260607-001-004
- Worktree: /home/xlib-standard/.worktree/workspaces/project-subagents
- Branch: codex/project-subagents
- Commit: pending before commit; final pushed commit is recorded in the task completion response because this evidence file is part of the same commit.

## Acceptance Criteria

- Release rules require every successful `main` merge to correspond to exactly one new stable semver patch release.
- Release policy binds merge-to-main versioning to `.github/workflows/release-auto-patch.yml`, duplicate-commit tag reuse, and the `release-auto-patch-main` concurrency group.
- Release docs and standard docs state the `vX.Y.(Z+1)` calculation, same-workflow release publication, and `already_released=true` rerun guard.
- Relevant docs/rules/governance gates pass.

## Changed Files

- `.agent/rules/release-rules.md`
- `.agent/policies/release.yaml`
- `docs/release.md`
- `docs/standard/release-standard.md`
- `docs/evidence/GOAL-20260607-001/merge-version-increment.md`

## Commands

- `git branch --show-current`
- `git status --short --branch --untracked-files=all`
- `git worktree list`
- `git diff --check`
- `GOWORK=off make docs-check`
- `GOWORK=off make rules-verify`
- `XLIB_CONTEXT=local_write GOWORK=off make governance-check`

## Results

- Current branch check returned `codex/project-subagents`.
- Status check before edits showed a clean feature worktree tracking `origin/codex/project-subagents`.
- Worktree check confirmed this task is running in `/home/xlib-standard/.worktree/workspaces/project-subagents`, not on `main`.
- `git diff --check` passed with no whitespace errors.
- `GOWORK=off make docs-check` passed.
- `GOWORK=off make rules-verify` passed with 419 total rules, 388 active rules, and valid `enforced_by` commands for active rules.
- `XLIB_CONTEXT=local_write GOWORK=off make governance-check` passed, including main guard, worktree guard, evidence-check, adoption-check, boundary, security secret scan, contracts, docs-check, CLI registry, issue registry, command registry, makefile baseline, audit-goal, rules-consistency-check, debt enforcement, and traceability-check.

## Risks And Gaps

- Live merge to `main` and GitHub Actions release publication were not executed from this feature branch.
- `make security` inside governance-check reported `govulncheck suspended`; no dependency vulnerability scan was claimed.
- The existing `.github/workflows/release-auto-patch.yml` implementation was inspected and left unchanged because it already computes the next patch version, serializes main-push releases, tags the current SHA, publishes the GitHub Release, and reuses an existing stable tag for same-commit reruns.
