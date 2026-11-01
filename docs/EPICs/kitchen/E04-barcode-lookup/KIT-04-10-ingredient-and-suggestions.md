# KIT-04-10 Product → ingredient, storage and use-by suggestion

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-04 Barcode lookup chain](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | KIT-04-09, KIT-03-04, KIT-02-03, KIT-10 (household timezone) |
| Unblocks | KIT-07 (stock confirm), PLN-06 (scanned items count for recipes) |

## Goal

After a hit, the product is linked to a catalogue ingredient so recipes can use it, and the app proposes where to
store it and until when (features.md §4 "After a hit, the product is matched to a catalogue ingredient, and the
storage place and use-by date are suggested from the product type").

## Business rules

1. If the product row (any source) has no `ingredient_id`, call `MatchNames` (`INTERACTIVE`, context `product`) with
   the product name, brand, and the language of the name.
2. A `matched` answer without `needs_review` is saved on the source row (shared, OFF or GS1 table), so every
   household gets it. A `matched` answer with `needs_review` is returned and saved, and updated later if the review
   changes it. `none`, `queued` or `provisional` → the product has no ingredient for now; the app lets the user pick
   one with autocomplete, which is saved only on the household's row.
3. Storage suggestion: OFF category `en:frozen-foods` or similar frozen tags → freezer; else the household's own
   storage for this product (level 1); else the ingredient's default storage; else pantry.
4. Use-by suggestion: a date from the 2D code (KIT-04-02) wins; else today in the household's timezone + the
   ingredient's shelf life for the suggested place (KIT-02-03); else none.
5. The suggestion carries `storage_unusual` from KIT-02-03 and `date_source` (`code`, `shelf_life`, `none`).
6. Product allergens and traces (from OFF, KIT-04-07) are returned with the product; planner's rule checker uses the
   ingredient's allergens plus these for stock items.

## Contract

`contracts/coolinary/kitchen/v1/barcode.proto`

```proto
message Suggestion { StoragePlace storage = 1; string use_by = 2; string date_source = 3; bool storage_unusual = 4; }
```

## Data

No schema change (writes `ingredient_id` on existing tables).

## Errors

Matching problems never fail the scan; the product is returned without an ingredient.

## Acceptance criteria

- [ ] A scanned "Chicken thighs 1 kg" product is linked to the chicken-thigh ingredient and suggested for the fridge
      with a 2-day date.
- [ ] A frozen-food OFF product is suggested for the freezer.
- [ ] A product with a 2D expiry date uses that date (`date_source = code`).
- [ ] A low-confidence match is saved and flagged for review (ADM-03 shows it).
