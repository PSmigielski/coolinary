# PLN-12-03 Portion and pot calculator

| | |
|---|---|
| Service | planner |
| Epic | [PLN-12 Per-person portions](README.md) |
| Release | R2 |
| Size | M |
| Depends on | PLN-12-02, PLN-04 |
| Unblocks | PLN-12-04, PLN-12-05, PLN-12-06, PLN-13-06 |

## Goal

For one dish and one meal, compute each eater's portion, the pot, who needs an extra and whether the dish is
unsuitable (features.md "How one meal is split", steps 3–5).

## Business rules

1. `K` = the dish's kcal per standard portion from PLN-04, unrounded. Display rounding (nearest 10 kcal) is never
   used in the arithmetic.
2. For an eater with a share `S` (meal target plus carry-in, PLN-12-05): `raw = S / K`; portion `p` = `raw` rounded
   to the nearest quarter, halves rounded up (`floor(raw × 4 + 0.5) / 4`), then clamped to `[0.5, 2.0]`.
3. An eater without a target gets `p = 1.0` and is never short or over.
4. **Short:** if `p = 2.0` and `S − 2K > 0.15 × S`, the eater needs an extra worth `gap = S − 2K` (PLN-12-04).
5. **Over:** if `0.5 × K > 1.15 × S`, the dish is unsuitable for this eater at this meal. The calculator returns
   `UNSUITABLE` with the member ids; the composer picks a different dish (PLN-13-06, PLN-14-06). Calories never
   cause a second dish.
6. Gaps within these limits are not fixed here; they carry over (PLN-12-05).
7. **Pot** = sum of all eaters' portions rounded up to a whole number (`ceil`). The difference is reported as
   `spare_portions` (e.g. 0.25).
8. No eaters → pot 0, no ingredients.
9. Dish with calories incomplete (PLN-04): every eater gets 1.0, the result is flagged `calories_incomplete`, nobody
   is short or over.
10. Output per eater: member id (or unnamed index), share, portion, portion kcal (`p × K`), gap, status
    (`ok` / `short` / `over_minimum`).

## Contract

Internal package `planner/internal/portions` (`Allocate(dishKcal, []EaterShare) Allocation`). No RPC.

## Data

No schema change.

## Errors

| Case | Behaviour |
|---|---|
| `K ≤ 0` and not flagged incomplete | error `dish_kcal_invalid`, caller returns `INTERNAL` |
| Share negative after carry-over | treated as 0, so the eater gets 0.5 and is checked for over |

## Acceptance criteria

- [ ] Worked example (chicken curry, K = 520): Adam 800 → 1½ (780 kcal), Ewa 720 → 1½ (780), Ola 560 → 1 (520),
      Kuba 1,280 → 2 (1,040, short 240 = 18.75%, needs an extra). Pot = 6. Only Kuba is short (Verification 17).
- [ ] Ola away: pot = 5.
- [ ] Rounding: raw 1.375 → 1.5; raw 1.374 → 1.25; raw 0.3 → 0.5; raw 2.6 → 2.0.
- [ ] A 1,300 kcal dish for a child with a 400 share: ½ portion = 650 > 460 → `UNSUITABLE` with the child's id.
- [ ] Incomplete calories: 4 eaters → pot 4, flagged.
- [ ] Pure function; table-driven tests cover every rule.

## Notes

- Decided (Q198): rule 9 (incomplete calories → everyone gets 1 standard portion, flagged "calories incomplete").
