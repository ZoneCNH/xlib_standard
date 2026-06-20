package debtcheck

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

type fakeDirReader struct {
	names    []string
	readErr  error
	closeErr error
}

func (f fakeDirReader) Readdirnames(int) ([]string, error) {
	return f.names, f.readErr
}

func (f fakeDirReader) Close() error {
	return f.closeErr
}

func TestRunPassesWithPolicyFilesAndCleanTree(t *testing.T) {
	root := t.TempDir()
	writePolicyFiles(t, root)
	writeDownstreamFiles(t, root)
	writeFile(t, root, "safe.go", "package fixture\n")

	report, err := Run(Options{Root: root, Mode: "enforce", MinScore: DefaultMinScore})
	if err != nil {
		t.Fatal(err)
	}

	if report.Status != "passed" {
		t.Fatalf("status = %q, want passed: %+v", report.Status, report.Summary)
	}
	if code := ExitCode(report); code != 0 {
		t.Fatalf("ExitCode = %d, want 0", code)
	}
	if problems := ValidateEvidence(EvidenceFromReport(report), DefaultMinScore); len(problems) != 0 {
		t.Fatalf("ValidateEvidence problems = %v, want none", problems)
	}
	if report.Digests.Report == "" || report.Digests.Rules == "missing" {
		t.Fatalf("digests = %+v, want populated policy and report digests", report.Digests)
	}
}

func TestRunPassesDownstreamSectionWithGovernanceFiles(t *testing.T) {
	root := t.TempDir()
	writePolicyFiles(t, root)
	writeDownstreamFiles(t, root)

	report, err := Run(Options{Root: root, Section: "downstream", Mode: "enforce", MinScore: DefaultMinScore})
	if err != nil {
		t.Fatal(err)
	}

	if report.Status != "passed" {
		t.Fatalf("status = %q, want passed: %+v", report.Status, report.Summary)
	}
	if len(report.Sections) != 1 || report.Sections[0].Name != "downstream" {
		t.Fatalf("sections = %+v, want only downstream", report.Sections)
	}
	if code := ExitCode(report); code != 0 {
		t.Fatalf("ExitCode = %d, want 0", code)
	}
}

func TestRunFailsDownstreamSectionWithoutRegistryCoverage(t *testing.T) {
	root := t.TempDir()
	writePolicyFiles(t, root)
	writeDownstreamFiles(t, root)
	writeFile(t, root, ".agent/registries/downstream-registry.yaml", `schema_version: "2.9.3"
downstreams:
  - repo: kernel/configx
    mode: patch-only
    status: unavailable_in_worker_workspace_gap_explicit
`)

	report, err := Run(Options{Root: root, Section: "downstream", Mode: "enforce", MinScore: DefaultMinScore})
	if err != nil {
		t.Fatal(err)
	}

	if report.Status != "failed" {
		t.Fatalf("status = %q, want failed", report.Status)
	}
	if report.Summary.P0 != 2 {
		t.Fatalf("P0 = %d, want 2", report.Summary.P0)
	}
	if !strings.Contains(ToMarkdown(report), "debt.downstream.registry-missing-repo") {
		t.Fatalf("markdown missing downstream registry finding: %s", ToMarkdown(report))
	}
}

func TestRunFailsDownstreamSectionOnFalseAdoptionClaim(t *testing.T) {
	root := t.TempDir()
	writePolicyFiles(t, root)
	writeDownstreamFiles(t, root)
	writeFile(t, root, ".agent/registries/downstream-adoption-status.yaml", `schema_version: "2.9.3"
current_registry:
  adoption_status: adopted
  proof_based_adoption: true
standard_target_libraries:
  - name: kernel
  - name: configx
  - name: observex
  - name: testkitx
  - name: postgresx
  - name: redisx
  - name: kafkax
  - name: natsx
  - name: taosx
  - name: ossx
  - name: clickhousex
`)

	report, err := Run(Options{Root: root, Section: "downstream", Mode: "enforce", MinScore: DefaultMinScore})
	if err != nil {
		t.Fatal(err)
	}

	if report.Status != "failed" {
		t.Fatalf("status = %q, want failed", report.Status)
	}
	markdown := ToMarkdown(report)
	if !strings.Contains(markdown, "debt.downstream.false-adoption-claim") {
		t.Fatalf("markdown missing false adoption finding: %s", markdown)
	}
	if !strings.Contains(markdown, "debt.downstream.false-proof-claim") {
		t.Fatalf("markdown missing false proof finding: %s", markdown)
	}
}

func TestRunFailsDownstreamSectionWithoutIntegrationDebtEvidenceGate(t *testing.T) {
	root := t.TempDir()
	writePolicyFiles(t, root)
	writeDownstreamFiles(t, root)
	writeFile(t, root, "scripts/run_integration.sh", `#!/usr/bin/env bash
TARGETS=(
  "kernel|github.com/ZoneCNH/kernel|kernel"
  "configx|github.com/ZoneCNH/configx|configx"
  "redisx|github.com/ZoneCNH/redisx|redisx"
)
`)

	report, err := Run(Options{Root: root, Section: "downstream", Mode: "enforce", MinScore: DefaultMinScore})
	if err != nil {
		t.Fatal(err)
	}

	if report.Status != "failed" {
		t.Fatalf("status = %q, want failed", report.Status)
	}
	if !strings.Contains(ToMarkdown(report), "debt.downstream.integration-missing-contract") {
		t.Fatalf("markdown missing downstream integration finding: %s", ToMarkdown(report))
	}
}

