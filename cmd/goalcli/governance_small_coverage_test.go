package main

import (
	"bytes"
	"strings"
	"testing"
)

// TestContextGateProfile covers all branches.
func TestContextGateProfile(t *testing.T) {
	cases := map[string]struct {
		profile string
		ok      bool
	}{
		"context-lite":         {"lite", true},
		"context-fast-check":   {"lite", true},
		"context-standard":     {"standard", true},
		"context-standard-check": {"standard", true},
		"context-full":         {"full", true},
		"context-full-check":   {"full", true},
		"context-release":      {"release", true},
		"bogus":                {"", false},
	}
	for gate, want := range cases {
		got, ok := contextGateProfile(gate)
		if got != want.profile || ok != want.ok {
			t.Errorf("contextGateProfile(%q) = (%q,%v); want (%q,%v)", gate, got, ok, want.profile, want.ok)
		}
	}
}

// TestMakefileTargetDependencies covers no-block, continuation lines, and tab-recipe break.
func TestMakefileTargetDependencies(t *testing.T) {
	t.Run("missing block returns nil", func(t *testing.T) {
		if got := makefileTargetDependencies("nothing here", "missing"); got != nil {
			t.Fatalf("got = %v; want nil", got)
		}
	})
	t.Run("continuation lines", func(t *testing.T) {
		content := "target: a b \\\n c d\n\tnot-a-dep\n"
		got := makefileTargetDependencies(content, "target")
		if !slicesContain(got, "a") || !slicesContain(got, "d") {
			t.Fatalf("got = %v; want a..d", got)
		}
	})
	t.Run("no continuation stops at first line", func(t *testing.T) {
		content := "target: a b\n\trecipe\n"
		got := makefileTargetDependencies(content, "target")
		if !slicesContain(got, "a") || !slicesContain(got, "b") {
			t.Fatalf("got = %v", got)
		}
	})
}

// TestAppendMakefileDuplicateGaps covers count==1 (no gap) and count!=1 (gap).
func TestAppendMakefileDuplicateGaps(t *testing.T) {
	var gaps []string
	// Single definition: no gap.
	appendMakefileDuplicateGaps("t: a\n", []string{"t"}, &gaps)
	if len(gaps) != 0 {
		t.Fatalf("single def gaps = %v; want none", gaps)
	}
	// Duplicate definition: gap.
	appendMakefileDuplicateGaps("t: a\nt: b\n", []string{"t"}, &gaps)
	if len(gaps) == 0 {
		t.Fatalf("duplicate def should add gap")
	}
}

// TestAppendMakefileTargetDependencyGaps covers missing block.
func TestAppendMakefileTargetDependencyGaps(t *testing.T) {
	var gaps []string
	appendMakefileTargetDependencyGaps("nothing", "missing-target", []string{"a"}, nil, &gaps)
	if !gapsContainSubstring(gaps, "missing target block missing-target") {
		t.Fatalf("gaps = %v; want missing block", gaps)
	}
}

// TestAppendMakefileTargetForbiddenReferenceGaps covers missing block + reference found.
func TestAppendMakefileTargetForbiddenReferenceGaps(t *testing.T) {
	t.Run("missing block", func(t *testing.T) {
		var gaps []string
		appendMakefileTargetForbiddenReferenceGaps("nothing", "missing", []string{"x"}, &gaps)
		if !gapsContainSubstring(gaps, "missing target block missing") {
			t.Fatalf("gaps = %v", gaps)
		}
	})
	t.Run("forbidden reference present", func(t *testing.T) {
		var gaps []string
		appendMakefileTargetForbiddenReferenceGaps("t: a\n\trelease-check\n", "t", []string{"release-check"}, &gaps)
		if !gapsContainSubstring(gaps, "must not reference release-check") {
			t.Fatalf("gaps = %v", gaps)
		}
	})
}

// TestAppendReleaseFinalDelegationGaps covers missing block + self-recursion + missing delegate.
func TestAppendReleaseFinalDelegationGaps(t *testing.T) {
	t.Run("missing block", func(t *testing.T) {
		var gaps []string
		appendReleaseFinalDelegationGaps("nothing", &gaps)
		if !gapsContainSubstring(gaps, "missing target block release-final-check") {
			t.Fatalf("gaps = %v", gaps)
		}
	})
	t.Run("self recursion via $(MAKE)", func(t *testing.T) {
		var gaps []string
		content := "release-final-check: a\n\t$(MAKE) release-final-check\n"
		appendReleaseFinalDelegationGaps(content, &gaps)
		if !gapsContainSubstring(gaps, "must not call itself") {
			t.Fatalf("gaps = %v", gaps)
		}
	})
	t.Run("no context-release delegate", func(t *testing.T) {
		var gaps []string
		content := "release-final-check: a\n\techo hi\n"
		appendReleaseFinalDelegationGaps(content, &gaps)
		if !gapsContainSubstring(gaps, "must call context-release") {
			t.Fatalf("gaps = %v", gaps)
		}
	})
}

