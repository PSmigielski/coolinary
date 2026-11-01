# AI-09-01 translate_ingredients

| | |
|---|---|
| Service | ai |
| Epic | [AI-09 Catalogue and library batch jobs](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | AI-04-03 |
| Unblocks | KIT-02 (catalogue translations) |

## Goal

Give every catalogue ingredient a name, synonyms and search forms in each launch language, so autocomplete and
matching work in Polish, German and English (features.md §2, "Ingredient catalogue translated and matched across
languages", once for everyone).

## Business rules

1. Batch kind `translate_ingredients`, Haiku 5.5, effort `low`, prompt `translate_ingredients/v1`.
2. Item input: `ingredient_id`, the names it already has (`en` always; `pl`/`de` when known), `category`,
   `typical_state` (raw, dry, cooked, canned, frozen, fresh), `target_languages` (1–5 of `pl`, `de`, `en`, `fr`,
   `es`, `it`, `nl`, `cs`).
3. Output per target language:
   - `name`: the everyday shop/kitchen name, singular, 1–60 chars; lower case except where the language
     capitalises nouns (German).
   - `synonyms`: 0–5 other common names, including regional ones ("ziemniak" / "kartofel"; "Kartoffel" /
     "Erdapfel").
   - `search_forms`: 0–6 inflected or plural forms people type ("papryki", "paprykę"; "Tomaten").
   - `confidence` 0–1.
4. Rules: no brand names; no translation that changes the food (e.g. "twaróg" is not "cottage cheese" — keep the
   specific food and its local name; when no target-language word exists, use the source name and lower
   confidence to ≤ 0.6).
5. Checks in ai: lengths, characters (letters, spaces, hyphens, apostrophes only), language codes in the allowed set,
   no duplicates between name, synonyms and search forms.
6. **Caller checks** (KIT-02): one result per requested target language; names unique within a language in its
   catalogue or flagged for review.
7. Results with confidence < 0.8 go to admin review; others are published and spot-checked (10% sample) by a native
   speaker before a language wave launches.

## Contract

```proto
message TranslateIngredientInput {
  string ingredient_id = 1;
  map<string, string> known_names = 2;   // language -> name; "en" required
  string category = 3;
  string typical_state = 4;
  repeated string target_languages = 5;
}

message IngredientNames {
  string language = 1;
  string name = 2;
  repeated string synonyms = 3;
  repeated string search_forms = 4;
  float confidence = 5;
}

message IngredientTranslation {
  string ingredient_id = 1;
  repeated IngredientNames names = 2;
}
```

## Data

No schema change.

## Errors

Per item, as AI-04-03.

## Acceptance criteria

- [ ] "potato" → pl ziemniak (synonym kartofel, forms ziemniaki), de Kartoffel (synonym Erdapfel).
- [ ] "twaróg" (en "quark-style curd cheese") → de "Quark", not "Hüttenkäse".
- [ ] Golden set thresholds pass (AI-09-06).
