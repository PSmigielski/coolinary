# MOB-07-03 Flyer inputs: photos, PDF, link

| | |
|---|---|
| Service | mobile |
| Epic | [MOB-07 Uploads](README.md) |
| Release | R2 |
| Size | S |
| Status | todo |
| Depends on | MOB-07-01, MOB-07-02, GW-10 (`POST /v1/flyers`), OFF-10 |
| Unblocks | MOB-03-08 |

## Goal

A user can photograph, share or link any flyer that Coolinary doesn't cover yet and get its offers (features.md
"Flyers and offers": own-flyer scan, works in every country).

## Business rules

1. Three inputs, one route `POST /v1/flyers` (GW-10-04). The hook checks `usePlusGate("own_flyer")`; the server
   decides. Own flyers are Plus only and never use the Free scan quota (docs/EPICs/README.md "Free and Plus").

   Inputs:

   | Input | Limits | Request |
   |---|---|---|
   | Photos | 1–10 pages, prepared as `flyer_page` (MOB-07-01) | multipart `image` parts + `meta { shopId }` |
   | PDF | ≤ 25 MB and ≤ 40 pages (25 MB per request), picked with `expo-document-picker` (`type: "application/pdf"`, `copyToCacheDirectory: true`) | multipart `pdf` + `meta { shopId }` |
   | Link | `https://` only, ≤ 2,000 characters | JSON `{ "shopId", "url": "https://…" }` |

   Limits follow offers (OFF-10-01) and the gateway (GW-01-05). Decided (Q022): own-flyer PDFs may have up to 40
   pages and 25 MB. The phone checks the size before uploading; the gateway checks the page count.
2. `shopId` is required (offers files a user flyer under a shop, OFF-10-01): a shop from `useShopsNear`, or one the
   user just added (`POST /v1/shops`).
3. Shared files: receiving a PDF or link shared from another app (share intent) is a frontend task; it calls the same
   function.
4. The PDF copy in the cache directory is deleted after the request, like prepared images.
5. The answer is always asynchronous: `202 { flyerId, state: "processing", pollAfterSeconds }`, then
   `GET /v1/flyers/{flyerId}` every `pollAfterSeconds` (3 s) for up to 3 minutes until `draft`, `failed` or
   `expired`. A `draft` is shown for review and confirmed with `POST /v1/flyers/{flyerId}/confirm` (dates and keep /
   drop / correct per offer) or discarded with `DELETE` (GW-10-04). A draft expires after 24 h.

## Contract

```ts
export function submitFlyer(input: { kind: "photos"; pages: PreparedImage[] } | { kind: "pdf"; uri: string } | { kind: "link"; url: string }, opts: { shopId: string }): UploadHandle<FlyerSubmitResult>;
```

## Data

None persisted.

## Errors

| Case | Behaviour |
|---|---|
| PDF over 25 MB | `pdf_too_large` before any upload |
| PDF over 40 pages | gateway 413 `payload_too_large`, shown as `pdf_too_large` |
| `http://` or malformed link | `validation` / `url_invalid` on the phone |

## Acceptance criteria

- [ ] A 24-page PDF of a local flyer returns offers or a job that finishes with offers.
- [ ] A 26 MB PDF is refused on the phone without a request; a 41-page PDF of 5 MB is refused by the gateway and
      shown as `pdf_too_large`.
