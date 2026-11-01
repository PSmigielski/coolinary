# ADM-06-08 Reference prices

| | |
|---|---|
| Service | admin |
| Epic | [ADM-06 Flyer sources, offers and shops review](README.md) |
| Release | R2 |
| Size | S |
| Status | todo |
| Depends on | ADM-02-03; OFF-11-03 (`ImportReferencePrices`, `ListReferencePrices`) |
| Unblocks | — |

## Goal

The founder can see the reference price per country and ingredient that the trip plan falls back to when a shop
has no flyer or receipt price, and override it by uploading a CSV.

## Business rules

1. The list is `ListReferencePrices`, filtered by country and ingredient: unit basis, price, whether it was computed
   or imported, the date computed and the number of observations behind it (OFF-11-03 rule 3).
2. Upload: a CSV with `country_code, ingredient_id, unit_basis, price_minor`, at most 2 MB, sent as is to
   `ImportReferencePrices`. offers validates it row by row and answers with imported and rejected rows; rows with
   unknown ingredients are reported, not imported (OFF-11-03 rule 4). Admin shows the rejected rows with offers'
   reasons.
3. Imported values override computed ones until replaced (OFF-11-03 rules 1 and 4). Computed reference prices are
   the weekly median of 90 days of observed prices, from at least 3 prices at 2 or more chains (Q174).
4. Audited (`offers.reference_prices.import`, file SHA-256, imported and rejected counts).
5. Role `admin` only.

## Contract

| Method | Path | Body / response |
|---|---|---|
| GET | `/admin/v1/offers/reference-prices?country=&ingredient_id=&page_token=` | `{ items, next_page_token }` |
| POST | `/admin/v1/offers/reference-prices/import` | multipart `file` → `{ imported, rejected: [{ row, reason }] }` |

## Data

None in admin.

## Errors

| Case | HTTP |
|---|---|
| File over 2 MB or not CSV | 413 / 400 |
| Every row rejected | 200 with `imported: 0` and the reasons |

## Acceptance criteria

- [ ] Uploading a CSV with one valid row and one unknown ingredient imports one row and lists the other as rejected.
- [ ] After the import, the list shows that ingredient's price as imported.
