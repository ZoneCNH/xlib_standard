package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestSchemaValidateFixtureAcceptsValidAndRejectsInvalid(t *testing.T) {
	fixture := writeSchemaCheckFixture(t, false)
	report := filepath.Join(t.TempDir(), "schema-check.json")
	var stdout, stderr bytes.Buffer

	got := run([]string{"schema", "validate", "--fixture", fixture, "--report", report}, strings.NewReader(""), &stdout, &stderr)

	if got != 0 {
		t.Fatalf("schema validate fixture exit = %d, stderr %q", got, stderr.String())
	}
	for _, needle := range []string{`"schema_version": "schema-check/v1"`, `"status": "passed"`, "valid fixture accepted", "invalid fixture rejected"} {
		if !strings.Contains(stdout.String(), needle) {
			t.Fatalf("stdout missing %q in:\n%s", needle, stdout.String())
		}
	}
	data, err := os.ReadFile(report)
	if err != nil {
		t.Fatalf("read schema-check report: %v", err)
	}
	if !strings.Contains(string(data), `"status": "passed"`) {
		t.Fatalf("report = %s; want passed status", string(data))
	}
}

func TestSchemaValidateFixtureFailsWhenInvalidFixturePasses(t *testing.T) {
	fixture := writeSchemaCheckFixture(t, true)
	report := filepath.Join(t.TempDir(), "schema-check.json")
	var stdout, stderr bytes.Buffer

	got := run([]string{"schema-check", "--fixture", fixture, "--report", report}, strings.NewReader(""), &stdout, &stderr)

	if got != 1 {
		t.Fatalf("schema-check invalid fixture exit = %d; want 1", got)
	}
	for _, needle := range []string{"schema-check found", "invalid fixture unexpectedly passed"} {
		combined := stdout.String() + stderr.String()
		if !strings.Contains(combined, needle) {
			t.Fatalf("output missing %q; stdout=%q stderr=%q", needle, stdout.String(), stderr.String())
		}
	}
	data, err := os.ReadFile(report)
	if err != nil {
		t.Fatalf("read schema-check report: %v", err)
	}
	if !strings.Contains(string(data), `"status": "failed"`) {
		t.Fatalf("report = %s; want failed status", string(data))
	}
}

func TestSchemaCheckAllWritesRepoReport(t *testing.T) {
	chdir(t, repoRoot(t))
	report := filepath.Join(t.TempDir(), "schema-check.json")
	var stdout, stderr bytes.Buffer

	got := run([]string{"schema-check", "--all", "--report", report}, strings.NewReader(""), &stdout, &stderr)

	if got != 0 {
		t.Fatalf("schema-check --all exit = %d, stderr %q, stdout:\n%s", got, stderr.String(), stdout.String())
	}
	for _, needle := range []string{`"status": "passed"`, ".agent/registries/command-registry.yaml", ".agent/registries/issue-registry.yaml", ".agent/policies/layer-governance.yaml", "layer governance semantics"} {
		if !strings.Contains(stdout.String(), needle) {
			t.Fatalf("stdout missing %q in:\n%s", needle, stdout.String())
		}
	}
	if _, err := os.Stat(report); err != nil {
		t.Fatalf("schema-check report missing: %v", err)
	}
}

func TestSchemaValidateRequiresExactlyOneMode(t *testing.T) {
	var stdout, stderr bytes.Buffer

	got := run([]string{"schema", "validate"}, strings.NewReader(""), &stdout, &stderr)

	if got != 2 {
		t.Fatalf("schema validate without mode exit = %d; want 2", got)
	}
	if !strings.Contains(stderr.String(), "requires exactly one of --all or --fixture") {
		t.Fatalf("stderr = %q; want mode error", stderr.String())
	}
}

