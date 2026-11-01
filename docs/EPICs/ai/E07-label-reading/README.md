# AI-07 Product label reading

| | |
|---|---|
| Service | ai |
| Release | R1 |
| Goal | When a barcode is not found anywhere, the user photographs the label and AI reads the product's name, brand and size once; the product then joins the shared catalogue (features.md §4 "How one scan is resolved", step 5). |

## Flow

```
scan -> kitchen ResolveBarcode: no hit (KIT-04) -> phone takes 1-2 label photos
  -> gateway (quota: counts as a photo scan) -> kitchen (KIT-05) -> ai.ReadProductLabel (Haiku 5.5, effort low)
  <- name, generic name, brand, pack size, best-before, storage hint
kitchen pre-fills the "name this product" form; the user confirms; the product is pending until a second household accepts it
```

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [AI-07-01](AI-07-01-read-product-label.md) | ReadProductLabel | R1 | M | AI-06-01 |
| [AI-07-02](AI-07-02-label-golden-set.md) | Golden set: 50 labels | R1 | S | AI-07-01, AI-03-03 |
