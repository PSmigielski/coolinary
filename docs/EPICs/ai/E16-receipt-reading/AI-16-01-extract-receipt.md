# AI-16-01 ExtractReceipt

| | |
|---|---|
| Service | ai |
| Epic | [AI-16 Receipt reading](README.md) |
| Release | R3 |
| Size | M |
| Depends on | AI-06-01 |
| Unblocks | AI-16-02; KIT-08 (receipt intake) |

## Goal

Read every line of a Polish or German supermarket receipt exactly as printed, with prices and discounts.

## Business rules

1. Job `extract_receipt`, Haiku 5.5 at effort `medium`. Prompt `extract_receipt/v1`.
2. Input: 1–3 images (long receipts in parts, top to bottom; limits as AI-06-01 rule 2); `country`; `language`.
   Digital receipts shared from loyalty apps come as images or a 1–3-page PDF (≤ 5 MiB).
3. Output:
   - `chain_name_printed`, `store_city_printed` (city only), `purchased_at` (ISO date and, if printed, time),
     `currency`.
   - `lines` (0–150): `text` as printed (≤ 80, e.g. "SER GOUDA PLAST. 150G"), `quantity` (number, default 1),
     `quantity_unit` (`PIECE` or `KILOGRAM` for weighed lines like "0,512 kg x 12,99"), `unit_price` and
     `line_total` (`Money`), `discount` (`Money`, ≥ 0, from the discount line that follows), `ean` only when printed
     on the line (rare), `vat_class` letter when printed.
   - `total` (`Money`) and `paid_total` when different (coupons, deposits).
   - `deposit_lines` (bottle deposits, "Pfand", "kaucja") marked `is_deposit=true` and never treated as food.
4. Multi-image receipts: lines are read across images in order; an overlap between photos is not counted twice
   (prompt rule; AI-16-02 checks totals).
5. Returns no names of people, card numbers, loyalty numbers, terminal ids or tax ids (AI-16-02).
6. Cache 1 day by image hashes; tracing: `images`, `lines`.

## Contract

```proto
rpc ExtractReceipt(ExtractReceiptRequest) returns (ExtractReceiptResponse);

enum ReceiptQuantityUnit { RECEIPT_QUANTITY_UNIT_UNSPECIFIED = 0; RECEIPT_QUANTITY_UNIT_PIECE = 1; RECEIPT_QUANTITY_UNIT_KILOGRAM = 2; }

message ReceiptLine {
  string text = 1;
  double quantity = 2;
  ReceiptQuantityUnit quantity_unit = 3;
  coolinary.common.v1.Money unit_price = 4;
  coolinary.common.v1.Money line_total = 5;
  coolinary.common.v1.Money discount = 6;
  string ean = 7;
  string vat_class = 8;
  bool is_deposit = 9;
}

message ExtractReceiptRequest {
  repeated InputImage images = 1;   // 1-3, or
  bytes pdf = 2;                    // 1-3 pages
  string country = 3;
  string language = 4;
  string household_id = 5;
}

message ExtractReceiptResponse {
  string chain_name_printed = 1;
  string store_city_printed = 2;
  string purchased_at = 3;          // YYYY-MM-DD or YYYY-MM-DDTHH:MM
  string currency = 4;
  repeated ReceiptLine lines = 5;
  coolinary.common.v1.Money total = 6;
  coolinary.common.v1.Money paid_total = 7;
  bool total_mismatch = 8;          // AI-16-02
  bool unreadable = 9;
  CallInfo call = 10;
}
```

## Data

No schema change.

## Errors

| Case | gRPC status |
|---|---|
| Bad image count, size or type | `INVALID_ARGUMENT`, `ai_bad_request` |
| Not a receipt / unreadable | OK, `unreadable=true` |

## Acceptance criteria

- [ ] A Biedronka receipt with a weighed line "0,512 x 12,99" returns quantity 0.512 `KILOGRAM`, unit price 1299, line total 665.
- [ ] A Lidl DE receipt with a "Pfand" line marks it `is_deposit`.
- [ ] A discount line ("Rabat -1,00") is attached to the line above as `discount` 100.
