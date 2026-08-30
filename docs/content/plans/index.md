# Plans

Past plans (done) and future plans (proposed). The roadmap lives here so
contributors know what is intended before it is built.

## Proposed

### Correctness

- [ ] **`add` should clear its destination** — `cloneRepo` only errors on
  `git.ErrRepositoryAlreadyExists`, which requires a `.git` directory, but `v`
  strips `.git` from every tree it writes. Re-adding over an existing vendored
  path therefore clones *into* the populated directory and merges new files
  over stale ones, so files deleted upstream survive. `update` already does
  `os.RemoveAll` first; `add` should match. (`cmd/add.go`)
- [ ] **`update <url> <ref>` loses the ref switch when the commit is
  unchanged** — the already-up-to-date branch returns early before the new ref
  is written back, so switching to a different ref that points at the same
  commit leaves the old ref recorded in `vendors.toml`. (`cmd/update.go`)
- [ ] **Guard the abbreviated-hash slices** — `update` slices `commit[:7]`
  unguarded, so a hand-written `vendors.toml` with a short or absent `commit`
  panics instead of erroring. (`cmd/update.go`)
- [ ] **Decide what SSH URL support means** — `defaultPath` runs `url.Parse` on
  the raw argument, so an scp-style `git@github.com:org/repo.git` has no
  parseable host and yields no sensible default path, yet the README and both
  reference pages advertise "HTTPS or SSH". Either normalise scp-style URLs or
  correct the claim. (`cmd/add.go`)

### Release and packaging

- [ ] **Replace the placeholder package description** — `.goreleaser.yml`
  ships `description: v is a CLI tool` to nfpm and Homebrew, so that string
  appears in `apt show v`, `rpm -qi v`, and `brew info v`. The nfpm
  `maintainer` also lacks the `Name <email>` form deb tooling expects.
- [ ] **Settle the tag prefix** — the README's release instructions say
  `git tag v0.x.0`; the repository's actual tags are unprefixed.
- [ ] **Stop pinning the `install.sh` URL to a release branch** — the README's
  one-liner needs editing every minor cycle, and published copies end up
  pointing at stale branches.

### Enhancements

- [ ] **Git round-trip support** — make it easy to temporarily restore the
  `.git` directory of a vendored dependency so you can make changes and
  contribute them back upstream, then re-vendor the updated code.
- [ ] **Selective vendoring** — allow individual dependencies to be
  excluded from the `vendor/` tree (e.g. for license-compatibility reasons)
  while still being tracked in `vendors.toml`.
- [ ] **`v status`** — verify each vendored tree still matches its recorded
  commit. Nothing currently detects local drift, which is the gap in the
  provenance argument the Explanation page makes. `v list` and `v remove` are
  natural companions.
- [ ] **Run `test:integrated` in a temporary directory** — it writes
  `vendors.toml` and `test-v` into the repository root, which is the only
  reason `vendors.toml` is gitignored, directly contradicting the
  documentation's "always commit `vendors.toml`".
- [ ] **Shallow-clone in `cloneRepo`** — it fetches full history and then
  deletes `.git`. A depth-limited or single-branch fetch is materially faster
  for large upstreams, at the cost of commits not reachable from an
  advertised ref.
- [ ] **Add a `.golangci.yml`** — lint currently runs whatever the pinned
  version defaults to; an explicit config keeps the gate reproducible across
  upgrades.
- [ ] **Automate dependency bumps** — `GOLANGCI_LINT_VERSION`,
  `MKDOCS_MATERIAL_VERSION`, `GORELEASER_IMAGE`, and the Go base image are all
  pinned by hand.

## Done

- [x] **Changelog reconciled with the tags** — the section published as
  `0.2.0` was never tagged, so its contents moved to `[Unreleased]`; the
  missing `0.1.0` and `0.1.1` entries were added; `update` was reattributed
  from `0.2.0` to `0.1.0`, where it actually shipped; and `0.0.1` and `0.1.1`
  are now marked as the re-tags they are, pointing at the same commits as
  `0.0.0` and `0.1.0`.
- [x] **Documentation accuracy pass** — corrected
  `.claude/docs/architectural_patterns.md`, which had described a two-stage
  `scratch` build that never existed and pointed at a `task release:snapshot`
  that is not a task; repointed the dead `main`-branch links in Quickstart,
  Explanation, and Changelog; described `task check` as the test + lint +
  `test:integrated` gate it actually is; moved the generated cobra pages out
  of `reference/api` into `reference/cli/` so they no longer collide with the
  `go doc` output, and gitignored them; completed the README task table; and
  added the two documentation steps to both "Adding a command" checklists.
- [x] **`add` command** — vendor a git repository and record it in
  `vendors.toml`.
- [x] **`update` command** — re-fetch vendored repositories at their current
  or a new ref.
- [x] **Release pipeline** — GoReleaser cross-builds, Linux packages, and a
  Homebrew formula.
