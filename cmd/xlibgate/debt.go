package main

import (
	"encoding/json"
	"errors"
	"flag"
	"io"

	"github.com/ZoneCNH/xlib-standard/internal/debtcheck"
)

func runDebt(args []string, stdout io.Writer, stderr io.Writer) int {
	flags := flag.NewFlagSet("xlibgate debt", flag.ContinueOnError)
	flags.SetOutput(stderr)
	flags.Bool("json", false, "emit JSON debt report")
	evidence := flags.Bool("evidence", false, "write release debt evidence artifacts")
	outDir := flags.String("out-dir", "release/debt", "debt evidence output directory")
	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}
	if flags.NArg() > 0 {
		write(stderr, "ERROR: debt invalid arguments: unexpected positional argument %q\n", flags.Arg(0))
		return 2
	}
	report := debtcheck.Evaluate(".", "debt")
	if *evidence {
		if err := debtcheck.WriteEvidence(".", *outDir, report); err != nil {
			write(stderr, "ERROR: %v\n", err)
			return 1
		}
	}
	if err := writeDebtReport(stdout, report); err != nil {
		write(stderr, "ERROR: %v\n", err)
		return 1
	}
	if report.Status == "passed" {
		return 0
	}
	return 1
}

func runDebtGate(command string, args []string, stdout io.Writer, stderr io.Writer) int {
	flags := flag.NewFlagSet("xlibgate "+command, flag.ContinueOnError)
	flags.SetOutput(stderr)
	flags.Bool("json", false, "emit JSON debt report")
	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}
	if flags.NArg() > 0 {
		write(stderr, "ERROR: %s invalid arguments: unexpected positional argument %q\n", command, flags.Arg(0))
		return 2
	}
	if !debtcheck.IsGate(command) {
		write(stderr, "ERROR: unknown debt gate %q\n", command)
		return 2
	}
	report := debtcheck.Evaluate(".", command)
	if err := writeDebtReport(stdout, report); err != nil {
		write(stderr, "ERROR: %v\n", err)
		return 1
	}
	if report.Status == "passed" {
		return 0
	}
	return 1
}

func writeDebtReport(stdout io.Writer, report debtcheck.Report) error {
	encoder := json.NewEncoder(stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(report)
}
