# OFF-07-07 Publish, "next week" and expiry

| | |
|---|---|
| Service | offers |
| Epic | [OFF-07 Offer extraction, validation and publishing](README.md) |
| Release | R2 |
| Size | L |
| Depends on | OFF-07-05, OFF-01-01, OFF-01-05 |
| Unblocks | OFF-09-*, OFF-11-01, PLN-17 (offers.published / offers.expired) |

## Goal

Offers go live on their start date, flyers published ahead show as "next week", and offers and their list badges
disappear after their end date (features.md §3 steps 6–7, Verification 6, 7).

## Business rules

1. A publish run every 15 minutes, and at 00:05 local time of each country.
2. A flyer version is publishable when no offer of it is still `extracted`. Its `ready` offers become `published`
   (`published_at = now()`), and the flyer `published`, only if the chain's coverage is `live` (or `feed`). For
   `syncing` chains the offers stay `ready`, invisible to users, until the chain goes live (OFF-08-02); then they are
   published at once if still valid.
3. Visibility to users (used by every read RPC): state `published`, `valid_to ≥ today` in the chain's time zone,
   chain `live` or `feed`. An offer with `valid_from > today` is shown as upcoming ("next week").
4. End dates are inclusive: an offer valid to 17 October is visible all of 17 October and gone from 00:00 on
   18 October local time. The 00:05 run sets those offers to `expired`. (Verification 7 "expired offers are gone on
   their end date" is read as: gone once the end date is over.)
5. Re-read flyers: when a new version of a flyer becomes publishable, the old version's published offers become
   `superseded` and the new ones `published`, in one transaction, so users never see both or neither.
6. Events, written through the outbox in the same transaction:
   - `offers.published` per version (and per released held offer): `offer_ids` (at most 1,000 per event; more are
     split), `chain_id`, `valid_from`, `valid_to` (the version's dates). For user-flyer offers the envelope
     `household_id` is set (OFF-10-04).
   - `offers.expired`: `offer_ids` of offers that expired, were superseded, rejected after publication, or removed
     (OFF-08-05).
7. Each run bumps the cache version of the affected editions so browse results refresh (OFF-09-01).
8. Metrics: `offers_published_total{chain}`, `offers_expired_total{chain}`, gauge `offers_visible{chain}`.

## Contract

No public RPC. Event payloads as in docs/EPICs/README.md "Events".

## Data

No schema change.

## Errors

| Case | Behaviour |
|---|---|
| Publish transaction fails | Rolled back and retried at the next run; outbox rows roll back with it, so events and states stay in step |

## Acceptance criteria

- [ ] A live chain's flyer for next Monday, extracted on Thursday, shows its offers as upcoming from Thursday and as
      current from Monday 00:05.
- [ ] On the day after the end date, the offers are `expired` by 00:05 and an `offers.expired` event lists them.
- [ ] A re-read version replaces the old offers in one transaction (no moment with zero or double offers for that
      flyer, checked by a concurrent reader in a test).
- [ ] Five offers added to a list show their badge; after the end date planner receives their ids in
      `offers.expired` (Verification 6, offers side).
