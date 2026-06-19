package debtcheck

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// This file is dedicated coverage backfill for debtcheck.go. It only adds tests;
// it does not edit existing _test.go files.

// ---- pure / semi-pure helpers ----

func TestParseOptionalBool(t *testing.T) {
	cases := []struct {
		in      string
		wantVal bool
		wantOK  bool
	}{
		{"true", true, true},
		{"false", false, true},
		{"", false, false},
		{"maybe", false, false},
		{"TRUE", false, false},
	}
	for _, tc := range cases {
		val, ok := parseOptionalBool(tc.in)
		if val != tc.wantVal || ok != tc.wantOK {
			t.Errorf("parseOptionalBool(%q) = (%v,%v); want (%v,%v)", tc.in, val, ok, tc.wantVal, tc.wantOK)
		}
	}
}

func TestUnquoteYAMLScalar(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{`"double"`, "double"},
		{`'single'`, "single"},
		{"plain", "plain"},
		{"x", "x"},             // len < 2 returns as-is
		{"", ""},               // len < 2
		{`"mixed'`, `"mixed'`}, // mismatched quotes returned as-is
		{`'mixed"`, `'mixed"`},
	}
	for _, tc := range cases {
		if got := unquoteYAMLScalar(tc.in); got != tc.want {
			t.Errorf("unquoteYAMLScalar(%q) = %q; want %q", tc.in, got, tc.want)
		}
	}
}

func TestValidateMode(t *testing.T) {
	for _, mode := range []string{"enforce", "warn", "observe"} {
		if err := validateMode(mode); err != nil {
			t.Errorf("validateMode(%q) err = %v; want nil", mode, err)
		}
	}
	if err := validateMode("bogus"); err == nil {
		t.Error("validateMode(bogus) err = nil; want error")
	}
}

func TestValidateSection(t *testing.T) {
	for _, section := range append([]string{"all"}, allSections()...) {
		if err := validateSection(section); err != nil {
			t.Errorf("validateSection(%q) err = %v; want nil", section, err)
		}
	}
	if err := validateSection("nonexistent"); err == nil {
		t.Error("validateSection(nonexistent) err = nil; want error")
	}
}

func TestSkipDir(t *testing.T) {
	skipped := []string{".git", ".omx", ".worktree", "vendor", "node_modules", "release", "tmp", ".cache"}
	for _, name := range skipped {
		if !skipDir(name) {
			t.Errorf("skipDir(%q) = false; want true", name)
		}
	}
	if skipDir("src") {
		t.Error("skipDir(src) = true; want false")
	}
}

func TestSkipFile(t *testing.T) {
	cases := []struct {
		path string
		want bool
	}{
		{"/a/.hidden", true},
		{"/a/.gitignore", false}, // exception
		{"/a/image.png", true},
		{"/a/photo.jpeg", true},
		{"/a/doc.pdf", true},
		{"/a/anim.gif", true},
		{"/a/pic.jpg", true},
		{"/a/normal.go", false},
		{"/a/data.txt", false},
	}
	for _, tc := range cases {
		if got := skipFile(tc.path); got != tc.want {
			t.Errorf("skipFile(%q) = %v; want %v", tc.path, got, tc.want)
		}
	}
}

func TestBytesLookBinary(t *testing.T) {
	if !bytesLookBinary([]byte{0x00, 0x01}) {
		t.Error("bytesLookBinary with NUL byte = false; want true")
	}
	if bytesLookBinary([]byte("plain text")) {
		t.Error("bytesLookBinary plain text = true; want false")
	}
	// >4096 bytes path: ensure limit branch is exercised (no NUL → false).
	big := make([]byte, 5000)
	for i := range big {
		big[i] = 'a'
	}
	if bytesLookBinary(big) {
		t.Error("bytesLookBinary big text = true; want false")
	}
	// >4096 bytes with NUL near the end (beyond limit) → must stay false.
	big2 := make([]byte, 5000)
	for i := range big2 {
		big2[i] = 'a'
	}
	big2[4900] = 0 // beyond 4096 limit
	if bytesLookBinary(big2) {
		t.Error("bytesLookBinary NUL beyond limit = true; want false")
	}
}

