# OFF-05-06 Source health, alerts and flyer status for users

| | |
|---|---|
| Service | offers |
| Epic | [OFF-05 Flyer source framework and sync scheduler](README.md) |
| Release | R2 |
| Size | M |
| Depends on | OFF-05-04, OFF-01-04 (notifier SendAlert) |
| Unblocks | GW-10 (flyer status route), PLN-17 |

## Goal

When a chain changes or blocks its flyer page, the founder hears about it the same day, and users of that chain see
"flyer not updated, scan it yourself" instead of silently old offers (features.md §3 step 9, Risks "A chain changes
its flyer page or blocks access", Verification 7).

## Business rules

1. Source status, recomputed after every check and every hour:
   - `disabled`: source not enabled, or chain `removed`.
   - `failing`: the last 3 checks failed, or one check failed with `ErrBlocked` or `ErrRobots`.
   - `stale`: not failing, but no successful check for 36 h, or (for `live` chains) no published flyer of the
     national edition is valid today.
   - `ok`: otherwise.
2. Each transition into `failing` or `stale` sends one founder alert through notifier `SendAlert`, key
   `offers.source.<source_id>`, with the last error text and the time of the last success. A transition back to `ok`
   sends one "recovered" alert with the same key. Notifier deduplicates repeats (NTF-05).
3. Per shop flyer status for users (`GetFlyerStatus`), from the shop's chain and editions:

   | Situation | `code` |
   |---|---|
   | Chain `user_flyers_only` or `syncing` | `not_covered` |
   | Chain `removed` | `not_covered` |
   | Per-store chain and shop not linked to a store | `store_not_linked` |
   | Source `failing`, `stale` or `disabled`, or the shop's edition has no published flyer valid today | `not_updated` |
   | Otherwise | `ok` |

   Every code except `ok` means the app offers "scan the flyer yourself" (OFF-10).
4. The status also gives `last_updated_at` (last successful check of the edition) and the end date of the latest
   published flyer.
5. A simulated failure (admin RPC `SimulateSourceFailure` in dev and staging only) forces the next check to fail with
   `ErrBlocked`, to test rule 2 and 3 end to end (Verification 7, last part).

## Contract

`contracts/coolinary/offers/v1/flyers.proto`

```proto
rpc GetFlyerStatus(GetFlyerStatusRequest) returns (GetFlyerStatusResponse);

message GetFlyerStatusRequest {
  string household_id = 1;
  repeated string shop_ids = 2;    // empty = the household's active my shops
}
message GetFlyerStatusResponse { repeated ShopFlyerStatus statuses = 1; }
message ShopFlyerStatus {
  string shop_id = 1;
  FlyerStatusCode code = 2;
  string last_updated_at = 3;      // RFC 3339, empty if never
  string latest_valid_to = 4;      // YYYY-MM-DD
}
enum FlyerStatusCode {
  FLYER_STATUS_CODE_UNSPECIFIED = 0;
  FLYER_STATUS_CODE_OK = 1;
  FLYER_STATUS_CODE_NOT_COVERED = 2;
  FLYER_STATUS_CODE_STORE_NOT_LINKED = 3;
  FLYER_STATUS_CODE_NOT_UPDATED = 4;
}
```

Admin: `SimulateSourceFailure(source_id)` (refused with `FAILED_PRECONDITION` when `ENVIRONMENT=prod`).

## Data

No schema change (uses `flyer_sources.status`, `flyer_checks`).

## Errors

| Case | Behaviour |
|---|---|
| Unknown shop id | `INVALID_ARGUMENT` |
| `SimulateSourceFailure` in prod | `FAILED_PRECONDITION` |

## Acceptance criteria

- [ ] Three failing fixture checks make the source `failing`, send exactly one alert to Mailpit, and
      `GetFlyerStatus` returns `not_updated` for its shops.
- [ ] The next successful check returns the source to `ok`, sends one "recovered" alert, and the status returns `ok`.
- [ ] A Biedronka shop while Biedronka is `user_flyers_only` returns `not_covered`.
