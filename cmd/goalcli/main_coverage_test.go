package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestRunSecretCommandHelp covers help subcommand.
func TestRunSecretCommandHelp(t *testing.T) {
	var stdout, stderr bytes.Buffer
	got := runSecretCommand([]string{"help"}, strings.NewReader(""), &stdout, &stderr)
	if got != 0 {
		t.Fatalf("got = %d; want 0", got)
	}
	if !strings.Contains(stdout.String(), "usage: goalcli secret check") {
		t.Fatalf("stdout = %q", stdout.String())
	}
}

// TestRunSecretCommandNoArgs covers len(args)==0.
func TestRunSecretCommandNoArgs(t *testing.T) {
	var stdout, stderr bytes.Buffer
	got := runSecretCommand(nil, strings.NewReader(""), &stdout, &stderr)
	if got != 2 {
		t.Fatalf("got = %d; want 2", got)
	}
}

// TestVulncheckInterval covers empty env, valid hours, invalid hours.
func TestVulncheckInterval(t *testing.T) {
	_ = os.Unsetenv(vulncheckIntervalHoursEnv)
	got, err := vulncheckInterval()
	if err != nil || got != defaultVulncheckInterval {
		t.Fatalf("default = %v %v; want %v nil", got, err, defaultVulncheckInterval)
	}
	t.Setenv(vulncheckIntervalHoursEnv, "24")
	got, err = vulncheckInterval()
	if err != nil || got != 24*time.Hour {
		t.Fatalf("24h = %v %v", got, err)
	}
	t.Setenv(vulncheckIntervalHoursEnv, "0")
	_, err = vulncheckInterval()
	if err == nil {
		t.Fatalf("0 hours should error")
	}
	t.Setenv(vulncheckIntervalHoursEnv, "notanumber")
	_, err = vulncheckInterval()
	if err == nil {
		t.Fatalf("non-numeric should error")
	}
}

// TestRecordVulncheckRun covers success and directory creation.
func TestRecordVulncheckRun(t *testing.T) {
	root := t.TempDir()
	statePath := filepath.Join(root, "sub", "state")
	if err := recordVulncheckRun(statePath, time.Now().UTC()); err != nil {
		t.Fatalf("err = %v", err)
	}
	if _, err := os.Stat(statePath); err != nil {
		t.Fatalf("state not written: %v", err)
	}
}

// TestFormatDuration covers hour-divisible and non-divisible.
func TestFormatDuration(t *testing.T) {
	if got := formatDuration(7 * 24 * time.Hour); got != "168h" {
		t.Fatalf("168h = %q", got)
	}
	// 90 minutes is not hour-divisible, so uses Duration.String().
	got := formatDuration(90 * time.Minute)
	if got == "" {
		t.Fatalf("90m = %q; want non-empty", got)
	}
}

// TestVulncheckDueError covers the read-state-error branch (non-not-exist).
func TestVulncheckDueError(t *testing.T) {
	root := t.TempDir()
	statePath := filepath.Join(root, "statefile")
	// Write a directory where a file is expected so ReadFile returns a non-ErrNotExist error.
	if err := os.MkdirAll(statePath, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	t.Setenv(vulncheckStateEnv, statePath)
	_ = os.Unsetenv(forceVulncheckEnv)
	_, _, _, err := vulncheckDue(time.Now().UTC())
	if err == nil {
		t.Fatalf("want error reading directory as file")
	}
}

// TestRunSecurityVulncheckDueAndForce covers the vulncheckDue + recordVulncheckRun integration.
func TestRunSecurityVulncheckDueAndForce(t *testing.T) {
	root, _ := setupSecurityFixture(t)
	t.Setenv(enableVulncheckEnv, "1")
	// First run: no state, due=true, govulncheck passes, records run.
	var stdout, stderr bytes.Buffer
	got := runSecurity(strings.NewReader(""), &stdout, &stderr)
	if got != 0 {
		t.Fatalf("first run got = %d; want 0", got)
	}
	// Second run immediately after: within interval, should skip.
	got = runSecurity(strings.NewReader(""), &stdout, &stderr)
	if got != 0 {
		t.Fatalf("second run got = %d; want 0", got)
	}
	if !strings.Contains(stderr.String(), "govulncheck skipped") {
		t.Fatalf("stderr = %q; want skipped", stderr.String())
	}
	_ = root
}

// TestRunSecurityVulncheckRecordError covers recordVulncheckRun failure during security run.
func TestRunSecurityVulncheckRecordError(t *testing.T) {
	root, _ := setupSecurityFixture(t)
	// Make the state path point to a location whose parent cannot be created.
	uncreatableDir := filepath.Join(root, "blocker")
	if err := os.WriteFile(uncreatableDir, []byte("x"), 0o644); err != nil {
		t.Fatalf("write blocker: %v", err)
	}
	t.Setenv(vulncheckStateEnv, filepath.Join(uncreatableDir, "state"))
	t.Setenv(enableVulncheckEnv, "1")
	var stdout, stderr bytes.Buffer
	got := runSecurity(strings.NewReader(""), &stdout, &stderr)
	// Either vulncheckDue errors (2) or recordVulncheckRun errors (1); both are non-zero.
	if got == 0 {
		t.Fatalf("got = %d; want non-zero", got)
	}
}

// TestRunSecurityVulncheckIntervalError covers vulncheckInterval error path in runSecurity.
func TestRunSecurityVulncheckIntervalError(t *testing.T) {
	setupSecurityFixture(t)
	t.Setenv(enableVulncheckEnv, "1")
	t.Setenv(vulncheckIntervalHoursEnv, "bad")
	var stdout, stderr bytes.Buffer
	got := runSecurity(strings.NewReader(""), &stdout, &stderr)
	if got != 2 {
		t.Fatalf("got = %d; want 2", got)
	}
}
