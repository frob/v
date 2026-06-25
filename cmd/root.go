// Package cmd implements the v command-line interface.
//
// Each subcommand lives in its own file (add.go, update.go) and registers
// itself on rootCmd via an init function. Subcommands return errors up the
// stack; only Execute, the single error boundary, calls os.Exit. The package
// exports only Execute — everything else (the cobra command values, their RunE
// handlers, and the vendoring helpers) is unexported because nothing here is
// meant to be imported by other packages.
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// version is the build version, injected via -ldflags at release time by
// GoReleaser. It defaults to "dev" for local and source builds.
var version = "dev"

var rootCmd = &cobra.Command{
	Use:   "v",
	Short: "Manage vendored git repositories",
	Long: `v manages vendored git repositories.

It resolves refs (branches, tags, or commit hashes) to exact commit hashes,
downloads repository contents into vendor/ (without the .git directory),
and records everything in vendors.toml.`,
	Version: version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
