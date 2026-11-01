# ADM-04-03 Dish review queue and editor

| | |
|---|---|
| Service | admin |
| Epic | [ADM-04 Recipe import and dish review](README.md) |
| Release | R1 |
| Size | L |
| Depends on | ADM-04-01 or ADM-04-02; PLN-03-07 (`ListDishes`, `ListReviewQueue`, `GetReviewItem`, `UpdateDraft`, `ApproveDish`, `RejectDish`), PLN-01-06 (`ListAdminChanges`) |
| Unblocks | ADM-04-04, ADM-04-05, ADM-04-07 |

## Goal

Every dish, imported or AI-drafted, is read, corrected and approved by a human before it can be published, with the
facts needed to judge it on one page: ingredients as matched, calories as computed, and the rule check.

## Business rules

1. Two lists:
   - the review queue, `ListReviewQueue` (default view): versions in `review`, oldest first, with planner's flags
     (`calories_incomplete`, `low_confidence_match`, `possible_duplicate`, `machine_translated`,
     `effort_tags_pending`, `dietitian_pending`);
   - all dishes, `ListDishes` (PLN-03-07 rule 10): statuses `draft`, `review`, `published`, `retired`, filters
     cuisine, meal, source (`import`, `ai_draft`; household AI dishes are never listed), flags, text search on the
     name. 50 per page.
2. Detail is `GetReviewItem`: names and steps in each language present, meal, cuisine, tools, total and active
   time, base portions, ingredients with their matched catalogue ingredient, amounts and units, effort tags (and
   pending suggestions), base/topping groups, derived diet and allergen flags, calories and nutrients per portion
   with the per-ingredient breakdown (or "calories incomplete" with the missing ingredients, features.md §5 rule
   6), planner's problems and warnings, and for an edit of a published dish a field-by-field diff. The change
   history comes from `ListAdminChanges(target = "dish:<id>")`.
3. Opening an item takes planner's 15-minute review claim (PLN-03-07 rule 8). Editing (`UpdateDraft`) and deciding
   work only while this admin holds it; another admin's claim → 409 `being_reviewed` with that admin's name. planner
   recomputes calories, flags and problems on save and returns them.
4. Review actions:
   - `approve` → `ApproveDish(dish_id, version, confirmed_flags)`: publishes the version. Every warning flag must be
     ticked by the admin and is sent in `confirmed_flags` (PLN-03-07 rule 5);
   - `reject` → `RejectDish(reason)`: back to `draft`, reason (≤ 500 characters) required and shown to the next
     reviewer.
5. Approving fails while planner's publication rules don't pass (unknown ingredients, rule-check errors): planner
   returns `FAILED_PRECONDITION` with the problems; admin shows them.
6. Role `admin` for everything; role `dietitian` may list and read (rules 1–2) but not edit or decide, and opening a
   dish as a dietitian takes no claim.
7. Audited: `planner.dish.update` (changed field names), `planner.dish.review` (action, reason).

## Contract

| Method | Path | Body / response |
|---|---|---|
| GET | `/admin/v1/planner/review-queue?source=&cuisine=&meal=&flags=&page_token=` | `{ items, next_page_token }` |
| GET | `/admin/v1/planner/dishes?status=&cuisine=&meal=&source=&flags=&q=&page_token=` | `{ items, next_page_token }` |
| GET | `/admin/v1/planner/dishes/{id}` | review item with problems, nutrition, diff, history |
| PATCH | `/admin/v1/planner/dishes/{id}` | `{ version, ...changes }` → updated item |
| POST | `/admin/v1/planner/dishes/{id}/review` | `{ action: approve\|reject, version, confirmed_flags?, reason? }` → item |

## Data

None in admin.

## Errors

| Case | HTTP |
|---|---|
| Claimed by another admin | 409 `being_reviewed` |
| Approve with an unconfirmed warning flag, unknown ingredients or rule errors | 409 with planner's list |
| Reject without a reason | 400 |

## Acceptance criteria

- [ ] An AI draft whose step needs an oven but whose tools list only "hob" shows the warning and cannot be approved.
- [ ] Replacing an ingredient updates the shown calories on save.
- [ ] A dietitian account can open a dish but gets 403 on PATCH.
