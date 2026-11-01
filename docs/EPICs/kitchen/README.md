# kitchen — epics and tasks

Go (`kitchen/`). Owns the ingredient catalogue and its translations, name matching (pg_trgm + Haiku), the barcode lookup chain, unknown products, nutrition tables and the nutrition calculation, the household's fridge stock and AI intake, expiry reminders, the kitchen profile (decision 3) with per-person targets, and the savings ledger. It produces `events.kitchen` and publishes `stock.changed` live notices.

What exists today: `KitchenService.ResolveBarcode` over a `products` table with one example row (pg_trgm installed), the Go platform packages (config, db, cache, log, rpc, admin) and the Go → Bun smoke test.

KIT-01 to KIT-06 and KIT-07 to KIT-13 were written by two writers; dependencies between the halves use epic ids.

## Epics

| Epic | Title | Rel. | Tasks |
|---|---|---|---|
| [KIT-01](E01-platform/README.md) | Platform: events, plan from identity, admin service, batch jobs | R1 | 8 |
| [KIT-02](E02-ingredient-catalogue/README.md) | Ingredient catalogue and translations | R1 | 10 |
| [KIT-03](E03-name-matching/README.md) | Name matching (trigram + Haiku) | R1 | 7 |
| [KIT-04](E04-barcode-lookup/README.md) | Barcode lookup chain | R1 | 11 |
| [KIT-05](E05-unknown-products/README.md) | Unknown products and shared-catalogue confirmation | R1 | 5 |
| [KIT-06](E06-nutrition/README.md) | Nutrition tables and calculation | R1 | 10 |
| [KIT-07](E07-fridge-stock/README.md) | Fridge stock | R1 | 9 |
| [KIT-08](E08-ai-intake/README.md) | AI intake: photo, text, receipt | R1–R3 | 6 |
| [KIT-09](E09-expiry-reminders/README.md) | Expiry reminders | R1 | 3 |
| [KIT-10](E10-kitchen-profile/README.md) | Kitchen profile and household settings | R1 | 9 |
| [KIT-11](E11-per-person-profiles/README.md) | Per-person profiles and calorie targets | R2 | 6 |
| [KIT-12](E12-savings/README.md) | Savings ledger and monthly report | R2–R3 | 5 |
| [KIT-13](E13-household-lifecycle/README.md) | Household lifecycle: event consumers and data export | R1–R2 | 9 |

98 tasks: 82 in R1, 11 in R2, 5 in R3.

## What other services use

RPCs defined in this service's tasks (first task that defines each). Admin RPCs are called only by admin.

