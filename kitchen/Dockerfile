# syntax=docker/dockerfile:1.7
# Stages: base -> dev (air + delve, hot reload) ; base -> build -> prod (distroless, non-root)
ARG GO_VERSION=1.27

FROM golang:${GO_VERSION}-alpine AS base
WORKDIR /app
ENV CGO_ENABLED=0 GOFLAGS=-buildvcs=false
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

FROM base AS dev
ARG AIR_VERSION=v1.67.4
ARG DELVE_VERSION=v1.27.2
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
    go install github.com/air-verse/air@${AIR_VERSION} \
 && go install github.com/go-delve/delve/cmd/dlv@${DELVE_VERSION}
COPY . .
EXPOSE 50051 8081 2345
HEALTHCHECK --interval=10s --timeout=3s --start-period=90s --retries=5 \
  CMD ["wget", "-qO-", "http://127.0.0.1:8081/healthz"]
CMD ["air", "-c", ".air.toml"]

FROM base AS build
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
    go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server

FROM gcr.io/distroless/static:nonroot AS prod
COPY --from=build /out/server /app/server
USER nonroot:nonroot
EXPOSE 50051 8081
HEALTHCHECK --interval=15s --timeout=3s --start-period=20s --retries=3 \
  CMD ["/app/server", "healthcheck"]
ENTRYPOINT ["/app/server"]
