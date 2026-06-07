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
ready for review.

The workflow uses:

- `anthropics/claude-code-action` pinned to a 40-character commit SHA with a tag
  comment;
- `ANTHROPIC_API_KEY` from repository or organization secrets;
- minimal repository permissions: `contents: read`, `pull-requests: write`, and
  `issues: write`;
- `id-token: write` so the pinned Claude action can mint the OIDC token it uses
  while setting up its GitHub token flow;
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
2. Configure `ANTHROPIC_API_KEY` as a repository or organization secret.
3. Reconcile the live `protect-main` ruleset with
   `.github/rulesets/protect-main.json` before applying it through GitHub's
   ruleset API or UI. Do not blindly replace the live ruleset if required status
   checks or bypass actors differ.
4. Confirm Copilot automatic code review is enabled for the branch ruleset.
5. After the first successful Claude run, decide whether `claude-review` should
   become a required status check. If it becomes required, update the ruleset,
   Evidence, and release notes together.

## Evidence Requirements

For every change to this automation, record:

- the changed workflow, ruleset, and instruction files;
- the exact validation commands and results;
- whether live GitHub settings were changed;
- whether `ANTHROPIC_API_KEY` was present;
- any gap between ruleset-as-code and the live repository ruleset.
