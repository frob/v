# Plans

Past plans (done) and future plans (proposed). The roadmap lives here so
contributors know what is intended before it is built.

## Proposed

- [ ] **Git round-trip support** — make it easy to temporarily restore the
  `.git` directory of a vendored dependency so you can make changes and
  contribute them back upstream, then re-vendor the updated code.
- [ ] **Selective vendoring** — allow individual dependencies to be
  excluded from the `vendor/` tree (e.g. for license-compatibility reasons)
  while still being tracked in `vendors.toml`.

## Done

- [x] **`add` command** — vendor a git repository and record it in
  `vendors.toml`.
- [x] **`update` command** — re-fetch vendored repositories at their current
  or a new ref.
- [x] **Release pipeline** — GoReleaser cross-builds, Linux packages, and a
  Homebrew formula.
