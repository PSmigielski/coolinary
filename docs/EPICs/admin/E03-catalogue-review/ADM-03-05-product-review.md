# ADM-03-05 Unknown products review

| | |
|---|---|
| Service | admin |
| Epic | [ADM-03 Catalogue, matching and nutrition review](README.md) |
| Release | R1 |
| Size | M |
| Depends on | ADM-02-03; KIT-05-05 (`ListProductReviews`, `GetSharedProduct`, `ResolveProduct`, `EditSharedProduct`) |
| Unblocks | — |

## Goal

Products that users named themselves (barcode not found anywhere, features.md §4 step 5) are checked when the
automatic rule cannot decide, so the shared catalogue stays clean.

## Business rules

1. A product named by one household becomes confirmed automatically when a second household accepts it unchanged
   (KIT-05). The admin queue is kitchen's `ListProductReviews` (KIT-05-05 rule 1): products kitchen marked
   `disputed` (conflicting submissions or reports, KIT-05-02), newest first.
2. Each item shows the code, every fingerprint's values (name, brand, size, matched ingredient, language), how many
   households hold each, first and last date, and report reasons with counts. Opening one calls
   `GetSharedProduct(code)` for the full row and its change history.
3. Actions, each one kitchen call:
   - `confirm` one fingerprint → `ResolveProduct(code, confirm_fingerprint(id))`;
   - `set_values` (the admin's own name, brand, size, ingredient) → `ResolveProduct(code, set_values(fields))`;
   - `retire` (no longer answered from the shared catalogue; each household keeps its own product) →
     `ResolveProduct(code, retire)`;
   - `edit` a product that is already confirmed (not in the queue, opened by code) → `EditSharedProduct(code,
     fields)`.
   Changing the ingredient also updates the stored name match (KIT-05-05 rule 6).
4. Names typed by users can contain personal text; admin never logs them and the audit row stores only the code
   and the action.
5. Audited (`kitchen.product.resolve`, code, action).
6. Role `admin` only.

## Contract

| Method | Path | Body / response |
|---|---|---|
| GET | `/admin/v1/kitchen/product-reviews?page_token=` | `{ items, next_page_token }` |
| GET | `/admin/v1/kitchen/products/{code}` | shared product with fingerprints and history |
| POST | `/admin/v1/kitchen/products/{code}/resolve` | `{ action: confirm\|set_values\|retire, fingerprint_id?, fields? }` → product |
| PATCH | `/admin/v1/kitchen/products/{code}` | `{ fields }` → product |

## Data

None in admin.

## Errors

| Case | HTTP |
|---|---|
| Unknown code | 404 |
| Invalid code or field value | 400 |

## Acceptance criteria

- [ ] `set_values` with an edited name makes the next scan of that code in another household return that name
      (features.md verification 14).
- [ ] `retire` keeps the product in the households that named it.
