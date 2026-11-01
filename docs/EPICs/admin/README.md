# admin — epics and tasks

New Bun + Elysia back-office service (decision 9, D2): `admin/`, address `10.77.0.23`, HTTP `3000`, admin port
`8081`, own `admin` database (audit log only), published at `admin.<domain>` for accounts identity marks as admins.
It serves an HTTP API with OpenAPI for the admin UI, which comes later with the frontend tasks, and calls each owning
service's `<Svc>AdminService`. Conventions: `docs/EPICs/README.md`, `docs/service-conventions.md`.

## Epics

| Epic | Title | Rel. | Tasks |
|---|---|---|---|
| [ADM-01](E01-service-skeleton/README.md) | Service skeleton (compose, Caddy, database, monitoring) | R1 | 5 |
| [ADM-02](E02-auth-and-audit/README.md) | Admin sign-in, roles and audit log | R1 | 4 |
| [ADM-03](E03-catalogue-review/README.md) | Catalogue, matching and nutrition review | R1 | 7 |
| [ADM-04](E04-recipe-review/README.md) | Recipe import and dish review | R1 (one R2 task) | 8 |
| [ADM-05](E05-households-support/README.md) | Households and plans support | R1 | 3 |
| [ADM-06](E06-flyer-review/README.md) | Flyer sources, offers and shops review | R2 (one R3 task) | 9 |
| [ADM-07](E07-ai-cost-and-overview/README.md) | AI cost and system overview | R1 | 4 |

## Build order

| # | Task | Title | Rel. | Size | Depends on |
|---|---|---|---|---|---|
| 1 | ADM-01-01 | Project scaffold (Bun + Elysia) | R1 | M | — |
| 2 | ADM-01-02 | Compose, Makefile and editor wiring | R1 | M | ADM-01-01 |
| 3 | ADM-01-03 | Database, migrations and secrets | R1 | S | ADM-01-01 |
| 4 | ADM-01-04 | Caddy route and Prometheus scrape | R1 | S | ADM-01-02 |
| 5 | ADM-01-05 | gRPC clients, error model and OpenAPI | R1 | M | ADM-01-01, GW-01 |
| 6 | ADM-02-01 | Sign-in through the admin origin | R1 | M | ADM-01-04, IDN-02 |
| 7 | ADM-02-02 | Admin check on every request | R1 | M | ADM-02-01, IDN-10-01 |
| 8 | ADM-02-03 | Audit log | R1 | M | ADM-01-03, ADM-02-02 |
| 9 | ADM-02-04 | Admin accounts list | R1 | S | ADM-02-03, IDN-10-01 |
| 10 | ADM-05-01 | Household lookup | R1 | M | ADM-02-03, IDN-10-02, IDN-04 |
| 11 | ADM-05-02 | Manual plan switch | R1 | S | ADM-05-01, IDN-10-03 |
| 12 | ADM-05-03 | Delete an account on request | R1 | S | ADM-05-01, IDN-09-02 |
| 13 | ADM-07-04 | Delivery log, alerts and test sends | R1 | S | ADM-02-03, NTF-06 |
| 14 | ADM-03-03 | Ingredient editor and translation review | R1 | M | ADM-02-03, KIT-02-05/06/10 |
| 15 | ADM-03-06 | Catalogue batch jobs | R1 | S | ADM-02-03, KIT-01-08, AI-02-04 |
| 16 | ADM-03-01 | Name-match review queue | R1 | M | ADM-02-03, KIT-03-06 |
| 17 | ADM-03-02 | Nutrition match approval | R1 | M | ADM-02-03, KIT-06-06 |
| 18 | ADM-03-04 | Ingredient merge | R1 | S | ADM-03-03, KIT-02-09 |
| 19 | ADM-03-05 | Unknown products review | R1 | M | ADM-02-03, KIT-05-05 |
| 20 | ADM-03-07 | Weighed-item barcode patterns | R1 | S | ADM-02-03, KIT-04-03 |
| 21 | ADM-04-01 | Recipe import upload | R1 | M | ADM-02-03, PLN-03-03 |
| 22 | ADM-04-02 | AI dish draft batches | R1 | M | ADM-02-03, PLN-03-04, AI-09 |
| 23 | ADM-04-03 | Dish review queue and editor | R1 | L | ADM-04-01 or ADM-04-02, PLN-03-07 |
| 24 | ADM-04-04 | Publish, retire and translate dishes | R1 | M | ADM-04-03, PLN-03-05 |
| 25 | ADM-04-05 | Dietitian sign-off | R1 | S | ADM-04-03, PLN-03-08 |
| 26 | ADM-04-06 | Effort tag suggestions review | R1 | S | ADM-04-03, PLN-03-06 |
| 27 | ADM-04-08 | Library settings: base pantries, no-cook extras, coverage | R1 | M | ADM-04-03, PLN-16-01, PLN-03-09 |
| 28 | ADM-07-03 | Accounts and plans statistics | R1 | S | IDN-10-04 |
| 29 | ADM-07-02 | AI cost and guardrails | R1 | M | ADM-07-03, AI-02-02, AI-02-04 |
| 30 | ADM-07-01 | System overview | R1 | M | KIT-01-07, PLN-01-06, OFF-01-06, NTF-06-02 |
| 31 | ADM-06-01 | Flyer source status board | R2 | M | OFF-05-08, OFF-08-01, OFF-08-04 |
| 32 | ADM-06-02 | Terms gate and chain status | R2 | M | ADM-06-01, OFF-08-02, OFF-08-05 |
| 33 | ADM-06-03 | Accuracy sample of 100 offers | R2 | L | ADM-06-02, OFF-08-01 |
| 34 | ADM-06-04 | Weekly spot-check of 20 offers per chain | R2 | S | ADM-06-03, OFF-08-03 |
| 35 | ADM-06-05 | Held-back offers queue | R2 | M | ADM-06-03, OFF-07-06 |
| 36 | ADM-06-06 | Shop additions and corrections from users | R2 | S | OFF-02-04, OFF-02-07 |
| 37 | ADM-06-07 | Link shops to chain store numbers | R2 | S | OFF-03-05 |
| 38 | ADM-06-08 | Reference prices | R2 | S | OFF-11-03 |
| 39 | ADM-04-07 | Promote well-rated AI dishes to the library | R2 | S | ADM-04-03, PLN-11-02/03 |
| 40 | ADM-06-09 | Licensed feed: retailer and store mapping | R3 | S | ADM-06-07, OFF-13-01 |

