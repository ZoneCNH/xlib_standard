package validation

import (
	"strings"
	"testing"
)

// This file is dedicated coverage backfill. It only adds table-driven tests for
// previously uncovered branches; it does not edit existing _test.go files.

func TestValidateRuntimeFileOwnershipEmptyContent(t *testing.T) {
	// Arrange
	path := ".agent/policies/runtime-file-ownership.yaml"
	// Act
	gaps := ValidateRuntimeFileOwnership(path, "   \n  \t\n")
	// Assert
	if !validationGapsContain(gaps, path+" must not be empty") {
		t.Fatalf("gaps = %#v; want must-not-be-empty gap", gaps)
	}
}

func TestValidateRuntimeFileOwnershipMissingSchemaVersionAndOwners(t *testing.T) {
	// Arrange
	path := ".agent/policies/runtime-file-ownership.yaml"
	fixture := `unrelated: "yes"
owners:
  ".agent/":
    owner: governance
    review_required: true
    review_rule: RULE-CHANGE
    rationale: Control plane manifests.
`
	// Act
	gaps := ValidateRuntimeFileOwnership(path, fixture)
	// Assert
	if !validationGapsContain(gaps, path+" missing schema_version") {
		t.Fatalf("gaps = %#v; want missing schema_version gap", gaps)
	}
}

func TestValidateRuntimeFileOwnershipMissingOwnersKey(t *testing.T) {
	// Arrange
	path := ".agent/policies/runtime-file-ownership.yaml"
	fixture := `schema_version: "2.9.3"
contexts: []
`
	// Act
	gaps := ValidateRuntimeFileOwnership(path, fixture)
	// Assert
	if !validationGapsContain(gaps, path+" missing owners") {
		t.Fatalf("gaps = %#v; want missing owners key gap", gaps)
	}
	if !validationGapsContain(gaps, path+" owners must include at least one path classification") {
		t.Fatalf("gaps = %#v; want empty owners gap", gaps)
	}
}

func TestValidateRuntimeFileOwnershipMissingOwnerAndRationale(t *testing.T) {
	// Arrange
	path := ".agent/policies/runtime-file-ownership.yaml"
	fixture := `schema_version: "2.9.3"
owners:
  ".agent/":
    review_required: false
`
	// Act
	gaps := ValidateRuntimeFileOwnership(path, fixture)
	// Assert
	if !validationGapsContain(gaps, path+" .agent/ missing owner") {
		t.Fatalf("gaps = %#v; want missing owner gap", gaps)
	}
	if !validationGapsContain(gaps, path+" .agent/ missing rationale") {
		t.Fatalf("gaps = %#v; want missing rationale gap", gaps)
	}
	// review_required false is valid; should NOT raise review-rule/required gap.
	for _, gap := range gaps {
		if strings.Contains(gap, "missing review_rule") || strings.Contains(gap, "review_required must be true or false") {
			t.Fatalf("unexpected review gap for false review_required: %s", gap)
		}
	}
}

func TestValidateRuntimeFileOwnershipReviewRequiredFalseNeedsNoRule(t *testing.T) {
	// Arrange: a NON-required path with review_required=false must not require review_rule.
	// The three required paths (.agent/, cmd/goalcli/, contracts/) are always forced true by
	// requireRuntimeOwner, so this exercises the per-entry "review_required != true && != false"
	// and "review_required==true && rule empty" branches via a benign extra path.
	gaps := ValidateRuntimeFileOwnership(".agent/policies/runtime-file-ownership.yaml", `schema_version: "2.9.3"
owners:
  ".agent/":
    owner: governance
    review_required: true
    review_rule: RULE-CHANGE
    rationale: Control plane manifests.
  "cmd/goalcli/":
    owner: gate-runtime
    review_required: true
    review_rule: GATE-RUNTIME-CHANGE
    rationale: Goalcli.
  "contracts/":
    owner: standard
    review_required: true
    review_rule: CONTRACT-CHANGE
    rationale: Contracts.
  "docs/":
    owner: governance
    review_required: false
    rationale: Documentation.
`)
	// Assert
	if len(gaps) != 0 {
		t.Fatalf("gaps = %#v; want none for review_required=false without rule on extra path", gaps)
	}
}

