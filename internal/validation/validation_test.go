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
	gaps := ValidateRuntimeFileOwnership(".agent/policies/runtime-file-ownership.yaml", runtimeFileOwnershipFixture())
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
	gaps := ValidateRuntimeFileOwnership(".agent/policies/runtime-file-ownership.yaml", fixture)
	if !validationGapsContain(gaps, ".agent/policies/runtime-file-ownership.yaml owners must include .agent/") {
		t.Fatalf("gaps = %#v; want .agent/ classification gap", gaps)
	}
}

func TestValidateRuntimeFileOwnershipRejectsInvalidReviewRequired(t *testing.T) {
	fixture := runtimeFileOwnershipFixture()
	fixture = strings.Replace(fixture, "review_required: true", "review_required: maybe", 1)

	gaps := ValidateRuntimeFileOwnership(".agent/policies/runtime-file-ownership.yaml", fixture)
	if !validationGapsContain(gaps, ".agent/policies/runtime-file-ownership.yaml .agent/ review_required must be true or false") {
		t.Fatalf("gaps = %#v; want boolean review_required gap", gaps)
	}
}

func TestValidateRuntimeFileOwnershipRejectsDuplicateEntries(t *testing.T) {
	fixture := runtimeFileOwnershipFixture() + `  ".agent/":
    owner: governance
    review_required: true
    review_rule: RULE-CHANGE
    rationale: Duplicate control plane.
`
	gaps := ValidateRuntimeFileOwnership(".agent/policies/runtime-file-ownership.yaml", fixture)
	if !validationGapsContain(gaps, ".agent/policies/runtime-file-ownership.yaml duplicate owner entry .agent/") {
		t.Fatalf("gaps = %#v; want duplicate owner gap", gaps)
	}
}

func TestValidateRuntimeFileOwnershipRejectsUnknownOwner(t *testing.T) {
	fixture := strings.Replace(runtimeFileOwnershipFixture(), "owner: governance", "owner: mystery-owner", 1)

	gaps := ValidateRuntimeFileOwnership(".agent/policies/runtime-file-ownership.yaml", fixture)
	if !validationGapsContain(gaps, ".agent/policies/runtime-file-ownership.yaml .agent/ unknown owner mystery-owner") {
		t.Fatalf("gaps = %#v; want unknown owner gap", gaps)
	}
}

func TestValidateRuntimeFileOwnershipRequiresReviewRule(t *testing.T) {
	fixture := strings.Replace(runtimeFileOwnershipFixture(), "    review_rule: RULE-CHANGE\n", "", 1)

	gaps := ValidateRuntimeFileOwnership(".agent/policies/runtime-file-ownership.yaml", fixture)
	if !validationGapsContain(gaps, ".agent/policies/runtime-file-ownership.yaml .agent/ missing review_rule") {
		t.Fatalf("gaps = %#v; want missing review_rule gap", gaps)
	}
}

func TestValidateRuntimeFileOwnershipRejectsAbsoluteOwnerPath(t *testing.T) {
	fixture := strings.Replace(runtimeFileOwnershipFixture(), "\".agent/\":", "\"/tmp/.agent/\":", 1)

	gaps := ValidateRuntimeFileOwnership(".agent/policies/runtime-file-ownership.yaml", fixture)
	if !validationGapsContain(gaps, ".agent/policies/runtime-file-ownership.yaml /tmp/.agent/ must be repository-relative") {
		t.Fatalf("gaps = %#v; want repository-relative path gap", gaps)
	}
}

func TestValidateRuntimeFileOwnershipRejectsEmptyManifest(t *testing.T) {
	gaps := ValidateRuntimeFileOwnership(".agent/policies/runtime-file-ownership.yaml", " \n\t")
	if len(gaps) != 1 || gaps[0] != ".agent/policies/runtime-file-ownership.yaml must not be empty" {
		t.Fatalf("gaps = %#v; want empty manifest gap", gaps)
	}
}

