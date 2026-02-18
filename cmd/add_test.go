package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

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
