# KIT-03-01 Name normaliser per language

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-03 Name matching](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | — |
| Unblocks | KIT-02-04, KIT-03-02, KIT-04-07 (pack size parsing), KIT-04-11, KIT-05-02 |

## Goal

One deterministic function turns a raw name into a matching key and pulls out any pack size, so "Mleko UHT 3,2% 1l",
"MLEKO UHT 1L" and "mleko uht" all land on the same key.

## Business rules

1. `Normalize(language, raw, brand?) → (key, quantity?, version)`. Pure function in `internal/kitchen/names`,
   `NORMALIZER_VERSION = 1`. Any rule change raises the version.
2. Steps, in order:
   1. Unicode NFKC; trim; collapse whitespace to one space.
   2. Lowercase with the language's rules (Polish, German, English).
   3. Pull out the first pack size and remove it from the text:
      - multipack `N x V U` (`x` or `×`, e.g. `6x0,5l`, `4 x 125 g`) → `N` units of `V U`;
      - amount `V U` with `U` in `kg g mg l ml cl dl` (comma or dot decimals);
      - count `N szt` / `N stk` / `N stück` / `N pcs` / `N st`.
      Convert to grams, millilitres or pieces (`kg` × 1000, `l` × 1000, `cl` × 10, `dl` × 100, `mg` ÷ 1000).
   4. Remove percentages (`3,2%`, `20 %`).
   5. Remove the brand's words if `brand` is given.
   6. Replace punctuation and hyphens with spaces; drop tokens that are only digits.
   7. Remove stop words per language, from a versioned list in `internal/kitchen/names/stopwords/<lang>.txt`
      (launch lists include: en `fresh organic bio premium classic pack the a of`; pl `świeży świeża świeże bio eko
      premium klasyczny op opak opakowanie`; de `frisch bio premium klassisch packung pckg`).
   8. Fold diacritics: `ą ć ę ł ń ó ś ź ż` → `a c e l n o s z z`; `ä ö ü` → `a o u`; `ß` → `ss`. Receipts often
      print without diacritics, so the key is folded.
   9. Remove repeated tokens, keep the order.
3. An empty key after normalisation is an error `empty_name`.
4. Raw names over 120 characters are an error `too_long`. Names containing `@`, `http`, or 7+ consecutive digits are
   an error `not_a_name` (protects the match store from personal data).
5. Unknown languages use the English stop words and keep their language code.
6. Decided (Q052): German compound words are not split ("Hähnchenschenkel" stays one token); Haiku resolves
   compounds in the AI path when text matching is unsure.

## Contract

No RPC (library).

## Data

No schema change.

## Errors

Library errors `empty_name`, `too_long`, `not_a_name`; RPCs that use it report them per name.

## Acceptance criteria

- [ ] At least 30 table tests, including:
      `de` "Hähnchenschenkel 1kg" → key `hahnchenschenkel`, 1000 g;
      `pl` "Mleko UHT 3,2% 1l" → `mleko uht`, 1000 ml;
      `pl` "MLEKO UHT 1L" → `mleko uht`, 1000 ml;
      `en` "Chicken thighs 1 kg" → `chicken thighs`, 1000 g;
      `pl` "Woda mineralna 6x0,5l" → `woda mineralna`, 6 × 500 ml;
      `de` "Bio Eier 10 Stück" → `eier`, 10 pieces;
      "call 600123456" → `not_a_name`.
- [ ] The same input always gives the same output (property test over 1,000 random names).
