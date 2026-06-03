package goalruntime

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEvaluateGoalRuntimeFinalPassesWithAuthority(t *testing.T) {
	root := t.TempDir()
	writeAuthorityFixture(t, root)

	report, err := Evaluate("goal-runtime-final", Options{Root: root})
	if err != nil {
		t.Fatalf("Evaluate returned error: %v", err)
	}
	if report.Status != "passed" {
		t.Fatalf("status = %q; gaps %#v", report.Status, report.Gaps)
	}
	if report.GoalID != DefaultGoalID {
		t.Fatalf("goal_id = %q; want %q", report.GoalID, DefaultGoalID)
	}
	if report.Gate != "G16" {
		t.Fatalf("gate = %q; want G16", report.Gate)
	}
	if !contains(report.Evidence, "evidence_ledger="+EvidenceLedgerPath) {
		t.Fatalf("evidence = %#v; want ledger path", report.Evidence)
	}
	if !contains(report.Evidence, "requires=goal-certify") {
		t.Fatalf("evidence = %#v; want final gate dependency", report.Evidence)
	}
	if !containsSubstring(report.Details, "不是全局 release blocking gates") {
		t.Fatalf("details = %#v; want false-completion boundary", report.Details)
	}
	if len(report.AuthorityPaths) == 0 {
		t.Fatalf("authority_paths is empty")
	}
}

func TestEvaluateRejectsUnknownCommand(t *testing.T) {
	if _, err := Evaluate("not-a-goalkit-command", Options{}); err == nil {
		t.Fatalf("Evaluate returned nil error for unknown command")
	}
}

func TestEvaluateReportsMissingAuthorityPaths(t *testing.T) {
	root := t.TempDir()
	report, err := Evaluate("goal-acceptance", Options{Root: root})
	if err != nil {
		t.Fatalf("Evaluate returned error: %v", err)
	}
	if report.Status != "failed" {
		t.Fatalf("status = %q; want failed", report.Status)
	}
	if len(report.Gaps) == 0 {
		t.Fatalf("gaps is empty; want missing authority paths")
	}
	if !containsSubstring(report.Gaps, ".worktree/goalkit-v0.1.0-plan.md") {
		t.Fatalf("gaps = %#v; want root plan gap", report.Gaps)
	}
}

func writeAuthorityFixture(t *testing.T, root string) {
	t.Helper()
	for _, path := range []string{
		".worktree/goalkit-v0.1.0-plan.md",
		".omx/context/goalkit-v0.1.0-team-20260603T005302Z.md",
		"docs/standard/xlibgate-cli-contract.md",
		".agent/harness.yaml",
		".agent/command-registry.yaml",
		"Makefile",
	} {
		full := filepath.Join(root, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatalf("mkdir fixture path %s: %v", path, err)
		}
		if err := os.WriteFile(full, []byte("fixture\n"), 0o644); err != nil {
			t.Fatalf("write fixture path %s: %v", path, err)
		}
	}
}

func contains(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func containsSubstring(values []string, want string) bool {
	for _, value := range values {
		if strings.Contains(value, want) {
			return true
		}
	}
	return false
}
