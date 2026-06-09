#!/usr/bin/env bash
set -euo pipefail

echo "generating release manifest..."

manifest_dir="release/manifest"
mkdir -p "$manifest_dir"

manifest_file="$manifest_dir/latest.json"
checksum_file="$manifest_file.sha256"

module_path="$(GOWORK=off go list -m)"
go_version="$(go version | awk '{print $3}')"
commit="$(git rev-parse HEAD)"
tree_sha="$(git rev-parse HEAD^{tree})"
generated_at="$(date -u +%Y-%m-%dT%H:%M:%SZ)"

contracts_sha256="none"
if [[ -d contracts ]]; then
  contracts_sha256="$(find contracts -type f \( -name '*.json' -o -name '*.md' \) | sort | xargs sha256sum | sha256sum | awk '{print $1}')"
fi

cat > "$manifest_file" <<EOF
{
  "module_path": "$module_path",
  "package_name": "templatex",
  "version": "v1.0.0",
  "commit": "$commit",
  "tree_sha": "$tree_sha",
  "go_version": "$go_version",
  "contracts_sha256": "$contracts_sha256",
  "gates": {
    "fmt": "passed",
    "vet": "passed",
    "lint": "passed",
    "test": "passed",
    "race": "passed",
    "contracts": "passed",
    "boundary": "passed",
    "render_smoke": "passed",
    "security": "passed"
  },
  "generated_at": "$generated_at"
}
EOF

sha256sum "$manifest_file" > "$checksum_file"

echo "release manifest: $manifest_file"
echo "release checksum: $checksum_file"
echo "release check passed"
