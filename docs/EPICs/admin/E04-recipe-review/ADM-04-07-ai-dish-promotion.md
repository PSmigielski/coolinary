# ADM-04-07 Promote well-rated AI dishes to the library

| | |
|---|---|
| Service | admin |
| Epic | [ADM-04 Recipe import and dish review](README.md) |
| Release | R2 |
| Size | S |
| Depends on | ADM-04-03; PLN-11-02 (`PromoteCandidate`), PLN-11-03 (`ListPromotionCandidates`, `DismissCandidate`) |
| Unblocks | — |

## Goal

AI dishes composed for households that people rated well are reviewed and added to the library, so the same dish
later costs nothing to suggest (features.md §2 "Dishes and the day plan", D12).

## Business rules

1. Candidates are chosen by planner (PLN-11-01 sets the rating threshold and minimum number of households). The
   list is `ListPromotionCandidates(status = open)`, most households first: the best dish of each candidate (texts,
   ingredients, steps), cuisine, number of households, average rating, and the closest library dish with its
   similarity score.
2. A candidate shows no household id, member names or free-text wishes, only the dish itself.
3. Action `promote` → `PromoteCandidate(candidate_id)`: planner copies the dish into the library with source
   `ai_promoted` and status `review` (PLN-11-02); it then goes through the normal review (ADM-04-03), translation
   and publish steps. Action `dismiss` → `DismissCandidate(candidate_id, reason)` with a reason (≤ 200
   characters); planner proposes the same candidate again only if its household count doubles (PLN-11-03 rule 2).
4. Audited (`planner.dish.promote` / `.dismiss`, candidate id, new dish id).
5. Role `admin` only.

## Contract

| Method | Path | Body / response |
|---|---|---|
| GET | `/admin/v1/planner/promotion-candidates?page_token=` | `{ items }` |
| POST | `/admin/v1/planner/promotion-candidates/{id}/promote` | → `{ dish_id }` |
| POST | `/admin/v1/planner/promotion-candidates/{id}/dismiss` | `{ reason }` → `{}` |

## Data

None in admin.

## Errors

| Case | HTTP |
|---|---|
| Already promoted or dismissed | 409 (planner's `candidate_closed`) |
| Dismiss without a reason, or over 200 characters | 400 |

## Acceptance criteria

- [ ] Promoting a candidate creates a dish in the review queue with source `ai_promoted` and no household data.
