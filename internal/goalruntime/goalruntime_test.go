package goalruntime

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEvaluateGoalRuntimeFinalPassesWithAuthority(t *testing.T) {
	root := t.TempDir()
	writeAuthorityFixture(t, root)
	writePrerequisiteLedgerFixture(t, root, DefaultGoalID)

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
	if report.Gate != "G12_G16_FINAL" {
		t.Fatalf("gate = %q; want G12_G16_FINAL", report.Gate)
	}
	if !report.Blocking {
		t.Fatalf("blocking = false; want final runtime evidence to be blocking")
	}
	if report.MVAStatus != "complete" {
		t.Fatalf("mva_status = %q; want complete", report.MVAStatus)
	}
	for _, gate := range report.Gates {
		if !gate.Blocking {
			t.Fatalf("gate = %#v; want goalcli MVA gates to be blocking", gate)
		}
	}
	if !contains(report.Evidence, "source_evidence_ledger="+SourceLedgerPath) {
		t.Fatalf("evidence = %#v; want source ledger path", report.Evidence)
	}
	if !contains(report.Evidence, "generated_evidence_pack="+EvidenceLedgerPath) {
		t.Fatalf("evidence = %#v; want generated evidence pack path", report.Evidence)
	}
	if !contains(report.Evidence, "requires=goal-certify") {
		t.Fatalf("evidence = %#v; want final gate dependency", report.Evidence)
	}
	for _, want := range downstreamAdoptionBoundaryEvidence() {
		if !contains(report.Evidence, want) {
			t.Fatalf("evidence = %#v; want downstream adoption boundary %s", report.Evidence, want)
		}
	}
	if !containsSubstring(report.Details, "完成状态由本地 authority 校验和 evidence 写入共同证明") {
		t.Fatalf("details = %#v; want completion evidence boundary", report.Details)
	}
	if len(report.AuthorityPaths) == 0 {
		t.Fatalf("authority_paths is empty")
	}
}

func TestEvaluateRejectsUnknownCommand(t *testing.T) {
	if _, err := Evaluate("not-a-goalcli-command", Options{}); err == nil {
		t.Fatalf("Evaluate returned nil error for unknown command")
	}
}

func TestCommandsReturnsSupportedGoalRuntimeCommands(t *testing.T) {
	want := []string{
		"goal-acceptance",
		"goal-delivery",
		"goal-handover",
		"goal-downstream-adoption",
		"goal-certify",
		"goal-runtime-final",
	}

	got := Commands()
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("Commands() = %#v; want %#v", got, want)
	}

	got[0] = "changed"
	if Commands()[0] != want[0] {
		t.Fatalf("Commands() returned mutable command storage")
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
	if !containsSubstring(report.Gaps, ".worktree/goalcli-v0.1.0-plan.md") {
		t.Fatalf("gaps = %#v; want root plan gap", report.Gaps)
	}
	if report.MVAStatus != "not-complete" {
		t.Fatalf("mva_status = %q; want not-complete when authority is missing", report.MVAStatus)
	}
}

func TestEvaluateRenderedDownstreamSkipsSourceOnlyAuthorityPaths(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module github.com/ZoneCNH/kernel\n"), 0o644); err != nil {
		t.Fatalf("write downstream go.mod: %v", err)
	}
	writeAuthorityPaths(t, root, portableAuthorityPaths)

	report, err := Evaluate("goal-acceptance", Options{Root: root})
	if err != nil {
		t.Fatalf("Evaluate returned error: %v", err)
	}
	if report.Status != "passed" {
		t.Fatalf("status = %q; gaps %#v", report.Status, report.Gaps)
	}
	for _, path := range sourceOnlyAuthorityPaths {
		if contains(report.AuthorityPaths, path) {
			t.Fatalf("authority_paths = %#v; want source-only path %s skipped", report.AuthorityPaths, path)
		}
		if containsSubstring(report.Gaps, path) {
			t.Fatalf("gaps = %#v; want source-only path %s skipped", report.Gaps, path)
		}
	}
	if !containsSubstring(report.Details, "rendered downstream") {
		t.Fatalf("details = %#v; want rendered downstream boundary", report.Details)
	}
}

