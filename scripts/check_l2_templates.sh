#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'USAGE'
Usage:
  scripts/check_l2_templates.sh [templates/l2]

Checks that the L2 template fixture contains the expected contract files and
that the contract test target is executable.
USAGE
}

if [[ $# -gt 1 ]]; then
  usage >&2
  exit 2
fi

template_dir="${1:-templates/l2}"
if [[ ! -d "$template_dir" ]]; then
  echo "ERROR: missing L2 template directory: $template_dir" >&2
  exit 1
fi

required_paths=(
  ".agent/evidence/README.md"
  ".agent/evidence/l2/.gitkeep"
  ".agent/gates/l2gate.yaml"
  ".agent/l2-capabilities.yaml"
  ".github/workflows/l2-gates.yml"
  "Makefile"
  "docker-compose.test.yml"
  "test/adoption/README.md"
  "test/benchmark/README.md"
  "test/chaos/README.md"
  "test/contract/l2_contract_test.go"
  "test/integration/README.md"
)

missing=()
for path in "${required_paths[@]}"; do
  if [[ ! -f "$template_dir/$path" ]]; then
    missing+=("$path")
  fi
done

if [[ ${#missing[@]} -gt 0 ]]; then
  printf 'ERROR: missing L2 template file(s):\n' >&2
  printf '  - %s\n' "${missing[@]}" >&2
  exit 1
fi

file_count="$(find "$template_dir" -type f | wc -l | tr -d '[:space:]')"
if [[ "$file_count" != "${#required_paths[@]}" ]]; then
  echo "ERROR: expected ${#required_paths[@]} L2 template files, found $file_count" >&2
  find "$template_dir" -type f | sort >&2
  exit 1
fi

required_targets=(
  "l2-capability-check"
  "l2-contract"
  "l2-integration"
  "l2-chaos"
  "l2-benchmark"
  "l2-adoption"
  "l2-evidence"
  "l2-release-readiness-check"
  "downstream-baseline"
  "downstream-adoption"
  "p2-runtime-check"
)

for target in "${required_targets[@]}"; do
  if ! grep -Eq "^\\.PHONY:.*[[:space:]]${target}([[:space:]]|$)|^${target}:" "$template_dir/Makefile"; then
    echo "ERROR: L2 template Makefile missing target: $target" >&2
    exit 1
  fi
done

GOWORK=off make -C "$template_dir" l2-contract
echo "L2 template check passed: $template_dir (${#required_paths[@]} files)"
