# ADM-06-09 Licensed feed: retailer and store mapping

| | |
|---|---|
| Service | admin |
| Epic | [ADM-06 Flyer sources, offers and shops review](README.md) |
| Release | R3 |
| Size | S |
| Depends on | ADM-06-07; OFF-13-01 (`ListFeedRetailers`, `MapFeedRetailer`, `ListFeedStoreQueue`, `LinkFeedStore`) |
| Unblocks | — |

## Goal

Offers from the licensed offers network (ShopFully) reach users only for retailers and stores the founder has mapped
to Coolinary chains and shops.

## Business rules

1. `GATE: signed feed agreement` (OFF-13-01 rule 1). The page is hidden until offers reports the feed as configured.
2. Retailers: `ListFeedRetailers` shows each feed retailer with its mapping. Action: map to one Coolinary chain or
   ignore → `MapFeedRetailer`. Nothing of an unmapped retailer is ingested (OFF-13-01 rule 3).
3. Stores: `ListFeedStoreQueue` lists feed stores of mapped retailers that offers couldn't link by itself (no single
   shop within 200 m). Action: link to a shop → `LinkFeedStore`, the same way as ADM-06-07. Offers creates a shop
   from a feed store only where OpenStreetMap has none, marked as coming from the feed (Q228, OFF-13-01 rule 4).
4. A chain fed by the network goes live only through the same 100-offer sample and `go_live` as a flyer chain
   (ADM-06-02, ADM-06-03; OFF-13-01 rule 6).
5. Audited (`offers.feed.map_retailer`, `offers.feed.link_store`).
6. Role `admin` only.

## Contract

| Method | Path | Body / response |
|---|---|---|
| GET | `/admin/v1/offers/feed/retailers` | `{ items: [{ feed_retailer_id, name, chain_id?, ignored }] }` |
| PUT | `/admin/v1/offers/feed/retailers/{feedRetailerId}` | `{ chain_id? , ignore? }` → retailer |
| GET | `/admin/v1/offers/feed/store-queue?page_token=` | `{ items, next_page_token }` |
| POST | `/admin/v1/offers/feed/store-links` | `{ feed_store_id, shop_id }` → link |

## Data

None in admin.

## Errors

| Case | HTTP |
|---|---|
| Feed not configured | 404 `feed_not_configured` |
| Shop already linked to another feed store | 409 |

## Acceptance criteria

- [ ] Mapping a feed retailer to "Rewe" lets its offers be ingested on the next daily refresh; an ignored retailer
      stays out.
- [ ] Linking a queued feed store removes it from the queue.
