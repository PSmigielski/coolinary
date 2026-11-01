# AI-15-03 Validation of sessions and leftovers

| | |
|---|---|
| Service | ai |
| Epic | [AI-15 Week plan](README.md) |
| Release | R3 |
| Size | M |
| Status | todo |
| Depends on | AI-15-02, AI-12-03 |
| Unblocks | AI-15-04 |

## Goal

A week plan that promises "cook once, eat twice" must be safe and consistent: leftovers come from a real earlier
batch, within its keeping time.

## Business rules

1. Every `COOK` and `NO_COOK` slot passes AI-10-03 (and AI-12-03 when eaters are present).
2. `SLOTS` — exactly one slot per requested day and meal.
3. `RHYTHM` — number of dates with a `COOK` slot that heats food: `EVERY_OTHER_DAY` ≤ ceil(days ÷ 2);
   `TWICE_A_WEEK` ≤ 2; `EVERY_DAY` ≤ days.
4. `LEFTOVER_REF` — `leftover_of` points to an earlier `COOK` slot with `double_batch=true`; each double batch feeds
   exactly one leftover slot.
5. `KEEPING` — days between cook and leftover ≤ the dish's keeping days (3 for `KEEPS_3_DAYS`, else `keeps_days`, else
   1); more only with `freeze=true` and the dish `FREEZES_WELL` / `freezes`.
6. `REPEATS` — a dish appears (cook + leftover) at most twice in the week; `new dishes` ≤ `max_new_dishes`.
7. `STOCK_USE` over the whole week: per stock line, sum ≤ available; `to_buy` lines ≤ `max_items_to_buy`.

## Contract

No change.

## Data

No schema change.

## Errors

| Case | gRPC status |
|---|---|
| Still invalid after the one repair | `UNAVAILABLE`, `ai_invalid_output` |

## Acceptance criteria

- [ ] One fake bad week per rule is caught with its rule code.
- [ ] A chili cooked Monday and eaten Thursday without `freeze` fails `KEEPING`; with `freeze=true` and `FREEZES_WELL` it passes.