// TestFlagValue covers all branches.
func TestFlagValue(t *testing.T) {
	args := []string{"--repo", "kernel", "--mode=enforce", "--context", "release_verify"}
	if got := flagValue(args, "repo", "fallback"); got != "kernel" {
		t.Fatalf("repo = %q", got)
	}
	if got := flagValue(args, "mode", "fallback"); got != "enforce" {
		t.Fatalf("mode = %q", got)
	}
	if got := flagValue(args, "missing", "fallback"); got != "fallback" {
		t.Fatalf("missing = %q", got)
	}
	// --name at end with no value -> fallback.
	if got := flagValue([]string{"--repo"}, "repo", "fallback"); got != "fallback" {
		t.Fatalf("repo at end = %q", got)
	}
}

// TestPlannedCommandVerifyRequested covers all branches.
func TestPlannedCommandVerifyRequested(t *testing.T) {
	cases := []struct {
		args []string
		want bool
	}{
		{[]string{"--verify"}, true},
		{[]string{"--strict"}, true},
		{[]string{"--context=release_verify"}, true},
		{[]string{"--context", "release_verify"}, true},
		{[]string{"--context", "local_write"}, false},
		{[]string{"--json"}, false},
		{nil, false},
	}
	for _, c := range cases {
		if got := plannedCommandVerifyRequested(c.args); got != c.want {
			t.Errorf("plannedCommandVerifyRequested(%v) = %v; want %v", c.args, got, c.want)
		}
	}
}

// TestDuplicateValues covers empty-value skipping and dedup of reports.
func TestDuplicateValues(t *testing.T) {
	got := duplicateValues([]string{"a", "", "a", "b", "b", "b", "c"})
	// a and b are duplicates (each reported once); c is not.
	if len(got) != 2 {
		t.Fatalf("got = %v; want 2 duplicates", got)
	}
}

// TestParseMakeEnforcerTarget covers all branches.
func TestParseMakeEnforcerTarget(t *testing.T) {
	cases := []struct {
		fields []string
		target string
		ok     bool
	}{
		{[]string{"make", "governance-check"}, "governance-check", true},
		{[]string{"goalcli", "make", "governance-check"}, "governance-check", true},
		{[]string{"make", "-s", "governance-check"}, "governance-check", true},
		{[]string{"make", "--silent", "governance-check"}, "governance-check", true},
		{[]string{"make", "VAR=val", "governance-check"}, "governance-check", true},
		{[]string{"make"}, "", false},
		{[]string{"goalcli", "evidence-check"}, "", false},
	}
	for _, c := range cases {
		got, ok := parseMakeEnforcerTarget(c.fields)
		if got != c.target || ok != c.ok {
			t.Errorf("parseMakeEnforcerTarget(%v) = (%q,%v); want (%q,%v)", c.fields, got, ok, c.target, c.ok)
		}
	}
}

// TestKnownEnforcementRef covers goalcli/make/hook branches.
func TestKnownEnforcementRef(t *testing.T) {
	commands := map[string]bool{"evidence-check": true}
	makeTargets := map[string]bool{"governance-check": true}
	cases := []struct {
		ref  string
		want bool
	}{
		{"goalcli evidence-check", true},
		{"make governance-check", true},
		{".githooks/pre-commit", false},
		{"make missing", false},
		{"goalcli missing", false},
		{"plain text", false},
	}
	for _, c := range cases {
		if got := knownEnforcementRef(c.ref, commands, makeTargets); got != c.want {
			t.Errorf("knownEnforcementRef(%q) = %v; want %v", c.ref, got, c.want)
		}
	}
}

// TestTrimYAMLScalar covers trimming quotes and whitespace.
func TestTrimYAMLScalar(t *testing.T) {
	cases := map[string]string{
		`"quoted"`: "quoted",
		`'single'`: "single",
		"  plain  ": "plain",
		"plain":    "plain",
	}
	for in, want := range cases {
		if got := trimYAMLScalar(in); got != want {
			t.Errorf("trimYAMLScalar(%q) = %q; want %q", in, got, want)
		}
	}
}

// TestBlockHasYAMLListItem covers present and absent.
func TestBlockHasYAMLListItem(t *testing.T) {
	block := "commands:\n  - name: a\n  - name: b\n"
	if !blockHasYAMLListItem(block, "commands") {
		t.Fatalf("commands list should be present")
	}
	if blockHasYAMLListItem(block, "missing") {
		t.Fatalf("missing should be absent")
	}
}

