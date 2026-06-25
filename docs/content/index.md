# v

`v` manages vendored git repositories. It resolves refs (branches, tags, or
commit hashes) to exact commit hashes, downloads repository contents
(without the `.git` directory) into `vendor/`, and records everything in
`vendors.toml`.

## Why

Sometimes a library needs to be included directly in a project — without a
package-manager integration — which means copying the source and effectively
forking it. That is error-prone: it is easy to lose track of *which* commit
you copied and *how* to refresh it. `v` documents the exact point at which a
codebase was imported and gives you a one-command path to keep it current.

## Where to go next

- **[Quickstart](quickstart.md)** — vendor your first repository in under
  five minutes.
- **[Tutorials](tutorials/index.md)** — learning-oriented walkthroughs.
- **[How-to](how-to/index.md)** — task-oriented recipes.
- **[Reference](reference/index.md)** — the `vendors.toml` format, CLI
  commands, and generated API docs.
- **[Explanation](explanation/index.md)** — the design and rationale behind
  `v`.
- **[Contribution](contributing/index.md)** — how to set up the project and
  submit changes.
