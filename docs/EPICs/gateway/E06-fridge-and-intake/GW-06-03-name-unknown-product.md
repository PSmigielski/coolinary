# GW-06-03 Name an unknown product, confirm a shared one

| | |
|---|---|
| Service | gateway |
| Epic | [GW-06 Fridge, barcode and intake routes](README.md) |
| Release | R1 |
| Size | S |
| Status | todo |
| Depends on | GW-06-02, KIT-05 (`NameProduct`, `ConfirmProposal`, `ReportProduct`) |
| Unblocks | MOB-03 |

## Goal

When a scan finds nothing, the user names the product once and it is known to everyone afterwards; it counts as
confirmed once a second household accepts it unchanged (features.md §4, verification 14).

## Business rules

1. `POST /v1/products` body `{ "ean", "name", "brand"?, "packSize"?, "unitsPerPack"?, "ingredientId"? }`.
   `name` 2–80 characters; `ean` digits only, 8–14. Calls kitchen `NameProduct(household_id, ...)`. The product joins
   this household's products immediately and the shared catalogue as `pending` (kitchen rules, KIT-05).
2. `POST /v1/products/{ean}/confirm` body `{ "fingerprintId", "accept": true, "edits"? }` — a household answers the
   `proposal` that `ResolveBarcode` returned for a pending product (KIT-05-02 rule 4). Calls
   `ConfirmProposal(household_id, code, fingerprint_id, accept, edits)`: accepting unchanged counts toward
   confirmation; `edits` (same fields as `POST /v1/products`) record a competing version. Confirming a product your
   own household named → 422 `own_product` (kitchen).
3. `POST /v1/products/{ean}/report` body `{ "reason": "wrong_name" | "wrong_size" | "not_food" | "other" }`
   calls `ReportProduct(household_id, code, reason)`. Three households reporting the same confirmed product send it
   to admin review (KIT-05-02 rule 7). Response 204.
4. Idempotency required on all three.
5. Names are user text: not logged.

## Contract

| Method | Path | Body | 2xx |
|---|---|---|---|
| POST | `/v1/products` | `{ "ean": "5901234123457", "name": "Twaróg półtłusty", "packSize": { "value": 250, "unit": "g" } }` | 201 product + suggested (as GW-06-02) |
| POST | `/v1/products/{ean}/confirm` | `{ "fingerprintId": 81, "accept": true }` | 200 `{ "product", "productConfirmed" }` |
| POST | `/v1/products/{ean}/report` | `{ "reason": "wrong_name" }` | 204 |

Tag `fridge`, `auth: true`.

## Data

No storage.

## Errors

| Case | HTTP | `code` |
|---|---|---|
| Bad EAN or name | 400 | `validation` |
| Product already confirmed with other data | 409 | `already_exists` |
| Confirming your own product | 422 | `own_product` |

## Acceptance criteria

- [ ] Account A names an unknown EAN; account B (another household) scanning it gets `found: true`,
      `status: "pending"` (verification 14).
- [ ] B confirms; a third household's scan returns `status: "confirmed"`.
- [ ] A confirming its own product gets 422 `own_product`.
- [ ] A report with reason `not_food` returns 204 and reaches kitchen; an unknown reason returns 400.
