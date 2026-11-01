# GW-10-04 Scan a flyer that is not covered (async)

| | |
|---|---|
| Service | gateway |
| Epic | [GW-10 Shops, offers, flyers and price matrix routes](README.md) |
| Release | R2 |
| Size | M |
| Status | todo |
| Depends on | GW-01-05, GW-03-03, OFF-10 (`UploadUserFlyer`, `SubmitUserFlyerLink`, `GetUserFlyer`, `ConfirmUserFlyer`, `DeleteUserFlyer`) |
| Unblocks | MOB-03, MOB-07 |

## Goal

A Plus user photographs or shares any flyer that is not covered yet (photos, a PDF, or a link) and gets its offers
extracted; this works in every country from release 2 (features.md §2, §3 source 1).

## Business rules

1. **Plus only** (`plus: "own_flyer"`, GW-03-03). Free → 403 `plus_required` before anything is read. Own flyers
   **never count towards the Free photo-scan quota** (they are not available on Free at all).
2. `POST /v1/flyers`, one of:
   - multipart: 1–10 `image` parts or one `pdf` of at most 40 pages, 25 MB in total (offers' limits, OFF-10-01;
     GW-01-05) plus `meta` `{ "shopId" }` → offers `UploadUserFlyer(household_id, user_id, shop_id, images | pdf)`;
   - JSON `{ "shopId", "url" }` with an `https` URL (≤ 2,000 chars) to a flyer file → offers
     `SubmitUserFlyerLink(household_id, user_id, shop_id, url)`; offers fetches it with its own safety checks
     (OFF-10-02).
   `shopId` is required: a shop from GW-10-01, or one the user just added with `POST /v1/shops`.
3. Extraction runs in the background in offers, so the route is **asynchronous**: offers returns the flyer id at
   once → 202 `{ "flyerId", "state": "processing" }`. Idempotency required (a retried upload must not start a second
   AI job).
4. `GET /v1/flyers/{flyerId}` → `GetUserFlyer`: `state` `processing` | `draft` (with the extracted offers and their
   warnings, and proposed dates) | `confirmed` | `failed` (with `failureReason`, e.g. `no_offers_found`,
   `ai_unavailable`, `interrupted`) | `expired`. Visible only to the household that submitted it. The app polls every
   3 s for up to 3 minutes.
5. `POST /v1/flyers/{flyerId}/confirm` body `{ "validFrom", "validTo", "decisions": [ { "draftId", "keep",
   "corrected"? } ] }` → `ConfirmUserFlyer`: the user checks the dates and keeps, drops or corrects each offer; only
   then do the offers appear (with `source: "user_flyer"`) in `GET /v1/offers` for that household. A draft not
   confirmed within 24 h expires (OFF-10-03 rule 7).
6. `DELETE /v1/flyers/{flyerId}` → `DeleteUserFlyer`: discards a draft, or ends a confirmed flyer's offers at once.
7. Decided (Q173): offers read from a user's flyer are private to the uploading household and never used as public
   prices (OFF-10-04); the gateway passes through what offers returns.
8. Images and PDFs are not stored (D6); offers keeps only the extracted facts.
9. Limits: offers allows 2 jobs processing at once and 10 uploads per household per day (OFF-10-01 rule 4); the
   gateway adds the `upload` rate-limit bucket.

## Contract

| Method | Path | 2xx |
|---|---|---|
| POST | `/v1/flyers` (multipart or JSON) | 202 `{ "flyerId": "fly_...", "state": "processing", "pollAfterSeconds": 3 }` |
| GET | `/v1/flyers/{flyerId}` | 200 `{ "flyerId", "shopId", "state": "draft", "draft": { "offers": [ { "draftId", "offer", "warnings": [] } ] }, "validFrom", "validTo" }` |
| POST | `/v1/flyers/{flyerId}/confirm` | 200 `{ "flyerId", "state": "confirmed", "offerCount": 37 }` |
| DELETE | `/v1/flyers/{flyerId}` | 204 |

Tag `flyers`, `auth: true`, `plus: "own_flyer"`.

## Data

No storage.

## Errors

| Case | HTTP | `code` |
|---|---|---|
| Free | 403 | `plus_required` |
| Neither files nor URL, bad URL, no `shopId` | 400 | `validation` |
| File limits | 413 / 415 | GW-01-05 codes |
| Unknown or other household's flyer | 404 | `not_found` |
| Confirming a flyer that is not a draft (expired, failed) | 422 | offers' code |
| More than 10 uploads today or 2 already processing | 429 | offers' code / `rate_limited` |

## Acceptance criteria

- [ ] A Free token gets 403 and the quota counters are unchanged.
- [ ] A 12-page PDF returns 202 at once; polling reaches `draft` with offers that each have a price and dates;
      after confirm they appear in `GET /v1/offers` for this household only.
- [ ] A retried `POST` with the same idempotency key returns the same `flyerId`.
- [ ] Another household asking for the flyer id gets 404.
- [ ] A 25-page PDF gets 413 before offers is called.
