# OFF-09-05 Find offers for ingredients (planner)

| | |
|---|---|
| Service | offers |
| Epic | [OFF-09 Offer browse, search and lookup by ingredient](README.md) |
| Release | R2 |
| Size | M |
| Status | todo |
| Depends on | OFF-07-07, OFF-04-01 |
| Unblocks | PLN-17 (offer badges), PLN-18 (plan from deals) |

## Goal

planner can ask, for a set of ingredients, which offers exist at the household's shops on given days, to put badges on
shopping-list lines ("−30% at Lidl until Saturday") and to compose dishes from deals (features.md §2 "Shopping list"
step 5, "Plan from deals").

## Business rules

1. Input: household id, ingredient ids (at most 300), shop ids nearest first (at most 40; empty = active my shops),
   and a date range (default: today to today). Same Free/Plus shop rule as OFF-09-01 rule 2.
2. Returns, per ingredient, the visible offers (OFF-07-07 rule 3, including the household's user-flyer offers) whose
   validity overlaps the range at any of the shops, each with `shop_ids`, price, regular price, discount, unit price,
   pack, structured conditions and dates.
3. Per ingredient, offers are ordered by unit price ascending, then by the position of their first shop.
4. Ingredients without offers are returned with an empty list (so planner can clear old badges).
5. Loyalty-only offers are included and flagged; planner decides how to show them.
6. p95 under 100 ms for 150 ingredients and 10 shops.

## Contract

```proto
rpc FindOffersForIngredients(FindOffersForIngredientsRequest) returns (FindOffersForIngredientsResponse);

message FindOffersForIngredientsRequest {
  string household_id = 1;
  repeated string ingredient_ids = 2;
  repeated string shop_ids = 3;
  string from = 4;               // YYYY-MM-DD
  string to = 5;
}
message FindOffersForIngredientsResponse {
  repeated IngredientOffers results = 1;
}
message IngredientOffers {
  string ingredient_id = 1;
  repeated Offer offers = 2;
}
```

## Data

No schema change (index `offers_ingredient` from OFF-07-01).

## Errors

Reason codes (in backticks) travel as trailing gRPC metadata `x-error-code` (README decision 18); version
conflicts use `ABORTED`.

| Case | gRPC status |
|---|---|
| More than 300 ingredients or 40 shops | `INVALID_ARGUMENT` |
| Free and a shop outside active my shops | `FAILED_PRECONDITION` (`plus_required`) |

## Acceptance criteria

- [ ] For a test list of 5 ingredients, 3 with offers at Lidl, the response holds 5 results and 3 non-empty lists
      with shop, price and end date (Verification 6, offers side).
- [ ] Two offers for rice (1 kg at 1.79 and 5 kg at 7.49) are ordered by unit price.
