# GW-06-02 Barcode and 2D code lookup

| | |
|---|---|
| Service | gateway |
| Epic | [GW-06 Fridge, barcode and intake routes](README.md) |
| Release | R1 |
| Size | S |
| Status | todo |
| Depends on | GW-01-02, KIT-04 (`ResolveBarcode` extended, `ListHouseholdProducts`) |
| Unblocks | MOB-03 (scan and unpack mode) |

## Goal

A scanned code — EAN-8, EAN-13, UPC-12, GTIN-14, a shop-internal weighed code (prefix 20–29) or a GS1 2D code
(DataMatrix / QR / Digital Link carrying the use-by date) — becomes a product with suggested storage place and
date (features.md §4 "How one scan is resolved").

## Business rules

1. New route `GET /v1/products/resolve?code=<raw scanned text>` (URL-encoded; 2D codes contain `/` and group
   separators, so they can't be a path segment). 1–512 characters.
2. The existing `GET /v1/products/{ean}` stays as an alias for digit-only codes (8–14 digits) and is marked
   deprecated in OpenAPI.
3. Calls kitchen `ResolveBarcode(code, household_id, language, user_scan = true)` (every call on this route is one
   real scan, which Open Food Facts' terms require for live lookups, KIT-04). All parsing (check digit, weighed-code patterns,
   GS1 element strings) is kitchen's (KIT-04); the gateway only checks length and printable characters
   (plus the GS1 separator `\x1d`).
4. The household id is now filled from the JWT (today it is sent empty), so household products are found first
   (lookup order step 1).
5. Not found → 200 with `found: false` and `next: "name_product"` (not 404): an unknown product is a normal step,
   the app then offers naming or a label photo (GW-06-03, GW-06-05).
6. Response cached nowhere in the gateway (kitchen caches, KIT-04).
7. **Offline copy:** `GET /v1/products/mine?updatedSince=&pageToken=` → kitchen `ListHouseholdProducts(household_id,
   updated_since, page)`: the household's own products changed since a time, 500 per page, so the phone resolves
   codes it scanned before without the network (features.md §4 lookup step 1, KIT-04-04).

## Contract

`GET /v1/products/resolve?code=...` — `auth: true`, tag `fridge`. `GET /v1/products/mine?updatedSince=` — `auth: true`,
tag `fridge`, 200 `{ "products": [...], "nextPageToken" }`.

```json
200 {
  "found": true, "source": "open_food_facts",
  "product": { "ean": "5900512300108", "name": "Mleko 3,2%", "brand": "Łaciate",
               "packSize": { "value": 1000, "unit": "ml" }, "unitsPerPack": 1, "status": "confirmed" },
  "ingredient": { "id": "ing_milk", "name": "mleko" },
  "suggested": { "storage": "fridge", "expiresOn": "2026-11-08" },
  "fromCode": { "expiresOn": null, "batch": null, "weightG": null, "priceMinor": null },
  "next": null
}
200 { "found": false, "source": null, "product": null, "ingredient": null, "suggested": null,
      "fromCode": { "weightG": 412 }, "next": "name_product" }
```

`fromCode` holds what the code itself carried: use-by date and batch from a 2D code, weight or price from a weighed
code. `product.status`: `confirmed` | `pending` (named by one household, not yet confirmed by a second).

## Data

No storage.

## Errors

| Case | HTTP | `code` |
|---|---|---|
| Empty or > 512 chars, non-printable | 400 | `validation` |
| Kitchen says the code is invalid (bad check digit) | 400 | `invalid_code` |
| Kitchen down | 503 | `service_unavailable` |

## Acceptance criteria

- [ ] `code=9520000000017` returns the example product with `source: "shared_catalogue"`.
- [ ] A GS1 Digital Link `https://id.gs1.org/01/09520000000017/10/ABC/17/261130` returns `fromCode.expiresOn:
      "2026-11-30"` and `batch: "ABC"` (kitchen parses).
- [ ] An unknown code returns 200 `found: false`, `next: "name_product"`.
- [ ] `GET /v1/products/9520000000017` still works and is flagged `deprecated` in OpenAPI.
- [ ] Verification 14 (8-digit, weighed and multipack codes) passes through this route.
