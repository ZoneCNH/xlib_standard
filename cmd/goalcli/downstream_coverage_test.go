package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestRunDownstreamSyncPlanBranches covers flag errors, positional, format, output validation.
func TestRunDownstreamSyncPlanBranches(t *testing.T) {
	t.Run("flag parse error", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		got := runDownstreamSyncPlan([]string{"--bad"}, &stdout, &stderr)
		if got != 2 {
			t.Fatalf("got = %d; want 2", got)
		}
	})
	t.Run("help", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		got := runDownstreamSyncPlan([]string{"-h"}, &stdout, &stderr)
		if got != 0 {
			t.Fatalf("got = %d; want 0", got)
		}
	})
	t.Run("positional arg", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		got := runDownstreamSyncPlan([]string{"positional"}, &stdout, &stderr)
		if got != 2 {
			t.Fatalf("got = %d; want 2", got)
		}
	})
	t.Run("invalid format", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		got := runDownstreamSyncPlan([]string{"--format", "xml"}, &stdout, &stderr)
		if got != 2 {
			t.Fatalf("got = %d; want 2", got)
		}
	})
	t.Run("invalid output path", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		got := runDownstreamSyncPlan([]string{"--output", "/abs/path"}, &stdout, &stderr)
		if got != 1 {
			t.Fatalf("got = %d; want 1", got)
		}
	})
	t.Run("missing impact report", func(t *testing.T) {
		root := t.TempDir()
		chdir(t, root)
		var stdout, stderr bytes.Buffer
		got := runDownstreamSyncPlan([]string{"--impact-report", "missing.md", "--output", "-"}, &stdout, &stderr)
		if got != 1 {
			t.Fatalf("got = %d; want 1", got)
		}
		if !strings.Contains(stderr.String(), "standard-impact-check") {
			t.Fatalf("stderr = %q", stderr.String())
		}
	})
	t.Run("valid markdown to stdout", func(t *testing.T) {
		root := t.TempDir()
		chdir(t, root)
		writeDownstreamImpactReport(t, root, false)
		var stdout, stderr bytes.Buffer
		got := runDownstreamSyncPlan([]string{"--impact-report", "impact.md", "--output", "-", "--format", "markdown"}, &stdout, &stderr)
		if got != 0 {
			t.Fatalf("got = %d; want 0; stderr=%q", got, stderr.String())
		}
		if !strings.Contains(stdout.String(), "# Downstream Sync Plan") {
			t.Fatalf("stdout = %q; want plan header", stdout.String())
		}
	})
	t.Run("valid json to file", func(t *testing.T) {
		root := t.TempDir()
		chdir(t, root)
		writeDownstreamImpactReport(t, root, true)
		var stdout, stderr bytes.Buffer
		got := runDownstreamSyncPlan([]string{"--impact-report", "impact.md", "--output", "plan.json", "--format", "json"}, &stdout, &stderr)
		if got != 0 {
			t.Fatalf("got = %d; want 0; stderr=%q", got, stderr.String())
		}
		if _, err := os.Stat(filepath.Join(root, "plan.json")); err != nil {
			t.Fatalf("plan.json not written: %v", err)
		}
	})
}

// TestValidateDownstreamSyncPlanOutputPath covers all branches.
func TestValidateDownstreamSyncPlanOutputPath(t *testing.T) {
	root := t.TempDir()
	chdir(t, root)
	cases := []struct {
		name      string
		output    string
		workspace string
		wantErr   bool
		errSubstr string
	}{
		{"dash", "-", "..", false, ""},
		{"empty", "  ", "..", true, "empty"},
		{"absolute", "/x/y", "..", true, "repository-relative"},
		{"parent traversal", "../escape", "..", true, "parent traversal"},
		{"dot only", ".", "..", true, "file"},
		{"protected adoption status", ".agent/registries/downstream-adoption-status.yaml", "..", true, "adoption truth"},
		{"protected truth state", ".agent/evidence/truth-state.yaml", "..", true, "adoption truth"},
		{"within workspace", "ws/plan.md", "ws", true, "downstream workspace"},
		{"valid file", "release/plan.md", "..", false, ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := validateDownstreamSyncPlanOutputPath(c.output, c.workspace)
			if c.wantErr {
				if err == nil {
					t.Fatalf("want error containing %q; got nil", c.errSubstr)
				}
				if c.errSubstr != "" && !strings.Contains(err.Error(), c.errSubstr) {
					t.Fatalf("err = %q; want containing %q", err.Error(), c.errSubstr)
				}
			} else {
				if err != nil {
					t.Fatalf("err = %v; want nil", err)
				}
			}
		})
	}
	// Existing directory error.
	dirPath := filepath.Join(root, "existingdir")
	if err := os.MkdirAll(dirPath, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := validateDownstreamSyncPlanOutputPath("existingdir", ".."); err == nil {
		t.Fatalf("directory output should error")
	}
}