func TestValidateRuntimeFileOwnershipOwnerMismatchAndReviewMismatch(t *testing.T) {
	// Arrange: swap owners so requireRuntimeOwner reports mismatches.
	path := ".agent/policies/runtime-file-ownership.yaml"
	fixture := `schema_version: "2.9.3"
owners:
  ".agent/":
    owner: gate-runtime
    review_required: false
    rationale: Control plane manifests.
  "cmd/goalcli/":
    owner: governance
    review_required: false
    rationale: Goalcli.
  "contracts/":
    owner: gate-runtime
    review_required: false
    rationale: Contracts.
`
	// Act
	gaps := ValidateRuntimeFileOwnership(path, fixture)
	// Assert
	if !validationGapsContain(gaps, path+" .agent/ owner must be governance") {
		t.Fatalf("gaps = %#v; want .agent/ owner mismatch", gaps)
	}
	if !validationGapsContain(gaps, path+" .agent/ review_required must be true") {
		t.Fatalf("gaps = %#v; want .agent/ review_required mismatch", gaps)
	}
	if !validationGapsContain(gaps, path+" cmd/goalcli/ owner must be gate-runtime") {
		t.Fatalf("gaps = %#v; want cmd/goalcli/ owner mismatch", gaps)
	}
	if !validationGapsContain(gaps, path+" contracts/ owner must be standard") {
		t.Fatalf("gaps = %#v; want contracts/ owner mismatch", gaps)
	}
}

func TestValidateRuntimeFileOwnershipMissingReviewRequired(t *testing.T) {
	// Arrange: extra path with no review_required field at all → empty-string branch.
	path := ".agent/policies/runtime-file-ownership.yaml"
	fixture := `schema_version: "2.9.3"
owners:
  ".agent/":
    owner: governance
    review_required: true
    review_rule: RULE-CHANGE
    rationale: Control plane manifests.
  "cmd/goalcli/":
    owner: gate-runtime
    review_required: true
    review_rule: GATE-RUNTIME-CHANGE
    rationale: Goalcli.
  "contracts/":
    owner: standard
    review_required: true
    review_rule: CONTRACT-CHANGE
    rationale: Contracts.
  "docs/":
    owner: governance
    rationale: Documentation.
`
	// Act
	gaps := ValidateRuntimeFileOwnership(path, fixture)
	// Assert
	if !validationGapsContain(gaps, path+" docs/ missing review_required") {
		t.Fatalf("gaps = %#v; want docs/ missing review_required", gaps)
	}
}

// NOTE: ValidateExecutionContext line 121-122 (`if context.name == "" { continue }`) is
// defensively guarded code that the parseExecutionContexts parser cannot reach with any
// well-formed YAML input — the parser's TrimSpace + "- " prefix matching always yields a
// non-empty name (a bare "-" line falls through to the indent<4 skip). Recorded as
// untestable; not covered.

func TestValidateRuntimeFileOwnershipMissingRequiredPaths(t *testing.T) {
	// Arrange: none of the three required paths present.
	path := ".agent/policies/runtime-file-ownership.yaml"
	fixture := `schema_version: "2.9.3"
owners:
  "docs/":
    owner: governance
    review_required: true
    review_rule: DOC-CHANGE
    rationale: Docs.
`
	// Act
	gaps := ValidateRuntimeFileOwnership(path, fixture)
	// Assert
	for _, want := range []string{
		path + " owners must include .agent/",
		path + " owners must include cmd/goalcli/",
		path + " owners must include contracts/",
	} {
		if !validationGapsContain(gaps, want) {
			t.Fatalf("gaps = %#v; want %q", gaps, want)
		}
	}
}

