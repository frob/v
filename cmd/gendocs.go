//go:build tools

package cmd

import "github.com/spf13/cobra/doc"

// GenMarkdownTree writes Markdown documentation for the entire command tree
// into dir, one file per command (v.md, v_add.md, v_update.md, ...).
//
// It is compiled only under the "tools" build tag so that cobra/doc and its
// Markdown dependencies never enter the release binary. The docs build invokes
// it through the root-level generator (see gendocs.go) as part of
// `task build:docs`.
func GenMarkdownTree(dir string) error {
	// Suppress the "Auto generated ... on <date>" footer so regenerated files
	// are stable and don't churn in version control.
	rootCmd.DisableAutoGenTag = true
	return doc.GenMarkdownTree(rootCmd, dir)
}
