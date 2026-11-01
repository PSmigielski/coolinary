# GW-10-05 Trip price matrix

| | |
|---|---|
| Service | gateway |
| Epic | [GW-10 Shops, offers, flyers and price matrix routes](README.md) |
| Release | R2 (receipt prices from R3) |
| Size | M |
| Status | todo |
| Depends on | GW-08-01, GW-03-03, PLN-19 (`GetTripPriceMatrix`) |
| Unblocks | MOB-03 (the phone's trip plan, a later frontend task) |

## Goal

The phone plans the shopping trip itself (decision 7): it needs, for every list line, the price at every candidate
shop. Planner builds that matrix from flyer offers, then receipt prices (R3), then a reference price the same for
every shop (features.md §2 "Where prices come from"). Plus only.

## Business rules

1. **Plus only** (`plus: "trip_plan"`). Body `{ "shopIds": [...] }`: 1–10 shop ids the phone has already filtered
   by range and walking limit, nearest first (planner's cap, PLN-19-01; decided, Q162). 10 shops with at most
   3 per trip is 175 sets, still trivial on the phone.
2. Calls planner `GetTripPriceMatrix(household_id, viewer_member_id, shop_ids)` with the `heavy` deadline (10 s).
   The matrix covers every open line of the household's list at `listVersion`.
3. Response, as planner builds it: `shops` (with chain and shop name, position, total for the whole list,
   unpriced cell count, `about`), `lines` (line id, ingredient, name, amount bought in whole packs), and `cells`
   row-major (lines × shops): cost or none, `source` `flyer` | `receipt` | `reference`, `about`, `offerId` when
   `flyer`, condition flags (e.g. loyalty card). Reference prices are flagged so the phone can say "until receipt
   prices build up, only flyer offers separate one shop from another".
4. The phone does the set search, the tie-breaks (fewer shops, then shorter walk) and the minimum-saving rule; the
   matrix must let it reproduce features.md Verification 13 exactly.
5. `ignoredShopIds` (shops planner could not price at all) and `unpricedLineIds` (free-text lines, or no price
   anywhere) are passed through.
6. Nothing about distances or positions is sent or returned (D13).

## Contract

`POST /v1/list/price-matrix` — `auth: true`, `plus: "trip_plan"`, tag `list`, idempotency none (read-only POST).

```json
request { "shopIds": ["osm-lidl-1", "osm-kaufland-2", "osm-aldi-3"] }
200 {
  "listVersion": 44, "currency": "EUR",
  "shops": [ { "shopId": "osm-lidl-1", "chainName": "Lidl", "position": 1,
               "total": { "amountMinor": 1065, "currency": "EUR" }, "unpricedCells": 0, "about": true } ],
  "lines": [ { "lineId": "lin_1", "ingredientId": "ing_chicken_thigh", "name": "Chicken thighs 1 kg" } ],
  "cells": [ { "lineId": "lin_1", "shopId": "osm-lidl-1", "cost": { "amountMinor": 389, "currency": "EUR" },
               "source": "flyer", "about": true, "offerId": "off_..", "conditionFlags": [] } ],
  "ignoredShopIds": [], "unpricedLineIds": []
}
```

## Data

No storage.

## Errors

| Case | HTTP | `code` |
|---|---|---|
| Free | 403 | `plus_required` |
| 0 or > 10 shops | 400 | `validation` |
| Empty list | 422 | `list_empty` |

## Acceptance criteria

- [ ] With the features.md §2 worked example data, totals are Lidl 10.65, Kaufland 11.25, Aldi 12.11, and the per-line
      prices let the phone find Lidl + Kaufland = 9.15 (verification 13, phone side tested later).
- [ ] A line on offer at one shop shows `source: "flyer"` there and `reference` elsewhere.
- [ ] The request schema has no coordinate or distance field.
