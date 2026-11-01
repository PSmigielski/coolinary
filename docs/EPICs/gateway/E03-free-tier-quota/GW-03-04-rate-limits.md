# GW-03-04 Rate limits and Plus fair use

| | |
|---|---|
| Service | gateway |
| Epic | [GW-03 Free-tier AI quota and rate limits](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | GW-01-01 |
| Unblocks | GW-02-02, GW-02-06, GW-04-04, upload routes |

## Goal

Protect the services, the AI budget and the sign-in flow from abuse and runaway clients, with limits a normal
user never meets.

## Business rules

1. Fixed-window counters in `redis-state`: `INCR` + `EXPIRE` on first use of a window. Subject is the user id for
   authenticated buckets and the client IP (GW-01-01) otherwise.
2. Buckets (config, `limit / window`):

   | Bucket | Subject | Limit | Applies to |
   |---|---|---|---|
   | `api` | user | 300 / 1 min | every authenticated `/v1` request |
   | `public` | IP | 120 / 1 min | unauthenticated `/v1` requests |
   | `otp_send` | IP | 5 / 15 min | `POST /api/auth/email-otp/send-verification-otp` |
   | `otp_send_email` | hashed e-mail | 3 / 15 min | same route; e-mail read from the JSON body, SHA-256 with `RATE_LIMIT_SALT` |
   | `otp_verify` | IP | 10 / 15 min | `POST /api/auth/sign-in/email-otp` |
   | `anon_signup` | IP | 5 / 1 h | `POST /api/auth/sign-in/anonymous` |
   | `ai` | user | 10 / 1 min | every route in the `ai` deadline class |
   | `upload` | user | 30 / 1 h | upload routes (GW-01-05) |
   | `invite_lookup` | user | 20 / 1 h | `GET /v1/invites/{code}`, `POST /v1/invites/{code}/accept` |
   | `export` | user | 3 / 24 h | `POST /v1/me/export` |
   | `ws_connect` | user | 30 / 5 min | `GET /v1/live` upgrades |
   | `shop_suggestion` | user | 10 / 24 h | `POST /v1/shops/suggestions` (GW-10-01) |
   | `own_flyer` | household | 10 / 24 h | `POST /v1/flyers` (GW-10-04) |

3. **Plus fair use:** per household, at most **60 AI responses per day** (UTC day), counted from GW-03-01's Plus
   counting across both kinds and every other AI route. Over it → 429 `fair_use_limit` with `Retry-After` to the next
   UTC midnight. Not shown in the Free/Plus marketing table; stated in the terms (decided, Q113).
4. Over a limit → 429 problem `rate_limited` (or `fair_use_limit`), `Retry-After` = seconds left in the window,
   `meta: { "bucket" }`. Rejected requests still count (no free retries in a loop).
5. Order: `public`/`api` first, then route buckets, then idempotency and quota. A replayed idempotent request counts
   only towards `api`.
6. `redis-state` unreachable: rate limits **fail open** (request continues, metric
   `gateway_ratelimit_errors_total` and a warn log once per minute). Quota keeps its own fail-closed rule (GW-03-01).
7. Health and metrics routes (`/healthz`, admin port) are not limited.

## Contract

Elysia plugin `rateLimit(bucket)`; OpenAPI lists 429 `Problem` on every route.

## Data (redis-state)

| Key | Value | TTL |
|---|---|---|
| `gw:rl:<bucket>:<subject>:<windowStartEpochS>` | request count | window length + 5 s |
| `gw:fair:<householdId>:<YYYY-MM-DD>` | AI responses today (Plus) | 2 days |

## Errors

| Case | HTTP | `code` |
|---|---|---|
| Bucket exceeded | 429 | `rate_limited` |
| Plus fair use exceeded | 429 | `fair_use_limit` |

## Acceptance criteria

- [ ] The sixth OTP request from one IP within 15 minutes gets 429 with `Retry-After` ≤ 900.
- [ ] The fourth OTP request for the same e-mail from different IPs within 15 minutes gets 429.
- [ ] The sixth anonymous sign-up from one IP in an hour gets 429 (identity's own limit is also 5, IDN-03-01).
- [ ] A Plus household's 61st AI response of the day gets 429 `fair_use_limit`.
- [ ] With redis-state stopped, ordinary routes keep working (fail open) while Free AI routes get 503 (GW-03-01).
- [ ] No e-mail address appears in Redis keys or logs (hashed only).

## Notes

- Decided (Q113): Plus fair use is **60 AI responses per household per UTC day**, stated in the terms (about 10× a
  heavy family's use; protects the 15% AI cost cap).
- Decided (Q034): **5 anonymous sign-ups per IP per hour**, the same number identity enforces (a reinstall gives a
  new anonymous account and a new Free quota).
