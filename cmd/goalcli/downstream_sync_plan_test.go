package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunDownstreamSyncPlanWritesRequiredPlan(t *testing.T) {
	t.Parallel()
	dir := localDownstreamSyncPlanTestDir(t)
	impactReport := filepath.Join(dir, "impact.md")
	output := filepath.Join(dir, "release", "downstream-sync", "latest.md")
	workspaceRoot := filepath.Join(dir, "workspace")
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
	dir := localDownstreamSyncPlanTestDir(t)
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
	dir := localDownstreamSyncPlanTestDir(t)
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
	if plan.AdoptionClaim != "not_claimed" {
		t.Fatalf("AdoptionClaim=%q, want not_claimed", plan.AdoptionClaim)
	}
	wantReview := downstreamSyncConsumerReview{
		Name:   "x.go",
		Role:   "consumer_only",
		Action: "consumer_only_review_required",
		Status: "review_pending_no_standard_write",
	}
	if plan.ConsumerReview != wantReview {
		t.Fatalf("ConsumerReview=%+v, want %+v", plan.ConsumerReview, wantReview)
	}
	for _, source := range []string{".agent/registries/downstream-adoption-status.yaml", ".agent/evidence/truth-state.yaml"} {
		if !containsString(plan.TruthSources, source) {
			t.Fatalf("TruthSources=%v, want %q", plan.TruthSources, source)
		}
	}
	for _, interpretation := range []string{"registered != adopted", "baseline_scanned != adopted", "patch_only != proof_based_adoption", "not_run != passed"} {
		if !containsString(plan.ForbiddenInterpretations, interpretation) {
			t.Fatalf("ForbiddenInterpretations=%v, want %q", plan.ForbiddenInterpretations, interpretation)
		}
	}
	for _, target := range plan.Targets {
		combined := strings.ToLower(strings.Join([]string{target.Name, target.ModulePath, target.PackageName}, " "))
		if strings.Contains(combined, "x.go") {
			t.Fatalf("target %+v includes x.go; want x.go consumer review only", target)
		}
		for _, command := range target.Commands {
			lowered := strings.ToLower(command)
			if strings.Contains(lowered, "x.go") {
				t.Fatalf("command %q includes x.go; want x.go consumer review only", command)
			}
			if strings.Contains(command, ".agent/registries/downstream-adoption-status.yaml") || strings.Contains(command, ".agent/evidence/truth-state.yaml") {
				t.Fatalf("command %q writes adoption truth source", command)
			}
		}
	}
}

func TestRunDownstreamSyncPlanReportsMarshalError(t *testing.T) {
	dir := localDownstreamSyncPlanTestDir(t)
	impactReport := filepath.Join(dir, "impact.md")
	if err := os.WriteFile(impactReport, []byte(requiredDownstreamImpactReportFixture()), 0o644); err != nil {
		t.Fatalf("write impact report fixture: %v", err)
	}

	old := downstreamSyncPlanMarshalIndent
	downstreamSyncPlanMarshalIndent = func(any, string, string) ([]byte, error) {
		return nil, errors.New("marshal failed")
	}
	t.Cleanup(func() { downstreamSyncPlanMarshalIndent = old })

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runDownstreamSyncPlan([]string{"--impact-report", impactReport, "--output", "-", "--format", "json"}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("runDownstreamSyncPlan returned %d, want 1; stderr=%s stdout=%s", code, stderr.String(), stdout.String())
	}
	if !strings.Contains(stderr.String(), "marshal failed") {
		t.Fatalf("stderr missing marshal error: %s", stderr.String())
	}
	if !strings.Contains(stdout.String(), `"command": "downstream-sync-plan"`) || !strings.Contains(stdout.String(), `"status": "failed"`) {
		t.Fatalf("stdout missing failed report: %s", stdout.String())
	}
}

