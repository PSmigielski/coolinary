# PLN-03-07 Review queue and decisions

| | |
|---|---|
| Service | planner |
| Epic | [PLN-03 Library import, AI drafting and review](README.md) |
| Release | R1 |
| Size | L |
| Depends on | PLN-02-04, PLN-02-03, PLN-04-01, PLN-01-06 |
| Unblocks | ADM-04, PLN-11-02 |

## Goal

One queue where a person checks every new dish or new version — imported, AI-drafted or promoted — edits it if needed,
and publishes or rejects it. Nothing reaches users without passing here.

## Business rules

1. Queue = all versions with status `review` (new dishes) or open draft versions in `review` (edits of published
   dishes). Sorted oldest first.
2. Filters: source, cuisine, meal, language missing, flags. Flags computed per item: `calories_incomplete`,
   `low_confidence_match` (PLN-03-02), `possible_duplicate`, `machine_translated`, `effort_tags_pending`,
   `dietitian_pending` (PLN-03-08).
3. `GetReviewItem` returns the version with all texts, problems (warnings included), nutrition with per-ingredient
   breakdown and missing ingredients, derived flags, pending effort-tag suggestions and, for new versions of published
   dishes, a field-by-field diff against the live version.
4. `UpdateDraft` edits the open version (any field); runs the validator; recomputes flags and nutrition.
5. `ApproveDish` publishes (PLN-02-04) when the publication rules (PLN-02-02 rules 11–13) pass. The reviewer must
   confirm each warning flag explicitly (`confirmed_flags`), so warnings can't be skipped by accident.
6. `RejectDish(reason)` moves the version back to `draft` with a reason (≤ 500 chars) visible to the next reviewer.
7. `RetireDish(reason)` and `RestoreDish` on published/retired dishes.
8. Reviewing is single-writer: `GetReviewItem` takes a 15-minute claim for that admin; another admin sees "being
   reviewed by X" and cannot approve until the claim ends or is released.
9. Every decision writes `admin_changes`.
10. `ListDishes` lists library dishes in any status (`draft`, `review`, `published`, `retired`) for admin, with
   the filters of rule 2 plus status, text query on the name, and the flags `dietitian_pending` and
   `effort_tags_pending` also for published dishes (the review queue holds only versions in review). 50 per page,
   newest change first. Household AI dishes (`source=household_ai`) are never listed.

## Contract

```proto
// PlannerAdminService
rpc ListReviewQueue(ListReviewQueueRequest) returns (ListReviewQueueResponse);
rpc GetReviewItem(GetReviewItemRequest) returns (GetReviewItemResponse);
rpc UpdateDraft(UpdateDraftRequest) returns (UpdateDraftResponse);
rpc ApproveDish(ApproveDishRequest) returns (ApproveDishResponse);
rpc RejectDish(RejectDishRequest) returns (RejectDishResponse);
rpc RetireDish(RetireDishRequest) returns (RetireDishResponse);
rpc RestoreDish(RestoreDishRequest) returns (RestoreDishResponse);
rpc ListDishes(ListDishesRequest) returns (ListDishesResponse);   // rule 10

message ApproveDishRequest { string dish_id = 1; int32 version = 2; repeated string confirmed_flags = 3; }
message ListDishesRequest { repeated string statuses = 1; string query = 2; string source = 3;
  coolinary.common.v1.Cuisine cuisine = 4; Meal meal = 5; repeated string flags = 6; coolinary.common.v1.PageRequest page = 7; }
message AdminDishSummary { string dish_id = 1; int32 version = 2; string status = 3; string name = 4; string source = 5;
  repeated string flags = 6; int32 kcal_per_portion = 7; string updated_at = 8; }
message ListDishesResponse { repeated AdminDishSummary dishes = 1; coolinary.common.v1.PageResponse page = 2; }
```

## Data

`planner/migrations/000NN_review.sql`

```sql
CREATE TABLE review_claims (
    dish_id        uuid NOT NULL,
    version        int  NOT NULL,
    admin_user_id  text NOT NULL,
    expires_at     timestamptz NOT NULL,
    PRIMARY KEY (dish_id, version)
);
ALTER TABLE dish_versions
    ADD COLUMN review_note text NOT NULL DEFAULT '',
    ADD COLUMN low_confidence_lines int[] NOT NULL DEFAULT '{}',
    ADD COLUMN possible_duplicate_of uuid;
```

## Errors

| Case | gRPC status | Message |
|---|---|---|
| Approve with an unconfirmed warning flag | FAILED_PRECONDITION | `confirm flags: …` |
| Claimed by someone else | FAILED_PRECONDITION | `being reviewed by <admin>` |
| Publication rule fails | FAILED_PRECONDITION | `dish not publishable`, problems attached |

## Acceptance criteria

- [ ] An AI draft with `calories_incomplete` can't be approved without confirming that flag (and only if PLN-02-02
      rule 12's decision allows it).
- [ ] Approving an edit of a published dish swaps the live version atomically; a concurrent `GetDish` sees either old
      or new, never a mix.
- [ ] Two admins: the second gets "being reviewed by" on approve during the first's claim.
- [ ] The diff view shows a changed amount on one ingredient line.
- [ ] `ListDishes(statuses=[published], flags=[dietitian_pending])` returns the unchecked starter dishes and no
      household AI dish.
