# gateway — epics and tasks

Bun + Elysia. The public HTTP API the app talks to: verifies JWTs, turns HTTP into gRPC calls on the services that
own the data, enforces the Free-tier AI quota, Plus-only routes and rate limits, makes mutations safe to retry, and
relays live shopping-list changes over a WebSocket. It owns no database; its state is a few `redis-state` keys
(`gw:*`) and short caches in `redis-cache`.

What exists today (`gateway/src`): JWT check against identity's JWKS, gRPC clients with deadlines and metrics,
`/v1/system/status`, `/v1/products/{ean}`, `/v1/me`, `/v1/household`, a catch-all `/api/auth/*` proxy, OpenAPI
export (`make openapi`). These tasks extend that code.

Shared rules: [../README.md](../README.md) (decisions, Free and Plus, call graph, events),
[../TEMPLATE.md](../TEMPLATE.md).

## Epics

| Epic | Title | Rel. | Tasks |
|---|---|---|---|
| [GW-01](E01-platform/README.md) | Platform: request context, error model, idempotency keys, OpenAPI rules, uploads | R1 | 5 |
| [GW-02](E02-auth-and-account/README.md) | Auth and account routes | R1 | 6 |
| [GW-03](E03-free-tier-quota/README.md) | Free-tier AI quota and rate limits | R1 | 5 |
| [GW-04](E04-household-and-invites/README.md) | Household, profiles and invite routes | R1–R2 | 5 |
| [GW-05](E05-kitchen-profile/README.md) | Kitchen profile and catalogue routes | R1–R2 | 4 |
| [GW-06](E06-fridge-and-intake/README.md) | Fridge, barcode and intake routes | R1–R3 | 7 |
| [GW-07](E07-recipes-and-plans/README.md) | Recipe library, plans and cooking routes | R1–R3 | 12 |
| [GW-08](E08-shopping-list/README.md) | Shopping list routes | R1 | 4 |
| [GW-09](E09-live-list/README.md) | Live shopping list (WebSocket) | R2 | 2 |
| [GW-10](E10-offers-and-shops/README.md) | Shops, offers, flyers and price matrix routes | R2–R3 | 6 |
| [GW-11](E11-devices-and-notifications/README.md) | Devices and notification preference routes | R1 | 2 |
| [GW-12](E12-savings/README.md) | Savings routes | R2–R3 | 2 |
| | **Total** | | **60** |

## Build order

A route task can start once its backend RPC exists (or against a fake gRPC server in tests); the order below follows
the gateway's own dependencies.

**R1 — launch (38 tasks)**

