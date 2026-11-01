# OFF-10-05 Shop for a user flyer, in any country

| | |
|---|---|
| Service | offers |
| Epic | [OFF-10 User-supplied flyers](README.md) |
| Release | R2 |
| Size | S |
| Status | todo |
| Depends on | OFF-02-05 |
| Unblocks | OFF-10-01 |

## Goal

A user flyer always belongs to a shop, and this works in every country, also where offers has no OSM shops or
postcode data yet (features.md §2: "This works in every country from the first flyer release").

## Business rules

1. The shop is one of: an active shop from `ListShopsNearArea` or my shops; the household's own pending shop; or a
   new household shop created with `SubmitShop` (OFF-02-05).
2. For countries other than PL and DE, `SubmitShop` accepts any EU country code with "other chain" and a name; the
   placeholder chain `other-<cc>` is created on first use with that country's currency and time zone from a static
   country table (EU 27). Such shops stay household-only and are not put in the moderation queue, because offers has
   no data to check them against.
3. A shop of a covered chain may be used too (the chain's flyer may be stale, or the user has a local edition).
4. The shop decides the currency of the flyer's prices and the time zone of its dates.

## Contract

Uses `SubmitShop`; no new RPC. The country table is code, not a contract.

## Data

Migration: country table (code, currency, time zone) for the EU 27 as static data.

## Errors

| Case | gRPC status |
|---|---|
| Non-EU country | `INVALID_ARGUMENT` |

## Acceptance criteria

- [ ] A household in Austria creates a private "Billa" shop (chain `other-at`, EUR, Europe/Vienna), uploads a flyer
      and saves its offers.
- [ ] That shop never appears in another household's shop list or in the moderation queue.
