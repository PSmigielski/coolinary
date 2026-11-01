# GW-01-01 Request context, metadata forwarding and deadline classes

| | |
|---|---|
| Service | gateway |
| Epic | [GW-01 Platform](README.md) |
| Release | R1 |
| Size | S |
| Depends on | — |
| Unblocks | every gateway route task |

## Goal

One place builds the call context for every route: request id, trace header, user, household, member, plan and
language. Services get the same metadata on every call, and every call has a deadline that fits what it does.

## Business rules

1. **Request id.** Use the client's `X-Request-Id` if it is 8–64 characters of `[A-Za-z0-9-]`; otherwise generate a
   UUIDv7. Echo it in the `X-Request-Id` response header on every response, errors included.
2. **Trace.** Forward `traceparent` unchanged when it matches the W3C format; drop it otherwise.
3. **Identity metadata** comes only from the verified JWT: `x-user-id` (`sub`), `x-household-id`, `x-plan`
   (`free` / `plus`). Any incoming HTTP header with one of these names, or with `x-member-id`, `x-language`,
   `x-admin-user-id`, is ignored (never forwarded).
4. **Language (D4).** Parse `Accept-Language`, take the first tag whose primary language is supported (`pl`, `de`,
   `en`; the list is config `SUPPORTED_LANGUAGES`), and forward it as `x-language` (lowercase primary tag, e.g. `pl`).
   No header or nothing supported → `en`. Respond with `Content-Language` set to the value used.
5. **Household and member ids in messages (D3).** Route handlers fill `household_id` and `member_id` request fields
   from the JWT (`householdId`, `memberId` claims, GW-02-01). Path or body ids naming another household are never
   used for authorisation; the owning service checks that the target belongs to `household_id`.
6. **Client IP.** In prod the gateway sits behind Caddy (`10.77.0.2`): take the last address of `X-Forwarded-For`
   only when the socket peer is Caddy's address (config `TRUSTED_PROXIES`); otherwise use the socket address. The IP
   is used for rate limits only and never logged in full (log the /24 for IPv4, /48 for IPv6).
7. **Deadline classes** (config, milliseconds), chosen per route in its definition:

   | Class | Default | Used for |
   |---|---|---|
   | `standard` | 2,000 | reads and simple writes |
   | `heavy` | 10,000 | list generation, export parts, price matrix, batch writes |
   | `ai` | 60,000 | every route whose service may call `ai` per request |
   | `ai_long` | 120,000 | week plan (GW-07-11), cheapest top-up (GW-10-06) |

8. When the client disconnects, the in-flight gRPC call is cancelled (existing `AbortSignal` wiring) and nothing is
   stored for idempotent replay.
9. **Client version (decided, Q020).** The app sends `x-client: coolinary-mobile/<version> (<ios|android>)` on every
   request (MOB-02-01). The gateway parses it into `client.app` and `client.platform`; a missing or malformed header
   is allowed (curl, tests). Config `MIN_APP_VERSION` (empty = off): an app version below it gets
   426 `app_update_required` on every route except `/healthz` and `/v1/system/status`, so too-old apps can be
   stopped. The header is not forwarded to services.
10. Logs: one line per request at `info` with `request_id`, `method`, `route` (template, not the raw path),
   `status`, `duration_ms`, `client` (app version and platform), `user_hash` (first 12 hex chars of SHA-256 of user id + `LOG_HASH_SALT`). Never the
   e-mail, query values, bodies, tokens or free text.

## Contract

No new route. Internal change to `callContext()` in `src/http/app.ts` (moved to `src/http/context.ts`):

```ts
interface RequestContext {
  requestId: string; traceparent?: string; language: "pl" | "de" | "en";
  clientIp: string; client?: { app: string; platform: "ios" | "android" }; user?: AuthUser; // AuthUser gains memberId, isAnonymous (GW-02-01)
}
function grpcContext(ctx: RequestContext, deadline: "standard" | "heavy" | "ai" | "ai_long"): CallContext;
```

Forwarded metadata: `x-request-id`, `traceparent`, `x-user-id`, `x-household-id`, `x-plan`, `x-language`.

## Data

No storage.

## Errors

| Case | Result |
|---|---|
| Malformed `X-Request-Id` | Replaced by a generated id; request continues |
| Unsupported `Accept-Language` | `en`; request continues |
| Malformed or missing `x-client` | request continues, `client` not logged |
| App older than `MIN_APP_VERSION` | 426 `app_update_required` |

## Acceptance criteria

- [ ] A request with `X-Request-Id: abc-12345678` gets the same id back; one with `X-Request-Id: <script>` gets a
      UUIDv7.
- [ ] A request with header `x-household-id: other` and a valid JWT forwards the JWT's household id (checked with a
      fake gRPC server in tests).
- [ ] `Accept-Language: de-AT,de;q=0.9,en;q=0.5` forwards `x-language: de`; `Accept-Language: fr` forwards `en`
      and responds `Content-Language: en`.
- [ ] With `TRUSTED_PROXIES=10.77.0.2`, a spoofed `X-Forwarded-For` from any other peer is ignored.
- [ ] With `MIN_APP_VERSION=1.3.0`, `x-client: coolinary-mobile/1.2.9 (ios)` gets 426 `app_update_required`;
      `1.3.0` passes; a request without `x-client` passes.
- [ ] An `ai`-class route sets a 60 s gRPC deadline; a `standard` route 2 s (asserted in a unit test).
- [ ] Request log lines contain no e-mail, token, query string or body (checked by a log-scanning test).
