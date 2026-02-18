# CLAUDE.md

## Project Overview

`v` is a Go CLI tool (`github.com/frob/v`) for managing vendored git repositories. It resolves refs to exact commit hashes, downloads repository contents (no `.git` directory) into a local `vendor/` tree, and records everything in `vendors.toml`.

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.24 |
| CLI framework | Cobra v1.8.1 (`vendor/github.com/spf13/cobra`) |
| TOML | go-toml v2 (`vendor/github.com/pelletier/go-toml`) |
| Git | go-git v5 (`vendor/github.com/go-git/go-git`) |
| Task runner | Taskfile v3 (`Taskfile.yml`) |
| Release pipeline | GoReleaser v2 (`.goreleaser.yml`) |
| Container | Docker (multi-stage, scratch final image) |

## Key Directories & Files

| Path | Purpose |
|---|---|
| `main.go` | Binary entry point — delegates immediately to `cmd` |
| `cmd/root.go` | Root Cobra command and `Execute()` |
| `cmd/add.go` | `add` command — resolves ref, downloads repo, writes `vendors.toml` |
| `cmd/add_test.go` | Tests for the `add` command |
| `dist/` | Compiled binary and release artifact output |
| `vendor/` | Vendored Go dependencies (checked in) |
| `.goreleaser.yml` | Cross-compilation and distribution config |
| `install.sh` | Curl-pipe installer for direct installation |
| `.claude/docs/` | Extended documentation for Claude |

## vendors.toml Format

Entries are keyed by remote URL. `ref` is always populated — defaults to the remote's default branch when not specified. `commit` is the exact resolved hash. `path` is where the repo was downloaded.

```toml
['https://github.com/example/repo']
url    = 'https://github.com/example/repo'
ref    = 'main'
commit = 'abc123...'
path   = 'vendor/github.com/example/repo'
```

## Commands

| Command | Description |
|---|---|
| `v add <url> [ref]` | Resolve ref, download repo (no `.git`), write `vendors.toml` |
| `v add <url> [ref] -d <dir>` | Download to a custom directory instead of `vendor/<host>/<path>` |

`ref` may be a branch, tag, or commit hash. When omitted, the remote's default branch is resolved automatically.

## Build & Test Tasks

All commands run through the Taskfile. If no task exists for an operation, add one first.

| Task | Action |
|---|---|
| `task build` | Compile binary to `dist/v` (current platform) |
| `task run` | `go run . [args]` — pass args with `-- <args>` |
| `task test` | `go test ./...` |
| `task clean` | Remove `dist/` |
| `task get -- <mod@version>` | Add a dependency |
| `task tidy` | `go mod tidy` (strips unused deps — run after `get`, before `vendor`) |
| `task vendor` | `go mod vendor` |
| `task docker:build` | Build Docker image tagged `v` |

## Release & Distribution Tasks

| Task | Action |
|---|---|
| `task release:check` | Validate `.goreleaser.yml` (requires reachable git remote) |
| `task release:snapshot` | Build all platforms locally into `dist/` — no publish |
| `task release` | Cut a release and publish (requires `GITHUB_TOKEN`) |

## Distribution

| Channel | Artifact | Notes |
|---|---|---|
| Homebrew | Formula in `frob/homebrew-v` | Pushed automatically on release |
| Debian/Ubuntu | `.deb` | Built by goreleaser/nfpm |
| Fedora/RHEL | `.rpm` | Built by goreleaser/nfpm |
| Arch Linux | `.pkg.tar.zst` | Built by goreleaser/nfpm |
| Any | `install.sh` | Detects OS/arch/package manager |

## Additional Documentation

- `.claude/docs/architectural_patterns.md` — command structure, error handling, Docker build, vendoring, Taskfile conventions, and release pipeline
