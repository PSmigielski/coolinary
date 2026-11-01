# OFF-09 Offer browse, search and lookup by ingredient

| | |
|---|---|
| Service | offers |
| Release | R2 |
| Goal | Users browse and search this week's and next week's offers of their shops by shop, category and end date, with conditions and a link to the original flyer; planner finds offers for list ingredients. Free sees 2 shops, Plus all its shops (features.md §2 "Flyers and offers", §6). |

## Flow

1. The phone sends shop ids nearest first (README D13); offers maps shops to editions and returns their visible
   offers, each tagged with the requested shops that carry it, in the phone's order.
2. Text search goes through catalogue ingredient names in the user's language (kitchen), plus product-name search.
3. planner asks for offers by ingredient for list badges and plan-from-deals.
4. `BrowseOffers` replaces the skeleton's `ListOffers` (postcode only, no shops, no editions).

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [OFF-09-01](OFF-09-01-browse-offers.md) | Browse offers of my shops | R2 | L | OFF-07-07, OFF-04-01, OFF-03-04 |
| [OFF-09-02](OFF-09-02-search-offers.md) | Search offers by text | R2 | M | OFF-09-01, KIT-02 (SearchIngredients) |
| [OFF-09-03](OFF-09-03-offer-detail.md) | Offer detail | R2 | S | OFF-09-01 |
| [OFF-09-04](OFF-09-04-offer-categories.md) | Offer categories with counts | R2 | S | OFF-09-01 |
| [OFF-09-05](OFF-09-05-find-offers-for-ingredients.md) | Find offers for ingredients (planner) | R2 | M | OFF-07-07, OFF-04-01 |
