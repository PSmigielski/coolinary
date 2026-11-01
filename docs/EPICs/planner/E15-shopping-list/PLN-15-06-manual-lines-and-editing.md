# PLN-15-06 Manual lines, editing, removing, clearing

| | |
|---|---|
| Service | planner |
| Epic | [PLN-15 Shopping list](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | PLN-15-01, PLN-15-08, KIT-03 (MatchNames) |
| Unblocks | GW-08 |

## Goal

Anyone in the household can add items by hand, change amounts, remove lines and clear what was bought (features.md
"Shopping list" step 4: "anything a household member adds by hand").

## Business rules

1. Add by ingredient (from autocomplete, KIT-02 `SearchIngredients` via the gateway) with an optional amount; if
   no amount is given, one pack (smallest known pack, else 1 piece) is used.
2. Add by free text: the text is matched once with KIT-03 `MatchNames` (context `list`, language from `x-language`);
   a confident match makes it an ingredient line, otherwise it stays free text. kitchen down → free text.
3. Adding an ingredient with an open line merges into it (adds a `manual` contribution); the response returns the
   existing line id.
4. Source of a manual line: kind `manual`, ref = member id ("added by Ewa" in the app).
5. Editing the amount sets `manual_amount = true`; from then on plan changes don't change the amount but flag the
   line `plan_changed` (PLN-15-10).
6. Removing a line tombstones it (state removed, kept 7 days for sync). Removing a line whose sources include plan
   meals does not change the plan.
7. `ClearTickedLines` removes all ticked lines (tombstones); ticked lines are also cleared automatically 48 hours
   after ticking.
8. Every change carries `client_mutation_id` and follows PLN-15-08.
9. Limits: 100 characters of free text, 300 active lines.

## Contract

`contracts/coolinary/planner/v1/list.proto`

```proto
rpc AddListLine(AddListLineRequest) returns (ListLine);
rpc UpdateListLine(UpdateListLineRequest) returns (ListLine);
rpc RemoveListLine(RemoveListLineRequest) returns (RemoveListLineResponse);
rpc ClearTickedLines(ClearTickedLinesRequest) returns (ShoppingList);

message AddListLineRequest {
  string household_id = 1; string member_id = 2; string client_mutation_id = 3;
  oneof item { string ingredient_id = 4; string free_text = 5; }
  coolinary.common.v1.Quantity amount = 6;
}
message UpdateListLineRequest {
  string household_id = 1; string member_id = 2; string client_mutation_id = 3;
  string line_id = 4; int64 base_version = 5;
  coolinary.common.v1.Quantity amount = 6;      // set = manual amount
}
message RemoveListLineRequest { string household_id = 1; string member_id = 2; string client_mutation_id = 3; string line_id = 4; }
message RemoveListLineResponse { int64 list_version = 1; }
message ClearTickedLinesRequest { string household_id = 1; string member_id = 2; string client_mutation_id = 3; }
```

## Data

No schema change (PLN-15-01 tables). The 48-hour auto-clear runs hourly.

## Errors

Reason codes (in backticks) travel as trailing gRPC metadata `x-error-code` (README decision 18); version
conflicts use `ABORTED`.

| Case | gRPC status |
|---|---|
| Text empty or > 100 characters | `INVALID_ARGUMENT` |
| Line not found or removed | `NOT_FOUND` |
| > 300 active lines | `RESOURCE_EXHAUSTED` |

## Acceptance criteria

- [ ] "milk" typed in Polish ("mleko") becomes the milk ingredient line; "birthday candles" stays free text.
- [ ] Adding milk twice gives one line with two manual contributions.
- [ ] A free-text line never appears in logs.
- [ ] Ticked lines disappear 48 hours after ticking.
