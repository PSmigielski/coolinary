# AI-13-01 Offer schema and checks

| | |
|---|---|
| Service | ai |
| Epic | [AI-13 Flyer offer extraction](README.md) |
| Release | R2 |
| Size | M |
| Status | todo |
| Depends on | AI-01-06, AI-06-01 |
| Unblocks | AI-13-02, AI-13-03; OFF-07, OFF-10 |

## Goal

One offer schema that captures what Polish and German flyers actually print — multibuys, loyalty-app prices, per-kg
prices, ranges — so offers can be validated and shown with their conditions.

## Business rules

1. One page in, 0–80 offers out (`FlyerPageOffers`). The page is untrusted content (AI-01-07).
2. Each `ExtractedOffer`:
   - `product_name` as printed (≤ 120 chars), `brand` (≤ 60), `variant` (≤ 120, e.g. "różne rodzaje",
     "verschiedene Sorten").
   - `pack_size` (`Quantity`; g, ml or pieces) and `pack_size_max` for ranges ("200–250 g"); `unit_count` for
     multipacks (6 × 0.5 l → 500 ml, 6).
   - `price`: `Money` in minor units and currency; `price_basis`: `PER_PACK`, `PER_KG`, `PER_LITRE`, `PER_PIECE`
     (weighed goods are often priced per kg).
   - `regular_price` (`Money`, only when printed), `discount_percent` (only when printed).
   - `conditions`: list of `{type, text ≤ 80, buy_qty, pay_qty}` with types `MULTIBUY` ("2+1" → buy 3 pay 2;
     "drugi za 1 zł" → `OTHER` with text), `LOYALTY_APP` (Lidl Plus, Moja Biedronka, Kaufland Card app),
     `LOYALTY_CARD`, `MIN_QUANTITY`, `LIMIT_PER_CUSTOMER`, `COUPON`, `DAYS_ONLY` (e.g. weekend only), `OTHER`.
     features.md §8: prices needing a loyalty app or card are shown as a condition.
   - `valid_from` / `valid_to` only when printed for this offer and different from the flyer's dates.
   - `category`: one of `fruit_veg`, `meat`, `fish`, `dairy_eggs`, `bakery`, `frozen`, `drinks`, `pantry`,
     `sweets_snacks`, `non_food`, `other`.
   - `bbox` (x, y, w, h normalised 0–1) of the offer on the page, for admin spot-checks (OFF-08).
   - `confidence` 0–1.
3. Prices are read exactly as printed: "4,99" → 499; "1.29 €" → 129; a price shown only as "−30%" without an amount
   → no price, the offer is dropped (`dropped_reason=NO_PRICE`).
4. Checks in ai (input-independent, so they run in batch too): price > 0 and ≤ 1,000,000 minor units; currency in
   `PLN`, `EUR`, `CZK`, `HUF`, `RON`; `regular_price` > `price` when both present; multibuy `pay_qty` < `buy_qty`;
   dates valid ISO; `bbox` inside 0–1; strings within limits. Failing offers are dropped (AI-01-06 rule 3) with
   `ai_output_dropped_total{job="extract_flyer_offers", rule}`.
5. **Caller checks** (OFF-07, features.md §3 step 5): currency matches the chain's country; dates inside the flyer's
   validity; unit present; product matched to the catalogue; duplicates across pages merged.
6. AI never outputs prices it computed (e.g. per-kg from pack price); only printed values.

## Contract

`contracts/coolinary/ai/v1/flyers.proto`:

```proto
enum PriceBasis { PRICE_BASIS_UNSPECIFIED = 0; PRICE_BASIS_PER_PACK = 1; PRICE_BASIS_PER_KG = 2;
  PRICE_BASIS_PER_LITRE = 3; PRICE_BASIS_PER_PIECE = 4; }
enum ConditionType { CONDITION_TYPE_UNSPECIFIED = 0; CONDITION_TYPE_MULTIBUY = 1; CONDITION_TYPE_LOYALTY_APP = 2;
  CONDITION_TYPE_LOYALTY_CARD = 3; CONDITION_TYPE_MIN_QUANTITY = 4; CONDITION_TYPE_LIMIT_PER_CUSTOMER = 5;
  CONDITION_TYPE_COUPON = 6; CONDITION_TYPE_DAYS_ONLY = 7; CONDITION_TYPE_OTHER = 8; }

message OfferCondition {
  ConditionType type = 1;
  string text = 2;
  int32 buy_qty = 3;
  int32 pay_qty = 4;
}

message BoundingBox { float x = 1; float y = 2; float w = 3; float h = 4; }

message ExtractedOffer {
  string product_name = 1;
  string brand = 2;
  string variant = 3;
  coolinary.common.v1.Quantity pack_size = 4;
  coolinary.common.v1.Quantity pack_size_max = 5;
  int32 unit_count = 6;
  coolinary.common.v1.Money price = 7;
  PriceBasis price_basis = 8;
  coolinary.common.v1.Money regular_price = 9;
  int32 discount_percent = 10;
  repeated OfferCondition conditions = 11;
  string valid_from = 12;
  string valid_to = 13;
  string category = 14;
  BoundingBox bbox = 15;
  float confidence = 16;
  int32 page_number = 17;
}

message FlyerMeta {
  string chain = 1;           // chain code or name, "" when unknown (user flyers)
  string country = 2;         // ISO 3166-1 alpha-2
  string currency = 3;        // expected currency, "" when unknown
  string valid_from = 4;      // flyer validity from the list page, "" when unknown
  string valid_to = 5;
  string language = 6;
}

message FlyerPageOffers {
  string page_ref = 1;
  int32 page_number = 2;
  repeated ExtractedOffer offers = 3;
  int32 dropped = 4;
}
```

## Data

No schema change.

## Errors

Offers failing a check are dropped and counted (rule 4); the page result is still returned.

## Acceptance criteria

- [ ] A Lidl DE page with "Lidl Plus" prices yields offers with a `LOYALTY_APP` condition; a Biedronka "2+1" offer yields `MULTIBUY` buy 3 pay 2.
- [ ] "−30%" without a printed price is dropped with `NO_PRICE`.
- [ ] A weighed meat offer at "19,99 zł/kg" has `PER_KG`.