func TestValidateRuntimeFileOwnershipRequiresTopLevelKeysAndEntries(t *testing.T) {
	gaps := ValidateRuntimeFileOwnership(".agent/policies/runtime-file-ownership.yaml", `schema_version_extra: "2.9.3"
# owners: hidden in comment
`)
	validationGapsContainAll(t, gaps,
		".agent/policies/runtime-file-ownership.yaml missing schema_version",
		".agent/policies/runtime-file-ownership.yaml missing owners",
		".agent/policies/runtime-file-ownership.yaml owners must include at least one path classification",
		".agent/policies/runtime-file-ownership.yaml owners must include .agent/",
		".agent/policies/runtime-file-ownership.yaml owners must include cmd/goalcli/",
		".agent/policies/runtime-file-ownership.yaml owners must include contracts/",
	)
}

func TestValidateRuntimeFileOwnershipRejectsMissingOwnerMetadata(t *testing.T) {
	fixture := `schema_version: "2.9.3"
owners:
  ".agent/":
    review_required: true
    review_rule: RULE-CHANGE
  "cmd/goalcli/":
    owner: gate-runtime
    review_rule: GATE-RUNTIME-CHANGE
    rationale: Goalcli validator surface.
  "contracts/":
    owner: standard
    review_required: true
    review_rule: CONTRACT-CHANGE
    rationale: Public contracts.
`
	gaps := ValidateRuntimeFileOwnership(".agent/policies/runtime-file-ownership.yaml", fixture)
	validationGapsContainAll(t, gaps,
		".agent/policies/runtime-file-ownership.yaml .agent/ missing owner",
		".agent/policies/runtime-file-ownership.yaml .agent/ missing rationale",
		".agent/policies/runtime-file-ownership.yaml cmd/goalcli/ missing review_required",
	)
}

func TestValidateRuntimeFileOwnershipEnforcesRequiredOwnerSemantics(t *testing.T) {
	fixture := runtimeFileOwnershipFixture()
	fixture = strings.Replace(fixture, "owner: governance", "owner: security", 1)
	fixture = strings.Replace(fixture, "review_required: true", "review_required: false", 1)

	gaps := ValidateRuntimeFileOwnership(".agent/policies/runtime-file-ownership.yaml", fixture)
	validationGapsContainAll(t, gaps,
		".agent/policies/runtime-file-ownership.yaml .agent/ owner must be governance",
		".agent/policies/runtime-file-ownership.yaml .agent/ review_required must be true",
	)
}

func TestValidateRuntimeFileOwnershipAllowsInlineComments(t *testing.T) {
	fixture := `schema_version: "2.9.3" # manifest version
owners: # runtime owner index
  ".agent/": # control plane
    owner: governance # owner
    review_required: true # bool
    review_rule: RULE-CHANGE # rule
    rationale: Control plane manifests. # rationale
  "cmd/goalcli/":
    owner: gate-runtime # owner
    review_required: true
    review_rule: GATE-RUNTIME-CHANGE
    rationale: Goalcli validator surface.
  "contracts/":
    owner: standard
    review_required: true
    review_rule: CONTRACT-CHANGE
    rationale: Public contracts.
`
	gaps := ValidateRuntimeFileOwnership(".agent/policies/runtime-file-ownership.yaml", fixture)
	if len(gaps) != 0 {
		t.Fatalf("gaps = %#v; want inline comments accepted", gaps)
	}
}

func TestValidateExecutionContextAcceptsSemanticManifest(t *testing.T) {
	gaps := ValidateExecutionContext(".agent/policies/execution-context.yaml", executionContextFixture(), executionContextsFixture())
	if len(gaps) != 0 {
		t.Fatalf("gaps = %#v; want none", gaps)
	}
}

func TestValidateExecutionContextRejectsUnknownContext(t *testing.T) {
	fixture := strings.Replace(executionContextFixture(), "release_verify:", "release_magic:", 1)

	gaps := ValidateExecutionContext(".agent/policies/execution-context.yaml", fixture, executionContextsFixture())
	if !validationGapsContain(gaps, ".agent/policies/execution-context.yaml unknown context release_magic") {
		t.Fatalf("gaps = %#v; want unknown context gap", gaps)
	}
}

