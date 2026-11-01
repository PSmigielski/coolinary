# AI-10-05 Model routing by plan and the cost decision

| | |
|---|---|
| Service | ai |
| Epic | [AI-10 Day plan composition and cuisine choices](README.md) |
| Release | R1 |
| Size | S |
| Depends on | AI-10-06, AI-02-05 |
| Unblocks | day plans going live |

## Goal

Route composition jobs by plan as decision 15 sets out: Free day plans on Haiku 5.5, Plus on Sonnet 5.5 only where
measured quality justifies it. Record the measured numbers next to the business-model cost caps.

## Business rules

1. The registry's `by_plan` overrides (AI-01-01 rule 3) pick model and effort per plan for every composition job
   (`compose_day_plan`, `modify_meal`, `rewrite_dish`, `compose_family_day`, `compose_from_deals`,
   `compose_week_plan`, `cheapest_top_up`). The plan comes from `x-plan` (AI-01-08).
2. Run the day-plan golden set (AI-10-06) on Haiku 5.5 at `medium` and `high`, and Sonnet 5.5 at `low` and
   `medium`. Record quality metrics, cost per call and p95 latency for each in one report.
3. Compute, from the measured cost per call:
   - Free monthly cost = 8.7 day plans + 8.7 photo scans; must be ≤ €0.10 (business model §6).
   - Plus monthly cost for a household planning 20 days a month plus 10 swaps; compare with ≈ €0.25 (15% of €19.7
     net a year ÷ 12).
4. Free: `compose_day_plan` runs on `claude-haiku-5-5` for `x-plan=free` (and for an empty or unknown plan). This
   is fixed by decision 15; the report in rule 2 only sets its effort level.
5. Plus: every composition job starts on Haiku 5.5 and moves to Sonnet 5.5 only if the report shows Sonnet clearly
   better (Notes). The other composition jobs are Plus-only, so their Free routing is never used.
6. Write the chosen routing into `jobs.toml` and the numbers into `docs/EPICs/ai/README.md` (job table).

## Contract

No change.

## Data

No schema change.

## Errors

None new.

## Acceptance criteria

- [ ] The comparison report is committed under `ai/evals/compose_day_plan/routing/`.
- [ ] A `compose_day_plan` call with `x-plan: free` (or no `x-plan`) runs on `claude-haiku-5-5`, checked in the
      Langfuse generation's model field and the `ai_cost` metric labels.
- [ ] The Free routing meets the €0.10 cap on measured cost; the Plus figure is written down next to the 15% cap.

## Notes

- Decided (README decision 15, 1 November 2026): **Free day plans compose on Haiku 5.5.** Plus compositions use
  Sonnet 5.5 only if it clearly beats Haiku on the golden set (AI-10-06); otherwise Haiku too. Why: the AI-02
  estimate (≈ $0.05 per Sonnet day plan) puts a Free household at ≈ $0.44 a month against the €0.10 cap, while
  Haiku costs ≈ $0.0025 a plan (≈ $0.02 a month). A daily Plus planner on Sonnet is ≈ $1.50 a month against the
  ≈ €0.25 cap, so the Plus choice is made on measured numbers from rule 3.
- "Clearly beats" means: on the AI-10-06 golden set, Sonnet passes validation first time on at least 10 percentage
  points more plans than Haiku, or Haiku fails a hard threshold of AI-10-06 that Sonnet meets.
  Decided (Q123): the margin is 10 percentage points.
- Decided (Q123): Sonnet may act as the repair model only for Plus, when Haiku's Plus output fails validation; Free
  repairs stay on Haiku.
