# Changelog

## [Unreleased]

## [0.2.0] - 2026-06-24

### Added
- `update` command to re-fetch vendored repositories at their current ref or a new one
- Optional second argument to `update` to switch a vendor to a different branch, tag, or commit hash
- Long-form help descriptions for all commands (`v --help`, `v add --help`, `v update --help`)
- `--version` flag reporting the build version
- `README.md` with installation, usage, and contributing instructions
- `CHANGELOG.md`
- `LICENSE` (MIT)
- GitHub Actions CI running tests and linting on pushes and pull requests

### Changed
- `task clean` no longer deletes `dist/config.yaml` or `dist/homebrew/`
- `dist/` is now fully excluded from version control; GoReleaser output is not tracked

## [0.0.1] - 2026-03-11

### Added
- GoReleaser pipeline producing binaries for Linux and macOS (amd64, arm64)
- `.deb`, `.rpm`, and `.pkg.tar.zst` packages via nfpm
- Homebrew formula published to `frob/homebrew-v` on release
- `install.sh` curl-pipe installer with auto-detection of OS and package manager

## [0.0.0] - 2026-03-11

### Added
- `add` command to vendor a git repository into `vendor/` and record it in `vendors.toml`
- `--destination` / `-d` flag to specify a custom vendor path
- Ref resolution: branch names, tags, annotated tags, and full commit hashes
- `vendors.toml` format with `url`, `ref`, `commit`, and `path` fields
