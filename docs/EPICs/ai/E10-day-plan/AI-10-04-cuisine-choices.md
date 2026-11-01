# AI-10-04 Cuisine fit and the three choices

| | |
|---|---|
| Service | ai |
| Epic | [AI-10 Day plan composition and cuisine choices](README.md) |
| Release | R1 |
| Size | M |
| Depends on | AI-10-03 |
| Unblocks | AI-10-06; PLN-08 |

## Goal

When the fridge doesn't fit the cuisine the user wants, offer three ways forward instead of a bad dish: the closest
dish possible now, the 1–3 items to buy, or a different cuisine for today (features.md §2 "Cuisine choice").

## Business rules

1. Every `ComposedMeal` has `cuisine_fit`:
   - `FITS` — the dish is in one of the meal's cuisines and needs nothing to buy. No alternatives.
   - `CLOSEST_NOW` — no dish in the cuisine is possible from stock alone; the meal's dish is the closest one in that
     cuisine from stock (simplified or partly authentic). Then exactly two alternatives follow.
   - `NOT_POSSIBLE` — no dish meets the hard limits at all (e.g. empty fridge and no staples). No dish, two
     alternatives may still be given (buy items; or none).
2. Alternatives (`MealAlternative`), in this order:
   1. `BUY_ITEMS` — a dish in the requested cuisine that needs 1–3 items from `buyable_ids` (the cuisine's base
      pantry and common items); `to_buy` lists them with amounts.
   2. `OTHER_CUISINE` — a dish fully from stock in another cuisine the household ticked (or any cuisine when none
      other is ticked), with its cuisine code.
3. Alternatives carry the dish (library id or a `NewDish` **without steps**) and their own `stock_uses`. Each
   alternative must be feasible together with the other meals' chosen dishes (STOCK_USE check runs per
   alternative). Steps are generated only when the user picks an alternative: planner calls `ModifyMeal` with mode
   `EXPAND_ALTERNATIVE` (AI-11-01), which keeps output tokens of the first call small.
4. `to_buy` in a `BUY_ITEMS` alternative is outside the day's `max_items_to_buy` limit, but limited to 3 lines.
5. Validation additions to AI-10-03: `CUISINE` — `FITS` and `CLOSEST_NOW` dishes have a cuisine in the meal's list;
   `OTHER_CUISINE` has one not in it; `BUY_ITEMS` has 1–3 `to_buy` lines; alternatives count is 0 for `FITS` and 2
   for `CLOSEST_NOW`.

## Contract

```proto
enum CuisineFitStatus { CUISINE_FIT_STATUS_UNSPECIFIED = 0; CUISINE_FIT_STATUS_FITS = 1;
  CUISINE_FIT_STATUS_CLOSEST_NOW = 2; CUISINE_FIT_STATUS_NOT_POSSIBLE = 3; }

message CuisineFit {
  CuisineFitStatus status = 1;
  string requested_cuisine = 2;   // the cuisine the meal aimed at
}

enum AlternativeKind { ALTERNATIVE_KIND_UNSPECIFIED = 0; ALTERNATIVE_KIND_BUY_ITEMS = 1; ALTERNATIVE_KIND_OTHER_CUISINE = 2; }

message MealAlternative {
  AlternativeKind kind = 1;
  oneof dish {
    string library_dish_id = 2;
    NewDish new_dish = 3;          // steps empty
  }
  repeated StockUse stock_uses = 4;
  repeated BuyItem to_buy = 5;
  coolinary.common.v1.Cuisine cuisine = 6;
}
```

## Data

No schema change.

## Errors

| Case | gRPC status |
|---|---|
| Still invalid after the one repair | `UNAVAILABLE`, `ai_invalid_output` |

## Acceptance criteria

- [ ] Verification 4: for each launch cuisine (Polish, German, Italian, Asian, Mediterranean), a fixed test fridge gives either `FITS` in that cuisine or `CLOSEST_NOW` with a `BUY_ITEMS` alternative naming 1–3 items and an `OTHER_CUISINE` alternative.
- [ ] A `BUY_ITEMS` alternative with 4 items fails validation.
- [ ] Alternatives have no steps; their stock uses plus the other meals' fit the stock.
