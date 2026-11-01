# KIT-05-04 Send confirmed products to Open Food Facts

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-05 Unknown products](README.md) |
| Release | R1 (GATE: licence decision and an Open Food Facts app account) |
| Size | M |
| Status | todo |
| Depends on | KIT-05-02, KIT-04-07, KIT-10-01 (`contribute_to_off` setting) |
| Unblocks | — |

## Goal

Products Coolinary users confirmed that Open Food Facts doesn't have are sent back to it, as its terms expect
(features.md §4 table "Product additions should be sent back").

## Business rules

1. A daily job at 04:00 UTC selects shared products with `status = 'confirmed'`, `source = 'user'`, no row in
   `off_products`, and not sent before; at most 200 per day.
2. Decided (Q076): each household can opt out. The kitchen profile setting `contribute_to_off` (KIT-10-01) is
   `true` by default. A product is sent only when at least two households whose setting is `true` (read at send
   time) submitted the confirmed name and brand. Submissions of households that opted out still count for
   confirmation inside Coolinary but never for sending. Decided (Q231): opting out stops future sends; products already sent are
   not withdrawn.
3. Sent fields: code, product name (with its language), brand, quantity text, and the countries of the contributing
   households that confirmed it (country codes only). Never household ids, photos, or anything else.
4. Sent through OFF's write API with Coolinary's app account (credentials as Docker secrets) and the app
   User-Agent (KIT-04-06 rule 1), tagged as coming from Coolinary.
5. Each attempt is recorded; failures are retried on the next 3 runs, then left with `last_error`.
6. Config `OFF_CONTRIBUTE_ENABLED` (default `false`).
7. GATE: enabled only after the licence decision in features.md §4 open points and after the OFF account is set up.

## Contract

No RPC.

## Data

`kitchen/migrations/000NN_off_contributions.sql`

```sql
CREATE TABLE off_contributions (
    gtin         char(14) PRIMARY KEY,
    attempts     integer NOT NULL DEFAULT 0,
    sent_at      timestamptz,
    last_error   text,
    updated_at   timestamptz NOT NULL DEFAULT now()
);
```

## Errors

None returned to users.

## Acceptance criteria

- [ ] With a stub OFF write endpoint, a newly confirmed product is sent once with exactly the listed fields.
- [ ] A product already in `off_products` is never sent.
- [ ] The daily cap of 200 holds.
- [ ] A product confirmed by two households, one of which set `contribute_to_off = false`, is not sent; after a
      third contributing household confirms it, it is sent with only the contributing households' countries.
- [ ] Turning the setting off after a product was sent changes nothing for that product and stops later sends.
