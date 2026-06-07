# EVID-TASK-GOAL-20260607-DOCS-ALIGNMENT-20260607-001

## Scope

- Goal: GOAL-20260607-DOCS-ALIGNMENT
- Task: TASK-GOAL-20260607-DOCS-ALIGNMENT-001
- Worktree: `/home/xlib-standard/.worktree/workspaces/update-alignment-docs`
- Branch: `codex/update-alignment-docs`
- Base commit: `216ef50cead9ab20437566845b3446d6dbd07ec9`
- Release anchor: `v0.6.1`
- Change class: docs / release metadata alignment

## Release Evidence

- `gh release view v0.6.1 --json tagName,targetCommitish,isDraft,isPrerelease,publishedAt,url`
- Result: `v0.6.1` published at `2026-06-07T05:33:38Z`, target `main`, draft `false`, prerelease `false`.
- URL: `https://github.com/ZoneCNH/xlib-standard/releases/tag/v0.6.1`

## Changed Files

- `AGENTS.md`
- `CHANGELOG.md`
- `README.md`
- `.agent/harness/harness.yaml`
- `.xlib/facts/xlib.yaml`
- `cmd/goalcli/main_test.go`
- `docs/generation.md`
- `docs/release.md`
- `docs/evidence/GOAL-20260607-DOCS-ALIGNMENT/EVID-TASK-GOAL-20260607-DOCS-ALIGNMENT-20260607-001.md`
- `internal/tools/releasemanifest/main.go`
- `internal/xlibfacts/facts.go`
- `internal/xlibfacts/facts_test.go`
- `pkg/templatex/version.go`
- `release/manifest/template.json`
- `scripts/render_template_test.go`

## Acceptance Criteria

- Release-facing docs reference the published `v0.6.1` baseline.
- The changelog preserves existing entries and adds a focused `v0.6.1` release note.
- Canonical release metadata, template version, manifest template, and harness release-preflight examples reference `v0.6.1`.
- Documentation, governance, and integration gates pass on the final branch state.

## Commands

- `git status --short --branch`
  - Result: work was performed on `codex/update-alignment-docs`, outside `main`.
- `git worktree list`
  - Result: docs alignment performed in `/home/xlib-standard/.worktree/workspaces/update-alignment-docs`, not on main.
- `GOWORK=off make fmt`
  - Result: passed; `go fmt ./...`.
- `GOWORK=off go test ./cmd/goalcli ./internal/xlibfacts ./scripts`
  - Result: first run failed because `.xlib/facts/xlib.yaml` still advertised `v0.6.0`; after aligning canonical facts to `v0.6.1`, rerun passed for all three packages.
- `git diff --check`
  - Result: passed with no output.
- `GOWORK=off make docs-check`
  - Result: passed; `docs-check passed`.
- `GOWORK=off make rules-verify`
  - Result: passed; `rules total: 419`, `rules active: 388`, and all active rules have valid `enforced_by` commands.
- `XLIB_CONTEXT=local_write GOWORK=off make governance-check`
  - Result: passed; main/worktree guard, evidence, adoption, boundary, contracts, docs, CLI contract, registries, rules, debt, security secret scan, and traceability gates passed.
- `GOWORK=off make integration`
  - Result: passed; rendered and checked kernel, configx, and redisx downstream fixtures.
- `GOWORK=off make release-check`
  - Result: passed; included fmt, vet, full test suite, race tests, release evidence generation/checksum verification, `fact audit --strict`, dependency check, score threshold, integration, docs debt, and governance dry-run contract checks.

## Risks and Gaps

- This is a post-release alignment. It does not rerun `release-final-check` and does not republish the GitHub Release.
- `release/manifest/template.json` is the editable manifest template, not the generated `release/manifest/latest.json` artifact.
- `make security` passed its secret scan path, but `govulncheck` was suspended by configuration; set `XLIB_ENABLE_VULNCHECK=1` to force vulnerability scanning evidence.
- Existing proof boundaries still apply: traceability proof depth is file-existence level where marked partial, and vulnerability scanning evidence depends on the configured `govulncheck` gate window.
