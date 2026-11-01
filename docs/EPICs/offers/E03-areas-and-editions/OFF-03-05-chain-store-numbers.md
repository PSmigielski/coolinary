# OFF-03-05 Link shops to chain store numbers

| | |
|---|---|
| Service | offers |
| Epic | [OFF-03 Postcode areas and flyer editions](README.md) |
| Release | R2 |
| Size | M |
| Depends on | OFF-03-04, OFF-02-03 |
| Unblocks | OFF-06-03, OFF-06-07, OFF-13-03 |

## Goal

Each Kaufland shop (and later each Rewe, Edeka, Netto DE shop) is linked to the store number its flyer uses
(`leaflets.kaufland.com/…/ar/[store]`), so its own edition can be read.

## Business rules

1. A chain with `edition_model = per_store` has a store list: store number, name, address, postcode, position. Its
   adapter provides `ListStores` (OFF-06-03, OFF-06-07; feed for OFF-13-03). Decided (Q150): the adapter reads the
   chain's public store finder, covered by the same retailer-terms gate as the flyers.
2. Store lists are refreshed monthly, after the OSM import.
3. Auto-link: an unlinked active shop of that chain is linked to a store when exactly one store lies within 200 m of
   it, and no other shop lies closer to that store. Otherwise it goes to the admin link queue.
4. One store links to at most one shop, and one shop to at most one store.
5. A link made by admin is never changed by auto-link.
6. When a store disappears from the store list twice in a row, its link is removed and its edition stops being read;
   the shop keeps existing (OSM decides whether it is closed).
7. Linking creates the store edition (`<chain>:store:<number>`) if it does not exist.

## Contract

`admin.proto`

```proto
rpc ListStoreLinkQueue(ListStoreLinkQueueRequest) returns (ListStoreLinkQueueResponse); // unlinked shops with candidate stores and distances
rpc LinkStore(LinkStoreRequest) returns (LinkStoreResponse);       // shop_id, store_ref
rpc UnlinkStore(UnlinkStoreRequest) returns (UnlinkStoreResponse);
```

## Data

`offers/migrations/000NN_chain_stores.sql`

```sql
CREATE TABLE chain_stores (
    chain_id   text NOT NULL REFERENCES chains (id),
    store_ref  text NOT NULL,
    name       text NOT NULL DEFAULT '',
    street     text NOT NULL DEFAULT '',
    postcode   text NOT NULL DEFAULT '',
    latitude   numeric(9,6),
    longitude  numeric(9,6),
    missing_refreshes smallint NOT NULL DEFAULT 0,
    PRIMARY KEY (chain_id, store_ref)
);
CREATE TABLE shop_store_links (
    shop_id   text PRIMARY KEY REFERENCES shops (id),
    chain_id  text NOT NULL,
    store_ref text NOT NULL,
    linked_by text NOT NULL CHECK (linked_by IN ('auto', 'admin')),
    linked_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (chain_id, store_ref),
    FOREIGN KEY (chain_id, store_ref) REFERENCES chain_stores (chain_id, store_ref)
);
```

## Errors

| Case | gRPC status |
|---|---|
| Store already linked to another shop | `ALREADY_EXISTS` |
| Shop's chain is not per-store | `FAILED_PRECONDITION` |

## Acceptance criteria

- [ ] With a fixture store list, a Kaufland shop 80 m from one store is auto-linked; one with two stores within
      200 m lands in the admin queue.
- [ ] An admin link survives the next refresh.
