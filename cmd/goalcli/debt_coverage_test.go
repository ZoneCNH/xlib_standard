package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ZoneCNH/xlib-standard/internal/debtcheck"
)

// TestRunDebtOutputFormatErrors covers invalid output and parse errors.
func TestRunDebtOutputFormatErrors(t *testing.T) {
	chdir(t, repoRoot(t))
	t.Run("invalid output format", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		got := runDebt([]string{"--output", "xml"}, &stdout, &stderr)
		if got != 2 {
			t.Fatalf("got = %d; want 2", got)
		}
		if !strings.Contains(stderr.String(), "unsupported debt output format") {
			t.Fatalf("stderr = %q; want unsupported format", stderr.String())
		}
	})
	t.Run("flag parse error", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		got := runDebt([]string{"--min-score", "notanumber"}, &stdout, &stderr)
		if got != 2 {
			t.Fatalf("got = %d; want 2", got)
		}
	})
	t.Run("markdown output", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		got := runDebt([]string{"--output", "markdown"}, &stdout, &stderr)
		if got != 0 {
			t.Fatalf("got = %d; want 0", got)
		}
		if !strings.Contains(stdout.String(), "# Debt Governance Report") {
			t.Fatalf("stdout = %q; want markdown header", stdout.String())
		}
	})
}

// TestRunDebtHelperBranches covers parse error, positional args, config error, mkdir error.
func TestRunDebtHelperBranches(t *testing.T) {
	t.Run("parse error", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		got := runDebtHelper("trend", []string{"--min-score", "x"}, &stdout, &stderr)
		if got != 2 {
			t.Fatalf("got = %d; want 2", got)
		}
	})
	t.Run("positional arg rejected", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		got := runDebtHelper("trend", []string{"positional"}, &stdout, &stderr)
		if got != 2 {
			t.Fatalf("got = %d; want 2", got)
		}
		if !strings.Contains(stderr.String(), "does not accept positional argument") {
			t.Fatalf("stderr = %q", stderr.String())
		}
	})
	t.Run("config error", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		got := runDebtHelper("trend", []string{"--config", "/nonexistent/path.yaml"}, &stdout, &stderr)
		if got == 0 {
			t.Fatalf("got = %d; want non-zero", got)
		}
	})
	t.Run("output mkdir error", func(t *testing.T) {
		// Point output at a path whose parent is a regular file.
		root := t.TempDir()
		blocker := filepath.Join(root, "blocker")
		if err := os.WriteFile(blocker, []byte("x"), 0o644); err != nil {
			t.Fatalf("write: %v", err)
		}
		chdir(t, repoRoot(t))
		var stdout, stderr bytes.Buffer
		got := runDebtHelper("trend", []string{"--output", filepath.Join(blocker, "out.json")}, &stdout, &stderr)
		if got != 2 {
			t.Fatalf("got = %d; want 2", got)
		}
	})
}

// TestBuildDebtHelperArtifactAllCommands verifies the switch branches.
func TestBuildDebtHelperArtifactAllCommands(t *testing.T) {
	report := debtcheck.Report{Status: "passed", Score: 8.5, MinScore: 7.0, Mode: "enforce"}
	for _, cmd := range []string{"register-update", "trend", "patch-suggest", "lifecycle-check", "unknown"} {
		artifact := buildDebtHelperArtifact(cmd, report)
		if artifact.SchemaVersion != "debt-helper/v1" {
			t.Errorf("cmd %s schema = %q", cmd, artifact.SchemaVersion)
		}
		if artifact.Command != cmd {
			t.Errorf("cmd %s command field = %q", cmd, artifact.Command)
		}
	}
}

