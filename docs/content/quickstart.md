# Quickstart

Get a repository vendored into your project in a few minutes.

## Install

=== "Homebrew"

    ```sh
    brew install frob/v/v
    ```

=== "Shell script"

    ```sh
    curl -sSf https://raw.githubusercontent.com/frob/v/0.2.x/install.sh | sh
    ```

See [Installation in the README](https://github.com/frob/v#installation) for
Debian/RPM/Arch packages and manual binaries.

## Vendor a repository

```sh
v add https://github.com/example/repo
```

This downloads the repository's contents into
`vendor/github.com/example/repo` (without the `.git` directory) and records
the resolved commit in `vendors.toml`.

## Pin to a tag or commit

```sh
v add https://github.com/example/repo v1.2.3
```

## Update later

```sh
# Update everything to the latest commit on each recorded ref
v update

# Update a single vendored repo
v update https://github.com/example/repo
```

Next: the [Tutorials](tutorials/index.md) walk through these in more depth.
