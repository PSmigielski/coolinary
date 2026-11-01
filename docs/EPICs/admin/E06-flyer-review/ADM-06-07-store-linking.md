# ADM-06-07 Link shops to chain store numbers

| | |
|---|---|
| Service | admin |
| Epic | [ADM-06 Flyer sources, offers and shops review](README.md) |
| Release | R2 |
| Size | S |
| Status | todo |
| Depends on | ADM-02-03; OFF-03-05 (`ListStoreLinkQueue`, `LinkStore`, `UnlinkStore`) |
| Unblocks | — |

## Goal

Chains whose flyer differs per store (Kaufland) show each user their own store's edition. When offers can't link an
OpenStreetMap shop to the chain's store number by itself, the founder does it here.

## Business rules

1. The queue is `ListStoreLinkQueue` (OFF-03-05 rule 3): active shops of a per-store chain that auto-link left
   unlinked, each with its candidate stores (store number, name, address) and their distance in metres.
2. Actions: `link` a shop to one store → `LinkStore(shop_id, store_ref)`; `unlink` → `UnlinkStore(shop_id)`. One
   store links to at most one shop and one shop to at most one store; offers refuses a second link (OFF-03-05 rule
   4). A link made here is never changed by auto-link (rule 5 there).
3. Linking creates the store edition if it doesn't exist (OFF-03-05 rule 7); the edition is read once a user has the
   shop in range (OFF-05-07).
4. Audited (`offers.store.link`, `offers.store.unlink`, shop id, store ref).
5. Role `admin` only.

## Contract

| Method | Path | Body / response |
|---|---|---|
| GET | `/admin/v1/offers/store-links?chain=&page_token=` | `{ items: [{ shop, candidates[{ store_ref, name, address, distance_m }] }], next_page_token }` |
| POST | `/admin/v1/offers/store-links` | `{ shop_id, store_ref }` → link |
| DELETE | `/admin/v1/offers/store-links/{shopId}` | → `{}` |

## Data

None in admin.

## Errors

| Case | HTTP |
|---|---|
| Store or shop already linked | 409 (offers' `ALREADY_EXISTS` / `FAILED_PRECONDITION`) |
| Unknown shop or store | 404 |

## Acceptance criteria

- [ ] Linking a Kaufland shop to store 1234 removes it from the queue, and a user with that shop in range sees store
      1234's edition (features.md Verification 8).
- [ ] Linking a second shop to the same store returns 409.
