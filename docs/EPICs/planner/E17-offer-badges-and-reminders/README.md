# PLN-17 Offer badges and offer-ends reminders

Release R2. List lines are matched to this week's offers at the household's shops and get a badge ("−30% at Lidl
until Saturday"); offers can be added to the list with one tap; badges disappear when offers end; an "offer ends
tomorrow" reminder goes out (features.md §2 "Flyers and offers", "Shopping list" step 5, Verification 6).

Offers are owned by the offers service: the household's shops come from `OFF-04 (ListMyShops)` (Free: 2 shops,
enforced there) and offers from `OFF-09 (FindOffersForIngredients)`.

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [PLN-17-01](PLN-17-01-offer-badges.md) | Offer badges on list lines | R2 | M | PLN-15-01, OFF-04 (ListMyShops), OFF-09 (FindOffersForIngredients), OFF-11 (GetPrices) |
| [PLN-17-02](PLN-17-02-add-offer-to-list.md) | Add an offer to the list | R2 | S | PLN-17-01, PLN-15-08, OFF-09 (GetOffers) |
| [PLN-17-03](PLN-17-03-badge-refresh-and-expiry.md) | Badge refresh and expiry | R2 | S | PLN-17-01, PLN-21-05 |
| [PLN-17-04](PLN-17-04-offer-ends-reminders.md) | "Offer ends tomorrow" reminders | R2 | M | PLN-17-01, NTF-03 (SendPush) |
