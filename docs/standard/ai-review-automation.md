# AI Review Automation

This document defines the repository control surfaces for automatic Copilot and
Claude pull request review.

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

### Claude

Claude review runs through `.github/workflows/claude-pr-review.yml` on
non-draft pull requests when they are opened, synchronized, reopened, or marked
ready for review. The workflow is limited to same-repository pull requests so
untrusted fork code is not executed on the self-hosted review runner.

The workflow uses:

- a self-hosted GitHub Actions runner with the `claude-review` label;
- the local `claude` CLI installed and already authenticated for the runner
  user;
- local Claude invocation through `claude --print`, without `--bare`, repository
  secrets, API keys, or the Claude Code GitHub App;
- disabled Claude tool access for the review invocation;
- minimal repository permissions: `contents: read`, `pull-requests: write`, and
  `issues: write`;
- an explicit review-only prompt that forbids pushing commits, creating
  branches, merging pull requests, closing pull requests, deleting branches,
  modifying repository settings, or weakening branch protection.

## Merge and Branch Deletion Boundary

AI review is advisory unless a repository ruleset explicitly promotes the
resulting workflow status into a required status check. Merge, branch deletion,
status check selection, and branch protection remain repository governance
decisions. Neither Copilot nor Claude may be configured as an autonomous merge
or branch deletion actor in this repository.

## Rollout Checklist

1. Merge the configuration PR after normal project gates pass.
2. Register a trusted self-hosted GitHub Actions runner for this repository or
   organization with the `claude-review` label.
3. Install the local `claude` CLI on that runner and authenticate it for the
   runner user. Do not use repository secrets, `ANTHROPIC_API_KEY`, `--bare`, or
   a committed token for Claude access.
4. Ensure the runner has `gh` available so the workflow can fetch the PR diff
   and publish the review comment through the workflow-provided GitHub token.
5. Reconcile the live `protect-main` ruleset with
   `.github/rulesets/protect-main.json` before applying it through GitHub's
   ruleset API or UI. Do not blindly replace the live ruleset if required status
   checks or bypass actors differ.
6. Confirm Copilot automatic code review is enabled for the branch ruleset.
7. After the first successful local Claude run, decide whether `claude-review` should
   become a required status check. If it becomes required, update the ruleset,
   Evidence, and release notes together.

## Evidence Requirements

For every change to this automation, record:

- the changed workflow, ruleset, and instruction files;
- the exact validation commands and results;
- whether live GitHub settings were changed;
- whether a trusted self-hosted runner with the `claude-review` label was
  available;
- whether the runner's local `claude` CLI was installed and authenticated;
- confirmation that no repository Claude API key, `--bare` mode, or Claude Code
  GitHub App was used;
- any gap between ruleset-as-code and the live repository ruleset.