func TestEvaluateAcceptsNonDefaultGoalID(t *testing.T) {
	root := t.TempDir()
	writeAuthorityFixture(t, root)

	report, err := Evaluate("goal-acceptance", Options{Root: root, GoalID: "GOAL-X"})
	if err != nil {
		t.Fatalf("Evaluate returned error: %v", err)
	}
	if report.GoalID != "GOAL-X" {
		t.Fatalf("goal_id = %q; want GOAL-X", report.GoalID)
	}
	if report.Status != "passed" {
		t.Fatalf("status = %q; gaps %#v", report.Status, report.Gaps)
	}
	if !containsSubstring(report.Details, "non-default goal_id accepted") {
		t.Fatalf("details = %#v; want non-default goal detail", report.Details)
	}
}

func TestEvaluateGoModWithoutModuleLineKeepsSourceAuthorityPaths(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("go 1.23\n"), 0o644); err != nil {
		t.Fatalf("write go.mod without module line: %v", err)
	}
	writeAuthorityFixture(t, root)

	report, err := Evaluate("goal-acceptance", Options{Root: root})
	if err != nil {
		t.Fatalf("Evaluate returned error: %v", err)
	}
	if report.Status != "passed" {
		t.Fatalf("status = %q; gaps %#v", report.Status, report.Gaps)
	}
	for _, path := range sourceOnlyAuthorityPaths {
		if !contains(report.AuthorityPaths, path) {
			t.Fatalf("authority_paths = %#v; want source-only path %s retained", report.AuthorityPaths, path)
		}
	}
	if containsSubstring(report.Details, "rendered downstream") {
		t.Fatalf("details = %#v; want source root classification", report.Details)
	}
}

func TestEvaluateGoalcliGatePassesAsBlockingMVAContract(t *testing.T) {
	root := t.TempDir()
	writeAuthorityFixture(t, root)

	report, err := Evaluate("goal-acceptance", Options{Root: root})
	if err != nil {
		t.Fatalf("Evaluate returned error: %v", err)
	}
	if report.Status != "passed" || report.MVAStatus != "complete" || !report.Blocking {
		t.Fatalf("report = %#v; want passed complete blocking goalcli gate", report)
	}
	if len(report.Gates) != 1 || !report.Gates[0].Blocking {
		t.Fatalf("gates = %#v; want one blocking gate", report.Gates)
	}
}

func TestEvaluateDownstreamAdoptionDeclaresLocalScope(t *testing.T) {
	root := t.TempDir()
	writeAuthorityFixture(t, root)

	report, err := Evaluate("goal-downstream-adoption", Options{Root: root})
	if err != nil {
		t.Fatalf("Evaluate returned error: %v", err)
	}
	if report.Status != "passed" {
		t.Fatalf("status = %q; gaps %#v", report.Status, report.Gaps)
	}
	for _, want := range downstreamAdoptionBoundaryEvidence() {
		if !contains(report.Evidence, want) {
			t.Fatalf("evidence = %#v; want downstream adoption boundary %s", report.Evidence, want)
		}
	}
	if !containsSubstring(report.Details, "不声明 proof-based downstream adoption") {
		t.Fatalf("details = %#v; want proof-based adoption boundary", report.Details)
	}
}

func TestEvaluateGoalRuntimeFinalRequiresPrerequisiteLedger(t *testing.T) {
	root := t.TempDir()
	writeAuthorityFixture(t, root)

	report, err := Evaluate("goal-runtime-final", Options{Root: root})
	if err != nil {
		t.Fatalf("Evaluate returned error: %v", err)
	}
	if report.Status != "failed" || report.MVAStatus != "not-complete" {
		t.Fatalf("report = %#v; want failed not-complete when prerequisite ledger is missing", report)
	}
	if !containsSubstring(report.Gaps, SourceLedgerPath) {
		t.Fatalf("gaps = %#v; want source ledger prerequisite gap", report.Gaps)
	}
}

