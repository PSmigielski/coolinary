# AI-02 Cost tracking and spend limit

| | |
|---|---|
| Service | ai |
| Release | R1 |
| Goal | Know what every AI call costs, per job and plan, and stop spending before the business-model caps are broken: a Free household ≤ €0.10 a month, AI ≤ 15% of a payer's net revenue (≈ €0.25 a month at the €19.7 net annual price). |

## Flow

```
completion usage -> price table (AI-02-01) -> cost in EUR micros
   -> CallInfo.cost_eur_micros in the response
   -> Prometheus counters by job / plan / model (AI-02-02)
   -> Langfuse generation cost details (AI-02-03)
   -> daily spend counters + limit check before the next call (AI-02-04)
   -> per-household monthly cost and cost per active household (AI-02-05)
```

## Cost estimates per call (to be replaced by measured usage)

Assumptions: prompt caching hits on the system prefix, effort as in AI-01-01, outputs at the schema's typical size.

| Job | Model | Input tok. | Output tok. (incl. thinking) | ≈ USD per call |
|---|---|---|---|---|
| text_to_items | Haiku 5.5 | 1,200 | 400 | 0.0003 |
| photo_to_items (2 photos) | Haiku 5.5 | 4,500 | 1,200 | 0.001 |
| read_label | Haiku 5.5 | 3,000 | 300 | 0.0005 |
| compose_day_plan | Sonnet 5.5 | 7,000 | 3,500 | 0.05 |
| compose_day_plan | Haiku 5.5 | 7,000 | 3,500 | 0.0025 |
| modify_meal | Sonnet 5.5 | 6,000 | 1,500 | 0.027 |
| compose_week_plan | Sonnet 5.5 | 10,000 | 10,000 | 0.12 |
| flyer page (batch) | Haiku 5.5 | 2,500 | 1,500 | 0.0005 |

**What this means for the caps.** A Free household may run 2 day plans and 2 photo scans a week (≈ 8.7 of each a
month). With Sonnet day plans that is ≈ $0.44 a month, about four times the €0.10 cap; with Haiku day plans it is
≈ $0.03. A Plus household planning every day with Sonnet costs ≈ $1.50 a month, about six times the 15% cap. Decision 15
and Q123 (AI-10-05) settle it: Free composes on Haiku, Plus uses Sonnet only if it wins on the golden set; the spend
limits below protect the budget either way.

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [AI-02-01](AI-02-01-price-table-and-cost.md) | Price table and cost calculator | R1 | S | AI-01-01 |
| [AI-02-02](AI-02-02-cost-metrics.md) | Cost and usage metrics, Grafana panels | R1 | S | AI-02-01 |
| [AI-02-03](AI-02-03-langfuse-fields.md) | Langfuse fields for every job | R1 | S | AI-02-01 |
| [AI-02-04](AI-02-04-daily-spend-limit.md) | Global daily spend limit | R1 | M | AI-02-01 |
| [AI-02-05](AI-02-05-household-cost.md) | Per-household monthly cost and cost per active household | R1 | M | AI-02-04 |
