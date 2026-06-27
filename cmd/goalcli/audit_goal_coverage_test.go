package main

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

// TestRunAuditGoalFlagBranches covers flag.Parse error, positional args.
func TestRunAuditGoalFlagBranches(t *testing.T) {
	t.Run("flag parse error", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		got := runAuditGoal([]string{"--bad"}, &stdout, &stderr)
		if got != 2 {
			t.Fatalf("got = %d; want 2", got)
		}
	})
	t.Run("help", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		got := runAuditGoal([]string{"-h"}, &stdout, &stderr)
		if got != 0 {
			t.Fatalf("got = %d; want 0", got)
		}
	})
	t.Run("positional arg", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		got := runAuditGoal([]string{"extra"}, &stdout, &stderr)
		if got != 2 {
			t.Fatalf("got = %d; want 2", got)
		}
	})
	t.Run("goal-id annotation", func(t *testing.T) {
		root := t.TempDir()
		chdir(t, root)
		var stdout, stderr bytes.Buffer
		got := runAuditGoal([]string{"--goal-id", "GOAL-X"}, &stdout, &stderr)
		// audit-goal runs checks that will fail in empty dir, but goal_id should appear.
		if !strings.Contains(stdout.String(), "goal_id=GOAL-X") {
			t.Fatalf("stdout = %q; want goal_id=GOAL-X", stdout.String())
		}
		_ = got
	})
}

// TestAuditGoalDefaultChecks verifies the returned check list contains the dry-run checks.
func TestAuditGoalDefaultChecks(t *testing.T) {
	checks := auditGoalDefaultChecks(traceabilityMatrixPath)
	names := make(map[string]bool, len(checks))
	for _, c := range checks {
		names[c.name] = true
	}
	for _, expected := range []string{
		"context-check", "spec-check", "design-check", "task-check",
		"evidence-check", "cli-contract", "issue-registry",
		"command-registry", "makefile-baseline", "traceability-check",
		"goal-acceptance:dry-run", "goal-delivery:dry-run",
		"goal-handover:dry-run", "goal-downstream-adoption:dry-run",
		"goal-certify:dry-run", "goal-runtime-final:dry-run",
	} {
		if !names[expected] {
			t.Errorf("auditGoalDefaultChecks missing %q", expected)
		}
	}
}

// TestAuditGoalRuntimeDryRunCheck constructs a dry-run check and exercises it.
func TestAuditGoalRuntimeDryRunCheck(t *testing.T) {
	root := t.TempDir()
	chdir(t, root)
	check := auditGoalRuntimeDryRunCheck("goal-acceptance")
	if check.name != "goal-acceptance:dry-run" {
		t.Fatalf("name = %q", check.name)
	}
	var stdout, stderr bytes.Buffer
	code := check.run(&stdout, &stderr)
	// In an empty temp dir, planned command will report failure/gap; code should be non-zero.
	_ = code
	// Either way the run() should have executed without panic.
}

// TestAuditGoalComponentSummary covers both stdout and stderr fallback plus truncation.
func TestAuditGoalComponentSummary(t *testing.T) {
	if got := auditGoalComponentSummary("  hello   world  ", ""); got != "hello world" {
		t.Fatalf("stdout summary = %q; want 'hello world'", got)
	}
	if got := auditGoalComponentSummary("", "err\nmsg"); got != "err msg" {
		t.Fatalf("stderr fallback = %q; want 'err msg'", got)
	}
	long := strings.Repeat("a", 400)
	got := auditGoalComponentSummary(long, "")
	if !strings.HasSuffix(got, "...") || len(got) != 303 {
		t.Fatalf("truncated len = %d; want 303", len(got))
	}
}

// TestNewAuditGoalChecksVariableOverride verifies the package-level var can be swapped.
func TestNewAuditGoalChecksVariableOverride(t *testing.T) {
	original := newAuditGoalChecks
	t.Cleanup(func() { newAuditGoalChecks = original })
	called := false
	newAuditGoalChecks = func(matrixPath string) []auditGoalCheck {
		called = true
		return []auditGoalCheck{auditGoalCheck{name: "custom", run: func(w1, w2 io.Writer) int { return 0 }}}
	}
	root := t.TempDir()
	chdir(t, root)
	var stdout, stderr bytes.Buffer
	_ = runAuditGoal(nil, &stdout, &stderr)
	if !called {
		t.Fatalf("custom newAuditGoalChecks was not invoked")
	}
}