func TestSchemaCommandRejectsUnknownSubcommand(t *testing.T) {
	for _, args := range [][]string{nil, {"check"}} {
		var stdout, stderr bytes.Buffer

		got := runSchemaCommand(args, &stdout, &stderr)

		if got != 2 {
			t.Fatalf("runSchemaCommand(%v) exit = %d; want 2", args, got)
		}
		if !strings.Contains(stderr.String(), "schema usage") {
			t.Fatalf("stderr = %q; want schema usage", stderr.String())
		}
	}
}

func TestSchemaTypeUnmarshalJSONVariants(t *testing.T) {
	var one schemaType
	if err := json.Unmarshal([]byte(`"object"`), &one); err != nil {
		t.Fatalf("unmarshal string schema type: %v", err)
	}
	if got := strings.Join(one, ","); got != "object" {
		t.Fatalf("single schema type = %q; want object", got)
	}

	var many schemaType
	if err := json.Unmarshal([]byte(`["object","null"]`), &many); err != nil {
		t.Fatalf("unmarshal array schema type: %v", err)
	}
	if got := strings.Join(many, ","); got != "object,null" {
		t.Fatalf("multi schema type = %q; want object,null", got)
	}

	var invalid schemaType
	if err := json.Unmarshal([]byte(`12`), &invalid); err == nil {
		t.Fatal("unmarshal numeric schema type succeeded; want error")
	}
}

func TestWriteSchemaCheckReportHandlesEmptyPathAndFilesystemErrors(t *testing.T) {
	if err := writeSchemaCheckReport("", []byte("ignored")); err != nil {
		t.Fatalf("empty report path error = %v; want nil", err)
	}

	dir := t.TempDir()
	fileParent := filepath.Join(dir, "file-parent")
	if err := os.WriteFile(fileParent, []byte("not a directory"), 0o644); err != nil {
		t.Fatalf("write file parent: %v", err)
	}
	if err := writeSchemaCheckReport(filepath.Join(fileParent, "report.json"), []byte("{}")); err == nil {
		t.Fatal("report under file parent succeeded; want mkdir error")
	}
	if err := writeSchemaCheckReport(dir, []byte("{}")); err == nil {
		t.Fatal("report path pointing at directory succeeded; want write error")
	}
}

func TestValidateYAMLArtifactAgainstSchemaBranches(t *testing.T) {
	dir := t.TempDir()
	schemaPath := filepath.Join(dir, "contract.schema.json")
	artifactPath := filepath.Join(dir, "artifact.yaml")
	writeSchemaCheckText(t, schemaPath, `{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "required": ["name"],
  "properties": {
    "name": {"type": "string", "minLength": 1},
    "enabled": {"type": "boolean"}
  }
}
`)
	writeSchemaCheckText(t, artifactPath, "name: ok\nenabled: true\n")

	check := validateYAMLArtifactAgainstSchema(artifactPath, schemaPath)
	if check.Status != "passed" {
		t.Fatalf("valid YAML check status = %s, gaps = %v", check.Status, check.Gaps)
	}
	if joined := strings.Join(check.Details, " "); !strings.Contains(joined, "artifact parsed") {
		t.Fatalf("valid YAML check details = %v; want parsed detail", check.Details)
	}

	missingSchema := validateYAMLArtifactAgainstSchema(artifactPath, filepath.Join(dir, "missing.schema.json"))
	if missingSchema.Status != "failed" || !strings.Contains(strings.Join(missingSchema.Gaps, "\n"), "missing") {
		t.Fatalf("missing schema check = %+v; want missing schema gap", missingSchema)
	}

	missingArtifact := validateYAMLArtifactAgainstSchema(filepath.Join(dir, "missing.yaml"), schemaPath)
	if missingArtifact.Status != "failed" || !strings.Contains(strings.Join(missingArtifact.Gaps, "\n"), "missing") {
		t.Fatalf("missing artifact check = %+v; want missing artifact gap", missingArtifact)
	}

	writeSchemaCheckText(t, artifactPath, "enabled: true\n")
	invalid := validateYAMLArtifactAgainstSchema(artifactPath, schemaPath)
	if invalid.Status != "failed" || !strings.Contains(strings.Join(invalid.Gaps, "\n"), `missing required field "name"`) {
		t.Fatalf("invalid YAML check = %+v; want required-field gap", invalid)
	}
}

