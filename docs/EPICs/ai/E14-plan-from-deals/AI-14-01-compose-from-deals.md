# AI-14-01 ComposeFromDeals

| | |
|---|---|
| Service | ai |
| Epic | [AI-14 Plan from deals](README.md) |
| Release | R2 |
| Size | L |
| Depends on | AI-10-03 |
| Unblocks | AI-14-02; PLN-18 |

## Goal

A day (or selected meals) built fridge-first and completed with what is on offer nearby this week.

## Business rules

1. Job `compose_from_deals`, Sonnet 5.5 at effort `medium`. Prompt `compose_from_deals/v1`.
2. Input: `CompositionContext` (its `buyable_ids` unused here), meals (1–3), and 1–60 `Deal`s: `offer_ref`,
   `ingredient_id` (matched by offers/kitchen), `shop_ref` (opaque), `shop_rank` (1 = nearest; from the phone's order,
   D13), `price`, `pack_size`, `valid_to`, `has_conditions` (loyalty app or multibuy), facts for each deal ingredient
   in `context.ingredients`.
3. Rules in priority order: hard limits; stock first — every stock line expiring within 2 days that fits is used
   before any deal; deals only for what stock and staples can't cover; never buy what is in stock in enough amount;
   at most 8 deal items for the day; prefer deals at lower `shop_rank` and from as few shops as possible (1 shop
   preferred, never more than 2 in one day); deals valid on the plan date (`valid_to ≥ context.date`).
4. Output: `ComposedMeal`s as AI-10-01 with `to_buy` lines that each carry an `offer_ref`, plus `MEAL_REASON_USES_OFFER`
   where a deal is used. No prices, savings or totals in output text; planner computes them.
5. Validation: AI-10-03 plus AI-14-02 rules. Cache 1 hour.

## Contract

```proto
rpc ComposeFromDeals(ComposeFromDealsRequest) returns (ComposeDayPlanResponse);

message Deal {
  string offer_ref = 1;
  string ingredient_id = 2;
  string shop_ref = 3;
  int32 shop_rank = 4;
  coolinary.common.v1.Money price = 5;
  coolinary.common.v1.Quantity pack_size = 6;
  string valid_to = 7;
  bool has_conditions = 8;
}

message ComposeFromDealsRequest {
  CompositionContext context = 1;
  repeated MealRequest meals = 2;
  repeated Deal deals = 3;            // 1-60
  repeated string exclude_dish_ids = 4;
  int32 variation = 5;
  bool bypass_cache = 6;
}
```

`BuyItem` (AI-10-01) gets `string offer_ref = 3;` (empty outside deal plans).

## Data

No schema change.

## Errors

| Case | gRPC status |
|---|---|
| Input limits broken, unknown ids or codes in the request | `INVALID_ARGUMENT`, `ai_bad_request` |
| Output still invalid after repair | `UNAVAILABLE`, `ai_invalid_output` |
| Others | AI-01-04; the caller falls back to its non-AI result |

## Acceptance criteria

- [ ] With chicken on offer at the rank-1 shop and rice in stock, a dinner uses stock rice and the chicken deal.
- [ ] Nothing in stock in enough amount is bought.
- [ ] Deals come from at most 2 shops in a day.
