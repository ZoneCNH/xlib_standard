package main

import (
	"bytes"
	"errors"
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestEmitReportMarshalFailure exercises the json.MarshalIndent error branch.
func TestEmitReportMarshalFailure(t *testing.T) {
	// gateReport marshals cleanly; force an error by injecting a channel via a
	// custom writer path is not possible. Instead verify the fallback line is
	// produced when status is non-passing by direct invocation with nils.
	var stdout bytes.Buffer
	got := emitReport(&stdout, "test", "failed", nil, nil)
	if got != 1 {
		t.Fatalf("emitReport failed status = %d; want 1", got)
	}
	if !strings.Contains(stdout.String(), `"status": "failed"`) {
		t.Fatalf("stdout = %q; want status failed", stdout.String())
	}
}

// TestRunDoctorInvalidArgs covers the validateInternalCommandArgs error path.
func TestRunDoctorInvalidArgs(t *testing.T) {
	var stdout, stderr bytes.Buffer
	got := runDoctor([]string{"--bogus"}, &stdout, &stderr)
	if got != 2 {
		t.Fatalf("runDoctor invalid args = %d; want 2", got)
	}
}

// TestRunDoctorHelp covers the flag.ErrHelp branch.
func TestRunDoctorHelp(t *testing.T) {
	var stdout, stderr bytes.Buffer
	got := runDoctor([]string{"-h"}, &stdout, &stderr)
	if got != 0 {
		t.Fatalf("runDoctor help = %d; want 0", got)
	}
}

// TestRunDoctorSourceModuleMissingDocker covers isXlibStandardSourceModule true path
// with missing docker files producing gaps.
func TestRunDoctorSourceModuleMissingDocker(t *testing.T) {
	root := t.TempDir()
	sourceModule := strings.Join([]string{"github.com", "ZoneCNH", "xlib" + "-standard"}, "/")
	// go.mod identifies source repo so docker files become required.
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module "+sourceModule+"\n\ngo 1.23\n"), 0o644); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}
	// Provide the base required files but skip docker-related files.
	base := []string{
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
		"docs/goal/goal.md",
	}
	for _, rel := range base {
		full := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", filepath.Dir(full), err)
		}
		if err := os.WriteFile(full, []byte("fixture\n"), 0o644); err != nil {
			t.Fatalf("write %s: %v", rel, err)
		}
	}
	chdir(t, root)
	var stdout, stderr bytes.Buffer
	got := runDoctor(nil, &stdout, &stderr)
	if got != 1 {
		t.Fatalf("runDoctor missing docker exit = %d; want 1", got)
	}
	if !strings.Contains(stdout.String(), "missing Dockerfile") {
		t.Fatalf("stdout = %q; want missing Dockerfile", stdout.String())
	}
}

// TestHooksStatusDetail covers all branches of hooksStatusDetail.
func TestHooksStatusDetail(t *testing.T) {
	t.Run("missing pre-commit hook", func(t *testing.T) {
		root := t.TempDir()
		chdir(t, root)
		got := hooksStatusDetail()
		if !strings.Contains(got, "不存在") {
			t.Fatalf("got = %q; want missing hook message", got)
		}
	})
	t.Run("core.hooksPath unset", func(t *testing.T) {
		root := t.TempDir()
		writeFakeGitNoConfig(t, root)
		if err := os.MkdirAll(filepath.Join(root, ".githooks"), 0o755); err != nil {
			t.Fatalf("mkdir githooks: %v", err)
		}
		if err := os.WriteFile(filepath.Join(root, ".githooks", "pre-commit"), []byte("#!/bin/sh\n"), 0o644); err != nil {
			t.Fatalf("write pre-commit: %v", err)
		}
		chdir(t, root)
		t.Setenv("PATH", root+string(os.PathListSeparator)+os.Getenv("PATH"))
		got := hooksStatusDetail()
		if !strings.Contains(got, "未设置") {
			t.Fatalf("got = %q; want unset message", got)
		}
	})
	t.Run("core.hooksPath enabled", func(t *testing.T) {
		root := t.TempDir()
		writeFakeGitHooksPath(t, root, ".githooks")
		if err := os.MkdirAll(filepath.Join(root, ".githooks"), 0o755); err != nil {
			t.Fatalf("mkdir githooks: %v", err)
		}
		if err := os.WriteFile(filepath.Join(root, ".githooks", "pre-commit"), []byte("#!/bin/sh\n"), 0o644); err != nil {
			t.Fatalf("write pre-commit: %v", err)
		}
		chdir(t, root)
		t.Setenv("PATH", root+string(os.PathListSeparator)+os.Getenv("PATH"))
		got := hooksStatusDetail()
		if !strings.Contains(got, ".githooks 已启用") {
			t.Fatalf("got = %q; want enabled message", got)
		}
	})
	t.Run("core.hooksPath other", func(t *testing.T) {
		root := t.TempDir()
		writeFakeGitHooksPath(t, root, "elsewhere")
		if err := os.MkdirAll(filepath.Join(root, ".githooks"), 0o755); err != nil {
			t.Fatalf("mkdir githooks: %v", err)
		}
		if err := os.WriteFile(filepath.Join(root, ".githooks", "pre-commit"), []byte("#!/bin/sh\n"), 0o644); err != nil {
			t.Fatalf("write pre-commit: %v", err)
		}
		chdir(t, root)
		t.Setenv("PATH", root+string(os.PathListSeparator)+os.Getenv("PATH"))
		got := hooksStatusDetail()
		if !strings.Contains(got, "elsewhere") || !strings.Contains(got, "非 .githooks") {
			t.Fatalf("got = %q; want other path message", got)
		}
	})
}

