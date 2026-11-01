# AI-09-03 draft_dishes

| | |
|---|---|
| Service | ai |
| Epic | [AI-09 Catalogue and library batch jobs](README.md) |
| Release | R1 |
| Size | L |
| Depends on | AI-04-03, AI-10-01 |
| Unblocks | PLN-03 (AI drafting into the review queue) |

## Goal

Draft library dishes from a brief (cuisine, meal, tools, diet, time), using only catalogue ingredients, so the
starter library of about 200 dishes in 5 cuisines and later growth cost a few euros and land in the same review
queue as hand-written dishes (decision 10).

## Business rules

1. Batch kind `draft_dishes`, **Sonnet 5.5, effort `high`** (composition, once for everyone; at batch prices a
   200-dish run costs about $3–4). Prompt `draft_dishes/v1`. No fallbacks in batches; a refusal is an item error.
2. Item input (`DishBrief`): `brief_ref`; `cuisine` (`coolinary.common.v1.Cuisine`, e.g. `CUISINE_POLISH`, `CUISINE_THAI`);
   `meals` (1–3 `MealSlot`); `tools` allowed (1–16); `diets` required (`VEGETARIAN`, `VEGAN`, `GLUTEN_FREE`,
   `LACTOSE_FREE`, `NUT_FREE`); `max_total_minutes` (10–240); `servings` (default 4); `effort_wanted` (optional
   effort tags such as `NO_COOKING`, `BASE_PLUS_TOPPINGS`); `allowed_ingredients` (≤ 400: id, English name,
   allergen codes, diet flags); `avoid_names` (≤ 200 existing dish names in that cuisine); `style_notes` (≤ 300
   chars from the founder, wrapped as untrusted text, AI-01-07).
3. The allowed-ingredient list for a cuisine is sent as a cached context block (AI-01-05; 1-hour TTL), so a
   200-brief run pays for it about once per cuisine.
4. Output `DishDraft`: a `NewDish` (AI-10-01) written in English, plus `original_name` (the dish's own name where it
   is a proper name, e.g. "Pierogi ruskie", "Käsespätzle"; else empty), `missing_ingredients` (0–5 items the dish
   needs that are not in the allowed list: English name, category, amount), and `confidence`.
5. Dish rules (prompt; checked as listed):
   - Amounts in g, ml or pieces for `servings` portions; `display_hint` like "2 tbsp" allowed.
   - Steps 1–25, each ≤ 400 chars, imperative, with `timer_seconds` where a step waits (simmer, bake, rest).
   - `total_minutes` ≤ `max_total_minutes`; `active_minutes` ≤ `total_minutes`.
   - Uses only the allowed tools; meets every required diet through its ingredients.
   - Effort tags are suggestions; suggest `BASE_PLUS_TOPPINGS` with ingredient roles `BASE` / `TOPPING` when the
     dish works that way (tacos, bowls, wraps; features.md §2 "One dish per meal").
   - No calorie or health claims, no "healthy" / "light" wording (features.md §2 "Calories and goals").
6. Checks in ai: calorie guard and claim words (`healthy`, `light`, `zdrowy`, `lekki`, `gesund`, `leicht` used as a
   claim) in name and steps; minutes ranges; step and amount ranges; tool codes in the global vocabulary; servings
   equal to the brief.
7. **Caller checks** (PLN-03): every `ingredient_id` is in that brief's allowed list; tools ⊆ brief tools; diets hold
   using catalogue flags; name not in `avoid_names`; calories computed by KIT-06 (a dish with unmatched ingredients
   is "calories incomplete" and can't be published); the draft then waits for admin review and translation.

## Contract

```proto
message DishBrief {
  string brief_ref = 1;
  coolinary.common.v1.Cuisine cuisine = 2;
  repeated MealSlot meals = 3;
  repeated coolinary.common.v1.KitchenTool tools = 4;
  repeated coolinary.common.v1.Diet diets = 5;
  int32 max_total_minutes = 6;
  int32 servings = 7;
  repeated coolinary.common.v1.EffortTag effort_wanted = 8;
  repeated IngredientFact allowed_ingredients = 9;   // AI-10-01
  repeated string avoid_names = 10;
  string style_notes = 11;
}

message MissingIngredient {
  string name_en = 1;
  string category = 2;
  coolinary.common.v1.Quantity amount = 3;
}

message DishDraft {
  string brief_ref = 1;
  NewDish dish = 2;                     // AI-10-01, language "en"
  string original_name = 3;
  repeated MissingIngredient missing_ingredients = 4;
  float confidence = 5;
}
```

## Data

No schema change.

## Errors

Per item, as AI-04-03: `ai_refused` (not retryable), `ai_truncated`, `ai_invalid_output`, `ai_item_errored`, `ai_canceled`, `ai_expired` (retryable). Whole-call errors as AI-04-01.

## Acceptance criteria

- [ ] A brief "polish, DINNER, tools hob+oven, VEGETARIAN, ≤ 60 min" yields a draft using only allowed ids, hob/oven only, no meat ingredients, steps with timers.
- [ ] A brief with the "student room" tools (microwave, kettle) yields dishes needing no other tool (verification 2 style, checked by the composition scorer AI-03-05).
- [ ] A planted "healthy" in a fake model answer is caught as an invalid item.
- [ ] Golden set thresholds pass (AI-09-06).

## Notes

- Decided (Q121): dishes are drafted in English, `original_name` keeps the proper dish name (e.g. "Pierogi ruskie"),
  then `translate_dishes` (AI-09-04) produces pl and de.