func TestRunFailsDownstreamSectionWithoutRenderedDebtEvidenceExclusions(t *testing.T) {
	root := t.TempDir()
	writePolicyFiles(t, root)
	writeDownstreamFiles(t, root)
	writeFile(t, root, "scripts/render_template.sh", `#!/usr/bin/env bash
rsync "$@"
`)

	report, err := Run(Options{Root: root, Section: "downstream", Mode: "enforce", MinScore: DefaultMinScore})
	if err != nil {
		t.Fatal(err)
	}

	if report.Status != "failed" {
		t.Fatalf("status = %q, want failed", report.Status)
	}
	if !strings.Contains(ToMarkdown(report), "debt.downstream.render-template-missing-exclusion") {
		t.Fatalf("markdown missing downstream render-template finding: %s", ToMarkdown(report))
	}
}

func TestRunFailsOnLegacyProductionImport(t *testing.T) {
	root := t.TempDir()
	writePolicyFiles(t, root)
	writeFile(t, root, "bad.go", "package fixture\n\nimport _ \"github.com/ZoneCNH/x.go\"\n")

	report, err := Run(Options{Root: root, Section: "architecture", Mode: "enforce", MinScore: DefaultMinScore})
	if err != nil {
		t.Fatal(err)
	}

	if report.Status != "failed" {
		t.Fatalf("status = %q, want failed", report.Status)
	}
	if report.Summary.P0 != 1 {
		t.Fatalf("P0 = %d, want 1", report.Summary.P0)
	}
	if code := ExitCode(report); code != 1 {
		t.Fatalf("ExitCode = %d, want 1", code)
	}
	if !strings.Contains(ToMarkdown(report), "legacy ZoneCNH x module") {
		t.Fatalf("markdown missing legacy import finding: %s", ToMarkdown(report))
	}
}

func TestRunAnnotatesFindingsWithOptionalRegistryMetadata(t *testing.T) {
	root := t.TempDir()
	writePolicyFiles(t, root)
	writeFile(t, root, DefaultRegistryPath, `schema_version: debt-rule-registry/v1
rules:
  - id: debt.docs.marker
    section: docs
    severity: P1
    invariant_id: INV-DEBT-DOCS-001
    release_blocking: false
    proof_depth: evidence_replay
    owner: standard
    expiry: 2026-07-01
    remediation: remove docs drift marker
    detector: debtcheck.scanTextMarker
`)
	writeFile(t, root, "docs/bad.md", "xlib-docs-drift\n")

	report, err := Run(Options{Root: root, Section: "docs", Mode: "warn", MinScore: DefaultMinScore})
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Sections) != 1 || len(report.Sections[0].Findings) != 1 {
		t.Fatalf("findings = %+v; want one docs finding", report.Sections)
	}
	finding := report.Sections[0].Findings[0]
	if finding.InvariantID != "INV-DEBT-DOCS-001" || finding.ProofDepth != "evidence_replay" || finding.Owner != "standard" ||
		finding.Expiry != "2026-07-01" || finding.Remediation != "remove docs drift marker" || finding.Detector != "debtcheck.scanTextMarker" {
		t.Fatalf("finding = %+v; want optional metadata from registry", finding)
	}
	if finding.ReleaseBlocking == nil || *finding.ReleaseBlocking != false {
		t.Fatalf("release_blocking = %v; want explicit false from registry", finding.ReleaseBlocking)
	}
	markdown := ToMarkdown(report)
	for _, want := range []string{
		"invariant_id=INV-DEBT-DOCS-001",
		"release_blocking=false",
		"proof_depth=evidence_replay",
		"owner=standard",
		"expiry=2026-07-01",
		"remediation=remove docs drift marker",
		"detector=debtcheck.scanTextMarker",
	} {
		if !strings.Contains(markdown, want) {
			t.Fatalf("markdown = %s; want optional metadata %q", markdown, want)
		}
	}
}

func TestFindingOptionalMetadataRoundTripAndMarkdown(t *testing.T) {
	releaseBlocking := false
	metadataFinding := Finding{
		ID:              "debt.docs.marker",
		Severity:        "P1",
		Path:            "docs/standard/debt-governance.md",
		Message:         "documentation debt marker is present",
		InvariantID:     "INV-DEBT-DOCS-001",
		ReleaseBlocking: &releaseBlocking,
		ProofDepth:      "evidence_replay",
		Owner:           "standard",
		Expiry:          "2026-07-01",
		Remediation:     "remove the marker after publishing replacement guidance",
		Detector:        "debtcheck.scanTextMarker",
	}
	metadataSection := SectionReport{
		Name:     "docs",
		Status:   "warning",
		P1:       1,
		Findings: []Finding{metadataFinding},
	}
	report := Report{
		SchemaVersion: SchemaVersion,
		Status:        "warning",
		Mode:          "warn",
		Sections:      []SectionReport{metadataSection},
	}

	data, err := json.Marshal(report)
	if err != nil {
		t.Fatal(err)
	}
	encoded := string(data)
	for _, want := range []string{
		`"invariant_id":"INV-DEBT-DOCS-001"`,
		`"release_blocking":false`,
		`"proof_depth":"evidence_replay"`,
		`"owner":"standard"`,
		`"expiry":"2026-07-01"`,
		`"remediation":"remove the marker after publishing replacement guidance"`,
		`"detector":"debtcheck.scanTextMarker"`,
	} {
		if !strings.Contains(encoded, want) {
			t.Fatalf("encoded report = %s; want %s", encoded, want)
		}
	}

	var roundTrip Report
	if err := json.Unmarshal(data, &roundTrip); err != nil {
		t.Fatal(err)
	}
	roundTripFinding := roundTrip.Sections[0].Findings[0]
	if roundTripFinding.ReleaseBlocking == nil || *roundTripFinding.ReleaseBlocking != false {
		t.Fatalf("release_blocking = %v; want explicit false", roundTripFinding.ReleaseBlocking)
	}
	if roundTripFinding.InvariantID != "INV-DEBT-DOCS-001" || roundTripFinding.ProofDepth != "evidence_replay" || roundTripFinding.Owner != "standard" ||
		roundTripFinding.Expiry != "2026-07-01" || roundTripFinding.Remediation == "" || roundTripFinding.Detector != "debtcheck.scanTextMarker" {
		t.Fatalf("round-trip finding = %+v; want optional metadata preserved", roundTripFinding)
	}

	markdown := ToMarkdown(roundTrip)
	for _, want := range []string{
		"invariant_id=INV-DEBT-DOCS-001",
		"release_blocking=false",
		"proof_depth=evidence_replay",
		"owner=standard",
		"expiry=2026-07-01",
		"remediation=remove the marker after publishing replacement guidance",
		"detector=debtcheck.scanTextMarker",
	} {
		if !strings.Contains(markdown, want) {
			t.Fatalf("markdown = %s; want optional metadata %q", markdown, want)
		}
	}
}

