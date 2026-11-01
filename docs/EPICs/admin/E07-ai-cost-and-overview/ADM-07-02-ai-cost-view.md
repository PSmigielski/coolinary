# ADM-07-02 AI cost and guardrails

| | |
|---|---|
| Service | admin |
| Epic | [ADM-07 AI cost and system overview](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | ADM-07-03; AI-02-02 (cost metrics), AI-02-04 (`AiAdminService.GetSpendStatus`) |
| Unblocks | — |

## Goal

The founder sees what AI costs per feature and per plan, and whether the business model's two limits hold: a Free
household costs at most €0.10 a month, and AI takes at most 15% of what a Plus household pays (business model §6,
features.md verification 11).

## Business rules

1. Cost by job and plan for the current ISO week, the previous week, the current month and the previous month,
   read from Prometheus (`PROMETHEUS_URL`) with AI-02-02's counter:
   `sum by (job, plan, mode) (increase(ai_cost_eur_micros_total[<range>])) / 1e6` (euros).
2. Batch jobs (once for everyone, `mode="batch"`) are shown separately from per-request jobs (`mode="sync"`).
3. Live spend limits: `AiAdminService.GetSpendStatus` (AI-02-04 rule 9) → for each bucket (`sync`, `sync_free`,
   `batch`) the spend in its period, the limit, the ratio and whether it is reached (D11).
4. Guardrail 1: per-request AI cost of Free households in the last 30 days ÷ Free households active in the last 30
   days (`households_active_30d_free`, ADM-07-03) ≤ €0.10. Shown as a number with a red mark when over.
5. Guardrail 2: per-request AI cost of Plus households in the last 30 days ÷ active Plus households
   (`households_active_30d_plus`) ≤ 15% of the net
   monthly revenue per Plus household. Until payments exist (decision 8) the net revenue is a config value
   `PLUS_NET_EUR_PER_MONTH`, default 1.64 (€19.7 a year after VAT and store fee, business model §5, ÷ 12).
   Decided (Q049): default 1.64, kept as config until real revenue exists.
6. When Prometheus is not running (dev without `make dev-full`), the route answers 200 with
   `metrics_unavailable: true` and still shows the spend status.
7. Role `admin` only; not audited.

## Contract

| Method | Path | Response |
|---|---|---|
| GET | `/admin/v1/ai/costs` | `{ periods: { week, last_week, month, last_month }: [{ job, plan, mode, eur }], spend: [{ bucket, period, spent_eur, limit_eur, ratio, limited }], guardrails: { free_per_household_eur, free_limit_eur: 0.10, plus_share_pct, plus_limit_pct: 15 }, metrics_unavailable }` |

## Data

New config `PROMETHEUS_URL`, `PLUS_NET_EUR_PER_MONTH`.

## Errors

| Case | HTTP |
|---|---|
| ai unreachable | spend part `null`, rest returned |

## Acceptance criteria

- [ ] With test metrics (Free cost €4.00 over 30 days, 50 active Free households) guardrail 1 shows €0.08, not red.
- [ ] With Prometheus down the route returns 200 and `metrics_unavailable: true`.
