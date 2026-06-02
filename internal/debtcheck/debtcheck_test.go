package debtcheck

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestEvaluateDebtPassesWithPolicyAndAnchors(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root)
	report := Evaluate(root, "debt")
	if report.Status != "passed" {
		t.Fatalf("status = %s, gaps = %#v", report.Status, report.Gaps)
	}
	if report.GateStatuses["debt"] != "passed" {
		t.Fatalf("debt gate status = %q", report.GateStatuses["debt"])
	}
	for _, gate := range GateNames {
		if report.GateStatuses[gate] != "passed" {
			t.Fatalf("gate %s status = %q", gate, report.GateStatuses[gate])
		}
	}
}

func TestEvaluateFailsClosedForMissingPolicy(t *testing.T) {
	root := t.TempDir()
	write(t, root, "go.mod", "module example.com/x\n")
	report := Evaluate(root, "dependency-debt")
	if report.Status != "failed" {
		t.Fatalf("status = %s, want failed", report.Status)
	}
	if len(report.Gaps) == 0 {
		t.Fatal("expected fail-closed gaps")
	}
}

func TestWriteEvidenceWritesDeterministicArtifacts(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root)
	report := Evaluate(root, "debt")
	outDir := filepath.Join(root, "release", "debt")
	if err := WriteEvidence(root, outDir, report); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(outDir, "latest.json"))
	if err != nil {
		t.Fatal(err)
	}
	var got Report
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if got.Runtime != Runtime || got.SchemaVersion != SchemaVersion || got.Status != "passed" {
		t.Fatalf("unexpected evidence: %#v", got)
	}
	for _, name := range []string{"latest.md", "latest.json.sha256"} {
		if _, err := os.Stat(filepath.Join(outDir, name)); err != nil {
			t.Fatalf("missing %s: %v", name, err)
		}
	}
}

func writeFixture(t *testing.T, root string) {
	t.Helper()
	for _, path := range policyFiles {
		content := "schema_version: \"1.0\"\n"
		if path == ".agent/debt/rules.yaml" {
			content += "p0_exceptions: forbidden\n"
		}
		if path == ".agent/debt/rule-registry.yaml" {
			content += "fail_closed: true\n"
		}
		write(t, root, path, content)
	}
	for _, anchors := range gateAnchors {
		for _, path := range anchors {
			write(t, root, path, "fixture\n")
		}
	}
}

func write(t *testing.T, root, path, content string) {
	t.Helper()
	full := filepath.Join(root, filepath.FromSlash(path))
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
