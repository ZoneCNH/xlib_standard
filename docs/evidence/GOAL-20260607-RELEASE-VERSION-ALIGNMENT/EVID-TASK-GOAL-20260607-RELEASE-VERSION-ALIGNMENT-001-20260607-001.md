# EVID-TASK-GOAL-20260607-RELEASE-VERSION-ALIGNMENT-001-20260607-001

Goal: GOAL-20260607-RELEASE-VERSION-ALIGNMENT
Task: release-version-align-131e4dc5
Date: 2026-06-07
Worktree: /home/xlib-standard/.worktree/workspaces/release-version-alignment-main
Branch: codex/release-version-alignment
HEAD at evidence time: 2c965ca
Tree state: dirty working tree with scoped release-version alignment and render-check propagation changes; no final commit created by leader in this evidence step.

## Summary

Aligned active release-version consumers from released `v0.6.1` to next locally untagged `v0.6.6`, added regression coverage so future active release-version drift is detected across facts, template constants, docs, harness commands, and release command examples, and repaired the missing `render-check` Makefile/registry/governance propagation gap found during follow-up review.

## Acceptance Evidence

- `git tag --list 'v0.6.6'` returned no output, so `v0.6.6` is locally untagged.
- `GOWORK=off go run ./cmd/goalcli fact audit --strict` passed with `current_release.version=v0.6.6`.
- `GOWORK=off go run ./cmd/goalcli version --json` reported details including `xlib-standard release v0.6.6`.
- `GOWORK=off go test ./cmd/goalcli ./internal/xlibfacts ./pkg/templatex ./scripts -count=1` passed.
- `GOWORK=off go test ./... -count=1` passed.
- `GOWORK=off make release-check` passed.
- `GOWORK=off make docs-check` passed after this evidence file was checked.
- `GOWORK=off make evidence-check` passed after this evidence file was checked.
- `git diff --check` passed after this evidence file was checked.
- `make -n render-check` was recorded only as dry-run target-shape evidence: it showed the new `render-check` target validates required `RENDER_CHECK_*` inputs before dispatching `goalcli render-check`, but it is not a runnable standalone proof.
- `GOWORK=off RENDER_CHECK_DIR=/tmp/tmp.YDQxYrypq2/rendered-kernel RENDER_CHECK_MODULE_NAME=kernel RENDER_CHECK_MODULE_PATH=github.com/ZoneCNH/kernel RENDER_CHECK_PACKAGE_NAME=kernel make render-check` passed against a temporary rendered fixture, proving the live parameterized Make target path.
- `GOWORK=off make makefile-baseline` passed after adding `render-check` to the Makefile baseline.
- `GOWORK=off make integration` passed and executed rendered-template checks for `kernel`, `configx`, and `redisx`.
- Stale-version scan outside historical/evidence paths found only the intentional `v0.6.1` regression sentinel in `cmd/goalcli/main_test.go`.

## Team Runtime Evidence

- `omx team status release-version-align-131e4dc5 --json --tail-lines 120` reported phase `complete`, tasks total `2`, completed `2`, failed `0`, pending `0`, blocked `0`, in_progress `0`, workers total `2`, dead `0`, non_reporting `0`.
- Task 1 result: worker-1 completed strict fact audit guard work and reported verification for goalcli, xlibfacts, full `go test ./...`, lint, and strict fact audit local-tag failure coverage.
- Task 2 result: worker-1 completed release consumer alignment; worker-2 remained blocked from edits by claim conflict and contributed read-only audit reports only.
- Leader reconciliation: actual local diff excludes generated release manifest artifacts.

## Follow-up Agent Team Evidence

- The persistent `release-version-align-131e4dc5` OMX team state was not available in the active repair worktree during follow-up orchestration, so the leader used bounded native agent lanes for the requested team execution.
- Fermat completed a read-only `render-check` propagation audit and confirmed the repair surface spans `Makefile`, `.agent/registries/makefile-target-registry.yaml`, `.agent/registries/makefile-baseline.yaml`, `cmd/goalcli/governance.go`, `cmd/goalcli/main.go`, `cmd/goalcli/main_test.go`, and `scripts/check_docs.sh`.
- Peirce completed a read-only stale-version audit and found no active current-fact anchors to `v0.6.1`; the remaining `v0.6.1` references are the intentional stale-output sentinel in `cmd/goalcli/main_test.go` plus historical/evidence/changelog text.
- Team recommendation accepted: keep the fix narrow, update this evidence file, and treat a future stale-version denylist gate as follow-up rather than expanding this repair.

## Commands

Passed:

