package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestEmitReportReportsMarshalFailure(t *testing.T) {
	old := governanceMarshalIndent
	governanceMarshalIndent = func(any, string, string) ([]byte, error) {
		return nil, errors.New("marshal failed")
	}
	t.Cleanup(func() { governanceMarshalIndent = old })

	var stdout bytes.Buffer
	got := emitReport(&stdout, "cmd", "failed", nil, []string{"gap"})
	if got != 1 {
		t.Fatalf("emitReport() = %d, want 1", got)
	}
	if stdout.String() != "{\"command\":\"cmd\",\"status\":\"failed\"}\n" {
		t.Fatalf("stdout = %q; want fallback report", stdout.String())
	}
}

func TestTrackedDocsMarkdownFilesFallsBackOutsideGit(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is required to exercise trackedDocsMarkdownFiles fallback")
	}

	root := t.TempDir()
	chdir(t, root)
	writeDebtCLIFile(t, root, "docs/z.md", "# Z\n")
	writeDebtCLIFile(t, root, "docs/nested/a.md", "# A\n")
	writeDebtCLIFile(t, root, "docs/ignore.txt", "not markdown\n")

	paths, err := trackedDocsMarkdownFiles()
	if err != nil {
		t.Fatalf("trackedDocsMarkdownFiles() error = %v; want filesystem fallback", err)
	}
	want := []string{"docs/nested/a.md", "docs/z.md"}
	if !reflect.DeepEqual(paths, want) {
		t.Fatalf("trackedDocsMarkdownFiles() = %#v; want %#v", paths, want)
	}
}

func TestFilesystemDocsMarkdownFilesReportsWalkError(t *testing.T) {
	paths, err := filesystemDocsMarkdownFiles("missing-docs-root")
	if err == nil {
		t.Fatalf("filesystemDocsMarkdownFiles() error = nil paths = %#v; want missing root error", paths)
	}
}

func TestCanonicalRepoPathBranches(t *testing.T) {
	tests := []struct {
		raw    string
		want   string
		wantOK bool
	}{
		{raw: "", want: "", wantOK: false},
		{raw: "docs/readme.md", want: "docs/readme.md", wantOK: true},
		{raw: "docs//guide.md", want: "docs/guide.md", wantOK: false},
		{raw: "../escape", want: "../escape", wantOK: false},
		{raw: ".", want: ".", wantOK: false},
		{raw: "/abs/path", want: "/abs/path", wantOK: false},
		{raw: `docs\guide.md`, want: "docs/guide.md", wantOK: false},
	}

	for _, tt := range tests {
		t.Run(tt.raw, func(t *testing.T) {
			got, ok := canonicalRepoPath(tt.raw)
			if got != tt.want || ok != tt.wantOK {
				t.Fatalf("canonicalRepoPath(%q) = %q, %v; want %q, %v", tt.raw, got, ok, tt.want, tt.wantOK)
			}
		})
	}
}

func TestEnvAndContextProfileHelpers(t *testing.T) {
	t.Setenv("XLIB_TEST_ENV_DEFAULT", "")
	if got := envDefault("XLIB_TEST_ENV_DEFAULT", "fallback"); got != "fallback" {
		t.Fatalf("envDefault() with empty env = %q; want fallback", got)
	}
	t.Setenv("XLIB_TEST_ENV_DEFAULT", "configured")
	if got := envDefault("XLIB_TEST_ENV_DEFAULT", "fallback"); got != "configured" {
		t.Fatalf("envDefault() with configured env = %q; want configured", got)
	}

	if got := normalizeContextProfile("fast"); got != "lite" {
		t.Fatalf("normalizeContextProfile(fast) = %q; want lite", got)
	}
	if got := normalizeContextProfile("release"); got != "release" {
		t.Fatalf("normalizeContextProfile(release) = %q; want release", got)
	}
	if !validContextProfileName("standard") || validContextProfileName("unknown") {
		t.Fatalf("validContextProfileName returned unexpected validity")
	}
	if got := mapContextAliasToProfile("context-fast-check"); got != "lite" {
		t.Fatalf("mapContextAliasToProfile(context-fast-check) = %q; want lite", got)
	}
	if got := mapContextAliasToProfile("context-release"); got != "release" {
		t.Fatalf("mapContextAliasToProfile(context-release) = %q; want release", got)
	}
	if got, ok := contextGateProfile("context-full-check"); got != "full" || !ok {
		t.Fatalf("contextGateProfile(context-full-check) = %q, %v; want full, true", got, ok)
	}
	if got, ok := contextGateProfile("context-release"); got != "release" || !ok {
		t.Fatalf("contextGateProfile(context-release) = %q, %v; want release, true", got, ok)
	}
	if got, ok := contextGateProfile("not-a-context-gate"); got != "" || ok {
		t.Fatalf("contextGateProfile(unknown) = %q, %v; want empty, false", got, ok)
	}
}

func TestAppendMakefileTargetDependencyAndReferenceGaps(t *testing.T) {
	content := strings.Join([]string{
		".PHONY: target target-ok forbidden-ref",
		"target: dep-a forbidden",
		"target-ok: dep-a",
		"forbidden-ref:",
		"\t$(MAKE) forbidden",
	}, "\n")

	var gaps []string
	appendMakefileTargetDependencyGaps(content, "missing", []string{"dep-a"}, nil, &gaps)
	appendMakefileTargetDependencyGaps(content, "target", []string{"dep-a", "dep-missing"}, []string{"forbidden"}, &gaps)
	appendMakefileTargetForbiddenReferenceGaps(content, "missing-ref", []string{"bad"}, &gaps)
	appendMakefileTargetForbiddenReferenceGaps(content, "forbidden-ref", []string{"forbidden"}, &gaps)

	for _, want := range []string{
		"Makefile missing target block missing",
		"Makefile target missing dependency dep-missing",
		"Makefile target must not depend on forbidden",
		"Makefile missing target block missing-ref",
		"Makefile forbidden-ref must not reference forbidden",
	} {
		if !containsString(gaps, want) {
			t.Fatalf("gaps = %#v; want %q", gaps, want)
		}
	}

	var okGaps []string
	appendMakefileTargetDependencyGaps(content, "target-ok", []string{"dep-a"}, []string{"forbidden"}, &okGaps)
	appendMakefileTargetForbiddenReferenceGaps(content, "target-ok", []string{"forbidden"}, &okGaps)
	if len(okGaps) != 0 {
		t.Fatalf("ok target gaps = %#v; want none", okGaps)
	}
}

func TestMakefileTargetDependenciesContinuationBranches(t *testing.T) {
	cases := []struct {
		name    string
		content string
		want    []string
	}{
		{
			name: "continued dependency line",
			content: strings.Join([]string{
				"target: dep-a \\",
				" dep-b dep-c",
				"\t@echo ignored",
				"empty:",
				"\t@true",
				"other: dep-z",
			}, "\n"),
			want: []string{"dep-a", "dep-b", "dep-c"},
		},
		{
			name: "recipe stops continuation",
			content: strings.Join([]string{
				"target: dep-a \\",
				"\t@true",
			}, "\n"),
			want: []string{"dep-a"},
		},
		{
			name: "blank continuation line continues",
			content: strings.Join([]string{
				"target: dep-a \\",
				"",
				" dep-b",
			}, "\n"),
			want: []string{"dep-a", "dep-b"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := makefileTargetDependencies(tc.content, "target")
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("makefileTargetDependencies() = %#v; want %#v", got, tc.want)
			}
		})
	}

	content := cases[0].content
	if got := makefileTargetDependencies(content, "empty"); len(got) != 0 {
		t.Fatalf("makefileTargetDependencies(empty) = %#v; want empty", got)
	}
	if got := makefileTargetDependencies(content, "missing"); got != nil {
		t.Fatalf("makefileTargetDependencies(missing) = %#v; want nil", got)
	}
}

func TestAppendMakefileDuplicateGapsBranches(t *testing.T) {
	content := strings.Join([]string{
		"once:",
		"\t@true",
		"twice:",
		"\t@true",
		"twice:",
		"\t@true",
	}, "\n")

	var gaps []string
	appendMakefileDuplicateGaps(content, []string{"once", "twice", "missing"}, &gaps)
	if len(gaps) != 2 {
		t.Fatalf("gaps = %#v; want only duplicate and missing target gaps", gaps)
	}
	assertGovernanceGapsContain(t, gaps, "Makefile target twice must be defined exactly once, found 2")
	assertGovernanceGapsContain(t, gaps, "Makefile target missing must be defined exactly once, found 0")
}

func TestAgentIndexParsingRequiredFieldsAndEnums(t *testing.T) {
	entries := parseAgentIndexEntries(`
schema_version: "1.0"
files:
  - path: "docs/one.md"
    owner: docs
    status: active
  - path: .agent/two.yaml
    owner: []
    status: draft
metadata:
  - path: ignored.md
`)
	if len(entries) != 2 {
		t.Fatalf("parseAgentIndexEntries() len = %d entries = %#v; want 2", len(entries), entries)
	}
	if entries[0].path != "docs/one.md" || entries[1].path != ".agent/two.yaml" {
		t.Fatalf("entry paths = %q, %q; want parsed and trimmed paths", entries[0].path, entries[1].path)
	}

	var gaps []string
	appendRequiredAgentIndexField("index.yaml", entries[0], "owner", &gaps)
	appendRequiredAgentIndexField("index.yaml", entries[0], "purpose", &gaps)
	appendRequiredAgentIndexField("index.yaml", entries[1], "owner", &gaps)
	appendAgentIndexEnumGap("index.yaml", entries[0], "status", map[string]bool{"active": true}, &gaps)
	appendAgentIndexEnumGap("index.yaml", entries[1], "status", map[string]bool{"active": true}, &gaps)

	for _, want := range []string{
		"index.yaml docs/one.md missing purpose:",
		"index.yaml .agent/two.yaml missing owner:",
		"index.yaml .agent/two.yaml invalid status draft",
	} {
		if !containsString(gaps, want) {
			t.Fatalf("gaps = %#v; want %q", gaps, want)
		}
	}

	before := len(gaps)
	appendAgentIndexEnumGap("index.yaml", agentIndexEntry{path: "docs/no-status.md", block: "- path: docs/no-status.md\n"}, "status", map[string]bool{"active": true}, &gaps)
	if len(gaps) != before {
		t.Fatalf("gaps = %#v; want no gap for missing optional enum field", gaps)
	}
}