func TestReadReportAcceptsFindingsWithoutOptionalMetadata(t *testing.T) {
	root := t.TempDir()
	reportPath := filepath.Join(root, "legacy-debt-report.json")
	writeFile(t, root, "legacy-debt-report.json", `{
  "schema_version": "debt-report/v1",
  "status": "warning",
  "mode": "warn",
  "active_profile": "xlib-standard-debt-v1",
  "score": 9.9,
  "min_score": 9.8,
  "digests": {},
  "summary": {"p0": 0, "p1": 1, "p2": 0},
  "sections": [
    {
      "name": "dependency",
      "status": "warning",
      "p0": 0,
      "p1": 1,
      "p2": 0,
      "findings": [
        {
          "id": "debt.dependency.unpinned-latest",
          "severity": "P1",
          "path": "Makefile",
          "message": "non-documentation file references @latest"
        }
      ]
    }
  ]
}`)

	report, err := ReadReport(reportPath)
	if err != nil {
		t.Fatal(err)
	}
	finding := report.Sections[0].Findings[0]
	if finding.ReleaseBlocking != nil || finding.InvariantID != "" || finding.ProofDepth != "" || finding.Owner != "" ||
		finding.Expiry != "" || finding.Remediation != "" || finding.Detector != "" {
		t.Fatalf("legacy finding = %+v; want absent optional metadata to stay zero-valued", finding)
	}
	markdown := ToMarkdown(report)
	if strings.Contains(markdown, "release_blocking=") || strings.Contains(markdown, "proof_depth=") || strings.Contains(markdown, "owner=") {
		t.Fatalf("legacy markdown = %s; want omitted optional metadata", markdown)
	}
}

func TestSkipPathSkipsMigratedInboxArchive(t *testing.T) {
	root := t.TempDir()

	if !skipPath(root, filepath.Join(root, ".agent", "archive", "inbox", "goal-patch-v1.0-to-v2.2.md")) {
		t.Fatal("skipPath should skip migrated inbox archive files")
	}
	if skipPath(root, filepath.Join(root, ".agent", "archive", "standard", "goal-runtime-canonical.md")) {
		t.Fatal("skipPath should not skip non-inbox archive files")
	}
}

func TestExitCodeCoversModesAndStatuses(t *testing.T) {
	tests := []struct {
		name   string
		report Report
		want   int
	}{
		{
			name:   "enforce passed exits zero",
			report: Report{Mode: "enforce", Status: "passed"},
			want:   0,
		},
		{
			name:   "enforce failed exits one",
			report: Report{Mode: "enforce", Status: "failed"},
			want:   1,
		},
		{
			name:   "enforce warning exits one",
			report: Report{Mode: "enforce", Status: "warning"},
			want:   1,
		},
		{
			name:   "warn failed exits zero",
			report: Report{Mode: "warn", Status: "failed"},
			want:   0,
		},
		{
			name:   "observe failed exits zero",
			report: Report{Mode: "observe", Status: "failed"},
			want:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExitCode(tt.report); got != tt.want {
				t.Fatalf("ExitCode(%+v) = %d, want %d", tt.report, got, tt.want)
			}
		})
	}
}

func TestRunRejectsInvalidModeAndSection(t *testing.T) {
	root := t.TempDir()
	tests := []struct {
		name string
		opts Options
		want string
	}{
		{
			name: "invalid mode",
			opts: Options{Root: root, Mode: "strict", Section: "docs"},
			want: "unsupported debt mode",
		},
		{
			name: "invalid section",
			opts: Options{Root: root, Mode: "enforce", Section: "unknown"},
			want: "unsupported debt section",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := Run(tt.opts); err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("Run() error = %v, want containing %q", err, tt.want)
			}
		})
	}
}

