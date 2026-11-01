# PLN-13-05 Second dish only when an allergy or diet forces it

| | |
|---|---|
| Service | planner |
| Epic | [PLN-13 Family day plan](README.md) |
| Release | R2 |
| Size | S |
| Depends on | PLN-13-01 |
| Unblocks | PLN-13-06, PLN-13-07, PLN-14-06 |

## Goal

A second dish for a meal is planned only when no single dish (base-plus-toppings included) can meet everyone's
allergies or diet (features.md "Risks": "a second dish only when an allergy or diet forces it").

## Business rules

1. Trigger: for a meal, no library candidate and no valid AI candidate passes the combined limits of all eaters,
   including base-plus-toppings dishes.
2. Before splitting, dislikes are relaxed: try again with dislikes removed from the combined limits. If that finds
   a dish, it is used and the dish card notes "contains <ingredient> that Ewa dislikes" for the household.
3. Still nothing: split eaters into two groups. Group A = the largest set of eaters a dish can satisfy (greedy:
   drop the narrowest eater from PLN-13-01 first, then the next). Group B = the rest. Each group gets its own dish.
4. Never more than two dishes per meal. If group B still cannot be served, the meal is returned with
   `no_dish_for` and the member ids, and nothing is planned for them.
5. Calories, portions or goals never cause a split.
6. Dish B prefers "no cooking" dishes; a cooked dish B counts as one extra session (PLN-13-03 key 2).
7. The pot and portions are computed per dish for its group.

## Contract

`contracts/coolinary/planner/v1/family.proto`: PLN-08's plan meal gains `repeated MealDish dishes` (one entry
normally, two when forced), each with its eater member ids. Existing single-dish fields stay for R1 clients.

## Data

`planner/migrations/000NN_plan_meal_dishes.sql`

```sql
CREATE TABLE plan_meal_dishes (
    plan_meal_id uuid NOT NULL REFERENCES plan_meals ON DELETE CASCADE,
    position     smallint NOT NULL CHECK (position IN (1, 2)),
    dish_id      uuid NOT NULL,
    eater_keys   text[] NOT NULL,
    forced_by    text NOT NULL DEFAULT '',   -- 'allergy' or 'diet' when position = 2
    PRIMARY KEY (plan_meal_id, position)
);
```

## Errors

| Case | gRPC status / behaviour |
|---|---|
| Group B cannot be served by any dish | not an error: the meal carries `no_dish_for` with those member ids (rule 4) |
| ai down while composing dish B | dish B comes from the library only (PLN-13 fallback); if none fits, rule 4 applies |
| A split would need a third dish | never produced; the remaining eaters go to `no_dish_for` |

## Acceptance criteria

- [ ] A family where Ola has a coeliac-safe requirement and the library has a gluten-free dinner: one dish for all.
- [ ] A household where one person is vegan and another has a soy allergy and no shared dinner exists: two dishes,
      the second one marked `forced_by = diet`.
- [ ] A dislike alone never produces a second dish.
- [ ] Verification 18: "one dish per meal (unless an allergy or diet forces a second)".
