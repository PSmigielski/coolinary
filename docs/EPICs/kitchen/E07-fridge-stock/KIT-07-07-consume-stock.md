# KIT-07-07 ConsumeStock on "cooked it"

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-07 Fridge stock](README.md) |
| Release | R1 |
| Size | L |
| Depends on | KIT-07-01, KIT-02 (unit conversion library, KIT-02-02), KIT-10-05 (staples) |
| Unblocks | PLN-10 ("cooked it"), KIT-12-02 |

## Goal

When the user taps "cooked it", the dish's ingredients leave the fridge on their own, oldest items first, so stock
stays true without typing (features.md §2: "Stock is reduced automatically when a dish is marked cooked").

## Business rules

1. Planner sends the ingredient lines already scaled to the pot that was cooked (PLN-10-02), 1–60 lines, with a
   `cook_event_id` (planner's `cooked_id`, a UUID made on the phone, the idempotency key end to end).
2. Lines with the same ingredient are summed first (converted to one unit with KIT-02-02 `Convert`; if they cannot be converted,
   they are handled one after the other).
3. **Which items are used**, for each ingredient:
   1. candidates are the household's matched items of that ingredient;
   2. items named in the line's `preferred_item_ids` (the items planner allocated when it made the plan, PLN-08-03)
      go first, in the order given, if they still exist and hold that ingredient; unknown or emptied ids are
      ignored;
   3. then fridge and pantry items, by use-by date ascending (items without a date after dated ones), then by
      `added_at`;
   4. freezer items only after all others. Decided (Q080): frozen food is assumed not to be cooked unless nothing
      else is left; within fridge and pantry the oldest date goes first.
4. **Units:** the line is converted to each item's unit with KIT-02-02 `Convert` (g ↔ ml by density, pieces ↔ g by typical
   weight; tablespoons and other measures through KIT-02-02). Items that cannot be converted are skipped; if no
   item can be converted, the line's outcome is `unconvertible` and nothing is taken. Lines KIT-02-02 reports as
   `not_counted` ("a pinch", "to taste") are skipped with outcome `not_counted`.
5. Items are reduced in order until the line is covered. A partly used item keeps the rest (version + 1).
6. **Rounding dust:** if what remains of an item is below 2% of its `initial_quantity`, or below 5 g, 5 ml or
   0.1 piece, the item is deleted. Decided (Q081): these thresholds.
7. **Staples:** if the ingredient is a staple (KIT-10-05) and the household holds no stock of it, the outcome is
   `staple` and nothing happens. If there is stock of it, it is used like any other ingredient. Water is always a
   staple.
8. **Short stock is not an error.** The cooking already happened: what exists is taken, and the line reports the
   amount that was missing (`partial` or `not_in_stock`). Optional lines (garnish) follow the same rule.
9. **Idempotency:** the result of a `cook_event_id` is stored for 30 days. A repeated call returns the stored result
   with `replayed = true` and changes nothing, even if the lines differ.
10. **Concurrency:** the household's candidate items are locked (`SELECT … FOR UPDATE`, ordered by id) for the
    transaction, so two "cooked it" taps or a parallel edit cannot take the same grams twice.
11. Each taken amount writes a `cooked` movement with `ref = cook_event_id` and the item's use-by date at that
    moment (KIT-12-02 uses it to count "used before expiry"). One `stock_version` bump and one `stock.changed`.
12. Unmatched items are never consumed.

## Contract

`contracts/coolinary/kitchen/v1/stock.proto`

```proto
rpc ConsumeStock(ConsumeStockRequest) returns (ConsumeStockResponse);

message ConsumeLine {
  string ingredient_id = 1;
  coolinary.common.v1.Quantity quantity = 2;
  bool optional = 3;
  repeated string preferred_item_ids = 4;   // from the plan allocation; used first
}

message ConsumeStockRequest {
  string household_id = 1;
  string cook_event_id = 2;       // idempotency key from planner, ≤ 64 chars
  string dish_id = 3;             // kept in the movement ref only
  repeated ConsumeLine lines = 4; // 1–60
}

enum ConsumeOutcome {
  CONSUME_OUTCOME_UNSPECIFIED = 0;
  CONSUME_OUTCOME_TAKEN = 1;
  CONSUME_OUTCOME_PARTIAL = 2;
  CONSUME_OUTCOME_NOT_IN_STOCK = 3;
  CONSUME_OUTCOME_STAPLE = 4;
  CONSUME_OUTCOME_UNCONVERTIBLE = 5;
  CONSUME_OUTCOME_NOT_COUNTED = 6;
}

message ConsumedLine {
  string ingredient_id = 1;
  coolinary.common.v1.Quantity taken = 2;      // in the line's unit
  coolinary.common.v1.Quantity missing = 3;
  ConsumeOutcome outcome = 4;
  repeated string emptied_item_ids = 5;
}

message ConsumeStockResponse {
  repeated ConsumedLine lines = 1;
  int64 stock_version = 2;
  bool replayed = 3;
}
```

## Data

`kitchen/migrations/000NN_cook_events.sql`

```sql
-- +goose Up
CREATE TABLE cook_events (
    household_id  uuid NOT NULL,
    cook_event_id text NOT NULL CHECK (char_length(cook_event_id) <= 64),
    result        jsonb NOT NULL,          -- the ConsumeStockResponse, for replays
    created_at    timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (household_id, cook_event_id)
);
-- Rows older than 30 days are deleted by KIT-13-08.

-- +goose Down
DROP TABLE cook_events;
```

## Errors

Reason codes (in backticks) travel as trailing gRPC metadata `x-error-code` (README decision 18); version
conflicts use `ABORTED`.

| Case | gRPC status | Behaviour |
|---|---|---|
| No `cook_event_id`, 0 or > 60 lines, quantity ≤ 0 | `INVALID_ARGUMENT` | |
| Unknown ingredient id | `INVALID_ARGUMENT` | names the line |
| Missing stock | — | outcome per line, never an error |

## Acceptance criteria

- [ ] Stock: eggs 4 pcs (expires in 2 days), eggs 6 pcs (in 9 days). Consuming 6 eggs empties the first item and
      leaves 4 in the second.
- [ ] Stock: rice 1,000 g in the pantry. Consuming 300 g leaves 700 g; consuming 996 g of a 1,000 g item deletes it
      (dust rule).
- [ ] A recipe line "2 onions" against onion stock in grams (typical weight 150 g) takes 300 g.
- [ ] Salt (staple, no stock) → outcome `staple`; flour (staple) with 500 g in stock → 200 g taken from stock.
- [ ] Stock of 100 g cheese, line 200 g → `partial`, `missing` 100 g, item deleted, no error.
- [ ] Repeating the call with the same `cook_event_id` returns `replayed = true` and the same lines; stock unchanged.
- [ ] Two concurrent calls with different ids on the same 6 eggs (each needs 4) → one gets 4, the other 2 + 2 missing.
- [ ] Kitchen takes exactly the amounts sent: lines scaled by planner to a 5-portion pot (features.md
      Verification 17) reduce stock by those amounts, with no scaling of its own.
