# OFF-04-03 Store-edition demand

| | |
|---|---|
| Service | offers |
| Epic | [OFF-04 My shops and shopping preferences](README.md) |
| Release | R2 |
| Size | S |
| Status | todo |
| Depends on | OFF-04-01, OFF-03-04 |
| Unblocks | OFF-05-07 |

## Goal

A Kaufland store edition (and any store or format edition) is read only when at least one household uses that store
(features.md §3 "Store editions on demand"), which keeps the fixed AI cost tied to real users.

## Business rules

1. An edition's `demand` = the number of households with an `active` my-shops entry for a shop that maps to that
   edition (`EditionsForShop`, OFF-03-04).
2. Demand is recomputed for the affected editions in the same transaction as every my-shops change, plan-change
   adjustment, shop merge and store link change.
3. features.md says "at least one user has that store in range"; offers cannot know ranges, so "has the store in my
   shops" is used instead.
4. An edition whose demand goes from 0 to 1 gets `wanted_since = now()` and an immediate sync request (OFF-05-07),
   so the household sees offers within the next sync run (at most 15 minutes later) instead of the next day.
5. An edition whose demand drops to 0 keeps being read for 30 more days (`unwanted_since`), so a user who removes and
   re-adds a shop does not lose offers. Decided (Q167): 30 days.
6. National editions of `live` and `syncing` chains are always read, whatever their demand.
7. Metric `offers_edition_demand{chain}` = number of editions with demand > 0 per chain.

## Contract

No public RPC. Admin `ListEditions` (OFF-05-08) shows demand.

## Data

`offers/migrations/000NN_edition_demand.sql`

```sql
ALTER TABLE editions
    ADD COLUMN wanted_since   timestamptz,
    ADD COLUMN unwanted_since timestamptz;
```

## Errors

| Case | Behaviour |
|---|---|
| Recompute fails | The my-shops change rolls back (`INTERNAL`) |

## Acceptance criteria

- [ ] A household adding a Kaufland shop linked to store 1234 sets that edition's demand to 1 and queues a sync.
- [ ] Removing it leaves the edition read for 30 days, then not.
- [ ] Two households with the same store give demand 2; one leaving gives 1.
