package cmd

import (
	"bytes"
	"os"
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

func TestAdd_URLOnly(t *testing.T) {
	chdir(t, t.TempDir())

	if err := runAdd(testCmd(), []string{"https://github.com/example/repo"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cfg := readVendors(t)
	entry, ok := cfg["https://github.com/example/repo"]
	if !ok {
		t.Fatal("entry not found in vendors.toml")
	}
	if entry.URL != "https://github.com/example/repo" {
		t.Errorf("url = %q, want %q", entry.URL, "https://github.com/example/repo")
	}
	if entry.Ref != "" {
		t.Errorf("ref = %q, want empty", entry.Ref)
	}
}

func TestAdd_WithRef(t *testing.T) {
	chdir(t, t.TempDir())

	if err := runAdd(testCmd(), []string{"https://github.com/example/repo", "v1.0.0"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entry := readVendors(t)["https://github.com/example/repo"]
	if entry.Ref != "v1.0.0" {
		t.Errorf("ref = %q, want %q", entry.Ref, "v1.0.0")
	}
}

func TestAdd_MultipleEntries(t *testing.T) {
	chdir(t, t.TempDir())

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
	if cfg["https://github.com/foo/one"].Ref != "v1.0.0" {
		t.Errorf("foo/one ref = %q, want %q", cfg["https://github.com/foo/one"].Ref, "v1.0.0")
	}
	if cfg["https://github.com/foo/two"].Ref != "abc123" {
		t.Errorf("foo/two ref = %q, want %q", cfg["https://github.com/foo/two"].Ref, "abc123")
	}
}

func TestAdd_UpdateExisting(t *testing.T) {
	chdir(t, t.TempDir())

	if err := runAdd(testCmd(), []string{"https://github.com/example/repo", "v1.0.0"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := runAdd(testCmd(), []string{"https://github.com/example/repo", "v2.0.0"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cfg := readVendors(t)
	if len(cfg) != 1 {
		t.Fatalf("len = %d, want 1 (duplicate entry created)", len(cfg))
	}
	if cfg["https://github.com/example/repo"].Ref != "v2.0.0" {
		t.Errorf("ref = %q, want %q", cfg["https://github.com/example/repo"].Ref, "v2.0.0")
	}
}
