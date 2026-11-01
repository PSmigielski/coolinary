# AI-11-05 RewriteDish: RESTYLE

| | |
|---|---|
| Service | ai |
| Epic | [AI-11 Swaps, wishes, tool rewrite, restyle](README.md) |
| Release | R3 |
| Size | S |
| Depends on | AI-11-04 |
| Unblocks | PLN-09 (restyle) |

## Goal

"Make this Italian-style": keep the dish's core and stock use, change seasoning and method to another cuisine
(features.md §2 "Cuisine choice", with AI, per request; release 3).

## Business rules

1. Same RPC and job as AI-11-04 with `mode=RESTYLE` and `target_cuisine` (a cuisine code the kitchen allows).
2. The main ingredients (role `MAIN` and `BASE`) stay; seasonings, sauces and garnishes may change, using only
   stock, staples or at most 2 `buyable_ids` items (listed with `to_buy` semantics in `changes=INGREDIENT`).
3. `cuisine` of the output = `target_cuisine`; name changes to fit ("Pierogi with sage butter and parmesan").
4. Hard limits as always; `not_possible=true` with reason `CUISINE` when the core can't fit the target cuisine.
5. Validation: AI-10-03 on the new dish; MAIN/BASE ids unchanged; at most 2 new non-stock ingredients.
6. Cache 1 hour (depends on stock).

## Contract

Uses `RewriteDishRequest.target_cuisine` from AI-11-04.

## Data

No schema change.

## Errors

| Case | gRPC status |
|---|---|
| Input limits broken, unknown ids or codes in the request | `INVALID_ARGUMENT`, `ai_bad_request` |
| Output still invalid after repair | `UNAVAILABLE`, `ai_invalid_output` |
| Others | AI-01-04; the caller falls back to its non-AI result |
| Core can't fit the target cuisine | OK, `not_possible=true`, reason `CUISINE` |

## Acceptance criteria

- [ ] A Polish potato dish restyled to Italian keeps potatoes as MAIN, adds no more than 2 items to buy, and has `cuisine` = Italian.
- [ ] Restyling into a cuisine the household has not ticked is still allowed when the user asked (`target_cuisine` is explicit).
