package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestSchemaTypeUnmarshalJSON covers both single-string and array forms + error.
func TestSchemaTypeUnmarshalJSON(t *testing.T) {
	var s schemaType
	if err := json.Unmarshal([]byte(`"string"`), &s); err != nil {
		t.Fatalf("string unmarshal: %v", err)
	}
	if len(s) != 1 || s[0] != "string" {
		t.Fatalf("s = %v", s)
	}
	var arr schemaType
	if err := json.Unmarshal([]byte(`["a","b"]`), &arr); err != nil {
		t.Fatalf("array unmarshal: %v", err)
	}
	if len(arr) != 2 {
		t.Fatalf("arr = %v", arr)
	}
	var bad schemaType
	if err := json.Unmarshal([]byte(`123`), &bad); err == nil {
		t.Fatalf("number should error")
	}
}

// TestRunSchemaCommand covers no-args and non-validate subcommand.
func TestRunSchemaCommand(t *testing.T) {
	t.Run("no args", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		got := runSchemaCommand(nil, &stdout, &stderr)
		if got != 2 {
			t.Fatalf("got = %d; want 2", got)
		}
	})
	t.Run("non-validate", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		got := runSchemaCommand([]string{"bogus"}, &stdout, &stderr)
		if got != 2 {
			t.Fatalf("got = %d; want 2", got)
		}
	})
}

// TestWithSchemaCheckDefaultAll covers the no-op and prepend branches.
func TestWithSchemaCheckDefaultAll(t *testing.T) {
	if got := withSchemaCheckDefaultAll([]string{"--all"}); len(got) != 1 || got[0] != "--all" {
		t.Fatalf("--all present: %v", got)
	}
	if got := withSchemaCheckDefaultAll([]string{"--fixture", "x"}); len(got) != 2 {
		t.Fatalf("--fixture present: %v", got)
	}
	got := withSchemaCheckDefaultAll([]string{"--json"})
	if len(got) != 2 || got[0] != "--all" {
		t.Fatalf("prepend: %v", got)
	}
}

// TestRunSchemaValidate covers all the flag/branch paths.
func TestRunSchemaValidate(t *testing.T) {
	t.Run("flag parse error", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		got := runSchemaValidate([]string{"--bad"}, &stdout, &stderr)
		if got != 2 {
			t.Fatalf("got = %d; want 2", got)
		}
	})
	t.Run("help", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		got := runSchemaValidate([]string{"-h"}, &stdout, &stderr)
		if got != 0 {
			t.Fatalf("got = %d; want 0", got)
		}
	})
	t.Run("positional arg", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		got := runSchemaValidate([]string{"positional"}, &stdout, &stderr)
		if got != 2 {
			t.Fatalf("got = %d; want 2", got)
		}
	})
	t.Run("neither all nor fixture", func(t *testing.T) {
		root := t.TempDir()
		chdir(t, root)
		var stdout, stderr bytes.Buffer
		got := runSchemaValidate(nil, &stdout, &stderr)
		if got != 2 {
			t.Fatalf("got = %d; want 2", got)
		}
	})
	t.Run("both all and fixture", func(t *testing.T) {
		root := t.TempDir()
		chdir(t, root)
		var stdout, stderr bytes.Buffer
		got := runSchemaValidate([]string{"--all", "--fixture", "x"}, &stdout, &stderr)
		if got != 2 {
			t.Fatalf("got = %d; want 2", got)
		}
	})
	t.Run("--all with no contract schemas fails", func(t *testing.T) {
		root := t.TempDir()
		chdir(t, root)
		var stdout, stderr bytes.Buffer
		got := runSchemaValidate([]string{"--all", "--report", ""}, &stdout, &stderr)
		if got != 1 {
			t.Fatalf("got = %d; want 1", got)
		}
	})
}

// TestWriteSchemaCheckReport covers empty path + mkdir error + success.
func TestWriteSchemaCheckReport(t *testing.T) {
	if err := writeSchemaCheckReport("", []byte("x")); err != nil {
		t.Fatalf("empty path err = %v", err)
	}
	root := t.TempDir()
	blocker := filepath.Join(root, "blocker")
	os.WriteFile(blocker, []byte("x"), 0o644)
	if err := writeSchemaCheckReport(filepath.Join(blocker, "report.json"), []byte("x")); err == nil {
		t.Fatalf("want mkdir error")
	}
	path := filepath.Join(root, "sub", "report.json")
	if err := writeSchemaCheckReport(path, []byte("x")); err != nil {
		t.Fatalf("err = %v", err)
	}
}

// TestLoadJSONSchema covers missing, invalid JSON.
func TestLoadJSONSchema(t *testing.T) {
	_, gaps := loadJSONSchema("nonexistent.json")
	if len(gaps) == 0 {
		t.Fatalf("want gaps for missing")
	}
	root := t.TempDir()
	path := filepath.Join(root, "bad.json")
	os.WriteFile(path, []byte("{bad"), 0o644)
	_, gaps = loadJSONSchema(path)
	if len(gaps) == 0 {
		t.Fatalf("want gaps for invalid JSON")
	}
}

