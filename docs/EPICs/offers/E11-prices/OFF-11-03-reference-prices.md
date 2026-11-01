# OFF-11-03 Reference prices per country

| | |
|---|---|
| Service | offers |
| Epic | [OFF-11 Prices](README.md) |
| Release | R2 |
| Size | M |
| Depends on | OFF-11-01, OFF-01-05 |
| Unblocks | OFF-11-04, OFF-07-05 (plausibility rule) |

## Goal

Every list line has some price for an "about" total, and shops without an offer on that line are treated the same, so
the reference price never sways which shop wins (features.md §2 "a reference price that is the same for every shop").

## Business rules

1. Source (Decided, Q174): computed every Sunday at 04:00 UTC per country, ingredient and unit basis as
   the median unit price over the last 90 days of non-outlier observations, using regular prices where flyers show
   them, else flyer prices, plus receipt prices from R3; at least 3 observations from at least 2 chains are needed.
   An admin-imported value (CSV) overrides the computed one.
2. Without enough observations and without an imported value, the ingredient has no reference price; `GetPrices`
   then returns "price unknown" for shops where nothing else is known.
3. Values are stored with the date computed and the number of observations behind them.
4. `ImportReferencePrices`: CSV with `country_code, ingredient_id, unit_basis, price_minor`; validated row by row;
   rows with unknown ingredients are reported, not imported. Imported values stay until replaced or deleted.
5. A computed value that moves more than 50% from the previous week is kept, but listed in an info alert for review.

## Contract

`admin.proto`

```proto
rpc ImportReferencePrices(ImportReferencePricesRequest) returns (ImportReferencePricesResponse); // csv bytes -> imported, rejected rows
rpc ListReferencePrices(ListReferencePricesRequest) returns (ListReferencePricesResponse);
```

## Data

`offers/migrations/000NN_reference_prices.sql`

```sql
CREATE TABLE reference_prices (
    country_code   text NOT NULL,
    ingredient_id  text NOT NULL,
    unit_basis     text NOT NULL CHECK (unit_basis IN ('kg', 'l', 'piece')),
    unit_price_minor bigint NOT NULL CHECK (unit_price_minor > 0),
    currency       text NOT NULL,
    source         text NOT NULL CHECK (source IN ('computed', 'imported')),
    observations   int NOT NULL DEFAULT 0,
    computed_on    date NOT NULL,
    PRIMARY KEY (country_code, ingredient_id, unit_basis)
);
```

## Errors

| Case | Behaviour |
|---|---|
| Invalid CSV row | Listed in `rejected`; the other rows are imported |
| CSV over 5 MB | `INVALID_ARGUMENT` |

## Acceptance criteria

- [ ] With 5 Polish flyer observations for basmati rice from 2 chains (regular prices 6.99, 7.49, 7.99, 8.49, 8.99 per
      kg), the reference price is 7.99 PLN per kg.
- [ ] With 3 observations from one chain only, no reference is computed.
- [ ] An imported value replaces the computed one and survives the next weekly run.
