# AI-14-02 Deals validation and golden set

| | |
|---|---|
| Service | ai |
| Epic | [AI-14 Plan from deals](README.md) |
| Release | R2 |
| Size | M |
| Depends on | AI-14-01, AI-03-05 |
| Unblocks | plan from deals going live |

## Goal

Deal plans keep the fridge-first promise and use only real, valid offers.

## Business rules

1. Validation rules added to AI-10-03 for this job:
   - `DEAL_REF` — every `offer_ref` is one of the request's deals; its `ingredient_id` equals the `to_buy` line's.
   - `DEAL_VALID` — `valid_to ≥ context.date`.
   - `DEAL_NEEDED` — no deal is bought for an ingredient whose stock line covers the dish's amount (same unit); with
     different units the line is allowed (planner checks with conversions).
   - `DEAL_LIMITS` — ≤ 8 deal lines a day; ≤ 2 distinct `shop_ref`s.
2. Golden set: 25 cases (pl/de), each with a fridge, 20–60 deals across 3–6 ranked shops.
3. Thresholds: violations 0; expiring-stock use ≥ 0.80; "bought though in stock" = 0; share of deal lines from the two
   nearest shops ≥ 0.70 where equal options exist; one-shop days ≥ 0.60.

## Contract

No change.

## Data

`ai/evals/compose_from_deals/`.

## Errors

| Case | gRPC status |
|---|---|
| Still invalid after the one repair | `UNAVAILABLE`, `ai_invalid_output` |

## Acceptance criteria

- [ ] A fake answer citing an unknown `offer_ref` fails `DEAL_REF`; one with three shops fails `DEAL_LIMITS`.
- [ ] Eval passes; baseline committed.
