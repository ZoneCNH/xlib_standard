package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/ZoneCNH/xlib-standard/internal/debtcheck"
)

func runDebtCommand(command string, args []string, stdout io.Writer, stderr io.Writer) int {
	flags := flag.NewFlagSet("xlibgate "+command, flag.ContinueOnError)
	flags.SetOutput(stderr)
	jsonOut := flags.Bool("json", false, "print debt report as JSON")
	root := flags.String("root", ".", "repository root")
	out := flags.String("out", "release/debt/latest.json", "debt evidence JSON path")
	markdown := flags.String("markdown", "release/debt/latest.md", "debt evidence Markdown path")
	checksum := flags.String("checksum", "release/debt/latest.json.sha256", "debt evidence checksum path")
	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}

	scopes := scopesForDebtCommand(command)
	report, err := debtcheck.Run(debtcheck.Options{Root: *root, Scopes: scopes, RunExternal: true})
	if err != nil {
		write(stderr, "ERROR: %v\n", err)
		return 1
	}
	if command == "debt-evidence" {
		if err := debtcheck.WriteEvidence(report, debtcheck.Options{OutPath: *out, MarkdownPath: *markdown, ChecksumPath: *checksum}); err != nil {
			write(stderr, "ERROR: %v\n", err)
			return 1
		}
	}
	if *jsonOut {
		data, err := json.MarshalIndent(report, "", "  ")
		if err != nil {
			write(stderr, "ERROR: %v\n", err)
			return 1
		}
		write(stdout, "%s\n", data)
	} else if command == "debt-evidence" {
		write(stdout, "generated debt evidence: %s\n", *out)
	} else {
		write(stdout, "debt gate %s: %s (score %.1f/min %.1f)\n", command, report.Status, report.Score, report.MinScore)
	}
	if err := debtcheck.StatusError(report); err != nil {
		write(stderr, "ERROR: %v\n", err)
		return 1
	}
	return 0
}

func scopesForDebtCommand(command string) []string {
	if command == "debt" || command == "debt-evidence" {
		return nil
	}
	if strings.HasSuffix(command, "-debt") || command == "architecture" || command == "domain" || command == "docs-drift" {
		return []string{command}
	}
	return nil
}

func debtCommandHelp() string {
	return fmt.Sprintf("debt scopes: %s", strings.Join(debtcheck.DefaultScopes, ", "))
}