| # | Task | Size |
|---|---|---|
| 1 | [GW-01-01](E01-platform/GW-01-01-request-context.md) Request context, metadata forwarding and deadline classes | S |
| 2 | [GW-01-02](E01-platform/GW-01-02-error-model.md) Problem JSON error model and gRPC → HTTP mapping | M |
| 3 | [GW-01-04](E01-platform/GW-01-04-openapi-rules.md) OpenAPI rules, tags and export check | S |
| 4 | [GW-01-03](E01-platform/GW-01-03-idempotency-keys.md) Idempotency keys for mutations | M |
| 5 | [GW-01-05](E01-platform/GW-01-05-uploads.md) Photo and PDF uploads | M |
| 6 | [GW-02-01](E02-auth-and-account/GW-02-01-jwt-claims.md) JWT claims: member id, anonymous flag, registered-only routes | S |
| 7 | [GW-02-02](E02-auth-and-account/GW-02-02-auth-proxy-allowlist.md) Better Auth proxy allowlist | M |
| 8 | [GW-02-03](E02-auth-and-account/GW-02-03-me.md) `GET /v1/me` | S |
| 9 | [GW-03-04](E03-free-tier-quota/GW-03-04-rate-limits.md) Rate limits and Plus fair use | M |
| 10 | [GW-03-01](E03-free-tier-quota/GW-03-01-quota-counters.md) Quota counters with reserve / commit / release | M |
| 11 | [GW-03-02](E03-free-tier-quota/GW-03-02-quota-responses.md) Quota in responses and `quota_reached` | S |
| 12 | [GW-03-03](E03-free-tier-quota/GW-03-03-plus-guard.md) Plus-only route guard | S |
| 13 | [GW-03-05](E03-free-tier-quota/GW-03-05-ai-metrics.md) AI usage metrics and dashboard panel | S |
| 14 | [GW-02-04](E02-auth-and-account/GW-02-04-entitlements.md) `GET /v1/entitlements` | S |
| 15 | [GW-04-01](E04-household-and-invites/GW-04-01-household-read-rename.md) `GET` and `PATCH /v1/household` | S |
| 16 | [GW-05-02](E05-kitchen-profile/GW-05-02-profile-options.md) `GET /v1/profile/options` | S |
| 17 | [GW-05-01](E05-kitchen-profile/GW-05-01-household-profile.md) `GET` and `PATCH /v1/profile` | M |
| 18 | [GW-05-03](E05-kitchen-profile/GW-05-03-ingredient-search.md) Ingredient autocomplete | S |
| 19 | [GW-06-01](E06-fridge-and-intake/GW-06-01-stock-routes.md) Stock: list, add, edit, remove | M |
| 20 | [GW-06-02](E06-fridge-and-intake/GW-06-02-barcode-resolve.md) Barcode and 2D code lookup | S |
| 21 | [GW-06-03](E06-fridge-and-intake/GW-06-03-name-unknown-product.md) Name an unknown product, confirm a shared one | S |
| 22 | [GW-06-04](E06-fridge-and-intake/GW-06-04-photo-intake.md) Fridge or shelf photo → proposals | M |
| 23 | [GW-06-05](E06-fridge-and-intake/GW-06-05-label-photo.md) Label photo for an unknown product | S |
| 24 | [GW-07-01](E07-recipes-and-plans/GW-07-01-library-browse.md) Library browse and dish detail | M |
| 25 | [GW-07-02](E07-recipes-and-plans/GW-07-02-favourites-ratings-history.md) Favourites, ratings, history | S |
| 26 | [GW-07-03](E07-recipes-and-plans/GW-07-03-suggestions.md) What can I make now | S |
| 27 | [GW-07-04](E07-recipes-and-plans/GW-07-04-ai-day-plan.md) AI day plan and cuisine choices | M |
| 28 | [GW-07-05](E07-recipes-and-plans/GW-07-05-plans-read-edit.md) Read, accept, edit and delete plans | S |
| 29 | [GW-07-06](E07-recipes-and-plans/GW-07-06-swaps-wishes.md) Swap a meal, replace an ingredient, wish | M |
| 30 | [GW-07-07](E07-recipes-and-plans/GW-07-07-cooking-and-cooked.md) Cooking mode and "cooked it" | M |
| 31 | [GW-08-01](E08-shopping-list/GW-08-01-list-lines.md) Read the list; add, edit, delete lines | M |
| 32 | [GW-08-02](E08-shopping-list/GW-08-02-proposals.md) Propose lines from the plan and accept them | M |
| 33 | [GW-08-03](E08-shopping-list/GW-08-03-tick.md) Tick in the shop → into the fridge | S |
| 34 | [GW-08-04](E08-shopping-list/GW-08-04-base-pantry.md) Cuisine base pantry | S |
| 35 | [GW-11-01](E11-devices-and-notifications/GW-11-01-devices.md) Register and unregister a device | S |
| 36 | [GW-11-02](E11-devices-and-notifications/GW-11-02-notification-preferences.md) Notification preferences and quiet hours | S |
| 37 | [GW-02-05](E02-auth-and-account/GW-02-05-delete-account.md) Account deletion with an e-mail code | M |
| 38 | [GW-02-06](E02-auth-and-account/GW-02-06-data-export.md) `POST /v1/me/export` | M |

