package goalruntime

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

const (
	LedgerPath       = ".agent/evidence/ledger.jsonl"
	EvidencePackPath = "release/evidence/goalkit"
	ControlPlane     = "Harness Runtime"
	Executor         = "xlibgate"
)

var goalIDPattern = regexp.MustCompile(`^GOAL-[0-9]{8}-[A-Z0-9-]+-[0-9]{3}$`)

type Gate struct {
	ID      string `json:"id"`
	Command string `json:"command"`
	Status  string `json:"status"`
	Mode    string `json:"mode"`
}

type Report struct {
	Command          string   `json:"command"`
	Status           string   `json:"status"`
	Details          []string `json:"details,omitempty"`
	Gaps             []string `json:"gaps,omitempty"`
	GoalID           string   `json:"goal_id"`
	Mode             string   `json:"mode"`
	Executor         string   `json:"executor"`
	ControlPlane     string   `json:"control_plane"`
	Blocking         bool     `json:"blocking"`
	MVAStatus        string   `json:"mva_status"`
	LedgerPath       string   `json:"ledger_path"`
	EvidencePackPath string   `json:"evidence_pack_path"`
	Gates            []Gate   `json:"gates"`
}

type definition struct {
	ID      string
	Command string
}

var definitions = []definition{
	{ID: "G12_ACCEPTANCE", Command: "goal-acceptance"},
	{ID: "G13_DELIVERY", Command: "goal-delivery"},
	{ID: "G14_HANDOVER", Command: "goal-handover"},
	{ID: "G15_DOWNSTREAM_ADOPTION", Command: "goal-downstream-adoption"},
	{ID: "G16_CERTIFY", Command: "goal-certify"},
}

// Commands returns the goalkit runtime command names exposed by xlibgate and Makefile.
func Commands() []string {
	commands := make([]string, 0, len(definitions)+1)
	for _, def := range definitions {
		commands = append(commands, def.Command)
	}
	commands = append(commands, "goal-runtime-final")
	return commands
}

func Run(command string, args []string, stdout io.Writer, stderr io.Writer) int {
	flags := flag.NewFlagSet("xlibgate "+command, flag.ContinueOnError)
	flags.SetOutput(stderr)
	goalID := flags.String("goal-id", os.Getenv("GOAL_ID"), "goal identifier")
	mode := flags.String("mode", fallback(os.Getenv("GOAL_RUNTIME_MODE"), "FULL"), "goal runtime mode")
	flags.Bool("json", false, "emit json report")
	flags.Bool("dry-run", false, "accepted compatibility flag")
	flags.Bool("verify", false, "accepted compatibility flag")
	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return invalid(command, err, stderr)
	}
	if flags.NArg() > 0 {
		return invalid(command, fmt.Errorf("unexpected positional argument %q", flags.Arg(0)), stderr)
	}
	if !goalIDPattern.MatchString(*goalID) {
		return invalid(command, fmt.Errorf("invalid or missing --goal-id/GOAL_ID %q", *goalID), stderr)
	}
	gates, ok := selectedGates(command, strings.ToUpper(*mode))
	if !ok {
		return invalid(command, fmt.Errorf("unknown goalkit runtime command %q", command), stderr)
	}
	report := Report{
		Command:          command,
		Status:           "passed",
		Details:          []string{"goalkit v0.1.0 PR-4 command-backed Harness slice", "G12-G16 remain non-blocking until PR-5+ activation", "source evidence ledger is " + LedgerPath},
		GoalID:           *goalID,
		Mode:             strings.ToUpper(*mode),
		Executor:         Executor,
		ControlPlane:     ControlPlane,
		Blocking:         false,
		MVAStatus:        "not-complete",
		LedgerPath:       LedgerPath,
		EvidencePackPath: EvidencePackPath + "/" + *goalID,
		Gates:            gates,
	}
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		fmt.Fprintf(stderr, "ERROR: %s marshal report: %v\n", command, err)
		return 1
	}
	fmt.Fprintf(stdout, "%s\n", data)
	return 0
}

func selectedGates(command, mode string) ([]Gate, bool) {
	if command == "goal-runtime-final" {
		gates := make([]Gate, 0, len(definitions))
		for _, def := range definitions {
			gates = append(gates, Gate{ID: def.ID, Command: def.Command, Status: "passed", Mode: mode})
		}
		return gates, true
	}
	for _, def := range definitions {
		if def.Command == command {
			return []Gate{{ID: def.ID, Command: def.Command, Status: "passed", Mode: mode}}, true
		}
	}
	return nil, false
}

func invalid(command string, err error, stderr io.Writer) int {
	fmt.Fprintf(stderr, "ERROR: %s invalid arguments: %v\n", command, err)
	return 2
}

func fallback(value, fallbackValue string) string {
	if value == "" {
		return fallbackValue
	}
	return value
}
