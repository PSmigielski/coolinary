# PLN-05-04 Favourites

| | |
|---|---|
| Service | planner |
| Epic | [PLN-05 Library browse](README.md) |
| Release | R1 |
| Size | S |
| Depends on | PLN-02-01 |
| Unblocks | PLN-05-01 (scope FAVOURITES), GW-07 |

## Goal

Mark dishes as favourites and find them again.

## Business rules

1. Favourites belong to the household: one shared list for everyone in it. Decided (Q101).
2. `SetFavourite(dish_id, favourite)` is idempotent. Any library dish (published or retired) or the household's own
   AI dish can be a favourite; another household's AI dish → NOT_FOUND.
3. At most 500 favourites per household; the 501st → FAILED_PRECONDITION.
4. Favourites point at the dish, not a version: they always show the current version.
5. Favourites raise a dish in suggestions and AI candidate lists (PLN-06-02, PLN-08-02).
6. Removed when the household is deleted (PLN-21).

## Contract

```proto
rpc SetFavourite(SetFavouriteRequest) returns (SetFavouriteResponse);
message SetFavouriteRequest { string household_id = 1; string dish_id = 2; bool favourite = 3; }
message SetFavouriteResponse {}
```

## Data

`planner/migrations/000NN_favourites.sql`

```sql
CREATE TABLE favourites (
    household_id uuid NOT NULL,
    dish_id      uuid NOT NULL REFERENCES dishes (id) ON DELETE CASCADE,
    created_at   timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (household_id, dish_id)
);
```

## Errors

| Case | gRPC status / behaviour |
|---|---|
| Missing `household_id` or `dish_id` | `INVALID_ARGUMENT` |
| Unknown dish, or another household's AI dish | `NOT_FOUND` |
| 501st favourite | `FAILED_PRECONDITION`, `x-error-code: favourites_limit` |
| Setting a state that already holds (favourite twice, unfavourite a non-favourite) | success, no change |

## Acceptance criteria

- [ ] Setting the same favourite twice leaves one row; unsetting a non-favourite succeeds.
- [ ] A retired favourite shows in scope FAVOURITES with `available=false`.
