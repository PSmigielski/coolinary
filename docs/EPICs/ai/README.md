# ai — epics and tasks

Python 3.14 + uv + FastAPI, gRPC on `50051`. Every LLM call in Coolinary goes through this service (infrastructure.md).
Provider: **Claude API direct** with the official `anthropic` Python SDK (decision 5); **Claude Haiku 5.5** reads and
matches, **Claude Sonnet 5.5** composes for Plus (decision 6); **Free day plans compose on Haiku 5.5** (decision 15).
The stub provider stays for dev, tests and CI.

Read first: `docs/EPICs/README.md` (decisions, conventions, cross-service RPCs), `docs/EPICs/TEMPLATE.md`.

## Epics

| Epic | Title | Rel. | Tasks |
|---|---|---|---|
| [AI-01](E01-claude-provider/README.md) | Claude provider foundation | R1 | 9 |
| [AI-02](E02-cost-and-spend-limit/README.md) | Cost tracking and spend limit | R1 | 5 |
| [AI-03](E03-prompts-and-evals/README.md) | Prompts, golden sets and evals | R1 | 5 |
| [AI-04](E04-batch-jobs/README.md) | Batch jobs (Message Batches) | R1 | 4 |
| [AI-05](E05-text-to-items/README.md) | Text to items | R1–R3 | 2 |
| [AI-06](E06-photo-to-items/README.md) | Photo to items | R1 | 2 |
| [AI-07](E07-label-reading/README.md) | Product label reading | R1 | 2 |
| [AI-08](E08-catalogue-match/README.md) | Catalogue match pick | R1 | 3 |
| [AI-09](E09-catalogue-library-batch-jobs/README.md) | Catalogue and library batch jobs | R1 | 6 |
| [AI-10](E10-day-plan/README.md) | Day plan composition and cuisine choices | R1 | 6 |
| [AI-11](E11-swaps-rewrites/README.md) | Swaps, wishes, tool rewrite, restyle | R1–R3 | 6 |
| [AI-12](E12-family-day-plan/README.md) | Family day plan | R2 | 4 |
| [AI-13](E13-flyer-offers/README.md) | Flyer offer extraction | R2 | 4 |
| [AI-14](E14-plan-from-deals/README.md) | Plan from deals | R2 | 2 |
| [AI-15](E15-week-plan/README.md) | Week plan | R3 | 4 |
| [AI-16](E16-receipt-reading/README.md) | Receipt reading | R3 | 3 |
| [AI-17](E17-cheapest-top-up/README.md) | Cheapest top-up | R3 | 2 |

69 tasks.

## Jobs, models and effort

Starting values in `ai/src/coolinary_ai/jobs.toml` (AI-01-01); each is tuned by its golden set (AI-03) and can be
overridden per plan (`by_plan`) and by env. Costs are estimates until measured (AI-02).

