# planner — epics and tasks

Go (`planner/`). Owns the recipe library (decision 2): dishes, import and AI drafting, review, calories per dish, browse and "what can I make now", the rule checker on every AI dish, day, family and week plans, swaps and rewrites, cooking mode and "cooked it", per-person portions, the shopping list with live notices, base pantries, offer badges and reminders, plan from deals, the trip price matrix and the cheapest top-up. It consumes events and publishes none.

What exists today: `PlannerService.SuggestDishes` over a skeleton `dishes` table with five example rows, replaced by PLN-02.

PLN-01 to PLN-11 and PLN-12 to PLN-21 were written by two writers; dependencies between the halves use epic ids.

## Epics

| Epic | Title | Rel. | Tasks |
|---|---|---|---|
| [PLN-01](E01-platform/README.md) | Platform: events, clients, snapshots, plan lookup | R1 | 6 |
| [PLN-02](E02-recipe-library-model/README.md) | Recipe library model and storage | R1 | 6 |
| [PLN-03](E03-library-import-and-review/README.md) | Library import, AI drafting and review | R1 | 9 |
| [PLN-04](E04-dish-nutrition/README.md) | Dish nutrition and calories | R1 | 4 |
| [PLN-05](E05-library-browse/README.md) | Library browse, favourites, ratings, history | R1 | 6 |
| [PLN-06](E06-what-can-i-make-now/README.md) | What can I make now | R1 | 3 |
| [PLN-07](E07-rule-checker/README.md) | Rule checker for AI and library dishes | R1 | 5 |
| [PLN-08](E08-ai-day-plan/README.md) | AI day plan and cuisine choice | R1 | 8 |
| [PLN-09](E09-swaps-and-rewrites/README.md) | Swaps, wishes, tool rewrite, restyle | R1–R3 | 5 |
| [PLN-10](E10-cooking-mode/README.md) | Cooking mode and "cooked it" | R1 | 3 |
| [PLN-11](E11-library-growth/README.md) | Library growth from well-rated AI dishes | R2 | 3 |
| [PLN-12](E12-per-person-portions/README.md) | Per-person portions (one dish, one pot) | R2 | 8 |
| [PLN-13](E13-family-day-plan/README.md) | Family day plan | R2 | 7 |
| [PLN-14](E14-week-plan/README.md) | Week plan with leftovers and fewest cooking sessions | R3 | 7 |
| [PLN-15](E15-shopping-list/README.md) | Shopping list | R1–R2 | 10 |
| [PLN-16](E16-cuisine-base-pantry/README.md) | Cuisine base pantry | R1 | 3 |
| [PLN-17](E17-offer-badges-and-reminders/README.md) | Offer badges and offer-ends reminders | R2 | 4 |
| [PLN-18](E18-plan-from-deals/README.md) | Plan from deals | R2 | 4 |
| [PLN-19](E19-trip-price-matrix/README.md) | Trip price matrix | R2–R3 | 4 |
| [PLN-20](E20-cheapest-top-up/README.md) | Cheapest top-up | R3 | 4 |
| [PLN-21](E21-household-lifecycle/README.md) | Household lifecycle: event consumers and data export | R1–R2 | 6 |

115 tasks: 68 in R1, 34 in R2, 13 in R3.

## What other services use

RPCs defined in this service's tasks (first task that defines each). Admin RPCs are called only by admin.

