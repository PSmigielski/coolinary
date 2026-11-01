# GW-01-02 Problem JSON error model and gRPC → HTTP mapping

| | |
|---|---|
| Service | gateway |
| Epic | [GW-01 Platform](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | GW-01-01 |
| Unblocks | every gateway route task, MOB-02 (error model) |

## Goal

Every error from the gateway has the same JSON shape and a stable machine code, so the app can show the paywall,
the quota message, "try again" or a field error without parsing text. Replaces today's `{ "error": "..." }` bodies
and the blanket `502`.

## Business rules

1. Every non-2xx response under `/v1` has `Content-Type: application/problem+json` (RFC 9457) and this body:

   ```json
   {
     "type": "https://api.coolinary.app/problems/plus_required",
     "title": "Plus required",
     "status": 403,
     "code": "plus_required",
     "detail": "Swapping a meal is part of Plus.",
     "requestId": "0192f0c4-...",
     "fields": [{ "path": "/items/0/quantity", "code": "too_small" }],
     "meta": { "feature": "swap_meal" }
   }
   ```

   `fields` only for validation errors; `meta` only where a task defines it (quota, rate limit, conflicts).
   `type` base URL is config `PROBLEM_TYPE_BASE`. `detail` is English, safe to show, never contains personal data;
   the app localises by `code`.
2. **Services name the reason** in trailing metadata `x-error-code` (snake_case, ≤ 40 chars). The gateway passes it
   through as `code` when it is in the allowlist `src/http/error-codes.ts`; unknown codes become the generic code of
   the status. (Cross-service convention, to be added to `docs/service-conventions.md`.)
3. **Mapping** (gRPC status → HTTP status, default `code`):

   | gRPC | HTTP | Default `code` | Special codes |
   |---|---|---|---|
   | `INVALID_ARGUMENT` | 400 | `invalid_argument` | |
   | `UNAUTHENTICATED` | 401 | `unauthorized` | |
   | `PERMISSION_DENIED` | 403 | `forbidden` | |
   | `NOT_FOUND` | 404 | `not_found` | |
   | `ALREADY_EXISTS` | 409 | `already_exists` | `already_member`, `household_full` |
   | `ABORTED` | 409 | `version_conflict` | (optimistic concurrency) |
   | `FAILED_PRECONDITION` | 422 | `rule_violated` | `plus_required` → **403**; `invite_expired`, `invite_used` → **410** |
   | `RESOURCE_EXHAUSTED` | 429 | `rate_limited` | `ai_spend_limit` → **503** `ai_unavailable` (spend is not revealed) |
   | `UNAVAILABLE` | 503 | `service_unavailable` | `ai_unavailable` |
   | `DEADLINE_EXCEEDED` | 504 | `timeout` | |
   | `CANCELLED` | 499 (not sent; client left) | | |
   | `INTERNAL`, `UNKNOWN`, `DATA_LOSS`, `UNIMPLEMENTED` | 500 | `internal` | |

4. Gateway's own errors: request validation → 400 `validation` with `fields` (Elysia `VALIDATION`); missing or bad
   JWT → 401 `unauthorized`; anonymous on a registered-only route → 403 `registration_required` (GW-02-01); unknown
   route → 404 `not_found`; body too large → 413 `payload_too_large`; wrong media type → 415
   `unsupported_media_type`; quota → 429 `quota_reached` (GW-03-02); rate limit → 429 `rate_limited` with
   `Retry-After` (GW-03-04); app version below `MIN_APP_VERSION` → 426 `app_update_required` (GW-01-01, Q020).
5. For 500 the `detail` is always "Something went wrong." The gRPC message is logged (`error` field) but never
   returned. For 4xx the gRPC message is returned as `detail` only when the service set `x-error-code` (services
   promise those messages are safe to show).
6. `503` and `504` carry `Retry-After: 5` (seconds); `ai_unavailable` carries `Retry-After: 60`.
7. Request validation rejects unknown body properties (`additionalProperties: false`) on every route. This also
   guards the privacy rule: a body carrying `weightKg`, `heightCm`, `age` or similar is refused with 400.
8. The four existing routes (`/v1/products/{ean}`, `/v1/me`, `/v1/household`, `/v1/system/status`) move to this
   model in the same change; their `400/404/502 { error }` responses disappear.
9. Metrics: `http_server_requests_total` keeps its labels; add `gateway_errors_total{code}` (code from the allowlist
   only, so cardinality stays bounded).

## Contract

Shared OpenAPI component `Problem` (all fields above; `fields` and `meta` optional), referenced by every route's
4xx/5xx responses. `src/http/errors.ts` exports `toProblem(err, ctx)` and `problem(status, code, detail, extra?)`.

## Data

No storage.

## Errors

This task defines them; see the mapping table.

## Acceptance criteria

- [ ] A fake service returning `FAILED_PRECONDITION` with `x-error-code: plus_required` produces HTTP 403,
      `code: plus_required`, `Content-Type: application/problem+json`.
- [ ] `FAILED_PRECONDITION` without `x-error-code` produces 422 `rule_violated`; with an unlisted code
      `foo_bar` it also produces `rule_violated`.
- [ ] `RESOURCE_EXHAUSTED` + `ai_spend_limit` produces 503 `ai_unavailable` with `Retry-After: 60`.
- [ ] `INTERNAL` with message "pq: duplicate key ..." returns `detail: "Something went wrong."` and logs the message.
- [ ] Every response includes `requestId` matching the `X-Request-Id` header.
- [ ] `PUT` with an extra property in the body returns 400 `validation` with the property's path in `fields`.
- [ ] `GET /v1/products/0000` returns 400 `validation`; an unknown EAN returns 404 `not_found` (existing route moved).
- [ ] The OpenAPI document lists `Problem` for every 4xx/5xx response (checked in GW-01-04).

## Notes

- `x-error-code` in trailing metadata and `ABORTED` → 409 `version_conflict` follow decision 18 in
  docs/EPICs/README.md; the `docs/service-conventions.md` update is listed there. The model choice per plan
  (decision 15) is not the gateway's concern: it only forwards `x-plan` from the JWT (GW-01-01).
