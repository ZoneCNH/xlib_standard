package debtcheck

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

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
	report := Report{
		SchemaVersion: SchemaVersion,
		Status:        "warning",
		Mode:          "warn",
		Sections: []SectionReport{{
			Name:   "docs",
			Status: "warning",
			P1:     1,
			Findings: []Finding{{
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
			}},
		}},
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
	finding := roundTrip.Sections[0].Findings[0]
	if finding.ReleaseBlocking == nil || *finding.ReleaseBlocking != false {
		t.Fatalf("release_blocking = %v; want explicit false", finding.ReleaseBlocking)
	}
	if finding.InvariantID != "INV-DEBT-DOCS-001" || finding.ProofDepth != "evidence_replay" || finding.Owner != "standard" ||
		finding.Expiry != "2026-07-01" || finding.Remediation == "" || finding.Detector != "debtcheck.scanTextMarker" {
		t.Fatalf("round-trip finding = %+v; want optional metadata preserved", finding)
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