| RPC | Defined in | Admin |
|---|---|---|
| `ListAdminChanges` | [PLN-01-06](E01-platform/PLN-01-06-admin-service.md) | yes |
| `GetQueueCounts` | [PLN-01-06](E01-platform/PLN-01-06-admin-service.md) | yes |
| `GetDish` | [PLN-02-06](E02-recipe-library-model/PLN-02-06-get-dish.md) |  |
| `BatchGetDishes` | [PLN-02-06](E02-recipe-library-model/PLN-02-06-get-dish.md) |  |
| `ImportDishes` | [PLN-03-03](E03-library-import-and-review/PLN-03-03-import-dishes-rpc.md) | yes |
| `DraftDishes` | [PLN-03-04](E03-library-import-and-review/PLN-03-04-ai-drafting-batch.md) | yes |
| `GetDraftBatch` | [PLN-03-04](E03-library-import-and-review/PLN-03-04-ai-drafting-batch.md) | yes |
| `ListDraftBatches` | [PLN-03-04](E03-library-import-and-review/PLN-03-04-ai-drafting-batch.md) | yes |
| `TranslateDishes` | [PLN-03-05](E03-library-import-and-review/PLN-03-05-dish-translation-batch.md) | yes |
| `SuggestEffortTags` | [PLN-03-06](E03-library-import-and-review/PLN-03-06-effort-tag-suggestions.md) | yes |
| `DecideEffortTag` | [PLN-03-06](E03-library-import-and-review/PLN-03-06-effort-tag-suggestions.md) | yes |
| `ListReviewQueue` | [PLN-03-07](E03-library-import-and-review/PLN-03-07-review-queue.md) | yes |
| `GetReviewItem` | [PLN-03-07](E03-library-import-and-review/PLN-03-07-review-queue.md) | yes |
| `UpdateDraft` | [PLN-03-07](E03-library-import-and-review/PLN-03-07-review-queue.md) | yes |
| `ApproveDish` | [PLN-03-07](E03-library-import-and-review/PLN-03-07-review-queue.md) | yes |
| `RejectDish` | [PLN-03-07](E03-library-import-and-review/PLN-03-07-review-queue.md) | yes |
| `RetireDish` | [PLN-03-07](E03-library-import-and-review/PLN-03-07-review-queue.md) | yes |
| `RestoreDish` | [PLN-03-07](E03-library-import-and-review/PLN-03-07-review-queue.md) | yes |
| `ListDishes` | [PLN-03-07](E03-library-import-and-review/PLN-03-07-review-queue.md) | yes |
| `RecordDietitianCheck` | [PLN-03-08](E03-library-import-and-review/PLN-03-08-dietitian-sign-off.md) | yes |
| `GetDietitianReport` | [PLN-03-08](E03-library-import-and-review/PLN-03-08-dietitian-sign-off.md) | yes |
| `GetLibraryCoverage` | [PLN-03-09](E03-library-import-and-review/PLN-03-09-library-coverage-report.md) | yes |
| `SearchDishes` | [PLN-05-01](E05-library-browse/PLN-05-01-search-dishes-rpc.md) |  |
| `SetFavourite` | [PLN-05-04](E05-library-browse/PLN-05-04-favourites.md) |  |
| `RateDish` | [PLN-05-05](E05-library-browse/PLN-05-05-ratings.md) |  |
| `ListCookedHistory` | [PLN-05-06](E05-library-browse/PLN-05-06-cooked-history.md) |  |
| `GetPlans` | [PLN-08-01](E08-ai-day-plan/PLN-08-01-plan-storage.md) |  |
| `AcceptPlan` | [PLN-08-01](E08-ai-day-plan/PLN-08-01-plan-storage.md) |  |
| `DeletePlan` | [PLN-08-01](E08-ai-day-plan/PLN-08-01-plan-storage.md) |  |
| `ChooseMealOption` | [PLN-08-06](E08-ai-day-plan/PLN-08-06-cuisine-choices.md) |  |
| `ComposeDayPlan` | [PLN-08-07](E08-ai-day-plan/PLN-08-07-compose-day-plan-rpc.md) |  |
| `SetPlanMeal` | [PLN-08-08](E08-ai-day-plan/PLN-08-08-set-plan-meal.md) |  |
| `SwapMeal` | [PLN-09-01](E09-swaps-and-rewrites/PLN-09-01-swap-meal.md) |  |
| `ReplaceIngredient` | [PLN-09-02](E09-swaps-and-rewrites/PLN-09-02-replace-ingredient.md) |  |
| `WishMeal` | [PLN-09-03](E09-swaps-and-rewrites/PLN-09-03-free-text-wish.md) |  |
| `RewriteDish` | [PLN-09-04](E09-swaps-and-rewrites/PLN-09-04-tool-rewrite.md) |  |
| `GetCookingSteps` | [PLN-10-01](E10-cooking-mode/PLN-10-01-cooking-steps.md) |  |
| `MarkCooked` | [PLN-10-02](E10-cooking-mode/PLN-10-02-mark-cooked.md) |  |
| `DismissRatingPrompt` | [PLN-10-03](E10-cooking-mode/PLN-10-03-rating-prompt.md) |  |
| `PromoteCandidate` | [PLN-11-02](E11-library-growth/PLN-11-02-anonymised-copy.md) | yes |
| `ListPromotionCandidates` | [PLN-11-03](E11-library-growth/PLN-11-03-promotion-admin-and-metrics.md) | yes |
| `DismissCandidate` | [PLN-11-03](E11-library-growth/PLN-11-03-promotion-admin-and-metrics.md) | yes |
| `SetMealAttendance` | [PLN-12-01](E12-per-person-portions/PLN-12-01-meal-attendance.md) |  |
| `ListMealAttendance` | [PLN-12-01](E12-per-person-portions/PLN-12-01-meal-attendance.md) |  |
| `SetAttendanceRule` | [PLN-12-01](E12-per-person-portions/PLN-12-01-meal-attendance.md) |  |
| `DeleteAttendanceRule` | [PLN-12-01](E12-per-person-portions/PLN-12-01-meal-attendance.md) |  |
| `ListNoCookExtras` | [PLN-12-04](E12-per-person-portions/PLN-12-04-no-cook-extras.md) | yes |
| `UpsertNoCookExtra` | [PLN-12-04](E12-per-person-portions/PLN-12-04-no-cook-extras.md) | yes |
| `GetMealPortions` | [PLN-12-06](E12-per-person-portions/PLN-12-06-portion-allocation-for-planned-meals.md) |  |
| `GetScaledIngredients` | [PLN-12-08](E12-per-person-portions/PLN-12-08-scaled-ingredients.md) |  |
| `CreateFamilyDayPlan` | [PLN-13-07](E13-family-day-plan/PLN-13-07-ai-family-day.md) |  |
| `GetWeekPlan` | [PLN-14-01](E14-week-plan/PLN-14-01-week-plan-model.md) |  |
| `ReplaceWeekMeal` | [PLN-14-01](E14-week-plan/PLN-14-01-week-plan-model.md) |  |
| `DeleteWeekPlan` | [PLN-14-01](E14-week-plan/PLN-14-01-week-plan-model.md) |  |
| `CreateWeekPlan` | [PLN-14-07](E14-week-plan/PLN-14-07-ai-week.md) |  |
| `GetList` | [PLN-15-01](E15-shopping-list/PLN-15-01-list-model.md) |  |
| `GenerateList` | [PLN-15-05](E15-shopping-list/PLN-15-05-generate-review-accept.md) |  |
| `AcceptListProposal` | [PLN-15-05](E15-shopping-list/PLN-15-05-generate-review-accept.md) |  |
| `AddListLine` | [PLN-15-06](E15-shopping-list/PLN-15-06-manual-lines-and-editing.md) |  |
| `UpdateListLine` | [PLN-15-06](E15-shopping-list/PLN-15-06-manual-lines-and-editing.md) |  |
| `RemoveListLine` | [PLN-15-06](E15-shopping-list/PLN-15-06-manual-lines-and-editing.md) |  |
| `ClearTickedLines` | [PLN-15-06](E15-shopping-list/PLN-15-06-manual-lines-and-editing.md) |  |
| `TickListLine` | [PLN-15-07](E15-shopping-list/PLN-15-07-tick-into-fridge.md) |  |
| `UntickListLine` | [PLN-15-07](E15-shopping-list/PLN-15-07-tick-into-fridge.md) |  |
| `GetListChanges` | [PLN-15-09](E15-shopping-list/PLN-15-09-live-notices-and-changes.md) |  |
| `ListBasePantries` | [PLN-16-01](E16-cuisine-base-pantry/PLN-16-01-base-pantry-lists.md) | yes |
| `SetBasePantry` | [PLN-16-01](E16-cuisine-base-pantry/PLN-16-01-base-pantry-lists.md) | yes |
| `GetBasePantry` | [PLN-16-02](E16-cuisine-base-pantry/PLN-16-02-base-pantry-checklist.md) |  |
| `AddBasePantryToList` | [PLN-16-03](E16-cuisine-base-pantry/PLN-16-03-add-missing-base-items.md) |  |
| `AddOfferToList` | [PLN-17-02](E17-offer-badges-and-reminders/PLN-17-02-add-offer-to-list.md) |  |
| `CreatePlanFromDeals` | [PLN-18-02](E18-plan-from-deals/PLN-18-02-ai-plan-from-deals.md) |  |
| `GetTripPriceMatrix` | [PLN-19-03](E19-trip-price-matrix/PLN-19-03-price-matrix-rpc.md) |  |
| `CreateCheapestTopUp` | [PLN-20-02](E20-cheapest-top-up/PLN-20-02-ai-cheapest-top-up.md) |  |
| `ApplyTopUp` | [PLN-20-04](E20-cheapest-top-up/PLN-20-04-apply-top-up.md) |  |
| `ExportHouseholdData` | [PLN-21-06](E21-household-lifecycle/PLN-21-06-export-household-data.md) |  |

