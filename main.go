// Command v manages vendored git repositories.
//
// It resolves refs (branches, tags, or commit hashes) to exact commit hashes,
// downloads repository contents into vendor/ (without the .git directory), and
// records everything in vendors.toml. All CLI logic lives in the cmd package;
// this file only wires up the entry point.
package main

import "github.com/frob/v/cmd"

func main() {
	cmd.Execute()
}
