# ADM-03-06 Catalogue batch jobs

| | |
|---|---|
| Service | admin |
| Epic | [ADM-03 Catalogue, matching and nutrition review](README.md) |
| Release | R1 |
| Size | S |
| Status | todo |
| Depends on | ADM-02-03; KIT-01-08 (`KitchenAdminService.StartCatalogueJob`, `ListCatalogueJobs`); AI-02-04 (`AiAdminService.GetSpendStatus`) |
| Unblocks | ADM-03-02, ADM-03-03 (their queues are filled by these jobs) |

## Goal

The founder starts the once-for-everyone AI jobs on the catalogue (translation, nutrition matching, bulk name
matching) and sees their cost before and after, because these jobs are a fixed cost the business model caps.

## Business rules

1. Kitchen runs its batch jobs daily on its own (KIT-02-08, KIT-03-05, KIT-06-05). This page starts one now. Job
   kinds: `translate_ingredients`, `match_nutrition_entries`, `pick_catalogue_matches` (README batch job kinds).
2. Scope: `all_missing` (the same selection kitchen's daily job makes) or, for `translate_ingredients` and
   `match_nutrition_entries`, a list of up to 1,000 ingredient ids. `pick_catalogue_matches` takes `all_missing`
   only.
3. Two steps: `dry_run = true` calls `StartCatalogueJob(..., dry_run = true)`, which returns the item count; admin
   turns it into an estimate with its config `AI_JOB_EUR_PER_ITEM` (one value per kind, batch prices, seeded from
   the AI-03 eval runs). The real start must send `confirm_cost_eur` equal to that estimate (rounded to cents); if a
   fresh dry run gives a different estimate → 409 with the new one.
4. At most one open batch per kind (KIT-01-08 rule 3): kitchen's `ALREADY_EXISTS` → 409 `job_running`.
5. The jobs list is `ListCatalogueJobs(kind, status)`: kind, status, item counts (total, succeeded, errored,
   parked), started by (`schedule` or an admin), submitted, ended and processed times. Kitchen stores the job state
   (D10: `ai` has no database). Cost per job is not tracked; the page shows the month's batch spend and limit from
   `AiAdminService.GetSpendStatus` (bucket `batch`).
6. Audited (`kitchen.job.start`, kind, scope size, estimated cost).
7. Role `admin` only.

## Contract

| Method | Path | Body / response |
|---|---|---|
| POST | `/admin/v1/kitchen/jobs` | `{ kind, scope: all_missing\|ids, ingredient_ids?, dry_run, confirm_cost_eur? }` → `{ job_id?, items, estimated_cost_eur }` |
| GET | `/admin/v1/kitchen/jobs?kind=&status=&page_token=` | `{ items, next_page_token, batch_spend: { month_eur, limit_eur } }` |

## Data

None in admin.

## Errors

| Case | HTTP |
|---|---|
| Missing or wrong `confirm_cost_eur` | 409 `confirm_cost` with the current estimate |
| A job of this kind is running | 409 `job_running` |
| `ai` over its daily spend limit | 429 (from `RESOURCE_EXHAUSTED`) |

## Acceptance criteria

- [ ] Dry run then start with the same estimate starts the job; start without the dry run → 409.
- [ ] A job started here shows `started_by` = the admin, and the month's batch spend rises once its results are
      processed.
