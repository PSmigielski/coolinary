# ADM-03-02 Nutrition match approval

| | |
|---|---|
| Service | admin |
| Epic | [ADM-03 Catalogue, matching and nutrition review](README.md) |
| Release | R1 |
| Size | M |
| Depends on | ADM-02-03; KIT-06-06 (`ListNutritionSuggestions`, `ApproveNutritionLink`, `RejectNutritionLink`, `SetNutritionLink`, `SearchNutritionEntries`, `CreateAdminEntry`, `NutritionCoverage`) |
| Unblocks | — |

## Goal

Each catalogue ingredient is tied to one nutrition-table entry (Ciqual, BLS, USDA) by a human, after the AI suggests
one (features.md §2 "Calories and goals": "checked by hand before use"; §5 rule 1). Calories are only as good as
these matches.

## Business rules

1. The queue is `ListNutritionSuggestions` (KIT-06-06 rule 1): ingredients with an AI suggestion waiting for
   approval. Each item shows the ingredient (names in pl, de, en; category), the suggested entry with all values,
   the other candidates, confidence and warnings, each entry with source, name and per-100 g energy (kcal),
   protein, carbohydrate, sugars, fat, saturates, fibre and salt. Ingredients parked with no suggestion ("needs
   manual match", KIT-06-05) come from `KitchenAdminService.ListIngredients(missing_nutrition=true)` (KIT-02-10).
2. Actions, each one kitchen call:
   - `approve` the suggestion → `ApproveNutritionLink(ingredient_id, entry_id)`;
   - `choose` another entry, found with `SearchNutritionEntries(query, source)` across the tables →
     `SetNutritionLink(ingredient_id, entry_id)`;
   - `reject` the suggestion and ask the AI again → `RejectNutritionLink(ingredient_id, entry_id)`;
   - `no_entry` (no table has it; dishes using it show "calories incomplete", features.md §5 rule 6) →
     `RejectNutritionLink(..., requeue = false)`;
   - `new_entry`: values a dietitian provides → `CreateAdminEntry(values, name)`, then `SetNutritionLink`;
   - `reopen` an approved link: the same `choose` or `no_entry` on an ingredient that already has one.
3. A note (≤ 300 characters) is required for `choose`, `no_entry` and `new_entry`. Kitchen doesn't store notes; the
   note is kept in the admin audit row only.
4. Approving or changing a match makes kitchen publish `catalogue.nutrition_changed`, and planner recalculates the
   affected dishes (KIT-06, PLN-04). admin shows "affected dishes will be recalculated".
5. Order as kitchen returns it: most used ingredients first. Filters: has warning, source of suggestion,
   ingredient category.
6. Audited (`kitchen.nutrition_match.resolve`, ingredient id, entry id, action).
7. Role `admin` only.
8. A coverage line above the queue: `NutritionCoverage()` for all active ingredients (count and share with an
   approved link). Coverage for only the ingredients of published dishes is not shown: no planner RPC lists them.

## Contract

| Method | Path | Body / response |
|---|---|---|
| GET | `/admin/v1/kitchen/nutrition-reviews?warning=&source=&category=&page_token=` | `{ items, next_page_token, coverage }` |
| GET | `/admin/v1/kitchen/nutrition-entries?q=&source=ciqual\|bls\|usda\|admin` | `{ items: [{ entry_id, source, name, per_100g }] }` (max 20) |
| POST | `/admin/v1/kitchen/nutrition-reviews/{ingredientId}/resolve` | `{ action: approve\|choose\|reject\|no_entry\|new_entry, entry_id?, entry?, note? }` → `{ ingredient_id, entry_id, status }` |

## Data

None in admin.

## Errors

| Case | HTTP |
|---|---|
| Missing note for `choose` / `no_entry` | 400 |
| Entry id unknown | 404 |

## Acceptance criteria

- [ ] Approving the suggestion for "basmati rice, dry" stores it, and a dish using rice gets a new calorie value.
- [ ] `no_entry` for an ingredient makes a dish that uses it show "calories incomplete" in planner.
