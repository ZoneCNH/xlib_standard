package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestContainsString covers both true/false branches.
func TestContainsString(t *testing.T) {
	values := []string{"a", "b", "~DEFAULT_BRANCH"}
	if !containsString(values, "b") {
		t.Fatalf("containsString should find b")
	}
	if containsString(values, "z") {
		t.Fatalf("containsString should not find z")
	}
	// empty slice
	if containsString(nil, "x") {
		t.Fatalf("containsString(nil) should be false")
	}
}

// TestAppendProtectMainRulesetGapsAllBranches covers every gap branch.
func TestAppendProtectMainRulesetGapsAllBranches(t *testing.T) {
	t.Run("missing file", func(t *testing.T) {
		root := t.TempDir()
		var gaps []string
		appendProtectMainRulesetGaps(root, &gaps)
		if !gapsContainSubstring(gaps, "missing .github/rulesets/protect-main.json") {
			t.Fatalf("gaps = %v; want missing file", gaps)
		}
	})
	t.Run("invalid JSON", func(t *testing.T) {
		root := t.TempDir()
		writeProtectMain(t, root, "{not json")
		var gaps []string
		appendProtectMainRulesetGaps(root, &gaps)
		if !gapsContainSubstring(gaps, "is not valid JSON") {
			t.Fatalf("gaps = %v; want invalid JSON", gaps)
		}
	})
	t.Run("wrong name", func(t *testing.T) {
		root := t.TempDir()
		body := strings.Replace(validProtectMainRulesetFixture(), `"name": "protect-main"`, `"name": "other"`, 1)
		writeProtectMain(t, root, body)
		var gaps []string
		appendProtectMainRulesetGaps(root, &gaps)
		if !gapsContainSubstring(gaps, "missing protect-main name") {
			t.Fatalf("gaps = %v; want missing name", gaps)
		}
	})
	t.Run("wrong target", func(t *testing.T) {
		root := t.TempDir()
		body := strings.Replace(validProtectMainRulesetFixture(), `"target": "branch"`, `"target": "tag"`, 1)
		writeProtectMain(t, root, body)
		var gaps []string
		appendProtectMainRulesetGaps(root, &gaps)
		if !gapsContainSubstring(gaps, "must target branch") {
			t.Fatalf("gaps = %v; want target branch", gaps)
		}
	})
	t.Run("not active", func(t *testing.T) {
		root := t.TempDir()
		body := strings.Replace(validProtectMainRulesetFixture(), `"enforcement": "active"`, `"enforcement": "disabled"`, 1)
		writeProtectMain(t, root, body)
		var gaps []string
		appendProtectMainRulesetGaps(root, &gaps)
		if !gapsContainSubstring(gaps, "enforcement must be active") {
			t.Fatalf("gaps = %v; want enforcement active", gaps)
		}
	})
	t.Run("bypass actors present", func(t *testing.T) {
		root := t.TempDir()
		body := strings.Replace(validProtectMainRulesetFixture(), `"bypass_actors": []`, `"bypass_actors": [{"actor_id": 5}]`, 1)
		writeProtectMain(t, root, body)
		var gaps []string
		appendProtectMainRulesetGaps(root, &gaps)
		if !gapsContainSubstring(gaps, "must not allow bypass actors") {
			t.Fatalf("gaps = %v; want bypass gap", gaps)
		}
	})
	t.Run("missing default branch in include", func(t *testing.T) {
		root := t.TempDir()
		body := strings.Replace(validProtectMainRulesetFixture(), `"~DEFAULT_BRANCH"`, `"refs/heads/main"`, 1)
		writeProtectMain(t, root, body)
		var gaps []string
		appendProtectMainRulesetGaps(root, &gaps)
		if !gapsContainSubstring(gaps, "must protect default branch") {
			t.Fatalf("gaps = %v; want protect default branch", gaps)
		}
	})
	t.Run("missing rule types", func(t *testing.T) {
		root := t.TempDir()
		// Remove pull_request rule block.
		body := validProtectMainRulesetFixture()
		// Replace "pull_request" with "nonexistent" to drop the rule.
		body = strings.Replace(body, `"type": "pull_request"`, `"type": "nonexistent_rule"`, 1)
		body = strings.Replace(body, `"type": "deletion"`, `"type": "other_unknown"`, 1)
		writeProtectMain(t, root, body)
		var gaps []string
		appendProtectMainRulesetGaps(root, &gaps)
		if !gapsContainSubstring(gaps, "missing pull_request rule") {
			t.Fatalf("gaps = %v; want missing pull_request rule", gaps)
		}
		if !gapsContainSubstring(gaps, "missing deletion rule") {
			t.Fatalf("gaps = %v; want missing deletion rule", gaps)
		}
	})
	t.Run("required_status_checks with invalid params", func(t *testing.T) {
		root := t.TempDir()
		body := validProtectMainRulesetFixture()
		// Replace the required_status_checks parameters with invalid JSON.
		body = strings.Replace(body, `"required_status_checks": [`, `"required_status_checks": "not-an-array`, 1)
		writeProtectMain(t, root, body)
		var gaps []string
		appendProtectMainRulesetGaps(root, &gaps)
		// We expect JSON parse failure or missing checks.
		if len(gaps) == 0 {
			t.Fatalf("gaps = %v; want some gap for malformed params", gaps)
		}
	})
	t.Run("missing required status checks", func(t *testing.T) {
		root := t.TempDir()
		body := validProtectMainRulesetFixture()
		// Remove the adoption-check, governance-check, release-check contexts.
		body = strings.Replace(body, `{ "context": "adoption-check" },\n`, "", 1)
		body = strings.Replace(body, `{ "context": "adoption-check" },`, "", 1)
		writeProtectMain(t, root, body)
		var gaps []string
		appendProtectMainRulesetGaps(root, &gaps)
		if !gapsContainSubstring(gaps, "required checks missing adoption-check") {
			t.Fatalf("gaps = %v; want missing adoption-check", gaps)
		}
	})
}

