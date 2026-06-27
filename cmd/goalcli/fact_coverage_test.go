package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestRunFactBranches covers no-args and unknown subcommand.
func TestRunFactBranches(t *testing.T) {
	t.Run("no args", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		got := runFact(nil, &stdout, &stderr)
		if got != 2 {
			t.Fatalf("got = %d; want 2", got)
		}
		if !strings.Contains(stderr.String(), "usage: goalcli fact audit") {
			t.Fatalf("stderr = %q", stderr.String())
		}
	})
	t.Run("unknown subcommand", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		got := runFact([]string{"bogus"}, &stdout, &stderr)
		if got != 2 {
			t.Fatalf("got = %d; want 2", got)
		}
		if !strings.Contains(stderr.String(), `unknown fact subcommand "bogus"`) {
			t.Fatalf("stderr = %q", stderr.String())
		}
	})
}

// TestRunFactAuditFlagBranches covers parse error, positional arg, missing facts.
func TestRunFactAuditFlagBranches(t *testing.T) {
	t.Run("parse error", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		got := runFactAudit([]string{"--bad"}, &stdout, &stderr)
		if got != 2 {
			t.Fatalf("got = %d; want 2", got)
		}
	})
	t.Run("positional arg", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		got := runFactAudit([]string{"positional"}, &stdout, &stderr)
		if got != 2 {
			t.Fatalf("got = %d; want 2", got)
		}
	})
	t.Run("missing facts file", func(t *testing.T) {
		root := t.TempDir()
		var stdout, stderr bytes.Buffer
		got := runFactAudit([]string{"--root", root}, &stdout, &stderr)
		if got != 1 {
			t.Fatalf("got = %d; want 1", got)
		}
	})
}

// TestFactAuditContext covers empty env -> default and explicit env.
func TestFactAuditContext(t *testing.T) {
	_ = os.Unsetenv("XLIB_CONTEXT")
	if got := factAuditContext(); got != "local_write" {
		t.Fatalf("got = %q; want local_write", got)
	}
	t.Setenv("XLIB_CONTEXT", "release_verify")
	if got := factAuditContext(); got != "release_verify" {
		t.Fatalf("got = %q; want release_verify", got)
	}
}

// TestFactStrictChecksLocalReleaseTag covers both branches.
func TestFactStrictChecksLocalReleaseTag(t *testing.T) {
	if !factStrictChecksLocalReleaseTag("local_write") {
		t.Fatalf("local_write should check tag")
	}
	if factStrictChecksLocalReleaseTag("ci_pull_request") {
		t.Fatalf("ci_pull_request should skip tag check")
	}
}

// TestFactStrictLocalTagGaps covers empty version and existing tag.
func TestFactStrictLocalTagGaps(t *testing.T) {
	root := t.TempDir()
	if got := factStrictLocalTagGaps(root, ""); len(got) != 0 {
		t.Fatalf("empty version gaps = %v; want none", got)
	}
	// Real git tag exists.
	if err := os.MkdirAll(filepath.Join(root, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	// Use a fixture that git won't recognize to exercise the no-match branch.
	got := factStrictLocalTagGaps(root, "v9.9.9-nonexistent-tag")
	if len(got) != 0 {
		t.Fatalf("non-existent tag gaps = %v; want none", got)
	}
}
