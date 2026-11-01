# KIT-06-09 ComputeNutrition RPC

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-06 Nutrition](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | KIT-06-08 |
| Unblocks | PLN-04, PLN-10 (cooked with scanned products), PLN-12, Verification 16 |

## Goal

Planner (and later other callers) get nutrition for many recipes in one call, with the versions they need to know
when to recalculate.

## Business rules

1. Input: up to 100 items; each item has up to 200 lines (KIT-06-08 rule 1), `portions` between 0.25 and 1,000, and
   an optional `item_key` the caller uses to match answers. `language` sets the names of missing ingredients.
2. Merged ingredient ids resolve to their target (KIT-02-09); unknown ids are missing lines with reason
   `unknown_ingredient`, not an error.
3. Output per item: `complete`, `missing` (line index, ingredient id, display name, reason), `per_portion_display`
   and `per_portion_exact` (only when complete), `total_exact`, `unknown_nutrients`, `sources` (codes of the sources
   used, for credits), and `nutrition_version` = the highest `nutrition_version` among the item's ingredients.
4. Approved values are cached in process per ingredient and refreshed when its `nutrition_version` changes
   (checked against the catalogue version every 30 s).
5. p95 ≤ 30 ms for one item of 20 lines; ≤ 300 ms for 100 items.
6. No household data is needed or stored; `household_id` is not part of the request.

## Contract

`contracts/coolinary/kitchen/v1/nutrition.proto`

```proto
rpc ComputeNutrition(ComputeNutritionRequest) returns (ComputeNutritionResponse);
message ComputeNutritionRequest { repeated NutritionItem items = 1; string language = 2; }
message NutritionItem { string item_key = 1; repeated NutritionLine lines = 2; double portions = 3; }
message NutritionLine { string ingredient_id = 1; string gtin = 2; coolinary.common.v1.Quantity quantity = 3;
  LineRole role = 4; optional double counted_g = 5; bool optional = 6; }
enum LineRole { LINE_ROLE_UNSPECIFIED = 0; LINE_ROLE_NORMAL = 1; LINE_ROLE_FRYING_OIL = 2; }
message NutritionResult { string item_key = 1; bool complete = 2; repeated MissingLine missing = 3;
  Nutrients per_portion_display = 4; Nutrients per_portion_exact = 5; Nutrients total_exact = 6;
  repeated string unknown_nutrients = 7; repeated string sources = 8; int32 nutrition_version = 9;
  repeated DisplayQualifier qualifiers = 10; } // e.g. sugars "< 0.5 g"
message MissingLine { int32 line_index = 1; string ingredient_id = 2; string name = 3; string reason = 4; }
message ComputeNutritionResponse { repeated NutritionResult results = 1; }
```

## Data

No schema change.

## Errors

| Case | Status | Message |
|---|---|---|
| More than 100 items or 200 lines | `INVALID_ARGUMENT` | `at most 100 items of 200 lines` |
| Portions out of range | `INVALID_ARGUMENT` | `portions must be 0.25-1000` |
| Negative or zero quantity | `INVALID_ARGUMENT` | `quantity must be positive` |

## Acceptance criteria

- [ ] Verification 16: for 20 starter dishes, per-portion calories are within 10% of a dietitian's calculation
      (checked in PLN-04 with this RPC); a dish with an unmatched ingredient returns `complete = false` and no number.
- [ ] With ai stopped, results are identical (no AI dependency).
- [ ] Latency targets met in a benchmark.
