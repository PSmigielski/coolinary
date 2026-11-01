# OFF-08-03 Weekly spot-checks

| | |
|---|---|
| Service | offers |
| Epic | [OFF-08 Accuracy gate and weekly spot-checks](README.md) |
| Release | R2 |
| Size | S |
| Depends on | OFF-08-01, OFF-01-05 |
| Unblocks | OFF-08-04 |

## Goal

Every live chain keeps being checked: 20 offers a week (features.md §3 step 5: "The founder spot-checks 20 offers per
chain each week").

## Business rules

1. Every Monday at 07:00 `Europe/Berlin`, a `weekly` sample of 20 offers is created for each `live` chain (OFF-08-01
   rule 2 window: versions published in the last 7 days). Chains with fewer than 20 offers in the window get a sample
   of all of them.
2. One reminder alert per Monday lists the open samples (`offers.spotcheck.due`).
3. A sample still open on Sunday 18:00 sends an overdue alert; on the next Monday it is `abandoned` and counted as
   missed in the accuracy record.
4. After closing:
   - score < 97.0: alert `offers.accuracy.<chain_id>` with the wrong items;
   - two weekly scores below 97.0 in a row, or one below 90.0 (Decided, Q158): a critical alert
     `offers.accuracy.<chain_id>.critical` with both scores and the wrong items. Nothing changes automatically: the
     chain stays `live` until the founder takes it back with `SetChainSyncing` (OFF-08-02 rule 5).
5. Chains with coverage `feed` are spot-checked the same way (the feed can also be wrong).

## Contract

No new RPC (samples use OFF-08-01 RPCs; the job appears in `ListJobs`).

## Data

No schema change.

## Errors

| Case | Behaviour |
|---|---|
| No offers in the window for a live chain | No sample; info alert (the chain may be stale) |

## Acceptance criteria

- [ ] On a Monday with 3 live chains, 3 samples of 20 are created and one reminder is sent.
- [ ] Two closed weekly samples at 95.0 and 96.0 send one critical alert; the chain is still `live` afterwards.
- [ ] One closed weekly sample at 89.0 sends the critical alert at once; no state change.