func TestSelectFixtureSchemaUsesExactStemThenSortedFallback(t *testing.T) {
	schemas := map[string]jsonSchema{
		"beta":  {Title: "beta schema"},
		"alpha": {Title: "alpha schema"},
	}

	path, schema := selectFixtureSchema(filepath.Join("fixtures", "valid", "beta.json"), schemas)
	if path != filepath.Join("schemas", "beta.schema.json") || schema.Title != "beta schema" {
		t.Fatalf("exact fixture schema = (%q, %q); want beta schema", path, schema.Title)
	}

	path, schema = selectFixtureSchema(filepath.Join("fixtures", "valid", "other.json"), schemas)
	if path != filepath.Join("schemas", "alpha.schema.json") || schema.Title != "alpha schema" {
		t.Fatalf("fallback fixture schema = (%q, %q); want sorted alpha schema", path, schema.Title)
	}
}

func TestValueMatchesTypeBranches(t *testing.T) {
	cases := []struct {
		name  string
		value any
		typ   string
		want  bool
	}{
		{name: "object", value: map[string]any{"key": "value"}, typ: "object", want: true},
		{name: "array", value: []any{"value"}, typ: "array", want: true},
		{name: "string", value: "value", typ: "string", want: true},
		{name: "integer", value: float64(2), typ: "integer", want: true},
		{name: "integer rejects fraction", value: float64(2.5), typ: "integer"},
		{name: "number", value: float64(2.5), typ: "number", want: true},
		{name: "boolean", value: true, typ: "boolean", want: true},
		{name: "null", value: nil, typ: "null", want: true},
		{name: "unknown type is permissive", value: struct{}{}, typ: "custom", want: true},
		{name: "wrong type", value: "value", typ: "object"},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if got := valueMatchesType(tt.value, tt.typ); got != tt.want {
				t.Fatalf("valueMatchesType(%#v, %q) = %v; want %v", tt.value, tt.typ, got, tt.want)
			}
		})
	}
}

func TestLayerGovernanceSemanticCheckAcceptsRepoFixture(t *testing.T) {
	chdir(t, repoRoot(t))

	check := validateLayerGovernanceSemantics(".agent/policies/layer-governance.yaml")

	if check.Status != "passed" {
		t.Fatalf("layer governance semantics status = %s, gaps = %v", check.Status, check.Gaps)
	}
}

func TestLayerGovernanceSemanticCheckRejectsBoundaryDrift(t *testing.T) {
	chdir(t, repoRoot(t))
	data, err := os.ReadFile(".agent/policies/layer-governance.yaml")
	if err != nil {
		t.Fatalf("read layer governance fixture: %v", err)
	}
	text := string(data)
	text = strings.Replace(text,
		"repos: [redisx, kafkax, postgresx, taosx, ossx, clickhousex, natsx]",
		"repos: [redisx, kafkax, postgresx, taosx, ossx, clickhousex, market-data]",
		1,
	)
	text = strings.Replace(text,
		"  - id: L3\n    visibility: private",
		"  - id: L3\n    visibility: public",
		1,
	)
	path := filepath.Join(t.TempDir(), "layer-governance.yaml")
	writeSchemaCheckText(t, path, text)

	check := validateLayerGovernanceSemantics(path)

	if check.Status != "failed" {
		t.Fatalf("layer governance semantics status = %s; want failed", check.Status)
	}
	combined := strings.Join(check.Gaps, "\n")
	for _, needle := range []string{
		"L2 repos missing natsx",
		"public layer L2 must not include private repo market-data",
		"L3 visibility must be private",
	} {
		if !strings.Contains(combined, needle) {
			t.Fatalf("gaps missing %q in:\n%s", needle, combined)
		}
	}
}

