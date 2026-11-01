# KIT-12-05 Shareable savings card data

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-12 Savings](README.md) |
| Release | R3 |
| Size | S |
| Status | todo |
| Depends on | KIT-12-03 |
| Unblocks | GW-12 (`GET /v1/savings/card`) |

## Goal

A household can share a card "€38 and 4 kg saved in October" — the built-in sharing channel of the business model
(§7 "shareable monthly savings card"). Kitchen gives the numbers and the headline; the phone draws the image.

## Business rules

1. `GetSavingsCard(household_id, month, language)`: default the last finished month. Plus only.
2. Returns: month, `money_saved` (rounded down to a whole currency unit), `kg_saved` (whole kg when ≥ 10, else one
   decimal), `meals_cooked`, and `headline` in `language` (pl, de, en) from fixed templates, e.g. "€38 and 4 kg
   saved in October" / "38 zł i 4 kg uratowane w październiku" / "38 € und 4 kg gespart im Oktober".
3. When money saved is 0, the headline names only the kilograms; when both are 0 → `NOT_FOUND` (nothing to share).
4. The card holds no names, no household id, no location, no item names.
5. The share link is the gateway's (GW-12-02); kitchen does not build URLs.

## Contract

`contracts/coolinary/kitchen/v1/savings.proto`

```proto
rpc GetSavingsCard(GetSavingsCardRequest) returns (GetSavingsCardResponse);

message GetSavingsCardRequest { string household_id = 1; string month = 2; string language = 3; }

message GetSavingsCardResponse {
  string month = 1;
  coolinary.common.v1.Money money_saved = 2;
  double kg_saved = 3;
  int32 meals_cooked = 4;
  string headline = 5;
}
```

## Data

No schema change.

## Errors

Reason codes (in backticks) travel as trailing gRPC metadata `x-error-code` (README decision 18); version
conflicts use `ABORTED`.

| Case | gRPC status | Reason |
|---|---|---|
| Free household | `FAILED_PRECONDITION` | `plus_required` |
| Nothing saved that month | `NOT_FOUND` | |

## Acceptance criteria

- [ ] October 38.70 € / 4.2 kg / 41 meals → `money_saved` 38 €, `kg_saved` 4.2, headline in the requested language.
- [ ] Every headline template exists in pl, de and en (test).
- [ ] The response contains no field that identifies the household.
