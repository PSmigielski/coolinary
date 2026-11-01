# OFF-09-04 Offer categories with counts

| | |
|---|---|
| Service | offers |
| Epic | [OFF-09 Offer browse, search and lookup by ingredient](README.md) |
| Release | R2 |
| Size | S |
| Depends on | OFF-09-01 |
| Unblocks | GW-10 |

## Goal

The offers screen shows category chips with how many offers each has in the user's shops, so filtering never leads to
an empty page.

## Business rules

1. `ListOfferCategories` takes the same shop ids and view as `BrowseOffers` (same Free/Plus rule) and returns each
   category id with its count, largest first; categories with 0 offers are left out.
2. Category ids are the catalogue's ingredient categories (OFF-07-04 rule 4). Labels in the user's language come from
   kitchen's catalogue data that the app already loads (GW-05); offers returns ids only.
3. Same cache rule as OFF-09-01 rule 8.

## Contract

```proto
rpc ListOfferCategories(ListOfferCategoriesRequest) returns (ListOfferCategoriesResponse);
message ListOfferCategoriesRequest { string household_id = 1; repeated string shop_ids = 2; OfferView view = 3; }
message ListOfferCategoriesResponse { repeated CategoryCount categories = 1; }
message CategoryCount { string category = 1; int32 count = 2; }
```

## Data

No schema change.

## Errors

Reason codes (in backticks) travel as trailing gRPC metadata `x-error-code` (README decision 18); version
conflicts use `ABORTED`.

| Case | Behaviour |
|---|---|
| Free and a shop outside active my shops | `FAILED_PRECONDITION` (`plus_required`) |

## Acceptance criteria

- [ ] Counts equal the number of offers `BrowseOffers` returns for each category with the same inputs.
