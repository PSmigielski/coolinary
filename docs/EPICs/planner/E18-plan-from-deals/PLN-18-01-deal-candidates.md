# PLN-18-01 Deal candidates for a household

| | |
|---|---|
| Service | planner |
| Epic | [PLN-18 Plan from deals](README.md) |
| Release | R2 |
| Size | M |
| Depends on | OFF-04 (ListMyShops), OFF-09 (FindOffersForIngredients), PLN-13-01 |
| Unblocks | PLN-18-02, PLN-18-03, PLN-20-01 |

## Goal

Build the set of offers a deals plan may use: current offers at the household's shops, matched to catalogue
ingredients, that fit the household's limits, nearest shops first.

## Business rules

1. Shops: `ListMyShops` order (nearest first, set by the phone, D13). Free households are refused before this step
   (PLN-18-02).
2. Ingredients asked for: the union of ingredients of published library dishes that pass the household's shared
   limits for the planned meals (PLN-13-01), capped at 500 ids, most-used first.
3. Offers: `FindOffersForIngredients(ids, shops, plan date)`; an offer counts if it is valid on the plan date.
4. Each offer's ingredient must pass the hard limits of everyone eating (allergens, diet, dislikes), using
   `GetIngredients` flags.
5. Rank: shop position (nearest first), then discount (PLN-17-01 rule 5) descending, then unit price ascending. Keep
   the top 40, and at most 15 from one shop.
6. Ingredients already at home in enough quantity for one standard dish are not "needed" and their offers are kept
   only if fewer than 40 candidates remain (fridge first).
7. Output per candidate: offer id, ingredient id, shop id, shop position, price, pack, unit price, discount,
   conditions, valid to.

## Contract

Internal (`deals.Candidates(household, date, meals) []DealCandidate`).

## Data

No schema change. Result cached in redis-cache for 30 minutes per (household, date, offers epoch).

## Errors

| Case | Behaviour |
|---|---|
| offers down | `UNAVAILABLE` from the calling RPC |
| No shops set | `FAILED_PRECONDITION` `no_shops` |

## Acceptance criteria

- [ ] Nut-free household: no nut offer appears.
- [ ] With Lidl first and Kaufland second, equal discounts rank Lidl first.
- [ ] Never more than 15 candidates from one shop.
