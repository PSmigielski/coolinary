# OFF-13-02 Ingest feed offers

| | |
|---|---|
| Service | offers |
| Epic | [OFF-13 Licensed offers feed](README.md) |
| Release | R3 |
| Size | L |
| Status | todo |
| Depends on | OFF-13-01, OFF-07-03, OFF-07-04, OFF-07-05, OFF-07-07 |
| Unblocks | OFF-13-03, OFF-13-04 |

## Goal

Feed offers reach users with the same checks and the same behaviour as flyer offers, without AI cost.

## Business rules

1. A job every 6 hours fetches changed offers for mapped retailers (incremental, using the provider's change marker).
2. Each feed offer becomes an offer with source kind `feed`, `feed_item_id` = the provider's id, the edition of its
   store (store edition) or retailer (national edition), and the provider's link as `source_url` (as the agreement
   says).
3. Fields arrive structured; they still go through normalisation (OFF-07-03: prices to minor units, packs, unit
   prices, conditions), matching (OFF-07-04), duplicate checks (OFF-07-08) and validation (OFF-07-05). Held feed
   offers go to the same hold-back queue.
4. A feed offer changed by the provider replaces the earlier row (old one `superseded`); one deleted by the provider
   becomes `expired` with `offers.expired`.
5. Publishing and expiry follow OFF-07-07; feed offers of a chain are visible only when its coverage is `feed`.
6. If the agreement forbids storing feed prices beyond the offer's validity, feed price observations (OFF-11-01) are
   deleted at expiry instead of after 400 days.
7. Metrics `offers_feed_ingested_total{chain,result}`.

## Contract

No public change.

## Data

`CREATE UNIQUE INDEX offers_feed_item ON offers (feed_item_id) WHERE source_kind = 'feed' AND state = 'published';`

## Errors

| Case | Behaviour |
|---|---|
| Feed unreachable | Run `failed`; after 24 h without success the feed chains' sources become `stale` (OFF-05-06 rules apply) |

## Acceptance criteria

- [ ] Sample feed offers for Rewe appear, once the chain is `feed`, with prices, packs and conditions normalised.
- [ ] A provider deletion removes the offer and sends `offers.expired`.
- [ ] No `SubmitBatch` or `ExtractFlyerOffers` call is made for feed offers.
