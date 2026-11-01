# GW-06-07 Receipt photo → proposals

| | |
|---|---|
| Service | gateway |
| Epic | [GW-06 Fridge, barcode and intake routes](README.md) |
| Release | R3 |
| Size | S |
| Depends on | GW-01-05, GW-03-01, KIT-08 (`IntakeFromReceipt`) |
| Unblocks | MOB-03 |

## Goal

A receipt photo (or a shared digital receipt image or PDF from a loyalty app) becomes items to confirm; kitchen
also records the prices for the receipt-price source (features.md §2, §4 "Receipts").

## Business rules

1. Multipart: 1–3 `image` parts (a long receipt photographed in parts) **or** one `pdf`; `meta` `{}`.
2. One `scan` for Free. Same order as GW-06-04, idempotency required.
3. Calls kitchen `IntakeFromReceipt(household_id, images | pdf, language)`. Response: proposals (as GW-06-04) plus
   `receipt { chainKey?, shopName, purchasedOn, dateUncertain, total?, leftOutLines }` (kitchen's `ReceiptInfo`;
   kitchen cannot call offers, so it names the chain by `chain_key`, not an offers chain id). Proposals carry the
   printed name and the price paid. The app confirms by sending the ticked items to `POST /v1/stock` with
   `source: "receipt"` and the `intakeRef`; kitchen then publishes `receipt.scanned` for the confirmed lines
   (KIT-08-04, KIT-08-06). The gateway does nothing extra.
4. Lines that are not food (bags, deposits) are not proposed (kitchen/ai rule).
5. AI unavailable → 503 `ai_unavailable`, quota released.

## Contract

`POST /v1/intake/receipt` (multipart) — `auth: true`, tag `intake`, idempotency required, quota `scan`.

```json
200 { "receipt": { "chainKey": "lidl", "shopName": "Lidl", "purchasedOn": "2026-11-01", "dateUncertain": false,
                   "total": { "amountMinor": 8743, "currency": "PLN" }, "leftOutLines": 1 },
      "intakeRef": "...",
      "proposals": [ { "proposalId": "...", "printedName": "UDKA Z KURCZ.", "item": { "name": "Udka z kurczaka",
                       "pricePaid": { "amountMinor": 1599, "currency": "PLN" } }, "...": "as GW-06-04" } ],
      "aiUsed": true, "quota": { ... } }
```

## Data

No storage.

## Errors

| Case | HTTP | `code` |
|---|---|---|
| > 3 images, both images and PDF | 400 | `validation` |
| Not a receipt | 422 | `not_a_receipt` (kitchen) |
| Free quota used up | 429 | `quota_reached` |
| AI unavailable | 503 | `ai_unavailable` |

## Acceptance criteria

- [ ] A test receipt with 12 food lines and a bag returns 12 proposals with prices.
- [ ] A PDF receipt from a loyalty app is accepted.
- [ ] A photo of a cat returns 422 `not_a_receipt`; whether it counts as a scan follows `ai_used` from kitchen.
