# PLN-05-05 Ratings

| | |
|---|---|
| Service | planner |
| Epic | [PLN-05 Library browse](README.md) |
| Release | R1 |
| Size | S |
| Status | todo |
| Depends on | PLN-02-01 |
| Unblocks | PLN-05-01 (POPULAR), PLN-10-03, PLN-11-01 |

## Goal

People rate dishes 1–5 stars, mostly right after cooking. Ratings sort the library and pick which household AI dishes
are worth adding to the library (PLN-11).

## Business rules

1. One rating per account per dish, 1–5 whole stars, changeable; `RateDish(dish_id, stars)`; stars 0 removes it.
2. Who: any account in the household that can read the dish (anonymous accounts too). The rating records the
   `dish_version` rated (the version live at the time) and the household.
3. Aggregates per dish: count and sum across all households, updated in the same transaction, used by POPULAR
   (PLN-05-01 rule 4).
4. A household's own rating of a dish = average of its members' ratings, used by PLN-11.
5. Ratings of a household AI dish count only towards that dish (private), never towards a library dish.
6. On `account.deleted` the account's ratings are anonymised (account id set to NULL) so aggregates stay; on
   `household.deleted` ratings of that household's AI dishes are deleted with the dishes (PLN-21).

## Contract

```proto
rpc RateDish(RateDishRequest) returns (RateDishResponse);
message RateDishRequest { string household_id = 1; string user_id = 2; string dish_id = 3; int32 stars = 4; }
message RateDishResponse { double average = 1; int32 count = 2; }
```

## Data

`planner/migrations/000NN_ratings.sql`

```sql
CREATE TABLE ratings (
    dish_id       uuid NOT NULL REFERENCES dishes (id) ON DELETE CASCADE,
    user_id       text,                         -- NULL after account deletion
    household_id  uuid NOT NULL,
    dish_version  int  NOT NULL,
    stars         smallint NOT NULL CHECK (stars BETWEEN 1 AND 5),
    rated_at      timestamptz NOT NULL DEFAULT now(),
    id            bigserial PRIMARY KEY
);
CREATE UNIQUE INDEX ratings_user_dish ON ratings (user_id, dish_id) WHERE user_id IS NOT NULL;
CREATE TABLE dish_rating_totals (
    dish_id  uuid PRIMARY KEY REFERENCES dishes (id) ON DELETE CASCADE,
    count    int NOT NULL DEFAULT 0,
    sum      int NOT NULL DEFAULT 0
);
```

## Errors

| Case | gRPC status | Message |
|---|---|---|
| Stars outside 0–5 | INVALID_ARGUMENT | `stars must be 1-5, or 0 to remove` |

## Acceptance criteria

- [ ] Rating 4 then 2 leaves one rating of 2 and totals count 1, sum 2.
- [ ] Deleting the account keeps totals unchanged and removes the user id.
