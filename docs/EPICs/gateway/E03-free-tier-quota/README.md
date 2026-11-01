# GW-03 Free-tier AI quota and rate limits

Release R1.

## Goal

Free households (anonymous included) get exactly 2 AI day plans and 2 AI scans per week; Plus gets unlimited AI
within a fair-use ceiling; Plus-only features answer with the paywall code; and every route is protected by rate
limits. All counters live in `redis-state` (never evicts). This is how the €0.10/month Free cost cap is enforced
(business model §6).

## Flow for an AI route

1. Idempotency lookup (GW-01-03): a replay returns the stored result and never touches the quota.
2. Plus guard (GW-03-03) for Plus-only routes.
3. Free: reserve one unit (GW-03-01). Over the limit → 429 `quota_reached` (GW-03-02).
4. Call the service. The service answers with `ai_used` true/false.
5. `ai_used: true` → commit the unit; `false`, error or timeout → release it.

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [GW-03-01](GW-03-01-quota-counters.md) | Quota counters with reserve / commit / release | R1 | M | GW-01-03 |
| [GW-03-02](GW-03-02-quota-responses.md) | Quota in responses and the `quota_reached` error | R1 | S | GW-03-01 |
| [GW-03-03](GW-03-03-plus-guard.md) | Plus-only route guard | R1 | S | GW-02-01 |
| [GW-03-04](GW-03-04-rate-limits.md) | Rate limits and Plus fair use | R1 | M | GW-01-01 |
| [GW-03-05](GW-03-05-ai-metrics.md) | AI usage metrics and dashboard panel | R1 | S | GW-03-01 |
