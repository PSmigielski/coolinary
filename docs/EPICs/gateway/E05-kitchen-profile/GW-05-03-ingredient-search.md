# GW-05-03 `GET /v1/catalogue/ingredients` (autocomplete)

| | |
|---|---|
| Service | gateway |
| Epic | [GW-05 Kitchen profile and catalogue routes](README.md) |
| Release | R1 |
| Size | S |
| Depends on | GW-01-01, KIT-02 (`SearchIngredients`, `GetIngredients`) |
| Unblocks | MOB-03 (fridge typing, staples, dislikes, manual list lines) |

## Goal

Typing into the fridge, staples, dislikes or the shopping list autocompletes from the ingredient catalogue in the
user's language (features.md §2 "Fridge and pantry").

## Business rules

1. `q` is trimmed, 2–50 characters; shorter → 200 with an empty list (no kitchen call), so the app can call on every
   keystroke without errors.
2. `limit` 1–20, default 10. Calls kitchen `SearchIngredients(q, language, limit)`.
3. Results cached in `redis-cache` for 10 minutes per `language:normalised q:limit` (normalised = lowercase, NFKC,
   collapsed spaces). Not per household: the catalogue is shared.
4. `GET /v1/catalogue/ingredients/{id}` returns one ingredient (`GetIngredients` with one id) for showing a stored
   item's name in another language.
5. Rate limit: the `api` bucket only (300/min covers fast typing).

## Contract

`GET /v1/catalogue/ingredients?q=pomid&limit=10` — `auth: true`, tag `catalogue`.

```json
200 { "items": [ { "id": "ing_tomato", "name": "pomidor", "category": "vegetables",
                   "defaultUnit": "piece", "defaultStorage": "fridge" } ] }
```

`GET /v1/catalogue/ingredients/{id}` → 200 one item, 404 `not_found`.

## Data (redis-cache)

| Key | Value | TTL |
|---|---|---|
| `gw:ingq:<language>:<sha1(q)>:<limit>` | response JSON | 10 min |

## Errors

| Case | HTTP | `code` |
|---|---|---|
| `q` longer than 50 | 400 | `validation` |
| Unknown id | 404 | `not_found` |

## Acceptance criteria

- [ ] `q=p` returns `{ "items": [] }` and kitchen is not called.
- [ ] `q=pomid` in Polish returns `pomidor` first (kitchen ranking).
- [ ] `q=Pomid ` and `q=pomid` share one cache entry.
