# KIT-06-08 Calculation rules (library)

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-06 Nutrition](README.md) |
| Release | R1 |
| Size | M |
| Depends on | KIT-06-01, KIT-02-02, KIT-06-07 |
| Unblocks | KIT-06-09 |

## Goal

One pure, tested function turns recipe lines and a portion count into nutrients per dish and per portion, exactly as
features.md §5 describes, and refuses to give a number when data is missing.

## Business rules

1. Input lines: `ingredient_id`, optional `gtin`, `quantity` (KIT-02-02 units and measures), `role` (`normal` or
   `frying_oil`), optional `counted_g` (for frying oil), `optional` flag. Plus `portions` (> 0).
2. Per line, the nutrition source is: plausible complete pack values of `gtin` (KIT-06-07) → the ingredient's
   approved table entry → none (`missing: no_nutrition`).
3. Grams from KIT-02-02 `ToGrams`; `no_weight` → `missing: no_weight`; `not_counted` (to taste) → 0 g, not missing.
4. Ingredients are counted as weighed in the recipe (dry pasta, raw meat); the approved entry is already the raw or
   dry one (KIT-06-05 rule 4), so no cooking factor is applied.
5. Decided (Q084): deep frying (features.md §5 step 3, "only the oil the food absorbs is counted"). A
   `frying_oil` line counts `counted_g` only; without `counted_g` the line is `missing: no_absorbed_amount`, so the
   dish is "calories incomplete" (planner then refuses to publish it, PLN-04-01). No default absorption is assumed.
6. Lines with `optional = true` (e.g. "to serve") are left out of the totals and listed as `not_counted`.
7. Totals: for each nutrient, Σ grams / 100 × value. Per portion = total / portions.
8. Complete = no line missing. If not complete, no energy or nutrient numbers are returned at all, only the missing
   lines (features.md §5 step 6, "rather than a number that is too low").
9. A nutrient other than energy that is unknown on any counted line is returned as unknown for the dish (not as a
   lower number); energy stays.
10. Output per portion, two forms:
    - exact (unrounded) values, for planner's portion arithmetic (PLN-12);
    - display values: energy rounded to the nearest 10 kcal, halves up (515 → 520) (features.md §5 step 5).
      Decided (Q085), EU guidance on nutrition labelling rounding: kJ to the nearest 10; protein,
      carbohydrate, sugars, fat, fibre: ≥ 10 g → 1 g, 0.5–10 g → 0.1 g, < 0.5 g → "< 0.5 g"; saturates: ≥ 10 g → 1 g,
      0.1–10 g → 0.1 g, < 0.1 g → "< 0.1 g"; salt: ≥ 1 g → 0.1 g, 0.0125–1 g → 0.01 g, < 0.0125 g → "< 0.01 g".
11. No AI anywhere in this library.

## Contract

No RPC (library `internal/kitchen/nutrition`).

## Data

No schema change.

## Errors

The library returns missing lines, not errors. Bad input (portions ≤ 0, negative quantity) is an error.

## Acceptance criteria

- [ ] Table tests: a 4-portion dish with pasta 400 g dry, olive oil 2 tbsp, onion 1 piece, tomatoes 400 g gives the
      hand-calculated totals from the fixture entries; per portion energy rounds to the nearest 10.
- [ ] One unmatched ingredient → complete false, no numbers, the ingredient listed.
- [ ] Frying oil with `counted_g = 30` counts 30 g regardless of the litre used for frying.
- [ ] Salt "to taste" adds 0 g and does not make the dish incomplete.
- [ ] Changing portions from 4 to 6 scales per-portion values by 4/6 (portion scaling, features.md §2).
