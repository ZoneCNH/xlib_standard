// SPDX-License-Identifier: Apache-2.0
package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// --- envBool coverage ---

func TestEnvBoolFallsBackForUnrecognizedValue(t *testing.T) {
	cases := []struct {
		name     string
		value    string
		fallback bool
		want     bool
	}{
		{name: "empty uses fallback true", value: "", fallback: true, want: true},
		{name: "empty uses fallback false", value: "", fallback: false, want: false},
		{name: "unrecognized uses fallback true", value: "maybe", fallback: true, want: true},
		{name: "unrecognized uses fallback false", value: "maybe", fallback: false, want: false},
		{name: "explicit true variants", value: "YES", fallback: false, want: true},
		{name: "explicit on", value: "On", fallback: false, want: true},
		{name: "explicit false variants", value: "NO", fallback: true, want: false},
		{name: "explicit off", value: "Off", fallback: true, want: false},
		{name: "numeric one", value: "1", fallback: false, want: true},
		{name: "numeric zero", value: "0", fallback: true, want: false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("ENBOOL_COVERAGE", tc.value)
			if got := envBool("ENBOOL_COVERAGE", tc.fallback); got != tc.want {
				t.Fatalf("envBool = %v; want %v", got, tc.want)
			}
		})
	}
}

// --- envCSVDefault coverage ---

func TestEnvCSVDefaultFallsBackForWhitespaceOnlyInput(t *testing.T) {
	t.Setenv("CSVCOVERAGE", "  ,  ,  ")
	got := envCSVDefault("CSVCOVERAGE", []string{"a", "b"})
	if len(got) != 2 || got[0] != "a" || got[1] != "b" {
		t.Fatalf("envCSVDefault whitespace-only = %v; want fallback [a b]", got)
	}
}

func TestEnvCSVDefaultFallsBackForEmptyInput(t *testing.T) {
	t.Setenv("CSVCOVERAGE", "")
	got := envCSVDefault("CSVCOVERAGE", []string{"x"})
	if len(got) != 1 || got[0] != "x" {
		t.Fatalf("envCSVDefault empty = %v; want fallback [x]", got)
	}
}

func TestEnvCSVDefaultTrimsAndFiltersParts(t *testing.T) {
	t.Setenv("CSVCOVERAGE", " a , ,b,")
	got := envCSVDefault("CSVCOVERAGE", nil)
	if len(got) != 2 || got[0] != "a" || got[1] != "b" {
		t.Fatalf("envCSVDefault = %v; want [a b]", got)
	}
}

// --- requireEnumValue coverage ---

func TestRequireEnumValueRejectsUnknownValue(t *testing.T) {
	var failures []string
	requireEnumValue(&failures, "field", "bogus", []string{"a", "b"})
	if len(failures) != 1 {
		t.Fatalf("failures = %v; want one", failures)
	}
	if !strings.Contains(failures[0], "must be one of a, b") {
		t.Fatalf("failure = %q; want enum list", failures[0])
	}
}

func TestRequireEnumValueAcceptsKnownValue(t *testing.T) {
	var failures []string
	requireEnumValue(&failures, "field", "a", []string{"a", "b"})
	if len(failures) != 0 {
		t.Fatalf("failures = %v; want none", failures)
	}
}

func TestRequireEnumValueRequiresNonEmpty(t *testing.T) {
	var failures []string
	requireEnumValue(&failures, "field", "  ", []string{"a"})
	if len(failures) != 1 || !strings.Contains(failures[0], "is required") {
		t.Fatalf("failures = %v; want required", failures)
	}
}

// --- validateDockerEvidence coverage ---

func TestValidateDockerEvidenceReportsAllBrokenFields(t *testing.T) {
	evidence := DockerEvidence{
		Enabled:              false,
		ContractVersion:      "bogus/v1",
		GoVersion:            "",
		GolangCILintVersion:  "",
		GovulncheckVersion:   "",
		BuildKitRequired:     false,
		CacheMounts:          nil,
		ValidatedBy:          nil,
		BaseImage:            "",
		ToolchainImage:       "",
		RuntimeImage:         "",
		WorkflowRunID:        "",
		ArtifactName:         "",
		ArtifactURL:          "",
		BaseImageDigest:      "not-a-sha",
		ToolchainImageDigest: "sha256:valid",
		RuntimeImageDigest:   "",
	}

	failures := validateDockerEvidence(evidence, true)
	message := strings.Join(failures, "\n")
	for _, want := range []string{
		"docker.contract_version",
		"docker.go_version is required",
		"docker.golangci_lint_version is required",
		"docker.govulncheck_version is required",
		"docker.contract_version must be docker-toolchain/v2",
		"docker.enabled must be true",
		"docker.buildkit_required must be true",
		"docker.cache_mounts must include go-build",
		"docker.validated_by must include",
		"docker.base_image is required",
		"docker.base_image_digest must start with sha256:",
		"docker.runtime_image_digest is required",
	} {
		if !strings.Contains(message, want) {
			t.Fatalf("failures = %q; want substring %q", message, want)
		}
	}
}

