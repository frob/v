# CLAUDE.md

## Project Overview

`v` is a Go CLI tool (`github.com/frob/v`). It is in early development — currently bootstrapped with a root command and help output only.

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.23 |
| CLI framework | Cobra v1.8.1 (`vendor/github.com/spf13/cobra`) |
| Task runner | Taskfile v3 (`Taskfile.yml`) |
| Release pipeline | GoReleaser v2 (`.goreleaser.yml`) |
| Container | Docker (multi-stage, scratch final image) |

## Key Directories

| Path | Purpose |
|---|---|
| `main.go` | Binary entry point — delegates immediately to `cmd` |
| `cmd/` | All CLI commands; one file per subcommand |
| `dist/` | Compiled binary and release artifact output (git-ignored) |
| `vendor/` | Vendored dependencies (checked in) |
| `.claude/docs/` | Extended documentation for Claude |

## Build & Test Commands

All commands run through the Taskfile. If no task exists for an operation, add one first.

| Task | Action |
|---|---|
| `task build` | Compile binary to `dist/v` (current platform) |
| `task run` | `go run . [args]` — pass args with `-- <args>` |
| `task test` | `go test ./...` |
| `task clean` | Remove `dist/` |
| `task tidy` | `go mod tidy` |
| `task vendor` | `go mod vendor` |
| `task docker:build` | Build Docker image tagged `v` |

## Release & Distribution Commands

| Task | Action |
|---|---|
| `task release:check` | Validate `.goreleaser.yml` config |
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

Check these files when relevant:

- `.claude/docs/architectural_patterns.md` — command structure, error handling, Docker build, vendoring, Taskfile conventions, and release pipeline
