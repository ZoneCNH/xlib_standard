#!/usr/bin/env bash
set -euo pipefail

tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT

./scripts/render_template.sh \
  --module-name foundationx \
  --module-path github.com/ZoneCNH/foundationx \
  --package-name foundationx \
  --out "$tmpdir/foundationx"

(
  cd "$tmpdir/foundationx"
  GOWORK=off go test ./...
)

echo "integration check passed"