func TestEvaluateGoalRuntimeFinalIgnoresGeneratedPackWithoutSourceLedger(t *testing.T) {
	root := t.TempDir()
	writeAuthorityFixture(t, root)

	packPath := filepath.Join(root, filepath.FromSlash(EvidenceLedgerPath+DefaultGoalID+".json"))
	if err := os.MkdirAll(filepath.Dir(packPath), 0o755); err != nil {
		t.Fatalf("mkdir generated evidence pack path: %v", err)
	}
	if err := os.WriteFile(packPath, []byte(`{"status":"passed","mva_status":"complete"}`), 0o644); err != nil {
		t.Fatalf("write generated evidence pack fixture: %v", err)
	}

	report, err := Evaluate("goal-runtime-final", Options{Root: root})
	if err != nil {
		t.Fatalf("Evaluate returned error: %v", err)
	}
	if report.Status != "failed" || report.MVAStatus != "not-complete" {
		t.Fatalf("report = %#v; want generated pack ignored without source ledger prerequisites", report)
	}
	if !containsSubstring(report.Gaps, SourceLedgerPath) {
		t.Fatalf("gaps = %#v; want source ledger prerequisite gap despite generated pack", report.Gaps)
	}
}

func TestWriteEvidenceWritesPackAndLedgerIdempotently(t *testing.T) {
	root := t.TempDir()
	writeAuthorityFixture(t, root)
	writePrerequisiteLedgerFixture(t, root, "GOAL-20260603-XLIB-GOALCLI-001")
	report, err := Evaluate("goal-runtime-final", Options{
		Root:   root,
		GoalID: "GOAL-20260603-XLIB-GOALCLI-001",
	})
	if err != nil {
		t.Fatalf("Evaluate returned error: %v", err)
	}
	if err := WriteEvidence(root, report); err != nil {
		t.Fatalf("WriteEvidence returned error: %v", err)
	}
	if err := WriteEvidence(root, report); err != nil {
		t.Fatalf("second WriteEvidence returned error: %v", err)
	}

	packPath := filepath.Join(root, filepath.FromSlash(report.EvidencePackPath))
	pack, err := os.ReadFile(packPath)
	if err != nil {
		t.Fatalf("read evidence pack: %v", err)
	}
	if !strings.Contains(string(pack), `"mva_status": "complete"`) || !strings.Contains(string(pack), `"blocking": true`) {
		t.Fatalf("evidence pack = %s; want complete blocking report", pack)
	}
	for _, want := range downstreamAdoptionBoundaryEvidence() {
		if !strings.Contains(string(pack), want) {
			t.Fatalf("evidence pack = %s; want downstream adoption boundary %s", pack, want)
		}
	}
	ledgerPath := filepath.Join(root, filepath.FromSlash(SourceLedgerPath))
	ledger, err := os.ReadFile(ledgerPath)
	if err != nil {
		t.Fatalf("read evidence ledger: %v", err)
	}
	if strings.Count(string(ledger), `"command":"goal-runtime-final"`) != 1 {
		t.Fatalf("ledger = %s; want one idempotent final entry for %s", ledger, report.EvidencePackPath)
	}
}

func TestWriteEvidenceRecordsPrerequisiteLedgerEntry(t *testing.T) {
	root := t.TempDir()
	writeAuthorityFixture(t, root)
	report, err := Evaluate("goal-acceptance", Options{Root: root})
	if err != nil {
		t.Fatalf("Evaluate returned error: %v", err)
	}
	if err := WriteEvidence(root, report); err != nil {
		t.Fatalf("WriteEvidence returned error: %v", err)
	}
	ledger, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(SourceLedgerPath)))
	if err != nil {
		t.Fatalf("read evidence ledger: %v", err)
	}
	if !strings.Contains(string(ledger), `"command":"goal-acceptance"`) || !strings.Contains(string(ledger), `"mva_status":"complete"`) {
		t.Fatalf("ledger = %s; want goal-acceptance complete entry", ledger)
	}
	if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(report.EvidencePackPath))); !os.IsNotExist(err) {
		t.Fatalf("non-final WriteEvidence must not write generated pack, stat err = %v", err)
	}
}

func TestWriteEvidenceRejectsIncompleteReport(t *testing.T) {
	root := t.TempDir()
	writeAuthorityFixture(t, root)
	report, err := Evaluate("goal-acceptance", Options{Root: root})
	if err != nil {
		t.Fatalf("Evaluate returned error: %v", err)
	}
	report.Status = "failed"
	if err := WriteEvidence(root, report); err == nil {
		t.Fatalf("WriteEvidence returned nil error for incomplete report")
	}
}

