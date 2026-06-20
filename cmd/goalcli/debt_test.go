package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ZoneCNH/xlib-standard/internal/debtcheck"
)

func TestRunDebtEvidenceWritesLatestArtifacts(t *testing.T) {
	root := t.TempDir()
	writeDebtCLICleanFixture(t, root)
	chdir(t, root)

	var stdout, stderr bytes.Buffer
	code := runDebtEvidence(nil, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("runDebtEvidence() code = %d stderr = %q", code, stderr.String())
	}
	for _, want := range []string{
		"wrote release/debt/latest.json",
		"wrote release/debt/latest.md",
		"wrote release/debt/latest.json.sha256",
	} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("stdout = %q; want %q", stdout.String(), want)
		}
	}

	data, err := os.ReadFile(filepath.FromSlash("release/debt/latest.json"))
	if err != nil {
		t.Fatalf("read latest debt json: %v", err)
	}
	var report debtcheck.Report
	if err := json.Unmarshal(data, &report); err != nil {
		t.Fatalf("unmarshal latest debt json: %v", err)
	}
	if report.Status != "passed" || report.SchemaVersion != debtcheck.SchemaVersion {
		t.Fatalf("report = %#v; want passed %s", report, debtcheck.SchemaVersion)
	}

	markdown, err := os.ReadFile(filepath.FromSlash("release/debt/latest.md"))
	if err != nil {
		t.Fatalf("read latest debt markdown: %v", err)
	}
	if !strings.Contains(string(markdown), "# Debt Governance Report") {
		t.Fatalf("markdown = %q; want report heading", markdown)
	}
	checksum, err := os.ReadFile(filepath.FromSlash("release/debt/latest.json.sha256"))
	if err != nil {
		t.Fatalf("read latest debt checksum: %v", err)
	}
	if !strings.Contains(string(checksum), "  release/debt/latest.json\n") {
		t.Fatalf("checksum = %q; want latest json path", checksum)
	}
}

func TestRunDebtEvidenceRejectsArguments(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := runDebtEvidence([]string{"unexpected"}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("runDebtEvidence() code = %d; want 2", code)
	}
	if !strings.Contains(stderr.String(), "does not accept arguments") {
		t.Fatalf("stderr = %q; want argument rejection", stderr.String())
	}
}

func TestRunDebtEvidenceReportsDebtCheckError(t *testing.T) {
	root := t.TempDir()
	writeDebtCLICleanFixture(t, root)
	chdir(t, root)

	old := debtCheckRun
	debtCheckRun = func(debtcheck.Options) (debtcheck.Report, error) {
		return debtcheck.Report{}, errors.New("debtcheck failed")
	}
	t.Cleanup(func() { debtCheckRun = old })

	var stdout, stderr bytes.Buffer
	code := runDebtEvidence(nil, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("runDebtEvidence() code = %d stdout = %q stderr = %q; want 2", code, stdout.String(), stderr.String())
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q; want empty output on debt check error", stdout.String())
	}
	if errText := stderr.String(); !strings.Contains(errText, "debtcheck failed") {
		t.Fatalf("stderr = %q; want debt check error", errText)
	}
}

func TestRunDebtReportsMarshalErrors(t *testing.T) {
	root := t.TempDir()
	writeDebtCLICleanFixture(t, root)
	chdir(t, root)

	old := debtMarshalIndent
	debtMarshalIndent = func(any, string, string) ([]byte, error) {
		return nil, errors.New("marshal failed")
	}
	t.Cleanup(func() { debtMarshalIndent = old })

	tests := []struct {
		name string
		run  func(*bytes.Buffer, *bytes.Buffer) int
	}{
		{name: "debt", run: func(stdout *bytes.Buffer, stderr *bytes.Buffer) int {
			return runDebt(nil, stdout, stderr)
		}},
		{name: "helper", run: func(stdout *bytes.Buffer, stderr *bytes.Buffer) int {
			return runDebtHelper("trend", nil, stdout, stderr)
		}},
		{name: "evidence", run: func(stdout *bytes.Buffer, stderr *bytes.Buffer) int {
			return runDebtEvidence(nil, stdout, stderr)
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer
			code := tt.run(&stdout, &stderr)
			if code != 2 {
				t.Fatalf("%s code = %d stdout = %q stderr = %q; want 2", tt.name, code, stdout.String(), stderr.String())
			}
			if stdout.Len() != 0 {
				t.Fatalf("%s stdout = %q, want empty", tt.name, stdout.String())
			}
			if !strings.Contains(stderr.String(), "marshal failed") {
				t.Fatalf("%s stderr = %q; want marshal failure", tt.name, stderr.String())
			}
		})
	}
}

