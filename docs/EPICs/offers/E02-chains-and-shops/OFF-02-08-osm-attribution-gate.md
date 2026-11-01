# OFF-02-08 OSM attribution and ODbL gate

| | |
|---|---|
| Service | offers |
| Epic | [OFF-02 Chains and shops](README.md) |
| Release | R2 |
| Size | S |
| Depends on | OFF-02-03 |
| Unblocks | release of OFF-03-03 to users |

## Goal

Coolinary credits OpenStreetMap wherever shop data is shown and settles the share-alike question before shop data
reaches users.

## Business rules

1. Every response that contains shops (`ListShopsNearArea`, `ListMyShops`, `BrowseOffers` shop list) carries
   `attribution = "© OpenStreetMap contributors"` and `attribution_url = "https://www.openstreetmap.org/copyright"`.
   The app shows it where the shops are shown (frontend task later).
2. User-added shops and corrections are kept in their own rows and tables (`source = 'user'`, `shop_overrides`), so the
   OSM-derived data can be told apart and published separately if the licence requires it.
3. `GATE:` before R2 goes live, a legal check answers: is Coolinary's shop table a "derivative database" of OSM under
   the ODbL when combined with user corrections, and if so what must be offered (the derived database, or only the
   changed rows) and how. The result is written into this task file and into docs/features.md "Open points".
4. Until the gate passes, shop data is served only in dev and staging.

## Contract

```proto
message Attribution {
  string text = 1;
  string url = 2;
}
// Added as field `attribution` to ListShopsNearAreaResponse, ListMyShopsResponse, BrowseOffersResponse.
```

## Data

No schema change.

## Errors

| Case | Behaviour |
|---|---|
| Gate not passed and `ENVIRONMENT=prod` | Shop RPCs return `FAILED_PRECONDITION` "shop data not released" |

## Acceptance criteria

- [ ] All three responses include the attribution.
- [ ] The legal answer is recorded and the gate ticked in offers/README.md.