// TestAppendRequiredAgentIndexField covers missing-field (gap added) and present-field (no gap).
func TestAppendRequiredAgentIndexField(t *testing.T) {
	t.Run("missing field", func(t *testing.T) {
		var gaps []string
		entry := agentIndexEntry{path: ".agent/x", block: "  path: .agent/x\n"}
		appendRequiredAgentIndexField(".agent/index.yaml", entry, "owner", &gaps)
		if len(gaps) == 0 {
			t.Fatalf("want gap for missing owner")
		}
	})
	t.Run("present field", func(t *testing.T) {
		var gaps []string
		entry := agentIndexEntry{path: ".agent/x", block: "  path: .agent/x\n  owner: governance\n"}
		appendRequiredAgentIndexField(".agent/index.yaml", entry, "owner", &gaps)
		if len(gaps) != 0 {
			t.Fatalf("want no gap for present owner")
		}
	})
}

// TestAppendAgentIndexEnumGap covers missing, invalid, and valid values.
func TestAppendAgentIndexEnumGap(t *testing.T) {
	allowed := map[string]bool{"registry": true, "policy": true}
	t.Run("missing value no gap", func(t *testing.T) {
		var gaps []string
		entry := agentIndexEntry{block: "  path: x\n"}
		appendAgentIndexEnumGap("idx", entry, "layer", allowed, &gaps)
		if len(gaps) != 0 {
			t.Fatalf("missing value should not add gap")
		}
	})
	t.Run("invalid value gap", func(t *testing.T) {
		var gaps []string
		entry := agentIndexEntry{block: "  layer: bogus\n"}
		appendAgentIndexEnumGap("idx", entry, "layer", allowed, &gaps)
		if len(gaps) == 0 {
			t.Fatalf("invalid value should add gap")
		}
	})
	t.Run("valid value no gap", func(t *testing.T) {
		var gaps []string
		entry := agentIndexEntry{block: "  layer: registry\n"}
		appendAgentIndexEnumGap("idx", entry, "layer", allowed, &gaps)
		if len(gaps) != 0 {
			t.Fatalf("valid value should not add gap")
		}
	})
}

// TestRunContextProfileBranches covers invalid args, help, positional, invalid profile.
func TestRunContextProfileBranches(t *testing.T) {
	t.Run("flag parse error", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		got := runContextProfile([]string{"--bad"}, &stdout, &stderr)
		if got != 2 {
			t.Fatalf("got = %d; want 2", got)
		}
	})
	t.Run("help", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		got := runContextProfile([]string{"-h"}, &stdout, &stderr)
		if got != 0 {
			t.Fatalf("got = %d; want 0", got)
		}
	})
	t.Run("positional arg", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		got := runContextProfile([]string{"positional"}, &stdout, &stderr)
		if got != 2 {
			t.Fatalf("got = %d; want 2", got)
		}
	})
	t.Run("invalid profile", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		got := runContextProfile([]string{"--profile", "bogus"}, &stdout, &stderr)
		if got != 2 {
			t.Fatalf("got = %d; want 2", got)
		}
	})
	t.Run("fast alias resolves to lite", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		got := runContextProfile([]string{"--profile", "fast"}, &stdout, &stderr)
		if got != 0 {
			t.Fatalf("got = %d; want 0", got)
		}
		if !strings.Contains(stdout.String(), "profile=lite") {
			t.Fatalf("stdout = %q; want profile=lite", stdout.String())
		}
	})
}

// TestRunContextProfileAliasInvalidArgs covers validateInternalCommandArgs error.
func TestRunContextProfileAliasInvalidArgs(t *testing.T) {
	var stdout, stderr bytes.Buffer
	got := runContextProfileAlias("context-lite", []string{"--bad"}, &stdout, &stderr)
	if got != 2 {
		t.Fatalf("got = %d; want 2", got)
	}
}

// TestRunContextProfileCheckBranches covers profile validation error + help.
func TestRunContextProfileCheckBranches(t *testing.T) {
	t.Run("invalid profile", func(t *testing.T) {
		root := t.TempDir()
		chdir(t, root)
		var stdout, stderr bytes.Buffer
		got := runContextProfileCheck("context-profile-check", []string{"--profile", "bogus"}, &stdout, &stderr)
		if got != 2 {
			t.Fatalf("got = %d; want 2", got)
		}
	})
	t.Run("positional arg via flag parse", func(t *testing.T) {
		root := t.TempDir()
		chdir(t, root)
		var stdout, stderr bytes.Buffer
		got := runContextProfileCheck("context-profile-check", []string{"positional"}, &stdout, &stderr)
		if got != 2 {
			t.Fatalf("got = %d; want 2", got)
		}
	})
}

// TestRulesRegistryGoalCLICommands verifies the map is non-empty and contains known commands.
func TestRulesRegistryGoalCLICommands(t *testing.T) {
	m := rulesRegistryGoalCLICommands()
	if len(m) == 0 {
		t.Fatalf("map empty; want entries")
	}
	if !m["evidence-check"] {
		t.Fatalf("evidence-check should be in map")
	}
}