**R2 — flyers and family (16 tasks)**

| # | Task | Size |
|---|---|---|
| 39 | [GW-04-02](E04-household-and-invites/GW-04-02-profiles-without-accounts.md) Profiles without accounts | S |
| 40 | [GW-04-03](E04-household-and-invites/GW-04-03-invites-create-list-cancel.md) Invites: create, list, cancel | M |
| 41 | [GW-04-04](E04-household-and-invites/GW-04-04-invite-preview-accept.md) Invite preview and accept | M |
| 42 | [GW-04-05](E04-household-and-invites/GW-04-05-leave-remove.md) Leave the household, remove a member | S |
| 43 | [GW-05-04](E05-kitchen-profile/GW-05-04-member-profiles.md) Per-person profiles and daily kcal targets | M |
| 44 | [GW-07-08](E07-recipes-and-plans/GW-07-08-portions-attendance.md) Per-person portions and "away" marking | S |
| 45 | [GW-09-01](E09-live-list/GW-09-01-socket-auth.md) `/v1/live` socket | M |
| 46 | [GW-09-02](E09-live-list/GW-09-02-pubsub-relay.md) Pub/sub relay and versioned change notices | M |
| 47 | [GW-10-01](E10-offers-and-shops/GW-10-01-shops-by-postcode.md) Shops around a postcode; suggest a change | S |
| 48 | [GW-10-02](E10-offers-and-shops/GW-10-02-my-shops.md) My shops and shopping preferences | S |
| 49 | [GW-10-03](E10-offers-and-shops/GW-10-03-offers-browse.md) Offer browse, search and flyer status | M |
| 50 | [GW-10-04](E10-offers-and-shops/GW-10-04-own-flyer.md) Scan a flyer that is not covered | M |
| 51 | [GW-10-05](E10-offers-and-shops/GW-10-05-price-matrix.md) Trip price matrix | M |
| 52 | [GW-07-09](E07-recipes-and-plans/GW-07-09-tool-rewrite.md) Rewrite a dish for the tools at hand | S |
| 53 | [GW-07-10](E07-recipes-and-plans/GW-07-10-plan-from-deals.md) Plan from deals | M |
| 54 | [GW-12-01](E12-savings/GW-12-01-savings-counters-report.md) Savings counters and monthly reports | S |

**R3 — depth (6 tasks)**

| # | Task | Size |
|---|---|---|
| 55 | [GW-06-06](E06-fridge-and-intake/GW-06-06-text-intake.md) Typed or spoken sentence → proposals | S |
| 56 | [GW-06-07](E06-fridge-and-intake/GW-06-07-receipt-intake.md) Receipt photo → proposals | S |
| 57 | [GW-07-11](E07-recipes-and-plans/GW-07-11-week-plan.md) Week plan | M |
| 58 | [GW-07-12](E07-recipes-and-plans/GW-07-12-restyle.md) Restyle a dish | S |
| 59 | [GW-10-06](E10-offers-and-shops/GW-10-06-cheapest-top-up.md) Cheapest top-up | S |
| 60 | [GW-12-02](E12-savings/GW-12-02-savings-card.md) Shareable savings card data | S |

## Route catalogue

The contract the mobile communication layer (MOB-02, MOB-03) builds on. `auth`: `—` public, `any` any account
(anonymous included), `reg` registered only. `Idem`: idempotency key `req`uired / `opt`ional (GW-01-03). `Quota`:
Free weekly kind (GW-03-01). `Plus`: Plus-only feature name (GW-03-03). Deadline classes: GW-01-01. Every backend
RPC name below is the one its owner's task defines.

