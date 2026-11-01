# OFF-07-01 Offers schema rework

| | |
|---|---|
| Service | offers |
| Epic | [OFF-07 Offer extraction, validation and publishing](README.md) |
| Release | R2 |
| Size | M |
| Depends on | OFF-05-01, OFF-03-04 |
| Unblocks | OFF-07-02 … 08, OFF-09-*, OFF-10-*, OFF-11-01, OFF-13-02 |

## Goal

An offer row holds everything features.md asks to extract (product, size, price, conditions, valid dates) plus what
the app and planner need: the edition it belongs to, the matched ingredient, the unit price and where it came from.

## Business rules

1. An offer of a chain flyer belongs to an edition, not to a shop: a national Lidl offer is one row valid for every
   Lidl PL shop. A user-supplied flyer's offer belongs to one shop and one household (OFF-10).
2. Source kinds: `chain_flyer`, `user_flyer`, `feed`. Exactly one of `flyer_version_id`, `user_flyer_id`,
   `feed_item_id` is set accordingly.
3. Money is stored in minor units with the chain's currency. `price_minor > 0` is enforced by the database only for
   states `ready` and `published` (a held offer may have a missing price).
4. Pack: `pack_count` (units in a multipack, default 1), `pack_value` + `pack_unit` (`g`, `ml`, `piece`) per unit,
   or `loose = true` with a price per `kg` or per `piece` (weighed counters).
5. Unit price (`unit_price_minor` per `unit_basis` `kg`, `l` or `piece`) is computed by OFF-07-03, never taken from AI.
6. Conditions are stored as text (as printed, for display) and structured: multi-buy, loyalty program, minimum
   quantity, maximum per customer.
7. Dates are local dates, inclusive.
8. The existing `offers` table and its two example rows are migrated: `shop_id` becomes nullable, example rows get
   source kind `chain_flyer`, edition `example-pl:national`, state `published`.
9. The existing `Offer` proto message gains the new fields; `shop_id` stays for user flyers and is deprecated for
   chain flyers (the shops of an offer come from its edition, OFF-09-01).

## Contract

`contracts/coolinary/offers/v1/offers.proto`

```proto
message Offer {
  string id = 1;
  string shop_id = 2 [deprecated = true];   // user flyers only; use shop_ids
  string product_name = 3;                  // as printed
  coolinary.common.v1.Money price = 4;
  string conditions = 5;                    // as printed: "2+1", "z aplikacją Lidl Plus"
  string valid_from = 6;
  string valid_to = 7;
  string source_url = 8;
  // new
  string chain_id = 9;
  string edition_id = 10;
  repeated string shop_ids = 11;            // the requested shops this offer applies to
  string brand = 12;
  OfferPack pack = 13;
  coolinary.common.v1.Money regular_price = 14;  // when the flyer shows it
  int32 discount_percent = 15;              // computed from regular price, 0 if unknown
  UnitPrice unit_price = 16;
  OfferConditions structured_conditions = 17;
  string ingredient_id = 18;
  string category = 19;
  int32 page_no = 20;
  OfferSourceKind source_kind = 21;
  bool upcoming = 22;                       // starts after today ("next week")
}
message OfferPack {
  int32 count = 1;
  coolinary.common.v1.Quantity per_unit = 2;
  bool loose = 3;
}
message UnitPrice {
  coolinary.common.v1.Money price = 1;
  string basis = 2;                         // "kg", "l", "piece"
}
message OfferConditions {
  int32 multi_buy_buy = 1;                  // "2+1" -> buy 3
  int32 multi_buy_pay = 2;                  //          pay 2
  coolinary.common.v1.Money multi_buy_total = 3;  // "2 for 5 zł" -> buy 2, total 5.00
  bool loyalty_required = 4;
  string loyalty_program = 5;               // "Lidl Plus", "Moja Biedronka"
  int32 min_quantity = 6;
  int32 max_per_customer = 7;
}
enum OfferSourceKind { OFFER_SOURCE_KIND_UNSPECIFIED = 0; OFFER_SOURCE_KIND_CHAIN_FLYER = 1; OFFER_SOURCE_KIND_USER_FLYER = 2; OFFER_SOURCE_KIND_FEED = 3; }
```