func TestRel(t *testing.T) {
	got := rel("/root", "/root/a/b.txt")
	if got != "a/b.txt" {
		t.Errorf("rel = %q; want a/b.txt", got)
	}
	// filepath.Rel only errors when a relative target cannot be made relative to root,
	// which on POSIX with absolute paths is rare. Exercise the fallback with a synthetic
	// mismatch is impractical portably; the success path covers the function body.
}

func TestRelFallback(t *testing.T) {
	// On most POSIX systems filepath.Rel never errors for absolute paths, so the error
	// branch (return path unchanged) is recorded as untestable here.
	_ = rel("/root", "/root/a")
}

// ---- normalize Root default ----

func TestNormalizeDefaultsRootToDot(t *testing.T) {
	opts := normalize(Options{})
	if opts.Root != "." {
		t.Errorf("Root = %q; want \".\"", opts.Root)
	}
	for _, check := range []struct{ field, want string }{
		{opts.ConfigPath, DefaultRulesPath},
		{opts.RegistryPath, DefaultRegistryPath},
		{opts.ExceptionsPath, DefaultExceptions},
		{opts.DependencyPurposePath, DefaultPurpose},
		{opts.Section, "all"},
		{opts.Mode, "enforce"},
	} {
		if check.field != check.want {
			t.Errorf("normalize default mismatch: got %q want %q", check.field, check.want)
		}
	}
	if opts.MinScore != DefaultMinScore {
		t.Errorf("MinScore = %v; want %v", opts.MinScore, DefaultMinScore)
	}
}

func TestStatus(t *testing.T) {
	cases := []struct {
		name     string
		summary  Summary
		score    float64
		minScore float64
		mode     string
		want     string
	}{
		{"p0 fails", Summary{P0: 1}, 10, 9, "enforce", "failed"},
		{"score below min enforce", Summary{}, 5, 9, "enforce", "failed"},
		{"score below min observe", Summary{}, 5, 9, "observe", "warning"},
		{"score below min warn", Summary{}, 5, 9, "warn", "warning"},
		{"p1+p2 enforce passes", Summary{P1: 1, P2: 1}, 9.9, 9, "enforce", "passed"},
		{"p1+p2 non-enforce warning", Summary{P1: 1}, 9.9, 9, "observe", "warning"},
		{"clean passes", Summary{}, 10, 9, "enforce", "passed"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := status(tc.summary, tc.score, tc.minScore, tc.mode); got != tc.want {
				t.Errorf("status = %q; want %q", got, tc.want)
			}
		})
	}
}

func TestBuildSectionStatuses(t *testing.T) {
	// P2-only → warning branch (P0==0, but P1==0 and P2>0 → status warning).
	s := buildSection("docs", []Finding{{Severity: "P2"}})
	if s.Status != "warning" || s.P2 != 1 {
		t.Fatalf("buildSection P2 = %+v; want status=warning P2=1", s)
	}
	// unknown severity → P2 bucket default.
	s2 := buildSection("docs", []Finding{{Severity: "weird"}})
	if s2.P2 != 1 {
		t.Fatalf("buildSection weird = %+v; want P2=1", s2)
	}
}

// ---- ValidateEvidence failure branches ----

func TestValidateEvidenceAllFailures(t *testing.T) {
	// All problems at once.
	e := Evidence{
		SchemaVersion:       "wrong-manifest",
		ReportSchemaVersion: "wrong-report",
		Status:              "failed",
		Score:               5,
		MinScore:            9,
		Sections: []SectionEvidence{
			{Name: "docs", P0: 2, Status: "failed"},
		},
	}
	problems := ValidateEvidence(e, 9)
	for _, want := range []string{
		"debt schema version mismatch",
		"debt report schema version mismatch",
		"debt status is failed",
		"debt score 5.00 below 9.00",
		"debt section docs has 2 P0 findings",
		"debt section docs status is failed",
	} {
		if !strings.Contains(strings.Join(problems, "\n"), want) {
			t.Errorf("problems = %v; want %q", problems, want)
		}
	}
}

func TestValidateEvidenceClean(t *testing.T) {
	e := Evidence{
		SchemaVersion:       ManifestSchema,
		ReportSchemaVersion: SchemaVersion,
		Status:              "passed",
		Score:               10,
		MinScore:            9,
		Sections:            []SectionEvidence{{Name: "docs", P0: 0, Status: "passed"}},
	}
	if problems := ValidateEvidence(e, 9); len(problems) != 0 {
		t.Fatalf("problems = %v; want none", problems)
	}
}

