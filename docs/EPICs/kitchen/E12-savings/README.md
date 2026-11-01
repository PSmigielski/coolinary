# KIT-12 Savings ledger and monthly report

| | |
|---|---|
| Service | kitchen |
| Release | R2 (counters, report); R3 (share card) |
| Depends on | KIT-07 (stock movements, prices), KIT-02-02 (grams), KIT-09-02 (scheduler pattern), NTF-03 (`SendPush`, template `savings_monthly`) |

## Goal

"€38 and 4 kg saved this month": counters of food used before it went off and money saved by buying on offer, shown
honestly next to what was thrown away, plus a monthly report for Plus households (business model §1 "Visible
savings"; features.md §2 "Savings and household"). Derived decision D8 puts this in kitchen.

## Flow

```
stock movement (used / cooked / thrown away / added with offer saving)
   └─► valuation: grams + value share (KIT-12-01) ─► ledger rules (KIT-12-02) ─► savings_monthly (+ meals cooked)
GetSavings (this month, Free and Plus) · ListSavingsReports (finished months, Plus) · GetSavingsCard (R3, Plus)
1st of the month ─► SendPush(savings_monthly) to Plus households with savings (KIT-12-04)
```

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [KIT-12-01](KIT-12-01-movement-valuation.md) | Grams and value of each stock movement | R2 | S | KIT-07-01, KIT-02-02 |
| [KIT-12-02](KIT-12-02-ledger-rules.md) | Ledger rules and monthly totals | R2 | M | KIT-12-01 |
| [KIT-12-03](KIT-12-03-get-savings.md) | GetSavings and ListSavingsReports | R2 | S | KIT-12-02 |
| [KIT-12-04](KIT-12-04-monthly-report-push.md) | Monthly savings push | R2 | S | KIT-12-03, KIT-09-02, NTF-03 |
| [KIT-12-05](KIT-12-05-share-card.md) | Shareable savings card data | R3 | S | KIT-12-03 |

## Decisions (taken 1 November 2026)

- Q209 (KIT-12-01): Count kilograms only; euros only where the price paid is known.
- Q210 (KIT-12-02): Yes, for everyone, so a household that upgrades sees its history.
- Q211 (KIT-12-02): Used within the 3 days before its date (the app's "soon").
- Q212 (KIT-12-03): This month's counters on Free; monthly report and past months on Plus.
- Q213 (KIT-12-04): The 1st of the month at 10:00 local time.
