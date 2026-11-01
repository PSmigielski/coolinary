---
paths:
  - "src/**/*.ts"
  - "scripts/**/*.ts"
---

# Gateway: Elysia patterns

The gateway is the app's only HTTP entry point. It owns no data and has no database: it checks the JWT, turns an
HTTP request into gRPC calls on the service that owns the data, and turns the reply into JSON. Business rules
belong in the owning service, not here. The app never reaches `ai`; the gateway only checks its health.

## Elysia version: 1.4, not 2

`package.json` pins `elysia` 1.4.x. The `elysiajs` skill in `.claude/skills/` describes **v2**. Use it for
concepts, but write v1 code, and don't upgrade unless asked. v1 docs: https://elysiajs.com/llms.txt

| Write (v1.4) | Not (v2 only) |
|---|---|
| `.get(path, handler, options)`: handler before options | `.get(path, options, handler)` |
| Macro `{ resolve() }`, `.resolve()` after validation, `.derive()` before | macro `derive` |
| Hook scope `{ as: "scoped" }` / `"global"` | `'plugin'` scope |
| `return status(404, { error })` | `problem()`, `HTTPError.id()`, `onError(Class, fn)` |
| `t` from `elysia` (TypeBox 0.34), `typeof Schema.static` | `typebox` 1.x, `t.Codec` |

## Routes

- Routes live in `createApp()` in `src/http/app.ts`. When the file grows, split by feature into named plugins
  (`new Elysia({ name: "kitchen", prefix: "/v1/products" })`) and `.use()` them in `createApp()`. Move the
  `auth` macro into its own named plugin first and `.use()` it in each feature, or `auth: true` won't type-check.
- Every route declares `t` schemas for `params` / `query` / `body`, **and a `response` schema for every status it
  returns**, plus `detail: { summary, tags }`. These become the OpenAPI document and the app's API types. After
  changing a route, run `make openapi` and commit `contracts/openapi/gateway.json` with it.
- Let `t` schemas do the input checks (format, length, pattern). Elysia rejects bad input with 422 on its own.
- Paths are versioned: `/v1/...`. JSON fields are camelCase. Proto enums go out as strings via the generated
  `xxxToJSON()` helpers, never as numbers.
- Routes that need a signed-in user set `auth: true`. The macro verifies the JWT locally against identity's JWKS
  (`src/lib/auth.ts`) and puts `user` (`id`, `householdId`, `plan`) in the context. Never call identity to
  authenticate a request.
- Register hooks (`onRequest`, `onAfterResponse`, `onError`) before the routes they apply to. Request metrics and
  error logging already exist in `createApp()`; don't add a second mechanism.

## Calling services

Always go through `unary()` from `src/grpc/unary.ts`, with `callContext(request, user)`:

```ts
const res = await unary<GetHouseholdResponse>("identity", "GetHousehold", (md, o, cb) =>
  clients.identity.getHousehold({ householdId }, md, o, cb), callContext(request, user));
```

`unary()` sets the deadline (2 s by default; pass `timeoutMs` for slower calls), forwards `x-request-id`,
`traceparent`, `x-user-id`, `x-household-id` and `x-plan`, cancels the call when the client disconnects, and
records client metrics. Don't call `clients.*` directly. Add a new service client in `src/grpc/clients.ts` with
`channelOptions` (keepalive). Bun services use unary RPCs only.

## Errors

Expected failures are returned, not thrown: `return status(code, { error: "..." })`, and the code is listed in
`response`. Map gRPC status codes with the `status` enum from `@grpc/grpc-js` (not raw numbers):

| gRPC | HTTP |
|---|---|
| `INVALID_ARGUMENT` | 400 |
| `UNAUTHENTICATED` / `PERMISSION_DENIED` | 401 / 403 |
| `NOT_FOUND` | 404 |
| anything else (`UNAVAILABLE`, `DEADLINE_EXCEEDED`, `INTERNAL`, ...) | 502 |

Pass the backend's `details` to the client only for 4xx; for 502, send a generic message and `log.error` the
details (they can contain database errors).

## Everything else

- Settings: add to `src/config.ts` with `env()` (supports `X_FILE`), and to `.env.example`, `.env.dev` and
  `.env.prod`.
- Logs: `log.info(msg, fields)` from `src/lib/log.ts`. Never log e-mail addresses, free text from users or health data.
- `src/gen/` is generated: never edit it; change `contracts/` and run `make proto`.
- Check: `bun run typecheck` in `gateway/` (Bun is at `~/.bun/bin/bun`). Test routes without a port with
  `createApp().handle(new Request("http://localhost/v1/..."))` under `bun test`.
