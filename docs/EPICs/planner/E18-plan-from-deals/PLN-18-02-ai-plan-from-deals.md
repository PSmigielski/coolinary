# PLN-18-02 Plan from deals with AI

| | |
|---|---|
| Service | planner |
| Epic | [PLN-18 Plan from deals](README.md) |
| Release | R2 |
| Size | L |
| Depends on | PLN-18-01, AI-14 (ComposeFromDeals), PLN-07, PLN-08, PLN-12-06 |
| Unblocks | PLN-18-04, GW-10 |

## Goal

"Plan from deals": a day plan from the fridge plus this week's offers, preferring the nearest shops.

## Business rules

1. Plus only (`plus_required`). Uses the family-day shape (PLN-13-07): one dish per meal for all eaters, portions per
   person.
2. Request to `AI-14 (ComposeFromDeals)`: stock oldest first, deal candidates (PLN-18-01) with shop position but
   no shop address, eaters as anonymous keys with kcal budgets and limits, tools, cuisines, time per meal, up to
   30 library candidates per meal. No names, ages, weights.
3. AI returns per meal a dish (library id or new dish) and the list of offer ids it relies on.
4. Checks without AI:
   - PLN-07 rule check per dish (limits, tools, time);
   - **fridge first:** every meal uses at least one stock ingredient when stock has any that passes the limits;
     across the day, at least 50% of the ingredients by weight come from stock or staples when stock allows it;
   - every bought ingredient is either covered by a referenced offer or is a non-offer purchase; at most 3 non-offer
     purchases per day;
   - offer ids must be from the candidate set;
   - portions via PLN-12 (UNSUITABLE fails).
5. One retry with violations; then the failing meals come from PLN-18-03 (`partial_fallback`); `ai` down, over
   budget or refusing → the whole plan from PLN-18-03.
6. Saved as a day plan with `source = deals`; each meal stores its offer items (offer id, ingredient id, shop id).
7. Response: the plan (PLN-13-07 shape, visibility-filtered), plus per meal the offer items with shop and price, and
   the day's offer cost estimate ("about €12.40 from offers"), labelled "about".

## Contract

`contracts/coolinary/planner/v1/deals.proto`

```proto
rpc CreatePlanFromDeals(CreatePlanFromDealsRequest) returns (DealsPlan);
message CreatePlanFromDealsRequest { string household_id = 1; string viewer_member_id = 2; string date = 3; bool replace = 4; bool use_ai = 5; }
message DealsPlan {
  FamilyDayPlan plan = 1;                       // PLN-13-07
  repeated DealItem items = 2;
  coolinary.common.v1.Money offer_cost_estimate = 3;
}
message DealItem { Meal meal = 1; string offer_id = 2; string ingredient_id = 3; string shop_id = 4; coolinary.common.v1.Money price = 5; coolinary.common.v1.Quantity pack = 6; string valid_to = 7; }
```

## Data

`planner/migrations/000NN_plan_meal_offer_items.sql`

```sql
CREATE TABLE plan_meal_offer_items (
    plan_meal_id  uuid NOT NULL,
    offer_id      text NOT NULL,
    ingredient_id text NOT NULL,
    shop_id       text NOT NULL,
    PRIMARY KEY (plan_meal_id, offer_id)
);
```

## Errors

Reason codes (in backticks) travel as trailing gRPC metadata `x-error-code` (README decision 18); version
conflicts use `ABORTED`.

| Case | gRPC status |
|---|---|
| Free household | `FAILED_PRECONDITION` `plus_required` |
| No shops | `FAILED_PRECONDITION` `no_shops` |
| kitchen or offers down | `UNAVAILABLE` |

## Acceptance criteria

- [ ] A fridge with chicken and an offer on basmati rice: dinner uses the chicken from stock and the rice offer.
- [ ] An AI plan buying everything from offers while the fridge has usable items is rejected by the fridge-first
      check.
- [ ] With `ai` stopped, a deals plan still comes back (`ai_used = false`).
- [ ] The AI request contains shop positions, never addresses or coordinates.

## Notes

- Decided (Q208): the fridge-first thresholds (rule 4: each meal uses ≥ 1 stock item, ≥ 50% by weight from stock,
  at most 3 non-offer purchases a day).
