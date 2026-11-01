# PLN-19-01 Candidate shops for a trip

| | |
|---|---|
| Service | planner |
| Epic | [PLN-19 Trip price matrix](README.md) |
| Release | R2 |
| Size | S |
| Depends on | OFF-04 (ListMyShops) |
| Unblocks | PLN-19-02, PLN-19-03 |

## Goal

Accept the shops the phone wants priced and check them against the household's shops.

## Business rules

1. The phone sends 1–10 shop ids, nearest first, already filtered by the travel range (features.md "How the plan is
   found" step 1). Planner never receives distances or positions.
2. Every id must be one of the household's shops (`ListMyShops`); unknown ids are dropped and reported in
   `ignored_shop_ids`.
3. If no valid shop remains: `FAILED_PRECONDITION` (`no_shops`).
4. Plus only: Free gets `FAILED_PRECONDITION` (`plus_required`) (features.md §6: "the trip plan across shops").
5. Duplicates are removed, keeping the first position.

## Contract

Part of `GetTripPriceMatrixRequest` (PLN-19-03): `repeated string shop_ids = 3;`.

## Data

No schema change.

## Errors

| Case | gRPC status / behaviour |
|---|---|
| More than 10 shop ids | `INVALID_ARGUMENT` |
| No valid shop left after rule 2 | `FAILED_PRECONDITION`, `x-error-code: no_shops` |
| Free household | `FAILED_PRECONDITION`, `x-error-code: plus_required` |
| offers down (`ListMyShops`) | `UNAVAILABLE`; no matrix without knowing the household's shops |
| Unknown or foreign shop ids | dropped and listed in `ignored_shop_ids` (not an error) |

## Acceptance criteria

- [ ] 12 ids sent → `INVALID_ARGUMENT` (more than 10).
- [ ] One id not in the household's shops → priced without it and listed in `ignored_shop_ids`.

## Notes

- Decided (Q162): at most 10 shops (175 three-shop sets, trivial on a phone).
