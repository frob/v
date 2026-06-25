# How-to / Examples

Task-oriented recipes. Each assumes you already know what `v` is — see the
[Quickstart](../quickstart.md) otherwise.

## Vendor to a custom path

```sh
v add https://github.com/example/repo main -d third_party/repo
```

By default a repository is placed at `vendor/<host>/<path>`. Use `-d` /
`--destination` to override.

## Vendor at an exact commit

```sh
v add https://github.com/example/repo abc123def456...
```

## Switch a vendored repo to a different branch

```sh
v update https://github.com/example/repo other-branch
```

## Update every vendored repository at once

```sh
v update
```

Updates all entries in `vendors.toml` to the latest commit on their current
ref.

!!! note "Placeholder"
    Add new recipes here as common workflows emerge.