40 tasks: 2 L, 20 M, 18 S. R1: 30, R2: 9, R3: 1. ADM-04-08's no-cook extras part is R2 (PLN-12-04).

## Admin RPCs this service calls

Names as the owning service's tasks define them (the task ID is where each RPC is specified). Admin also calls the
public `IdentityService.GetHousehold` (ADM-05-01) and `KitchenService.SearchIngredients` / `GetIngredients` (ADM-03).

| Owner | RPC | Defined in | Used in |
|---|---|---|---|
| identity | `IdentityAdminService.IsAdmin` (with roles), `ListAdmins` | IDN-10-01 | ADM-02-02, ADM-02-04 |
| identity | `FindAccount` | IDN-10-02 | ADM-05-01 |
| identity | `SetPlan` | IDN-10-03 | ADM-05-02 |
| identity | `GetAccountStats` | IDN-10-04 | ADM-07-03, ADM-07-02 |
| identity | `AdminDeleteAccount` | IDN-09-02 (listed in IDN-10-01) | ADM-05-03 |
| notifier | `NotifierAdminService.ListDeliveries`, `GetDeliveryStats` | NTF-06-01 | ADM-07-04 |
| notifier | `ListAlerts`, `AcknowledgeAlert`, `SendTestEmail`, `SendTestPush` | NTF-06-01, NTF-06-02 | ADM-07-01, ADM-07-04 |
| kitchen | `KitchenAdminService.ListAdminChanges`, `GetQueueCounts` | KIT-01-07 | ADM-03-01, ADM-03-03, ADM-07-01 |
| kitchen | `StartCatalogueJob`, `ListCatalogueJobs` | KIT-01-08 | ADM-03-06 |
| kitchen | `MergeIngredients` (with `dry_run`) | KIT-02-09 | ADM-03-04 |
| kitchen | `ListIngredients`, `CreateIngredient`, `UpdateIngredient`, `SetIngredientNames`, `ApproveMachineNames`, `SetMeasures`, `SetPackSizes`, `RetireIngredient` | KIT-02-10 | ADM-03-02, ADM-03-03 |
| kitchen | `ListMatchReviews`, `DecideMatch` | KIT-03-06 | ADM-03-01, ADM-06-05 |
| kitchen | `ListInternalCodePatterns`, `SetInternalCodePattern`, `DecodeInternalCode` | KIT-04-03 | ADM-03-07 |
| kitchen | `ListProductReviews`, `GetSharedProduct`, `ResolveProduct`, `EditSharedProduct` | KIT-05-05 | ADM-03-05 |
| kitchen | `ListNutritionSuggestions`, `ApproveNutritionLink`, `RejectNutritionLink`, `SetNutritionLink`, `SearchNutritionEntries`, `CreateAdminEntry`, `NutritionCoverage` | KIT-06-06 | ADM-03-02 |
| planner | `PlannerAdminService.ListAdminChanges`, `GetQueueCounts` | PLN-01-06 | ADM-04-03, ADM-07-01 |
| planner | `ImportDishes` | PLN-03-03 | ADM-04-01 |
| planner | `DraftDishes`, `GetDraftBatch`, `ListDraftBatches` (all batch kinds) | PLN-03-04 | ADM-04-02, ADM-04-04, ADM-04-06 |
| planner | `TranslateDishes` (with `dry_run`) | PLN-03-05 | ADM-04-04 |
| planner | `SuggestEffortTags` (with `dry_run`), `DecideEffortTag` | PLN-03-06 | ADM-04-06 |
| planner | `ListReviewQueue`, `ListDishes`, `GetReviewItem`, `UpdateDraft`, `ApproveDish`, `RejectDish`, `RetireDish`, `RestoreDish` | PLN-03-07 | ADM-04-03 … 06 |
| planner | `RecordDietitianCheck`, `GetDietitianReport` | PLN-03-08 | ADM-04-05 |
| planner | `GetLibraryCoverage` | PLN-03-09 | ADM-04-08 |
| planner | `PromoteCandidate` / `ListPromotionCandidates`, `DismissCandidate` | PLN-11-02 / PLN-11-03 | ADM-04-07 |
| planner | `ListNoCookExtras`, `UpsertNoCookExtra` | PLN-12-04 | ADM-04-08 |
| planner | `ListBasePantries`, `SetBasePantry` | PLN-16-01 | ADM-04-08 |
| offers | `OffersAdminService.GetQueueCounts` | OFF-01-06 | ADM-07-01 |
| offers | `ListShopMerges`, `MarkNotDuplicate` | OFF-02-04 | ADM-06-06 |
| offers | `ListShopReviews`, `ReviewShop`, `ReviewShopCorrection` | OFF-02-07 | ADM-06-06 |
| offers | `ListStoreLinkQueue`, `LinkStore`, `UnlinkStore` | OFF-03-05 | ADM-06-07 |
| offers | `ListSources`, `UpdateSource`, `ApproveSourceTerms`, `RevokeSourceTerms`, `TriggerSync`, `ListFlyers`, `GetFlyerPage`, `ListEditions` | OFF-05-08 | ADM-06-01 … 03, ADM-06-05 |
| offers | `ListHeldOffers`, `UpdateHeldOffer`, `ReleaseHeldOffer`, `RejectHeldOffer`, `RejectHeldOffers` | OFF-07-06 | ADM-06-05 |
| offers | `CreateReviewSample`, `GetReviewSample`, `RecordReviewVerdict`, `AddMissedOffer`, `ListReviewSamples` | OFF-08-01 | ADM-06-01, ADM-06-03, ADM-06-04 |
| offers | `SetChainLive`, `SetChainSyncing` | OFF-08-02 | ADM-06-02 |
| offers | `GetChainAccuracy` | OFF-08-04 | ADM-06-01, ADM-06-04 |
| offers | `RemoveChain`, `RestoreChain` | OFF-08-05 | ADM-06-02 |
| offers | `ImportReferencePrices`, `ListReferencePrices` | OFF-11-03 | ADM-06-08 |
| offers | `ListFeedRetailers`, `MapFeedRetailer`, `ListFeedStoreQueue`, `LinkFeedStore` | OFF-13-01 | ADM-06-09 |
| ai | `AiAdminService.GetSpendStatus` | AI-02-04 | ADM-03-06, ADM-07-02 |