| Method | Path | Auth | Idem | Quota | Plus | Backend RPC | Task |
|---|---|---|---|---|---|---|---|
| GET | `/v1/system/status` | — | | | | health checks | exists |
| GET | `/v1/me` | any | | | | (JWT only) | GW-02-03 |
| GET | `/v1/entitlements` | any | | | | identity `GetEntitlements` | GW-02-04 |
| POST | `/v1/me/deletion-code` | reg | | | | identity `RequestAccountDeletion` | GW-02-05 |
| DELETE | `/v1/me` | any | req | | | identity `DeleteAccount` | GW-02-05 |
| POST | `/v1/me/export` | reg | | | | identity, notifier `ExportUserData`; kitchen, planner, offers `ExportHouseholdData` | GW-02-06 |
| * | `/api/auth/*` (allowlist) | — | | | | Better Auth HTTP | GW-02-02 |
| GET | `/v1/household` | any | | | | identity `GetHousehold` | GW-04-01 |
| PATCH | `/v1/household` | any | opt | | | identity `UpdateHousehold` | GW-04-01 |
| POST | `/v1/household/profiles` | reg | req | | `profiles` | identity `CreateProfile` | GW-04-02 |
| PATCH, DELETE | `/v1/household/profiles/{memberId}` | reg | opt | | | identity `UpdateMember`, `DeleteProfile` | GW-04-02 |
| POST | `/v1/household/invites` | reg | req | | `shared_household` | identity `CreateInvite` | GW-04-03 |
| GET | `/v1/household/invites` | reg | | | | identity `ListInvites` | GW-04-03 |
| DELETE | `/v1/household/invites/{inviteId}` | reg | opt | | | identity `CancelInvite` | GW-04-03 |
| GET | `/v1/invites/{code}` | reg | | | | identity `GetInvitePreview` | GW-04-04 |
| POST | `/v1/invites/{code}/accept` | reg | req | | | identity `AcceptInvite` | GW-04-04 |
| POST | `/v1/household/leave` | reg | req | | | identity `LeaveHousehold` | GW-04-05 |
| POST | `/v1/household/owner` | reg | opt | | | identity `TransferOwnership` | GW-04-05 |
| DELETE | `/v1/household/members/{memberId}` | reg | opt | | | identity `RemoveMember` | GW-04-05 |
| GET, PATCH | `/v1/profile` | any | opt | | | kitchen `GetProfile`, `UpdateProfile` | GW-05-01 |
| GET | `/v1/profile/options` | — | | | | kitchen `GetProfileOptions` | GW-05-02 |
| GET | `/v1/catalogue/nutrition-sources` | — | | | | kitchen `ListNutritionSources` | GW-05-02 |
| GET | `/v1/profile/members` | any | | | | kitchen `GetMemberProfiles` | GW-05-04 |
| PUT | `/v1/profile/members/{memberId}` | any | opt | | | kitchen `UpdateMemberProfile` | GW-05-04 |
| GET | `/v1/catalogue/ingredients`, `/{id}` | any | | | | kitchen `SearchIngredients`, `GetIngredients` | GW-05-03 |
| GET | `/v1/stock` | any | | | | kitchen `ListStock` | GW-06-01 |
| POST | `/v1/stock` | any | req | | | kitchen `AddStockItems` | GW-06-01 |
| PATCH, DELETE | `/v1/stock/{itemId}` (`If-Match`) | any | req | | | kitchen `UpdateStockItem`, `RemoveStockItems` | GW-06-01 |
| POST | `/v1/stock/preview` | any | | | | kitchen `PreviewStockItems` | GW-06-01 |
| GET | `/v1/stock/defaults` | any | | | | kitchen `SuggestStockDefaults` | GW-06-01 |
| GET | `/v1/products/resolve?code=` | any | | | | kitchen `ResolveBarcode` | GW-06-02 |
| GET | `/v1/products/{ean}` (deprecated alias) | any | | | | kitchen `ResolveBarcode` | GW-06-02 |
| GET | `/v1/products/mine` | any | | | | kitchen `ListHouseholdProducts` | GW-06-02 |
| POST | `/v1/products` | any | req | | | kitchen `NameProduct` | GW-06-03 |
| POST | `/v1/products/{ean}/confirm` | any | req | | | kitchen `ConfirmProposal` | GW-06-03 |
| POST | `/v1/products/{ean}/report` | any | req | | | kitchen `ReportProduct` | GW-06-03 |
| POST | `/v1/intake/photo` (multipart) | any | req | scan | | kitchen `IntakeFromPhoto` | GW-06-04 |
| POST | `/v1/intake/label` (multipart) | any | req | scan | | kitchen `ReadLabel` | GW-06-05 |
| POST | `/v1/intake/text` | any | req | scan | | kitchen `IntakeFromText` | GW-06-06 |
| POST | `/v1/intake/receipt` (multipart) | any | req | scan | | kitchen `IntakeFromReceipt` | GW-06-07 |
| GET | `/v1/dishes`, `/v1/dishes/{dishId}` | any | | | | planner `SearchDishes`, `GetDish` | GW-07-01 |
| PUT, DELETE | `/v1/dishes/{dishId}/favourite` | any | opt | | | planner `SetFavourite` | GW-07-02 |
| PUT, DELETE | `/v1/dishes/{dishId}/rating` | any | opt | | | planner `RateDish` | GW-07-02 |
| GET | `/v1/history` | any | | | | planner `ListCookedHistory` | GW-07-02 |
| GET | `/v1/suggestions` | any | | | | planner `SuggestDishes` | GW-07-03 |
| POST | `/v1/plans/day` | any | req | dayPlan | | planner `ComposeDayPlan`; Plus from R2 `CreateFamilyDayPlan` | GW-07-04 |
| POST | `/v1/plans/{planId}/meals/{meal}/choice` | any | req | | | planner `ChooseMealOption` | GW-07-04 |
| GET | `/v1/plans?from=&to=` | any | | | | planner `GetPlans` | GW-07-05 |
| POST | `/v1/plans/{planId}/accept` | any | req | | | planner `AcceptPlan` | GW-07-05 |
| PUT | `/v1/plans/{planId}/meals/{meal}` | any | opt | | | planner `SetPlanMeal` | GW-07-05 |
| DELETE | `/v1/plans/{planId}` | any | opt | | | planner `DeletePlan` | GW-07-05 |
| POST | `/v1/plans/{planId}/meals/{meal}/swap` | any | req | | | planner `SwapMeal` | GW-07-06 |
| POST | `/v1/plans/{planId}/meals/{meal}/replace-ingredient` | any | req | | `swap_meal` | planner `ReplaceIngredient` | GW-07-06 |
| POST | `/v1/plans/wish` | any | req | | `swap_meal` | planner `WishMeal` | GW-07-06 |
| GET | `/v1/cooking/{dishId}` | any | | | | planner `GetCookingSteps` | GW-07-07 |
| POST | `/v1/cooked` | any | req | | | planner `MarkCooked` | GW-07-07 |
| POST | `/v1/dishes/{dishId}/rating-prompt/dismiss` | any | opt | | | planner `DismissRatingPrompt` | GW-07-07 |
| GET | `/v1/plan-meals/{planMealId}/portions` | any | | | | planner `GetMealPortions` | GW-07-08 |
| GET | `/v1/plan-meals/{planMealId}/ingredients` | any | | | | planner `GetScaledIngredients` | GW-07-08 |
| PUT | `/v1/attendance` | any | opt | | | planner `SetMealAttendance` | GW-07-08 |
| GET | `/v1/attendance` | any | | | | planner `ListMealAttendance` | GW-07-08 |
| POST, PUT, DELETE | `/v1/attendance/rules`, `/{ruleId}` | any | req / opt | | | planner `SetAttendanceRule`, `DeleteAttendanceRule` | GW-07-08 |
| POST | `/v1/dishes/{dishId}/rewrite` | any | req | | `tool_rewrite` | planner `RewriteDish` (TOOLS) | GW-07-09 |
| POST | `/v1/plans/from-deals` | any | req | | `plan_from_deals` | planner `CreatePlanFromDeals` | GW-07-10 |
| POST, GET | `/v1/plans/week` | any | req / — | | `week_plan` | planner `CreateWeekPlan`, `GetWeekPlan` | GW-07-11 |
| PUT | `/v1/plans/week/{weekPlanId}/meals` | any | opt | | `week_plan` | planner `ReplaceWeekMeal` | GW-07-11 |
| DELETE | `/v1/plans/week/{weekPlanId}` | any | opt | | | planner `DeleteWeekPlan` | GW-07-11 |
| POST | `/v1/dishes/{dishId}/restyle` | any | req | | `restyle` | planner `RewriteDish` (RESTYLE) | GW-07-12 |
| GET | `/v1/list` | any | | | | planner `GetList` | GW-08-01 |
| POST | `/v1/list/lines` | any | req | | | planner `AddListLine` | GW-08-01 |
| PATCH, DELETE | `/v1/list/lines/{lineId}` | any | req | | | planner `UpdateListLine`, `RemoveListLine` | GW-08-01 |
| DELETE | `/v1/list/lines?state=ticked` | any | req | | | planner `ClearTickedLines` | GW-08-01 |
| GET | `/v1/list/changes?since=` | any | | | | planner `GetListChanges` | GW-08-01 |
| POST | `/v1/list/proposals` | any | | | | planner `GenerateList` | GW-08-02 |
| POST | `/v1/list/proposals/{proposalId}/accept` | any | req | | | planner `AcceptListProposal` | GW-08-02 |
| POST | `/v1/list/lines/{lineId}/tick`, `/untick` | any | req | | | planner `TickListLine`, `UntickListLine` | GW-08-03 |
| GET | `/v1/cuisines/{cuisine}/base-pantry` | any | | | | planner `GetBasePantry` | GW-08-04 |
| POST | `/v1/list/base-pantry` | any | req | | | planner `AddBasePantryToList` | GW-08-04 |
| GET (WS) | `/v1/live` | token in first message | | | | Redis pub/sub `live.household.<id>` | GW-09 |
| GET | `/v1/areas?country=&q=` | any | | | | offers `ResolveArea` | GW-10-01 |
| GET | `/v1/shops` | any | | | | offers `ListShopsNearArea` | GW-10-01 |
| POST | `/v1/shops`, `/v1/shops/{shopId}/corrections` | any | opt | | | offers `SubmitShop`, `ProposeShopCorrection` | GW-10-01 |
| GET, PUT | `/v1/my-shops` | any | opt | | | offers `ListMyShops`, `SetMyShops` | GW-10-02 |
| GET, PUT | `/v1/shopping-preferences` | any | opt | | | offers `GetShoppingPreferences`, `SetShoppingPreferences` | GW-10-02 |
| GET | `/v1/offers`, `/v1/offers/{offerId}`, `/v1/offers/categories` | any | | | | offers `BrowseOffers`, `GetOffer`, `ListOfferCategories` | GW-10-03 |
| POST | `/v1/list/offer-lines` | any | req | | | planner `AddOfferToList` | GW-10-03 |
| GET | `/v1/flyers/status?shopIds=` | any | | | | offers `GetFlyerStatus` | GW-10-03 |
| POST | `/v1/flyers` (multipart or JSON) | any | req | | `own_flyer` | offers `UploadUserFlyer`, `SubmitUserFlyerLink` | GW-10-04 |
| GET, DELETE | `/v1/flyers/{flyerId}` | any | — / opt | | `own_flyer` | offers `GetUserFlyer`, `DeleteUserFlyer` | GW-10-04 |
| POST | `/v1/flyers/{flyerId}/confirm` | any | req | | `own_flyer` | offers `ConfirmUserFlyer` | GW-10-04 |
| POST | `/v1/list/price-matrix` | any | | | `trip_plan` | planner `GetTripPriceMatrix` | GW-10-05 |
| POST | `/v1/list/cheapest-top-up` | any | req | | `cheapest_top_up` | planner `CreateCheapestTopUp` | GW-10-06 |
| POST | `/v1/list/cheapest-top-up/{topUpId}/apply` | any | req | | `cheapest_top_up` | planner `ApplyTopUp` | GW-10-06 |
| PUT, DELETE | `/v1/devices/{installationId}` | any | opt | | | notifier `RegisterDevice`, `UnregisterDevice` | GW-11-01 |
| GET, PUT | `/v1/notification-preferences` | any | opt | | | notifier `GetPreferences`, `SetPreferences` | GW-11-02 |
| GET | `/v1/savings` | any | | | | kitchen `GetSavings` | GW-12-01 |
| GET | `/v1/savings/reports` | any | | | `savings_report` | kitchen `ListSavingsReports` | GW-12-01 |
| GET | `/v1/savings/card?month=` | any | | | `savings_report` | kitchen `GetSavingsCard` | GW-12-02 |

