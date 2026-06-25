# Explanation

Understanding-oriented background. Why `v` works the way it does.

## The problem with hand-vendoring

Copying a library's source into your repository is a legitimate strategy:
no package-manager coupling, no transitive surprises, full control. But done
by hand it loses the two facts that matter most when the upstream changes:
*which commit you copied* and *how to reproduce the copy*. Months later, a
`vendor/` directory is an undated snapshot of unknown provenance.

## What `v` records

`v` makes the provenance explicit. Every vendored repository has a
`vendors.toml` entry capturing the URL, the requested ref, the **resolved
commit hash**, and the destination path. The resolved commit is the key
detail: a branch name drifts, but the recorded hash pins exactly what is in
`vendor/`.

## No `.git` directory

`v` downloads repository *contents*, not the git history. The vendored tree
is plain files — it does not become a nested repository or a submodule, and
it does not interfere with your own version control. `vendors.toml` carries
the metadata that a `.git` directory otherwise would.

## Architecture

The internal command structure, error-handling boundary, Docker build, and
release pipeline are documented in
[`.claude/docs/architectural_patterns.md`](https://github.com/frob/v/blob/main/.claude/docs/architectural_patterns.md)
and summarized for contributors in the
[Contribution guide](../contributing/index.md).
