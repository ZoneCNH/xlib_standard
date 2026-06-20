package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func setupRetroFixture(t *testing.T, withPatch bool) string {
	t.Helper()
	root := t.TempDir()
	must := func(rel, body string) {
		full := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}
		if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
			t.Fatalf("write %s: %v", rel, err)
		}
	}
	must(".agent/archive/retrospective.md", "# 复盘\n## 失败项\n- foo\n## 补丁\n- bar\n")
	must(".agent/docs/retrospective-template.md",
		"# Retrospective Template\n## Failure\n## Root Cause\n## Patch\n### Prompt Patch\n### Harness Patch\n### Rule Patch\n")
	if withPatch {
		must(".agent/harness/harness-patches.yaml",
			"schema: harness-patches\nentries:\n  - patch_id: PATCH-TEST-001\n    status: PROPOSED\n")
	} else {
		must(".agent/harness/harness-patches.yaml", "schema: harness-patches\nentries: []\n")
	}
	must(".agent/policies/prompt-patches.yaml", "schema: prompt-patches\nentries: []\n")
	must(".agent/policies/rule-patches.yaml", "schema: rule-patches\nentries: []\n")
	must(".agent/harness/gates/retro-gate.yaml", "gate_id: GATE-RETRO\n")
	return root
}

func TestSelfImprovingCheck_Lenient_Passes(t *testing.T) {
	root := setupRetroFixture(t, false)
	var out, errBuf bytes.Buffer
	code := runSelfImprovingCheck("self-improving-check", []string{"--root", root}, &out, &errBuf)
	if code != 0 {
		t.Fatalf("expected 0, got %d; stdout=%s", code, out.String())
	}
}

func TestSelfImprovingCheckSkipsUnreadablePatchRegistry(t *testing.T) {
	root := setupRetroFixture(t, false)
	registry := filepath.Join(root, ".agent/policies/rule-patches.yaml")
	if err := os.Remove(registry); err != nil {
		t.Fatalf("remove rule patches registry: %v", err)
	}
	if err := os.Mkdir(registry, 0o755); err != nil {
		t.Fatalf("mkdir rule patches registry: %v", err)
	}

	var out, errBuf bytes.Buffer
	code := runSelfImprovingCheck("self-improving-check", []string{"--root", root}, &out, &errBuf)
	if code != 0 {
		t.Fatalf("expected 0, got %d; stdout=%s stderr=%s", code, out.String(), errBuf.String())
	}
}

func TestSelfImprovingCheckArgumentBranches(t *testing.T) {
	root := setupRetroFixture(t, false)
	cases := []struct {
		name     string
		args     []string
		wantCode int
		wantOut  string
		wantErr  string
	}{
		{
			name:     "help",
			args:     []string{"--help"},
			wantCode: 0,
		},
		{
			name:     "invalid flag",
			args:     []string{"--root", root, "--missing"},
			wantCode: 2,
			wantErr:  "flag provided but not defined",
		},
		{
			name:     "positional argument",
			args:     []string{"--root", root, "extra"},
			wantCode: 2,
			wantErr:  "unexpected positional argument",
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			var out, errBuf bytes.Buffer
			code := runSelfImprovingCheck("self-improving-check", tt.args, &out, &errBuf)
			if code != tt.wantCode {
				t.Fatalf("exit=%d want %d; stdout=%s stderr=%s", code, tt.wantCode, out.String(), errBuf.String())
			}
			if tt.wantOut != "" && !bytes.Contains(out.Bytes(), []byte(tt.wantOut)) {
				t.Fatalf("stdout missing %q: %s", tt.wantOut, out.String())
			}
			if tt.wantErr != "" && !bytes.Contains(errBuf.Bytes(), []byte(tt.wantErr)) {
				t.Fatalf("stderr missing %q: %s", tt.wantErr, errBuf.String())
			}
		})
	}
}

func TestSelfImprovingCheck_Strict_FailsWithoutEntries(t *testing.T) {
	root := setupRetroFixture(t, false)
	var out, errBuf bytes.Buffer
	code := runSelfImprovingCheck("retro-check", []string{"--root", root, "--strict"}, &out, &errBuf)
	if code != 1 {
		t.Fatalf("expected 1, got %d; stdout=%s", code, out.String())
	}
	if !bytes.Contains(out.Bytes(), []byte("RULE-SI-001")) {
		t.Fatalf("expected RULE-SI-001 in output: %s", out.String())
	}
}

func TestSelfImprovingCheck_Strict_PassesWithEntry(t *testing.T) {
	root := setupRetroFixture(t, true)
	var out, errBuf bytes.Buffer
	code := runSelfImprovingCheck("self-improving-check", []string{"--root", root, "--strict"}, &out, &errBuf)
	if code != 0 {
		t.Fatalf("expected 0, got %d; stdout=%s", code, out.String())
	}
}

