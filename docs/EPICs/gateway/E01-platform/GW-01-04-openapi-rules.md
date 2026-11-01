# GW-01-04 OpenAPI rules, tags and export check

| | |
|---|---|
| Service | gateway |
| Epic | [GW-01 Platform](README.md) |
| Release | R1 |
| Size | S |
| Status | todo |
| Depends on | GW-01-02 |
| Unblocks | MOB-02, MOB-03 (generated types) |

## Goal

The OpenAPI document is complete and stable enough that the app's generated client (`mobile/src/api/schema.d.ts`)
never drifts from the gateway, and accidental breaking changes are caught before merge.

## Business rules

1. Every `/v1` route has: one tag, a `summary`, a stable `operationId` (Elysia's generated
   `<method><PathInPascalCase>` form, e.g. `postV1PlansDay`; never renamed once released), a request schema, and
   response schemas for every status it can return, 4xx/5xx using the `Problem` component (GW-01-02).
2. Tags, in this order: `system`, `account`, `household`, `kitchen-profile`, `catalogue`, `fridge`, `intake`,
   `recipes`, `plans`, `cooking`, `list`, `live`, `shops`, `offers`, `flyers`, `notifications`, `savings`.
3. Security: a `bearer` (JWT) scheme; every authenticated route lists it. Public routes: `/v1/system/status`,
   `/v1/profile/presets`.
4. Shared components: `Problem`, `Money` (`{ amountMinor: integer, currency: string }`), `Quantity`
   (`{ value: number, unit: "g" | "ml" | "piece" }`), `Page` (`{ nextPageToken?: string }`), `Quota` (GW-03-02),
   headers `Idempotency-Key`, `X-Request-Id`.
5. Paging: list routes take `pageSize` (1–100, default 20) and `pageToken`, and return `nextPageToken` when more exist.
6. Dates are ISO strings: `YYYY-MM-DD` for calendar dates, RFC 3339 UTC for instants. Money always minor units.
7. **Compatibility rule:** inside `/v1` only additive changes (new routes, new optional request fields, new response
   fields). Removing or renaming a field, a route, an enum value, or making a field required needs a new route
   version. Enum values the app does not know must be handled by the app (MOB-02).
8. `bun run openapi:check` regenerates the document in memory and fails if it differs from
   `contracts/openapi/gateway.json`; `make openapi` writes it (exists). The check runs in the gateway's lint script.
9. `bun run openapi:breaking` compares the document with the one on `main` (using `oasdiff` in the protogen image,
   pinned) and fails on breaking changes; wired into `make proto-breaking`.

## Contract

`/openapi` (UI) and `/openapi/json` keep working; `/api/auth/*` stays excluded.

## Data

No storage.

## Errors

None.

## Acceptance criteria

- [ ] `bun run openapi:check` fails when a route schema changes without re-running `make openapi`.
- [ ] Removing a response field makes `bun run openapi:breaking` fail; adding an optional one passes.
- [ ] Every operation in the exported document has a tag from the list, a summary and an `operationId`.
- [ ] `bun run api:types` in `mobile/` generates types with the `Problem` and `Quota` components.