- `git branch --show-current`
- `git status --short --branch`
- `git worktree list`
- `git rev-parse --short HEAD`
- `git log --oneline -n 8`
- `git tag --list 'v0.6.6'`
- `rg -n "v0\\.6\\.1"`
- `gofmt -w internal/tools/releasemanifest/main.go cmd/goalcli/main_test.go`
- `git diff --check`
- `GOWORK=off go run ./cmd/goalcli fact audit --strict`
- `GOWORK=off go run ./cmd/goalcli version --json`
- `GOWORK=off go test ./cmd/goalcli ./internal/xlibfacts ./pkg/templatex ./scripts -count=1`
- `GOWORK=off go test ./... -count=1`
- `GOWORK=off make lint`
- `GOWORK=off make docs-check`
- `GOWORK=off make rules-verify`
- `XLIB_CONTEXT=local_write GOWORK=off make governance-check`
- `GOWORK=off make standard-impact-check`
- `GOWORK=off make integration`
- `GOWORK=off make fmt`
- `GOWORK=off make vet`
- `GOWORK=off make release-check`
- `GOWORK=off make evidence-check`
- `gofmt -w cmd/goalcli/governance.go cmd/goalcli/main.go cmd/goalcli/main_test.go`
- `make -n render-check` (dry-run target-shape evidence only; runnable `make render-check` requires explicit `RENDER_CHECK_DIR`, `RENDER_CHECK_MODULE_NAME`, `RENDER_CHECK_MODULE_PATH`, and `RENDER_CHECK_PACKAGE_NAME` values)
- `GOWORK=off RENDER_CHECK_DIR=/tmp/tmp.YDQxYrypq2/rendered-kernel RENDER_CHECK_MODULE_NAME=kernel RENDER_CHECK_MODULE_PATH=github.com/ZoneCNH/kernel RENDER_CHECK_PACKAGE_NAME=kernel make render-check`
- `GOWORK=off go test ./cmd/goalcli -run 'TestRunDispatchesExternalCommands|TestMakefileBaseline|TestFactStrictProjectionGaps|TestFactAuditStrictPassesCanonicalFacts|TestVersionConstantsTrackChangelogRelease|TestVersionCommandReportsCurrentReleaseVersion'`
- `GOWORK=off go test ./cmd/goalcli`
- `GOWORK=off make fact-audit`
- `GOWORK=off make makefile-baseline`
- `GOWORK=off make test`
- `rg -n --hidden --glob '!.git/**' --glob '!CHANGELOG.md' --glob '!docs/evidence/**' --glob '!docs/v0.6.0/**' --glob '!.agent/evidence/**' 'v0\.6\.1|0\.6\.1'`
- `git ls-files release/manifest/latest.json release/manifest/latest.json.sha256 release/manifest/template.json`
- `git check-ignore -v release/manifest/latest.json release/manifest/latest.json.sha256`

Parameterized `render-check` proof output:

```text
fixture=/tmp/tmp.YDQxYrypq2/rendered-kernel
go run ./cmd/goalcli render-check "/tmp/tmp.YDQxYrypq2/rendered-kernel" "kernel" "github.com/ZoneCNH/kernel" "kernel"
rendered template check passed: kernel
```

Failed or unavailable:

- `GOWORK=off make render-check` without `RENDER_CHECK_*` values is intentionally not a standalone gate; the target requires explicit rendered directory/module inputs. Functional rendered-template coverage was verified through the live parameterized temporary fixture proof above and `GOWORK=off make integration`.
- `XLIB_CONTEXT=release_verify GOWORK=off make release-preflight VERSION=v0.6.6` failed because release preflight must run on `main`; current branch is `codex/release-version-alignment`.

## Changed Files

- `.agent/harness/harness.yaml`
- `.agent/registries/makefile-baseline.yaml`
- `.agent/registries/makefile-target-registry.yaml`
- `.agent/release/release-required-gates.yaml`
- `.xlib/facts/xlib.yaml`
- `AGENTS.md`
- `CHANGELOG.md`
- `Makefile`
- `README.md`
- `cmd/goalcli/governance.go`
- `cmd/goalcli/main.go`
- `cmd/goalcli/main_test.go`
- `docs/generation.md`
- `docs/release.md`
- `internal/tools/releasemanifest/main.go`
- `internal/xlibfacts/facts.go`
- `internal/xlibfacts/facts_test.go`
- `pkg/templatex/version.go`
- `scripts/check_docs.sh`
- `scripts/render_template_test.go`
- `docs/evidence/GOAL-20260607-RELEASE-VERSION-ALIGNMENT/EVID-TASK-GOAL-20260607-RELEASE-VERSION-ALIGNMENT-001-20260607-001.md`

## Known Proof Boundaries and Risks

- `release-preflight` was not passable on this branch because the gate requires `main`.
- `render-check` now exists as a Make target, but direct use requires explicit `RENDER_CHECK_DIR`, `RENDER_CHECK_MODULE_NAME`, `RENDER_CHECK_MODULE_PATH`, and `RENDER_CHECK_PACKAGE_NAME` values; the live temporary fixture proof validates the parameterized path and integration remains the higher-level reproducible proof path.
- There is no repo-wide stale-version denylist gate yet; the ad-hoc scan found only the intentional test sentinel outside historical/evidence docs.
- `make security` evidence inside `release-check` covers secret scan; `govulncheck` was suspended unless `XLIB_ENABLE_VULNCHECK=1` is set, so dependency vulnerability safety is not claimed.
- Traceability evidence remains partial: `traceability_status=partial_implemented`, `proof_depth=file_exists`, `proof_depth_level=D3`; this does not prove a complete lifecycle graph.
- Downstream adoption proof is local template/integration proof only; no external downstream CI or repository adoption proof was produced.
- Working tree remained dirty at evidence time; final commit and PR publication were not performed in this step.

## Follow-up

- Run release preflight from `main` or the release-authorized branch when the release process requires it.
- Add a stale-version denylist gate if future releases require automated enforcement beyond the current targeted regression tests and ad-hoc scan.
- Run `XLIB_ENABLE_VULNCHECK=1 GOWORK=off make security` if vulnerability-scan evidence is required for release.