func TestRunDownstreamSyncPlanReportsOutputWriteErrors(t *testing.T) {
	dir := localDownstreamSyncPlanTestDir(t)
	impactReport := filepath.Join(dir, "impact.md")
	if err := os.WriteFile(impactReport, []byte(requiredDownstreamImpactReportFixture()), 0o644); err != nil {
		t.Fatalf("write impact report fixture: %v", err)
	}
	tests := []struct {
		name  string
		patch func(*testing.T)
		want  string
	}{
		{name: "mkdir", patch: func(t *testing.T) {
			old := downstreamSyncPlanMkdirAll
			downstreamSyncPlanMkdirAll = func(string, os.FileMode) error {
				return errors.New("mkdir failed")
			}
			t.Cleanup(func() { downstreamSyncPlanMkdirAll = old })
		}, want: "mkdir failed"},
		{name: "write", patch: func(t *testing.T) {
			old := downstreamSyncPlanWriteFile
			downstreamSyncPlanWriteFile = func(string, []byte, os.FileMode) error {
				return errors.New("write failed")
			}
			t.Cleanup(func() { downstreamSyncPlanWriteFile = old })
		}, want: "write failed"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.patch(t)
			var stdout bytes.Buffer
			var stderr bytes.Buffer
			output := filepath.Join(dir, tt.name, "latest.md")
			code := runDownstreamSyncPlan([]string{"--impact-report", impactReport, "--output", output}, &stdout, &stderr)
			if code != 1 {
				t.Fatalf("runDownstreamSyncPlan returned %d, want 1; stderr=%s stdout=%s", code, stderr.String(), stdout.String())
			}
			if !strings.Contains(stderr.String(), tt.want) {
				t.Fatalf("stderr = %s; want %q", stderr.String(), tt.want)
			}
			if !strings.Contains(stdout.String(), `"status": "failed"`) {
				t.Fatalf("stdout missing failed report: %s", stdout.String())
			}
		})
	}
}

func TestRunDownstreamSyncPlanRendersMarkdownToStdout(t *testing.T) {
	t.Parallel()
	dir := localDownstreamSyncPlanTestDir(t)
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

func TestRunDownstreamSyncPlanRejectsUnsafeOutputPaths(t *testing.T) {
	t.Parallel()
	dir := localDownstreamSyncPlanTestDir(t)
	impactReport := filepath.Join(dir, "impact.md")
	if err := os.WriteFile(impactReport, []byte(requiredDownstreamImpactReportFixture()), 0o644); err != nil {
		t.Fatalf("write impact report fixture: %v", err)
	}
	workspaceRoot := filepath.Join(dir, "workspace")
	absoluteOutput := filepath.Join(t.TempDir(), "downstream-sync-plan.md")
	cases := []struct {
		name   string
		output string
	}{
		{name: "empty", output: ""},
		{name: "absolute", output: absoluteOutput},
		{name: "parentTraversal", output: "../downstream-sync-plan.md"},
		{name: "truthState", output: ".agent/evidence/truth-state.yaml"},
		{name: "adoptionStatus", output: ".agent/registries/downstream-adoption-status.yaml"},
		{name: "workspaceRoot", output: filepath.Join(workspaceRoot, "latest.md")},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := runDownstreamSyncPlan([]string{
				"--impact-report", impactReport,
				"--output", tc.output,
				"--workspace-root", workspaceRoot,
			}, &stdout, &stderr)
			if code != 1 {
				t.Fatalf("runDownstreamSyncPlan returned %d, want 1; stderr=%s stdout=%s", code, stderr.String(), stdout.String())
			}
			if !strings.Contains(stderr.String(), "invalid downstream sync plan output") {
				t.Fatalf("stderr missing invalid output message: %s", stderr.String())
			}
			if !strings.Contains(stdout.String(), `"status": "failed"`) {
				t.Fatalf("stdout missing failed report: %s", stdout.String())
			}
		})
	}
	if _, err := os.Stat(absoluteOutput); !os.IsNotExist(err) {
		t.Fatalf("absolute output should not be written, stat err = %v", err)
	}
}

func TestRunDownstreamSyncPlanRejectsInvalidArguments(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		args []string
		want int
	}{
		{name: "help", args: []string{"--help"}, want: 0},
		{name: "unknownFlag", args: []string{"--missing"}, want: 2},
		{name: "positional", args: []string{"unexpected"}, want: 2},
		{name: "invalidFormat", args: []string{"--format", "yaml"}, want: 2},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := runDownstreamSyncPlan(tc.args, &stdout, &stderr)
			if code != tc.want {
				t.Fatalf("runDownstreamSyncPlan(%v) returned %d, want %d; stderr=%s stdout=%s", tc.args, code, tc.want, stderr.String(), stdout.String())
			}
		})
	}
}

