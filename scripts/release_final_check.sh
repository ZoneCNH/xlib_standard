#!/usr/bin/env bash
set -euo pipefail

echo "running release final check..."

echo ""
echo "=== Step 1: CI gate ==="
GOWORK=off make ci

echo ""
echo "=== Step 2: Release check ==="
GOWORK=off make release-check

echo ""
echo "=== Step 3: Verify manifest ==="
manifest="release/manifest/latest.json"
checksum="release/manifest/latest.json.sha256"

if [[ ! -f "$manifest" ]]; then
  echo "ERROR: manifest not found: $manifest" >&2
  exit 1
fi

if [[ ! -f "$checksum" ]]; then
  echo "ERROR: checksum not found: $checksum" >&2
  exit 1
fi

echo "verifying checksum..."
sha256sum -c "$checksum"

echo ""
echo "=== Step 4: Search for forbidden terms ==="
for term in "Goal Runtime" "Context Runtime" "Evidence Runtime" "goalcli" "score --min" "adoption-check" "downstream matrix"; do
  if grep -R --fixed-strings "$term" . --exclude-dir=.git --exclude-dir=.worktree >/dev/null 2>&1; then
    echo "WARNING: found forbidden term: $term"
  fi
done

echo ""
echo "release final check passed"
