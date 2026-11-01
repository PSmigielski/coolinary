# GW-03-05 AI usage metrics and dashboard panel

| | |
|---|---|
| Service | gateway |
| Epic | [GW-03 Free-tier AI quota and rate limits](README.md) |
| Release | R1 |
| Size | S |
| Status | todo |
| Depends on | GW-03-01 |
| Unblocks | ADM-07, business-model guardrail tracking |

## Goal

The founder can see weekly how often Free and Plus households use AI, how often they hit the quota or the paywall,
and how often services fall back to non-AI results.

## Business rules

1. Counter `gateway_ai_requests_total{feature, plan, outcome}`:
   - `feature`: `day_plan`, `scan_photo`, `scan_label`, `scan_text`, `scan_receipt`, `swap_meal`, `tool_rewrite`,
     `restyle`, `week_plan`, `plan_from_deals`, `cheapest_top_up`, `own_flyer`.
   - `plan`: `free`, `plus` (anonymous counts as `free`; extra label `anonymous` `true`/`false`).
   - `outcome`: `ai` (AI result), `fallback` (`ai_used: false`), `quota_reached`, `plus_required`, `fair_use`,
     `error`.
2. Gauge per week is not kept in Prometheus (no household labels); weekly per-household numbers come from the
   quota keys (GW-03-01) and from Langfuse (AI-02).
3. Grafana "Coolinary Services" dashboard gets a row "AI use": requests per feature and outcome, quota_reached per
   day, fallback share per feature.

## Contract

Metric names above; dashboard JSON in `platform/grafana/dashboards/services.json`.

## Data

No storage.

## Errors

None.

## Acceptance criteria

- [ ] A Free day plan answered from the library increments `{feature="day_plan",plan="free",outcome="fallback"}`.
- [ ] A 429 quota response increments `outcome="quota_reached"`.
- [ ] The dashboard row renders with data from `make dev-full`.
