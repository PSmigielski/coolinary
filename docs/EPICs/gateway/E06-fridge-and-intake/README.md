# GW-06 Fridge, barcode and intake routes

Releases R1 (stock, barcode, naming unknown products, fridge photo, label photo) and R3 (typed or spoken sentence,
receipt photo).

## Goal

Every way of stocking the fridge (features.md §4 "Ways to stock the fridge") reaches kitchen through one consistent
set of routes. AI intake never writes stock directly: it returns **proposals** the user confirms, and the confirmed
list is sent with `POST /v1/stock` (the same route as unpack mode).

## Flow

```
typing        ──────────────────────────────────────────────► POST /v1/stock
barcode scan  ── GET /v1/products/resolve ── (unknown) POST /v1/products ─► POST /v1/stock (unpack: many at once)
photo / text / receipt ── POST /v1/intake/* ── proposals ── user edits ──► POST /v1/stock
shopping-list tick ──── planner → kitchen AddStockItems (GW-08-03)
```

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [GW-06-01](GW-06-01-stock-routes.md) | Stock: list, add, edit, remove | R1 | M | GW-01-03, KIT-07 |
| [GW-06-02](GW-06-02-barcode-resolve.md) | Barcode and 2D code lookup | R1 | S | GW-01-02, KIT-04 |
| [GW-06-03](GW-06-03-name-unknown-product.md) | Name an unknown product, confirm a shared one | R1 | S | GW-06-02, KIT-05 |
| [GW-06-04](GW-06-04-photo-intake.md) | Fridge or shelf photo → proposals | R1 | M | GW-01-05, GW-03-01, KIT-08 |
| [GW-06-05](GW-06-05-label-photo.md) | Label photo for an unknown product | R1 | S | GW-01-05, GW-03-01, KIT-05 |
| [GW-06-06](GW-06-06-text-intake.md) | Typed or spoken sentence → proposals | R3 | S | GW-03-01, KIT-08 |
| [GW-06-07](GW-06-07-receipt-intake.md) | Receipt photo → proposals | R3 | S | GW-01-05, GW-03-01, KIT-08 |
