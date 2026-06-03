package goalruntime

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestRunGoalRuntimeFinalIncludesAllNonBlockingGates(t *testing.T) {
	var stdout, stderr bytes.Buffer
	got := Run("goal-runtime-final", []string{"--goal-id", "GOAL-20260603-XLIB-RUNTIME-001", "--mode", "FULL", "--json"}, &stdout, &stderr)
	if got != 0 {
		t.Fatalf("Run exit = %d, stderr %q, stdout %q; want 0", got, stderr.String(), stdout.String())
	}
	var report Report
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("stdout is not Report JSON: %v; stdout %q", err, stdout.String())
	}
	if report.Command != "goal-runtime-final" || report.Status != "passed" || report.Blocking || report.MVAStatus != "not-complete" {
		t.Fatalf("report = %#v; want passed, non-blocking, not-complete final report", report)
	}
	if report.Executor != Executor || report.ControlPlane != ControlPlane || report.LedgerPath != LedgerPath {
		t.Fatalf("report authority = %#v; want xlibgate/Harness/%s", report, LedgerPath)
	}
	if len(report.Gates) != 5 {
		t.Fatalf("gates = %#v; want 5 G12-G16 gates", report.Gates)
	}
}

func TestRunGoalRuntimeRequiresValidGoalID(t *testing.T) {
	var stdout, stderr bytes.Buffer
	got := Run("goal-acceptance", []string{"--goal-id", "bad", "--json"}, &stdout, &stderr)
	if got != 2 {
		t.Fatalf("Run invalid goal id exit = %d, stdout %q, stderr %q; want 2", got, stdout.String(), stderr.String())
	}
	if !strings.Contains(stderr.String(), "invalid or missing --goal-id/GOAL_ID") {
		t.Fatalf("stderr = %q; want invalid goal id error", stderr.String())
	}
}
