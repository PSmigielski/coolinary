# notifier

Bun. Sends e-mail (SMTP; Mailpit in dev at http://10.77.0.3:8025) and keeps a delivery log. Push notifications (Expo Push API) come with their RPC.

- gRPC `50051`, admin `8081`.
- RPC: `NotifierService.SendEmail`.
- Schema in `src/db/schema.ts`; `bun run db:generate` writes a migration to `drizzle/`. Migrations run at start-up.

Runs with the rest of the stack: `make dev` from the repo root (hot reload, debugger). Settings: `.env.example`
(all variables), `.env.dev` (dev), `.env.prod` (prod, secrets via `*_FILE`). Shared rules:
`../docs/service-conventions.md`. Generated gRPC code is committed; regenerate with `make proto`.
