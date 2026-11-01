# KIT-04-09 ResolveBarcode: the lookup chain

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-04 Barcode lookup chain](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | KIT-04-02, KIT-04-03, KIT-04-04, KIT-04-05, KIT-04-06, KIT-04-07 (KIT-04-08 optional) |
| Unblocks | KIT-04-10, KIT-05-01, GW-06 (`GET /v1/products/{code}`), KIT-07 (unpack mode) |

## Goal

One RPC answers any scan by walking the sources in order and stopping at the first hit, never failing because an
outside source is slow or down.

## Business rules

1. Input: `code` (raw scan: digits, element string or URI), `symbology`, `household_id`, `language`, `country`,
   `user_scan` (true when a person just scanned it). The old `ean` field still works and means `code` with
   symbology unknown.
2. Parse with KIT-04-01 / KIT-04-02. `not_food` → found false, `code_class = NOT_FOOD`. `demo` outside prod →
   normal lookup; in prod → found false.
3. `restricted_internal` → KIT-04-03 only.
4. Otherwise, in order, stopping at the first hit:
   1. household products (KIT-04-04) → `HOUSEHOLD`;
   2. shared catalogue, `confirmed` only (KIT-04-05) → `SHARED_CATALOGUE`;
   3. if no hit and a `pending` submission exists from another household → found false with a `proposal`
      (KIT-05-02);
   4. Open Food Facts (KIT-04-06/07) → `OPEN_FOOD_FACTS`;
   5. GS1 (KIT-04-08) → `GS1`;
   6. miss → found false, `needs_naming = true`.
5. A hit is enriched by KIT-04-10 (ingredient, storage, use-by) and carries any 2D extras (use-by, batch, weight).
6. Total deadline 3 s. An outside source that fails or is skipped sets `degraded = true`; the chain continues.
7. Results of levels 2–5 for a GTIN are cached together in redis-cache for 24 hours; level 1 is never cached across
   households.
8. Verification 14 metric: `kitchen_barcode_resolve_total{source, first_scan, country}`, where `first_scan` is true
   when the household had no level-1 row for the code before this call. The share of first scans that are hits per
   country is the recognition rate (target ≥ 80%).
9. No personal data in logs: the code may be logged; the household id only hashed (KIT-01-01 rule 7).

## Contract

`contracts/coolinary/kitchen/v1/barcode.proto` (extends the existing messages)

```proto
message ResolveBarcodeRequest { string ean = 1; string household_id = 2; string code = 3; Symbology symbology = 4;
  string language = 5; string country = 6; bool user_scan = 7; }
message ResolveBarcodeResponse { bool found = 1; Product product = 2; ResolveSource source = 3;
  ParsedCode parsed = 4; Suggestion suggestion = 5; bool needs_naming = 6; ProductProposal proposal = 7;
  bool degraded = 8; string attribution = 9; string attribution_url = 10; bool gtin_registered = 11;
  CodeClass code_class = 12; repeated InternalCodeChoice internal_choices = 13; }
```

## Data

No schema change.

## Errors

| Case | Status | Message |
|---|---|---|
| Unparseable code | `INVALID_ARGUMENT` | from KIT-04-01 / 02 |
| Household mismatch | `PERMISSION_DENIED` | `household mismatch` |
| Postgres down | `UNAVAILABLE` | `catalogue unavailable` |

## Acceptance criteria

- [ ] Each level answers when the ones before it miss (stubs per level).
- [ ] OFF timing out returns a GS1 or miss answer within 3 s with `degraded = true`.
- [ ] Verification 14: scanning the 50-product Polish and German test sets yields the recognition rate per country
      from the metric; an unknown product named once is recognised on the next scan from a second account (with
      KIT-05).
- [ ] The existing gateway route `GET /v1/products/{ean}` keeps working unchanged.
