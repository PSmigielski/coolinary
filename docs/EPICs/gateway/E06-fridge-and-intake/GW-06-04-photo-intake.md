# GW-06-04 Fridge or shelf photo → proposals

| | |
|---|---|
| Service | gateway |
| Epic | [GW-06 Fridge, barcode and intake routes](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | GW-01-05, GW-03-01, GW-03-02, KIT-08 (`IntakeFromPhoto`) |
| Unblocks | MOB-03, MOB-07 |

## Goal

A photo of the fridge or a shelf becomes a list of items the user confirms (features.md §2 "Fridge and pantry",
AI per request). Free households get 2 scans a week (shared `scan` quota).

## Business rules

1. Multipart: 1–4 `image` parts (several shelves in one scan count as **one** scan) and `meta`
   `{ "storage": "fridge" | "freezer" | "pantry" }` (where the photo was taken, used for defaults).
2. Order: rate limit (`ai`, `upload`) → idempotency (required) → quota reserve (`scan`) → kitchen
   `IntakeFromPhoto(household_id, images, storage, language)` with the `ai` deadline (60 s) → commit if
   `ai_used`, else release.
3. Response is kitchen's `IntakeResponse`; **nothing is written to stock and nothing is stored** (KIT-08-01 rule 9).
   It carries `intakeRef` and proposals; each proposal has a `proposalId`, a ready-made item (name, matched
   ingredient or none, quantity, storage, date), whether it is pre-ticked and notes (`low_confidence`,
   `not_matched`, `estimated_quantity`). Proposals have no server-side expiry. The app sends the ticked (possibly
   edited) items to `POST /v1/stock` with `source: "photo"` and the `intakeRef` (GW-06-01). `truncated: true` means
   kitchen kept only the first 60 items.
4. No AI fallback exists for photos: if kitchen answers `ai_unavailable`, the gateway returns 503
   `ai_unavailable`, the quota is not used, and the app suggests typing or barcodes.
5. Zero items recognised is a success (`proposals: []`, `ai_used: true`, the scan counts).
6. Images are not stored (D6, GW-01-05).

## Contract

`POST /v1/intake/photo` (multipart) — `auth: true`, tag `intake`, idempotency required, quota `scan`.

```json
200 {
  "proposals": [
    { "proposalId": "k3f9...:0", "ticked": true, "notes": [],
      "item": { "name": "Papryka czerwona", "ingredientId": "ing_bell_pepper", "quantity": { "value": 2, "unit": "piece" },
                "storage": "fridge", "expiresOn": "2026-11-06" } }
  ],
  "intakeRef": "k3f9...", "truncated": false,
  "aiUsed": true,
  "quota": { "kind": "scan", "limit": 2, "used": 1, "remaining": 1, "nextFreeAt": "2026-11-05T18:42:10Z" }
}
```

## Data

No storage in the gateway (quota keys from GW-03-01).

## Errors

| Case | HTTP | `code` |
|---|---|---|
| No image, > 4 images, too large, wrong type | 400 / 413 / 415 | GW-01-05 codes |
| Free quota used up | 429 | `quota_reached` |
| AI down, refused or over budget | 503 | `ai_unavailable` |

## Acceptance criteria

- [ ] A Free household's first two photo scans succeed; the third returns 429 before kitchen is called.
- [ ] A scan with 3 images uses one unit.
- [ ] `ai_unavailable` from kitchen returns 503 and `remaining` is unchanged.
- [ ] Stock is unchanged after the call; posting the proposals to `/v1/stock` adds them.
