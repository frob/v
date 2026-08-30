# Documentation toolchain container — MkDocs + Material theme.
#
# One container per concern: this image knows nothing about the Go
# toolchain. The reference pages are generated separately in the build
# container — api.md by `go doc` (docs/gen-api.sh) and cli/ from the cobra
# command tree (docs/gen-cli.sh) — as gitignored Markdown that this image
# then renders into the site.
#
# ENTRYPOINT is `mkdocs`, so `docker run <image> build` runs `mkdocs build`.
FROM python:3.12-alpine

ARG MKDOCS_MATERIAL_VERSION=9.5.44
RUN pip install --no-cache-dir "mkdocs-material==${MKDOCS_MATERIAL_VERSION}"

WORKDIR /work/docs
ENTRYPOINT ["mkdocs"]
