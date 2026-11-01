# ADM-06-04 Weekly spot-check of 20 offers per chain

| | |
|---|---|
| Service | admin |
| Epic | [ADM-06 Flyer sources, offers and shops review](README.md) |
| Release | R2 |
| Size | S |
| Depends on | ADM-06-03; OFF-08-03 (weekly samples), OFF-08-01 (`ListReviewSamples`, `GetReviewSample`), OFF-08-04 (`GetChainAccuracy`) |
| Unblocks | — |

## Goal

Each week the founder checks 20 offers per live chain (features.md §3 step 5), so a slow drop in extraction quality
is caught.

## Business rules

1. offers creates a `weekly` sample of 20 offers for every `live` chain each Monday at 07:00 Europe/Berlin (all of
   them when the chain has fewer than 20 in the window; OFF-08-03 rule 1). admin lists open weekly samples with
   `ListReviewSamples(kind = weekly, status = open)`, showing chain, created on and progress.
2. Reviewing uses the same item screen, verdict rules and correctness rule as ADM-06-03.
3. Result: offers scores the closed sample and alerts the founder when it is below 97.0 (OFF-08-03 rule 4).
   Decided (Q158): bad weeks only alert; a chain never moves back to `syncing` by itself. The founder does it with
   ADM-06-02 `back_to_syncing` when needed.
4. A sample still open on Sunday 18:00 makes offers send an overdue alert; on the next Monday offers marks it
   `abandoned` and counts it as missed (OFF-08-03 rule 3). The board (ADM-06-01) shows missed weeks from
   `GetChainAccuracy`.
5. Audited as ADM-06-03.
6. Role `admin` only.

## Contract

| Method | Path | Response |
|---|---|---|
| GET | `/admin/v1/offers/samples?kind=weekly&status=open` | `{ items: [{ id, chain, created_on, done, total }] }` |

Item and verdict routes as ADM-06-03.

## Data

None in admin.

## Errors

As ADM-06-03.

## Acceptance criteria

- [ ] With 3 live chains, Monday shows 3 open weekly samples of 20 items.
- [ ] A closed weekly sample with 18 of 20 correct (score 90.0) appears in the chain's accuracy record and offers
      sends one founder alert.
