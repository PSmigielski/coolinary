# GW-03-02 Quota in responses and the `quota_reached` error

| | |
|---|---|
| Service | gateway |
| Epic | [GW-03 Free-tier AI quota and rate limits](README.md) |
| Release | R1 |
| Size | S |
| Status | todo |
| Depends on | GW-03-01 |
| Unblocks | GW-02-04, MOB-02 (QuotaReached error), MOB-03 |

## Goal

The app always knows how much Free AI is left: every counted route returns the current quota, and running out
gives a clear paywall error.

## Business rules

1. Every 2xx response of a counted route (GW-03-01 rule 1) includes a top-level `quota` object for its kind,
   computed after commit/release: `{ "kind", "limit", "used", "remaining", "nextFreeAt" }`. `used` counts the last
   7 days; `nextFreeAt` is when the oldest counted use frees up (rolling window, Q111), `null` when nothing is used.
   Plus: `limit`, `remaining` are `null`.
2. Response headers on the same routes: `X-Quota-Limit`, `X-Quota-Remaining` (omitted for Plus), `X-Quota-Next-Free-At`
   (RFC 3339; omitted when nothing is used).
3. Over the limit → 429 problem `quota_reached`, `Retry-After` = seconds until `nextFreeAt`,
   `meta: { "kind", "limit", "nextFreeAt", "paywall": true }`. Distinct from `rate_limited` so the app shows the Plus
   offer, not "slow down".
4. `GET /v1/entitlements` carries both kinds (GW-02-04).
5. A replayed idempotent response returns the quota as it was stored (it is a replay); the app re-reads
   `/v1/entitlements` after reconnecting.

## Contract

OpenAPI component `Quota`:

```json
{ "kind": "scan", "limit": 2, "used": 2, "remaining": 0, "nextFreeAt": "2026-11-05T18:42:10Z" }
```

429 example:

```json
{ "type": ".../quota_reached", "title": "Weekly limit reached", "status": 429, "code": "quota_reached",
  "detail": "You've used your 2 free photo scans for the last 7 days.", "requestId": "...",
  "meta": { "kind": "scan", "limit": 2, "nextFreeAt": "2026-11-05T18:42:10Z", "paywall": true } }
```

## Data

Reads GW-03-01 keys.

## Errors

| Case | HTTP | `code` |
|---|---|---|
| Free limit reached | 429 | `quota_reached` |

## Acceptance criteria

- [ ] After the first scan in 7 days, the response has `quota.remaining: 1` and `X-Quota-Remaining: 1`.
- [ ] The third scan gets 429 `quota_reached` with `meta.paywall: true` and `Retry-After` matching `nextFreeAt`.
- [ ] Plus responses carry `quota.limit: null` and no `X-Quota-Remaining` header.
