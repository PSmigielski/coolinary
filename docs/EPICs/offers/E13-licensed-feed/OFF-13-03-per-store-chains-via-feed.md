# OFF-13-03 Per-store chains through the feed (Rewe, Edeka, Netto DE)

| | |
|---|---|
| Service | offers |
| Epic | [OFF-13 Licensed offers feed](README.md) |
| Release | R3 |
| Size | M |
| Status | todo |
| Depends on | OFF-13-02, OFF-05-07 |
| Unblocks | — |

## Goal

German households near Rewe, Edeka or Netto Marken-Discount get those shops' own per-store offers, which these chains
do not allow to be read from their sites (features.md §3: "all three sites refused automated access").

## Business rules

1. Chains `rewe-de`, `edeka-de`, `netto-de` have `edition_model = per_store`; their store editions come from linked feed
   stores (OFF-13-01).
2. Store editions are fetched from the feed only when wanted (OFF-05-07 rules, same 30-day tail), if the provider
   allows per-store requests; otherwise all stores are ingested and only wanted editions are kept.
3. Edeka's independent merchants may run different offers per store; nothing is shared between stores unless the feed
   says so.
4. Users of these chains' shops see `not_covered` until the chain's coverage is `feed`.

## Contract

No change.

## Data

No schema change.

## Errors

| Case | Behaviour |
|---|---|
| Wanted store not in the feed | Shop status `store_not_linked` (OFF-05-06) |

## Acceptance criteria

- [ ] Two Edeka shops linked to different feed stores show different offers.
- [ ] An unwanted Rewe store's offers are not stored.
