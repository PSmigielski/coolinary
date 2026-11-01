# MOB-03-08 Shops, offers, flyers and price matrix hooks

| | |
|---|---|
| Service | mobile |
| Epic | [MOB-03 Query keys, hooks and cache rules](README.md) |
| Release | R2 (shops, my shops, offers, own flyer, price matrix, plan from deals), R3 (cheapest top-up) |
| Size | M |
| Status | todo |
| Depends on | MOB-03-01, MOB-03-02, MOB-07-03, GW-10 |
| Unblocks | later offers screens and the later on-device trip planner |

## Goal

Hooks for flyers and offers that keep the user's location on the phone: the server gets a postcode at most, and the
phone sorts shops by distance and plans the trip itself (README decision 7, D13).

## Business rules

1. Hooks and routes (GW-10 route catalogue in docs/EPICs/gateway/README.md):

   | Hook | HTTP | Offline |
   |---|---|---|
   | `useAreas(country, q)` | `GET /v1/areas?country=&q=` → postcode areas or towns (typed postcode or town) | not offline |
   | `useShopsNear(area)` | `GET /v1/shops?areaId=` or `?postcode=&country=` → shops with chain, name, address, latitude, longitude, editions | cached 1 h |
   | `useMyShops()` / `useSetMyShops()` | `GET` / `PUT /v1/my-shops` `{ shopIds }` (nearest first) | queued |
   | `useShoppingPreferences()` / `useSetShoppingPreferences()` | `GET` / `PUT /v1/shopping-preferences` `{ rangeKm: 1\|3\|10, walkBetweenM: 0..1500, shopsPerTrip: 1..3, minExtraShopSaving: Money }` | queued |
   | `useOffers(filters)` | `GET /v1/offers?shopIds=&category=&q=&endingBy=&view=current\|upcoming\|all&pageToken=` | current view cached |
   | `useOffer(id)` / `useOfferCategories()` | `GET /v1/offers/{id}` / `GET /v1/offers/categories?shopIds=&view=` | cached |
   | `useFlyerStatus()` | `GET /v1/flyers/status?shopIds=` (per shop: `ok`, `not_covered`, `store_not_linked`, `not_updated`) | cached |
   | `useUploadFlyer()` / `useUserFlyer(id)` | `POST /v1/flyers` (multipart or `{ shopId, url }`, MOB-07-03) / `GET /v1/flyers/{flyerId}` | never (AI) |
   | `useConfirmFlyer()` / `useDeleteFlyer()` | `POST /v1/flyers/{flyerId}/confirm` `{ validFrom, validTo, decisions }` / `DELETE /v1/flyers/{flyerId}` | never |
   | `usePriceMatrix()` | `POST /v1/list/price-matrix` `{ shopIds }` (1–10) → shops, lines and cells with cost, source (`flyer`, `receipt`, `reference`) and `about` | never cached |
   | `usePlanFromDeals()` | `POST /v1/plans/from-deals` `{ date, replace?, useAi? }` (planner uses the saved my-shops order) | never (AI) |
   | `useCheapestTopUp()` (R3) | `POST /v1/list/cheapest-top-up` `{ weekPlanId, shopIds }`, then `POST .../{topUpId}/apply` `{ swapKeys }` | never (AI) |
2. The postcode comes from the user's typing or from the phone's own reverse geocoding (a later frontend task);
   these hooks accept only `{ postcode, country }`, never coordinates.
3. Shop coordinates come from the server; distances, the range filter and the nearest-first order are computed on the
   phone by a later module. Hooks that need an order take `shopIds` already sorted nearest first.
4. Free households may keep 2 shops (`limits.shops` from MOB-02-05); `useSetMyShops()` sends at most that many and
   the server answers `plus_required` otherwise.
5. `usePriceMatrix()` is Plus only; its result feeds the later on-device trip planner (features.md "Planning the
   shopping trip"). The hook returns it unchanged and never persists it.
6. `useOffers` pages of 50; offers disappear on their end date, so the cache for `view=current` is dropped every day
   at local midnight.

## Contract

```ts
interface PriceMatrix { lines: { lineId: string; prices: { shopId: string; priceMinor: number | null; currency: string; source: "flyer" | "receipt" | "reference"; offerId?: string; approx: boolean }[] }[] }
```

## Data

Query cache per MOB-03-01; price matrix and flyer uploads never persisted.

## Errors

| Case | Behaviour |
|---|---|
| Third shop on Free | `plusRequired` |
| Postcode unknown | `validation` / `postcode_unknown` |

## Acceptance criteria

- [ ] Verification 10 from the app side: with location permission denied, a typed postcode returns shops and
      offers; no request contains coordinates of the user.
- [ ] A network capture of every offers call contains no latitude/longitude field.
