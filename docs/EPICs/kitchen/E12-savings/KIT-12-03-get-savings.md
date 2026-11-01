# KIT-12-03 GetSavings and ListSavingsReports

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-12 Savings](README.md) |
| Release | R2 |
| Size | S |
| Depends on | KIT-12-02 |
| Unblocks | GW-12 (`GET /v1/savings`, `GET /v1/savings/reports`), KIT-12-04, KIT-12-05 |

## Goal

The app shows this month's counters to everyone and the monthly report with history to Plus households
(features.md §6: savings report is Plus).

## Business rules

1. `GetSavings(household_id, month)`: one month, default the current local month. Figures:
   1. `used_before_expiry`: money (`rescued_minor`) and kg (`rescued_g` ÷ 1,000, one decimal);
   2. `bought_on_offer`: money (`offer_minor`);
   3. `money_saved` = used before expiry + bought on offer; `kg_saved` = used-before-expiry kg;
   4. `thrown_away`: kg and money;
   5. `meals_cooked`;
   6. `about = true` always (the figures rest on prices users paid; features.md "totals are labelled about").
2. Available on Free and Plus for the current month. Decided (Q212, shared with GW-12-01): Free sees the current
   month's counters; reports and past months are Plus.
3. A past month through `GetSavings` on Free → `plus_required`.
4. `ListSavingsReports(household_id, page)`: finished months, newest first, 12 per page, each with the figures of
   rule 1 and `vs_previous_month` = percentage change of `money_saved` (whole number; empty when the previous month
   had 0). Plus only.
5. Several currencies in one month: one figure set per currency, the household's current currency first.
6. A month with no row returns zeros, not an error.

## Contract

`contracts/coolinary/kitchen/v1/savings.proto`

```proto
rpc GetSavings(GetSavingsRequest) returns (GetSavingsResponse);
rpc ListSavingsReports(ListSavingsReportsRequest) returns (ListSavingsReportsResponse);

message SavingsFigures {
  string currency = 1;
  coolinary.common.v1.Money money_saved = 2;
  double kg_saved = 3;
  coolinary.common.v1.Money used_before_expiry_money = 4;
  double used_before_expiry_kg = 5;
  coolinary.common.v1.Money bought_on_offer_money = 6;
  double thrown_away_kg = 7;
  coolinary.common.v1.Money thrown_away_money = 8;
}

message MonthSavings {
  string month = 1;                    // YYYY-MM
  bool about = 2;
  repeated SavingsFigures figures = 3;
  int32 meals_cooked = 4;
  optional int32 vs_previous_month_pct = 5;
}

message GetSavingsRequest { string household_id = 1; string month = 2; }
message GetSavingsResponse { MonthSavings savings = 1; }

message ListSavingsReportsRequest { string household_id = 1; coolinary.common.v1.PageRequest page = 2; }
message ListSavingsReportsResponse { repeated MonthSavings reports = 1; coolinary.common.v1.PageResponse page = 2; }
```

## Data

No schema change.

## Errors

Reason codes (in backticks) travel as trailing gRPC metadata `x-error-code` (README decision 18); version
conflicts use `ABORTED`.

| Case | gRPC status | Reason |
|---|---|---|
| Free asking for a past month or reports | `FAILED_PRECONDITION` | `plus_required` |
| Month not `YYYY-MM`, or in the future | `INVALID_ARGUMENT` | |

## Acceptance criteria

- [ ] A Plus household with October rescued 3.2 kg / 31.00 € and offers 7.00 € shows money saved 38.00 €, 3.2 kg.
- [ ] A Free household gets the current month and `plus_required` for reports.
- [ ] October 38 € after September 34 € → `vs_previous_month_pct = 12`.
