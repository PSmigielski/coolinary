# offers — epics and tasks

Go service. Shops, flyer sync, offers, prices and flyer editions by area. Offers arrive in release 2; release 3 adds
receipt prices, the next chains and the licensed feed. The server never learns where a user is: it gets a postcode,
returns shops with positions, and the phone computes distances (docs/EPICs/README.md D13).

79 tasks: 24 S, 45 M, 10 L (roughly 180 developer-days). 68 in R2, 11 in R3.

## Epics

| Epic | Title | Rel. | Tasks |
|---|---|---|---|
| [OFF-01](E01-platform/README.md) | Platform: events, clients, plan lookup, jobs, admin service | R2 | 6 |
| [OFF-02](E02-chains-and-shops/README.md) | Chains and shops from OpenStreetMap | R2 | 8 |
| [OFF-03](E03-areas-and-editions/README.md) | Postcode areas and flyer editions | R2 | 5 |
| [OFF-04](E04-my-shops/README.md) | My shops and shopping preferences | R2 | 5 |
| [OFF-05](E05-flyer-sync/README.md) | Flyer source framework and sync scheduler | R2 | 8 |
| [OFF-06](E06-launch-chains/README.md) | Launch chain adapters | R2 | 7 |
| [OFF-07](E07-offer-pipeline/README.md) | Offer extraction, validation and publishing | R2 | 8 |
| [OFF-08](E08-accuracy-gate/README.md) | Accuracy gate and weekly spot-checks | R2 | 5 |
| [OFF-09](E09-offer-browse/README.md) | Offer browse, search and lookup by ingredient | R2 | 5 |
| [OFF-10](E10-user-flyers/README.md) | User-supplied flyers (Plus only) | R2 | 5 |
| [OFF-11](E11-prices/README.md) | Prices: flyer, receipt and reference | R2–R3 | 4 |
| [OFF-12](E12-next-chains/README.md) | Next chain adapters | R3 | 6 |
| [OFF-13](E13-licensed-feed/README.md) | Licensed offers feed (ShopFully) | R3 | 4 |
| [OFF-14](E14-household-lifecycle/README.md) | Household lifecycle: event consumers and data export | R2 | 3 |

## What other services use

| RPC / event | Task | Used by |
|---|---|---|
| `ResolveArea`, `ListShopsNearArea` | OFF-03-02, OFF-03-03 | gateway (GW-10) |
| `SubmitShop`, `ProposeShopCorrection` | OFF-02-05, OFF-02-06 | gateway |
| `SetMyShops`, `ListMyShops` | OFF-04-01 | gateway; `ListMyShops` also planner (PLN-17, 18, 19) |
| `GetShoppingPreferences`, `SetShoppingPreferences` | OFF-04-02 | gateway |
| `GetFlyerStatus` | OFF-05-06 | gateway |
| `BrowseOffers` (replaces `ListOffers`), `GetOffer`, `ListOfferCategories` | OFF-09-01, 03, 04 | gateway |
| `FindOffersForIngredients` | OFF-09-05 | planner (PLN-17, PLN-18) |
| `GetPrices` | OFF-11-04 | planner (PLN-19) |
| `UploadUserFlyer`, `SubmitUserFlyerLink`, `GetUserFlyer`, `ConfirmUserFlyer`, `DeleteUserFlyer` | OFF-10-01 … 04 | gateway |
| `RecordFeedUsage` | OFF-13-04 | gateway (R3, only if the feed agreement requires it) |
| `ExportHouseholdData` | OFF-14-03 | gateway (GW-02) |
| `OffersAdminService` (jobs, chains, sources, flyers, pages, shop review, store links, held offers, samples, go-live, remove chain, accuracy, reference prices, feed mapping) | OFF-01-06 and the admin task of each epic | admin (ADM-06, ADM-07) |
| `offers.published`, `offers.expired` | OFF-07-07, OFF-08-05, OFF-10-04 | planner |

offers calls: identity `GetEntitlements` (plan checks, Decided Q004); kitchen `MatchNames`, `SearchIngredients`,
`GetIngredients`; ai `SubmitBatch`, `GetBatch`, `GetBatchResults`, `ExtractFlyerOffers`; notifier `SendAlert`. It
consumes `household.deleted`, `account.deleted`, `catalogue.ingredient_merged` and, from R3, `receipt.scanned`.

## Build order

**A. Foundation**

