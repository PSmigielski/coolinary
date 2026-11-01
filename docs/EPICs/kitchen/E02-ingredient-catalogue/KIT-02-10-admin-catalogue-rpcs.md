# KIT-02-10 Admin catalogue editing RPCs

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-02 Ingredient catalogue](README.md) |
| Release | R1 |
| Size | M |
| Depends on | KIT-01-07, KIT-02-01, KIT-02-02, KIT-02-03, KIT-02-04 |
| Unblocks | ADM-03, KIT-03-07 (re-match on new names) |

## Goal

The founder can add and fix ingredients, names, weights, storage and packs from the back office, and find what is
missing.

## Business rules

1. RPCs on `KitchenAdminService`: `CreateIngredient`, `UpdateIngredient` (category, base unit, allergens, diet
   flags, staple candidate, density, storage, shelf lives), `SetIngredientNames` (full replacement of one
   language's names), `ApproveMachineNames` (set `reviewed = true` for given name ids), `SetMeasures`,
   `SetPackSizes`, `RetireIngredient`, `ListIngredients`.
2. All validation rules of KIT-02-01 … KIT-02-04 apply; a failed rule returns `INVALID_ARGUMENT` naming the field.
3. Admin edits set `origin = 'admin'`, `reviewed = true` on names they touch.
4. Every write raises `catalogue_version`, writes `admin_changes`, and raises `nutrition_version` when measures or
   density change.
5. Decided (Q054): a change of allergens or diet flags writes the event `catalogue.ingredient_flags_changed`
   `{ingredient_ids}` (README event table), so planner re-derives the allergens of every dish that uses the
   ingredient at once (safety).
6. Decided (Q055): `RetireIngredient` hides the ingredient from search and from new matches but keeps every
   reference (stock, recipes, products) working.
7. `ListIngredients` filters: text query, category, status, `missing_language` (no primary in L),
   `missing_measures` (no piece weight for `PIECE` ingredients, no density for `MILLILITRE` ingredients),
   `missing_nutrition` (no approved match, KIT-06), `unreviewed_names`. Sorted by `popularity` descending; 50 per
   page, maximum 200.

## Contract

`contracts/coolinary/kitchen/v1/admin.proto`

```proto
rpc CreateIngredient(CreateIngredientRequest) returns (Ingredient);
rpc UpdateIngredient(UpdateIngredientRequest) returns (Ingredient);
rpc SetIngredientNames(SetIngredientNamesRequest) returns (Ingredient);
rpc ApproveMachineNames(ApproveMachineNamesRequest) returns (ApproveMachineNamesResponse);
rpc SetMeasures(SetMeasuresRequest) returns (Ingredient);
rpc SetPackSizes(SetPackSizesRequest) returns (Ingredient);
rpc RetireIngredient(RetireIngredientRequest) returns (Ingredient);
rpc ListIngredients(ListIngredientsRequest) returns (ListIngredientsResponse);
message UpdateIngredientRequest { string ingredient_id = 1; google.protobuf.FieldMask update_mask = 2; Ingredient ingredient = 3; }
```

## Data

No schema change.

## Events

May publish `catalogue.ingredient_flags_changed` v1 (pending the decision above) and
`catalogue.nutrition_changed` v1 (via KIT-06-10).

## Errors

| Case | Status | Message |
|---|---|---|
| Rule from KIT-02-01 … 04 broken | `INVALID_ARGUMENT` | names the field and rule |
| Unknown ingredient | `NOT_FOUND` | `ingredient not found` |
| Editing a merged ingredient | `FAILED_PRECONDITION` | `ingredient is merged` |

## Acceptance criteria

- [ ] Adding `gluten` to an ingredient records the change, raises the catalogue version, and (if accepted) emits
      `catalogue.ingredient_flags_changed`.
- [ ] `ListIngredients(missing_language=de)` lists exactly the ingredients without a German primary.
- [ ] A retired ingredient no longer appears in `SearchIngredients` but `GetIngredients` still returns it.
