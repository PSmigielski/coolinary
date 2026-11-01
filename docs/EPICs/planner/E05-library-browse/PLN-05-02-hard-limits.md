# PLN-05-02 Hard limits from the kitchen profile

| | |
|---|---|
| Service | planner |
| Epic | [PLN-05 Library browse](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | PLN-05-01, PLN-01-03, PLN-02-03 |
| Unblocks | PLN-03-09, PLN-06, PLN-08, Verification 2 and 3 |

## Goal

No dish is ever shown that needs an unticked tool or breaks the allergens, diet or dislikes of the people it is for
(features.md §1: tools, diet and allergens are hard limits).

## Business rules

1. **Tools:** a dish is shown only if its tools ⊆ the household's tools from the profile.
2. **Eaters:** limits are applied for the people in `eater_member_ids`; empty = every person in the profile.
   In R1 (one profile) this is that profile. Decided (Q099): browsing filters for everyone in the household by
   default; the app's "for me" switch sends only the user's member id.
3. **Allergens:** a dish is hidden if any non-optional line (or every option of a topping group) contains an allergen
   of any eater. Dishes with `allergens_complete=false` are hidden when any eater has an allergen.
4. **Diets** (vegetarian, vegan, gluten-free, lactose-free, nut-free): a dish is shown only if every eater's diets are
   in its `diets`, or — for `base_plus_toppings` dishes — in `diets_capable` (each eater can build a plate that fits).
5. **Dislikes** (ingredient ids, hard limit per features.md §1): a dish is hidden if a disliked ingredient is in a
   non-optional `main` or `base` line. In an optional line or in a topping group with another option, the dish stays
   and the ingredient id goes to `leave_out_ingredient_ids`.
6. Optional lines with an eater's allergen or diet conflict also go to `leave_out_ingredient_ids` (the app says "leave
   out the walnuts").
7. Limits are applied in SQL on `library_search` (tools `<@`, allergens `&&`, diets `@>`), not after paging, so pages
   are full and `total` is exact.
8. Hard limits can't be turned off by any request field. The one-off tool filter (PLN-05-03) only narrows further.

## Contract

Fields `eater_member_ids` and `leave_out_ingredient_ids` in PLN-05-01.

## Data

`library_search` gains `tools text[]`, `allergens text[]`, `diets text[]`, `diets_capable text[]`, `main_ingredient_ids
text[]` (non-optional main/base lines), with GIN indexes.

## Errors

| Case | gRPC status | Message |
|---|---|---|
| Eater id not in the profile | INVALID_ARGUMENT | `unknown member` |

## Acceptance criteria

- [ ] With the "student room" preset (microwave + kettle), 50 results contain no other tool (Verification 2, library side).
- [ ] With "no gluten, no nuts", 50 results contain zero dishes with gluten or nuts in a non-optional line
      (Verification 3, library side).
- [ ] A vegetarian eater sees tacos with beef and bean toppings, not beef stew.
- [ ] A dislike of coriander in an optional garnish keeps the dish and lists coriander to leave out.
