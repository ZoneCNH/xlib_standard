# Evidence: Project Subagents

Goal: GOAL-20260607-001
Task: TASK-GOAL-20260607-001-001
Evidence: EVID-TASK-GOAL-20260607-001-001-20260607-001
Date: 2026-06-07
Worktree: `/home/xlib-standard/.worktree/workspaces/project-subagents`
Branch: `codex/project-subagents`
Base commit: `726db7d`

## Scope

配置 `xlib-standard` 项目级 Codex native subagents，使仓库内可用的 subagent 明确继承本仓库的宪法、worktree、分层、Gate 和 Evidence 约束。

## Changed Files

- `.gitignore`
- `.codex/agents/README.md`
- `.codex/agents/xlib-explore.toml`
- `.codex/agents/xlib-executor.toml`
- `.codex/agents/xlib-verifier.toml`
- `.codex/agents/xlib-release-reviewer.toml`
- `docs/evidence/GOAL-20260607-001/project-subagents.md`

## Acceptance Criteria

- Project-level native subagent TOMLs exist under `.codex/agents/`.
- Subagents cover repository exploration, scoped execution, verification, and release readiness review.
- Each subagent explicitly preserves repository laws: Constitution first, no main development, Evidence, Harness Gates, layer boundaries, generated artifact policy, and no secret exposure.
- `.codex/agents` is reviewable as source-controlled project configuration despite the local `.git/info/exclude` default.
- Configuration parses as valid TOML.

## Commands and Results

| Command | Result | Summary |
| --- | --- | --- |
| `git branch --show-current` | passed | Current branch is `codex/project-subagents`. |
| `git status --short --untracked-files=all` | passed | Shows `.gitignore`, `.codex/agents/*`, and this Evidence file as the task changes. |
| `git worktree list` | passed | Confirms the task worktree is separate from `/home/xlib-standard` on `main`. |
| `python3 - <<'PY' ... tomllib.load(...) ... PY` | passed | Parsed 4 project agent TOMLs: `xlib-executor`, `xlib-explore`, `xlib-release-reviewer`, `xlib-verifier`. |
| `git diff --check -- .gitignore` | passed | No whitespace errors in tracked `.gitignore` diff. |
| `python3 - <<'PY' ... trailing whitespace check ... PY` | passed | No trailing whitespace in `.gitignore` or `.codex/agents/*`. |
| `GOWORK=off make worktree-guard` | passed | `goalcli worktree-guard --context local_write` passed for `codex/project-subagents`. |
| `GOWORK=off make docs-check` | passed | `docs-check passed`. |
| `GOWORK=off make rules-verify` | passed | 419 total rules, 388 active; all active rules have valid `enforced_by` commands. |
| `GOWORK=off make security` | passed | Secret check passed; `govulncheck` remains suspended unless `XLIB_ENABLE_VULNCHECK=1`. |
| `XLIB_CONTEXT=local_write GOWORK=off make governance-check` | passed | `main-guard`, `worktree-guard`, `evidence-check`, `adoption-check`, `boundary`, `security`, `contracts`, `docs-check`, registry checks, debt checks, audit checks, and `traceability-check` passed. `traceability-check` still reports the repository's known D3 file-existence proof boundary. |

## Risks and Gaps

- `govulncheck` was not run because the repository security target keeps it suspended by default without `XLIB_ENABLE_VULNCHECK=1`.
- `traceability-check` passed inside `governance-check`, but its current proof depth remains the repository's known D3 file-existence boundary rather than full lifecycle graph proof.
- This task configures Codex native subagents only. It does not configure OMX team/swarm workers, skills, or external plugins.
- The current Codex session may need to be reopened or reloaded before the new project agent TOMLs are picked up by the runtime.