Headers on every response: `X-Request-Id`, `Content-Language`. On quota routes: `X-Quota-Limit`,
`X-Quota-Remaining`, `X-Quota-Next-Free-At`. On replays: `Idempotent-Replayed: true`. Errors: `application/problem+json`
with a stable `code` (GW-01-02). WebSocket: control messages carry `type` (`ready`, `pong`, `degraded`); change
notices carry `kind` (`list.changed`, `stock.changed`) in the README shape (GW-09-02). Close codes: 4401 auth, 4403
household changed, 4429 limits, 1001 shutdown.

## What the gateway needs from other services

- **Error reasons** in trailing gRPC metadata `x-error-code` and `ABORTED` for version conflicts (README decision 18).
- **`ai_used`** on every AI-capable response: kitchen intake RPCs and `ReadLabel`; planner `ComposeDayPlan`
  (`quota_consumed` for the Free count), `SwapMeal`, `ReplaceIngredient`, `WishMeal`, `RewriteDish`,
  `CreateFamilyDayPlan`, `CreateWeekPlan`, `CreatePlanFromDeals`, `CreateCheapestTopUp`. All are in the owners'
  contract sketches.
- **JWT claims from identity:** `memberId` and `isAnonymous` (GW-02-01).
- **Revocation key** `auth:revoked:<user_id>:<household_id>` written by identity (IDN-08-04), read on every request.
- **Message size:** kitchen and offers gRPC servers accept 26 MB messages for uploads (GW-01-05; offers allows 32 MB
  for user flyers).
