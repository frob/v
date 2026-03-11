package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update [url [ref]]",
	Short: "Update vendored repositories to the latest commit for their ref",
	Args:  cobra.MaximumNArgs(2),
	RunE:  runUpdate,
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func runUpdate(cmd *cobra.Command, args []string) error {
	data, err := os.ReadFile(vendorsFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("%s not found; nothing to update", vendorsFile)
		}
		return fmt.Errorf("reading %s: %w", vendorsFile, err)
	}

	cfg := make(vendorsConfig)
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("parsing %s: %w", vendorsFile, err)
	}

	var newRef string
	urls := make([]string, 0, len(cfg))
	if len(args) >= 1 {
		if _, ok := cfg[args[0]]; !ok {
			return fmt.Errorf("%q not found in %s", args[0], vendorsFile)
		}
		urls = append(urls, args[0])
		if len(args) == 2 {
			newRef = args[1]
		}
	} else {
		for u := range cfg {
			urls = append(urls, u)
		}
	}

	for _, u := range urls {
		v := cfg[u]
		ref := v.Ref
		if newRef != "" {
			ref = newRef
		}
		_, commit, err := resolveCommit(v.URL, ref)
		if err != nil {
			return fmt.Errorf("resolving %s: %w", v.URL, err)
		}

		if commit == v.Commit {
			fmt.Fprintf(cmd.OutOrStdout(), "%s is already up to date (%s)\n", v.URL, commit[:7])
			continue
		}

		if err := os.RemoveAll(v.Path); err != nil {
			return fmt.Errorf("removing %s: %w", v.Path, err)
		}
		if err := downloadRepo(v.URL, commit, v.Path); err != nil {
			return fmt.Errorf("downloading %s: %w", v.URL, err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "updated %s %s -> %s\n", v.URL, v.Commit[:7], commit[:7])
		v.Commit = commit
		v.Ref = ref
		cfg[u] = v
	}

	out, err := toml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("encoding %s: %w", vendorsFile, err)
	}
	if err := os.WriteFile(vendorsFile, out, 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", vendorsFile, err)
	}

	return nil
}
