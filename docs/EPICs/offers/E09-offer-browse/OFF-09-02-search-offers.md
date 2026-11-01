# OFF-09-02 Search offers by text

| | |
|---|---|
| Service | offers |
| Epic | [OFF-09 Offer browse, search and lookup by ingredient](README.md) |
| Release | R2 |
| Size | M |
| Depends on | OFF-09-01, KIT-02 (SearchIngredients) |
| Unblocks | GW-10 |

## Goal

A user types "chicken" in Germany, or "kurczak" in Poland, and finds the matching offers, whatever language the flyer
is printed in.

## Business rules

1. `BrowseOffers` with `query` (2 to 60 characters, trimmed) searches within the same shops, view and filters.
2. Step 1: kitchen `SearchIngredients(query, language = x-language)` returns up to 10 ingredient ids; offers matched
   to those ingredients are hits (score 1.0 for the first ingredient, minus 0.05 per rank).
3. Step 2: trigram similarity between the query and the offer's product name and brand, threshold 0.3; score =
   similarity × 0.9.
4. Hits from both steps are merged (an offer keeps its best score), ordered by score, then by the OFF-09-01 order.
5. kitchen unavailable: step 2 only, and the response sets `partial = true`.
6. Queries are not logged and not stored.
7. Search results are not cached.

## Contract

Adds to `BrowseOffersResponse`: `bool partial = 5;`.

## Data

No schema change (`offers_name_trgm` index from OFF-07-01).

## Errors

| Case | Behaviour |
|---|---|
| Query under 2 or over 60 characters | `INVALID_ARGUMENT` |
| kitchen down | `partial = true`, not an error |

## Acceptance criteria

- [ ] With `x-language: en` and query "chicken thigh", a German offer "Hähnchenschenkel 1 kg" matched to chicken thigh
      is the first hit.
- [ ] Query "Milka" finds Milka offers by brand even if their ingredient is "chocolate".
- [ ] With kitchen stopped, the search still answers with `partial = true`.