func TestAppendAgentIndexGapsBranches(t *testing.T) {
	root := t.TempDir()
	chdir(t, root)

	var gaps []string
	writeGovernanceHelperFile(t, ".agent/index.yaml", "files:\n")
	appendAgentIndexGaps(".agent/index.yaml", &gaps)
	for _, want := range []string{
		"missing schema_version:",
		"missing module: xlib-standard",
		"missing control_plane:",
		"must define files entries",
	} {
		assertGovernanceGapsContain(t, gaps, want)
	}

	if err := os.MkdirAll(filepath.Join(root, ".agent", "dir"), 0o755); err != nil {
		t.Fatalf("mkdir agent dir: %v", err)
	}
	writeGovernanceHelperFile(t, ".agent/dir/.keep", "keep\n")
	writeGovernanceHelperFile(t, ".agent/unindexed.md", "manual\n")
	gaps = nil
	writeGovernanceHelperFile(t, ".agent/index.yaml", `
schema_version: "1"
module: xlib-standard
control_plane: true
files:
  - path:
    layer: runtime_contract
    authority: source_of_truth
    mutability: hand_written
    owner: standard
    purpose: missing path
    status: active
  - path: .agent/missing.md
    layer: invalid
    authority: invalid
    mutability: invalid
    owner: standard
    purpose: missing file
    status: active
  - path: .agent/missing.md
    layer: runtime_contract
    authority: source_of_truth
    mutability: hand_written
    owner: standard
    purpose: duplicate
    status: active
  - path: docs/outside.md
    layer: runtime_contract
    authority: source_of_truth
    mutability: hand_written
    owner: standard
    purpose: outside
    status: active
  - path: .agent/dir
    layer: runtime_contract
    authority: source_of_truth
    mutability: hand_written
    owner: standard
    purpose: directory
    status: active
`)
	appendAgentIndexGaps(".agent/index.yaml", &gaps)
	for _, want := range []string{
		"contains file entry without path",
		"duplicate file entry .agent/missing.md",
		"docs/outside.md must stay under .agent/",
		"references missing .agent/missing.md",
		".agent/dir must be a file",
		"invalid layer invalid",
		"invalid authority invalid",
		"invalid mutability invalid",
		"missing file entry .agent/unindexed.md",
	} {
		assertGovernanceGapsContain(t, gaps, want)
	}
}

func TestShouldScanDocsFromFilesystemRejectsNonExitErrors(t *testing.T) {
	if shouldScanDocsFromFilesystem(errors.New("not git"), []byte("not a git repository")) {
		t.Fatalf("shouldScanDocsFromFilesystem() = true for non-exit error; want false")
	}
}

func TestTrackedDocsMarkdownFilesGitBranches(t *testing.T) {
	t.Run("uses git ls-files output", func(t *testing.T) {
		root := t.TempDir()
		chdir(t, root)
		fakeBin := filepath.Join(root, "bin")
		writeFakeExecutable(t, fakeBin, "git", `#!/bin/sh
printf '%b' 'docs/b.txt\000docs/z.md\000docs/a.md\000'
`)
		t.Setenv("PATH", fakeBin+string(os.PathListSeparator)+os.Getenv("PATH"))

		paths, err := trackedDocsMarkdownFiles()
		if err != nil {
			t.Fatalf("trackedDocsMarkdownFiles() error = %v; want nil", err)
		}
		want := []string{"docs/a.md", "docs/z.md"}
		if !reflect.DeepEqual(paths, want) {
			t.Fatalf("trackedDocsMarkdownFiles() = %#v; want %#v", paths, want)
		}
	})

	t.Run("returns git output with error", func(t *testing.T) {
		root := t.TempDir()
		chdir(t, root)
		fakeBin := filepath.Join(root, "bin")
		writeFakeExecutable(t, fakeBin, "git", `#!/bin/sh
printf '%s\n' 'fatal custom' >&2
exit 1
`)
		t.Setenv("PATH", fakeBin+string(os.PathListSeparator)+os.Getenv("PATH"))

		paths, err := trackedDocsMarkdownFiles()
		if err == nil || !strings.Contains(err.Error(), "fatal custom") {
			t.Fatalf("trackedDocsMarkdownFiles() paths = %#v error = %v; want git output error", paths, err)
		}
	})

	t.Run("returns silent git error", func(t *testing.T) {
		root := t.TempDir()
		chdir(t, root)
		fakeBin := filepath.Join(root, "bin")
		writeFakeExecutable(t, fakeBin, "git", `#!/bin/sh
exit 1
`)
		t.Setenv("PATH", fakeBin+string(os.PathListSeparator)+os.Getenv("PATH"))

		paths, err := trackedDocsMarkdownFiles()
		if err == nil {
			t.Fatalf("trackedDocsMarkdownFiles() paths = %#v error = nil; want git error", paths)
		}
	})
}

func TestWorktreeGateArgumentBranches(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := runWorktreeGate("worktree-check", []string{"-h"}, &stdout, &stderr); code != 0 {
		t.Fatalf("runWorktreeGate(-h) code = %d stderr = %q; want 0", code, stderr.String())
	}

	stdout.Reset()
	stderr.Reset()
	if code := runWorktreeGate("worktree-check", []string{"--bad"}, &stdout, &stderr); code != 2 {
		t.Fatalf("runWorktreeGate(invalid flag) code = %d; want 2", code)
	}

	stdout.Reset()
	stderr.Reset()
	if code := runWorktreeGate("worktree-check", []string{"unexpected"}, &stdout, &stderr); code != 2 {
		t.Fatalf("runWorktreeGate(unexpected) code = %d; want 2", code)
	}
	if !strings.Contains(stderr.String(), "unexpected positional argument") {
		t.Fatalf("stderr = %q; want positional argument error", stderr.String())
	}

	stdout.Reset()
	stderr.Reset()
	if code := runWorktreeGate("worktree-check", []string{"--context", "unknown"}, &stdout, &stderr); code != 2 {
		t.Fatalf("runWorktreeGate(invalid context) code = %d; want 2", code)
	}
	if !strings.Contains(stderr.String(), "invalid context") {
		t.Fatalf("stderr = %q; want invalid context error", stderr.String())
	}
}

func TestWorktreeGateReportsNonWorkerTree(t *testing.T) {
	root := t.TempDir()
	chdir(t, root)

	var stdout, stderr bytes.Buffer
	code := runWorktreeGate("worktree-check", []string{"--context", "local_write"}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("runWorktreeGate exit = %d; want 1; stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), "local_write requires a worker worktree") {
		t.Fatalf("stdout = %s; want worker worktree gap", stdout.String())
	}
}

func TestPRCheckArgumentBranches(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := runPRCheck([]string{"-h"}, strings.NewReader(""), &stdout, &stderr); code != 0 {
		t.Fatalf("runPRCheck(-h) code = %d stderr = %q; want 0", code, stderr.String())
	}

	stdout.Reset()
	stderr.Reset()
	if code := runPRCheck([]string{"--bad"}, strings.NewReader(""), &stdout, &stderr); code != 2 {
		t.Fatalf("runPRCheck(invalid flag) code = %d; want 2", code)
	}

	stdout.Reset()
	stderr.Reset()
	if code := runPRCheck([]string{"unexpected"}, strings.NewReader(""), &stdout, &stderr); code != 2 {
		t.Fatalf("runPRCheck(unexpected) code = %d; want 2", code)
	}
	if !strings.Contains(stderr.String(), "unexpected positional argument") {
		t.Fatalf("stderr = %q; want positional argument error", stderr.String())
	}

	stdout.Reset()
	stderr.Reset()
	if code := runPRCheck([]string{"--context", "unknown"}, strings.NewReader(""), &stdout, &stderr); code != 2 {
		t.Fatalf("runPRCheck(invalid context) code = %d; want 2", code)
	}
	if !strings.Contains(stderr.String(), "invalid context") {
		t.Fatalf("stderr = %q; want invalid context error", stderr.String())
	}

	stdout.Reset()
	stderr.Reset()
	if code := runPRCheck([]string{"--dry-run", "--context", "local_readonly"}, strings.NewReader(""), &stdout, &stderr); code != 0 {
		t.Fatalf("runPRCheck(dry-run) code = %d stdout = %q stderr = %q; want 0", code, stdout.String(), stderr.String())
	}
	for _, want := range []string{"mode=dry-run", "context=local_readonly", "delegates=worktree-check,lint,test"} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("stdout = %q; want %q", stdout.String(), want)
		}
	}
}

func TestContextProfileArgumentBranches(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := runContextProfile([]string{"-h"}, &stdout, &stderr); code != 0 {
		t.Fatalf("runContextProfile(-h) code = %d stderr = %q; want 0", code, stderr.String())
	}

	stdout.Reset()
	stderr.Reset()
	if code := runContextProfile([]string{"--bad"}, &stdout, &stderr); code != 2 {
		t.Fatalf("runContextProfile(invalid flag) code = %d; want 2", code)
	}

	stdout.Reset()
	stderr.Reset()
	if code := runContextProfile([]string{"unexpected"}, &stdout, &stderr); code != 2 {
		t.Fatalf("runContextProfile(unexpected) code = %d; want 2", code)
	}
	if !strings.Contains(stderr.String(), "unexpected positional argument") {
		t.Fatalf("stderr = %q; want positional argument error", stderr.String())
	}

	stdout.Reset()
	stderr.Reset()
	if code := runContextProfile([]string{"--profile", "unknown"}, &stdout, &stderr); code != 2 {
		t.Fatalf("runContextProfile(invalid profile) code = %d; want 2", code)
	}
	if !strings.Contains(stderr.String(), "invalid context profile") {
		t.Fatalf("stderr = %q; want invalid profile error", stderr.String())
	}

	stdout.Reset()
	stderr.Reset()
	if code := runContextProfile([]string{"--profile", "fast"}, &stdout, &stderr); code != 0 {
		t.Fatalf("runContextProfile(fast) code = %d stdout = %q stderr = %q; want 0", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), `"command": "context-profile"`) || !strings.Contains(stdout.String(), "profile=lite") {
		t.Fatalf("stdout = %q; want context-profile report with lite profile", stdout.String())
	}
}