## Build order

Task order follows [ROADMAP.md](../ROADMAP.md); `#` is the task's number there and `Milestone` its milestone.

### Release 1

| # | Task | Title | Size | Milestone |
|---|---|---|---|---|
| 8 | [PLN-01-01](E01-platform/PLN-01-01-event-consumer.md) | Event consumer (Redis Streams, Go) | M | M1 |
| 9 | [PLN-01-02](E01-platform/PLN-01-02-service-clients.md) | Clients for kitchen, offers, ai, notifier, identity | M | M1 |
| 10 | [PLN-01-06](E01-platform/PLN-01-06-admin-service.md) | PlannerAdminService scaffold and change log | S | M1 |
| 75 | [PLN-01-05](E01-platform/PLN-01-05-household-plan-copy.md) | Household plan lookup from identity | S | M2 |
| 126 | [PLN-01-03](E01-platform/PLN-01-03-profile-snapshot.md) | Planning-profile snapshot cache | M | M4 |
| 166 | [PLN-01-04](E01-platform/PLN-01-04-stock-snapshot.md) | Stock snapshot cache | S | M5 |
| 187 | [PLN-02-01](E02-recipe-library-model/PLN-02-01-dish-schema.md) | Dish schema (replaces skeleton `dishes`) | M | M6 |
| 188 | [PLN-02-02](E02-recipe-library-model/PLN-02-02-dish-validation.md) | Dish validation rules | M | M6 |
| 189 | [PLN-02-03](E02-recipe-library-model/PLN-02-03-derived-flags.md) | Derived diet and allergen flags | M | M6 |
| 190 | [PLN-02-04](E02-recipe-library-model/PLN-02-04-dish-lifecycle.md) | Dish lifecycle and versions | M | M6 |
| 191 | [PLN-02-05](E02-recipe-library-model/PLN-02-05-portion-scaling.md) | Portion scaling rules | S | M6 |
| 192 | [PLN-02-06](E02-recipe-library-model/PLN-02-06-get-dish.md) | GetDish and BatchGetDishes | M | M6 |
| 193 | [PLN-03-01](E03-library-import-and-review/PLN-03-01-import-file-format.md) | Import file format and parser | M | M6 |
| 194 | [PLN-03-02](E03-library-import-and-review/PLN-03-02-ingredient-name-resolution.md) | Ingredient name resolution | M | M6 |
| 195 | [PLN-03-03](E03-library-import-and-review/PLN-03-03-import-dishes-rpc.md) | ImportDishes admin RPC | M | M6 |
| 196 | [PLN-04-01](E04-dish-nutrition/PLN-04-01-compute-dish-nutrition.md) | Compute and store nutrition per dish version | M | M6 |
| 197 | [PLN-03-07](E03-library-import-and-review/PLN-03-07-review-queue.md) | Review queue and decisions | L | M6 |
| 198 | [PLN-03-08](E03-library-import-and-review/PLN-03-08-dietitian-sign-off.md) | Dietitian sign-off | S | M6 |
| 199 | [PLN-04-02](E04-dish-nutrition/PLN-04-02-recompute-on-catalogue-events.md) | Recompute on catalogue events | M | M6 |
| 200 | [PLN-04-03](E04-dish-nutrition/PLN-04-03-nutrition-in-responses.md) | Nutrition in dish responses (kcal/kJ, rounding, incomplete) | S | M6 |
| 201 | [PLN-04-04](E04-dish-nutrition/PLN-04-04-scanned-product-override.md) | Pack values of the household's scanned products | M | M6 |
| 202 | [PLN-05-01](E05-library-browse/PLN-05-01-search-dishes-rpc.md) | SearchDishes RPC (text search, scopes, sorting, paging) | M | M6 |
| 203 | [PLN-05-02](E05-library-browse/PLN-05-02-hard-limits.md) | Hard limits from the kitchen profile | M | M6 |
| 204 | [PLN-03-09](E03-library-import-and-review/PLN-03-09-library-coverage-report.md) | Library coverage report | M | M6 |
| 205 | [PLN-05-03](E05-library-browse/PLN-05-03-filters-and-facets.md) | Filters and facet counts | M | M6 |
| 206 | [PLN-05-04](E05-library-browse/PLN-05-04-favourites.md) | Favourites | S | M6 |
| 207 | [PLN-05-05](E05-library-browse/PLN-05-05-ratings.md) | Ratings | S | M6 |
| 208 | [PLN-05-06](E05-library-browse/PLN-05-06-cooked-history.md) | Cooked history | S | M6 |
| 209 | [PLN-06-01](E06-what-can-i-make-now/PLN-06-01-stock-coverage.md) | Stock coverage per dish | M | M6 |
| 210 | [PLN-06-02](E06-what-can-i-make-now/PLN-06-02-ranking.md) | Ranking: share on hand, expiry, tie-breaks | M | M6 |
| 211 | [PLN-06-03](E06-what-can-i-make-now/PLN-06-03-suggest-dishes-rpc.md) | Extend SuggestDishes | M | M6 |
| 212 | [PLN-07-01](E07-rule-checker/PLN-07-01-checker-core.md) | Checker core, candidate dish and verdict | M | M6 |
| 213 | [PLN-07-02](E07-rule-checker/PLN-07-02-tool-check.md) | Tool check | S | M6 |
| 214 | [PLN-07-03](E07-rule-checker/PLN-07-03-allergen-diet-dislike-check.md) | Allergen, diet and dislike check per eater | M | M6 |
| 215 | [PLN-07-04](E07-rule-checker/PLN-07-04-stock-check.md) | Stock and allowed-ingredient check | M | M6 |
| 216 | [PLN-07-05](E07-rule-checker/PLN-07-05-ai-output-checks.md) | AI output checks, time budget, metrics, 50-dish harness | M | M6 |
| 244 | [PLN-08-01](E08-ai-day-plan/PLN-08-01-plan-storage.md) | Plan storage, accept and read | M | M8 |
| 245 | [PLN-08-03](E08-ai-day-plan/PLN-08-03-plan-allocation.md) | Ingredient allocation across meals | M | M8 |
| 246 | [PLN-08-05](E08-ai-day-plan/PLN-08-05-library-fallback.md) | Library day plan (non-AI fallback) | M | M8 |
| 247 | [PLN-08-08](E08-ai-day-plan/PLN-08-08-set-plan-meal.md) | SetPlanMeal: put a chosen dish into a plan meal | M | M8 |
| 249 | [PLN-08-02](E08-ai-day-plan/PLN-08-02-compose-request.md) | Compose request assembly | M | M8 |
| 250 | [PLN-08-04](E08-ai-day-plan/PLN-08-04-validate-repair-store.md) | Validate, repair once, store AI dishes | L | M8 |
| 251 | [PLN-08-06](E08-ai-day-plan/PLN-08-06-cuisine-choices.md) | Cuisine mismatch: three choices | M | M8 |
| 252 | [PLN-08-07](E08-ai-day-plan/PLN-08-07-compose-day-plan-rpc.md) | ComposeDayPlan RPC, quota signal, idempotency | M | M8 |
| 254 | [PLN-03-04](E03-library-import-and-review/PLN-03-04-ai-drafting-batch.md) | AI drafting batch (`draft_dishes`) | L | M8 |
| 256 | [PLN-03-05](E03-library-import-and-review/PLN-03-05-dish-translation-batch.md) | Dish translation batch (`translate_dishes`) | M | M8 |
| 258 | [PLN-03-06](E03-library-import-and-review/PLN-03-06-effort-tag-suggestions.md) | Effort tag suggestions (`suggest_effort_tags`) | S | M8 |
| 266 | [PLN-09-01](E09-swaps-and-rewrites/PLN-09-01-swap-meal.md) | Swap one meal | M | M8 |
| 267 | [PLN-09-02](E09-swaps-and-rewrites/PLN-09-02-replace-ingredient.md) | Replace a missing ingredient | M | M8 |
| 268 | [PLN-09-03](E09-swaps-and-rewrites/PLN-09-03-free-text-wish.md) | Free-text wish for a meal | M | M8 |
| 280 | [PLN-10-01](E10-cooking-mode/PLN-10-01-cooking-steps.md) | Cooking mode payload | S | M9 |
| 281 | [PLN-10-02](E10-cooking-mode/PLN-10-02-mark-cooked.md) | "Cooked it": consume stock, history, plan status | M | M9 |
| 282 | [PLN-10-03](E10-cooking-mode/PLN-10-03-rating-prompt.md) | Rating prompt after cooking | S | M9 |
| 283 | [PLN-15-01](E15-shopping-list/PLN-15-01-list-model.md) | List model and GetList | M | M9 |
| 284 | [PLN-15-02](E15-shopping-list/PLN-15-02-needs-from-plan.md) | Ingredient needs from the plan | M | M9 |
| 285 | [PLN-15-03](E15-shopping-list/PLN-15-03-subtract-stock-and-staples.md) | Subtract stock and staples | M | M9 |
| 286 | [PLN-15-04](E15-shopping-list/PLN-15-04-merge-and-round-to-packs.md) | Merge lines and round to packs | M | M9 |
| 287 | [PLN-15-05](E15-shopping-list/PLN-15-05-generate-review-accept.md) | Generate, review and accept proposed lines | M | M9 |
| 288 | [PLN-15-08](E15-shopping-list/PLN-15-08-versions-idempotency-conflicts.md) | Versions, idempotent offline edits, conflict rules | M | M9 |
| 289 | [PLN-15-06](E15-shopping-list/PLN-15-06-manual-lines-and-editing.md) | Manual lines, editing, removing, clearing | M | M9 |
| 290 | [PLN-15-07](E15-shopping-list/PLN-15-07-tick-into-fridge.md) | Tick into the fridge (and undo) | M | M9 |
| 291 | [PLN-15-09](E15-shopping-list/PLN-15-09-live-notices-and-changes.md) | Live change notices and GetListChanges | S | M9 |
| 292 | [PLN-16-01](E16-cuisine-base-pantry/PLN-16-01-base-pantry-lists.md) | Base pantry lists and admin edits | S | M9 |
| 293 | [PLN-16-02](E16-cuisine-base-pantry/PLN-16-02-base-pantry-checklist.md) | Base pantry checklist for a household | S | M9 |
| 294 | [PLN-16-03](E16-cuisine-base-pantry/PLN-16-03-add-missing-base-items.md) | Add missing base items to the list | S | M9 |
| 339 | [PLN-21-01](E21-household-lifecycle/PLN-21-01-household-and-account-deletion.md) | Household and account deletion | M | M11 |
| 340 | [PLN-21-04](E21-household-lifecycle/PLN-21-04-ingredient-merge-for-lists.md) | Ingredient merge for lists, pantries and extras | S | M11 |
| 341 | [PLN-21-06](E21-household-lifecycle/PLN-21-06-export-household-data.md) | ExportHouseholdData | M | M11 |