func TestWriteEvidenceRejectsFinalBeforePrerequisites(t *testing.T) {
	root := t.TempDir()
	writeAuthorityFixture(t, root)
	report, err := Evaluate("goal-runtime-final", Options{Root: root, GoalID: "GOAL-X"})
	if err != nil {
		t.Fatalf("Evaluate returned error: %v", err)
	}
	report.Status = "passed"
	report.MVAStatus = "complete"
	report.Blocking = true

	err = WriteEvidence(root, report)
	if err == nil {
		t.Fatalf("WriteEvidence returned nil error for missing final prerequisites")
	}
	if !strings.Contains(err.Error(), "before prerequisites") {
		t.Fatalf("WriteEvidence error = %v; want final prerequisite error", err)
	}
}

func TestWriteEvidenceRejectsUnsupportedCommand(t *testing.T) {
	report := Report{
		Command:   "not-supported",
		Status:    "passed",
		MVAStatus: "complete",
		Blocking:  true,
	}
	if err := WriteEvidence(t.TempDir(), report); err == nil {
		t.Fatal("WriteEvidence returned nil error for unsupported command")
	}
}

func TestWriteEvidenceReportsEvidencePackDirectoryError(t *testing.T) {
	root := t.TempDir()
	report := passedFinalReport(t, root, "GOAL-X")
	if err := os.MkdirAll(filepath.Join(root, "release"), 0o755); err != nil {
		t.Fatalf("mkdir release parent: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "release", "evidence"), []byte("not a directory\n"), 0o644); err != nil {
		t.Fatalf("write evidence path blocker: %v", err)
	}

	err := WriteEvidence(root, report)
	if err == nil {
		t.Fatal("WriteEvidence returned nil error for evidence directory blocker")
	}
	if !strings.Contains(err.Error(), "create evidence pack directory") {
		t.Fatalf("WriteEvidence error = %v; want evidence pack directory error", err)
	}
}

func TestWriteEvidenceReportsEvidencePackWriteError(t *testing.T) {
	root := t.TempDir()
	report := passedFinalReport(t, root, "GOAL-X")
	packPath := filepath.Join(root, filepath.FromSlash(report.EvidencePackPath))
	if err := os.MkdirAll(packPath, 0o755); err != nil {
		t.Fatalf("mkdir evidence pack file path as directory: %v", err)
	}

	err := WriteEvidence(root, report)
	if err == nil {
		t.Fatal("WriteEvidence returned nil error for evidence pack write blocker")
	}
	if !strings.Contains(err.Error(), "write evidence pack") {
		t.Fatalf("WriteEvidence error = %v; want evidence pack write error", err)
	}
}

func TestWriteEvidenceReportsMarshalFailure(t *testing.T) {
	root := t.TempDir()
	report := passedFinalReport(t, root, "GOAL-X")
	old := goalruntimeMarshalIndent
	goalruntimeMarshalIndent = func(any, string, string) ([]byte, error) {
		return nil, errors.New("marshal failed")
	}
	t.Cleanup(func() { goalruntimeMarshalIndent = old })

	err := WriteEvidence(root, report)
	if err == nil {
		t.Fatal("WriteEvidence returned nil error for marshal failure")
	}
	if !strings.Contains(err.Error(), "marshal evidence pack") {
		t.Fatalf("WriteEvidence error = %v; want marshal evidence pack error", err)
	}
}

func TestWriteEvidenceDefaultsRootToCurrentDirectory(t *testing.T) {
	root := t.TempDir()
	writeAuthorityFixture(t, root)
	old, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatalf("chdir %s: %v", root, err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(old); err != nil {
			t.Fatalf("restore cwd %s: %v", old, err)
		}
	})

	report, err := Evaluate("goal-acceptance", Options{})
	if err != nil {
		t.Fatalf("Evaluate returned error: %v", err)
	}
	if err := WriteEvidence("", report); err != nil {
		t.Fatalf("WriteEvidence returned error: %v", err)
	}
	if _, err := os.Stat(filepath.FromSlash(SourceLedgerPath)); err != nil {
		t.Fatalf("stat default-root source ledger: %v", err)
	}
}

