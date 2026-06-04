package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunDownstreamSyncPlanWritesRequiredPlan(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	impactReport := filepath.Join(dir, "impact.md")
	output := filepath.Join(dir, "release", "downstream-sync", "latest.md")
	workspaceRoot := filepath.Join(dir, "workspace")
	if err := os.MkdirAll(filepath.Dir(output), 0o755); err != nil {
		t.Fatalf("create output dir: %v", err)
	}
	if err := os.WriteFile(impactReport, []byte(requiredDownstreamImpactReportFixture()), 0o644); err != nil {
		t.Fatalf("write impact report fixture: %v", err)
	}

	var stdout, stderr bytes.Buffer
	code := runDownstreamSyncPlan([]string{
		"--impact-report", impactReport,
		"--output", output,
		"--workspace-root", workspaceRoot,
	}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("runDownstreamSyncPlan returned %d; stderr=%s stdout=%s", code, stderr.String(), stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
	if !strings.Contains(stdout.String(), `"status": "passed"`) {
		t.Fatalf("stdout missing passed report: %s", stdout.String())
	}
	data, err := os.ReadFile(output)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	plan := string(data)
	for _, want := range []string{
		"downstream_sync_required: `true`",
		"adoption_claim: `not_claimed`",
		"| `kernel` | `L0` | `P0` | `primary_sync_required` | `blocked_pending_downstream_workspace` |",
		"scripts/render_template.sh --module-name kernel --module-path github.com/ZoneCNH/kernel --package-name kernel",
		"GOWORK=off go test ./...",
		"RELEASE_EVIDENCE_REQUIRE_PASSED=1 GOWORK=off make release-evidence-check",
		"`x.go`: `consumer_only_review_required` / `review_pending_no_standard_write`.",
		"This command must not modify downstream repositories or adoption truth files.",
	} {
		if !strings.Contains(plan, want) {
			t.Fatalf("plan missing %q:\n%s", want, plan)
		}
	}
}

func TestRunDownstreamSyncPlanWritesNotRequiredPlan(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	impactReport := filepath.Join(dir, "impact.md")
	output := filepath.Join(dir, "latest.md")
	if err := os.WriteFile(impactReport, []byte(notRequiredDownstreamImpactReportFixture()), 0o644); err != nil {
		t.Fatalf("write impact report fixture: %v", err)
	}

	var stdout, stderr bytes.Buffer
	code := runDownstreamSyncPlan([]string{"--impact-report", impactReport, "--output", output}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("runDownstreamSyncPlan returned %d; stderr=%s stdout=%s", code, stderr.String(), stdout.String())
	}
	data, err := os.ReadFile(output)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	plan := string(data)
	for _, want := range []string{
		"downstream_sync_required: `false`",
		"| `kernel` | `L0` | `P0` | `sync_not_required` | `not_required_by_standard_impact` |",
		"No downstream write commands are generated because standard impact does not require sync.",
		"`x.go`: `consumer_only_no_write` / `not_required_by_standard_impact`.",
	} {
		if !strings.Contains(plan, want) {
			t.Fatalf("plan missing %q:\n%s", want, plan)
		}
	}
	if strings.Contains(plan, "scripts/render_template.sh") {
		t.Fatalf("not-required plan should not contain render commands:\n%s", plan)
	}
}

func TestRunDownstreamSyncPlanReportsMissingImpactReport(t *testing.T) {
	t.Parallel()
	var stdout, stderr bytes.Buffer
	code := runDownstreamSyncPlan([]string{"--impact-report", filepath.Join(t.TempDir(), "missing.md")}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("runDownstreamSyncPlan returned %d, want 1", code)
	}
	if !strings.Contains(stderr.String(), "run GOWORK=off make standard-impact-check first") {
		t.Fatalf("stderr missing remediation: %s", stderr.String())
	}
	if !strings.Contains(stdout.String(), `"status": "failed"`) {
		t.Fatalf("stdout missing failed report: %s", stdout.String())
	}
}

func TestRunDownstreamSyncPlanRendersJSONToStdout(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	impactReport := filepath.Join(dir, "impact.md")
	if err := os.WriteFile(impactReport, []byte(requiredDownstreamImpactReportFixture()), 0o644); err != nil {
		t.Fatalf("write impact report fixture: %v", err)
	}

	var stdout, stderr bytes.Buffer
	code := runDownstreamSyncPlan([]string{
		"--impact-report", impactReport,
		"--output", "-",
		"--format", "json",
	}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("runDownstreamSyncPlan returned %d; stderr=%s stdout=%s", code, stderr.String(), stdout.String())
	}
	var plan downstreamSyncPlan
	if err := json.Unmarshal(stdout.Bytes(), &plan); err != nil {
		t.Fatalf("stdout is not plan JSON: %v\n%s", err, stdout.String())
	}
	if !plan.DownstreamSyncRequired {
		t.Fatalf("DownstreamSyncRequired=false, want true")
	}
	if len(plan.Targets) != 11 {
		t.Fatalf("len(Targets)=%d, want 11", len(plan.Targets))
	}
	if plan.ConsumerReview.Name != "x.go" {
		t.Fatalf("ConsumerReview.Name=%q, want x.go", plan.ConsumerReview.Name)
	}
}

func TestRunDownstreamSyncPlanRendersMarkdownToStdout(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	impactReport := filepath.Join(dir, "impact.md")
	if err := os.WriteFile(impactReport, []byte(requiredDownstreamImpactReportFixture()), 0o644); err != nil {
		t.Fatalf("write impact report fixture: %v", err)
	}

	var stdout, stderr bytes.Buffer
	code := runDownstreamSyncPlan([]string{
		"--impact-report", impactReport,
		"--output", "-",
	}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("runDownstreamSyncPlan returned %d; stderr=%s stdout=%s", code, stderr.String(), stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
	for _, want := range []string{
		"# Downstream Sync Plan",
		"downstream_sync_required: `true`",
		"| `kernel` | `L0` | `P0` | `primary_sync_required` | `blocked_pending_downstream_workspace` |",
	} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("stdout missing %q:\n%s", want, stdout.String())
		}
	}
	if strings.Contains(stdout.String(), `"status": "passed"`) {
		t.Fatalf("markdown stdout must not be replaced by status report:\n%s", stdout.String())
	}
}

func TestRunDispatchesDownstreamSyncPlan(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	impactReport := filepath.Join(dir, "impact.md")
	output := filepath.Join(dir, "latest.md")
	if err := os.WriteFile(impactReport, []byte(requiredDownstreamImpactReportFixture()), 0o644); err != nil {
		t.Fatalf("write impact report fixture: %v", err)
	}

	var stdout, stderr bytes.Buffer
	code := run([]string{"downstream-sync-plan", "--impact-report", impactReport, "--output", output}, strings.NewReader(""), &stdout, &stderr)
	if code != 0 {
		t.Fatalf("run returned %d; stderr=%s stdout=%s", code, stderr.String(), stdout.String())
	}
	if _, err := os.Stat(output); err != nil {
		t.Fatalf("output missing: %v", err)
	}
}

func requiredDownstreamImpactReportFixture() string {
	return `# Standard Impact Report

- generated_at: ` + "`" + `2026-06-04T08:00:00Z` + "`" + `
- downstream_sync_required: ` + "`" + `true` + "`" + `
- context_runtime_change: ` + "`" + `true` + "`" + `
- governance_registry_change: ` + "`" + `true` + "`" + `
- downstream_release_decision: ` + "`" + `required` + "`" + `
- repository_rules_release_decision: ` + "`" + `audit_required` + "`" + `
- primary_downstream: ` + "`" + `github.com/ZoneCNH/kernel` + "`" + `
- changed_file_count: ` + "`" + `3` + "`" + `

## Downstream

- ` + "`" + `github.com/ZoneCNH/kernel` + "`" + `
- ` + "`" + `github.com/ZoneCNH/configx` + "`" + `

## contracts

- ` + "`" + `contracts/example.schema.json` + "`" + `

## context_runtime

- ` + "`" + `.agent/harness.yaml` + "`" + `

## governance_registry

- ` + "`" + `.agent/command-registry.yaml` + "`" + `

## harness

- 无变化

## repository_rules

- 无变化

## generator

- 无变化

## downstream_context

- 无变化

## evidence

- 无变化

## docs

- 无变化

## other

- 无变化

## Sync Decision

- ` + "`" + `required` + "`" + `
- 原因：检测到标准契约变更。
`
}

func notRequiredDownstreamImpactReportFixture() string {
	return `# Standard Impact Report

- generated_at: ` + "`" + `2026-06-04T08:00:00Z` + "`" + `
- downstream_sync_required: ` + "`" + `false` + "`" + `
- context_runtime_change: ` + "`" + `false` + "`" + `
- governance_registry_change: ` + "`" + `false` + "`" + `
- downstream_release_decision: ` + "`" + `not_required` + "`" + `
- repository_rules_release_decision: ` + "`" + `not_required` + "`" + `
- primary_downstream: ` + "`" + `github.com/ZoneCNH/kernel` + "`" + `
- changed_file_count: ` + "`" + `1` + "`" + `

## Downstream

- ` + "`" + `github.com/ZoneCNH/kernel` + "`" + `

## docs

- ` + "`" + `docs/README.md` + "`" + `

## contracts

- 无变化

## context_runtime

- 无变化

## governance_registry

- 无变化

## harness

- 无变化

## repository_rules

- 无变化

## generator

- 无变化

## downstream_context

- 无变化

## evidence

- 无变化

## other

- 无变化

## Sync Decision

- ` + "`" + `not_required` + "`" + `
- 原因：仅文档变更。
`
}
