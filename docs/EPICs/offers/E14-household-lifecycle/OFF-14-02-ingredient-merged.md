# OFF-14-02 Follow merged catalogue ingredients

| | |
|---|---|
| Service | offers |
| Epic | [OFF-14 Household lifecycle](README.md) |
| Release | R2 |
| Size | S |
| Depends on | OFF-01-02, KIT-02 (catalogue.ingredient_merged) |
| Unblocks | — |

## Goal

When kitchen merges two catalogue ingredients, offers, prices and searches keep working under the surviving id.

## Business rules

1. On `catalogue.ingredient_merged` (`from_ingredient_id`, `into_ingredient_id`), in one transaction: rewrite
   `ingredient_id` in `offers`, `price_observations` and `reference_prices`.
2. Reference prices: if both ids had a value for the same country and unit basis, the surviving id keeps the value
   with more observations; the next weekly run recomputes it anyway.
3. Offers held only for `no_match` are not affected (they have no id).
4. Browse caches of affected editions are invalidated.

## Contract

Consumes `catalogue.ingredient_merged`.

## Data

No schema change.

## Errors

| Case | Behaviour |
|---|---|
| Unknown `into` id | Rewrite anyway (kitchen is the owner); logged |

## Acceptance criteria

- [ ] After merging "spring onion" into "green onion", `FindOffersForIngredients` for the green onion id returns the
      offers that were matched to spring onion.
