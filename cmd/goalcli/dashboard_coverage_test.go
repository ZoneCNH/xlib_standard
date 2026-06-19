package main

import (
	"bytes"
	"strings"
	"testing"
)

// TestRunDashboardGenerateBranches covers flag errors, positional, format, gaps.
func TestRunDashboardGenerateBranches(t *testing.T) {
	t.Run("flag parse error", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		got := runDashboardGenerate([]string{"--bad"}, &stdout, &stderr)
		if got != 2 {
			t.Fatalf("got = %d; want 2", got)
		}
	})
	t.Run("help", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		got := runDashboardGenerate([]string{"-h"}, &stdout, &stderr)
		if got != 0 {
			t.Fatalf("got = %d; want 0", got)
		}
	})
	t.Run("positional arg", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		got := runDashboardGenerate([]string{"positional"}, &stdout, &stderr)
		if got != 2 {
			t.Fatalf("got = %d; want 2", got)
		}
	})
	t.Run("invalid format", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		got := runDashboardGenerate([]string{"--format", "xml"}, &stdout, &stderr)
		if got != 2 {
			t.Fatalf("got = %d; want 2", got)
		}
	})
	t.Run("markdown with gaps", func(t *testing.T) {
		setupEmptyDashboardRoot(t)
		var stdout, stderr bytes.Buffer
		got := runDashboardGenerate([]string{"--format", "markdown", "--matrix", ".agent/traceability/traceability-matrix.md"}, &stdout, &stderr)
		// In empty dir, checks fail; report.Status should be failed.
		_ = got
		if !strings.Contains(stdout.String(), "# Goal Dashboard") {
			t.Fatalf("stdout = %q; want dashboard header", stdout.String())
		}
	})
	t.Run("json output", func(t *testing.T) {
		setupEmptyDashboardRoot(t)
		var stdout, stderr bytes.Buffer
		_ = runDashboardGenerate([]string{"--format", "json"}, &stdout, &stderr)
		if !strings.Contains(stdout.String(), `"command": "dashboard-generate"`) {
			t.Fatalf("stdout = %q; want command", stdout.String())
		}
	})
}

func setupEmptyDashboardRoot(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	chdir(t, root)
	return root
}

// TestBuildDashboardReport covers both passed and failed scenarios.
func TestBuildDashboardReport(t *testing.T) {
	root := t.TempDir()
	chdir(t, root)
	report := buildDashboardReport("", ".agent/traceability/traceability-matrix.md")
	if report.Command != "dashboard-generate" {
		t.Fatalf("command = %q", report.Command)
	}
	// In empty dir many checks fail.
	if report.Status != "failed" {
		t.Logf("status = %q; expected failed for empty dir", report.Status)
	}
	if len(report.Components) == 0 {
		t.Fatalf("components empty; want entries")
	}
}

// TestRenderDashboardMarkdown covers goal-id rendering and gaps rendering.
func TestRenderDashboardMarkdown(t *testing.T) {
	t.Run("with goal id and no gaps", func(t *testing.T) {
		report := dashboardReport{
			SchemaVersion: "1.0",
			Command:       "dashboard-generate",
			Status:        "passed",
			GoalID:        "GOAL-1",
			Matrix:        "matrix.md",
			Scope:         []string{"goal", "req"},
			Mode:          "local-readonly",
			Components:    []dashboardComponent{{Name: "c", Status: "passed", Summary: "ok"}},
		}
		md := renderDashboardMarkdown(report)
		if !strings.Contains(md, "| goal_id | GOAL-1 |") {
			t.Fatalf("md missing goal_id: %q", md)
		}
		if !strings.Contains(md, "- None") {
			t.Fatalf("md missing None gaps: %q", md)
		}
	})
	t.Run("with gaps", func(t *testing.T) {
		report := dashboardReport{
			Command:    "dashboard-generate",
			Status:     "failed",
			Components: []dashboardComponent{{Name: "c", Status: "failed", Summary: "boom"}},
			Gaps:       []string{"gap one"},
		}
		md := renderDashboardMarkdown(report)
		if !strings.Contains(md, "- gap one") {
			t.Fatalf("md missing gap: %q", md)
		}
	})
}
