# GW-07-02 Favourites, ratings, history

| | |
|---|---|
| Service | gateway |
| Epic | [GW-07 Recipe library, plans and cooking routes](README.md) |
| Release | R1 |
| Size | S |
| Status | todo |
| Depends on | GW-01-03, PLN-05 (`SetFavourite`, `RateDish`, `ListCookedHistory`) |
| Unblocks | MOB-03 |

## Goal

Household members mark favourites, rate dishes (ratings also decide which AI dishes join the library, PLN-11), and
see what they cooked before (features.md §2 "Dishes").

## Business rules

1. Favourites are per household. `PUT` sets, `DELETE` clears; both idempotent by nature (setting twice is fine).
2. Rating is per account and dish: `PUT /v1/dishes/{dishId}/rating` body `{ "stars": 1–5 }` calls planner
   `RateDish(household_id, user_id, dish_id, stars)`; `DELETE` calls it with `stars = 0`, which removes the rating
   (PLN-05-05 rule 1).
   Household-private AI dishes can be rated too.
3. History lists "cooked it" entries newest first, paged (default 20): dish, date, meal, portions, who marked it.
4. Calls planner `SetFavourite`, `RateDish`, `ListCookedHistory`.

## Contract

| Method | Path | Body | 2xx |
|---|---|---|---|
| PUT | `/v1/dishes/{dishId}/favourite` | — | 204 |
| DELETE | `/v1/dishes/{dishId}/favourite` | — | 204 |
| PUT | `/v1/dishes/{dishId}/rating` | `{ "stars": 4 }` | 200 `{ "stars": 4, "householdAverage": 4.5 }` |
| DELETE | `/v1/dishes/{dishId}/rating` | — | 204 |
| GET | `/v1/history?pageSize=&pageToken=` | — | 200 `{ "entries": [ { "id", "dishId", "dishName", "cookedOn", "meal", "portions", "byMemberId" } ], "nextPageToken" }` |

Tag `recipes`, `auth: true`, idempotency optional.

## Data

No storage.

## Errors

| Case | HTTP | `code` |
|---|---|---|
| Stars outside 1–5 | 400 | `validation` |
| Unknown dish | 404 | `not_found` |

## Acceptance criteria

- [ ] Favourite then `GET /v1/dishes?favourites=true` lists the dish.
- [ ] Two members rating 4 and 5 give `householdAverage: 4.5`.
- [ ] After "cooked it", the history's first entry is that dish.
