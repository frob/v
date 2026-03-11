# v

`v` manages vendored git repositories. It resolves refs to exact commit hashes, downloads repository contents (without the `.git` directory) into `vendor/`, and records everything in `vendors.toml`.

## Installation

**Homebrew:**
```sh
brew install frob/v/v
```

**Shell script (auto-detects OS and package manager):**
```sh
curl -sSf https://raw.githubusercontent.com/frob/v/main/install.sh | sh
```

**Manual:** Download a binary from the [releases page](https://github.com/frob/v/releases).

## Usage

### Add a repository

```sh
v add <url> [ref]
```

Downloads the repository at the given URL into `vendor/` and records it in `vendors.toml`.

- `url` — any git remote URL (HTTPS or SSH)
- `ref` — branch, tag, or commit hash (optional; defaults to the remote's default branch)

**Examples:**

```sh
# Vendor at the default branch
v add https://github.com/example/repo

# Vendor at a specific tag
v add https://github.com/example/repo v1.2.3

# Vendor at a specific commit
v add https://github.com/example/repo abc123...

# Vendor to a custom path
v add https://github.com/example/repo main -d third_party/repo
```

By default the repository is placed at `vendor/<host>/<path>` (e.g. `vendor/github.com/example/repo`). Use `-d` / `--destination` to override.

### vendors.toml

`v` keeps a `vendors.toml` in your project root. Each entry is keyed by URL:

```toml
['https://github.com/example/repo']
url    = 'https://github.com/example/repo'
ref    = 'main'
commit = 'abc123...'
path   = 'vendor/github.com/example/repo'
```

Re-running `v add` for an existing URL updates the entry in place.

### Update a repository

```sh
v update [url [ref]]
```

Re-fetches a vendored repository and updates `vendors.toml` with the new commit hash.

- With no arguments, updates all vendors to the latest commit on their current ref.
- With a URL, updates only that entry.
- With a URL and ref, switches to the new ref and fetches its latest commit.

**Examples:**

```sh
# Update all vendors
v update

# Update one vendor to the latest commit on its current ref
v update https://github.com/example/repo

# Switch a vendor to a different branch or tag
v update https://github.com/example/repo v2.0.0
```

## Contributing

### Prerequisites

- [Go](https://go.dev/) 1.24+
- [Task](https://taskfile.dev/) (`brew install go-task`)
- [GoReleaser](https://goreleaser.com/) (for releases only)

### Development workflow

All commands run through the Taskfile:

| Task | Description |
|------|-------------|
| `task build` | Build the binary to `dist/v` |
| `task run -- <args>` | Run the CLI without building first |
| `task test` | Run the test suite |
| `task clean` | Remove build artifacts |
| `task tidy` | Tidy Go modules |
| `task vendor` | Re-vendor dependencies |

### Adding a command

1. Create a new file in `cmd/` (one file per subcommand).
2. Register it on `rootCmd` via `init()` or an explicit `AddCommand` call.
3. Return errors up the call stack — do not call `os.Exit` from subcommands.

### Adding a dependency

```sh
task get -- github.com/example/pkg@v1.0.0
task tidy
task vendor
```

### Releasing

Tag a commit and run:

```sh
git tag v0.x.0
git push origin v0.x.0
GITHUB_TOKEN=<token> task release
```

To do a local dry-run first: `task release:snapshot`
