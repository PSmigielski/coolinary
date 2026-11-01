# KIT-04 Barcode lookup chain

| | |
|---|---|
| Service | kitchen |
| Release | R1 |
| Unblocks | GW-06 (scan routes), KIT-05 (unknown products), KIT-07 (unpack mode), Verification 14 |

## Goal

A scanned code becomes a product, an ingredient, a storage place and a use-by date, trying sources in the order of
features.md §4 "How one scan is resolved": this household's own products, Coolinary's shared catalogue, Open Food
Facts, GS1, then "no hit". It also reads GS1 2D codes (use-by date in one scan), shop-internal codes for weighed
goods (prefixes 20–29) and multipacks. Target: ≥ 80% of 50 packaged products per country recognised on first scan
(Verification 14).

## Flow

```
raw scan ──parse (GTIN / 2D / Digital Link / internal code)──┐
                                                             ├─ 20–29 → household internal codes
                                                             └─ GTIN → household products → shared catalogue
                                                                       → Open Food Facts (GATE) → GS1 (GATE) → no hit
hit ──MatchNames──> ingredient ──> storage + use-by suggestion ──> app confirms (KIT-07 / KIT-05)
```

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [KIT-04-01](KIT-04-01-gtin-parsing.md) | GTIN parsing, check digits and prefix classes | R1 | M | — |
| [KIT-04-02](KIT-04-02-gs1-2d-codes.md) | GS1 2D codes and Digital Link | R1 | M | KIT-04-01 |
| [KIT-04-03](KIT-04-03-internal-codes.md) | Shop-internal codes (prefixes 20–29) | R1 | M | KIT-04-01, KIT-10 |
| [KIT-04-04](KIT-04-04-household-products.md) | Household products (level 1) | R1 | M | KIT-04-01 |
| [KIT-04-05](KIT-04-05-shared-catalogue.md) | Shared catalogue (level 2) | R1 | M | KIT-04-01 |
| [KIT-04-06](KIT-04-06-open-food-facts-client.md) | Open Food Facts client (level 3) | R1 | M | KIT-04-01 |
| [KIT-04-07](KIT-04-07-open-food-facts-mapping.md) | Open Food Facts mapping and storage | R1 | M | KIT-04-06, KIT-03-01 |
| [KIT-04-08](KIT-04-08-gs1-lookup.md) | GS1 register lookup (level 4) | R1 | M | KIT-04-01 |
| [KIT-04-09](KIT-04-09-resolve-chain.md) | ResolveBarcode: the lookup chain | R1 | M | KIT-04-02 … KIT-04-07 |
| [KIT-04-10](KIT-04-10-ingredient-and-suggestions.md) | Product → ingredient, storage and use-by suggestion | R1 | M | KIT-04-09, KIT-03-04, KIT-02-03 |
| [KIT-04-11](KIT-04-11-multipacks.md) | Multipacks | R1 | S | KIT-04-01, KIT-04-05, KIT-03-01 |

## Gates

- KIT-04-06/07: Open Food Facts goes live only after the licence check (features.md §4 open points: how the shared
  catalogue may be combined with ODbL data).
- KIT-04-08: GS1 goes live only once access terms and price are agreed for Poland and Germany.
- KIT-04-03: value patterns per chain are filled in only where the chain's pattern is known.

## Decisions (taken 1 November 2026)

- Q067 (KIT-04-01): Yes: 14-digit GTIN as the key; display 8, 12 or 13 digits.
- Q068 (KIT-04-03): The first 7 digits (prefix plus 5 item digits), the common European layout.
- Q069 (KIT-04-03): Per household, keyed by country, item digits and an optional chain.
- Q070 (KIT-04-04): Up to 5,000; a nightly job drops the least recently scanned.
- Q071 (KIT-04-07): Yes: a separate `off_products` table, never mixed with our catalogue.
- Q072 (KIT-04-11): One line with the total (3,000 ml, or 6 pieces).
