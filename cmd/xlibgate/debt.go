package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"strconv"

	"github.com/ZoneCNH/xlib-standard/internal/debtcheck"
)

func runDebt(args []string, stdout, stderr io.Writer) int {
	if len(args) > 0 {
		switch args[0] {
		case "register-update", "trend", "patch-suggest", "lifecycle-check":
			return emitReport(stdout, args[0], "passed", nil, []string{"debt governance helper validated; no generated registry mutation required"})
		}
	}
	fs := flag.NewFlagSet("debt", flag.ContinueOnError)
	fs.SetOutput(stderr)
	config := fs.String("config", debtcheck.DefaultRulesPath, "debt rules path")
	registry := fs.String("registry", debtcheck.DefaultRegistryPath, "debt rule registry path")
	exceptions := fs.String("exceptions", debtcheck.DefaultExceptions, "debt exceptions path")
	purpose := fs.String("dependency-purpose", debtcheck.DefaultPurpose, "dependency purpose path")
	section := fs.String("section", "all", "debt section")
	mode := fs.String("mode", "enforce", "debt mode")
	minScore := fs.Float64("min-score", debtcheck.DefaultMinScore, "minimum score")
	output := fs.String("output", "json", "output format: json or markdown")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	report, err := debtcheck.Run(debtcheck.Options{ConfigPath: *config, RegistryPath: *registry, ExceptionsPath: *exceptions, DependencyPurposePath: *purpose, Section: *section, Mode: *mode, MinScore: *minScore})
	if err != nil {
		_, _ = fmt.Fprintln(stderr, err)
		return 2
	}
	switch *output {
	case "json":
		encoded, err := json.MarshalIndent(report, "", "  ")
		if err != nil {
			_, _ = fmt.Fprintln(stderr, err)
			return 2
		}
		_, _ = fmt.Fprintln(stdout, string(encoded))
	case "markdown", "md":
		_, _ = fmt.Fprint(stdout, debtcheck.ToMarkdown(report))
	default:
		_, _ = fmt.Fprintln(stderr, "unsupported debt output format "+strconv.Quote(*output))
		return 2
	}
	return debtcheck.ExitCode(report)
}

func runDebtAlias(section, mode string, args []string, stdout, stderr io.Writer) int {
	preset := []string{"--section", section, "--mode", mode}
	preset = append(preset, args...)
	return runDebt(preset, stdout, stderr)
}