func TestRunMarksMinScoreShortfallByMode(t *testing.T) {
	tests := []struct {
		mode       string
		wantStatus string
		wantExit   int
	}{
		{mode: "enforce", wantStatus: "failed", wantExit: 1},
		{mode: "warn", wantStatus: "warning", wantExit: 0},
		{mode: "observe", wantStatus: "warning", wantExit: 0},
	}

	for _, tt := range tests {
		t.Run(tt.mode, func(t *testing.T) {
			root := t.TempDir()
			writePolicyFiles(t, root)
			writeFile(t, root, "docs/drift.md", "xlib-docs-drift\n")

			report, err := Run(Options{Root: root, Section: "docs", Mode: tt.mode, MinScore: 10})
			if err != nil {
				t.Fatal(err)
			}
			if report.Score >= report.MinScore {
				t.Fatalf("score = %.2f, min_score = %.2f; want shortfall", report.Score, report.MinScore)
			}
			if report.Status != tt.wantStatus {
				t.Fatalf("status = %q, want %q", report.Status, tt.wantStatus)
			}
			if code := ExitCode(report); code != tt.wantExit {
				t.Fatalf("ExitCode = %d, want %d", code, tt.wantExit)
			}
		})
	}
}

func TestRunReportsMissingPolicyFiles(t *testing.T) {
	root := t.TempDir()

	report, err := Run(Options{Root: root, Section: "docs", Mode: "enforce", MinScore: DefaultMinScore})
	if err != nil {
		t.Fatal(err)
	}

	if report.Status != "failed" {
		t.Fatalf("status = %q, want failed", report.Status)
	}
	if report.Summary.P0 != 4 {
		t.Fatalf("summary P0 = %d, want 4", report.Summary.P0)
	}
	if report.Digests.Rules != "missing" || report.Digests.RuleRegistry != "missing" ||
		report.Digests.Exceptions != "missing" || report.Digests.DependencyPurpose != "missing" {
		t.Fatalf("digests = %+v, want missing policy digests", report.Digests)
	}
}

func TestToMarkdownIncludesNoFindingsForCleanSection(t *testing.T) {
	root := t.TempDir()
	writePolicyFiles(t, root)

	report, err := Run(Options{Root: root, Section: "docs", Mode: "enforce", MinScore: DefaultMinScore})
	if err != nil {
		t.Fatal(err)
	}
	markdown := ToMarkdown(report)
	if !strings.Contains(markdown, "No findings.") {
		t.Fatalf("markdown = %s; want clean section no-findings marker", markdown)
	}

	finding := Finding{
		ID:       "debt.policy.missing",
		Severity: "P0",
		Message:  "policy finding without a file path",
	}
	section := SectionReport{Name: "docs", Findings: []Finding{finding}}
	policyMarkdown := ToMarkdown(Report{Sections: []SectionReport{section}})
	if !strings.Contains(policyMarkdown, "debt.policy.missing policy: policy finding without a file path") {
		t.Fatalf("markdown = %s; want empty finding path rendered as policy", policyMarkdown)
	}
}

func TestReadReportRejectsMissingMalformedAndUnsupportedSchema(t *testing.T) {
	root := t.TempDir()

	if _, err := ReadReport(filepath.Join(root, "missing.json")); err == nil {
		t.Fatal("ReadReport missing file error = nil, want error")
	}

	malformedPath := filepath.Join(root, "malformed.json")
	writeFile(t, root, "malformed.json", `{"schema_version":`)
	if _, err := ReadReport(malformedPath); err == nil {
		t.Fatal("ReadReport malformed JSON error = nil, want error")
	}

	unsupportedPath := filepath.Join(root, "unsupported.json")
	writeFile(t, root, "unsupported.json", `{"schema_version":"debt-report/v0"}`)
	if _, err := ReadReport(unsupportedPath); err == nil || !strings.Contains(err.Error(), "unsupported debt report schema") {
		t.Fatalf("ReadReport unsupported schema error = %v, want unsupported schema", err)
	}
}

func TestValidateEvidenceReportsAllBoundaryProblems(t *testing.T) {
	evidence := Evidence{
		SchemaVersion:       "wrong-manifest",
		ReportSchemaVersion: "wrong-report",
		Status:              "failed",
		Score:               9.7,
		Sections: []SectionEvidence{
			{Name: "architecture", Status: "failed", P0: 2},
			{Name: "docs", Status: "warning"},
		},
	}

	got := ValidateEvidence(evidence, 9.8)
	want := []string{
		"debt schema version mismatch",
		"debt report schema version mismatch",
		"debt status is failed",
		"debt score 9.70 below 9.80",
		"debt section architecture has 2 P0 findings",
		"debt section architecture status is failed",
		"debt section docs status is warning",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ValidateEvidence problems = %#v, want %#v", got, want)
	}
}

func TestParseHelpersHandleBoundaryValues(t *testing.T) {
	boolTests := []struct {
		value string
		want  bool
		ok    bool
	}{
		{value: "true", want: true, ok: true},
		{value: "false", want: false, ok: true},
		{value: "TRUE", want: false, ok: false},
		{value: "", want: false, ok: false},
	}
	for _, tt := range boolTests {
		got, ok := parseOptionalBool(tt.value)
		if got != tt.want || ok != tt.ok {
			t.Fatalf("parseOptionalBool(%q) = (%t, %t), want (%t, %t)", tt.value, got, ok, tt.want, tt.ok)
		}
	}

	scalarTests := []struct {
		value string
		want  string
	}{
		{value: "", want: ""},
		{value: "x", want: "x"},
		{value: `"quoted"`, want: "quoted"},
		{value: `'single quoted'`, want: "single quoted"},
		{value: `"unterminated`, want: `"unterminated`},
		{value: `plain`, want: "plain"},
	}
	for _, tt := range scalarTests {
		if got := unquoteYAMLScalar(tt.value); got != tt.want {
			t.Fatalf("unquoteYAMLScalar(%q) = %q, want %q", tt.value, got, tt.want)
		}
	}

	text := `
# adoption_status: adopted
adoption_status: not_adopted
proof_based_adoption: true
proof_based_adoption: "true"
`
	if hasYAMLScalarLine(text, "adoption_status", "adopted") {
		t.Fatal("hasYAMLScalarLine matched commented scalar, want false")
	}
	if !hasYAMLScalarLine(text, "proof_based_adoption", "true") {
		t.Fatal("hasYAMLScalarLine did not match plain true scalar")
	}
	if hasYAMLScalarLine(text, "proof_based_adoption", "missing") {
		t.Fatal("hasYAMLScalarLine matched missing scalar, want false")
	}
}

