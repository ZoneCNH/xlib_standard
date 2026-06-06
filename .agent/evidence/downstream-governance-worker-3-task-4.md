# Downstream Governance Worker 3 Task 4 Evidence

## Task

Task 4: docs/contract/release-manifest integration for downstream governance.

## Boundaries Preserved

- Downstream adoption remains `not_claimed` unless downstream-generated proof and accepted ledger evidence are both present.
- Release manifest evidence is local contract evidence only (`local_contract_only`).
- `x.go` remains consumer-review-only; no downstream patch/render/write command targets it.
- No new dependencies were added.

## Integrated Subagent Findings

- Review probe `019e9bab-0e29-76f1-b6a2-5c47081ec0e6`: release docs and manifest must avoid treating registered downstreams or local plans as adoption truth.
- Test probe `019e9bab-0fff-71a2-8be9-1ccec2c11b8a`: add regression coverage for release manifest defaults/validation and x.go consumer-only invariants.
- Change-slice probe `019e9bab-11ee-7ad2-857e-966091a6b68a`: fingerprint the downstream adoption proof contract in the release manifest and keep adoption fields defaulted to local evidence.

## Changed Surfaces

- `internal/tools/releasemanifest/{types.go,docker.go,util.go,vars.go,verify.go,main_test.go}`: downstream adoption evidence field, proof-contract fingerprint, validation, and regression tests.
- `cmd/goalcli/downstream_sync_plan_test.go`: regression checks that `x.go` is consumer-only and truth/adoption paths are not write targets.
- `release/manifest/template.json`: downstream adoption proof contract fingerprint placeholder and local-only adoption evidence defaults.
- `docs/standard/evidence-protocol.md`: downstream adoption boundary and release manifest requirements.
- `docs/standard/downstream-registry.md`: registry/adoption truth boundary and x.go consumer-only rule.
- `.agent/release/release.md`: release evidence explicitly records no downstream adoption claim.
- `.agent/index.yaml`: indexes worker evidence artifacts for governance checks.

## Verification Evidence

Focused checks before index fix:

- PASS: `GOWORK=off go test ./internal/tools/releasemanifest -run 'Test(RunCLIGeneratesManifestToOut|RunCLIVerifyReportsDrift|VerifyManifestAcceptsFreshManifestAndRejectsDrift|VerifyManifestRejectsCorruptedManifestFields|BuildDownstreamAdoptionEvidenceDefaultsToNotClaimed|ValidateDownstreamAdoptionEvidenceRequiresProofAndLedgerForClaims|BuildManifestRecordsFixtureRepositoryFacts)' -count=1`
- PASS: `GOWORK=off go test ./cmd/goalcli -run 'TestRunDownstreamSyncPlan' -count=1`
- PASS: `GOWORK=off go test ./contracts -run 'TestDownstreamAdoptionProofContractRequiredFields|TestExecutionEvidenceContractMatchesEvidenceManifest' -count=1`
- FAIL, fixed here: `GOWORK=off go test ./internal/tools/releasemanifest ./cmd/goalcli ./contracts -count=1` initially reported `.agent/index.yaml missing file entry .agent/evidence/downstream-governance-worker-3-task-3.md`.

Final checks are recorded in the Task 4 completion transition after re-running the full verification sequence.

## Residual Gaps

- No external downstream repository was modified or validated.
- No downstream-generated adoption proof artifact or accepted ledger evidence was supplied, so adoption remains `not_claimed`.
