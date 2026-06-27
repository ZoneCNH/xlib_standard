#!/usr/bin/env bash
set -euo pipefail

if [[ $# -lt 1 ]]; then
  echo "Usage: scripts/docker/docker_gate.sh <make-target> [extra make args...]" >&2
  exit 2
fi

target="$1"
shift

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
default_image_name="$(
  basename "$repo_root" |
    tr '[:upper:]' '[:lower:]' |
    sed -E 's/[^a-z0-9]+/-/g; s/^-+//; s/-+$//'
)"
if [[ -z "$default_image_name" ]]; then
  default_image_name="xlib-standard"
fi
image="${DOCKER_IMAGE:-${default_image_name}-toolchain:local}"
go_version="${GO_VERSION:-1.23}"
golangci_lint_version="${GOLANGCI_LINT_VERSION:-v2.1.6}"
govulncheck_version="${GOVULNCHECK_VERSION:-v1.1.4}"
git_mount_args=()

append_git_mount() {
  local path="$1"
  local mount

  if [[ -z "$path" || "$path" != /* || ! -e "$path" || "$path" == "$repo_root/.git" ]]; then
    return
  fi
  mount="$path:$path"
  for existing in "${git_mount_args[@]}"; do
    if [[ "$existing" == "$mount" ]]; then
      return
    fi
  done
  git_mount_args+=(--volume "$mount")
}

if git_dir="$(git -C "$repo_root" rev-parse --git-dir 2>/dev/null)" &&
  git_common_dir="$(git -C "$repo_root" rev-parse --git-common-dir 2>/dev/null)"; then
  if [[ "$git_dir" != /* ]]; then
    git_dir="$(cd "$repo_root" && cd "$git_dir" && pwd)"
  fi
  if [[ "$git_common_dir" != /* ]]; then
    git_common_dir="$(cd "$repo_root" && cd "$git_common_dir" && pwd)"
  fi

  if [[ "$git_dir" == "$git_common_dir"/* ]]; then
    append_git_mount "$git_common_dir"
  else
    append_git_mount "$git_common_dir"
    append_git_mount "$git_dir"
  fi
fi

export DOCKER_BUILDKIT=1

"$repo_root/scripts/docker/check_toolchain.sh"

docker buildx build \
  --load \
  --target toolchain \
  --build-arg "GO_VERSION=$go_version" \
  --build-arg "GOLANGCI_LINT_VERSION=$golangci_lint_version" \
  --build-arg "GOVULNCHECK_VERSION=$govulncheck_version" \
  --tag "$image" \
  "$repo_root"

docker run --rm \
  --workdir /workspace \
  --volume "$repo_root:/workspace" \
  "${git_mount_args[@]}" \
  --volume go-build-cache:/root/.cache/go-build \
  --volume go-mod-cache:/go/pkg/mod \
  --env "GOWORK=${GOWORK:-off}" \
  --env "XLIB_CONTEXT=${XLIB_CONTEXT:-ci_pull_request}" \
  --env "VERSION=${VERSION:-}" \
  --env "DOWNSTREAM=${DOWNSTREAM:-}" \
  --env "XLIB_ENABLE_VULNCHECK=${XLIB_ENABLE_VULNCHECK:-}" \
  --env "CI=${CI:-}" \
  --env "GITHUB_ACTIONS=${GITHUB_ACTIONS:-}" \
  --env "GIT_CONFIG_COUNT=1" \
  --env "GIT_CONFIG_KEY_0=safe.directory" \
  --env "GIT_CONFIG_VALUE_0=/workspace" \
  "$image" make "$target" "$@"
