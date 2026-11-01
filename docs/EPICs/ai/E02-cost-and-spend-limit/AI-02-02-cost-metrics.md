# AI-02-02 Cost and usage metrics, Grafana panels

| | |
|---|---|
| Service | ai |
| Epic | [AI-02 Cost tracking and spend limit](README.md) |
| Release | R1 |
| Size | S |
| Depends on | AI-02-01 |
| Unblocks | ADM-07 |

## Goal

See AI spend per job, plan and model in Grafana every day, so a cost problem shows up in hours, not at the end of
the month.

## Business rules

1. Prometheus metrics (no household or user label, to keep cardinality low):
   - `ai_cost_eur_micros_total{job, plan, model, mode="sync"|"batch"}` (counter)
   - `ai_tokens_total{job, model, kind="input"|"output"|"cache_read"|"cache_write"}` (counter)
   - `ai_requests_total{job, plan, outcome="ok"|"cached"|"refused"|"invalid"|"error"|"budget"}` (counter)
   - `ai_request_duration_seconds{job, model}` (histogram, buckets 0.5 s to 120 s)
2. `plan` label values: `free`, `plus`, `none` (batch jobs, which are for everyone).
3. Grafana: add a row "AI" to `platform/grafana/dashboards/services.json` with: € per day by job; € per day by plan;
   requests and outcomes by job; cache-read share by job (AI-01-05); p95 duration by job.
4. Prometheus alert rules (in `platform/prometheus`): spend today above 80% of the daily limit (AI-02-04);
   refusal share of a job above 5% over 1 hour; `outcome="error"` share above 10% over 15 minutes.

## Contract

No gRPC change.

## Data

No schema change.

## Errors

None new.

## Acceptance criteria

- [ ] After one stub call and one live call, `/metrics` shows the counters with the right labels; the stub adds 0 cost.
- [ ] The dashboard JSON loads in Grafana 13 without errors and shows the five panels.
- [ ] No metric carries a household id, user id or text.
