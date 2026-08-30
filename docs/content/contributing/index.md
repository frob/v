# Contribution guide

## Tooling model

Every build, test, and lint step runs inside a container. The only host
prerequisites are a **container runtime** (Docker, Podman, or OrbStack) and
the **[Task](https://taskfile.dev/)** runner. You do not install Go,
`golangci-lint`, or GoReleaser on your machine — they live in the
`Dockerfile` toolchain image and are invoked through Taskfile targets.

## Getting started

```sh
task init     # build the tooling + docs containers, prime the module cache
task check    # run tests, linters, and the integrated smoke test (needs network)
```

Run `task` with no arguments to list every available task.

## Common tasks

| Task | What it does |
|------|--------------|
| `task build` | Compile the binary to `dist/v`. |
| `task test` | Run the test suite. |
| `task lint` | Run `go vet` and `golangci-lint`. |
| `task check` | Quality gate: test + lint + `test:integrated`. Needs network. |
| `task test:integrated` | End-to-end smoke test against a real remote. Needs network. |
| `task build:docs` | Generate this documentation site into `site/`. |
| `task serve:docs` | Build the docs and serve them at <http://localhost:8080>. |
| `task shell` | Open an interactive shell inside the tooling container. |
| `task build:release` | Cross-compile all release artifacts locally (no publish). |

## Adding a command

1. Create a new file in `cmd/` — one file per subcommand. See `cmd/add.go`
   for a reference implementation.
2. Register it on `rootCmd` with `rootCmd.AddCommand(...)` inside the
   file's `init()`.
3. Return errors up the call stack — never call `os.Exit` from a
   subcommand. `cmd/root.go` owns the single error boundary.
4. Document the command in [CLI commands](../reference/commands.md).
5. Add its generated cobra page to the `nav` in `docs/mkdocs.yml` under
   **Reference → CLI (cobra)**. `task build:docs` writes the page into
   `docs/content/reference/cli/`, but the nav entry is added by hand.

## Adding a dependency

```sh
task get -- github.com/example/pkg@v1.0.0
task tidy
task vendor   # vendor/ is committed, so builds need no network
```

## Adding a task

If an operation will be run more than once, add it to `Taskfile.yml` with
`verb:subject` naming and a multi-sentence `desc`. Any task that needs the
toolchain container must `deps:` on `build:container`. Never expose `go`,
`docker`, or `goreleaser` as the user-facing interface — wrap them in a
task.

## Releasing

See the [Plans](../plans/index.md) and the project README. Releases are cut
with GoReleaser via `task deploy:release` (requires `GITHUB_TOKEN` and a
pushed tag).

## Publishing the docs

This site publishes to [GitHub Pages](https://frob.github.io/v/) automatically.
The `Docs` workflow (`.github/workflows/docs.yml`) runs on every push to the
default branch, rebuilds the site with the same `task build:docs` used locally,
and deploys the `site/` output. There is nothing to run by hand — merge to the
default branch and the live site updates. Use `task serve:docs` to preview
changes before they merge.
