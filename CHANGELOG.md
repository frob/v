# Changelog

Versions correspond to git tags. `0.0.1` and `0.1.1` point at the same commits
as `0.0.0` and `0.1.0` respectively — they are re-tags and carry no changes of
their own.

## [Unreleased]

### Added
- `--version` flag reporting the build version, injected at release time via
  `-ldflags`
- Long-form help descriptions for all commands (`v --help`, `v add --help`,
  `v update --help`)
- `LICENSE` (MIT)
- `CONTRIBUTORS.md` setting out the human-review and human-communication
  requirements for contributions
- GitHub Actions CI running the `task check` quality gate on every push and
  pull request
- Diataxis documentation site (MkDocs + Material) under `docs/content/`, built
  with `task build:docs` and previewed with `task serve:docs`
- GitHub Actions workflow publishing the documentation site to GitHub Pages on
  every push to the default branch
- Generated reference pages, rebuilt on every docs build: the `go doc` API
  output and the per-command cobra pages
- `task test:integrated`, an end-to-end smoke test against a real remote, now
  part of `task check`
- `task test:force`, `task test:coverage`, and `task release:check`
- Tests for the `update` command
- Linux installation instructions (`.deb`, `.rpm`, Arch) and a roadmap in the
  README

### Changed
- `install.sh` resolves the release tag by following the `/releases/latest`
  redirect, so package artifacts that embed the version in their file names
  download correctly
- `install.sh` is fetched from the `0.2.x` branch rather than `main`; this
  project does not use a `main` branch
- Pinned the GoReleaser container image version used by `task build:release`
  and `task deploy:release`
- `task clean` no longer deletes `dist/config.yaml` or `dist/homebrew/`
- `dist/` is now fully excluded from version control; GoReleaser output is not
  tracked

### Fixed
- `.claude/docs/architectural_patterns.md` described a two-stage
  `golang:1.23-alpine` → `scratch` Docker build that never existed, and
  referenced a `task release:snapshot` that is not a task
- Documentation links pointing at a nonexistent `main` branch
- `task check` was documented as "test + lint" in three places despite also
  running the network-dependent `test:integrated`
- Generated cobra CLI pages moved from `docs/content/reference/api/` to
  `docs/content/reference/cli/`, where they no longer collide with the
  generated `go doc` output in `reference/api.md`

## [0.1.1] - 2026-03-11

Re-tag of 0.1.0. No changes.

## [0.1.0] - 2026-03-11

### Added
- `update` command to re-fetch vendored repositories at their current ref or a
  new one
- Optional second argument to `update` to switch a vendor to a different
  branch, tag, or commit hash
- `CHANGELOG.md`

## [0.0.1] - 2026-03-10

Re-tag of 0.0.0. No changes.

## [0.0.0] - 2026-03-10

### Added
- `add` command to vendor a git repository into `vendor/` and record it in
  `vendors.toml`
- `--destination` / `-d` flag to specify a custom vendor path
- Ref resolution: branch names, tags, annotated tags, and full commit hashes
- `vendors.toml` format with `url`, `ref`, `commit`, and `path` fields
- GoReleaser pipeline producing binaries for Linux and macOS (amd64, arm64)
- `.deb`, `.rpm`, and `.pkg.tar.zst` packages via nfpm
- Homebrew formula published to `frob/homebrew-v` on release
- `install.sh` curl-pipe installer with auto-detection of OS and package
  manager
- `README.md`