func TestValidateDockerEvidenceAllowsValidDigestsWhenNotRequired(t *testing.T) {
	// requireDigests=false with non-sha digest still flags the prefix violation.
	evidence := validDockerEvidence()
	evidence.BaseImageDigest = "not-a-sha"
	failures := validateDockerEvidence(evidence, false)
	if !contains(failures, "docker.base_image_digest must start with sha256:") {
		t.Fatalf("failures = %v; want digest prefix violation when digest present", failures)
	}
}

func TestValidateDockerEvidenceSkipsMissingDigestsWhenNotRequired(t *testing.T) {
	evidence := validDockerEvidence()
	evidence.BaseImageDigest = ""
	evidence.ToolchainImageDigest = ""
	evidence.RuntimeImageDigest = ""
	failures := validateDockerEvidence(evidence, false)
	for _, f := range failures {
		if strings.Contains(f, "digest") {
			t.Fatalf("unexpected digest failure when not required: %q", f)
		}
	}
}

func validDockerEvidence() DockerEvidence {
	return DockerEvidence{
		Enabled:              true,
		ContractVersion:      "docker-toolchain/v2",
		GoVersion:            "1.23",
		GolangCILintVersion:  "golangci-lint v2.1.6",
		GovulncheckVersion:   "govulncheck v1.1.4",
		BuildKitRequired:     true,
		CacheMounts:          []string{"go-build", "go-mod", "golangci-lint"},
		BaseImage:            "golang:1.23",
		BaseImageDigest:      "sha256:abc",
		ToolchainImage:       "toolchain:local",
		ToolchainImageDigest: "sha256:def",
		RuntimeImage:         "runtime:local",
		RuntimeImageDigest:   "sha256:ghi",
		ValidatedBy:          append([]string(nil), dockerEvidenceValidators...),
		WorkflowRunID:        "run",
		ArtifactName:         "artifact",
		ArtifactURL:          "url",
	}
}

// --- validateDownstreamAdoptionEvidence coverage ---

func TestValidateDownstreamAdoptionEvidenceRejectsRepoWrite(t *testing.T) {
	evidence := buildDownstreamAdoptionEvidence()
	evidence.DownstreamRepoWrite = true
	failures := validateDownstreamAdoptionEvidence(evidence)
	if !contains(failures, "downstream_adoption.downstream_repo_write must be false") {
		t.Fatalf("failures = %v; want repo_write failure", failures)
	}
}

func TestValidateDownstreamAdoptionEvidenceRejectsScopeMismatchForNotClaimed(t *testing.T) {
	evidence := buildDownstreamAdoptionEvidence()
	evidence.DownstreamAdoptionScope = "downstream_generated"
	failures := validateDownstreamAdoptionEvidence(evidence)
	want := `downstream_adoption.downstream_adoption_scope must be local_contract_only when adoption_claim is not_claimed, got "downstream_generated"`
	if !contains(failures, want) {
		t.Fatalf("failures = %v; want scope mismatch failure %q", failures, want)
	}
}

func TestValidateDownstreamAdoptionEvidenceRejectsProofBasedForNotClaimed(t *testing.T) {
	evidence := buildDownstreamAdoptionEvidence()
	evidence.ProofBasedAdoption = true
	failures := validateDownstreamAdoptionEvidence(evidence)
	if !contains(failures, "downstream_adoption.proof_based_adoption must be false when adoption_claim is not_claimed") {
		t.Fatalf("failures = %v; want proof_based failure", failures)
	}
}

func TestValidateDownstreamAdoptionEvidenceAcceptsClaimWithProofAndLedger(t *testing.T) {
	evidence := buildDownstreamAdoptionEvidence()
	evidence.AdoptionClaim = "adopted"
	evidence.DownstreamAdoptionScope = "downstream_generated"
	evidence.ProofBasedAdoption = true
	evidence.ProofArtifactPath = "release/adoption/proof.json"
	evidence.AcceptedLedgerEvidencePath = "release/adoption/accepted.jsonl"
	failures := validateDownstreamAdoptionEvidence(evidence)
	for _, f := range failures {
		if strings.Contains(f, "require downstream-generated proof") {
			t.Fatalf("unexpected proof requirement failure when proof and ledger present: %q", f)
		}
	}
}

func TestValidateDownstreamAdoptionEvidenceRejectsClaimWithMissingProof(t *testing.T) {
	evidence := buildDownstreamAdoptionEvidence()
	evidence.AdoptionClaim = "adopted"
	evidence.DownstreamAdoptionScope = "local_contract_only"
	evidence.ProofArtifactPath = "release/adoption/proof.json"
	// AcceptedLedgerEvidencePath intentionally empty.
	failures := validateDownstreamAdoptionEvidence(evidence)
	if !contains(failures, "downstream adoption claims require downstream-generated proof and accepted ledger evidence") {
		t.Fatalf("failures = %v; want proof+ledger failure", failures)
	}
}

