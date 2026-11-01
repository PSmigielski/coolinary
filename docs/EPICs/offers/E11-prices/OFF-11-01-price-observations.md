# OFF-11-01 Price observations from flyers

| | |
|---|---|
| Service | offers |
| Epic | [OFF-11 Prices](README.md) |
| Release | R2 |
| Size | M |
| Status | todo |
| Depends on | OFF-07-07 |
| Unblocks | OFF-11-02, OFF-11-03, OFF-11-04 |

## Goal

offers keeps a history of the prices it has seen per chain and ingredient, so reference prices can be computed and the
trip plan has prices beyond today's flyer.

## Business rules

1. When an offer of a chain flyer or the feed is published with an ingredient id, one observation is written in the
   same transaction: country, chain, edition, ingredient, product name, pack, price, regular price, unit price and
   basis, loyalty flag, valid from/to, source `flyer` (or `feed`), and the offer id.
2. User-flyer offers are not recorded (they stay private, OFF-10-04 rule 5).
3. A superseded offer's observation is replaced by the new version's (same flyer, newer reading).
4. Observations are kept 400 days, then deleted (enough for a year of reference prices).
5. When a chain is removed (OFF-08-05), its flyer and feed observations are deleted within 24 h, like its pages.
6. Observations hold no household or user data.

## Contract

No RPC.

## Data

`offers/migrations/000NN_price_observations.sql`

```sql
CREATE TABLE price_observations (
    id            bigserial PRIMARY KEY,
    country_code  text NOT NULL,
    chain_id      text NOT NULL REFERENCES chains (id),
    edition_id    text REFERENCES editions (id),
    ingredient_id text NOT NULL,
    ean           text,
    product_name  text NOT NULL DEFAULT '',
    pack_count    int NOT NULL DEFAULT 1,
    pack_value    double precision,
    pack_unit     text,
    loose         boolean NOT NULL DEFAULT false,
    price_minor   bigint NOT NULL CHECK (price_minor > 0),
    regular_price_minor bigint,
    unit_price_minor bigint NOT NULL,
    unit_basis    text NOT NULL CHECK (unit_basis IN ('kg', 'l', 'piece')),
    currency      text NOT NULL,
    loyalty_required boolean NOT NULL DEFAULT false,
    source        text NOT NULL CHECK (source IN ('flyer', 'feed', 'receipt')),
    offer_id      text,
    observed_on   date NOT NULL,
    valid_from    date,
    valid_to      date,
    outlier       boolean NOT NULL DEFAULT false
);
CREATE INDEX price_obs_lookup ON price_observations (chain_id, ingredient_id, source, observed_on DESC);
CREATE INDEX price_obs_country ON price_observations (country_code, ingredient_id, observed_on DESC);
```

## Errors

| Case | Behaviour |
|---|---|
| Offer without a unit price | No observation (it cannot be compared) |

## Acceptance criteria

- [ ] Publishing a version with 120 matched offers adds 120 observations; a re-read replaces them, leaving 120.
- [ ] A household's user-flyer offers add no observation.
