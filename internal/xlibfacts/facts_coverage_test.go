// SPDX-License-Identifier: Apache-2.0
package xlibfacts

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadDefaultsRootToCwdAndParsesFacts(t *testing.T) {
	root := t.TempDir()
	factsPath := filepath.Join(root, filepath.FromSlash(Path))
	if err := os.MkdirAll(filepath.Dir(factsPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(factsPath, []byte(expectedFactsYAML()), 0o644); err != nil {
		t.Fatal(err)
	}

	// chdir into root so empty root arg resolves to "." (covers `if root == ""` branch).
	previous, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(previous) })

	facts, err := Load("")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if facts.Module != Module {
		t.Fatalf("module = %q; want %q", facts.Module, Module)
	}
}

func TestLoadExplicitRootReadsFacts(t *testing.T) {
	root := t.TempDir()
	factsPath := filepath.Join(root, filepath.FromSlash(Path))
	if err := os.MkdirAll(filepath.Dir(factsPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(factsPath, []byte(expectedFactsYAML()), 0o644); err != nil {
		t.Fatal(err)
	}

	facts, err := Load(root)
	if err != nil {
		t.Fatalf("Load(%q) error = %v", root, err)
	}
	if gaps := facts.Validate(); len(gaps) > 0 {
		t.Fatalf("Validate() gaps = %v", gaps)
	}
}

func TestLoadReportsMissingFile(t *testing.T) {
	root := t.TempDir()
	if _, err := Load(root); err == nil {
		t.Fatal("Load() returned nil error for missing facts file")
	}
}

func TestParseRejectsLineWithoutColon(t *testing.T) {
	if _, err := Parse([]byte("not a yaml mapping line\n")); err == nil {
		t.Fatal("Parse() returned nil error for line without colon")
	}
}

func TestParseRejectsCommentOnlyAndBlankLines(t *testing.T) {
	// Comments and blank lines are skipped, so a minimal valid document with
	// only a schema_version top-level key parses successfully.
	data := []byte("# leading comment\n\nschema_version: xlib-facts/v1\n")
	facts, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if facts.SchemaVersion != "xlib-facts/v1" {
		t.Fatalf("schema_version = %q; want xlib-facts/v1", facts.SchemaVersion)
	}
}

func TestParseRejectsUnknownTopLevelKey(t *testing.T) {
	if _, err := Parse([]byte("bogus_top_level: value\n")); err == nil {
		t.Fatal("Parse() returned nil error for unknown top-level key")
	}
}

func TestParseRejectsUnsupportedIndentation(t *testing.T) {
	data := []byte("current_release:\n    version: v1.0.1\n")
	if _, err := Parse(data); err == nil {
		t.Fatal("Parse() returned nil error for unsupported indentation")
	}
}

func TestParseRejectsUnknownCurrentReleaseKey(t *testing.T) {
	data := []byte("current_release:\n  bogus_release_key: value\n")
	if _, err := Parse(data); err == nil {
		t.Fatal("Parse() returned nil error for unknown current_release key")
	}
}

func TestParseRejectsUnknownRuntimeKey(t *testing.T) {
	data := []byte("runtime:\n  bogus_runtime_key: value\n")
	if _, err := Parse(data); err == nil {
		t.Fatal("Parse() returned nil error for unknown runtime key")
	}
}

func TestParseRejectsUnknownToolsKey(t *testing.T) {
	data := []byte("tools:\n  bogus_tools_key: value\n")
	if _, err := Parse(data); err == nil {
		t.Fatal("Parse() returned nil error for unknown tools key")
	}
}

func TestParseRejectsKeyOutsideSupportedSection(t *testing.T) {
	// Top-level section header followed by an indented key, but the section
	// header is unknown so the indented key falls through to the default case.
	data := []byte("unknown_section:\n  key: value\n")
	if _, err := Parse(data); err == nil {
		t.Fatal("Parse() returned nil error for key outside supported section")
	}
}

func TestValidateFlagsMalformedReleasedAt(t *testing.T) {
	facts := Expected()
	facts.CurrentRelease.ReleasedAt = "2026-06-18 not-rfc3339"

	gaps := facts.Validate()
	found := false
	for _, gap := range gaps {
		if strings.Contains(gap, "current_release.released_at must be RFC3339") {
			found = true
		}
	}
	if !found {
		t.Fatalf("gaps = %v; want RFC3339 violation", gaps)
	}
}

func TestValidateFlagsMissingRequiredFields(t *testing.T) {
	// Empty Facts: every require() call should append a missing-* gap,
	// covering the require() closure body.
	gaps := Facts{}.Validate()
	requiredFields := []string{
		"missing schema_version",
		"missing module",
		"missing current_release.version",
		"missing current_release.commit",
		"missing current_release.released_at",
		"missing runtime.goal_runtime_version",
		"missing runtime.governance_runtime_version",
		"missing tools.go",
		"missing tools.golangci_lint",
		"missing tools.govulncheck",
	}
	for _, want := range requiredFields {
		if !containsGap(gaps, want) {
			t.Fatalf("gaps = %v; want %q", gaps, want)
		}
	}
}

// TestParseReportsScannerBufferOverflow covers the scanner.Err() branch by
// feeding a single line longer than bufio.Scanner's default 64 KiB token limit.
func TestParseReportsScannerBufferOverflow(t *testing.T) {
	huge := make([]byte, 0, 70*1024)
	huge = append(huge, []byte("schema_version: ")...)
	for len(huge) < 70*1024 {
		huge = append(huge, 'x')
	}
	if _, err := Parse(huge); err == nil {
		t.Fatal("Parse() returned nil error for oversized line")
	}
}

func containsGap(gaps []string, want string) bool {
	for _, gap := range gaps {
		if strings.Contains(gap, want) {
			return true
		}
	}
	return false
}

func expectedFactsYAML() string {
	return strings.Join([]string{
		"schema_version: " + SchemaVersion,
		"module: " + Module,
		"current_release:",
		"  version: " + CurrentReleaseVersion,
		"  commit: " + CurrentReleaseCommit,
		"  released_at: " + CurrentReleaseReleasedAt,
		"runtime:",
		"  goal_runtime_version: " + GoalRuntimeVersion,
		"  governance_runtime_version: " + GovernanceRuntimeVersion,
		"tools:",
		`  go: "` + GoVersion + `"`,
		"  golangci_lint: " + GolangCILintVersion,
		"  govulncheck: " + GovulncheckVersion,
		"",
	}, "\n")
}
