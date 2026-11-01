# PLN-20-02 Cheapest top-up with AI

| | |
|---|---|
| Service | planner |
| Epic | [PLN-20 Cheapest top-up](README.md) |
| Release | R3 |
| Size | L |
| Status | todo |
| Depends on | PLN-20-01, AI-17 (CheapestTopUp), PLN-19-02 |
| Unblocks | PLN-20-04, GW-10 |

## Goal

Ask AI for the smallest, cheapest basket that completes the week, using offers and allowed swaps; check and price it
without AI.

## Business rules

1. Plus only (`plus_required`). Requires an active week plan for the current or next week (`no_week_plan`
   otherwise).
2. Request to `AI-17 (CheapestTopUp)`: needs with amounts, allowed swaps (PLN-20-01, nothing else), candidate offers
   with shop positions and prices. No household names or personal data.
3. AI returns a basket: per need either the original ingredient or one allowed swap, and an offer id or "no offer".
   Any swap not in the allowed set, or an offer id not in the candidates, is dropped and the original is used.
4. The basket is priced without AI with PLN-19-02's rules at the shop of each chosen offer (no-offer lines at the
   household's first shop, reference price).
5. The basket is compared with the no-swap basket (PLN-20-03); if it is not cheaper, the no-swap basket is returned
   (`ai_basket_rejected`).
6. Result: lines (ingredient, amount, shop, offer, cost, source), swaps (from, to, dishes, saving), total, saving
   against the no-swap basket. Labelled "about".
7. The result is stored for 24 hours so PLN-20-04 can apply it.
8. `ai` down, over budget or refusing → PLN-20-03 result with `ai_used = false`.

## Contract

`contracts/coolinary/planner/v1/topup.proto`

```proto
rpc CreateCheapestTopUp(CreateCheapestTopUpRequest) returns (TopUp);
message CreateCheapestTopUpRequest { string household_id = 1; string week_plan_id = 2; repeated string shop_ids = 3; bool use_ai = 4; }
message TopUp {
  string id = 1; string expires_at = 2;
  repeated TopUpLine lines = 3;
  repeated TopUpSwap swaps = 4;
  coolinary.common.v1.Money total = 5;
  coolinary.common.v1.Money saving = 6;
  bool ai_used = 7; string fallback_reason = 8;
}
message TopUpLine { string ingredient_id = 1; coolinary.common.v1.Quantity amount = 2; string shop_id = 3; string offer_id = 4; coolinary.common.v1.Money cost = 5; PriceSource source = 6; }
message TopUpSwap { string key = 1; string from_ingredient_id = 2; string to_ingredient_id = 3; repeated string dish_ids = 4; coolinary.common.v1.Money saving = 5; int32 kcal_change_percent = 6; }
```

`shop_ids` are the shops the phone allows (nearest first), validated as in PLN-19-01.

## Data

`planner/migrations/000NN_topups.sql`

```sql
CREATE TABLE topups (
    id           uuid PRIMARY KEY,
    household_id uuid NOT NULL,
    week_plan_id uuid NOT NULL,
    payload      jsonb NOT NULL,
    created_at   timestamptz NOT NULL DEFAULT now(),
    expires_at   timestamptz NOT NULL
);
```

## Errors

Reason codes (in backticks) travel as trailing gRPC metadata `x-error-code` (README decision 18); version
conflicts use `ABORTED`.

| Case | gRPC status |
|---|---|
| Free household | `FAILED_PRECONDITION` `plus_required` |
| No active week plan | `FAILED_PRECONDITION` `no_week_plan` |
| offers down | `UNAVAILABLE` |

## Acceptance criteria

- [ ] An AI swap outside the allowed set never reaches the response.
- [ ] Every price in the response equals PLN-19-02's price for that shop and line.
- [ ] With AI off, a basket without swaps is returned.
