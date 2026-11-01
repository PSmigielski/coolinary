# AI-17-01 CheapestTopUp

| | |
|---|---|
| Service | ai |
| Epic | [AI-17 Cheapest top-up](README.md) |
| Release | R3 |
| Size | L |
| Status | todo |
| Depends on | AI-14-01 |
| Unblocks | AI-17-02; PLN-20 |

## Goal

Choose, for each missing line of the week plan, the offer to buy or a discounted equivalent that keeps the dishes
working, so the basket is as small and cheap as the offers allow.

## Business rules

1. Job `cheapest_top_up`, Sonnet 5.5 at effort `medium`. Prompt `cheapest_top_up/v1`.
2. Input: `CompositionContext` (rules, facts; stock not needed beyond what planner already subtracted), 1–80 `Need`s
   (`line_ref`, `ingredient_id`, `amount`, the dishes it is for with their other ingredient ids), and 1–200 `Deal`s
   (AI-14-01) for the needed ingredients and for same-category alternatives.
3. Per need, one `TopUpChoice`:
   - `OFFER` — buy `packs` of a deal with the same `ingredient_id`.
   - `SWAP` — buy a deal with another `ingredient_id` that works in every dish the need is for (same culinary role),
     with `note` ≤ 80 chars ("use turkey instead of chicken"), and the dishes it changes.
   - `NO_OFFER` — no suitable deal; planner prices it at the receipt or reference price.
4. Rules: hard limits of every dish's eaters hold for swaps; fewest shops (prefer one shop, never more than 2 for the
   whole basket); fewest packs that cover the amount (planner rounds to packs; AI gives its estimate); prefer lower
   `shop_rank` when prices are about equal; a swap must also be a real saving, which planner checks by prices.
5. AI states no prices, totals or savings in text.
6. Validation (AI-17-02). Cache 1 hour.

## Contract

```proto
rpc CheapestTopUp(CheapestTopUpRequest) returns (CheapestTopUpResponse);

message Need {
  string line_ref = 1;
  string ingredient_id = 2;
  coolinary.common.v1.Quantity amount = 3;
  repeated string dish_refs = 4;
  repeated string dish_ingredient_ids = 5;  // other ingredients of those dishes, for context
}

enum TopUpKind { TOP_UP_KIND_UNSPECIFIED = 0; TOP_UP_KIND_OFFER = 1; TOP_UP_KIND_SWAP = 2; TOP_UP_KIND_NO_OFFER = 3; }

message TopUpChoice {
  string line_ref = 1;
  TopUpKind kind = 2;
  string offer_ref = 3;
  int32 packs = 4;
  string swap_ingredient_id = 5;
  repeated string changed_dish_refs = 6;
  string note = 7;
}

message CheapestTopUpRequest {
  CompositionContext context = 1;
  repeated Need needs = 2;      // 1-80
  repeated Deal deals = 3;      // 1-200
  bool bypass_cache = 4;
}

message CheapestTopUpResponse {
  repeated TopUpChoice choices = 1;
  CallInfo call = 2;
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

## Acceptance criteria

- [ ] With chicken thighs on offer and chicken breast needed for a curry, the result may be a `SWAP` with a note; for a dish where the swap breaks a diet it is never offered.
- [ ] The basket uses at most 2 shops.
