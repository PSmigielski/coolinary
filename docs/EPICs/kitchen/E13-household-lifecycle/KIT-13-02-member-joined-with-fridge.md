# KIT-13-02 Joining with a fridge: move stock, carry the person's data

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-13 Household lifecycle](README.md) |
| Release | R1 (stock move, used by the anonymous-account merge IDN-03-03); R2 (person data, invites IDN-07) |
| Size | M |
| Depends on | KIT-13-01, KIT-07-03 (merge rule), KIT-11-01 (member rows), IDN-01-03 |
| Unblocks | IDN-07-04 (join with fridge), IDN-03-03 (anonymous merge) |

## Goal

A person who joins a household can bring their fridge along, and their own diet and target come with them, so
nothing has to be typed twice (features.md §2 "Linking people into a household": "chooses whether to bring their own
fridge items along").

## Business rules

1. Part of the `household.member_added` handler (after KIT-11-01 has created the member row).
2. **`brought_fridge_from = X`** (X ≠ the new household):
   1. every stock item of X moves to the new household in one transaction: it merges into an identical item there
      (KIT-07-03 rule 6) or is re-keyed to the new household;
   2. movements: `moved_out` in X, `moved_in` in the new household (no savings effect, KIT-12-02 rule 5);
   3. if the new household would pass 1,000 items, items move in order of use-by date (soonest first, no date last)
      until the limit; the rest stay in X and the count is logged (Decided, Q140);
   4. both households get one `stock_version` bump and one `stock.changed`;
   5. X has no stock afterwards (or only what did not fit); identity deletes X afterwards and always sends this event
      first (IDN-01-03 rule 3).
3. **`previous_member_id = P`** (R2): the person's own data is copied from P's row to the new row: `target_kind`,
   `goal_set_by_parent`, target with its consent fields, `calories_hidden`, personal diet, allergens and dislikes.
   If P had no personal restrictions but P's household had household-level ones and P was its only account
   member, those household-level values are copied as P's personal values. Decided (Q139): the old household-level
   diet travels with them, since it was in fact that person's diet.
4. The new household's own profile sections are never overwritten (its cuisines, tools, staples stay).
5. If X or P no longer exists (deleted earlier), nothing is moved and the handler succeeds.
6. With `brought_fridge_from = null`, no stock moves (the default when the person chose not to bring it).

## Contract

No RPC. Consumes `household.member_added` (fields `brought_fridge_from`, `previous_member_id`).

## Data

No schema change.

## Events

Publishes `stock.changed` for both households and `profile.changed` (`members`) for the new one.

## Errors

| Case | Behaviour |
|---|---|
| Handler takes longer than 10 s (KIT-01-04 deadline) | retried; the move is one transaction, so a retry starts clean |

## Acceptance criteria

- [ ] Ewa (8 items) joins Adam's household (12 items, 2 identical to Ewa's) with her fridge → Adam's household holds
      18 items, Ewa's old household none (features.md Verification 15: both phones then show the same fridge).
- [ ] Joining without the fridge moves nothing.
- [ ] An anonymous account merged into an existing account (IDN-03-03) brings its stock.
- [ ] (R2) Kuba's nut-free and target come with him into the new household; the new household's tools are unchanged.
- [ ] Replaying the event moves nothing twice.
