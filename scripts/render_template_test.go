package scripts_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestRenderTemplateExcludesGeneratedReleaseArtifacts(t *testing.T) {
	contents, err := os.ReadFile("render_template.sh")
	if err != nil {
		t.Fatalf("read render_template.sh: %v", err)
	}

	script := string(contents)
	for _, required := range []string{
		"--exclude='./.omc'",
		"--exclude='./.omx'",
		"--exclude='./.worktree'",
		"--exclude='./release/manifest/latest.json'",
		"--exclude='./release/manifest/latest.json.sha256'",
		"--exclude='./release/standard-impact/latest.md'",
		"--exclude='./release/downstream-sync/latest.md'",
		"--exclude='./release/debt/latest.json'",
		"--exclude='./release/debt/latest.md'",
		"--exclude='./release/debt/latest.json.sha256'",
		`rm -rf "$out_dir/.omc"`,
		`rm -rf "$out_dir/.omx"`,
		`rm -rf "$out_dir/.worktree"`,
		`rm -f "$out_dir/release/manifest/latest.json"`,
		`rm -f "$out_dir/release/manifest/latest.json.sha256"`,
		`rm -f "$out_dir/release/standard-impact/latest.md"`,
		`rm -f "$out_dir/release/downstream-sync/latest.md"`,
		`rm -f "$out_dir/release/debt/latest.json"`,
		`rm -f "$out_dir/release/debt/latest.md"`,
		`rm -f "$out_dir/release/debt/latest.json.sha256"`,
	} {
		if !strings.Contains(script, required) {
			t.Fatalf("render_template.sh missing render omission rule %q", required)
		}
	}
}

func TestRenderTemplatePrunesOmittedAgentInboxIndexEntries(t *testing.T) {
	outDir := filepath.Join(t.TempDir(), "kernel")
	cmd := exec.Command(
		"bash",
		"render_template.sh",
		"--module-name",
		"kernel",
		"--module-path",
		"github.com/ZoneCNH/kernel",
		"--package-name",
		"kernel",
		"--out",
		outDir,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("render template: %v\n%s", err, output)
	}

	index, err := os.ReadFile(filepath.Join(outDir, ".agent", "index.yaml"))
	if err != nil {
		t.Fatalf("read rendered agent index: %v", err)
	}
	if strings.Contains(string(index), ".agent/inbox/") {
		t.Fatalf("rendered agent index still references omitted inbox entries")
	}
	if _, err := os.Stat(filepath.Join(outDir, ".agent", "inbox")); !os.IsNotExist(err) {
		t.Fatalf("rendered agent inbox should be omitted, stat err=%v", err)
	}
}

func TestRenderTemplateGitArchivePrunesRuntimeState(t *testing.T) {
	outDir := filepath.Join(t.TempDir(), "redisx")
	cmd := exec.Command(
		"bash",
		"render_template.sh",
		"--module-name",
		"redisx",
		"--module-path",
		"github.com/ZoneCNH/redisx",
		"--package-name",
		"redisx",
		"--out",
		outDir,
	)
	cmd.Env = append(os.Environ(), "XLIB_RENDER_FORCE_GIT_ARCHIVE=1")

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("render template from git archive: %v\n%s", err, output)
	}

	for _, omitted := range []string{
		".omc",
		".omx",
		".worktree",
		filepath.Join("release", "manifest", "latest.json"),
		filepath.Join("release", "manifest", "latest.json.sha256"),
		filepath.Join("release", "standard-impact", "latest.md"),
		filepath.Join("release", "downstream-sync", "latest.md"),
		filepath.Join("release", "debt", "latest.json"),
		filepath.Join("release", "debt", "latest.md"),
		filepath.Join("release", "debt", "latest.json.sha256"),
	} {
		if _, err := os.Stat(filepath.Join(outDir, omitted)); !os.IsNotExist(err) {
			t.Fatalf("rendered git archive should omit %s, stat err=%v", omitted, err)
		}
	}
}

func TestRenderTemplateGitArchiveSkipsUntrackedFiles(t *testing.T) {
	markerPath := filepath.Join("..", ".xlib-render-untracked-marker-test")
	if err := os.WriteFile(markerPath, []byte("untracked marker"), 0o600); err != nil {
		t.Fatalf("write untracked marker: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Remove(markerPath)
	})

	outDir := filepath.Join(t.TempDir(), "kernel")
	cmd := exec.Command(
		"bash",
		"render_template.sh",
		"--module-name",
		"kernel",
		"--module-path",
		"github.com/ZoneCNH/kernel",
		"--package-name",
		"kernel",
		"--out",
		outDir,
	)
	cmd.Env = append(os.Environ(), "XLIB_RENDER_FORCE_GIT_ARCHIVE=1")

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("render template: %v\n%s", err, output)
	}

	if _, err := os.Stat(filepath.Join(outDir, ".xlib-render-untracked-marker-test")); !os.IsNotExist(err) {
		t.Fatalf("expected git archive render to skip untracked marker, stat err=%v", err)
	}
}