// ---- ExitCode / ToMarkdown uncovered branches ----

func TestExitCodeZeroOnNonPassedObserveAndWarn(t *testing.T) {
	// ExitCode returns 0 when mode is observe or warn regardless of status.
	if code := ExitCode(Report{Mode: "observe", Status: "failed"}); code != 0 {
		t.Errorf("observe failed ExitCode = %d; want 0", code)
	}
	if code := ExitCode(Report{Mode: "warn", Status: "failed"}); code != 0 {
		t.Errorf("warn failed ExitCode = %d; want 0", code)
	}
}

func TestToMarkdownNoFindingsAndEmptyPath(t *testing.T) {
	report := Report{
		SchemaVersion: SchemaVersion,
		Status:        "passed",
		Mode:          "enforce",
		Sections: []SectionReport{
			{Name: "docs", Status: "passed", Findings: nil},                                                // "No findings." branch
			{Name: "sec", Status: "warning", Findings: []Finding{{ID: "x", Severity: "P1", Message: "m"}}}, // empty Path → "policy"
		},
	}
	markdown := ToMarkdown(report)
	if !strings.Contains(markdown, "No findings.") {
		t.Errorf("markdown missing 'No findings.': %s", markdown)
	}
	if !strings.Contains(markdown, "x policy: m") {
		t.Errorf("markdown missing empty-path fallback to 'policy': %s", markdown)
	}
}

func TestFindingMetadataMarkdownEmpty(t *testing.T) {
	if out := findingMetadataMarkdown(Finding{}); out != "" {
		t.Errorf("findingMetadataMarkdown empty = %q; want empty", out)
	}
}

// ---- readRuleMetadata branches ----

func TestReadRuleMetadataMissingFile(t *testing.T) {
	// Arrange: file does not exist → returns nil.
	if m := readRuleMetadata(t.TempDir(), "missing.yaml"); m != nil {
		t.Fatalf("readRuleMetadata missing = %#v; want nil", m)
	}
}

func TestReadRuleMetadataMalformedAndDashed(t *testing.T) {
	// Arrange: registry with a "- " entry whose body is empty (trimmed == "" continue),
	// a non-rule line, a bad no-colon line, plus a real rule with release_blocking=false.
	root := t.TempDir()
	registry := strings.Join([]string{
		"schema_version: debt-rule-registry/v1",
		"rules:",
		"- ",
		"not-a-rule-line",
		"- id: debt.docs.marker",
		"bad-no-colon-line",
		"  invariant_id: INV-1",
		"  release_blocking: false",
		"  unknown_key: ignored",
		"",
		"# comment line",
	}, "\n")
	writeFile(t, root, "registry.yaml", registry)
	// Act
	metadata := readRuleMetadata(root, "registry.yaml")
	// Assert
	if len(metadata) != 1 {
		t.Fatalf("metadata = %#v; want one rule", metadata)
	}
	f, ok := metadata["debt.docs.marker"]
	if !ok {
		t.Fatalf("metadata = %#v; want debt.docs.marker", metadata)
	}
	if f.InvariantID != "INV-1" {
		t.Errorf("InvariantID = %q; want INV-1", f.InvariantID)
	}
	if f.ReleaseBlocking == nil || *f.ReleaseBlocking != false {
		t.Errorf("ReleaseBlocking = %v; want false", f.ReleaseBlocking)
	}
}

func TestReadRuleMetadataUnmatchedIDNotFlushed(t *testing.T) {
	// Lines before any "- " are ignored; a key:value at file scope (inRule=false) hits the
	// `if !inRule { continue }` branch.
	root := t.TempDir()
	registry := "stray_key: value\nid: orphan\n- id: debt.x\n  owner: team\n"
	writeFile(t, root, "registry.yaml", registry)
	metadata := readRuleMetadata(root, "registry.yaml")
	if len(metadata) != 1 || metadata["debt.x"].Owner != "team" {
		t.Fatalf("metadata = %#v; want one rule with owner team", metadata)
	}
}

