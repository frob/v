package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
)

func chdir(t *testing.T, dir string) {
	t.Helper()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(orig); err != nil {
			t.Errorf("restore chdir: %v", err)
		}
	})
}

func testCmd() *cobra.Command {
	cmd := &cobra.Command{}
	cmd.SetOut(&bytes.Buffer{})
	return cmd
}

// mockCommit replaces resolveCommit with a stub returning the given ref and hash.
func mockCommit(t *testing.T, ref, hash string) {
	t.Helper()
	orig := resolveCommit
	resolveCommit = func(_, _ string) (string, string, error) { return ref, hash, nil }
	t.Cleanup(func() { resolveCommit = orig })
}

// mockDownload replaces downloadRepo with a no-op stub.
func mockDownload(t *testing.T) {
	t.Helper()
	orig := downloadRepo
	downloadRepo = func(_, _, _ string) error { return nil }
	t.Cleanup(func() { downloadRepo = orig })
}

func readVendors(t *testing.T) vendorsConfig {
	t.Helper()
	data, err := os.ReadFile(vendorsFile)
	if err != nil {
		t.Fatalf("read %s: %v", vendorsFile, err)
	}
	cfg := make(vendorsConfig)
	if err := toml.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("parse %s: %v", vendorsFile, err)
	}
	return cfg
}

const fakeHash = "aabbccddaabbccddaabbccddaabbccddaabbccdd"