### Release 2

| # | Task | Title | Size | Milestone |
|---|---|---|---|---|
| 402 | [PLN-19-01](E19-trip-price-matrix/PLN-19-01-candidate-shops.md) | Candidate shops for a trip | S | M14 |
| 423 | [PLN-17-01](E17-offer-badges-and-reminders/PLN-17-01-offer-badges.md) | Offer badges on list lines | M | M14 |
| 424 | [PLN-17-02](E17-offer-badges-and-reminders/PLN-17-02-add-offer-to-list.md) | Add an offer to the list | S | M14 |
| 425 | [PLN-17-04](E17-offer-badges-and-reminders/PLN-17-04-offer-ends-reminders.md) | "Offer ends tomorrow" reminders | M | M14 |
| 426 | [PLN-19-02](E19-trip-price-matrix/PLN-19-02-price-per-line-and-shop.md) | Price per line and shop | M | M14 |
| 427 | [PLN-19-03](E19-trip-price-matrix/PLN-19-03-price-matrix-rpc.md) | Price matrix RPC and shop totals | M | M14 |
| 469 | [PLN-12-01](E12-per-person-portions/PLN-12-01-meal-attendance.md) | Meal attendance (marking people away) | M | M16 |
| 470 | [PLN-12-02](E12-per-person-portions/PLN-12-02-meal-shares.md) | Meal shares per person | S | M16 |
| 471 | [PLN-12-03](E12-per-person-portions/PLN-12-03-portion-and-pot-calculator.md) | Portion and pot calculator | M | M16 |
| 472 | [PLN-12-04](E12-per-person-portions/PLN-12-04-no-cook-extras.md) | No-cook extras | M | M16 |
| 473 | [PLN-12-05](E12-per-person-portions/PLN-12-05-carry-over-and-day-balance.md) | Carry-over and day balance | M | M16 |
| 474 | [PLN-12-06](E12-per-person-portions/PLN-12-06-portion-allocation-for-planned-meals.md) | Portion allocation for planned meals | M | M16 |
| 475 | [PLN-12-07](E12-per-person-portions/PLN-12-07-calorie-visibility.md) | Who sees which numbers | S | M16 |
| 476 | [PLN-12-08](E12-per-person-portions/PLN-12-08-scaled-ingredients.md) | Scaled ingredients for a planned meal | S | M16 |
| 477 | [PLN-13-01](E13-family-day-plan/PLN-13-01-shared-meal-limits.md) | Hard limits for a shared meal | S | M16 |
| 478 | [PLN-13-02](E13-family-day-plan/PLN-13-02-base-plus-toppings.md) | Base-plus-toppings plates | M | M16 |
| 479 | [PLN-13-03](E13-family-day-plan/PLN-13-03-plan-comparator.md) | Plan comparator (priority order) | M | M16 |
| 480 | [PLN-13-04](E13-family-day-plan/PLN-13-04-eat-lighter-as-a-family.md) | Eat lighter as a family | S | M16 |
| 481 | [PLN-13-05](E13-family-day-plan/PLN-13-05-second-dish-when-forced.md) | Second dish only when an allergy or diet forces it | S | M16 |
| 482 | [PLN-13-06](E13-family-day-plan/PLN-13-06-non-ai-family-day.md) | Family day without AI | L | M16 |
| 483 | [PLN-15-10](E15-shopping-list/PLN-15-10-list-follows-plan.md) | List follows plan changes | M | M16 |
| 484 | [PLN-18-01](E18-plan-from-deals/PLN-18-01-deal-candidates.md) | Deal candidates for a household | M | M16 |
| 485 | [PLN-18-02](E18-plan-from-deals/PLN-18-02-ai-plan-from-deals.md) | Plan from deals with AI | L | M16 |
| 486 | [PLN-18-03](E18-plan-from-deals/PLN-18-03-non-ai-plan-from-deals.md) | Plan from deals without AI | M | M16 |
| 487 | [PLN-18-04](E18-plan-from-deals/PLN-18-04-deals-plan-to-list.md) | Deals plan to the shopping list | S | M16 |
| 489 | [PLN-13-07](E13-family-day-plan/PLN-13-07-ai-family-day.md) | Family day with AI | L | M16 |
| 500 | [PLN-11-01](E11-library-growth/PLN-11-01-candidate-selection.md) | Promotion candidate job | M | M17 |
| 501 | [PLN-11-02](E11-library-growth/PLN-11-02-anonymised-copy.md) | Anonymised copy into review | M | M17 |
| 502 | [PLN-11-03](E11-library-growth/PLN-11-03-promotion-admin-and-metrics.md) | Promotion admin RPCs and savings metric | S | M17 |
| 503 | [PLN-21-02](E21-household-lifecycle/PLN-21-02-member-added-and-removed.md) | Member added and removed | S | M17 |
| 504 | [PLN-21-03](E21-household-lifecycle/PLN-21-03-plan-downgrade.md) | Plus → Free downgrade | S | M17 |
| 505 | [PLN-21-05](E21-household-lifecycle/PLN-21-05-offers-events.md) | Offers events | S | M17 |
| 506 | [PLN-17-03](E17-offer-badges-and-reminders/PLN-17-03-badge-refresh-and-expiry.md) | Badge refresh and expiry | S | M17 |
| 508 | [PLN-09-04](E09-swaps-and-rewrites/PLN-09-04-tool-rewrite.md) | Rewrite a recipe for my tools | M | M17 |

