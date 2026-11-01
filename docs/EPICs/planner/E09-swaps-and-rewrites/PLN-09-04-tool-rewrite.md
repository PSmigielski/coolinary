# PLN-09-04 Rewrite a recipe for my tools

| | |
|---|---|
| Service | planner |
| Epic | [PLN-09 Swaps, wishes, tool rewrite, restyle](README.md) |
| Release | R2 |
| Size | M |
| Status | todo |
| Depends on | PLN-07-05, PLN-02-04, PLN-04-01, AI-11 (RewriteDish) |
| Unblocks | GW-07 (`POST /v1/dishes/{id}/rewrite`) |

## Goal

A library dish needs an oven the household doesn't have: rewrite it for the tools at hand — an oven dish as a pan
version, air-fryer times, cooking-robot steps (features.md "Kitchen tools").

## Business rules

1. `RewriteDish(dish_id, kind=TOOLS, target_tools?)` for any dish the household can read; default target tools = the
   household's tools; `target_tools` narrows them ("air fryer only").
2. Plus only (Free → FAILED_PRECONDITION `plus_required`).
3. If the dish already fits the target tools → INVALID_ARGUMENT `dish already fits your tools`.
4. Planner sends the dish (ingredient ids, amounts, steps, times, tools), target tools and language. AI may change
   steps, timers, times and tools, and amounts of fat or liquid by up to ±50%; it may not add or remove other
   ingredients. Checked: tools ⊆ target (PLN-07-02), ingredient set unchanged except fat/liquid amounts, structure.
5. The result is a household AI dish with `parent_dish_id` = original, `source=household_ai`; nutrition recomputed.
6. The ai service's result cache makes the same rewrite (same dish version, tool set, language) free for the next
   household; planner still stores one copy per household.
7. Failing answer → FAILED_PRECONDITION `rewrite not possible`.
8. The response says `ai_used = true` when the dish came from AI and `cached = true` when ai's result cache
   answered (rule 6), so the gateway counts Plus fair use (GW-03-04) and AI metrics (GW-03-05).

## Contract

```proto
rpc RewriteDish(RewriteDishRequest) returns (RewriteDishResponse);
message RewriteDishRequest {
  string household_id = 1;
  string dish_id = 2;
  RewriteKind kind = 3;
  repeated coolinary.common.v1.KitchenTool target_tools = 4;    // TOOLS
  coolinary.common.v1.Cuisine target_cuisine = 5;                // RESTYLE, PLN-09-05
  string request_id = 6;
  string language = 7;
}
message RewriteDishResponse { Dish dish = 1; bool ai_used = 3; bool cached = 4; }
enum RewriteKind { REWRITE_KIND_UNSPECIFIED = 0; REWRITE_KIND_TOOLS = 1; REWRITE_KIND_RESTYLE = 2; }
```

## Data

No schema change (`dishes.parent_dish_id`, PLN-02-01).

## Errors

Reason codes (in backticks) travel as trailing gRPC metadata `x-error-code` (README decision 18); version
conflicts use `ABORTED`.

| Case | gRPC status | Message |
|---|---|---|
| Free | FAILED_PRECONDITION | `plus_required` |
| Already fits | INVALID_ARGUMENT | `dish already fits your tools` |
| Fails checks | FAILED_PRECONDITION | `rewrite not possible` |

## Acceptance criteria

- [ ] An oven lasagne rewritten for hob + microwave has no oven in tools or step texts (PLN-07-02 word list).
- [ ] The rewrite keeps the same ingredient ids; only oil or water amounts may differ, by ≤ 50%.
- [ ] The rewritten dish shows its own calories, computed by kitchen.
- [ ] A rewrite answered from ai's result cache (rule 6) returns `ai_used = true, cached = true`.
