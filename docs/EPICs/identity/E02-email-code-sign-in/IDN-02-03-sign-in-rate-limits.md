# IDN-02-03 Sign-in rate limits and throttles

| | |
|---|---|
| Service | identity |
| Epic | [IDN-02 Sign-in with e-mail code](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | IDN-02-02, GW-02 (forwards the client IP in `x-forwarded-for`) |
| Unblocks | IDN-03-01 |

## Goal

Codes can't be guessed, inboxes can't be flooded and e-mail cost stays small, without blocking normal people who
mistype or ask for a second code.

## Business rules

1. Turn on Better Auth's built-in rate limiter in every environment, with `customStorage` backed by redis-state
   (keys `identity:rl:*`). Do not configure Better Auth `secondaryStorage`: it would move sessions out of Postgres.
2. Client IP is read from the first entry of `x-forwarded-for` (`advanced.ipAddress.ipAddressHeaders`). Caddy sets
   it in prod; the gateway proxy must pass it on unchanged (GW-02). Without it every request would share the
   gateway's IP.
3. Per IP (Better Auth custom rules):

   | Path | Limit |
   |---|---|
   | `/email-otp/send-verification-otp` | 10 per hour |
   | `/sign-in/email-otp` | 20 per hour |
   | `/sign-in/anonymous` | 5 per hour (IDN-03-01) |
   | everything else | Better Auth default (100 per 10 s) |

4. Per e-mail, counted in `sendVerificationOTP` with redis-state counters (key = SHA-256 of the normalised e-mail,
   never the address): at most 1 code per 60 s, 5 per hour, 10 per day.
5. Per e-mail, failed verifications: after 10 wrong codes in one hour (across several codes), every sign-in attempt
   for that e-mail is refused for 1 hour, even with a correct new code. With rule 4 this caps guessing at well under
   1 in 10,000 per day.
6. Over any limit: HTTP 429 with `{"code": "rate_limited", "retry_after_seconds": n}` and a `Retry-After` header.
7. If redis-state is down: sending a code returns 503 `code_not_sent` (fail closed, it costs money and is the abuse
   target); verifying a code still works, because the code itself is in Postgres.
8. Metric `auth_rate_limited_total{rule}` with `rule` in `ip_send`, `ip_verify`, `ip_anonymous`, `email_send`,
   `email_fail_lock`.

## Contract

No new endpoint. 429 responses on the Better Auth endpoints listed above.

## Data

No schema change. redis-state keys:

| Key | TTL | Holds |
|---|---|---|
| `identity:rl:<better-auth key>` | per rule | Better Auth counters |
| `identity:otp:send:<sha256>:{m,h,d}` | 60 s / 1 h / 24 h | sends per e-mail |
| `identity:otp:fail:<sha256>` | 1 h | wrong codes per e-mail |

## Errors

| Case | HTTP | Body `code` |
|---|---|---|
| Any limit hit | 429 | `rate_limited` |
| redis-state down when sending | 503 | `code_not_sent` |

## Acceptance criteria

- [ ] A second "send code" for the same e-mail within 60 s returns 429 with `retry_after_seconds` ≤ 60.
- [ ] Six sends in one hour (61 s apart) → the sixth returns 429.
- [ ] Ten wrong codes in one hour lock the e-mail; a correct new code is refused until the hour has passed.
- [ ] Two clients with different `x-forwarded-for` values are counted separately.
- [ ] No redis-state key contains an e-mail address in clear text.
