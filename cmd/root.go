package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "v",
	Short: "Manage vendored git repositories",
	Long: `v manages vendored git repositories.

It resolves refs (branches, tags, or commit hashes) to exact commit hashes,
downloads repository contents into vendor/ (without the .git directory),
and records everything in vendors.toml.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
