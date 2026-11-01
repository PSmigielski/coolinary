# KIT-02-02 Typical weights, density and unit conversion

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-02 Ingredient catalogue](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | KIT-02-01 |
| Unblocks | KIT-06-08 (nutrition), KIT-07 (stock consumption), PLN-12 / PLN-15 (scaling, list) |

## Goal

Turn any recipe or stock amount ("1 onion", "2 tbsp oil", "200 ml milk") into grams, and back into the stock's
unit, so calories, stock and shopping lists use the same numbers (features.md §5 step 2).

## Business rules

1. Each ingredient may have measures with a weight in grams: `piece`, `piece_small`, `piece_medium`,
   `piece_large`, `slice`, `clove`, `bunch`, `leaf`, `sprig`, `handful`, `can`, `jar`, `tbsp`, `tsp`, `cup`,
   `glass`, `pinch`. `piece` is the default piece; if it is missing, `piece_medium` is used.
2. Each ingredient may have `density_g_per_ml` between 0.2 and 2.5.
3. Volume measures without an explicit gram weight convert through density: `tbsp` = 15 ml, `tsp` = 5 ml,
   `cup` = 250 ml, `glass` = 250 ml (Polish "szklanka").
4. Decided (Q051): `pinch` without its own weight counts 0.5 g; `to taste` counts 0 g and is reported as
   "not counted", never as missing.
5. Decided (Q050): there is no default density. Millilitres of an ingredient without a density can't become
   grams and are reported as `no_weight` (nutrition then says "calories incomplete") until a density is added.
6. `common.v1.Unit` gets the measure values recipes need: `UNIT_TABLESPOON = 4`, `UNIT_TEASPOON = 5`,
   `UNIT_CUP = 6`, `UNIT_PINCH = 7`, `UNIT_TO_TASTE = 8`. Piece sizes and other measures (`clove`, `slice` …) are
   passed as `Quantity.measure` (new string field) with `UNIT_PIECE`.
7. Library `internal/kitchen/units`:
   - `ToGrams(ingredient, Quantity) (grams float64, status)` with status `ok`, `not_counted`, `no_weight`.
   - `FromGrams(ingredient, grams, targetUnit) (value, status)` for stock: grams → ml via density; grams → pieces
     via the default piece weight.
   - `Convert(ingredient, Quantity, targetUnit)` combines both. No rounding inside the library; callers round.
8. Changing a measure or density raises the ingredient's `nutrition_version` (KIT-06-10), because calories of
   every dish using it may change.

## Contract

`contracts/coolinary/common/v1/common.proto` (shared change; the one definition of `Unit` for every service,
README decision 16; planner and ai import it)

```proto
enum Unit { UNIT_UNSPECIFIED = 0; UNIT_GRAM = 1; UNIT_MILLILITRE = 2; UNIT_PIECE = 3;
  UNIT_TABLESPOON = 4; UNIT_TEASPOON = 5; UNIT_CUP = 6; UNIT_PINCH = 7; UNIT_TO_TASTE = 8; }
message Quantity { double value = 1; Unit unit = 2; string measure = 3; } // measure: "clove", "piece_large", …
```

`contracts/coolinary/kitchen/v1/catalogue.proto`

```proto
message Measure { string measure = 1; double grams = 2; }
```

## Data

`kitchen/migrations/000NN_ingredient_measures.sql`

```sql
ALTER TABLE ingredients ADD COLUMN density_g_per_ml numeric(5,3) CHECK (density_g_per_ml BETWEEN 0.2 AND 2.5);
CREATE TABLE ingredient_measures (
    ingredient_id  uuid NOT NULL REFERENCES ingredients (id) ON DELETE CASCADE,
    measure        text NOT NULL CHECK (measure IN ('piece','piece_small','piece_medium','piece_large','slice',
                    'clove','bunch','leaf','sprig','handful','can','jar','tbsp','tsp','cup','glass','pinch')),
    grams          numeric(8,2) NOT NULL CHECK (grams > 0 AND grams <= 5000),
    PRIMARY KEY (ingredient_id, measure)
);
```

## Errors

Library returns statuses, not errors. Bad input in RPCs that use it is `INVALID_ARGUMENT` (negative or zero value,
unknown measure).

## Acceptance criteria

- [ ] Onion with `piece` = 110 g: 1 piece → 110 g; 330 g → 3 pieces.
- [ ] Olive oil with density 0.91 and no `tbsp` weight: 2 tbsp → 27.3 g.
- [ ] Flour with `tbsp` = 10 g: 3 tbsp → 30 g (explicit weight beats density).
- [ ] Milk without density: 200 ml → `no_weight`.
- [ ] Salt "to taste" → 0 g with `not_counted`.
- [ ] Adding a measure raises `nutrition_version` by 1.
