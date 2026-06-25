# Tutorials

Learning-oriented, hand-held walkthroughs. Start here if you are new to `v`.

## Vendor your first dependency

1. Pick a repository you want to copy into your project.
2. Run `v add <url>`. `v` resolves the default branch to an exact commit,
   downloads the contents into `vendor/`, and writes a `vendors.toml`
   entry.
3. Commit both `vendor/` and `vendors.toml` to your repository.

## Pin and later refresh a dependency

1. `v add <url> v1.0.0` — vendor a specific tag.
2. Some weeks later, `v update <url> v1.1.0` — switch to a newer tag and
   re-fetch.
3. Review the diff in `vendor/` and the new `commit` in `vendors.toml`.

!!! note "Placeholder"
    More step-by-step tutorials will be added here. See
    [How-to](../how-to/index.md) for shorter, task-focused recipes.
