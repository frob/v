package cmd

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
)

const vendorsFile = "vendors.toml"

type vendor struct {
	URL    string `toml:"url"`
	Ref    string `toml:"ref,omitempty"`
	Commit string `toml:"commit,omitempty"`
	Path   string `toml:"path"`
}

type vendorsConfig = map[string]vendor

// resolveCommit resolves a URL and optional ref to (resolvedRef, commitHash).
// When ref is empty the default branch name and its HEAD commit are returned.
// It is a variable so tests can replace the implementation.
var resolveCommit = remoteCommit

// downloadRepo is a variable so tests can replace the implementation.
var downloadRepo = cloneRepo

var addDestination string

var addCmd = &cobra.Command{
	Use:   "add <url> [ref]",
	Short: "Add a git repository to vendors.toml",
	Args:  cobra.RangeArgs(1, 2),
	RunE:  runAdd,
}

func init() {
	addCmd.Flags().StringVarP(&addDestination, "destination", "d", "", "directory to clone into (default: vendor/<host>/<path>)")
	rootCmd.AddCommand(addCmd)
}

func runAdd(cmd *cobra.Command, args []string) error {
	repoURL := args[0]
	var ref string
	if len(args) == 2 {
		ref = args[1]
	}

	resolvedRef, commit, err := resolveCommit(repoURL, ref)
	if err != nil {
		return fmt.Errorf("resolving commit: %w", err)
	}

	dest := addDestination
	if dest == "" {
		dest, err = defaultPath(repoURL)
		if err != nil {
			return err
		}
	}

	if err := downloadRepo(repoURL, commit, dest); err != nil {
		return fmt.Errorf("downloading repository: %w", err)
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

	cfg[repoURL] = vendor{URL: repoURL, Ref: resolvedRef, Commit: commit, Path: dest}

	out, err := toml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("encoding %s: %w", vendorsFile, err)
	}
	if err := os.WriteFile(vendorsFile, out, 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", vendorsFile, err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "added %s to %s\n", repoURL, vendorsFile)
	return nil
}

func defaultPath(repoURL string) (string, error) {
	u, err := url.Parse(repoURL)
	if err != nil {
		return "", fmt.Errorf("parsing url: %w", err)
	}
	path := filepath.FromSlash(strings.TrimSuffix(u.Path, ".git"))
	return filepath.Join("vendor", u.Host, path), nil
}

func cloneRepo(repoURL, commit, dest string) error {
	repo, err := git.PlainClone(dest, false, &git.CloneOptions{
		URL:      repoURL,
		Progress: os.Stderr,
	})
	if err != nil {
		if errors.Is(err, git.ErrRepositoryAlreadyExists) {
			return fmt.Errorf("directory already exists at %s; remove it before re-adding", dest)
		}
		return fmt.Errorf("cloning: %w", err)
	}

	w, err := repo.Worktree()
	if err != nil {
		return err
	}

	if err := w.Checkout(&git.CheckoutOptions{
		Hash:  plumbing.NewHash(commit),
		Force: true,
	}); err != nil {
		return err
	}

	return os.RemoveAll(filepath.Join(dest, ".git"))
}

func remoteCommit(repoURL, ref string) (resolvedRef, commit string, err error) {
	rem := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: "origin",
		URLs: []string{repoURL},
	})

	refs, err := rem.List(&git.ListOptions{})
	if err != nil {
		return "", "", fmt.Errorf("listing remote refs: %w", err)
	}

	if ref == "" {
		// HEAD is a symbolic ref pointing to the default branch (e.g. refs/heads/main).
		// Its hash is always zero — we must follow the symref to get the real commit.
		var headTarget plumbing.ReferenceName
		for _, r := range refs {
			if r.Name() == plumbing.HEAD {
				headTarget = r.Target()
				break
			}
		}
		if headTarget == "" {
			return "", "", fmt.Errorf("HEAD not found")
		}
		for _, r := range refs {
			if r.Name() == headTarget {
				branch := strings.TrimPrefix(headTarget.String(), "refs/heads/")
				return branch, r.Hash().String(), nil
			}
		}
		return "", "", fmt.Errorf("default branch %q not found", headTarget)
	}

	if isCommitHash(ref) {
		return ref, ref, nil
	}

	// For annotated tags, the remote advertises a peeled ref (e.g. refs/tags/v1.0.0^{})
	// that points directly to the commit rather than the tag object.
	peeled := plumbing.ReferenceName("refs/tags/" + ref + "^{}")
	for _, r := range refs {
		if r.Name() == peeled {
			return ref, r.Hash().String(), nil
		}
	}

	candidates := []plumbing.ReferenceName{
		plumbing.NewBranchReferenceName(ref),
		plumbing.NewTagReferenceName(ref),
	}
	for _, r := range refs {
		for _, c := range candidates {
			if r.Name() == c {
				return ref, r.Hash().String(), nil
			}
		}
	}

	return "", "", fmt.Errorf("ref %q not found in remote", ref)
}

func isCommitHash(s string) bool {
	if len(s) != 40 {
		return false
	}
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}
