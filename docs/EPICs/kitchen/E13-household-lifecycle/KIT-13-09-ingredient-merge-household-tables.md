# KIT-13-09 Register household tables for ingredient merges

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-13 Household lifecycle](README.md) |
| Release | R1 |
| Size | S |
| Depends on | KIT-02-09 (`ingredientRefs` registry), KIT-07-01, KIT-08-05, KIT-10-04, KIT-10-05 |
| Unblocks | — |

## Goal

When an admin merges two catalogue ingredients (zucchini and courgette), every household table that points to the
old one is moved to the new one in the same transaction, without duplicates.

## Business rules

1. The following columns are added to the KIT-02-09 `ingredientRefs` registry, with a conflict rule each:

   | Table.column | On a duplicate after re-pointing |
   |---|---|
   | `stock_items.ingredient_id` | none possible (no unique key); identical items are left side by side |
   | `stock_movements.ingredient_id` | none possible |
   | `profile_staples.ingredient_id` | the duplicate row is dropped |
   | `profile_dislikes.ingredient_id` | the duplicate row is dropped |
   | `receipt_name_links.ingredient_id` | none possible (key is the printed name) |

2. Every household whose staples or dislikes changed gets a `planning_version` bump and one `profile.changed`
   (`staples` or `diet`); every household whose stock changed gets a `stock_version` bump and one `stock.changed`.
3. The merge transaction of KIT-02-09 runs these updates; it may touch many households — the updates are set-based
   (one `UPDATE … WHERE ingredient_id = $from` per table), and events are written per household in the same
   transaction.
4. A registry test fails when a kitchen table has an `ingredient_id` column that is not registered.

## Contract

No RPC.

## Data

No schema change.

## Errors

A failure aborts the whole merge (KIT-02-09).

## Acceptance criteria

- [ ] Merge courgette into zucchini: a household with both as staples keeps one; stock items of courgette now point to
      zucchini; one `profile.changed` and one `stock.changed` for that household.
- [ ] The registry test lists every `ingredient_id` column in the kitchen schema.