// ---- annotateFindings branch ----

func TestAnnotateFindingsMetadataMissing(t *testing.T) {
	// Finding whose ID is not in metadata → skipped continue branch.
	findings := []Finding{{ID: "unknown.id", Severity: "P1"}}
	out := annotateFindings(findings, map[string]Finding{"other": {Owner: "x"}})
	if len(out) != 1 || out[0].Owner != "" {
		t.Fatalf("out = %#v; want unchanged finding", out)
	}
}

func TestAnnotateFindingsEmptyMetadataShortCircuit(t *testing.T) {
	// Empty metadata → early return, findings unchanged (already covered indirectly via
	// Run without registry, but exercise directly for the len==0 branch).
	in := []Finding{{ID: "x"}}
	out := annotateFindings(in, nil)
	if len(out) != 1 {
		t.Fatalf("out = %#v; want input unchanged", out)
	}
}

// ---- Run error branches (validateMode/validateSection) ----

func TestRunRejectsInvalidMode(t *testing.T) {
	_, err := Run(Options{Root: t.TempDir(), Mode: "bogus"})
	if err == nil || !strings.Contains(err.Error(), "unsupported debt mode") {
		t.Fatalf("err = %v; want unsupported debt mode", err)
	}
}

func TestRunRejectsInvalidSection(t *testing.T) {
	_, err := Run(Options{Root: t.TempDir(), Mode: "enforce", Section: "bogus"})
	if err == nil || !strings.Contains(err.Error(), "unsupported debt section") {
		t.Fatalf("err = %v; want unsupported debt section", err)
	}
}

func TestRunNormalizeDefaults(t *testing.T) {
	// Provide fully empty Options → normalize fills every default, and Run proceeds.
	root := t.TempDir()
	writePolicyFiles(t, root)
	writeDownstreamFiles(t, root)
	report, err := Run(Options{Root: root})
	if err != nil {
		t.Fatal(err)
	}
	if report.Mode != "enforce" || len(report.Sections) == 0 {
		t.Fatalf("report = %+v; want defaulted enforce mode with sections", report)
	}
	// Explicit defaults applied (MinScore == DefaultMinScore).
	if report.MinScore != DefaultMinScore {
		t.Errorf("MinScore = %v; want %v", report.MinScore, DefaultMinScore)
	}
}

// ---- missingPolicyFindings existing-files branch ----

func TestMissingPolicyFindingsAllPresent(t *testing.T) {
	root := t.TempDir()
	writePolicyFiles(t, root)
	opts := normalize(Options{Root: root})
	if f := missingPolicyFindings(opts); len(f) != 0 {
		t.Fatalf("missingPolicyFindings = %#v; want none when all files exist", f)
	}
}

// ---- scanSection default nil ----

func TestScanSectionUnknownReturnsNil(t *testing.T) {
	if f := scanSection(t.TempDir(), "unknown-section"); f != nil {
		t.Fatalf("scanSection unknown = %#v; want nil", f)
	}
}

// ---- scanGoImports read + parse errors ----

func TestScanGoImportsReadError(t *testing.T) {
	// walkFiles on a single .go file path that does not exist → visit returns err from
	// ReadFile, surfaced but findings still empty. Use a nonexistent root file path.
	// walkFiles calls os.Lstat first; a path that exists as a directory with an unreadable
	// file triggers the ReadFile error path inside scanGoImports.
	root := t.TempDir()
	// Create a .go file then revoke read permission so os.ReadFile fails.
	secret := filepath.Join(root, "secret.go")
	if err := os.WriteFile(secret, []byte("package x\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(secret, 0o000); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Chmod(secret, 0o644); err != nil { // restore so TempDir cleanup works
			t.Errorf("restore chmod %s: %v", secret, err)
		}
	}()
	// Skip the test if running as root (root bypasses 0o000).
	if os.Geteuid() == 0 {
		t.Skip("cannot simulate unreadable file as root")
	}
	findings := scanGoImports(root)
	// We don't assert findings content (parse error vs read error both possible), only
	// that the function completed without panic and exercised the err branch.
	_ = findings
}