func TestLayerGovernanceSemanticGapsReportsStructuralDrift(t *testing.T) {
	value := map[string]any{
		"dependency_direction": "wrong",
		"layers": []any{
			map[string]any{},
			map[string]any{
				"id":            "Standard",
				"visibility":    "private",
				"repos":         []any{"x.go", "extra"},
				"may_depend_on": []any{"L0"},
				"forbids":       []any{"business_semantics"},
			},
			map[string]any{"id": "Standard"},
			map[string]any{"id": "Unknown", "visibility": "public"},
		},
		"rules": []any{
			map[string]any{},
			map[string]any{
				"id":       "LAYER-P0-PRIVATE-BOUNDARY",
				"level":    "P2",
				"evidence": []any{"docs-check"},
			},
			map[string]any{"id": "LAYER-P0-PRIVATE-BOUNDARY"},
			map[string]any{"id": "UNKNOWN-RULE"},
		},
	}

	gaps := layerGovernanceSemanticGaps("layers.yaml", value)
	combined := strings.Join(gaps, "\n")

	for _, want := range []string{
		`layers.yaml dependency_direction = "wrong"; want L3>L2>L1>L0>Standard`,
		"layers.yaml layers contains item missing id",
		"layers.yaml duplicate layer Standard",
		"layers.yaml unexpected layer Unknown",
		"layers.yaml missing required layer L0",
		"layers.yaml Standard visibility must be public",
		"layers.yaml Standard repos missing xlib-standard",
		"layers.yaml Standard repos unexpected extra",
		"layers.yaml public layer Standard must not include private repo x.go",
		"layers.yaml rules contains item missing id",
		"layers.yaml duplicate rule LAYER-P0-PRIVATE-BOUNDARY",
		"layers.yaml unexpected rule UNKNOWN-RULE",
		"layers.yaml LAYER-P0-PRIVATE-BOUNDARY level must be P0",
		"layers.yaml LAYER-P0-PRIVATE-BOUNDARY evidence missing schema-check",
	} {
		if !strings.Contains(combined, want) {
			t.Fatalf("gaps missing %q in:\n%s", want, combined)
		}
	}
}

func TestSchemaCheckListHelpersIgnoreNonMapAndNonStringValues(t *testing.T) {
	maps := schemaCheckMapList([]any{map[string]any{"id": "ok"}, "skip", nil})
	if len(maps) != 1 || schemaCheckStringField(maps[0], "id") != "ok" {
		t.Fatalf("schemaCheckMapList() = %#v; want one map item", maps)
	}
	if maps := schemaCheckMapList("not-a-list"); maps != nil {
		t.Fatalf("schemaCheckMapList(non-list) = %#v; want nil", maps)
	}

	values := schemaCheckStringList([]any{"one", 2, "two", nil})
	if len(values) != 2 || values[0] != "one" || values[1] != "two" {
		t.Fatalf("schemaCheckStringList() = %#v; want string-only values", values)
	}
	if values := schemaCheckStringList("not-a-list"); values != nil {
		t.Fatalf("schemaCheckStringList(non-list) = %#v; want nil", values)
	}
}

func TestSchemaValidateArgumentAndDefaultBranches(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := runSchemaValidate([]string{"-h"}, &stdout, &stderr); code != 0 {
		t.Fatalf("runSchemaValidate(-h) code = %d stderr = %q; want 0", code, stderr.String())
	}

	stdout.Reset()
	stderr.Reset()
	if code := runSchemaValidate([]string{"--fixture"}, &stdout, &stderr); code != 2 {
		t.Fatalf("runSchemaValidate(parse error) code = %d; want 2", code)
	}

	stdout.Reset()
	stderr.Reset()
	if code := runSchemaValidate([]string{"unexpected"}, &stdout, &stderr); code != 2 {
		t.Fatalf("runSchemaValidate(unexpected) code = %d; want 2", code)
	}
	if !strings.Contains(stderr.String(), "unexpected positional argument") {
		t.Fatalf("stderr = %q; want positional argument error", stderr.String())
	}

	stdout.Reset()
	stderr.Reset()
	if code := runSchemaValidate([]string{"--all", "--fixture", "fixtures"}, &stdout, &stderr); code != 2 {
		t.Fatalf("runSchemaValidate(mode mismatch) code = %d; want 2", code)
	}
	if !strings.Contains(stderr.String(), "requires exactly one of --all or --fixture") {
		t.Fatalf("stderr = %q; want mode mismatch error", stderr.String())
	}

	if got := withSchemaCheckDefaultAll([]string{"--report", "out.json"}); len(got) == 0 || got[0] != "--all" {
		t.Fatalf("withSchemaCheckDefaultAll(report) = %#v; want --all prefix", got)
	}
	if got := withSchemaCheckDefaultAll([]string{"--fixture=fixtures"}); strings.Join(got, " ") != "--fixture=fixtures" {
		t.Fatalf("withSchemaCheckDefaultAll(fixture=) = %#v; want unchanged fixture args", got)
	}
}