func TestValidateDownstreamAdoptionEvidenceRejectsProofScopeWithoutProof(t *testing.T) {
	evidence := buildDownstreamAdoptionEvidence()
	evidence.DownstreamAdoptionScope = "downstream_generated"
	// ProofBasedAdoption=false, AdoptionClaim=not_claimed, but scope implies adoption.
	failures := validateDownstreamAdoptionEvidence(evidence)
	if !contains(failures, "downstream adoption claims require downstream-generated proof and accepted ledger evidence") {
		t.Fatalf("failures = %v; want proof+ledger failure for proof scope", failures)
	}
}

func TestValidateDownstreamAdoptionEvidenceRequiresCoreFields(t *testing.T) {
	failures := validateDownstreamAdoptionEvidence(DownstreamAdoptionEvidence{})
	for _, want := range []string{
		"downstream_adoption.adoption_claim is required",
		"downstream_adoption.downstream_adoption_scope is required",
		"downstream_adoption.source is required",
	} {
		if !contains(failures, want) {
			t.Fatalf("failures = %v; want %q", failures, want)
		}
	}
}

// --- buildDebtEvidence + buildManifest coverage ---

// TestBuildDebtEvidenceReportsInvalidJSONCoverage covers the json.Unmarshal error
// branch in buildDebtEvidence.
func TestBuildDebtEvidenceReportsInvalidJSONCoverage(t *testing.T) {
	root := t.TempDir()
	reportPath := filepath.Join(root, filepath.FromSlash(debtReportPath))
	if err := os.MkdirAll(filepath.Dir(reportPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(reportPath, []byte("{not-json"), 0o644); err != nil {
		t.Fatal(err)
	}
	chdir(t, root)

	if _, err := buildDebtEvidence(); err == nil {
		t.Fatal("buildDebtEvidence returned nil error for invalid JSON")
	}
}

// TestBuildDebtEvidenceCarriesReportMinScore covers the
// `if report.MinScore != 0` branch.
func TestBuildDebtEvidenceCarriesReportMinScore(t *testing.T) {
	root := t.TempDir()
	reportPath := filepath.Join(root, filepath.FromSlash(debtReportPath))
	if err := os.MkdirAll(filepath.Dir(reportPath), 0o755); err != nil {
		t.Fatal(err)
	}
	// min_score set to a non-zero value different from the default; score from
	// sections list (no checks) to also exercise the CheckCount fallback.
	data := []byte(`{
  "status": "passed",
  "score": 9.9,
  "min_score": 9.5,
  "sections": [{"id": "only"}]
}` + "\n")
	if err := os.WriteFile(reportPath, data, 0o644); err != nil {
		t.Fatal(err)
	}
	chdir(t, root)

	got, err := buildDebtEvidence()
	if err != nil {
		t.Fatalf("buildDebtEvidence error = %v", err)
	}
	if got.MinScore != 9.5 {
		t.Fatalf("min_score = %.1f; want 9.5", got.MinScore)
	}
	if got.CheckCount != 1 {
		t.Fatalf("check_count = %d; want 1 (sections fallback)", got.CheckCount)
	}
}

// TestBuildManifestReportsDebtEvidenceFailure covers buildManifest line 307-309
// (debt evidence error propagation).
func TestBuildManifestReportsDebtEvidenceFailure(t *testing.T) {
	t.Setenv("GOWORK", "off")
	repo := releaseManifestFixtureRepo(t)
	// Replace debt report with invalid JSON.
	reportPath := filepath.Join(repo, filepath.FromSlash(debtReportPath))
	if err := os.WriteFile(reportPath, []byte("{not-json"), 0o644); err != nil {
		t.Fatal(err)
	}
	chdir(t, repo)

	if _, err := buildManifest(); err == nil {
		t.Fatal("buildManifest succeeded for invalid debt JSON, want error")
	}
}

// --- validateChecks coverage (debt.score < min + check_count == 0) ---

// TestVerifyManifestRejectsDebtScoreBelowMinAndZeroCheckCount covers verify.go
// lines 139-144 inside the requirePassed debt block.
func TestVerifyManifestRejectsDebtScoreBelowMinAndZeroCheckCount(t *testing.T) {
	t.Setenv("GOWORK", "off")
	t.Setenv("CHECK_STATUS", "passed")
	setDockerDigestEvidence(t)
	repo := releaseManifestFixtureRepo(t)
	writeStandardImpactReportFixture(t, repo)
	chdir(t, repo)

	manifest, err := buildManifest()
	if err != nil {
		t.Fatal(err)
	}
	// Force debt.score below min_score and check_count == 0.
	manifest.Debt.Score = 9.0
	manifest.Debt.MinScore = 9.8
	manifest.Debt.CheckCount = 0

	path := filepath.Join(t.TempDir(), "debt.json")
	if err := writeManifest(path, manifest); err != nil {
		t.Fatal(err)
	}

	err = verifyManifest(path, true, false, "", 0)
	if err == nil {
		t.Fatal("verifyManifest succeeded for debt score below min and zero check count")
	}
	message := err.Error()
	for _, want := range []string{
		"debt.score 9.0 is below minimum 9.8",
		"debt.check_count must be greater than zero",
	} {
		if !strings.Contains(message, want) {
			t.Fatalf("error = %q; want substring %q", message, want)
		}
	}
}
