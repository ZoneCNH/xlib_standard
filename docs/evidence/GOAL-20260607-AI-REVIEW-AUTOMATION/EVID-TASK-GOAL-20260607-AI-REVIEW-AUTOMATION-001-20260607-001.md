# Evidence: AI Review Automation

Goal: GOAL-20260607-AI-REVIEW-AUTOMATION
Task: TASK-GOAL-20260607-AI-REVIEW-AUTOMATION-001
Worktree: `/home/xlib-standard/.worktree/workspaces/ai-review-automation`
Branch: `codex/ai-review-automation`
Date: 2026-06-07

## Scope

Configure repository-as-code surfaces for automatic Copilot review and local
Claude pull request review while keeping merge, branch deletion, and protected
branch settings under repository governance.

## Acceptance Criteria

- Copilot automatic review is represented in the main branch ruleset-as-code.
- Copilot repository instructions align review findings with xlib-standard
  constitution, harness, evidence, generated-artifact, and layer-boundary rules.
- Local Claude review can run on demand or through a local watcher for changed
  non-draft pull request heads.
- Claude cannot merge, close, delete branches, push commits, or modify
  repository settings by default script configuration or prompt.
- Documentation records the local Claude and `gh` prerequisites, live ruleset
  rollout, dry-run/publish mode, and residual risk.
- Matching documentation and governance gates are run or explicitly recorded as
  gaps.

## Changed Files

- `.github/copilot-instructions.md`
- `.github/rulesets/protect-main.json`
- `scripts/local_claude_pr_review.sh`
- `scripts/local_claude_pr_watch.sh`
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
- PR #101 previously proved that the `claude-review` workflow queued on PR
  branch pushes, but live runner inspection found no self-hosted runner
  available. The runner workflow was removed and replaced with local scripts.
- The first attempted Claude integration used the hosted Anthropic action and
  failed because repository/API-key/App activation was not present. The workflow
  was changed after user direction to use local Claude instead of key-based
  Anthropic Action authentication.
- Local CLI inspection on this machine found `claude` at
  `/home/zone/.npm-global/bin/claude`; `claude --help` documents `--print` as a
  non-interactive output mode and documents `--bare` as API-key mode. The
  local review script intentionally uses `claude --print` and does not use
  `--bare`.

## Command Results

The following checks were run after switching Claude review away from a
self-hosted runner and into local scripts:

- `git diff --check`: passed.
- `ruby -e 'require "json"; JSON.parse(File.read(".github/rulesets/protect-main.json")); puts "json ok"'`: passed (`json ok`).
- `bash -n scripts/local_claude_pr_review.sh scripts/local_claude_pr_watch.sh`:
  passed.
- `command -v claude`: passed (`/home/zone/.npm-global/bin/claude`).
- `claude --version`: passed (`2.1.157 (Claude Code)`).
- `command -v gh`: passed (`/usr/bin/gh`).
- `gh auth status`: passed for `github.com` as account `ZoneCNH`; the displayed
  token was masked by `gh` and is not recorded here.
- `CLAUDE_REVIEW_DRY_RUN=1 CLAUDE_REVIEW_ARTIFACT_DIR=/tmp/xlib-local-claude-review-test scripts/local_claude_pr_review.sh 101`:
  reached the local Claude CLI, wrote
  `/tmp/xlib-local-claude-review-test/local-claude-pr-review-pr101-20260607T080351Z.md`,
  and exited non-zero because local Claude returned `API Error: 402 You have
  reached your subscription quota limit`. Because dry-run mode was enabled, the
  script skipped publishing the `local-claude-review=error` status and skipped
  PR comment mutation. This validates the local invocation and failure-artifact
  path, but it is not a passing Claude review.
- `GOWORK=off make docs-check`: passed (`docs-check passed`).
- `GOWORK=off make rules-verify`: passed (`all active rules have valid enforced_by commands`).
- `XLIB_CONTEXT=local_write GOWORK=off make governance-check`: passed. Key gates passed:
  `main-guard`, `worktree-guard`, `evidence-check`, `adoption-check --verify`,
  `boundary`, `security` secret check, `contracts`, `docs-check`,
  `cli-contract`, `issue-registry`, `command-registry`, `makefile-baseline`,
  `audit-goal`, `rules-consistency-check`, and `traceability-check`.
- The pushed PR state must be checked after this evidence update is committed
  and pushed. The new local Claude review path has no GitHub Actions runner
  dependency.

## Live GitHub Settings

No live ruleset update was applied in this task. The observed live `protect-main`
ruleset differs from the repository ruleset-as-code in required status checks
and bypass actors, so blindly replacing it would be a governance risk. The PR
contains the desired ruleset-as-code change and rollout documentation.

## Risks and Gaps

- Local Claude review only runs while a trusted local operator or watcher process
  runs it.
- The local machine must have the `claude` CLI installed and authenticated for
  the local user. This repository code change cannot provision or log in that
  local account.
- The local Claude account/session must have available subscription quota. The
  dry-run invocation on 2026-06-07 reached local Claude but returned 402 quota
  exhaustion; rerun the local review after local quota refresh or local account
  repair before treating `local-claude-review` as a review signal.
- The local machine must have `gh` authenticated with permission to read PR
  diffs, write PR comments, publish commit statuses, and optionally request
  auto-merge.
- Fork pull request handling is governed by the local operator's `gh` access and
  the review prompt. There is no self-hosted runner executing fork code.
- Copilot automatic review requires the live GitHub ruleset to be reconciled and
  applied after the configuration PR is accepted.
- `local-claude-review` is not added as a required status check in this task;
  making it required is a separate governance decision after the first
  successful published local run.
- `CLAUDE_REVIEW_AUTO_MERGE` is disabled by default. If a trusted local operator
  enables it, the script only requests auto-merge and branch deletion after
  Claude returns `Review result: PASS`; GitHub branch protection still controls
  merge readiness.
- `govulncheck` was not forced; the governance security step reported it
  suspended because `XLIB_ENABLE_VULNCHECK=1` was not set. This change does not
  modify Go dependencies.