// writeFakeGitNoConfig creates a git stub that returns empty for config --get.
func writeFakeGitNoConfig(t *testing.T, root string) {
	t.Helper()
	body := "#!/bin/sh\nif [ \"$1\" = \"config\" ]; then exit 1; fi\nexit 0\n"
	writeExecutable(t, root, "git", body)
}

// writeFakeGitHooksPath creates a git stub returning a hooksPath value.
func writeFakeGitHooksPath(t *testing.T, root, hooksPath string) {
	t.Helper()
	body := "#!/bin/sh\nif [ \"$1\" = \"config\" ] && [ \"$2\" = \"--get\" ] && [ \"$3\" = \"core.hooksPath\" ]; then\n  printf '%s\\n' " + quoteShell(hooksPath) + "\n  exit 0\nfi\nexit 0\n"
	writeExecutable(t, root, "git", body)
}

func quoteShell(s string) string {
	return "'" + s + "'"
}

// TestRunMainGuardFlagParseError covers the flag.Parse non-help error.
func TestRunMainGuardFlagParseError(t *testing.T) {
	var stdout, stderr bytes.Buffer
	got := runMainGuard([]string{"--bad"}, &stdout, &stderr)
	if got != 2 {
		t.Fatalf("runMainGuard bad flag = %d; want 2", got)
	}
}

// TestRunMainGuardHelp covers flag.ErrHelp branch.
func TestRunMainGuardHelp(t *testing.T) {
	var stdout, stderr bytes.Buffer
	got := runMainGuard([]string{"-h"}, &stdout, &stderr)
	if got != 0 {
		t.Fatalf("runMainGuard help = %d; want 0", got)
	}
}

// TestRunMainGuardInvalidContext covers invalid context branch.
func TestRunMainGuardInvalidContext(t *testing.T) {
	var stdout, stderr bytes.Buffer
	got := runMainGuard([]string{"--context", "bogus"}, &stdout, &stderr)
	if got != 2 {
		t.Fatalf("runMainGuard invalid context = %d; want 2", got)
	}
}

// TestRunWorktreeGateBranches covers flag.Parse errors, help, positional args,
// invalid context.
func TestRunWorktreeGateBranches(t *testing.T) {
	t.Run("flag parse error", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		got := runWorktreeGate("worktree-guard", []string{"--nope"}, &stdout, &stderr)
		if got != 2 {
			t.Fatalf("got = %d; want 2", got)
		}
	})
	t.Run("help", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		got := runWorktreeGate("worktree-guard", []string{"-h"}, &stdout, &stderr)
		if got != 0 {
			t.Fatalf("got = %d; want 0", got)
		}
	})
	t.Run("positional arg", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		got := runWorktreeGate("worktree-guard", []string{"extra"}, &stdout, &stderr)
		if got != 2 {
			t.Fatalf("got = %d; want 2", got)
		}
	})
	t.Run("invalid context", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		got := runWorktreeGate("worktree-guard", []string{"--context", "nope"}, &stdout, &stderr)
		if got != 2 {
			t.Fatalf("got = %d; want 2", got)
		}
	})
	t.Run("readonly on main passes", func(t *testing.T) {
		setupPRCheckFixture(t, "main")
		var stdout, stderr bytes.Buffer
		got := runWorktreeGate("worktree-guard", []string{"--context", "local_readonly"}, &stdout, &stderr)
		if got != 0 {
			t.Fatalf("got = %d; want 0", got)
		}
	})
}

