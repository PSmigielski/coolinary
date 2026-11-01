# syntax=docker/dockerfile:1.7
# Stages: base -> deps -> dev (bun --watch) ; deps -> build (bundle) -> prod (distroless, non-root)
ARG BUN_VERSION=1.4.4

FROM oven/bun:${BUN_VERSION} AS base
WORKDIR /app

FROM base AS deps
COPY package.json bun.lock ./
RUN --mount=type=cache,target=/root/.bun/install/cache bun install --frozen-lockfile

FROM base AS dev
COPY --from=deps /app/node_modules ./node_modules
COPY . .
ENV NODE_ENV=development
EXPOSE 50051 8081 6499
HEALTHCHECK --interval=10s --timeout=3s --start-period=30s --retries=5 \
  CMD ["bun", "src/healthcheck.ts"]
CMD ["bun", "--watch", "src/index.ts"]

FROM deps AS build
COPY . .
RUN bun run typecheck \
 && bun build src/index.ts src/healthcheck.ts --target=bun --outdir=dist --sourcemap=linked

FROM oven/bun:${BUN_VERSION}-distroless AS prod
WORKDIR /app
COPY --from=build /app/dist ./dist
COPY --from=build /app/drizzle ./drizzle
ENV NODE_ENV=production
USER 65532:65532
EXPOSE 50051 8081
HEALTHCHECK --interval=15s --timeout=3s --start-period=10s --retries=3 \
  CMD ["/usr/local/bin/bun", "dist/healthcheck.js"]
CMD ["dist/index.js"]
