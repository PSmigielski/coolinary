# AI-05-01 ExtractItemsFromText on Claude

| | |
|---|---|
| Service | ai |
| Epic | [AI-05 Text to items](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | AI-01-03, AI-01-06, AI-01-07, AI-01-09, AI-03-01 |
| Unblocks | AI-06-01 (shared item rules, R1); KIT-08-03 (typed and voice intake, R3) |

## Goal

Replace the stub's regex parser in production with Haiku 5.5, keeping the existing RPC and adding the hints kitchen
needs to propose storage places and match names.

## Business rules

1. Job `text_to_items`, model and effort from the registry (Haiku 5.5, `low`). Prompt `text_to_items/v1`.
2. Input: `text` 1–2,000 characters after trimming (existing limit), `language` (BCP 47, default `en`). The text is
   untrusted (AI-01-07). Accepted languages are any; prompts and golden sets cover `pl`, `de`, `en`.
3. Output: 0–50 items. Each item:
   - `name`: the food as the user named it, in the user's language, singular where the language allows
     ("papryka", not "papryki"), 1–80 characters, without quantity words.
   - `quantity`: number > 0, or absent when no amount was said ("some cheese").
   - `unit`: `GRAM`, `MILLILITRE`, `PIECE` or `UNSPECIFIED`. Conversions: kg → ×1000 g, dag/dkg → ×10 g, l → ×1000 ml,
     cl → ×10 ml; "dozen" / "tuzin" / "Dutzend" → 12 pieces; "half a kilo" / "pół kilo" / "ein halbes Kilo" → 500 g.
     A count of packs ("two packs of butter") → `PIECE` with `package_hint=PACK`.
   - `package_hint`: `NONE`, `PACK`, `BOTTLE`, `CAN`, `JAR`, `BAG`, `BOX`, `LOOSE`.
   - `storage_hint`: `UNSPECIFIED`, `FRIDGE`, `FREEZER`, `PANTRY` — only when the text says so ("frozen peas" →
     `FREEZER`).
   - `confidence` 0–1.
4. Same name and unit said twice ("eggs … and 6 more eggs") are merged by summing quantities; if one has no quantity,
   the merged item has none.
5. Ranges: quantity in grams or ml ≤ 100,000; pieces ≤ 1,000. Out of range → quantity dropped (item kept,
   `confidence` capped at 0.5).
6. Non-food (detergent, batteries, medicines) is left out and counted in `ignored_count`.
7. Spoken input from phone speech recognition often lacks punctuation; "and"/"i"/"und" and pauses separate items.
8. Cache: 7 days, key per AI-01-09 (household not part of the key).
9. Tracing: `chars`, `items`, `language` only.
10. The stub provider stays the dev default and keeps its current tests.

## Contract

`contracts/coolinary/ai/v1/ai.proto` (extend the existing messages; field numbers continue):

```proto
message ExtractItemsFromTextResponse {
  repeated ExtractedItem items = 1;
  string model = 2;          // kept for compatibility; same as call.model
  bool cached = 3;           // kept for compatibility
  int32 ignored_count = 4;
  CallInfo call = 5;
}

message ExtractedItem {
  string name = 1;
  coolinary.common.v1.Quantity quantity = 2;   // absent when no amount was given
  float confidence = 3;
  PackageHint package_hint = 4;
  StorageHint storage_hint = 5;
}

enum PackageHint { PACKAGE_HINT_UNSPECIFIED = 0; PACKAGE_HINT_PACK = 1; PACKAGE_HINT_BOTTLE = 2;
  PACKAGE_HINT_CAN = 3; PACKAGE_HINT_JAR = 4; PACKAGE_HINT_BAG = 5; PACKAGE_HINT_BOX = 6; PACKAGE_HINT_LOOSE = 7; }
enum StorageHint { STORAGE_HINT_UNSPECIFIED = 0; STORAGE_HINT_FRIDGE = 1; STORAGE_HINT_FREEZER = 2; STORAGE_HINT_PANTRY = 3; }
```

`PackageHint` and `StorageHint` are shared with AI-06 and AI-07.

Output schema (pydantic, outline): `{items: [{name: str(1..80), quantity: float|null, unit: enum, package_hint:
enum, storage_hint: enum, confidence: float(0..1)}] (max 50), ignored_count: int}`.

## Data

No schema change.

## Errors

| Case | gRPC status |
|---|---|
| Empty text or over 2,000 characters | `INVALID_ARGUMENT`, `ai_bad_request` (existing) |
| Others | AI-01-04 |

## Acceptance criteria

- [ ] `"kupiłam dwie papryki, pół kilo mąki i mrożony groszek"` → papryka 2 PIECE; mąka 500 GRAM; groszek, no quantity, `FREEZER`.
- [ ] `"2 Packungen Butter und 1,5 l Milch"` → Butter 2 PIECE `PACK`; Milch 1500 MILLILITRE.
- [ ] `"eggs and six eggs"` → one item, eggs 6 PIECE.
- [ ] Response carries `call.prompt_version = "text_to_items/v1"` and a non-zero cost on a live call.
- [ ] Golden set thresholds pass (AI-05-02).

## Notes

- Decided (Q219): typed and voice intake ship in release 3 (KIT-08-03). This task stays in R1 because photo intake
  (AI-06-01) reuses its item schema, quantity ranges and validation; the RPC is built and tested in R1 and called
  by kitchen from R3.
