package cmd

import (
	"fmt"
	"os"
	"testing"

	"github.com/pelletier/go-toml/v2"
)

// writeVendors writes a vendorsConfig to vendors.toml in the current directory.
func writeVendors(t *testing.T, cfg vendorsConfig) {
	t.Helper()
	data, err := toml.Marshal(cfg)
	if err != nil {
		t.Fatalf("marshal vendors: %v", err)
	}
	if err := os.WriteFile(vendorsFile, data, 0o644); err != nil {
		t.Fatalf("write %s: %v", vendorsFile, err)
	}
}

const (
	hash1 = "1111111111111111111111111111111111111111"
	hash2 = "2222222222222222222222222222222222222222"
)

func TestUpdate_SingleEntry(t *testing.T) {
	chdir(t, t.TempDir())
	mockDownload(t)

	writeVendors(t, vendorsConfig{
		"https://github.com/example/repo": {
			URL: "https://github.com/example/repo", Ref: "main", Commit: hash1,
			Path: "vendor/github.com/example/repo",
		},
	})

	orig := resolveCommit
	resolveCommit = func(_, _ string) (string, string, error) { return "main", hash2, nil }
	t.Cleanup(func() { resolveCommit = orig })

	if err := runUpdate(testCmd(), []string{"https://github.com/example/repo"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entry := readVendors(t)["https://github.com/example/repo"]
	if entry.Commit != hash2 {
		t.Errorf("commit = %q, want %q", entry.Commit, hash2)
	}
	if entry.Ref != "main" {
		t.Errorf("ref = %q, want %q", entry.Ref, "main")
	}
}

func TestUpdate_AllEntries(t *testing.T) {
	chdir(t, t.TempDir())
	mockDownload(t)

	writeVendors(t, vendorsConfig{
		"https://github.com/foo/one": {
			URL: "https://github.com/foo/one", Ref: "main", Commit: hash1,
			Path: "vendor/github.com/foo/one",
		},
		"https://github.com/foo/two": {
			URL: "https://github.com/foo/two", Ref: "main", Commit: hash1,
			Path: "vendor/github.com/foo/two",
		},
	})

	orig := resolveCommit
	resolveCommit = func(_, _ string) (string, string, error) { return "main", hash2, nil }
	t.Cleanup(func() { resolveCommit = orig })

	if err := runUpdate(testCmd(), []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, key := range []string{"https://github.com/foo/one", "https://github.com/foo/two"} {
		entry := readVendors(t)[key]
		if entry.Commit != hash2 {
			t.Errorf("%s commit = %q, want %q", key, entry.Commit, hash2)
		}
	}
}

func TestUpdate_NewRef(t *testing.T) {
	chdir(t, t.TempDir())
	mockDownload(t)

	writeVendors(t, vendorsConfig{
		"https://github.com/example/repo": {
			URL: "https://github.com/example/repo", Ref: "main", Commit: hash1,
			Path: "vendor/github.com/example/repo",
		},
	})

	var gotRef string
	orig := resolveCommit
	resolveCommit = func(_, ref string) (string, string, error) { gotRef = ref; return ref, hash2, nil }
	t.Cleanup(func() { resolveCommit = orig })

	if err := runUpdate(testCmd(), []string{"https://github.com/example/repo", "v2.0.0"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotRef != "v2.0.0" {
		t.Errorf("resolved ref = %q, want %q", gotRef, "v2.0.0")
	}

	entry := readVendors(t)["https://github.com/example/repo"]
	if entry.Ref != "v2.0.0" {
		t.Errorf("ref = %q, want %q", entry.Ref, "v2.0.0")
	}
	if entry.Commit != hash2 {
		t.Errorf("commit = %q, want %q", entry.Commit, hash2)
	}
}

func TestUpdate_AlreadyUpToDate(t *testing.T) {
	chdir(t, t.TempDir())

	writeVendors(t, vendorsConfig{
		"https://github.com/example/repo": {
			URL: "https://github.com/example/repo", Ref: "main", Commit: hash1,
			Path: "vendor/github.com/example/repo",
		},
	})

	var downloadCalled bool
	orig := downloadRepo
	downloadRepo = func(_, _, _ string) error { downloadCalled = true; return nil }
	t.Cleanup(func() { downloadRepo = orig })

	mockCommit(t, "main", hash1) // same commit — nothing to do

	if err := runUpdate(testCmd(), []string{"https://github.com/example/repo"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if downloadCalled {
		t.Error("downloadRepo should not be called when already up to date")
	}
}

func TestUpdate_URLNotFound(t *testing.T) {
	chdir(t, t.TempDir())

	writeVendors(t, vendorsConfig{
		"https://github.com/example/repo": {
			URL: "https://github.com/example/repo", Ref: "main", Commit: hash1,
			Path: "vendor/github.com/example/repo",
		},
	})

	if err := runUpdate(testCmd(), []string{"https://github.com/other/repo"}); err == nil {
		t.Fatal("expected error for unknown URL, got nil")
	}
}

func TestUpdate_NoVendorsFile(t *testing.T) {
	chdir(t, t.TempDir())

	if err := runUpdate(testCmd(), []string{}); err == nil {
		t.Fatal("expected error when vendors.toml does not exist, got nil")
	}
}

func TestUpdate_ResolveError(t *testing.T) {
	chdir(t, t.TempDir())

	writeVendors(t, vendorsConfig{
		"https://github.com/example/repo": {
			URL: "https://github.com/example/repo", Ref: "main", Commit: hash1,
			Path: "vendor/github.com/example/repo",
		},
	})

	orig := resolveCommit
	resolveCommit = func(_, _ string) (string, string, error) { return "", "", fmt.Errorf("no network") }
	t.Cleanup(func() { resolveCommit = orig })

	if err := runUpdate(testCmd(), []string{"https://github.com/example/repo"}); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestUpdate_DownloadError(t *testing.T) {
	chdir(t, t.TempDir())

	writeVendors(t, vendorsConfig{
		"https://github.com/example/repo": {
			URL: "https://github.com/example/repo", Ref: "main", Commit: hash1,
			Path: "vendor/github.com/example/repo",
		},
	})

	orig := resolveCommit
	resolveCommit = func(_, _ string) (string, string, error) { return "main", hash2, nil }
	t.Cleanup(func() { resolveCommit = orig })

	origDl := downloadRepo
	downloadRepo = func(_, _, _ string) error { return fmt.Errorf("network failure") }
	t.Cleanup(func() { downloadRepo = origDl })

	if err := runUpdate(testCmd(), []string{"https://github.com/example/repo"}); err == nil {
		t.Fatal("expected error, got nil")
	}
}
