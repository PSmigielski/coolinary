# AI-11-04 RewriteDish: TOOLS

| | |
|---|---|
| Service | ai |
| Epic | [AI-11 Swaps, wishes, tool rewrite, restyle](README.md) |
| Release | R2 |
| Size | M |
| Depends on | AI-11-01 |
| Unblocks | PLN-09 (tool rewrite) |

## Goal

Rewrite a recipe for the tools at hand: an oven dish as a pan version, air-fryer times, cooking-robot steps
(features.md §2 "Kitchen tools", with AI, per request; Plus).

## Business rules

1. Job `rewrite_dish`, Sonnet 5.5 at effort `low`. Prompt `rewrite_dish/v1`.
2. Input: `CompositionContext` (tools, rules, facts for the dish's ingredients; stock optional), the dish (library
   dish as a full `NewDish` copy, or a new dish), `mode=TOOLS`, `target_tools` (the tools to use, ⊆
   `rules.tools`; empty = any the kitchen has), `avoid_tools` (e.g. `oven`).
3. Output: a `NewDish` with the same ingredients (ids and amounts; small changes allowed only for fat or liquid,
   max ±30% per line, flagged in `changes`), new steps and times, `tools` ⊆ target or kitchen tools and
   disjoint from `avoid_tools`; `changes` lists what changed (`TOOLS`, `TIME`, `AMOUNT`, `STEP_ORDER`).
4. If no version is possible with those tools (e.g. bread without oven or pan), return `not_possible=true` with
   reason `TOOLS` — not an error.
5. Validation: AI-10-03 rules 3–8, 11–13 on the new dish; ingredient set equal to the original's (same ids); amount
   changes within ±30% and only on ingredients of category fat or liquid (from facts).
6. Cache: 30 days (a rewrite of the same dish for the same tools is the same for everyone); key without household
   id.

## Contract

```proto
rpc RewriteDish(RewriteDishRequest) returns (RewriteDishResponse);

enum RewriteMode { REWRITE_MODE_UNSPECIFIED = 0; REWRITE_MODE_TOOLS = 1; REWRITE_MODE_RESTYLE = 2; }
enum RewriteChange { REWRITE_CHANGE_UNSPECIFIED = 0; REWRITE_CHANGE_TOOLS = 1; REWRITE_CHANGE_TIME = 2;
  REWRITE_CHANGE_AMOUNT = 3; REWRITE_CHANGE_STEP_ORDER = 4; REWRITE_CHANGE_INGREDIENT = 5; REWRITE_CHANGE_CUISINE = 6; }

message RewriteDishRequest {
  CompositionContext context = 1;
  NewDish dish = 2;
  string source_dish_id = 3;          // library id when rewriting a library dish
  RewriteMode mode = 4;
  repeated coolinary.common.v1.KitchenTool target_tools = 5;
  repeated coolinary.common.v1.KitchenTool avoid_tools = 6;
  coolinary.common.v1.Cuisine target_cuisine = 7;          // RESTYLE (AI-11-05)
  bool bypass_cache = 8;
}

message RewriteDishResponse {
  NewDish dish = 1;
  repeated RewriteChange changes = 2;
  bool not_possible = 3;
  string not_possible_reason = 4;     // TOOLS | CUISINE | RULES
  CallInfo call = 5;
}
```

## Data

No schema change.

## Errors

| Case | gRPC status |
|---|---|
| Input limits broken, unknown ids or codes in the request | `INVALID_ARGUMENT`, `ai_bad_request` |
| Output still invalid after repair | `UNAVAILABLE`, `ai_invalid_output` |
| Others | AI-01-04; the caller falls back to its non-AI result |
| No version possible with those tools | OK, `not_possible=true`, reason `TOOLS` |

## Acceptance criteria

- [ ] An oven-baked chicken dish rewritten with `avoid_tools=[oven]` for a hob+pan kitchen uses only hob/pan and keeps all ingredients.
- [ ] An air-fryer rewrite gives air-fryer times and temperature in the steps.
- [ ] A rewrite that adds an ingredient fails validation.
