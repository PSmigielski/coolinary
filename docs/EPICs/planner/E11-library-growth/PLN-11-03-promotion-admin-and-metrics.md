# PLN-11-03 Promotion admin RPCs and savings metric

| | |
|---|---|
| Service | planner |
| Epic | [PLN-11 Library growth from well-rated AI dishes](README.md) |
| Release | R2 |
| Size | S |
| Status | todo |
| Depends on | PLN-11-01, PLN-11-02, PLN-01-06 |
| Unblocks | ADM-04 (promotion queue) |

## Goal

Admin sees and decides promotion candidates; the founder can see whether library growth actually replaces AI work.

## Business rules

1. `ListPromotionCandidates(status=open, page)` sorted by households desc, then avg rating desc; each item shows the
   best dish (texts, ingredients, steps), cuisine, households, average rating and the closest library dish with its
   Jaccard score.
2. `DismissCandidate(candidate_id, reason)` with reason ≤ 200 chars; dismissed clusters are not proposed again unless
   their households count doubles.
3. Every decision writes `admin_changes` (PLN-01-06).
4. Metric `planner_plan_meals_total{source}` counts day-plan meals by `ai` (new dish), `ai_library_pick` (AI chose a
   library dish id), `library` (non-AI) — the share of `ai_library_pick` + `library` shows how often a plan reused the
   library. Decided (Q216): planner does not skip the AI call yet; this metric is measured first and the choice is
   made with launch data.

## Contract

```proto
// PlannerAdminService
rpc ListPromotionCandidates(ListPromotionCandidatesRequest) returns (ListPromotionCandidatesResponse);
rpc DismissCandidate(DismissCandidateRequest) returns (DismissCandidateResponse);
```

## Data

No schema change.

## Errors

| Case | gRPC status / behaviour |
|---|---|
| Unknown candidate | `NOT_FOUND` |
| Candidate already promoted or dismissed | `FAILED_PRECONDITION`, `x-error-code: candidate_closed` |
| Reason longer than 200 characters | `INVALID_ARGUMENT` |
| Called without `x-admin-user-id` | `PERMISSION_DENIED` |

## Acceptance criteria

- [ ] Candidates come back ordered by households, then rating.
- [ ] A dismissed cluster with 2 households reappears only when it reaches 4.
- [ ] The metric increments per meal of each stored day plan with the right `source`.
