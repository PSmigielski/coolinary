# IDN-10-03 Set a plan by hand

| | |
|---|---|
| Service | identity |
| Epic | [IDN-10 Admin RPCs and admin role](README.md) |
| Release | R1 |
| Size | S |
| Depends on | IDN-10-01, IDN-05-01 |
| Unblocks | ADM-05; testing every Plus feature before payments exist |

## Goal

Until payments are built, the founder can give Plus to testers, founding-plan buyers and support cases, optionally
until a date.

## Business rules

1. `SetPlan(household_id, plan, expires_at, reason)` calls `changePlan` with `source = manual` for Plus and
   `source = default` for Free, `changedBy` = the admin id from `x-admin-user-id`.
2. `expires_at` is optional for Plus (empty = until changed); ignored for Free. Validation as IDN-05-01 rule 3.
3. `reason` is required, 3–200 characters (e.g. "founding annual plan, paid 19 EUR on 2026-10-12").
4. A household with a store subscription (IDN-05-05, later) can't be set to Free by hand while the store says Plus.
5. The answer carries the new plan state and `token_stale: true` (the household's members see it on their next
   token, at most 15 minutes).
6. Setting Free on a shared household keeps every member; only new invites are blocked (IDN-05-04, Q187, R2).

## Contract

```proto
message SetPlanRequest {
  string household_id = 1;
  Plan plan = 2;
  string expires_at = 3;   // RFC 3339 or ""
  string reason = 4;
}
message SetPlanResponse {
  Plan plan = 1; string plan_source = 2; string plan_expires_at = 3; bool token_stale = 4;
}
```

## Data

No schema change.

## Events

`household.plan_changed` through `changePlan` (when the plan changes).

## Errors

| Case | gRPC status |
|---|---|
| Missing reason, bad expiry | `INVALID_ARGUMENT` |
| Unknown household | `NOT_FOUND` |

## Acceptance criteria

- [ ] Setting Plus until next month: `GetEntitlements` shows Plus, and IDN-05-03 turns it back to Free after the date.
- [ ] The plan history row names the admin and the reason.
