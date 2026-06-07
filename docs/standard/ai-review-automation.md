# AI Review Automation

This document defines the repository control surfaces for automatic Copilot
review and local Claude pull request review.

## Control Planes

### Copilot

Copilot review is configured through GitHub repository or organization rulesets.
The repository ruleset-as-code file is `.github/rulesets/protect-main.json` and
includes a `copilot_code_review` rule:

- `review_on_push: true` asks Copilot to review updated PR heads.
- `review_draft_pull_requests: false` keeps draft pull requests out of automatic
  review.

Repository-specific review guidance lives in `.github/copilot-instructions.md`.
Those instructions align Copilot findings with the constitution, agent rules,
harness gates, layer boundaries, generated-artifact registry, and evidence
requirements.

### Local Claude

Claude review is intentionally local. This repository does not require a
GitHub Actions runner, repository Claude API key, `--bare` mode, or Claude Code
GitHub App for Claude review.

The local control surface is:

- `scripts/local_claude_pr_review.sh <pr-number>` for one pull request review;
- `scripts/local_claude_pr_watch.sh` for a local watcher that reviews changed
  non-draft PR heads.

Both scripts use:

- the local `claude` CLI, already authenticated for the local operator user;
- the local `gh` CLI, already authenticated for GitHub operations;
- local Claude invocation through `claude --print`, without `--bare`,
  repository secrets, API keys, or the Claude Code GitHub App;
- disabled Claude tool access for the review invocation;
- a review-only prompt that forbids pushing commits, creating branches, closing
  pull requests, deleting branches, changing repository settings, weakening
  branch protection, or requesting API keys.

The local `claude` CLI must have an available local subscription/session quota.
If local Claude returns a command error, such as quota exhaustion, the review
script writes an artifact, records the redacted Claude stderr, and publishes an
`error` commit status when publishing is enabled. The remediation is to rerun
from a usable local Claude session after quota refresh or local account repair;
do not add a repository API key to bypass this boundary.

The review script writes a local artifact under `.omx/artifacts/`, upserts a
sticky PR comment marked with `<!-- xlib-standard-local-claude-review -->`, and
publishes a commit status whose default context is `local-claude-review`.

Example commands:

```bash
CLAUDE_REVIEW_DRY_RUN=1 scripts/local_claude_pr_review.sh 101
scripts/local_claude_pr_review.sh 101
CLAUDE_REVIEW_WATCH_ONCE=1 scripts/local_claude_pr_watch.sh
scripts/local_claude_pr_watch.sh
```

## Merge and Branch Deletion Boundary

AI review is advisory unless a repository ruleset explicitly promotes the
resulting status context into a required status check. Merge, branch deletion,
status check selection, and branch protection remain repository governance
decisions.

Local Claude auto-merge is disabled by default. It is only available when a
trusted local operator explicitly sets `CLAUDE_REVIEW_AUTO_MERGE=1`. Even then,
the script only requests GitHub auto-merge after Claude returns
`Review result: PASS`; GitHub branch protection and required checks still decide
whether the PR can merge. Branch deletion is only requested through
`gh pr merge --auto --delete-branch` after that same PASS path.

## Rollout Checklist

1. Merge the configuration PR after normal project gates pass.
2. Install the local `claude` CLI on the operator machine and authenticate it for
   that local user. Do not use repository secrets, `ANTHROPIC_API_KEY`,
   `--bare`, or a committed token for Claude access.
3. Confirm the local Claude session has enough available quota to run a PR
   review. A quota error is recorded as a review infrastructure failure, not a
   passing review.
4. Authenticate `gh` for the same operator account with permission to read PR
   diffs, write PR comments, publish commit statuses, and optionally request
   auto-merge.
5. Run `scripts/local_claude_pr_review.sh <pr-number>` for one review, or run
   `scripts/local_claude_pr_watch.sh` as a local long-lived watcher.
6. Reconcile the live `protect-main` ruleset with
   `.github/rulesets/protect-main.json` before applying it through GitHub's
   ruleset API or UI. Do not blindly replace the live ruleset if required status
   checks or bypass actors differ.
7. Confirm Copilot automatic code review is enabled for the branch ruleset.
8. After the first successful local Claude run, decide whether
   `local-claude-review` should become a required status check. If it becomes
   required, update the ruleset, Evidence, and release notes together.

## Evidence Requirements

For every change to this automation, record:

- the changed scripts, ruleset, and instruction files;
- the exact validation commands and results;
- whether live GitHub settings were changed;
- whether the local `claude` CLI was installed and authenticated;
- whether the local `gh` CLI was authenticated for comment/status publication;
- confirmation that no repository Claude API key, `--bare` mode, or Claude Code
  GitHub App was used;
- whether local Claude review ran in dry-run or published mode;
- whether `CLAUDE_REVIEW_AUTO_MERGE` was enabled;
- any gap between ruleset-as-code and the live repository ruleset.
