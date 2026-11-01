# ADM-03 Catalogue, matching and nutrition review

**Goal.** The human half of "AI proposes, rules check" for kitchen's catalogue: resolve names Haiku could not match
(decision 13), approve ingredient-to-nutrition-table matches, edit ingredients and approve translations, merge
duplicates, review user-named products, start the catalogue batch jobs with a cost confirmation, and maintain the
weighed-item barcode layouts.

**Flow.** kitchen fills review queues (from batch jobs and live matching) → admin lists them →
`KitchenAdminService.Resolve…` → kitchen stores the answer for everyone and publishes events where data changed
(`catalogue.nutrition_changed`, `catalogue.ingredient_merged`).

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [ADM-03-01](ADM-03-01-match-review-queue.md) | Name-match review queue | R1 | M | ADM-02-03, KIT-03-06 |
| [ADM-03-02](ADM-03-02-nutrition-match-review.md) | Nutrition match approval | R1 | M | ADM-02-03, KIT-06-06, KIT-02-10 |
| [ADM-03-03](ADM-03-03-ingredient-editor.md) | Ingredient editor and translation review | R1 | M | ADM-02-03, KIT-02-05, KIT-02-06, KIT-02-10 |
| [ADM-03-04](ADM-03-04-ingredient-merge.md) | Ingredient merge | R1 | S | ADM-03-03, KIT-02-09 |
| [ADM-03-05](ADM-03-05-product-review.md) | Unknown products review | R1 | M | ADM-02-03, KIT-05-05 |
| [ADM-03-06](ADM-03-06-catalogue-jobs.md) | Catalogue batch jobs | R1 | S | ADM-02-03, KIT-01-08, AI-02-04 |
| [ADM-03-07](ADM-03-07-weighed-code-patterns.md) | Weighed-item barcode patterns | R1 | S | ADM-02-03, KIT-04-03 |
