# KIT-01-06 Household plan from identity

| | |
|---|---|
| Service | kitchen (shared platform code; copy to planner and offers) |
| Epic | [KIT-01 Platform](README.md) |
| Release | R1 |
| Size | S |
| Depends on | KIT-01-01, IDN-05 (`GetEntitlements`) |
| Unblocks | KIT-09 (expiry reminders are Plus only), KIT-10-08, KIT-11, KIT-12, KIT-13-06 |

## Goal

Every kitchen rule that depends on Free or Plus, at request time and in background jobs, gets the plan from identity,
the one place that knows it (Decided, Q004). Kitchen keeps no plan of its own.

## Business rules

1. Decided (Q004): kitchen keeps no local copy of any household's plan. Package `internal/platform/plan` exposes
   `PlanOf(ctx, householdID) (Plan, error)`, which calls `IdentityService.GetEntitlements(household_id,
   actor_user_id = "")` (IDN-05-02) with a 300 ms deadline and returns `free` or `plus`.
2. **Memoised for one unit of work only:** within one gRPC request (a value on the request context) or one job run
   (a map owned by the run). Never across requests or runs; nothing is written to Postgres or Redis.
3. **Requests:** every Plus rule in kitchen uses `PlanOf`, not `x-plan` from the JWT (which can be up to 15 minutes
   stale). `x-plan` is used only as a fallback (rule 5) and as a log and metric label.
4. **Jobs** (expiry reminders KIT-09, savings report KIT-12) call `PlanOf` once per candidate household per run.
5. **identity unreachable** (`UNAVAILABLE`, `DEADLINE_EXCEEDED`), decided (Q229):
   1. a request falls back to `x-plan`; when that is missing too, the household is treated as `free`;
   2. a job skips that household for this run and tries again at the next tick;
   3. both cases count `kitchen_plan_lookups_total{outcome="fallback_jwt" | "skipped"}`. Successful lookups count
      `outcome="ok"`.
6. `household.plan_changed` is still consumed, but only to act on a change (KIT-13-06 bumps the planning version;
   KIT-11-04 merges per-person allergens on a downgrade). The event stores nothing about the plan.
7. Configuration: `IDENTITY_GRPC_ADDR=identity:50051`. This adds the call kitchen → identity (no cycle: identity
   calls only notifier).

## Contract

Calls `IdentityService.GetEntitlements` (IDN-05-02). No kitchen proto change.

## Data

No schema change (no plan table).

## Errors

| Case | Behaviour |
|---|---|
| identity down during a request | `x-plan` used, else `free`; metric `fallback_jwt` |
| identity down during a job | household skipped this run; metric `skipped` |
| `NOT_FOUND` from identity (household deleted) | treated as `free`; jobs drop the household |

## Acceptance criteria

- [ ] With identity stubbed to answer `plus`, a Plus-only kitchen RPC succeeds even when `x-plan=free`; stubbed to
      `free`, it returns `plus_required` even when `x-plan=plus`.
- [ ] Two plan checks inside one request make one identity call; two requests make two.
- [ ] With identity stopped, a request uses `x-plan` and counts `fallback_jwt`; the reminder job skips the household,
      counts `skipped`, and sends at the next tick once identity is back.
- [ ] No table stores a household's plan (`\dt` in kitchen's database).