func TestRunDebtEvidenceReportsFilesystemErrors(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(t *testing.T, root string)
		wantErrSub string
	}{
		{
			name: "release path is file",
			setup: func(t *testing.T, root string) {
				writeDebtCLIFile(t, root, "release", "not a directory")
			},
			wantErrSub: "not a directory",
		},
		{
			name: "json path is directory",
			setup: func(t *testing.T, root string) {
				if err := os.MkdirAll(filepath.Join(root, filepath.FromSlash("release/debt/latest.json")), 0o755); err != nil {
					t.Fatalf("mkdir latest.json dir: %v", err)
				}
			},
			wantErrSub: "latest.json",
		},
		{
			name: "markdown path is directory",
			setup: func(t *testing.T, root string) {
				if err := os.MkdirAll(filepath.Join(root, filepath.FromSlash("release/debt/latest.md")), 0o755); err != nil {
					t.Fatalf("mkdir latest.md dir: %v", err)
				}
			},
			wantErrSub: "latest.md",
		},
		{
			name: "checksum path is directory",
			setup: func(t *testing.T, root string) {
				if err := os.MkdirAll(filepath.Join(root, filepath.FromSlash("release/debt/latest.json.sha256")), 0o755); err != nil {
					t.Fatalf("mkdir checksum dir: %v", err)
				}
			},
			wantErrSub: "latest.json.sha256",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			writeDebtCLICleanFixture(t, root)
			tt.setup(t, root)
			chdir(t, root)

			var stdout, stderr bytes.Buffer
			code := runDebtEvidence(nil, &stdout, &stderr)
			if code != 2 {
				t.Fatalf("runDebtEvidence() code = %d stdout = %q stderr = %q; want 2", code, stdout.String(), stderr.String())
			}
			if !strings.Contains(stderr.String(), tt.wantErrSub) {
				t.Fatalf("stderr = %q; want %q", stderr.String(), tt.wantErrSub)
			}
		})
	}
}

func TestRunDebtFormatsAndUnsupportedOutput(t *testing.T) {
	root := t.TempDir()
	writeDebtCLICleanFixture(t, root)
	chdir(t, root)

	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "json", args: []string{"--output", "json"}, want: `"schema_version": "debt-report/v1"`},
		{name: "markdown", args: []string{"--output", "markdown"}, want: "# Debt Governance Report"},
		{name: "md alias", args: []string{"--output", "md"}, want: "# Debt Governance Report"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := runDebt(tt.args, &stdout, &stderr)
			if code != 0 {
				t.Fatalf("runDebt(%v) code = %d stderr = %q", tt.args, code, stderr.String())
			}
			if !strings.Contains(stdout.String(), tt.want) {
				t.Fatalf("stdout = %q; want %q", stdout.String(), tt.want)
			}
		})
	}

	var stdout, stderr bytes.Buffer
	code := runDebt([]string{"--output", "xml"}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("runDebt unsupported output code = %d; want 2", code)
	}
	if !strings.Contains(stderr.String(), "unsupported debt output format") {
		t.Fatalf("stderr = %q; want unsupported output error", stderr.String())
	}
}

