# IDN-05-03 Manual plan expiry

| | |
|---|---|
| Service | identity |
| Epic | [IDN-05 Plan and entitlements](README.md) |
| Release | R1 |
| Size | S |
| Depends on | IDN-05-01 |
| Unblocks | ADM-05 (time-limited Plus for testers and founding plans) |

## Goal

Plus given by hand until a date (testers, founding annual plans sold in the validation phase) ends on that date
without anyone remembering to switch it off.

## Business rules

1. A job runs every 15 minutes under a Postgres advisory lock.
2. It selects households with `plan = plus`, `plan_source = manual` and `plan_expires_at <= now()`, at most 1,000
   per run, and calls `changePlan(to: free, source: default, reason: "expired", changedBy: "system")` for each.
3. A household whose expiry was moved into the future before the job ran is not touched (it no longer matches).
4. The downgrade itself follows IDN-05-04 for shared households.
5. Metric `identity_plans_expired_total`.

## Contract

No RPC.

## Data

Index `households (plan_expires_at) WHERE plan_source = 'manual'`.

## Events

`household.plan_changed` via `changePlan`.

## Errors

| Case | Behaviour |
|---|---|
| Lock held | Skip the run |

## Acceptance criteria

- [ ] A manual Plus with expiry 10 minutes ago becomes Free within 15 minutes, with a history row
      (`changed_by = system`, `reason = expired`) and one event.
- [ ] A Plus without expiry is never touched.
