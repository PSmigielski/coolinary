# OFF-12-03 Carrefour PL adapter (per store format)

| | |
|---|---|
| Service | offers |
| Epic | [OFF-12 Next chain adapters](README.md) |
| Release | R3 |
| Size | L |
| Depends on | OFF-05-02, OFF-05-03, OFF-03-04 |
| Unblocks | OFF-08-02 for `carrefour-pl` |

## Goal

Carrefour Poland's flyers, which differ by store format, reach the right shops (features.md §3: "PDF offered; differs
by store format").

## Business rules

1. Source `carrefour-pl`: adapter `carrefour_pl`, list URL `https://www.carrefour.pl/gazetka-handlowa`, allowed hosts
   `www.carrefour.pl` plus the PDF host from the fixtures.
2. Editions by format (OFF-03-04 rule 2, `format`): `hipermarket`, `market`, `express`. The adapter maps each flyer on
   the list page to one or more formats from its title or section (e.g. "Gazetka Carrefour Market"). A flyer whose
   format cannot be read is assigned to all three and flagged in the check record.
3. PDFs are preferred over page images (cheaper to read).
4. Format of each shop: from OSM (`shop=hypermarket` or "Hipermarket" in the name → hipermarket; "Express" in the name
   → express; otherwise market). Admin can override per shop (`shop_overrides`).
5. Daily check 06:00 `Europe/Warsaw`; extra checks set after the first shadow week.
6. `GATE: retailer terms checked` for carrefour.pl.

## Contract

No gRPC change.

## Data

Migration inserts the source, disabled, terms not approved.

## Errors

| Case | Behaviour |
|---|---|
| Expected list markers missing | `ErrLayoutChanged` (OFF-05-02) |

## Acceptance criteria

- [ ] Fixtures with a Market and an Express flyer map them to the two format editions; a Carrefour Express shop gets
      only the Express offers.
- [ ] Two-week shadow run detects every flyer before its start date.
