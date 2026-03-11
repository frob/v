# CLAUDE.md

`v` manages vendored git repositories: it resolves refs to exact commit hashes, downloads repo contents (no `.git` directory) into `vendor/`, and records everything in `vendors.toml`.

## vendors.toml Format

Entries are keyed by remote URL. `ref` defaults to the remote's default branch when omitted. `commit` is the resolved hash.

```toml
['https://github.com/example/repo']
url    = 'https://github.com/example/repo'
ref    = 'main'
commit = 'abc123...'
path   = 'vendor/github.com/example/repo'
```

## Workflow

All commands run through the Taskfile (`task <name>`). If no task exists for an operation, add one first, then run it via `task`.

## Distribution

Homebrew formula lives in `frob/homebrew-v` (pushed automatically on release).

## Additional Documentation

- `.claude/docs/architectural_patterns.md` — command structure, error handling, Docker build, vendoring, Taskfile conventions, and release pipeline
