# IDN-05-02 `GetEntitlements`

| | |
|---|---|
| Service | identity |
| Epic | [IDN-05 Plan and entitlements](README.md) |
| Release | R1 |
| Size | M |
| Depends on | IDN-05-01, IDN-04-04 |
| Unblocks | GW-02 (`GET /v1/entitlements`), GW-03 (quota limits), MOB-03 (paywall data), KIT-01, PLN-01, OFF-01 (plan checks, Q004) |

## Goal

The app shows exactly what the household's plan includes, from one list kept in one place, and the gateway's
quota uses the same numbers.

## Business rules

1. `src/entitlements.ts` holds one static table: for each feature key, its Free and Plus value. It mirrors
   features.md §6 and the Free/Plus table in docs/EPICs/README.md:

   | Key | Free | Plus |
   |---|---|---|
   | `ai_day_plans` | 2 per week | unlimited |
   | `ai_photo_scans` | 2 per week | unlimited |
   | `ai_other` (swaps, rewrites, week and family plans, deals, top-up, own flyer) | not allowed | unlimited |
   | `expiry_reminders` | not allowed | allowed |
   | `calorie_goals` | 1 | unlimited |
   | `per_person_diets` | not allowed | allowed |
   | `shops_with_offers` | 2 | unlimited |
   | `trip_plan` | not allowed | allowed |
   | `shared_household` | not allowed | 5 accounts |
   | `near_shop_reminders` | not allowed | allowed |
   | `savings_report` | not allowed | allowed |

2. `GetEntitlements(household_id, actor_user_id)` returns the plan, `plan_source`, `plan_expires_at`,
   `is_anonymous` (of the actor), and one entry per key: `allowed`, `limit` (0 = no limit), `period` (`week` or
   none; `week` is a rolling 7 days, Q111). The actor must be a member (IDN-04-04). Service callers (kitchen,
   planner, offers, gateway jobs) may leave `actor_user_id` empty; then the membership check is skipped (internal
   network only).
3. Remaining quota is not identity's data; the gateway adds usage from its counters (GW-03) when it answers the app.
4. New keys are only ever added; a key is never renamed or reused (the app keeps an old build's mapping).
5. Decided (Q025): the gateway reads quota limits from `GetEntitlements` and caches them for 5 minutes; it keeps no
   constants of its own.
6. Anonymous accounts get the Free values (business-model §3).
7. Decided (Q004): this is the one place every service learns a household's plan. The gateway, kitchen, planner and
   offers call `GetEntitlements` whenever a request needs the plan, and background jobs call it per household per
   run. No service keeps a local plan copy; a caller may memoise the answer only for the length of one request or
   one job run. `household.plan_changed` is still emitted (IDN-05-01) for consumers that react to a change.
8. Cost: one primary-key read of `households` plus the static table, no joins. Target: p99 ≤ 10 ms inside identity
   at 200 calls per second.

## Contract

```proto
rpc GetEntitlements(GetEntitlementsRequest) returns (GetEntitlementsResponse);

message GetEntitlementsRequest { string household_id = 1; string actor_user_id = 2; }  // actor empty = service caller
message GetEntitlementsResponse {
  Plan plan = 1;
  string plan_source = 2;          // "default", "manual", "store"
  string plan_expires_at = 3;      // RFC 3339 or ""
  bool is_anonymous = 4;
  repeated Entitlement entitlements = 5;
}
message Entitlement {
  string key = 1;
  bool allowed = 2;
  int32 limit = 3;                 // 0 = no limit
  string period = 4;               // "week" or ""
}
```

## Data

No schema change.

## Errors

| Case | gRPC status |
|---|---|
| Actor not a member | `PERMISSION_DENIED` |
| Unknown household | `NOT_FOUND` |

## Acceptance criteria

- [ ] A Free household gets `ai_day_plans` `{allowed: true, limit: 2, period: "week"}` and `trip_plan`
      `{allowed: false}`.
- [ ] A Plus household gets `shared_household` `{allowed: true, limit: 5}` and `ai_day_plans` `{limit: 0}`.
- [ ] A unit test compares the table with the Free/Plus table in docs/EPICs/README.md (keys and values).
- [ ] A call with an empty `actor_user_id` returns the household's entitlements without a membership check.
- [ ] A load test at 200 calls per second keeps p99 at or under 10 ms inside identity.
