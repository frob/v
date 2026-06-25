//go:build ignore

// Command gendocs regenerates the Markdown CLI reference from the cobra command
// tree. It is excluded from normal builds by the "ignore" constraint and is run
// explicitly by the docs build:
//
//	go run -tags tools gendocs.go <output-dir>
//
// The "ignore" tag keeps it out of `go build ./...`; the "-tags tools" applied
// to its dependencies pulls in cmd.GenMarkdownTree (and thus cobra/doc), which
// is likewise excluded from the release binary. See docs/gen-cli.sh.
package main

import (
	"fmt"
	"os"

	"github.com/frob/v/cmd"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run -tags tools gendocs.go <output-dir>")
		os.Exit(1)
	}
	if err := cmd.GenMarkdownTree(os.Args[1]); err != nil {
		fmt.Fprintln(os.Stderr, "gendocs:", err)
		os.Exit(1)
	}
}
