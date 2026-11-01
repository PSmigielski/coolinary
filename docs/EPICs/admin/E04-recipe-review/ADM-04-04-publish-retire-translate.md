# ADM-04-04 Publish, retire and translate dishes

| | |
|---|---|
| Service | admin |
| Epic | [ADM-04 Recipe import and dish review](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | ADM-04-03; PLN-03-07 (`ApproveDish`, `RetireDish`, `RestoreDish`), PLN-03-05 (`TranslateDishes`), PLN-03-04 (`ListDraftBatches`, `GetDraftBatch`); AI-09 (`translate_dishes`) |
| Unblocks | — |

## Goal

Approved dishes become visible to users only when they are complete: calories known, all three languages present,
and (for the starter library) checked by the dietitian.

## Business rules

1. Publishing is `ApproveDish` (PLN-03-07); planner checks the publication rules and returns the unmet ones
   (PLN-02-02 rules 11–13): texts in `pl`, `de` and `en`; nutrition complete (Q086: the 200 starter dishes must
   have complete calories; a later dish may be published as "calories incomplete" only when an admin confirms it);
   derived flags computed. A dietitian check is not required to publish; the launch waits until every starter dish
   is checked (Q093, PLN-03-08 rule 3, ADM-04-05).
2. Retire: `RetireDish(reason)` hides a published dish from the library and from new plans; existing plans and
   history keep showing it. Reason required (≤ 300 characters). `RestoreDish` brings it back.
3. Translation job: `TranslateDishes(dish_ids | all_missing, target_languages)` for selected dish ids (≤ 500
   versions per batch, planner's limit) or `all_missing`, with the same `dry_run` / `confirm_cost_eur` steps as
   ADM-03-06 (`TranslateDishes(dry_run = true)` returns the item count). Planner writes accepted translations
   straight onto the dish as machine-translated texts (PLN-03-05 rule 5); the admin can correct them in the editor
   (ADM-04-03), which clears the machine-translated mark; saving without changes accepts the text and clears the
   mark too. Machine-translated dishes go live once a native speaker has spot-checked 10% of the dishes in that
   language (Q122, PLN-03-05 rule 7); the jobs list shows per language how many dishes are still marked.
4. The translation jobs list is `ListDraftBatches(kind = translate_dishes)`; one job's report (results dropped by
   planner's checks) is `GetDraftBatch`.
5. Bulk publish of up to 50 dishes: one `ApproveDish` per dish, only for dishes without warning flags (warnings must
   be confirmed one by one, ADM-04-03 rule 4); each dish gets its own result.
6. Audited: `planner.dish.publish`, `planner.dish.retire`, `planner.dish.restore`, `planner.translation_job.start`.
7. Role `admin` only.

## Contract

| Method | Path | Body / response |
|---|---|---|
| POST | `/admin/v1/planner/dishes/publish` | `{ items: [{ id, version }] }` (≤ 50) → `{ results: [{ id, ok, unmet[] }] }` |
| POST | `/admin/v1/planner/dishes/{id}/retire` | `{ reason }` → dish |
| POST | `/admin/v1/planner/dishes/{id}/restore` | → dish |
| POST | `/admin/v1/planner/translation-jobs` | `{ dish_ids?, all_missing?, target_languages, dry_run, confirm_cost_eur? }` → `{ batch_id?, items, estimated_cost_eur }` |
| GET | `/admin/v1/planner/translation-jobs?page_token=` | `{ items, next_page_token }` |

## Data

None in admin.

## Errors

| Case | HTTP |
|---|---|
| Unmet publish condition, or a dish with warning flags in a bulk publish | per-item `unmet` list (HTTP 200 for the bulk call) |
| Missing reason | 400 |
| Cost not confirmed | 409 |

## Acceptance criteria

- [ ] A dish missing its `de` texts comes back from bulk publish with `unmet` naming the missing language.
- [ ] A retired dish disappears from `GET /v1/dishes` (gateway) but still shows in a household's past plan.
