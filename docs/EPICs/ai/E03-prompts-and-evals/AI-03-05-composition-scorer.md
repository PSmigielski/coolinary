# AI-03-05 Composition scorer

| | |
|---|---|
| Service | ai |
| Epic | [AI-03 Prompts, golden sets and evals](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | AI-03-03 |
| Unblocks | AI-10-06, AI-11-06, AI-12-04, AI-14-02, AI-15-04, AI-17-02 |

## Goal

Score composed plans the way the planner's rule checker (PLN-07) will judge them, so plan quality is measured in the
ai service before planner exists, and the features.md verification checks 2, 3 and 4 can be run on AI output alone.

## Business rules

1. Eval-only code in `ai/evals/_scoring/compose.py`. It is not used at runtime; PLN-07 stays the authority.
2. Fixture catalogue `ai/evals/_fixtures/catalogue.json`: about 400 ingredients with id, names (pl, de, en), EU
   allergen codes, vegan/vegetarian flags, gluten/lactose flags and a typical-unit hint; tool codes and cuisine codes
   as defined in KIT-10.
3. Per composed meal it checks:
   - **Tools**: every tool the dish needs is in the request's tools (verification 2: with the "microwave + kettle"
     preset, 50 composed dishes need no other tool).
   - **Allergens and diet**: no ingredient carries an excluded allergen or breaks a requested diet (verification 3:
     "no gluten, no nuts" → 0 violations in 50 dishes).
   - **Stock**: every `ingredient_id` is in stock, staples or the allowed purchase list; total use per ingredient
     across the day ≤ stock amount in the same unit.
   - **Cuisine**: the dish's cuisine is one of the meal's requested cuisines, or the response used the cuisine-choice
     path with alternatives (verification 4).
   - **Time**: `total_minutes` ≤ the meal's limit.
   - **Expiry use**: share of stock items with `days_to_expiry ≤ 2` that are used somewhere in the day (target in
     AI-10-06).
   - **Calorie guard** (AI-01-06 rule 2) on all text.
4. Output per case: violations by rule, pass/fail, expiry-use share; totals as the thresholds in AI-03-04 need.

## Contract

No gRPC change.

## Data

No schema change.

## Errors

None at runtime; this task is eval or tooling work. A failed threshold blocks the release of the job (AI-03-04).

## Acceptance criteria

- [ ] Hand-made bad plans (one per rule) are each caught with the right rule code.
- [ ] A hand-made good plan passes all checks.
