# ADM-06-06 Shop additions and corrections from users

| | |
|---|---|
| Service | admin |
| Epic | [ADM-06 Flyer sources, offers and shops review](README.md) |
| Release | R2 |
| Size | S |
| Status | todo |
| Depends on | ADM-02-03; OFF-02-07 (`ListShopReviews`, `ReviewShop`, `ReviewShopCorrection`), OFF-02-04 (`MarkNotDuplicate`, `ListShopMerges`) |
| Unblocks | — |

## Goal

Users can add a missing shop or correct a misplaced one (features.md §3 "Location"); the founder confirms these
before they change what other users see.

## Business rules

1. The queue is `ListShopReviews` (OFF-02-07 rule 1): pending user-added shops and open corrections (`move`,
   `rename`, `closed`), oldest first. Each item shows the shop, the proposed values and the current ones for a
   correction, the nearest 3 active shops of any chain with their distance in metres, and how many households have
   the shop in their shops. No household ids are shown.
2. Actions on a user-added shop → `ReviewShop`: `approve`, `edit_and_approve` (name, chain or position first), or
   `reject` (offers then follows its OFF-02-05 rule 9). An approved shop goes through offers' duplicate check
   (OFF-02-04).
3. Actions on a correction → `ReviewShopCorrection`: `approve` (applied to the shop; `closed` closes it) or `reject`
   (for a "closed" report the flag is cleared).
4. Duplicates: the merges offers made are listed with `ListShopMerges`; two shops that are really different (e.g.
   two shops in one mall) are marked with `MarkNotDuplicate(shop_id_a, shop_id_b)` so offers never merges them
   again (OFF-02-04 rule 5).
5. Audited (`offers.shop_review.<action>`, `offers.shop.not_duplicate`).
6. Role `admin` only.

## Contract

| Method | Path | Body / response |
|---|---|---|
| GET | `/admin/v1/offers/shop-reviews?kind=&page_token=` | `{ items, next_page_token }` |
| POST | `/admin/v1/offers/shop-reviews/{id}/resolve` | `{ action, fields?, reason? }` → item |
| GET | `/admin/v1/offers/shop-merges?page_token=` | `{ items, next_page_token }` |
| POST | `/admin/v1/offers/shop-merges/not-duplicate` | `{ shop_id_a, shop_id_b }` |

## Data

None in admin.

## Errors

| Case | HTTP |
|---|---|
| Already resolved | 409 |

## Acceptance criteria

- [ ] Approving a user-added shop makes it appear in `GET /v1/shops` for its postcode area.

## Notes

Decided (Q159): every correction is reviewed by hand until volumes are known; none is accepted automatically.
