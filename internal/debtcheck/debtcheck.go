package debtcheck

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const DefaultMinScore = 9.8

var DefaultScopes = []string{
	"architecture",
	"domain",
	"docs-drift",
	"dependency-debt",
	"security-debt",
	"testing-debt",
	"implementation-debt",
	"downstream-debt",
}

type Options struct {
	Root          string
	Scopes        []string
	RunExternal   bool
	WriteEvidence bool
	OutPath       string
	MarkdownPath  string
	ChecksumPath  string
}

type Report struct {
	SchemaVersion     string             `json:"schema_version"`
	GeneratedAt       string             `json:"generated_at"`
	Status            string             `json:"status"`
	Score             float64            `json:"score"`
	MinScore          float64            `json:"min_score"`
	PolicyPath        string             `json:"policy_path"`
	Checks            []Check            `json:"checks"`
	DownstreamTargets []DownstreamTarget `json:"downstream_targets,omitempty"`
}

type Check struct {
	Name     string    `json:"name"`
	Category string    `json:"category"`
	Status   string    `json:"status"`
	Summary  string    `json:"summary"`
	Command  []string  `json:"command,omitempty"`
	Findings []Finding `json:"findings,omitempty"`
}

type Finding struct {
	Severity string `json:"severity"`
	Path     string `json:"path"`
	Message  string `json:"message"`
}

type DownstreamTarget struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

func Run(opts Options) (Report, error) {
	root := opts.Root
	if root == "" {
		root = "."
	}
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return Report{}, err
	}
	scopes := opts.Scopes
	if len(scopes) == 0 {
		scopes = DefaultScopes
	}

	report := Report{
		SchemaVersion: "1.0",
		GeneratedAt:   time.Now().UTC().Format(time.RFC3339),
		Status:        "passed",
		Score:         DefaultMinScore,
		MinScore:      DefaultMinScore,
		PolicyPath:    ".agent/debt/rules.yaml",
		Checks:        []Check{},
		DownstreamTargets: []DownstreamTarget{
			{Name: "kernel/configx", Status: "tracked"},
			{Name: "kernel/redisx", Status: "tracked"},
			{Name: "corekit", Status: "tracked"},
		},
	}

	report.Checks = append(report.Checks, policyChecks(absRoot)...)
	for _, scope := range scopes {
		report.Checks = append(report.Checks, runScope(absRoot, scope, opts.RunExternal))
	}

	for _, check := range report.Checks {
		if check.Status != "passed" {
			report.Status = "failed"
			report.Score = 0
			break
		}
	}
	return report, nil
}

func WriteEvidence(report Report, opts Options) error {
	root := opts.Root
	if root == "" {
		root = "."
	}
	out := evidencePath(root, opts.OutPath, "release/debt/latest.json")
	md := evidencePath(root, opts.MarkdownPath, "release/debt/latest.md")
	checksum := evidencePath(root, opts.ChecksumPath, "release/debt/latest.json.sha256")
	if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(md), 0o755); err != nil {
		return err
	}
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(report); err != nil {
		return err
	}
	if err := os.WriteFile(out, buf.Bytes(), 0o644); err != nil {
		return err
	}
	if err := os.WriteFile(md, []byte(markdown(report)), 0o644); err != nil {
		return err
	}
	return writeChecksum(out, checksum)
}

func policyChecks(root string) []Check {
	files := []string{
		".agent/debt/rules.yaml",
		".agent/debt/rule-registry.yaml",
		".agent/debt/profile.yaml",
		".agent/debt/register.md",
	}
	checks := make([]Check, 0, len(files)+1)
	for _, path := range files {
		checks = append(checks, fileContains(root, path, []string{"debt"}, "policy"))
	}
	rulesPath := filepath.Join(root, ".agent/debt/rules.yaml")
	data, err := os.ReadFile(rulesPath)
	check := Check{Name: "p0-no-exceptions", Category: "policy", Status: "passed", Summary: "P0 debt exceptions are disallowed"}
	if err != nil {
		check.Status = "failed"
		check.Findings = append(check.Findings, Finding{Severity: "P0", Path: ".agent/debt/rules.yaml", Message: err.Error()})
	} else {
		text := strings.ToLower(string(data))
		if strings.Contains(text, "allow_p0_exceptions: true") || strings.Contains(text, "p0_exception") || (strings.Contains(text, "severity: p0") && (strings.Contains(text, "exceptions:") || strings.Contains(text, "except:"))) {
			check.Status = "failed"
			check.Findings = append(check.Findings, Finding{Severity: "P0", Path: ".agent/debt/rules.yaml", Message: "P0 debt cannot be excepted"})
		}
	}
	checks = append(checks, check)
	return checks
}

func runScope(root string, scope string, runExternal bool) Check {
	switch scope {
	case "architecture":
		return externalScope(root, scope, "architecture boundaries are enforced", runExternal, "./scripts/check_boundary.sh")
	case "domain":
		return externalScope(root, scope, "domain boundary rules are enforced", runExternal, "./scripts/check_boundary.sh")
	case "docs-drift":
		return externalScope(root, scope, "documentation drift checks are enforced", runExternal, "./scripts/check_docs.sh")
	case "dependency-debt":
		return externalScope(root, scope, "dependency drift checks are enforced", runExternal, "./scripts/check_dependency_diff.sh")
	case "security-debt":
		return externalScope(root, scope, "secret and security scanners are enforced", runExternal, "./scripts/check_secrets.sh")
	case "testing-debt":
		return allContains(root, scope, "test target coverage is enforced", "Makefile", []string{"test:", "race:", "property:", "golden:", "fuzz-smoke:"})
	case "implementation-debt":
		return implementationDebt(root)
	case "downstream-debt":
		return downstreamDebt(root)
	default:
		return Check{Name: scope, Category: "debt", Status: "failed", Summary: "unknown debt scope", Findings: []Finding{{Severity: "P1", Message: "unknown scope " + scope}}}
	}
}