// TestIsWithinDownstreamWorkspace covers all branches.
func TestIsWithinDownstreamWorkspace(t *testing.T) {
	cases := []struct {
		output, workspace string
		want              bool
	}{
		{"release/a.md", "", false},
		{"release/a.md", ".", false},
		{"release/a.md", "/abs", false},
		{"release/a.md", "../parent", false}, // parent traversal rejected
		{"ws/a.md", "ws", true},
		{"ws/sub/a.md", "ws", true},
		{"other/a.md", "ws", false},
		{"ws/a.md", "ws/", true}, // trailing slash cleaned
	}
	for _, c := range cases {
		if got := isWithinDownstreamWorkspace(c.output, c.workspace); got != c.want {
			t.Errorf("isWithinDownstreamWorkspace(%q,%q) = %v; want %v", c.output, c.workspace, got, c.want)
		}
	}
}

// TestParseDownstreamImpactReport covers valid + all error branches.
func TestParseDownstreamImpactReport(t *testing.T) {
	t.Run("missing file", func(t *testing.T) {
		_, err := parseDownstreamImpactReport("nonexistent.md")
		if err == nil {
			t.Fatalf("want error")
		}
	})
	t.Run("missing downstream_sync_required", func(t *testing.T) {
		root := t.TempDir()
		path := filepath.Join(root, "bad.md")
		if err := os.WriteFile(path, []byte("- primary_downstream: `k`\n- changed_file_count: `1`\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		_, err := parseDownstreamImpactReport(path)
		if err == nil || !strings.Contains(err.Error(), "missing downstream_sync_required") {
			t.Fatalf("err = %v", err)
		}
	})
	t.Run("invalid bool", func(t *testing.T) {
		root := t.TempDir()
		path := filepath.Join(root, "bad.md")
		if err := os.WriteFile(path, []byte("- downstream_sync_required: `maybe`\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		_, err := parseDownstreamImpactReport(path)
		if err == nil || !strings.Contains(err.Error(), "invalid downstream_sync_required") {
			t.Fatalf("err = %v", err)
		}
	})
	t.Run("invalid release decision", func(t *testing.T) {
		root := t.TempDir()
		path := filepath.Join(root, "bad.md")
		if err := os.WriteFile(path, []byte("- downstream_sync_required: `true`\n- downstream_release_decision: `bogus`\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		_, err := parseDownstreamImpactReport(path)
		if err == nil || !strings.Contains(err.Error(), "invalid downstream_release_decision") {
			t.Fatalf("err = %v", err)
		}
	})
	t.Run("sync true but decision not required", func(t *testing.T) {
		root := t.TempDir()
		path := filepath.Join(root, "bad.md")
		if err := os.WriteFile(path, []byte("- downstream_sync_required: `true`\n- downstream_release_decision: `not_required`\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		_, err := parseDownstreamImpactReport(path)
		if err == nil || !strings.Contains(err.Error(), "requires downstream_release_decision=required") {
			t.Fatalf("err = %v", err)
		}
	})
	t.Run("sync false but decision required", func(t *testing.T) {
		root := t.TempDir()
		path := filepath.Join(root, "bad.md")
		if err := os.WriteFile(path, []byte("- downstream_sync_required: `false`\n- downstream_release_decision: `required`\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		_, err := parseDownstreamImpactReport(path)
		if err == nil || !strings.Contains(err.Error(), "requires downstream_release_decision=not_required") {
			t.Fatalf("err = %v", err)
		}
	})
	t.Run("invalid repo rules decision", func(t *testing.T) {
		root := t.TempDir()
		path := filepath.Join(root, "bad.md")
		if err := os.WriteFile(path, []byte("- downstream_sync_required: `false`\n- downstream_release_decision: `not_required`\n- repository_rules_release_decision: `bogus`\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		_, err := parseDownstreamImpactReport(path)
		if err == nil || !strings.Contains(err.Error(), "invalid repository_rules_release_decision") {
			t.Fatalf("err = %v", err)
		}
	})
	t.Run("missing primary_downstream", func(t *testing.T) {
		root := t.TempDir()
		path := filepath.Join(root, "bad.md")
		if err := os.WriteFile(path, []byte("- downstream_sync_required: `false`\n- downstream_release_decision: `not_required`\n- repository_rules_release_decision: `not_required`\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		_, err := parseDownstreamImpactReport(path)
		if err == nil || !strings.Contains(err.Error(), "missing primary_downstream") {
			t.Fatalf("err = %v", err)
		}
	})
	t.Run("missing changed_file_count", func(t *testing.T) {
		root := t.TempDir()
		path := filepath.Join(root, "bad.md")
		if err := os.WriteFile(path, []byte("- downstream_sync_required: `false`\n- downstream_release_decision: `not_required`\n- repository_rules_release_decision: `not_required`\n- primary_downstream: `k`\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		_, err := parseDownstreamImpactReport(path)
		if err == nil || !strings.Contains(err.Error(), "missing changed_file_count") {
			t.Fatalf("err = %v", err)
		}
	})
	t.Run("invalid changed_file_count", func(t *testing.T) {
		root := t.TempDir()
		path := filepath.Join(root, "bad.md")
		if err := os.WriteFile(path, []byte("- downstream_sync_required: `false`\n- downstream_release_decision: `not_required`\n- repository_rules_release_decision: `not_required`\n- primary_downstream: `k`\n- changed_file_count: `abc`\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		_, err := parseDownstreamImpactReport(path)
		if err == nil || !strings.Contains(err.Error(), "invalid changed_file_count") {
			t.Fatalf("err = %v", err)
		}
	})
}

// TestShellQuote covers all branches.
func TestShellQuote(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"", "''"},
		{"plain", "plain"},
		{"a_b-c.d/e:f", "a_b-c.d/e:f"},
		{"has space", "'has space'"},
		{"it's", "'it'\"'\"'s'"},
	}
	for _, c := range cases {
		if got := shellQuote(c.in); got != c.want {
			t.Errorf("shellQuote(%q) = %q; want %q", c.in, got, c.want)
		}
	}
}

// TestIsDownstreamImpactCategory covers known + unknown.
func TestIsDownstreamImpactCategory(t *testing.T) {
	if !isDownstreamImpactCategory("contracts") {
		t.Fatalf("contracts should be known")
	}
	if isDownstreamImpactCategory("bogus") {
		t.Fatalf("bogus should be unknown")
	}
}

// TestSortedDownstreamImpactCategories covers known + extras ordering.
func TestSortedDownstreamImpactCategories(t *testing.T) {
	counts := map[string]int{
		"docs":      1,
		"contracts": 2,
		"zebra":     3, // extra - sorts after known
	}
	got := sortedDownstreamImpactCategories(counts)
	// Known categories come in downstreamImpactCategories order first.
	if got[0] != "contracts" || got[1] != "docs" {
		t.Fatalf("known order = %v; want contracts then docs", got)
	}
	if got[2] != "zebra" {
		t.Fatalf("extra = %v; want zebra last", got)
	}
}

// writeDownstreamImpactReport writes a valid impact report fixture.
func writeDownstreamImpactReport(t *testing.T, root string, syncRequired bool) {
	t.Helper()
	required := "false"
	decision := "not_required"
	if syncRequired {
		required = "true"
		decision = "required"
	}
	body := "- downstream_sync_required: `" + required + "`\n" +
		"- downstream_release_decision: `" + decision + "`\n" +
		"- repository_rules_release_decision: `not_required`\n" +
		"- primary_downstream: `kernel`\n" +
		"- changed_file_count: `5`\n\n" +
		"## contracts\n\n- `a.go`\n"
	path := filepath.Join(root, "impact.md")
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
}
