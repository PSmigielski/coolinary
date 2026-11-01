# AI-02-05 Per-household monthly cost and cost per active household

| | |
|---|---|
| Service | ai |
| Epic | [AI-02 Cost tracking and spend limit](README.md) |
| Release | R1 |
| Size | M |
| Depends on | AI-02-04 |
| Unblocks | ADM-07 |

## Goal

Track the business-model guardrails directly: what an average Free and an average Plus household cost each week
and month (business model §6: "track cost per active user weekly"), and stop a single Free household from going far
past its cap.

## Business rules

1. Per household and calendar month (UTC), keep `ai:hh:<household_hash>:<YYYY-MM>` = micro-euros spent (sync calls
   only; batch jobs are for everyone). The key uses the salted hash from `tracing.household_hash`, never the raw
   id. Expires 40 days after month end.
2. Active households: a HyperLogLog `ai:active:<plan>:<ISO-week>` gets the household hash on every live or cached
   AI call. Weekly spend per plan: `ai:spend:week:<plan>:<ISO-week>`.
3. A scheduled task in the service (every 15 minutes) sets gauges:
   `ai_cost_per_active_household_eur{plan, period="week"}` = weekly spend ÷ HLL count, and the same for the current
   month.
4. Free household guard: a `free` call whose household has already spent more than `AI_FREE_HOUSEHOLD_MONTHLY_EUR`
   (default 0.15, i.e. the €0.10 cap plus 50% margin) this month is refused with `RESOURCE_EXHAUSTED
   ai_household_budget`. The gateway's weekly quota (GW-03) is the normal limit; this guard only catches what the
   quota doesn't (e.g. repeated swaps if they were ever opened to Free, or very large inputs).
5. Plus households are never refused by this guard; their cost shows in the gauge and in Langfuse by hashed user.
6. Admin overview (ADM-07) reads the gauges through Prometheus; no new RPC.

## Contract

No gRPC change. RPCs can return `RESOURCE_EXHAUSTED ai_household_budget`.

## Data

No schema change. Redis keys as above (same instance as AI-02-04).

## Errors

| Case | gRPC status |
|---|---|
| Free household over its monthly guard | `RESOURCE_EXHAUSTED`, `ai_household_budget` |

## Acceptance criteria

- [ ] Three Free households with costs €0.02, €0.04, €0.06 in one week give a gauge of €0.04 for `plan="free"`.
- [ ] A Free household at €0.16 this month gets `ai_household_budget`; a Plus household at €5 does not.
- [ ] No Redis key or metric contains a raw household id.

## Notes

- Decided (Q019): a hard guard of €0.15 per Free household per month. Once reached, that household's per-request AI
  calls return `RESOURCE_EXHAUSTED` until the next calendar month.