| RPC | Defined in | Admin |
|---|---|---|
| `ListDeadLetters` | [KIT-01-05](E01-platform/KIT-01-05-dead-letters-lag.md) | yes |
| `ReplayDeadLetter` | [KIT-01-05](E01-platform/KIT-01-05-dead-letters-lag.md) | yes |
| `ListAdminChanges` | [KIT-01-07](E01-platform/KIT-01-07-admin-service-scaffold.md) | yes |
| `GetQueueCounts` | [KIT-01-07](E01-platform/KIT-01-07-admin-service-scaffold.md) | yes |
| `StartCatalogueJob` | [KIT-01-08](E01-platform/KIT-01-08-ai-batch-tracker.md) | yes |
| `ListCatalogueJobs` | [KIT-01-08](E01-platform/KIT-01-08-ai-batch-tracker.md) | yes |
| `SearchIngredients` | [KIT-02-05](E02-ingredient-catalogue/KIT-02-05-search-ingredients.md) |  |
| `GetIngredients` | [KIT-02-06](E02-ingredient-catalogue/KIT-02-06-get-ingredients.md) | yes |
| `GetCatalogueVersion` | [KIT-02-06](E02-ingredient-catalogue/KIT-02-06-get-ingredients.md) |  |
| `MergeIngredients` | [KIT-02-09](E02-ingredient-catalogue/KIT-02-09-merge-ingredients.md) | yes |
| `CreateIngredient` | [KIT-02-10](E02-ingredient-catalogue/KIT-02-10-admin-catalogue-rpcs.md) | yes |
| `UpdateIngredient` | [KIT-02-10](E02-ingredient-catalogue/KIT-02-10-admin-catalogue-rpcs.md) | yes |
| `SetIngredientNames` | [KIT-02-10](E02-ingredient-catalogue/KIT-02-10-admin-catalogue-rpcs.md) | yes |
| `ApproveMachineNames` | [KIT-02-10](E02-ingredient-catalogue/KIT-02-10-admin-catalogue-rpcs.md) | yes |
| `SetMeasures` | [KIT-02-10](E02-ingredient-catalogue/KIT-02-10-admin-catalogue-rpcs.md) | yes |
| `SetPackSizes` | [KIT-02-10](E02-ingredient-catalogue/KIT-02-10-admin-catalogue-rpcs.md) | yes |
| `RetireIngredient` | [KIT-02-10](E02-ingredient-catalogue/KIT-02-10-admin-catalogue-rpcs.md) | yes |
| `ListIngredients` | [KIT-02-10](E02-ingredient-catalogue/KIT-02-10-admin-catalogue-rpcs.md) | yes |
| `MatchNames` | [KIT-03-03](E03-name-matching/KIT-03-03-match-names-deterministic.md) |  |
| `ListMatchReviews` | [KIT-03-06](E03-name-matching/KIT-03-06-match-review-queue.md) | yes |
| `DecideMatch` | [KIT-03-06](E03-name-matching/KIT-03-06-match-review-queue.md) | yes |
| `ListInternalCodePatterns` | [KIT-04-03](E04-barcode-lookup/KIT-04-03-internal-codes.md) | yes |
| `SetInternalCodePattern` | [KIT-04-03](E04-barcode-lookup/KIT-04-03-internal-codes.md) | yes |
| `DecodeInternalCode` | [KIT-04-03](E04-barcode-lookup/KIT-04-03-internal-codes.md) | yes |
| `ListHouseholdProducts` | [KIT-04-04](E04-barcode-lookup/KIT-04-04-household-products.md) |  |
| `NameProduct` | [KIT-05-01](E05-unknown-products/KIT-05-01-name-product.md) |  |
| `ConfirmProposal` | [KIT-05-02](E05-unknown-products/KIT-05-02-submissions-and-confirmation.md) |  |
| `ReportProduct` | [KIT-05-02](E05-unknown-products/KIT-05-02-submissions-and-confirmation.md) |  |
| `ReadLabel` | [KIT-05-03](E05-unknown-products/KIT-05-03-label-photo.md) |  |
| `ListProductReviews` | [KIT-05-05](E05-unknown-products/KIT-05-05-admin-product-review.md) | yes |
| `GetSharedProduct` | [KIT-05-05](E05-unknown-products/KIT-05-05-admin-product-review.md) | yes |
| `ResolveProduct` | [KIT-05-05](E05-unknown-products/KIT-05-05-admin-product-review.md) | yes |
| `EditSharedProduct` | [KIT-05-05](E05-unknown-products/KIT-05-05-admin-product-review.md) | yes |
| `ListNutritionSources` | [KIT-06-01](E06-nutrition/KIT-06-01-nutrition-schema.md) |  |
| `ListNutritionSuggestions` | [KIT-06-06](E06-nutrition/KIT-06-06-admin-nutrition-approval.md) | yes |
| `ApproveNutritionLink` | [KIT-06-06](E06-nutrition/KIT-06-06-admin-nutrition-approval.md) | yes |
| `RejectNutritionLink` | [KIT-06-06](E06-nutrition/KIT-06-06-admin-nutrition-approval.md) | yes |
| `SetNutritionLink` | [KIT-06-06](E06-nutrition/KIT-06-06-admin-nutrition-approval.md) | yes |
| `SearchNutritionEntries` | [KIT-06-06](E06-nutrition/KIT-06-06-admin-nutrition-approval.md) | yes |
| `CreateAdminEntry` | [KIT-06-06](E06-nutrition/KIT-06-06-admin-nutrition-approval.md) | yes |
| `NutritionCoverage` | [KIT-06-06](E06-nutrition/KIT-06-06-admin-nutrition-approval.md) | yes |
| `ComputeNutrition` | [KIT-06-09](E06-nutrition/KIT-06-09-compute-nutrition-rpc.md) |  |
| `SuggestStockDefaults` | [KIT-07-02](E07-fridge-stock/KIT-07-02-storage-and-expiry-defaults.md) |  |
| `AddStockItems` | [KIT-07-03](E07-fridge-stock/KIT-07-03-add-stock-items.md) |  |
| `UpdateStockItem` | [KIT-07-04](E07-fridge-stock/KIT-07-04-update-and-remove.md) |  |
| `RemoveStockItems` | [KIT-07-04](E07-fridge-stock/KIT-07-04-update-and-remove.md) |  |
| `ListStock` | [KIT-07-05](E07-fridge-stock/KIT-07-05-list-stock.md) |  |
| `PreviewStockItems` | [KIT-07-06](E07-fridge-stock/KIT-07-06-proposals-and-unpack.md) |  |
| `ConsumeStock` | [KIT-07-07](E07-fridge-stock/KIT-07-07-consume-stock.md) |  |
| `IntakeFromPhoto` | [KIT-08-02](E08-ai-intake/KIT-08-02-photo-intake.md) |  |
| `IntakeFromText` | [KIT-08-03](E08-ai-intake/KIT-08-03-text-intake.md) |  |
| `IntakeFromReceipt` | [KIT-08-04](E08-ai-intake/KIT-08-04-receipt-intake.md) |  |
| `GetProfile` | [KIT-10-01](E10-kitchen-profile/KIT-10-01-profile-record-and-defaults.md) |  |
| `UpdateProfile` | [KIT-10-01](E10-kitchen-profile/KIT-10-01-profile-record-and-defaults.md) |  |
| `GetProfileOptions` | [KIT-10-07](E10-kitchen-profile/KIT-10-07-profile-options.md) |  |
| `GetPlanningProfile` | [KIT-10-08](E10-kitchen-profile/KIT-10-08-planning-profile-snapshot.md) |  |
| `UpdateMemberProfile` | [KIT-11-02](E11-per-person-profiles/KIT-11-02-edit-rights.md) |  |
| `GetMemberProfiles` | [KIT-11-06](E11-per-person-profiles/KIT-11-06-member-visibility-and-snapshot.md) |  |
| `GetSavings` | [KIT-12-03](E12-savings/KIT-12-03-get-savings.md) |  |
| `ListSavingsReports` | [KIT-12-03](E12-savings/KIT-12-03-get-savings.md) |  |
| `GetSavingsCard` | [KIT-12-05](E12-savings/KIT-12-05-share-card.md) |  |
| `ExportHouseholdData` | [KIT-13-07](E13-household-lifecycle/KIT-13-07-export-household-data.md) |  |

