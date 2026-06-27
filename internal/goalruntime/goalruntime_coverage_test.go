// SPDX-License-Identifier: Apache-2.0
package goalruntime

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCommandsReturnsSupportedGoalcliMVACommands(t *testing.T) {
	got := Commands()
	want := []string{
		"goal-acceptance",
		"goal-delivery",
		"goal-handover",
		"goal-downstream-adoption",
		"goal-certify",
		"goal-runtime-final",
	}
	if len(got) != len(want) {
		t.Fatalf("Commands() = %v; want %v", got, want)
	}
	for i, command := range got {
		if command != want[i] {
			t.Fatalf("Commands()[%d] = %q; want %q", i, command, want[i])
		}
	}
}

// TestEvaluateDefaultsEmptyRootToCwd covers the `if root == ""` branch by
// leaving Options.Root empty and chdir-ing into a fixture root.
func TestEvaluateDefaultsEmptyRootToCwd(t *testing.T) {
	root := t.TempDir()
	writeAuthorityFixture(t, root)
	previous, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(previous) })

	report, err := Evaluate("goal-acceptance", Options{})
	if err != nil {
		t.Fatalf("Evaluate returned error: %v", err)
	}
	if report.Status != "passed" {
		t.Fatalf("status = %q; gaps %#v", report.Status, report.Gaps)
	}
}

// TestEvaluateAppendsNonDefaultGoalIDDetail covers the goalID != DefaultGoalID
// detail-append branch.
func TestEvaluateAppendsNonDefaultGoalIDDetail(t *testing.T) {
	root := t.TempDir()
	writeAuthorityFixture(t, root)

	report, err := Evaluate("goal-acceptance", Options{Root: root, GoalID: "GOAL-CUSTOM-001"})
	if err != nil {
		t.Fatalf("Evaluate returned error: %v", err)
	}
	if report.GoalID != "GOAL-CUSTOM-001" {
		t.Fatalf("goal_id = %q; want GOAL-CUSTOM-001", report.GoalID)
	}
	if !containsSubstring(report.Details, "non-default goal_id accepted for local contract replay") {
		t.Fatalf("details = %#v; want non-default goal_id detail", report.Details)
	}
}

// TestModulePathForRootHandlesMissingModuleLine covers the branch where go.mod
// is readable but contains no `module` directive.
func TestModulePathForRootHandlesMissingModuleLine(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("go 1.23\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	path, ok := modulePathForRoot(root)
	if ok || path != "" {
		t.Fatalf("modulePathForRoot = %q, %t; want \"\", false", path, ok)
	}
}

// TestModulePathForRootReportsMissingGoMod covers the os.ReadFile error branch.
func TestModulePathForRootReportsMissingGoMod(t *testing.T) {
	root := t.TempDir()
	path, ok := modulePathForRoot(root)
	if ok || path != "" {
		t.Fatalf("modulePathForRoot = %q, %t; want \"\", false", path, ok)
	}
}

// TestWriteEvidenceRejectsUnsupportedCommandCoverage covers the unsupported-command
// rejection branch.
func TestWriteEvidenceRejectsUnsupportedCommandCoverage(t *testing.T) {
	root := t.TempDir()
	if err := WriteEvidence(root, Report{Command: "not-a-goalcli-command"}); err == nil {
		t.Fatal("WriteEvidence returned nil error for unsupported command")
	}
}

// TestWriteEvidenceDefaultsEmptyRootToCwd covers the `if root == ""` branch in
// WriteEvidence for a non-final command.
func TestWriteEvidenceDefaultsEmptyRootToCwd(t *testing.T) {
	root := t.TempDir()
	writeAuthorityFixture(t, root)
	report, err := Evaluate("goal-acceptance", Options{Root: root})
	if err != nil {
		t.Fatalf("Evaluate returned error: %v", err)
	}
	// Pre-seed ledger into cwd-relative location so the WriteEvidence with empty
	// root appends idempotently rather than creating at ".".
	previous, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(previous) })

	if err := WriteEvidence("", report); err != nil {
		t.Fatalf("WriteEvidence with empty root returned error: %v", err)
	}
	if _, err := os.Stat(filepath.FromSlash(SourceLedgerPath)); err != nil {
		t.Fatalf("expected ledger at cwd-relative path: %v", err)
	}
}

// TestWriteEvidenceFinalRejectsMissingPrerequisites covers the
// validateFinalPrerequisites gap branch in WriteEvidence for the final command.
func TestWriteEvidenceFinalRejectsMissingPrerequisites(t *testing.T) {
	root := t.TempDir()
	writeAuthorityFixture(t, root)
	report, err := Evaluate("goal-runtime-final", Options{Root: root})
	if err != nil {
		t.Fatalf("Evaluate returned error: %v", err)
	}
	if err := WriteEvidence(root, report); err == nil {
		t.Fatal("WriteEvidence returned nil error for final command without prerequisites")
	}
}