| Job | RPC / batch kind | Model | Effort | Mode | Rel. | ≈ USD per call |
|---|---|---|---|---|---|---|
| text_to_items | `ExtractItemsFromText` | claude-haiku-5-5 | low | sync | R1 (built), used from R3 (Q219) | 0.0003 |
| photo_to_items | `ExtractItemsFromPhoto` | claude-haiku-5-5 | medium | sync | R1 | 0.001 |
| read_label | `ReadProductLabel` | claude-haiku-5-5 | low | sync | R1 | 0.0005 |
| pick_matches | `PickCatalogueMatches` / `pick_catalogue_matches` | claude-haiku-5-5 | low | sync + batch | R1 | 0.0005 per 10 names |
| translate_ingredients | `translate_ingredients` | claude-haiku-5-5 | low | batch | R1 | < 0.0002 |
| match_nutrition_entries | `match_nutrition_entries` | claude-haiku-5-5 | medium | batch | R1 | < 0.0003 |
| draft_dishes | `draft_dishes` | claude-sonnet-5-5 | high | batch | R1 | ≈ 0.02 |
| translate_dishes | `translate_dishes` | claude-haiku-5-5 | low | batch | R1 | < 0.0005 |
| suggest_effort_tags | `suggest_effort_tags` | claude-haiku-5-5 | low | batch | R1 | < 0.0002 |
| compose_day_plan | `ComposeDayPlan` | Free: claude-haiku-5-5 (decision 15). Plus: claude-haiku-5-5, moved to claude-sonnet-5-5 only if it clearly wins (AI-10-05) | medium | sync | R1 | Haiku 0.0025, Sonnet 0.05 |
| modify_meal | `ModifyMeal` | claude-sonnet-5-5 | low | sync | R1 | 0.027 |
| rewrite_dish | `RewriteDish` | claude-sonnet-5-5 | low | sync | R2–R3 | 0.025 |
| compose_family_day | `ComposeFamilyDay` | claude-sonnet-5-5 | medium | sync | R2 | 0.06 |
| extract_flyer_offers | `ExtractFlyerOffers` / `extract_flyer_offers` | claude-haiku-5-5 | medium | sync + batch | R2 | 0.001 a page (batch 0.0005) |
| compose_from_deals | `ComposeFromDeals` | claude-sonnet-5-5 | medium | sync | R2 | 0.06 |
| compose_week_plan | `ComposeWeekPlan` | claude-sonnet-5-5 | high | sync, streamed | R3 | 0.12 |
| extract_receipt | `ExtractReceipt` | claude-haiku-5-5 | medium | sync | R3 | 0.001 |
| cheapest_top_up | `CheapestTopUp` | claude-sonnet-5-5 | medium | sync | R3 | 0.05 |

Composition jobs other than `compose_day_plan` are Plus-only, so the Sonnet rows above apply to Plus households;
each still starts from the AI-10-05 comparison (Sonnet only where it clearly beats Haiku).

Composition jobs on Sonnet 5.5 use server-side refusal fallbacks on the Claude API (AI-01-03, Q017);
batches never do (not supported there).

## Rules for callers (kitchen, planner, offers)

1. Set a deadline on every call; send `x-plan`, `x-household-id`, `x-language`, `x-request-id` (AI-01-08). Missing
   `x-plan` is treated as `free`. `x-plan` is the plan the caller got from identity `GetEntitlements` for this request
   (Decided, Q004); ai itself never calls identity.
2. **Fall back to the non-AI result on every status except `INVALID_ARGUMENT`.** Reason tokens (first word of the
   status message): `ai_refused`, `ai_invalid_output`, `ai_rate_limited`, `ai_overloaded`, `ai_unreachable`,
   `ai_misconfigured`, `ai_request_bug`, `ai_budget`, `ai_household_budget`, `ai_disabled`, `ai_circuit_open`,
   `ai_bad_request` (AI-01-04).
3. Every response carries `CallInfo` (model, prompt version, cached, cost in micro-euros, repaired).
4. Composition calls send ingredient facts for every id they mention, stock in its own units, and only catalogue
   names — never user text except a wish, never health data (AI-10-01). Stock past its date is sent with
   `past_date = true`; its uses come back with `check_first = true` (Q103). Family calls use opaque eater refs
   `e1`–`e8` and meal kcal targets only (AI-12-01).
5. Batch jobs: store `batch_id`, `kind`, `prompt_version` and your inputs; poll `GetBatch` every 10 minutes; fetch
   results with `GetBatchResults(batch_id, kind, prompt_version)` before `results_available_until`; join by
   `custom_id`; **run the input-dependent checks listed in each kind's task** (ai keeps no inputs, D10); resubmit
   items with `retryable=true` at most twice; if more than 20% of a batch fails, resubmit nothing and send a founder
   alert (AI-04-04).
6. `SubmitBatch` limits: 1–2,000 items and 30 MiB per call; split larger work.
7. Images: JPEG/PNG/WebP, ≤ 3 MiB each, long edge ≤ 2576 px (the phone sends 2048 px). ai never fetches URLs and
   never stores photos (D6).

## Build order

### R1 — foundation (provider, cost, evals, batches)

