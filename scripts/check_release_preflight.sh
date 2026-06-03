#!/usr/bin/env bash
set -euo pipefail

version="${1:-${VERSION:-}}"

if [[ -z "$version" ]]; then
  echo "ERROR: set VERSION=vX.Y.Z when running release preflight"
  exit 1
fi

if [[ ! "$version" =~ ^v[0-9]+\.[0-9]+\.[0-9]+([-+][0-9A-Za-z.-]+)?$ ]]; then
  echo "ERROR: release version must look like vX.Y.Z: $version"
  exit 1
fi

if ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  echo "ERROR: release preflight must run inside a git worktree"
  exit 1
fi

branch="$(git rev-parse --abbrev-ref HEAD)"
if [[ "$branch" != "main" ]]; then
  echo "ERROR: release preflight must run on main; current branch is $branch"
  exit 1
fi

if [[ -n "$(git status --porcelain)" ]]; then
  echo "ERROR: release preflight requires a clean git worktree"
  git status --short
  exit 1
fi

git fetch --quiet origin main --tags

head_sha="$(git rev-parse HEAD)"
origin_main_sha="$(git rev-parse origin/main)"
if [[ "$head_sha" != "$origin_main_sha" ]]; then
  echo "ERROR: local main is not aligned with origin/main"
  echo "HEAD=$head_sha"
  echo "origin/main=$origin_main_sha"
  exit 1
fi

if git rev-parse -q --verify "refs/tags/$version" >/dev/null; then
  echo "ERROR: local tag already exists: $version"
  exit 1
fi

if git ls-remote --exit-code --tags origin "refs/tags/$version" >/dev/null 2>&1; then
  echo "ERROR: remote tag already exists: $version"
  exit 1
fi

if ! grep -Eq "^## \\[?$version\\]?( |$)" CHANGELOG.md; then
  echo "ERROR: CHANGELOG.md must contain a release heading for $version"
  exit 1
fi

<<<<<<< HEAD
if ! command -v golangci-lint >/dev/null 2>&1; then
  echo "ERROR: golangci-lint not installed"
  exit 1
fi

if [[ "${XLIB_ENABLE_VULNCHECK:-0}" == "1" ]] && ! command -v govulncheck >/dev/null 2>&1; then
  echo "ERROR: govulncheck not installed (required when XLIB_ENABLE_VULNCHECK=1)"
  exit 1
fi
=======
required_tools=(golangci-lint)
if [[ "${XLIB_ENABLE_VULNCHECK:-}" == "1" ]]; then
  required_tools+=(govulncheck)
fi

for tool in "${required_tools[@]}"; do
  if ! command -v "$tool" >/dev/null 2>&1; then
    echo "ERROR: $tool not installed"
    exit 1
  fi
done
>>>>>>> f52a64f (omx(team): auto-checkpoint worker-1 [1])

echo "release preflight metadata checks passed for $version"
