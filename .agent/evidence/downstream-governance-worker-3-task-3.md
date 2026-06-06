# Downstream Governance — Worker 3 Task 3 Verification Evidence

Team: `execute-home-xlib-sta-131e4dc5`
Worker: `worker-3`
Task: `3` — verification and evidence lane with focused tests and harness gate summaries
Timestamp: `2026-06-06T14:41:23+08:00`
Base HEAD: `1f3f873`

## Exact files changed by worker-3

- `.agent/evidence/downstream-governance-worker-3-task-3.md` — added this verification and harness-gate evidence report.

This report does not claim true downstream adoption. Per `.agent/evidence/README.md`, adoption evidence remains `adoption_claim=not_claimed`, `downstream_adoption_scope=local_contract_only`, `proof_based_adoption=false`, and `downstream_repo_write=false` unless accepted downstream-generated proof exists.

## Full-suite and diagnostics evidence

| Check | Command | Result |
| --- | --- | --- |
| Full test suite | `GOWORK=off go test ./...` | PASS (all packages ok; slowest observed package `internal/tools/releasemanifest` ok in `96.909s`) |
| Focused release/evidence package suite | `GOWORK=off go test ./internal/tools/releasemanifest ./scripts ./contracts ./cmd/goalcli` | PASS (`releasemanifest` ok in `117.032s`, `scripts` ok in `7.500s`, `contracts` ok in `0.011s`, `cmd/goalcli` ok in `1.620s`) |
| Focused downstream/evidence regex suite | `GOWORK=off go test ./... -run 'Test.*(Downstream\|Evidence\|Harness\|Release\|Manifest\|Contract\|Adoption\|Integration\|StandardImpact)'` | PASS (all selected packages ok) |
| Static diagnostics | `GOWORK=off go vet ./...` | PASS (exit `0`) |
| Contracts gate | `GOWORK=off make contracts` | PASS (`contract check passed`) |
| Evidence registry gate | `GOWORK=off make evidence-check` | PASS (`registry contract satisfied`) |
| Docs gate | `GOWORK=off make docs-check` | PASS (`docs-check passed`) |
| Lint | `GOWORK=off make lint` | PASS (`0 issues.`) |
| Build/typecheck | `GOWORK=off make build` | PASS (`go build ./...`) |
| Governance end-to-end | `GOWORK=off make governance-check` | PASS (main/worktree guards, evidence/adoption/boundary/debt/contracts/docs/registry/audit/traceability gates passed) |

## Focused downstream governance regression evidence

| Check | Command | Result |
| --- | --- | --- |
| goalcli downstream/evidence/adoption tests | `GOWORK=off go test ./cmd/goalcli -run 'TestRunDownstreamSyncPlan\|TestRunDispatchesDownstreamSyncPlan\|TestDownstream\|TestEvidenceReplay\|TestAdoption' -count=1` | PASS (`ok .../cmd/goalcli 0.020s`) |
| standard impact policy tests | `GOWORK=off go test ./scripts -run 'TestStandardImpact\|TestDependencyAutomationPolicy' -count=1` | PASS (`ok .../scripts 0.529s`) |
| evidence and adoption proof contract tests | `GOWORK=off go test ./contracts -run 'TestDownstreamAdoptionProofContractRequiredFields\|TestExecutionEvidenceContractRequiredFields' -count=1` | PASS (`ok .../contracts 0.004s`) |
| downstream section/debt policy tests | `GOWORK=off go test ./internal/debtcheck -run 'TestRunFailsDownstreamSection\|TestRunPassesDownstreamSection' -count=1` | PASS (`ok .../internal/debtcheck 0.005s`) |
| release manifest downstream/rules tests | `GOWORK=off go test ./internal/tools/releasemanifest -run 'Test.*StandardImpact\|Test.*Downstream.*Decision\|Test.*RepositoryRules.*Decision' -count=1` | PASS (`ok .../internal/tools/releasemanifest 20.052s`) |
| direct x.go dependency smoke | `GOWORK=off go list -deps ./... \| rg 'github.com/(bytechainx\|ZoneCNH)/x\.go'` | PASS: no forbidden `x.go` dependency matched |

## Release evidence harness sequence

Initial release evidence validation correctly failed before generated release artifacts were present:

- `GOWORK=off make release-evidence-check` → FAIL (`open release/manifest/latest.json: no such file or directory`).

After generating reproducible evidence outputs, the release evidence gates passed:

| Check | Command | Result |
| --- | --- | --- |
| Generate release evidence and verify manifest | `CHECK_STATUS=passed GOWORK=off make evidence && GOWORK=off make release-evidence-check` | PASS (`release/manifest/latest.json` generated and accepted) |
| Missing checksum guard | `GOWORK=off make release-evidence-checksum-check` before hash generation | FAIL as expected (`release evidence checksum is missing or empty`) |
| Generate checksum and verify it | `GOWORK=off make release-evidence-hash && GOWORK=off make release-evidence-checksum-check` | PASS (`release/manifest/latest.json: OK`) |

The generated `release/manifest/latest.json` and `.sha256` artifacts are reproducible harness outputs and remained untracked/ignored in this worker worktree.

## Harness gate summary

`GOWORK=off make governance-check` exercised and passed the aggregate governance path:

- main guard and worktree guard
- evidence check and adoption check
- boundary, architecture, domain, security, and debt checks
- contracts and docs checks
- CLI contract, issue registry, command registry, makefile baseline
- audit-goal, rules consistency, debt, and traceability gates

## Delegation evidence and integrated findings

Subagent spawn evidence: 1, child task "Test probe" (`019e9ba2-4f62-75e0-94ab-b99cd4b1d619`; child reported sub-run `019e9ba2-ec1f-7dd3-b269-59ec9e1b7035`), integrated findings that existing coverage spans downstream sync planning, evidence replay, adoption proof contracts, standard-impact matching, debt checks, and release-manifest decisions; focused verification should include goalcli/scripts/contracts/debtcheck/releasemanifest package tests; standalone x.go no-dependency smoke was missing and was added to this evidence run.

Serial searches before spawn: `0`.

## Residual gaps and risks

- No external downstream repository proof was provided or generated by this task; downstream adoption remains local-contract-only and not claimed.
- `x.go` remains consumer-review-only/no-write by policy and focused coverage, but the subagent found no single assertion locking the full role/action/status/reason tuple together.
- The subagent found no explicit negative test proving downstream sync plans never generate patch/render/write commands for `x.go`.
- The subagent found no single end-to-end test for `standard-impact-check -> downstream-sync-plan -> release manifest` asserting `x.go` remains consumer-only and `adoption_claim=not_claimed`.
- `GOWORK=off make governance-check` reported `govulncheck suspended`; vulnerability scanning still requires `XLIB_ENABLE_VULNCHECK=1` or the configured force/window conditions.
- Traceability passed under the current contract while reporting requirement lifecycle depth gaps (`partial_implemented`, `proof_depth=file_exists`, `full_lifecycle_graph=gap`).