func TestAdd_URLOnly_DefaultBranch(t *testing.T) {
	chdir(t, t.TempDir())
	mockCommit(t, "main", fakeHash) // resolveCommit returns the default branch name
	mockDownload(t)

	if err := runAdd(testCmd(), []string{"https://github.com/example/repo"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entry, ok := readVendors(t)["https://github.com/example/repo"]
	if !ok {
		t.Fatal("entry not found in vendors.toml")
	}
	if entry.Ref != "main" {
		t.Errorf("ref = %q, want %q", entry.Ref, "main")
	}
	if entry.Commit != fakeHash {
		t.Errorf("commit = %q, want %q", entry.Commit, fakeHash)
	}
}

func TestAdd_WithRef(t *testing.T) {
	chdir(t, t.TempDir())
	mockCommit(t, "v1.0.0", fakeHash)
	mockDownload(t)

	if err := runAdd(testCmd(), []string{"https://github.com/example/repo", "v1.0.0"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entry := readVendors(t)["https://github.com/example/repo"]
	if entry.Ref != "v1.0.0" {
		t.Errorf("ref = %q, want %q", entry.Ref, "v1.0.0")
	}
	if entry.Commit != fakeHash {
		t.Errorf("commit = %q, want %q", entry.Commit, fakeHash)
	}
}

func TestAdd_DefaultPath(t *testing.T) {
	chdir(t, t.TempDir())
	mockCommit(t, "main", fakeHash)

	var gotDest string
	orig := downloadRepo
	downloadRepo = func(_, _, dest string) error { gotDest = dest; return nil }
	t.Cleanup(func() { downloadRepo = orig })

	if err := runAdd(testCmd(), []string{"https://github.com/example/repo"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := filepath.Join("vendor", "github.com", "example", "repo")
	if gotDest != want {
		t.Errorf("dest = %q, want %q", gotDest, want)
	}
	if readVendors(t)["https://github.com/example/repo"].Path != want {
		t.Errorf("path in vendors.toml = %q, want %q", readVendors(t)["https://github.com/example/repo"].Path, want)
	}
}

func TestAdd_CustomDestination(t *testing.T) {
	chdir(t, t.TempDir())
	mockCommit(t, "main", fakeHash)

	var gotDest string
	orig := downloadRepo
	downloadRepo = func(_, _, dest string) error { gotDest = dest; return nil }
	t.Cleanup(func() { downloadRepo = orig })

	addDestination = "libs/myrepo"
	t.Cleanup(func() { addDestination = "" })

	if err := runAdd(testCmd(), []string{"https://github.com/example/repo"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotDest != "libs/myrepo" {
		t.Errorf("dest = %q, want %q", gotDest, "libs/myrepo")
	}
	if readVendors(t)["https://github.com/example/repo"].Path != "libs/myrepo" {
		t.Errorf("path in vendors.toml = %q, want %q", readVendors(t)["https://github.com/example/repo"].Path, "libs/myrepo")
	}
}

func TestAdd_MultipleEntries(t *testing.T) {
	chdir(t, t.TempDir())
	mockCommit(t, "main", fakeHash)
	mockDownload(t)

	if err := runAdd(testCmd(), []string{"https://github.com/foo/one", "v1.0.0"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := runAdd(testCmd(), []string{"https://github.com/foo/two", "abc123"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cfg := readVendors(t)
	if len(cfg) != 2 {
		t.Fatalf("len = %d, want 2", len(cfg))
	}
	if cfg["https://github.com/foo/one"].Ref != "main" {
		t.Errorf("foo/one ref = %q, want %q", cfg["https://github.com/foo/one"].Ref, "main")
	}
	if cfg["https://github.com/foo/two"].Ref != "main" {
		t.Errorf("foo/two ref = %q, want %q", cfg["https://github.com/foo/two"].Ref, "main")
	}
}

func TestAdd_UpdateExisting(t *testing.T) {
	chdir(t, t.TempDir())
	mockDownload(t)

	hash1 := "1111111111111111111111111111111111111111"
	hash2 := "2222222222222222222222222222222222222222"

	orig := resolveCommit
	t.Cleanup(func() { resolveCommit = orig })

	resolveCommit = func(_, _ string) (string, string, error) { return "v1.0.0", hash1, nil }
	if err := runAdd(testCmd(), []string{"https://github.com/example/repo", "v1.0.0"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	resolveCommit = func(_, _ string) (string, string, error) { return "v2.0.0", hash2, nil }
	if err := runAdd(testCmd(), []string{"https://github.com/example/repo", "v2.0.0"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cfg := readVendors(t)
	if len(cfg) != 1 {
		t.Fatalf("len = %d, want 1 (duplicate entry created)", len(cfg))
	}
	entry := cfg["https://github.com/example/repo"]
	if entry.Ref != "v2.0.0" {
		t.Errorf("ref = %q, want %q", entry.Ref, "v2.0.0")
	}
	if entry.Commit != hash2 {
		t.Errorf("commit = %q, want %q", entry.Commit, hash2)
	}
}

func TestAdd_CommitResolutionError(t *testing.T) {
	chdir(t, t.TempDir())

	orig := resolveCommit
	resolveCommit = func(_, _ string) (string, string, error) { return "", "", fmt.Errorf("no network") }
	t.Cleanup(func() { resolveCommit = orig })

	if err := runAdd(testCmd(), []string{"https://github.com/example/repo"}); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAdd_DownloadError(t *testing.T) {
	chdir(t, t.TempDir())
	mockCommit(t, "main", fakeHash)

	orig := downloadRepo
	downloadRepo = func(_, _, _ string) error { return fmt.Errorf("network failure") }
	t.Cleanup(func() { downloadRepo = orig })

	if err := runAdd(testCmd(), []string{"https://github.com/example/repo"}); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDefaultPath(t *testing.T) {
	cases := []struct {
		url  string
		want string
	}{
		{
			url:  "https://github.com/foo/bar",
			want: filepath.Join("vendor", "github.com", "foo", "bar"),
		},
		{
			url:  "https://github.com/foo/bar.git",
			want: filepath.Join("vendor", "github.com", "foo", "bar"),
		},
		{
			url:  "https://gitlab.com/org/project",
			want: filepath.Join("vendor", "gitlab.com", "org", "project"),
		},
	}

	for _, tc := range cases {
		got, err := defaultPath(tc.url)
		if err != nil {
			t.Errorf("defaultPath(%q) error: %v", tc.url, err)
			continue
		}
		if got != tc.want {
			t.Errorf("defaultPath(%q) = %q, want %q", tc.url, got, tc.want)
		}
	}
}

func TestIsCommitHash(t *testing.T) {
	cases := []struct {
		name string
		s    string
		want bool
	}{
		{name: "valid lowercase", s: "aabbccddaabbccddaabbccddaabbccddaabbccdd", want: true},
		{name: "valid uppercase", s: "AABBCCDDAABBCCDDAABBCCDDAABBCCDDAABBCCDD", want: true},
		{name: "valid mixed case", s: "AaBbCcDd1234567890aabbccddeeff0011223344", want: true},
		{name: "valid all digits", s: "0123456789012345678901234567890123456789", want: true},
		{name: "empty", s: "", want: false},
		{name: "too short", s: "aabbccdd", want: false},
		{name: "too long", s: "aabbccddaabbccddaabbccddaabbccddaabbccdd0", want: false},
		{name: "non-hex letter g", s: "gabbccddaabbccddaabbccddaabbccddaabbccdd", want: false},
		{name: "non-hex symbol", s: "-abbccddaabbccddaabbccddaabbccddaabbccdd", want: false},
		{name: "ref name not a hash", s: "main", want: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := isCommitHash(tc.s); got != tc.want {
				t.Errorf("isCommitHash(%q) = %v, want %v", tc.s, got, tc.want)
			}
		})
	}
}

// initSourceRepo creates a real single-commit git repository in a temp dir using
// go-git's in-process API (no network, no git binary) and returns its path and
// the commit hash. It serves as a local "remote" for cloneRepo tests.
func initSourceRepo(t *testing.T) (dir, commit string) {
	t.Helper()
	dir = t.TempDir()
	repo, err := git.PlainInit(dir, false)
	if err != nil {
		t.Fatalf("init: %v", err)
	}
	w, err := repo.Worktree()
	if err != nil {
		t.Fatalf("worktree: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("hello\n"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	if _, err := w.Add("README.md"); err != nil {
		t.Fatalf("add: %v", err)
	}
	h, err := w.Commit("initial commit", &git.CommitOptions{
		Author: &object.Signature{Name: "Test", Email: "test@example.com", When: time.Unix(0, 0).UTC()},
	})
	if err != nil {
		t.Fatalf("commit: %v", err)
	}
	return dir, h.String()
}

func TestCloneRepo(t *testing.T) {
	src, commit := initSourceRepo(t)
	dest := filepath.Join(t.TempDir(), "clone")

	if err := cloneRepo(src, commit, dest); err != nil {
		t.Fatalf("cloneRepo: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dest, "README.md")); err != nil {
		t.Errorf("expected README.md in clone: %v", err)
	}
	// cloneRepo strips the .git directory so vendored trees carry no history.
	if _, err := os.Stat(filepath.Join(dest, ".git")); !os.IsNotExist(err) {
		t.Errorf(".git should be removed, stat err = %v", err)
	}
}

func TestCloneRepo_DestinationExists(t *testing.T) {
	src, commit := initSourceRepo(t)
	dest := t.TempDir()
	// Pre-populate dest with a repo so PlainClone reports it already exists.
	if _, err := git.PlainInit(dest, false); err != nil {
		t.Fatalf("init dest: %v", err)
	}

	err := cloneRepo(src, commit, dest)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("error = %q, want it to mention %q", err, "already exists")
	}
}

func TestCloneRepo_CloneError(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "does-not-exist")
	dest := filepath.Join(t.TempDir(), "clone")

	if err := cloneRepo(missing, fakeHash, dest); err == nil {
		t.Fatal("expected error cloning from a nonexistent source, got nil")
	}
}

// stubRemoteRefs replaces listRemoteRefs with a fixture for the duration of a test.
func stubRemoteRefs(t *testing.T, refs []*plumbing.Reference, listErr error) {
	t.Helper()
	orig := listRemoteRefs
	listRemoteRefs = func(string) ([]*plumbing.Reference, error) { return refs, listErr }
	t.Cleanup(func() { listRemoteRefs = orig })
}

func TestRemoteCommit(t *testing.T) {
	const (
		mainHash = "1111111111111111111111111111111111111111"
		devHash  = "2222222222222222222222222222222222222222"
		tagHash  = "3333333333333333333333333333333333333333"
		otherHsh = "9999999999999999999999999999999999999999"
	)
	hashRef := func(name, hash string) *plumbing.Reference {
		return plumbing.NewHashReference(plumbing.ReferenceName(name), plumbing.NewHash(hash))
	}
	headTo := func(target string) *plumbing.Reference {
		return plumbing.NewSymbolicReference(plumbing.HEAD, plumbing.ReferenceName(target))
	}

	cases := []struct {
		name     string
		inputRef string
		refs     []*plumbing.Reference
		listErr  error
		wantRef  string
		wantHash string
		wantErr  bool
	}{
		{
			name:     "default branch via HEAD symref",
			inputRef: "",
			refs:     []*plumbing.Reference{headTo("refs/heads/main"), hashRef("refs/heads/main", mainHash)},
			wantRef:  "main",
			wantHash: mainHash,
		},
		{
			name:     "branch ref",
			inputRef: "dev",
			refs:     []*plumbing.Reference{hashRef("refs/heads/dev", devHash)},
			wantRef:  "dev",
			wantHash: devHash,
		},
		{
			name:     "lightweight tag",
			inputRef: "v2.0.0",
			refs:     []*plumbing.Reference{hashRef("refs/tags/v2.0.0", tagHash)},
			wantRef:  "v2.0.0",
			wantHash: tagHash,
		},
		{
			name:     "annotated tag prefers peeled ref",
			inputRef: "v1.0.0",
			refs:     []*plumbing.Reference{hashRef("refs/tags/v1.0.0", otherHsh), hashRef("refs/tags/v1.0.0^{}", tagHash)},
			wantRef:  "v1.0.0",
			wantHash: tagHash,
		},
		{
			name:     "explicit commit hash passes through",
			inputRef: mainHash,
			wantRef:  mainHash,
			wantHash: mainHash,
		},
		{
			name:     "ref not found",
			inputRef: "nope",
			refs:     []*plumbing.Reference{hashRef("refs/heads/main", mainHash)},
			wantErr:  true,
		},
		{
			name:     "list error",
			inputRef: "",
			listErr:  fmt.Errorf("no network"),
			wantErr:  true,
		},
		{
			name:     "HEAD missing for default branch",
			inputRef: "",
			refs:     []*plumbing.Reference{hashRef("refs/heads/main", mainHash)},
			wantErr:  true,
		},
		{
			name:     "HEAD target branch missing",
			inputRef: "",
			refs:     []*plumbing.Reference{headTo("refs/heads/ghost")},
			wantErr:  true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			stubRemoteRefs(t, tc.refs, tc.listErr)

			gotRef, gotHash, err := remoteCommit("https://example.com/org/repo", tc.inputRef)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got ref=%q hash=%q", gotRef, gotHash)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if gotRef != tc.wantRef {
				t.Errorf("ref = %q, want %q", gotRef, tc.wantRef)
			}
			if gotHash != tc.wantHash {
				t.Errorf("hash = %q, want %q", gotHash, tc.wantHash)
			}
		})
	}
}
