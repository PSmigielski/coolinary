# PLN-08-07 ComposeDayPlan RPC, quota signal, idempotency

| | |
|---|---|
| Service | planner |
| Epic | [PLN-08 AI day plan](README.md) |
| Release | R1 |
| Size | M |
| Depends on | PLN-08-01 … PLN-08-06 |
| Unblocks | GW-07 (`POST /v1/plans/day`), GW-03 (quota counting), MOB-03 |

## Goal

The public RPC behind "plan my day". The gateway decides whether AI is allowed (Free quota, GW-03); planner reports
whether AI was actually used, so only real AI plans count, and a retried request never costs twice.

## Business rules

1. Request: `household_id`, `date` (default today in the household's time zone), `meals` (default all three),
   `cuisine_overrides` per meal, `ai_allowed` (set by the gateway, never by the app), `request_id` (the gateway's
   idempotency key, GW-01), `language`.
2. `ai_allowed=false` → library plan (PLN-08-05) with `AI_NOT_ALLOWED`.
3. Idempotency: a second call with the same `(household_id, request_id)` within 24 h returns the stored plan with
   `replayed=true` and makes no AI call. `quota_consumed` is false on replays.
4. Concurrency: one compose per household at a time (Postgres advisory lock on the household id); a second concurrent
   call waits up to 50 s, then returns ABORTED "plan in progress".
5. Response: the plan (status `proposed`), `ai_used`, `ai_meals_kept` (PLN-08-04 rule 8), `fallback_reason`,
   `quota_consumed` = `ai_used && ai_meals_kept > 0 && !replayed`. The gateway increments the Free counter only when
   `quota_consumed` is true.
6. A new proposal for a date replaces the previous proposal (old one `superseded`), never an accepted plan.
7. Anonymous accounts behave like Free (README "Free and Plus"); Plus has no limit here (the gateway passes
   `ai_allowed=true`).
8. End-to-end deadline: 55 s (two AI calls at most, PLN-08-04 rule 3).

## Contract

```proto
rpc ComposeDayPlan(ComposeDayPlanRequest) returns (ComposeDayPlanResponse);

message ComposeDayPlanRequest {
  string household_id = 1;
  string date = 2;
  repeated Meal meals = 3;
  map<int32, coolinary.common.v1.Cuisine> cuisine_overrides = 4;   // key = Meal
  bool ai_allowed = 5;
  string request_id = 6;
  string language = 7;
}
message ComposeDayPlanResponse {
  Plan plan = 1;
  bool ai_used = 2;
  int32 ai_meals_kept = 3;
  FallbackReason fallback_reason = 4;
  bool quota_consumed = 5;
  bool replayed = 6;
}
```

## Data

`plans.request_id` with a unique index per household (PLN-08-01).

## Errors

| Case | gRPC status | Message |
|---|---|---|
| Date more than 14 days ahead or in the past | INVALID_ARGUMENT | `date out of range` |
| Another compose running > 50 s | ABORTED | `plan in progress` |
| Profile or stock unavailable | UNAVAILABLE | no plan possible |

## Acceptance criteria

- [ ] Free household, gateway sends `ai_allowed=false`: library plan, `quota_consumed=false`.
- [ ] Same `request_id` twice: one AI call, the second response `replayed=true, quota_consumed=false`.
- [ ] AI down: library plan with `AI_UNAVAILABLE`, `quota_consumed=false`.
- [ ] An AI plan where every meal fell back reports `ai_meals_kept=0, quota_consumed=false`.
