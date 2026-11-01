# MOB-02-01 Request middleware

| | |
|---|---|
| Service | mobile |
| Epic | [MOB-02 API client and error model](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | MOB-01-04, GW-01 (request context, idempotency keys) |
| Unblocks | MOB-02-02, MOB-03-*, MOB-04-02, MOB-07-02 |

## Goal

Every call to the gateway carries the same headers, has a timeout, and can be retried safely, without each feature
repeating that code.

## Business rules

1. Extend the existing `client` in `src/api/client.ts` (openapi-fetch) with one middleware `src/api/middleware.ts`
   registered by `client.use(...)`. The `$api` hooks keep working unchanged.
2. Headers on every `/v1/*` request:

   | Header | Value |
   |---|---|
   | `Authorization` | `Bearer <jwt>` from `tokenProvider.getToken()` (MOB-01-04); omitted for public routes (`/v1/system/status`, `/healthz`) |
   | `x-request-id` | UUIDv4 created once per logical call and reused on its retries |
   | `Accept-Language` | the app language: device locale from `expo-localization`, reduced to `pl`, `de` or `en`; anything else → `en` |
   | `Idempotency-Key` | on `POST`, `PUT`, `PATCH`, `DELETE` only, see rule 4 |
   | `x-client` | `coolinary-mobile/<app version> (<ios|android>)` from `expo-constants` |

   Decided (Q020): `x-client` is sent on every request; the gateway logs it and answers 426 `app_update_required`
   to versions below its `MIN_APP_VERSION` (GW-01-01). The app shows an "update the app" screen on 426.
3. The language comes from one function `appLanguage()`; a later settings screen may override it. The gateway turns
   `Accept-Language` into `x-language` (README D4).
4. **Idempotency keys** are created by the caller, once per user action (a hook's `mutationFn` or an outbox entry,
   MOB-04-01), and passed as a header parameter. The middleware adds a fresh key only if the caller forgot one and
   logs a dev warning. Retries and offline replays reuse the same key, so the gateway (GW-01) returns the stored
   response instead of doing the work twice.
5. **Timeouts** per route class, applied as an abort signal combined with the caller's own signal:

   | Class | Routes | Timeout |
   |---|---|---|
   | default | everything else | 15 s |
   | AI | `POST /v1/plans/day`, `/v1/plans/week`, `/v1/plans/from-deals`, `.../swap`, `.../replace-ingredient`, `/v1/dishes/{id}/rewrite`, `/v1/dishes/{id}/restyle`, `/v1/intake/text`, `/v1/list/cheapest-top-up` | 60 s |
   | upload | `/v1/intake/photo`, `/v1/intake/label`, `/v1/intake/receipt`, `/v1/flyers` | 120 s (MOB-07-02) |

   The table lives in one constant `ROUTE_TIMEOUTS`; Hermes support for `AbortSignal.timeout` and `AbortSignal.any`
   is checked at implementation, with an `AbortController` + `setTimeout` fallback.
6. AI routes with a timeout get one automatic retry after 2 s with the same `Idempotency-Key`. If the gateway answers
   `409 idempotency_in_progress`, the client waits 3 s and asks again, until 60 s in total have passed.
7. Responses: the middleware reads quota headers (MOB-02-05) and passes every non-2xx response to the error mapper
   (MOB-02-02). It never logs request or response bodies; dev logs show method, path template, status, duration and
   request id only.
8. `API_URL` must be `https://` in release builds; the app refuses to start a release build with an `http://` URL.

## Contract

```ts
export interface CallOptions { idempotencyKey?: string; signal?: AbortSignal; timeoutClass?: "default" | "ai" | "upload" }
export function newIdempotencyKey(): string; // expo-crypto randomUUID()
```

New packages: `expo-localization`, `expo-crypto` (both in Expo Go; install with `bunx expo install`).

## Data

No storage.

## Errors

| Case | Behaviour |
|---|---|
| Timeout | `TimeoutError` (MOB-02-02); AI routes per rule 6 |
| No token obtainable | The request is not sent; the token error is returned |

## Acceptance criteria

- [ ] A request log from the gateway shows the five headers on an authenticated `POST` and no `Authorization` on
      `/v1/system/status`.
- [ ] Two retries of the same mutation reach the gateway with the same `Idempotency-Key` and `x-request-id`.
- [ ] A simulated 70 s day-plan call aborts at 60 s, retries once, and the gateway sees one AI call.
- [ ] With the phone set to French, `Accept-Language` is `en`; set to German, it is `de`.
