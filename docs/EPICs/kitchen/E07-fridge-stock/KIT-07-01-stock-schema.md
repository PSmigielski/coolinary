# KIT-07-01 Stock items and movements schema

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-07 Fridge stock](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | KIT-01 (migrations, outbox), KIT-02 (`ingredients` table), KIT-04 (`products` table) |
| Unblocks | KIT-07-02 … KIT-07-09, KIT-12-01 |

## Goal

Store each household's fridge, freezer and pantry as stock items, and record every change as a movement, so later
tasks can add, consume, value and report on stock without guessing what happened.

## Business rules

1. A stock item belongs to exactly one household. Every query filters by `household_id`; an item id of another
   household behaves exactly like an unknown id (`NOT_FOUND`), so ids never leak across households.
2. Quantity is stored in one of the three base units of `coolinary.common.v1.Unit`: grams, millilitres or pieces.
   Callers convert kg and l (× 1,000) before sending. Quantity is > 0 and at most 100,000 g, 100,000 ml or
   1,000 pieces per item.
3. Each item has exactly one storage place: fridge, freezer or pantry.
4. `expires_on` is a calendar date in the household's timezone (KIT-10-01), or empty for items without a date.
   `expiry_source` says where the date came from: `none` (no date), `default` (shelf life), `user`, `pack_code`
   (GS1 2D code read at scan).
5. `ingredient_id` links the item to the catalogue. It is empty only while the item is "not matched" (typed free
   text, KIT-07-09). Unmatched items are listed and editable but never count for recipes or the shopping list.
6. `product_ean` is set when the item came from a barcode, a receipt line linked to a product, or a list line that
   names a product.
7. `name` is a display snapshot in the household's language, 1–120 characters: the product name, else the
   ingredient name, else the user's text.
8. Price is optional: `price_paid_minor` + `currency` for the whole item when it was added; `offer_saving_minor` is
   the saving against the regular price when it was bought on offer. Both ≥ 0, and `currency` is required when
   either is set.
9. `initial_quantity` keeps the quantity at creation (plus merged additions). It is used for the rounding-dust rule
   (KIT-07-07) and to split the price over partial use (KIT-12-01).
10. `version` starts at 1 and increases by 1 on every change to the item. Edits must send the version they saw
    (KIT-07-04).
11. Decided (Q077): a household holds at most 1,000 stock items, as a guard against abuse.
12. Every insert, quantity change, move or removal writes exactly one row to `stock_movements` in the same
    transaction. Movements are append-only and hold no free text and no user ids.
13. `household_stock_state.stock_version` increases by 1 in every transaction that changes the household's stock,
    however many items it touches. Responses return it so clients and planner can cache by version.
14. Movement reasons: `added`, `cooked`, `used`, `thrown_away`, `mistake`, `moved_in`, `moved_out`, `adjusted`.
    `grams` and `value_minor` on a movement are filled by KIT-12-01; until then they stay empty.

## Contract

`contracts/coolinary/kitchen/v1/stock.proto` — shared messages, used by the RPCs in KIT-07-03 … KIT-07-07.

```proto
enum StoragePlace {
  STORAGE_PLACE_UNSPECIFIED = 0;
  STORAGE_PLACE_FRIDGE = 1;
  STORAGE_PLACE_FREEZER = 2;
  STORAGE_PLACE_PANTRY = 3;
}

enum ExpirySource {
  EXPIRY_SOURCE_UNSPECIFIED = 0;
  EXPIRY_SOURCE_NONE = 1;
  EXPIRY_SOURCE_DEFAULT = 2;
  EXPIRY_SOURCE_USER = 3;
  EXPIRY_SOURCE_PACK_CODE = 4;
}

enum StockSource {
  STOCK_SOURCE_UNSPECIFIED = 0;
  STOCK_SOURCE_TYPED = 1;
  STOCK_SOURCE_BARCODE = 2;
  STOCK_SOURCE_UNPACK = 3;
  STOCK_SOURCE_PHOTO = 4;
  STOCK_SOURCE_TEXT = 5;
  STOCK_SOURCE_VOICE = 6;
  STOCK_SOURCE_RECEIPT = 7;
  STOCK_SOURCE_LIST_TICK = 8;
  STOCK_SOURCE_MOVED = 9;
}

// Computed at read time in the household's timezone (KIT-07-05).
enum ExpiryStatus {
  EXPIRY_STATUS_UNSPECIFIED = 0;
  EXPIRY_STATUS_NO_DATE = 1;
  EXPIRY_STATUS_OK = 2;
  EXPIRY_STATUS_SOON = 3;      // 2–3 days ahead
  EXPIRY_STATUS_TOMORROW = 4;
  EXPIRY_STATUS_TODAY = 5;
  EXPIRY_STATUS_EXPIRED = 6;
}

message StockItem {
  string id = 1;
  string ingredient_id = 2;              // empty while not matched
  string product_ean = 3;
  string name = 4;
  coolinary.common.v1.Quantity quantity = 5;
  StoragePlace storage = 6;
  string expires_on = 7;                 // YYYY-MM-DD; empty = no date
  ExpirySource expiry_source = 8;
  ExpiryStatus expiry_status = 9;
  StockSource source = 10;
  coolinary.common.v1.Money price_paid = 11;
  string added_at = 12;                  // RFC 3339
  int32 version = 13;
  bool matched = 14;
}
```

## Data

`kitchen/migrations/000NN_stock.sql`

```sql
-- +goose Up
CREATE TABLE stock_items (
    id                 uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    household_id       uuid NOT NULL,
    ingredient_id      uuid REFERENCES ingredients (id),      -- column type follows KIT-02
    product_ean        text REFERENCES products (ean),
    name               text NOT NULL CHECK (char_length(name) BETWEEN 1 AND 120),
    quantity           numeric(12,3) NOT NULL CHECK (quantity > 0),
    initial_quantity   numeric(12,3) NOT NULL CHECK (initial_quantity > 0),
    unit               smallint NOT NULL CHECK (unit IN (1, 2, 3)),     -- common.v1.Unit
    storage            smallint NOT NULL CHECK (storage IN (1, 2, 3)),  -- StoragePlace
    expires_on         date,
    expiry_source      smallint NOT NULL CHECK (expiry_source BETWEEN 1 AND 4),
    source             smallint NOT NULL CHECK (source BETWEEN 1 AND 9),
    price_paid_minor   bigint CHECK (price_paid_minor >= 0),
    offer_saving_minor bigint CHECK (offer_saving_minor >= 0),
    currency           text CHECK (currency ~ '^[A-Z]{3}$'),
    added_at           timestamptz NOT NULL DEFAULT now(),
    updated_at         timestamptz NOT NULL DEFAULT now(),
    version            integer NOT NULL DEFAULT 1,
    CHECK (quantity <= CASE unit WHEN 3 THEN 1000 ELSE 100000 END),
    CHECK ((price_paid_minor IS NULL AND offer_saving_minor IS NULL) OR currency IS NOT NULL),
    CHECK ((expires_on IS NULL) = (expiry_source = 1))
);
CREATE INDEX stock_items_household_expiry ON stock_items (household_id, expires_on NULLS LAST, added_at);
CREATE INDEX stock_items_household_ingredient ON stock_items (household_id, ingredient_id);
CREATE INDEX stock_items_unmatched ON stock_items (household_id) WHERE ingredient_id IS NULL;

CREATE TABLE stock_movements (
    id            bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    household_id  uuid NOT NULL,
    stock_item_id uuid,                     -- kept after the item is deleted, no FK
    ingredient_id uuid,
    delta         numeric(12,3) NOT NULL,   -- + added / moved in, − taken
    unit          smallint NOT NULL,
    reason        smallint NOT NULL CHECK (reason BETWEEN 1 AND 8),
    expires_on    date,                     -- the item's date at the time of the movement
    grams         numeric(12,3),            -- KIT-12-01
    value_minor   bigint,                   -- KIT-12-01
    currency      text,
    occurred_at   timestamptz NOT NULL DEFAULT now(),
    occurred_on   date NOT NULL,            -- local date in the household's timezone
    ref           text                      -- cook event id, receipt ref, move ref
);
CREATE INDEX stock_movements_household_time ON stock_movements (household_id, occurred_at);

CREATE TABLE household_stock_state (
    household_id  uuid PRIMARY KEY,
    stock_version bigint NOT NULL DEFAULT 0,
    item_count    integer NOT NULL DEFAULT 0 CHECK (item_count BETWEEN 0 AND 1000)
);

-- +goose Down
DROP TABLE household_stock_state;
DROP TABLE stock_movements;
DROP TABLE stock_items;
```

`item_count` is kept in step by the insert and delete queries, so the 1,000 limit is checked without counting rows.

## Errors

Schema only; RPC errors are defined in KIT-07-03 … KIT-07-07.

## Acceptance criteria

- [ ] The migration applies on the skeleton database and on an empty one; `goose down` removes the three tables.
- [ ] sqlc queries exist: insert item, get item by id and household, list by household, update with version check
      (`WHERE id = $1 AND household_id = $2 AND version = $3`), delete, insert movement, bump stock version
      (returns the new value), adjust item count.
- [ ] Constraint tests: quantity 0, 1,001 pieces, unit 4, storage 0, a price without currency, a date with
      `expiry_source = none` are all rejected.
- [ ] Inserting the 1,001st item of a household fails on the `item_count` check; another household is unaffected.
- [ ] A query run with household A's id never returns household B's items (test with two households).
