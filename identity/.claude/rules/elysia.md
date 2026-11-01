---
paths:
  - "src/**/*.ts"
---

# Identity: Elysia + Better Auth patterns

Elysia has one job here: serve Better Auth's routes on port 3000. That port is not published; the app reaches it
only through the gateway's `/api/auth/*` proxy. Everything else identity offers goes over gRPC (see
`bun-grpc.md`).

## Elysia version: 1.4, not 2

`package.json` pins `elysia` 1.4.x. The `elysiajs` skill in `.claude/skills/` describes **v2** (its
`references/better-auth.md` is still useful for concepts). Write v1 code and don't upgrade unless asked. v1 docs:
https://elysiajs.com/llms.txt

## Rules

- `src/http.ts` stays a single `.mount(auth.handler)`. **Don't add app-facing routes here.** A new endpoint for
  the app goes in the gateway, which calls an identity gRPC method (add it in `contracts/coolinary/identity/v1/`).
- Better Auth is configured once in `src/auth.ts`. Side effects of auth events (for example, creating the
  personal household on sign-up) go in `databaseHooks` and call plain functions in their own module
  (`src/households.ts`), so gRPC handlers can reuse them.
- JWT claims come from `jwt({ definePayload })`. The gateway reads them in `gateway/src/lib/auth.ts` (`sub`,
  `email`, `householdId`, `plan`). If you change a claim, change the gateway in the same commit. Keep tokens
  short-lived (15 min) and claims small: ids and the plan, no personal data beyond the e-mail.
- Adding a Better Auth plugin: add it in `src/auth.ts`, run `bun run auth:schema` (regenerates
  `src/db/auth-schema.ts`, never edit that file by hand), then `bun run db:generate` and commit the migration.
- Origins the app calls from go in `TRUSTED_ORIGINS`, not in code.
- Check: `bun run lint` in `identity/` (Biome + `tsc`; Bun is at `~/.bun/bin/bun`), `bun run lint:fix` to format.