func TestRunSchemaCheckDefaultsAllAndReportsFailure(t *testing.T) {
	chdir(t, t.TempDir())
	var stdout, stderr bytes.Buffer

	code := runSchemaCheck([]string{"--report", ""}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("runSchemaCheck(default all failure) code = %d stdout = %q stderr = %q; want 1", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), `"status": "failed"`) {
		t.Fatalf("stdout = %q; want failed report", stdout.String())
	}
	if !strings.Contains(stderr.String(), "schema-check found") {
		t.Fatalf("stderr = %q; want failure summary", stderr.String())
	}
}

func TestRunSchemaCheckReportsMarshalError(t *testing.T) {
	chdir(t, repoRoot(t))
	old := schemaCheckMarshalIndent
	schemaCheckMarshalIndent = func(any, string, string) ([]byte, error) {
		return nil, errors.New("marshal failed")
	}
	t.Cleanup(func() { schemaCheckMarshalIndent = old })

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	got := runSchemaCheck([]string{"--all", "--report", ""}, &stdout, &stderr)
	if got != 1 {
		t.Fatalf("runSchemaCheck(marshal error) = %d stdout = %q stderr = %q; want 1", got, stdout.String(), stderr.String())
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}
	if !strings.Contains(stderr.String(), "marshal failed") {
		t.Fatalf("stderr = %q; want marshal failure", stderr.String())
	}
}

