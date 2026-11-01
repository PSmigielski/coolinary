# OFF-07-08 Duplicate and conflicting offers

| | |
|---|---|
| Service | offers |
| Epic | [OFF-07 Offer extraction, validation and publishing](README.md) |
| Release | R2 |
| Size | S |
| Status | todo |
| Depends on | OFF-07-03 |
| Unblocks | OFF-07-05 |

## Goal

The same product printed twice (on the cover and inside, or in two overlapping flyers) appears once, and two
different prices for the same product on the same days are caught.

## Business rules

1. Key of an offer: edition, normalised product name (lower-cased, without diacritics and punctuation), brand, pack
   (count, value, unit, loose).
2. Within one version: offers with the same key and the same price are duplicates; the one on the lowest page is kept,
   the others are `rejected` with reason `duplicate`.
3. Across flyers of the same edition with overlapping dates: same key and same price → the offer of the flyer that
   started earlier is kept, the newer one `rejected` (`duplicate`); same key and a different price → the newer one is
   held with `conflicting_price` (a weekend-only price, for example, is legitimate and the founder releases it).
4. Runs after normalisation and before validation; duplicates are not counted as errors in accuracy figures.

## Contract

No gRPC change.

## Data

No schema change.

## Errors

| Case | Behaviour |
|---|---|
| Duplicate or conflict found | Duplicate rejected, conflict held; no caller-visible error |

## Acceptance criteria

- [ ] A product on page 1 and page 7 at the same price leaves one offer from page 1.
- [ ] The same product at 4,99 in Monday's flyer and 3,99 in an overlapping weekend flyer: the second is held with
      `conflicting_price`.