func TestValidateExecutionContextRequiresDistinctLocalWriteAndReleaseVerify(t *testing.T) {
	fixture := strings.Replace(executionContextFixture(), "mutates_files: false\n    release_evidence: true", "mutates_files: true\n    release_evidence: false", 1)

	gaps := ValidateExecutionContext(".agent/policies/execution-context.yaml", fixture, executionContextsFixture())
	if !validationGapsContain(gaps, ".agent/policies/execution-context.yaml release_verify mutates_files must be false") ||
		!validationGapsContain(gaps, ".agent/policies/execution-context.yaml release_verify release_evidence must be true") {
		t.Fatalf("gaps = %#v; want release_verify semantic gaps", gaps)
	}
}

func TestValidateExecutionContextRejectsEmptyManifest(t *testing.T) {
	gaps := ValidateExecutionContext(".agent/policies/execution-context.yaml", " \n\t", executionContextsFixture())
	if len(gaps) != 1 || gaps[0] != ".agent/policies/execution-context.yaml must not be empty" {
		t.Fatalf("gaps = %#v; want empty manifest gap", gaps)
	}
}

func TestValidateExecutionContextRequiresTopLevelKeysAndEntries(t *testing.T) {
	gaps := ValidateExecutionContext(".agent/policies/execution-context.yaml", `schema_version_extra: "2.9.3"
# contexts: hidden in comment
`, []string{"local_write"})
	validationGapsContainAll(t, gaps,
		".agent/policies/execution-context.yaml missing schema_version",
		".agent/policies/execution-context.yaml missing contexts",
		".agent/policies/execution-context.yaml contexts must include at least one execution context",
		".agent/policies/execution-context.yaml missing context local_write",
	)
}

func TestValidateExecutionContextRejectsDuplicateContext(t *testing.T) {
	fixture := executionContextFixture() + `  local_write:
    write_scope: worktree
    mutates_files: true
    release_evidence: false
    requires_gowork: off
`
	gaps := ValidateExecutionContext(".agent/policies/execution-context.yaml", fixture, executionContextsFixture())
	if !validationGapsContain(gaps, ".agent/policies/execution-context.yaml duplicate context local_write") {
		t.Fatalf("gaps = %#v; want duplicate context gap", gaps)
	}
}

func TestValidateExecutionContextRejectsMissingAndInvalidRequiredFields(t *testing.T) {
	fixture := `schema_version: "2.9.3"
contexts:
  local_write:
    mutates_files: maybe
    release_evidence: also-maybe
  release_verify:
    write_scope: release_read_only
    mutates_files: false
    release_evidence: true
    requires_gowork: on
`
	gaps := ValidateExecutionContext(".agent/policies/execution-context.yaml", fixture, []string{"local_write", "release_verify"})
	validationGapsContainAll(t, gaps,
		".agent/policies/execution-context.yaml local_write missing write_scope",
		".agent/policies/execution-context.yaml local_write mutates_files must be true or false",
		".agent/policies/execution-context.yaml local_write release_evidence must be true or false",
		".agent/policies/execution-context.yaml local_write missing requires_gowork",
		".agent/policies/execution-context.yaml local_write mutates_files must be true",
		".agent/policies/execution-context.yaml local_write release_evidence must be false",
		".agent/policies/execution-context.yaml release_verify requires_gowork must be off",
	)
}

func TestValidateExecutionContextRequiresBoolFieldsWhenAbsent(t *testing.T) {
	fixture := `schema_version: "2.9.3"
contexts:
  local_write:
    write_scope: worktree
    requires_gowork: off
`
	gaps := ValidateExecutionContext(".agent/policies/execution-context.yaml", fixture, []string{"local_write"})
	validationGapsContainAll(t, gaps,
		".agent/policies/execution-context.yaml local_write missing mutates_files",
		".agent/policies/execution-context.yaml local_write missing release_evidence",
	)
}

func TestValidateExecutionContextRequiresRelativePathLikeFields(t *testing.T) {
	fixture := strings.Replace(executionContextFixture(), "requires_gowork: off", `requires_gowork: off
    config_path: /tmp/config.yaml
    workspace_root: /tmp/work
    evidence_manifest: /tmp/manifest.json`, 1)

	gaps := ValidateExecutionContext(".agent/policies/execution-context.yaml", fixture, executionContextsFixture())
	validationGapsContainAll(t, gaps,
		".agent/policies/execution-context.yaml local_write config_path must be repository-relative",
		".agent/policies/execution-context.yaml local_write workspace_root must be repository-relative",
		".agent/policies/execution-context.yaml local_write evidence_manifest must be repository-relative",
	)
}

