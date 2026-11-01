# AI-04-04 Failed items, retry rules and batch metrics

| | |
|---|---|
| Service | ai |
| Epic | [AI-04 Batch jobs](README.md) |
| Release | R1 |
| Size | S |
| Depends on | AI-04-03 |
| Unblocks | callers' batch pipelines |

## Goal

One documented rule for what callers do with failed batch items, and metrics that show when a batch kind degrades.

## Business rules

1. Caller rule (written into `docs/EPICs/ai/README.md` and the ai service README): resubmit items whose error is
   `retryable=true` in a new batch, at most 2 times per item; after that, put the item in the caller's own review
   queue (admin). Never resubmit `ai_refused`.
2. If more than 20% of a batch's items fail, the caller resubmits nothing automatically and raises a founder alert
   (`SendAlert` via its own notifier edge: KIT, PLN, OFF all have one), because the prompt or input is likely
   wrong.
3. Metrics in ai: `ai_batch_items_total{kind, outcome}` (outcome = `ok` or the error reason),
   `ai_batch_submitted_total{kind, caller}`, `ai_batch_cost_eur_micros_total{kind}`.
4. Grafana: panel "Batch items by outcome" in the AI row (AI-02-02).

## Contract

No change.

## Data

No schema change.

## Errors

None new.

## Acceptance criteria

- [ ] Metrics reflect the 5-item test of AI-04-03.
- [ ] The retry rule is in the ai README "Rules for callers" section.
