# PLN-10-02 "Cooked it": consume stock, history, plan status

| | |
|---|---|
| Service | planner |
| Epic | [PLN-10 Cooking mode and "cooked it"](README.md) |
| Release | R1 |
| Size | M |
| Depends on | PLN-05-06, PLN-08-03, PLN-02-05, KIT-07 (ConsumeStock) |
| Unblocks | PLN-10-03, PLN-11-01, KIT-12 (used-before-expiry savings, via kitchen) |

## Goal

One tap after cooking takes the used ingredients out of the fridge list, oldest first, records the dish in the history
and marks the plan meal cooked — and is safe to send twice from the offline queue (MOB-04).

## Business rules

1. `MarkCooked(cooked_id, dish_id, version, servings, plan_meal_id?, chosen_toppings?, left_out?)`. `cooked_id` is a
   UUID made on the phone; it is the idempotency key end to end.
2. If a history entry with that `cooked_id` exists, return it unchanged (`replayed=true`) and call nothing.
3. Lines to consume: every non-optional `main`/`base` line, plus for each topping group the options chosen
   (`chosen_toppings`; default: the first option), plus optional lines the user ticked. Lines in `left_out` are skipped.
   Staples are never consumed (they are not tracked as stock, features.md §1).
4. Amounts = exact scaled amounts (PLN-02-05 rule 5) for `servings`.
5. Planner calls `ConsumeStock(household_id, lines, reason=COOKED, idempotency_key=cooked_id, dish_id)`; for plan meals
   it passes the stock item ids from the plan allocation (PLN-08-03) as preferred items, so the oldest items chosen at
   planning time go first. Kitchen subtracts and returns what it could not find (shortfalls).
6. Shortfalls are not errors: the user may have used something not on the list. They are returned so the app can say
   "2 items weren't in your fridge list".
7. After a successful kitchen call, in one transaction: write `cooked_history` (PLN-05-06) with `id = cooked_id`, set
   the plan meal to `cooked`, and delete the stock snapshot cache (PLN-01-04 rule 3).
8. Kitchen unavailable → UNAVAILABLE and nothing written; the phone's offline queue retries with the same `cooked_id`.
   Because kitchen is idempotent on the key, a retry after a timeout never consumes twice.
9. `cooked_at` = the time sent by the phone (it may have been queued offline), not later than server time + 5 min and
   not earlier than 7 days ago; else server time.
10. Free and Plus alike.

## Contract

```proto
rpc MarkCooked(MarkCookedRequest) returns (MarkCookedResponse);
message MarkCookedRequest {
  string household_id = 1;
  string user_id = 2;
  string cooked_id = 3;
  string dish_id = 4;
  int32 version = 5;
  double servings = 6;
  string plan_meal_id = 7;
  map<string, StringList> chosen_toppings = 8;    // topping_group -> ingredient ids
  repeated string left_out_ingredient_ids = 9;
  repeated string optional_used_ingredient_ids = 10;
  string cooked_at = 11;                          // RFC 3339, from the phone
}
message Shortfall { string ingredient_id = 1; string name = 2; coolinary.common.v1.Quantity missing = 3; }
message MarkCookedResponse { CookedEntry entry = 1; repeated Shortfall shortfalls = 2; bool replayed = 3; bool rating_prompt = 4; }
```

Decided (Q126): KIT-07 `ConsumeStock` accepts an idempotency key and preferred stock item ids (both needed for
offline retries and "oldest first" from the plan).

## Data

Uses `cooked_history` (PLN-05-06) and `plan_meals.status` (PLN-08-01).

## Errors

| Case | gRPC status | Message |
|---|---|---|
| Kitchen unavailable | UNAVAILABLE | nothing written; retry with the same `cooked_id` |
| Plan meal of another household | NOT_FOUND | `plan meal not found` |
| Bad servings | INVALID_ARGUMENT | from PLN-02-05 |

## Acceptance criteria

- [ ] Cooking the curry for 4 with 600 g chicken in stock leaves 0 g chicken and 1 history entry.
- [ ] Sending the same request twice (including after a simulated timeout between kitchen and planner) consumes once.
- [ ] Salt and oil (staples) are not sent to `ConsumeStock`.
- [ ] A left-out topping is not consumed.
- [ ] Cooking with AI off works unchanged (features.md rule 1).
