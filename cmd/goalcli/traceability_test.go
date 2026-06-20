package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTraceFile(t *testing.T, dir, rel, content string) {
	t.Helper()
	full := filepath.Join(dir, rel)
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", full, err)
	}
	if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", full, err)
	}
}

const traceabilityValidMatrix = `# Traceability

| REQ | 需求摘要 | 主要产物 | 验证/Evidence | 收敛 owner |
| --- | --- | --- | --- | --- |
| REQ-001 | 示例需求 A | docs/standard/foo.md; .agent/state.yaml | docs/release.md; cmd/goalcli/main.go | agent-runtime |
| REQ-002 | 示例需求 B 含说明文字 | docs/spec.md; object-model; harness yaml | docs/release.md | agent-runtime |
`

func TestTraceabilityCheckPasses(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)
	writeTraceFile(t, dir, ".agent/traceability/traceability-matrix.md", traceabilityValidMatrix)
	writeTraceFile(t, dir, "docs/standard/foo.md", "x")
	writeTraceFile(t, dir, ".agent/state.yaml", "x")
	writeTraceFile(t, dir, "docs/release.md", "x")
	writeTraceFile(t, dir, "cmd/goalcli/main.go", "x")
	writeTraceFile(t, dir, "docs/spec.md", "x")

	var stdout, stderr bytes.Buffer
	code := runTraceabilityCheck(nil, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit=%d want 0; stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), `"status": "passed"`) {
		t.Fatalf("expected passed JSON, got: %s", stdout.String())
	}
}

func TestTraceabilityCheckHelpAndInvalidFlag(t *testing.T) {
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
			args:     []string{"--missing"},
			wantCode: 2,
			wantErr:  "ERROR: traceability-check invalid arguments",
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := runTraceabilityCheck(tt.args, &stdout, &stderr)
			if code != tt.wantCode {
				t.Fatalf("exit=%d want %d; stdout=%s stderr=%s", code, tt.wantCode, stdout.String(), stderr.String())
			}
			if tt.wantOut != "" && !strings.Contains(stdout.String(), tt.wantOut) {
				t.Fatalf("stdout missing %q: %s", tt.wantOut, stdout.String())
			}
			if tt.wantErr != "" && !strings.Contains(stderr.String(), tt.wantErr) {
				t.Fatalf("stderr missing %q: %s", tt.wantErr, stderr.String())
			}
		})
	}
}

func TestTraceabilityCheckCustomMatrixTextSuccess(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)
	matrix := `# Traceability

| REQ | 需求摘要 | 主要产物 | 验证/Evidence | 收敛 owner |
| --- | --- | --- | --- | --- |
| REQ-001 | 示例需求 A | docs/product.md | docs/evidence.md | agent-runtime |
`
	writeTraceFile(t, dir, "custom/matrix.md", matrix)
	writeTraceFile(t, dir, "docs/product.md", "x")
	writeTraceFile(t, dir, "docs/evidence.md", "x")

	var stdout, stderr bytes.Buffer
	code := runTraceabilityCheck([]string{"--matrix", "custom/matrix.md", "--json=false"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit=%d want 0; stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), "traceability-check passed: 1 REQ rows verified") {
		t.Fatalf("expected text success summary, got: %s", stdout.String())
	}
}

func TestTraceabilityCheckIgnoresNarrativeEvidenceRefs(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)
	matrix := `# Traceability

| REQ | 需求摘要 | 主要产物 | 验证/Evidence | 收敛 owner |
| --- | --- | --- | --- | --- |
| REQ-001 | 示例需求 A | docs/product.md | manual verification; docs/evidence.md | agent-runtime |
`
	writeTraceFile(t, dir, "custom/matrix.md", matrix)
	writeTraceFile(t, dir, "docs/product.md", "x")
	writeTraceFile(t, dir, "docs/evidence.md", "x")

	var stdout, stderr bytes.Buffer
	code := runTraceabilityCheck([]string{"--matrix", "custom/matrix.md"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit=%d want 0; stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), `"status": "passed"`) {
		t.Fatalf("expected passed JSON, got: %s", stdout.String())
	}
}

func TestTraceabilityCheckNoRequirementRowsExits2(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)
	matrix := `# Traceability

| REQ | 需求摘要 | 主要产物 | 验证/Evidence | 收敛 owner |
| --- | --- | --- | --- | --- |
| NOTE-001 | 非需求 | docs/product.md | docs/evidence.md | agent-runtime |
`
	writeTraceFile(t, dir, ".agent/traceability/traceability-matrix.md", matrix)
	writeTraceFile(t, dir, "docs/product.md", "x")
	writeTraceFile(t, dir, "docs/evidence.md", "x")

	var stdout, stderr bytes.Buffer
	code := runTraceabilityCheck(nil, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("exit=%d want 2; stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stderr.String(), "has no REQ rows") {
		t.Fatalf("expected no REQ rows error, got: %s", stderr.String())
	}
}

