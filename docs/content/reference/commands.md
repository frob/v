# CLI commands

## `v add`

```sh
v add <url> [ref]
```

Downloads the repository at `<url>` into `vendor/` and records it in
`vendors.toml`.

| Argument / flag        | Description                                                        |
|------------------------|--------------------------------------------------------------------|
| `url`                  | Any git remote URL (HTTPS or SSH). Required.                       |
| `ref`                  | Branch, tag, or commit. Optional; defaults to the remote default.  |
| `-d`, `--destination`  | Override the destination path (default `vendor/<host>/<path>`).    |

## `v update`

```sh
v update [url [ref]]
```

Re-fetches a vendored repository and updates `vendors.toml` with the new
commit hash.

| Form                       | Effect                                                       |
|----------------------------|--------------------------------------------------------------|
| `v update`                 | Update every vendored repo to the latest commit on its ref.  |
| `v update <url>`           | Update only that entry on its current ref.                   |
| `v update <url> <ref>`     | Switch the entry to `<ref>` and fetch its latest commit.     |

## `v --version`

Prints the build version. Injected at release time via `-ldflags`; reports
`dev` for local source builds.
