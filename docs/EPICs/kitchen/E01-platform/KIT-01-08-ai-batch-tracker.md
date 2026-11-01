# KIT-01-08 AI client and batch job tracker

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-01 Platform](README.md) |
| Release | R1 |
| Size | M |
| Depends on | KIT-01-01, AI-04 (`SubmitBatch`, `GetBatch`, `GetBatchResults`) |
| Unblocks | KIT-02-08, KIT-03-05, KIT-06-05 |

## Goal

Kitchen can hand large once-for-everyone jobs (translations, matching, nutrition suggestions) to ai as Message
Batches at half price, keep track of them across restarts, and process results exactly once (README D10: ai is
stateless, the caller stores the batch id).

## Business rules

1. The ai client is built once in `main.go` (`AI_GRPC_ADDR`). Per-request calls use a 5 s deadline unless the task
   says otherwise; batch RPCs use 10 s.
2. A batch job has a `kind` from README's list (`translate_ingredients`, `pick_catalogue_matches`,
   `match_nutrition_entries`) and up to 5,000 items. Each item has a `custom_id` of the form `<kind>:<item key>`,
   unique within the batch.
3. At most one open batch per kind. A feature task that wants to submit while one is open waits for the next tick.
4. Tracker loop every 5 minutes: for each open batch call `GetBatch`; when it reports `ended`, page through
   `GetBatchResults` (500 per page) and pass each page to the kind's registered handler inside one transaction that
   also records the page number in `ai_batches.results_page`. A restart resumes after the last committed page.
5. Per item results: `succeeded` → handler; `errored`, `canceled`, `expired` → the item goes back to its feature's
   queue with `attempts + 1`; after 3 attempts the item is parked and counted (feature tasks decide what that means).
6. A batch that has not ended 26 hours after submission is marked `stuck` and alerted via `SendAlert`
   (`kitchen.ai_batch_stuck.<kind>`). Anthropic keeps results for 29 days; results not processed within 7 days of
   ending raise the same alert.
7. When ai returns `RESOURCE_EXHAUSTED` (spend limit, README D11) or `UNAVAILABLE` on submit, the tracker retries on
   the next tick; nothing is lost because items stay in their feature queue until submitted.
8. Metrics: `ai_batches_total{kind,status}`, `ai_batch_items_total{kind,result}`.
9. Admin start (ADM-03-06): `KitchenAdminService.StartCatalogueJob(kind, scope, dry_run)` runs a kind's queue now
   instead of waiting for its daily job. `scope` is `all_missing` (the same selection the daily job makes) or a list
   of up to 1,000 ingredient ids. `dry_run = true` returns the number of items it would submit and writes nothing;
   otherwise it submits one batch at once (rule 3 still applies: an open batch of that kind → `ALREADY_EXISTS`). The
   admin service turns the item count into a cost estimate and asks the founder to confirm it; kitchen doesn't
   estimate cost.
10. `ListCatalogueJobs(kind?, status?)` lists `ai_batches` rows newest first (50 per page) with kind, status, item
   count, succeeded / errored / parked counts, submitted, ended and processed times, and `started_by`
   (`schedule` or the admin user id).

## Contract

Uses `coolinary.ai.v1.AiService` `SubmitBatch`, `GetBatch`, `GetBatchResults` (AI-04).

`contracts/coolinary/kitchen/v1/admin.proto` (`KitchenAdminService`), rules 9–10:

```proto
rpc StartCatalogueJob(StartCatalogueJobRequest) returns (StartCatalogueJobResponse);
rpc ListCatalogueJobs(ListCatalogueJobsRequest) returns (ListCatalogueJobsResponse);
message StartCatalogueJobRequest { string kind = 1; bool all_missing = 2; repeated string ingredient_ids = 3; bool dry_run = 4; }
message StartCatalogueJobResponse { int32 items = 1; int64 job_id = 2; bool dry_run = 3; }
message CatalogueJob { int64 job_id = 1; string kind = 2; string status = 3; int32 items = 4; int32 succeeded = 5;
  int32 errored = 6; int32 parked = 7; string started_by = 8; string submitted_at = 9; string ended_at = 10; string processed_at = 11; }
message ListCatalogueJobsRequest { string kind = 1; string status = 2; coolinary.common.v1.PageRequest page = 3; }
message ListCatalogueJobsResponse { repeated CatalogueJob jobs = 1; coolinary.common.v1.PageResponse page = 2; }
```

## Data

`kitchen/migrations/000NN_ai_batches.sql`

```sql
CREATE TABLE ai_batches (
    id              bigserial PRIMARY KEY,
    kind            text NOT NULL,
    ai_batch_id     text NOT NULL UNIQUE,
    status          text NOT NULL CHECK (status IN ('submitted','in_progress','ended','processed','stuck','failed')),
    item_count      integer NOT NULL,
    results_page    integer NOT NULL DEFAULT 0,
    submitted_at    timestamptz NOT NULL DEFAULT now(),
    ended_at        timestamptz,
    processed_at    timestamptz,
    last_error      text,
    started_by      text NOT NULL DEFAULT 'schedule'   -- 'schedule' or the admin user id (rule 10)
);
CREATE UNIQUE INDEX ai_batches_one_open ON ai_batches (kind)
    WHERE status IN ('submitted','in_progress','ended');
```

## Errors

None returned to users. Admin RPCs: unknown `kind` → `INVALID_ARGUMENT`; an open batch of that kind →
`ALREADY_EXISTS`; more than 1,000 ids → `INVALID_ARGUMENT`.

## Acceptance criteria

- [ ] With a stub ai returning a finished batch of 1,200 results, the handler sees every item once, also when the
      process is killed after page 1.
- [ ] A second submit of the same kind while one is open is refused by the unique index and retried next tick.
- [ ] Failed items return to the feature queue with `attempts` increased.
- [ ] A batch open for 26 hours raises one alert.
- [ ] `StartCatalogueJob(kind=translate_ingredients, all_missing, dry_run)` returns the same item count the daily job
      would submit and creates no batch; the real call creates one `ai_batches` row with `started_by` = the admin id.