Owner admin RPCs no admin task uses yet: kitchen `ListDeadLetters`, `ReplayDeadLetter` (KIT-01-05); identity
`RepublishHouseholdSnapshots`, `GetSnapshotJob` (IDN-01-04); offers `ListJobs`, `TriggerJob` (OFF-01-06). They can be
called with `make grpc` until a page needs them.

## Changes needed in other services

- identity: `TRUSTED_ORIGINS` includes `https://admin.<domain>` and `http://10.77.0.23:3000` (ADM-02-01). Admin
  roles `admin` and `dietitian` and `ListAdmins` are in IDN-10-01; granting stays in identity's CLI.
- notifier: metadata key `x-admin-email` on `SendTestEmail` (NTF-06-02).
- `docs/service-conventions.md`: HTTP metrics are no longer gateway-only (admin serves HTTP too).

## Creating the database on an existing volume

`platform/postgres/initdb` runs only on an empty volume. Once, on a machine that already has data
(`make psql s=postgres`):

```sql
CREATE ROLE admin LOGIN PASSWORD 'admin-dev';   -- prod: the content of secrets/admin_db_password
CREATE DATABASE admin OWNER admin;
REVOKE ALL ON DATABASE admin FROM PUBLIC;
```

## Decisions

Settled with the founder on 1–2 November 2026 (full list: [DECISIONS.md](../DECISIONS.md)). Each task states its rule as `Decided (Qnnn)`; decisions owned by another service are listed with that service's task.

