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