| # | Task | Title | Rel. | Size |
|---|---|---|---|---|
| 1 | [OFF-01-01](E01-platform/OFF-01-01-event-outbox-publisher.md) | Event outbox and Redis Streams publisher | R2 | M |
| 2 | [OFF-01-02](E01-platform/OFF-01-02-event-consumer.md) | Event consumer group and idempotent handlers | R2 | M |
| 3 | [OFF-01-04](E01-platform/OFF-01-04-service-clients.md) | Clients for identity, kitchen, ai and notifier | R2 | S |
| 4 | [OFF-01-03](E01-platform/OFF-01-03-household-plan-copy.md) | Household plan lookup through identity | R2 | S |
| 5 | [OFF-01-05](E01-platform/OFF-01-05-scheduled-jobs.md) | Scheduled job runner | R2 | M |
| 6 | [OFF-01-06](E01-platform/OFF-01-06-admin-service-shell.md) | OffersAdminService shell and change log | R2 | S |
| 7 | [OFF-02-01](E02-chains-and-shops/OFF-02-01-chains-seed.md) | Chains table and seed | R2 | S |
| 8 | [OFF-02-02](E02-chains-and-shops/OFF-02-02-shops-schema.md) | Shops schema rework | R2 | S |
| 9 | [OFF-03-04](E03-areas-and-editions/OFF-03-04-editions-model.md) | Flyer editions and shop-to-edition rules | R2 | M |
| 10 | [OFF-03-01](E03-areas-and-editions/OFF-03-01-postcode-dataset.md) | Postcode and town dataset import | R2 | M |
| 11 | [OFF-03-02](E03-areas-and-editions/OFF-03-02-resolve-area.md) | Resolve a postcode or town to an area | R2 | S |
| 12 | [OFF-05-03](E05-flyer-sync/OFF-05-03-polite-fetcher.md) | Polite fetcher | R2 | M |
| 13 | [OFF-05-01](E05-flyer-sync/OFF-05-01-sources-and-flyers-schema.md) | Flyer sources, flyers and versions | R2 | M |
| 14 | [OFF-05-02](E05-flyer-sync/OFF-05-02-adapter-interface.md) | Adapter interface and fixture harness | R2 | M |

**B. Shops and my shops**

| # | Task | Title | Rel. | Size |
|---|---|---|---|---|
| 15 | [OFF-02-03](E02-chains-and-shops/OFF-02-03-osm-import.md) | OSM shop import job | R2 | L |
| 16 | [OFF-02-04](E02-chains-and-shops/OFF-02-04-shop-dedupe.md) | Merge duplicate shops | R2 | S |
| 17 | [OFF-02-08](E02-chains-and-shops/OFF-02-08-osm-attribution-gate.md) | OSM attribution and ODbL gate | R2 | S |
| 18 | [OFF-03-03](E03-areas-and-editions/OFF-03-03-shops-near-area.md) | Shops near an area | R2 | M |
| 19 | [OFF-03-05](E03-areas-and-editions/OFF-03-05-chain-store-numbers.md) | Link shops to chain store numbers | R2 | M |
| 20 | [OFF-04-01](E04-my-shops/OFF-04-01-my-shops.md) | My shops: set and list | R2 | M |
| 21 | [OFF-04-02](E04-my-shops/OFF-04-02-shopping-preferences.md) | Shopping preferences ("How you shop") | R2 | S |
| 22 | [OFF-04-03](E04-my-shops/OFF-04-03-edition-demand.md) | Store-edition demand | R2 | S |
| 23 | [OFF-04-04](E04-my-shops/OFF-04-04-plan-change-rules.md) | My shops on plan changes | R2 | S |
| 24 | [OFF-04-05](E04-my-shops/OFF-04-05-closed-merged-shops.md) | Closed and merged shops in my shops | R2 | S |
| 25 | [OFF-02-05](E02-chains-and-shops/OFF-02-05-user-adds-shop.md) | User adds a missing shop | R2 | M |
| 26 | [OFF-02-06](E02-chains-and-shops/OFF-02-06-user-corrects-shop.md) | User proposes a shop correction | R2 | S |
| 27 | [OFF-02-07](E02-chains-and-shops/OFF-02-07-shop-moderation-admin.md) | Shop moderation admin RPCs | R2 | S |

**C. Flyer pipeline**

