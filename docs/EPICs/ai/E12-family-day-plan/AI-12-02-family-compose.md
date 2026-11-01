# AI-12-02 Family composition prompt and call

| | |
|---|---|
| Service | ai |
| Epic | [AI-12 Family day plan](README.md) |
| Release | R2 |
| Size | L |
| Depends on | AI-12-01, AI-10-02 |
| Unblocks | AI-12-03; PLN-13 |

## Goal

One dish per meal that everyone at the table can eat, sized well for very different appetites, cooked as few times
as possible.

## Business rules

1. Job `compose_family_day`, Sonnet 5.5 at effort `medium`. Prompt `compose_family_day/v1`.
2. All rules of AI-10-02 apply, with these on top (features.md §2 "One dish per meal, few cooking sessions"):
   1. One main dish per meal for all present eaters. It must meet the hard limits of **every** present eater.
   2. When one eater's diet narrows the choice too much (one vegetarian, one gluten-free), prefer a
      base-plus-toppings dish (tacos, bowls, wraps): one shared base, toppings chosen per eater; toppings are no-cook
      or share the base's pan (a topping needing its own pan counts as extra cooking and is avoided).
   3. A second dish only when no single dish (base plus toppings included) can meet every present eater's hard
      limits; it lists the eaters who need it.
   4. Calorie fit: prefer dishes whose standard portion lets every eater's portion land between ½ and 2 portions —
      for library candidates, ratio `meal target ÷ kcal_per_portion` between 0.5 and 2 for every present eater. New
      dishes aim at a standard portion near the median eater's meal target. Planner checks the real numbers.
   5. `eat_lighter`: prefer dishes rich in vegetables and with fewer calories per portion; for breakfast and the
      non-main meal prefer fruit, vegetables, milk and yoghurt over bread, cheese and nuts. Each eater keeps their own
      target.
   6. Meals with `prefer_no_cooking` come from no-cooking dishes, as on single days.
3. New dishes' `servings` = number of present eaters at that meal (standard portions); planner rescales the pot.
4. Expected size ≈ 9,000 input and 4,000 output tokens.

## Contract

As AI-12-01.

## Data

No schema change.

## Errors

| Case | gRPC status |
|---|---|
| Input limits broken, unknown ids or codes in the request | `INVALID_ARGUMENT`, `ai_bad_request` |
| Output still invalid after repair | `UNAVAILABLE`, `ai_invalid_output` |
| Others | AI-01-04; the caller falls back to its non-AI result |

## Acceptance criteria

- [ ] Worked example (Adam, Ewa, Ola, Kuba as e1–e4; dinner the main meal; chicken curry with rice as a 520-kcal candidate): curry is chosen for dinner for all four, no second dish.
- [ ] One vegetarian eater among four: a base-plus-toppings dish with a vegetarian topping set for that eater in ≥ 7 of 10 runs; never a second dish unless a hard limit forces it.
- [ ] `eat_lighter=true` shifts breakfast extras to fruit/yoghurt in the eval (AI-12-04).
