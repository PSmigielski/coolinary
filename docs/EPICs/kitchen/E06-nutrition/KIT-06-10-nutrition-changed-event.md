# KIT-06-10 Nutrition version and `catalogue.nutrition_changed`

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-06 Nutrition](README.md) |
| Release | R1 |
| Size | S |
| Status | todo |
| Depends on | KIT-06-01, KIT-01-02, KIT-02-02 |
| Unblocks | PLN-04 (recalculate dish calories), KIT-06-06 |

## Goal

When anything that changes an ingredient's calories changes, planner learns which ingredients are affected and
recalculates only the dishes that use them.

## Business rules

1. An ingredient's `nutrition_version` is raised by 1, in the same transaction, when:
   - its approved link is set, changed or removed (KIT-06-06, KIT-02-09);
   - the values of its approved entry change on re-import (KIT-06-02/03/04);
   - its measures or density change (KIT-02-02, KIT-02-10).
2. Each such transaction writes one outbox event `catalogue.nutrition_changed` `{ingredient_ids}` with every
   ingredient it touched, split into events of at most 500 ids.
3. Re-imports touching many ingredients write the events at the end of the import transaction.
4. Pack value changes (KIT-06-07) raise nothing: library dishes don't use pack values.
5. `GetIngredients` returns `nutrition_version` (KIT-02-06), so a consumer that missed events can compare.

## Contract

No RPC.

## Data

No schema change (column added in KIT-02-01).

## Events

Publishes `catalogue.nutrition_changed` v1 `{ingredient_ids: [uuid]}` on `events.kitchen`.

## Errors

None.

## Acceptance criteria

- [ ] Approving a link emits one event naming that ingredient.
- [ ] A re-import changing 1,200 linked entries emits 3 events covering all 1,200 ids.
- [ ] Changing a pack value emits nothing.