| # | Task | Title | Rel. | Size |
|---|---|---|---|---|
| 28 | [OFF-05-04](E05-flyer-sync/OFF-05-04-sync-and-detection.md) | Sync checks and change detection | R2 | L |
| 29 | [OFF-05-05](E05-flyer-sync/OFF-05-05-page-storage.md) | Flyer page storage and retention | R2 | M |
| 30 | [OFF-05-06](E05-flyer-sync/OFF-05-06-source-health-status.md) | Source health, alerts and flyer status for users | R2 | M |
| 31 | [OFF-05-07](E05-flyer-sync/OFF-05-07-on-demand-editions.md) | Read store editions only on demand | R2 | S |
| 32 | [OFF-05-08](E05-flyer-sync/OFF-05-08-sources-admin.md) | Sources and flyers admin RPCs | R2 | M |
| 33 | [OFF-07-01](E07-offer-pipeline/OFF-07-01-offers-schema.md) | Offers schema rework | R2 | M |
| 34 | [OFF-07-02](E07-offer-pipeline/OFF-07-02-batch-extraction.md) | Submit flyers for extraction and collect results | R2 | L |
| 35 | [OFF-07-03](E07-offer-pipeline/OFF-07-03-normalise-offers.md) | Normalise prices, packs and conditions | R2 | L |
| 36 | [OFF-07-08](E07-offer-pipeline/OFF-07-08-duplicate-offers.md) | Duplicate and conflicting offers | R2 | S |
| 37 | [OFF-07-04](E07-offer-pipeline/OFF-07-04-match-products.md) | Match offer products to catalogue ingredients | R2 | M |
| 38 | [OFF-07-05](E07-offer-pipeline/OFF-07-05-validation-rules.md) | Validation without AI | R2 | M |
| 39 | [OFF-07-06](E07-offer-pipeline/OFF-07-06-hold-back-queue.md) | Hold-back queue and admin review | R2 | M |
| 40 | [OFF-07-07](E07-offer-pipeline/OFF-07-07-publish-expire.md) | Publish, "next week" and expiry | R2 | L |

**D. Launch chains and the accuracy gate**

| # | Task | Title | Rel. | Size |
|---|---|---|---|---|
| 41 | [OFF-06-01](E06-launch-chains/OFF-06-01-biedronka-pl.md) | Biedronka (PL) adapter | R2 | M |
| 42 | [OFF-06-02](E06-launch-chains/OFF-06-02-lidl-pl.md) | Lidl PL adapter | R2 | M |
| 43 | [OFF-06-03](E06-launch-chains/OFF-06-03-kaufland-pl.md) | Kaufland PL adapter (per store) | R2 | L |
| 44 | [OFF-06-04](E06-launch-chains/OFF-06-04-lidl-de.md) | Lidl DE adapter | R2 | M |
| 45 | [OFF-06-05](E06-launch-chains/OFF-06-05-aldi-sued-de.md) | Aldi Süd adapter | R2 | M |
| 46 | [OFF-06-06](E06-launch-chains/OFF-06-06-aldi-nord-de.md) | Aldi Nord adapter (PDF) | R2 | M |
| 47 | [OFF-06-07](E06-launch-chains/OFF-06-07-kaufland-de.md) | Kaufland DE adapter (per store) | R2 | S |
| 48 | [OFF-08-01](E08-accuracy-gate/OFF-08-01-review-samples.md) | Review samples and verdicts | R2 | M |
| 49 | [OFF-08-02](E08-accuracy-gate/OFF-08-02-go-live-gate.md) | Go-live gate per chain | R2 | M |
| 50 | [OFF-08-03](E08-accuracy-gate/OFF-08-03-weekly-spot-checks.md) | Weekly spot-checks | R2 | S |
| 51 | [OFF-08-04](E08-accuracy-gate/OFF-08-04-accuracy-record.md) | Accuracy record and metrics | R2 | S |
| 52 | [OFF-08-05](E08-accuracy-gate/OFF-08-05-remove-chain.md) | Remove a chain within a day | R2 | M |

**E. Read side, prices and lifecycle**

