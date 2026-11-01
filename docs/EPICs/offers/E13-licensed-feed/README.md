# OFF-13 Licensed offers feed (ShopFully)

| | |
|---|---|
| Service | offers |
| Release | R3 |
| Goal | Take offers from a licensed offers network such as ShopFully (formerly Offerista; 25+ countries) once there is traffic to show, to cover per-store chains that block automated reading (Rewe, Edeka, Netto DE) and more countries, and possibly earn from flyer views (features.md §3 source 3, business model §4 "Offers network"). |

`GATE:` nothing in this epic goes live before a signed agreement that settles data use, attribution, display rules,
reporting duties and price. Until then the tasks run only against the provider's sandbox or sample files.

## Flow

1. Feed retailers and stores are mapped to Coolinary chains and shops.
2. Feed offers arrive already structured, so no AI reading is needed; they go through the same normalisation,
   matching, validation and publishing as flyer offers.
3. Per-store chains get store editions from the feed, read on demand like Kaufland.
4. Views and clicks are reported back if the agreement requires it.

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [OFF-13-01](OFF-13-01-feed-client-mapping.md) | Feed client and retailer/store mapping | R3 | L | OFF-02-01, OFF-03-05 |
| [OFF-13-02](OFF-13-02-feed-ingest.md) | Ingest feed offers | R3 | L | OFF-13-01, OFF-07-03, OFF-07-04, OFF-07-05, OFF-07-07 |
| [OFF-13-03](OFF-13-03-per-store-chains-via-feed.md) | Per-store chains through the feed (Rewe, Edeka, Netto DE) | R3 | M | OFF-13-02, OFF-05-07 |
| [OFF-13-04](OFF-13-04-feed-reporting.md) | View and click reporting to the network | R3 | M | OFF-13-02 |
