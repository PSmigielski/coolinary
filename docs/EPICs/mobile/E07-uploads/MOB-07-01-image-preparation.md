# MOB-07-01 Image preparation

| | |
|---|---|
| Service | mobile |
| Epic | [MOB-07 Uploads](README.md) |
| Release | R1 (fridge, shelf, label), R2 (flyer pages), R3 (receipt) |
| Size | M |
| Status | todo |
| Depends on | — |
| Unblocks | MOB-07-02, MOB-07-03 |

## Goal

Every photo sent for AI reading is small, fast to upload, cheap to process and free of location data, while still
readable for the model.

## Business rules

1. Use `expo-image-manipulator`'s current API (`ImageManipulator.manipulate(uri)` → `resize` → `renderAsync()` →
   `saveAsync({ format: SaveFormat.JPEG, compress })`; the old `manipulateAsync` is deprecated). It is in Expo Go.
2. Limits per kind, in one constant `IMAGE_RULES`:

   | Kind | Long edge (px) | JPEG quality | Max size after encoding | Images per request |
   |---|---|---|---|---|
   | `fridge`, `shelf` | 1568 | 0.75 | 1.5 MB | 1 |
   | `label` | 1568 | 0.8 | 1.5 MB | 1 |
   | `receipt` (R3) | 1568 per part | 0.8 | 1.5 MB each | 1–3 parts (top to bottom) |
   | `flyer_page` (R2) | 1568 | 0.75 | 1.5 MB each | 1–12 pages |

   1568 px is the long edge above which Claude scales images down anyway, so larger images only cost upload time.
   Decided (Q114): a long receipt is photographed in up to 3 parts, so small print stays readable.
3. Images smaller than the long edge are not enlarged. If the encoded file is above the max size, re-encode once at
   quality 0.55; still too large → fail with `image_too_large`.
4. Location and other metadata must not leave the phone: the re-encoded JPEG is checked to carry no EXIF GPS block.
   If the manipulator keeps EXIF on a platform, strip it before upload (acceptance check below decides).
5. Pickers are frontend work; this module accepts any local `file://` or `content://` URI from the camera or the
   library and returns prepared files.
6. Prepared files are written by `saveAsync` to the cache directory; they are deleted after the upload finishes,
   fails or is cancelled (MOB-07-02). The original photo the picker gave is not touched; camera captures made only
   for the scan are deleted too.
7. Images are never written to AsyncStorage, the query cache or the outbox (README D6).

## Contract

```ts
type ImageKind = "fridge" | "shelf" | "label" | "receipt" | "flyer_page";
interface PreparedImage { uri: string; width: number; height: number; bytes: number; mime: "image/jpeg" }
export function prepareImage(sourceUri: string, kind: ImageKind): Promise<PreparedImage>;
export function discardPrepared(images: PreparedImage[]): void;
```

## Data

Temporary files in `Paths.cache` only.

## Errors

| Case | Behaviour |
|---|---|
| Unreadable source | `image_unreadable` |
| Too large after re-encode | `image_too_large` |

## Acceptance criteria

- [ ] A 12 MP fridge photo becomes ≤ 1568 px on the long edge and ≤ 1.5 MB.
- [ ] A photo taken with location tagging on has no GPS EXIF after preparation (checked with an EXIF reader on both
      platforms).
- [ ] After an upload, no prepared file remains in the cache directory.