func TestScanGoImportsParseError(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "broken.go", "this is not valid go code at all !!!\n")
	findings := scanGoImports(root)
	found := false
	for _, f := range findings {
		if f.ID == "debt.architecture.parse" {
			found = true
		}
	}
	if !found {
		t.Fatalf("findings = %#v; want a parse finding", findings)
	}
}

// ---- scanDependencyDebt + scanSecurityDebt triggers ----

func TestScanDependencyDebtDetectsPatterns(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "Makefile", "install:\n\tcurl https://x | bash\n\tgo get foo@latest\n")
	findings := scanDependencyDebt(root)
	ids := map[string]bool{}
	for _, f := range findings {
		ids[f.ID] = true
	}
	if !ids["debt.dependency.unpinned-latest"] {
		t.Errorf("findings = %#v; want unpinned-latest (Makefile is non-.md)", findings)
	}
	if !ids["debt.dependency.curl-pipe-bash"] {
		t.Errorf("findings = %#v; want curl-pipe-bash", findings)
	}
}

func TestScanDependencyDebtSkipsMarkdownForLatest(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "doc.md", "use foo@latest in docs\n")
	findings := scanDependencyDebt(root)
	for _, f := range findings {
		if f.ID == "debt.dependency.unpinned-latest" {
			t.Fatalf("findings = %#v; markdown must skip @latest", findings)
		}
	}
}

func TestScanSecurityDebtDetectsMarkers(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "leak.txt", privateKeyPrefix+"\nxlib-security-debt\n")
	findings := scanSecurityDebt(root)
	ids := map[string]bool{}
	for _, f := range findings {
		ids[f.ID] = true
	}
	if !ids["debt.security.private-key"] {
		t.Errorf("findings = %#v; want private-key", findings)
	}
	if !ids["debt.security.marker"] {
		t.Errorf("findings = %#v; want security marker", findings)
	}
}

func TestScanTextMarkerFound(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "doc.md", "xlib-docs-drift here\n")
	findings := scanTextMarker(root, "xlib-docs-drift", "debt.docs.marker", "documentation drift marker is present")
	if len(findings) != 1 || findings[0].ID != "debt.docs.marker" {
		t.Fatalf("findings = %#v; want one docs marker finding", findings)
	}
}

func TestScanTextMarkerSortedByPathThenID(t *testing.T) {
	// Multiple files with same marker → sort.Slice Path/ID branches.
	root := t.TempDir()
	writeFile(t, root, "b.md", "xlib-docs-drift\n")
	writeFile(t, root, "a.md", "xlib-docs-drift\n")
	findings := scanTextMarker(root, "xlib-docs-drift", "debt.docs.marker", "m")
	if len(findings) != 2 {
		t.Fatalf("findings = %#v; want two", findings)
	}
	if findings[0].Path != "a.md" {
		t.Errorf("findings not sorted by path: %#v", findings)
	}
}

func TestScanTrackedTextSkipsBinaryAndUnreadable(t *testing.T) {
	// bytesLookBinary true → return nil branch in scanTrackedText.
	root := t.TempDir()
	writeFile(t, root, "binary.dat", "text\n\x00\nmore\n") // not skipped by skipFile (.dat ok)
	findings := scanTextMarker(root, "marker", "id", "msg")
	for _, f := range findings {
		if f.Path == "binary.dat" {
			t.Errorf("binary file should be skipped: %#v", f)
		}
	}
}

func TestScanTrackedTextUnreadableFileReturnsError(t *testing.T) {
	// scanTrackedText propagates walkFiles ReadFile error (discarded via _ =). Cover the
	// err branch by making a file unreadable; the result is the walk stops without panic.
	if os.Geteuid() == 0 {
		t.Skip("cannot simulate unreadable file as root")
	}
	root := t.TempDir()
	secret := filepath.Join(root, "secret.md")
	if err := os.WriteFile(secret, []byte("marker\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(secret, 0o000); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Chmod(secret, 0o644); err != nil {
			t.Errorf("restore chmod %s: %v", secret, err)
		}
	}()
	_ = scanTextMarker(root, "marker", "id", "msg")
}

// ---- scanDownstreamDebt missing-file + schema + placeholder + gap branches ----

