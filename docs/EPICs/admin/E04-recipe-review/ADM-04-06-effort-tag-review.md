# ADM-04-06 Effort tag suggestions review

| | |
|---|---|
| Service | admin |
| Epic | [ADM-04 Recipe import and dish review](README.md) |
| Release | R1 |
| Size | S |
| Status | todo |
| Depends on | ADM-04-03; PLN-03-06 (`SuggestEffortTags`, `DecideEffortTag`), PLN-03-07 (`ListDishes`, `GetReviewItem`), PLN-03-04 (`ListDraftBatches`); AI-09 (`suggest_effort_tags`) |
| Unblocks | — |

## Goal

Effort tags on library dishes (no cooking, keeps 3 days in the fridge, freezes well, good cold the next day, base
plus toppings; features.md §2 "Calories and goals") are suggested by AI and checked by hand before use.

## Business rules

1. Start a suggestion job with `SuggestEffortTags(dish_ids | all_untagged)` (≤ 500 versions per batch, planner's
   limit), with `dry_run` / `confirm_cost_eur` as ADM-03-06 (`dry_run = true` returns the item count). The job
   list is `ListDraftBatches(kind = suggest_effort_tags)`.
2. The review list is `ListDishes(flags = [effort_tags_pending])` (published dishes included); each item opened with
   `GetReviewItem` shows the dish, its current tags, and each suggested tag with the AI's yes/no and one-line reason.
3. Action per suggested tag: accept or reject → `DecideEffortTag(dish_id, version, tag, accept)` (PLN-03-06 rule 4).
   Accepting adds the tag to that version without a new version; nothing is applied without a decision.
4. Planner keeps only suggestions that pass its own checks (`no_cook` and `base_plus_toppings` need the dish's
   structure, PLN-03-06 rule 3); a dish missing its base/topping groups never shows a `base_plus_toppings`
   suggestion. The admin adds that tag by hand in the editor (ADM-04-03) once the groups exist.
5. Bulk accept of suggestions as-is for up to 50 dishes: one `DecideEffortTag` per suggested tag; each dish gets its
   own result. Sending the same decision twice is harmless (planner treats it as no change).
6. Audited (`planner.effort_tags.decide`, dish id, version, tag, accept).
7. Role `admin` only.

## Contract

| Method | Path | Body / response |
|---|---|---|
| POST | `/admin/v1/planner/effort-tag-jobs` | `{ dish_ids?, all_untagged?, dry_run, confirm_cost_eur? }` → `{ batch_id?, items, estimated_cost_eur }` |
| GET | `/admin/v1/planner/effort-tag-reviews?page_token=` | `{ items, next_page_token }` |
| POST | `/admin/v1/planner/effort-tag-reviews/decide` | `{ items: [{ dish_id, version, tag, accept }] }` (≤ 250) → `{ results }` |

## Data

None in admin.

## Errors

| Case | HTTP |
|---|---|
| Unknown tag | 400 |
| Version retired or no longer current | per-item 409 (planner's `version_not_current`) |
| No stored suggestion for that tag | per-item 404 |

## Acceptance criteria

- [ ] Accepted tags show on the dish and the library filter "freezes well" finds it.
