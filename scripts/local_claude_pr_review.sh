#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'USAGE'
Usage: scripts/local_claude_pr_review.sh <pr-number>

Runs a local Claude Code review for a GitHub pull request and, by default,
publishes a sticky PR comment plus a commit status named local-claude-review.

Environment:
  GITHUB_REPOSITORY                 owner/repo, default ZoneCNH/xlib-standard
  CLAUDE_REVIEW_PUBLISH             1 to publish comment/status, default 1
  CLAUDE_REVIEW_DRY_RUN             1 to avoid GitHub mutations, default 0
  CLAUDE_REVIEW_ARTIFACT_DIR        local artifact dir, default .omx/artifacts
  CLAUDE_REVIEW_STATUS_CONTEXT      status context, default local-claude-review
  CLAUDE_REVIEW_AUTO_MERGE          1 to request gh auto-merge after PASS
  CLAUDE_REVIEW_MERGE_METHOD        squash, merge, or rebase; default squash
USAGE
}

if [[ "${1:-}" == "-h" || "${1:-}" == "--help" ]]; then
  usage
  exit 0
fi

pr_number="${1:-}"
if [[ -z "$pr_number" || ! "$pr_number" =~ ^[0-9]+$ ]]; then
  usage >&2
  exit 64
fi

repo="${GITHUB_REPOSITORY:-ZoneCNH/xlib-standard}"
publish="${CLAUDE_REVIEW_PUBLISH:-1}"
dry_run="${CLAUDE_REVIEW_DRY_RUN:-0}"
artifact_dir="${CLAUDE_REVIEW_ARTIFACT_DIR:-.omx/artifacts}"
status_context="${CLAUDE_REVIEW_STATUS_CONTEXT:-local-claude-review}"
auto_merge="${CLAUDE_REVIEW_AUTO_MERGE:-0}"
merge_method="${CLAUDE_REVIEW_MERGE_METHOD:-squash}"
marker="<!-- xlib-standard-local-claude-review -->"

require_command() {
  local name="$1"
  if ! command -v "$name" >/dev/null 2>&1; then
    echo "ERROR: required command not found: $name" >&2
    exit 127
  fi
}

require_command gh
require_command claude
require_command git
require_command date

repo_root="$(git rev-parse --show-toplevel)"
cd "$repo_root"

tmp_dir="$(mktemp -d)"
cleanup() {
  rm -rf "$tmp_dir"
}
trap cleanup EXIT

prompt_file="$tmp_dir/prompt.md"
review_file="$tmp_dir/review.md"
comment_file="$tmp_dir/comment.md"
diff_file="$tmp_dir/pr.diff"
claude_stderr_file="$tmp_dir/claude.stderr"

title="$(gh pr view "$pr_number" --repo "$repo" --json title --jq '.title')"
body="$(gh pr view "$pr_number" --repo "$repo" --json body --jq '.body // ""')"
author="$(gh pr view "$pr_number" --repo "$repo" --json author --jq '.author.login')"
url="$(gh pr view "$pr_number" --repo "$repo" --json url --jq '.url')"
state="$(gh pr view "$pr_number" --repo "$repo" --json state --jq '.state')"
is_draft="$(gh pr view "$pr_number" --repo "$repo" --json isDraft --jq '.isDraft')"
base_ref="$(gh pr view "$pr_number" --repo "$repo" --json baseRefName --jq '.baseRefName')"
head_ref="$(gh pr view "$pr_number" --repo "$repo" --json headRefName --jq '.headRefName')"
head_sha="$(gh pr view "$pr_number" --repo "$repo" --json headRefOid --jq '.headRefOid')"

publish_status() {
  local state_value="$1"
  local description="$2"

  if [[ "$publish" != "1" || "$dry_run" == "1" ]]; then
    echo "status skipped: $status_context=$state_value ($description)"
    return 0
  fi

  gh api -X POST "repos/$repo/statuses/$head_sha" \
    -f state="$state_value" \
    -f context="$status_context" \
    -f description="$description" \
    -f target_url="$url" >/dev/null
}

fail_with_status() {
  local description="$1"
  publish_status "error" "$description" || true
  echo "ERROR: $description" >&2
  exit 1
}

if [[ "$state" != "OPEN" ]]; then
  fail_with_status "PR is not open"
fi

if [[ "$is_draft" == "true" ]]; then
  publish_status "pending" "draft PR skipped by local Claude review"
  echo "PR #$pr_number is draft; skipped local Claude review."
  exit 0
fi

gh pr diff "$pr_number" --repo "$repo" > "$diff_file"
if [[ ! -s "$diff_file" ]]; then
  fail_with_status "PR diff is empty or unavailable"
fi

cat > "$prompt_file" <<PROMPT
You are reviewing a pull request for xlib-standard.

Output contract:
- Start with exactly one of these lines:
  - Review result: PASS
  - Review result: BLOCKED
  - Review result: NEEDS-HUMAN
- Then list findings by severity.
- Every actionable finding must include a concrete file path and line or hunk reference when the diff provides one.
- Prefer concise, review-grade findings over broad advice.

