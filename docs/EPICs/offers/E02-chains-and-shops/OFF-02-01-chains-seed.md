# OFF-02-01 Chains table and seed

| | |
|---|---|
| Service | offers |
| Epic | [OFF-02 Chains and shops](README.md) |
| Release | R2 |
| Size | S |
| Status | todo |
| Depends on | — |
| Unblocks | OFF-02-02, OFF-02-03, OFF-03-04, OFF-05-01, OFF-11-02 |

## Goal

One table says which grocery chains Coolinary knows, in which country, how to recognise their shops in OSM and on
receipts, how their flyers are split into editions, and how far their coverage has got.

## Business rules

1. A chain is one company in one country: "Lidl PL" and "Lidl DE" are two chains; "Aldi Nord" and "Aldi Süd" are two
   chains (features.md §3: Aldi is two companies split by region). Aldi in Poland is run by Aldi Nord and is its own
   chain "Aldi PL".
2. Seeded chains (id, country, currency, time zone):

   | id | Name | Country | Edition model | Coverage at seed |
   |---|---|---|---|---|
   | `biedronka-pl` | Biedronka | PL | national | `user_flyers_only` |
   | `lidl-pl` | Lidl | PL | national | `user_flyers_only` |
   | `kaufland-pl` | Kaufland | PL | per_store | `user_flyers_only` |
   | `lidl-de` | Lidl | DE | national | `user_flyers_only` |
   | `aldi-sued-de` | Aldi Süd | DE | national | `user_flyers_only` |
   | `aldi-nord-de` | Aldi Nord | DE | regional | `user_flyers_only` |
   | `kaufland-de` | Kaufland | DE | per_store | `user_flyers_only` |
   | `dino-pl` | Dino | PL | national | `user_flyers_only` |
   | `netto-pl` | Netto | PL | national | `user_flyers_only` |
   | `carrefour-pl` | Carrefour | PL | format | `user_flyers_only` |
   | `auchan-pl` | Auchan | PL | national | `user_flyers_only` |
   | `aldi-pl` | Aldi | PL | national | `user_flyers_only` |
   | `penny-de` | Penny | DE | national | `user_flyers_only` |
   | `rewe-de` | Rewe | DE | per_store | `user_flyers_only` |
   | `edeka-de` | Edeka | DE | per_store | `user_flyers_only` |
   | `netto-de` | Netto Marken-Discount | DE | per_store | `user_flyers_only` |

   PL chains use `PLN` and `Europe/Warsaw`; DE chains use `EUR` and `Europe/Berlin`.
3. Coverage states: `user_flyers_only` (shops known, offers only from flyers users scan), `syncing` (adapter runs, not
   yet shown to users), `live` (offers shown; set only by the go-live gate OFF-08-02), `feed` (offers from the
   licensed feed, OFF-13), `removed` (taken down, OFF-08-05). Only OFF-08 tasks and OFF-13 change it.
4. Each chain has OSM match rules: a list of `brand:wikidata` ids and, as fallback, `brand` names (case-insensitive,
   exact after trimming), plus the OSM `shop` values accepted (`supermarket`, `convenience`, `discount`... per chain).
   Wikidata ids are written in the seed and checked by hand at implementation; this task does not invent them.
5. Each chain has receipt aliases: the chain names as printed on receipt headers (e.g. "JERONIMO MARTINS" for
   Biedronka), used by OFF-11-02.
6. The seed is a migration with `INSERT … ON CONFLICT (id) DO NOTHING`; later changes go through admin RPCs.
7. The existing `shops.chain` free-text column is replaced by `chain_id` in OFF-02-02; the example chain
   "Example Market" becomes chain `example-pl` with coverage `live` in dev seed data only.

## Contract

`contracts/coolinary/offers/v1/shops.proto`

```proto
message Chain {
  string id = 1;              // "lidl-pl"
  string name = 2;            // "Lidl"
  string country_code = 3;    // ISO 3166-1 alpha-2
  string currency = 4;        // ISO 4217
  ChainCoverage coverage = 5;
}

enum ChainCoverage {
  CHAIN_COVERAGE_UNSPECIFIED = 0;
  CHAIN_COVERAGE_USER_FLYERS_ONLY = 1;
  CHAIN_COVERAGE_SYNCING = 2;
  CHAIN_COVERAGE_LIVE = 3;
  CHAIN_COVERAGE_FEED = 4;
  CHAIN_COVERAGE_REMOVED = 5;
}
```

Admin (`admin.proto`): `ListChains`, `UpdateChain` (name, OSM match rules, receipt aliases; not coverage).

## Data

`offers/migrations/000NN_chains.sql`

```sql
CREATE TABLE chains (
    id               text PRIMARY KEY,
    name             text NOT NULL,
    country_code     text NOT NULL CHECK (country_code ~ '^[A-Z]{2}$'),
    currency         text NOT NULL CHECK (currency ~ '^[A-Z]{3}$'),
    time_zone        text NOT NULL,
    edition_model    text NOT NULL CHECK (edition_model IN ('national', 'regional', 'per_store', 'format')),
    coverage         text NOT NULL DEFAULT 'user_flyers_only'
                     CHECK (coverage IN ('user_flyers_only', 'syncing', 'live', 'feed', 'removed')),
    osm_wikidata     text[] NOT NULL DEFAULT '{}',
    osm_brand_names  text[] NOT NULL DEFAULT '{}',
    osm_shop_values  text[] NOT NULL DEFAULT '{supermarket}',
    receipt_aliases  text[] NOT NULL DEFAULT '{}',
    created_at       timestamptz NOT NULL DEFAULT now()
);
```

## Errors

| Case | gRPC status |
|---|---|
| `UpdateChain` tries to change `coverage` | `INVALID_ARGUMENT` "coverage changes through the go-live gate" |
| Unknown chain id | `NOT_FOUND` |

## Acceptance criteria

- [ ] After migration, 16 chains exist with the countries, currencies and edition models above, all
      `user_flyers_only`.
- [ ] `UpdateChain` adds a receipt alias and writes an `admin_changes` row.
- [ ] Running the migration twice adds nothing.
