# CLAUDE.md

## Project Overview

`v` is a Go CLI tool (`github.com/frob/v`) for managing vendored git repositories. It reads and writes a `vendors.toml` file in the current directory.

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
| `cmd/add.go` | `add` command — writes entries to `vendors.toml` |
| `cmd/add_test.go` | Tests for the `add` command |
| `dist/` | Compiled binary and release artifact output |
| `vendor/` | Vendored dependencies (checked in) |
| `.goreleaser.yml` | Cross-compilation and distribution config |
| `install.sh` | Curl-pipe installer for direct installation |
| `.claude/docs/` | Extended documentation for Claude |

## vendors.toml Format

Entries are keyed by remote URL. `ref` is optional and may be a commit hash, tag, or branch.

```toml
['https://github.com/example/repo']
url = 'https://github.com/example/repo'
ref = 'v1.0.0'
```

## Commands

| Command | Description |
|---|---|
| `v add <url> [ref]` | Add or update a repository in `vendors.toml` |

## Build & Test Tasks

All commands run through the Taskfile. If no task exists for an operation, add one first.

| Task | Action |
|---|---|
| `task build` | Compile binary to `dist/v` (current platform) |
| `task run` | `go run . [args]` — pass args with `-- <args>` |
| `task test` | `go test ./...` |
| `task clean` | Remove `dist/` |
| `task get -- <mod@version>` | Add a dependency |
| `task tidy` | `go mod tidy` (strips unused deps) |
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
