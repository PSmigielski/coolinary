# GW-12-01 Savings counters and monthly reports

| | |
|---|---|
| Service | gateway |
| Epic | [GW-12 Savings routes](README.md) |
| Release | R2 |
| Size | S |
| Status | todo |
| Depends on | GW-03-03, KIT-12 (`GetSavings`, `ListSavingsReports`) |
| Unblocks | MOB-03 |

## Goal

Everyone sees this month's savings counters; Plus households also get the monthly report and its history ("€38 and
4 kg saved this month", business model §1).

## Business rules

1. `GET /v1/savings` → kitchen `GetSavings(household_id, month)` for the current month (household time zone, kitchen
   rule): `moneySaved` (Money), `kgSaved`, broken down into `usedBeforeExpiry` and `boughtOnOffer`, plus
   `thrownAwayKg` (shown honestly; never subtracted silently). Available on Free and Plus.
2. `GET /v1/savings/reports?pageSize=&pageToken=` → `ListSavingsReports`: finished months, newest first, each with
   the same figures and the change against the previous month. **Plus only** (`plus: "savings_report"`).
3. Figures are labelled `about: true` (they rest on reference and paid prices).
4. Currency is the household's (PLN, EUR, …); kitchen does not convert.

## Contract

| Method | Path | 2xx |
|---|---|---|
| GET | `/v1/savings` | 200 `{ "month": "2026-11", "about": true, "moneySaved": Money, "kgSaved": 4.1, "usedBeforeExpiry": { "money": Money, "kg": 3.2 }, "boughtOnOffer": { "money": Money }, "thrownAwayKg": 0.4 }` |
| GET | `/v1/savings/reports` | 200 `{ "reports": [ { "month": "2026-10", "...": "same fields", "vsPreviousMonth": { "moneyPct": 12 } } ], "nextPageToken" }` |

Tag `savings`, `auth: true`.

## Data

No storage.

## Errors

| Case | HTTP | `code` |
|---|---|---|
| Free on `/reports` | 403 | `plus_required` |

## Acceptance criteria

- [ ] Using an item the day before its expiry date increases `usedBeforeExpiry` (kitchen rule, end to end).
- [ ] A Free household gets counters and 403 on `/reports`.

## Notes

- Decided (Q212): Free sees the **current month's counters**; the **monthly report and past months** are Plus
  (features.md lists the "savings report" under Plus and the counters under "all without AI").
