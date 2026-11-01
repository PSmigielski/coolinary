# ADM-04-02 AI dish draft batches

| | |
|---|---|
| Service | admin |
| Epic | [ADM-04 Recipe import and dish review](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | ADM-02-03; PLN-03-04 (`DraftDishes`, `GetDraftBatch`, `ListDraftBatches`); AI-04, AI-09 (`draft_dishes`) |
| Unblocks | ADM-04-03 |

## Goal

The founder asks the AI for batches of draft dishes from short briefs (cuisine, meal, tools, diet), to fill the
~200-dish starter library in five cuisines faster; every draft still goes through human review (decision 10).

## Business rules

1. A brief has planner's fields (`DishBrief`, PLN-03-04): `cuisine` (with optional sub-cuisine; one of the profile
   cuisines, features.md §1), `meal` (breakfast, dinner, supper), `tools` (the dish may use only these), `diets`
   (optional: vegetarian, vegan, gluten-free, lactose-free, nut-free), `max_total_minutes` (5–180), optional effort
   tags, `language` (`pl`, `de` or `en`; translations come later, ADM-04-04), and `count` (1–10). The admin may add
   a free note (≤ 300 characters) for the audit row only; planner has no note field.
2. A batch holds 1–50 briefs and at most 200 dishes in total (planner's limits, PLN-03-04).
3. Two steps: `dry_run = true` returns admin's cost estimate (dishes requested × `AI_JOB_EUR_PER_ITEM` for
   `draft_dishes`, as ADM-03-06 rule 3); the start needs `confirm_cost_eur` equal to it, then admin calls
   `DraftDishes(briefs)`.
4. At most 3 drafting batches open at once: admin counts unprocessed ones with
   `ListDraftBatches(kind = draft_dishes)` before starting, and refuses a fourth with 409.
5. The batch list is `ListDraftBatches(kind = draft_dishes, status)`; one batch's report is `GetDraftBatch`: per
   brief requested, received, kept, dropped by reason (unknown ingredients, tools outside the brief, ...), plus the
   unresolved ingredient names. Started by and when come from planner's batch row.
6. Drafts land in the review queue with source `ai_draft` and are labelled "AI draft" there.
7. Audited (`planner.draft_batch.start`, brief count, dish count, estimated cost).
8. Role `admin` only.

## Contract

| Method | Path | Body / response |
|---|---|---|
| POST | `/admin/v1/planner/draft-batches` | `{ briefs[], note?, dry_run, confirm_cost_eur? }` → `{ batch_id?, dishes, estimated_cost_eur }` |
| GET | `/admin/v1/planner/draft-batches?status=&page_token=` | `{ items, next_page_token }` |
| GET | `/admin/v1/planner/draft-batches/{id}` | the batch report |

## Data

None in admin.

## Errors

| Case | HTTP |
|---|---|
| Over 50 briefs or 200 dishes, count outside 1–10, invalid cuisine/meal/tool | 400 |
| Missing or stale `confirm_cost_eur` | 409 `confirm_cost` |
| 3 batches running | 409 |
| AI spend limit reached | 429 |

## Acceptance criteria

- [ ] A batch "Polish, dinner, hob + oven, 10 dishes" ends with up to 10 drafts in the review queue, each using only
      hob and oven.
- [ ] With 3 unprocessed drafting batches, starting a fourth returns 409 and calls nothing in planner.
