# GW-08-04 Cuisine base pantry

| | |
|---|---|
| Service | gateway |
| Epic | [GW-08 Shopping list routes](README.md) |
| Release | R1 |
| Size | S |
| Depends on | GW-08-01, PLN-16 (`GetBasePantry`, `AddBasePantryToList`) |
| Unblocks | MOB-03 |

## Goal

Each cuisine has a "base pantry" checklist (Asian: soy sauce, rice, ginger, sesame oil; Italian: olive oil, passata,
pasta, hard cheese); one tap adds the missing base items to the list (features.md §2 "Cuisine choice").

## Business rules

1. `GET /v1/cuisines/{cuisine}/base-pantry` → planner `GetBasePantry(household_id, cuisine)`: items with `atHome`
   (from stock and staples) and `onList`.
2. `POST /v1/list/base-pantry` body `{ "cuisine" }` → `AddBasePantryToList(household_id, member_id, cuisine,
   client_mutation_id)`: adds only items neither at home nor on the list; lines get source `cuisine_base`.
   Idempotency required (the key is the `client_mutation_id`). Response: list `version` and the added lines (empty
   when nothing is missing).
3. `cuisine` codes from `GET /v1/profile/options`.

## Contract

| Method | Path | 2xx |
|---|---|---|
| GET | `/v1/cuisines/{cuisine}/base-pantry` | 200 `{ "cuisine": "asian", "items": [ { "ingredientId", "name", "atHome": true, "onList": false } ] }` |
| POST | `/v1/list/base-pantry` | 200 `{ "version", "added": [line] }` |

Tag `list`, `auth: true`.

## Data

No storage.

## Errors

| Case | HTTP | `code` |
|---|---|---|
| Unknown cuisine | 404 | `not_found` |

## Acceptance criteria

- [ ] With soy sauce and rice at home, adding the Asian base adds only ginger and sesame oil.
- [ ] Calling it again adds nothing.
