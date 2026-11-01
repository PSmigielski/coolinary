# KIT-07-06 Stock proposals and unpack preview

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-07 Fridge stock](README.md) |
| Release | R1 |
| Size | M |
| Depends on | KIT-07-02, KIT-04 (barcode resolution incl. weighed and multipack codes), KIT-05 (unknown products), KIT-02 (pack sizes) |
| Unblocks | GW-06 (unpack mode), KIT-08-01 |

## Goal

Before anything is stored, the app shows a ready-to-confirm list: "unpack mode" turns a bag of scanned barcodes into
one list the user confirms once, and the AI intakes (KIT-08) reuse the same builder (features.md §2 "Unpack mode",
§4 "Ways to stock the fridge").

## Business rules

1. `PreviewStockItems` takes 1–60 candidates of one `source`; each candidate is an EAN, an ingredient id or a
   name, with an optional quantity and confidence. Nothing is written to the database.
2. **Barcodes** are resolved with the KIT-04 chain. Repeated scans of the same EAN in one call collapse into one
   proposal: quantity = count × pack size, or count pieces when the pack size is unknown. Weighed codes (prefix
   20–29) carry their own weight from KIT-04 and are never collapsed. Multipack codes use the unit count KIT-04
   returns.
3. An EAN that no source knows becomes a proposal with note `unknown_product`, not ticked, and the app starts the
   naming flow (KIT-05) for it.
4. **Ingredient ids** become proposals with the ingredient's name and the given quantity, or the ingredient's
   typical pack size (KIT-02) with note `estimated_quantity`, or 1 piece.
5. **Names** without an id become proposals with note `not_matched`, not ticked (KIT-08 matches names before
   calling this builder; this rule covers what is left).
6. Each proposal is pre-filled with storage place and use-by date (KIT-07-02), so confirming needs no typing.
7. `pre_ticked`: barcodes found → true; AI items → true only when confidence ≥ 0.5 and matched (KIT-08-01);
   unknown and unmatched → false.
8. Note `already_in_stock` when the household already holds the same ingredient (shown as a hint, not a block).
9. The proposal's `item` is a complete `NewStockItem`; the app sends back the ticked (and possibly edited)
   proposals to `AddStockItems` with the same `source`.

## Contract

`contracts/coolinary/kitchen/v1/stock.proto`

```proto
rpc PreviewStockItems(PreviewStockItemsRequest) returns (PreviewStockItemsResponse);

message StockCandidate {
  oneof what {
    string ean = 1;
    string ingredient_id = 2;
    string name = 3;
  }
  coolinary.common.v1.Quantity quantity = 4;
  float confidence = 5;           // 0 = not from AI
  string printed_name = 6;        // receipts only
}

enum ProposalNote {
  PROPOSAL_NOTE_UNSPECIFIED = 0;
  PROPOSAL_NOTE_UNKNOWN_PRODUCT = 1;
  PROPOSAL_NOTE_ESTIMATED_QUANTITY = 2;
  PROPOSAL_NOTE_NOT_MATCHED = 3;
  PROPOSAL_NOTE_ALREADY_IN_STOCK = 4;
  PROPOSAL_NOTE_LOW_CONFIDENCE = 5;
}

message StockProposal {
  NewStockItem item = 1;          // ready for AddStockItems
  string display_name = 2;
  bool matched = 3;
  float confidence = 4;
  bool pre_ticked = 5;
  repeated ProposalNote notes = 6;
  string printed_name = 7;
  int32 scan_count = 8;           // unpack: how many scans this line collapses
}

message PreviewStockItemsRequest {
  string household_id = 1;
  StockSource source = 2;
  repeated StockCandidate candidates = 3;   // 1–60
}

message PreviewStockItemsResponse { repeated StockProposal proposals = 1; }
```

## Data

No schema change.

## Errors

Reason codes (in backticks) travel as trailing gRPC metadata `x-error-code` (README decision 18); version
conflicts use `ABORTED`.

| Case | gRPC status | Behaviour |
|---|---|---|
| 0 or more than 60 candidates | `INVALID_ARGUMENT` | |
| Malformed EAN (check digit) | — | proposal with note `unknown_product`, not an error, so one bad scan does not lose the bag |
| Open Food Facts or GS1 down | — | KIT-04 falls back; affected codes become `unknown_product` |

## Acceptance criteria

- [ ] 12 scans: 3× the same yoghurt (150 g), 2× milk 1 l, one weighed cheese code (0.312 kg), one unknown EAN,
      5 other known products → 9 proposals: yoghurt 450 g (`scan_count` 3), milk 2,000 ml, cheese 312 g, the unknown
      one not ticked with `unknown_product`.
- [ ] Every proposal has a storage place and, where the catalogue knows a shelf life, a use-by date.
- [ ] Confirming the proposals unchanged through `AddStockItems` stores exactly the ticked ones.
- [ ] The call writes nothing (row counts unchanged).
- [ ] features.md Verification 14 (barcode handling: weighed code, 8-digit code, multipack) passes through unpack.