func TestSchemaValidateReportWriteError(t *testing.T) {
	root := writeSchemaCheckFixture(t, false)
	reportPath := filepath.Join(t.TempDir(), "report-dir")
	if err := os.MkdirAll(reportPath, 0o755); err != nil {
		t.Fatalf("mkdir report path: %v", err)
	}

	var stdout, stderr bytes.Buffer
	code := runSchemaValidate([]string{"--fixture", root, "--report", reportPath}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("runSchemaValidate(report dir) code = %d stdout = %q stderr = %q; want 1", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stderr.String(), "schema-check write report") {
		t.Fatalf("stderr = %q; want report write error", stderr.String())
	}
}

func TestValidateRepoSchemaArtifactsNoSchemas(t *testing.T) {
	chdir(t, t.TempDir())

	checks := validateRepoSchemaArtifacts()

	assertSchemaCheckGap(t, checks, "no contract schemas discovered")
}

func TestValidateRepoSchemaArtifactsReportsSchemaDrift(t *testing.T) {
	root := t.TempDir()
	chdir(t, root)
	if err := os.MkdirAll(filepath.Join(root, "contracts"), 0o755); err != nil {
		t.Fatalf("mkdir contracts: %v", err)
	}
	writeSchemaCheckText(t, filepath.Join(root, "contracts", "broken.schema.json"), `{`)
	writeSchemaCheckText(t, filepath.Join(root, "contracts", "missing-meta.schema.json"), `{"type":"object"}`)
	writeSchemaCheckText(t, filepath.Join(root, "contracts", "wrong-root.schema.json"), `{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "string"
}`)

	checks := validateRepoSchemaArtifacts()

	for _, needle := range []string{
		"invalid schema JSON",
		"missing-meta.schema.json missing $schema",
		"wrong-root.schema.json root type is not object",
		"missing contracts/command-registry.schema.json",
		"missing .agent/policies/layer-governance.yaml",
	} {
		assertSchemaCheckGap(t, checks, needle)
	}
}

func TestValidateFixtureSchemasDiscoveryAndJSONGroupBranches(t *testing.T) {
	checks := validateFixtureSchemas(t.TempDir())
	assertSchemaCheckGap(t, checks, "no fixture schemas discovered")

	badSchemaRoot := t.TempDir()
	if err := os.MkdirAll(filepath.Join(badSchemaRoot, "schemas"), 0o755); err != nil {
		t.Fatalf("mkdir schemas: %v", err)
	}
	writeSchemaCheckText(t, filepath.Join(badSchemaRoot, "schemas", "broken.schema.json"), `{`)
	checks = validateFixtureSchemas(badSchemaRoot)
	assertSchemaCheckGap(t, checks, "invalid schema JSON")

	noCaseRoot := t.TempDir()
	if err := os.MkdirAll(filepath.Join(noCaseRoot, "schemas"), 0o755); err != nil {
		t.Fatalf("mkdir schemas: %v", err)
	}
	writeSchemaCheckText(t, filepath.Join(noCaseRoot, "schemas", "sample.schema.json"), `{"type":"object"}`)
	checks = validateFixtureSchemas(noCaseRoot)
	assertSchemaCheckGap(t, checks, "no valid fixtures discovered")
	assertSchemaCheckGap(t, checks, "no invalid fixtures discovered")

	schema := jsonSchema{
		Type:     schemaType{"object"},
		Required: []string{"name"},
		Properties: map[string]jsonSchema{
			"name": {Type: schemaType{"string"}},
		},
	}
	schemas := map[string]jsonSchema{"sample": schema}
	fixtureRoot := t.TempDir()
	validDir := filepath.Join(fixtureRoot, "valid")
	invalidDir := filepath.Join(fixtureRoot, "invalid")
	for _, dir := range []string{validDir, invalidDir} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", dir, err)
		}
	}

	validChecks := validateFixtureJSONGroup(filepath.Join(validDir, "*.json"), schemas, true)
	assertSchemaCheckGap(t, validChecks, "no valid fixtures discovered")

	writeSchemaCheckText(t, filepath.Join(validDir, "sample.json"), `{"name":"ok"}`)
	writeSchemaCheckText(t, filepath.Join(validDir, "broken.json"), `{`)
	validChecks = validateFixtureJSONGroup(filepath.Join(validDir, "*.json"), schemas, true)
	assertSchemaCheckGap(t, validChecks, "invalid JSON")

	if err := os.Remove(filepath.Join(validDir, "broken.json")); err != nil {
		t.Fatalf("remove broken valid fixture: %v", err)
	}
	writeSchemaCheckText(t, filepath.Join(validDir, "missing.json"), `{}`)
	validChecks = validateFixtureJSONGroup(filepath.Join(validDir, "*.json"), schemas, true)
	assertSchemaCheckGap(t, validChecks, `missing required field "name"`)

	writeSchemaCheckText(t, filepath.Join(invalidDir, "sample.json"), `{"name":"ok"}`)
	invalidChecks := validateFixtureJSONGroup(filepath.Join(invalidDir, "*.json"), schemas, false)
	assertSchemaCheckGap(t, invalidChecks, "invalid fixture unexpectedly passed")

	writeSchemaCheckText(t, filepath.Join(invalidDir, "sample.json"), `{}`)
	invalidChecks = validateFixtureJSONGroup(filepath.Join(invalidDir, "*.json"), schemas, false)
	assertSchemaCheckDetail(t, invalidChecks, "invalid fixture rejected")
}