func TestValidateFinalPrerequisitesReportsIncompleteAndMissing(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, filepath.FromSlash(SourceLedgerPath))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir ledger path: %v", err)
	}
	entry := `{"schema_version":"goalcli-mva/v1","goal_id":"GOAL-X","command":"goal-acceptance","status":"failed","mva_status":"not-complete","blocking":false,"evidence_pack_path":"release/evidence/goalcli/GOAL-X.json"}`
	if err := os.WriteFile(path, []byte(entry+"\n"), 0o644); err != nil {
		t.Fatalf("write ledger fixture: %v", err)
	}

	gaps := validateFinalPrerequisites(root, "GOAL-X")
	if !containsSubstring(gaps, "incomplete prerequisite evidence for goal_id GOAL-X: goal-acceptance") {
		t.Fatalf("gaps = %#v; want incomplete goal-acceptance", gaps)
	}
	if !containsSubstring(gaps, "missing prerequisite evidence for goal_id GOAL-X: goal-delivery") {
		t.Fatalf("gaps = %#v; want missing goal-delivery", gaps)
	}
}

func TestReadLedgerEntriesRejectsInvalidLine(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "ledger.jsonl")
	if err := os.WriteFile(path, []byte(`{"goal_id":"ok"}`+"\nnot-json\n"), 0o644); err != nil {
		t.Fatalf("write ledger fixture: %v", err)
	}
	if _, err := readLedgerEntries(path); err == nil || !strings.Contains(err.Error(), "line 2") {
		t.Fatalf("readLedgerEntries() error = %v; want invalid line 2", err)
	}
}

func TestUpsertLedgerEntryReplacesMatchingEntryAndPreservesOthers(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, filepath.FromSlash(SourceLedgerPath))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir ledger path: %v", err)
	}
	initial := strings.Join([]string{
		`{"schema_version":"goalcli-mva/v1","goal_id":"GOAL-X","command":"goal-acceptance","status":"failed","mva_status":"not-complete","blocking":false,"evidence_pack_path":"release/evidence/goalcli/GOAL-X.json"}`,
		`{"schema_version":"goalcli-mva/v1","goal_id":"GOAL-X","command":"goal-delivery","status":"passed","mva_status":"complete","blocking":true,"evidence_pack_path":"release/evidence/goalcli/GOAL-X.json"}`,
		`not-json-but-preserved`,
		"",
	}, "\n")
	if err := os.WriteFile(path, []byte(initial), 0o644); err != nil {
		t.Fatalf("write initial ledger: %v", err)
	}

	entry := LedgerEntry{
		SchemaVersion:    "goalcli-mva/v1",
		GoalID:           "GOAL-X",
		Command:          "goal-acceptance",
		Status:           "passed",
		MVAStatus:        "complete",
		Blocking:         true,
		EvidencePackPath: "release/evidence/goalcli/GOAL-X.json",
	}
	if err := upsertLedgerEntry(path, entry); err != nil {
		t.Fatalf("upsertLedgerEntry returned error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read ledger: %v", err)
	}
	ledger := string(data)
	if strings.Count(ledger, `"command":"goal-acceptance"`) != 1 {
		t.Fatalf("ledger = %s; want one replaced goal-acceptance entry", ledger)
	}
	if !strings.Contains(ledger, `"status":"passed"`) ||
		!strings.Contains(ledger, `"command":"goal-delivery"`) ||
		!strings.Contains(ledger, "not-json-but-preserved") {
		t.Fatalf("ledger = %s; want replacement plus preserved entries", ledger)
	}
}

