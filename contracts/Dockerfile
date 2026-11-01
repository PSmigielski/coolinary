# syntax=docker/dockerfile:1.7
# Code generators used by `make proto`, all pinned and installed locally, so
# generation needs no network and no account on the Buf registry.
#   Go:         protoc-gen-go + protoc-gen-connect-go (via buf)
#   TypeScript: ts-proto, run by Bun                  (via buf)
#   Python:     grpcio-tools (same version as the ai service's grpcio runtime)
ARG GO_VERSION=1.27
ARG BUF_VERSION=1.75.0
ARG BUN_VERSION=1.4.4
ARG PYTHON_VERSION=3.14

FROM golang:${GO_VERSION}-alpine AS go-plugins
ENV CGO_ENABLED=0
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.12 \
 && go install connectrpc.com/connect/cmd/protoc-gen-connect-go@v1.21.1

FROM oven/bun:${BUN_VERSION} AS ts-plugin
WORKDIR /opt/ts-proto
RUN echo '{"name":"protogen","private":true}' > package.json && bun add ts-proto@2.18.3

FROM bufbuild/buf:${BUF_VERSION} AS buf

FROM oven/bun:${BUN_VERSION} AS bun

FROM python:${PYTHON_VERSION}-slim AS final
RUN pip install --no-cache-dir grpcio-tools==1.85.1
COPY --from=buf /usr/local/bin/buf /usr/local/bin/buf
COPY --from=go-plugins /go/bin/ /usr/local/bin/
COPY --from=bun /usr/local/bin/bun /usr/local/bin/bun
COPY --from=ts-plugin /opt/ts-proto /opt/ts-proto
# ts-proto's plugin script starts with `#!/usr/bin/env node`; Bun runs it as node.
RUN ln -s /usr/local/bin/bun /usr/local/bin/node \
 && ln -s /opt/ts-proto/node_modules/.bin/protoc-gen-ts_proto /usr/local/bin/protoc-gen-ts_proto
ENV BUF_CACHE_DIR=/tmp/buf-cache
WORKDIR /workspace
