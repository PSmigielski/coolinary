# PLN-05-06 Cooked history

| | |
|---|---|
| Service | planner |
| Epic | [PLN-05 Library browse](README.md) |
| Release | R1 |
| Size | S |
| Status | todo |
| Depends on | PLN-02-01 |
| Unblocks | PLN-10-02 (writes it), PLN-06-02, PLN-08-02, PLN-11-01 |

## Goal

A household sees what it cooked and when; planner uses it to avoid repeating dishes too soon.

## Business rules

1. A history entry = household, dish, dish version, servings (the pot), cooked at (UTC) and local date, the account
   that tapped "cooked it", optional plan meal id. Written only by `MarkCooked` (PLN-10-02).
2. `ListCookedHistory(household, from?, to?, page)` newest first, page size ≤ 50.
3. Entries are kept 2 years, then deleted by a monthly job. Decided (Q102).
4. Derived signals used elsewhere: `last_cooked_on(dish)` and `cooked_count(dish)` per household.
5. History never changes when a dish gets a new version or is retired; it shows the version cooked.

## Contract

```proto
rpc ListCookedHistory(ListCookedHistoryRequest) returns (ListCookedHistoryResponse);
message CookedEntry { string id = 1; string dish_id = 2; int32 dish_version = 3; string name = 4;
                      double servings = 5; string cooked_on = 6; string plan_meal_id = 7; }
```

## Data

`planner/migrations/000NN_cooked_history.sql`

```sql
CREATE TABLE cooked_history (
    id            uuid PRIMARY KEY,               -- client-generated cooked_id, PLN-10-02
    household_id  uuid NOT NULL,
    dish_id       uuid NOT NULL REFERENCES dishes (id),
    dish_version  int  NOT NULL,
    servings      numeric(5,2) NOT NULL,
    cooked_at     timestamptz NOT NULL,
    cooked_on     date NOT NULL,                  -- household local date
    user_id       text,
    plan_meal_id  uuid
);
CREATE INDEX cooked_history_household ON cooked_history (household_id, cooked_at DESC);
CREATE INDEX cooked_history_dish ON cooked_history (household_id, dish_id, cooked_on DESC);
```

## Errors

| Case | gRPC status / behaviour |
|---|---|
| `from` after `to`, or a date not in `YYYY-MM-DD` | `INVALID_ARGUMENT` |
| Page size above 50 | clamped to 50 |
| Invalid or expired page token | `INVALID_ARGUMENT` (`bad_page_token`) |
| No entries | success with an empty list |

## Acceptance criteria

- [ ] After cooking a dish twice, `cooked_count` is 2 and `last_cooked_on` is the later local date.
- [ ] An entry at 23:30 UTC for a household in Europe/Warsaw (UTC+1 in winter) has the next day as `cooked_on`.
