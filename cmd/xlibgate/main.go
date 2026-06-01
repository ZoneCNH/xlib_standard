package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/ZoneCNH/baselib-template/internal/releasequality"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "usage: xlibgate score [--min <score>]")
		return 2
	}
	switch args[0] {
	case "score":
		return runScore(args[1:], stdout, stderr)
	case "help", "-h", "--help":
		fmt.Fprintln(stdout, "usage: xlibgate score [--min <score>]")
		return 0
	default:
		fmt.Fprintf(stderr, "unknown command %q\n", args[0])
		return 2
	}
}

func runScore(args []string, stdout io.Writer, stderr io.Writer) int {
	flags := flag.NewFlagSet("xlibgate score", flag.ContinueOnError)
	flags.SetOutput(stderr)
	min := flags.Float64("min", releasequality.DefaultMinimum, "minimum acceptable release score")
	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}
}

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = os.Environ()
	return cmd.Run()
}

type scoreCheck struct {
	name string
	file string
	text string
}

var scoreChecks = []scoreCheck{
	{"xlibgate command surface", "cmd/xlibgate/main.go", "docs-check"},
	{"xlibgate boundary command", "cmd/xlibgate/main.go", "boundary"},
	{"xlibgate contracts command", "cmd/xlibgate/main.go", "contracts"},
	{"xlibgate render command", "cmd/xlibgate/main.go", "render-check"},
	{"xlibgate release command", "cmd/xlibgate/main.go", "release-final-check"},
	{"xlibgate score command", "cmd/xlibgate/main.go", "score"},
	{"Makefile xlibgate variable", "Makefile", "XLIBGATE ?= go run ./cmd/xlibgate"},
	{"Makefile docs gate", "Makefile", "$(XLIBGATE) docs-check"},
	{"Makefile boundary gate", "Makefile", "$(XLIBGATE) boundary"},
	{"Makefile contract gate", "Makefile", "$(XLIBGATE) contracts"},
	{"Makefile integration gate", "Makefile", "$(XLIBGATE) integration"},
	{"Makefile score target", "Makefile", "$(XLIBGATE) score --min 9.8"},
	{"CI score gate", ".github/workflows/ci.yml", "go run ./cmd/xlibgate score --min 9.8"},
	{"release workflow score gate", ".github/workflows/release.yml", "go run ./cmd/xlibgate score --min 9.8"},
	{"docs-check enforces xlibgate", "scripts/check_docs.sh", "cmd/xlibgate/main.go"},
	{"kernel downstream integration", "scripts/run_integration.sh", "github.com/ZoneCNH/kernel"},
	{"release docs mention xlibgate", "docs/release.md", "cmd/xlibgate"},
	{"harness docs mention xlibgate", "docs/standard/harness-gates.md", "cmd/xlibgate"},
	{"downstream docs mention kernel", "docs/standard/downstream-compatibility.md", "kernel"},
	{"contracts still verified", "scripts/check_contracts.sh", "go test ./contracts"},
}

func runScore(args []string) error {
	fs := flag.NewFlagSet("score", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	minimumScore := fs.Float64("min", 0, "minimum required score")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 0 {
		return fmt.Errorf("score received unexpected args: %s", strings.Join(fs.Args(), " "))
	}

	passed := 0
	failed := make([]string, 0)
	for _, check := range scoreChecks {
		ok, err := fileContains(check.file, check.text)
		if err != nil || !ok {
			msg := check.name
			if err != nil && !errors.Is(err, os.ErrNotExist) {
				msg += ": " + err.Error()
			}
			failed = append(failed, msg)
			continue
		}
		passed++
	}

	score := 10 * float64(passed) / float64(len(scoreChecks))
	fmt.Printf("xlibgate score: %.2f/10 (%d/%d checks passed)\n", score, passed, len(scoreChecks))
	if len(failed) > 0 {
		fmt.Println("failed checks:")
		for _, name := range failed {
			fmt.Println("- " + name)
		}
	}
	if score+1e-9 < *minimumScore {
		return fmt.Errorf("score %.2f is below minimum %s", score, strconv.FormatFloat(*minimumScore, 'f', -1, 64))
	}
	return nil
}

func fileContains(path, needle string) (bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(stderr, "ERROR: %v\n", err)
		return 1
	}
	fmt.Fprintln(stdout, string(data))
	if err := releasequality.Verify(report, *min); err != nil {
		fmt.Fprintf(stderr, "ERROR: %v\n", err)
		return 1
	}
	return 0
}
