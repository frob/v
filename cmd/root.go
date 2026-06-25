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