func TestAppendGeneratedArtifactsGapsBranches(t *testing.T) {
	root := t.TempDir()
	chdir(t, root)

	var gaps []string
	appendGeneratedArtifactsGaps(".agent/registries/generated-artifacts.yaml", ".agent/index.yaml", &gaps)
	assertGovernanceGapsContain(t, gaps, "missing .agent/registries/generated-artifacts.yaml")

	gaps = nil
	writeGovernanceHelperFile(t, ".agent/registries/generated-artifacts.yaml", `
classification:
  artifact_class:
artifacts:
`)
	appendGeneratedArtifactsGaps(".agent/registries/generated-artifacts.yaml", ".agent/index.yaml", &gaps)
	for _, want := range []string{
		"missing schema_version:",
		"classification missing artifact_class",
		"classification missing authority",
		"classification missing validated_by",
		"must define generated artifact entries",
	} {
		assertGovernanceGapsContain(t, gaps, want)
	}

	gaps = nil
	writeGovernanceHelperFile(t, ".agent/registries/generated-artifacts.yaml", `
schema_version: "1"
classification:
  artifact_class: registry
  authority: source_of_truth
  validated_by: schema-check
artifacts:
  - path: ../bad
    classification: wrong
    source_control: tracked
    validated_by: unknown-validator
  - path:
    classification: generated_artifact
  - path: .agent/generated.txt
    classification: generated_artifact
    source_control: generated-only
    generated_by: goalcli dashboard-generate
`)
	appendGeneratedArtifactsGaps(".agent/registries/generated-artifacts.yaml", ".agent/index.yaml", &gaps)
	for _, want := range []string{
		"../bad must use canonical repo-relative slash path",
		"must set classification: generated_artifact",
		"must set source_control: generated-only",
		"missing generated_by",
		"validated_by unknown-validator is not a known goalcli or Makefile gate",
		"contains artifact without path",
		"missing validated_by",
		"missing .agent/index.yaml",
	} {
		assertGovernanceGapsContain(t, gaps, want)
	}

	gaps = nil
	writeGovernanceHelperFile(t, ".agent/index.yaml", `
files:
  - path: .agent/registries/generated-artifacts.yaml
    layer: wrong
    authority: wrong
    mutability: generated
`)
	appendGeneratedArtifactsGaps(".agent/registries/generated-artifacts.yaml", ".agent/index.yaml", &gaps)
	for _, want := range []string{
		"must set layer: registry",
		"must set authority: source_of_truth",
		"must set mutability: hand_written",
	} {
		assertGovernanceGapsContain(t, gaps, want)
	}

	gaps = nil
	writeGovernanceHelperFile(t, ".agent/index.yaml", `
files:
  - path: .agent/other.yaml
    layer: registry
    authority: source_of_truth
    mutability: hand_written
`)
	appendGeneratedArtifactsGaps(".agent/registries/generated-artifacts.yaml", ".agent/index.yaml", &gaps)
	assertGovernanceGapsContain(t, gaps, ".agent/index.yaml missing file entry .agent/registries/generated-artifacts.yaml")
}

