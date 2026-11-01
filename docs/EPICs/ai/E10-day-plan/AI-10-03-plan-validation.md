# AI-10-03 Validation and repair of composed meals

| | |
|---|---|
| Service | ai |
| Epic | [AI-10 Day plan composition and cuisine choices](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | AI-10-02, AI-01-06 |
| Unblocks | AI-10-04, AI-11, AI-12, AI-14, AI-15 (same validator, extended) |

## Goal

Catch every composed meal that breaks a hard rule before it leaves ai, so the planner's rule checker (PLN-07, the
authority) rarely has to reject one and users never see one (features.md rule 2, "AI proposes, rules check").

## Business rules

The validator `validate_composed(request, meals)` is shared by all composition jobs. It returns violations with a
rule code and a path. Rules:

1. `MEALS` — one result per requested meal slot, none extra.
2. `LIBRARY_ID` — `library_dish_id` is one of the candidates, lists the meal slot in its `meals`, and is not
   excluded.
3. `INGREDIENT_ID` — every new-dish ingredient id is in stock, staples or `buyable_ids` and has a fact.
4. `TOOLS` — new dish `tools` ⊆ `rules.tools`; library dish tools ⊆ `rules.tools`.
5. `ALLERGEN` — no ingredient (new dish, or library candidate ingredients) has an excluded allergen.
6. `DIET` — `VEGAN` ⇒ every ingredient `vegan`; `VEGETARIAN` ⇒ every ingredient `vegetarian`; `GLUTEN_FREE` ⇒
   every ingredient `gluten_free`; `LACTOSE_FREE` ⇒ every ingredient `lactose_free`; `NUT_FREE` ⇒ no ingredient with
   allergen `nuts` or `peanuts`. Optional ingredients count too.
7. `DISLIKED` — no disliked ingredient id.
8. `TIME` — `total_minutes` ≤ the meal's limit; `active_minutes` ≤ `total_minutes`.
9. `STOCK_USE` — every `stock_uses` id is a stock line, unit equals the line's unit, amount > 0; per line, the sum
   over all meals ≤ `available`. Every stock-line ingredient of a dish has a `stock_uses` entry.
10. `BUY` — `to_buy` ids ⊆ `buyable_ids`; total `to_buy` lines in the day ≤ `max_items_to_buy` (alternatives
    excepted, AI-10-04).
11. `SERVINGS` — new dish `servings` = request `portions`.
12. `TEXT` — calorie guard and claim words (AI-01-06, AI-09-03 rule 6); name ≤ 80, steps 1–20 (25 for drafts), step
    ≤ 400 chars; name not in `exclude_dish_names` (case-insensitive).
13. `AMOUNT` — amounts within 0 < g/ml ≤ 5,000 per portion × portions, pieces ≤ 50 per portion × portions.

After validation passes, ai sets `check_first = true` on every `stock_uses` entry whose stock line has
`past_date = true` and `false` on all others, overwriting the model's value (Decided, Q103).

Outcome: any violation → one repair call (AI-01-06 rule 4) with the violation list; still failing → `UNAVAILABLE
ai_invalid_output`. Violations are counted per rule in `ai_compose_violations_total{job, rule}`.

## Contract

No change.

## Data

No schema change.

## Errors

| Case | gRPC status |
|---|---|
| Still invalid after repair | `UNAVAILABLE`, `ai_invalid_output` |

## Acceptance criteria

- [ ] One hand-made bad response per rule (13 cases) is caught with the right rule code and path.
- [ ] Two meals using 4 and 3 eggs of a 6-egg stock line fail `STOCK_USE`; 3 and 3 pass.
- [ ] A plan using a `past_date` yoghurt line returns that use with `check_first = true` even when the model's output
      left it false; uses of other lines are `false` (Q103).
- [ ] A vegan request with a new dish containing honey (fact `vegan=false`) fails `DIET`.
- [ ] The repair message lists only rule codes and paths.
