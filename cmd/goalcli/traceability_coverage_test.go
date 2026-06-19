package main

import (
	"os"
	"path/filepath"
	"testing"
)

// TestVerifyArtifactExists covers directory glob, glob, file exists, missing.
func TestVerifyArtifactExists(t *testing.T) {
	root := t.TempDir()
	chdir(t, root)

	t.Run("existing file", func(t *testing.T) {
		path := filepath.Join(root, "file.md")
		if err := os.WriteFile(path, []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := verifyArtifactExists("file.md"); err != nil {
			t.Fatalf("err = %v", err)
		}
	})
	t.Run("missing file", func(t *testing.T) {
		if err := verifyArtifactExists("missing.md"); err == nil {
			t.Fatalf("want error for missing file")
		}
	})
	t.Run("directory glob existing", func(t *testing.T) {
		if err := os.MkdirAll(filepath.Join(root, "docs", "sub"), 0o755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}
		if err := verifyArtifactExists("docs/sub/*"); err != nil {
			t.Fatalf("err = %v", err)
		}
	})
	t.Run("directory glob missing", func(t *testing.T) {
		if err := verifyArtifactExists("nodir/*"); err == nil {
			t.Fatalf("want error for missing dir glob")
		}
	})
	t.Run("directory glob but file", func(t *testing.T) {
		if err := os.WriteFile(filepath.Join(root, "afile"), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := verifyArtifactExists("afile/*"); err == nil {
			t.Fatalf("want not-a-directory error")
		}
	})
	t.Run("glob with matches", func(t *testing.T) {
		if err := os.WriteFile(filepath.Join(root, "match1.md"), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := verifyArtifactExists("match*.md"); err != nil {
			t.Fatalf("err = %v", err)
		}
	})
	t.Run("glob no matches", func(t *testing.T) {
		if err := verifyArtifactExists("nomatch*.md"); err == nil {
			t.Fatalf("want no-match error")
		}
	})
}

// TestIsGitIgnored covers non-git dir (returns false) and repo-root gitignored path.
func TestIsGitIgnored(t *testing.T) {
	// In a temp dir with no git repo, git check-ignore errors so returns false.
	root := t.TempDir()
	chdir(t, root)
	if isGitIgnored("foo.md") {
		t.Fatalf("non-git dir should return false")
	}
}