func TestSchemaValidationPrimitiveBranches(t *testing.T) {
	if !schemaAllowsType(jsonSchema{}, "object") {
		t.Fatalf("schemaAllowsType(empty, object) = false; want true")
	}
	if schemaAllowsType(jsonSchema{Type: schemaType{"string"}}, "object") {
		t.Fatalf("schemaAllowsType(string, object) = true; want false")
	}
	if !valueMatchesAnyType("value", []string{"number", "string"}) {
		t.Fatalf("valueMatchesAnyType(string, number|string) = false; want true")
	}
	if valueMatchesAnyType(true, []string{"string"}) {
		t.Fatalf("valueMatchesAnyType(bool, string) = true; want false")
	}

	minLength := 4
	minItems := 2
	minimum := 10.0
	schema := jsonSchema{
		Type:     schemaType{"object"},
		Required: []string{"missing"},
		Properties: map[string]jsonSchema{
			"const":       {Const: "expected"},
			"enum":        {Enum: []any{"allowed"}},
			"items":       {Type: schemaType{"array"}, MinItems: &minItems, Items: &jsonSchema{Type: schemaType{"string"}}},
			"minLength":   {Type: schemaType{"string"}, MinLength: &minLength},
			"badPattern":  {Type: schemaType{"string"}, Pattern: "["},
			"pattern":     {Type: schemaType{"string"}, Pattern: "^ok$"},
			"minimum":     {Type: schemaType{"number"}, Minimum: &minimum},
			"nestedArray": {Type: schemaType{"array"}, Items: &jsonSchema{Type: schemaType{"object"}}},
		},
	}
	value := map[string]any{
		"const":       "actual",
		"enum":        "other",
		"items":       []any{"ok"},
		"minLength":   "no",
		"badPattern":  "value",
		"pattern":     "bad",
		"minimum":     1.0,
		"nestedArray": []any{"not-object"},
	}

	var gaps []string
	appendSchemaValidationGaps(value, schema, "root", &gaps)
	combined := strings.Join(gaps, "\n")
	for _, needle := range []string{
		`root missing required field "missing"`,
		"root.const expected const expected",
		"root.enum expected enum [allowed]",
		"root.items expected at least 2 item(s)",
		"root.minLength expected minLength 4",
		`root.badPattern invalid schema pattern "["`,
		`root.pattern value does not match pattern "^ok$"`,
		"root.minimum expected minimum 10",
		"root.nestedArray[0] expected type object",
	} {
		if !strings.Contains(combined, needle) {
			t.Fatalf("schema gaps missing %q in:\n%s", needle, combined)
		}
	}

	gaps = nil
	appendSchemaValidationGaps("wrong", jsonSchema{Type: schemaType{"object"}}, "root", &gaps)
	if len(gaps) != 1 || !strings.Contains(gaps[0], "root expected type object") {
		t.Fatalf("type mismatch gaps = %#v; want root object mismatch", gaps)
	}
}

func TestBaselineYAMLParserBranches(t *testing.T) {
	data, err := parseBaselineYAML(`
name: sample # stripped comment
enabled: true
threshold: 12
items:
  - path: docs/a.md
    tags: [alpha, beta]
    quoted: "hash # kept"
`, "fixture.yaml")
	if err != nil {
		t.Fatalf("parseBaselineYAML(valid) error = %v", err)
	}
	if data["name"] != "sample" || data["enabled"] != true || data["threshold"] != 12.0 {
		t.Fatalf("parseBaselineYAML(valid) = %#v; want parsed scalar values", data)
	}
	items, ok := data["items"].([]any)
	if !ok || len(items) != 1 {
		t.Fatalf("items = %#v; want one parsed list item", data["items"])
	}

	if _, err := parseBaselineYAML("invalid", "fixture.yaml"); err == nil || !strings.Contains(err.Error(), "invalid YAML line") {
		t.Fatalf("parseBaselineYAML(invalid top-level) error = %v; want invalid line", err)
	}
	if _, err := parseBaselineYAML("items:\n  - invalid-list-item", "fixture.yaml"); err == nil || !strings.Contains(err.Error(), "invalid YAML list item") {
		t.Fatalf("parseBaselineYAML(invalid list item) error = %v; want invalid list item", err)
	}
	if _, err := parseBaselineYAML("items:\n  - path: docs/a.md\n    invalid-mapping", "fixture.yaml"); err == nil || !strings.Contains(err.Error(), "invalid YAML mapping") {
		t.Fatalf("parseBaselineYAML(invalid mapping) error = %v; want invalid mapping", err)
	}
	if _, err := parseBaselineYAMLFile(filepath.Join(t.TempDir(), "missing.yaml")); err == nil || !strings.Contains(err.Error(), "missing") {
		t.Fatalf("parseBaselineYAMLFile(missing) error = %v; want missing", err)
	}
}

