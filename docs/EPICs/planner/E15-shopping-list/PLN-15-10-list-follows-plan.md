# PLN-15-10 List follows plan changes

| | |
|---|---|
| Service | planner |
| Epic | [PLN-15 Shopping list](README.md) |
| Release | R2 |
| Size | M |
| Depends on | PLN-15-05, PLN-12-06 |
| Unblocks | PLN-21-02 |

## Goal

When the plan changes after its lines are on the list (someone marked away, a dish swapped, a meal deleted, a week
edited), the list follows: "the pot shrinks and the shopping list follows" (features.md §2, Verification 17).

## Business rules

1. Trigger: a planned meal's pot, dish or existence changes (PLN-12-06 recompute, PLN-09 swap, PLN-14-01 edit, plan
   deletion), for meals not yet cooked.
2. For each open (unticked) line with a `plan_meal` or `extra` contribution from that meal, the contribution is
   recomputed with PLN-15-02 and PLN-15-03 for that meal (current stock).
3. **Decreases and removals apply automatically:** the contribution shrinks or is removed; the line's amount and packs
   are recomputed (PLN-15-04); a line with no contributions left is removed.
4. **Increases and new ingredients are proposed, not added:** they go into a pending plan-change proposal
   (PLN-15-05 model) that the app shows as "Plan changed: 3 items to review"; accepting it adds them.
5. Lines with `manual_amount` keep their amount; they get the flag `plan_changed` and a suggested amount.
6. Ticked lines never change.
7. All automatic changes for one trigger are one list change (one version bump, one notice).
8. Meals in the past or cooked never trigger changes.

## Contract

`contracts/coolinary/planner/v1/list.proto`: `ShoppingList` gains `string pending_proposal_id = 4;` set when a
plan-change proposal is waiting. Accepting uses `AcceptListProposal`.

## Data

`ALTER TABLE list_proposals ADD COLUMN kind smallint NOT NULL DEFAULT 1;` (1 generated, 2 plan change); at most one
open plan-change proposal per household (partial unique index).

## Errors

None of its own: runs inside the change that triggered it (PLN-12-06, PLN-09, PLN-14-01), in the same transaction. If kitchen is down when recomputing stock (rule 2), decreases are still applied from the old stock snapshot and increases wait in the pending proposal; the line gets `recheck_pending = true` and is recomputed on the next list read.

## Acceptance criteria

- [ ] Verification 17: marking Ola away for dinner makes the pot 5; the curry lines on the list drop from × 1.5 to
      × 1.25 and are re-rounded to packs, automatically.
- [ ] Swapping dinner from curry to pierogi: curry-only lines are removed; pierogi ingredients appear as a pending
      proposal.
- [ ] A line whose amount was typed by hand keeps it and shows "plan changed".
- [ ] A ticked line is never changed.

## Notes

- Decided (Q207): rule 4 (increases go into a "plan changed, review" proposal; decreases apply directly).
