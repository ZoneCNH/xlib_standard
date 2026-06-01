package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type command struct {
	description string
	run         func([]string) error
}

var commands = map[string]command{
	"docs-check":                      {"run documentation contract checks", script("./scripts/check_docs.sh")},
	"boundary":                        {"run module and dependency boundary checks", script("./scripts/check_boundary.sh")},
	"secrets":                         {"run secret scan", script("./scripts/check_secrets.sh")},
	"contracts":                       {"run public contract checks", script("./scripts/check_contracts.sh")},
	"integration":                     {"run rendered downstream integration checks", script("./scripts/run_integration.sh")},
	"render":                          {"render the template into a downstream module", script("./scripts/render_template.sh")},
	"render-check":                    {"verify a rendered downstream module", script("./scripts/check_rendered_template.sh")},
	"evidence":                        {"generate release evidence manifest", goRun("./internal/tools/releasemanifest", "--out", "release/manifest/latest.json")},
	"release-evidence-hash":           {"write release evidence checksum", script("./scripts/hash_release_evidence.sh")},
	"release-evidence-check":          {"verify release evidence manifest", script("./scripts/check_release_evidence.sh")},
	"release-evidence-checksum-check": {"verify release evidence checksum", script("./scripts/hash_release_evidence.sh", "--check")},
	"release-check":                   {"run the release-check make gate", makeTarget("release-check")},
	"release-final-check":             {"run the release-final-check make gate", makeTarget("release-final-check")},
	"score":                           {"score Full Goal Runtime v3.1 gate readiness", runScore},
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "xlibgate:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		usage()
		return nil
	}

	cmd, ok := commands[args[0]]
	if !ok {
		usage()
		return fmt.Errorf("unknown command %q", args[0])
	}
	return cmd.run(args[1:])
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: xlibgate <command> [args]")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "commands:")
	for _, name := range []string{
		"docs-check", "boundary", "secrets", "contracts", "integration",
		"render", "render-check", "evidence", "release-evidence-hash",
		"release-evidence-check", "release-evidence-checksum-check",
		"release-check", "release-final-check", "score",
	} {
		fmt.Fprintf(os.Stderr, "  %-32s %s\n", name, commands[name].description)
	}
}

func script(path string, fixedArgs ...string) func([]string) error {
	return func(args []string) error {
		return runCommand(path, append(append([]string{}, fixedArgs...), args...)...)
	}
}

func goRun(pkg string, fixedArgs ...string) func([]string) error {
	return func(args []string) error {
		goArgs := append([]string{"run", pkg}, fixedArgs...)
		goArgs = append(goArgs, args...)
		return runCommand("go", goArgs...)
	}
}

func makeTarget(target string) func([]string) error {
	return func(args []string) error {
		if len(args) > 0 {
			return fmt.Errorf("%s does not accept extra args: %s", target, strings.Join(args, " "))
		}
		return runCommand("make", target)
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
		return false, err
	}
	return strings.Contains(string(data), needle), nil
}