func TestValidateExecutionContextEmptyAndMissingKeys(t *testing.T) {
	path := ".agent/policies/execution-context.yaml"
	// empty content
	if gaps := ValidateExecutionContext(path, "  \n\t\n", nil); !validationGapsContain(gaps, path+" must not be empty") {
		t.Fatalf("gaps = %#v; want must-not-be-empty", gaps)
	}
	// missing schema_version + contexts keys
	gaps := ValidateExecutionContext(path, "unrelated: true\n", nil)
	if !validationGapsContain(gaps, path+" missing schema_version") {
		t.Fatalf("gaps = %#v; want missing schema_version", gaps)
	}
	if !validationGapsContain(gaps, path+" missing contexts") {
		t.Fatalf("gaps = %#v; want missing contexts", gaps)
	}
	if !validationGapsContain(gaps, path+" contexts must include at least one execution context") {
		t.Fatalf("gaps = %#v; want empty contexts", gaps)
	}
}

func TestValidateExecutionContextDuplicateName(t *testing.T) {
	path := ".agent/policies/execution-context.yaml"
	fixture := `schema_version: "2.9.3"
contexts:
  local_write:
    write_scope: worktree
    mutates_files: true
    release_evidence: false
    requires_gowork: off
  local_write:
    write_scope: worktree
    mutates_files: true
    release_evidence: false
    requires_gowork: off
`
	gaps := ValidateExecutionContext(path, fixture, []string{"local_write"})
	if !validationGapsContain(gaps, path+" duplicate context local_write") {
		t.Fatalf("gaps = %#v; want duplicate context", gaps)
	}
}

func TestValidateExecutionContextMissingExpectedAndMissingFields(t *testing.T) {
	path := ".agent/policies/execution-context.yaml"
	// context with empty body → all required fields missing, plus bool defaults missing.
	fixture := `schema_version: "2.9.3"
contexts:
  local_write:
    note: empty
`
	gaps := ValidateExecutionContext(path, fixture, []string{"local_write", "release_verify"})
	for _, want := range []string{
		path + " local_write missing write_scope",
		path + " local_write missing mutates_files",
		path + " local_write missing release_evidence",
		path + " local_write missing requires_gowork",
		path + " missing context release_verify",
	} {
		if !validationGapsContain(gaps, want) {
			t.Fatalf("gaps = %#v; want %q", gaps, want)
		}
	}
}

func TestValidateExecutionContextInvalidBoolField(t *testing.T) {
	path := ".agent/policies/execution-context.yaml"
	fixture := `schema_version: "2.9.3"
contexts:
  local_write:
    write_scope: worktree
    mutates_files: maybe
    release_evidence: nope
    requires_gowork: off
`
	gaps := ValidateExecutionContext(path, fixture, []string{"local_write"})
	if !validationGapsContain(gaps, path+" local_write mutates_files must be true or false") {
		t.Fatalf("gaps = %#v; want mutates_files bool gap", gaps)
	}
	if !validationGapsContain(gaps, path+" local_write release_evidence must be true or false") {
		t.Fatalf("gaps = %#v; want release_evidence bool gap", gaps)
	}
}

func TestValidateExecutionContextAbsoluteRelativeField(t *testing.T) {
	path := ".agent/policies/execution-context.yaml"
	// contextFieldMustBeRelative matches path/root/manifest substrings only.
	// write_scope contains none of those, so it is NOT checked; use manifest_root + path_root.
	fixture := `schema_version: "2.9.3"
contexts:
  local_write:
    write_scope: worktree
    mutates_files: true
    release_evidence: false
    requires_gowork: off
    manifest_root: /abs/manifest
    archive_path: /abs/archive
`
	gaps := ValidateExecutionContext(path, fixture, []string{"local_write"})
	if !validationGapsContain(gaps, path+" local_write manifest_root must be repository-relative") {
		t.Fatalf("gaps = %#v; want manifest_root relative gap", gaps)
	}
	if !validationGapsContain(gaps, path+" local_write archive_path must be repository-relative") {
		t.Fatalf("gaps = %#v; want archive_path relative gap", gaps)
	}
}