// TestDebtTrendCoverageBranches covers missing prior, invalid prior, valid prior.
func TestDebtTrendCoverageBranches(t *testing.T) {
	report := debtcheck.Report{Status: "passed", Score: 8.0, MinScore: 7.0}

	t.Run("no prior evidence", func(t *testing.T) {
		root := t.TempDir()
		chdir(t, root)
		details := debtTrendDetails(report)
		if !slicesContain(details, "no prior debt evidence found at release/debt/latest.json") {
			t.Fatalf("details = %v; want no-prior", details)
		}
	})

	t.Run("invalid prior evidence", func(t *testing.T) {
		root := t.TempDir()
		chdir(t, root)
		writeLatestDebt(t, root, "{not valid json")
		details := debtTrendDetails(report)
		if !gapsContainSubstring(details, "is not a debt report") {
			t.Fatalf("details = %v; want not-a-debt-report", details)
		}
	})

	t.Run("valid prior evidence produces delta", func(t *testing.T) {
		root := t.TempDir()
		chdir(t, root)
		writeLatestDebt(t, root, `{"status":"failed","score":6.0}`)
		details := debtTrendDetails(report)
		found := false
		for _, d := range details {
			if strings.Contains(d, "score delta") {
				found = true
			}
		}
		if !found {
			t.Fatalf("details = %v; want score delta", details)
		}
	})
}

func writeLatestDebt(t *testing.T, root, content string) {
	t.Helper()
	path := filepath.Join(root, "release", "debt", "latest.json")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
}

// TestDebtPatchSuggestions covers empty and non-empty findings.
func TestDebtPatchSuggestions(t *testing.T) {
	t.Run("no findings", func(t *testing.T) {
		report := debtcheck.Report{}
		suggestions := debtPatchSuggestions(report)
		if len(suggestions) != 1 || !strings.Contains(suggestions[0], "no patch suggestions") {
			t.Fatalf("suggestions = %v; want no-patch message", suggestions)
		}
	})
	t.Run("with findings", func(t *testing.T) {
		report := debtcheck.Report{
			Sections: []debtcheck.SectionReport{
				{
					Name: "architecture",
					Findings: []debtcheck.Finding{
						{ID: "F1", Severity: "high", Path: "docs/x.md", Message: "fix it"},
					},
				},
			},
		}
		suggestions := debtPatchSuggestions(report)
		if len(suggestions) == 0 {
			t.Fatalf("suggestions empty; want at least one")
		}
		if !strings.Contains(suggestions[0], "F1") {
			t.Fatalf("suggestions = %v; want F1", suggestions)
		}
	})
}

// TestRunDebtEvidence covers the evidence command happy path and arg rejection.
func TestRunDebtEvidence(t *testing.T) {
	t.Run("rejects args", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		got := runDebtEvidence([]string{"extra"}, &stdout, &stderr)
		if got != 2 {
			t.Fatalf("got = %d; want 2", got)
		}
	})
	t.Run("writes evidence files", func(t *testing.T) {
		root := t.TempDir()
		chdir(t, root)
		var stdout, stderr bytes.Buffer
		got := runDebtEvidence(nil, &stdout, &stderr)
		// Exit code reflects debtcheck.ExitCode which may be non-zero in a temp dir;
		// the important contract is that all three evidence files are written.
		_ = got
		if !strings.Contains(stdout.String(), "wrote release/debt/latest.json") {
			t.Fatalf("stdout = %q; want write confirmation", stdout.String())
		}
		for _, rel := range []string{"release/debt/latest.json", "release/debt/latest.md", "release/debt/latest.json.sha256"} {
			if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(rel))); err != nil {
				t.Fatalf("missing %s: %v", rel, err)
			}
		}
	})
}

// TestDefaultDebtHelperOutput covers the helper.
func TestDefaultDebtHelperOutput(t *testing.T) {
	got := defaultDebtHelperOutput("trend")
	if !strings.HasSuffix(filepath.ToSlash(got), "release/debt/trend.json") {
		t.Fatalf("got = %q; want release/debt/trend.json", got)
	}
}

// TestRunDebtAlias covers the alias wrapper.
func TestRunDebtAlias(t *testing.T) {
	chdir(t, repoRoot(t))
	var stdout, stderr bytes.Buffer
	got := runDebtAlias("architecture", "enforce", nil, &stdout, &stderr)
	if got != 0 {
		t.Fatalf("got = %d; want 0", got)
	}
	if !strings.Contains(stdout.String(), `"name": "architecture"`) {
		t.Fatalf("stdout = %q; want architecture section", stdout.String())
	}
}
