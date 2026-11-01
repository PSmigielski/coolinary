# GW-06-05 Label photo for an unknown product

| | |
|---|---|
| Service | gateway |
| Epic | [GW-06 Fridge, barcode and intake routes](README.md) |
| Release | R1 |
| Size | S |
| Depends on | GW-01-05, GW-03-01, KIT-05 (`ReadLabel`) |
| Unblocks | MOB-03 |

## Goal

When a barcode is unknown, the user can photograph the label instead of typing; AI reads name, brand and size,
and the user confirms them with `POST /v1/products` (features.md §4 lookup step 5).

## Business rules

1. Multipart: exactly 1 `image` and `meta` `{ "ean" }` (digits, 8–14).
2. Counts as one `scan` for Free households (decided, Q112). Same order as GW-06-04.
3. Calls kitchen `ReadLabel(household_id, ean, image, language)`; returns a **suggestion** only (name, brand, pack
   size, units per pack, confidence). Nothing is saved until the user sends `POST /v1/products` (GW-06-03).
4. AI unavailable → 503 `ai_unavailable`, quota released; the app falls back to typing the name.

## Contract

`POST /v1/intake/label` (multipart) — `auth: true`, tag `intake`, idempotency required, quota `scan`.

```json
200 { "suggestion": { "ean": "5901234123457", "name": "Twaróg półtłusty", "brand": "Piątnica",
                      "packSize": { "value": 250, "unit": "g" }, "unitsPerPack": 1, "confidence": 0.9 },
      "aiUsed": true, "quota": { ... } }
```

## Data

No storage.

## Errors

| Case | HTTP | `code` |
|---|---|---|
| Missing `ean` or image | 400 | `validation` |
| Free quota used up | 429 | `quota_reached` |
| AI unavailable | 503 | `ai_unavailable` |

## Acceptance criteria

- [ ] A label photo returns a suggestion and creates no product.
- [ ] Sending the suggestion to `POST /v1/products` makes the next scan of that EAN resolve.
- [ ] A Free household with 2 scans used gets 429.
