# GW-10-03 Offer browse, search and flyer status

| | |
|---|---|
| Service | gateway |
| Epic | [GW-10 Shops, offers, flyers and price matrix routes](README.md) |
| Release | R2 |
| Size | M |
| Status | todo |
| Depends on | GW-10-02, OFF-09 (`BrowseOffers`, `GetOffer`, `ListOfferCategories`), OFF-05 (`GetFlyerStatus`), PLN-17 (`AddOfferToList`) |
| Unblocks | MOB-03 |

## Goal

Browse and search this week's (and next week's) offers by shop, category and end date; tap an offer to add it to
the list; every offer links to its original flyer; chains whose flyer could not be read say so (features.md §2
"Flyers and offers", §3 "Keeping flyers in sync").

## Business rules

1. `GET /v1/offers` query: `shopIds` (repeatable, nearest first; default the household's saved shops), `view`
   (`current` | `upcoming` | `all`, default `current`), `category` (one), `q` (2–60 chars), `endingBy` (date),
   `pageSize`, `pageToken`. Calls offers `BrowseOffers(household_id, shop_ids, view, category, ending_by, query,
   page)`, which replaces the skeleton `ListOffers` (OFF-09-01).
2. Free: offers returns only offers of the household's 2 saved shops, whatever `shopIds` says (OFF-04/09). The
   gateway does not filter.
3. Ordering by distance is the phone's job: the phone passes `shopIds` nearest first and offers keeps that order.
4. Offer fields as offers returns them: id, shop and chain, product name, size, price, regular price if printed,
   `conditions` (e.g. "2+1", "with app or card"), valid dates, upcoming flag, `sourceUrl` (the original flyer page),
   matched ingredient (or null), category. The response also carries each shop's flyer status and the attribution.
5. `GET /v1/offers/{offerId}` → `GetOffer`: the offer with its chain and flyer.
6. `GET /v1/offers/categories?shopIds=&view=` → `ListOfferCategories`: categories with counts, for the filter.
7. **Adding an offer to the list:** `POST /v1/list/offer-lines` body `{ "offerId" }` → planner
   `AddOfferToList(household_id, member_id, offer_id, client_mutation_id)` (PLN-17-02); the line carries the offer
   badge. Idempotency required.
8. `GET /v1/flyers/status?shopIds=` → offers `GetFlyerStatus(household_id, shop_ids)` (empty = the household's
   shops): per shop `ok` | `not_covered` | `store_not_linked` | `not_updated` with `lastUpdatedAt`, so the app can show
   "flyer not updated, scan it yourself".
9. Expired offers are never returned (offers removes them on their end date).

## Contract

| Method | Path | 2xx |
|---|---|---|
| GET | `/v1/offers?shopIds=&view=&category=&q=&endingBy=` | 200 `{ "offers": [offer], "shopStatuses": [...], "nextPageToken", "attribution" }` |
| GET | `/v1/offers/{offerId}` | 200 `{ "offer", "chain", "flyer" }` |
| GET | `/v1/offers/categories?shopIds=&view=` | 200 `{ "categories": [ { "category": "meat", "count": 14 } ] }` |
| POST | `/v1/list/offer-lines` | 201 `{ "version", "line": line }` |
| GET | `/v1/flyers/status?shopIds=` | 200 `{ "statuses": [ { "shopId", "code": "ok", "lastUpdatedAt": "2026-11-04T06:00:00Z" } ] }` |

```json
offer { "id": "off_...", "chainId": "kaufland-pl", "shopId": "osm-w99", "productName": "Udka z kurczaka",
        "size": { "value": 1000, "unit": "g" }, "price": { "amountMinor": 1599, "currency": "PLN" },
        "regularPrice": { "amountMinor": 2199, "currency": "PLN" }, "conditions": "z kartą Kaufland Card",
        "validFrom": "2026-11-05", "validTo": "2026-11-11", "upcoming": false,
        "sourceUrl": "https://leaflets.kaufland.com/pl-PL/...", "ingredientId": "ing_chicken_thigh",
        "category": "meat" }
```

Tags `offers`, `flyers`, `auth: true`.

## Data

No storage.

## Errors

| Case | HTTP | `code` |
|---|---|---|
| Bad `view`, `q` too short or long | 400 | `validation` |
| No saved shops and no `shopIds` | 422 | `no_shops` |

## Acceptance criteria

- [ ] Every offer carries a `sourceUrl` (features.md §2 "Every offer links to the original flyer").
- [ ] `view=upcoming` returns flyers published ahead, marked `upcoming: true`.
- [ ] A Free household asking for 5 shop ids receives offers of its 2 saved shops only.
- [ ] After an offer's end date it is gone from the response (verification 6, with OFF-07).
- [ ] A shop whose chain sync failed shows a not-updated status (verification 7).
- [ ] Adding an offer to the list twice with the same idempotency key creates one line with the offer badge.
