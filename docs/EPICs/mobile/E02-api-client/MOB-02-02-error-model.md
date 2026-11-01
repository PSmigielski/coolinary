# MOB-02-02 Typed error model and retry policy

| | |
|---|---|
| Service | mobile |
| Epic | [MOB-02 API client and error model](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | MOB-02-01, MOB-02-04, GW-01 (error model) |
| Unblocks | MOB-02-05, MOB-03-*, MOB-04-03 |

## Goal

Screens get one small set of error types they can explain to the user ("You've used both free day plans for the last 7 days",
"This needs Plus", "You're offline") instead of raw HTTP statuses.

## Business rules

1. Gateway error body (GW-01-02): `application/problem+json` (RFC 9457) with a stable `code`; the app localises by
   `code`, never shows `detail` as is:

   ```json
   { "type": "…/problems/quota_reached", "title": "…", "status": 429, "code": "quota_reached",
     "detail": "…", "requestId": "…", "fields": [ { "path": "/items/0/quantity", "code": "too_small" } ],
     "meta": { "kind": "scan", "limit": 2, "nextFreeAt": "…", "paywall": true } }
   ```
   The server's quota kinds `dayPlan` and `scan` map to the app's gate names `day_plan` and `photo_scan`.
2. Mapping to one discriminated union `ApiError` (`src/api/errors.ts`):

   | Kind | When | Extra fields |
   |---|---|---|
   | `offline` | fetch failed and MOB-02-04 says no connection | — |
   | `serverDown` | fetch failed while online, or `502`/`503`/`504`, or gateway probe fails | `retryAfterMs?` |
   | `timeout` | abort by the route timeout | `timeoutClass` |
   | `unauthorized` | `401` after the retry in MOB-01-06 | — |
   | `notMember` | `403` code `not_member` | — |
   | `plusRequired` | `403` code `plus_required` | `feature` |
   | `forbidden` | other `403` | — |
   | `notFound` | `404` | — |
   | `conflict` | `409` (except `idempotency_in_progress`) | `current?` (server copy of the resource), `code` |
   | `validation` | `400`/`422` | `fields: { path: string; code: string }[]` |
   | `quotaReached` | `429` code `quota_reached` | `feature: "day_plan" \| "photo_scan"` (from `meta.kind`), `limit`, `nextFreeAt` (when the oldest use of the last 7 days frees up, Q111), `libraryPlanAvailable?` |
   | `rateLimited` | other `429` | `retryAfterMs` from `Retry-After` (default 30 s) |
   | `aiUnavailable` | `503` code `ai_unavailable` (only if a route has no non-AI fallback) | — |
   | `appUpdateRequired` | `426` code `app_update_required` (GW-01-01, Q020) | — |
   | `internal` | `500` or an unknown body | `requestId` |
3. Every `ApiError` carries `requestId` (from the body or the request) so support can find the call.
4. Errors never include the request body, the user's free text, e-mail addresses or tokens.
5. **Query retries** (TanStack Query `retry` function in `src/lib/query.ts`, replacing `retry: 1`): retry only
   `offline` (after reconnect, via `onlineManager`), `serverDown` and `timeout`; at most 2 retries; delays 1 s, then
   4 s, each ± 20 % jitter. Never retry any 4xx.
6. **Mutation retries**: 0 by default. Offline-capable mutations go through the outbox (MOB-04) instead; AI calls
   follow MOB-02-01 rule 6.
7. `rateLimited` pauses that route class for `retryAfterMs`; calls during the pause fail fast with the same error.
8. Helper `isRetryable(error)` and `userMessageKey(error)` (a translation key, not a sentence) are exported; screens
   use the key.

## Contract

```ts
export type ApiError =
  | { kind: "offline"; requestId: string }
  | { kind: "quotaReached"; requestId: string; feature: "day_plan" | "photo_scan"; limit: number; used: number; nextFreeAt: string }
  | { kind: "plusRequired"; requestId: string; feature: string }
  | { kind: "conflict"; requestId: string; code: string; current?: unknown }
  | { kind: "validation"; requestId: string; fields: { path: string; code: string }[] }
  // … one member per row of the table above
  ;
export function toApiError(input: { response?: Response; body?: unknown; cause?: unknown; requestId: string }): ApiError;
```

## Data

No storage.

## Errors

This task defines them.

## Acceptance criteria

- [ ] Unit tests map one recorded gateway response per row of the table to the right `kind`.
- [ ] Airplane mode gives `offline`, not `serverDown`; gateway stopped with Wi-Fi on gives `serverDown`.
- [ ] A `400` is requested exactly once; a `503` is requested three times with growing delays.
- [ ] No error object serialised in dev logs contains a body field.
