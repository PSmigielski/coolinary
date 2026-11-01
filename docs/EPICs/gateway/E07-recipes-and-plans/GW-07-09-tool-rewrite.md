# GW-07-09 Rewrite a dish for the tools at hand

| | |
|---|---|
| Service | gateway |
| Epic | [GW-07 Recipe library, plans and cooking routes](README.md) |
| Release | R2 |
| Size | S |
| Depends on | GW-03-03, PLN-09 (`RewriteDish`) |
| Unblocks | MOB-03 |

## Goal

An oven dish becomes a pan version, air-fryer times, or cooking-robot steps (features.md §2 "Kitchen tools", AI per
request, Plus).

## Business rules

1. Plus-only (`plus: "tool_rewrite"`). Body `{ "tools": [...] }` (1–10 tool codes; empty = the household's tools),
   `ai` deadline, idempotency required (the key is passed as `request_id`).
2. Planner `RewriteDish(household_id, dish_id, kind = TOOLS, target_tools, request_id, language)` returns a
   household-private dish (new id, `parentDishId` = the original), already rule-checked; calories recomputed
   without AI (PLN-09-04). `aiUsed` and `cached` come from planner (a cached rewrite still says `aiUsed: true`).
3. A dish that already fits the tools → planner `INVALID_ARGUMENT` `dish already fits your tools` → 400
   `already_fits`. A rewrite that fails the checks → 422 `rewrite_not_possible`.
4. If AI is unavailable → 503 `ai_unavailable` (there is no non-AI rewrite).
5. Putting the rewritten dish into a plan meal is not part of this call (planner has no such RPC; see GW-07-05
   rule 4).

## Contract

`POST /v1/dishes/{dishId}/rewrite` — `auth: true`, `plus: "tool_rewrite"`, tag `recipes`.

```json
request { "tools": ["hob", "microwave"] }
200     { "dish": { "...detail as GW-07-01": "", "parentDishId": "lasagne" }, "aiUsed": true, "cached": false }
```

## Data

No storage.

## Errors

| Case | HTTP | `code` |
|---|---|---|
| Free | 403 | `plus_required` |
| Dish already fits the tools | 400 | `already_fits` |
| Not possible with these tools | 422 | `rewrite_not_possible` (planner) |
| AI unavailable | 503 | `ai_unavailable` |

## Acceptance criteria

- [ ] Rewriting an oven lasagne for `["hob"]` returns a dish whose tools are only `hob`.
- [ ] A Free token gets 403.
