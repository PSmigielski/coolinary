# KIT-07 Fridge stock

| | |
|---|---|
| Service | kitchen |
| Release | R1 (live notice R2) |
| Depends on | KIT-01 (outbox, events), KIT-02 (ingredients, shelf life, pack sizes, unit conversion KIT-02-02), KIT-04 (products, barcode resolution), KIT-10-01 (household timezone) |

## Goal

The household's fridge, freezer and pantry as one shared list that is easy to fill, always shows what to use first,
and shrinks on its own when a dish is cooked. Planner reads it to rank dishes and to build the shopping list.

## Flow

```
typing / barcode / unpack / photo / text / receipt ──► PreviewStockItems (proposals, nothing stored)
                                                         │ user confirms
                                                         ▼
planner: list line ticked ──────────────────────────► AddStockItems ──► stock_items + stock_movements
                                                                          │
app edits ─► UpdateStockItem / RemoveStockItem (used, thrown away, mistake)│
planner: "cooked it" ─► ConsumeStock (oldest expiry first) ───────────────┤
                                                                          ▼
                                     stock.changed event (planner) + live notice (shared phones)
ListStock ◄── app, planner (with stock_version for cheap re-checks)
```

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [KIT-07-01](KIT-07-01-stock-schema.md) | Stock items and movements schema | R1 | M | KIT-01, KIT-02, KIT-04 |
| [KIT-07-02](KIT-07-02-storage-and-expiry-defaults.md) | Default storage place and use-by date | R1 | M | KIT-07-01, KIT-02, KIT-04, KIT-10-01 |
| [KIT-07-03](KIT-07-03-add-stock-items.md) | AddStockItems: add, merge, idempotent adds | R1 | M | KIT-07-01, KIT-07-02 |
| [KIT-07-04](KIT-07-04-update-and-remove.md) | Edit and remove items (used, thrown away, mistake) | R1 | M | KIT-07-01, KIT-07-02 |
| [KIT-07-05](KIT-07-05-list-stock.md) | ListStock with expiry status and totals | R1 | M | KIT-07-01, KIT-02 |
| [KIT-07-06](KIT-07-06-proposals-and-unpack.md) | Stock proposals and unpack preview | R1 | M | KIT-07-02, KIT-04 |
| [KIT-07-07](KIT-07-07-consume-stock.md) | ConsumeStock on "cooked it" | R1 | L | KIT-07-01, KIT-02, KIT-10-05 |
| [KIT-07-08](KIT-07-08-stock-changed.md) | stock.changed event and live notice | R1 (live R2) | S | KIT-07-03, KIT-07-04, KIT-07-07, KIT-01 |
| [KIT-07-09](KIT-07-09-link-unmatched-items.md) | Link items typed without a catalogue match | R1 | M | KIT-07-03, KIT-03 |

## Decisions (taken 1 November 2026)

- Q077 (KIT-07-01): Yes, 1,000 items, as an abuse guard.
- Q078 (KIT-07-02): At most 2 days, or the fridge shelf life if shorter.
- Q079 (KIT-07-03): Yes: stored as "not matched" and matched later.
- Q080 (KIT-07-07): Yes: freezer last; otherwise oldest date first.
- Q081 (KIT-07-07): Below 2% of the original, or under 5 g, 5 ml or 0.1 piece.
