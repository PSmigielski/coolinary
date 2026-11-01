# PLN-21-04 Ingredient merge for lists, pantries and extras

| | |
|---|---|
| Service | planner |
| Epic | [PLN-21 Household lifecycle](README.md) |
| Release | R1 |
| Size | S |
| Status | todo |
| Depends on | PLN-01, PLN-15-01 |
| Unblocks | — |

## Goal

When kitchen merges two duplicate catalogue ingredients (`catalogue.ingredient_merged`), every planner reference
outside recipes moves to the surviving id.

## Business rules

1. Rewrite `from_ingredient_id` → `into_ingredient_id` in list lines, list line sources, proposals (they expire anyway;
   rewrite open ones), base pantries, no-cook extras, plan meal offer items and top-ups.
2. If a list then has two open lines for the same ingredient, they are merged into the older line: contributions
   summed, packs re-rounded (PLN-15-04), the newer line tombstoned. One list version bump and live notice per
   affected list.
3. In base pantries, a duplicate after rewrite is removed (keep the earlier position).
4. Idempotent: a second run finds nothing to rewrite.

## Contract

Event consumer only.

## Data

No schema change.

## Errors

Handler errors roll back and are retried by the consumer framework; after 5 failures the event goes to `dead_events` (PLN-01-01). A merge whose source id no longer appears anywhere is a success.

## Acceptance criteria

- [ ] A list with "Zwiebel" and "onion" lines that kitchen merges ends with one onion line with both contributions.
- [ ] The Italian base pantry keeps one entry after merging two "passata" ids.
