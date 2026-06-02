package debtcheck

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const Runtime = "debt-governance-runtime-v3.6"
const SchemaVersion = "1.0.0"

var GateNames = []string{
	"architecture",
	"domain",
	"docs-drift",
	"dependency-debt",
	"security-debt",
	"testing-debt",
	"implementation-debt",
}

var DownstreamTargets = []string{"kernel", "configx", "redisx"}

var policyFiles = []string{
	".agent/debt/rules.yaml",
	".agent/debt/rule-registry.yaml",
	".agent/debt/profile.yaml",
	".agent/debt/register.md",
}

var gateAnchors = map[string][]string{
	"architecture":        {".agent/harness.yaml", ".agent/traceability-matrix.md", "docs/standard/harness-gates.md"},
	"domain":              {"docs/goal.md", "contracts/xlibgate-report.schema.json"},
	"docs-drift":          {"docs/standard/xlibgate-cli-contract.md", ".agent/command-registry.yaml"},
	"dependency-debt":     {"go.mod", "scripts/check_dependency_diff.sh"},
	"security-debt":       {".agent/security.yaml", "scripts/check_secrets.sh"},
	"testing-debt":        {"cmd/xlibgate/main_test.go", "internal/tools/releasemanifest/main_test.go"},
	"implementation-debt": {"cmd/xlibgate/main.go", "Makefile"},
}

type Report struct {
	SchemaVersion     string            `json:"schema_version"`
	Runtime           string            `json:"runtime"`
	Command           string            `json:"command"`
	Status            string            `json:"status"`
	PolicyFiles       []string          `json:"policy_files"`
	DownstreamTargets []string          `json:"downstream_targets"`
	Gates             []GateReport      `json:"gates"`
	GateStatuses      map[string]string `json:"gate_statuses"`
	Details           []string          `json:"details,omitempty"`
	Gaps              []string          `json:"gaps,omitempty"`
}

type GateReport struct {
	Name    string   `json:"name"`
	Status  string   `json:"status"`
	Details []string `json:"details,omitempty"`
	Gaps    []string `json:"gaps,omitempty"`
}

func IsGate(command string) bool {
	for _, gate := range GateNames {
		if command == gate {
			return true
		}
	}
	return false
}

func Evaluate(root, command string) Report {
	if root == "" {
		root = "."
	}
	if command == "" {
		command = "debt"
	}
	report := Report{
		SchemaVersion:     SchemaVersion,
		Runtime:           Runtime,
		Command:           command,
		PolicyFiles:       append([]string(nil), policyFiles...),
		DownstreamTargets: append([]string(nil), DownstreamTargets...),
		GateStatuses:      map[string]string{},
	}

	policyGaps := policyGaps(root)
	selected := []string{command}
	if command == "debt" {
		selected = append([]string(nil), GateNames...)
	} else if !IsGate(command) {
		report.Status = "failed"
		report.Gaps = []string{"unknown debt gate " + command}
		report.GateStatuses[command] = "failed"
		return report
	}

	allPassed := len(policyGaps) == 0
	if command == "debt" {
		report.GateStatuses["debt"] = statusFromGaps(policyGaps)
	}
	for _, name := range selected {
		if name == "debt" {
			continue
		}
		gate := evaluateGate(root, name, policyGaps)
		report.Gates = append(report.Gates, gate)
		report.GateStatuses[name] = gate.Status
		if gate.Status != "passed" {
			allPassed = false
		}
	}
	if command != "debt" {
		report.GateStatuses[command] = statusFromGaps(append(policyGaps, gateGaps(report.Gates)...))
	}
	if allPassed {
		report.Status = "passed"
		report.Details = []string{"debt policy files present", "scanner failure is fail-closed", "P0 debt exceptions are forbidden"}
	} else {
		report.Status = "failed"
		report.Gaps = append(report.Gaps, policyGaps...)
		report.Gaps = append(report.Gaps, gateGaps(report.Gates)...)
	}
	return report
}

func WriteEvidence(root, outDir string, report Report) error {
	if outDir == "" {
		outDir = "release/debt"
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}
	jsonBytes, err := marshalReport(report)
	if err != nil {
		return err
	}
	jsonPath := filepath.Join(outDir, "latest.json")
	if err := os.WriteFile(jsonPath, jsonBytes, 0o644); err != nil {
		return err
	}
	sum := sha256.Sum256(jsonBytes)
	if err := os.WriteFile(jsonPath+".sha256", []byte(hex.EncodeToString(sum[:])+"  latest.json\n"), 0o644); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(outDir, "latest.md"), markdownReport(report), 0o644)
}

func marshalReport(report Report) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	if err := enc.Encode(report); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func markdownReport(report Report) []byte {
	var b strings.Builder
	fmt.Fprintf(&b, "# Debt Governance Evidence\n\n")
	fmt.Fprintf(&b, "runtime: %s\n", report.Runtime)
	fmt.Fprintf(&b, "schema_version: %s\n", report.SchemaVersion)
	fmt.Fprintf(&b, "status: %s\n", report.Status)
	fmt.Fprintf(&b, "command: %s\n", report.Command)
	fmt.Fprintf(&b, "downstream_targets: %s\n\n", strings.Join(report.DownstreamTargets, ","))
	fmt.Fprintf(&b, "## Gate statuses\n\n")
	keys := make([]string, 0, len(report.GateStatuses))
	for key := range report.GateStatuses {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		fmt.Fprintf(&b, "- %s: %s\n", key, report.GateStatuses[key])
	}
	if len(report.Gaps) > 0 {
		fmt.Fprintf(&b, "\n## Gaps\n\n")
		for _, gap := range report.Gaps {
			fmt.Fprintf(&b, "- %s\n", gap)
		}
	}
	return []byte(b.String())
}

func policyGaps(root string) []string {
	var gaps []string
	for _, path := range policyFiles {
		if !fileExists(root, path) {
			gaps = append(gaps, "missing "+path)
		}
	}
	if text, ok := readText(root, ".agent/debt/rules.yaml"); ok && !strings.Contains(text, "p0_exceptions: forbidden") {
		gaps = append(gaps, ".agent/debt/rules.yaml must forbid P0 exceptions")
	}
	if text, ok := readText(root, ".agent/debt/rule-registry.yaml"); ok && !strings.Contains(text, "fail_closed: true") {
		gaps = append(gaps, ".agent/debt/rule-registry.yaml must fail closed")
	}
	return gaps
}

func evaluateGate(root, name string, policyGaps []string) GateReport {
	gate := GateReport{Name: name}
	for _, path := range gateAnchors[name] {
		if fileExists(root, path) {
			gate.Details = append(gate.Details, "found "+path)
		} else {
			gate.Gaps = append(gate.Gaps, "missing "+path)
		}
	}
	gate.Gaps = append(gate.Gaps, policyGaps...)
	gate.Status = statusFromGaps(gate.Gaps)
	return gate
}

func gateGaps(gates []GateReport) []string {
	var gaps []string
	for _, gate := range gates {
		for _, gap := range gate.Gaps {
			gaps = append(gaps, gate.Name+": "+gap)
		}
	}
	return gaps
}

func statusFromGaps(gaps []string) string {
	if len(gaps) == 0 {
		return "passed"
	}
	return "failed"
}

func fileExists(root, path string) bool {
	info, err := os.Stat(filepath.Join(root, filepath.FromSlash(path)))
	return err == nil && !info.IsDir()
}

func readText(root, path string) (string, bool) {
	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(path)))
	if err != nil {
		return "", false
	}
	return string(data), true
}
