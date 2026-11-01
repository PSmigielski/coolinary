# OFF-07-03 Normalise prices, packs and conditions

| | |
|---|---|
| Service | offers |
| Epic | [OFF-07 Offer extraction, validation and publishing](README.md) |
| Release | R2 |
| Size | L |
| Status | todo |
| Depends on | OFF-07-02 |
| Unblocks | OFF-07-04, OFF-07-08, OFF-10-03 |

## Goal

Code, not AI, turns what was printed into numbers the app and the trip plan can compare: price in minor units, pack
size, price per kg, litre or piece, and the conditions ("2+1", "with app or card") as structured fields.

## Business rules

1. Price: parsed from the printed price text with the flyer language's format (`4,99`, `4,99 zł`, `4.99 €`,
   `-,99`, `1.-`). If the parsed value and the AI's value differ, the parsed one is kept and the offer gets the held
   reason `price_mismatch`. If no parse is possible, the AI's value is used and the offer gets `price_unparsed`.
2. Regular price ("zamiast 6,99", "statt 2,49", struck-through price): parsed the same way; kept only if greater than
   the price; `discount_percent` = round((regular − price) / regular × 100).
3. Pack, parsed from the pack text: `500 g`, `0,5 kg`, `1 l`, `330 ml`, `6 x 0,5 l` (count 6, 500 ml each), `10 szt.`,
   `10 Stück`, `ca. 1 kg`/`ok. 1 kg` (approximate, kept as 1 kg), `1 kg` with "cena za kg" / "je kg" (loose, per kg).
   Units are converted to `g`, `ml` or `piece`.
4. Unit price: total quantity = count × per-unit value; price per kg (grams), per litre (ml) or per piece. Rounded to
   the nearest minor unit. Loose offers: the price is already per kg or per piece.
5. Multi-buy, from the conditions text: "2+1" / "1+1 gratis" → buy 3 pay 2 / buy 2 pay 1; "3 za 2" / "3 für 2" →
   buy 3 pay 2; "2 za 5 zł" / "2 für 5 €" → buy 2, total 5.00; "druga sztuka -50%" / "2. Stück -50%" → buy 2 pay 1.5
   (stored as buy 4 pay 3). Unknown patterns are kept as text only.
6. Loyalty, from the conditions text and per-chain keyword lists: "z aplikacją Lidl Plus", "z Kartą Moja Biedronka",
   "mit Lidl Plus", "mit Kaufland Card", "with app or card" → `loyalty_required = true` and the program name.
   features.md Risks: "Prices that need a loyalty app or card: Shown as a condition on the offer".
7. "max. 2 szt. na klienta" / "max. 2 pro Kunde" → `max_per_customer = 2`; "przy zakupie 2 szt." → `min_quantity = 2`.
8. Non-food (Decided, Q155): offers the AI marks `non_food` (electronics, clothes, garden) are dropped here and
   never stored.
9. Text fields are trimmed and whitespace collapsed; product names keep their original language and spelling.
10. Normalisation is deterministic and has a table-driven test with at least 60 real examples per language
    (written from the fixtures).

## Contract

No gRPC change.

## Data

No schema change (fills OFF-07-01 columns). Raw rows stay in `extracted_offers` for comparison.

## Errors

| Case | Behaviour |
|---|---|
| Pack text not parseable | No pack; validation holds the offer with `no_unit` |
| Unknown condition pattern | Text kept, no structured terms; not an error |

## Acceptance criteria

- [ ] "6 x 0,5 l" at 9,99 zł → count 6, 500 ml each, 3,33 zł per litre.
- [ ] "2+1" → buy 3 pay 2; "2 für 5 €" → buy 2 total 5.00 €.
- [ ] "4,99 zł z Kartą Moja Biedronka" → price 499 PLN minor, loyalty program "Moja Biedronka".
- [ ] AI price 4.99 with printed "5,99" gives price 599 and held reason `price_mismatch`.
