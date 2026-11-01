# OFF-13-01 Feed client and retailer/store mapping

| | |
|---|---|
| Service | offers |
| Epic | [OFF-13 Licensed offers feed](README.md) |
| Release | R3 |
| Size | L |
| Status | todo |
| Depends on | OFF-02-01, OFF-03-05 |
| Unblocks | OFF-13-02 |

## Goal

offers can read the network's retailers, stores and offers, and knows which Coolinary chain and shop each of them is.

## Business rules

1. `GATE: signed feed agreement` (epic README). Credentials come from Docker secrets (`OFFERS_FEED_API_KEY_FILE`).
2. The client reads the provider's retailers, stores (with addresses and positions) and offers through its documented
   API, within its published rate limits. The provider's own field names stay inside the client; the rest of offers
   sees Coolinary types.
3. Retailer mapping: each feed retailer is mapped by admin to one Coolinary chain (or ignored). Unmapped retailers are
   listed for review; nothing of theirs is ingested.
4. Store mapping: each feed store of a mapped retailer is linked to a Coolinary shop with the same rule as chain store
   numbers (OFF-03-05 rule 3: exactly one shop within 200 m), else queued for admin. Feed stores with no shop can
   create a shop with `source = 'feed'` if the agreement allows using their positions. Decided (Q228): shops
   are created from feed stores only where OSM has none, flagged `source = 'feed'`.
5. Mappings refresh daily at 05:00 `Europe/Berlin`.
6. A chain fed by the network gets coverage `feed` only through the gate of OFF-08-02 (a 100-offer sample also applies
   to feed offers).

## Contract

`admin.proto`

```proto
rpc ListFeedRetailers(ListFeedRetailersRequest) returns (ListFeedRetailersResponse);
rpc MapFeedRetailer(MapFeedRetailerRequest) returns (MapFeedRetailerResponse);   // feed retailer -> chain_id or ignore
rpc ListFeedStoreQueue(ListFeedStoreQueueRequest) returns (ListFeedStoreQueueResponse);
rpc LinkFeedStore(LinkFeedStoreRequest) returns (LinkFeedStoreResponse);
```

## Data

`offers/migrations/000NN_feed.sql`

```sql
CREATE TABLE feed_retailers (
    feed_retailer_id text PRIMARY KEY,
    name        text NOT NULL,
    country_code text NOT NULL,
    chain_id    text REFERENCES chains (id),
    ignored     boolean NOT NULL DEFAULT false
);
CREATE TABLE feed_stores (
    feed_store_id    text PRIMARY KEY,
    feed_retailer_id text NOT NULL REFERENCES feed_retailers (feed_retailer_id),
    name      text NOT NULL DEFAULT '',
    street    text NOT NULL DEFAULT '',
    postcode  text NOT NULL DEFAULT '',
    latitude  numeric(9,6),
    longitude numeric(9,6),
    shop_id   text REFERENCES shops (id),
    linked_by text CHECK (linked_by IN ('auto', 'admin'))
);
ALTER TABLE shops DROP CONSTRAINT shops_source_check,
    ADD CONSTRAINT shops_source_check CHECK (source IN ('osm', 'user', 'feed'));
```

## Errors

| Case | Behaviour |
|---|---|
| Feed API down or 401 | Refresh `failed`, alert `offers.feed`; existing mappings kept |
| Rate limit hit | Wait as the provider's response says; resume |

## Acceptance criteria

- [ ] With the provider's sample data, retailers appear unmapped until admin maps them.
- [ ] A feed store 50 m from an OSM Rewe shop is auto-linked; one with no shop nearby is queued.
