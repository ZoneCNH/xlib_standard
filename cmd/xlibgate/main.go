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
		write(stderr, "usage: xlibgate score [--min <score>]\n")
		return 2
	}
	switch args[0] {
	case "score":
		return runScore(args[1:], stdout, stderr)
	case "help", "-h", "--help":
		write(stdout, "usage: xlibgate score [--min <score>]\n")
		return 0
	default:
		write(stderr, "unknown command %q\n", args[0])
		return 2
	}
}

func runScore(args []string, stdout io.Writer, stderr io.Writer) int {
	flags := flag.NewFlagSet("xlibgate score", flag.ContinueOnError)
	flags.SetOutput(stderr)
	minimum := flags.Float64("min", releasequality.DefaultMinimum, "minimum acceptable release score")
	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}
	report := releasequality.Compute(*minimum)
	data, err := releasequality.Marshal(report)
	if err != nil {
		write(stderr, "ERROR: %v\n", err)
		return 1
	}
	write(stdout, "%s\n", data)
	if err := releasequality.Verify(report, *minimum); err != nil {
		write(stderr, "ERROR: %v\n", err)
		return 1
	}
	return 0
}

func write(writer io.Writer, format string, args ...any) {
	_, _ = fmt.Fprintf(writer, format, args...)
}
