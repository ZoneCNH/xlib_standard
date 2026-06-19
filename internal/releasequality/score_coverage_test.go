// SPDX-License-Identifier: Apache-2.0
package releasequality

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestScoreGateDimensionReportsMissingReleaseFinalCheckNeedle exercises the
// per-needle missing-detection branch (lines covered only when one needle is
// present and the other is absent).
func TestScoreGateDimensionReportsMissingReleaseFinalCheckNeedle(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "Makefile")
	// score-check present, release-final-check absent, score minimum satisfied.
	content := `score-check:
		score --min 9.8
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write Makefile fixture: %v", err)
	}

	dimension := scoreGateDimension(path)
	if dimension.Passed {
		t.Fatal("expected dimension to fail when release-final-check is missing")
	}
	if !strings.Contains(dimension.Detail, "release-final-check") {
		t.Fatalf("dimension detail = %q; want release-final-check missing", dimension.Detail)
	}
}

// TestScoreGateDimensionReportsMissingScoreCheckNeedle covers the second needle
// miss branch to fully exercise the loop body.
func TestScoreGateDimensionReportsMissingScoreCheckNeedle(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "Makefile")
	content := `release-final-check:
		score --min 9.8
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write Makefile fixture: %v", err)
	}

	dimension := scoreGateDimension(path)
	if dimension.Passed {
		t.Fatal("expected dimension to fail when score-check is missing")
	}
	if !strings.Contains(dimension.Detail, "score-check") {
		t.Fatalf("dimension detail = %q; want score-check missing", dimension.Detail)
	}
}

// TestScoreGateDimensionReportsMissingPath covers the os.ReadFile error branch.
func TestScoreGateDimensionReportsMissingPath(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "does-not-exist", "Makefile")

	dimension := scoreGateDimension(path)
	if dimension.Passed {
		t.Fatal("expected dimension to fail when Makefile is missing")
	}
	if !strings.Contains(dimension.Detail, "missing "+path) {
		t.Fatalf("dimension detail = %q; want missing path", dimension.Detail)
	}
}