- **Proxy trust:** identity trusts `X-Forwarded-For` only from the gateway (GW-02-02).
- **Live notices** in the README shape `{kind, household_id, version}`, published by planner (PLN-15-09) and
  kitchen (KIT-07-08).
- **planner `SetPlanMeal`** (PLN-08-01) for manual placement (Q117); **notifier `UnregisterDevice` by
  `installation_id`** (NTF-04-01, Q130); **identity `owner_must_transfer`** on deletion (IDN-09, Q137).

## Decisions

Settled on 1 November 2026 (full list: [../DECISIONS.md](../DECISIONS.md)).

| ID | Decision | Tasks |
|---|---|---|
| Q020 | Parse `x-client`; versions below `MIN_APP_VERSION` get 426 `app_update_required` | [GW-01-01](E01-platform/GW-01-01-request-context.md), [GW-01-02](E01-platform/GW-01-02-error-model.md) |
| Q021 | Idempotent responses are replayable for 24 hours | [GW-01-03](E01-platform/GW-01-03-idempotency-keys.md) |
| Q022 | Own-flyer PDFs up to 40 pages and 25 MB | [GW-01-05](E01-platform/GW-01-05-uploads.md), [GW-10-04](E10-offers-and-shops/GW-10-04-own-flyer.md) |
| Q025 | Free limits come from identity's entitlements, never older than 5 minutes | [GW-03-01](E03-free-tier-quota/GW-03-01-quota-counters.md) |
| Q034 | 5 anonymous sign-ups per IP per hour | [GW-03-04](E03-free-tier-quota/GW-03-04-rate-limits.md) |
| Q035 | Only registered accounts can join a household | [GW-04-04](E04-household-and-invites/GW-04-04-invite-preview-accept.md) |
| Q076 | `household.contributeToOpenFoodFacts` in the profile, default `true` | [GW-05-01](E05-kitchen-profile/GW-05-01-household-profile.md) |
| Q111 | Free quota is a rolling 7 days; responses carry `nextFreeAt` and `X-Quota-Next-Free-At` | [GW-03-01](E03-free-tier-quota/GW-03-01-quota-counters.md), [GW-03-02](E03-free-tier-quota/GW-03-02-quota-responses.md) |
| Q112 | Every AI reading request counts as one Free scan (label photos and typed text included) | [GW-03-01](E03-free-tier-quota/GW-03-01-quota-counters.md), GW-06-05, GW-06-06 |
| Q113 | Plus fair use: 60 AI responses per household per UTC day, stated in the terms | [GW-03-04](E03-free-tier-quota/GW-03-04-rate-limits.md) |
| Q117 | `PUT /v1/plans/{planId}/meals/{meal}` puts a library dish into a day plan via planner `SetPlanMeal` | [GW-07-05](E07-recipes-and-plans/GW-07-05-plans-read-edit.md) |
| Q125 | From R2 every Plus household uses the family day plan | [GW-07-04](E07-recipes-and-plans/GW-07-04-ai-day-plan.md) |
| Q130 | Devices are unregistered by installation id; `DELETE /v1/devices/{installationId}` has no body | [GW-11-01](E11-devices-and-notifications/GW-11-01-devices.md) |
| Q137, Q186 | An owner names the new owner before leaving or deleting the account (`owner_must_transfer`) | [GW-02-05](E02-auth-and-account/GW-02-05-delete-account.md), [GW-04-05](E04-household-and-invites/GW-04-05-leave-remove.md) |
| Q138, Q142, Q143 | Export: JSON per service nested by the gateway, a direct download, none for anonymous accounts | [GW-02-06](E02-auth-and-account/GW-02-06-data-export.md) |
| Q162 | The price matrix compares up to 10 shops | [GW-10-05](E10-offers-and-shops/GW-10-05-price-matrix.md) |
| Q173 | Offers from a user's flyer stay private to that household | [GW-10-04](E10-offers-and-shops/GW-10-04-own-flyer.md) |
| Q181 | Invite links are `https://coolinary.app/i/<code>` | [GW-04-03](E04-household-and-invites/GW-04-03-invites-create-list-cancel.md) |
| Q187 | After Plus → Free a shared household keeps its members; invites and joining are blocked | GW-04-03, [GW-04-04](E04-household-and-invites/GW-04-04-invite-preview-accept.md) |
| Q191 | A removed member loses access at once (revocation key checked on every request) | [GW-04-05](E04-household-and-invites/GW-04-05-leave-remove.md) |
| Q212 | Free sees the current month's savings counters; report and history are Plus | [GW-12-01](E12-savings/GW-12-01-savings-counters-report.md) |
| Q219 | Typed-sentence intake ships in R3 with voice | [GW-06-06](E06-fridge-and-intake/GW-06-06-text-intake.md) |
| Q225 | The savings card image is drawn on the phone | [GW-12-02](E12-savings/GW-12-02-savings-card.md) |