Repository rules to enforce:
- Constitution and AGENTS.md rules are binding.
- No secrets, tokens, passwords, or private connection strings may be introduced.
- No direct main development is allowed.
- Claims of completion require evidence and matching gates.
- Layer boundaries must not be weakened.
- Generated artifacts must stay out of commits unless explicitly allowed by the registry.

Hard limits for this review:
- Do not push commits.
- Do not create, close, merge, or delete branches or pull requests.
- Do not change repository settings.
- Do not weaken branch protection or required checks.
- Do not request or use API keys.

PR metadata:
- Repository: $repo
- PR: #$pr_number
- URL: $url
- Title: $title
- Author: $author
- Base: $base_ref
- Head: $head_ref
- Head SHA: $head_sha

PR body:
$body

Diff:
$(cat "$diff_file")
PROMPT

timestamp="$(date -u +%Y%m%dT%H%M%SZ)"
mkdir -p "$artifact_dir"
artifact_path="$artifact_dir/local-claude-pr-review-pr${pr_number}-${timestamp}.md"

write_artifact() {
  {
    echo "# Local Claude PR Review"
    echo
    echo "Repository: $repo"
    echo "PR: #$pr_number"
    echo "URL: $url"
    echo "Head SHA: $head_sha"
    echo "Status context: $status_context"
    echo "Published: $publish"
    echo "Dry run: $dry_run"
    echo "Timestamp: $timestamp"
    echo
    echo "## Prompt"
    echo
    cat "$prompt_file"
    echo
    echo "## Raw Output"
    echo
    if [[ -s "$review_file" ]]; then
      cat "$review_file"
    else
      echo "(empty)"
    fi
    if [[ -s "$claude_stderr_file" ]]; then
      echo
      echo "## Claude stderr"
      echo
      sed -E 's/(gho_|github_pat_|sk-ant-api03-|xox[baprs]-)[A-Za-z0-9_:\.-]+/[REDACTED]/g' "$claude_stderr_file"
    fi
  } > "$artifact_path"
}

if ! claude --print --input-format text --output-format text --allowedTools "" < "$prompt_file" > "$review_file" 2> "$claude_stderr_file"; then
  write_artifact
  publish_status "error" "local Claude review command failed" || true
  echo "ERROR: local Claude review command failed" >&2
  echo "Artifact: $artifact_path" >&2
  exit 1
fi

first_line="$(sed -n '1p' "$review_file" | tr -d '\r')"
case "$first_line" in
  "Review result: PASS"*)
    review_state="success"
    review_description="local Claude review passed"
    ;;
  "Review result: BLOCKED"*)
    review_state="failure"
    review_description="local Claude found blockers"
    ;;
  "Review result: NEEDS-HUMAN"*)
    review_state="failure"
    review_description="local Claude requires human review"
    ;;
  *)
    review_state="error"
    review_description="local Claude output missing review result"
    ;;
esac

write_artifact

{
  echo "$marker"
  echo "## Local Claude Review"
  echo
  echo "- PR: #$pr_number"
  echo "- Head: \`$head_sha\`"
  echo "- Status: \`$review_state\`"
  echo "- Local artifact: \`$artifact_path\`"
  echo
  cat "$review_file"
} > "$comment_file"

if [[ "$publish" == "1" && "$dry_run" != "1" ]]; then
  existing_comment_id="$(
    gh pr view "$pr_number" --repo "$repo" --json comments \
      --jq '.comments[] | select(.body | contains("'"$marker"'")) | .id' \
      | sed -n '1p'
  )"

  if [[ -n "$existing_comment_id" ]]; then
    gh api graphql \
      -f id="$existing_comment_id" \
      -F body=@"$comment_file" \
      -f query='mutation($id: ID!, $body: String!) { updateIssueComment(input: {id: $id, body: $body}) { issueComment { id } } }' >/dev/null
  else
    gh pr comment "$pr_number" --repo "$repo" --body-file "$comment_file" >/dev/null
  fi
else
  echo "comment skipped; local review comment body: $comment_file"
fi

publish_status "$review_state" "$review_description"

if [[ "$review_state" == "error" ]]; then
  echo "Review output did not follow the required result header." >&2
  echo "Artifact: $artifact_path" >&2
  exit 1
fi

if [[ "$auto_merge" == "1" ]]; then
  if [[ "$review_state" != "success" ]]; then
    echo "auto-merge skipped because local Claude review did not PASS"
  elif [[ "$dry_run" == "1" ]]; then
    echo "auto-merge skipped because CLAUDE_REVIEW_DRY_RUN=1"
  else
    case "$merge_method" in
      merge|squash|rebase)
        gh pr merge "$pr_number" --repo "$repo" --auto --delete-branch "--$merge_method"
        ;;
      *)
        echo "ERROR: unsupported CLAUDE_REVIEW_MERGE_METHOD: $merge_method" >&2
        exit 64
        ;;
    esac
  fi
fi

echo "local Claude review complete: $review_state"
echo "artifact: $artifact_path"
