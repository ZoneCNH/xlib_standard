# Evidence: AI Review Automation

Goal: GOAL-20260607-AI-REVIEW-AUTOMATION
Task: TASK-GOAL-20260607-AI-REVIEW-AUTOMATION-001
Worktree: `/home/xlib-standard/.worktree/workspaces/ai-review-automation`
Branch: `codex/ai-review-automation`
Date: 2026-06-07

## Scope

Configure repository-as-code surfaces for automatic Copilot and Claude pull
request review while keeping merge, branch deletion, and protected branch
settings under repository governance.

## Acceptance Criteria

- Copilot automatic review is represented in the main branch ruleset-as-code.
- Copilot repository instructions align review findings with xlib-standard
  constitution, harness, evidence, generated-artifact, and layer-boundary rules.
- Claude runs automatically on non-draft pull request updates as an advisory
  reviewer for same-repository pull requests.
- Claude cannot merge, close, delete branches, push commits, or modify
  repository settings by workflow design or prompt.
- Documentation records the local Claude runner prerequisite, live ruleset
  rollout, and residual risk.
- Matching documentation and governance gates are run or explicitly recorded as
  gaps.

## Changed Files

- `.github/copilot-instructions.md`
- `.github/rulesets/protect-main.json`
- `.github/workflows/claude-pr-review.yml`
- `docs/standard/ai-review-automation.md`
- `docs/standard/README.md`
- `docs/evidence/GOAL-20260607-AI-REVIEW-AUTOMATION/EVID-TASK-GOAL-20260607-AI-REVIEW-AUTOMATION-001-20260607-001.md`

## External Source Confirmation

- GitHub documents Copilot automatic review as a repository or organization
  ruleset configuration.
- GitHub repository rulesets expose `copilot_code_review` parameters
  `review_on_push` and `review_draft_pull_requests`.
- GitHub documents `.github/copilot-instructions.md` as a repository custom
  instructions surface for Copilot.
- PR #101 previously proved that the `claude-review` workflow auto-triggers on
  PR branch pushes.
- The first attempted Claude integration used the hosted Anthropic action and
  failed because repository/API-key/App activation was not present. The workflow
  was changed after user direction to use local Claude instead of key-based
  Anthropic Action authentication.
- Local CLI inspection on this machine found `claude` at
  `/home/zone/.npm-global/bin/claude`; `claude --help` documents `--print` as a
  non-interactive output mode and documents `--bare` as API-key mode. The
  workflow intentionally uses `claude --print` and does not use `--bare`.

## Command Results

The following checks were run after switching Claude review to a local
self-hosted runner and recording the local activation gate:

- `git diff --check`: passed.
- `ruby -e 'require "json"; JSON.parse(File.read(".github/rulesets/protect-main.json")); puts "json ok"'`: passed (`json ok`).
- `ruby -e 'require "yaml"; YAML.load_file(".github/workflows/claude-pr-review.yml"); puts "yaml ok"'`: passed (`yaml ok`).
- `command -v claude`: passed (`/home/zone/.npm-global/bin/claude`).
- `claude --version`: passed (`2.1.157 (Claude Code)`).
- `GOWORK=off make docs-check`: passed (`docs-check passed`).
- `GOWORK=off make rules-verify`: passed (`all active rules have valid enforced_by commands`).
- `XLIB_CONTEXT=local_write GOWORK=off make governance-check`: passed. Key gates passed:
  `main-guard`, `worktree-guard`, `evidence-check`, `adoption-check --verify`,
  `boundary`, `security` secret check, `contracts`, `docs-check`,
  `cli-contract`, `issue-registry`, `command-registry`, `makefile-baseline`,
  `audit-goal`, `rules-consistency-check`, and `traceability-check`.
- The pushed PR workflow state must be checked after this evidence update is
  committed and pushed. If a `claude-review` self-hosted runner is not online,
  GitHub will leave that job queued or pending until the external runner exists.

## Live GitHub Settings

No live ruleset update was applied in this task. The observed live `protect-main`
ruleset differs from the repository ruleset-as-code in required status checks
and bypass actors, so blindly replacing it would be a governance risk. The PR
contains the desired ruleset-as-code change and rollout documentation.

## Risks and Gaps

- A trusted self-hosted GitHub Actions runner with the `claude-review` label
  must be online before Claude review can run successfully.
- The runner must have the local `claude` CLI installed and authenticated for the
  runner user. This repository code change cannot provision or log in that local
  account.
- The workflow intentionally skips fork pull requests to avoid executing
  untrusted fork code on a self-hosted runner.
- Copilot automatic review requires the live GitHub ruleset to be reconciled and
  applied after the configuration PR is accepted.
- `claude-review` is not added as a required status check in this task; making
  it required is a separate governance decision after the first successful run.
- `govulncheck` was not forced; the governance security step reported it
  suspended because `XLIB_ENABLE_VULNCHECK=1` was not set. This change does not
  modify Go dependencies.