## Build order

Task order follows [ROADMAP.md](../ROADMAP.md); `#` is the task's number there and `Milestone` its milestone.

### Release 1

| # | Task | Title | Size | Milestone |
|---|---|---|---|---|
| 3 | [KIT-01-01](E01-platform/KIT-01-01-request-metadata.md) | Request metadata and household guard | S | M1 |
| 4 | [KIT-01-02](E01-platform/KIT-01-02-event-envelope-outbox.md) | Event envelope and outbox writer | M | M1 |
| 5 | [KIT-01-03](E01-platform/KIT-01-03-outbox-relay.md) | Outbox relay to Redis Streams | M | M1 |
| 6 | [KIT-01-04](E01-platform/KIT-01-04-stream-consumer.md) | Stream consumer with idempotent handlers | L | M1 |
| 7 | [KIT-01-07](E01-platform/KIT-01-07-admin-service-scaffold.md) | KitchenAdminService scaffold and change records | M | M1 |
| 33 | [KIT-01-08](E01-platform/KIT-01-08-ai-batch-tracker.md) | AI client and batch job tracker | M | M1 |
| 73 | [KIT-01-05](E01-platform/KIT-01-05-dead-letters-lag.md) | Dead letters, lag metrics and alerts | S | M2 |
| 74 | [KIT-01-06](E01-platform/KIT-01-06-household-plan-copy.md) | Household plan from identity | S | M2 |
| 104 | [KIT-02-01](E02-ingredient-catalogue/KIT-02-01-ingredient-core-schema.md) | Ingredient core: categories, allergens, diet flags | M | M4 |
| 105 | [KIT-02-02](E02-ingredient-catalogue/KIT-02-02-measures-and-units.md) | Typical weights, density and unit conversion | M | M4 |
| 106 | [KIT-02-03](E02-ingredient-catalogue/KIT-02-03-storage-shelf-life-packs.md) | Storage place, shelf life and pack sizes | S | M4 |
| 107 | [KIT-03-01](E03-name-matching/KIT-03-01-name-normaliser.md) | Name normaliser per language | M | M4 |
| 108 | [KIT-02-04](E02-ingredient-catalogue/KIT-02-04-names-and-synonyms.md) | Names, synonyms and plurals per language | M | M4 |
| 109 | [KIT-02-05](E02-ingredient-catalogue/KIT-02-05-search-ingredients.md) | SearchIngredients (autocomplete) | M | M4 |
| 110 | [KIT-02-06](E02-ingredient-catalogue/KIT-02-06-get-ingredients.md) | GetIngredients and catalogue version | M | M4 |
| 111 | [KIT-02-07](E02-ingredient-catalogue/KIT-02-07-starter-catalogue-seed.md) | Starter catalogue seed and seed command | M | M4 |
| 112 | [KIT-02-10](E02-ingredient-catalogue/KIT-02-10-admin-catalogue-rpcs.md) | Admin catalogue editing RPCs | M | M4 |
| 113 | [KIT-03-02](E03-name-matching/KIT-03-02-match-store.md) | Stored matches | S | M4 |
| 114 | [KIT-02-09](E02-ingredient-catalogue/KIT-02-09-merge-ingredients.md) | Merge duplicate ingredients | M | M4 |
| 115 | [KIT-03-03](E03-name-matching/KIT-03-03-match-names-deterministic.md) | MatchNames: stored, exact and trigram paths | M | M4 |
| 116 | [KIT-03-06](E03-name-matching/KIT-03-06-match-review-queue.md) | Match review queue (admin) | M | M4 |
| 117 | [KIT-10-01](E10-kitchen-profile/KIT-10-01-profile-record-and-defaults.md) | Profile record, country defaults, GetProfile and UpdateProfile | M | M4 |
| 118 | [KIT-10-02](E10-kitchen-profile/KIT-10-02-cuisines-per-meal.md) | Cuisines per meal | S | M4 |
| 119 | [KIT-10-03](E10-kitchen-profile/KIT-10-03-kitchen-tools-and-presets.md) | Kitchen tools and presets | S | M4 |
| 120 | [KIT-10-04](E10-kitchen-profile/KIT-10-04-household-diet-allergens-dislikes.md) | Household diet, allergens and dislikes | M | M4 |
| 121 | [KIT-10-05](E10-kitchen-profile/KIT-10-05-staples.md) | Staples ("always at home") | S | M4 |
| 122 | [KIT-10-06](E10-kitchen-profile/KIT-10-06-household-rhythm-and-meal-shares.md) | Household rhythm and meal shares | M | M4 |
| 123 | [KIT-10-07](E10-kitchen-profile/KIT-10-07-profile-options.md) | Profile options and onboarding defaults | S | M4 |
| 124 | [KIT-10-08](E10-kitchen-profile/KIT-10-08-planning-profile-snapshot.md) | GetPlanningProfile snapshot for planner | M | M4 |
| 125 | [KIT-10-09](E10-kitchen-profile/KIT-10-09-profile-changed.md) | profile.changed event and planning version | S | M4 |
| 128 | [KIT-03-04](E03-name-matching/KIT-03-04-match-names-ai-pick.md) | MatchNames: Haiku pick and fallbacks | M | M4 |
| 129 | [KIT-03-05](E03-name-matching/KIT-03-05-bulk-matching-batch.md) | Bulk matching through Message Batches | M | M4 |
| 130 | [KIT-03-07](E03-name-matching/KIT-03-07-rematch-on-catalogue-change.md) | Re-match after catalogue changes | S | M4 |
| 134 | [KIT-02-08](E02-ingredient-catalogue/KIT-02-08-ai-translation-batch.md) | AI translation batch job | M | M4 |
| 142 | [KIT-04-01](E04-barcode-lookup/KIT-04-01-gtin-parsing.md) | GTIN parsing, check digits and prefix classes | M | M5 |
| 143 | [KIT-04-02](E04-barcode-lookup/KIT-04-02-gs1-2d-codes.md) | GS1 2D codes and Digital Link | M | M5 |
| 144 | [KIT-04-03](E04-barcode-lookup/KIT-04-03-internal-codes.md) | Shop-internal codes (prefixes 20–29) | M | M5 |
| 145 | [KIT-04-04](E04-barcode-lookup/KIT-04-04-household-products.md) | Household products (level 1) | M | M5 |
| 146 | [KIT-04-05](E04-barcode-lookup/KIT-04-05-shared-catalogue.md) | Shared catalogue (level 2) | M | M5 |
| 147 | [KIT-04-06](E04-barcode-lookup/KIT-04-06-open-food-facts-client.md) | Open Food Facts client (level 3) | M | M5 |
| 148 | [KIT-04-07](E04-barcode-lookup/KIT-04-07-open-food-facts-mapping.md) | Open Food Facts mapping and storage | M | M5 |
| 149 | [KIT-04-08](E04-barcode-lookup/KIT-04-08-gs1-lookup.md) | GS1 register lookup (level 4) | M | M5 |
| 150 | [KIT-04-09](E04-barcode-lookup/KIT-04-09-resolve-chain.md) | ResolveBarcode: the lookup chain | M | M5 |
| 151 | [KIT-04-10](E04-barcode-lookup/KIT-04-10-ingredient-and-suggestions.md) | Product → ingredient, storage and use-by suggestion | M | M5 |
| 152 | [KIT-04-11](E04-barcode-lookup/KIT-04-11-multipacks.md) | Multipacks | S | M5 |
| 153 | [KIT-05-01](E05-unknown-products/KIT-05-01-name-product.md) | Name an unknown product once | M | M5 |
| 154 | [KIT-05-02](E05-unknown-products/KIT-05-02-submissions-and-confirmation.md) | Submissions, proposals and confirmation by a second household | L | M5 |
| 155 | [KIT-05-04](E05-unknown-products/KIT-05-04-contribute-to-off.md) | Send confirmed products to Open Food Facts | M | M5 |
| 156 | [KIT-05-05](E05-unknown-products/KIT-05-05-admin-product-review.md) | Admin product review | M | M5 |
| 157 | [KIT-07-01](E07-fridge-stock/KIT-07-01-stock-schema.md) | Stock items and movements schema | M | M5 |
| 158 | [KIT-07-02](E07-fridge-stock/KIT-07-02-storage-and-expiry-defaults.md) | Default storage place and use-by date | M | M5 |
| 159 | [KIT-07-03](E07-fridge-stock/KIT-07-03-add-stock-items.md) | AddStockItems: add, merge, idempotent adds | M | M5 |
| 160 | [KIT-07-04](E07-fridge-stock/KIT-07-04-update-and-remove.md) | Edit and remove items (used, thrown away, mistake) | M | M5 |
| 161 | [KIT-07-05](E07-fridge-stock/KIT-07-05-list-stock.md) | ListStock with expiry status and totals | M | M5 |
| 162 | [KIT-07-06](E07-fridge-stock/KIT-07-06-proposals-and-unpack.md) | Stock proposals and unpack preview | M | M5 |
| 163 | [KIT-07-07](E07-fridge-stock/KIT-07-07-consume-stock.md) | ConsumeStock on "cooked it" | L | M5 |
| 164 | [KIT-07-08](E07-fridge-stock/KIT-07-08-stock-changed.md) | stock.changed event and live notice | S | M5 |
| 165 | [KIT-07-09](E07-fridge-stock/KIT-07-09-link-unmatched-items.md) | Link items typed without a catalogue match | M | M5 |
| 177 | [KIT-06-01](E06-nutrition/KIT-06-01-nutrition-schema.md) | Nutrition entries, links and sources | M | M6 |
| 178 | [KIT-06-02](E06-nutrition/KIT-06-02-import-ciqual.md) | Import Ciqual | M | M6 |
| 179 | [KIT-06-03](E06-nutrition/KIT-06-03-import-bls.md) | Import Bundeslebensmittelschlüssel 4.0 | M | M6 |
| 180 | [KIT-06-04](E06-nutrition/KIT-06-04-import-usda.md) | Import USDA FoodData Central | M | M6 |
| 181 | [KIT-06-05](E06-nutrition/KIT-06-05-match-suggestions.md) | Ingredient → nutrition entry suggestions | M | M6 |
| 182 | [KIT-06-07](E06-nutrition/KIT-06-07-pack-nutrition.md) | Pack nutrition values of scanned products | M | M6 |
| 183 | [KIT-06-08](E06-nutrition/KIT-06-08-calculation-rules.md) | Calculation rules (library) | M | M6 |
| 184 | [KIT-06-09](E06-nutrition/KIT-06-09-compute-nutrition-rpc.md) | ComputeNutrition RPC | M | M6 |
| 185 | [KIT-06-10](E06-nutrition/KIT-06-10-nutrition-changed-event.md) | Nutrition version and `catalogue.nutrition_changed` | S | M6 |
| 186 | [KIT-06-06](E06-nutrition/KIT-06-06-admin-nutrition-approval.md) | Admin approval of nutrition links | M | M6 |
| 223 | [KIT-08-01](E08-ai-intake/KIT-08-01-intake-pipeline.md) | Shared intake pipeline: match, propose, metrics | M | M7 |
| 226 | [KIT-08-02](E08-ai-intake/KIT-08-02-photo-intake.md) | Photo intake (fridge or shelf) | M | M7 |
| 229 | [KIT-05-03](E05-unknown-products/KIT-05-03-label-photo.md) | Read a product label photo | M | M7 |
| 313 | [KIT-09-01](E09-expiry-reminders/KIT-09-01-reminder-selection.md) | Which items a reminder lists | S | M10 |
| 314 | [KIT-09-02](E09-expiry-reminders/KIT-09-02-reminder-scheduler.md) | Daily reminder scheduler per household timezone | M | M10 |
| 315 | [KIT-09-03](E09-expiry-reminders/KIT-09-03-reminder-push.md) | Send the reminder and record it | S | M10 |
| 331 | [KIT-13-01](E13-household-lifecycle/KIT-13-01-household-created.md) | household.created: default profile and stock state | S | M11 |
| 332 | [KIT-13-02](E13-household-lifecycle/KIT-13-02-member-joined-with-fridge.md) | Joining with a fridge: move stock, carry the person's data | M | M11 |
| 333 | [KIT-13-04](E13-household-lifecycle/KIT-13-04-household-deleted.md) | household.deleted: erase the household | M | M11 |
| 334 | [KIT-13-05](E13-household-lifecycle/KIT-13-05-account-deleted.md) | account.deleted: erase the person | S | M11 |
| 335 | [KIT-13-06](E13-household-lifecycle/KIT-13-06-plan-changed.md) | Plan changes: planning version and effects | S | M11 |
| 336 | [KIT-13-07](E13-household-lifecycle/KIT-13-07-export-household-data.md) | ExportHouseholdData | M | M11 |
| 337 | [KIT-13-08](E13-household-lifecycle/KIT-13-08-retention-job.md) | Retention clean-up job | S | M11 |
| 338 | [KIT-13-09](E13-household-lifecycle/KIT-13-09-ingredient-merge-household-tables.md) | Register household tables for ingredient merges | S | M11 |

