# KIT-08-06 receipt.scanned event for shop prices

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-08 AI intake](README.md) |
| Release | R3 |
| Size | M |
| Status | todo |
| Depends on | KIT-08-04, KIT-01 (outbox) |
| Unblocks | OFF-11 (receipt prices, `receipt.scanned` consumer) |

## Goal

Prices that users actually paid feed the shop comparison: from release 3, an item that is not on offer is priced at
"the last price seen on receipts that users scanned at that chain" (features.md §2 "Where prices come from").

## Business rules

1. On confirmation of a receipt (`AddStockItems`, source `receipt`, receipt copy present), kitchen writes one
   `receipt.scanned` event to the outbox in the same transaction.
2. Only confirmed lines are included, and only those with a matched ingredient and a price > 0.
3. Per line: `ingredient_id`, `ean` (if a product matched), `quantity` in base units (the whole line), `item_quantity`
   (number of packs or pieces on the receipt line, e.g. 2 for "2 × 500 g"; 1 when the receipt shows no count, and
   for weighed goods), `unit_price` = line total ÷ quantity, expressed per kg, per litre or per piece, in minor
   units, rounded half up. Offers derives the pack price as line total ÷ `item_quantity` (OFF-11-02).
4. Event fields: `chain_key`, `shop_name`, `purchased_on`, `country` (household's), `currency`, lines.
5. **No household link:** the envelope's `household_id` is empty and the payload holds no household, user or
   receipt id. Decided (Q222): prices are shared facts and do not need a household link; the README envelope allows
   an empty `household_id` for such events.
6. **Duplicates:** a fingerprint = SHA-256 of (chain key, purchase date, sorted line totals). The same household
   confirming the same fingerprint again within 90 days emits no second event (fingerprints kept per household,
   90 days, deleted with the household).
7. Receipts with `date_uncertain = true` are still sent; offers decides how to weigh them.
8. Price sanity (outliers, wrong units) is checked by offers (OFF-11), not here.

## Contract

Event `receipt.scanned` (docs/EPICs/README.md), payload extended with `chain_key`, `shop_name`, `currency`.

```json
{
  "type": "receipt.scanned", "version": 1, "producer": "kitchen", "household_id": "",
  "payload": {
    "chain_key": "lidl", "shop_name": "Lidl", "purchased_on": "2027-03-04", "country": "PL", "currency": "PLN",
    "lines": [
      { "ingredient_id": "…", "ean": "5900000000000", "quantity": { "value": 1000, "unit": "MILLILITRE" },
        "item_quantity": 1, "unit_price_minor": 349, "price_basis": "per_litre" }
    ]
  }
}
```

## Data

`kitchen/migrations/000NN_receipt_fingerprints.sql`

```sql
-- +goose Up
CREATE TABLE receipt_fingerprints (
    household_id uuid NOT NULL,
    fingerprint  bytea NOT NULL,
    created_at   timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (household_id, fingerprint)
);
-- Rows older than 90 days are deleted by KIT-13-08.

-- +goose Down
DROP TABLE receipt_fingerprints;
```

## Errors

| Case | Behaviour |
|---|---|
| No matched, priced line | no event |
| Receipt copy evicted | no event; items still added |

## Acceptance criteria

- [ ] Confirming 10 of 14 lines (8 matched with prices) emits one event with 8 lines.
- [ ] Milk 1 l at 3.49 PLN → `unit_price_minor = 349`, `price_basis = per_litre`; cheese 312 g at 7.49 →
      2,401 per kg.
- [ ] Confirming the same receipt twice emits one event.
- [ ] The event payload contains no household id, user id or printed names.
