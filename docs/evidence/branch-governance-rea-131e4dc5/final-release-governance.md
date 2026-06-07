# EVID-branch-governance-rea-131e4dc5-20260607-v0.6.0

## Scope

- Goal: branch-governance-rea-131e4dc5
- Release: v0.6.0
- Worktree: /home/xlib-standard/.worktree/workspaces/governance-release-v060-final-20260607
- Branch: governance/release-v0.6.0-final-20260607
- Baseline main/origin/main: cc184e85991e53c79b3d94b05ad024f61583ccb1
- Requested mode: unattended branch governance, push, and release
- Timestamp: 2026-06-07T01:27:01Z

## Protected State

- Main worktree before release metadata: clean.
- Local and remote non-main branches before release metadata: none, except the isolated release worktree branch created for AGENTS.md no-main-development compliance.
- Existing stash entries were preserved and not dropped.
- Prior branch backup bundle: .git/branch-governance-backups/20260606T231903Z/branches/codex__v060-docs-analysis/codex-v060-docs-analysis.bundle
- Prior branch patch backup: .git/branch-governance-backups/20260606T231903Z/branches/codex__v060-docs-analysis/codex-v060-docs-analysis.patch
- Prior backup ref: refs/backup/branch-governance-rea-131e4dc5/worker-2-906a545

## Release Metadata Changes

- Current release version is synchronized to v0.6.0.
- Current release facts retain the governed baseline commit cc184e85991e53c79b3d94b05ad024f61583ccb1 as the source anchor for branch governance completion.
- Changelog contains a v0.6.0 entry dated 2026-06-07.
- Release preflight examples and harness entry use VERSION=v0.6.0.

## Validation Results

- PASS: `GOWORK=off make docs-check`
- PASS: `GOWORK=off make rules-verify`
- PASS: `GOWORK=off make fmt`
- PASS: `GOWORK=off make vet`
- PASS: `GOWORK=off make test`
- PASS: `GOWORK=off make release-check`
- PASS: release evidence hash `18461abe3c7d794d84876cd7fd734c7ff0f26937074de71e1ee06355dca9c6bb`
- NOTE: `GOWORK=off make release-check` reported `govulncheck suspended`; secret, debt, boundary, contracts, evidence, integration, governance, and release checks passed.
- PENDING-POST-LAND: `XLIB_CONTEXT=release_verify GOWORK=off make release-preflight VERSION=v0.6.0` after main is fast-forwarded and pushed.

## Known Boundaries

- This evidence records pre-landing validation; post-landing release-preflight must run on clean main after origin/main receives the v0.6.0 metadata commit.
- Dependency vulnerability proof remains bounded unless `XLIB_ENABLE_VULNCHECK=1` is used; this run validated secret/security policy gates but did not run govulncheck.
- release/manifest/latest.json is generated and must not be manually edited or committed.
- The release tag should point at the final v0.6.0 metadata commit; the current_release.commit fact records the governed baseline source anchor.
