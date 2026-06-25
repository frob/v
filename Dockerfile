# Build / test / lint toolchain container.
#
# The host needs only Docker and go-task; every Go command (build, test,
# vet, lint) runs inside this image, driven by Taskfile targets. Pinned for
# reproducibility — a toolchain upgrade is a diff to this file, not a
# "brew upgrade" in someone's setup notes.
#
# This is NOT the distributable artifact. Release binaries are cross-compiled
# by GoReleaser (see `task build:release` / `task deploy:release`).
FROM golang:1.24-alpine

# git       — go-git's test fixtures and `go` module fetching need it.
# bash      — interactive `task shell`.
# curl      — installing pinned tools below.
# ca-certs  — HTTPS module/proxy fetches.
RUN apk add --no-cache git bash curl ca-certificates

# golangci-lint — pinned. Installed to /usr/local/bin so it is on PATH for
# any UID the container runs as.
ARG GOLANGCI_LINT_VERSION=v1.64.8
RUN curl -sSfL "https://raw.githubusercontent.com/golangci/golangci-lint/${GOLANGCI_LINT_VERSION}/install.sh" \
      | sh -s -- -b /usr/local/bin "${GOLANGCI_LINT_VERSION}"

WORKDIR /work

# The Taskfile runs this image as the host UID with GOPATH / GOCACHE /
# GOMODCACHE / GOLANGCI_LINT_CACHE redirected under /work/.cache (a
# gitignored, host-owned dir) so build outputs and caches are never
# root-owned and survive across container runs without a named volume.
