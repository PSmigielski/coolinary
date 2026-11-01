# AI-16 Receipt reading

| | |
|---|---|
| Service | ai |
| Release | R3 |
| Goal | A receipt photo becomes the items bought (to stock the fridge) and the prices paid at that chain (for the trip plan's "last price seen on receipts", features.md §2 "Planning the shopping trip", "Where prices come from"; §4 "Receipts"). |

## Flow

```
phone photo(s) of a receipt -> gateway (quota: photo scan) -> kitchen KIT-08 -> ai.ExtractReceipt (Haiku 5.5 medium)
  <- chain as printed, date, currency, lines {text as printed, quantity, unit price, line total, discount}, total
kitchen matches lines to products this household scanned at that shop, then the catalogue (KIT-03);
proposed items to confirm; prices go to offers through the receipt.scanned event
```

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [AI-16-01](AI-16-01-extract-receipt.md) | ExtractReceipt | R3 | M | AI-06-01 |
| [AI-16-02](AI-16-02-receipt-privacy-and-totals.md) | Receipt privacy filter and totals check | R3 | S | AI-16-01 |
| [AI-16-03](AI-16-03-receipt-golden-set.md) | Golden set: 40 receipts | R3 | M | AI-16-02, AI-03-03 |