func TestRunDebtRejectsFlagParseAndRunErrors(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := runDebt([]string{"--bogus"}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("runDebt unknown flag code = %d; want 2", code)
	}
	if !strings.Contains(stderr.String(), "flag provided but not defined") {
		t.Fatalf("stderr = %q; want flag parse error", stderr.String())
	}

	stdout.Reset()
	stderr.Reset()
	code = runDebt([]string{"--section", "unknown"}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("runDebt invalid section code = %d stdout = %q stderr = %q; want 2", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stderr.String(), "unsupported debt section") {
		t.Fatalf("stderr = %q; want unsupported section error", stderr.String())
	}

	root := t.TempDir()
	chdir(t, root)
	stdout.Reset()
	stderr.Reset()
	code = runDebt(nil, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("runDebt missing fixture code = %d stdout = %q stderr = %q; want 1", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), "debt.rules.missing") {
		t.Fatalf("stdout = %q; want missing policy finding", stdout.String())
	}
}

func TestRunDebtReturnsFailureForFindings(t *testing.T) {
	root := t.TempDir()
	writeDebtCLICleanFixture(t, root)
	writeDebtCLIFile(t, root, "bad.go", "package bad\n\nimport _ \"github.com/ZoneCNH/x.go\"\n")
	chdir(t, root)

	var stdout, stderr bytes.Buffer
	code := runDebt([]string{"--section", "architecture", "--output", "json"}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("runDebt() code = %d stderr = %q stdout = %q; want failing findings", code, stderr.String(), stdout.String())
	}
	if !strings.Contains(stdout.String(), "debt.architecture.legacy-import") {
		t.Fatalf("stdout = %q; want legacy import finding", stdout.String())
	}
}

func TestRunDebtHelperRejectsPositionalArgument(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := runDebtHelper("trend", []string{"unexpected"}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("runDebtHelper() code = %d; want 2", code)
	}
	if !strings.Contains(stderr.String(), "does not accept positional argument") {
		t.Fatalf("stderr = %q; want positional argument rejection", stderr.String())
	}
}

func TestRunDebtHelperRejectsFlagParseRunAndOutputErrors(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := runDebtHelper("trend", []string{"--bogus"}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("runDebtHelper unknown flag code = %d; want 2", code)
	}
	if !strings.Contains(stderr.String(), "flag provided but not defined") {
		t.Fatalf("stderr = %q; want flag parse error", stderr.String())
	}

	stdout.Reset()
	stderr.Reset()
	code = runDebtHelper("trend", []string{"--mode", "strict"}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("runDebtHelper invalid mode code = %d stdout = %q stderr = %q; want 2", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stderr.String(), "unsupported debt mode") {
		t.Fatalf("stderr = %q; want unsupported mode error", stderr.String())
	}

	root := t.TempDir()
	chdir(t, root)
	stdout.Reset()
	stderr.Reset()
	code = runDebtHelper("trend", nil, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("runDebtHelper missing fixture code = %d stdout = %q stderr = %q; want 1", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), "wrote release/debt/trend.json") {
		t.Fatalf("stdout = %q; want helper artifact write", stdout.String())
	}

	writeDebtCLICleanFixture(t, root)
	writeDebtCLIFile(t, root, "blocked", "not a directory")
	stdout.Reset()
	stderr.Reset()
	code = runDebtHelper("trend", []string{"--output", filepath.FromSlash("blocked/artifact.json")}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("runDebtHelper mkdir conflict code = %d stdout = %q stderr = %q; want 2", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stderr.String(), "blocked") {
		t.Fatalf("stderr = %q; want blocked output path error", stderr.String())
	}

	if err := os.MkdirAll(filepath.FromSlash("outdir"), 0o755); err != nil {
		t.Fatalf("mkdir outdir: %v", err)
	}
	stdout.Reset()
	stderr.Reset()
	code = runDebtHelper("trend", []string{"--output", "outdir"}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("runDebtHelper write directory code = %d stdout = %q stderr = %q; want 2", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stderr.String(), "outdir") {
		t.Fatalf("stderr = %q; want output directory write error", stderr.String())
	}
}

func TestRunDebtHelperWritesTrendAndPatchSuggestionArtifacts(t *testing.T) {
	root := t.TempDir()
	writeDebtCLICleanFixture(t, root)
	writeDebtCLIFile(t, root, "bad.go", "package bad\n\nimport _ \"github.com/ZoneCNH/x.go\"\n")
	previous := debtcheck.Report{SchemaVersion: debtcheck.SchemaVersion, Status: "failed", Score: 3.0}
	previousData, err := json.Marshal(previous)
	if err != nil {
		t.Fatalf("marshal previous report: %v", err)
	}
	writeDebtCLIFile(t, root, "release/debt/latest.json", string(previousData))
	chdir(t, root)

	var stdout, stderr bytes.Buffer
	trendPath := filepath.FromSlash("out/trend.json")
	code := runDebtHelper("trend", []string{"--mode", "warn", "--output", trendPath}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("trend helper code = %d stderr = %q", code, stderr.String())
	}
	assertDebtHelperArtifactContains(t, trendPath, "score delta")

	stdout.Reset()
	stderr.Reset()
	patchPath := filepath.FromSlash("out/patch.json")
	code = runDebtHelper("patch-suggest", []string{"--mode", "warn", "--output", patchPath}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("patch helper code = %d stderr = %q", code, stderr.String())
	}
	data, err := os.ReadFile(patchPath)
	if err != nil {
		t.Fatalf("read patch helper artifact: %v", err)
	}
	if !strings.Contains(string(data), "debt.architecture.legacy-import") ||
		!strings.Contains(string(data), "derived patch suggestions") {
		t.Fatalf("patch helper artifact = %s; want suggestion details", data)
	}
}

func TestBuildDebtHelperArtifactRegisterAndLifecycleDetails(t *testing.T) {
	section := debtcheck.SectionReport{Name: "architecture"}
	report := debtcheck.Report{
		SchemaVersion: debtcheck.SchemaVersion,
		Status:        "passed",
		Mode:          "enforce",
		ActiveProfile: "test",
		Score:         10,
		MinScore:      9.8,
		Digests:       debtcheck.Digests{Rules: "rules", RuleRegistry: "registry", Exceptions: "exceptions", DependencyPurpose: "purpose", Report: "report"},
		Summary:       debtcheck.Summary{},
		Sections:      []debtcheck.SectionReport{section},
	}

	register := buildDebtHelperArtifact("register-update", report)
	if register.Command != "register-update" || register.Status != "passed" || register.Digests.Rules != "rules" {
		t.Fatalf("register artifact = %#v; want copied report fields", register)
	}
	for _, want := range []string{
		"captured debt governance registry state",
		"digest rules=rules",
		"digest rule_registry=registry",
		"digest exceptions=exceptions",
		"digest dependency_purpose=purpose",
	} {
		if !strings.Contains(strings.Join(register.Details, "\n"), want) {
			t.Fatalf("register details = %#v; want %q", register.Details, want)
		}
	}

	lifecycle := buildDebtHelperArtifact("lifecycle-check", report)
	for _, want := range []string{
		"validated debt policy inputs and current report lifecycle",
		"score 10.00 minimum 9.80",
		"digest report=report",
	} {
		if !strings.Contains(strings.Join(lifecycle.Details, "\n"), want) {
			t.Fatalf("lifecycle details = %#v; want %q", lifecycle.Details, want)
		}
	}
}

func TestDebtTrendDetailsBranches(t *testing.T) {
	root := t.TempDir()
	chdir(t, root)
	report := debtcheck.Report{Status: "passed", Score: 7.5}
	if got := strings.Join(debtTrendDetails(report), "\n"); !strings.Contains(got, "no prior debt evidence") {
		t.Fatalf("debtTrendDetails missing prior = %q", got)
	}

	writeDebtCLIFile(t, root, "release/debt/latest.json", "not-json")
	if got := strings.Join(debtTrendDetails(report), "\n"); !strings.Contains(got, "is not a debt report") {
		t.Fatalf("debtTrendDetails invalid prior = %q", got)
	}

	previous := debtcheck.Report{Status: "failed", Score: 3.0}
	data, err := json.Marshal(previous)
	if err != nil {
		t.Fatalf("marshal previous report: %v", err)
	}
	writeDebtCLIFile(t, root, "release/debt/latest.json", string(data))
	got := strings.Join(debtTrendDetails(report), "\n")
	if !strings.Contains(got, "previous status failed score 3.00") ||
		!strings.Contains(got, "score delta 4.50") {
		t.Fatalf("debtTrendDetails valid prior = %q", got)
	}
}

func TestDebtPatchSuggestionsCapsAndFallback(t *testing.T) {
	if got := debtPatchSuggestions(debtcheck.Report{}); len(got) != 1 ||
		got[0] != "no patch suggestions; current debt report has no findings" {
		t.Fatalf("debtPatchSuggestions empty = %#v", got)
	}

	findings := make([]debtcheck.Finding, 25)
	for i := range findings {
		findings[i] = debtcheck.Finding{
			ID:       fmt.Sprintf("finding-%02d", i),
			Severity: "P2",
			Path:     fmt.Sprintf("path-%02d.go", i),
			Message:  "fixture finding",
		}
	}
	section := debtcheck.SectionReport{
		Name:     "architecture",
		Findings: findings,
	}
	got := debtPatchSuggestions(debtcheck.Report{Sections: []debtcheck.SectionReport{section}})
	if len(got) != 20 {
		t.Fatalf("debtPatchSuggestions len = %d; want 20", len(got))
	}
	if !strings.Contains(got[0], "architecture: address P2 finding finding-00 at path-00.go") {
		t.Fatalf("first suggestion = %q", got[0])
	}
	if strings.Contains(strings.Join(got, "\n"), "finding-20") {
		t.Fatalf("suggestions = %#v; want capped before finding-20", got)
	}
}

func TestRunDebtAliasPrependsSectionAndMode(t *testing.T) {
	root := t.TempDir()
	writeDebtCLICleanFixture(t, root)
	writeDebtCLIFile(t, root, "bad.go", "package bad\n\nimport _ \"github.com/ZoneCNH/x.go\"\n")
	chdir(t, root)

	var stdout, stderr bytes.Buffer
	code := runDebtAlias("architecture", "warn", []string{"--output", "json"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("runDebtAlias() code = %d stderr = %q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"mode": "warn"`) ||
		!strings.Contains(stdout.String(), "debt.architecture.legacy-import") {
		t.Fatalf("stdout = %q; want warn-mode architecture finding", stdout.String())
	}
}

func assertDebtHelperArtifactContains(t *testing.T, path string, want string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read helper artifact %s: %v", path, err)
	}
	if !strings.Contains(string(data), want) {
		t.Fatalf("helper artifact %s = %s; want %q", path, data, want)
	}
}

func writeDebtCLICleanFixture(t *testing.T, root string) {
	t.Helper()
	writeDebtCLIPolicyFiles(t, root)
	writeDebtCLIDownstreamFiles(t, root)
	writeDebtCLIFile(t, root, "go.mod", "module github.com/ZoneCNH/xlib-standard\n")
}

func writeDebtCLIPolicyFiles(t *testing.T, root string) {
	t.Helper()
	writeDebtCLIFile(t, root, debtcheck.DefaultRulesPath, "schema_version: debt-rules/v1\nprofile: test\n")
	writeDebtCLIFile(t, root, debtcheck.DefaultExceptions, "schema_version: debt-exceptions/v1\nexceptions: []\n")
	writeDebtCLIFile(t, root, debtcheck.DefaultPurpose, "schema_version: debt-dependency-purpose/v1\npurposes: []\n")
	writeDebtCLIFile(t, root, debtcheck.DefaultRegistryPath, `schema_version: debt-rule-registry/v1
rules:
  - id: debt.architecture.legacy-import
    invariant_id: INV-ARCH-001
    release_blocking: true
    proof_depth: runtime
    owner: architecture
    expiry: 2026-12-31
    remediation: remove legacy import
    detector: go-ast
`)
}

func writeDebtCLIDownstreamFiles(t *testing.T, root string) {
	t.Helper()
	writeDebtCLIFile(t, root, ".agent/registries/downstream-registry.yaml", `schema_version: "2.9.3"
downstreams:
  - repo: kernel/configx
    mode: patch-only
    status: repo_unavailable_gap
  - repo: kernel/redisx
    mode: patch-only
    status: repo_unavailable_gap
  - repo: corekit
    mode: patch-only
    status: repo_unavailable_gap
`)
	writeDebtCLIFile(t, root, ".agent/registries/downstream-baseline-scan.yaml", `schema_version: "2.9.3"
repo: kernel/configx
mode: patch-only
status: repo_missing_gap
`)
	writeDebtCLIFile(t, root, ".agent/registries/downstream-adoption-modes.yaml", `schema_version: "2.9.3"
modes: [patch-only, dry-run]
forbidden: [direct_downstream_write_without_repo]
`)
	writeDebtCLIFile(t, root, ".agent/registries/downstream-adoption-status.yaml", `schema_version: "2.9.3"
targets:
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
status: adoption_status/proof_based
`)
	writeDebtCLIFile(t, root, "docs/downstream-matrix.md", "`kernel` `configx` `observex` `testkitx` `postgresx` `redisx` `kafkax` `natsx` `taosx` `ossx` `clickhousex`\n")
	writeDebtCLIFile(t, root, "docs/standard/downstream-compatibility.md", "`kernel` `corekit` GOWORK=off make integration\n")
	writeDebtCLIFile(t, root, "scripts/run_integration.sh", `#!/usr/bin/env bash
set -euo pipefail
kernel|github.com/ZoneCNH/kernel|kernel
configx|github.com/ZoneCNH/configx|configx
redisx|github.com/ZoneCNH/redisx|redisx
GOWORK=off make debt
GOWORK=off make debt-evidence
GOWORK=off make debt-evidence-checksum-check
`)
	writeDebtCLIFile(t, root, "scripts/render_template.sh", `#!/usr/bin/env bash
set -euo pipefail
release/debt/latest.json
release/debt/latest.md
release/debt/latest.json.sha256
`)
}

func writeDebtCLIFile(t *testing.T, root, rel, content string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", rel, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", rel, err)
	}
}