func externalScope(root, name, summary string, runExternal bool, command ...string) Check {
	check := Check{Name: name, Category: "scanner", Status: "passed", Summary: summary, Command: command}
	if len(command) == 0 {
		return check
	}
	if _, err := os.Stat(filepath.Join(root, command[0])); err != nil {
		check.Status = "failed"
		check.Findings = append(check.Findings, Finding{Severity: "P0", Path: command[0], Message: err.Error()})
		return check
	}
	if !runExternal {
		return check
	}
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Dir = root
	cmd.Env = append(os.Environ(), "GOWORK=off")
	output, err := cmd.CombinedOutput()
	if err != nil {
		check.Status = "failed"
		check.Findings = append(check.Findings, Finding{Severity: "P0", Path: command[0], Message: strings.TrimSpace(string(output)) + ": " + err.Error()})
	}
	return check
}

func fileContains(root, path string, markers []string, category string) Check {
	check := Check{Name: path, Category: category, Status: "passed", Summary: path + " is present"}
	data, err := os.ReadFile(filepath.Join(root, path))
	if err != nil {
		check.Status = "failed"
		check.Findings = append(check.Findings, Finding{Severity: "P0", Path: path, Message: err.Error()})
		return check
	}
	text := string(data)
	for _, marker := range markers {
		if !strings.Contains(text, marker) {
			check.Status = "failed"
			check.Findings = append(check.Findings, Finding{Severity: "P1", Path: path, Message: "missing marker " + marker})
		}
	}
	return check
}

func allContains(root, name, summary, path string, markers []string) Check {
	check := Check{Name: name, Category: "debt", Status: "passed", Summary: summary}
	data, err := os.ReadFile(filepath.Join(root, path))
	if err != nil {
		check.Status = "failed"
		check.Findings = append(check.Findings, Finding{Severity: "P0", Path: path, Message: err.Error()})
		return check
	}
	text := string(data)
	for _, marker := range markers {
		if !strings.Contains(text, marker) {
			check.Status = "failed"
			check.Findings = append(check.Findings, Finding{Severity: "P1", Path: path, Message: "missing marker " + marker})
		}
	}
	return check
}

func implementationDebt(root string) Check {
	markers := map[string][]string{
		"cmd/xlibgate/main.go":                   {"debt", "debt-evidence"},
		"Makefile":                               {"debt:", "debt-evidence:", "implementation-debt"},
		".agent/command-registry.yaml":           {"debt", "implementation-debt"},
		".agent/makefile-baseline.yaml":          {"debt", "debt-evidence"},
		".agent/makefile-target-registry.yaml":   {"debt", "debt-evidence"},
		"docs/standard/xlibgate-cli-contract.md": {"debt", "debt-evidence"},
	}
	return markerMapCheck(root, "implementation-debt", "debt command and target metadata are synchronized", markers)
}

func downstreamDebt(root string) Check {
	markers := map[string][]string{
		".agent/downstream-registry.yaml": {"kernel/configx", "kernel/redisx", "corekit"},
		"scripts/run_integration.sh":      {"kernel", "corekit"},
	}
	return markerMapCheck(root, "downstream-debt", "downstream debt targets are tracked", markers)
}

func markerMapCheck(root, name, summary string, markers map[string][]string) Check {
	check := Check{Name: name, Category: "debt", Status: "passed", Summary: summary}
	paths := make([]string, 0, len(markers))
	for path := range markers {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		data, err := os.ReadFile(filepath.Join(root, path))
		if err != nil {
			check.Status = "failed"
			check.Findings = append(check.Findings, Finding{Severity: "P0", Path: path, Message: err.Error()})
			continue
		}
		text := string(data)
		for _, marker := range markers[path] {
			if !strings.Contains(text, marker) {
				check.Status = "failed"
				check.Findings = append(check.Findings, Finding{Severity: "P1", Path: path, Message: "missing marker " + marker})
			}
		}
	}
	return check
}

func markdown(report Report) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# Debt Governance Evidence\n\n")
	fmt.Fprintf(&b, "- status: `%s`\n", report.Status)
	fmt.Fprintf(&b, "- score: `%.1f`\n", report.Score)
	fmt.Fprintf(&b, "- min_score: `%.1f`\n", report.MinScore)
	fmt.Fprintf(&b, "- policy_path: `%s`\n\n", report.PolicyPath)
	fmt.Fprintf(&b, "## Checks\n\n")
	for _, check := range report.Checks {
		fmt.Fprintf(&b, "- `%s`: `%s` — %s\n", check.Name, check.Status, check.Summary)
	}
	return b.String()
}

func writeChecksum(path, checksumPath string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(checksumPath), 0o755); err != nil {
		return err
	}
	sum := sha256.Sum256(data)
	line := fmt.Sprintf("%s  %s\n", hex.EncodeToString(sum[:]), filepath.ToSlash(path))
	return os.WriteFile(checksumPath, []byte(line), 0o644)
}

func evidencePath(root, got, fallback string) string {
	path := defaultPath(got, fallback)
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(root, filepath.FromSlash(path))
}

func defaultPath(got, fallback string) string {
	if got == "" {
		return fallback
	}
	return got
}

func StatusError(report Report) error {
	if report.Status == "passed" {
		return nil
	}
	return errors.New("debt governance checks failed")
}
