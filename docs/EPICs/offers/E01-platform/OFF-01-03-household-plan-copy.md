# OFF-01-03 Household plan lookup through identity

| | |
|---|---|
| Service | offers |
| Epic | [OFF-01 Platform](README.md) |
| Release | R2 |
| Size | S |
| Status | todo |
| Depends on | OFF-01-04, IDN-05 (`GetEntitlements`) |
| Unblocks | OFF-04-01, OFF-04-04, OFF-09-01, OFF-10-01 |

## Goal

offers knows whether a household is Free or Plus by asking identity, the single owner of the plan, instead of keeping
its own copy (Decided, Q004).

## Business rules

1. Decided (Q004): offers keeps no local plan copy and does not consume `household.plan_changed`. Every Free/Plus
   rule in offers asks identity through `PlanFor(ctx, householdID)`, which calls
   `IdentityService.GetEntitlements` (IDN-05) and returns `free` or `plus`.
2. In a request, `PlanFor` is called at most once and its answer is memoised for the rest of that request only;
   nothing is cached across requests. `x-plan` metadata is forwarded for logs and ai tracing, and decides only when
   identity can't be reached (rule 4).
3. Background jobs call `PlanFor` once per household per run (e.g. the my-shops plan rules, OFF-04-04).
4. Decided (Q229): identity unavailable or past its deadline: a request falls back to the `x-plan` metadata from the
   signed-in user's token (at most 15 minutes old), and to Free when that is missing; counted in
   `offers_plan_lookups_total{outcome="fallback_jwt"}`. A job skips that household and retries it in the next run.
5. A household identity doesn't know (`NOT_FOUND`) is treated as deleted: the request fails `NOT_FOUND`; jobs skip it
   (OFF-14-01 erases it on `household.deleted`).
6. Deadline 2 s; one retry on `UNAVAILABLE` after 200 ms (OFF-01-04 rule 4).

## Contract

No gRPC change in offers. Uses `IdentityService.GetEntitlements` (IDN-05). Internal Go function
`PlanFor(ctx, householdID) (Plan, error)`.

## Data

No schema change (no `household_plans` table).

## Errors

| Case | Behaviour |
|---|---|
| identity down or slow in a request | no error: the token's `x-plan` decides, Free when missing (rule 4) |
| identity down during a job | Household skipped, retried next run; `offers_plan_lookup_failures_total` + 1 |
| Household unknown to identity | Request `NOT_FOUND`; job skips it |

## Acceptance criteria

- [ ] A request that checks the plan twice (e.g. browse with shop limits) calls `GetEntitlements` once.
- [ ] Two requests in a row call `GetEntitlements` twice (no cross-request cache).
- [ ] A household upgraded in identity is treated as Plus on the very next request, with no event involved.
- [ ] With identity stopped, `SetMyShops` with `x-plan: plus` accepts more than 2 shops; without `x-plan` it
      applies the Free limit of 2.
- [ ] offers registers no handler for `household.plan_changed`.