### Release 2

| # | Task | Title | Size | Milestone |
|---|---|---|---|---|
| 462 | [KIT-11-01](E11-per-person-profiles/KIT-11-01-member-rows.md) | Member rows from identity events | M | M16 |
| 463 | [KIT-11-02](E11-per-person-profiles/KIT-11-02-edit-rights.md) | UpdateMemberProfile and who may change what | M | M16 |
| 464 | [KIT-11-03](E11-per-person-profiles/KIT-11-03-daily-kcal-target.md) | Daily kcal target: limits, consent, deletion | M | M16 |
| 465 | [KIT-11-04](E11-per-person-profiles/KIT-11-04-per-person-restrictions.md) | Per-person diet, allergens and dislikes | M | M16 |
| 466 | [KIT-11-05](E11-per-person-profiles/KIT-11-05-display-and-family-flags.md) | Calories hidden and "eat lighter as a family" | S | M16 |
| 467 | [KIT-11-06](E11-per-person-profiles/KIT-11-06-member-visibility-and-snapshot.md) | GetMemberProfiles visibility and members in the planning snapshot | M | M16 |
| 468 | [KIT-13-03](E13-household-lifecycle/KIT-13-03-member-removed.md) | member_removed: delete the person's data | S | M16 |
| 496 | [KIT-12-01](E12-savings/KIT-12-01-movement-valuation.md) | Grams and value of each stock movement | S | M17 |
| 497 | [KIT-12-02](E12-savings/KIT-12-02-ledger-rules.md) | Ledger rules and monthly totals | M | M17 |
| 498 | [KIT-12-03](E12-savings/KIT-12-03-get-savings.md) | GetSavings and ListSavingsReports | S | M17 |
| 499 | [KIT-12-04](E12-savings/KIT-12-04-monthly-report-push.md) | Monthly savings push | S | M17 |

