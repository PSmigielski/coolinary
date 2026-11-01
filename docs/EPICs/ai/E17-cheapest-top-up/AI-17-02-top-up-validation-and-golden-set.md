# AI-17-02 Top-up validation and golden set

| | |
|---|---|
| Service | ai |
| Epic | [AI-17 Cheapest top-up](README.md) |
| Release | R3 |
| Size | M |
| Depends on | AI-17-01, AI-03-05 |
| Unblocks | cheapest top-up going live |

## Goal

Every choice refers to a real offer, swaps never break a dish's rules, and the AI basket is cheaper than the plain one
often enough to be worth its cost.

## Business rules

1. Validation:
   - `CHOICES` — exactly one choice per need.
   - `OFFER_REF` — `offer_ref` is a request deal; for `OFFER` its ingredient equals the need's; for `SWAP` it equals
     `swap_ingredient_id` and differs from the need's.
   - `SWAP_LIMITS` — the swap ingredient keeps every hard limit of the household rules (facts); `changed_dish_refs`
     ⊆ the need's `dish_refs`.
   - `SHOPS` — at most 2 distinct `shop_ref`s across all choices.
   - `PACKS` — 1–20 packs.
   - Text: calorie guard; no currency amounts in `note` (regex for digits followed by a currency sign or code).
2. Golden set: 20 week baskets (pl/de) with 30–80 deals each. The scorer prices the AI basket and the plain basket
   (cheapest same-ingredient offer per line, or reference price) with the deal prices.
3. Thresholds: violations 0; AI basket cheaper or equal in ≥ 0.90 of cases; median saving vs plain basket recorded in
   the report (no threshold); swaps accepted by reviewer ≥ 0.80.

## Contract

No change.

## Data

`ai/evals/cheapest_top_up/`.

## Errors

| Case | gRPC status |
|---|---|
| Still invalid after the one repair | `UNAVAILABLE`, `ai_invalid_output` |

## Acceptance criteria

- [ ] One fake bad answer per rule is caught.
- [ ] Eval passes; baseline committed.
