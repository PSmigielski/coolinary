# ADM-06-01 Flyer source status board

| | |
|---|---|
| Service | admin |
| Epic | [ADM-06 Flyer sources, offers and shops review](README.md) |
| Release | R2 |
| Size | M |
| Status | todo |
| Depends on | ADM-01-05; OFF-05-08 (`ListSources`, `TriggerSync`, `ListFlyers`, `ListEditions`), OFF-08-01 (`ListReviewSamples`), OFF-08-04 (`GetChainAccuracy`) |
| Unblocks | ADM-06-02, ADM-06-03 |

## Goal

One page shows, per chain, whether its flyers are being read, when the last new flyer arrived, and what failed, so a
blocked or changed source is noticed the same day (features.md §3 "Keeping flyers in sync", step 9).

## Business rules

1. The board is `ListSources` (OFF-05-08 rule 1): per source its chain, the chain's coverage (`user_flyers_only`,
   `syncing`, `live`, `removed`, OFF-08-02), the source status, terms approved or not, last success, last error,
   consecutive failures and the last 10 checks.
2. A source's detail page adds: its flyers from `ListFlyers(chain)` (state, valid dates, version and page counts,
   offers by state), its editions from `ListEditions(chain)` (e.g. Kaufland stores with users in range: demand,
   wanted state, last check), its review samples from `ListReviewSamples(chain)` (last go-live result, last 4
   weekly results) and its accuracy record from `GetChainAccuracy(chain)`.
3. Sort: sources with consecutive failures first, then `live` chains by oldest last success, then the rest by chain
   name.
4. "Check now" calls `TriggerSync` for the source, or for one edition of it; offers refuses while a check of that
   source is running (409).
5. Read routes need role `admin`; "check now" is audited (`offers.source.check_now`).

## Contract

| Method | Path | Body / response |
|---|---|---|
| GET | `/admin/v1/offers/sources` | `{ items }` |
| GET | `/admin/v1/offers/sources/{id}` | source with flyers, editions, samples, accuracy |
| POST | `/admin/v1/offers/sources/{id}/check-now` | `{ edition_id? }` → `{ check_id }` |

## Data

None in admin.

## Errors

| Case | HTTP |
|---|---|
| Unknown source | 404 |
| A check of this source is running | 409 `check_running` |

## Acceptance criteria

- [ ] A source whose list page returns 403 shows its last error and failure count at the top of the board.
- [ ] "Check now" while that source's check is running → 409.
