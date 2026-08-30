# Architectural Patterns

## Thin Entry Point / cmd Package Delegation

`main.go:11-13` contains only a call to `cmd.Execute()`. All CLI logic lives in the `cmd` package. This keeps the binary entry point trivial and makes the `cmd` package independently testable.

## Single Error Boundary

`cmd/root.go:32-36` is the sole call site for `os.Exit`. Subcommands return errors up the call stack; the root `Execute()` function handles the exit. Do not call `os.Exit` from within subcommands.

## One File Per Subcommand in cmd/

New subcommands are added as individual files in `cmd/` and registered on `rootCmd` via `init()` or an explicit `AddCommand` call. `cmd/root.go` owns the root command definition only.

## Taskfile as the Command Interface

All shell operations are run through `task <name>` (Taskfile.yml). If a new shell command is needed, add a task for it first, then invoke it via `task`. Direct shell invocation is avoided.

## Vendored Dependencies

Dependencies are checked into `vendor/` (`go mod vendor`). Builds do not require network access. When adding or updating dependencies: run `task tidy` then `task vendor`.

## Containerized Toolchain, Not a Distributable Image

`Dockerfile` is a single-stage `golang:1.24-alpine` image carrying the build, test, and lint toolchain (git, bash, curl, and a pinned `golangci-lint`). It is **not** the distributable artifact — nothing ships from it. The Taskfile runs it as the host UID with `GOPATH`, `GOCACHE`, `GOMODCACHE`, and `GOLANGCI_LINT_CACHE` redirected under `/work/.cache` so build outputs are never root-owned and caches survive between runs.

`docs/Dockerfile` is a separate `python:3.12-alpine` MkDocs image — one container per concern, so neither image carries the other's toolchain.

Release binaries are cross-compiled by GoReleaser with `CGO_ENABLED=0` (see below). Keep CGO disabled; avoid dependencies that require cgo.

## GoReleaser as Release Pipeline

`.goreleaser.yml` is the single source of truth for all distribution. It cross-compiles for `linux/amd64`, `linux/arm64`, `darwin/amd64`, `darwin/arm64` — all with `CGO_ENABLED=0`. nfpm produces `.deb`, `.rpm`, and `.pkg.tar.zst` (Arch) packages from the same build. The Homebrew formula is pushed to `frob/homebrew-v` automatically. Use `task build:release` for a local dry-run before publishing, and `task release:check` to validate the configuration without building. The `install.sh` script (`install.sh:1`) is the curl-pipe entry point that auto-detects the package manager and falls back to binary extraction.