## Data

`offers/migrations/000NN_offers_rework.sql`

```sql
ALTER TABLE offers DROP CONSTRAINT IF EXISTS offers_price_minor_check;
ALTER TABLE offers ALTER COLUMN shop_id DROP NOT NULL, ALTER COLUMN price_minor DROP NOT NULL;
ALTER TABLE offers
    ADD COLUMN source_kind      text NOT NULL DEFAULT 'chain_flyer' CHECK (source_kind IN ('chain_flyer', 'user_flyer', 'feed')),
    ADD COLUMN chain_id         text REFERENCES chains (id),
    ADD COLUMN edition_id       text REFERENCES editions (id),
    ADD COLUMN flyer_version_id uuid REFERENCES flyer_versions (id),
    ADD COLUMN user_flyer_id    uuid,
    ADD COLUMN feed_item_id     text,
    ADD COLUMN household_id     text,               -- user flyers only
    ADD COLUMN language         text NOT NULL DEFAULT 'pl',
    ADD COLUMN brand            text NOT NULL DEFAULT '',
    ADD COLUMN pack_count       int NOT NULL DEFAULT 1 CHECK (pack_count >= 1),
    ADD COLUMN pack_value       double precision,
    ADD COLUMN pack_unit        text CHECK (pack_unit IN ('g', 'ml', 'piece')),
    ADD COLUMN loose            boolean NOT NULL DEFAULT false,
    ADD COLUMN regular_price_minor bigint,
    ADD COLUMN unit_price_minor bigint,
    ADD COLUMN unit_basis       text CHECK (unit_basis IN ('kg', 'l', 'piece')),
    ADD COLUMN cond_buy         int, ADD COLUMN cond_pay int, ADD COLUMN cond_total_minor bigint,
    ADD COLUMN loyalty_required boolean NOT NULL DEFAULT false,
    ADD COLUMN loyalty_program  text NOT NULL DEFAULT '',
    ADD COLUMN min_quantity     int, ADD COLUMN max_per_customer int,
    ADD COLUMN page_no          int,
    ADD COLUMN ingredient_id    text,
    ADD COLUMN category         text NOT NULL DEFAULT '',
    ADD COLUMN match_confidence real,
    ADD COLUMN ai_confidence    real,
    ADD COLUMN kind             text NOT NULL DEFAULT 'food' CHECK (kind IN ('food', 'drink', 'non_food')),
    ADD COLUMN state            text NOT NULL DEFAULT 'extracted'
               CHECK (state IN ('extracted', 'ready', 'held_back', 'published', 'expired', 'rejected', 'superseded')),
    ADD COLUMN held_reasons     text[] NOT NULL DEFAULT '{}',
    ADD COLUMN published_at     timestamptz,
    ADD COLUMN created_at       timestamptz NOT NULL DEFAULT now(),
    ADD CONSTRAINT offers_price_when_live CHECK (state NOT IN ('ready', 'published') OR price_minor > 0);
CREATE INDEX offers_visible ON offers (edition_id, valid_to) WHERE state = 'published';
CREATE INDEX offers_ingredient ON offers (ingredient_id, valid_to) WHERE state = 'published';
CREATE INDEX offers_household ON offers (household_id) WHERE household_id IS NOT NULL;
CREATE INDEX offers_name_trgm ON offers USING gin (product_name gin_trgm_ops);
```

## Errors

| Case | Behaviour |
|---|---|
| `ready` or `published` offer with price ≤ 0 | Refused by the database (a bug in the caller, `INTERNAL`) |

## Acceptance criteria

- [ ] Migration keeps the two example offers listable through the existing `ListOffers`.
- [ ] A `ready` offer with `price_minor = 0` is refused by the database.
- [ ] `buf breaking` passes: only fields added and one field marked deprecated.