func TestValidateExecutionContextRequiresConfiguredContexts(t *testing.T) {
	gaps := ValidateExecutionContext(".agent/policies/execution-context.yaml", executionContextFixture(), append(executionContextsFixture(), "release_candidate"))
	if !validationGapsContain(gaps, ".agent/policies/execution-context.yaml missing context release_candidate") {
		t.Fatalf("gaps = %#v; want missing expected context gap", gaps)
	}
}

func TestValidateExecutionContextAllowsListSyntaxAndInlineComments(t *testing.T) {
	fixture := `schema_version: "2.9.3" # manifest version
contexts: # accepted context list
  - local_write # local writes
    write_scope: worktree # scope
    mutates_files: true # bool
    release_evidence: false # bool
    requires_gowork: off # gowork mode
  - release_verify
    write_scope: release_read_only
    mutates_files: false
    release_evidence: true
    requires_gowork: off
`
	gaps := ValidateExecutionContext(".agent/policies/execution-context.yaml", fixture, []string{"local_write", "release_verify"})
	if len(gaps) != 0 {
		t.Fatalf("gaps = %#v; want list syntax with inline comments accepted", gaps)
	}
}

func TestValidateExecutionContextDetectsIdenticalLocalWriteAndReleaseVerifyFields(t *testing.T) {
	fixture := strings.Replace(executionContextFixture(), `write_scope: release_read_only
    mutates_files: false
    release_evidence: true
    requires_gowork: off`, `write_scope: worktree
    mutates_files: true
    release_evidence: false
    requires_gowork: off`, 1)

	gaps := ValidateExecutionContext(".agent/policies/execution-context.yaml", fixture, executionContextsFixture())
	if !validationGapsContain(gaps, ".agent/policies/execution-context.yaml local_write and release_verify must have distinct semantics") {
		t.Fatalf("gaps = %#v; want distinct semantics gap", gaps)
	}
}

func TestContainsYAMLKeyBoundaries(t *testing.T) {
	tests := []struct {
		name    string
		content string
		key     string
		want    bool
	}{
		{
			name: "key with scalar and comment",
			content: `schema_version: "2.9.3" # comment
owners:
`,
			key:  "schema_version",
			want: true,
		},
		{
			name: "key with comment only value",
			content: `owners: # comment
`,
			key:  "owners",
			want: true,
		},
		{
			name: "commented key ignored",
			content: `# contexts:
`,
			key:  "contexts",
			want: false,
		},
		{
			name: "prefix key ignored",
			content: `schema_version_extra: "2.9.3"
`,
			key:  "schema_version",
			want: false,
		},
		{
			name: "missing key false",
			content: `contexts:
`,
			key:  "owners",
			want: false,
		},
		{
			name: "compact scalar without space ignored",
			content: `schema_version:"2.9.3"
`,
			key:  "schema_version",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := containsYAMLKey(tt.content, tt.key); got != tt.want {
				t.Fatalf("containsYAMLKey(%q, %q) = %t; want %t", tt.content, tt.key, got, tt.want)
			}
		})
	}
}

func TestParseRuntimeFileOwnersIgnoresMalformedLinesAndStopsAtTopLevel(t *testing.T) {
	owners := parseRuntimeFileOwners(`metadata: ignored
owners:
  - not a mapping entry
  ".agent/":
    owner: "governance"
    review_required: 'true'
    review_rule: RULE-CHANGE
    rationale: Control plane manifests.
    malformed field line
    unknown: ignored
next_section:
  "cmd/goalcli/":
    owner: gate-runtime
`)
	if len(owners) != 1 {
		t.Fatalf("owners = %#v; want one parsed owner", owners)
	}
	if owners[0].path != ".agent/" ||
		owners[0].owner != "governance" ||
		owners[0].reviewRequired != "true" ||
		owners[0].reviewRule != "RULE-CHANGE" ||
		owners[0].rationale != "Control plane manifests." {
		t.Fatalf("owner = %#v; want parsed and unquoted fields", owners[0])
	}
}

