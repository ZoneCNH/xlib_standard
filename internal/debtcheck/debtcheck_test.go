package debtcheck

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunPassesFixtureAndWritesEvidence(t *testing.T) {
	repo := writeDebtFixtureRepo(t)

	report, err := Run(Options{Root: repo})
	if err != nil {
		t.Fatal(err)
	}
	if report.Status != "passed" {
		t.Fatalf("status = %q, want passed: %+v", report.Status, report.Checks)
	}
	if report.Score != DefaultMinScore || len(report.Checks) == 0 {
		t.Fatalf("report = %+v, want default score and checks", report)
	}
	if err := StatusError(report); err != nil {
		t.Fatalf("StatusError = %v, want nil", err)
	}

	if err := WriteEvidence(report, Options{Root: repo}); err != nil {
		t.Fatal(err)
	}
	for _, rel := range []string{
		"release/debt/latest.json",
		"release/debt/latest.md",
		"release/debt/latest.json.sha256",
	} {
		if _, err := os.Stat(filepath.Join(repo, filepath.FromSlash(rel))); err != nil {
			t.Fatalf("stat %s: %v", rel, err)
		}
	}
	checksum, err := os.ReadFile(filepath.Join(repo, "release", "debt", "latest.json.sha256"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(checksum), "release/debt/latest.json") {
		t.Fatalf("checksum = %q, want evidence path", checksum)
	}
}

func TestRunFailsMissingPolicy(t *testing.T) {
	repo := writeDebtFixtureRepo(t)
	if err := os.Remove(filepath.Join(repo, ".agent", "debt", "rules.yaml")); err != nil {
		t.Fatal(err)
	}

	report, err := Run(Options{Root: repo})
	if err != nil {
		t.Fatal(err)
	}
	if report.Status != "failed" {
		t.Fatalf("status = %q, want failed", report.Status)
	}
	if err := StatusError(report); err == nil {
		t.Fatal("StatusError returned nil, want failure")
	}
}

func TestRunFailsP0ExceptionMarker(t *testing.T) {
	repo := writeDebtFixtureRepo(t)
	if err := os.WriteFile(filepath.Join(repo, ".agent", "debt", "rules.yaml"), []byte("allow_p0_exceptions: true\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	report, err := Run(Options{Root: repo})
	if err != nil {
		t.Fatal(err)
	}
	if report.Status != "failed" {
		t.Fatalf("status = %q, want failed", report.Status)
	}
}

func writeDebtFixtureRepo(t *testing.T) string {
	t.Helper()
	repo := t.TempDir()
	write := func(rel, contents string) {
		t.Helper()
		path := filepath.Join(repo, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	write(".agent/debt/rules.yaml", "schema_version: 1\ndebt_policy: strict\nseverity_floor: p0\n")
	write(".agent/debt/rule-registry.yaml", "rules:\n  - id: debt.p0.no_exceptions\n")
	write(".agent/debt/profile.yaml", "profile: default\ndebt: governed\n")
	write(".agent/debt/register.md", "# Debt Register\n\ndebt entries are tracked here.\n")
	write("Makefile", strings.Join([]string{
		"test:",
		"race:",
		"property:",
		"golden:",
		"fuzz-smoke:",
		"debt:",
		"debt-evidence:",
		"implementation-debt:",
		"",
	}, "\n"))
	write("cmd/xlibgate/main.go", "package main\n// debt debt-evidence\n")
	write(".agent/command-registry.yaml", "commands:\n  - debt\n  - implementation-debt\n")
	write(".agent/makefile-baseline.yaml", "targets:\n  - debt\n  - debt-evidence\n")
	write(".agent/makefile-target-registry.yaml", "targets:\n  - debt\n  - debt-evidence\n")
	write("docs/standard/xlibgate-cli-contract.md", "# CLI\n\ndebt\ndebt-evidence\n")
	write(".agent/downstream-registry.yaml", "targets:\n  - kernel/configx\n  - kernel/redisx\n  - corekit\n")
	write("scripts/check_boundary.sh", "#!/bin/sh\nexit 0\n")
	write("scripts/check_docs.sh", "#!/bin/sh\nexit 0\n")
	write("scripts/check_dependency_diff.sh", "#!/bin/sh\nexit 0\n")
	write("scripts/check_secrets.sh", "#!/bin/sh\nexit 0\n")
	write("scripts/run_integration.sh", "#!/bin/sh\n# kernel corekit\n")
	return repo
}
