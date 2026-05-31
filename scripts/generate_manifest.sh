#!/usr/bin/env bash
set -euo pipefail

mkdir -p release/manifest

MODULE="$(go list -m)"
VERSION="${VERSION:-v0.1.0}"
COMMIT="$(git rev-parse HEAD 2>/dev/null || echo unknown)"
GO_VERSION="$(go version | awk '{print $3}')"
GENERATED_AT="$(date -u +"%Y-%m-%dT%H:%M:%SZ")"
GENERATED_BY="${GENERATED_BY:-scripts/generate_manifest.sh}"
CHECK_STATUS="${CHECK_STATUS:-unknown}"
FMT_STATUS="${FMT_STATUS:-$CHECK_STATUS}"
VET_STATUS="${VET_STATUS:-$CHECK_STATUS}"
LINT_STATUS="${LINT_STATUS:-$CHECK_STATUS}"
UNIT_TEST_STATUS="${UNIT_TEST_STATUS:-$CHECK_STATUS}"
RACE_TEST_STATUS="${RACE_TEST_STATUS:-$CHECK_STATUS}"
BOUNDARY_STATUS="${BOUNDARY_STATUS:-$CHECK_STATUS}"
SECRET_SCAN_STATUS="${SECRET_SCAN_STATUS:-$CHECK_STATUS}"
SECURITY_STATUS="${SECURITY_STATUS:-$CHECK_STATUS}"
CONTRACT_STATUS="${CONTRACT_STATUS:-$CHECK_STATUS}"
INTEGRATION_STATUS="${INTEGRATION_STATUS:-$CHECK_STATUS}"
if [[ -z "$(git status --porcelain --untracked-files=all 2>/dev/null)" ]]; then
  TREE_STATE="clean"
else
  TREE_STATE="dirty"
fi

cat > release/manifest/latest.json <<JSON
{
  "module": "${MODULE}",
  "version": "${VERSION}",
  "commit": "${COMMIT}",
  "go_version": "${GO_VERSION}",
  "generated_at": "${GENERATED_AT}",
  "generated_by": "${GENERATED_BY}",
  "tree_state": "${TREE_STATE}",
  "checks": {
    "fmt": "${FMT_STATUS}",
    "vet": "${VET_STATUS}",
    "lint": "${LINT_STATUS}",
    "unit_test": "${UNIT_TEST_STATUS}",
    "race_test": "${RACE_TEST_STATUS}",
    "boundary": "${BOUNDARY_STATUS}",
    "secret_scan": "${SECRET_SCAN_STATUS}",
    "security": "${SECURITY_STATUS}",
    "contract": "${CONTRACT_STATUS}",
    "integration": "${INTEGRATION_STATUS}"
  },
  "artifacts": [
    "release/manifest/latest.json"
  ],
  "notes": {
    "breaking_changes": "none",
    "known_risks": []
  }
}
JSON

echo "generated release/manifest/latest.json"
