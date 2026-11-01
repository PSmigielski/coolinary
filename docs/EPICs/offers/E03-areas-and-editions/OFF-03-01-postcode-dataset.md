# OFF-03-01 Postcode and town dataset import

| | |
|---|---|
| Service | offers |
| Epic | [OFF-03 Postcode areas and flyer editions](README.md) |
| Release | R2 |
| Size | M |
| Depends on | OFF-01-05 |
| Unblocks | OFF-03-02, OFF-02-05 (postcode fill) |

## Goal

offers knows the centre point of every postcode and town in Poland and Germany, so a typed postcode or town gives
the same shop list as the phone's location (features.md §3 "Location", Verification 10).

## Business rules

1. Dataset (Decided, Q147): GeoNames postal codes, one file per country, CC BY 4.0, with the visible credit
   "GeoNames": country, postcode, place name, region names, latitude, longitude.
2. Import per country on a yearly schedule (1 February, 04:00 UTC) and on `TriggerJob postcode_import:<cc>`.
3. One postcode may have several rows (several places). The postcode's centre is the mean latitude and mean
   longitude of its rows; its places are kept as a list.
4. Postcodes are stored normalised: PL `NN-NNN`, DE `NNNNN`. Rows that don't match the format are skipped and counted.
5. Towns are stored once per (country, place name, region) with the mean centre of their postcodes and the list of
   their postcodes. A search key is the lower-cased name without diacritics (`łódź` → `lodz`).
6. Import replaces the country's rows in one transaction; if the new file has less than 90% of the previous row count,
   nothing is replaced and an alert is sent (`offers.postcodes.<cc>`).
7. Shops without a postcode get the postcode whose centre is nearest (OFF-02-02 rule 5); this runs after each
   postcode import and after each OSM import.
8. The credit "Postcode data © GeoNames (CC BY 4.0)" is returned with area responses (OFF-03-02).

## Contract

No public RPC; admin `TriggerJob`.

## Data

`offers/migrations/000NN_postcodes.sql`

```sql
CREATE EXTENSION IF NOT EXISTS unaccent;
CREATE TABLE postcode_areas (
    country_code text NOT NULL,
    postcode     text NOT NULL,
    places       text[] NOT NULL,
    region       text NOT NULL DEFAULT '',
    latitude     numeric(9,6) NOT NULL,
    longitude    numeric(9,6) NOT NULL,
    PRIMARY KEY (country_code, postcode)
);
CREATE TABLE towns (
    id           bigserial PRIMARY KEY,
    country_code text NOT NULL,
    name         text NOT NULL,
    search_key   text NOT NULL,
    region       text NOT NULL DEFAULT '',
    postcodes    text[] NOT NULL,
    latitude     numeric(9,6) NOT NULL,
    longitude    numeric(9,6) NOT NULL,
    UNIQUE (country_code, name, region)
);
CREATE INDEX towns_search ON towns USING gin (search_key gin_trgm_ops);
```

(`pg_trgm` is enabled in this database by the migration if not present.)

## Errors

| Case | Behaviour |
|---|---|
| Download fails | Run `failed`, old data kept; alert after 3 failed runs (OFF-01-05 rule 8) |
| New file under 90% of the old row count | Nothing replaced, alert |

## Acceptance criteria

- [ ] After import, `00-001` (Warsaw) and `80331` (Munich) exist with centres inside those cities.
- [ ] A file with 50% of the rows is refused with an alert and the old data stays.
- [ ] Town search key for "Łódź" is `lodz`.