func TestRunDownstreamSyncPlanRejectsDirectoryAndRelativeWorkspaceOutputs(t *testing.T) {
	t.Parallel()
	dir := localDownstreamSyncPlanTestDir(t)
	impactReport := filepath.Join(dir, "impact.md")
	if err := os.WriteFile(impactReport, []byte(requiredDownstreamImpactReportFixture()), 0o644); err != nil {
		t.Fatalf("write impact report fixture: %v", err)
	}
	outputDir := filepath.Join(dir, "plans")
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		t.Fatalf("create output dir: %v", err)
	}
	workspaceRoot := filepath.Join(dir, "workspace")

	for _, output := range []string{".", outputDir, filepath.Join(workspaceRoot, "latest.md")} {
		t.Run(output, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := runDownstreamSyncPlan([]string{
				"--impact-report", impactReport,
				"--output", output,
				"--workspace-root", workspaceRoot,
			}, &stdout, &stderr)
			if code != 1 {
				t.Fatalf("runDownstreamSyncPlan returned %d, want 1; stderr=%s stdout=%s", code, stderr.String(), stdout.String())
			}
			if !strings.Contains(stderr.String(), "invalid downstream sync plan output") {
				t.Fatalf("stderr missing invalid output message: %s", stderr.String())
			}
			if !strings.Contains(stdout.String(), `"status": "failed"`) {
				t.Fatalf("stdout missing failed report: %s", stdout.String())
			}
		})
	}
}

func TestValidateDownstreamSyncPlanOutputPathReportsStatErrors(t *testing.T) {
	root := t.TempDir()
	chdir(t, root)
	if err := os.WriteFile("blocked", []byte("not a directory"), 0o644); err != nil {
		t.Fatalf("write blocking file: %v", err)
	}

	err := validateDownstreamSyncPlanOutputPath("blocked/latest.md", "workspace")
	if err == nil || !strings.Contains(err.Error(), "check output path") {
		t.Fatalf("validateDownstreamSyncPlanOutputPath() error = %v, want check output path", err)
	}
}

func TestRunDownstreamSyncPlanReportsOutputPathStatError(t *testing.T) {
	root := t.TempDir()
	impactReport := filepath.Join(root, "impact.md")
	if err := os.WriteFile(impactReport, []byte(requiredDownstreamImpactReportFixture()), 0o644); err != nil {
		t.Fatalf("write impact report fixture: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "blocked"), []byte("not a directory"), 0o644); err != nil {
		t.Fatalf("write blocking file: %v", err)
	}
	chdir(t, root)

	var stdout, stderr bytes.Buffer
	code := runDownstreamSyncPlan([]string{
		"--impact-report", impactReport,
		"--output", "blocked/latest.md",
		"--workspace-root", "workspace",
	}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("runDownstreamSyncPlan returned %d, want 1; stderr=%s stdout=%s", code, stderr.String(), stdout.String())
	}
	if !strings.Contains(stderr.String(), "invalid downstream sync plan output") ||
		!strings.Contains(stderr.String(), "check output path") {
		t.Fatalf("stderr missing output stat error: %s", stderr.String())
	}
	if !strings.Contains(stdout.String(), `"status": "failed"`) {
		t.Fatalf("stdout missing failed report: %s", stdout.String())
	}
}

func TestParseDownstreamImpactReportRejectsInvalidMetadata(t *testing.T) {
	t.Parallel()
	required := requiredDownstreamImpactReportFixture()
	notRequired := notRequiredDownstreamImpactReportFixture()
	cases := []struct {
		name    string
		report  string
		wantErr string
	}{
		{
			name:    "missingBool",
			report:  strings.Replace(required, "- downstream_sync_required: `true`\n", "", 1),
			wantErr: "missing downstream_sync_required",
		},
		{
			name:    "invalidBool",
			report:  strings.Replace(required, "- downstream_sync_required: `true`", "- downstream_sync_required: `sometimes`", 1),
			wantErr: "invalid downstream_sync_required",
		},
		{
			name:    "requiredDecisionMismatch",
			report:  strings.Replace(required, "- downstream_release_decision: `required`", "- downstream_release_decision: `not_required`", 1),
			wantErr: "downstream_sync_required=true requires",
		},
		{
			name:    "notRequiredDecisionMismatch",
			report:  strings.Replace(notRequired, "- downstream_release_decision: `not_required`", "- downstream_release_decision: `required`", 1),
			wantErr: "downstream_sync_required=false requires",
		},
		{
			name:    "invalidDownstreamDecision",
			report:  strings.Replace(required, "- downstream_release_decision: `required`", "- downstream_release_decision: `maybe`", 1),
			wantErr: "invalid downstream_release_decision",
		},
		{
			name:    "invalidRepositoryRulesDecision",
			report:  strings.Replace(required, "- repository_rules_release_decision: `audit_required`", "- repository_rules_release_decision: `maybe`", 1),
			wantErr: "invalid repository_rules_release_decision",
		},
		{
			name:    "missingPrimary",
			report:  strings.Replace(required, "- primary_downstream: `github.com/ZoneCNH/kernel`\n", "", 1),
			wantErr: "missing primary_downstream",
		},
		{
			name:    "missingChangedCount",
			report:  strings.Replace(required, "- changed_file_count: `3`\n", "", 1),
			wantErr: "missing changed_file_count",
		},
		{
			name:    "invalidChangedCount",
			report:  strings.Replace(required, "- changed_file_count: `3`", "- changed_file_count: `many`", 1),
			wantErr: "invalid changed_file_count",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "impact.md")
			if err := os.WriteFile(path, []byte(tc.report), 0o644); err != nil {
				t.Fatalf("write impact report: %v", err)
			}
			_, err := parseDownstreamImpactReport(path)
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("parseDownstreamImpactReport() error = %v, want containing %q", err, tc.wantErr)
			}
		})
	}
}