func TestValidateExecutionContextLocalWriteAndReleaseVerifySemanticGaps(t *testing.T) {
	path := ".agent/policies/execution-context.yaml"
	// local_write with release_evidence=true and mutates_files=false (wrong).
	// release_verify with requires_gowork set to a non-"off" value → triggers
	// requireContextValue's "got != want" branch.
	fixture := `schema_version: "2.9.3"
contexts:
  local_write:
    write_scope: worktree
    mutates_files: false
    release_evidence: true
    requires_gowork: on
  release_verify:
    write_scope: release_read_only
    mutates_files: true
    release_evidence: false
    requires_gowork: on
`
	gaps := ValidateExecutionContext(path, fixture, []string{"local_write", "release_verify"})
	for _, want := range []string{
		path + " local_write mutates_files must be true",
		path + " local_write release_evidence must be false",
		path + " release_verify mutates_files must be false",
		path + " release_verify release_evidence must be true",
		path + " release_verify requires_gowork must be off",
	} {
		if !validationGapsContain(gaps, want) {
			t.Fatalf("gaps = %#v; want %q", gaps, want)
		}
	}
}

func TestValidateExecutionContextIdenticalLocalWriteReleaseVerify(t *testing.T) {
	path := ".agent/policies/execution-context.yaml"
	fields := `    write_scope: worktree
    mutates_files: true
    release_evidence: false
    requires_gowork: off
`
	fixture := "schema_version: \"2.9.3\"\ncontexts:\n  local_write:\n" + fields + "  release_verify:\n" + fields
	gaps := ValidateExecutionContext(path, fixture, []string{"local_write", "release_verify"})
	if !validationGapsContain(gaps, path+" local_write and release_verify must have distinct semantics") {
		t.Fatalf("gaps = %#v; want distinct-semantics gap", gaps)
	}
}

func TestContainsYAMLKeyTableDriven(t *testing.T) {
	cases := []struct {
		name    string
		content string
		key     string
		want    bool
	}{
		{"exact key colon", "schema_version: \"2.9.3\"\n", "schema_version", true},
		{"key with inline value", "owners: []\n", "owners", true},
		{"key with leading spaces and comment", "  # schema_version: hidden\n  schema_version: \"1\"\n", "schema_version", true},
		{"commented out only", "# schema_version: hidden\n", "schema_version", false},
		{"missing key", "unrelated: value\n", "schema_version", false},
		{"prefix match rejected", "schema_versionx: 1\n", "schema_version", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := containsYAMLKey(tc.content, tc.key); got != tc.want {
				t.Fatalf("containsYAMLKey(%q, %q) = %v; want %v", tc.content, tc.key, got, tc.want)
			}
		})
	}
}

func TestStripInlineYAMLComment(t *testing.T) {
	cases := []struct {
		name string
		line string
		want string
	}{
		{"no comment", "owners: []", "owners: []"},
		{"trailing comment", "owners: [] # list of owners", "owners: [] "},
		// strings.Cut splits at the first "#"; the function does NOT protect quoted values.
		{"hash splits even inside quotes", `url: "http://x#y"`, `url: "http://x`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := stripInlineYAMLComment(tc.line); got != tc.want {
				t.Fatalf("stripInlineYAMLComment(%q) = %q; want %q", tc.line, got, tc.want)
			}
		})
	}
}

func TestParseRuntimeFileOwnersEdgeCases(t *testing.T) {
	// Indent==0 line after owners should flush and break.
	content := "schema_version: \"2.9.3\"\nowners:\n  \"docs/\":\n    owner: standard\n    review_required: true\n    review_rule: R\n    rationale: x\nnot_owners: true\n"
	owners := parseRuntimeFileOwners(content)
	if len(owners) != 1 || owners[0].path != "docs/" || owners[0].owner != "standard" || owners[0].reviewRule != "R" {
		t.Fatalf("owners = %#v; want one docs/ owner", owners)
	}

	// Field line at indent<4 without current owner is dropped; bad field format dropped.
	content2 := "owners:\n    owner: orphan\n  \"docs/\":\n    owner: standard\n    not_a_field\n    review_required: true\n    rationale: x\n"
	owners2 := parseRuntimeFileOwners(content2)
	if len(owners2) != 1 || owners2[0].owner != "standard" {
		t.Fatalf("owners2 = %#v; want one owner with bad field dropped", owners2)
	}

	// Owners never opened → no entries.
	if o := parseRuntimeFileOwners("schema_version: \"2.9.3\"\n"); len(o) != 0 {
		t.Fatalf("o = %#v; want none", o)
	}
}

