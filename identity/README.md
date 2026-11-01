# identity

Bun + Elysia + Better Auth. Accounts, sessions and JWTs (Better Auth), households and their members. Every new account gets its own household. JWTs carry `householdId` and `plan`; the gateway verifies them against `/api/auth/jwks`.

- Better Auth HTTP on `3000` (reached only through the gateway), gRPC `50051`, admin `8081`.
- RPC: `IdentityService.GetHousehold`.
- Events: changes other services follow are written with `emit(tx, type, householdId, payload)`
  (`src/events/outbox.ts`) to `outbox_events`, in the same transaction as the change. Definitions: `src/events/types.ts`.
- Schema: `src/db/schema.ts` (households, outbox) + `src/db/auth-schema.ts` (Better Auth's tables; regenerate with `bun run auth:schema` after adding a Better Auth plugin), then `bun run db:generate` writes a migration to `drizzle/`. Migrations run at start-up.

Runs with the rest of the stack: `make dev` from the repo root (hot reload, debugger). Settings: `.env.example`
(all variables), `.env.dev` (dev), `.env.prod` (prod, secrets via `*_FILE`). Shared rules:
`../docs/service-conventions.md`. Generated gRPC code is committed; regenerate with `make proto`.
