package debtcheck

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunPassesWithPolicyFilesAndCleanTree(t *testing.T) {
	root := t.TempDir()
	writePolicyFiles(t, root)
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