func TestRulesConsistencyCheckReadAndArgumentBranches(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := runRulesConsistencyCheck([]string{"unexpected"}, &stdout, &stderr); code != 2 {
		t.Fatalf("runRulesConsistencyCheck(unexpected) code = %d; want 2", code)
	}
	if !strings.Contains(stderr.String(), "invalid arguments") {
		t.Fatalf("stderr = %q; want invalid arguments", stderr.String())
	}

	root := t.TempDir()
	chdir(t, root)
	stdout.Reset()
	stderr.Reset()
	if code := runRulesConsistencyCheck(nil, &stdout, &stderr); code != 1 {
		t.Fatalf("runRulesConsistencyCheck(missing canonical) code = %d stdout = %q stderr = %q; want 1", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stderr.String(), "read .agent/runtime/standard/goal-runtime-canonical.md") {
		t.Fatalf("stderr = %q; want canonical read error", stderr.String())
	}
}

func TestAppendRulesEnforcedByGapsBranches(t *testing.T) {
	root := t.TempDir()
	chdir(t, root)

	var gaps []string
	appendRulesEnforcedByGaps(".agent/registries/rules.yaml", &gaps)
	assertGovernanceGapsContain(t, gaps, "missing .agent/registries/rules.yaml")

	gaps = nil
	writeGovernanceHelperFile(t, ".agent/registries/rules.yaml", "rules:\n")
	appendRulesEnforcedByGaps(".agent/registries/rules.yaml", &gaps)
	assertGovernanceGapsContain(t, gaps, ".agent/registries/rules.yaml missing generated_from:")
	assertGovernanceGapsContain(t, gaps, ".agent/registries/rules.yaml must define rules")

	gaps = nil
	writeGovernanceHelperFile(t, ".agent/registries/rules.yaml", `
generated_from: governance
rules:
  - id: active-known-command
    status: active
    enforced_by: goalcli audit-goal
  - id: active-known-make
    status: active
    enforced_by: make fmt
  - id: active-known-script
    status: active
    enforced_by: scripts/check.sh
  - id: active-missing-enforcer
    status: active
  - id: active-unknown-enforcer
    status: active
    enforced_by: tool unknown
  - id: indexed-with-enforcer
    status: indexed
    enforced_by: goalcli audit-goal
  - id: deprecated-no-enforcer
    status: deprecated
  - id: invalid-status
    status: paused
`)
	writeGovernanceHelperFile(t, "scripts/check.sh", "#!/bin/sh\n")
	appendRulesEnforcedByGaps(".agent/registries/rules.yaml", &gaps)
	for _, want := range []string{
		".agent/registries/rules.yaml active-missing-enforcer active rule missing enforced_by",
		".agent/registries/rules.yaml active-unknown-enforcer enforced_by tool unknown is not tied to a known goalcli command, Makefile target, script, or hook",
		".agent/registries/rules.yaml indexed-with-enforcer indexed rule must not set enforced_by",
		".agent/registries/rules.yaml invalid-status invalid status paused",
	} {
		assertGovernanceGapsContain(t, gaps, want)
	}
	joined := strings.Join(gaps, "\n")
	for _, unexpected := range []string{"active-known-command", "active-known-make", "active-known-script", "deprecated-no-enforcer"} {
		if strings.Contains(joined, unexpected) {
			t.Fatalf("gaps = %#v; did not expect gap for %s", gaps, unexpected)
		}
	}
}

func TestKnownEnforcementRefBranches(t *testing.T) {
	root := t.TempDir()
	chdir(t, root)
	writeGovernanceHelperFile(t, "cmd/goalcli/main.go", "package main\n")
	writeGovernanceHelperFile(t, "scripts/check.sh", "#!/bin/sh\n")
	writeGovernanceHelperFile(t, ".githooks/pre-commit", "#!/bin/sh\n")

	commands := map[string]bool{"audit-goal": true, "bare-command": true}
	targets := map[string]bool{"fmt": true, "bare-target": true}

	tests := []struct {
		ref  string
		want bool
	}{
		{ref: "", want: false},
		{ref: "goalcli", want: true},
		{ref: "goalcli audit-goal", want: true},
		{ref: "goalcli fmt", want: true},
		{ref: "goalcli unknown", want: false},
		{ref: "make", want: false},
		{ref: "make fmt", want: true},
		{ref: "make unknown", want: false},
		{ref: "scripts/check.sh", want: true},
		{ref: "scripts/missing.sh", want: false},
		{ref: ".githooks/pre-commit", want: true},
		{ref: "bare-command", want: true},
		{ref: "bare-target", want: true},
		{ref: "unknown", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.ref, func(t *testing.T) {
			if got := knownEnforcementRef(tt.ref, commands, targets); got != tt.want {
				t.Fatalf("knownEnforcementRef(%q) = %v; want %v", tt.ref, got, tt.want)
			}
		})
	}
}

func TestKnownEnforcementRefRejectsDirectoriesAndTrimsFields(t *testing.T) {
	root := t.TempDir()
	chdir(t, root)
	if err := os.MkdirAll("cmd/goalcli/main.go", 0o755); err != nil {
		t.Fatalf("mkdir fake goalcli entry: %v", err)
	}
	if err := os.MkdirAll("scripts/check.sh", 0o755); err != nil {
		t.Fatalf("mkdir fake script entry: %v", err)
	}
	if err := os.MkdirAll(".githooks/pre-commit", 0o755); err != nil {
		t.Fatalf("mkdir fake hook entry: %v", err)
	}

	commands := map[string]bool{"audit-goal": true}
	targets := map[string]bool{"fmt": true}

	tests := []struct {
		ref  string
		want bool
	}{
		{ref: " goalcli audit-goal ", want: true},
		{ref: "make -s fmt", want: true},
		{ref: "make ", want: false},
		{ref: "goalcli", want: false},
		{ref: "scripts/check.sh", want: false},
		{ref: ".githooks/pre-commit", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.ref, func(t *testing.T) {
			if got := knownEnforcementRef(tt.ref, commands, targets); got != tt.want {
				t.Fatalf("knownEnforcementRef(%q) = %v; want %v", tt.ref, got, tt.want)
			}
		})
	}
}

func TestValidateIssueRegistryEntriesBranches(t *testing.T) {
	emptyGaps := validateIssueRegistryEntries("issues.yaml", nil)
	assertGovernanceGapsContain(t, emptyGaps, "issues.yaml must contain issue entries")

	entries := []issueRegistryEntry{
		{
			id: "P0-001",
			block: `- id: P0-001
  title: First issue
  status: implemented
  command: make test
  evidence:
    - reports/p0.md`,
		},
		{
			id: "P0-001",
			block: `- id: P0-001
  status: planned`,
		},
		{
			id: "BAD-001",
			block: `- id: BAD-001
  title: Bad issue
  status: implemented
  command: make test
  evidence: reports/bad.md`,
		},
		{
			id: "P1-002",
			block: `- id: P1-002
  title: P1 issue
  status: implemented
  command: make test
  evidence: []`,
		},
		{
			id: "P2-001",
			block: `- id: P2-001
  title: P2 issue
  status: implemented
  command: make test
  evidence: reports/p2.md`,
		},
		{
			id: "CTX-001",
			block: `- id: CTX-001
  title: Context issue
  status: implemented
  command: make test
  evidence: reports/ctx.md`,
		},
	}

	gaps := validateIssueRegistryEntries("issues.yaml", entries)
	for _, want := range []string{
		"issues.yaml duplicate issue id P0-001",
		"issues.yaml P0-001 missing title",
		"issues.yaml P0-001 status must be implemented",
		"issues.yaml P0-001 missing command",
		"issues.yaml P0-001 missing evidence",
		"issues.yaml invalid issue id BAD-001",
		"issues.yaml P1-002 missing evidence",
		"issues.yaml P1 ids must be contiguous; missing P1-001",
	} {
		assertGovernanceGapsContain(t, gaps, want)
	}

	validWithoutCTX := []issueRegistryEntry{
		{
			id: "P0-001",
			block: `- id: P0-001
  title: P0 issue
  status: implemented
  command: make p0
  evidence: reports/p0.md`,
		},
		{
			id: "P1-001",
			block: `- id: P1-001
  title: P1 issue
  status: implemented
  command: make p1
  evidence: reports/p1.md`,
		},
		{
			id: "P2-001",
			block: `- id: P2-001
  title: P2 issue
  status: implemented
  command: make p2
  evidence: reports/p2.md`,
		},
	}
	gaps = validateIssueRegistryEntries("issues.yaml", validWithoutCTX)
	assertGovernanceGapsContain(t, gaps, "issues.yaml missing CTX-001")
}

func TestYAMLSectionAndDuplicateValueHelperBranches(t *testing.T) {
	root := t.TempDir()
	chdir(t, root)
	writeGovernanceHelperFile(t, "registry.yaml", `
rules:
  - RULE-002
  - RULE-001
next:
  - ignored
map:
  first:
    value: x
  second:
other:
  ignored:
`)

	sequenceValues, err := yamlSequenceValuesInSection("registry.yaml", "rules")
	if err != nil {
		t.Fatalf("yamlSequenceValuesInSection() error = %v", err)
	}
	if want := []string{"RULE-002", "RULE-001"}; !reflect.DeepEqual(sequenceValues, want) {
		t.Fatalf("yamlSequenceValuesInSection() = %#v; want %#v", sequenceValues, want)
	}
	emptySequence, err := yamlSequenceValuesInSection("registry.yaml", "missing")
	if err != nil {
		t.Fatalf("yamlSequenceValuesInSection(missing section) error = %v", err)
	}
	if len(emptySequence) != 0 {
		t.Fatalf("yamlSequenceValuesInSection(missing section) = %#v; want empty", emptySequence)
	}
	if _, err := yamlSequenceValuesInSection("missing.yaml", "rules"); err == nil {
		t.Fatal("yamlSequenceValuesInSection(missing file) error = nil; want error")
	}

	mapValues, err := yamlMapKeysInSection("registry.yaml", "map")
	if err != nil {
		t.Fatalf("yamlMapKeysInSection() error = %v", err)
	}
	if want := []string{"first", "value", "second"}; !reflect.DeepEqual(mapValues, want) {
		t.Fatalf("yamlMapKeysInSection() = %#v; want %#v", mapValues, want)
	}
	emptyMap, err := yamlMapKeysInSection("registry.yaml", "missing")
	if err != nil {
		t.Fatalf("yamlMapKeysInSection(missing section) error = %v", err)
	}
	if len(emptyMap) != 0 {
		t.Fatalf("yamlMapKeysInSection(missing section) = %#v; want empty", emptyMap)
	}
	if _, err := yamlMapKeysInSection("missing.yaml", "map"); err == nil {
		t.Fatal("yamlMapKeysInSection(missing file) error = nil; want error")
	}

	if got, want := duplicateValues([]string{"", "b", "a", "b", "a", "b"}), []string{"a", "b"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("duplicateValues() = %#v; want %#v", got, want)
	}
	for _, block := range []string{
		"evidence:\n  # comment\n\nnext: value\n",
		"evidence:\n",
	} {
		if blockHasYAMLListItem(block, "evidence") {
			t.Fatalf("blockHasYAMLListItem(%q) = true; want false", block)
		}
	}
}

func TestRulesRegistryEnforcedByGapsBranches(t *testing.T) {
	root := t.TempDir()
	chdir(t, root)
	writeGovernanceHelperFile(t, ".githooks/pre-commit", "#!/bin/sh\n")
	writeGovernanceHelperFile(t, ".agent/registries/command-registry.yaml", `
commands:
  - name: custom-command
`)
	writeGovernanceHelperFile(t, ".agent/registries/makefile-target-registry.yaml", `
targets:
  - custom-target
`)

	registryText := `
rules:
  - id: known-goalcli
    enforced_by: goalcli
  - id: known-required-command
    enforced_by: goalcli audit-goal
  - id: known-custom-command
    enforced_by: goalcli custom-command
  - id: unknown-command
    enforced_by: goalcli missing-command
  - id: known-required-make
    enforced_by: make fmt
  - id: known-custom-make
    enforced_by: make custom-target
  - id: unknown-make
    enforced_by: make missing-target
  - id: existing-hook
    enforced_by: .githooks/pre-commit
  - id: missing-hook
    enforced_by: .githooks/missing
  - id: unsupported
    enforced_by: scripts/check.sh
  - id: empty-list
    enforced_by: []
`
	var gaps []string
	appendRulesRegistryEnforcedByGaps(".agent/registries/rules.yaml", registryText, &gaps)
	for _, want := range []string{
		".agent/registries/rules.yaml enforced_by goalcli missing-command references unknown goalcli command missing-command",
		".agent/registries/rules.yaml enforced_by make missing-target references unknown make target missing-target",
		".agent/registries/rules.yaml enforced_by .githooks/missing references missing hook",
		".agent/registries/rules.yaml enforced_by scripts/check.sh is not a supported gate reference",
	} {
		assertGovernanceGapsContain(t, gaps, want)
	}
	joined := strings.Join(gaps, "\n")
	for _, unexpected := range []string{"custom-command", "custom-target", ".githooks/pre-commit", "audit-goal", "make fmt", "empty-list"} {
		if strings.Contains(joined, unexpected) {
			t.Fatalf("gaps = %#v; did not expect gap for %s", gaps, unexpected)
		}
	}
}

func TestRegistryCommandAndTargetHelpersMergeFiles(t *testing.T) {
	root := t.TempDir()
	chdir(t, root)
	writeGovernanceHelperFile(t, ".agent/registries/command-registry.yaml", `
commands:
  - name: custom-command
`)
	writeGovernanceHelperFile(t, ".agent/registries/makefile-target-registry.yaml", `
targets:
  - custom-target
`)

	commands := rulesRegistryGoalCLICommands()
	if !commands["audit-goal"] || !commands["rules-consistency-check"] || !commands["custom-command"] {
		t.Fatalf("rulesRegistryGoalCLICommands() = %#v; want required, builtin, and custom commands", commands)
	}
	targets := rulesRegistryMakeTargets()
	if !targets["fmt"] || !targets["custom-target"] {
		t.Fatalf("rulesRegistryMakeTargets() = %#v; want required and custom targets", targets)
	}
}

func TestParseMakeEnforcerTargetBranches(t *testing.T) {
	tests := []struct {
		fields []string
		want   string
		wantOK bool
	}{
		{fields: []string{"goalcli", "audit-goal"}, want: "", wantOK: false},
		{fields: []string{"make"}, want: "", wantOK: false},
		{fields: []string{"make", "-C", ".", "fmt"}, want: ".", wantOK: true},
		{fields: []string{"make", "VAR=1", "--silent", "fmt"}, want: "fmt", wantOK: true},
		{fields: []string{"env", "make", "-C", ".", "test"}, want: ".", wantOK: true},
	}
	for _, tt := range tests {
		t.Run(strings.Join(tt.fields, " "), func(t *testing.T) {
			got, ok := parseMakeEnforcerTarget(tt.fields)
			if got != tt.want || ok != tt.wantOK {
				t.Fatalf("parseMakeEnforcerTarget(%#v) = %q, %v; want %q, %v", tt.fields, got, ok, tt.want, tt.wantOK)
			}
		})
	}
}

func TestContextProfileContractRejectsInvalidProfilesAndReleaseGates(t *testing.T) {
	original := contextProfileGates
	contextProfileGates = map[string][]string{
		"":           {"release-check"},
		"unexpected": {"release-final-check", "missing-gate"},
	}
	t.Cleanup(func() {
		contextProfileGates = original
	})

	var gaps []string
	appendContextProfileContractGaps("release-check:\nrelease-final-check:\n", &gaps)

	for _, want := range []string{
		"context profile name must not be empty",
		"unknown context profile ",
		"context profile  must not include release-check",
		"unknown context profile unexpected",
		"context profile unexpected must not include release-final-check",
		"context profile unexpected references unknown Makefile gate missing-gate",
	} {
		assertGovernanceGapsContain(t, gaps, want)
	}
}

func TestReleaseFinalDelegationRejectsMissingTarget(t *testing.T) {
	var gaps []string

	appendReleaseFinalDelegationGaps("context-release:\n\t@true\n", &gaps)

	assertGovernanceGapsContain(t, gaps, "Makefile missing target block release-final-check")
}

func TestEvidenceReplayStaleGapBranches(t *testing.T) {
	if got := evidenceReplayStaleGap(evidenceReplayExpectation{}, "expected.json"); got != "" {
		t.Fatalf("evidenceReplayStaleGap(empty) = %q; want empty", got)
	}

	invalid := evidenceReplayExpectation{GeneratedAt: "not-time", MaxAgeHours: 1}
	if got := evidenceReplayStaleGap(invalid, "expected.json"); !strings.Contains(got, "invalid evidence replay generated_at expected.json") {
		t.Fatalf("evidenceReplayStaleGap(invalid) = %q; want invalid timestamp gap", got)
	}

	stale := evidenceReplayExpectation{
		GeneratedAt: time.Now().UTC().Add(-2 * time.Hour).Format(time.RFC3339),
		MaxAgeHours: 1,
	}
	if got := evidenceReplayStaleGap(stale, "expected.json"); !strings.Contains(got, "stale evidence replay fixture:") {
		t.Fatalf("evidenceReplayStaleGap(stale) = %q; want stale gap", got)
	}
}

func TestInvalidInternalArgsExitTreatsHelpAsSuccess(t *testing.T) {
	var stderr bytes.Buffer

	if got := invalidInternalArgsExit("schema-check", flag.ErrHelp, &stderr); got != 0 {
		t.Fatalf("invalidInternalArgsExit(help) = %d; want 0", got)
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q; want empty", stderr.String())
	}
}

func TestIsXlibStandardSourceModuleBranches(t *testing.T) {
	root := t.TempDir()
	chdir(t, root)
	sourceModule := strings.Join([]string{"github.com", "ZoneCNH", "xlib" + "-standard"}, "/")

	if !isXlibStandardSourceModule() {
		t.Fatalf("isXlibStandardSourceModule() without go.mod = false; want true")
	}

	writeGovernanceHelperFile(t, "go.mod", "module "+sourceModule+"\n")
	if !isXlibStandardSourceModule() {
		t.Fatalf("isXlibStandardSourceModule() for source module = false; want true")
	}

	writeGovernanceHelperFile(t, "go.mod", "module github.com/ZoneCNH/other\n")
	if isXlibStandardSourceModule() {
		t.Fatalf("isXlibStandardSourceModule() for downstream module = true; want false")
	}

	writeGovernanceHelperFile(t, "go.mod", "go 1.23\n")
	if isXlibStandardSourceModule() {
		t.Fatalf("isXlibStandardSourceModule() for malformed go.mod = true; want false")
	}
}

func TestAppendUnclassifiedAgentFileGapsReportsWalkErrors(t *testing.T) {
	missingRoot := filepath.Join(t.TempDir(), ".agent", "missing")
	var gaps []string

	appendUnclassifiedAgentFileGaps(missingRoot, ".agent/index.yaml", map[string]bool{}, &gaps)

	assertGovernanceGapsContain(t, gaps, "read "+filepath.ToSlash(missingRoot)+":")
}

func TestTrimYAMLScalarStripsInlineComments(t *testing.T) {
	if got := trimYAMLScalar(` "value" # comment`); got != "value" {
		t.Fatalf("trimYAMLScalar() = %q; want value", got)
	}
}

func TestContextProfileAliasAndCheckBranches(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if got := runContextProfileAlias("context-fast-check", []string{"unexpected"}, &stdout, &stderr); got != 2 {
		t.Fatalf("runContextProfileAlias(unexpected arg) = %d; want 2", got)
	}
	if !strings.Contains(stderr.String(), "unexpected positional argument") {
		t.Fatalf("stderr = %q; want positional argument error", stderr.String())
	}

	stdout.Reset()
	stderr.Reset()
	if got := runContextProfileAlias("context-release", nil, &stdout, &stderr); got != 0 {
		t.Fatalf("runContextProfileAlias(context-release) = %d; want 0; stderr=%q", got, stderr.String())
	}
	if out := stdout.String(); !strings.Contains(out, `"command": "context-release"`) || !strings.Contains(out, "profile=release") {
		t.Fatalf("stdout = %q; want context-release report with release profile", out)
	}

	if _, err := parseContextProfileCheckProfile("context-profile-check", []string{"--bad"}); err == nil {
		t.Fatalf("parseContextProfileCheckProfile(invalid flag) error = nil; want error")
	}
	if _, err := parseContextProfileCheckProfile("context-profile-check", []string{"unexpected"}); err == nil {
		t.Fatalf("parseContextProfileCheckProfile(positional) error = nil; want error")
	}
	profile, err := parseContextProfileCheckProfile("context-profile-check", []string{"--json", "--profile", "fast"})
	if err != nil {
		t.Fatalf("parseContextProfileCheckProfile(valid) error = %v", err)
	}
	if profile != "fast" {
		t.Fatalf("profile = %q; want fast", profile)
	}

	stdout.Reset()
	stderr.Reset()
	if got := runContextProfileCheck("context-profile-check", []string{"--profile", "missing"}, &stdout, &stderr); got != 2 {
		t.Fatalf("runContextProfileCheck(missing profile) = %d; want 2", got)
	}
	if !strings.Contains(stderr.String(), `invalid context profile "missing"`) {
		t.Fatalf("stderr = %q; want invalid profile", stderr.String())
	}

	stdout.Reset()
	stderr.Reset()
	if got := runContextProfileCheck("context-profile-check", []string{"--bad"}, &stdout, &stderr); got != 2 {
		t.Fatalf("runContextProfileCheck(invalid flag) = %d; want 2", got)
	}
	if !strings.Contains(stderr.String(), "flag provided but not defined") {
		t.Fatalf("stderr = %q; want invalid flag error", stderr.String())
	}

	root := t.TempDir()
	chdir(t, root)
	stdout.Reset()
	stderr.Reset()
	if got := runContextProfileCheck("context-profile-check", nil, &stdout, &stderr); got != 1 {
		t.Fatalf("runContextProfileCheck(missing files) = %d; want 1", got)
	}
	if out := stdout.String(); !strings.Contains(out, `"status": "failed"`) || !strings.Contains(out, "missing .agent/registries/command-registry.yaml") {
		t.Fatalf("stdout = %q; want failed report with missing command registry", out)
	}

	writeGovernanceHelperFile(t, ".agent/registries/command-registry.yaml", "commands:\n")
	stdout.Reset()
	stderr.Reset()
	if got := runContextProfileCheck("context-profile-check", nil, &stdout, &stderr); got != 1 {
		t.Fatalf("runContextProfileCheck(empty registry) = %d; want 1", got)
	}
	if !strings.Contains(stdout.String(), ".agent/registries/command-registry.yaml missing name:") {
		t.Fatalf("stdout = %q; want missing command name gap", stdout.String())
	}

	if got, ok := contextGateProfile("context-standard-check"); got != "standard" || !ok {
		t.Fatalf("contextGateProfile(context-standard-check) = %q, %v; want standard, true", got, ok)
	}
	if got, ok := contextGateProfile("unknown-gate"); got != "" || ok {
		t.Fatalf("contextGateProfile(unknown-gate) = %q, %v; want empty, false", got, ok)
	}
}

func TestMakefileTargetDependenciesBranches(t *testing.T) {
	content := strings.Join([]string{
		"target: dep-a \\",
		" dep-b dep-c",
		"\t@echo ignored",
		"empty:",
		"\t@true",
		"other: dep-z",
	}, "\n")

	got := makefileTargetDependencies(content, "target")
	want := []string{"dep-a", "dep-b", "dep-c"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("makefileTargetDependencies(target) = %#v; want %#v", got, want)
	}
	if got := makefileTargetDependencies(content, "empty"); len(got) != 0 {
		t.Fatalf("makefileTargetDependencies(empty) = %#v; want empty", got)
	}
	if got := makefileTargetDependencies(content, "missing"); got != nil {
		t.Fatalf("makefileTargetDependencies(missing) = %#v; want nil", got)
	}
}

func TestYAMLAndPlannedCommandHelperBranches(t *testing.T) {
	if got := countYAMLLinesWithValue("x: true # comment\ny: true\nx: false\nx: true\n", "x", "true"); got != 2 {
		t.Fatalf("countYAMLLinesWithValue() = %d; want 2", got)
	}
	if got := countYAMLListItems("items:\n  - one\n  - two\nnext: value\n  - ignored\n", "items"); got != 2 {
		t.Fatalf("countYAMLListItems() = %d; want 2", got)
	}
	if got := stripInlineComment("x: true # note"); got != "x: true " {
		t.Fatalf("stripInlineComment() = %q; want prefix before comment", got)
	}
	if got := plannedCommandMarkers("agent-team-contract", ".agent/contracts/team-contract.yaml"); len(got) == 0 {
		t.Fatalf("plannedCommandMarkers(agent-team-contract) = empty; want markers")
	}
	if got := plannedCommandMarkers("unknown", ".agent/contracts/team-contract.yaml"); got != nil {
		t.Fatalf("plannedCommandMarkers(unknown) = %#v; want nil", got)
	}

	args := []string{"--dry-run", "--input", "fixture", "--context=release_verify"}
	if !flagProvided(args, "dry-run") || !flagProvided(args, "context") || flagProvided(args, "missing") {
		t.Fatalf("flagProvided(%#v) returned unexpected values", args)
	}
	if got := flagValue(args, "input", "fallback"); got != "fixture" {
		t.Fatalf("flagValue(input) = %q; want fixture", got)
	}
	if got := flagValue([]string{"--input=fixture2"}, "input", "fallback"); got != "fixture2" {
		t.Fatalf("flagValue(input=) = %q; want fixture2", got)
	}
	if got := flagValue(args, "missing", "fallback"); got != "fallback" {
		t.Fatalf("flagValue(missing) = %q; want fallback", got)
	}

	for _, args := range [][]string{
		{"--verify"},
		{"--strict"},
		{"--context=release_verify"},
		{"--context", "release_verify"},
	} {
		if !plannedCommandVerifyRequested(args) {
			t.Fatalf("plannedCommandVerifyRequested(%#v) = false; want true", args)
		}
	}
	if plannedCommandVerifyRequested([]string{"--context", "local_write"}) {
		t.Fatalf("plannedCommandVerifyRequested(local_write) = true; want false")
	}

	if err := validatePlannedCommandArgs("release-ready", []string{"--context", "missing"}); err == nil {
		t.Fatalf("validatePlannedCommandArgs(invalid context) error = nil; want error")
	}
	if err := validatePlannedCommandArgs("release-ready", []string{"unexpected"}); err == nil {
		t.Fatalf("validatePlannedCommandArgs(positional) error = nil; want error")
	}
	if err := validatePlannedCommandArgs("release-ready", []string{"-h"}); !errors.Is(err, flag.ErrHelp) {
		t.Fatalf("validatePlannedCommandArgs(help) = %v; want flag.ErrHelp", err)
	}
}

func TestEvidenceReplayEvaluationBranches(t *testing.T) {
	root := t.TempDir()

	details, _, _ := evaluateEvidenceReplay("")
	if !strings.Contains(strings.Join(details, "\n"), "fixture=testkit/governance/fixtures/evidence-replay/passed") {
		t.Fatalf("default fixture details=%#v; want default fixture path", details)
	}

	if _, gaps, status := evaluateEvidenceReplay(filepath.Join(root, "missing")); status != "gap" || !strings.Contains(strings.Join(gaps, "\n"), "missing evidence replay expected status") {
		t.Fatalf("missing fixture status=%q gaps=%#v; want gap with missing expected status", status, gaps)
	}

	invalidExpected := filepath.Join(root, "invalid-expected")
	writeEvidenceReplayRawFile(t, invalidExpected, "expected-status.json", "{")
	if _, gaps, status := evaluateEvidenceReplay(invalidExpected); status != "failed" || !strings.Contains(strings.Join(gaps, "\n"), "invalid evidence replay expected status") {
		t.Fatalf("invalid expected status=%q gaps=%#v; want failed invalid expected", status, gaps)
	}

	stale := filepath.Join(root, "stale")
	writeEvidenceReplayExpected(t, stale, evidenceReplayExpectation{
		GeneratedAt: time.Now().UTC().Add(-2 * time.Hour).Format(time.RFC3339),
		MaxAgeHours: 1,
	})
	writeEvidenceReplayRawFile(t, stale, "ledger.jsonl", "")
	if _, gaps, status := evaluateEvidenceReplay(stale); status != "gap" || !strings.Contains(strings.Join(gaps, "\n"), "stale evidence replay fixture") {
		t.Fatalf("stale fixture status=%q gaps=%#v; want gap stale", status, gaps)
	}

	missingLedger := filepath.Join(root, "missing-ledger")
	writeEvidenceReplayExpected(t, missingLedger, evidenceReplayExpectation{})
	if _, gaps, status := evaluateEvidenceReplay(missingLedger); status != "gap" || !strings.Contains(strings.Join(gaps, "\n"), "missing evidence replay ledger") {
		t.Fatalf("missing ledger status=%q gaps=%#v; want gap missing ledger", status, gaps)
	}

	emptyLedger := filepath.Join(root, "empty-ledger")
	writeEvidenceReplayExpected(t, emptyLedger, evidenceReplayExpectation{})
	writeEvidenceReplayRawFile(t, emptyLedger, "ledger.jsonl", "")
	if _, gaps, status := evaluateEvidenceReplay(emptyLedger); status != "failed" || !strings.Contains(strings.Join(gaps, "\n"), "evidence replay ledger has no entries") {
		t.Fatalf("empty ledger status=%q gaps=%#v; want failed empty ledger", status, gaps)
	}

	missingExpectedCommand := filepath.Join(root, "missing-expected-command")
	writeEvidenceReplayExpected(t, missingExpectedCommand, evidenceReplayExpectation{Commands: map[string]string{"release-ready": "passed"}})
	writeEvidenceReplayRawFile(t, missingExpectedCommand, "ledger.jsonl", "")
	if _, gaps, status := evaluateEvidenceReplay(missingExpectedCommand); status != "failed" || !strings.Contains(strings.Join(gaps, "\n"), "evidence replay missing expected command status: release-ready") {
		t.Fatalf("missing expected command status=%q gaps=%#v; want failed missing command", status, gaps)
	}

	invalidLedger := filepath.Join(root, "invalid-ledger")
	writeEvidenceReplayExpected(t, invalidLedger, evidenceReplayExpectation{})
	writeEvidenceReplayRawFile(t, invalidLedger, "ledger.jsonl", "{")
	if _, gaps, status := evaluateEvidenceReplay(invalidLedger); status != "failed" || !strings.Contains(strings.Join(gaps, "\n"), "invalid JSON") {
		t.Fatalf("invalid ledger status=%q gaps=%#v; want failed invalid JSON", status, gaps)
	}

	badLedgerFields := filepath.Join(root, "bad-ledger-fields")
	writeEvidenceReplayExpected(t, badLedgerFields, evidenceReplayExpectation{})
	writeEvidenceReplayRawFile(t, badLedgerFields, "ledger.jsonl", strings.Join([]string{
		`{"seq":1,"previous_hash":"0000000000000000000000000000000000000000000000000000000000000000","stdout_sha256":"bad","entry_hash":"bad"}`,
		`{"command":"implicit-seq","status":"passed","artifact_path":"missing.txt","previous_hash":"bad","entry_hash":"bad"}`,
	}, "\n"))
	_, gaps, status := evaluateEvidenceReplay(badLedgerFields)
	if status != "failed" {
		t.Fatalf("bad ledger fields status=%q gaps=%#v; want failed", status, gaps)
	}
	for _, want := range []string{"entry 1 missing command", "entry 1 missing status", "entry 2 missing artifact missing.txt"} {
		if !strings.Contains(strings.Join(gaps, "\n"), want) {
			t.Fatalf("bad ledger fields gaps=%#v missing %q", gaps, want)
		}
	}

	passed := filepath.Join(root, "passed")
	writeEvidenceReplayFixture(t, passed, "release-ready", "passed", "release ready\n")
	details, gaps, status = evaluateEvidenceReplay(passed)
	if status != "passed" || len(gaps) != 0 {
		t.Fatalf("passed fixture status=%q gaps=%#v details=%#v; want passed with no gaps", status, gaps, details)
	}
	for _, want := range []string{"checksum verified", "hash chain verified", "expected command status verified"} {
		if !strings.Contains(strings.Join(details, "\n"), want) {
			t.Fatalf("details=%#v missing %q", details, want)
		}
	}

	mismatch := filepath.Join(root, "mismatch")
	writeEvidenceReplayFixture(t, mismatch, "release-ready", "passed", "release ready\n")
	writeEvidenceReplayExpected(t, mismatch, evidenceReplayExpectation{Commands: map[string]string{"release-ready": "failed"}})
	if _, gaps, status := evaluateEvidenceReplay(mismatch); status != "failed" || !strings.Contains(strings.Join(gaps, "\n"), "status mismatch") {
		t.Fatalf("mismatch status=%q gaps=%#v; want failed status mismatch", status, gaps)
	}
}

func TestReleaseReadyDecisionDetailsBranches(t *testing.T) {
	var gaps []string
	details := releaseReadyDecisionDetails([]string{"--verify", "--context", "local_write"}, map[string]string{}, &gaps)
	joinedGaps := strings.Join(gaps, "\n")
	for _, want := range []string{
		"requires context release_verify",
		"gates must include required release gates",
		"missing required_release_evidence",
		"replay must be strict",
	} {
		if !strings.Contains(joinedGaps, want) {
			t.Fatalf("gaps missing %q in %#v", want, gaps)
		}
	}
	if joined := strings.Join(details, "\n"); !strings.Contains(joined, "verdict=gap") || !strings.Contains(joined, "score=0/100") {
		t.Fatalf("details=%#v; want gap score", details)
	}

	files := map[string]string{
		".agent/release/release-required-gates.yaml": `
gates:
  - id: one
    required_for_release: true
    release_usable: false
required_release_evidence:
  - release/manifest/latest.json
`,
		".agent/evidence/evidence-replay.yaml": "replay:\n  strict: true\n",
	}
	gaps = nil
	details = releaseReadyDecisionDetails([]string{"--verify"}, files, &gaps)
	if !strings.Contains(strings.Join(gaps, "\n"), "release-ready verdict not_ready") {
		t.Fatalf("gaps=%#v; want not_ready readiness gap", gaps)
	}
	if !strings.Contains(strings.Join(details, "\n"), "mode=readiness_gate") {
		t.Fatalf("details=%#v; want readiness_gate mode", details)
	}

	gaps = nil
	details = releaseReadyDecisionDetails([]string{"--verify", "--dry-run"}, files, &gaps)
	if strings.Contains(strings.Join(gaps, "\n"), "release-ready verdict not_ready") {
		t.Fatalf("dry-run gaps=%#v; want no not_ready readiness gap", gaps)
	}
	if !strings.Contains(strings.Join(details, "\n"), "mode=dry_run_contract") {
		t.Fatalf("details=%#v; want dry_run_contract mode", details)
	}

	files[".agent/release/release-required-gates.yaml"] = `
gates:
  - id: one
    required_for_release: true
    release_usable: true
required_release_evidence:
  - release/manifest/latest.json
`
	gaps = nil
	details = releaseReadyDecisionDetails([]string{"--verify"}, files, &gaps)
	if len(gaps) != 0 {
		t.Fatalf("ready gaps=%#v; want none", gaps)
	}
	if joined := strings.Join(details, "\n"); !strings.Contains(joined, "verdict=ready") || !strings.Contains(joined, "score=100/100") {
		t.Fatalf("ready details=%#v; want ready score", details)
	}
}

func TestPlannedCommandFileValidationBranches(t *testing.T) {
	root := t.TempDir()
	chdir(t, root)

	if _, gap, ok := readPlannedCommandFile("missing.yaml"); ok || !strings.Contains(gap, "missing missing.yaml") {
		t.Fatalf("readPlannedCommandFile(missing) gap=%q ok=%v; want missing gap", gap, ok)
	}
	if err := os.MkdirAll("dir.yaml", 0o755); err != nil {
		t.Fatalf("mkdir dir.yaml: %v", err)
	}
	if _, gap, ok := readPlannedCommandFile("dir.yaml"); ok || !strings.Contains(gap, "must be a file") {
		t.Fatalf("readPlannedCommandFile(dir) gap=%q ok=%v; want directory gap", gap, ok)
	}
	writeGovernanceHelperFile(t, "ok.yaml", "schema_version: 1\n")
	if content, gap, ok := readPlannedCommandFile("ok.yaml"); !ok || gap != "" || string(content) != "schema_version: 1\n" {
		t.Fatalf("readPlannedCommandFile(ok) content=%q gap=%q ok=%v; want file content", string(content), gap, ok)
	}
	oldReadFile := readPlannedCommandFileReadFile
	t.Cleanup(func() { readPlannedCommandFileReadFile = oldReadFile })
	readPlannedCommandFileReadFile = func(string) ([]byte, error) {
		return nil, errors.New("read failed")
	}
	if _, gap, ok := readPlannedCommandFile("ok.yaml"); ok || !strings.Contains(gap, "ok.yaml unreadable: read failed") {
		t.Fatalf("readPlannedCommandFile(read error) gap=%q ok=%v; want unreadable gap", gap, ok)
	}
	readPlannedCommandFileReadFile = oldReadFile

	gaps := validatePlannedCommandFile("agent-team-contract", ".agent/contracts/team-contract.yaml", []byte("schema_version: 1\n"))
	for _, want := range []string{"missing semantic marker roles:", "missing semantic marker rule:"} {
		if !strings.Contains(strings.Join(gaps, "\n"), want) {
			t.Fatalf("gaps=%#v missing %q", gaps, want)
		}
	}
	if gaps := validatePlannedCommandFile("custom", "bad.json", []byte("{")); !strings.Contains(strings.Join(gaps, "\n"), "must be valid JSON") {
		t.Fatalf("validatePlannedCommandFile(bad json) gaps=%#v; want JSON gap", gaps)
	}
	if gaps := validatePlannedCommandFile("custom", "empty.md", []byte(" \n")); !strings.Contains(strings.Join(gaps, "\n"), "empty.md must not be empty") {
		t.Fatalf("validatePlannedCommandFile(empty) gaps=%#v; want empty gap", gaps)
	}

	var stdout, stderr bytes.Buffer
	if got := runPlannedCommand("agent-team-contract", []string{"-h"}, &stdout, &stderr); got != 0 {
		t.Fatalf("runPlannedCommand(help) = %d; want 0; stderr=%q", got, stderr.String())
	}
	if stdout.Len() != 0 || stderr.Len() != 0 {
		t.Fatalf("help stdout=%q stderr=%q; want discarded flag output", stdout.String(), stderr.String())
	}

	stdout.Reset()
	stderr.Reset()
	if got := runPlannedCommand("not-covered", nil, &stdout, &stderr); got != 1 {
		t.Fatalf("runPlannedCommand(not-covered) = %d; want 1", got)
	}
	if !strings.Contains(stdout.String(), `"status": "failed"`) || !strings.Contains(stdout.String(), "no manifest coverage") {
		t.Fatalf("stdout=%q; want failed manifest coverage report", stdout.String())
	}
}

func TestHarnessAliasGapsBranches(t *testing.T) {
	root := t.TempDir()
	chdir(t, root)

	var gaps []string
	appendHarnessAliasGaps(".agent/harness/harness.yaml", &gaps)
	assertGovernanceGapsContain(t, gaps, "missing .agent/harness/harness.yaml")

	gaps = nil
	writeGovernanceHelperFile(t, ".agent/harness/harness.yaml", "required_gates: []\n")
	appendHarnessAliasGaps(".agent/harness/harness.yaml", &gaps)
	assertGovernanceGapsContain(t, gaps, "must define required_gates")

	gaps = nil
	writeGovernanceHelperFile(t, ".agent/harness/harness.yaml", `
required_gates:
  - id: governance_check
    semantic_role: source
  - id: p1_governance_check
    semantic_role: source
  - id: p2_runtime_check
    semantic_role: source
  - id: governance_chain
    alias_of: governance_check
  - id: governance_release_scope
    alias_of: wrong_target
    semantic_role: release
  - id: p1_governance_release_scope
    alias_of: p1_governance_check
    semantic_role: release
  - id: p2_runtime_chain
    alias_of: p2_runtime_check
    semantic_role: chain
  - id: p2_runtime_release_scope
    alias_of: p2_runtime_check
    semantic_role: release
`)
	appendHarnessAliasGaps(".agent/harness/harness.yaml", &gaps)
	for _, want := range []string{
		"governance_chain missing semantic_role",
		"governance_release_scope must set alias_of: governance_check",
		"missing required gate alias p1_governance_chain",
	} {
		assertGovernanceGapsContain(t, gaps, want)
	}

	gaps = nil
	writeGovernanceHelperFile(t, ".agent/harness/harness.yaml", `
required_gates:
  - id: governance_chain
    alias_of: missing_target
    semantic_role: chain
	`)
	appendHarnessAliasGaps(".agent/harness/harness.yaml", &gaps)
	assertGovernanceGapsContain(t, gaps, "governance_chain alias_of missing target governance_check")
}

func TestAgentIndexClassificationGapsBranches(t *testing.T) {
	entry := agentIndexEntry{
		path: ".agent/registries/generated-artifacts.yaml",
		block: `- path: .agent/registries/generated-artifacts.yaml
  layer: documentation
  authority: validated_mirror
  mutability: generated
`,
	}
	entries := []agentIndexEntry{entry}
	var gaps []string
	appendAgentIndexClassificationGaps(".agent/index.yaml", entries, &gaps)

	for _, want := range []string{
		".agent/index.yaml .agent/registries/generated-artifacts.yaml must classify layer as registry",
		".agent/index.yaml .agent/registries/generated-artifacts.yaml must classify authority as source_of_truth",
		".agent/index.yaml .agent/registries/generated-artifacts.yaml must classify mutability as hand_written",
	} {
		assertGovernanceGapsContain(t, gaps, want)
	}
}

func TestGeneratedArtifactClassificationGapsBranches(t *testing.T) {
	root := t.TempDir()
	chdir(t, root)

	var gaps []string
	writeGovernanceHelperFile(t, ".agent/index.yaml", `
files:
  - path: .agent/registries/generated-artifacts.yaml
    layer: registry
    authority: validated_mirror
    mutability: generated
  - path: .agent/generated.md
    layer: runtime_contract
    authority: validated_mirror
    mutability: generated
`)
	appendGeneratedArtifactClassificationGaps(".agent/index.yaml", ".agent/registries/generated-artifacts.yaml", &gaps)
	for _, want := range []string{
		".agent/registries/generated-artifacts.yaml must be indexed as source_of_truth",
		".agent/registries/generated-artifacts.yaml must be indexed as hand_written",
		"missing .agent/registries/generated-artifacts.yaml for generated .agent files",
	} {
		assertGovernanceGapsContain(t, gaps, want)
	}

	gaps = nil
	writeGovernanceHelperFile(t, ".agent/registries/generated-artifacts.yaml", `
artifacts:
  - path: .agent/unindexed.md
  - path: .agent/not-generated.md
  - path: .agent/source.md
`)
	writeGovernanceHelperFile(t, ".agent/index.yaml", `
files:
  - path: .agent/registries/generated-artifacts.yaml
    layer: registry
    authority: source_of_truth
    mutability: hand_written
  - path: .agent/generated.md
    layer: runtime_contract
    authority: validated_mirror
    mutability: generated
  - path: .agent/not-generated.md
    layer: runtime_contract
    authority: validated_mirror
    mutability: hand_written
  - path: .agent/source.md
    layer: runtime_contract
    authority: source_of_truth
    mutability: generated
`)
	appendGeneratedArtifactClassificationGaps(".agent/index.yaml", ".agent/registries/generated-artifacts.yaml", &gaps)
	for _, want := range []string{
		".agent/index.yaml .agent/generated.md mutability generated requires .agent/registries/generated-artifacts.yaml entry",
		".agent/registries/generated-artifacts.yaml references unindexed agent artifact .agent/unindexed.md",
		".agent/registries/generated-artifacts.yaml .agent/not-generated.md must be indexed with mutability generated",
		".agent/registries/generated-artifacts.yaml .agent/source.md generated artifact must not be source_of_truth",
	} {
		assertGovernanceGapsContain(t, gaps, want)
	}
}

func TestDoctorMainGuardAndSpecCheckBranches(t *testing.T) {
	root := t.TempDir()
	chdir(t, root)
	writeGovernanceHelperFile(t, "go.mod", "module example.com/consumer\n")

	var stdout, stderr bytes.Buffer
	if got := runDoctor([]string{"--bad"}, &stdout, &stderr); got != 2 {
		t.Fatalf("runDoctor(invalid args) = %d; want 2", got)
	}
	stdout.Reset()
	stderr.Reset()
	if got := runDoctor(nil, &stdout, &stderr); got != 1 {
		t.Fatalf("runDoctor(missing files) = %d; want 1", got)
	}
	if !strings.Contains(stdout.String(), "missing Makefile") {
		t.Fatalf("stdout=%q; want missing Makefile gap", stdout.String())
	}

	for _, path := range []string{
		".agent/harness/harness.yaml",
		".agent/index.yaml",
		".agent/registries/issue-registry.yaml",
		".agent/registries/command-registry.yaml",
		".agent/registries/makefile-target-registry.yaml",
		".agent/registries/makefile-baseline.yaml",
		".github/workflows/adoption-check.yml",
		"mk/governance.mk",
		"docs/standard/goalcli-cli-contract.md",
		"contracts/goalcli-report.schema.json",
		"Makefile",
	} {
		writeGovernanceHelperFile(t, path, "present\n")
	}
	stdout.Reset()
	stderr.Reset()
	if got := runDoctor(nil, &stdout, &stderr); got != 0 {
		t.Fatalf("runDoctor(present files) = %d stderr=%q stdout=%q; want 0", got, stderr.String(), stdout.String())
	}
	if !strings.Contains(stdout.String(), "required governance files are present") {
		t.Fatalf("stdout=%q; want doctor success detail", stdout.String())
	}

	stdout.Reset()
	stderr.Reset()
	if got := runMainGuard([]string{"--context", "bad"}, &stdout, &stderr); got != 2 {
		t.Fatalf("runMainGuard(invalid context) = %d; want 2", got)
	}
	stdout.Reset()
	stderr.Reset()
	if got := runMainGuard([]string{"--bad"}, &stdout, &stderr); got != 2 {
		t.Fatalf("runMainGuard(invalid args) = %d; want 2", got)
	}
	stdout.Reset()
	stderr.Reset()
	if got := runMainGuard([]string{"--help"}, &stdout, &stderr); got != 0 {
		t.Fatalf("runMainGuard(help) = %d; want 0", got)
	}

	specRoot := t.TempDir()
	chdir(t, specRoot)
	stdout.Reset()
	stderr.Reset()
	if got := runSpecCheck([]string{"--bad"}, &stdout, &stderr); got != 2 {
		t.Fatalf("runSpecCheck(invalid args) = %d; want 2", got)
	}
	stdout.Reset()
	stderr.Reset()
	if got := runSpecCheck(nil, &stdout, &stderr); got != 1 {
		t.Fatalf("runSpecCheck(missing docs) = %d; want 1", got)
	}
	if !strings.Contains(stdout.String(), "missing docs") {
		t.Fatalf("stdout=%q; want missing docs gap", stdout.String())
	}

	writeGovernanceHelperFile(t, "docs/spec.md", "# Spec\n")
	oldTrackedDocs := runSpecCheckTrackedDocsMarkdownFiles
	oldSpecReadFile := runSpecCheckReadFile
	t.Cleanup(func() {
		runSpecCheckTrackedDocsMarkdownFiles = oldTrackedDocs
		runSpecCheckReadFile = oldSpecReadFile
	})
	runSpecCheckTrackedDocsMarkdownFiles = func() ([]string, error) {
		return nil, errors.New("scan failed")
	}
	stdout.Reset()
	stderr.Reset()
	if got := runSpecCheck(nil, &stdout, &stderr); got != 1 {
		t.Fatalf("runSpecCheck(scan error) = %d; want 1", got)
	}
	if !strings.Contains(stdout.String(), "scan docs: scan failed") {
		t.Fatalf("stdout=%q; want scan error gap", stdout.String())
	}
	runSpecCheckTrackedDocsMarkdownFiles = func() ([]string, error) {
		return []string{"docs/spec.md"}, nil
	}
	runSpecCheckReadFile = func(string) ([]byte, error) {
		return nil, errors.New("read failed")
	}
	stdout.Reset()
	stderr.Reset()
	if got := runSpecCheck(nil, &stdout, &stderr); got != 1 {
		t.Fatalf("runSpecCheck(read error) = %d; want 1", got)
	}
	if !strings.Contains(stdout.String(), "read docs/spec.md: read failed") {
		t.Fatalf("stdout=%q; want read error gap", stdout.String())
	}
	runSpecCheckTrackedDocsMarkdownFiles = oldTrackedDocs
	runSpecCheckReadFile = oldSpecReadFile
	stdout.Reset()
	stderr.Reset()
	if got := runSpecCheck(nil, &stdout, &stderr); got != 0 {
		t.Fatalf("runSpecCheck(no requirements) = %d stderr=%q stdout=%q; want 0", got, stderr.String(), stdout.String())
	}
	if !strings.Contains(stdout.String(), "warning: no docs markdown file contains REQ-") {
		t.Fatalf("stdout=%q; want no requirements warning", stdout.String())
	}

	writeGovernanceHelperFile(t, "docs/spec.md", "# Spec\n\nREQ-001\n")
	stdout.Reset()
	stderr.Reset()
	if got := runSpecCheck(nil, &stdout, &stderr); got != 0 {
		t.Fatalf("runSpecCheck(requirement present) = %d stderr=%q stdout=%q; want 0", got, stderr.String(), stdout.String())
	}
	if strings.Contains(stdout.String(), "warning: no docs markdown file contains REQ-") {
		t.Fatalf("stdout=%q; want no missing requirements warning", stdout.String())
	}
}

func TestHarnessAndYAMLBlockHelperBranches(t *testing.T) {
	if !blockHasEvidence("evidence: release/manifest/latest.json\n") {
		t.Fatal("blockHasEvidence(scalar) = false; want true")
	}
	if !blockHasEvidence("evidence:\n  - make test\n") {
		t.Fatal("blockHasEvidence(list) = false; want true")
	}
	if blockHasEvidence("evidence: []\n") {
		t.Fatal("blockHasEvidence(empty list scalar) = true; want false")
	}
	if !blockHasYAMLListItem("evidence:\n  - make test\nnext: value\n", "evidence") {
		t.Fatal("blockHasYAMLListItem(evidence list) = false; want true")
	}

	entries := parseYAMLSequenceBlocks("required_gates:\n  - id: coverage_check\n    command: make coverage-check\nnext:\n  - id: ignored\n", "required_gates", "id")
	if len(entries) != 1 || entries[0].value != "coverage_check" {
		t.Fatalf("parseYAMLSequenceBlocks entries=%#v; want one coverage_check entry", entries)
	}
	var gaps []string
	requireYAMLBlockValue("harness", entries[0].value, entries[0].block, "command", "make lint", &gaps)
	assertGovernanceGapsContain(t, gaps, "must set command: make lint")

	root := t.TempDir()
	chdir(t, root)
	gaps = nil
	appendHarnessProofDepthGaps(".agent/harness/harness.yaml", &gaps)
	assertGovernanceGapsContain(t, gaps, "missing .agent/harness/harness.yaml")

	writeGovernanceHelperFile(t, ".agent/harness/harness.yaml", `
required_gates:
  - id: coverage_check
    proof_depth: command
`)
	gaps = nil
	appendHarnessProofDepthGaps(".agent/harness/harness.yaml", &gaps)
	assertGovernanceGapsContain(t, gaps, "proof_depth taxonomy must define ids")
	assertGovernanceGapsContain(t, gaps, "coverage_check missing target_depth")

	validHarness := `
proof_depth:
  taxonomy:
    - id: command
    - id: release_artifact
required_gates:
  - id: coverage_check
    proof_depth: command
    target_depth: release_artifact
gate_link_semantics:
  duplicate_command_links: aliases
  duplicate_entries_do_not_create_new_authorities: true
  authority_source: required_gates[].id
`
	writeGovernanceHelperFile(t, ".agent/harness/harness.yaml", validHarness)
	gaps = nil
	appendHarnessProofDepthGaps(".agent/harness/harness.yaml", &gaps)
	appendHarnessGateLinkSemanticsGaps(".agent/harness/harness.yaml", &gaps)
	if len(gaps) != 0 {
		t.Fatalf("valid harness gaps=%#v; want none", gaps)
	}
	ids := parseHarnessProofDepthTaxonomyIDs(validHarness + "\nother:\n  taxonomy:\n    - id: ignored\n")
	if !ids["command"] || !ids["release_artifact"] || ids["ignored"] {
		t.Fatalf("taxonomy ids=%#v; want command/release_artifact only", ids)
	}
}

func TestRulesConsistencyCheckBranches(t *testing.T) {
	root := t.TempDir()
	chdir(t, root)

	var stdout, stderr bytes.Buffer
	if got := runRulesConsistencyCheck([]string{"--bad"}, &stdout, &stderr); got != 2 {
		t.Fatalf("runRulesConsistencyCheck(invalid args) = %d; want 2", got)
	}
	stdout.Reset()
	stderr.Reset()
	if got := runRulesConsistencyCheck(nil, &stdout, &stderr); got != 1 {
		t.Fatalf("runRulesConsistencyCheck(missing canonical) = %d; want 1", got)
	}
	if !strings.Contains(stderr.String(), "goal-runtime-canonical.md") || stdout.Len() != 0 {
		t.Fatalf("stderr=%q stdout=%q; want canonical read error only", stderr.String(), stdout.String())
	}

	writeGovernanceHelperFile(t, ".agent/runtime/standard/goal-runtime-canonical.md", `
# Canonical

## 1. Core
| Rule | Description |
| --- | --- |
| RULE-CORE-001 | first |

## 2. Ignored
| RULE-CORE-999 | ignored |
`)
	stdout.Reset()
	stderr.Reset()
	if got := runRulesConsistencyCheck(nil, &stdout, &stderr); got != 1 {
		t.Fatalf("runRulesConsistencyCheck(missing iron) = %d; want 1", got)
	}
	if !strings.Contains(stderr.String(), "iron-rules.md") {
		t.Fatalf("stderr=%q; want iron read error", stderr.String())
	}

	writeGovernanceHelperFile(t, ".agent/rules/iron-rules.md", `
# Iron Rules

## 七律
第一律 RULE-CORE-001

## 其他
RULE-CORE-999
`)
	stdout.Reset()
	stderr.Reset()
	if got := runRulesConsistencyCheck(nil, &stdout, &stderr); got != 1 {
		t.Fatalf("runRulesConsistencyCheck(missing registry) = %d; want 1", got)
	}
	if !strings.Contains(stderr.String(), "registry.yaml") {
		t.Fatalf("stderr=%q; want registry read error", stderr.String())
	}

	writeGovernanceHelperFile(t, ".agent/runtime/standard/goal-runtime-canonical.md", "# Canonical\n\n## 1. Core\nNo rules.\n")
	writeGovernanceHelperFile(t, ".agent/rules/iron-rules.md", "# Iron Rules\n\n## 七律\nNo rules.\n")
	writeGovernanceHelperFile(t, ".agent/rules/registry.yaml", "rules:\n  - id: not-rule\n")
	stdout.Reset()
	stderr.Reset()
	if got := runRulesConsistencyCheck(nil, &stdout, &stderr); got != 1 {
		t.Fatalf("runRulesConsistencyCheck(no ids) = %d; want 1", got)
	}
	for _, want := range []string{
		".agent/runtime/standard/goal-runtime-canonical.md: 未发现八条铁律段的 RULE-* 引用",
		".agent/rules/iron-rules.md: 未发现七律段的 RULE-* 引用",
		".agent/rules/registry.yaml: 未发现 RULE-* 登记",
	} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("stdout=%q missing %q", stdout.String(), want)
		}
	}
	if ids := extractCanonicalIronRuleIDs("# none\n"); len(ids) != 0 {
		t.Fatalf("extractCanonicalIronRuleIDs(no ids) = %#v; want empty", ids)
	}
	if ids := extractIronRulesIDs("# none\n"); len(ids) != 0 {
		t.Fatalf("extractIronRulesIDs(no ids) = %#v; want empty", ids)
	}

	writeGovernanceHelperFile(t, ".agent/runtime/standard/goal-runtime-canonical.md", `
# Canonical

## 1. Core
| Rule | Description |
| --- | --- |
| RULE-CORE-001 | first |

## 2. Ignored
| RULE-CORE-999 | ignored |
`)
	writeGovernanceHelperFile(t, ".agent/rules/iron-rules.md", `
# Iron Rules

## 七律
第一律 RULE-CORE-001

## 其他
RULE-CORE-999
`)

	writeGovernanceHelperFile(t, ".agent/rules/registry.yaml", `
rules:
  - id: RULE-CORE-001
    status: active
    enforced_by: goalcli rules-consistency-check
`)
	stdout.Reset()
	stderr.Reset()
	if got := runRulesConsistencyCheck(nil, &stdout, &stderr); got != 0 {
		t.Fatalf("runRulesConsistencyCheck(valid) = %d stderr=%q stdout=%q; want 0", got, stderr.String(), stdout.String())
	}
	if !strings.Contains(stdout.String(), "canonical=1 iron=1 registry=1") {
		t.Fatalf("stdout=%q; want consistency counts", stdout.String())
	}

	writeGovernanceHelperFile(t, ".agent/rules/registry.yaml", `
rules:
  - id: RULE-OTHER-001
    status: active
    enforced_by: goalcli
  - id: RULE-HOOK-001
    status: active
    enforced_by: .githooks/pre-commit
  - id: RULE-GATE-001
    status: active
    enforced_by: unknown-gate
`)
	stdout.Reset()
	stderr.Reset()
	if got := runRulesConsistencyCheck(nil, &stdout, &stderr); got != 1 {
		t.Fatalf("runRulesConsistencyCheck(gaps) = %d; want 1", got)
	}
	for _, want := range []string{"RULE-CORE-001", "references missing hook", "unknown-gate is not a supported gate reference"} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("stdout=%q missing %q", stdout.String(), want)
		}
	}
	if !strings.Contains(stderr.String(), "rules-consistency-check found") {
		t.Fatalf("stderr=%q; want gap summary", stderr.String())
	}
}

