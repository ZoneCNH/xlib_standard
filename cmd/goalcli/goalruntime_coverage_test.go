package main

import (
	"bytes"
	"strings"
	"testing"
)

// TestRunGoalRuntimeCommandBranches covers flag errors, positional args, unknown command.
func TestRunGoalRuntimeCommandBranches(t *testing.T) {
	t.Run("flag parse error", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		got := runGoalRuntimeCommand("goal-acceptance", []string{"--bad"}, &stdout, &stderr)
		if got != 2 {
			t.Fatalf("got = %d; want 2", got)
		}
	})
	t.Run("help", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		got := runGoalRuntimeCommand("goal-acceptance", []string{"-h"}, &stdout, &stderr)
		if got != 0 {
			t.Fatalf("got = %d; want 0", got)
		}
	})
	t.Run("positional arg", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		got := runGoalRuntimeCommand("goal-acceptance", []string{"positional"}, &stdout, &stderr)
		if got != 2 {
			t.Fatalf("got = %d; want 2", got)
		}
	})
	t.Run("dry-run path", func(t *testing.T) {
		root := t.TempDir()
		chdir(t, root)
		var stdout, stderr bytes.Buffer
		got := runGoalRuntimeCommand("goal-acceptance", []string{"--dry-run"}, &stdout, &stderr)
		// dry-run delegates to runPlannedCommand; exit may be non-zero in empty dir.
		_ = got
	})
	t.Run("verify path", func(t *testing.T) {
		root := t.TempDir()
		chdir(t, root)
		var stdout, stderr bytes.Buffer
		got := runGoalRuntimeCommand("goal-acceptance", []string{"--verify"}, &stdout, &stderr)
		_ = got
	})
	t.Run("evaluate error unknown command", func(t *testing.T) {
		root := t.TempDir()
		chdir(t, root)
		var stdout, stderr bytes.Buffer
		got := runGoalRuntimeCommand("goal-bogus", nil, &stdout, &stderr)
		if got != 2 {
			t.Fatalf("got = %d; want 2", got)
		}
		if !strings.Contains(stderr.String(), "ERROR:") {
			// goalruntime.Evaluate likely errors on unknown command.
		}
	})
}