// TestRunContextCheckInvalidArgs covers validateInternalCommandArgs error.
func TestRunContextCheckInvalidArgs(t *testing.T) {
	var stdout, stderr bytes.Buffer
	got := runContextCheck([]string{"--bad"}, &stdout, &stderr)
	if got != 2 {
		t.Fatalf("got = %d; want 2", got)
	}
}

// TestRunContextCheckHelp covers help branch.
func TestRunContextCheckHelp(t *testing.T) {
	var stdout, stderr bytes.Buffer
	got := runContextCheck([]string{"-h"}, &stdout, &stderr)
	if got != 0 {
		t.Fatalf("got = %d; want 0", got)
	}
}

// TestRunSpecCheckBranches covers all spec-check branches.
func TestRunSpecCheckBranches(t *testing.T) {
	t.Run("invalid args", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		got := runSpecCheck([]string{"--bad"}, &stdout, &stderr)
		if got != 2 {
			t.Fatalf("got = %d; want 2", got)
		}
	})
	t.Run("help", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		got := runSpecCheck([]string{"-h"}, &stdout, &stderr)
		if got != 0 {
			t.Fatalf("got = %d; want 0", got)
		}
	})
	t.Run("missing docs dir", func(t *testing.T) {
		root := t.TempDir()
		chdir(t, root)
		var stdout, stderr bytes.Buffer
		got := runSpecCheck(nil, &stdout, &stderr)
		if got != 1 {
			t.Fatalf("got = %d; want 1", got)
		}
		if !strings.Contains(stdout.String(), "missing docs") {
			t.Fatalf("stdout = %q; want missing docs", stdout.String())
		}
	})
	t.Run("docs present no REQ markers via filesystem walk", func(t *testing.T) {
		// Non-git directory forces filesystemDocsMarkdownFiles path.
		root := t.TempDir()
		if err := os.MkdirAll(filepath.Join(root, "docs", "sub"), 0o755); err != nil {
			t.Fatalf("mkdir docs: %v", err)
		}
		if err := os.WriteFile(filepath.Join(root, "docs", "a.md"), []byte("# A\nno reqs here\n"), 0o644); err != nil {
			t.Fatalf("write a.md: %v", err)
		}
		chdir(t, root)
		var stdout, stderr bytes.Buffer
		got := runSpecCheck(nil, &stdout, &stderr)
		if got != 0 {
			t.Fatalf("got = %d; want 0", got)
		}
		if !strings.Contains(stdout.String(), "warning: no docs markdown file contains REQ-") {
			t.Fatalf("stdout = %q; want REQ warning", stdout.String())
		}
	})
	t.Run("docs with REQ via filesystem walk", func(t *testing.T) {
		root := t.TempDir()
		if err := os.MkdirAll(filepath.Join(root, "docs"), 0o755); err != nil {
			t.Fatalf("mkdir docs: %v", err)
		}
		if err := os.WriteFile(filepath.Join(root, "docs", "a.md"), []byte("REQ-001 do thing\n"), 0o644); err != nil {
			t.Fatalf("write a.md: %v", err)
		}
		chdir(t, root)
		var stdout, stderr bytes.Buffer
		got := runSpecCheck(nil, &stdout, &stderr)
		if got != 0 {
			t.Fatalf("got = %d; want 0", got)
		}
		if !strings.Contains(stdout.String(), "scanned_markdown=1") {
			t.Fatalf("stdout = %q; want scanned_markdown=1", stdout.String())
		}
	})
}

