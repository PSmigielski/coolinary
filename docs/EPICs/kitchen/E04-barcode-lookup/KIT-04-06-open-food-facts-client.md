# KIT-04-06 Open Food Facts client (level 3)

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-04 Barcode lookup chain](README.md) |
| Release | R1 (GATE: licence check before production) |
| Size | M |
| Status | todo |
| Depends on | KIT-04-01 |
| Unblocks | KIT-04-07, KIT-04-09, KIT-06-07 |

## Goal

Codes Coolinary doesn't know yet are looked up live in Open Food Facts, within its terms: one request per real scan,
identified, rate-limited and credited (features.md §4 table "Open Food Facts … Live lookups are allowed when one
call is one real scan").

## Business rules

1. Config: `OFF_ENABLED` (default `false`), `OFF_BASE_URL` (`https://world.openfoodfacts.org`), `OFF_USER_AGENT`
   (`Coolinary/<version> (<contact address>)`), `OFF_RATE_PER_MIN` (default 80).
2. A request is made only from `ResolveBarcode` with `user_scan = true` (the user scanned it now). Never from jobs,
   loops, prefetching or backfills. Bulk data would use OFF's downloads, which R1 does not do.
3. Request: `GET /api/v2/product/<code>?fields=code,product_name,product_name_pl,product_name_de,product_name_en,
   generic_name,brands,quantity,product_quantity,product_quantity_unit,categories_tags,countries_tags,
   allergens_tags,traces_tags,nutriments,nutrition_data_per,lang` with the User-Agent header. The code sent is the
   13-digit form (8 digits for GTIN-8).
4. Timeout 1.5 s; no retry inside the request. `status = 0` or HTTP 404 → miss. HTTP 429, 5xx or a timeout → the
   level is skipped and the response is marked `degraded`.
5. Rate limit shared by all kitchen instances: a token bucket in redis-state (`kitchen:off:bucket`), refilled to
   `OFF_RATE_PER_MIN` per minute (OFF's documentation allows 100 product reads per minute; confirm before go-live).
   No token → the level is skipped (`degraded`).
6. Circuit breaker: 5 failures in a row → skip OFF for 60 s.
7. Cache in redis-cache `kitchen:off:<gtin>`: a hit for 7 days, a miss for 24 hours.
8. GATE: production enables `OFF_ENABLED` only after the licence check (features.md §4 open points, §8 risk "The
   open product database must be shared alike").
9. Metrics: `kitchen_off_requests_total{result=hit|miss|error|rate_limited|circuit_open}`,
   `kitchen_off_request_seconds`.

## Contract

No RPC (client used by KIT-04-09).

## Data

No schema change (storage in KIT-04-07).

## Errors

OFF problems never fail `ResolveBarcode`; they set `degraded = true`.

## Acceptance criteria

- [ ] With a stub OFF, a scan of an unknown code makes exactly one request; the next scan within 7 days makes none.
- [ ] 100 scans in one minute with `OFF_RATE_PER_MIN=80` make at most 80 requests; the rest are `degraded`.
- [ ] Every request carries the configured User-Agent.
- [ ] With `OFF_ENABLED=false` no request is ever made.
