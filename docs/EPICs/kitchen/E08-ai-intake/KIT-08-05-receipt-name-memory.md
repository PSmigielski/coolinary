# KIT-08-05 Receipt-name memory per household and chain

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-08 AI intake](README.md) |
| Release | R3 |
| Size | M |
| Depends on | KIT-08-04, KIT-07-03 (`AddStockItems` with `intake_ref`) |
| Unblocks | KIT-08-04 rule 6.1 |

## Goal

Once a user has confirmed that "MLEKO UHT 3,2% 1L" at Lidl is their usual milk, every later Lidl receipt matches
that line at once, with no AI and no guessing (features.md §4: "lines are matched to products this household
scanned before at that shop").

## Business rules

1. On `AddStockItems` with source `receipt` and an `intake_ref` whose receipt copy is still in `redis-cache`
   (KIT-08-04 rule 10), every confirmed item that came from a printed name updates the memory for
   (household, chain key, normalised printed name):
   1. the product and ingredient as finally confirmed (the user's correction wins over the proposal);
   2. `confirmed_count` + 1, `last_seen_on` = the purchase date.
2. A proposal the user unticked teaches nothing; an existing link is not removed by one untick.
3. A user who changes the ingredient of a remembered line replaces the link (count restarts at 1).
4. Normalising a printed name: lower case, whitespace collapsed, trailing tokens that are only a price, a quantity
   ("2x", "0,312kg") or a tax letter ("A", "B") removed.
5. At most 5,000 links per household; adding the 5,001st deletes the one with the oldest `last_seen_on`.
6. Links unused for 365 days are deleted by KIT-13-08.
7. Links are this household's data: never shared with other households, included in the export (KIT-13-07),
   deleted with the household (KIT-13-04).
8. Ingredient and product merges (KIT-13-09) update links.

## Contract

No new RPC (side effect of `AddStockItems`).

## Data

`kitchen/migrations/000NN_receipt_name_links.sql`

```sql
-- +goose Up
CREATE TABLE receipt_name_links (
    household_id      uuid NOT NULL,
    chain_key         text NOT NULL,
    printed_name_norm text NOT NULL CHECK (char_length(printed_name_norm) <= 120),
    product_ean       text REFERENCES products (ean),
    ingredient_id     uuid REFERENCES ingredients (id),
    confirmed_count   integer NOT NULL DEFAULT 1,
    last_seen_on      date NOT NULL,
    PRIMARY KEY (household_id, chain_key, printed_name_norm),
    CHECK (product_ean IS NOT NULL OR ingredient_id IS NOT NULL)
);
CREATE INDEX receipt_name_links_last_seen ON receipt_name_links (household_id, last_seen_on);

-- +goose Down
DROP TABLE receipt_name_links;
```

## Errors

| Case | Behaviour |
|---|---|
| Receipt copy evicted from Redis | items still added; no memory update; metric `kitchen_receipt_context_missing_total` + 1 |

## Acceptance criteria

- [ ] Confirm "MLEKO UHT 3,2% 1L" → milk product at `lidl`; the next Lidl receipt matches it with confidence 1.0 and
      ticked, without a `MatchNames` call for that line.
- [ ] The same printed name at `biedronka` is not matched by the Lidl link.
- [ ] "JOGHURT NAT. 2x 0,79 A" and "JOGHURT NAT." normalise to the same key.
- [ ] A user correcting the ingredient replaces the link.
