# GW-10-01 Shops around a postcode; suggest a shop change

| | |
|---|---|
| Service | gateway |
| Epic | [GW-10 Shops, offers, flyers and price matrix routes](README.md) |
| Release | R2 |
| Size | S |
| Status | todo |
| Depends on | GW-01-02, OFF-02 (`SubmitShop`, `ProposeShopCorrection`), OFF-03 (`ResolveArea`, `ListShopsNearArea`) |
| Unblocks | GW-10-02, MOB-03 |

## Goal

The phone gets every shop the app has data for around the user's postcode, with positions, so it can compute
distances itself, sort nearest first and filter by travel range (features.md §1 "My shops and location", §3
"Location"). Users can add or correct a shop.

## Business rules

1. `GET /v1/areas?country=PL&q=00-950` → offers `ResolveArea(country_code, query)`: matching postcode areas or towns
   (label, area id), with the GeoNames credit. `country` ISO alpha-2, launch values `PL`, `DE` (config); other
   countries → 400 `country_not_covered`. Used when the user types a postcode or town instead of sharing location.
2. `GET /v1/shops?areaId=` or `?country=&postcode=` → offers `ListShopsNearArea(household_id, area_id | country_code
   + postcode)`: every shop within the largest travel range (10 km) plus a margin around the area (OFF-03), each with
   its flyer edition ids, the chains, `truncated`, and the OpenStreetMap attribution. The phone computes distances
   and filters by range (D13).
3. The postcode is the only location the server receives; it is not stored by the gateway or logged in full (log
   the first 2 characters).
4. Public shop data: cached in `redis-cache` 24 h per area; `Cache-Control: private, max-age=3600`.
5. **A missing shop:** `POST /v1/shops` body `{ "chainId" | "otherChainName", "latitude", "longitude", "street"?,
   "postcode", "country" }` → offers `SubmitShop` (moderated there, OFF-02-05). The position is the shop's, chosen by
   the user on a map, not the user's own.
6. **A wrong shop:** `POST /v1/shops/{shopId}/corrections` body one of `{ "newPosition": { "latitude", "longitude" } }`,
   `{ "newChainId" }`, `{ "closed": true }` → offers `ProposeShopCorrection` (OFF-02-06).
7. Rate limit for 5 and 6: 10 per user per day (bucket `shop_suggestion`).

## Contract

| Method | Path | 2xx |
|---|---|---|
| GET | `/v1/areas?country=&q=` | 200 `{ "areas": [ { "id", "kind", "label" } ], "attribution" }` |
| GET | `/v1/shops?areaId=` (or `country` + `postcode`) | 200 `{ "shops": [shop], "chains": [...], "truncated": false, "attributions": [...] }` |
| POST | `/v1/shops` | 202 `{ "shopId", "status": "pending" }` |
| POST | `/v1/shops/{shopId}/corrections` | 202 `{ "status": "pending" }` |

```json
shop { "id": "osm-n2384", "chainId": "lidl-pl", "name": "Lidl Marszałkowska", "address": "ul. Marszałkowska 104",
       "latitude": 52.2319, "longitude": 21.0067, "editionIds": ["lidl-pl-national"] }
```

Tags `shops`, `auth: true`.

## Data (redis-cache)

| Key | Value | TTL |
|---|---|---|
| `gw:shops:<country>:<postcode>` | response JSON | 24 h |

## Errors

| Case | HTTP | `code` |
|---|---|---|
| Bad postcode format | 400 | `validation` |
| Country not covered | 400 | `country_not_covered` |
| Unknown postcode or area | 404 | `not_found` |
| Correction with none or several of the three fields | 400 | `validation` |
| Too many suggestions | 429 | `rate_limited` |

## Acceptance criteria

- [ ] A Warsaw postcode returns shops with coordinates and no distance field.
- [ ] Hamburg and Munich postcodes return Aldi Nord and Aldi Süd shops respectively (verification 8, data from offers).
- [ ] With location permission denied, the postcode path gives the app the same shop list (verification 10).
- [ ] Request logs show only the first 2 characters of the postcode.
- [ ] Typing "Warsz" in `GET /v1/areas` returns Warsaw areas; the chosen area id then lists the same shops as its postcode.