| # | Task | Title | Rel. | Size |
|---|---|---|---|---|
| 53 | [OFF-09-01](E09-offer-browse/OFF-09-01-browse-offers.md) | Browse offers of my shops | R2 | L |
| 54 | [OFF-09-02](E09-offer-browse/OFF-09-02-search-offers.md) | Search offers by text | R2 | M |
| 55 | [OFF-09-03](E09-offer-browse/OFF-09-03-offer-detail.md) | Offer detail | R2 | S |
| 56 | [OFF-09-04](E09-offer-browse/OFF-09-04-offer-categories.md) | Offer categories with counts | R2 | S |
| 57 | [OFF-09-05](E09-offer-browse/OFF-09-05-find-offers-for-ingredients.md) | Find offers for ingredients (planner) | R2 | M |
| 58 | [OFF-11-01](E11-prices/OFF-11-01-price-observations.md) | Price observations from flyers | R2 | M |
| 59 | [OFF-11-03](E11-prices/OFF-11-03-reference-prices.md) | Reference prices per country | R2 | M |
| 60 | [OFF-11-04](E11-prices/OFF-11-04-get-prices.md) | GetPrices for the trip price matrix | R2 (receipt step R3) | M |
| 61 | [OFF-14-01](E14-household-lifecycle/OFF-14-01-erase-on-deletion.md) | Erase data on household or account deletion | R2 | M |
| 62 | [OFF-14-02](E14-household-lifecycle/OFF-14-02-ingredient-merged.md) | Follow merged catalogue ingredients | R2 | S |
| 63 | [OFF-14-03](E14-household-lifecycle/OFF-14-03-export-household-data.md) | Export household data | R2 | S |

**F. User-supplied flyers**

| # | Task | Title | Rel. | Size |
|---|---|---|---|---|
| 64 | [OFF-10-05](E10-user-flyers/OFF-10-05-shop-for-user-flyer.md) | Shop for a user flyer, in any country | R2 | S |
| 65 | [OFF-10-01](E10-user-flyers/OFF-10-01-upload-and-jobs.md) | Upload photos or a PDF | R2 | M |
| 66 | [OFF-10-02](E10-user-flyers/OFF-10-02-flyer-link.md) | Flyer from a link | R2 | M |
| 67 | [OFF-10-03](E10-user-flyers/OFF-10-03-extract-draft.md) | Extract and check a draft | R2 | M |
| 68 | [OFF-10-04](E10-user-flyers/OFF-10-04-confirm-private-offers.md) | Confirm and save private offers | R2 | M |

**G. Release 3**

| # | Task | Title | Rel. | Size |
|---|---|---|---|---|
| 69 | [OFF-11-02](E11-prices/OFF-11-02-receipt-prices.md) | Receipt prices from scanned receipts | R3 | M |
| 70 | [OFF-12-01](E12-next-chains/OFF-12-01-dino-pl.md) | Dino (PL) adapter | R3 | M |
| 71 | [OFF-12-02](E12-next-chains/OFF-12-02-netto-pl.md) | Netto PL adapter | R3 | M |
| 72 | [OFF-12-03](E12-next-chains/OFF-12-03-carrefour-pl.md) | Carrefour PL adapter (per store format) | R3 | L |
| 73 | [OFF-12-04](E12-next-chains/OFF-12-04-auchan-pl.md) | Auchan PL adapter | R3 | M |
| 74 | [OFF-12-05](E12-next-chains/OFF-12-05-aldi-pl.md) | Aldi PL adapter | R3 | M |
| 75 | [OFF-12-06](E12-next-chains/OFF-12-06-penny-de.md) | Penny DE adapter | R3 | M |
| 76 | [OFF-13-01](E13-licensed-feed/OFF-13-01-feed-client-mapping.md) | Feed client and retailer/store mapping | R3 | L |
| 77 | [OFF-13-02](E13-licensed-feed/OFF-13-02-feed-ingest.md) | Ingest feed offers | R3 | L |
| 78 | [OFF-13-03](E13-licensed-feed/OFF-13-03-per-store-chains-via-feed.md) | Per-store chains through the feed (Rewe, Edeka, Netto DE) | R3 | M |
| 79 | [OFF-13-04](E13-licensed-feed/OFF-13-04-feed-reporting.md) | View and click reporting to the network | R3 | M |

## Gates

Work that must not go live until an outside condition is met.

| Task | Gate |
|---|---|
| OFF-02-08 | Legal check: is the shop table an ODbL "derivative database" with user corrections, and what must be shared. Shop data stays in dev/staging until answered. |
| OFF-05-01 | No source can be enabled without its retailer-terms approval recorded in admin. |
| OFF-06-01 | Biedronka terms checked |
| OFF-06-02 | Lidl PL terms checked |
| OFF-06-03 | Kaufland PL terms checked (also covers using its store finder) |
| OFF-06-04 | Lidl DE terms checked |
| OFF-06-05 | Aldi Süd terms checked |
| OFF-06-06 | Aldi Nord terms checked |
| OFF-06-07 | Kaufland DE terms checked (also covers its store finder) |
| OFF-08-02 | Each chain goes live only with terms approved and a 100-offer sample at ≥ 97%. |
| OFF-12-01 … 06 | Dino, Netto PL, Carrefour PL, Auchan PL, Aldi PL, Penny DE terms checked |
| OFF-13-01 | Signed feed agreement (data use, attribution, display, reporting, price) |
| OFF-13-04 | Built only if the feed agreement requires usage reporting |

