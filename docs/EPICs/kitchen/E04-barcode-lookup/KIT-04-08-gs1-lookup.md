# KIT-04-08 GS1 register lookup (level 4)

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-04 Barcode lookup chain](README.md) |
| Release | R1 (GATE: GS1 access terms and price for Poland and Germany) |
| Size | M |
| Status | todo |
| Depends on | KIT-04-01 |
| Unblocks | KIT-04-09 |

## Goal

Codes neither Coolinary nor Open Food Facts know get their official name, brand and brand owner from the GS1
register, and a code can be confirmed as genuinely registered (features.md §4 step 4).

## Business rules

1. Adapter interface `gs1.Client.Lookup(ctx, gtin) (Record, error)` so the access route agreed with GS1 (GS1 Polska
   eProdukty, or Verified by GS1 through a member organisation) can be swapped without touching the chain.
2. Config: `GS1_ENABLED` (default `false`), `GS1_BASE_URL`, `GS1_API_KEY` (secret), `GS1_DAILY_LIMIT`.
3. Called only on a real scan, after Open Food Facts missed. Timeout 2 s; failures skip the level (`degraded`).
4. A daily counter in redis-state (`kitchen:gs1:<UTC date>`) stops calls at `GS1_DAILY_LIMIT`. The public Verified by
   GS1 web page (30 lookups a day) is never scraped.
5. Stored in `gs1_products` (separate table, licence terms): product description, brand, brand owner company name,
   GPC category code, net content, `registered = true`. Refreshed after 180 days on the next scan.
6. A GS1 record with no usable description still sets `registered = true`; the chain then asks the user to name the
   product (KIT-05), and the naming flow can say "this is a registered product".
7. GATE: stays disabled in production until access terms and price are agreed (features.md §4 open points).

## Contract

`ResolveSource` already has `RESOLVE_SOURCE_GS1 = 4`. `ResolveBarcodeResponse` gains `gtin_registered`.

## Data

`kitchen/migrations/000NN_gs1_products.sql`

```sql
CREATE TABLE gs1_products (
    gtin          char(14) PRIMARY KEY,
    description   text NOT NULL DEFAULT '',
    brand         text NOT NULL DEFAULT '',
    brand_owner   text NOT NULL DEFAULT '',
    gpc_code      text NOT NULL DEFAULT '',
    net_value     double precision,
    net_unit      smallint NOT NULL DEFAULT 0,
    registered    boolean NOT NULL,
    fetched_at    timestamptz NOT NULL DEFAULT now()
);
```

## Errors

GS1 problems never fail `ResolveBarcode`.

## Acceptance criteria

- [ ] With a stub adapter, a code missed by OFF resolves from GS1 and is stored.
- [ ] Calls stop at the daily limit and the response is `degraded`.
- [ ] `GS1_ENABLED=false` makes no call.