func writeProtectMain(t *testing.T, root, body string) {
	t.Helper()
	path := filepath.Join(root, ".github", "rulesets", "protect-main.json")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
}

// TestRunAdoptionCheckFlagBranches covers flag.Parse error, help, positional.
func TestRunAdoptionCheckFlagBranches(t *testing.T) {
	t.Run("flag parse error", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		got := runAdoptionCheck([]string{"--bad"}, &stdout, &stderr)
		if got != 2 {
			t.Fatalf("got = %d; want 2", got)
		}
	})
	t.Run("help", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		got := runAdoptionCheck([]string{"-h"}, &stdout, &stderr)
		if got != 0 {
			t.Fatalf("got = %d; want 0", got)
		}
	})
	t.Run("positional arg", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		got := runAdoptionCheck([]string{"extra"}, &stdout, &stderr)
		if got != 2 {
			t.Fatalf("got = %d; want 2", got)
		}
	})
}

// TestIsAdoptionSourceRepository covers both branches.
func TestIsAdoptionSourceRepository(t *testing.T) {
	root := t.TempDir()
	// no go.mod
	if isAdoptionSourceRepository(root) {
		t.Fatalf("missing go.mod should be false")
	}
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module github.com/ZoneCNH/kernel\n\ngo 1.23\n"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	if isAdoptionSourceRepository(root) {
		t.Fatalf("kernel should not be source")
	}
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module github.com/ZoneCNH/xlib-standard\n\ngo 1.23\n"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	if !isAdoptionSourceRepository(root) {
		t.Fatalf("xlib-standard should be source")
	}
}

// TestEvaluateAdoptionCheckSourceRepo covers the source-repo path that deletes the lock requirement.
func TestEvaluateAdoptionCheckSourceRepo(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module github.com/ZoneCNH/xlib-standard\n\ngo 1.23\n"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	details, gaps := evaluateAdoptionCheck(root)
	if !slicesContain(details, "source repository governance pack present") {
		t.Fatalf("details = %v; want source repo detail", details)
	}
	// Source repo deletes lockPath so missing lock should not be a gap.
	if gapsContainSubstring(gaps, "xlib-standard.lock") {
		t.Fatalf("source repo should not require lock; gaps = %v", gaps)
	}
}
