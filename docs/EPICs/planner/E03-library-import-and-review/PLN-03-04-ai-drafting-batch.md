# PLN-03-04 AI drafting batch (`draft_dishes`)

| | |
|---|---|
| Service | planner |
| Epic | [PLN-03 Library import, AI drafting and review](README.md) |
| Release | R1 |
| Size | L |
| Depends on | PLN-03-02, PLN-02-02, PLN-02-03, AI-04 (SubmitBatch, GetBatch, GetBatchResults), AI-09 (`draft_dishes`) |
| Unblocks | PLN-03-05, PLN-03-06, ADM-04 (AI drafts) |

## Goal

Admin describes what the library lacks ("10 Polish breakfasts, microwave + kettle, vegetarian, ≤ 15 minutes"); an AI
batch job (50% cheaper, decision 5) drafts them; planner checks every draft without AI and puts the good ones in the
review queue. Cost is fixed per batch, not per user.

## Business rules

1. `DraftDishes(briefs)`: 1–50 briefs per call; each brief = cuisine (+ sub-cuisine), meal, allowed tools, required
   diets, max total minutes, effort tags wanted, language to write in (`pl`, `de` or `en`), count 1–10. ≤ 200
   dishes per call.
2. Before submitting, planner adds to each brief the names of existing library dishes in that cuisine and meal (up to
   100) so AI avoids repeats, and the catalogue ingredient names allowed (top 500 by use in that cuisine, from
   kitchen) so drafts use known ingredients.
3. One `SubmitBatch(kind="draft_dishes")` call; item `custom_id` = `<brief_index>`. The returned batch id is stored in
   `planner_batches` with kind, status `submitted`, counts, admin user (D10: planner stores the id, ai stays stateless).
4. A poll job every 10 minutes calls `GetBatch` for each batch not `ended`; on `ended` it fetches `GetBatchResults`,
   processes them and sets status `processed`. Batches not ended after 26 hours are marked `expired` and reported.
5. Each draft is checked, in order, without AI:
   a. name resolution (PLN-03-02); any unresolved name → draft dropped with reason `unresolved_ingredients`;
   b. validator (PLN-02-02); blocking problem → dropped;
   c. brief constraints: tools ⊆ brief tools, every required diet in derived `diets_capable`, total minutes ≤ brief max,
      meal and cuisine equal the brief → else dropped with the failed constraint;
   d. duplicate check: same cuisine and (name trigram similarity ≥ 0.8 in the draft's language, or ingredient-id set
      Jaccard ≥ 0.8) against library and review dishes → kept, flagged `possible_duplicate_of=<dish_id>`.
6. Kept drafts become dishes with `source=ai_draft`, status `review`, texts in the brief language with
   `machine_translated=false`. Nutrition and flags are computed (PLN-04-01, PLN-02-03); AI never provides calories.
7. The batch report (`GetDraftBatch`) lists counts per brief: requested, received, kept, dropped by reason.
8. Item errors (`errored`, `expired` in batch results) are counted, never retried automatically.
9. `ListDraftBatches(kind?, status?)` lists every `planner_batches` row (kinds `draft_dishes`,
   `translate_dishes`, `suggest_effort_tags`), newest first, 50 per page, with kind, status, item count, admin user
   id, submitted and processed times and the report counts; without `kind` it lists all kinds. `GetDraftBatch`
   returns the report of any kind. Admin (ADM-04-02, ADM-04-04, ADM-04-06) uses these as its one jobs view.

## Contract

```proto
// PlannerAdminService
rpc DraftDishes(DraftDishesRequest) returns (DraftDishesResponse);
rpc GetDraftBatch(GetDraftBatchRequest) returns (GetDraftBatchResponse);
rpc ListDraftBatches(ListDraftBatchesRequest) returns (ListDraftBatchesResponse);   // rule 9, all batch kinds
message ListDraftBatchesRequest { string kind = 1; string status = 2; coolinary.common.v1.PageRequest page = 3; }

message DishBrief {
  coolinary.common.v1.Cuisine cuisine = 1;
  coolinary.common.v1.Cuisine sub_cuisine = 2;
  Meal meal = 3;
  repeated coolinary.common.v1.KitchenTool tools = 4;
  repeated coolinary.common.v1.Diet diets = 5;
  int32 max_total_minutes = 6;
  repeated coolinary.common.v1.EffortTag effort_tags = 7;
  string language = 8;
  int32 count = 9;
}
message DraftDishesRequest { repeated DishBrief briefs = 1; }
message DraftDishesResponse { string batch_id = 1; }
message GetDraftBatchRequest { string batch_id = 1; }
message BriefOutcome {
  int32 brief_index = 1;
  int32 requested = 2; int32 received = 3; int32 kept = 4;
  map<string, int32> dropped_by_reason = 5;
  repeated string dish_ids = 6;
}
message GetDraftBatchResponse { string status = 1; repeated BriefOutcome briefs = 2; repeated UnresolvedName unresolved = 3; }
```

## Data

`planner/migrations/000NN_planner_batches.sql` (shared by PLN-03-05 and PLN-03-06)

```sql
CREATE TABLE planner_batches (
    id             text PRIMARY KEY,                 -- ai batch id
    kind           text NOT NULL CHECK (kind IN ('draft_dishes','translate_dishes','suggest_effort_tags')),
    status         text NOT NULL CHECK (status IN ('submitted','ended','processed','expired','failed')),
    admin_user_id  text NOT NULL,
    request        jsonb NOT NULL,                   -- briefs or dish ids, no user data
    report         jsonb,
    submitted_at   timestamptz NOT NULL DEFAULT now(),
    processed_at   timestamptz
);
CREATE INDEX planner_batches_open ON planner_batches (status) WHERE status IN ('submitted','ended');
```

## Errors

| Case | gRPC status | Message |
|---|---|---|
| > 50 briefs, > 200 dishes, count outside 1–10 | INVALID_ARGUMENT | limits |
| ai unavailable or over its spend limit | UNAVAILABLE / RESOURCE_EXHAUSTED | nothing stored |

## Acceptance criteria

- [ ] A brief "polish, breakfast, tools [microwave, kettle], vegetarian, ≤ 15 min, count 5" yields only kept drafts
      whose tools ⊆ {microwave, kettle} and whose derived flags include vegetarian (Verification 2 at draft level).
- [ ] A draft containing an unknown ingredient name is dropped with `unresolved_ingredients` and its name appears in
      `unresolved`.
- [ ] A draft named like an existing dish is kept with `possible_duplicate_of` set.
- [ ] Restarting planner during polling processes each batch exactly once (`processed` status guards it).
- [ ] No dish written by this task carries a calorie number from AI (checked by PLN-02-02 rule 10).
- [ ] After one drafting and one translation batch, `ListDraftBatches(kind="translate_dishes")` returns only the
      translation batch and `ListDraftBatches()` returns both.
