#!/usr/bin/env bash
set -euo pipefail

repo="${GITHUB_REPOSITORY:-ZoneCNH/xlib-standard}"
interval="${CLAUDE_REVIEW_WATCH_INTERVAL_SECONDS:-120}"
limit="${CLAUDE_REVIEW_WATCH_LIMIT:-20}"
once="${CLAUDE_REVIEW_WATCH_ONCE:-0}"
state_dir="${CLAUDE_REVIEW_WATCH_STATE_DIR:-.omx/state/local-claude-review}"

if [[ ! "$interval" =~ ^[0-9]+$ || "$interval" -lt 1 ]]; then
  echo "ERROR: CLAUDE_REVIEW_WATCH_INTERVAL_SECONDS must be a positive integer" >&2
  exit 64
fi

if [[ ! "$limit" =~ ^[0-9]+$ || "$limit" -lt 1 ]]; then
  echo "ERROR: CLAUDE_REVIEW_WATCH_LIMIT must be a positive integer" >&2
  exit 64
fi

require_command() {
  local name="$1"
  if ! command -v "$name" >/dev/null 2>&1; then
    echo "ERROR: required command not found: $name" >&2
    exit 127
  fi
}

require_command gh
require_command git

repo_root="$(git rev-parse --show-toplevel)"
cd "$repo_root"

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
review_script="$script_dir/local_claude_pr_review.sh"
if [[ ! -x "$review_script" ]]; then
  echo "ERROR: review script is not executable: $review_script" >&2
  exit 126
fi

mkdir -p "$state_dir"

while true; do
  gh pr list --repo "$repo" --state open --limit "$limit" \
    --json number,isDraft,headRefOid \
    --jq '.[] | select(.isDraft | not) | "\(.number) \(.headRefOid)"' |
  while read -r pr_number head_sha; do
    [[ -n "${pr_number:-}" && -n "${head_sha:-}" ]] || continue

    state_file="$state_dir/pr-${pr_number}.sha"
    if [[ -f "$state_file" && "$(tr -d '\r\n' < "$state_file")" == "$head_sha" ]]; then
      continue
    fi

    echo "reviewing PR #$pr_number at $head_sha"
    if "$review_script" "$pr_number"; then
      printf '%s\n' "$head_sha" > "$state_file"
    else
      echo "local Claude review failed for PR #$pr_number" >&2
    fi
  done

  if [[ "$once" == "1" ]]; then
    break
  fi

  sleep "$interval"
done