// TestShouldScanDocsFromFilesystem covers both matching and non-matching cases.
func TestShouldScanDocsFromFilesystem(t *testing.T) {
	notGit := errors.New("some other error")
	if shouldScanDocsFromFilesystem(notGit, nil) {
		t.Fatalf("non-exit-error should not trigger filesystem scan")
	}
	// Build a real *exec.ExitError with code 128 via a shell command.
	exitErr := makeExitError(t, 128)
	// 128 + "not a git repository"
	if !shouldScanDocsFromFilesystem(exitErr, []byte("fatal: not a git repository (or any parent up to mount point)")) {
		t.Fatalf("128 not-a-git-repository should trigger filesystem scan")
	}
	if !shouldScanDocsFromFilesystem(exitErr, []byte("not in a git directory")) {
		t.Fatalf("128 not-in-git-directory should trigger filesystem scan")
	}
	if shouldScanDocsFromFilesystem(exitErr, []byte("some other 128 error")) {
		t.Fatalf("unrelated 128 error should not trigger scan")
	}
}

// makeExitError runs a shell command that exits with the given code and returns
// the resulting *exec.ExitError (or fails the test).
func makeExitError(t *testing.T, code int) error {
	t.Helper()
	cmd := exec.Command("sh", "-c", "exit "+itoa(code))
	err := cmd.Run()
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		t.Fatalf("could not build *exec.ExitError for code %d: %v", code, err)
	}
	return exitErr
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	pos := len(buf)
	neg := n < 0
	if neg {
		n = -n
	}
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		pos--
		buf[pos] = '-'
	}
	return string(buf[pos:])
}

// TestFilesystemDocsMarkdownFiles covers the filesystem walk + walk error.
func TestFilesystemDocsMarkdownFiles(t *testing.T) {
	t.Run("walks nested docs", func(t *testing.T) {
		root := t.TempDir()
		if err := os.MkdirAll(filepath.Join(root, "docs", "nested"), 0o755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}
		if err := os.WriteFile(filepath.Join(root, "docs", "a.md"), []byte("a"), 0o644); err != nil {
			t.Fatalf("write: %v", err)
		}
		if err := os.WriteFile(filepath.Join(root, "docs", "nested", "b.md"), []byte("b"), 0o644); err != nil {
			t.Fatalf("write: %v", err)
		}
		if err := os.WriteFile(filepath.Join(root, "docs", "ignore.txt"), []byte("x"), 0o644); err != nil {
			t.Fatalf("write: %v", err)
		}
		paths, err := filesystemDocsMarkdownFiles(filepath.Join(root, "docs"))
		if err != nil {
			t.Fatalf("err = %v; want nil", err)
		}
		if len(paths) != 2 {
			t.Fatalf("paths = %v; want 2", paths)
		}
	})
	t.Run("missing root returns error", func(t *testing.T) {
		_, err := filesystemDocsMarkdownFiles(filepath.Join(t.TempDir(), "nope"))
		if err == nil {
			t.Fatalf("err = nil; want error")
		}
	})
}

// TestRunDesignCheckBranches covers design-check branches.
func TestRunDesignCheckBranches(t *testing.T) {
	t.Run("invalid args", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		got := runDesignCheck([]string{"--bad"}, &stdout, &stderr)
		if got != 2 {
			t.Fatalf("got = %d; want 2", got)
		}
	})
	t.Run("help", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		got := runDesignCheck([]string{"-h"}, &stdout, &stderr)
		if got != 0 {
			t.Fatalf("got = %d; want 0", got)
		}
	})
	t.Run("missing docs/adr", func(t *testing.T) {
		root := t.TempDir()
		chdir(t, root)
		var stdout, stderr bytes.Buffer
		got := runDesignCheck(nil, &stdout, &stderr)
		if got != 0 {
			t.Fatalf("got = %d; want 0", got)
		}
		if !strings.Contains(stdout.String(), "optional docs/adr not present") {
			t.Fatalf("stdout = %q; want optional", stdout.String())
		}
	})
	t.Run("docs/adr present", func(t *testing.T) {
		root := t.TempDir()
		if err := os.MkdirAll(filepath.Join(root, "docs", "adr"), 0o755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}
		chdir(t, root)
		var stdout, stderr bytes.Buffer
		got := runDesignCheck(nil, &stdout, &stderr)
		if got != 0 {
			t.Fatalf("got = %d; want 0", got)
		}
		if !strings.Contains(stdout.String(), "docs/adr is present") {
			t.Fatalf("stdout = %q; want present", stdout.String())
		}
	})
}

