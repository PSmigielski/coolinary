# PLN-03-06 Effort tag suggestions (`suggest_effort_tags`)

| | |
|---|---|
| Service | planner |
| Epic | [PLN-03 Library import, AI drafting and review](README.md) |
| Release | R1 |
| Size | S |
| Depends on | PLN-03-04 (planner_batches, poll job), AI-04, AI-09 (`suggest_effort_tags`) |
| Unblocks | PLN-03-07, PLN-12, PLN-14 (keeps/freezes tags) |

## Goal

Effort tags (no cooking, keeps 3 days, freezes well, good cold the next day, base plus toppings) drive the family and
week plans. AI suggests them once for everyone; a person accepts or rejects each one (features.md "Calories and goals",
"both are checked by hand before use").

## Business rules

1. `SuggestEffortTags(dish_ids | all_untagged=true)` submits a batch for library versions in review or published,
   ≤ 500 per batch.
2. Suggestions are stored per (dish, version, tag) with AI's yes/no and a one-line reason; they never change the
   dish's `effort_tags` directly.
3. Checks without AI before storing: `no_cook` is suggested only if PLN-02-02 rule 8 holds; `base_plus_toppings` only if
   rule 9 holds; others are stored as given.
4. In review (PLN-03-07) the reviewer accepts or rejects each suggestion; accepting adds the tag to the version
   (texts and tags on a published version change without a new version, like translations).
5. Suggestions not decided within 60 days stay pending; nothing is applied automatically.
6. `dry_run = true` selects exactly as rule 1 and returns the number of versions it would submit, without
   submitting or writing anything (ADM-04-06 cost confirmation).

## Contract

```proto
// PlannerAdminService
rpc SuggestEffortTags(SuggestEffortTagsRequest) returns (SuggestEffortTagsResponse);
rpc DecideEffortTag(DecideEffortTagRequest) returns (DecideEffortTagResponse);
message SuggestEffortTagsRequest { repeated string dish_ids = 1; bool all_untagged = 2; bool dry_run = 3; }
message SuggestEffortTagsResponse { string batch_id = 1; int32 items = 2; }
message DecideEffortTagRequest { string dish_id = 1; int32 version = 2; coolinary.common.v1.EffortTag tag = 3; bool accept = 4; }
```

## Data

`planner/migrations/000NN_effort_tag_suggestions.sql`

```sql
CREATE TABLE effort_tag_suggestions (
    dish_id     uuid NOT NULL,
    version     int  NOT NULL,
    tag         text NOT NULL,
    suggested   boolean NOT NULL,
    reason      text NOT NULL DEFAULT '',
    decision    text CHECK (decision IN ('accepted','rejected')),
    decided_by  text,
    decided_at  timestamptz,
    PRIMARY KEY (dish_id, version, tag),
    FOREIGN KEY (dish_id, version) REFERENCES dish_versions (dish_id, version) ON DELETE CASCADE
);
```

## Errors

| Case | gRPC status / behaviour |
|---|---|
| Neither `dish_ids` nor `all_untagged`, or both | `INVALID_ARGUMENT` |
| More than 500 versions selected | `INVALID_ARGUMENT` (`batch_too_large`) |
| Unknown dish or version | `NOT_FOUND` |
| Decision on a version that is retired or no longer current | `FAILED_PRECONDITION` (`version_not_current`) |
| Decision on a tag with no stored suggestion | `NOT_FOUND` |
| Same decision sent twice | success, no change (idempotent) |
| ai down or over its spend limit at submission | `UNAVAILABLE` / `RESOURCE_EXHAUSTED` passed through |
| A suggestion fails rule 3 | not stored; counted in the batch report |

## Acceptance criteria

- [ ] AI suggesting `no_cook` for an oven dish is discarded by rule 3.
- [ ] Accepting `freezes_well` adds it to the version's tags and writes `admin_changes`.
- [ ] No tag reaches `dish_versions.effort_tags` without a `decision='accepted'` row or an import/admin edit.
- [ ] `dry_run = true` returns the item count and creates no batch.
