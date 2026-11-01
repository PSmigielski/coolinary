# GW-10-02 My shops and shopping preferences

| | |
|---|---|
| Service | gateway |
| Epic | [GW-10 Shops, offers, flyers and price matrix routes](README.md) |
| Release | R2 |
| Size | S |
| Status | todo |
| Depends on | GW-10-01, OFF-04 (`ListMyShops`, `SetMyShops`, `GetShoppingPreferences`, `SetShoppingPreferences`) |
| Unblocks | GW-10-03, GW-07-10, GW-10-05, MOB-03 |

## Goal

The household saves the shops it uses (nearest first, as the phone computed) and its "how you shop" settings, so
every household member and every phone sees the same, and offers knows which shops' offers to show (Free: 2 shops).

## Business rules

1. `PUT /v1/my-shops` body `{ "shopIds": [...] }` — in the order the phone sorted them, nearest first (D13). Calls
   offers `SetMyShops(household_id, shop_ids)`, which replaces the list; duplicates keep their first position. Free
   households may save at most **2** shops and Plus at most 40 (offers' limit, returned as `limit`); more → offers
   `FAILED_PRECONDITION` `plus_required` (Free) → 403, or its own code (Plus cap) → 422. The phone sends its 2
   nearest first on Free (features.md §6: "offers from the 2 nearest shops").
2. The phone sends the list again when the user's position changes enough to change the order; planner uses this
   order for "plan from deals" (GW-07-10).
3. `GET /v1/my-shops` → `ListMyShops`: active shops in order, then inactive and closed ones, each with chain,
   coverage, editions, `state` and position, plus `limit`.
4. `PUT /v1/shopping-preferences` body `{ "rangeKm": 1 | 3 | 10, "walkBetweenM": 0–1500 (step 100), "shopsPerTrip":
   1–3, "minExtraShopSaving": Money }` → `SetShoppingPreferences`. `GET /v1/shopping-preferences?country=` →
   `GetShoppingPreferences(household_id, country_code)`, which returns offers' defaults per country when never set
   (e.g. €2 or its local equivalent).
5. Preferences are used by the phone (range filter, trip plan); the server only stores them for sync.
6. Offers keeps no version on these resources: both `PUT`s replace the whole value, last writer wins; no `If-Match`.

## Contract

| Method | Path | 2xx |
|---|---|---|
| GET | `/v1/my-shops` | 200 `{ "limit": 2, "shops": [ { "shop", "coverage", "state": "active", "position": 1 } ] }` |
| PUT | `/v1/my-shops` | 200 same |
| GET | `/v1/shopping-preferences?country=PL` | 200 `{ "rangeKm": 3, "walkBetweenM": 900, "shopsPerTrip": 2, "minExtraShopSaving": { "amountMinor": 800, "currency": "PLN" }, "updatedAt": "..." }` |
| PUT | `/v1/shopping-preferences` | 200 same |

Tag `shops`, `auth: true`, idempotency optional.

## Data

No storage.

## Errors

| Case | HTTP | `code` |
|---|---|---|
| Bad range, walk or shops-per-trip value | 400 | `validation` |
| Free with > 2 shops | 403 | `plus_required` |
| Plus with more than 40 shops, or an unknown shop id | 422 | offers' code |

## Acceptance criteria

- [ ] A Free household saving 3 shops gets 403 `plus_required`; saving 2 works.
- [ ] A Plus household saves 8 shops in the given order and `GET` returns the same order.
- [ ] `walkBetweenM: 2000` is refused with 400.
