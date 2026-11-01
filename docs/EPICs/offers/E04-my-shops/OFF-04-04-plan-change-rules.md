# OFF-04-04 My shops on plan changes

| | |
|---|---|
| Service | offers |
| Epic | [OFF-04 My shops and shopping preferences](README.md) |
| Release | R2 |
| Size | S |
| Status | todo |
| Depends on | OFF-04-01, OFF-01-03 (`PlanFor`), OFF-01-05 (job runner) |
| Unblocks | — |

## Goal

Moving from Plus to Free keeps the 2 nearest shops working and remembers the rest, so upgrading again restores the
user's full list.

## Business rules

1. Decided (Q004): offers learns the plan from identity (`PlanFor`, OFF-01-03), not from events. The rules below run
   (a) at the start of every request that reads or writes the household's my shops (`ListMyShops`, `SetMyShops`,
   `BrowseOffers`), and (b) in a daily job at 03:30 UTC over households that have more than 2 `active` entries or any
   `inactive` entry, calling `PlanFor` once per household.
2. Plan is Free and more than 2 entries are `active`: the first 2 `active` entries by position stay `active`; every
   other `active` entry becomes `inactive` with `inactive_since = now()`. Edition demand is recomputed.
3. Plan is Plus and `inactive` entries exist: they become `active` again after the existing active ones, keeping
   their old relative order, as long as their shop is still active and the 40-shop limit is not passed. Closed
   entries stay closed.
4. Decided (Q168): inactive entries older than 90 days are deleted by the daily job.
5. The household's private user-flyer offers already published stay until they expire (OFF-10-04); new uploads need
   Plus again.
6. No notification is sent from offers; the app reads the new state through `ListMyShops`.
7. Rules 2–3 are idempotent: running them twice in a row changes nothing the second time. A request and the job
   working on the same household serialise on a row lock on the household's my-shops rows.

## Contract

No contract change.

## Data

No schema change.

## Errors

| Case | Behaviour |
|---|---|
| identity unavailable in a request | The token's `x-plan` decides, Free when missing (OFF-01-03 rule 4) |
| identity unavailable in the job | Household skipped, retried in the next run |

## Acceptance criteria

- [ ] A Plus household with 6 shops, downgraded in identity, calls `ListMyShops`: shops 1 and 2 are active, 3–6
      inactive.
- [ ] The same downgrade with no request from that household is applied by the next daily job run.
- [ ] Upgrading 10 days later restores all 6 in the original order on the next request.
- [ ] After 91 days on Free the 4 inactive entries are gone.
- [ ] Running the job twice in a row changes nothing the second time.
