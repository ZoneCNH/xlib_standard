// SPDX-License-Identifier: Apache-2.0
package xlibfacts

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const expectedFactsYAML = `schema_version: xlib-facts/v1
module: github.com/ZoneCNH/xlib-standard
current_release:
  version: v1.0.2
  commit: 26792dc01317794fb337a0dc81bd732285e49100
  released_at: 2026-06-20T00:00:00Z
runtime:
  goal_runtime_version: v3.1
  governance_runtime_version: v2.9.3
tools:
  go: "1.23.0"
  golangci_lint: "v2.1.6"
  govulncheck: "v1.1.4"
`

func TestParseExpectedFacts(t *testing.T) {
	facts, err := Parse([]byte(expectedFactsYAML))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if gaps := facts.Validate(); len(gaps) > 0 {
		t.Fatalf("Validate() gaps = %v", gaps)
	}
	if gaps := DriftGaps(facts, Expected()); len(gaps) > 0 {
		t.Fatalf("DriftGaps() = %v", gaps)
	}
}

func TestParseIgnoresBlankAndCommentLines(t *testing.T) {
	input := "\n# generated fixture\n" + expectedFactsYAML + "\n# trailing comment\n\n"

	facts, err := Parse([]byte(input))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if gaps := DriftGaps(facts, Expected()); len(gaps) > 0 {
		t.Fatalf("DriftGaps() = %v", gaps)
	}
}

func TestDriftGapsReportsReleaseVersion(t *testing.T) {
	facts := Expected()
	facts.CurrentRelease.Version = "v0.4.14"

	gaps := DriftGaps(facts, Expected())
	if len(gaps) != 1 || !strings.Contains(gaps[0], "current_release.version") {
		t.Fatalf("DriftGaps() = %v; want current_release.version drift", gaps)
	}
}

func TestLoadReadsFactsFile(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, filepath.FromSlash(Path))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir facts path: %v", err)
	}
	if err := os.WriteFile(path, []byte(expectedFactsYAML), 0o644); err != nil {
		t.Fatalf("write facts fixture: %v", err)
	}

	facts, err := Load(root)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if gaps := facts.Validate(); len(gaps) > 0 {
		t.Fatalf("Validate() gaps = %v", gaps)
	}
	if gaps := DriftGaps(facts, Expected()); len(gaps) > 0 {
		t.Fatalf("DriftGaps() = %v", gaps)
	}
}

func TestLoadDefaultsToCurrentDirectory(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, filepath.FromSlash(Path))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir facts path: %v", err)
	}
	if err := os.WriteFile(path, []byte(expectedFactsYAML), 0o644); err != nil {
		t.Fatalf("write facts fixture: %v", err)
	}
	old, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatalf("chdir %s: %v", root, err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(old); err != nil {
			t.Fatalf("restore cwd %s: %v", old, err)
		}
	})

	facts, err := Load("")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if gaps := DriftGaps(facts, Expected()); len(gaps) > 0 {
		t.Fatalf("DriftGaps() = %v", gaps)
	}
}

func TestLoadReturnsReadError(t *testing.T) {
	_, err := Load(t.TempDir())
	if err == nil {
		t.Fatal("Load() error = nil; want missing facts file error")
	}
}

func TestParseRejectsMalformedFacts(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "non mapping line",
			input: "not a mapping\n",
			want:  "invalid facts line",
		},
		{
			name:  "unknown top level key",
			input: "unexpected: value\n",
			want:  "unknown top-level facts key",
		},
		{
			name: "unsupported indentation",
			input: `current_release:
    version: v1
`,
			want: "unsupported indentation",
		},
		{
			name: "unknown current release key",
			input: `current_release:
  unexpected: value
`,
			want: "unknown current_release facts key",
		},
		{
			name: "unknown runtime key",
			input: `runtime:
  unexpected: value
`,
			want: "unknown runtime facts key",
		},
		{
			name: "unknown tools key",
			input: `tools:
  unexpected: value
`,
			want: "unknown tools facts key",
		},
		{
			name:  "nested key without section",
			input: "  version: v1\n",
			want:  "outside a supported section",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse([]byte(tt.input))
			if err == nil {
				t.Fatalf("Parse() error = nil; want containing %q", tt.want)
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("Parse() error = %q; want containing %q", err.Error(), tt.want)
			}
		})
	}
}

func TestParseReportsScannerError(t *testing.T) {
	_, err := Parse([]byte("schema_version: " + strings.Repeat("a", 70*1024) + "\n"))
	if err == nil {
		t.Fatal("Parse() error = nil; want scanner token error")
	}
	if !strings.Contains(err.Error(), "token too long") {
		t.Fatalf("Parse() error = %q; want scanner token error", err.Error())
	}
}

func TestValidateReportsMissingFieldsAndInvalidTimestamp(t *testing.T) {
	missing := Facts{}
	for _, want := range []string{
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
	} {
		if !hasGap(missing.Validate(), want) {
			t.Fatalf("Validate() gaps missing %q: %v", want, missing.Validate())
		}
	}

	facts := Expected()
	facts.CurrentRelease.ReleasedAt = "not-rfc3339"
	if !hasGap(facts.Validate(), "current_release.released_at must be RFC3339") {
		t.Fatalf("Validate() gaps = %v; want invalid timestamp gap", facts.Validate())
	}
}

func hasGap(gaps []string, want string) bool {
	for _, gap := range gaps {
		if gap == want {
			return true
		}
	}
	return false
}
