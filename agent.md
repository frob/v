# agent.md

`v` manages vendored git repositories: it resolves refs to exact commit
hashes, downloads repo contents (no `.git` directory) into `vendor/`, and
records everything in `vendors.toml`. This file tells agents and new
contributors how to work in the repo. Keep it short; deeper detail lives in
the docs site.

## Tooling model

Every build, test, lint, and docs step runs **inside a container**. The only
host prerequisites are a container runtime (Docker, Podman, OrbStack) and the
[Task](https://taskfile.dev/) runner. Do not install Go, `golangci-lint`, or
GoReleaser on the host — they live in the `Dockerfile` toolchain image and the
docs image, invoked through Taskfile targets. See the `containerized-tooling`
conventions.

## Task runner

`Taskfile.yml` is the single entry point. Run `task` with no arguments to list
every task. Never invoke `go`, `docker`, `goreleaser`, or `mkdocs` directly as
the user-facing interface — they are implementation details wrapped by tasks.

Common tasks:

- `task init` — one-time bootstrap (builds the container images, primes caches).
- `task check` — quality gate (test + lint); run before opening a PR.
- `task build` — compile `dist/v`.
- `task build:release` — local cross-compile of all release artifacts.
- `task shell` — interactive shell in the tooling container.

Task naming follows `verb:subject` (`build:docs`, `serve:docs`,
`deploy:release`). See the `taskfile-conventions` rules.

## Documentation

Docs source is Markdown under `docs/content/` (Diataxis structure). Build the
site with `task build:docs` (output `./site`, gitignored) and preview it with
`task serve:docs` → <http://localhost:8080>. The API reference under
`docs/content/reference/api.md` is regenerated from `go doc` on every docs
build — do not hand-edit it.

## Architecture

Command structure, the single `os.Exit` error boundary, the multi-stage
Docker build, vendoring, and the release pipeline are documented in
`.claude/docs/architectural_patterns.md`. The short version:

- `main.go` only calls `cmd.Execute()`. All CLI logic is in the `cmd` package.
- One file per subcommand in `cmd/`, registered on `rootCmd` via `init()`.
- Subcommands return errors up the stack; only `cmd/root.go` calls `os.Exit`.

## Adding a new task

If an operation will be run more than once, add it to `Taskfile.yml` with
`verb:subject` naming and a multi-sentence `desc`. Any task that needs the
toolchain container must `deps:` on `build:container`.

## What NOT to do

- Do not install toolchains on the host.
- Do not invoke `docker run`, `go`, or `mkdocs` from docs or examples — wrap
  them in a task.
- Do not add language-native runners (`go`, `make`) as user-facing interfaces.
- Do not call `os.Exit` from a subcommand.
- Do not commit generated output (`dist/`, `site/`, `.cache/`).
