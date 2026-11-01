# MOB-07 Uploads

Goal: photos and files for AI reading (fridge and shelf photos, product labels, receipts, flyers) are prepared on the
phone (≤ 1568 px, JPEG, no location data), uploaded with progress and cancel, retried without using a second free
scan, and never kept on the phone.

## Flow

```
camera / library / document picker (frontend) ──> prepareImage(kind) ──> cache dir JPEG ≤ 1.5 MB, no GPS EXIF
   ──> plus gate + online ──> multipart upload (Bearer, Idempotency-Key, 120 s, 1 retry)
   ──> proposals / offers ──> delete temporary files
```

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [MOB-07-01](MOB-07-01-image-preparation.md) | Image preparation | R1–R3 | M | — |
| [MOB-07-02](MOB-07-02-upload-transport.md) | Upload transport | R1 | M | MOB-07-01, MOB-02-01, MOB-02-02, GW-06 |
| [MOB-07-03](MOB-07-03-flyer-inputs.md) | Flyer inputs: photos, PDF, link | R2 | S | MOB-07-01, MOB-07-02, GW-10, OFF-10 |