// TestWriteEvidenceFinalMkdirFailure covers the os.MkdirAll error branch in
// WriteEvidence when the evidence pack directory cannot be created.
func TestWriteEvidenceFinalMkdirFailure(t *testing.T) {
	root := t.TempDir()
	writeAuthorityFixture(t, root)
	writePrerequisiteLedgerFixture(t, root, DefaultGoalID)
	report, err := Evaluate("goal-runtime-final", Options{Root: root})
	if err != nil {
		t.Fatalf("Evaluate returned error: %v", err)
	}
	// Make the evidence pack parent directory a file so MkdirAll fails.
	blocker := filepath.Join(root, filepath.FromSlash(EvidenceLedgerPath))
	if err := os.MkdirAll(filepath.Dir(blocker), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(blocker, []byte("blocker"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := WriteEvidence(root, report); err == nil {
		t.Fatal("WriteEvidence returned nil error for MkdirAll failure")
	}
}

// TestValidateFinalPrerequisitesReportsMissingAndIncomplete covers the missing
// and incomplete prerequisite branches directly.
func TestValidateFinalPrerequisitesReportsMissingAndIncomplete(t *testing.T) {
	root := t.TempDir()

	// Missing ledger entirely.
	gaps := validateFinalPrerequisites(root, DefaultGoalID)
	if !containsSubstring(gaps, "missing prerequisite evidence ledger") {
		t.Fatalf("gaps = %v; want missing ledger gap", gaps)
	}

	// Incomplete prerequisite entry: present but not passed.
	ledgerDir := filepath.Join(root, filepath.Dir(filepath.FromSlash(SourceLedgerPath)))
	if err := os.MkdirAll(ledgerDir, 0o755); err != nil {
		t.Fatal(err)
	}
	incomplete := LedgerEntry{
		SchemaVersion: "goalcli-mva/v1",
		GoalID:        DefaultGoalID,
		Command:       "goal-acceptance",
		Status:        "failed",
		MVAStatus:     "complete",
		Blocking:      true,
	}
	data, err := json.Marshal(incomplete)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(ledgerDir, "ledger.jsonl"), append(data, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}

	gaps = validateFinalPrerequisites(root, DefaultGoalID)
	if !containsSubstring(gaps, "incomplete prerequisite evidence") {
		t.Fatalf("gaps = %v; want incomplete prerequisite gap", gaps)
	}
	if !containsSubstring(gaps, "missing prerequisite evidence") {
		t.Fatalf("gaps = %v; want missing prerequisite gap for remaining commands", gaps)
	}
}

// TestReadLedgerEntriesRejectsInvalidJSON covers the json.Unmarshal error
// branch in readLedgerEntries.
func TestReadLedgerEntriesRejectsInvalidJSON(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "ledger.jsonl")
	if err := os.WriteFile(path, []byte("{not json\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := readLedgerEntries(path); err == nil {
		t.Fatal("readLedgerEntries returned nil error for invalid JSON")
	}
}

// TestWriteEvidenceFinalRejectsPrerequisitesWithPassingReport covers the
// validateFinalPrerequisites gap branch in WriteEvidence by feeding a passing
// final report into a root that lacks the prerequisite ledger.
func TestWriteEvidenceFinalRejectsPrerequisitesWithPassingReport(t *testing.T) {
	root := t.TempDir()
	writeAuthorityFixture(t, root)
	report := Report{
		SchemaVersion:    "goalcli-mva/v1",
		Command:          finalRuntimeCommand,
		Status:           "passed",
		GoalID:           DefaultGoalID,
		MVAStatus:        "complete",
		Blocking:         true,
		LedgerPath:       SourceLedgerPath,
		EvidencePackPath: EvidenceLedgerPath + DefaultGoalID + ".json",
	}
	if err := WriteEvidence(root, report); err == nil {
		t.Fatal("WriteEvidence returned nil error for passing final report without prerequisites")
	}
}

// TestWriteEvidenceFinalReportsPackWriteFailure covers the os.WriteFile error
// branch when the evidence pack path is a directory.
func TestWriteEvidenceFinalReportsPackWriteFailure(t *testing.T) {
	root := t.TempDir()
	writeAuthorityFixture(t, root)
	writePrerequisiteLedgerFixture(t, root, DefaultGoalID)
	report, err := Evaluate("goal-runtime-final", Options{Root: root})
	if err != nil {
		t.Fatalf("Evaluate returned error: %v", err)
	}
	// Create the evidence pack path as a directory so os.WriteFile fails.
	packPath := filepath.Join(root, filepath.FromSlash(report.EvidencePackPath))
	if err := os.MkdirAll(packPath, 0o755); err != nil {
		t.Fatal(err)
	}

	if err := WriteEvidence(root, report); err == nil {
		t.Fatal("WriteEvidence returned nil error for pack write failure")
	}
}

// TestUpsertLedgerEntryReportsMkdirFailure covers the os.MkdirAll error branch
// when the ledger parent directory cannot be created.
func TestUpsertLedgerEntryReportsMkdirFailure(t *testing.T) {
	root := t.TempDir()
	// Make a parent that is a file so MkdirAll fails when creating its child dir.
	blocker := filepath.Join(root, "blocker")
	if err := os.WriteFile(blocker, []byte("file"), 0o644); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(blocker, "evidence", "ledger.jsonl")

	err := upsertLedgerEntry(path, ledgerEntryForReport(Report{
		SchemaVersion: "goalcli-mva/v1",
		GoalID:        DefaultGoalID,
		Command:       "goal-acceptance",
		Status:        "passed",
		MVAStatus:     "complete",
		Blocking:      true,
	}))
	if err == nil {
		t.Fatal("upsertLedgerEntry returned nil error for MkdirAll failure")
	}
}

// TestUpsertLedgerEntryReportsReadFailure covers the non-IsNotExist read
// error branch in upsertLedgerEntry.
func TestUpsertLedgerEntryReportsReadFailure(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, filepath.FromSlash("release/evidence"))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "ledger.jsonl")
	// Create the ledger path as a directory so os.ReadFile fails with a non-IsNotExist error.
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := upsertLedgerEntry(path, ledgerEntryForReport(Report{
		SchemaVersion: "goalcli-mva/v1",
		GoalID:        DefaultGoalID,
		Command:       "goal-acceptance",
		Status:        "passed",
		MVAStatus:     "complete",
		Blocking:      true,
	})); err == nil {
		t.Fatal("upsertLedgerEntry returned nil error for non-IsNotExist read failure")
	}
}

// TestUpsertLedgerEntryPreservesUnrelatedEntriesAndDedupes covers the existing-
// entry dedupe path and the append of unrelated entries.
func TestUpsertLedgerEntryDedupesAndPreservesUnrelated(t *testing.T) {
	root := t.TempDir()
	ledgerDir := filepath.Join(root, "evidence")
	if err := os.MkdirAll(ledgerDir, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(ledgerDir, "ledger.jsonl")

	target := ledgerEntryForReport(Report{
		SchemaVersion:    "goalcli-mva/v1",
		GoalID:           DefaultGoalID,
		Command:          "goal-acceptance",
		Status:           "passed",
		MVAStatus:        "complete",
		Blocking:         true,
		EvidencePackPath: "release/evidence/pack.json",
	})
	other := LedgerEntry{
		SchemaVersion:    "goalcli-mva/v1",
		GoalID:           DefaultGoalID,
		Command:          "goal-delivery",
		Status:           "passed",
		MVAStatus:        "complete",
		Blocking:         true,
		EvidencePackPath: "release/evidence/delivery.json",
	}
	otherData, err := json.Marshal(other)
	if err != nil {
		t.Fatal(err)
	}
	targetData, err := json.Marshal(target)
	if err != nil {
		t.Fatal(err)
	}
	seed := append(append(otherData, '\n'), append(targetData, '\n')...)
	if err := os.WriteFile(path, seed, 0o644); err != nil {
		t.Fatal(err)
	}

	// Re-upsert target: should dedupe (drop existing target) and re-append once.
	if err := upsertLedgerEntry(path, target); err != nil {
		t.Fatalf("upsertLedgerEntry returned error: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(string(data), `"command":"goal-acceptance"`) != 1 {
		t.Fatalf("ledger = %s; want one goal-acceptance entry after dedupe", data)
	}
	if !strings.Contains(string(data), `"command":"goal-delivery"`) {
		t.Fatalf("ledger = %s; want unrelated goal-delivery entry preserved", data)
	}
}

// TestUpsertLedgerEntryReportsWriteFailure covers the os.WriteFile error branch.
func TestUpsertLedgerEntryReportsWriteFailure(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "ledger.jsonl")

	old := goalruntimeWriteFile
	goalruntimeWriteFile = func(name string, data []byte, perm os.FileMode) error {
		if name == path {
			return errors.New("write failed")
		}
		return old(name, data, perm)
	}
	t.Cleanup(func() { goalruntimeWriteFile = old })

	err := upsertLedgerEntry(path, ledgerEntryForReport(Report{
		SchemaVersion: "goalcli-mva/v1",
		GoalID:        DefaultGoalID,
		Command:       "goal-acceptance",
		Status:        "passed",
		MVAStatus:     "complete",
		Blocking:      true,
	}))
	if err == nil {
		t.Fatal("upsertLedgerEntry returned nil error for write failure")
	}
	if !strings.Contains(err.Error(), "write evidence ledger") {
		t.Fatalf("upsertLedgerEntry error = %v; want write evidence ledger error", err)
	}
}
