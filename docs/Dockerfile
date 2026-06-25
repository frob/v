# Documentation toolchain container — MkDocs + Material theme.
#
# One container per concern: this image knows nothing about the Go
# toolchain. The API reference is generated separately by `go doc` in the
# build container (see docs/gen-api.sh) and committed as Markdown that this
# image renders into the site.
#
# ENTRYPOINT is `mkdocs`, so `docker run <image> build` runs `mkdocs build`.
FROM python:3.12-alpine

ARG MKDOCS_MATERIAL_VERSION=9.5.44
RUN pip install --no-cache-dir "mkdocs-material==${MKDOCS_MATERIAL_VERSION}"

WORKDIR /work/docs
ENTRYPOINT ["mkdocs"]
