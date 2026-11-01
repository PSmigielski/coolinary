# PLN-11-02 Anonymised copy into review

| | |
|---|---|
| Service | planner |
| Epic | [PLN-11 Library growth from well-rated AI dishes](README.md) |
| Release | R2 |
| Size | M |
| Depends on | PLN-11-01, PLN-03-07, PLN-03-05, PLN-03-06, PLN-02-04 |
| Unblocks | PLN-11-03 |

## Goal

Copy a candidate dish into the library review queue with nothing that ties it to the household it came from.

## Business rules

1. `PromoteCandidate(candidate_id)` creates a new library dish: `source=ai_promoted`, status `review`, a new slug
   from the English name (lower case, hyphens, a numeric suffix if taken), version 1 copied from the best dish's current
   version (meals, cuisine, tools, times, servings, ingredients, steps, texts).
2. The copy has no `household_id` and no `parent_dish_id`; the link lives only in `promotion_candidates.library_dish_id`
   (admin-only). On `household.deleted`, cluster member rows of that household's dishes are deleted with the dishes; the
   library copy stays.
3. Texts are copied as written by AI in that household's language; the review flags `ai_promoted_text` so the reviewer
   reads them for anything personal (AI texts may echo a wish). Missing languages are queued for translation
   (PLN-03-05), effort tags for suggestion (PLN-03-06); nutrition and flags are recomputed (PLN-04-01, PLN-02-03).
4. Ratings and history of the household dish are not copied.
5. Decided (Q215): users are told by a sentence in the terms and on the AI-dish card that their AI dishes may join
   the shared library in anonymised form; no opt-in.

## Contract

```proto
// PlannerAdminService
rpc PromoteCandidate(PromoteCandidateRequest) returns (PromoteCandidateResponse);
message PromoteCandidateRequest { string candidate_id = 1; }
message PromoteCandidateResponse { string library_dish_id = 1; }
```

## Data

Uses `promotion_candidates.library_dish_id` (PLN-11-01).

## Errors

| Case | gRPC status | Message |
|---|---|---|
| Candidate not open | FAILED_PRECONDITION | `candidate already decided` |
| Best dish deleted meanwhile | NOT_FOUND | `candidate dish gone` |

## Acceptance criteria

- [ ] The new library dish has no household id and no parent id, and appears in the review queue with
      `ai_promoted_text`.
- [ ] Deleting the source household leaves the library dish intact.
- [ ] Promoting twice returns FAILED_PRECONDITION the second time.
