# OFF-03-04 Flyer editions and shop-to-edition rules

| | |
|---|---|
| Service | offers |
| Epic | [OFF-03 Postcode areas and flyer editions](README.md) |
| Release | R2 |
| Size | M |
| Depends on | OFF-02-01, OFF-02-02 |
| Unblocks | OFF-03-03, OFF-03-05, OFF-04-03, OFF-05-01, OFF-07-01, OFF-09-01 |

## Goal

Every offer belongs to a flyer edition, and every shop knows which editions apply to it, so location picks the right
flyer where a chain's flyer differs by store or region (features.md §2 "Flyers and offers", §3 "Kaufland's flyer
depends on the store").

## Business rules

1. An edition is one stream of flyers of one chain: `national` (one per chain), `region` (one per chain region),
   `store` (one per store, keyed by the chain's store number), `format` (one per store format, e.g. Carrefour
   Hipermarket, Market, Express).
2. Shop → editions, by the chain's `edition_model`:
   - `national`: the chain's national edition.
   - `regional` (Aldi Nord): the national edition. Decided (Q148): Aldi Nord's regional extras are not read in R2;
     only the national flyer is (the mapping of shops to Aldi Nord's regions is not known).
   - `per_store` (Kaufland, later Rewe, Edeka, Netto DE): the store edition of the shop's linked store number
     (OFF-03-05). A shop without a link has no edition; its flyer status is "flyer for this store not available yet".
   - `format` (Carrefour PL, R3): the edition of the shop's format, derived from the OSM `shop` value and name
     (`Carrefour Express` → `express`, `shop=supermarket` → `market`, `shop=hypermarket` or name containing
     "Hipermarket" → `hipermarket`).
3. Aldi Nord and Aldi Süd are separate chains with their own national editions; which one a user sees follows only
   from which shops are near them.
4. Editions are created lazily: national editions with the chain seed; store editions when a store number is linked;
   format editions on first use.
5. An edition has a `demand` count (OFF-04-03) used to decide whether to read it (OFF-05-07); national editions of
   live chains are always read.
6. Function `EditionsForShop(shopID)` is the only place these rules live; browse, prices and status all call it.

## Contract

```proto
message Edition {
  string id = 1;                // "lidl-pl:national", "kaufland-de:store:1234", "carrefour-pl:format:express"
  string chain_id = 2;
  EditionKind kind = 3;
  string key = 4;               // "" / region code / store number / format
  string name = 5;
}
enum EditionKind {
  EDITION_KIND_UNSPECIFIED = 0;
  EDITION_KIND_NATIONAL = 1;
  EDITION_KIND_REGION = 2;
  EDITION_KIND_STORE = 3;
  EDITION_KIND_FORMAT = 4;
}
```

## Data

`offers/migrations/000NN_editions.sql`

```sql
CREATE TABLE editions (
    id         text PRIMARY KEY,
    chain_id   text NOT NULL REFERENCES chains (id),
    kind       text NOT NULL CHECK (kind IN ('national', 'region', 'store', 'format')),
    key        text NOT NULL DEFAULT '',
    name       text NOT NULL,
    demand     int NOT NULL DEFAULT 0,
    created_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (chain_id, kind, key)
);
INSERT INTO editions (id, chain_id, kind, name)
SELECT id || ':national', id, 'national', name FROM chains WHERE edition_model IN ('national', 'regional');
```

## Errors

| Case | Behaviour |
|---|---|
| Per-store shop without a store link | No edition; not an error |
| Unknown `edition_model` | `INTERNAL` (bug) |

## Acceptance criteria

- [ ] Every Lidl PL shop maps to `lidl-pl:national`.
- [ ] Two Kaufland DE shops linked to stores 1234 and 5678 map to two different editions (Verification 8).
- [ ] An unlinked Kaufland shop maps to no edition.
