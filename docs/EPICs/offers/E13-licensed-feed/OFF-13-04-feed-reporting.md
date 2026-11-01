# OFF-13-04 View and click reporting to the network

| | |
|---|---|
| Service | offers |
| Epic | [OFF-13 Licensed offers feed](README.md) |
| Release | R3 |
| Size | M |
| Status | todo |
| Depends on | OFF-13-02 |
| Unblocks | the offers-network revenue line (business model §4) |

## Goal

If the agreement pays for flyer views or requires usage reports, offers reports them without sending anything that
identifies a user.

## Business rules

1. `GATE:` only if the signed agreement asks for reporting; otherwise this task is not built.
2. The gateway records, for feed offers only, "shown in a list", "opened" and "source link opened" events per offer
   (GW-10), and calls `RecordFeedUsage` in batches.
3. offers aggregates per offer, per day and per country, and sends daily totals to the provider. No user id, household
   id, device id, IP address or position is stored or sent; postcodes only if the agreement requires them and only
   at the 2-digit level (PL `00`, DE `80`).
4. Totals are kept 13 months for reconciliation, then deleted.

## Contract

```proto
rpc RecordFeedUsage(RecordFeedUsageRequest) returns (RecordFeedUsageResponse);
message RecordFeedUsageRequest { repeated FeedUsage usage = 1; }
message FeedUsage { string offer_id = 1; string kind = 2; int32 count = 3; string date = 4; }
```

## Data

`offers/migrations/000NN_feed_usage.sql`

```sql
CREATE TABLE feed_usage_daily (
    offer_id text NOT NULL,
    day      date NOT NULL,
    kind     text NOT NULL CHECK (kind IN ('shown', 'opened', 'source_opened')),
    count    int NOT NULL,
    reported_at timestamptz,
    PRIMARY KEY (offer_id, day, kind)
);
```

## Errors

| Case | Behaviour |
|---|---|
| Provider report endpoint down | Retry hourly; alert after 48 h unreported |
| Usage for a non-feed offer | Ignored |

## Acceptance criteria

- [ ] A day of test usage is sent as totals; the sent payload contains no user or household fields.