func TestDownstreamSyncPlanPureHelpers(t *testing.T) {
	t.Parallel()
	if !isWithinDownstreamWorkspace("workspace/latest.md", "workspace") {
		t.Fatalf("relative file under workspace should be detected")
	}
	if !isWithinDownstreamWorkspace("workspace", "workspace") {
		t.Fatalf("workspace root should be detected")
	}
	for _, tc := range []struct {
		output    string
		workspace string
	}{
		{output: "workspace/latest.md"},
		{output: "workspace/latest.md", workspace: "."},
		{output: "workspace/latest.md", workspace: "./"},
		{output: "workspace/latest.md", workspace: filepath.Join(t.TempDir(), "workspace")},
		{output: "workspace/latest.md", workspace: "../workspace"},
	} {
		if isWithinDownstreamWorkspace(tc.output, tc.workspace) {
			t.Fatalf("isWithinDownstreamWorkspace(%q, %q) = true, want false", tc.output, tc.workspace)
		}
	}
	for input, want := range map[string]string{
		"":            "''",
		"safe/path-1": "safe/path-1",
		"needs space": "'needs space'",
		"it's":        "'it'\"'\"'s'",
	} {
		if got := shellQuote(input); got != want {
			t.Fatalf("shellQuote(%q) = %q, want %q", input, got, want)
		}
	}
	categories := sortedDownstreamImpactCategories(map[string]int{"zz": 1, "aa": 2, "docs": 3})
	if got, want := strings.Join(categories, ","), "docs,aa,zz"; got != want {
		t.Fatalf("sortedDownstreamImpactCategories() = %q, want %q", got, want)
	}
}

func TestRunDispatchesDownstreamSyncPlan(t *testing.T) {
	t.Parallel()
	dir := localDownstreamSyncPlanTestDir(t)
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

func localDownstreamSyncPlanTestDir(t *testing.T) string {
	t.Helper()
	name := strings.NewReplacer("/", "_", "\\", "_", " ", "_", ":", "_").Replace(t.Name())
	dir := filepath.Join(".downstream-sync-plan-test", name)
	if err := os.RemoveAll(dir); err != nil {
		t.Fatalf("remove stale downstream sync plan test dir: %v", err)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("create local test dir %s: %v", dir, err)
	}
	t.Cleanup(func() {
		if err := os.RemoveAll(dir); err != nil {
			t.Errorf("remove downstream sync plan test dir: %v", err)
		}
	})
	return dir
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

- ` + "`" + `.agent/harness/harness.yaml` + "`" + `

## governance_registry

- ` + "`" + `.agent/registries/command-registry.yaml` + "`" + `

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
