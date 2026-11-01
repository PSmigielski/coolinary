# PLN-18 Plan from deals

Release R2, Plus only. Compose dishes from the fridge plus this week's offers, preferring the nearest shops
(features.md §2 "Flyers and offers" with AI). Coolinary stays fridge-first: deals only fill the gaps left after using
what is at home (features.md §8 "Flyva" risk).

## Flow

1. Collect this week's offers at the household's shops that fit the household's limits (PLN-18-01).
2. Compose with AI through `AI-14 (ComposeFromDeals)` (PLN-18-02), or without AI from the library (PLN-18-03); rule
   check, portions.
3. Accepting the plan puts the offer items on the list with pinned badges (PLN-18-04).

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [PLN-18-01](PLN-18-01-deal-candidates.md) | Deal candidates for a household | R2 | M | OFF-04 (ListMyShops), OFF-09 (FindOffersForIngredients), PLN-13-01 |
| [PLN-18-02](PLN-18-02-ai-plan-from-deals.md) | Plan from deals with AI | R2 | L | PLN-18-01, AI-14 (ComposeFromDeals), PLN-07, PLN-08, PLN-12-06 |
| [PLN-18-03](PLN-18-03-non-ai-plan-from-deals.md) | Plan from deals without AI | R2 | M | PLN-18-01, PLN-13-06 |
| [PLN-18-04](PLN-18-04-deals-plan-to-list.md) | Deals plan to the shopping list | R2 | S | PLN-18-02, PLN-15-05, PLN-17-01 |
