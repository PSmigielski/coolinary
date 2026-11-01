# OFF-05-07 Read store editions only on demand

| | |
|---|---|
| Service | offers |
| Epic | [OFF-05 Flyer source framework and sync scheduler](README.md) |
| Release | R2 |
| Size | S |
| Depends on | OFF-05-04, OFF-04-03 |
| Unblocks | OFF-06-03, OFF-06-07 |

## Goal

The fixed AI cost of per-store chains grows with the stores users actually use, not with the number of stores in
the country (features.md §3 step 8, business model "Fixed from the flyer release").

## Business rules

1. A store or format edition is "wanted" when `demand > 0`, or when `unwanted_since` is less than 30 days ago
   (OFF-04-03).
2. Sync checks of per-store and per-format chains cover only wanted editions.
3. When an edition becomes wanted (demand 0 → 1), a one-off check of that edition runs within 15 minutes, outside
   the daily schedule.
4. At most 300 store editions per chain are checked per day (Decided, Q169). Above it, editions are taken in order
   of demand (highest first) and the founder gets an alert `offers.editions.<chain_id>` once a day.
5. Editions that stop being wanted keep their already published offers until those end; nothing new is read.
6. Two store editions whose flyer versions have the same fingerprint (identical pages) share one extraction: the
   offers of the first are copied to the second instead of a second AI run.

## Contract

No contract change.

## Data

No schema change.

## Errors

| Case | Behaviour |
|---|---|
| More than 300 wanted editions for a chain | The rest skipped that day, alert once a day |

## Acceptance criteria

- [ ] With 3 Kaufland DE store editions wanted out of 50 linked stores, a check reads exactly 3 editions.
- [ ] A household adding a fourth store triggers a one-off check of it within 15 minutes.
- [ ] Two stores with identical flyer pages cause one extraction batch, and both editions get the offers.