| # | Task | Title | Rel. | Size | Depends on |
|---|---|---|---|---|---|
| 1 | [AI-01-01](E01-claude-provider/AI-01-01-sdk-config-and-job-registry.md) | SDK dependency, settings and job registry | R1 | M | — |
| 2 | [AI-01-02](E01-claude-provider/AI-01-02-provider-interface.md) | Provider interface for structured jobs | R1 | M | AI-01-01 |
| 3 | [AI-01-08](E01-claude-provider/AI-01-08-grpc-limits-and-context.md) | gRPC limits, request context and concurrency | R1 | M | AI-01-01 |
| 4 | [AI-03-01](E03-prompts-and-evals/AI-03-01-prompt-registry.md) | Prompt registry and versions | R1 | S | AI-01-02 |
| 5 | [AI-01-03](E01-claude-provider/AI-01-03-claude-provider.md) | Claude provider | R1 | L | AI-01-02 |
| 6 | [AI-01-04](E01-claude-provider/AI-01-04-error-mapping-and-circuit-breaker.md) | Error mapping, retries, deadlines and circuit breaker | R1 | M | AI-01-03 |
| 7 | [AI-01-05](E01-claude-provider/AI-01-05-prompt-caching-layout.md) | Prompt caching and request layout | R1 | M | AI-01-03 |
| 8 | [AI-01-06](E01-claude-provider/AI-01-06-output-validation-and-repair.md) | Output validation, calorie guard and one repair call | R1 | M | AI-01-02 |
| 9 | [AI-01-07](E01-claude-provider/AI-01-07-untrusted-input.md) | Untrusted input handling (prompt injection) | R1 | S | AI-01-02 |
| 10 | [AI-01-09](E01-claude-provider/AI-01-09-result-cache-rules.md) | Result cache rules per job | R1 | S | AI-01-02 |
| 11 | [AI-02-01](E02-cost-and-spend-limit/AI-02-01-price-table-and-cost.md) | Price table and cost calculator | R1 | S | AI-01-01 |
| 12 | [AI-02-02](E02-cost-and-spend-limit/AI-02-02-cost-metrics.md) | Cost and usage metrics, Grafana panels | R1 | S | AI-02-01 |
| 13 | [AI-02-03](E02-cost-and-spend-limit/AI-02-03-langfuse-fields.md) | Langfuse fields for every job | R1 | S | AI-02-01 |
| 14 | [AI-02-04](E02-cost-and-spend-limit/AI-02-04-daily-spend-limit.md) | Global daily spend limit | R1 | M | AI-02-01 |
| 15 | [AI-02-05](E02-cost-and-spend-limit/AI-02-05-household-cost.md) | Per-household monthly cost and cost per active household | R1 | M | AI-02-04 |
| 16 | [AI-03-02](E03-prompts-and-evals/AI-03-02-golden-set-format.md) | Golden set format, storage and privacy | R1 | S | AI-03-01 |
| 17 | [AI-03-03](E03-prompts-and-evals/AI-03-03-eval-runner.md) | Eval runner CLI | R1 | M | AI-03-02, AI-02-01 |
| 18 | [AI-03-04](E03-prompts-and-evals/AI-03-04-thresholds-and-policy.md) | Thresholds and release policy | R1 | S | AI-03-03 |
| 19 | [AI-03-05](E03-prompts-and-evals/AI-03-05-composition-scorer.md) | Composition scorer | R1 | M | AI-03-03 |
| 20 | [AI-04-01](E04-batch-jobs/AI-04-01-submit-batch.md) | SubmitBatch contract and request building | R1 | L | AI-01-03, AI-01-05, AI-02-04, AI-01-08 |
| 21 | [AI-04-02](E04-batch-jobs/AI-04-02-get-and-cancel-batch.md) | GetBatch and CancelBatch | R1 | S | AI-04-01 |
| 22 | [AI-04-03](E04-batch-jobs/AI-04-03-batch-results-stream.md) | GetBatchResults stream, validation and cost | R1 | M | AI-04-02, AI-01-06, AI-02-03 |
| 23 | [AI-04-04](E04-batch-jobs/AI-04-04-batch-failures-and-metrics.md) | Failed items, retry rules and batch metrics | R1 | S | AI-04-03 |

### R1 — reading and matching jobs