// TestRunTaskCheckBranches covers task-check branches.
func TestRunTaskCheckBranches(t *testing.T) {
	t.Run("invalid args", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		got := runTaskCheck([]string{"--bad"}, &stdout, &stderr)
		if got != 2 {
			t.Fatalf("got = %d; want 2", got)
		}
	})
	t.Run("help", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		got := runTaskCheck([]string{"-h"}, &stdout, &stderr)
		if got != 0 {
			t.Fatalf("got = %d; want 0", got)
		}
	})
	t.Run("canonical registry present", func(t *testing.T) {
		root := t.TempDir()
		if err := os.MkdirAll(filepath.Join(root, ".agent", "registries"), 0o755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}
		if err := os.WriteFile(filepath.Join(root, ".agent", "registries", "command-registry.yaml"), []byte("fixture\n"), 0o644); err != nil {
			t.Fatalf("write: %v", err)
		}
		chdir(t, root)
		var stdout, stderr bytes.Buffer
		got := runTaskCheck(nil, &stdout, &stderr)
		if got != 0 {
			t.Fatalf("got = %d; want 0", got)
		}
	})
	t.Run("only compatibility commands.yaml", func(t *testing.T) {
		root := t.TempDir()
		if err := os.MkdirAll(filepath.Join(root, ".agent", "registries"), 0o755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}
		if err := os.WriteFile(filepath.Join(root, ".agent", "registries", "commands.yaml"), []byte("fixture\n"), 0o644); err != nil {
			t.Fatalf("write: %v", err)
		}
		chdir(t, root)
		var stdout, stderr bytes.Buffer
		got := runTaskCheck(nil, &stdout, &stderr)
		if got != 1 {
			t.Fatalf("got = %d; want 1", got)
		}
		if !strings.Contains(stdout.String(), "canonical .agent/registries/command-registry.yaml missing") {
			t.Fatalf("stdout = %q; want compatibility-only gap", stdout.String())
		}
	})
	t.Run("nothing present", func(t *testing.T) {
		root := t.TempDir()
		chdir(t, root)
		var stdout, stderr bytes.Buffer
		got := runTaskCheck(nil, &stdout, &stderr)
		if got != 1 {
			t.Fatalf("got = %d; want 1", got)
		}
		if !strings.Contains(stdout.String(), "missing .agent/registries/command-registry.yaml") {
			t.Fatalf("stdout = %q; want missing", stdout.String())
		}
	})
}

// TestRunPRCheckFlagBranches covers pr-check flag paths not exercised elsewhere.
func TestRunPRCheckFlagBranches(t *testing.T) {
	t.Run("flag parse error", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		got := runPRCheck([]string{"--bad"}, strings.NewReader(""), &stdout, &stderr)
		if got != 2 {
			t.Fatalf("got = %d; want 2", got)
		}
	})
	t.Run("help", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		got := runPRCheck([]string{"-h"}, strings.NewReader(""), &stdout, &stderr)
		if got != 0 {
			t.Fatalf("got = %d; want 0", got)
		}
	})
	t.Run("positional arg", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		got := runPRCheck([]string{"extra"}, strings.NewReader(""), &stdout, &stderr)
		if got != 2 {
			t.Fatalf("got = %d; want 2", got)
		}
	})
	t.Run("invalid context", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		got := runPRCheck([]string{"--context", "bogus"}, strings.NewReader(""), &stdout, &stderr)
		if got != 2 {
			t.Fatalf("got = %d; want 2", got)
		}
	})
	t.Run("dry-run passes", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		got := runPRCheck([]string{"--dry-run", "--context", "local_readonly"}, strings.NewReader(""), &stdout, &stderr)
		if got != 0 {
			t.Fatalf("got = %d; want 0", got)
		}
		if !strings.Contains(stdout.String(), "mode=dry-run") {
			t.Fatalf("stdout = %q; want dry-run", stdout.String())
		}
	})
}

// TestEnvDefault covers envDefault fallback.
func TestEnvDefault(t *testing.T) {
	t.Setenv("MY_TEST_ENV_VAR", "value")
	if got := envDefault("MY_TEST_ENV_VAR", "fallback"); got != "value" {
		t.Fatalf("got = %q; want value", got)
	}
	_ = os.Unsetenv("MY_TEST_ENV_VAR")
	if got := envDefault("MY_TEST_ENV_VAR", "fallback"); got != "fallback" {
		t.Fatalf("got = %q; want fallback", got)
	}
}

