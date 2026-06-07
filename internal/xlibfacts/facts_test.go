// SPDX-License-Identifier: Apache-2.0
package xlibfacts

import (
	"strings"
	"testing"
)

func TestParseExpectedFacts(t *testing.T) {
	facts, err := Parse([]byte(`schema_version: xlib-facts/v1
module: github.com/ZoneCNH/xlib-standard
current_release:
  version: v0.6.0
  commit: cc184e85991e53c79b3d94b05ad024f61583ccb1
  released_at: 2026-06-07T01:27:01Z
runtime:
  goal_runtime_version: v3.1
  governance_runtime_version: v2.9.3
tools:
  go: "1.23.0"
  golangci_lint: "v2.1.6"
  govulncheck: "v1.1.4"
`))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if gaps := facts.Validate(); len(gaps) > 0 {
		t.Fatalf("Validate() gaps = %v", gaps)
	}
	if gaps := DriftGaps(facts, Expected()); len(gaps) > 0 {
		t.Fatalf("DriftGaps() = %v", gaps)
	}
}

func TestDriftGapsReportsReleaseVersion(t *testing.T) {
	facts := Expected()
	facts.CurrentRelease.Version = "v0.4.14"

	gaps := DriftGaps(facts, Expected())
	if len(gaps) != 1 || !strings.Contains(gaps[0], "current_release.version") {
		t.Fatalf("DriftGaps() = %v; want current_release.version drift", gaps)
	}
}