func TestTraceabilityCheckReportsPartialProofDepthAndLifecycleGap(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)
	writeTraceFile(t, dir, ".agent/traceability/traceability-matrix.md", traceabilityValidMatrix)
	writeTraceFile(t, dir, "docs/standard/foo.md", "x")
	writeTraceFile(t, dir, ".agent/state.yaml", "x")
	writeTraceFile(t, dir, "docs/release.md", "x")
	writeTraceFile(t, dir, "cmd/goalcli/main.go", "x")
	writeTraceFile(t, dir, "docs/spec.md", "x")

	var stdout, stderr bytes.Buffer
	code := runTraceabilityCheck(nil, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit=%d want 0; stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	out := stdout.String()
	for _, needle := range []string{
		"traceability_status=partial_implemented",
		"proof_depth=file_exists",
		"proof_depth_level=D3",
		"full_lifecycle_graph=gap",
	} {
		if !strings.Contains(out, needle) {
			t.Fatalf("expected traceability detail %q, got: %s", needle, out)
		}
	}
	if strings.Contains(out, "traceability_status=implemented") {
		t.Fatalf("traceability report must not overstate full implementation, got: %s", out)
	}
}

func TestTraceabilityCheckEmptyArtifactColumnExits9(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)
	matrix := `# Traceability

| REQ | 需求摘要 | 主要产物 | 验证/Evidence | 收敛 owner |
| --- | --- | --- | --- | --- |
| REQ-001 | 示例需求 A |  | docs/evidence.md | agent-runtime |
`
	writeTraceFile(t, dir, ".agent/traceability/traceability-matrix.md", matrix)
	writeTraceFile(t, dir, "docs/evidence.md", "x")

	var stdout, stderr bytes.Buffer
	code := runTraceabilityCheck(nil, &stdout, &stderr)
	if code != 9 {
		t.Fatalf("exit=%d want 9 (gap); stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), "主要产物 column is empty") {
		t.Fatalf("expected empty artifact gap, got: %s", stdout.String())
	}
}

func TestTraceabilityCheckMissingEvidencePathExits9(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)
	matrix := `# Traceability

| REQ | 需求摘要 | 主要产物 | 验证/Evidence | 收敛 owner |
| --- | --- | --- | --- | --- |
| REQ-001 | 示例需求 A | docs/standard/foo.md | docs/missing-evidence.md | agent-runtime |
`
	writeTraceFile(t, dir, ".agent/traceability/traceability-matrix.md", matrix)
	writeTraceFile(t, dir, "docs/standard/foo.md", "x")

	var stdout, stderr bytes.Buffer
	code := runTraceabilityCheck(nil, &stdout, &stderr)
	if code != 9 {
		t.Fatalf("exit=%d want 9 (gap); stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), "missing evidence ref: docs/missing-evidence.md") {
		t.Fatalf("expected missing evidence gap, got: %s", stdout.String())
	}
}

func TestTraceabilityCheckRejectsUnexpectedArgs(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)
	var stdout, stderr bytes.Buffer

	code := runTraceabilityCheck([]string{"unexpected"}, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("exit=%d want 2; stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stderr.String(), "accepts no positional arguments") {
		t.Fatalf("expected positional argument error, got: %s", stderr.String())
	}
}

func TestTraceabilityCheckTextOutputForFailures(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)
	matrix := `# Traceability

| REQ | 需求摘要 | 主要产物 | 验证/Evidence | 收敛 owner |
| --- | --- | --- | --- | --- |
| REQ-001 | 示例需求 A | docs/missing.md | docs/missing-evidence.md | agent-runtime |
`
	writeTraceFile(t, dir, ".agent/traceability/traceability-matrix.md", matrix)
	var stdout, stderr bytes.Buffer

	code := runTraceabilityCheck([]string{"--json=false"}, &stdout, &stderr)

	if code != 9 {
		t.Fatalf("exit=%d want 9; stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), "traceability-check failed: 2 gap(s)") {
		t.Fatalf("expected text failure summary, got: %s", stdout.String())
	}
	for _, needle := range []string{
		"missing artifact: docs/missing.md",
		"missing evidence ref: docs/missing-evidence.md",
	} {
		if !strings.Contains(stderr.String(), needle) {
			t.Fatalf("stderr missing %q in: %s", needle, stderr.String())
		}
	}
}

func TestVerifyArtifactExistsBranches(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)
	writeTraceFile(t, dir, "docs/readme.md", "x")
	writeTraceFile(t, dir, "docs/plain.txt", "x")

	cases := []struct {
		name     string
		artifact string
		wantErr  string
	}{
		{name: "file exists", artifact: "docs/readme.md"},
		{name: "directory wildcard", artifact: "docs/*"},
		{name: "wildcard not directory", artifact: "docs/plain.txt/*", wantErr: "not a directory: docs/plain.txt"},
		{name: "missing directory wildcard", artifact: "missing/*", wantErr: "directory not found: missing"},
		{name: "glob matches", artifact: "docs/*.md"},
		{name: "glob missing", artifact: "docs/*.none", wantErr: "glob matched no files: docs/*.none"},
		{name: "invalid glob", artifact: "docs/[", wantErr: "invalid glob"},
		{name: "missing file", artifact: "docs/missing.md", wantErr: "missing artifact: docs/missing.md"},
		{name: "directory stat error", artifact: "bad\x00path/*", wantErr: "directory stat error"},
		{name: "artifact stat error", artifact: "bad\x00path", wantErr: "artifact stat error"},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := verifyArtifactExists(tt.artifact)
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("verifyArtifactExists(%q) error = %v", tt.artifact, err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("verifyArtifactExists(%q) error = %v, want contains %q", tt.artifact, err, tt.wantErr)
			}
		})
	}
}

