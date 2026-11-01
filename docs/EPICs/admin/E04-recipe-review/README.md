# ADM-04 Recipe import and dish review

**Goal.** Fill and curate planner's recipe library (decision 2) from both sources (decision 10): a hand-written import
file and AI draft batches. Every dish is reviewed in one editor, gets checked effort tags and translations, and the
starter library is checked by a dietitian (the launch gate). From R2, well-rated AI dishes from households
are promoted into the library.

**Flow.** import file (validate → commit) or AI draft batch (dry run → confirm cost) → dishes in `review` → editor
(ingredients, calories, rule check) → approve (= publish) → translations → dietitian check (starter dishes). Base pantries, no-cook extras and the
coverage report sit next to it (ADM-04-08).

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [ADM-04-01](ADM-04-01-import-upload.md) | Recipe import upload | R1 | M | ADM-01-05, ADM-02-03, PLN-03-01, PLN-03-03 |
| [ADM-04-02](ADM-04-02-ai-draft-batches.md) | AI dish draft batches | R1 | M | ADM-02-03, PLN-03-04, AI-04, AI-09 |
| [ADM-04-03](ADM-04-03-review-queue-and-editor.md) | Dish review queue and editor | R1 | L | ADM-04-01 or ADM-04-02, PLN-03-07, PLN-01-06 |
| [ADM-04-04](ADM-04-04-publish-retire-translate.md) | Publish, retire and translate dishes | R1 | M | ADM-04-03, PLN-03-04/05/07, AI-09 |
| [ADM-04-05](ADM-04-05-dietitian-sign-off.md) | Dietitian sign-off | R1 | S | ADM-02-02, ADM-04-03, PLN-03-07/08 |
| [ADM-04-06](ADM-04-06-effort-tag-review.md) | Effort tag suggestions review | R1 | S | ADM-04-03, PLN-03-04/06/07, AI-09 |
| [ADM-04-07](ADM-04-07-ai-dish-promotion.md) | Promote well-rated AI dishes to the library | R2 | S | ADM-04-03, PLN-11-02/03 |
| [ADM-04-08](ADM-04-08-library-settings.md) | Library settings: base pantries, no-cook extras, coverage | R1 | M | ADM-04-03, PLN-16-01, PLN-12-04, PLN-03-09 |
