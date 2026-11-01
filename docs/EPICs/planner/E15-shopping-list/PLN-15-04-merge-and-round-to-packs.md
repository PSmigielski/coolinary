# PLN-15-04 Merge lines and round to packs

| | |
|---|---|
| Service | planner |
| Epic | [PLN-15 Shopping list](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | PLN-15-03, KIT-02 (GetIngredients) |
| Unblocks | PLN-15-05, PLN-15-10, PLN-16-03, PLN-19-02 |

## Goal

Step 3: the same ingredient from several dishes becomes one line, rounded up to whole packs where the catalogue
knows the pack size (features.md "Shopping list" step 3).

## Business rules

1. Remaining needs are grouped by ingredient id; amounts in the same base unit are summed. Needs left in an
   unconverted unit form a separate line flagged `unit_unconverted`.
2. Each merged line keeps every source with its own contribution (PLN-15-01 rule 4).
3. Pack sizes come from `GetIngredients` (list of pack quantities, for the household's country when kitchen has
   per-country sizes).
4. One pack size `p`: `pack_count = ceil(need / p)`.
5. Several pack sizes: choose the combination of at most 2 distinct sizes and at most 20 packs that covers the need
   with the least excess; ties → fewer packs → larger pack. The line shows the main pack and count ("2 × 500 g");
   a two-size result is returned as two pack entries.
6. No pack size known: round the amount up — g and ml to the next 10 when ≤ 100, to the next 50 when > 100;
   pieces to the next whole piece.
7. A need below 5% of the smallest pack (e.g. 10 g of a 1 kg pack) is still rounded to one pack; the line gets the
   flag `small_need` so the app can say "you need 10 g".
8. Optional-only lines (every source optional) keep the flag `optional`.
9. Free-text lines are never merged with ingredient lines; free text merges only with identical normalised text
   (lowercase, trimmed, single spaces).

## Contract

Internal (`list.MergeAndRound(needs, ingredients) []ProposedLine`). `ListLine.pack_size`/`pack_count` (PLN-15-01);
a second pack entry is added as `repeated PackPart packs = 17` if two sizes are used.

## Data

`planner/migrations/000NN_list_extra_packs.sql`: a second pack size, when used, is stored as
`ALTER TABLE list_lines ADD COLUMN extra_packs jsonb;` (`[{"value": 700, "unit": 1, "count": 1}]`).

## Errors

None of its own; internal logic called by PLN-15-05, PLN-15-10, PLN-16-03 and PLN-19-02. An ingredient without a pack size keeps its exact amount (no rounding), not an error.

## Acceptance criteria

- [ ] Onions 300 g from curry and 150 g from salad → one line 450 g with two sources (Verification 12: "each once,
      with amounts merged across dishes").
- [ ] Rice need 600 g, packs 1 kg only → 1 × 1 kg.
- [ ] Passata need 900 g, packs 500 g and 700 g → 2 × 500 g (excess 100) beats 1 × 700 + 1 × 500 (excess 300).
- [ ] Unknown pack, need 137 g → 150 g; need 63 g → 70 g; 2.3 pieces → 3.

## Notes

- Decided (Q128): at most 2 pack sizes and 20 packs per line (rule 5); without a pack size, round up to 10 g
  (≤ 100 g) or 50 g (rule 6).
