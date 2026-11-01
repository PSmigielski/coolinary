# KIT-02-09 Merge duplicate ingredients

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-02 Ingredient catalogue](README.md) |
| Release | R1 |
| Size | M |
| Depends on | KIT-02-04, KIT-01-02, KIT-01-07, KIT-03-02 |
| Unblocks | ADM-03 (merge screen), PLN-21 and OFF-14 (consume `catalogue.ingredient_merged`) |

## Goal

When two ingredients turn out to be the same thing, an admin merges them and every reference — in kitchen and in
other services — ends up pointing at one of them.

## Business rules

1. Admin RPC `MergeIngredients(from_id, into_id, reason)`. Both must be `active` and different.
2. In one transaction:
   1. `from` gets `status = 'merged'`, `merged_into = into`.
   2. `from`'s names move to `into` as `synonym` (kind kept as `plural` for plurals); names that collide with an
      existing name of `into` are dropped.
   3. Measures: `into`'s values win; measures only `from` has are copied.
   4. Storage, shelf lives, pack sizes, flags: `into`'s values stay.
   5. Nutrition: `into`'s approved match stays; if `into` has none and `from` has one, it moves over.
   6. Every kitchen table with an `ingredient_id` column is re-pointed through one registry function
      (`ingredientRefs`), so tables added later (stock KIT-07, profile KIT-10/11, matches KIT-03, products KIT-04)
      are covered by adding them to the registry.
   7. Earlier merges into `from` are re-pointed to `into`, so `merged_into` always names an active ingredient.
   8. An outbox event `catalogue.ingredient_merged` `{from_ingredient_id, into_ingredient_id}` is written.
   9. If `into`'s nutrition changed (rule 5), `nutrition_version` is raised and `catalogue.nutrition_changed` is
      written for `into`.
   10. `admin_changes` records both ingredients' before and after state.
3. If allergens or diet flags of `from` and `into` differ, the merge is refused unless `force = true`, because a
   recipe written with `from` would silently change its allergens.
4. `dry_run = true` runs every check (rules 1 and 3) and returns the counts the merge would produce
   (`references_moved` per table in `references_by_table`, `names_moved`, `names_dropped`, plus
   `flags_differ`) without writing anything. The admin preview (ADM-03-04) uses it.
5. Decided (Q056): merges cannot be undone. If a merge was wrong, an admin creates a new ingredient and re-points
   what needs it.

## Contract

`contracts/coolinary/kitchen/v1/admin.proto`

```proto
rpc MergeIngredients(MergeIngredientsRequest) returns (MergeIngredientsResponse);
message MergeIngredientsRequest { string from_ingredient_id = 1; string into_ingredient_id = 2; string reason = 3; bool force = 4; bool dry_run = 5; }
message MergeIngredientsResponse { int32 references_moved = 1; int32 names_moved = 2; int32 names_dropped = 3;
  map<string, int32> references_by_table = 4; bool flags_differ = 5; bool dry_run = 6; }
```

## Data

No new table; uses `ingredients.merged_into`.

## Events

Publishes `catalogue.ingredient_merged` v1; may publish `catalogue.nutrition_changed` v1.

## Errors

| Case | Status | Message |
|---|---|---|
| Same id twice, or either not active | `FAILED_PRECONDITION` | `both ingredients must be active and different` |
| Allergens or diet flags differ without `force` | `FAILED_PRECONDITION` | `allergens or diet flags differ` |
| Unknown id | `NOT_FOUND` | `ingredient not found` |

## Acceptance criteria

- [ ] After merging "courgette" into "zucchini", search for "courgette" returns zucchini, `GetIngredients` for the
      old id returns zucchini with `requested_id`, and a stock item that used the old id now uses the new one.
- [ ] A then B merged into C leaves A and B both pointing at C.
- [ ] Exactly one `catalogue.ingredient_merged` event per merge reaches `events.kitchen`.
- [ ] `dry_run = true` returns the same counts as the real merge that follows, and writes no row, change record or
      event.