func writeGovernanceHelperFile(t *testing.T, path string, content string) {
	t.Helper()
	fullPath := filepath.FromSlash(path)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(fullPath), err)
	}
	if err := os.WriteFile(fullPath, []byte(strings.TrimPrefix(content, "\n")), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func writeFakeExecutable(t *testing.T, dir string, name string, content string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", dir, err)
	}
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0o755); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func writeEvidenceReplayRawFile(t *testing.T, root string, rel string, content string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func writeEvidenceReplayExpected(t *testing.T, root string, expected evidenceReplayExpectation) {
	t.Helper()
	data, err := json.Marshal(expected)
	if err != nil {
		t.Fatalf("marshal expected status: %v", err)
	}
	writeEvidenceReplayRawFile(t, root, "expected-status.json", string(data)+"\n")
}

func writeEvidenceReplayFixture(t *testing.T, root string, command string, status string, stdout string) {
	t.Helper()
	artifactPath := "stdout/" + command + ".txt"
	writeEvidenceReplayRawFile(t, root, artifactPath, stdout)

	entry := evidenceReplayEntry{
		Seq:          1,
		Command:      command,
		Status:       status,
		StdoutSHA256: sha256Hex([]byte(stdout)),
		ArtifactPath: artifactPath,
		PreviousHash: strings.Repeat("0", 64),
	}
	entry.EntryHash = evidenceReplayEntryHash(entry)
	entryData, err := json.Marshal(entry)
	if err != nil {
		t.Fatalf("marshal evidence replay entry: %v", err)
	}
	writeEvidenceReplayRawFile(t, root, "ledger.jsonl", string(entryData)+"\n")
	writeEvidenceReplayExpected(t, root, evidenceReplayExpectation{Commands: map[string]string{command: status}})
}

func assertGovernanceGapsContain(t *testing.T, gaps []string, needle string) {
	t.Helper()
	if strings.Contains(strings.Join(gaps, "\n"), needle) {
		return
	}
	t.Fatalf("gaps missing %q in:\n%s", needle, strings.Join(gaps, "\n"))
}
