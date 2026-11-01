# KIT-02-05 SearchIngredients (autocomplete)

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-02 Ingredient catalogue](README.md) |
| Release | R1 |
| Size | M |
| Depends on | KIT-02-04 |
| Unblocks | GW-05 (catalogue routes), KIT-07 (add by typing), KIT-10 (staples, dislikes) |

## Goal

Typing two or three letters suggests the right ingredient in the user's language, fast enough for every keystroke
(features.md §2 "Add items by typing, with autocomplete from an ingredient catalogue").

## Business rules

1. Input: `query` (1–60 characters), `language` (from the request, else `x-language`, else `en`), `limit` (default
   10, maximum 25), optional `categories` filter.
2. The query is normalised like names (KIT-03-01). A normalised query shorter than 2 characters returns no results.
3. Search covers names in the requested language and in `en`. Only `active` ingredients are returned.
4. Ranking, best first:
   1. exact match of a primary name, then of a synonym or plural;
   2. a name that starts with the query, or any word in it that starts with the query;
   3. trigram `word_similarity(query, normalized)` ≥ 0.3, highest first.
   Within a tier: higher `popularity`, then shorter name, then alphabetical.
5. One result per ingredient: the best-ranked name wins. Each result returns the display name in the requested
   language (KIT-02-04 rule 6) and `matched_name` when the hit was a synonym ("courgette" → "zucchini").
6. `popularity` counts how often households add the ingredient to stock (KIT-07 increments it, at most once per
   household per ingredient per day).
7. Responses are cached in redis-cache under `kitchen:search:<catalogue_version>:<lang>:<limit>:<query>` for
   10 minutes.
8. p95 latency ≤ 50 ms with 10,000 ingredients × 3 languages in the database.

## Contract

`contracts/coolinary/kitchen/v1/catalogue.proto`

```proto
rpc SearchIngredients(SearchIngredientsRequest) returns (SearchIngredientsResponse);
message SearchIngredientsRequest { string query = 1; string language = 2; int32 limit = 3; repeated Category categories = 4; }
message SearchIngredientsResponse { repeated IngredientHit hits = 1; int64 catalogue_version = 2; }
message IngredientHit { string ingredient_id = 1; string name = 2; string matched_name = 3; Category category = 4;
  coolinary.common.v1.Unit base_unit = 5; bool translated = 6; }
```

## Data

No schema change (uses KIT-02-04 indexes and `ingredients.popularity`).

## Errors

| Case | Status | Message |
|---|---|---|
| Query empty or over 60 characters | `INVALID_ARGUMENT` | `query must be 1-60 characters` |
| `limit` over 25 | clamped to 25 | — |

## Acceptance criteria

- [ ] "pom" in `pl` returns "pomidor" before "pomarańcza" when tomatoes are more popular.
- [ ] "zucch" in `en` returns zucchini; "courg" returns zucchini with `matched_name = courgette`.
- [ ] "hähn" in `de` and "hahn" both find Hähnchen (folded matching).
- [ ] Merged and retired ingredients never appear.
- [ ] Load test: p95 ≤ 50 ms at 50 requests/second.