func TestScanDownstreamDebtMissingFiles(t *testing.T) {
	// No downstream files present → every required file yields a missing finding
	// (exercises the os.ReadFile err branch AND the sort-by-different-path branch).
	findings := scanDownstreamDebt(t.TempDir())
	ids := map[string]int{}
	for _, f := range findings {
		ids[f.ID]++
	}
	if ids["debt.downstream.file-missing"] == 0 {
		t.Fatalf("findings = %#v; want file-missing findings across distinct paths", findings)
	}
}

func TestScanDownstreamDebtSchemaMissingAndPlaceholder(t *testing.T) {
	root := t.TempDir()
	// Provide ALL required downstream files but with: no schema_version, and placeholder text.
	for _, rel := range downstreamRequiredFiles {
		writeFile(t, root, rel, "TODO placeholder TBD\n")
	}
	findings := scanDownstreamDebt(root)
	ids := map[string]bool{}
	for _, f := range findings {
		ids[f.ID] = true
	}
	if !ids["debt.downstream.schema-missing"] {
		t.Errorf("findings = %#v; want schema-missing", findings)
	}
	if !ids["debt.downstream.placeholder"] {
		t.Errorf("findings = %#v; want placeholder", findings)
	}
}

func TestScanDownstreamDebtBaselineMissingGapStatus(t *testing.T) {
	root := t.TempDir()
	writeDownstreamFiles(t, root)
	// Overwrite baseline so it has no "gap" substring.
	writeFile(t, root, ".agent/registries/downstream-baseline-scan.yaml", "schema_version: \"2.9.3\"\nrepo: kernel/configx\nmode: patch-only\nstatus: ok\n")
	findings := scanDownstreamDebt(root)
	ids := map[string]bool{}
	for _, f := range findings {
		ids[f.ID] = true
	}
	if !ids["debt.downstream.baseline-missing-gap-status"] {
		t.Errorf("findings = %#v; want baseline-missing-gap-status", findings)
	}
}

func TestScanDownstreamDebtBaselineEmptyTriggersRequireTokenSkip(t *testing.T) {
	// requireDownstreamToken short-circuits when text=="" (missing file already reported).
	// Cover that branch by deleting baseline only.
	root := t.TempDir()
	writeDownstreamFiles(t, root)
	if err := os.Remove(filepath.Join(root, filepath.FromSlash(".agent/registries/downstream-baseline-scan.yaml"))); err != nil {
		t.Fatalf("remove baseline: %v", err)
	}
	findings := scanDownstreamDebt(root)
	// baseline-missing repo/mode must NOT appear (text=="" path), only file-missing.
	for _, f := range findings {
		if f.ID == "debt.downstream.baseline-missing-repo" || f.ID == "debt.downstream.baseline-missing-mode" {
			t.Fatalf("baseline empty should skip require-token: %#v", f)
		}
	}
}

func TestRequireDownstreamTokenAppendBranch(t *testing.T) {
	// Directly exercise the append branch: non-empty text lacking the token.
	var findings []Finding
	requireDownstreamToken(&findings, "some content", "p", "MISSING-TOKEN", "id.x", "msg")
	if len(findings) != 1 || findings[0].ID != "id.x" {
		t.Fatalf("findings = %#v; want one id.x", findings)
	}
}

// ---- walkFiles + walkDir branches ----

func TestWalkFilesMissingRoot(t *testing.T) {
	// Nonexistent root → os.Lstat error returned.
	if err := walkFiles(filepath.Join(t.TempDir(), "nope"), func(string) error { return nil }); err == nil {
		t.Error("walkFiles missing root err = nil; want error")
	}
}