func TestParseExecutionContextsEdgeCases(t *testing.T) {
	// ListItem form at indent 2.
	content := "contexts:\n  - local_write\n    write_scope: worktree\n    mutates_files: true\n"
	ctx := parseExecutionContexts(content)
	if len(ctx) != 1 || ctx[0].name != "local_write" || ctx[0].fields["write_scope"] != "worktree" {
		t.Fatalf("ctx = %#v; want list-item form parsed", ctx)
	}

	// indent==0 flushes and breaks.
	content2 := "contexts:\n  local_write:\n    write_scope: worktree\nnot_contexts: true\n"
	ctx2 := parseExecutionContexts(content2)
	if len(ctx2) != 1 || ctx2[0].name != "local_write" {
		t.Fatalf("ctx2 = %#v; want one context", ctx2)
	}

	// contexts never opened → empty.
	if c := parseExecutionContexts("schema_version: \"2.9.3\"\n"); len(c) != 0 {
		t.Fatalf("c = %#v; want none", c)
	}

	// field at indent<4 dropped.
	content3 := "contexts:\n  local_write:\n  write_scope: worktree\n"
	ctx3 := parseExecutionContexts(content3)
	if len(ctx3) != 1 || len(ctx3[0].fields) != 0 {
		t.Fatalf("ctx3 = %#v; want context with no fields", ctx3)
	}

	// bad field (no colon) dropped.
	content4 := "contexts:\n  local_write:\n    just_text_no_colon\n"
	ctx4 := parseExecutionContexts(content4)
	if len(ctx4) != 1 || len(ctx4[0].fields) != 0 {
		t.Fatalf("ctx4 = %#v; want fields empty", ctx4)
	}
}

func TestRequireContextFieldDirect(t *testing.T) {
	var gaps []string
	// field present → no gap.
	requireContextField("p", executionContext{name: "n", fields: map[string]string{"f": "v"}}, "f", &gaps)
	if len(gaps) != 0 {
		t.Fatalf("gaps = %#v; want none when field present", gaps)
	}
	// field missing → gap.
	requireContextField("p", executionContext{name: "n", fields: map[string]string{}}, "f", &gaps)
	if len(gaps) != 1 || gaps[0] != "p n missing f" {
		t.Fatalf("gaps = %#v; want missing-f gap", gaps)
	}
}

func TestRequireBoolContextFieldDirect(t *testing.T) {
	var gaps []string
	// present and valid → no gap.
	requireBoolContextField("p", executionContext{name: "n", fields: map[string]string{"f": "true"}}, "f", &gaps)
	requireBoolContextField("p", executionContext{name: "n", fields: map[string]string{"f": "false"}}, "f", &gaps)
	if len(gaps) != 0 {
		t.Fatalf("gaps = %#v; want none for valid bools", gaps)
	}
	// missing → gap.
	requireBoolContextField("p", executionContext{name: "n", fields: map[string]string{}}, "f", &gaps)
	if !validationGapsContain(gaps, "p n missing f") {
		t.Fatalf("gaps = %#v; want missing-f gap", gaps)
	}
	// invalid → gap.
	gaps = nil
	requireBoolContextField("p", executionContext{name: "n", fields: map[string]string{"f": "maybe"}}, "f", &gaps)
	if !validationGapsContain(gaps, "p n f must be true or false") {
		t.Fatalf("gaps = %#v; want bool gap", gaps)
	}
}
