# OFF-08-04 Accuracy record and metrics

| | |
|---|---|
| Service | offers |
| Epic | [OFF-08 Accuracy gate and weekly spot-checks](README.md) |
| Release | R2 |
| Size | S |
| Status | todo |
| Depends on | OFF-08-01 |
| Unblocks | ADM-06, ADM-07 |

## Goal

The founder sees, per chain, how well offers are read over time and where errors come from (product, price, dates,
conditions), to decide on prompts, adapters or taking a chain down.

## Business rules

1. Per chain and ISO week: sample score, counts per field of right and wrong, missed offers, held-back share
   (held / extracted), rejected share, published count.
2. Rolling 4-week score = all correct / all judged over the last 4 closed weekly samples.
3. `GetChainAccuracy(chain_id, weeks)` returns the weekly rows and the rolling score.
4. Prometheus gauges: `offers_chain_accuracy_percent{chain,kind="last_sample|rolling_4w"}`,
   `offers_held_share{chain}`. A Grafana panel is added to *Coolinary Services*.
5. Figures contain no user data.

## Contract

`admin.proto`: `rpc GetChainAccuracy(GetChainAccuracyRequest) returns (GetChainAccuracyResponse);`

## Data

Computed from `review_samples`, `review_items`, `offers`; no new table.

## Errors

| Case | Behaviour |
|---|---|
| Unknown chain | `NOT_FOUND` |

## Acceptance criteria

- [ ] After 4 weekly samples of 97, 98, 96, 99 (each of 20 offers), the rolling score is 97.5.
- [ ] The gauges appear in Prometheus and on the dashboard.
