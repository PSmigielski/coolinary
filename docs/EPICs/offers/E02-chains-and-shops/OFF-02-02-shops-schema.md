# OFF-02-02 Shops schema rework

| | |
|---|---|
| Service | offers |
| Epic | [OFF-02 Chains and shops](README.md) |
| Release | R2 |
| Size | S |
| Depends on | OFF-02-01 |
| Unblocks | OFF-02-03, OFF-02-05, OFF-03-03, OFF-04-01 |

## Goal

Shops carry everything the phone needs to compute distance and show a shop (chain, name, address, position) and
everything offers needs to keep them in sync with OSM and user corrections.

## Business rules

1. A shop belongs to one chain. Its id is stable for ever (my shops, editions and prices point at it): `osm-n123`,
   `osm-w456` for OSM nodes and ways, `usr-<uuidv7>` for user-added shops.
2. Position is stored as latitude and longitude with 6 decimals (about 0.1 m), WGS 84. Positions outside the chain's
   country bounding box are rejected.
3. Status: `active`, `closed` (gone from OSM or reported closed and approved), `pending` (user-added, waiting for
   review; visible only to the adding household), `rejected` (hidden), `merged` (duplicate; `merged_into` points to
   the kept shop).
4. Source: `osm` or `user`. For OSM shops `osm_type` (`node`, `way`) and `osm_id` are unique together.
5. `postcode` and `city` come from OSM `addr:*` tags when present; when missing, OFF-03-01 fills `postcode` from the
   nearest postcode centroid.
6. The existing `shops` table and its example row are migrated: `chain` text becomes `chain_id`, the example shop
   gets coordinates in Warsaw (52.229700, 21.012200), postcode `00-001`.
7. A shop's display name is the OSM `name`, else the chain name plus street (`Lidl, ul. Prosta 2`).

## Contract

`contracts/coolinary/offers/v1/shops.proto`

```proto
message Shop {
  string id = 1;
  string chain_id = 2;
  string chain_name = 3;
  string name = 4;
  string street = 5;           // street and house number
  string postcode = 6;
  string city = 7;
  double latitude = 8;
  double longitude = 9;
  ShopStatus status = 10;
  bool household_only = 11;    // true for the household's own pending shop
}

enum ShopStatus {
  SHOP_STATUS_UNSPECIFIED = 0;
  SHOP_STATUS_ACTIVE = 1;
  SHOP_STATUS_CLOSED = 2;
  SHOP_STATUS_PENDING = 3;
}
```

## Data

`offers/migrations/000NN_shops_rework.sql`

```sql
ALTER TABLE shops ADD COLUMN chain_id text REFERENCES chains (id);
UPDATE shops SET chain_id = 'example-pl';           -- dev example only
ALTER TABLE shops ALTER COLUMN chain_id SET NOT NULL, DROP COLUMN chain;
ALTER TABLE shops
    ADD COLUMN street      text NOT NULL DEFAULT '',
    ADD COLUMN postcode    text NOT NULL DEFAULT '',
    ADD COLUMN city        text NOT NULL DEFAULT '',
    ADD COLUMN latitude    numeric(9,6),
    ADD COLUMN longitude   numeric(9,6),
    ADD COLUMN status      text NOT NULL DEFAULT 'active'
               CHECK (status IN ('active', 'closed', 'pending', 'rejected', 'merged')),
    ADD COLUMN merged_into text REFERENCES shops (id),
    ADD COLUMN source      text NOT NULL DEFAULT 'osm' CHECK (source IN ('osm', 'user')),
    ADD COLUMN osm_type    text CHECK (osm_type IN ('node', 'way')),
    ADD COLUMN osm_id      bigint,
    ADD COLUMN missing_refreshes smallint NOT NULL DEFAULT 0,
    ADD COLUMN added_by_household text,                 -- user-added shops only
    ADD COLUMN updated_at  timestamptz NOT NULL DEFAULT now();
CREATE UNIQUE INDEX shops_osm ON shops (osm_type, osm_id) WHERE source = 'osm';
CREATE INDEX shops_position ON shops (country_code, latitude, longitude) WHERE status = 'active';
CREATE INDEX shops_postcode ON shops (country_code, postcode);
```

## Errors

| Case | Behaviour |
|---|---|
| Position outside the chain's country | `INVALID_ARGUMENT` in every RPC or job that writes shops |
| Duplicate `(osm_type, osm_id)` | Unique index; the import upserts instead |

## Acceptance criteria

- [ ] Migration runs on the skeleton database; the example offers still list.
- [ ] Inserting a shop at latitude 48.1 / longitude 11.5 (Munich) for a PL chain fails.
- [ ] Two OSM shops with the same `(osm_type, osm_id)` cannot exist.