func TestReadJSONValueAndParseYAMLScalarBranches(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "value.json")
	writeSchemaCheckText(t, path, `{"name":"ok"}`)

	value, err := readJSONValue(path)
	if err != nil {
		t.Fatalf("readJSONValue(valid) error = %v", err)
	}
	object, ok := value.(map[string]any)
	if !ok || object["name"] != "ok" {
		t.Fatalf("readJSONValue(valid) = %#v; want object name", value)
	}
	if _, err := readJSONValue(filepath.Join(dir, "missing.json")); err == nil || !strings.Contains(err.Error(), "missing") {
		t.Fatalf("readJSONValue(missing) error = %v; want missing", err)
	}

	cases := []struct {
		input string
		want  any
	}{
		{input: "[]", want: []any{}},
		{input: "[one, 2, false]", want: []any{"one", 2.0, false}},
		{input: "'quoted'", want: "quoted"},
		{input: "-1", want: "-1"},
		{input: "1_000", want: "1_000"},
		{input: "", want: ""},
	}
	for _, tt := range cases {
		t.Run(tt.input, func(t *testing.T) {
			if got := parseYAMLScalar(tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("parseYAMLScalar(%q) = %#v; want %#v", tt.input, got, tt.want)
			}
		})
	}
}

func assertSchemaCheckGap(t *testing.T, checks []schemaCheckItem, needle string) {
	t.Helper()
	if strings.Contains(joinSchemaCheckItems(checks), needle) {
		return
	}
	t.Fatalf("schema checks missing gap/detail %q in:\n%s", needle, joinSchemaCheckItems(checks))
}

func assertSchemaCheckDetail(t *testing.T, checks []schemaCheckItem, needle string) {
	t.Helper()
	if strings.Contains(joinSchemaCheckItems(checks), needle) {
		return
	}
	t.Fatalf("schema checks missing detail %q in:\n%s", needle, joinSchemaCheckItems(checks))
}

func joinSchemaCheckItems(checks []schemaCheckItem) string {
	var lines []string
	for _, check := range checks {
		lines = append(lines, check.Name, check.Status)
		lines = append(lines, check.Details...)
		lines = append(lines, check.Gaps...)
	}
	return strings.Join(lines, "\n")
}

func writeSchemaCheckFixture(t *testing.T, invalidFixturePasses bool) string {
	t.Helper()
	root := t.TempDir()
	for _, dir := range []string{"schemas", "valid", "invalid"} {
		if err := os.MkdirAll(filepath.Join(root, dir), 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", dir, err)
		}
	}
	writeSchemaCheckText(t, filepath.Join(root, "schemas", "sample.schema.json"), `{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "required": ["name", "count", "tags"],
  "properties": {
    "name": {"type": "string", "minLength": 1},
    "count": {"type": "integer", "minimum": 0},
    "tags": {"type": "array", "minItems": 1, "items": {"type": "string"}}
  }
}
`)
	writeSchemaCheckText(t, filepath.Join(root, "valid", "sample.json"), `{"name":"ok","count":1,"tags":["fixture"]}`)
	invalid := `{"count":1,"tags":["fixture"]}`
	if invalidFixturePasses {
		invalid = `{"name":"ok","count":1,"tags":["fixture"]}`
	}
	writeSchemaCheckText(t, filepath.Join(root, "invalid", "sample.json"), invalid)
	return root
}

func writeSchemaCheckText(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