### Release 3

| # | Task | Title | Size | Milestone |
|---|---|---|---|---|
| 512 | [PLN-14-01](E14-week-plan/PLN-14-01-week-plan-model.md) | Week plan model, read and edit | M | M18 |
| 513 | [PLN-14-02](E14-week-plan/PLN-14-02-cooking-days.md) | Cooking days from the rhythm | S | M18 |
| 514 | [PLN-14-03](E14-week-plan/PLN-14-03-cook-once-eat-twice.md) | Cook once, eat twice | M | M18 |
| 515 | [PLN-14-04](E14-week-plan/PLN-14-04-effort-summary.md) | Effort summary (sessions, stove time) | S | M18 |
| 516 | [PLN-14-05](E14-week-plan/PLN-14-05-week-balance.md) | Week balance | M | M18 |
| 517 | [PLN-14-06](E14-week-plan/PLN-14-06-non-ai-week.md) | Week plan without AI | L | M18 |
| 519 | [PLN-14-07](E14-week-plan/PLN-14-07-ai-week.md) | Week plan with AI | L | M18 |
| 527 | [PLN-19-04](E19-trip-price-matrix/PLN-19-04-receipt-prices.md) | Receipt prices in the matrix | S | M19 |
| 537 | [PLN-09-05](E09-swaps-and-rewrites/PLN-09-05-restyle.md) | Restyle a dish ("make this Italian-style") | M | M20 |
| 538 | [PLN-20-01](E20-cheapest-top-up/PLN-20-01-needs-and-allowed-swaps.md) | Top-up needs and allowed swaps | M | M20 |
| 539 | [PLN-20-03](E20-cheapest-top-up/PLN-20-03-non-ai-top-up.md) | Top-up without AI | S | M20 |
| 543 | [PLN-20-02](E20-cheapest-top-up/PLN-20-02-ai-cheapest-top-up.md) | Cheapest top-up with AI | L | M20 |
| 544 | [PLN-20-04](E20-cheapest-top-up/PLN-20-04-apply-top-up.md) | Applying a top-up | M | M20 |

