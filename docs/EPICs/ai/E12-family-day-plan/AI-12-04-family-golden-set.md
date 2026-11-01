# AI-12-04 Golden set: families including the worked example

| | |
|---|---|
| Service | ai |
| Epic | [AI-12 Family day plan](README.md) |
| Release | R2 |
| Size | M |
| Depends on | AI-12-03, AI-03-05 |
| Unblocks | family plans going live |

## Goal

Measure family plans on realistic households before release 2.

## Business rules

1. 35 cases: the features.md worked example; 10 families with one vegetarian or gluten-free member; 8 with an
   allergy that forces a second dish; 8 with very different targets (1,400–3,200 kcal days); 4 with
   `eat_lighter=true`; 4 with eaters away for a meal.
2. The scorer adds a portion check that replays the planner's rule (features.md §2 "How one meal is split" step 3)
   on library dishes: every present eater's portion = target ÷ kcal per portion lands in ½–2 after rounding to
   quarters.
3. Thresholds: hard-limit violations 0; second dish only when forced (100%); portion check passes for library dishes
   ≥ 0.90; base-plus-toppings chosen in ≥ 0.70 of the one-diverging-diet cases; `eat_lighter` extras correct ≥ 0.80.

## Contract

No change.

## Data

`ai/evals/compose_family_day/`.

## Errors

None at runtime; this task is eval or tooling work. A failed threshold blocks the release of the job (AI-03-04).

## Acceptance criteria

- [ ] Eval passes; baseline committed.
