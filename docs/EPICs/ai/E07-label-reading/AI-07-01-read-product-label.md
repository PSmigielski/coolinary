# AI-07-01 ReadProductLabel

| | |
|---|---|
| Service | ai |
| Epic | [AI-07 Product label reading](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | AI-06-01 |
| Unblocks | KIT-05 (naming unknown products) |

## Goal

Read what is printed on a pack so naming an unknown product takes one tap to confirm.

## Business rules

1. Job `read_label`, Haiku 5.5, effort `low`. Prompt `read_label/v1`.
2. Input: 1–2 images (front, back) with the limits of AI-06-01 rule 2; `language` (the user's); optional `ean`
   (13/8 digits) as context only — the model is told the code, never asked to output it.
3. Output fields, each with its own confidence 0–1:
   - `product_name`: as printed, in the pack's language, 1–120 chars ("Jogurt naturalny 2% tłuszczu").
   - `generic_name`: the food in the user's language for catalogue matching ("natural yoghurt"), 1–80 chars.
   - `brand`: as printed, 0–60 chars.
   - `pack_size`: `Quantity` (g, ml or pieces), only when printed; multipacks give `unit_count` (e.g. 4 × 125 g →
     `pack_size` 125 g, `unit_count` 4).
   - `best_before`: ISO date only when printed and readable.
   - `storage_hint`: from printed instructions ("store below 6 °C" → `FRIDGE`, "keep frozen" → `FREEZER`), else
     `UNSPECIFIED`; `refrigerate_after_opening` boolean.
   - `pack_language`: BCP 47 code of the label text.
4. Validation: `pack_size` g/ml in 1–25,000, pieces 1–200; `unit_count` 1–48; date range as AI-06-01 rule 7; at
   least one of `product_name` / `generic_name` non-empty, else the result is `unreadable=true` with all fields
   empty (not an error).
5. Nutrition values printed on the label are **not** read in v1 (see Notes). Pack nutrition comes from Open Food
   Facts where available (KIT-06).
6. Cache: 30 days by image hashes; tracing: `images`, language.

## Contract

```proto
rpc ReadProductLabel(ReadProductLabelRequest) returns (ReadProductLabelResponse);

message ReadProductLabelRequest {
  repeated InputImage images = 1;   // 1-2
  string language = 2;
  string ean = 3;                   // optional context
  string household_id = 4;
}

message LabelField { string value = 1; float confidence = 2; }

message ReadProductLabelResponse {
  LabelField product_name = 1;
  LabelField generic_name = 2;
  LabelField brand = 3;
  coolinary.common.v1.Quantity pack_size = 4;
  float pack_size_confidence = 5;
  int32 unit_count = 6;
  string best_before = 7;            // YYYY-MM-DD or empty
  StorageHint storage_hint = 8;
  bool refrigerate_after_opening = 9;
  string pack_language = 10;
  bool unreadable = 11;
  CallInfo call = 12;
}
```

## Data

No schema change.

## Errors

| Case | gRPC status |
|---|---|
| 0 or more than 2 images, size or type invalid | `INVALID_ARGUMENT`, `ai_bad_request` |
| Nothing readable | OK, `unreadable=true` |
| Others | AI-01-04 |

## Acceptance criteria

- [ ] A Polish yoghurt label returns product name, brand, 150 g, `FRIDGE`; a German 4 × 125 g multipack returns 125 g and `unit_count` 4.
- [ ] A blank or blurred photo returns `unreadable=true`.
- [ ] No output field contains a calorie number (calorie guard).

## Notes

- Decided (Q110): the printed nutrition table is not read in R1. The label job returns name, brand and pack size
  only; calories keep coming from nutrition tables or Open Food Facts pack values.
