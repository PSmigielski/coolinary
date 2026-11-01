# KIT-08-02 Photo intake (fridge or shelf)

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-08 AI intake](README.md) |
| Release | R1 |
| Size | M |
| Depends on | KIT-08-01, AI-06 (`ExtractItemsFromPhoto`) |
| Unblocks | GW-06 (`POST /v1/intake/photo`) |

## Goal

A photo of the fridge or a shelf becomes an item list the user confirms (features.md §2, AI per request).

## Business rules

1. 1–4 images per call (several shelves of one fridge count as one scan); each JPEG, PNG, WebP or HEIC,
   10 KB–5 MB; all together ≤ 20 MB. The phone resizes before upload (MOB-07).
2. The kitchen gRPC server accepts messages up to 26 MB (the gateway client limit, GW-01-05).
3. Images go to `ai` `ExtractItemsFromPhoto` in one call with the language and the storage hint; deadline 45 s
   (inside the gateway's 60 s).
4. Images are never stored, never logged, never written to any cache by kitchen (D6). Only the AI result passes
   through KIT-08-01.
5. AI result → KIT-08-01 pipeline. Zero recognised items is a success: `ai_used = true`, no proposals.
6. AI unavailable (down, deadline, spend limit, refusal) → `UNAVAILABLE` with reason `ai_unavailable`. There is no
   non-AI fallback for photos; the app offers typing and barcodes instead.
7. Free and Plus are treated the same here; the Free weekly limit of 2 scans is enforced by the gateway (GW-03).
8. A photo item that names a branded product does not get an EAN; it is matched to an ingredient only.

## Contract

`contracts/coolinary/kitchen/v1/intake.proto`

```proto
rpc IntakeFromPhoto(IntakeFromPhotoRequest) returns (IntakeResponse);

message IntakeImage {
  bytes data = 1;
  string media_type = 2;          // image/jpeg, image/png, image/webp, image/heic
}

message IntakeFromPhotoRequest {
  string household_id = 1;
  repeated IntakeImage images = 2;  // 1–4
  StoragePlace storage = 3;         // where the photo was taken; optional
  string language = 4;
}
```

## Data

No schema change.

## Errors

| Case | gRPC status | Behaviour |
|---|---|---|
| 0 or > 4 images, unknown media type, image < 10 KB or > 5 MB, total > 20 MB | `INVALID_ARGUMENT` | |
| ai unavailable | `UNAVAILABLE` | reason `ai_unavailable` |
| ai says the photo shows no food at all | — | success, `ai_used = true`, empty list |

## Acceptance criteria

- [ ] A fixed test photo of a fridge (golden set, AI-03) gives proposals with storage `fridge` when the hint is
      fridge, each with a use-by date where the catalogue has a shelf life.
- [ ] Two images in one call make exactly one `ExtractItemsFromPhoto` call.
- [ ] With `ai` stopped, the RPC returns `UNAVAILABLE` / `ai_unavailable` within 2 s (connection refused), and
      nothing is stored.
- [ ] A 6 MB image is rejected without calling `ai`.
- [ ] No image bytes appear in logs, Redis or Postgres (checked in a test with a marker byte sequence).
