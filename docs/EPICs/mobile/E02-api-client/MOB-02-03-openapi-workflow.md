# MOB-02-03 OpenAPI types workflow

| | |
|---|---|
| Service | mobile |
| Epic | [MOB-02 API client and error model](README.md) |
| Release | R1 |
| Size | S |
| Depends on | GW-01 (OpenAPI rules) |
| Unblocks | MOB-03-* |

## Goal

The app's request and response types always match the gateway, so a changed route breaks the app's typecheck
instead of breaking users.

## Business rules

1. The only source of API types is `src/api/schema.d.ts`, generated from `contracts/openapi/gateway.json` by
   `bun run api:types`. `make openapi` (repo root) exports the document from the gateway and regenerates the file.
2. Nobody edits `schema.d.ts` by hand and no feature declares its own response interfaces; they use type aliases
   derived from `paths`/`operations`, collected per area in `src/api/types/<area>.ts` (e.g.
   `export type StockItem = operations["getV1Stock"]["responses"][200]["content"]["application/json"]["items"][number]`).
3. Workflow for a gateway change: change the route → `make openapi` → `bun run typecheck` in `mobile/` → fix every
   error → commit `gateway.json`, `schema.d.ts` and the app change together.
4. `bun run check:api` (new script) regenerates the types into a temporary file and fails if they differ from the
   committed `schema.d.ts`; it runs in `make lint` for mobile.
5. The persisted query cache buster (MOB-03-01) changes with every app release, so cached data from an older schema
   is never read by newer code.

## Contract

New `package.json` script `check:api`. No HTTP.

## Data

No storage.

## Errors

| Case | Behaviour |
|---|---|
| `gateway.json` changed but `schema.d.ts` not regenerated | `check:api` fails with the diff |

## Acceptance criteria

- [ ] Renaming a response field in the gateway and running `make openapi` makes `bun run typecheck` fail in mobile.
- [ ] `bun run check:api` passes on a clean checkout and fails after a manual edit of `schema.d.ts`.
