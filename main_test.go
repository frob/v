package main

import (
	"os"
	"testing"
)

// TestEntrypoint is a smoke test that main() wires into cmd.Execute() correctly.
// Run with no subcommand, cobra prints usage and returns nil (no error), so
// Execute() does not call os.Exit and main() returns normally. os.Args and the
// output streams are saved and restored so the test leaves no global state
// behind, and stdout/stderr are redirected to discard the usage text.
func TestEntrypoint(t *testing.T) {
	origArgs := os.Args
	origStdout, origStderr := os.Stdout, os.Stderr

	devnull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		t.Fatalf("open %s: %v", os.DevNull, err)
	}

	os.Args = []string{"v"}
	os.Stdout, os.Stderr = devnull, devnull
	t.Cleanup(func() {
		os.Args = origArgs
		os.Stdout, os.Stderr = origStdout, origStderr
		devnull.Close()
	})

	main()
}
