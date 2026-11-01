# PLN-16-02 Base pantry checklist for a household

| | |
|---|---|
| Service | planner |
| Epic | [PLN-16 Cuisine base pantry](README.md) |
| Release | R1 |
| Size | S |
| Status | todo |
| Depends on | PLN-16-01, KIT-07 (ListStock), KIT-10 (GetPlanningProfile) |
| Unblocks | PLN-16-03, GW-08 |

## Goal

Show, for a cuisine, which base items the household has and which are missing.

## Business rules

1. An item counts as at home when the household's stock has any amount of it (not expired) or it is one of the
   household's staples.
2. Response per item: ingredient id, name in `x-language`, `at_home`, and `on_list` when an open list line exists.
3. Works for any cuisine code, not only the household's favourites.
4. kitchen down → `UNAVAILABLE` (the checklist would be wrong).

## Contract

`contracts/coolinary/planner/v1/pantry.proto`

```proto
rpc GetBasePantry(GetBasePantryRequest) returns (BasePantryChecklist);
message GetBasePantryRequest { string household_id = 1; coolinary.common.v1.Cuisine cuisine = 2; }
message BasePantryChecklist { coolinary.common.v1.Cuisine cuisine = 1; repeated BasePantryItem items = 2; int32 missing = 3; }
message BasePantryItem { string ingredient_id = 1; string name = 2; bool at_home = 3; bool on_list = 4; }
```

## Data

No schema change.

## Errors

| Case | gRPC status / behaviour |
|---|---|
| `cuisine` unspecified | `INVALID_ARGUMENT` |
| No base pantry defined for that cuisine | `NOT_FOUND` |
| kitchen down | `UNAVAILABLE` (rule 4) |

## Acceptance criteria

- [ ] Fridge with rice and ginger: Asian checklist shows soy sauce and sesame oil missing, `missing = 2`.
- [ ] Olive oil as a staple counts as at home for Italian.
