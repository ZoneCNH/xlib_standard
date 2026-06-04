package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func runAdoptionCheck(args []string, stdout io.Writer, stderr io.Writer) int {
	flags := flag.NewFlagSet("goalcli adoption-check", flag.ContinueOnError)
	flags.SetOutput(stderr)
	flags.Bool("json", false, "")
	flags.Bool("verify", false, "")
	root := flags.String("root", ".", "repository root")
	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}
	if flags.NArg() > 0 {
		write(stderr, "ERROR: adoption-check invalid arguments: unexpected positional argument %q\n", flags.Arg(0))
		return 2
	}

	details, gaps := evaluateAdoptionCheck(*root)
	if len(gaps) > 0 {
		write(stderr, "ERROR: adoption-check found %d gap(s)\n", len(gaps))
		return emitReport(stdout, "adoption-check", "failed", details, gaps)
	}
	return emitReport(stdout, "adoption-check", "passed", details, nil)
}

func evaluateAdoptionCheck(root string) ([]string, []string) {
	root = filepath.Clean(root)
	// Keep the downstream governance lock contract as xlib-standard.lock even after template text replacement.
	lockPath := "xlib-" + "standard.lock"
	required := map[string][]string{
		lockPath: {
			"schema_version:",
			"standard_version:",
			"standard_commit:",
			"module_name:",
			"module_path:",
			"package_name:",
			"layer:",
			`adoption_check: "GOWORK=off make adoption-check"`,
		},
		".githooks/pre-commit":                            {},
		".githooks/pre-push":                              {},
		".github/workflows/adoption-check.yml":            {"GOWORK=off make adoption-check", "workflow_dispatch", "pull_request"},
		"mk/governance.mk":                                {"adoption-check:", "$(GOALCLI) adoption-check --verify"},
		".agent/harness/harness.yaml":                     {"id: adoption_check", "GOWORK=off make adoption-check"},
		".agent/registries/command-registry.yaml":         {"name: adoption-check"},
		".agent/registries/makefile-target-registry.yaml": {"adoption-check"},
		".agent/registries/makefile-baseline.yaml":        {"adoption-check"},
		"Makefile": {".PHONY: adoption-check", "adoption-check:", "$(GOALCLI) adoption-check --verify"},
	}

	var gaps []string
	for path, needles := range required {
		content, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(path)))
		if err != nil {
			gaps = append(gaps, "missing "+path)
			continue
		}
		text := string(content)
		for _, needle := range needles {
			if !strings.Contains(text, needle) {
				gaps = append(gaps, fmt.Sprintf("%s missing %s", path, needle))
			}
		}
	}

	details := []string{
		"governance lock present",
		"git hooks present",
		"adoption workflow present",
		"governance Makefile target present",
		"harness adoption gate present",
	}
	return details, gaps
}