| Q | Decision | Tasks |
|---|---|---|
| Q010 | Flyer pages kept 90 days — User photos are still never stored; chain flyer pages are kept until their offers expire plus 90 days. |  (owner) |
| Q028 | Two roles: admin (everything) and dietitian. | ADM-02-02 |
| Q046 | Proxy through the admin service. | ADM-02-01 |
| Q047 | 2 years. | ADM-02-03 |
| Q048 | No; only identity's command line grants them. | ADM-02-04 |
| Q049 | €1.64 per month, as a config value. | ADM-07-02 |
| Q054 | Yes: new event catalogue.ingredient_flags_changed. | ADM-03-03 |
| Q086 | Not for the 200 starter dishes; later yes, if an admin confirms. | PLN-02-02 (owner) |
| Q093 | No: publish, but launch waits until every starter dish is checked. | PLN-03-08 (owner) |
| Q094 | 95%, with every starter dish checked. | PLN-03-08 (owner) |
| Q122 | Yes, after a native speaker spot-checks 10% of dishes per language. | PLN-03-05 (owner) |
| Q158 | Alert only — Bad weekly spot-checks only alert the founder; a chain is never taken offline automatically. | OFF-08-03 (owner) |
| Q159 | No; review every correction by hand until volumes are known. | ADM-06-06 |
| Q174 | Weekly median of 90 days of observed prices (≥ 3 from ≥ 2 chains); admin CSV overrides. | OFF-11-03 (owner) |
| Q228 | Only where OSM has none, marked as coming from the feed. | OFF-13-01 (owner) |

## Gates

- ADM-06-02: a chain's terms decision records a legal check done outside the system; no chain samples or goes live
  without `approved`.
