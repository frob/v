## v update

Update vendored repositories to the latest commit for their ref

### Synopsis

Update re-fetches vendored repositories and updates vendors.toml with the
latest commit hash for each entry's ref.

With no arguments all vendors are updated. Provide a URL to update a single
entry. Provide a URL and a ref to switch that entry to a different branch,
tag, or commit hash.

```
v update [url [ref]] [flags]
```

### Options

```
  -h, --help   help for update
```

### SEE ALSO

* [v](v.md)	 - Manage vendored git repositories

