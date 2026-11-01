# ADM-06-05 Held-back offers queue

| | |
|---|---|
| Service | admin |
| Epic | [ADM-06 Flyer sources, offers and shops review](README.md) |
| Release | R2 |
| Size | M |
| Status | todo |
| Depends on | ADM-06-03 (page images); OFF-07-06 (`ListHeldOffers`, `UpdateHeldOffer`, `ReleaseHeldOffer`, `RejectHeldOffer`, `RejectHeldOffers`); KIT-03-06 (`DecideMatch`) |
| Unblocks | — |

## Goal

Offers that failed the non-AI checks (price not above zero, dates outside the flyer's validity, unit missing,
product not matched; features.md §3 step 5) are fixed or discarded by hand instead of silently lost.

## Business rules

1. The queue is `ListHeldOffers` (OFF-07-06 rule 1): each held offer with chain, edition, flyer page reference
   (`version_id`, `page_no`, shown through ADM-06-03 rule 7), the extracted fields, the reasons it was held, and
   the raw AI result. Order: soonest start date first, then flyer, then page. Filters: chain, edition, reason,
   "starts within N days".
2. Actions:
   - `fix`: corrected product name, brand, pack, price, regular price, conditions, dates or ingredient →
     `UpdateHeldOffer`. offers re-runs its normalisation and checks; with no reason left the offer becomes `ready`
     and is published by the next publish run, otherwise the remaining reasons come back;
   - `match_ingredient`: an `UpdateHeldOffer` with the ingredient id, plus `KitchenAdminService.DecideMatch`
     (`set_ingredient`) on the product name's match so everyone's next offer with that name matches (OFF-07-06
     rule 3);
   - `release`: publish despite the remaining reasons → `ReleaseHeldOffer`; offers allows it only for
     `low_confidence` and `implausible_price` (OFF-07-06 rule 4);
   - `reject` with a reason text → `RejectHeldOffer`.
3. Offers whose end date passes without a decision become `expired` (offers) and leave the queue.
4. Bulk: reject every held offer of one flyer version (a misread page) → `RejectHeldOffers(version_id, reason)`.
5. Audited (`offers.held.<action>`, offer id or version id, changed fields).
6. Role `admin` only.

## Contract

| Method | Path | Body / response |
|---|---|---|
| GET | `/admin/v1/offers/held?chain=&edition=&reason=&starts_within_days=&page_token=` | `{ items, next_page_token }` |
| POST | `/admin/v1/offers/held/{id}/resolve` | `{ action: fix\|match_ingredient\|release\|reject, fields?, ingredient_id?, reason? }` → `{ state, reasons[] }` |
| POST | `/admin/v1/offers/held/reject-version` | `{ version_id, reason }` → `{ rejected }` |

## Data

None in admin.

## Errors

| Case | HTTP |
|---|---|
| `release` with a reason that must be fixed | 409 with the reasons (offers' `FAILED_PRECONDITION`) |
| Offer no longer held | 409 |

## Acceptance criteria

- [ ] A held offer with price 0 fixed to 3.89 becomes `ready` and is published by the next publish run.
- [ ] Matching an unmatched product stores the name match in kitchen, and the next offer with the same name is not
      held for it.
