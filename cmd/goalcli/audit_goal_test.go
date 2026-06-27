package main

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func overrideAuditGoalChecks(t *testing.T, checks func(string) []auditGoalCheck) {
	t.Helper()
	old := newAuditGoalChecks
	t.Cleanup(func() { newAuditGoalChecks = old })
	newAuditGoalChecks = checks
}

func TestAuditGoalComponentSummaryBranches(t *testing.T) {
	if got := auditGoalComponentSummary("  line one\n\tline   two  ", "fallback"); got != "line one line two" {
		t.Fatalf("summary from stdout = %q; want collapsed stdout", got)
	}
	if got := auditGoalComponentSummary("", "  stderr\n output  "); got != "stderr output" {
		t.Fatalf("summary from stderr = %q; want collapsed stderr", got)
	}

	long := strings.Repeat("x", 301)
	got := auditGoalComponentSummary(long, "")
	if len(got) != 303 || !strings.HasSuffix(got, "...") {
		t.Fatalf("truncated summary length = %d value suffix %q; want 300 chars plus ellipsis", len(got), got[len(got)-3:])
	}
}

func TestRunAuditGoalReportsComponentFailures(t *testing.T) {
	overrideAuditGoalChecks(t, func(matrixPath string) []auditGoalCheck {
		if matrixPath != "custom-matrix.md" {
			t.Fatalf("matrixPath = %q; want custom-matrix.md", matrixPath)
		}
		return []auditGoalCheck{
			{name: "ok", run: func(stdout, stderr io.Writer) int { return 0 }},
			{name: "stderr-only", run: func(stdout, stderr io.Writer) int {
				_, _ = stderr.Write([]byte("component failed\nwith detail"))
				return 5
			}},
			{name: "silent", run: func(stdout, stderr io.Writer) int { return 6 }},
		}
	})

	var stdout, stderr bytes.Buffer
	code := runAuditGoal([]string{"--goal-id", "GOAL-1", "--matrix", "custom-matrix.md"}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("runAuditGoal() code = %d stderr = %q stdout = %q; want 1", code, stderr.String(), stdout.String())
	}
	for _, want := range []string{
		`"command": "audit-goal"`,
		`"status": "failed"`,
		"goal_id=GOAL-1",
		"ok: passed",
		"stderr-only: exit code 5: component failed with detail",
		"silent: exit code 6: no component output",
	} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("stdout = %q; want %q", stdout.String(), want)
		}
	}
	if !strings.Contains(stderr.String(), "audit-goal found 2 gap(s)") {
		t.Fatalf("stderr = %q; want gap count", stderr.String())
	}
}

func TestRunAuditGoalArgumentBranches(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantCode int
		wantErr  string
	}{
		{name: "help", args: []string{"--help"}, wantCode: 0},
		{name: "parse error", args: []string{"--missing"}, wantCode: 2, wantErr: "invalid arguments"},
		{name: "positional", args: []string{"unexpected"}, wantCode: 2, wantErr: "accepts no positional arguments"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := runAuditGoal(tt.args, &stdout, &stderr)
			if code != tt.wantCode {
				t.Fatalf("runAuditGoal(%v) code = %d stderr = %q; want %d", tt.args, code, stderr.String(), tt.wantCode)
			}
			if tt.wantErr != "" && !strings.Contains(stderr.String(), tt.wantErr) {
				t.Fatalf("stderr = %q; want %q", stderr.String(), tt.wantErr)
			}
		})
	}
}

func TestAuditGoalDefaultChecksConstructsComponentList(t *testing.T) {
	checks := auditGoalDefaultChecks("custom-traceability.md")
	got := make([]string, 0, len(checks))
	for _, check := range checks {
		got = append(got, check.name)
	}
	want := []string{
		"context-check",
		"spec-check",
		"design-check",
		"task-check",
		"evidence-check",
		"cli-contract",
		"issue-registry",
		"command-registry",
		"makefile-baseline",
		"traceability-check",
		"goal-acceptance:dry-run",
		"goal-delivery:dry-run",
		"goal-handover:dry-run",
		"goal-downstream-adoption:dry-run",
		"goal-certify:dry-run",
		"goal-runtime-final:dry-run",
	}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Fatalf("auditGoalDefaultChecks names = %#v; want %#v", got, want)
	}

	var stdout, stderr bytes.Buffer
	code := checks[9].run(&stdout, &stderr)
	if code != 2 {
		t.Fatalf("traceability check code = %d stdout = %q stderr = %q; want missing matrix error", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stderr.String(), "custom-traceability.md") {
		t.Fatalf("traceability stderr = %q; want custom matrix path", stderr.String())
	}
}

func TestAuditGoalDefaultGovernanceChecksRun(t *testing.T) {
	root := t.TempDir()
	chdir(t, root)
	writeDebtCLIFile(t, root, "docs/goal/goal.md", "REQ-1: goal context\n")
	writeDebtCLIFile(t, root, "docs/requirements.md", "REQ-2: requirement\n")
	writeDebtCLIFile(t, root, ".agent/registries/commands.yaml", "commands:\n  - name: legacy\n")

	for _, check := range auditGoalDefaultChecks("missing-traceability.md")[:9] {
		t.Run(check.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := check.run(&stdout, &stderr)
			if code != 0 && code != 1 {
				t.Fatalf("%s code = %d stdout = %q stderr = %q; want report exit", check.name, code, stdout.String(), stderr.String())
			}
			if !strings.Contains(stdout.String(), `"command": "`+check.name+`"`) {
				t.Fatalf("%s stdout = %q stderr = %q; want structured command report", check.name, stdout.String(), stderr.String())
			}
		})
	}
}

func TestAuditGoalRuntimeDryRunCheckRunsVerifier(t *testing.T) {
	root := t.TempDir()
	chdir(t, root)
	writeDebtCLIFile(t, root, ".agent/harness/harness.yaml", `schema_version: "2.9.3"
goalcli_mva_gates:
  G12_ACCEPTANCE:
    command: goal-acceptance
    planned_command: goal-acceptance --dry-run --verify
`)

	check := auditGoalRuntimeDryRunCheck("goal-acceptance")
	if check.name != "goal-acceptance:dry-run" {
		t.Fatalf("runtime dry-run check name = %q; want goal-acceptance:dry-run", check.name)
	}

	var stdout, stderr bytes.Buffer
	code := check.run(&stdout, &stderr)
	if code != 0 {
		t.Fatalf("runtime dry-run code = %d stdout = %q stderr = %q; want 0", code, stdout.String(), stderr.String())
	}
	for _, want := range []string{
		`"command": "goal-acceptance"`,
		`"status": "passed"`,
		"local dry-run verifier satisfied manifest coverage",
	} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("stdout = %q; want %q", stdout.String(), want)
		}
	}
}
