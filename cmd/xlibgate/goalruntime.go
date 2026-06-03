package main

import (
	"io"

	"github.com/ZoneCNH/xlib-standard/internal/goalruntime"
)

func runGoalRuntime(command string, args []string, stdout io.Writer, stderr io.Writer) int {
	return goalruntime.Run(command, args, stdout, stderr)
}

func goalkitRuntimeTargets() []string {
	return goalruntime.Commands()
}