func TestRuleMetadataParserBoundariesAndAnnotationPreservesExistingFields(t *testing.T) {
	root := t.TempDir()
	if metadata := readRuleMetadata(root, "missing.yaml"); metadata != nil {
		t.Fatalf("readRuleMetadata missing file = %+v, want nil", metadata)
	}

	writeFile(t, root, "registry.yaml", `# comments and malformed lines should be ignored
schema_version: debt-rule-registry/v1
rules:
  -
  - id: 'debt.testing.marker'
    no colon here
    release_blocking: maybe
    owner: "metadata-owner"
    proof_depth: evidence_replay
    unknown: ignored
  - id: debt.docs.marker
    release_blocking: true
    invariant_id: INV-DOCS
`)

	metadata := readRuleMetadata(root, "registry.yaml")
	testingRule := metadata["debt.testing.marker"]
	if testingRule.ID != "debt.testing.marker" || testingRule.Owner != "metadata-owner" || testingRule.ProofDepth != "evidence_replay" {
		t.Fatalf("testing metadata = %+v, want parsed quoted scalar fields", testingRule)
	}
	if testingRule.ReleaseBlocking != nil {
		t.Fatalf("testing release_blocking = %v, want nil for invalid bool", testingRule.ReleaseBlocking)
	}

	existingReleaseBlocking := false
	findings := annotateFindings([]Finding{
		{
			ID:              "debt.docs.marker",
			InvariantID:     "existing-invariant",
			ReleaseBlocking: &existingReleaseBlocking,
			Owner:           "existing-owner",
		},
		{
			ID:      "debt.unknown.marker",
			Path:    "docs/unknown.md",
			Message: "metadata should not be required for every finding",
		},
	}, metadata)
	got := findings[0]
	if got.InvariantID != "existing-invariant" || got.Owner != "existing-owner" ||
		got.ReleaseBlocking == nil || *got.ReleaseBlocking != false {
		t.Fatalf("annotated finding = %+v, want existing fields preserved", got)
	}
	if findings[1].ID != "debt.unknown.marker" || findings[1].Path != "docs/unknown.md" || findings[1].Owner != "" {
		t.Fatalf("unknown finding = %+v, want unchanged finding without metadata", findings[1])
	}
}

func TestStatusTransitions(t *testing.T) {
	tests := []struct {
		name     string
		summary  Summary
		score    float64
		minScore float64
		mode     string
		want     string
	}{
		{name: "P0 always fails", summary: Summary{P0: 1}, score: 10, minScore: 9.8, mode: "observe", want: "failed"},
		{name: "min score shortfall enforces failure", score: 9.7, minScore: 9.8, mode: "enforce", want: "failed"},
		{name: "min score shortfall warns in warn mode", score: 9.7, minScore: 9.8, mode: "warn", want: "warning"},
		{name: "min score shortfall warns in observe mode", score: 9.7, minScore: 9.8, mode: "observe", want: "warning"},
		{name: "P1 is passed in enforce mode above threshold", summary: Summary{P1: 1}, score: 9.9, minScore: 9.8, mode: "enforce", want: "passed"},
		{name: "P2 warns outside enforce mode", summary: Summary{P2: 1}, score: 9.95, minScore: 9.8, mode: "observe", want: "warning"},
		{name: "clean summary passes", score: 10, minScore: 9.8, mode: "enforce", want: "passed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := status(tt.summary, tt.score, tt.minScore, tt.mode); got != tt.want {
				t.Fatalf("status(%+v, %.2f, %.2f, %q) = %q, want %q", tt.summary, tt.score, tt.minScore, tt.mode, got, tt.want)
			}
		})
	}
}

func TestNormalizeValidationAndSectionBoundaries(t *testing.T) {
	defaults := normalize(Options{})
	if defaults.Root != "." || defaults.ConfigPath != DefaultRulesPath || defaults.RegistryPath != DefaultRegistryPath ||
		defaults.ExceptionsPath != DefaultExceptions || defaults.DependencyPurposePath != DefaultPurpose ||
		defaults.Section != "all" || defaults.Mode != "enforce" || defaults.MinScore != DefaultMinScore {
		t.Fatalf("normalize defaults = %+v, want package defaults", defaults)
	}

	custom := Options{
		Root:                  "root",
		ConfigPath:            "rules.yaml",
		RegistryPath:          "registry.yaml",
		ExceptionsPath:        "exceptions.yaml",
		DependencyPurposePath: "purpose.yaml",
		Section:               "docs",
		Mode:                  "warn",
		MinScore:              9.1,
	}
	if got := normalize(custom); !reflect.DeepEqual(got, custom) {
		t.Fatalf("normalize custom = %+v, want %+v", got, custom)
	}
	if err := validateMode("observe"); err != nil {
		t.Fatalf("validateMode observe error = %v, want nil", err)
	}
	if err := validateMode("strict"); err == nil {
		t.Fatal("validateMode strict error = nil, want error")
	}
	if err := validateSection("downstream"); err != nil {
		t.Fatalf("validateSection downstream error = %v, want nil", err)
	}
	if err := validateSection("invalid"); err == nil {
		t.Fatal("validateSection invalid error = nil, want error")
	}
	if sections := selectedSections("security"); !reflect.DeepEqual(sections, []string{"security"}) {
		t.Fatalf("selectedSections(security) = %+v, want security only", sections)
	}
	if got := scanSection("unused", "unknown"); got != nil {
		t.Fatalf("scanSection unknown = %+v, want nil", got)
	}
}

