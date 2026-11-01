# KIT-12-02 Ledger rules and monthly totals

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-12 Savings](README.md) |
| Release | R2 |
| Size | M |
| Depends on | KIT-12-01, KIT-07-07 (cook events) |
| Unblocks | KIT-12-03, KIT-12-04, KIT-12-05 |

## Goal

Turn movements into honest monthly figures: food rescued before it went off, money saved on offers, food thrown
away, and meals cooked.

## Business rules

1. **Used before expiry ("rescued"):** a `cooked` or `used` movement of an item with a date, where the movement's
   local date is on or before the date and at most 3 days before it. It adds the movement's grams and value.
   Everyday use of food far from its date does not count, or every meal would be "saved". Decided (Q211): the
   3-day window, the same window the app calls "soon".
2. **Bought on offer:** an `added` movement with `saving_minor` adds that saving.
3. **Thrown away:** a `thrown_away` movement adds grams and value to the wasted counters; shown next to the savings,
   never subtracted from them silently.
4. **Meals cooked:** each new `cook_event_id` (KIT-07-07, not replays) adds 1.
5. `mistake`, `adjusted`, `moved_in`, `moved_out` and `added` without saving change nothing.
6. Month = the movement's local date (household timezone) truncated to the month.
7. Totals are kept per (household, month, currency) and updated in the same transaction as the movement, so a
   figure never runs ahead of or behind the stock.
8. Counters are collected for every household, Free and anonymous included, so a household that upgrades sees its
   history. Decided (Q210): counters are collected for all households.
9. Money is never converted between currencies; a household that moved from PLN to EUR has two rows for that month.

## Contract

No RPC.

## Data

`kitchen/migrations/000NN_savings_monthly.sql`

```sql
-- +goose Up
CREATE TABLE savings_monthly (
    household_id   uuid NOT NULL,
    month          date NOT NULL CHECK (extract(day FROM month) = 1),
    currency       text NOT NULL DEFAULT '',      -- '' for rows with grams only
    rescued_g      numeric(14,1) NOT NULL DEFAULT 0,
    rescued_minor  bigint NOT NULL DEFAULT 0,
    offer_minor    bigint NOT NULL DEFAULT 0,
    wasted_g       numeric(14,1) NOT NULL DEFAULT 0,
    wasted_minor   bigint NOT NULL DEFAULT 0,
    meals_cooked   integer NOT NULL DEFAULT 0,
    PRIMARY KEY (household_id, month, currency)
);

-- +goose Down
DROP TABLE savings_monthly;
```

## Errors

None of its own.

## Acceptance criteria

- [ ] Milk dated 10 March, cooked on 8 March → rescued; the same milk cooked on 2 March → not rescued; on 11 March
      (after the date) → not rescued.
- [ ] A ticked offer with a 2.00 saving adds 200 to `offer_minor`.
- [ ] Throwing away 300 g of spinach worth 4.50 adds to the wasted counters only.
- [ ] A replayed `cook_event_id` does not add a second meal.
- [ ] A movement at 00:30 on 1 April Warsaw time counts for April even though it is 31 March in UTC.
