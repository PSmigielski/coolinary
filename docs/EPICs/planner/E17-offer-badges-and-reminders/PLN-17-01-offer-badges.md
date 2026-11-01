# PLN-17-01 Offer badges on list lines

| | |
|---|---|
| Service | planner |
| Epic | [PLN-17 Offer badges and offer-ends reminders](README.md) |
| Release | R2 |
| Size | M |
| Depends on | PLN-15-01, OFF-04 (ListMyShops), OFF-09 (FindOffersForIngredients), OFF-11 (GetPrices) |
| Unblocks | PLN-17-02, PLN-17-03, PLN-17-04, PLN-15-07, PLN-18-04 |

## Goal

Match each list line to this week's offers at the household's shops and show the best one as a badge
(features.md "Shopping list" step 5).

## Business rules

1. Only open ingredient lines get badges; free-text lines and ticked lines don't.
2. Shops: the household's shops from `ListMyShops`, in the stored order (nearest first, set by the phone, D13).
   Free households get at most 2 shops from offers; planner does not enforce it again.
3. Offers: `FindOffersForIngredients(ingredient ids, shop ids, today)` in the household's time zone. Offers valid
   today count as current. Offers starting within the next 3 days count as `upcoming` and are shown only when no
   current offer exists for that line.
4. Per line and shop keep the offer with the lowest unit price (price ÷ pack amount in base unit). Up to 3 badges per
   line (best shop first); the primary badge is the lowest unit price overall; ties go to the earlier shop in the
   household's order.
5. Discount: if the offer states a regular price, `discount = 1 − offer / regular`; else, if OFF-11 `GetPrices`
   gives a reference unit price, `discount = 1 − offer unit price / reference unit price`. Shown only when ≥ 5%,
   rounded down to a whole percent. Otherwise the badge shows the price only.
6. Badge fields: offer id, shop id, chain and shop name, price, pack, unit price, discount (optional), conditions
   text and flags (`app_or_card`, `multi_buy`), valid from/to, source flyer URL, `upcoming`, `pinned`.
7. A pinned badge (from PLN-17-02 or a deals plan, PLN-18-04) is kept as the primary until its offer ends, even if a
   cheaper offer appears; the cheaper one is listed second.
8. Badges are computed after each list change that adds or changes ingredient lines (background, within 5 seconds)
   and stored on the line; `GetList` returns stored badges and never blocks on offers.
9. offers down: badges stay as last computed; the list works without them.
10. Badge text is built by the app from the fields, in the user's language ("−30% at Lidl until Saturday").

## Contract

`contracts/coolinary/planner/v1/list.proto`

```proto
message OfferBadge {
  string offer_id = 1; string shop_id = 2; string chain_name = 3; string shop_name = 4;
  coolinary.common.v1.Money price = 5; coolinary.common.v1.Quantity pack = 6;
  coolinary.common.v1.Money unit_price = 7;           // per kg, l or piece
  optional int32 discount_percent = 8;
  string conditions = 9; repeated string condition_flags = 10;
  string valid_from = 11; string valid_to = 12; string source_url = 13;
  bool upcoming = 14; bool pinned = 15; bool primary = 16;
}
```

## Data

`planner/migrations/000NN_list_line_offers.sql`

```sql
CREATE TABLE list_line_offers (
    line_id          uuid NOT NULL REFERENCES list_lines ON DELETE CASCADE,
    offer_id         text NOT NULL,
    shop_id          text NOT NULL,
    position         smallint NOT NULL,
    payload          jsonb NOT NULL,          -- OfferBadge fields
    valid_to         date NOT NULL,
    pinned           boolean NOT NULL DEFAULT false,
    computed_at      timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (line_id, offer_id)
);
CREATE INDEX list_line_offers_valid_to ON list_line_offers (valid_to);
CREATE INDEX list_line_offers_offer ON list_line_offers (offer_id);
ALTER TABLE shopping_lists ADD COLUMN badges_stale boolean NOT NULL DEFAULT false;
```

## Errors

| Case | Behaviour |
|---|---|
| offers down | badges kept; `badges_stale = true`, retried by PLN-17-03 |

## Acceptance criteria

- [ ] Chicken thighs on the list, Lidl offer 3.89 for 1 kg with regular price 5.59: badge "−30% at Lidl until
      Saturday" data (discount 30).
- [ ] Two shops with offers for rice: the lower unit price is primary; equal unit prices → the nearer shop.
- [ ] Free household with 2 shops: no badge from a third shop ever appears.
- [ ] Five offers added from a flyer appear on the list with shop, price and end date (Verification 6).

## Notes

- Decided (Q177, Q176): offers starting within 3 days show as `upcoming` badges (rule 3); the smallest discount shown
  is 5%, rounded down to a whole percent, else the price only (rule 5).