| # | Task | Title | Rel. | Size | Depends on |
|---|---|---|---|---|---|
| 24 | [AI-05-01](E05-text-to-items/AI-05-01-text-to-items-claude.md) | ExtractItemsFromText on Claude | R1 | M | AI-01-03, AI-01-06, AI-01-07, AI-01-09, AI-03-01 |
| 25 | [AI-05-02](E05-text-to-items/AI-05-02-text-golden-set.md) | Golden set: typed and spoken item lists | R3 | S | AI-05-01, AI-03-03 |
| 26 | [AI-06-01](E06-photo-to-items/AI-06-01-photo-to-items.md) | ExtractItemsFromPhoto | R1 | L | AI-05-01, AI-01-08 |
| 27 | [AI-06-02](E06-photo-to-items/AI-06-02-photo-golden-set.md) | Golden set: 60 fridge and shelf photos | R1 | M | AI-06-01, AI-03-03 |
| 28 | [AI-07-01](E07-label-reading/AI-07-01-read-product-label.md) | ReadProductLabel | R1 | M | AI-06-01 |
| 29 | [AI-07-02](E07-label-reading/AI-07-02-label-golden-set.md) | Golden set: 50 labels | R1 | S | AI-07-01, AI-03-03 |
| 30 | [AI-08-01](E08-catalogue-match/AI-08-01-pick-catalogue-matches.md) | PickCatalogueMatches (sync) | R1 | M | AI-01-03, AI-01-06, AI-01-09 |
| 31 | [AI-08-02](E08-catalogue-match/AI-08-02-batch-pick-matches.md) | Batch kind pick_catalogue_matches | R1 | S | AI-08-01, AI-04-03 |
| 32 | [AI-08-03](E08-catalogue-match/AI-08-03-match-golden-set.md) | Golden set: 500 names | R1 | M | AI-08-01, AI-03-03 |

### R1 — catalogue and library batch jobs

| # | Task | Title | Rel. | Size | Depends on |
|---|---|---|---|---|---|
| 33 | [AI-10-01](E10-day-plan/AI-10-01-composition-contract.md) | Composition contract and shared dish types | R1 | M | AI-02-01 |
| 34 | [AI-09-01](E09-catalogue-library-batch-jobs/AI-09-01-translate-ingredients.md) | translate_ingredients | R1 | M | AI-04-03 |
| 35 | [AI-09-02](E09-catalogue-library-batch-jobs/AI-09-02-match-nutrition-entries.md) | match_nutrition_entries | R1 | M | AI-04-03 |
| 36 | [AI-09-03](E09-catalogue-library-batch-jobs/AI-09-03-draft-dishes.md) | draft_dishes | R1 | L | AI-04-03, AI-10-01 |
| 37 | [AI-09-04](E09-catalogue-library-batch-jobs/AI-09-04-translate-dishes.md) | translate_dishes | R1 | M | AI-09-03 |
| 38 | [AI-09-05](E09-catalogue-library-batch-jobs/AI-09-05-suggest-effort-tags.md) | suggest_effort_tags | R1 | S | AI-04-03, AI-10-01 |
| 39 | [AI-09-06](E09-catalogue-library-batch-jobs/AI-09-06-batch-golden-sets.md) | Golden sets for catalogue and library jobs | R1 | M | AI-09-01 … AI-09-05, AI-03-03 |

### R1 — day plans, swaps and wishes

| # | Task | Title | Rel. | Size | Depends on |
|---|---|---|---|---|---|
| 40 | [AI-10-02](E10-day-plan/AI-10-02-compose-day-plan.md) | ComposeDayPlan prompt and call | R1 | L | AI-10-01, AI-01-03, AI-01-05, AI-01-09 |
| 41 | [AI-10-03](E10-day-plan/AI-10-03-plan-validation.md) | Validation and repair of composed meals | R1 | M | AI-10-02, AI-01-06 |
| 42 | [AI-10-04](E10-day-plan/AI-10-04-cuisine-choices.md) | Cuisine fit and the three choices | R1 | M | AI-10-03 |
| 43 | [AI-10-06](E10-day-plan/AI-10-06-day-plan-golden-set.md) | Golden set and eval for day plans | R1 | M | AI-10-04, AI-03-05 |
| 44 | [AI-10-05](E10-day-plan/AI-10-05-model-routing-by-plan.md) | Model routing by plan and the cost decision | R1 | S | AI-10-06, AI-02-05 |
| 45 | [AI-11-01](E11-swaps-rewrites/AI-11-01-modify-meal-swap-expand.md) | ModifyMeal: contract, SWAP and EXPAND_ALTERNATIVE | R1 | M | AI-10-03, AI-10-04 |
| 46 | [AI-11-02](E11-swaps-rewrites/AI-11-02-replace-ingredient.md) | ModifyMeal: REPLACE_INGREDIENT | R1 | M | AI-11-01 |
| 47 | [AI-11-03](E11-swaps-rewrites/AI-11-03-free-text-wish.md) | ModifyMeal: WISH (free text) | R1 | M | AI-11-01, AI-01-07 |
| 48 | [AI-11-06](E11-swaps-rewrites/AI-11-06-modify-golden-sets.md) | Golden sets for swaps, wishes and rewrites | R1–R3 | M | AI-11-01 … AI-11-05, AI-03-05 |

