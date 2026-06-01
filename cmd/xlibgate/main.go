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
	report := releasequality.Compute(*min)
	data, err := releasequality.Marshal(report)
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
