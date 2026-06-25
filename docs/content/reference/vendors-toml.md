# `vendors.toml`

`v` keeps a `vendors.toml` file in your project root. It is the source of
truth for what has been vendored and at which commit. Each entry is keyed by
the remote URL.

```toml
['https://github.com/example/repo']
url    = 'https://github.com/example/repo'
ref    = 'main'
commit = 'abc123...'
path   = 'vendor/github.com/example/repo'
```

## Fields

| Field    | Meaning                                                                 |
|----------|-------------------------------------------------------------------------|
| `url`    | The git remote URL (HTTPS or SSH). Also the table key.                  |
| `ref`    | Branch, tag, or commit requested. Defaults to the remote default branch when omitted. |
| `commit` | The exact commit hash `ref` resolved to. Written by `v`.                |
| `path`   | Destination directory for the downloaded contents (no `.git`).          |

## Behavior

- Re-running `v add` for an existing URL updates the entry in place.
- `v update` re-resolves `ref` to its latest `commit` and re-downloads.
- The destination `path` is taken from the existing entry on update; use
  `v add -d` to change it.
