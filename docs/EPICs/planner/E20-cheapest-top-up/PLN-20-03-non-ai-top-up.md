# PLN-20-03 Top-up without AI

| | |
|---|---|
| Service | planner |
| Epic | [PLN-20 Cheapest top-up](README.md) |
| Release | R3 |
| Size | S |
| Status | todo |
| Depends on | PLN-20-01, PLN-19-02 |
| Unblocks | PLN-20-02 |

## Goal

The no-swap basket: every need priced at its cheapest allowed shop. It is the baseline for PLN-20-02 and the
fallback.

## Business rules

1. For each need, price it at every allowed shop with PLN-19-02 and take the cheapest cell; ties → earlier shop.
2. No swaps.
3. Same response shape as PLN-20-02, `ai_used = false`, `saving = 0`.
4. The basket may span more shops than the household's "shops per trip" setting; the app shows it per shop and the
   phone's trip plan decides the final set.

## Contract

Used by `CreateCheapestTopUp` (PLN-20-02).

## Data

No schema change.

## Errors

None of its own; the baseline and fallback inside `CreateCheapestTopUp` (PLN-20-02), whose errors apply.

## Acceptance criteria

- [ ] Worked example prices: onions, soy sauce and ginger go to Kaufland, chicken and rice to Lidl; total 9.15.