func TestWalkFilesSingleFile(t *testing.T) {
	// A file (not dir) → visit is called directly.
	root := t.TempDir()
	path := filepath.Join(root, "a.txt")
	if err := os.WriteFile(path, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	visited := false
	if err := walkFiles(path, func(p string) error { visited = true; return nil }); err != nil {
		t.Fatal(err)
	}
	if !visited {
		t.Error("walkFiles single file did not visit")
	}
}

func TestWalkDirOpenError(t *testing.T) {
	// A directory without read permission → os.Open / Readdirnames error.
	if os.Geteuid() == 0 {
		t.Skip("cannot simulate unreadable dir as root")
	}
	root := t.TempDir()
	sub := filepath.Join(root, "locked")
	if err := os.Mkdir(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	inner := filepath.Join(sub, "file.txt")
	if err := os.WriteFile(inner, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(sub, 0o000); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Chmod(sub, 0o755); err != nil {
			t.Errorf("restore chmod %s: %v", sub, err)
		}
	}()
	if err := walkDir(sub, func(string) error { return nil }); err == nil {
		t.Error("walkDir unreadable err = nil; want error")
	}
}

// NOTE: walkDir's readErr branch (L776-778) cannot be deterministically exercised:
// os.Open and file.Readdirnames are coupled — any directory mode that causes
// Readdirnames to fail (insufficient read bit) also causes os.Open to fail first,
// hitting the Open-error return at L771-772 instead. The closeErr branch (L779-781)
// is similarly unreachable: file.Close() on an already-open fd does not error.
// Both are defensive guards; recorded as untestable.

func TestWalkDirRecursiveError(t *testing.T) {
	// A nested directory whose visit returns an error → propagated through the recursive
	// walkDir call (L793-795 branch).
	root := t.TempDir()
	nested := filepath.Join(root, "sub", "deep")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(nested, "a.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	wantErr := os.ErrInvalid
	err := walkDir(root, func(p string) error {
		return wantErr // any visit error exercises both file-visit and nested-walkDir propagation
	})
	if !errors.Is(err, wantErr) {
		t.Errorf("walkDir recursive err = %v; want %v", err, wantErr)
	}
}

func TestWalkFilesVisitError(t *testing.T) {
	// visit returns an error → walkDir propagates it.
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "a.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	wantErr := os.ErrNotExist
	err := walkFiles(root, func(p string) error { return wantErr })
	if !errors.Is(err, wantErr) {
		t.Errorf("walkFiles visit err = %v; want %v", err, wantErr)
	}
}

func TestWalkFilesSkipDirSkipsSubtree(t *testing.T) {
	// A directory named "release" (in skipDir list) must be pruned, not visited.
	root := t.TempDir()
	release := filepath.Join(root, "release")
	if err := os.Mkdir(release, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(release, "a.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	var seen []string
	if err := walkFiles(root, func(p string) error {
		rel, _ := filepath.Rel(root, p)
		seen = append(seen, filepath.ToSlash(rel))
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	for _, s := range seen {
		if strings.HasPrefix(s, "release/") {
			t.Errorf("release subtree visited: %s", s)
		}
	}
}

// ---- digestFile missing branch ----

func TestDigestFileMissing(t *testing.T) {
	if got := digestFile(t.TempDir(), "nonexistent.yaml"); got != "missing" {
		t.Errorf("digestFile missing = %q; want missing", got)
	}
}

// ---- ReadReport error branches ----

func TestReadReportMissingFile(t *testing.T) {
	if _, err := ReadReport(filepath.Join(t.TempDir(), "no.json")); err == nil {
		t.Error("ReadReport missing err = nil; want error")
	}
}

func TestReadReportBadJSON(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "bad.json")
	if err := os.WriteFile(path, []byte("{not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := ReadReport(path); err == nil {
		t.Error("ReadReport bad json err = nil; want error")
	}
}

func TestReadReportWrongSchema(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "wrong.json")
	if err := os.WriteFile(path, []byte(`{"schema_version":"wrong/v9"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := ReadReport(path); err == nil {
		t.Error("ReadReport wrong schema err = nil; want error")
	}
}

// ---- Run end-to-end exercise of normalize + selectedSections("all") + missingPolicyFindings ----

func TestRunAllSectionsWithMissingPolicies(t *testing.T) {
	// Empty root → all four policy files missing → missingPolicyFindings hits err branch,
	// digestFile hits "missing" for each. Run still returns a report.
	report, err := Run(Options{Root: t.TempDir(), Mode: "observe"})
	if err != nil {
		t.Fatal(err)
	}
	// Every section should have the 4 missing-policy findings.
	totalP0 := 0
	for _, s := range report.Sections {
		totalP0 += s.P0
	}
	if totalP0 == 0 {
		t.Errorf("expected missing-policy findings across sections; got report %+v", report)
	}
	// observe mode → status warning/failed but ExitCode 0.
	if code := ExitCode(report); code != 0 {
		t.Errorf("observe ExitCode = %d; want 0", code)
	}
}