func TestSelfImprovingCheckReportsRetrospectiveTemplateAndSchemaGaps(t *testing.T) {
	root := setupRetroFixture(t, false)
	plainRetro := "# Notes\nNo keyword coverage here.\n"
	if err := os.WriteFile(filepath.Join(root, ".agent/archive/retrospective.md"), []byte(plainRetro), 0o644); err != nil {
		t.Fatalf("write retrospective: %v", err)
	}
	shortTemplate := "# Retrospective Template\n## Failure\n"
	if err := os.WriteFile(filepath.Join(root, ".agent/docs/retrospective-template.md"), []byte(shortTemplate), 0o644); err != nil {
		t.Fatalf("write template: %v", err)
	}
	badPatch := "status: PROPOSED\n"
	if err := os.WriteFile(filepath.Join(root, ".agent/policies/rule-patches.yaml"), []byte(badPatch), 0o644); err != nil {
		t.Fatalf("write patch registry: %v", err)
	}

	var out, errBuf bytes.Buffer
	code := runSelfImprovingCheck("self-improving-check", []string{"--root", root}, &out, &errBuf)
	if code != 1 {
		t.Fatalf("expected 1, got %d; stdout=%s stderr=%s", code, out.String(), errBuf.String())
	}
	for _, needle := range []string{
		"RULE-RETRO-001",
		"RULE-RETRO-002",
		"RULE-RETRO-CHECK-001",
		"RULE-SI-003",
	} {
		if !bytes.Contains(out.Bytes(), []byte(needle)) {
			t.Fatalf("expected %s in output: %s", needle, out.String())
		}
	}
}

func TestSelfImprovingCheckAcceptsKnownPatchStatuses(t *testing.T) {
	root := setupRetroFixture(t, false)
	registries := map[string]string{
		".agent/harness/harness-patches.yaml": `schema: harness-patches
entries:
  - patch_id: PATCH-ACCEPTED
    status: ACCEPTED
  - patch_id: PATCH-REJECTED
    status: REJECTED
`,
		".agent/policies/prompt-patches.yaml": `schema: prompt-patches
entries:
  - patch_id: PATCH-SUPERSEDED
    status: SUPERSEDED
  - patch_id: PATCH-IMPLEMENTED
    status: IMPLEMENTED
`,
		".agent/policies/rule-patches.yaml": `schema: rule-patches
entries:
  - patch_id: PATCH-RECONCILED
    status: reconciled_stub
`,
	}
	for rel, content := range registries {
		if err := os.WriteFile(filepath.Join(root, rel), []byte(content), 0o644); err != nil {
			t.Fatalf("write %s: %v", rel, err)
		}
	}
	var out, errBuf bytes.Buffer
	code := runSelfImprovingCheck("self-improving-check", []string{"--root", root, "--strict"}, &out, &errBuf)
	if code != 0 {
		t.Fatalf("expected 0, got %d; stdout=%s stderr=%s", code, out.String(), errBuf.String())
	}
	if !bytes.Contains(out.Bytes(), []byte("5 patch entries")) {
		t.Fatalf("expected patch entry count in output: %s", out.String())
	}
}

func TestSelfImprovingCheck_MissingFile_Fails(t *testing.T) {
	root := setupRetroFixture(t, false)
	if err := os.Remove(filepath.Join(root, ".agent/archive/retrospective.md")); err != nil {
		t.Fatalf("rm: %v", err)
	}
	var out, errBuf bytes.Buffer
	code := runSelfImprovingCheck("self-improving-check", []string{"--root", root}, &out, &errBuf)
	if code != 1 {
		t.Fatalf("expected 1, got %d; stdout=%s", code, out.String())
	}
	if !bytes.Contains(out.Bytes(), []byte("RULE-RETRO-001")) {
		t.Fatalf("expected RULE-RETRO-001 in output: %s", out.String())
	}
}

func TestSelfImprovingCheck_BadStatus_Fails(t *testing.T) {
	root := setupRetroFixture(t, false)
	bad := "schema: harness-patches\nentries:\n  - patch_id: PATCH-X\n    status: INVALID_X\n"
	if err := os.WriteFile(filepath.Join(root, ".agent/harness/harness-patches.yaml"), []byte(bad), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	var out, errBuf bytes.Buffer
	code := runSelfImprovingCheck("self-improving-check", []string{"--root", root}, &out, &errBuf)
	if code != 1 {
		t.Fatalf("expected 1, got %d; stdout=%s", code, out.String())
	}
	if !bytes.Contains(out.Bytes(), []byte("RULE-SI-002")) {
		t.Fatalf("expected RULE-SI-002 in output: %s", out.String())
	}
}
