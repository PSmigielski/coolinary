# GW-06-01 Stock: list, add, edit, remove

| | |
|---|---|
| Service | gateway |
| Epic | [GW-06 Fridge, barcode and intake routes](README.md) |
| Release | R1 |
| Size | M |
| Depends on | GW-01-03, KIT-07 (`ListStock`, `AddStockItems`, `UpdateStockItem`, `RemoveStockItems`, `PreviewStockItems`, `SuggestStockDefaults`) |
| Unblocks | MOB-03, MOB-04 (offline stock edits) |

## Goal

The household's fridge, freezer and pantry: list it sorted by what expires first, add one or many items (typing,
unpack mode, confirmed proposals), change quantity, place or date, and remove items saying whether they were used
or thrown away (the savings counters need that, KIT-12).

## Business rules

1. `GET /v1/stock`: optional `storage` (`fridge` | `freezer` | `pantry`, repeatable). Calls kitchen
   `ListStock(household_id, storage)`; kitchen's order is what expires first (KIT-07-05). Other orders (name,
   added) are done on the phone. The whole stock is returned (kitchen caps a household at 1,000 items) with
   `version` (kitchen's `stock_version`, increases on every change).
2. `POST /v1/stock` body `{ "items": [...], "source": "typed" | "barcode" | "unpack" | "photo" | "receipt" |
   "text", "intakeRef"? }`, 1–100 items. Each item: `ingredientId`, `ean` or `name` (one required), optional
   `quantity`, `storage`, `expiresOn`, `noExpiry`, `packCodeExpiresOn`, `pricePaid`. `intakeRef` is the reference an
   intake response returned (GW-06-04/06/07); kitchen stores no proposals, so the ticked proposals are simply sent
   back as items (KIT-08-01 rule 9). Missing storage or date: kitchen fills defaults (features.md §2). Calls
   `AddStockItems(household_id, source, items, intake_ref)`. Idempotency **required**. Returns the created items
   with kitchen's `merged` / `duplicate` flags.
3. `PATCH /v1/stock/{itemId}` body any of `quantity`, `storage`, `expiresOn`, `clearExpiry`, `ingredientId`, `name`;
   `If-Match: "<item version>"` (required). Calls `UpdateStockItem(household_id, item_id, version, ...)`. `ABORTED` → 409 `version_conflict`
   (another member edited it). Idempotency required.
4. `DELETE /v1/stock/{itemId}?reason=used|thrown|mistake` (default `used`; optional `quantity` to remove part),
   `If-Match: "<item version>"` (required: kitchen checks the version on removals too, KIT-07-04 rule 1). Calls
   `RemoveStockItems` with one entry (`item_id`, `version`, `reason`, `quantity`). Idempotency required. Removing
   an already removed item → 404 (the offline queue treats 404 on delete as done, MOB-04).
5. Unpack mode and the add form: `POST /v1/stock/preview` body `{ "source", "candidates": [...] }` (1–60 EANs,
   ingredient ids or names) calls `PreviewStockItems` and returns ready-to-confirm proposals; nothing is stored
   (KIT-07-06). `GET /v1/stock/defaults?ingredientId=|ean=&storage=&packCodeExpiresOn=` calls
   `SuggestStockDefaults` and returns the suggested storage, date and shelf life per place (KIT-07-02).
6. Quantities use the `Quantity` component (`g`, `ml`, `piece`); `quantity.value` > 0 and ≤ 100,000.
7. `expiresOn` is a date (`YYYY-MM-DD`), not more than 5 years ahead; past dates are allowed (already expired item).
8. Stock edits by "cooked it" and list ticks go through planner (GW-07-07, GW-08-03), not these routes.

## Contract

| Method | Path | 2xx |
|---|---|---|
| GET | `/v1/stock?storage=` | 200 `{ "version": 41, "items": [item] }` |
| POST | `/v1/stock` | 201 `{ "version": 42, "items": [item] }` |
| PATCH | `/v1/stock/{itemId}` | 200 item |
| DELETE | `/v1/stock/{itemId}?reason=used&quantity=` | 204 |
| POST | `/v1/stock/preview` | 200 `{ "proposals": [proposal] }` |
| GET | `/v1/stock/defaults?ingredientId=&storage=` | 200 `{ "storage", "expiresOn", "expirySource", "shelfLifeDaysByStorage" }` |

```json
item { "id": "stk_...", "ingredientId": "ing_egg", "name": "Jajka", "ean": "590...", "quantity": { "value": 6,
       "unit": "piece" }, "storage": "fridge", "expiresOn": "2026-11-12", "addedAt": "2026-11-01T09:00:00Z",
       "source": "barcode", "version": 1, "expiresInDays": 11 }
```

Tag `fridge`, `auth: true`.

## Data

No storage.

## Errors

| Case | HTTP | `code` |
|---|---|---|
| 0 or > 100 items, bad quantity or date | 400 | `validation` |
| Item not in this household | 404 | `not_found` |
| Version mismatch | 409 | `version_conflict` |
| Stock cap reached | 422 | `stock_limit` |

## Acceptance criteria

- [ ] Unpack mode: one `POST` with 25 items creates 25 items and one new stock version.
- [ ] Retrying the same `POST` with the same idempotency key creates nothing new.
- [ ] `GET` returns the item expiring tomorrow first (kitchen's order).
- [ ] `DELETE ?reason=thrown` reaches kitchen with reason `thrown`.
- [ ] Two members editing the same item: the second `PATCH` with the old version gets 409; the same holds for
      `DELETE` with an old version.
- [ ] `POST /v1/stock/preview` with 3 EANs (one unknown) returns 3 proposals, the unknown one marked
      `unknown_product`, and the stock version is unchanged.