// TestReadJSONValue covers missing, invalid, valid.
func TestReadJSONValue(t *testing.T) {
	if _, err := readJSONValue("nonexistent.json"); err == nil {
		t.Fatalf("want error")
	}
	root := t.TempDir()
	path := filepath.Join(root, "bad.json")
	os.WriteFile(path, []byte("{bad"), 0o644)
	if _, err := readJSONValue(path); err == nil {
		t.Fatalf("want error")
	}
	path2 := filepath.Join(root, "ok.json")
	os.WriteFile(path2, []byte(`{"a":1}`), 0o644)
	if _, err := readJSONValue(path2); err != nil {
		t.Fatalf("err = %v", err)
	}
}

// TestSelectFixtureSchema covers direct stem match and fallback to sorted first key.
func TestSelectFixtureSchema(t *testing.T) {
	schemas := map[string]jsonSchema{
		"thing":   {Schema: "http://x"},
		"apple":   {Schema: "http://y"},
	}
	t.Run("stem match", func(t *testing.T) {
		ref, _ := selectFixtureSchema(filepath.Join("valid", "thing.json"), schemas)
		if !strings.Contains(ref, "thing.schema.json") {
			t.Fatalf("ref = %q; want thing", ref)
		}
	})
	t.Run("fallback sorted first", func(t *testing.T) {
		ref, _ := selectFixtureSchema(filepath.Join("valid", "unknown.json"), schemas)
		if !strings.Contains(ref, "apple.schema.json") {
			t.Fatalf("ref = %q; want apple (sorted first)", ref)
		}
	})
}

// TestValueMatchesType covers all type cases.
func TestValueMatchesType(t *testing.T) {
	cases := []struct {
		val  any
		typ  string
		want bool
	}{
		{"x", "string", true},
		{1, "string", false},
		{float64(1), "integer", true},
		{float64(1.5), "integer", false},
		{float64(1.5), "number", true},
		{true, "boolean", true},
		{[]any{1}, "array", true},
		{map[string]any{}, "object", true},
		{nil, "null", true},
		{"x", "unknowntype", true}, // default returns true
	}
	for _, c := range cases {
		if got := valueMatchesType(c.val, c.typ); got != c.want {
			t.Errorf("valueMatchesType(%v,%q) = %v; want %v", c.val, c.typ, got, c.want)
		}
	}
}

func TestSchemaAllowsType(t *testing.T) {
	s := jsonSchema{Type: schemaType{"object"}}
	if !schemaAllowsType(s, "object") {
		t.Fatalf("object should be allowed")
	}
	if schemaAllowsType(s, "string") {
		t.Fatalf("string should not be allowed")
	}
	empty := jsonSchema{}
	if !schemaAllowsType(empty, "object") {
		t.Fatalf("empty type should allow any")
	}
}

func TestValueMatchesAnyType(t *testing.T) {
	if !valueMatchesAnyType("x", []string{"string", "number"}) {
		t.Fatalf("string should match")
	}
	if valueMatchesAnyType("x", []string{"number"}) {
		t.Fatalf("string should not match number-only")
	}
}

// TestParseBaselineYAMLFile covers missing and valid.
func TestParseBaselineYAMLFile(t *testing.T) {
	if _, err := parseBaselineYAMLFile("nonexistent.yaml"); err == nil {
		t.Fatalf("want error for missing")
	}
	root := t.TempDir()
	path := filepath.Join(root, "ok.yaml")
	os.WriteFile(path, []byte("key: value\n"), 0o644)
	v, err := parseBaselineYAMLFile(path)
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	if v["key"] != "value" {
		t.Fatalf("v = %v", v)
	}
}

// TestStripYAMLComment covers plain, commented, and quoted-hash.
func TestStripYAMLComment(t *testing.T) {
	if got := stripYAMLComment("plain"); got != "plain" {
		t.Fatalf("plain = %q", got)
	}
	if got := stripYAMLComment("key: val # comment"); strings.Contains(got, "comment") {
		t.Fatalf("comment stripped = %q", got)
	}
	// Hash inside quotes should be preserved.
	if got := stripYAMLComment(`key: "a#b"`); !strings.Contains(got, "a#b") {
		t.Fatalf("quoted hash = %q", got)
	}
}

// TestParseYAMLScalar covers quoted, unquoted, bool, array, number.
func TestParseYAMLScalar(t *testing.T) {
	if got := parseYAMLScalar(`"quoted"`); got != "quoted" {
		t.Fatalf("quoted = %v", got)
	}
	if got := parseYAMLScalar("plain"); got != "plain" {
		t.Fatalf("plain = %v", got)
	}
	if got := parseYAMLScalar("true"); got != true {
		t.Fatalf("bool = %v", got)
	}
	if got := parseYAMLScalar("false"); got != false {
		t.Fatalf("bool false = %v", got)
	}
	if got := parseYAMLScalar("123"); got != float64(123) {
		t.Fatalf("number = %v", got)
	}
	arr := parseYAMLScalar("[a, b]")
	if _, ok := arr.([]any); !ok {
		t.Fatalf("array = %v; want []any", arr)
	}
	emptyArr := parseYAMLScalar("[]")
	if _, ok := emptyArr.([]any); !ok {
		t.Fatalf("empty array = %v", emptyArr)
	}
}