func TestPathBinaryAndSectionHelpers(t *testing.T) {
	for _, name := range []string{".git", ".omx", ".worktree", "vendor", "node_modules", "release", "tmp", ".cache"} {
		if !skipDir(name) {
			t.Fatalf("skipDir(%q) = false, want true", name)
		}
	}
	if skipDir("src") {
		t.Fatal("skipDir(src) = true, want false")
	}

	skipFileTests := []struct {
		path string
		want bool
	}{
		{path: ".secret", want: true},
		{path: ".gitignore", want: false},
		{path: "logo.png", want: true},
		{path: "photo.jpg", want: true},
		{path: "photo.jpeg", want: true},
		{path: "anim.gif", want: true},
		{path: "paper.pdf", want: true},
		{path: "docs/readme.md", want: false},
	}
	for _, tt := range skipFileTests {
		if got := skipFile(tt.path); got != tt.want {
			t.Fatalf("skipFile(%q) = %t, want %t", tt.path, got, tt.want)
		}
	}

	if bytesLookBinary([]byte("plain text")) {
		t.Fatal("bytesLookBinary(plain text) = true, want false")
	}
	if !bytesLookBinary([]byte{'a', 0, 'b'}) {
		t.Fatal("bytesLookBinary(with NUL) = false, want true")
	}
	long := make([]byte, 5000)
	for i := range long {
		long[i] = 'a'
	}
	if bytesLookBinary(long) {
		t.Fatal("bytesLookBinary(long zero-valued buffer) = true, want false before NUL marker test")
	}
	long[4095] = 1
	long[4096] = 0
	if bytesLookBinary(long) {
		t.Fatal("bytesLookBinary should inspect only first 4096 bytes")
	}
	long[4095] = 0
	if !bytesLookBinary(long) {
		t.Fatal("bytesLookBinary should detect NUL within first 4096 bytes")
	}

	root := t.TempDir()
	path := filepath.Join(root, "nested", "file.txt")
	if got := rel(root, path); got != "nested/file.txt" {
		t.Fatalf("rel(root, path) = %q, want nested/file.txt", got)
	}
	if got := rel(root, root); got != "." {
		t.Fatalf("rel(root, root) = %q, want .", got)
	}
	if got := rel("/absolute/root", "relative/path"); got != "relative/path" {
		t.Fatalf("rel(abs, relative) = %q, want original path on Rel error", got)
	}

	section := buildSection("custom", []Finding{
		{Severity: "P0"},
		{Severity: "P1"},
		{Severity: "P3"},
	})
	if section.Status != "failed" || section.P0 != 1 || section.P1 != 1 || section.P2 != 1 {
		t.Fatalf("buildSection = %+v, want failed with P0/P1/P2 counts", section)
	}
}

