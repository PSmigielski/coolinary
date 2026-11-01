# ADM-04-08 Library settings: base pantries, no-cook extras, coverage

| | |
|---|---|
| Service | admin |
| Epic | [ADM-04 Recipe import and dish review](README.md) |
| Release | R1 (base pantries, coverage); R2 (no-cook extras) |
| Size | M |
| Depends on | ADM-04-03; PLN-16-01 (`ListBasePantries`, `SetBasePantry`), PLN-12-04 (`ListNoCookExtras`, `UpsertNoCookExtra`), PLN-03-09 (`GetLibraryCoverage`); KIT-02-05 (`SearchIngredients`) |
| Unblocks | — |

## Goal

The founder can maintain the lists the planner uses next to the dishes: each cuisine's base pantry, the no-cook
extras used to top up portions, and see where the starter library is still thin before writing new briefs.

## Business rules

1. Base pantries (R1): `ListBasePantries` shows one list per pantry key (a cuisine, or `asian/<sub-cuisine>`,
   PLN-16-01 rule 1) with its ingredients in order, who changed it last and when. Editing sends the whole ordered
   list with `SetBasePantry(cuisine, ingredient_ids)`; ingredients are picked with `KitchenService.SearchIngredients`.
   Planner's limits apply: 4–12 ingredients, every id known to kitchen (PLN-16-01 rules 2 and 4).
2. No-cook extras (R2): `ListNoCookExtras` shows ingredient, one serving (amount and unit), label key, default rank,
   lighter rank (empty = not offered when the family eats lighter), `child_ok` and active. `UpsertNoCookExtra`
   creates or changes one row. Ranks are whole numbers from 1; rows may share a rank (banana and apple are both
   rank 2 in planner's seed). Turning `child_ok` on for nuts asks for a confirmation (PLN-12-04 rule 4 skips them
   for children on purpose).
3. Library coverage (R1): `GetLibraryCoverage` shows published dishes per cuisine × meal × tool preset × diet, the
   cells below target with their gap, and the Verification 1 simulation on the test fridges (PLN-03-09). Each
   cell below target has a "draft dishes for this" link that opens ADM-04-02 with a brief filled in (cuisine, meal,
   tools of the preset, diet, `count` = the gap, at most 10).
4. Read access for roles `admin` and `dietitian` on coverage; edits for `admin` only.
5. Audited: `planner.base_pantry.set` (key, ingredient ids), `planner.no_cook_extra.upsert` (ingredient id, changed
   fields). Coverage views are not audited.

## Contract

| Method | Path | Body / response |
|---|---|---|
| GET | `/admin/v1/planner/base-pantries` | `{ items: [{ key, ingredient_ids, ingredients, updated_by, updated_at }] }` |
| PUT | `/admin/v1/planner/base-pantries/{key}` | `{ ingredient_ids }` → pantry |
| GET | `/admin/v1/planner/no-cook-extras` | `{ items }` |
| PUT | `/admin/v1/planner/no-cook-extras/{ingredientId}` | `{ serving, label_key, rank_default, rank_lighter?, child_ok, active }` → extra |
| GET | `/admin/v1/planner/coverage` | `{ cells, gaps, fridge_simulation }` |

Calls `PlannerAdminService.ListBasePantries`, `SetBasePantry` (PLN-16-01), `ListNoCookExtras`, `UpsertNoCookExtra`
(PLN-12-04), `GetLibraryCoverage` (PLN-03-09), and `KitchenService.SearchIngredients` for the pickers.

## Data

None in admin.

## Errors

| Case | HTTP |
|---|---|
| Fewer than 4 or more than 12 pantry ingredients, unknown ingredient id | 400 (planner's `INVALID_ARGUMENT`) |
| Rank below 1 | 400 |
| Role `dietitian` on an edit | 403 `role_required` |

## Acceptance criteria

- [ ] Setting the Italian pantry to olive oil, passata, pasta, hard cheese stores it in that order, and the one-tap
      "add base items" in the app (PLN-16-03) then adds those that are missing.
- [ ] A pantry with 3 ingredients is refused with planner's message.
- [ ] A coverage cell below target opens the drafting form with cuisine, meal, tools and count filled in.
