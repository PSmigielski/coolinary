# KIT-03-05 Bulk matching through Message Batches

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-03 Name matching](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | KIT-03-04, KIT-01-08 |
| Unblocks | OFF-07 (weekly flyer offers), PLN-03 (recipe import), KIT-08 receipts (R3) |

## Goal

Thousands of names a week (flyer offers, imports) are matched at batch price instead of one by one.

## Business rules

1. In `BULK` mode, names that reach KIT-03-03 step 7 are stored as `queued` and answered `queued` at once.
2. A job every 15 minutes: if queued rows exist and no `pick_catalogue_matches` batch is open (KIT-01-08), take up to
   5,000 queued rows, highest `hit_count` first, build candidates as in KIT-03-04 rule 2, and submit one batch.
3. Results go through the same checks and bands as KIT-03-04 rules 4–5. Rows decided meanwhile by an admin are left
   alone.
4. Failed items stay `queued` with `attempts + 1`; after 3 attempts they become `none` with `review = needed`.
5. A queued row older than 48 hours raises `SendAlert` `kitchen.match_queue_stale` (once per 6 hours).
6. Decided (Q065): callers learn about results by asking again: offers calls `MatchNames` in `BULK` mode every
   30 minutes for offers it holds back. There is no `catalogue.names_matched` event.

## Contract

No change to `MatchNames`. Uses ai batch kind `pick_catalogue_matches`.

## Data

No schema change (uses `name_matches.outcome = 'queued'`).

## Errors

None returned to callers.

## Acceptance criteria

- [ ] 2,000 bulk names with no stored match are answered `queued`, submitted in one batch, and answered from the
      store after the batch ends (stub ai).
- [ ] An admin decision made while the batch runs is not overwritten.
- [ ] Gauge `kitchen_match_queue_size`.