func TestWalkFilesVisitsSortedFilesSkipsDirsAndPropagatesErrors(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "b.txt", "b")
	writeFile(t, root, "a.txt", "a")
	writeFile(t, root, "nested/c.txt", "c")
	writeFile(t, root, "vendor/skip.txt", "skip")

	var visited []string
	if err := walkFiles(root, func(path string) error {
		visited = append(visited, rel(root, path))
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	want := []string{"a.txt", "b.txt", "nested/c.txt"}
	if !reflect.DeepEqual(visited, want) {
		t.Fatalf("visited = %+v, want sorted non-skipped files %+v", visited, want)
	}

	singlePath := filepath.Join(root, "single.txt")
	writeFile(t, root, "single.txt", "single")
	var singleVisited []string
	if err := walkFiles(singlePath, func(path string) error {
		singleVisited = append(singleVisited, filepath.Base(path))
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(singleVisited, []string{"single.txt"}) {
		t.Fatalf("single visited = %+v, want single file", singleVisited)
	}

	sentinel := errors.New("stop walk")
	if err := walkFiles(root, func(path string) error {
		if rel(root, path) == "b.txt" {
			return sentinel
		}
		return nil
	}); !errors.Is(err, sentinel) {
		t.Fatalf("walkFiles propagated error = %v, want %v", err, sentinel)
	}

	if err := walkFiles(filepath.Join(root, "missing"), func(string) error { return nil }); err == nil {
		t.Fatal("walkFiles missing root error = nil, want error")
	}
	if err := walkDir(filepath.Join(root, "missing-dir"), func(string) error { return nil }); err == nil {
		t.Fatal("walkDir missing directory error = nil, want error")
	}
	if err := walkDir(singlePath, func(string) error { return nil }); err == nil {
		t.Fatal("walkDir file-as-directory error = nil, want error")
	}
}

func TestWalkDirPropagatesUnreadableChildDirectory(t *testing.T) {
	root := t.TempDir()
	child := filepath.Join(root, "child")
	if err := os.Mkdir(child, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(child, 0); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Chmod(child, 0o755); err != nil {
			t.Fatalf("restore child permissions: %v", err)
		}
	}()

	if err := walkDir(root, func(string) error { return nil }); err == nil {
		t.Skip("permission model allowed opening unreadable child directory")
	}
}

func TestWalkDirReportsInjectedDirectoryErrors(t *testing.T) {
	root := t.TempDir()
	oldOpen, oldLstat := debtcheckOpenDir, debtcheckLstat
	t.Cleanup(func() {
		debtcheckOpenDir = oldOpen
		debtcheckLstat = oldLstat
	})

	debtcheckOpenDir = func(string) (dirReader, error) {
		return fakeDirReader{readErr: errors.New("read failed"), closeErr: errors.New("close failed")}, nil
	}
	if err := walkDir(root, func(string) error { return nil }); err == nil || !strings.Contains(err.Error(), "read failed") {
		t.Fatalf("walkDir read error = %v; want read failed", err)
	}

	debtcheckOpenDir = func(string) (dirReader, error) {
		return fakeDirReader{closeErr: errors.New("close failed")}, nil
	}
	if err := walkDir(root, func(string) error { return nil }); err == nil || !strings.Contains(err.Error(), "close failed") {
		t.Fatalf("walkDir close error = %v; want close failed", err)
	}

	debtcheckOpenDir = func(string) (dirReader, error) {
		return fakeDirReader{names: []string{"child.go"}}, nil
	}
	debtcheckLstat = func(string) (os.FileInfo, error) {
		return nil, errors.New("lstat failed")
	}
	if err := walkDir(root, func(string) error { return nil }); err == nil || !strings.Contains(err.Error(), "lstat failed") {
		t.Fatalf("walkDir lstat error = %v; want lstat failed", err)
	}
}

func TestScanHelpersDetectAndSkipBoundaryFiles(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "visible-marker.txt", "xlib-domain-forbidden\n")
	writeFile(t, root, ".hidden", "xlib-domain-forbidden\n")
	writeFile(t, root, "image.png", "xlib-domain-forbidden\n")
	writeBytesFile(t, root, "binary.dat", []byte("xlib-domain-forbidden\x00\n"))
	writeFile(t, root, ".agent/debt/ignored.txt", "xlib-domain-forbidden\n")

	markerFindings := scanTextMarker(root, "xlib-domain-forbidden", "debt.domain.marker", "domain debt marker is present")
	if len(markerFindings) != 1 || markerFindings[0].Path != "visible-marker.txt" {
		t.Fatalf("marker findings = %+v, want only visible text marker", markerFindings)
	}

	writeFile(t, root, "scripts/install.sh", "go install example.com/tool@latest\ncurl https://example.com/install.sh | bash\n")
	writeFile(t, root, "docs/install.md", "go install example.com/tool@latest\n")
	dependencyFindings := scanDependencyDebt(root)
	if len(dependencyFindings) != 2 {
		t.Fatalf("dependency findings = %+v, want unpinned latest and curl-pipe-bash findings", dependencyFindings)
	}
	dependencyIDs := []string{dependencyFindings[0].ID, dependencyFindings[1].ID}
	if !reflect.DeepEqual(dependencyIDs, []string{"debt.dependency.curl-pipe-bash", "debt.dependency.unpinned-latest"}) {
		t.Fatalf("dependency IDs = %+v, want sorted dependency findings", dependencyIDs)
	}

	writeFile(t, root, "secrets.txt", privateKeyPrefix+"\nxlib-security-debt\n")
	securityFindings := scanSecurityDebt(root)
	if len(securityFindings) != 2 {
		t.Fatalf("security findings = %+v, want private-key and marker findings", securityFindings)
	}
}

func TestScanTrackedTextKeepsCollectedFindingsWhenLaterReadFails(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "a-marker.txt", "xlib-docs-drift\n")
	writeFile(t, root, "b-marker.txt", "xlib-docs-drift\n")
	if err := os.Symlink("missing-target", filepath.Join(root, "z-broken.txt")); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}

	findings := scanTextMarker(root, "xlib-docs-drift", "debt.docs.marker", "documentation debt marker is present")
	gotPaths := make([]string, 0, len(findings))
	for _, finding := range findings {
		gotPaths = append(gotPaths, finding.Path)
	}
	wantPaths := []string{"a-marker.txt", "b-marker.txt"}
	if !reflect.DeepEqual(gotPaths, wantPaths) {
		t.Fatalf("marker paths = %+v, want collected findings before read failure %+v", gotPaths, wantPaths)
	}
}

func TestScanGoImportsReportsParseErrorsAndIgnoresTests(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "bad.go", "package")
	writeFile(t, root, "ignored_test.go", "package")
	if err := os.Symlink("missing-target", filepath.Join(root, "broken.go")); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}

	findings := scanGoImports(root)
	if len(findings) != 1 {
		t.Fatalf("findings = %+v, want one parse finding", findings)
	}
	if findings[0].ID != "debt.architecture.parse" || findings[0].Path != "bad.go" {
		t.Fatalf("finding = %+v, want bad.go parse finding", findings[0])
	}
}

func TestScanDownstreamDebtReportsMissingRequiredFiles(t *testing.T) {
	findings := scanDownstreamDebt(t.TempDir())
	if !hasFindingID(findings, "debt.downstream.file-missing") {
		t.Fatalf("findings = %+v, want missing downstream file finding", findings)
	}
}

