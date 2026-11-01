# GW-12-02 Shareable savings card data

| | |
|---|---|
| Service | gateway |
| Epic | [GW-12 Savings routes](README.md) |
| Release | R3 |
| Size | S |
| Depends on | GW-12-01, KIT-12 (`GetSavingsCard`) |
| Unblocks | MOB-03 |

## Goal

A household shares a card with last month's savings (business model §7 "Built-in sharing": "a shareable monthly
savings card"), which also brings new installs.

## Business rules

1. `GET /v1/savings/card?month=YYYY-MM` (default last finished month) → kitchen `GetSavingsCard(household_id, month,
   language)`: month, `moneySaved`, `kgSaved`, `mealsCooked` and the localized `headline` (kitchen builds it in the
   requested language). The gateway adds only `shareUrl` (the store/landing link with a campaign tag
   `utm_source=savings_card`).
2. No names, no household id, no e-mail on the card or in `shareUrl`.
3. The image is drawn on the phone (decided, Q225); the gateway returns data only.
4. Available to Plus households (the report is Plus); Free → 403 `plus_required`.

## Contract

`GET /v1/savings/card?month=2026-10` — `auth: true`, `plus: "savings_report"`, tag `savings`.

```json
200 { "month": "2026-10", "moneySaved": { "amountMinor": 3800, "currency": "EUR" }, "kgSaved": 4,
      "mealsCooked": 41, "headline": "€38 and 4 kg saved in October",
      "shareUrl": "https://coolinary.app/?utm_source=savings_card&utm_campaign=2026-10" }
```

## Data

No storage.

## Errors

| Case | HTTP | `code` |
|---|---|---|
| Free | 403 | `plus_required` |
| Month without data | 404 | `not_found` |

## Acceptance criteria

- [ ] The response contains no member name, household id or e-mail.
- [ ] The headline is in `x-language`.

## Notes

- Decided (Q225): the card image is drawn on the phone; there is no image service on the server.
