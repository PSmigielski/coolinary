# MOB-07-02 Upload transport

| | |
|---|---|
| Service | mobile |
| Epic | [MOB-07 Uploads](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | MOB-07-01, MOB-02-01, MOB-02-02, GW-06 (intake routes) |
| Unblocks | MOB-03-05, MOB-07-03 |

## Goal

Photos reach the gateway with progress, can be cancelled, are retried safely, and count against the AI quota once.

## Business rules

1. One image: `expo-file-system`'s `File.createUploadTask(url, { uploadType: UploadType.MULTIPART, fieldName,
   mimeType, parameters, headers, onProgress })` and `task.uploadAsync()` (the legacy `createUploadTask` is
   deprecated). Progress = `bytesSent / totalBytes` (includes multipart framing, fine for a progress bar).
2. Several images (receipt parts, flyer pages): `fetch` from `expo/fetch` with a `FormData` that appends each `File`;
   progress is shown per request, not per byte.
   Decided (Q115): one multipart request carries all pages (one AI call, one quota count).
3. Headers are the same as MOB-02-01 (Bearer, `x-request-id`, `Accept-Language`, `Idempotency-Key`, `x-client`);
   upload code gets them from the middleware's header builder, because upload tasks bypass openapi-fetch.
4. Form fields (assumed GW-06/GW-10): `image` for single uploads, `images` repeated for parts, `pages` repeated for
   flyer pages, plus plain fields (`kind`, `ean`, `shopId`). Responses are parsed and typed with the `schema.d.ts`
   aliases.
5. Timeout 120 s per request. On timeout or network failure: one automatic retry after 3 s with the same
   `Idempotency-Key` (the gateway then counts the photo scan once, GW-03); then the error goes to the caller.
6. Cancel: `task.cancel()` (or aborting the fetch) stops the upload; the prepared files are discarded. A cancel after
   the server has the whole body may still use a free scan; the UI says so only when it happens
   (`X-Quota-Remaining` decreased).
7. Before starting: `usePlusGate("photo_scan")` must allow it (MOB-02-05), and the phone must be online. Uploads are
   never queued offline.
8. Prepared files are deleted in a `finally` block (MOB-07-01 rule 6).

## Contract

```ts
export interface UploadHandle<T> { progress: () => number; cancel: () => void; result: Promise<T> }
export function uploadImages<T>(route: "/v1/intake/photo" | "/v1/intake/label" | "/v1/intake/receipt" | "/v1/flyers", images: PreparedImage[], fields: Record<string, string>, idempotencyKey: string): UploadHandle<T>;
```

## Data

None persisted.

## Errors

Mapped by MOB-02-02; plus `upload_cancelled`.

## Acceptance criteria

- [ ] A fridge photo upload on a throttled network (3G profile) shows progress from 0 to 1 and returns proposals.
- [ ] Dropping the connection mid-upload retries once; the gateway counts one photo scan.
- [ ] Cancelling at 50 % leaves no temporary file and no stock change.
