# KIT-08-03 Text and voice intake with a non-AI fallback

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-08 AI intake](README.md) |
| Release | R3 |
| Size | M |
| Status | todo |
| Depends on | KIT-08-01, AI-05 (`ExtractItemsFromText`), KIT-02 (`SearchIngredients`) |
| Unblocks | GW-06 (`POST /v1/intake/text`) |

## Goal

A typed or spoken sentence ("kupiłam jajka, dwie papryki i trochę sera") becomes items to confirm. Voice is turned
into text on the phone (D7), so kitchen only ever receives text. When AI is not available, a simple rule-based
parser still gives a usable list.

## Business rules

1. Text is 1–2,000 characters after trimming (the `ai` limit). `kind` is `text` or `voice`; it is used for metrics
   only.
2. With AI available: `ExtractItemsFromText` (Haiku 5.5, AI-05) with the language; deadline 15 s; result →
   KIT-08-01.
3. **Fallback without AI** (ai down, deadline, spend limit, refusal): kitchen's own parser, `ai_used = false`,
   `rules_fallback = true`:
   1. split on commas, semicolons, new lines and the words "and", "i", "oraz", "und" (whole words, any case);
   2. at most 30 parts; each trimmed; empty parts dropped;
   3. a leading number ("2", "1,5") or number word (one–twelve in English, Polish, German; "a", "an", "ein",
      "eine" = 1) becomes the quantity in pieces;
   4. "500 g", "500g", "1 kg", "250 ml", "1,5 l" become grams or millilitres;
   5. filler words are removed: "bought", "some", "a bit of", "kupiłem", "kupiłam", "trochę", "gekauft", "etwas",
      "ich habe";
   6. the rest is looked up with `SearchIngredients` (KIT-02) in the language; the top hit is used when its
      similarity is ≥ 0.6, else the part stays unmatched;
   7. every fallback proposal has confidence 0.4 (so matched ones are unticked under KIT-08-01 rule 8).
4. The text is user free text: never logged, never stored, never put in metrics.
5. Decided (Q219): typed-sentence intake ships in release 3, together with voice, as features.md groups them; the
   gateway route is R3 (GW-06).
6. Free and Plus: decided (Q112), every AI reading request counts as one Free scan, a typed or spoken list included;
   the gateway counts it (GW-03-01) when `ai_used` is true. Kitchen applies no plan rule.

## Contract

`contracts/coolinary/kitchen/v1/intake.proto`

```proto
rpc IntakeFromText(IntakeFromTextRequest) returns (IntakeResponse);

message IntakeFromTextRequest {
  string household_id = 1;
  string text = 2;                // 1–2,000 characters
  IntakeKind kind = 3;            // TEXT or VOICE
  StoragePlace storage = 4;       // optional hint
  string language = 5;
}
```

## Data

No schema change.

## Errors

| Case | gRPC status | Behaviour |
|---|---|---|
| Empty text or > 2,000 characters | `INVALID_ARGUMENT` | |
| kind not TEXT/VOICE | `INVALID_ARGUMENT` | |
| ai unavailable | — | rule-based fallback, OK response |

## Acceptance criteria

- [ ] With `ai` up: "bought eggs, two peppers, some cheese" → eggs, peppers 2 pcs, cheese (estimated quantity).
- [ ] With `ai` stopped: "2 jajka, 500 g mąki i mleko" → eggs 2 pcs, flour 500 g, milk, all unticked, `ai_used = false`,
      `rules_fallback = true`.
- [ ] "zwei Paprika und 1,5 l Milch" (German, fallback) → peppers 2 pcs, milk 1,500 ml.
- [ ] A 2,001-character text is rejected without calling `ai`.
- [ ] No part of the text appears in logs (test with a marker word).
