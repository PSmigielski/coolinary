# PLN-07-03 Allergen, diet and dislike check per eater

| | |
|---|---|
| Service | planner |
| Epic | [PLN-07 Rule checker](README.md) |
| Release | R1 |
| Size | M |
| Depends on | PLN-07-01, PLN-02-03, KIT-02 (GetIngredients flags) |
| Unblocks | PLN-07-05, PLN-13, Verification 3 |

## Goal

The shared dish meets the allergens and diet of everyone eating it (features.md "One dish per meal"), and nobody is
served what they marked as a dislike.

## Business rules

1. Flags per ingredient come from `GetIngredients` (pre-fetched). An ingredient whose flags are unknown counts as
   containing every allergen and fitting no diet (conservative).
2. **Main and base lines** (non-optional): for every eater, any allergen of the eater in the line → `ALLERGEN`
   (member, ingredient, allergen); a line that breaks an eater's diet → `DIET`; a disliked ingredient → `DISLIKE`.
3. **Topping groups**: for each eater, the allowed options are those passing rule 2 for that eater. If an eater has no
   allowed option in a group, → violation for that group (`ALLERGEN`/`DIET`/`DISLIKE` with `topping_group`). Otherwise
   the verdict's `plates` lists, per eater, the allowed options per group (used by cooking mode: "Ola: beans").
4. **Optional lines** breaking an eater's limits → no violation; they go into `leave_out` for that eater ("leave out the
   cheese for Ewa"). If the eater's limit is an allergen, the leave-out is marked `required=true`.
5. Cross-contact ("may contain") is not modelled. Decided (Q106): the app and the terms state that cross-contact
   isn't checked.
6. The same rules apply to library dishes when a plan uses them for a set of eaters (PLN-08-05, PLN-13).

## Contract

Internal; `plates` and `leave_out` are carried into plan meals (PLN-08-01) and cooking mode (PLN-10-01).

## Data

No schema change.

## Errors

None of its own; internal check run by the checker core (PLN-07-01).

## Acceptance criteria

- [ ] Eaters "no gluten" and "no nuts": a pasta dish with wheat pasta fails with `ALLERGEN gluten`.
- [ ] Taco bowl with toppings beef / beans and a vegetarian eater passes, with that eater's plate = beans only.
- [ ] Optional parmesan with a lactose-free eater passes with a leave-out for that eater.
- [ ] An ingredient with unknown flags fails for an eater with any allergen.
