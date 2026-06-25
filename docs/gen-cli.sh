#!/usr/bin/env sh
# Generate the per-command CLI reference from the cobra command tree, embedded
# under the docs site's Reference section. Runs inside the build (Go) container
# as part of `task build:docs` — NOT on the host, NOT in the docs container.
# Wipes and rewrites docs/content/reference/cli/.
set -eu

OUT=docs/content/reference/api

rm -rf "${OUT}"
mkdir -p "${OUT}"

# gendocs.go is excluded from normal builds; -tags tools pulls in the cobra/doc
# generator without linking it into the release binary.
go run -tags tools gendocs.go "${OUT}"

echo "Wrote CLI docs to ${OUT}"