func TestVerifyArtifactExistsSkipsGitIgnoredPaths(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)
	binDir := filepath.Join(dir, "bin")
	writeFakeExecutable(t, binDir, "git", `#!/bin/sh
if [ "$1" = "check-ignore" ]; then
  exit 0
fi
exit 1
`)
	t.Setenv("PATH", binDir+string(os.PathListSeparator)+os.Getenv("PATH"))

	if !isGitIgnored("release/generated.json") {
		t.Fatal("expected fake git check-ignore to mark path ignored")
	}
	if err := verifyArtifactExists("release/generated.json"); err != nil {
		t.Fatalf("ignored artifact should be skipped, got error: %v", err)
	}
}

func TestTraceabilityCheckEmptyEvidenceColumnExits9(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)
	matrix := `# Traceability

| REQ | 需求摘要 | 主要产物 | 验证/Evidence | 收敛 owner |
| --- | --- | --- | --- | --- |
| REQ-001 | 示例需求 A | docs/standard/foo.md |  | agent-runtime |
`
	writeTraceFile(t, dir, ".agent/traceability/traceability-matrix.md", matrix)
	writeTraceFile(t, dir, "docs/standard/foo.md", "x")

	var stdout, stderr bytes.Buffer
	code := runTraceabilityCheck(nil, &stdout, &stderr)
	if code != 9 {
		t.Fatalf("exit=%d want 9 (gap); stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), "验证/Evidence column is empty") {
		t.Fatalf("expected empty evidence gap, got: %s", stdout.String())
	}
}

func TestTraceabilityCheckGapExits9(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)
	writeTraceFile(t, dir, ".agent/traceability/traceability-matrix.md", traceabilityValidMatrix)
	// 故意只创建部分产物, 其余缺失
	writeTraceFile(t, dir, "docs/standard/foo.md", "x")

	var stdout, stderr bytes.Buffer
	code := runTraceabilityCheck(nil, &stdout, &stderr)
	if code != 9 {
		t.Fatalf("exit=%d want 9 (gap); stdout=%s", code, stdout.String())
	}
	if !strings.Contains(stdout.String(), `"status": "failed"`) {
		t.Fatalf("expected failed status JSON, got: %s", stdout.String())
	}
	if !strings.Contains(stdout.String(), "missing artifact") {
		t.Fatalf("expected gap detail, got: %s", stdout.String())
	}
}

func TestTraceabilityCheckMissingMatrixExits2(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)
	var stdout, stderr bytes.Buffer
	code := runTraceabilityCheck(nil, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("exit=%d want 2 (parse error); stderr=%s", code, stderr.String())
	}
}

func TestSplitCellTokensPreservesLeadingDot(t *testing.T) {
	tokens := splitCellTokens(".agent/foo.yaml; docs/goal.md v2.9.3 Complete; .github/workflows/ci.yml")
	want := []string{".agent/foo.yaml", "docs/goal.md", ".github/workflows/ci.yml"}
	if len(tokens) != len(want) {
		t.Fatalf("len=%d want %d; tokens=%#v", len(tokens), len(want), tokens)
	}
	for i := range tokens {
		if tokens[i] != want[i] {
			t.Errorf("tokens[%d]=%q want %q", i, tokens[i], want[i])
		}
	}
}

func TestLooksLikePath(t *testing.T) {
	cases := map[string]bool{
		"":                          false,
		".agent/foo.yaml":           true,
		".github/workflows/ci.yml":  true,
		".gitignore":                true,
		".golangci.yml":             true,
		"Makefile":                  true,
		"README.md":                 true,
		"cmd/goalcli":               true,
		"contracts/foo.schema.json": true,
		"docs/spec.md":              true,
		"examples/demo.yaml":        true,
		"foo/bar":                   false,
		"foo/bar.md":                true,
		"foo/bar.mk":                true,
		"foo/bar.unknown":           false,
		"harness yaml":              false,
		"internal/foo":              true,
		"object-model":              false,
		"pkg/foo":                   true,
		"release/foo.json":          true,
		"renovate.json":             true,
		"render worker":             false,
		"scripts/build.sh":          true,
		"testkit/foo.go":            true,
	}
	for in, want := range cases {
		if got := looksLikePath(in); got != want {
			t.Errorf("looksLikePath(%q)=%v want %v", in, got, want)
		}
	}
}
