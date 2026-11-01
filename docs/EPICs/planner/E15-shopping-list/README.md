# PLN-15 Shopping list

Release R1 (live sharing and plan-follow details used from R2). The list is made without AI from the plan: take
every ingredient, scale it, subtract what is at home and the staples, merge, round to packs, add extras, let the
user review, and move ticked items into the fridge (features.md §2 "Shopping list", Verification 12). Offer badges
(PLN-17) and the trip price matrix (PLN-19) build on it in R2.

## Flow

1. `GenerateList` builds a proposal from a date range of the plan: needs (PLN-15-02) → minus stock and staples
   (PLN-15-03) → merged and rounded to packs (PLN-15-04).
2. The user unticks lines they don't want and accepts the proposal (PLN-15-05). Manual lines can be added at any
   time (PLN-15-06).
3. In the shop, ticking a line adds it to the fridge (PLN-15-07).
4. Every change bumps the list version, is idempotent for offline replays (PLN-15-08) and is announced on Redis
   pub/sub for the live list (PLN-15-09).
5. When the plan changes (someone away, a swap), the list follows (PLN-15-10).

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [PLN-15-01](PLN-15-01-list-model.md) | List model and `GetList` | R1 | M | PLN-01, KIT-02 (GetIngredients) |
| [PLN-15-02](PLN-15-02-needs-from-plan.md) | Ingredient needs from the plan | R1 | M | PLN-15-01, PLN-08, PLN-02 |
| [PLN-15-03](PLN-15-03-subtract-stock-and-staples.md) | Subtract stock and staples | R1 | M | PLN-15-02, KIT-07 (ListStock), KIT-10 (GetPlanningProfile) |
| [PLN-15-04](PLN-15-04-merge-and-round-to-packs.md) | Merge lines and round to packs | R1 | M | PLN-15-03, KIT-02 (GetIngredients) |
| [PLN-15-05](PLN-15-05-generate-review-accept.md) | Generate, review and accept proposed lines | R1 | M | PLN-15-04 |
| [PLN-15-06](PLN-15-06-manual-lines-and-editing.md) | Manual lines, editing, removing, clearing | R1 | M | PLN-15-01, PLN-15-08, KIT-03 (MatchNames) |
| [PLN-15-07](PLN-15-07-tick-into-fridge.md) | Tick into the fridge (and undo) | R1 | M | PLN-15-01, PLN-15-08, KIT-07 (AddStockItems, RemoveStockItems) |
| [PLN-15-08](PLN-15-08-versions-idempotency-conflicts.md) | Versions, idempotent offline edits, conflict rules | R1 | M | PLN-15-01 |
| [PLN-15-09](PLN-15-09-live-notices-and-changes.md) | Live change notices and `GetListChanges` | R1 | S | PLN-15-08 |
| [PLN-15-10](PLN-15-10-list-follows-plan.md) | List follows plan changes | R2 | M | PLN-15-05, PLN-12-06 |
