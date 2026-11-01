# KIT-06-03 Import Bundeslebensmittelschlüssel 4.0

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-06 Nutrition](README.md) |
| Release | R1 (GATE: terms for a commercial app) |
| Size | M |
| Status | todo |
| Depends on | KIT-06-01 |
| Unblocks | KIT-06-05 |

## Goal

Load the German BLS 4.0 (Max Rubner-Institut; about 7,140 foods and dishes, up to 138 nutrients, free since December
2025) for German foods and as the second source where Ciqual has no entry (features.md §5 table).

## Business rules

1. Command `server import-nutrition --source bls --file <path> --version <label> [--dry-run]`, same behaviour as
   KIT-06-02 rules 1, 5, 6, 7, 8.
2. Mapping: the BLS 4.0 component codes for energy (kcal and kJ), protein, fat, carbohydrate, total sugars,
   saturated fatty acids, dietary fibre and sodium (salt = sodium × 2.5 if salt isn't given). The exact column codes
   are taken from the BLS 4.0 documentation and kept in one mapping table in the importer.
3. Name: the German food name, plus the English name where the release has one.
4. `state` from the German name: "roh" → `raw`; "getrocknet" → `dry`; "gekocht", "gegart", "gebraten", "gebacken",
   "frittiert" → `cooked`; "zubereitet" → `prepared`; else `other`.
5. BLS units: values given per 100 g edible portion are used as is; mg and µg values are converted to g.

## Contract

No RPC (CLI).

## Data

No schema change.

## Errors

As KIT-06-02.

## Acceptance criteria

- [ ] Import creates about 7,140 entries; a second run changes nothing.
- [ ] 10 sample foods match the published values.
- [ ] Sodium-only entries get salt computed.
