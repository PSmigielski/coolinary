# GW-08-02 Propose lines from the plan and accept them

| | |
|---|---|
| Service | gateway |
| Epic | [GW-08 Shopping list routes](README.md) |
| Release | R1 |
| Size | M |
| Depends on | GW-08-01, GW-07-05, PLN-15 (`GenerateList`, `AcceptListProposal`) |
| Unblocks | MOB-03 |

## Goal

The list is made from the plan without AI (features.md §2 "How the list is made", steps 1–5): every ingredient of
the chosen days scaled to the pot, minus stock and staples, merged, rounded up to packs; the user reviews and
unticks lines before they join the list (verification 12).

## Business rules

1. `POST /v1/list/proposals` body `{ "from": "YYYY-MM-DD", "to": "YYYY-MM-DD", "meals"?: [...] }` (at most 14 days)
   → planner `GenerateList(household_id, from, to, meals)` with the `heavy` deadline (10 s).
2. Response: planner's `ListProposal`: `proposalId`, `expiresAt`, proposed lines (each with a `key`, ingredient,
   amount, pack size and count, sources, `wanted` pre-set, flags) and `coveredIngredientIds` (fully at home, for the
   "nothing already at home appears" check). Planner keeps a proposal for 24 h (PLN-15-05 rule 2).
3. `POST /v1/list/proposals/{proposalId}/accept` body `{ "wantedKeys": [...] }` — the lines the user kept (unticked
   lines are left out) → `AcceptListProposal(household_id, member_id, proposal_id, wanted_keys,
   client_mutation_id)`. Lines merge into the list (same ingredient → one line). Idempotency required (the key is
   the `client_mutation_id`).
4. Accepting an expired proposal → planner `FAILED_PRECONDITION` `proposal_expired` → 410; the app proposes again.
5. Never calls AI; free for Free and Plus.

## Contract

| Method | Path | 2xx |
|---|---|---|
| POST | `/v1/list/proposals` | 201 `{ "proposalId", "expiresAt", "lines": [proposedLine], "coveredIngredientIds": ["ing_rice"] }` |
| POST | `/v1/list/proposals/{proposalId}/accept` | 200 `{ "version", "lines": [line] }` |

`proposedLine` = `{ "key", "ingredientId", "name", "category", "quantity", "packSize", "packCount", "sources",
"wanted", "flags" }`. Tag `list`, `auth: true`.

## Data

No storage.

## Errors

| Case | HTTP | `code` |
|---|---|---|
| Range > 14 days, no plans in range | 400 / 422 | `validation` / `no_plans` |
| Proposal expired | 410 | `proposal_expired` |
| Unknown key in accept | 422 | planner's code |

## Acceptance criteria

- [ ] For the fixed test plan and fridge, the proposal holds exactly the missing ingredients, each once, amounts
      merged across dishes, and nothing already at home (verification 12).
- [ ] Unticked lines are not added on accept.
- [ ] Accepting twice with the same idempotency key adds the lines once.