### Release 3

| # | Task | Title | Size | Milestone |
|---|---|---|---|---|
| 525 | [KIT-08-03](E08-ai-intake/KIT-08-03-text-intake.md) | Text and voice intake with a non-AI fallback | M | M19 |
| 526 | [KIT-08-05](E08-ai-intake/KIT-08-05-receipt-name-memory.md) | Receipt-name memory per household and chain | M | M19 |
| 529 | [KIT-08-04](E08-ai-intake/KIT-08-04-receipt-intake.md) | Receipt intake and line matching | L | M19 |
| 530 | [KIT-08-06](E08-ai-intake/KIT-08-06-receipt-scanned-event.md) | receipt.scanned event for shop prices | M | M19 |
| 536 | [KIT-12-05](E12-savings/KIT-12-05-share-card.md) | Shareable savings card data | S | M20 |

## Gates

| Task | Gate |
|---|---|
| [KIT-04-06](E04-barcode-lookup/KIT-04-06-open-food-facts-client.md) | licence check before production) \| |
| [KIT-04-06](E04-barcode-lookup/KIT-04-06-open-food-facts-client.md) | production enables `OFF_ENABLED` only after the licence check (features.md §4 open points, §8 risk "The open product database must be shared alike"). |
| [KIT-04-07](E04-barcode-lookup/KIT-04-07-open-food-facts-mapping.md) | licence check before production) \| |
| [KIT-04-08](E04-barcode-lookup/KIT-04-08-gs1-lookup.md) | GS1 access terms and price for Poland and Germany) \| |
| [KIT-04-08](E04-barcode-lookup/KIT-04-08-gs1-lookup.md) | stays disabled in production until access terms and price are agreed (features.md §4 open points). |
| [KIT-05-04](E05-unknown-products/KIT-05-04-contribute-to-off.md) | licence decision and an Open Food Facts app account) \| |
| [KIT-05-04](E05-unknown-products/KIT-05-04-contribute-to-off.md) | enabled only after the licence decision in features.md §4 open points and after the OFF account is set up. |
| [KIT-06-02](E06-nutrition/KIT-06-02-import-ciqual.md) | licence and credit wording for a paid app) \| |
| [KIT-06-03](E06-nutrition/KIT-06-03-import-bls.md) | terms for a commercial app) \| |

