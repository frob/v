# Architectural Patterns

## Thin Entry Point / cmd Package Delegation

`main.go:1-7` contains only a call to `cmd.Execute()`. All CLI logic lives in the `cmd` package. This keeps the binary entry point trivial and makes the `cmd` package independently testable.

## Single Error Boundary

`cmd/root.go:14-17` is the sole call site for `os.Exit`. Subcommands return errors up the call stack; the root `Execute()` function handles the exit. Do not call `os.Exit` from within subcommands.

## One File Per Subcommand in cmd/

New subcommands are added as individual files in `cmd/` and registered on `rootCmd` via `init()` or an explicit `AddCommand` call. `cmd/root.go` owns the root command definition only.

## Taskfile as the Command Interface

All shell operations are run through `task <name>` (Taskfile.yml). If a new shell command is needed, add a task for it first, then invoke it via `task`. Direct shell invocation is avoided.

## Vendored Dependencies

Dependencies are checked into `vendor/` (`go mod vendor`). Builds do not require network access. When adding or updating dependencies: run `task tidy` then `task vendor`.

## Multi-Stage Docker Build / Static Binary

`Dockerfile:1-15` uses a two-stage build: `golang:1.23-alpine` as builder, `scratch` as the final image. `CGO_ENABLED=0` produces a fully static binary. Keep CGO disabled; avoid dependencies that require cgo.