func TestScanDownstreamDebtReportsSchemaPlaceholderAndGapProblems(t *testing.T) {
	root := t.TempDir()
	writeDownstreamFiles(t, root)
	writeFile(t, root, ".agent/registries/downstream-baseline-scan.yaml", `schema_version: "2.9.3"
repo: kernel/configx
mode: patch-only
status: explicit_missing_repo
`)
	writeFile(t, root, ".agent/registries/downstream-adoption-modes.yaml", `modes: []
# TODO document explicit modes before adoption
`)
	writeFile(t, root, "docs/standard/downstream-compatibility.md", "# Downstream Compatibility\n")

	findings := scanDownstreamDebt(root)
	for _, want := range []string{
		"debt.downstream.baseline-missing-gap-status",
		"debt.downstream.schema-missing",
		"debt.downstream.placeholder",
		"debt.downstream.mode-missing-patch-only",
		"debt.downstream.mode-missing-write-guard",
		"debt.downstream.compatibility-missing-contract",
	} {
		if !hasFindingID(findings, want) {
			t.Fatalf("findings = %+v, want finding id %s", findings, want)
		}
	}
}

func writePolicyFiles(t *testing.T, root string) {
	t.Helper()
	files := map[string]string{
		DefaultRulesPath:    "schema_version: debt-rules/v1\nprofile: test\n",
		DefaultRegistryPath: "schema_version: debt-rule-registry/v1\nrules: []\n",
		DefaultExceptions:   "schema_version: debt-exceptions/v1\nexceptions: []\n",
		DefaultPurpose:      "schema_version: debt-dependency-purpose/v1\npurposes: []\n",
	}
	for path, content := range files {
		writeFile(t, root, path, content)
	}
}

func writeDownstreamFiles(t *testing.T, root string) {
	t.Helper()
	files := map[string]string{
		".agent/registries/downstream-registry.yaml": `schema_version: "2.9.3"
downstreams:
  - repo: kernel/configx
    mode: patch-only
    status: unavailable_in_worker_workspace_gap_explicit
  - repo: kernel/redisx
    mode: patch-only
    status: unavailable_in_worker_workspace_gap_explicit
  - repo: corekit
    mode: patch-only
    status: unavailable_in_worker_workspace_gap_explicit
`,
		".agent/registries/downstream-baseline-scan.yaml": `schema_version: "2.9.3"
repo: kernel/configx
mode: patch-only
status: gap_explicit_when_repo_missing
`,
		".agent/registries/downstream-adoption-modes.yaml": `schema_version: "2.9.3"
modes: [patch-only, dry-run]
forbidden: [direct_downstream_write_without_repo]
`,
		".agent/registries/downstream-adoption-status.yaml": `schema_version: "2.9.3"
current_registry:
  adoption_status: not_adopted
  proof_based_adoption: false
first_pr_mva_assertions:
  no_proof_based_adoption: true
standard_target_libraries:
  - name: kernel
  - name: configx
  - name: observex
  - name: testkitx
  - name: postgresx
  - name: redisx
  - name: kafkax
  - name: natsx
  - name: taosx
  - name: ossx
  - name: clickhousex
`,
		"docs/downstream-matrix.md": `# Downstream Matrix

| Library | Adoption |
| --- | --- |
| ` + "`kernel`" + ` | not_adopted |
| ` + "`configx`" + ` | not_adopted |
| ` + "`observex`" + ` | not_adopted |
| ` + "`testkitx`" + ` | not_adopted |
| ` + "`postgresx`" + ` | not_adopted |
| ` + "`redisx`" + ` | not_adopted |
| ` + "`kafkax`" + ` | not_adopted |
| ` + "`natsx`" + ` | not_adopted |
| ` + "`taosx`" + ` | not_adopted |
| ` + "`ossx`" + ` | not_adopted |
| ` + "`clickhousex`" + ` | not_adopted |
`,
		"docs/standard/downstream-compatibility.md": `# Downstream Compatibility

默认 downstream 为 ` + "`kernel`" + ` 与 ` + "`corekit`" + `。
发布验证命令必须包含 GOWORK=off make integration。
`,
		"scripts/run_integration.sh": `#!/usr/bin/env bash
TARGETS=(
  "kernel|github.com/ZoneCNH/kernel|kernel"
  "configx|github.com/ZoneCNH/configx|configx"
  "redisx|github.com/ZoneCNH/redisx|redisx"
)
GOWORK=off make debt
GOWORK=off make debt-evidence
GOWORK=off make debt-evidence-checksum-check
`,
		"scripts/render_template.sh": `#!/usr/bin/env bash
rsync \
  --exclude release/debt/latest.json \
  --exclude release/debt/latest.md \
  --exclude release/debt/latest.json.sha256
`,
	}
	for path, content := range files {
		writeFile(t, root, path, content)
	}
}

func writeFile(t *testing.T, root, path, content string) {
	t.Helper()
	fullPath := filepath.Join(root, filepath.FromSlash(path))
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func writeBytesFile(t *testing.T, root, path string, content []byte) {
	t.Helper()
	fullPath := filepath.Join(root, filepath.FromSlash(path))
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(fullPath, content, 0o644); err != nil {
		t.Fatal(err)
	}
}

func hasFindingID(findings []Finding, id string) bool {
	for _, finding := range findings {
		if finding.ID == id {
			return true
		}
	}
	return false
}
