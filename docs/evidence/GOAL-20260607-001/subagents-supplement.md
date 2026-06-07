# Evidence: Project Subagents Supplement

Goal: GOAL-20260607-001
Task: TASK-GOAL-20260607-001-001
Date: 2026-06-07
Worktree: /home/xlib-standard/.worktree/workspaces/project-subagents
Branch: codex/project-subagents
Base commit: 8c1c2e0ac847722a2a04ffac9242bf588b889cd9

## Scope

Supplement project-local Codex native subagents under `.codex/agents` after repo analysis identified missing coverage for external review, downstream adoption proof, docs-contract drift, security/dependency evidence, and layer boundary review.

No source code, public API, runtime config, storage, dependency, release manifest, or generated artifact was intentionally changed.

## Changed Files

- `.codex/agents/README.md`
- `.codex/agents/xlib-claude-reviewer.toml`
- `.codex/agents/xlib-docs-contract-drift-auditor.toml`
- `.codex/agents/xlib-downstream-adoption-auditor.toml`
- `.codex/agents/xlib-layer-boundary-reviewer.toml`
- `.codex/agents/xlib-security-dependency-auditor.toml`
- `docs/evidence/GOAL-20260607-001/subagents-supplement.md`

## Acceptance Criteria

- Add a bounded Claude CLI review coordinator that is fail-closed and advisory only.
- Add downstream adoption proof auditing that separates local proof from external proof.
- Add docs/contracts/examples/templates drift auditing.
- Add security/dependency evidence auditing with secret-safe output.
- Add layer boundary review coverage for reverse dependencies, L2 coupling, and contract bypass.
- Update the project-local subagent README routing table.
- Validate TOML syntax and required documentation/rule Gates.

## Command Evidence

| Command | Result | Summary |
| --- | --- | --- |
| `git branch --show-current` | pass | `codex/project-subagents` |
| `git status --short` | pass | Expected `.codex/agents` changes plus this Evidence file. |
| `git worktree list` | pass | Main worktree is separate; work performed in `/home/xlib-standard/.worktree/workspaces/project-subagents`. |
| `python3 -c 'import pathlib,tomllib; files=sorted(pathlib.Path(".codex/agents").glob("*.toml")); [tomllib.loads(p.read_text()) for p in files]; print("toml ok", len(files))'` | pass | `toml ok 12` |
| `git diff --check` | pass | No whitespace errors reported for the latest worktree diff. |
| `GOWORK=off make docs-check` | pass | Re-run after Evidence creation; `docs-check passed`. |
| `GOWORK=off make rules-verify` | pass | Re-run after Evidence creation; `all active rules have valid enforced_by commands`. |
| `omx_agent_validate` | not applicable | This validates committed OMX template names, not project-local `.codex/agents`; it reported missing generic templates (`architect`, `planner`, `researcher`, `executor`, `reviewer`, `operator`). |

## Risks And Gaps

- Native Codex session reload was not exercised in this shell. A new or reloaded Codex session is required before these project-local agents are available for routing.
- Claude CLI execution itself was not run; the new `xlib-claude-reviewer` agent records how to run and evidence-gate that review when requested.
- No release was executed. The project release gate remains separate and requires release manifest, checksums, score, context release, and final release evidence.

## Secret And Boundary Check

No secret values, tokens, private connection strings, dependency updates, public API changes, generated release manifests, or storage/config behavior changes were introduced.