func TestParseExecutionContextsIgnoresMalformedLinesAndStopsAtTopLevel(t *testing.T) {
	contexts := parseExecutionContexts(`metadata: ignored
contexts:
  stray_without_colon
  - "local_write"
    write_scope: worktree
    mutates_files: true
    release_evidence: false
    requires_gowork: off
    malformed field line
    unknown_field: retained
next_section:
  release_verify:
    write_scope: release_read_only
`)
	if len(contexts) != 1 {
		t.Fatalf("contexts = %#v; want one parsed context", contexts)
	}
	if contexts[0].name != "local_write" {
		t.Fatalf("context name = %q; want local_write", contexts[0].name)
	}
	if got := contexts[0].fields["unknown_field"]; got != "retained" {
		t.Fatalf("unknown field = %q; want retained", got)
	}
}

func TestValidateExecutionContextSkipsEmptyContextNames(t *testing.T) {
	fixture := `schema_version: "2.9.3"
contexts:
  "":
    write_scope: worktree
    mutates_files: true
    release_evidence: false
    requires_gowork: off
`
	gaps := ValidateExecutionContext(".agent/policies/execution-context.yaml", fixture, []string{"local_write"})
	if !validationGapsContain(gaps, ".agent/policies/execution-context.yaml missing context local_write") {
		t.Fatalf("gaps = %#v; want missing expected context after empty name is skipped", gaps)
	}
}

func TestStripInlineYAMLCommentBoundaries(t *testing.T) {
	tests := []struct {
		name string
		line string
		want string
	}{
		{name: "inline comment", line: "owner: governance # comment", want: "owner: governance "},
		{name: "full line comment", line: "# comment", want: ""},
		{name: "no comment", line: "owner: governance", want: "owner: governance"},
		{name: "quoted hash follows current simple parser", line: `rationale: "value # with hash"`, want: `rationale: "value `},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stripInlineYAMLComment(tt.line); got != tt.want {
				t.Fatalf("stripInlineYAMLComment(%q) = %q; want %q", tt.line, got, tt.want)
			}
		})
	}
}

func TestContextFieldMustBeRelativeBoundaries(t *testing.T) {
	tests := []struct {
		field string
		want  bool
	}{
		{field: "config_path", want: true},
		{field: "workspace_root", want: true},
		{field: "evidence_manifest", want: true},
		{field: "write_scope", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.field, func(t *testing.T) {
			if got := contextFieldMustBeRelative(tt.field); got != tt.want {
				t.Fatalf("contextFieldMustBeRelative(%q) = %t; want %t", tt.field, got, tt.want)
			}
		})
	}
}

func runtimeFileOwnershipFixture() string {
	return `schema_version: "2.9.3"
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
    rationale: Goalcli validator surface.
  "contracts/":
    owner: standard
    review_required: true
    review_rule: CONTRACT-CHANGE
    rationale: Public contracts.
`
}

func executionContextFixture() string {
	return `schema_version: "2.9.3"
contexts:
  local_write:
    write_scope: worktree
    mutates_files: true
    release_evidence: false
    requires_gowork: off
  local_readonly:
    write_scope: read_only
    mutates_files: false
    release_evidence: false
    requires_gowork: off
  ci_pull_request:
    write_scope: read_only
    mutates_files: false
    release_evidence: false
    requires_gowork: off
  ci_main_verify:
    write_scope: read_only
    mutates_files: false
    release_evidence: false
    requires_gowork: off
  release_verify:
    write_scope: release_read_only
    mutates_files: false
    release_evidence: true
    requires_gowork: off
`
}

func executionContextsFixture() []string {
	return []string{"local_write", "local_readonly", "ci_pull_request", "ci_main_verify", "release_verify"}
}

func validationGapsContain(gaps []string, want string) bool {
	for _, gap := range gaps {
		if gap == want {
			return true
		}
	}
	return false
}

func validationGapsContainAll(t *testing.T, gaps []string, wants ...string) {
	t.Helper()
	for _, want := range wants {
		if !validationGapsContain(gaps, want) {
			t.Fatalf("gaps = %#v; want %q", gaps, want)
		}
	}
}
