# PLN-13-01 Hard limits for a shared meal

| | |
|---|---|
| Service | planner |
| Epic | [PLN-13 Family day plan](README.md) |
| Release | R2 |
| Size | S |
| Status | todo |
| Depends on | PLN-07, KIT-11 (GetPlanningProfile) |
| Unblocks | PLN-13-02, PLN-13-05, PLN-13-06 |

## Goal

The shared dish must meet the allergens and diet of everyone eating it (features.md "One dish per meal, few
cooking sessions"). This task builds the combined limits per meal and measures how much one person narrows the
choice.

## Business rules

1. Eaters of a meal: PLN-12-01 attendance applied to the snapshot members.
2. Combined limits = union of all eaters' allergens, the strictest diet (vegan ⊃ vegetarian), union of
   gluten-free, lactose-free and nut-free flags, union of dislikes. Household tools and the meal's time limit apply
   to everyone.
3. Free households (one profile): the household profile is the only limit set (features.md §6, per-person diet
   profiles are Plus). After a Plus → Free downgrade, kitchen has merged every person's allergens into that household
   list (Q195), so they keep applying; per-person diets don't.
4. The combined limits are passed to the rule checker (PLN-07) unchanged, so a library dish and an AI dish are
   judged the same way.
5. **Narrowness:** for each eater, count library dishes for this meal that pass everyone's limits (`n_all`) and
   those that pass everyone's limits except this eater's (`n_without`). An eater narrows the choice too much when
   `n_all < 5` or `n_all < 0.3 × n_without`. The narrowest such eater is reported; composers then prefer
   base-plus-toppings dishes (PLN-13-02).
6. Dislikes are hard limits for the shared dish, except that a disliked ingredient used only as a topping the
   person skips is allowed (PLN-13-02).
7. Counts are cached per (household, meal, profile version) in redis-cache for 1 hour.

## Contract

Internal package `planner/internal/family` (`SharedLimits(meal, eaters) Limits`, `Narrowness(...)`). No RPC.

## Data

No schema change.

## Errors

| Case | Behaviour |
|---|---|
| Snapshot unavailable (kitchen down) | composer stops with `UNAVAILABLE`; no plan without limits |

## Acceptance criteria

- [ ] Ewa vegetarian, Adam nut allergy: combined = vegetarian + nut-free.
- [ ] Ola away for dinner: her lactose-free flag does not apply to dinner but applies to supper.
- [ ] With 200 library dishes, 12 vegetarian dinners out of 60 dinners overall: Ewa is reported as narrowing
      (12 < 0.3 × 60).
- [ ] No dish passing the combined limits is ever offered by the composers (Verification 3).

## Notes

- Decided (Q204): the narrowness thresholds in rule 5 (fewer than 5 dishes fit everyone, or under 30% of those
  without that person's limits).
