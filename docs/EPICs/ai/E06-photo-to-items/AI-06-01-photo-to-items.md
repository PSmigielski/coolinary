# AI-06-01 ExtractItemsFromPhoto

| | |
|---|---|
| Service | ai |
| Epic | [AI-06 Photo to items](README.md) |
| Release | R1 |
| Size | L |
| Status | todo |
| Depends on | AI-05-01, AI-01-08 |
| Unblocks | KIT-08 (photo intake) |

## Goal

Read 1–4 photos of a storage place and list the food items visible, with honest estimates of amounts, so the user
only has to correct, not type.

## Business rules

1. Job `photo_to_items`, Haiku 5.5, effort `medium`. Prompt `photo_to_items/v1`.
2. Input:
   - 1–4 images, each JPEG, PNG or WebP, at most 3 MiB and at most 2576 px on the long edge (the phone sends
     2048 px). Total at most 10 MiB. Type is checked from magic bytes, not the declared type.
   - `context`: `FRIDGE`, `FREEZER`, `PANTRY`, `SHELF`, `UNPACKED_BAG`.
   - `language` for item names.
3. Output: 0–60 items. Each item:
   - `name`: generic food name in the requested language, specific enough to match the catalogue
     ("natural yoghurt", "red pepper", "smoked ham"), 1–80 chars.
   - `brand`: only if clearly readable on the pack, else empty.
   - `quantity`: estimate with unit — count of pieces for countable items (`PIECE`), pack size in g/ml only when
     printed and readable; else absent. `quantity_confidence` 0–1 separately from `confidence`.
   - `package_hint` and `storage_hint` (from the context: FRIDGE → `FRIDGE`, FREEZER → `FREEZER`, PANTRY and SHELF
     → `PANTRY`; `UNPACKED_BAG` → hint only if obvious, e.g. frozen goods).
   - `opened`: true when visibly opened (cut cheese, open jar).
   - `use_by_visible`: ISO date only when a date is clearly printed and readable on the pack; never guessed from how
     food looks.
   - `image_index` (0-based) of the photo where it is best seen.
   - `confidence` 0–1.
4. Each physical item is listed once even if seen in two photos (prompt rule). ai merges entries with the same
   `name`, `brand` and `image_index` by summing counts.
5. Non-food items (cleaning products, medicines, cosmetics) are left out and counted in `ignored_count`. People,
   faces, documents and screens are never described; if a photo shows mainly a person or a document, its items
   are empty and `unusable_images` includes its index.
6. A photo too dark or blurred to read gives `unusable_images` with that index and no invented items.
7. Validation (AI-01-06): names non-empty; `image_index` < number of images; dates valid and within today − 30 days
   to today + 3 years; quantities in the ranges of AI-05-01 rule 5.
8. Photos are never written to disk, Redis, logs or traces; the cache key uses the SHA-256 of each image (TTL 1 day)
   so a retry of the same upload costs nothing (D6, AI-01-09).
9. Tracing: `images`, `items`, `language`, context.

## Contract

```proto
rpc ExtractItemsFromPhoto(ExtractItemsFromPhotoRequest) returns (ExtractItemsFromPhotoResponse);

enum PhotoContext {
  PHOTO_CONTEXT_UNSPECIFIED = 0;
  PHOTO_CONTEXT_FRIDGE = 1;
  PHOTO_CONTEXT_FREEZER = 2;
  PHOTO_CONTEXT_PANTRY = 3;
  PHOTO_CONTEXT_SHELF = 4;
  PHOTO_CONTEXT_UNPACKED_BAG = 5;
}

message InputImage {
  bytes data = 1;
  string media_type = 2;   // image/jpeg | image/png | image/webp
}

message ExtractItemsFromPhotoRequest {
  repeated InputImage images = 1;   // 1-4
  PhotoContext context = 2;
  string language = 3;
  string household_id = 4;
}

message PhotoItem {
  string name = 1;
  string brand = 2;
  coolinary.common.v1.Quantity quantity = 3;
  float quantity_confidence = 4;
  PackageHint package_hint = 5;
  StorageHint storage_hint = 6;
  bool opened = 7;
  string use_by_visible = 8;        // YYYY-MM-DD or empty
  int32 image_index = 9;
  float confidence = 10;
}

message ExtractItemsFromPhotoResponse {
  repeated PhotoItem items = 1;
  int32 ignored_count = 2;
  repeated int32 unusable_images = 3;
  CallInfo call = 4;
}
```

`InputImage` is shared by AI-07, AI-13 and AI-16.

## Data

No schema change.

## Errors

| Case | gRPC status |
|---|---|
| 0 or more than 4 images, image over 3 MiB, bad type, total over 10 MiB | `INVALID_ARGUMENT`, `ai_bad_request` |
| All images unusable | OK with empty items and all indexes in `unusable_images` (kitchen shows "photo too dark") |
| Others | AI-01-04 |

## Acceptance criteria

- [ ] A test fridge photo with 12 known items returns ≥ 10 of them with plausible names; a printed date on a yoghurt lid is returned as `use_by_visible`.
- [ ] A photo of a person returns no items and its index in `unusable_images`.
- [ ] The same request twice: second one is cached, cost 0.
- [ ] Captured logs and traces contain no image bytes and no item names.
- [ ] Golden set thresholds pass (AI-06-02); cost per 2-photo call ≈ $0.001.