func TestUpsertLedgerEntryReportsFilesystemErrors(t *testing.T) {
	entry := LedgerEntry{
		SchemaVersion:    "goalcli-mva/v1",
		GoalID:           "GOAL-X",
		Command:          "goal-acceptance",
		Status:           "passed",
		MVAStatus:        "complete",
		Blocking:         true,
		EvidencePackPath: "release/evidence/goalcli/GOAL-X.json",
	}

	t.Run("create ledger directory", func(t *testing.T) {
		root := t.TempDir()
		blocker := filepath.Join(root, "blocked")
		if err := os.WriteFile(blocker, []byte("not a directory\n"), 0o644); err != nil {
			t.Fatalf("write ledger directory blocker: %v", err)
		}

		err := upsertLedgerEntry(filepath.Join(blocker, "ledger.jsonl"), entry)
		if err == nil {
			t.Fatal("upsertLedgerEntry returned nil error for directory blocker")
		}
		if !strings.Contains(err.Error(), "create evidence ledger directory") {
			t.Fatalf("upsertLedgerEntry error = %v; want directory creation error", err)
		}
	})

	t.Run("read ledger", func(t *testing.T) {
		root := t.TempDir()
		path := filepath.Join(root, "ledger.jsonl")
		if err := os.MkdirAll(path, 0o755); err != nil {
			t.Fatalf("mkdir ledger file path as directory: %v", err)
		}

		err := upsertLedgerEntry(path, entry)
		if err == nil {
			t.Fatal("upsertLedgerEntry returned nil error for unreadable ledger path")
		}
		if !strings.Contains(err.Error(), "read evidence ledger") {
			t.Fatalf("upsertLedgerEntry error = %v; want read ledger error", err)
		}
	})

	t.Run("write ledger", func(t *testing.T) {
		root := t.TempDir()
		path := filepath.Join(root, "ledger.jsonl")
		if err := os.WriteFile(path, []byte(`{"goal_id":"other"}`+"\n"), 0o644); err != nil {
			t.Fatalf("write readable ledger fixture: %v", err)
		}
		old := goalruntimeWriteFile
		goalruntimeWriteFile = func(name string, data []byte, perm os.FileMode) error {
			if name == path {
				return errors.New("write failed")
			}
			return old(name, data, perm)
		}
		t.Cleanup(func() { goalruntimeWriteFile = old })

		err := upsertLedgerEntry(path, entry)
		if err == nil {
			t.Fatal("upsertLedgerEntry returned nil error for readonly ledger")
		}
		if !strings.Contains(err.Error(), "write evidence ledger") {
			t.Fatalf("upsertLedgerEntry error = %v; want write ledger error", err)
		}
	})
}

func TestUpsertLedgerEntryReportsMarshalFailure(t *testing.T) {
	old := goalruntimeMarshal
	goalruntimeMarshal = func(any) ([]byte, error) {
		return nil, errors.New("marshal failed")
	}
	t.Cleanup(func() { goalruntimeMarshal = old })

	entry := LedgerEntry{
		SchemaVersion:    "goalcli-mva/v1",
		GoalID:           "GOAL-X",
		Command:          "goal-acceptance",
		Status:           "passed",
		MVAStatus:        "complete",
		Blocking:         true,
		EvidencePackPath: "release/evidence/goalcli/GOAL-X.json",
	}
	err := upsertLedgerEntry(filepath.Join(t.TempDir(), "ledger.jsonl"), entry)
	if err == nil {
		t.Fatal("upsertLedgerEntry returned nil error for marshal failure")
	}
	if !strings.Contains(err.Error(), "marshal evidence ledger entry") {
		t.Fatalf("upsertLedgerEntry error = %v; want marshal evidence ledger entry error", err)
	}
}

func writeAuthorityFixture(t *testing.T, root string) {
	t.Helper()
	writeAuthorityPaths(t, root, requiredAuthorityPaths(true))
}

func writeAuthorityPaths(t *testing.T, root string, paths []string) {
	t.Helper()
	for _, path := range paths {
		full := filepath.Join(root, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatalf("mkdir fixture path %s: %v", path, err)
		}
		if err := os.WriteFile(full, []byte("fixture\n"), 0o644); err != nil {
			t.Fatalf("write fixture path %s: %v", path, err)
		}
	}
}

func writePrerequisiteLedgerFixture(t *testing.T, root string, goalID string) {
	t.Helper()
	for _, command := range finalPrerequisiteCommands {
		report, err := Evaluate(command, Options{Root: root, GoalID: goalID})
		if err != nil {
			t.Fatalf("Evaluate prerequisite %s returned error: %v", command, err)
		}
		if err := WriteEvidence(root, report); err != nil {
			t.Fatalf("WriteEvidence prerequisite %s returned error: %v", command, err)
		}
	}
}

func passedFinalReport(t *testing.T, root string, goalID string) Report {
	t.Helper()
	writeAuthorityFixture(t, root)
	writePrerequisiteLedgerFixture(t, root, goalID)
	report, err := Evaluate("goal-runtime-final", Options{Root: root, GoalID: goalID})
	if err != nil {
		t.Fatalf("Evaluate final returned error: %v", err)
	}
	if report.Status != "passed" {
		t.Fatalf("final report status = %q; gaps %#v", report.Status, report.Gaps)
	}
	return report
}

func downstreamAdoptionBoundaryEvidence() []string {
	return []string{
		downstreamAdoptionClaimEvidence,
		downstreamAdoptionScopeEvidence,
		downstreamAdoptionProofEvidence,
		downstreamRepoWriteEvidence,
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
