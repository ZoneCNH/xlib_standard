// SPDX-License-Identifier: Apache-2.0
package main

import (
	"strings"
	"testing"
)

func TestValidateDockerEvidenceReportsMissingRequiredFields(t *testing.T) {
	failures := validateDockerEvidence(DockerEvidence{}, false)

	for _, want := range []string{
		"docker.contract_version is required",
		"docker.go_version is required",
		"docker.golangci_lint_version is required",
		"docker.govulncheck_version is required",
		"docker.enabled must be true",
		"docker.buildkit_required must be true",
		"docker.cache_mounts must include go-build",
		"docker.cache_mounts must include go-mod",
		"docker.cache_mounts must include golangci-lint",
		"docker.validated_by must include " + dockerEvidenceValidators[0],
		"docker.base_image is required",
		"docker.toolchain_image is required",
		"docker.runtime_image is required",
		"docker.workflow_run_id is required",
		"docker.artifact_name is required",
		"docker.artifact_url is required",
	} {
		requireTestFailure(t, failures, want)
	}
}

func TestValidateDockerEvidenceReportsInvalidContractAndDigestFields(t *testing.T) {
	evidence := validDockerEvidenceForValidationTest()
	evidence.ContractVersion = "docker-toolchain/v1"
	evidence.BaseImageDigest = "not-sha"
	evidence.ToolchainImageDigest = ""
	evidence.RuntimeImageDigest = "digest:runtime"

	failures := validateDockerEvidence(evidence, true)

	for _, want := range []string{
		`docker.contract_version must be docker-toolchain/v2, got "docker-toolchain/v1"`,
		"docker.base_image_digest must start with sha256:",
		"docker.toolchain_image_digest is required",
		"docker.runtime_image_digest must start with sha256:",
	} {
		requireTestFailure(t, failures, want)
	}
}

func TestValidateDownstreamAdoptionEvidenceReportsRequiredAndForbiddenFields(t *testing.T) {
	failures := validateDownstreamAdoptionEvidence(DownstreamAdoptionEvidence{})
	for _, want := range []string{
		"downstream_adoption.adoption_claim is required",
		"downstream_adoption.downstream_adoption_scope is required",
		"downstream_adoption.source is required",
	} {
		requireTestFailure(t, failures, want)
	}

	evidence := DownstreamAdoptionEvidence{
		AdoptionClaim:           "not_claimed",
		DownstreamAdoptionScope: "downstream_generated",
		ProofBasedAdoption:      true,
		DownstreamRepoWrite:     true,
		Source:                  "test",
	}

	failures = validateDownstreamAdoptionEvidence(evidence)
	for _, want := range []string{
		"downstream_adoption.downstream_repo_write must be false",
		`downstream_adoption.downstream_adoption_scope must be local_contract_only when adoption_claim is not_claimed, got "downstream_generated"`,
		"downstream_adoption.proof_based_adoption must be false when adoption_claim is not_claimed",
		"downstream adoption claims require downstream-generated proof and accepted ledger evidence",
	} {
		requireTestFailure(t, failures, want)
	}
}

func TestValidateDownstreamAdoptionEvidenceRequiresBothProofAndLedger(t *testing.T) {
	cases := []struct {
		name        string
		evidence    DownstreamAdoptionEvidence
		wantFailure bool
	}{
		{
			name: "claim with proof missing ledger",
			evidence: DownstreamAdoptionEvidence{
				AdoptionClaim:           "adopted",
				DownstreamAdoptionScope: "local_contract_only",
				ProofArtifactPath:       "release/evidence/downstream-proof.json",
				Source:                  "test",
			},
			wantFailure: true,
		},
		{
			name: "scope with ledger missing proof",
			evidence: DownstreamAdoptionEvidence{
				AdoptionClaim:              "adopted",
				DownstreamAdoptionScope:    "downstream_generated",
				AcceptedLedgerEvidencePath: "release/evidence/accepted-ledger.json",
				Source:                     "test",
			},
			wantFailure: true,
		},
		{
			name: "proof based adoption has proof and ledger",
			evidence: DownstreamAdoptionEvidence{
				AdoptionClaim:              "adopted",
				DownstreamAdoptionScope:    "downstream_generated",
				ProofBasedAdoption:         true,
				ProofArtifactPath:          "release/evidence/downstream-proof.json",
				AcceptedLedgerEvidencePath: "release/evidence/accepted-ledger.json",
				Source:                     "test",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			failures := validateDownstreamAdoptionEvidence(tc.evidence)
			gotFailure := testFailuresContain(failures, "downstream adoption claims require downstream-generated proof and accepted ledger evidence")
			if gotFailure != tc.wantFailure {
				t.Fatalf("proof/ledger failure = %v, want %v; failures = %v", gotFailure, tc.wantFailure, failures)
			}
		})
	}
}

func TestRequireEnumValueReportsRequiredValidAndInvalid(t *testing.T) {
	allowed := []string{"required", "not_required"}

	var failures []string
	requireEnumValue(&failures, "release_decision", "", allowed)
	requireTestFailure(t, failures, "release_decision is required")

	failures = nil
	requireEnumValue(&failures, "release_decision", "required", allowed)
	if len(failures) != 0 {
		t.Fatalf("valid enum produced failures: %v", failures)
	}

	requireEnumValue(&failures, "release_decision", "skipped", allowed)
	requireTestFailure(t, failures, `release_decision must be one of required, not_required, got "skipped"`)
}

func validDockerEvidenceForValidationTest() DockerEvidence {
	return DockerEvidence{
		Enabled:              true,
		ContractVersion:      "docker-toolchain/v2",
		GoVersion:            "1.23",
		GolangCILintVersion:  "golangci-lint v2.1.6",
		GovulncheckVersion:   "govulncheck v1.1.4",
		BuildKitRequired:     true,
		CacheMounts:          []string{"go-build", "go-mod", "golangci-lint"},
		BaseImage:            "golang:1.23-bookworm",
		BaseImageDigest:      "sha256:base",
		ToolchainImage:       "xlib-standard-toolchain:local",
		ToolchainImageDigest: "sha256:toolchain",
		RuntimeImage:         "xlib-standard-goalcli-runtime:local",
		RuntimeImageDigest:   "sha256:runtime",
		ValidatedBy:          append([]string(nil), dockerEvidenceValidators...),
		WorkflowRunID:        "123",
		ArtifactName:         "release-manifest-123",
		ArtifactURL:          "https://example.test/actions/runs/123",
	}
}

func requireTestFailure(t *testing.T, failures []string, want string) {
	t.Helper()

	if !testFailuresContain(failures, want) {
		t.Fatalf("failures = %v, want substring %q", failures, want)
	}
}

func testFailuresContain(failures []string, want string) bool {
	for _, failure := range failures {
		if strings.Contains(failure, want) {
			return true
		}
	}
	return false
}
