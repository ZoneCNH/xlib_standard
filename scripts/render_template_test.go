package scripts_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestRenderTemplateExcludesGeneratedDebtArtifacts(t *testing.T) {
	contents, err := os.ReadFile("render_template.sh")
	if err != nil {
		t.Fatalf("read render_template.sh: %v", err)
	}

	script := string(contents)
	for _, exclude := range []string{
		"--exclude='./release/debt/latest.json'",
		"--exclude='./release/debt/latest.md'",
		"--exclude='./release/debt/latest.json.sha256'",
	} {
		if !strings.Contains(script, exclude) {
			t.Fatalf("render_template.sh missing generated debt artifact exclude %q", exclude)
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

func TestRenderTemplateGitArchiveSkipsUntrackedFiles(t *testing.T) {
	const markerName = ".xlib-render-untracked-marker-test"

	markerPath := filepath.Join("..", markerName)
	if err := os.WriteFile(markerPath, []byte("do not render\n"), 0o600); err != nil {
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
		t.Fatalf("render template with git archive: %v\n%s", err, output)
	}
	if _, err := os.Stat(filepath.Join(outDir, markerName)); !os.IsNotExist(err) {
		t.Fatalf("git archive render should skip untracked marker, stat err=%v", err)
	}
	if _, err := os.Stat(filepath.Join(outDir, "docs", "goal.md")); !os.IsNotExist(err) {
		t.Fatalf("git archive render should prune docs/goal.md, stat err=%v", err)
	}
}