## Decisions

Final answers (docs/EPICs/DECISIONS.md), written into the tasks as `Decided (Qnnn)`.

| Decision | Rule | Tasks |
|---|---|---|
| Q004 | Plan checks ask identity `GetEntitlements` per request (memoised within it) and per household per job run; no plan copy, no `household.plan_changed` handler. | OFF-01-03, OFF-01-04, OFF-04-01, OFF-04-04, OFF-09-01, OFF-10-01 |
| Q010 | Chain flyer pages kept until their version's end date (or withdrawal) + 90 days. | OFF-05-05, OFF-08-01 |
| Q022 | Own-flyer PDFs up to 40 pages, upload up to 25 MB. | OFF-10-01 |
| Q144 | Shops from the Overpass API, one query per chain per country, monthly. | OFF-02-03 |
| Q145 | A rejected user-added shop is hidden for everyone and removed from the adder's my shops. | OFF-02-05 |
| Q146 | Shop corrections are never sent to OpenStreetMap automatically; the founder edits OSM by hand. | OFF-02-06 |
| Q147 | GeoNames postal codes, CC BY 4.0, visible "GeoNames" credit, refreshed yearly. | OFF-03-01 |
| Q148 | Aldi Nord: national flyer only in R2; regional extras not read. | OFF-03-04, OFF-06-06 |
| Q149 | At most 1,000 shops per area response. | OFF-03-03 |
| Q150 | Kaufland store numbers from its public store finder, under the retailer-terms gate. | OFF-03-05 |
| Q151 | 3 s between requests to one host, or the site's `Crawl-delay` if larger. | OFF-05-03 |
| Q152 | A flyer missing twice before its start is withdrawn; once started it is kept to its end date. | OFF-05-04 |
| Q153 | Flyer pages on the Docker volume `offers-flyers` at `/data/flyers`. | OFF-05-05 |
| Q154 | Flyers fetched < 36 h before start are read at once at full price. | OFF-07-02 |
| Q155 | Non-food offers are dropped, never stored. | OFF-07-03 |
| Q156 | Offer categories are the catalogue's ingredient categories. | OFF-07-04 |
| Q157 | Hold offers below 0.6 AI confidence or outside 0.2×–5× of the reference price. | OFF-07-05 |
| Q158 | Bad weekly spot-checks (two weeks < 97% or one < 90%) only send a critical alert; no automatic step back. | OFF-08-03 |
| Q160 | Biedronka flyers with only a start date end at start + 6 days; offer dates may narrow it. | OFF-06-01 |
| Q161 | Plus may save at most 40 shops; Free 2. | OFF-04-01 |
| Q163 | Default 2 shops per trip. | OFF-04-02 |
| Q164 | Default travel range 3 km. | OFF-04-02 |
| Q165 | Walk between shops 0–1,500 m in 100 m steps; default 500 m. | OFF-04-02 |
| Q166 | Minimum saving for an extra shop: PLN 8.00 in Poland, €2.00 in euro countries. | OFF-04-02 |
| Q167 | An unwanted store edition is still read for 30 days. | OFF-04-03 |
| Q168 | Inactive my-shops entries are deleted after 90 days. | OFF-04-04 |
| Q169 | At most 300 store editions per chain per day, highest demand first, daily alert above it. | OFF-05-07 |
| Q170 | Offers show the product name as printed plus the ingredient name in the user's language. | OFF-09-01 |
| Q171 | At most 10 own-flyer uploads per household per day. | OFF-10-01 |
| Q172 | Only direct PDF or image links in R2. | OFF-10-02 |
| Q173 | Offers from a user's flyer are private to that household and never used as public prices. | OFF-10-04 |
| Q174 | Reference price: weekly median of 90 days of observations (≥ 3 from ≥ 2 chains); admin CSV overrides. | OFF-11-03 |
| Q175 | Loyalty-only offers count in the price matrix, flagged "with app or card". | OFF-11-04 |
| Q226 | Auchan PL: national flyers only in R3. | OFF-12-04 |
| Q227 | Penny DE: national offers only. | OFF-12-06 |
| Q228 | Shops are created from the feed only where OSM has none, marked `source = 'feed'`. | OFF-13-01 |