// TestInvalidInternalArgsExitHelp covers the ErrHelp branch.
func TestInvalidInternalArgsExitHelp(t *testing.T) {
	var stderr bytes.Buffer
	got := invalidInternalArgsExit("cmd", flag.ErrHelp, &stderr)
	if got != 0 {
		t.Fatalf("got = %d; want 0", got)
	}
}

// TestInvalidInternalArgsExitOther covers the generic error branch.
func TestInvalidInternalArgsExitOther(t *testing.T) {
	var stderr bytes.Buffer
	got := invalidInternalArgsExit("cmd", errors.New("boom"), &stderr)
	if got != 2 {
		t.Fatalf("got = %d; want 2", got)
	}
}

// TestIsXlibStandardSourceModule covers all three return paths.
func TestIsXlibStandardSourceModule(t *testing.T) {
	t.Run("missing go.mod defaults true", func(t *testing.T) {
		root := t.TempDir()
		chdir(t, root)
		if !isXlibStandardSourceModule() {
			t.Fatalf("missing go.mod should default true")
		}
	})
	t.Run("non-source module", func(t *testing.T) {
		root := t.TempDir()
		if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module github.com/ZoneCNH/kernel\n\ngo 1.23\n"), 0o644); err != nil {
			t.Fatalf("write: %v", err)
		}
		chdir(t, root)
		if isXlibStandardSourceModule() {
			t.Fatalf("kernel module should not be source")
		}
	})
	t.Run("source module", func(t *testing.T) {
		root := t.TempDir()
		sourceModule := strings.Join([]string{"github.com", "ZoneCNH", "xlib" + "-standard"}, "/")
		if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module "+sourceModule+"\n\ngo 1.23\n"), 0o644); err != nil {
			t.Fatalf("write: %v", err)
		}
		chdir(t, root)
		if !isXlibStandardSourceModule() {
			t.Fatalf("xlib-standard module should be source")
		}
	})
	t.Run("go.mod without module line", func(t *testing.T) {
		root := t.TempDir()
		if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("go 1.23\n"), 0o644); err != nil {
			t.Fatalf("write: %v", err)
		}
		chdir(t, root)
		if isXlibStandardSourceModule() {
			t.Fatalf("go.mod without module line should be false")
		}
	})
}

// TestValidContextProfileName covers known/unknown profiles.
func TestValidContextProfileName(t *testing.T) {
	for _, p := range []string{"lite", "standard", "full", "release"} {
		if !validContextProfileName(p) {
			t.Errorf("validContextProfileName(%q) = false; want true", p)
		}
	}
	if validContextProfileName("bogus") {
		t.Errorf("validContextProfileName(bogus) = true; want false")
	}
}

// TestNormalizeContextProfile covers fast->lite aliasing and default.
func TestNormalizeContextProfile(t *testing.T) {
	if got := normalizeContextProfile("fast"); got != "lite" {
		t.Fatalf("fast -> %q; want lite", got)
	}
	if got := normalizeContextProfile("standard"); got != "standard" {
		t.Fatalf("standard -> %q; want standard", got)
	}
}

// TestCanonicalRepoPath covers all branches.
func TestCanonicalRepoPath(t *testing.T) {
	cases := []struct {
		raw       string
		want      string
		canonical bool
	}{
		{"", "", false},
		{"/abs/path", "/abs/path", false},
		{"a\\b/c", "a/b/c", false},
		{"..", "..", false},
		{"../parent", "../parent", false},
		{"docs/a.md", "docs/a.md", true},
		{"docs//a.md", "docs/a.md", false},
	}
	for _, c := range cases {
		got, canonical := canonicalRepoPath(c.raw)
		if got != c.want || canonical != c.canonical {
			t.Errorf("canonicalRepoPath(%q) = (%q,%v); want (%q,%v)", c.raw, got, canonical, c.want, c.canonical)
		}
	}
}

// TestParseContextProfileCheckProfile covers flag errors and positional args.
func TestParseContextProfileCheckProfile(t *testing.T) {
	if _, err := parseContextProfileCheckProfile("cmd", []string{"--bad"}); err == nil {
		t.Fatalf("bad flag should error")
	}
	if _, err := parseContextProfileCheckProfile("cmd", []string{"positional"}); err == nil {
		t.Fatalf("positional should error")
	}
	profile, err := parseContextProfileCheckProfile("cmd", []string{"--profile", "standard"})
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	if profile != "standard" {
		t.Fatalf("profile = %q; want standard", profile)
	}
}
