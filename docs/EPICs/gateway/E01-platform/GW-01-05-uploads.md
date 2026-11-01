# GW-01-05 Photo and PDF uploads

| | |
|---|---|
| Service | gateway |
| Epic | [GW-01 Platform](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | GW-01-02 |
| Unblocks | GW-06-04, GW-06-05, GW-06-07, GW-10-04, MOB-07 |

## Goal

One upload path for fridge, shelf, label, receipt and flyer files: checked for type and size, held in memory only,
and passed to the owning service over gRPC.

## Business rules

1. Upload routes take `multipart/form-data` with parts named `image` (repeatable) or `pdf` (single), plus a `meta`
   part holding the route's JSON fields.
2. **Allowed types**, checked by magic bytes, not by the declared content type: JPEG (`FF D8 FF`), PNG
   (`89 50 4E 47`), WebP (`RIFF....WEBP`) and PDF (`%PDF-`). HEIC/HEIF and anything else → 415
   `unsupported_media_type` (the app converts to JPEG, MOB-07).
3. **Limits** (config):

   | | Per file | Files per request | Whole request |
   |---|---|---|---|
   | Images | 5 MB | route-specific, at most 20 | 25 MB |
   | PDF | 25 MB, ≤ 40 pages | 1 | 25 MB |

   The app is expected to resize images to at most 1,568 px on the long edge (MOB-07); the server does not resize.
   Page count is read from the PDF trailer (`/Count`); unreadable → 400 `pdf_unreadable`.
4. Over a limit → 413 `payload_too_large` with `meta: { limitBytes, limitFiles }`. The gateway stops reading the
   body as soon as the whole-request limit is passed (Bun `maxRequestBodySize` for the server, plus a per-route
   check).
5. **Privacy (D6).** Bytes stay in memory for the length of the request: never written to disk, never cached, never
   logged, never put in metrics or idempotency entries (only their SHA-256 for the fingerprint, GW-01-03).
6. EXIF metadata is not inspected or forwarded beyond the raw bytes; the app strips GPS EXIF before upload (MOB-07).
7. gRPC message size: the gateway's clients for kitchen and offers set `grpc.max_send_message_length` to 26 MB.
   Kitchen and offers must accept the same on their servers (noted as a dependency on KIT-08 and OFF-10).
8. Upload routes use the `ai` deadline class and are rate limited (GW-03-04, bucket `upload`).

## Contract

Helper `parseUpload(request, { maxImages, allowPdf })` → `{ images: Uint8Array[], pdf?: Uint8Array, meta: unknown }`,
with `meta` validated against the route's schema. OpenAPI documents each upload route as `multipart/form-data`.

## Data

No storage.

## Errors

| Case | HTTP | `code` |
|---|---|---|
| Not multipart, or `meta` missing | 400 | `validation` |
| Wrong file type | 415 | `unsupported_media_type` |
| File or request too large, too many files | 413 | `payload_too_large` |
| PDF without readable page count, or > 40 pages | 400 / 413 | `pdf_unreadable` / `payload_too_large` |

## Acceptance criteria

- [ ] A 4 MB JPEG passes; a 6 MB JPEG gets 413 with `limitBytes: 5242880`.
- [ ] A PNG renamed `.jpg` with `Content-Type: image/jpeg` is accepted as PNG; a HEIC file gets 415.
- [ ] A 41-page PDF gets 413; a 40-page PDF of 24 MB passes; a 26 MB PDF gets 413.
- [ ] A 30 MB request is cut off without reading the rest (memory stays flat in a test that streams 100 MB).
- [ ] No file bytes appear in logs or in `gw:idem:*` values (test inspects Redis after an upload with a key).
- [ ] A fake kitchen receives the exact bytes over gRPC with a 26 MB max message size.

## Notes

- Decided (Q022): own-flyer PDFs may have up to **40 pages** and 25 MB (typical chain flyers are 16–40 pages);
  OFF-10-01 uses the same limits.
