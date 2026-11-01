# OFF-02-06 User proposes a shop correction

| | |
|---|---|
| Service | offers |
| Epic | [OFF-02 Chains and shops](README.md) |
| Release | R2 |
| Size | S |
| Depends on | OFF-02-02 |
| Unblocks | OFF-02-07 |

## Goal

A user can report that a shop is in the wrong place, belongs to another chain, or has closed, so the map and the
"nearest first" order stay right.

## Business rules

1. Any account may propose for an active shop: a new position, a different chain, or "closed".
2. A proposal does not change the shop until an admin approves it (OFF-02-07). The proposing household does not see
   the change early (keeps one truth for distances).
3. One open proposal per household per shop; a new one replaces the old.
4. A moved position more than 2 km from the current one is refused (likely a different shop; add it instead).
5. When 3 different households propose "closed" for the same shop within 30 days, the shop is hidden from shop lists
   (status stays `active` but `reported_closed = true`) until an admin decides. It stays in existing my shops.
6. For OSM shops, an approved correction is stored as an override (`shop_overrides`) so the next OSM import does not
   undo it. Decided (Q146): overrides are never sent to OpenStreetMap automatically; the founder edits OSM by hand
   when worthwhile.

## Contract

```proto
rpc ProposeShopCorrection(ProposeShopCorrectionRequest) returns (ProposeShopCorrectionResponse);

message ProposeShopCorrectionRequest {
  string household_id = 1;
  string shop_id = 2;
  oneof change {
    LatLng new_position = 3;
    string new_chain_id = 4;
    bool closed = 5;
  }
}
message LatLng { double latitude = 1; double longitude = 2; }
```

## Data

`offers/migrations/000NN_shop_corrections.sql`

```sql
CREATE TABLE shop_corrections (
    id           uuid PRIMARY KEY,
    shop_id      text NOT NULL REFERENCES shops (id),
    household_id text NOT NULL,
    kind         text NOT NULL CHECK (kind IN ('position', 'chain', 'closed')),
    new_latitude  numeric(9,6),
    new_longitude numeric(9,6),
    new_chain_id text REFERENCES chains (id),
    status       text NOT NULL DEFAULT 'open' CHECK (status IN ('open', 'approved', 'rejected', 'replaced')),
    created_at   timestamptz NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX shop_corrections_open ON shop_corrections (shop_id, household_id) WHERE status = 'open';
CREATE TABLE shop_overrides (
    shop_id   text PRIMARY KEY REFERENCES shops (id),
    latitude  numeric(9,6),
    longitude numeric(9,6),
    chain_id  text REFERENCES chains (id),
    closed    boolean
);
ALTER TABLE shops ADD COLUMN reported_closed boolean NOT NULL DEFAULT false;
```

## Errors

| Case | gRPC status |
|---|---|
| Shop not active | `FAILED_PRECONDITION` |
| Move > 2 km | `INVALID_ARGUMENT` "add a new shop instead" |

## Acceptance criteria

- [ ] A position correction approved by admin moves the shop; a following OSM import keeps the corrected position.
- [ ] Three "closed" reports from three households hide the shop from `ListShopsNearArea`; two do not.
