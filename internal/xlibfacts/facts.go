// SPDX-License-Identifier: Apache-2.0
package xlibfacts

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	Path          = ".xlib/facts/xlib.yaml"
	SchemaVersion = "xlib-facts/v1"
	Module        = "github.com/ZoneCNH/xlib-standard"

	CurrentReleaseVersion    = "v0.4.15"
	CurrentReleaseCommit     = "c0fc3813e156cf35a37ddd0033432a78943bb32b"
	CurrentReleaseReleasedAt = "2026-06-05T10:24:00Z"

	GoalRuntimeVersion       = "v3.1"
	GovernanceRuntimeVersion = "v2.9.3"

	GoVersion           = "1.23.0"
	GolangCILintVersion = "v2.1.6"
	GovulncheckVersion  = "v1.1.4"
)

type Facts struct {
	SchemaVersion  string
	Module         string
	CurrentRelease ReleaseFacts
	Runtime        RuntimeFacts
	Tools          ToolFacts
}

type ReleaseFacts struct {
	Version    string
	Commit     string
	ReleasedAt string
}

type RuntimeFacts struct {
	GoalRuntimeVersion       string
	GovernanceRuntimeVersion string
}

type ToolFacts struct {
	Go            string
	GolangCILint  string
	Govulncheck   string
}

func Expected() Facts {
	return Facts{
		SchemaVersion: SchemaVersion,
		Module:        Module,
		CurrentRelease: ReleaseFacts{
			Version:    CurrentReleaseVersion,
			Commit:     CurrentReleaseCommit,
			ReleasedAt: CurrentReleaseReleasedAt,
		},
		Runtime: RuntimeFacts{
			GoalRuntimeVersion:       GoalRuntimeVersion,
			GovernanceRuntimeVersion: GovernanceRuntimeVersion,
		},
		Tools: ToolFacts{
			Go:           GoVersion,
			GolangCILint: GolangCILintVersion,
			Govulncheck:  GovulncheckVersion,
		},
	}
}

func Load(root string) (Facts, error) {
	if root == "" {
		root = "."
	}
	data, err := os.ReadFile(filepath.Join(root, Path))
	if err != nil {
		return Facts{}, err
	}
	return Parse(data)
}

func Parse(data []byte) (Facts, error) {
	var facts Facts
	section := ""
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := strings.TrimRight(scanner.Text(), " \t")
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		indent := len(line) - len(strings.TrimLeft(line, " "))
		if !strings.Contains(trimmed, ":") {
			return Facts{}, fmt.Errorf("invalid facts line %q", trimmed)
		}
		parts := strings.SplitN(trimmed, ":", 2)
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if indent == 0 && value == "" {
			section = key
			continue
		}
		value = unquoteYAMLScalar(value)
		if indent == 0 {
			section = ""
			switch key {
			case "schema_version":
				facts.SchemaVersion = value
			case "module":
				facts.Module = value
			default:
				return Facts{}, fmt.Errorf("unknown top-level facts key %q", key)
			}
			continue
		}
		if indent != 2 {
			return Facts{}, fmt.Errorf("unsupported indentation for facts key %q", key)
		}
		switch section {
		case "current_release":
			switch key {
			case "version":
				facts.CurrentRelease.Version = value
			case "commit":
				facts.CurrentRelease.Commit = value
			case "released_at":
				facts.CurrentRelease.ReleasedAt = value
			default:
				return Facts{}, fmt.Errorf("unknown current_release facts key %q", key)
			}
		case "runtime":
			switch key {
			case "goal_runtime_version":
				facts.Runtime.GoalRuntimeVersion = value
			case "governance_runtime_version":
				facts.Runtime.GovernanceRuntimeVersion = value
			default:
				return Facts{}, fmt.Errorf("unknown runtime facts key %q", key)
			}
		case "tools":
			switch key {
			case "go":
				facts.Tools.Go = value
			case "golangci_lint":
				facts.Tools.GolangCILint = value
			case "govulncheck":
				facts.Tools.Govulncheck = value
			default:
				return Facts{}, fmt.Errorf("unknown tools facts key %q", key)
			}
		default:
			return Facts{}, fmt.Errorf("facts key %q is outside a supported section", key)
		}
	}
	if err := scanner.Err(); err != nil {
		return Facts{}, err
	}
	return facts, nil
}

func (facts Facts) Validate() []string {
	var gaps []string
	require := func(name, value string) {
		if strings.TrimSpace(value) == "" {
			gaps = append(gaps, "missing "+name)
		}
	}
	require("schema_version", facts.SchemaVersion)
	require("module", facts.Module)
	require("current_release.version", facts.CurrentRelease.Version)
	require("current_release.commit", facts.CurrentRelease.Commit)
	require("current_release.released_at", facts.CurrentRelease.ReleasedAt)
	require("runtime.goal_runtime_version", facts.Runtime.GoalRuntimeVersion)
	require("runtime.governance_runtime_version", facts.Runtime.GovernanceRuntimeVersion)
	require("tools.go", facts.Tools.Go)
	require("tools.golangci_lint", facts.Tools.GolangCILint)
	require("tools.govulncheck", facts.Tools.Govulncheck)
	if facts.CurrentRelease.ReleasedAt != "" {
		if _, err := time.Parse(time.RFC3339, facts.CurrentRelease.ReleasedAt); err != nil {
			gaps = append(gaps, "current_release.released_at must be RFC3339")
		}
	}
	return gaps
}

func DriftGaps(actual, expected Facts) []string {
	var gaps []string
	compare := func(name, got, want string) {
		if got != want {
			gaps = append(gaps, fmt.Sprintf("%s drift: got %q want %q", name, got, want))
		}
	}
	compare("schema_version", actual.SchemaVersion, expected.SchemaVersion)
	compare("module", actual.Module, expected.Module)
	compare("current_release.version", actual.CurrentRelease.Version, expected.CurrentRelease.Version)
	compare("current_release.commit", actual.CurrentRelease.Commit, expected.CurrentRelease.Commit)
	compare("current_release.released_at", actual.CurrentRelease.ReleasedAt, expected.CurrentRelease.ReleasedAt)
	compare("runtime.goal_runtime_version", actual.Runtime.GoalRuntimeVersion, expected.Runtime.GoalRuntimeVersion)
	compare("runtime.governance_runtime_version", actual.Runtime.GovernanceRuntimeVersion, expected.Runtime.GovernanceRuntimeVersion)
	compare("tools.go", actual.Tools.Go, expected.Tools.Go)
	compare("tools.golangci_lint", actual.Tools.GolangCILint, expected.Tools.GolangCILint)
	compare("tools.govulncheck", actual.Tools.Govulncheck, expected.Tools.Govulncheck)
	return gaps
}

func unquoteYAMLScalar(value string) string {
	value = strings.TrimSpace(value)
	if len(value) >= 2 {
		quote := value[0]
		if (quote == '\'' || quote == '"') && value[len(value)-1] == quote {
			return value[1 : len(value)-1]
		}
	}
	return value
}
