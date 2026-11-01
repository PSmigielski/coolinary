# ADM-03-04 Ingredient merge

| | |
|---|---|
| Service | admin |
| Epic | [ADM-03 Catalogue, matching and nutrition review](README.md) |
| Release | R1 |
| Size | S |
| Depends on | ADM-03-03; KIT-02-09 (`MergeIngredients` with `dry_run`, event `catalogue.ingredient_merged`) |
| Unblocks | — |

## Goal

Duplicate ingredients (e.g. "spring onion" and "scallion" created separately) become one, everywhere, without
breaking recipes, stock or offer matches.

## Business rules

1. Step 1, preview: the admin picks `from` and `into`; admin calls `MergeIngredients(from, into, reason, dry_run =
   true)` and shows kitchen's counts: references moved per table (`references_by_table`: products, stock lines,
   stored name matches, profiles, ...), names moved, names dropped, and `flags_differ`.
2. If allergen or diet flags differ, the merge is refused until the admin aligns them (ADM-03-03), because a merge
   must never silently drop an allergen. Admin never sends kitchen's `force = true` (KIT-02-09 rule 3); kitchen's
   `FAILED_PRECONDITION` is shown as 409 `flags_differ`.
3. Step 2, merge: requires `confirm_name` equal to the English name of `into`, typed by the admin, and a reason
   (3–200 characters, passed to kitchen); otherwise 400. Admin then calls `MergeIngredients(..., dry_run = false)`.
4. The merge cannot be undone. kitchen publishes `catalogue.ingredient_merged`; planner and offers update their own
   references (PLN-21, OFF-14). admin only reports kitchen's counts.
5. `from` = `into` → 400.
6. Audited (`kitchen.ingredient.merge`, both ids, the preview counts).
7. Role `admin` only.

## Contract

| Method | Path | Body / response |
|---|---|---|
| POST | `/admin/v1/kitchen/ingredient-merges/preview` | `{ from_id, into_id }` → `{ references_by_table, names_moved, names_dropped, flags_differ }` |
| POST | `/admin/v1/kitchen/ingredient-merges` | `{ from_id, into_id, confirm_name, reason }` → `{ references_moved, names_moved, names_dropped }` |

Both call `KitchenAdminService.MergeIngredients` (KIT-02-09), the preview with `dry_run = true`.

## Data

None in admin.

## Errors

| Case | HTTP |
|---|---|
| Same id, wrong `confirm_name`, missing reason | 400 |
| Flags differ | 409 `flags_differ` with the list |
| Either id unknown | 404 |

## Acceptance criteria

- [ ] Merging "scallion" into "spring onion" moves its products and stock lines; kitchen emits the event once.
- [ ] Merging two ingredients whose gluten flag differs is refused.
