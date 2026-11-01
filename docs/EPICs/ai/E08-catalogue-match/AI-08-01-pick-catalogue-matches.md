# AI-08-01 PickCatalogueMatches (sync)

| | |
|---|---|
| Service | ai |
| Epic | [AI-08 Catalogue match pick](README.md) |
| Release | R1 |
| Size | M |
| Depends on | AI-01-03, AI-01-06, AI-01-09 |
| Unblocks | KIT-03 (MatchNames) |

## Goal

Haiku chooses the right catalogue ingredient among a few candidates, or none, so recipes, stock and offers all point
to the same ingredient.

## Business rules

1. Job `pick_matches`, Haiku 5.5, effort `low`. Prompt `pick_matches/v1`.
2. Input: 1–50 `MatchItem`s per call. Each has `ref` (caller's id, ≤ 64 chars, unique in the call), `name`
   (normalised by kitchen, 1–160 chars), `language`, `context` (`PRODUCT`, `RECEIPT_LINE`, `OFFER`,
   `RECIPE_IMPORT`, `USER_TEXT`), optional `brand`, `pack_size`, `category_hint`, and 1–10 `candidates`
   (`ingredient_id`, canonical name in the item's language, English name, up to 5 synonyms, category).
3. Picking rules (prompt and eval):
   - Pick the candidate that **is** the same food. A brand or own-brand product maps to its generic ingredient
     ("Łaciate mleko 3,2%" → milk 3.2%).
   - Prefer the most specific candidate that is still the same food ("smoked salmon" → smoked salmon, not salmon,
     when both exist).
   - If only a broader candidate exists, pick it with reason `GENERIC_FORM` and confidence at most 0.8.
   - A composite product (ready meal, pizza, sandwich) is `none` unless a candidate is that exact product type.
   - Receipt abbreviations ("SER GOUDA PLAST.") are read as the food they name.
   - When unsure between two candidates, return `none` rather than guess (confidence < 0.5 means none).
4. Output per item: `ref`, `ingredient_id` (empty for none), `confidence` 0–1, `reason` (`EXACT`, `SYNONYM`,
   `BRANDED_VARIANT`, `GENERIC_FORM`, `ABBREVIATION`, `NONE`).
5. Validation in ai (sync mode has the inputs): every `ref` answered exactly once; `ingredient_id` is one of that
   item's candidates or empty; `reason=NONE` ⇔ empty id. A wrong id is not fixable → repair (AI-01-06).
6. Kitchen's thresholds (KIT-03, stated here so both sides agree): ≥ 0.90 auto-accept, 0.60–0.89 accept and queue
   for admin review, < 0.60 treat as none.
7. Cache: 30 days per item (key over name, language, context and the sorted candidate ids), so a call with 50 names
   where 40 are cached sends only 10 to Claude.
8. Tracing: `items`, `candidates` counts.

## Contract

`contracts/coolinary/ai/v1/matching.proto`:

```proto
rpc PickCatalogueMatches(PickCatalogueMatchesRequest) returns (PickCatalogueMatchesResponse);

enum MatchContext {
  MATCH_CONTEXT_UNSPECIFIED = 0;
  MATCH_CONTEXT_PRODUCT = 1;
  MATCH_CONTEXT_RECEIPT_LINE = 2;
  MATCH_CONTEXT_OFFER = 3;
  MATCH_CONTEXT_RECIPE_IMPORT = 4;
  MATCH_CONTEXT_USER_TEXT = 5;
}

message MatchCandidate {
  string ingredient_id = 1;
  string name = 2;              // in the item's language
  string name_en = 3;
  repeated string synonyms = 4; // up to 5
  string category = 5;
}

message MatchItem {
  string ref = 1;
  string name = 2;
  string language = 3;
  MatchContext context = 4;
  string brand = 5;
  coolinary.common.v1.Quantity pack_size = 6;
  string category_hint = 7;
  repeated MatchCandidate candidates = 8;   // 1-10
}

enum MatchReason {
  MATCH_REASON_UNSPECIFIED = 0; MATCH_REASON_EXACT = 1; MATCH_REASON_SYNONYM = 2;
  MATCH_REASON_BRANDED_VARIANT = 3; MATCH_REASON_GENERIC_FORM = 4; MATCH_REASON_ABBREVIATION = 5; MATCH_REASON_NONE = 6;
}

message MatchPick {
  string ref = 1;
  string ingredient_id = 2;     // empty = none
  float confidence = 3;
  MatchReason reason = 4;
}

message PickCatalogueMatchesRequest { repeated MatchItem items = 1; }
message PickCatalogueMatchesResponse { repeated MatchPick picks = 1; CallInfo call = 2; }
```

## Data

No schema change.

## Errors

| Case | gRPC status |
|---|---|
| 0 or > 50 items, an item with 0 or > 10 candidates, duplicate `ref` | `INVALID_ARGUMENT`, `ai_bad_request` |
| Others | AI-01-04 |

## Acceptance criteria

- [ ] "Hähnchenschenkel 1kg" with candidates chicken thigh, chicken breast, chicken (whole) → chicken thigh, `EXACT` or `SYNONYM`, ≥ 0.9.
- [ ] "Pizza Margherita mrożona" with only cheese/tomato/flour candidates → none.
- [ ] A fake model answer with an id not in the candidates triggers one repair.
- [ ] A 50-item call with 40 cached items sends 10 items to the provider (asserted on the request).
- [ ] Golden set thresholds pass (AI-08-03).
