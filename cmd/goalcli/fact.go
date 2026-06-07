// SPDX-License-Identifier: Apache-2.0
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/ZoneCNH/xlib-standard/internal/xlibfacts"
)

func runFact(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) == 0 {
		write(stderr, "usage: goalcli fact audit [--strict] [--root <path>] [--json]\n")
		return 2
	}
	switch args[0] {
	case "audit":
		return runFactAudit(args[1:], stdout, stderr)
	default:
		write(stderr, "ERROR: unknown fact subcommand %q\n", args[0])
		write(stderr, "usage: goalcli fact audit [--strict] [--root <path>] [--json]\n")
		return 2
	}
}

func runFactAudit(args []string, stdout io.Writer, stderr io.Writer) int {
	flags := flag.NewFlagSet("goalcli fact audit", flag.ContinueOnError)
	flags.SetOutput(stderr)
	strict := flags.Bool("strict", false, "verify local consumers of canonical facts")
	jsonOut := flags.Bool("json", false, "emit JSON report")
	root := flags.String("root", ".", "repository root")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	_ = jsonOut
	if flags.NArg() > 0 {
		write(stderr, "ERROR: fact audit invalid arguments: unexpected positional argument %q\n", flags.Arg(0))
		return 2
	}

	facts, err := xlibfacts.Load(*root)
	if err != nil {
		write(stderr, "ERROR: fact audit could not load %s: %v\n", xlibfacts.Path, err)
		return emitReport(stdout, "fact audit", "failed", nil, []string{fmt.Sprintf("missing or unreadable %s", xlibfacts.Path)})
	}

	gaps := facts.Validate()
	gaps = append(gaps, xlibfacts.DriftGaps(facts, xlibfacts.Expected())...)
	details := []string{
		"canonical_facts=" + xlibfacts.Path,
		"current_release.version=" + facts.CurrentRelease.Version,
		"current_release.commit=" + facts.CurrentRelease.Commit,
		"current_release.released_at=" + facts.CurrentRelease.ReleasedAt,
		"runtime.goal_runtime_version=" + facts.Runtime.GoalRuntimeVersion,
		"runtime.governance_runtime_version=" + facts.Runtime.GovernanceRuntimeVersion,
	}
	if *strict {
		details = append(details, "strict=true")
		gaps = append(gaps, factStrictProjectionGaps(*root)...)
	}
	if len(gaps) > 0 {
		write(stderr, "ERROR: fact audit found %d gap(s)\n", len(gaps))
		return emitReport(stdout, "fact audit", "failed", details, gaps)
	}
	return emitReport(stdout, "fact audit", "passed", details, nil)
}

func factStrictProjectionGaps(root string) []string {
	checks := []struct {
		path    string
		needles []string
	}{
		{path: "cmd/goalcli/governance.go", needles: []string{"xlibfacts.CurrentReleaseVersion", "xlibfacts.GovernanceRuntimeVersion", "\"fact\""}},
		{path: "internal/tools/releasemanifest/main.go", needles: []string{"xlibfacts.CurrentReleaseVersion"}},
		{path: ".agent/harness/harness.yaml", needles: []string{"release-preflight VERSION=" + xlibfacts.CurrentReleaseVersion}},
		{path: ".agent/release/release-required-gates.yaml", needles: []string{"release-preflight VERSION=" + xlibfacts.CurrentReleaseVersion}},
		{path: ".agent/registries/makefile-baseline.yaml", needles: []string{"fact-audit: \"$(GOALCLI) fact audit --strict\""}},
		{path: ".agent/registries/makefile-target-registry.yaml", needles: []string{"- fact-audit"}},
		{path: ".agent/registries/command-registry.yaml", needles: []string{"name: fact"}},
		{path: "docs/standard/goalcli-cli-contract.md", needles: []string{"fact audit"}},
	}
	var gaps []string
	for _, check := range checks {
		content, err := os.ReadFile(filepath.Join(root, check.path))
		if err != nil {
			gaps = append(gaps, "missing "+check.path)
			continue
		}
		text := string(content)
		for _, needle := range check.needles {
			if !strings.Contains(text, needle) {
				gaps = append(gaps, check.path+" missing "+needle)
			}
		}
	}
	makefile, err := os.ReadFile(filepath.Join(root, "Makefile"))
	if err != nil {
		return append(gaps, "missing Makefile")
	}
	makefileText := string(makefile)
	factBlock := makefileTargetBlock(makefileText, "fact-audit")
	if factBlock == "" {
		gaps = append(gaps, "Makefile missing fact-audit target")
	} else if !strings.Contains(factBlock, "$(GOALCLI) fact audit --strict") {
		gaps = append(gaps, "Makefile fact-audit must run $(GOALCLI) fact audit --strict")
	}
	for _, target := range []string{"context-release", "release-check", "release-check-extended"} {
		if !makefileDependencyHasToken(makefileTargetDependencies(makefileText, target), "fact-audit") {
			gaps = append(gaps, "Makefile "+target+" must depend on fact-audit")
		}
	}
	return gaps
}
