# KIT-08-04 Receipt intake and line matching

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-08 AI intake](README.md) |
| Release | R3 |
| Size | L |
| Depends on | KIT-08-01, AI-16 (`ExtractReceipt`), KIT-04 (household products), KIT-03 (`MatchNames`), KIT-08-05 (memory lookup) |
| Unblocks | GW-06 (`POST /v1/intake/receipt`), KIT-08-05, KIT-08-06 |

## Goal

A receipt photo, or a digital receipt image or PDF shared from a loyalty app, becomes the items to put away, with
the price actually paid, matched first to products this household has bought before (features.md §2, §4 "Receipts").

## Business rules

1. Input: 1–3 images (a long receipt in parts; same limits as KIT-08-02) or one PDF (≤ 5 MB, ≤ 10 pages), the
   language and the household's country (KIT-10-01).
2. `ExtractReceipt` (AI-16, Haiku 5.5), deadline 45 s, returns: shop name, store address hint, purchase date,
   currency, total, and lines with printed name, quantity, unit (pieces or kg), unit price, line total after
   discounts, `is_food`.
3. **Lines left out:** `is_food = false`; deposits ("Pfand", "Leergut", "kaucja"); carrier bags; totals,
   subtotals, payment and tax lines. Discount lines that AI did not fold into a line are subtracted from the line
   directly above. The count of left-out lines is returned.
4. **Chain:** the shop name is normalised to a `chain_key` (lower case, ASCII, no legal suffix: "Lidl sp. z o.o."
   → `lidl`, "ALDI SÜD" → `aldi-sued`, "Kaufland" → `kaufland`). Kitchen does not call offers (no edge); offers maps
   keys to its chains (OFF-11). Unknown shops keep their normalised name as key.
5. **Purchase date:** the receipt's date; if missing, in the future, or more than 30 days ago → today, with
   `date_uncertain = true`.
6. **Matching, per line, first hit wins:**
   1. this household's receipt-name memory for this chain (KIT-08-05) → product and/or ingredient, confidence 1.0;
   2. this household's own products (KIT-04): trigram similarity ≥ 0.5 between the printed name and the product
      name (receipt names are abbreviated, e.g. "MLEKO UHT 3,2% 1L") → product, confidence = similarity;
   3. `MatchNames` (KIT-03, context `receipt`) → ingredient only.
7. **Quantity:** weighed lines (0.734 kg) become 734 g; counted lines become count × the product's pack size when a
   product matched, else count pieces.
8. **Price paid:** the line total after discounts, in the receipt currency, goes into each proposal's
   `price_paid`. No offer saving is guessed.
9. Proposals go through KIT-08-01 (pre-tick rule, storage defaults) and carry `printed_name`.
10. The extracted receipt (chain key, date, currency, per-line printed name, price, quantity and match) is kept in
    `redis-cache` under `kitchen:receipt:<intake_ref>` for 2 hours, so confirmation can teach the memory and emit
    prices (KIT-08-05, KIT-08-06); the key is also added to the set `kitchen:receipts:<household_id>` (same TTL) so
    household deletion can remove it (KIT-13-04). No image is kept. If the entry is evicted, the items are still added; only the
    memory update and the price event are skipped. Decided (Q221): the server keeps a 2-hour copy of the extracted
    lines (never the image), so prices can't be forged by the app.
11. AI unavailable → `UNAVAILABLE` / `ai_unavailable` (no fallback). AI says "not a receipt" →
    `INVALID_ARGUMENT` with reason `not_a_receipt`; `ai_used = true` is reported in the error details so the gateway
    counts the scan (GW-06).

## Contract

`contracts/coolinary/kitchen/v1/intake.proto`

```proto
rpc IntakeFromReceipt(IntakeFromReceiptRequest) returns (ReceiptIntakeResponse);

message IntakeFromReceiptRequest {
  string household_id = 1;
  repeated IntakeImage images = 2;   // 1–3, or empty when pdf is set
  bytes pdf = 3;                     // ≤ 5 MB, ≤ 10 pages
  string language = 4;
}

message ReceiptInfo {
  string chain_key = 1;
  string shop_name = 2;
  string purchased_on = 3;           // YYYY-MM-DD
  bool date_uncertain = 4;
  coolinary.common.v1.Money total = 5;
  int32 left_out_lines = 6;
}

message ReceiptIntakeResponse {
  IntakeResponse intake = 1;         // proposals carry printed_name and price_paid
  ReceiptInfo receipt = 2;
}
```

## Data

No schema change (memory table in KIT-08-05; short-lived copy in `redis-cache`).

## Errors

| Case | gRPC status | Reason |
|---|---|---|
| Both images and PDF, neither, too many, too large | `INVALID_ARGUMENT` | |
| Not a receipt | `INVALID_ARGUMENT` | `not_a_receipt` (details: `ai_used = true`) |
| ai unavailable | `UNAVAILABLE` | `ai_unavailable` |

## Acceptance criteria

- [ ] Golden Lidl receipt (AI-03): 14 lines, 1 bag and 2 deposits left out (`left_out_lines = 3`), chain key `lidl`.
- [ ] A line "-0,50 Rabatt" under "Joghurt 1,29" gives a yoghurt proposal at 0.79.
- [ ] "Käse lose 0,312 kg" gives 312 g.
- [ ] A household that scanned the barcode of "Mleko UHT 3,2% 1 l" before gets the receipt line "MLEKO UHT 3,2% 1L"
      matched to that product.
- [ ] A receipt dated in the future is set to today with `date_uncertain = true`.
- [ ] After 2 hours the Redis copy is gone (TTL checked).
