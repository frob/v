# CLAUDE.md

## Project Overview

`v` is a Go CLI tool (`github.com/frob/v`). It is in early development — currently bootstrapped with a root command and help output only.

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.23 |
| CLI framework | Cobra v1.8.1 (`vendor/github.com/spf13/cobra`) |
| Task runner | Taskfile v3 (`Taskfile.yml`) |
| Container | Docker (multi-stage, scratch final image) |

## Key Directories

| Path | Purpose |
|---|---|
| `main.go` | Binary entry point — delegates immediately to `cmd` |
| `cmd/` | All CLI commands; one file per subcommand |
| `dist/` | Compiled binary output (git-ignored) |
| `vendor/` | Vendored dependencies (checked in) |
| `.claude/docs/` | Extended documentation for Claude |

## Build & Test Commands

All commands run through the Taskfile. If no task exists for an operation, add one first.

| Task | Action |
|---|---|
| `task build` | Compile binary to `dist/v` |
| `task run` | `go run . [args]` — pass args with `-- <args>` |
| `task test` | `go test ./...` |
| `task clean` | Remove `dist/` |
| `task tidy` | `go mod tidy` |
| `task vendor` | `go mod vendor` |
| `task docker:build` | Build Docker image tagged `v` |

## Additional Documentation

Check these files when relevant:

- `.claude/docs/architectural_patterns.md` — command structure, error handling, Docker build, vendoring, and Taskfile conventions