## Decisions

Settled answers (full list in [DECISIONS.md](../DECISIONS.md)); each task states its rule as `Decided (Qnnn)`.

| Task | Decision |
|---|---|
| [KIT-01-02](E01-platform/KIT-01-02-event-envelope-outbox.md) | Q015 — each stream entry has two fields, `type` (for metrics and skipping) and `envelope` (the full JSON). The identity (Bun) publisher uses the same layout (README "Events"). |
| [KIT-01-05](E01-platform/KIT-01-05-dead-letters-lag.md) | Q042 — no automatic resync in R1. The alert links to a runbook: ask identity to republish household snapshots (IDN-01-04) by hand. |
| [KIT-01-06](E01-platform/KIT-01-06-household-plan-copy.md) | Q004 — kitchen keeps no local copy of any household's plan. Package `internal/platform/plan` exposes `PlanOf(ctx, householdID) (Plan, error)`, which calls `IdentityService.GetEntitlements(household_id, actor_user_id = "")` (IDN-05-02) with a 300 ms deadline … |
| [KIT-02-02](E02-ingredient-catalogue/KIT-02-02-measures-and-units.md) | Q051 — `pinch` without its own weight counts 0.5 g; `to taste` counts 0 g and is reported as "not counted", never as missing. |
| [KIT-02-02](E02-ingredient-catalogue/KIT-02-02-measures-and-units.md) | Q050 — there is no default density. Millilitres of an ingredient without a density can't become grams and are reported as `no_weight` (nutrition then says "calories incomplete") until a density is added. |
| [KIT-02-07](E02-ingredient-catalogue/KIT-02-07-starter-catalogue-seed.md) | Q053 — the starter content is about 700 ingredients (everything the 200 starter dishes use plus common fridge items in Poland and Germany), drafted with AI help outside the product, then checked by the founder (names, allergens, diet flags) and by a native … |
| [KIT-02-08](E02-ingredient-catalogue/KIT-02-08-ai-translation-batch.md) | Q066 — machine names go live before review; the native-speaker check per language (business model) follows and marks them `reviewed`. |
| [KIT-02-09](E02-ingredient-catalogue/KIT-02-09-merge-ingredients.md) | Q056 — merges cannot be undone. If a merge was wrong, an admin creates a new ingredient and re-points what needs it. |
| [KIT-02-10](E02-ingredient-catalogue/KIT-02-10-admin-catalogue-rpcs.md) | Q054 — a change of allergens or diet flags writes the event `catalogue.ingredient_flags_changed` `{ingredient_ids}` (README event table), so planner re-derives the allergens of every dish that uses the ingredient at once (safety). |
| [KIT-02-10](E02-ingredient-catalogue/KIT-02-10-admin-catalogue-rpcs.md) | Q055 — `RetireIngredient` hides the ingredient from search and from new matches but keeps every reference (stock, recipes, products) working. |
| [KIT-03-01](E03-name-matching/KIT-03-01-name-normaliser.md) | Q052 — German compound words are not split ("Hähnchenschenkel" stays one token); Haiku resolves compounds in the AI path when text matching is unsure. |
| [KIT-03-03](E03-name-matching/KIT-03-03-match-names-deterministic.md) | Q057 — best ≥ 0.85 and best − second ≥ 0.15 → `matched`, method `trigram`, confidence = best, no review. Stored. The thresholds are config values, tuned from review data after launch. |
| [KIT-03-04](E03-name-matching/KIT-03-04-match-names-ai-pick.md) | Q064 — confidence bands: ≥ 0.8 → `matched`, no review; 0.5–0.8 → `matched`, `review = needed`; below 0.5 or "none" → `none`, `review = needed`. Stored with method `ai`. |
| [KIT-03-05](E03-name-matching/KIT-03-05-bulk-matching-batch.md) | Q065 — callers learn about results by asking again: offers calls `MatchNames` in `BULK` mode every 30 minutes for offers it holds back. There is no `catalogue.names_matched` event. |
| [KIT-03-06](E03-name-matching/KIT-03-06-match-review-queue.md) | Q058 — when an admin changes an existing `matched` answer, kitchen writes the event `catalogue.match_changed` `{language, key, from_ingredient_id, to_ingredient_id}` (README event table) so offers and planner fix references they stored. |
| [KIT-04-01](E04-barcode-lookup/KIT-04-01-gtin-parsing.md) | Q067 — canonical key = the GTIN-14 form (left-padded with zeros, the GS1 standard), stored in every table. Display form = 8, 12 or 13 digits, or 14 when the first digit isn't 0. |
| [KIT-04-03](E04-barcode-lookup/KIT-04-03-internal-codes.md) | Q068 — with no pattern, the item key is the first 7 digits (2-digit prefix + 5 item digits), the common European layout 2 + 5 + 5 + 1, so a new weight of a known item doesn't ask for a name again. |
| [KIT-04-03](E04-barcode-lookup/KIT-04-03-internal-codes.md) | Q069 — names for internal codes are kept per household, keyed by (household, country, item key), optionally with the chain the user picked (the same digits mean different things at different chains). They are never shared with other households. When two … |
| [KIT-04-04](E04-barcode-lookup/KIT-04-04-household-products.md) | Q070 — at most 5,000 products per household; beyond that the least recently scanned are deleted by a nightly job (a typical household scans a few hundred). |
| [KIT-04-07](E04-barcode-lookup/KIT-04-07-open-food-facts-mapping.md) | Q071 — OFF-derived products are stored in their own table `off_products`, never mixed into `products`, which keeps the share-alike ODbL database separate from Coolinary's catalogue and from household data (features.md §8). |
| [KIT-04-11](E04-barcode-lookup/KIT-04-11-multipacks.md) | Q072 — a multipack is stocked as one line: for a `PIECE` ingredient the quantity is N pieces; otherwise N × unit size in the base unit (6 × 0.5 l water → 3,000 ml). |
| [KIT-05-01](E05-unknown-products/KIT-05-01-name-product.md) | Q073 — at most 50 namings per household per day (redis-state counter `kitchen:naming:<household>:<date>`), then `RESOURCE_EXHAUSTED` until the next day. This stops abuse of the shared catalogue. |
| [KIT-05-02](E05-unknown-products/KIT-05-02-submissions-and-confirmation.md) | Q075 — profanity word lists per language (pl, de, en) in `internal/kitchen/names/profanity/<lang>.txt`; a name or brand containing a listed word is never proposed to other households. |
| [KIT-05-02](E05-unknown-products/KIT-05-02-submissions-and-confirmation.md) | Q074 — disputes. Three distinct fingerprints within 30 days without one reaching two households → `status = 'disputed'`, sent to admin review (KIT-05-05). A confirmed product reported wrong (`ReportProduct`) by three distinct households → `disputed`. A … |
| [KIT-05-03](E05-unknown-products/KIT-05-03-label-photo.md) | Q110 — no nutrition values are read from labels in R1; only name, brand and size. A misread number would break "AI never states a calorie number" in spirit. |
| [KIT-05-04](E05-unknown-products/KIT-05-04-contribute-to-off.md) | Q076 — each household can opt out. The kitchen profile setting `contribute_to_off` (KIT-10-01) is `true` by default. A product is sent only when at least two households whose setting is `true` (read at send time) submitted the confirmed name and brand. … |
| [KIT-06-01](E06-nutrition/KIT-06-01-nutrition-schema.md) | Q082 — "traces" is stored as 0 and "< x" as x / 2, each with a qualifier kept in `qualifiers` (common practice; keeps the numbers conservative). |
| [KIT-06-07](E06-nutrition/KIT-06-07-pack-nutrition.md) | Q083 — plausibility rules: each value 0–100 g; fat + carbohydrate + protein + fibre + salt ≤ 101 g; saturates ≤ fat; sugars ≤ carbohydrate; stated kcal within ±20% or ±20 kcal (whichever is larger) of 4 × protein + 4 × carbohydrate + 9 × fat + 2 × fibre. … |
| [KIT-06-08](E06-nutrition/KIT-06-08-calculation-rules.md) | Q084 — deep frying (features.md §5 step 3, "only the oil the food absorbs is counted"). A `frying_oil` line counts `counted_g` only; without `counted_g` the line is `missing: no_absorbed_amount`, so the dish is "calories incomplete" (planner then refuses to … |
| [KIT-06-08](E06-nutrition/KIT-06-08-calculation-rules.md) | Q085 — EU guidance on nutrition labelling rounding: kJ to the nearest 10; protein, carbohydrate, sugars, fat, fibre: ≥ 10 g → 1 g, 0.5–10 g → 0.1 g, < 0.5 g → "< 0.5 g"; saturates: ≥ 10 g → 1 g, |
| [KIT-07-01](E07-fridge-stock/KIT-07-01-stock-schema.md) | Q077 — a household holds at most 1,000 stock items, as a guard against abuse. |
| [KIT-07-02](E07-fridge-stock/KIT-07-02-storage-and-expiry-defaults.md) | Q078 — thawed food gets at most 2 days, or the fridge shelf life if shorter; |
| [KIT-07-03](E07-fridge-stock/KIT-07-03-add-stock-items.md) | Q079 — free-text items are allowed and matched later. |
| [KIT-07-05](E07-fridge-stock/KIT-07-05-list-stock.md) | Q103 — items past their date stay in the list and in the totals. Planner counts them as at home and marks the dish or plan line "check it first" (PLN-06-01, PLN-08); the status `expired` and `includes_expired` tell it which. Kitchen never drops them on its own. |
| [KIT-07-07](E07-fridge-stock/KIT-07-07-consume-stock.md) | Q080 — frozen food is assumed not to be cooked unless nothing else is left; within fridge and pantry the oldest date goes first. |
| [KIT-07-07](E07-fridge-stock/KIT-07-07-consume-stock.md) | Q081 — these thresholds. |
| [KIT-08-01](E08-ai-intake/KIT-08-01-intake-pipeline.md) | Q109 — the pre-tick threshold is 0.5. |
| [KIT-08-03](E08-ai-intake/KIT-08-03-text-intake.md) | Q219 — typed-sentence intake ships in release 3, together with voice, as features.md groups them; the gateway route is R3 (GW-06). |
| [KIT-08-04](E08-ai-intake/KIT-08-04-receipt-intake.md) | Q221 — the server keeps a 2-hour copy of the extracted lines (never the image), so prices can't be forged by the app. |
| [KIT-08-06](E08-ai-intake/KIT-08-06-receipt-scanned-event.md) | Q222 — prices are shared facts and do not need a household link; the README envelope allows an empty `household_id` for such events. |
| [KIT-09-02](E09-expiry-reminders/KIT-09-02-reminder-scheduler.md) | Q134 — 09:00 local time, in time to plan the day's meals. Notifier delays anything that falls in a user's quiet hours (NTF-03-02). |
| [KIT-10-01](E10-kitchen-profile/KIT-10-01-profile-record-and-defaults.md) | Q059 — Austria defaults like Germany; every other country gets English, EUR, Europe/Brussels and no home cuisine. Onboarding is never refused for a country. |
| [KIT-10-03](E10-kitchen-profile/KIT-10-03-kitchen-tools-and-presets.md) | Q061 — the kettle is a tool that dishes can require. |
| [KIT-10-06](E10-kitchen-profile/KIT-10-06-household-rhythm-and-meal-shares.md) | Q062 — bounds breakfast 10–40, main 25–60, other 10–40, snacks 0–20. |
| [KIT-10-06](E10-kitchen-profile/KIT-10-06-household-rhythm-and-meal-shares.md) | Q063 — one default, 25 / 40 / 25 / 10, for every country for now; a per-country default (features.md §5 open point) can be added later as config. |
| [KIT-11-02](E11-per-person-profiles/KIT-11-02-edit-rights.md) | Q192 — owner-only for accounts, so one adult cannot lock another. |
| [KIT-11-03](E11-per-person-profiles/KIT-11-03-daily-kcal-target.md) | Q193 — upper bound 5,000 kcal. |
| [KIT-11-03](E11-per-person-profiles/KIT-11-03-daily-kcal-target.md) | Q194 — the server's absolute child bounds are 500–4,000 kcal; a paediatric dietitian reviews them before R2 (features.md §5 open points) and may narrow them (config values). |
| [KIT-11-04](E11-per-person-profiles/KIT-11-04-per-person-restrictions.md) | Q195 — **after a downgrade to Free** (`household.plan_changed`, run by KIT-13-06) |
| [KIT-11-04](E11-per-person-profiles/KIT-11-04-per-person-restrictions.md) | Q230 — the health-protecting diet flags gluten-free, lactose-free and nut-free are merged the same way, as household diet flags; |
| [KIT-12-01](E12-savings/KIT-12-01-movement-valuation.md) | Q209 — items added by typing or barcode without a price add kilograms but no euros; euros are counted only where the price paid is known. |
| [KIT-12-02](E12-savings/KIT-12-02-ledger-rules.md) | Q211 — the 3-day window, the same window the app calls "soon". |
| [KIT-12-02](E12-savings/KIT-12-02-ledger-rules.md) | Q210 — counters are collected for all households. |
| [KIT-12-03](E12-savings/KIT-12-03-get-savings.md) | Q212, shared with GW-12-01): Free sees the current month's counters; reports and past months are Plus. |
| [KIT-12-04](E12-savings/KIT-12-04-monthly-report-push.md) | Q213 — the 1st of the month at 10:00 local time. |
| [KIT-13-02](E13-household-lifecycle/KIT-13-02-member-joined-with-fridge.md) | Q139 — the old household-level diet travels with them, since it was in fact that person's diet. |

