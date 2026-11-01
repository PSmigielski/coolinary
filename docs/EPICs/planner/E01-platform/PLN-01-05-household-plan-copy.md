# PLN-01-05 Household plan lookup from identity

| | |
|---|---|
| Service | planner |
| Epic | [PLN-01 Platform](README.md) |
| Release | R1 |
| Size | S |
| Status | todo |
| Depends on | PLN-01-02, IDN-05 (GetEntitlements) |
| Unblocks | PLN-09, PLN-11, PLN-12, PLN-13, PLN-14, PLN-17, PLN-18, PLN-19, PLN-20, PLN-21-03 |

## Goal

Every Free/Plus rule in planner uses the household's plan as identity knows it right now. Decided (Q004): planner asks
identity each time instead of keeping a local copy fed by `household.plan_changed`. This adds the call
planner → identity (no cycle: identity calls only notifier).

## Business rules

1. `plans.Get(ctx, household_id)` calls `IdentityService.GetEntitlements(household_id)` (IDN-05-02) and returns the
   plan (`free` / `plus`) and the entitlements planner needs (family goals, week plan, deals, price matrix, top-up).
2. Requests: the first Plus check in a request makes the call; the result is memoised in the request context for the
   rest of that request and never shared between requests. `x-plan` is forwarded for tracing and to ai, but no Plus
   rule in planner is decided from it.
3. Background jobs (PLN-11 promotion candidates, PLN-17 offer-ends reminders): one call per household per job run.
4. Deadline 1 s (PLN-01-02 rule 2). If identity is unavailable:
   - Decided (Q229): a request falls back to the `x-plan` metadata from the signed-in user's token (at most 15 minutes
     old); when that is missing too, the household is treated as Free. Counted in
     `planner_plan_lookup_failures_total{context="request"}`. Requests that work the same on Free and Plus don't call
     identity at all;
   - a background job skips that household for this run and counts it in
     `planner_plan_lookup_failures_total{context="job"}`; the next run tries again.
5. An unknown household (NOT_FOUND from identity) is treated as Free.
6. No table, no event consumption for the plan. `household.plan_changed` is still consumed by PLN-21-03 to recompute
   future meals after a downgrade or upgrade.

## Contract

No planner contract change. Uses `IdentityService.GetEntitlements` (IDN-05-02).

## Data

No schema change.

## Errors

| Case | gRPC status | Message or behaviour |
|---|---|---|
| identity unavailable during a Plus-only request | no error: the token's `x-plan` decides (Q229); Free when missing | none |
| identity unavailable during a job | — | household skipped this run, metric + 1 |
| identity answers NOT_FOUND | — | treated as Free |

## Acceptance criteria

- [ ] A Plus-only RPC calls `GetEntitlements` once, even when it checks Plus three times in one request.
- [ ] Two requests in a row both call identity (no cache between requests).
- [ ] A household upgraded a second ago gets Plus behaviour on its next request, regardless of the `x-plan` header.
- [ ] With identity stopped, `CreateWeekPlan` with `x-plan: plus` succeeds, without `x-plan` it is refused as Free;
      `SuggestDishes` still works.
- [ ] The offer-ends reminder job skips a household whose lookup fails and processes it on the next run.