## Decisions

Settled answers (full list in [DECISIONS.md](../DECISIONS.md)); each task states its rule as `Decided (Qnnn)`.

| Task | Decision |
|---|---|
| [PLN-01-05](E01-platform/PLN-01-05-household-plan-copy.md) | Q004 — planner asks identity each time instead of keeping a local copy fed by `household.plan_changed`. This adds the call planner → identity (no cycle: identity calls only notifier). |
| [PLN-01-05](E01-platform/PLN-01-05-household-plan-copy.md) | Q229 — a request falls back to the `x-plan` metadata from the signed-in user's token (at most 15 minutes old); when that is missing too, the household is treated as Free. Counted in `planner_plan_lookup_failures_total{context="request"}`. Requests that work … |
| [PLN-01-06](E01-platform/PLN-01-06-admin-service.md) | Q016 — . |
| [PLN-02-02](E02-recipe-library-model/PLN-02-02-dish-validation.md) | Q087 — a dish needing only a microwave, kettle, toaster or blender counts as "no cooking" (features.md: breakfast and the non-main meal come from "no cooking" dishes; the student preset is microwave + kettle). |
| [PLN-02-02](E02-recipe-library-model/PLN-02-02-dish-validation.md) | Q086 — the 200 starter dishes must have complete calories to be published; later dishes may be published with "calories incomplete" only after an admin confirms it. |
| [PLN-02-03](E02-recipe-library-model/PLN-02-03-derived-flags.md) | Q054 — kitchen publishes `catalogue.ingredient_flags_changed` with `ingredient_ids` (KIT-02-10), and planner re-derives the flags of every dish using those ingredients at once; the daily job stays as a safety net. |
| [PLN-02-05](E02-recipe-library-model/PLN-02-05-portion-scaling.md) | Q089 — kitchen's `GetIngredients` returns a `whole_units_only` flag (KIT-02) that drives this rule. |
| [PLN-02-05](E02-recipe-library-model/PLN-02-05-portion-scaling.md) | Q088 — when servings > 2 × base, the response carries the hint "cooking times may be longer". |
| [PLN-03-02](E03-library-import-and-review/PLN-03-02-ingredient-name-resolution.md) | Q091 — accept ≥ 0.90, review 0.60–0.90, reject < 0.60; the thresholds are config values, tuned after the first 200 dishes. |
| [PLN-03-02](E03-library-import-and-review/PLN-03-02-ingredient-name-resolution.md) | Q090 — kitchen (KIT-03) accepts the context value `recipe` next to product, receipt and offer. |
| [PLN-03-05](E03-library-import-and-review/PLN-03-05-dish-translation-batch.md) | Q122 — dishes may be published with machine-translated texts after a native speaker spot-checks 10% of dishes per language, before each language goes live (business model "one-off per language"). |
| [PLN-03-08](E03-library-import-and-review/PLN-03-08-dietitian-sign-off.md) | Q093 — checking is not required to publish, but launch waits until every starter dish is checked (rule 4). |
| [PLN-03-08](E03-library-import-and-review/PLN-03-08-dietitian-sign-off.md) | Q094 — 95%, with every starter dish checked (features.md's 20-dish check, Verification 16, is part of it). |
| [PLN-03-09](E03-library-import-and-review/PLN-03-09-library-coverage-report.md) | Q060 — kitchen (KIT-10) owns the preset lists and the report reads them from kitchen; "basic kitchen" = hob, oven, kettle, toaster, hand blender, baking tins. |
| [PLN-03-09](E03-library-import-and-review/PLN-03-09-library-coverage-report.md) | Q100 — these targets (they leave room for ≥ 3 matches after stock filtering, Verification 1). |
| [PLN-04-01](E04-dish-nutrition/PLN-04-01-compute-dish-nutrition.md) | Q084 — the recipe must state the absorbed oil; without it the dish shows "calories incomplete" and cannot be published. No default share. |
| [PLN-04-01](E04-dish-nutrition/PLN-04-01-compute-dish-nutrition.md) | Q092 — the default plate is base + the first option of each topping group, plus the min–max range. |
| [PLN-04-02](E04-dish-nutrition/PLN-04-02-recompute-on-catalogue-events.md) | Q095 — . |
| [PLN-04-03](E04-dish-nutrition/PLN-04-03-nutrition-in-responses.md) | Q096 — the app hides them; planner never strips values. |
| [PLN-04-04](E04-dish-nutrition/PLN-04-04-scanned-product-override.md) | Q098 — . |
| [PLN-04-04](E04-dish-nutrition/PLN-04-04-scanned-product-override.md) | Q097 — KIT-06 accepts a product EAN per line (features.md §5 step 4 is kitchen's rule). |
| [PLN-05-02](E05-library-browse/PLN-05-02-hard-limits.md) | Q099 — browsing filters for everyone in the household by default; the app's "for me" switch sends only the user's member id. |
| [PLN-05-04](E05-library-browse/PLN-05-04-favourites.md) | Q101 — . |
| [PLN-05-06](E05-library-browse/PLN-05-06-cooked-history.md) | Q102 — . |
| [PLN-06-01](E06-what-can-i-make-now/PLN-06-01-stock-coverage.md) | Q103 — items past their date count as at home; every line they cover is flagged `check_first` with those stock item ids, and the flag travels to suggestions, plans and cooking mode so the app shows "check it first". |
| [PLN-06-02](E06-what-can-i-make-now/PLN-06-02-ranking.md) | Q104 — the weights above (0.7/0.3, +0.05 favourite, −0.20/−0.10 recently cooked, 0.30 cut-off) are config values, tuned with the concierge-test data. |
| [PLN-07-02](E07-rule-checker/PLN-07-02-tool-check.md) | Q105 — no substitution table at launch. |
| [PLN-07-03](E07-rule-checker/PLN-07-03-allergen-diet-dislike-check.md) | Q106 — the app and the terms state that cross-contact isn't checked. |
| [PLN-07-04](E07-rule-checker/PLN-07-04-stock-check.md) | Q107 — . |
| [PLN-07-05](E07-rule-checker/PLN-07-05-ai-output-checks.md) | Q108 — . |
| [PLN-08-01](E08-ai-day-plan/PLN-08-01-plan-storage.md) | Q118 — . |
| [PLN-08-02](E08-ai-day-plan/PLN-08-02-compose-request.md) | Q119 — this applies as a soft preference from R1, not only in family plans. |
| [PLN-08-04](E08-ai-day-plan/PLN-08-04-validate-repair-store.md) | Q120 — such a plan does not use Free quota; the gateway counts only when at least one AI meal survived (field `ai_meals_kept`). |
| [PLN-08-08](E08-ai-day-plan/PLN-08-08-set-plan-meal.md) | Q117 — users can put a library dish they picked (from browse, favourites or "what can I make now") into a day-plan meal themselves. It is checked by rules, uses no AI and no Free quota. |
| [PLN-09-02](E09-swaps-and-rewrites/PLN-09-02-replace-ingredient.md) | Q124 — no non-AI substitution table for Free at launch. |
| [PLN-09-05](E09-swaps-and-rewrites/PLN-09-05-restyle.md) | Q223 — the allowed lines are those in KIT-02's seasoning, herb, sauce and fat category codes; the code list is fixed in this task at implementation. |
| [PLN-10-02](E10-cooking-mode/PLN-10-02-mark-cooked.md) | Q126 — KIT-07 `ConsumeStock` accepts an idempotency key and preferred stock item ids (both needed for offline retries and "oldest first" from the plan). |
| [PLN-11-01](E11-library-growth/PLN-11-01-candidate-selection.md) | Q214 — . |
| [PLN-11-02](E11-library-growth/PLN-11-02-anonymised-copy.md) | Q215 — users are told by a sentence in the terms and on the AI-dish card that their AI dishes may join the shared library in anonymised form; no opt-in. |
| [PLN-11-03](E11-library-growth/PLN-11-03-promotion-admin-and-metrics.md) | Q216 — planner does not skip the AI call yet; this metric is measured first and the choice is made with launch data. |
| [PLN-12-01](E12-per-person-portions/PLN-12-01-meal-attendance.md) | Q196 — weekly rules exist ("school lunch" repeats every weekday), at most 20 per household. |
| [PLN-12-02](E12-per-person-portions/PLN-12-02-meal-shares.md) | Q197 — extra eaters beyond the household's people count as unnamed eaters with 1 standard portion each (rule 7). |
| [PLN-12-03](E12-per-person-portions/PLN-12-03-portion-and-pot-calculator.md) | Q198 — rule 9 (incomplete calories → everyone gets 1 standard portion, flagged "calories incomplete"). |
| [PLN-12-04](E12-per-person-portions/PLN-12-04-no-cook-extras.md) | Q199 — the default order (rule 2) is bread, fruit, yoghurt, milk, cheese, nuts, raw vegetables, with the lighter order for "eat lighter as a family". |
| [PLN-12-04](E12-per-person-portions/PLN-12-04-no-cook-extras.md) | Q200 — nuts are never offered to child profiles (choking risk for small children; age is not on the server). |
| [PLN-12-05](E12-per-person-portions/PLN-12-05-carry-over-and-day-balance.md) | Q202 — the day check compares planned meals with 90% of the target (target minus the 10% snacks share), because snacks are not planned. |
| [PLN-12-05](E12-per-person-portions/PLN-12-05-carry-over-and-day-balance.md) | Q201 — the ±25% carry cap (rule 4). |
| [PLN-12-07](E12-per-person-portions/PLN-12-07-calorie-visibility.md) | Q203 — identity has no parent–child link; every adult account member of the household sees a child's calories and target when they are not hidden (rule 4). |
| [PLN-13-01](E13-family-day-plan/PLN-13-01-shared-meal-limits.md) | Q204 — the narrowness thresholds in rule 5 (fewer than 5 dishes fit everyone, or under 30% of those without that person's limits). |
| [PLN-13-03](E13-family-day-plan/PLN-13-03-plan-comparator.md) | Q205 — "eat lighter" sits at step 4, after each person's day target and before using up stock. |
| [PLN-13-03](E13-family-day-plan/PLN-13-03-plan-comparator.md) | Q206 — the urgency formula in step 5 (share used × 1/(1 + days to expiry); undated items 0.05). |
| [PLN-14-02](E14-week-plan/PLN-14-02-cooking-days.md) | Q217 — `twice_a_week` cooks on days 1 and 4, double amounts only; the other main meals are no-cook. |
| [PLN-14-04](E14-week-plan/PLN-14-04-effort-summary.md) | Q218 — a double amount takes ×1.25 the time; +10 minutes per topping that needs its own pan. Replace with per-dish values if PLN-02 stores them later. |
| [PLN-15-03](E15-shopping-list/PLN-15-03-subtract-stock-and-staples.md) | Q127 — rule 4 (stock expiring before the meal date does not count, so the item goes on the list). Items already past their date are still offered in plans with "check it first" (Q103); the list buys a replacement. |
| [PLN-15-04](E15-shopping-list/PLN-15-04-merge-and-round-to-packs.md) | Q128 — at most 2 pack sizes and 20 packs per line (rule 5); without a pack size, round up to 10 g (≤ 100 g) or 50 g (rule 6). |
| [PLN-15-07](E15-shopping-list/PLN-15-07-tick-into-fridge.md) | Q129 — rule 4 (free-text lines are ticked on the list only; they don't go into the fridge). |
| [PLN-15-10](E15-shopping-list/PLN-15-10-list-follows-plan.md) | Q207 — rule 4 (increases go into a "plan changed, review" proposal; decreases apply directly). |
| [PLN-17-01](E17-offer-badges-and-reminders/PLN-17-01-offer-badges.md) | Q177, Q176): offers starting within 3 days show as `upcoming` badges (rule 3); the smallest discount shown is 5%, rounded down to a whole percent, else the price only (rule 5). |
| [PLN-17-04](E17-offer-badges-and-reminders/PLN-17-04-offer-ends-reminders.md) | Q179, Q178): reminders are sent at 18:00 local time, to Free and Plus households. |
| [PLN-18-02](E18-plan-from-deals/PLN-18-02-ai-plan-from-deals.md) | Q208 — the fridge-first thresholds (rule 4: each meal uses ≥ 1 stock item, ≥ 50% by weight from stock, at most 3 non-offer purchases a day). |
| [PLN-19-01](E19-trip-price-matrix/PLN-19-01-candidate-shops.md) | Q162 — at most 10 shops (175 three-shop sets, trivial on a phone). |
| [PLN-19-04](E19-trip-price-matrix/PLN-19-04-receipt-prices.md) | Q220 — receipt prices older than 60 days are ignored. |
| [PLN-20-01](E20-cheapest-top-up/PLN-20-01-needs-and-allowed-swaps.md) | Q224 — calories change ≤ 15%, a swap saves ≥ 10% or €0.50, at most 3 swaps a week. |

