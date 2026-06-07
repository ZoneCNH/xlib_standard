# Evidence: Harness Subagents

Goal: GOAL-20260607-001
Task: TASK-GOAL-20260607-001-002
Evidence: EVID-TASK-GOAL-20260607-001-002-20260607-001
Date: 2026-06-07
Worktree: `/home/xlib-standard/.worktree/workspaces/project-subagents`
Branch: `codex/project-subagents`
Parent commit: `8084a09`

## Scope

配置 `xlib-standard` 项目级 Harness subagents，使 Gate 选择、Gate 执行和 Harness 契约审计职责可被单独委派、审查和验证。

## Changed Files

- `.codex/agents/README.md`
- `.codex/agents/xlib-harness-selector.toml`
- `.codex/agents/xlib-harness-runner.toml`
- `.codex/agents/xlib-harness-auditor.toml`
- `docs/evidence/GOAL-20260607-001/harness-subagents.md`

## Acceptance Criteria

- Harness subagents exist under `.codex/agents/` as project-level native Codex agent TOMLs.
- Harness responsibilities are split into gate selection, bounded gate execution, and read-only Harness contract audit.
- Each Harness subagent preserves repository laws: Constitution first, no main development, matching Harness Gates, Evidence, generated artifact policy, no secret exposure, and known proof boundaries.
- `.codex/agents/README.md` documents routing for Harness-related delegation.
- Configuration parses as valid TOML and passes matching project checks.

## Commands and Results

| Command | Result | Summary |
| --- | --- | --- |
| `git branch --show-current` | passed | Reported `codex/project-subagents`. |
| `git status --short --untracked-files=all` | passed | Reported only task-scoped changes: `.codex/agents/README.md`, three Harness subagent TOMLs, and this Evidence file. |
| `git worktree list` | passed | Confirmed `/home/xlib-standard/.worktree/workspaces/project-subagents` on branch `codex/project-subagents`; main worktree remains separate. |
| `python3 - <<'PY' ... tomllib.loads(...) ... PY` | passed | Parsed all 7 `.codex/agents/*.toml` files; each has `name`, `description`, `model`, `model_reasoning_effort`, and `developer_instructions`. |
| `git diff --check` | passed | No whitespace errors. |
| `rg -n "[[:blank:]]$" .codex/agents docs/evidence/GOAL-20260607-001/harness-subagents.md` | passed | No trailing whitespace matches. |
| `GOWORK=off make worktree-guard` | passed | `goalcli worktree-guard --context local_write` passed for branch `codex/project-subagents`. |
| `GOWORK=off make docs-check` | passed | `docs-check passed`. Required for docs / prompt / config narrative changes. |
| `GOWORK=off make rules-verify` | passed | `rules total: 419`, `rules active: 388`; all active rules have valid `enforced_by` commands. |
| `GOWORK=off make security` | passed | Secret check passed. Output also stated `govulncheck suspended; set XLIB_ENABLE_VULNCHECK=1 to run vulnerability scan`. |
| `XLIB_CONTEXT=local_write GOWORK=off make governance-check` | passed | Passed main guard, worktree guard, evidence-check, adoption-check, boundary, security, contracts, docs-check, CLI/issue/command registry, makefile baseline, audit-goal, rules-consistency-check, debt checks, and traceability-check. |

## Risks and Gaps

- This task configures Codex native subagent definitions only; it does not prove runtime reload pickup inside the already-running Codex session.
- `make security` default behavior proves secret scanning unless vulnerability scanning is forced by repository environment.
- `traceability-check` may pass with the repository's known D3 file-existence proof boundary rather than full lifecycle graph proof.
