package validation

import (
	"strings"
	"testing"
)

func TestRequireNonEmptyRejectsEmptyValue(t *testing.T) {
	if err := RequireNonEmpty("name", ""); err == nil {
		t.Fatal("expected empty value to fail")
	}
}

func TestRequireNonEmptyAcceptsValue(t *testing.T) {
	if err := RequireNonEmpty("name", "templatex"); err != nil {
		t.Fatalf("expected value to pass: %v", err)
	}
}

func TestValidateRuntimeFileOwnershipAcceptsControlPlaneIndex(t *testing.T) {
	gaps := ValidateRuntimeFileOwnership(".agent/runtime-file-ownership.yaml", runtimeFileOwnershipFixture())
	if len(gaps) != 0 {
		t.Fatalf("gaps = %#v; want none", gaps)
	}
}

func TestValidateRuntimeFileOwnershipRejectsMissingControlPlaneClassification(t *testing.T) {
	fixture := `schema_version: "2.9.3"
owners:
  "cmd/goalcli/":
    owner: gate-runtime
    review_required: true
    rationale: CLI validators.
  "contracts/":
    owner: standard
    review_required: true
    rationale: Schema contracts.
`
	gaps := ValidateRuntimeFileOwnership(".agent/runtime-file-ownership.yaml", fixture)
	if !validationGapsContain(gaps, ".agent/runtime-file-ownership.yaml owners must include .agent/") {
		t.Fatalf("gaps = %#v; want .agent/ classification gap", gaps)
	}
}

func TestValidateRuntimeFileOwnershipRejectsInvalidReviewRequired(t *testing.T) {
	fixture := runtimeFileOwnershipFixture()
	fixture = strings.Replace(fixture, "review_required: true", "review_required: maybe", 1)

	gaps := ValidateRuntimeFileOwnership(".agent/runtime-file-ownership.yaml", fixture)
	if !validationGapsContain(gaps, ".agent/runtime-file-ownership.yaml .agent/ review_required must be true or false") {
		t.Fatalf("gaps = %#v; want boolean review_required gap", gaps)
	}
}

func TestValidateRuntimeFileOwnershipRejectsDuplicateEntries(t *testing.T) {
	fixture := runtimeFileOwnershipFixture() + `  ".agent/":
    owner: governance
    review_required: true
    rationale: Duplicate control plane.
`
	gaps := ValidateRuntimeFileOwnership(".agent/runtime-file-ownership.yaml", fixture)
	if !validationGapsContain(gaps, ".agent/runtime-file-ownership.yaml duplicate owner entry .agent/") {
		t.Fatalf("gaps = %#v; want duplicate owner gap", gaps)
	}
}

func runtimeFileOwnershipFixture() string {
	return `schema_version: "2.9.3"
owners:
  ".agent/":
    owner: governance
    review_required: true
    rationale: Control plane manifests.
  "cmd/goalcli/":
    owner: gate-runtime
    review_required: true
    rationale: Goalcli validator surface.
  "contracts/":
    owner: standard
    review_required: true
    rationale: Public contracts.
`
}

func validationGapsContain(gaps []string, want string) bool {
	for _, gap := range gaps {
		if gap == want {
			return true
		}
	}
	return false
}
