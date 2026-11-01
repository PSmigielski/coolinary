# OFF-02-04 Merge duplicate shops

| | |
|---|---|
| Service | offers |
| Epic | [OFF-02 Chains and shops](README.md) |
| Release | R2 |
| Size | S |
| Depends on | OFF-02-03 |
| Unblocks | OFF-03-03 |

## Goal

One real shop appears once in the app, even when OSM maps it both as a point and as a building outline, or when a
user adds a shop that OSM later gets.

## Business rules

1. Two active shops are duplicates when they have the same `chain_id` and their positions are at most 50 m apart.
   (50 m: two branches of one chain are practically never that close; a point inside its own building usually is.)
2. Which one is kept, in this order: the one referenced by more households in my shops; then the OSM way over the
   node; then the OSM shop over the user-added one; then the older id.
3. The other becomes `merged` with `merged_into` set. Every reference moves to the kept shop in the same
   transaction: my shops (OFF-04-01, duplicates collapse to one entry keeping the better order position), edition
   demand (OFF-04-03), store-number links (OFF-03-05), price observations with a shop (OFF-11-01).
4. Runs right after every OSM import and after an admin approves a user-added shop.
5. A pair the admin marked `not_duplicate` (two shops in one mall) is never merged again.
6. Merges are logged in `shop_merges`.

## Contract

Admin (`admin.proto`): `MarkNotDuplicate(shop_id_a, shop_id_b)`; `ListShopMerges(page)`.

## Data

`offers/migrations/000NN_shop_merges.sql`

```sql
CREATE TABLE shop_merges (
    merged_id  text NOT NULL REFERENCES shops (id),
    kept_id    text NOT NULL REFERENCES shops (id),
    distance_m numeric(6,1) NOT NULL,
    merged_at  timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (merged_id)
);
CREATE TABLE shop_not_duplicates (
    shop_a text NOT NULL REFERENCES shops (id),
    shop_b text NOT NULL REFERENCES shops (id),
    PRIMARY KEY (shop_a, shop_b),
    CHECK (shop_a < shop_b)
);
```

## Errors

| Case | Behaviour |
|---|---|
| Pair marked not-duplicate | Skipped |
| Moving a reference fails | Whole merge rolled back; retried after the next import |

## Acceptance criteria

- [ ] A node and a way of the same Lidl 20 m apart merge into the way; a household that had the node in my shops now
      has the way.
- [ ] Two Lidl shops 300 m apart are not merged.
- [ ] A pair marked not-duplicate stays separate after the next import.
