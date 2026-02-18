package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
)

const vendorsFile = "vendors.toml"

type vendor struct {
	URL string `toml:"url"`
	Ref string `toml:"ref,omitempty"`
}

type vendorsConfig = map[string]vendor

var addCmd = &cobra.Command{
	Use:   "add <url> [ref]",
	Short: "Add a git repository to vendors.toml",
	Args:  cobra.RangeArgs(1, 2),
	RunE:  runAdd,
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func runAdd(cmd *cobra.Command, args []string) error {
	url := args[0]
	var ref string
	if len(args) == 2 {
		ref = args[1]
	}

	cfg := make(vendorsConfig)
	data, err := os.ReadFile(vendorsFile)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("reading %s: %w", vendorsFile, err)
	}
	if err == nil {
		if err := toml.Unmarshal(data, &cfg); err != nil {
			return fmt.Errorf("parsing %s: %w", vendorsFile, err)
		}
	}

	cfg[url] = vendor{URL: url, Ref: ref}

	out, err := toml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("encoding %s: %w", vendorsFile, err)
	}
	if err := os.WriteFile(vendorsFile, out, 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", vendorsFile, err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "added %s to %s\n", url, vendorsFile)
	return nil
}
