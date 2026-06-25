## v add

Add a git repository to vendors.toml

### Synopsis

Add resolves a git repository ref to an exact commit hash, downloads the
repository contents into vendor/ (without the .git directory), and records
the entry in vendors.toml.

The ref can be a branch name, tag, or full commit hash. If omitted, the
remote's default branch is used.

By default the repository is placed at vendor/<host>/<path>. Use --destination
to override.

```
v add <url> [ref] [flags]
```

### Options

```
  -d, --destination string   directory to clone into (default: vendor/<host>/<path>)
  -h, --help                 help for add
```

### SEE ALSO

* [v](v.md)	 - Manage vendored git repositories

