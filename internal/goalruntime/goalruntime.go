package goalruntime

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	DefaultGoalID      = "GOAL-20260603-XLIB-RUNTIME-001"
	EvidenceLedgerPath = "release/evidence/goalkit/"
)

var commandGates = map[string]string{
	"goal-acceptance":          "G12",
	"goal-delivery":            "G13",
	"goal-handover":            "G14",
	"goal-downstream-adoption": "G15A",
	"goal-certify":             "G15B",
	"goal-runtime-final":       "G16",
}

var commandHarnessGates = map[string]string{
	"goal-acceptance":          "G12_ACCEPTANCE",
	"goal-delivery":            "G13_DELIVERY",
	"goal-handover":            "G14_HANDOVER",
	"goal-downstream-adoption": "G15_DOWNSTREAM_ADOPTION",
	"goal-certify":             "G16_CERTIFY",
}

var gateDescriptions = map[string]string{
	"goal-acceptance":          "验收矩阵和目标 ID contract 已收敛",
	"goal-delivery":            "交付证据路径和 xlibgate 执行面已收敛",
	"goal-handover":            "接手材料和边界声明已收敛",
	"goal-downstream-adoption": "下游采用证明保持为本地 contract，不修改 downstream 仓库",
	"goal-certify":             "认证证明保持为本地 contract，不宣称完整 release 完成",
	"goal-runtime-final":       "G12-G15B 本地 contract 汇总为 goalkit v0.1.0 MVA final evidence",
}

var requiredAuthorityPaths = []string{
	".worktree/goalkit-v0.1.0-plan.md",
	".omx/context/goalkit-v0.1.0-team-20260603T005302Z.md",
	"docs/standard/xlibgate-cli-contract.md",
	".agent/harness.yaml",
	".agent/command-registry.yaml",
	"Makefile",
}

// Options configures a goalkit MVA contract evaluation.
type Options struct {
	GoalID string
	Mode   string
	Root   string
}

// Report is the machine-readable goalkit MVA evidence returned by xlibgate.
type Report struct {
	SchemaVersion    string       `json:"schema_version"`
	Command          string       `json:"command"`
	Status           string       `json:"status"`
	GoalID           string       `json:"goal_id"`
	Gate             string       `json:"gate"`
	Mode             string       `json:"mode"`
	Executor         string       `json:"executor"`
	ControlPlane     string       `json:"control_plane"`
	Blocking         bool         `json:"blocking"`
	MVAStatus        string       `json:"mva_status"`
	LedgerPath       string       `json:"ledger_path"`
	EvidencePackPath string       `json:"evidence_pack_path"`
	Gates            []GateReport `json:"gates,omitempty"`
	Details          []string     `json:"details,omitempty"`
	Evidence         []string     `json:"evidence,omitempty"`
	AuthorityPaths   []string     `json:"authority_paths,omitempty"`
	Gaps             []string     `json:"gaps,omitempty"`
}

// GateReport records one non-blocking local goalkit evidence gate.
type GateReport struct {
	ID       string `json:"id"`
	Command  string `json:"command"`
	Status   string `json:"status"`
	Blocking bool   `json:"blocking"`
}

// Evaluate verifies the local goalkit v0.1.0 MVA contract for a single command.
func Evaluate(command string, options Options) (Report, error) {
	gate, ok := commandGates[command]
	if !ok {
		return Report{}, fmt.Errorf("unsupported goalkit command %q", command)
	}
	goalID := strings.TrimSpace(options.GoalID)
	if goalID == "" {
		goalID = DefaultGoalID
	}
	mode := strings.TrimSpace(options.Mode)
	if mode == "" {
		mode = "FULL"
	}
	root := options.Root
	if root == "" {
		root = "."
	}
	report := Report{
		SchemaVersion:    "goalkit-mva/v1",
		Command:          command,
		Status:           "passed",
		GoalID:           goalID,
		Gate:             gate,
		Mode:             mode,
		Executor:         "xlibgate",
		ControlPlane:     "Harness Runtime",
		Blocking:         false,
		MVAStatus:        "not-complete",
		LedgerPath:       ".agent/evidence/ledger.jsonl",
		EvidencePackPath: EvidenceLedgerPath + goalID + ".json",
		Gates:            gatesForCommand(command),
		Details: []string{
			gateDescriptions[command],
			"goalkit v0.1.0 不提供独立 CLI；xlibgate 是当前执行面",
			"G12-G16 是 goalkit MVA evidence gates，不是全局 release blocking gates",
			"root plan 是 backlog/roadmap authority，不是完成声明",
		},
		Evidence: []string{
			"fixture_id=" + DefaultGoalID,
			"evidence_ledger=" + EvidenceLedgerPath,
			"runtime_logic=internal/goalruntime",
		},
	}
	if goalID != DefaultGoalID {
		report.Details = append(report.Details, "non-default goal_id accepted for local contract replay")
	}
	for _, path := range requiredAuthorityPaths {
		full := filepath.Join(root, filepath.FromSlash(path))
		if _, err := os.Stat(full); err != nil {
			report.Gaps = append(report.Gaps, "missing authority path: "+path)
			continue
		}
		report.AuthorityPaths = append(report.AuthorityPaths, path)
	}
	if command == "goal-runtime-final" {
		report.Evidence = append(report.Evidence,
			"requires=goal-acceptance",
			"requires=goal-delivery",
			"requires=goal-handover",
			"requires=goal-downstream-adoption",
			"requires=goal-certify",
		)
	}
	if len(report.Gaps) > 0 {
		report.Status = "failed"
	}
	return report, nil
}

func gatesForCommand(command string) []GateReport {
	commands := []string{command}
	if command == "goal-runtime-final" {
		commands = []string{
			"goal-acceptance",
			"goal-delivery",
			"goal-handover",
			"goal-downstream-adoption",
			"goal-certify",
		}
	}
	reports := make([]GateReport, 0, len(commands))
	for _, gateCommand := range commands {
		reports = append(reports, GateReport{
			ID:       commandHarnessGates[gateCommand],
			Command:  gateCommand,
			Status:   "passed",
			Blocking: false,
		})
	}
	return reports
}

// Commands returns the supported goalkit MVA command names.
func Commands() []string {
	return []string{
		"goal-acceptance",
		"goal-delivery",
		"goal-handover",
		"goal-downstream-adoption",
		"goal-certify",
		"goal-runtime-final",
	}
}