### R2 — tool rewrite, family plans, flyers, deals

| # | Task | Title | Rel. | Size | Depends on |
|---|---|---|---|---|---|
| 49 | [AI-11-04](E11-swaps-rewrites/AI-11-04-rewrite-for-tools.md) | RewriteDish: TOOLS | R2 | M | AI-11-01 |
| 50 | [AI-12-01](E12-family-day-plan/AI-12-01-family-contract.md) | ComposeFamilyDay contract | R2 | S | AI-10-01 |
| 51 | [AI-12-02](E12-family-day-plan/AI-12-02-family-compose.md) | Family composition prompt and call | R2 | L | AI-12-01, AI-10-02 |
| 52 | [AI-12-03](E12-family-day-plan/AI-12-03-family-validation.md) | Validation for eaters, toppings and second dishes | R2 | M | AI-12-02, AI-10-03 |
| 53 | [AI-12-04](E12-family-day-plan/AI-12-04-family-golden-set.md) | Golden set: families including the worked example | R2 | M | AI-12-03, AI-03-05 |
| 54 | [AI-13-01](E13-flyer-offers/AI-13-01-offer-schema.md) | Offer schema and checks | R2 | M | AI-01-06, AI-06-01 |
| 55 | [AI-13-02](E13-flyer-offers/AI-13-02-extract-flyer-offers-sync.md) | ExtractFlyerOffers for user flyers (sync, chunked) | R2 | L | AI-13-01, AI-01-08 |
| 56 | [AI-13-03](E13-flyer-offers/AI-13-03-flyer-batch.md) | Batch kind extract_flyer_offers for chain flyers | R2 | S | AI-13-01, AI-04-03 |
| 57 | [AI-13-04](E13-flyer-offers/AI-13-04-flyer-golden-set.md) | Golden set per launch chain | R2 | M | AI-13-02, AI-13-03, AI-03-03 |
| 58 | [AI-14-01](E14-plan-from-deals/AI-14-01-compose-from-deals.md) | ComposeFromDeals | R2 | L | AI-10-03 |
| 59 | [AI-14-02](E14-plan-from-deals/AI-14-02-deals-validation-and-golden-set.md) | Deals validation and golden set | R2 | M | AI-14-01, AI-03-05 |

### R3 — restyle, week plans, receipts, top-up

| # | Task | Title | Rel. | Size | Depends on |
|---|---|---|---|---|---|
| 60 | [AI-11-05](E11-swaps-rewrites/AI-11-05-restyle.md) | RewriteDish: RESTYLE | R3 | S | AI-11-04 |
| 61 | [AI-15-01](E15-week-plan/AI-15-01-week-plan-contract.md) | ComposeWeekPlan contract | R3 | S | AI-12-01 |
| 62 | [AI-15-02](E15-week-plan/AI-15-02-week-plan-compose.md) | Week composition prompt and call | R3 | L | AI-15-01 |
| 63 | [AI-15-03](E15-week-plan/AI-15-03-week-plan-validation.md) | Validation of sessions and leftovers | R3 | M | AI-15-02, AI-12-03 |
| 64 | [AI-15-04](E15-week-plan/AI-15-04-week-plan-golden-set.md) | Golden set: family weeks | R3 | M | AI-15-03, AI-03-05 |
| 65 | [AI-16-01](E16-receipt-reading/AI-16-01-extract-receipt.md) | ExtractReceipt | R3 | M | AI-06-01 |
| 66 | [AI-16-02](E16-receipt-reading/AI-16-02-receipt-privacy-and-totals.md) | Receipt privacy filter and totals check | R3 | S | AI-16-01 |
| 67 | [AI-16-03](E16-receipt-reading/AI-16-03-receipt-golden-set.md) | Golden set: 40 receipts | R3 | M | AI-16-02, AI-03-03 |
| 68 | [AI-17-01](E17-cheapest-top-up/AI-17-01-cheapest-top-up.md) | CheapestTopUp | R3 | L | AI-14-01 |
| 69 | [AI-17-02](E17-cheapest-top-up/AI-17-02-top-up-validation-and-golden-set.md) | Top-up validation and golden set | R3 | M | AI-17-01, AI-03-05 |
## Decided after writing (README decisions 15–17)

