# AI-02-04 Global daily spend limit

| | |
|---|---|
| Service | ai |
| Epic | [AI-02 Cost tracking and spend limit](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | AI-02-01 |
| Unblocks | AI-02-05, AI-04-01 |

## Goal

A bug, a loop or abuse can never spend more than a fixed amount a day. Above the limit, AI calls stop and every
caller shows its non-AI result (D11).

## Business rules

1. Three counters in **redis-state** (`noeviction`, AOF; README decision 17), in micro-euros, keyed by UTC day or
   month. ai gets a new setting `REDIS_STATE_URL` (also readable as `REDIS_STATE_URL_FILE`) in `config.py`,
   `.env.example`, `.env.dev`, `.env.prod` and the compose files; redis-cache stays for the result cache only:
   - `ai:spend:sync:<YYYY-MM-DD>` — all per-request calls.
   - `ai:spend:sync:free:<YYYY-MM-DD>` — the Free-plan part of it.
   - `ai:spend:batch:<YYYY-MM>` — all batch jobs.
   Keys expire 3 days (daily) or 40 days (monthly) after their period ends.
2. Limits (config, euros):
   - `AI_DAILY_LIMIT_EUR` — default 20.
   - `AI_DAILY_FREE_SHARE` — default 0.4: Free traffic may use at most 40% of the daily limit, so Free can't starve
     Plus.
   - `AI_BATCH_MONTHLY_LIMIT_EUR` — default 60.
3. Before a sync call: if `spend_today + estimated_max_cost > limit` (and the Free share for Free calls), return
   `RESOURCE_EXHAUSTED ai_budget` without calling Claude. Estimated max cost = input estimate (characters ÷ 3.5
   for text, 1,600 tokens per image, 2,000 per PDF page) + `max_tokens` output, at the job's model price.
4. After a call: add the actual cost (AI-02-01) with `INCRBY`. A repair adds its own cost.
5. Before a batch submission: same check against the monthly batch limit using the estimate for all items; the
   actual cost is added once when results are fetched (AI-04-03 rule 5).
6. At 80% of any limit: log `warn` once per period and expose `ai_spend_ratio{bucket}` (0–1) for the alert rule
   (AI-02-02).
7. If Redis is unreachable for the check, allow the call (availability first) and count
   `ai_spend_check_skipped_total`; more than 100 skips in an hour raises the existing error-rate alert.
8. The limits reset at 00:00 UTC (daily) and on the 1st of the month UTC (batch).
9. `AiAdminService.GetSpendStatus` (new service in `contracts/coolinary/ai/v1/admin.proto`, called only by admin,
   D1) returns the three counters for the current period, their limits, the ratio of each, and whether each limit
   is reached, read straight from redis-state. It requires `x-admin-user-id` metadata (`PERMISSION_DENIED`
   otherwise), changes nothing and is not traced in Langfuse. When redis-state is unreachable it returns
   `UNAVAILABLE`.

## Contract

Every RPC can return `RESOURCE_EXHAUSTED ai_budget`. New admin service (rule 9):

```proto
service AiAdminService {
  rpc GetSpendStatus(GetSpendStatusRequest) returns (GetSpendStatusResponse);
}
message GetSpendStatusRequest {}
message SpendBucket { string bucket = 1; int64 spent_eur_micros = 2; int64 limit_eur_micros = 3; double ratio = 4; bool limited = 5; string period = 6; }
message GetSpendStatusResponse { repeated SpendBucket buckets = 1; string computed_at = 2; }  // buckets: sync, sync_free, batch
```

## Data

No schema change. Redis keys as above.

## Errors

| Case | gRPC status |
|---|---|
| Limit reached | `RESOURCE_EXHAUSTED`, `ai_budget` |
| `GetSpendStatus` without `x-admin-user-id` | `PERMISSION_DENIED` |
| `GetSpendStatus` with redis-state down | `UNAVAILABLE` |

## Acceptance criteria

- [ ] With `AI_DAILY_LIMIT_EUR=0.001`, the second live call of the day returns `ai_budget` and makes no HTTP call.
- [ ] With the Free share full, a `free` call is refused while a `plus` call still succeeds.
- [ ] Counters reset at UTC midnight (test with a fake clock).
- [ ] Redis down: calls succeed and `ai_spend_check_skipped_total` increases.
- [ ] After two live calls costing €0.004 in total, `GetSpendStatus` reports the `sync` bucket at 4,000 micro-euros
      with `limited = false`.

## Notes

- Decided (README decision 17): the counters live in redis-state so an eviction can't silently reset the limit.
- Decided (Q018): default limits are €20 a day overall, Free at most 40% of it, and €60 a month for batch jobs. They
  are sized for the first 1,000 households and are config values to raise with growth.