- Free day plans compose on Haiku 5.5; Plus compositions use Sonnet 5.5 only if it clearly beats Haiku (AI-10-05).
- Spend counters live in redis-state; ai gets `REDIS_STATE_URL` (AI-02-04).
- Shared enums (`Allergen`, `Diet`, `Cuisine`, `KitchenTool`, `EffortTag`) and extra `Unit` values come from
  `coolinary.common.v1`; ai defines no copies.

## Decisions

Final answers (docs/EPICs/DECISIONS.md), written into the tasks as `Decided (Qnnn)`.

| Decision | Rule | Tasks |
|---|---|---|
| Q017 | Server-side refusal fallbacks on for every Sonnet job; the answering model's real price is recorded. Haiku refusals go to the caller's non-AI fallback. | [AI-01-03](E01-claude-provider/AI-01-03-claude-provider.md) |
| Q018 | Spend limits: €20 a day overall, Free at most 40% of it, €60 a month for batch jobs (config values). | [AI-02-04](E02-cost-and-spend-limit/AI-02-04-daily-spend-limit.md) |
| Q019 | Hard guard of €0.15 per Free household per month. | [AI-02-05](E02-cost-and-spend-limit/AI-02-05-household-cost.md) |
| Q022 | User flyer PDFs up to 40 pages and 25 MiB. | [AI-13-02](E13-flyer-offers/AI-13-02-extract-flyer-offers-sync.md) |
| Q103 | Past-date stock is usable in compositions; every use of it comes back with `check_first = true`, set by ai. | [AI-10-01](E10-day-plan/AI-10-01-composition-contract.md), [AI-10-02](E10-day-plan/AI-10-02-compose-day-plan.md), [AI-10-03](E10-day-plan/AI-10-03-plan-validation.md) |
| Q110 | Label photos give name, brand and pack size only in R1; no nutrition table. | [AI-07-01](E07-label-reading/AI-07-01-read-product-label.md) |
| Q121 | Dishes drafted in English with `original_name` kept, then translated to pl and de. | [AI-09-03](E09-catalogue-library-batch-jobs/AI-09-03-draft-dishes.md) |
| Q123 | Plus moves to Sonnet only with 10 percentage points more plans passing first time; Sonnet repairs Plus output only. | [AI-10-05](E10-day-plan/AI-10-05-model-routing-by-plan.md) |
| Q219 | Typed and voice intake reach users in R3; AI-05-01 is still built in R1 (photo intake reuses it), AI-05-02 moves to R3. | [AI-05-01](E05-text-to-items/AI-05-01-text-to-items-claude.md), [AI-05-02](E05-text-to-items/AI-05-02-text-golden-set.md) |

## Notes for other services

- `CallInfo`, `InputImage`, `PackageHint`, `StorageHint` and the composition types live in `coolinary.ai.v1`
  (`common.proto`, `compose.proto`, `matching.proto`, `batch.proto`, `flyers.proto`).
- Recipe amounts from AI are only g, ml or pieces (+ a display hint like "2 tbsp"); spoons and cups are converted by
  kitchen's typical weights (KIT-06) when they come from hand-written recipes.
- `CancelBatch` (AI-04-02) is an extra RPC not listed in `docs/EPICs/README.md`; `GetBatchResults` is
  server-streaming, so only Go services call it.
