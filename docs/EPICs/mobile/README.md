# mobile — communication layer

Scope: how the Expo app talks to the back end: auth session, typed API client and errors, query hooks and cache
rules, offline queue, live list socket, push registration, uploads. Screens and on-device calculators (calorie target
from age/height/weight, distances to shops, the trip-plan search) are later frontend tasks; the hooks here leave
typed inputs for them (`dailyKcalTarget` only, `shopIds` nearest first, the price matrix).

Builds on the existing app: `src/api/client.ts` (openapi-fetch + openapi-react-query), `src/lib/query.ts` (TanStack
Query with AsyncStorage persistence), `src/lib/secure.ts` (SecureStore for health data), `src/stores/settings.ts`.

## Epics

| Epic | Title | Rel. | Tasks |
|---|---|---|---|
| [MOB-01](E01-auth-session/README.md) | Auth session (anonymous start, e-mail code, JWT refresh) | R1 | 8 |
| [MOB-02](E02-api-client/README.md) | API client and error model | R1 | 5 |
| [MOB-03](E03-query-hooks/README.md) | Query keys, hooks and cache rules | R1–R3 | 9 |
| [MOB-04](E04-offline-queue/README.md) | Offline mutation queue | R1 | 4 |
| [MOB-05](E05-live-list/README.md) | Live shopping list client | R2 | 3 |
| [MOB-06](E06-push/README.md) | Push notifications | R1 | 4 |
| [MOB-07](E07-uploads/README.md) | Uploads (photos, PDFs) | R1–R3 | 3 |

## Build order

| # | Task | Title | Rel. | Size |
|---|---|---|---|---|
| 1 | [MOB-01-01](E01-auth-session/MOB-01-01-auth-client-setup.md) | Better Auth client setup | R1 | S |
| 2 | [MOB-02-04](E02-api-client/MOB-02-04-connectivity.md) | Connectivity and server reachability | R1 | S |
| 3 | [MOB-01-02](E01-auth-session/MOB-01-02-auth-state.md) | Auth state and app start | R1 | M |
| 4 | [MOB-01-04](E01-auth-session/MOB-01-04-jwt-provider.md) | JWT provider and refresh | R1 | M |
| 5 | [MOB-02-01](E02-api-client/MOB-02-01-request-middleware.md) | Request middleware | R1 | M |
| 6 | [MOB-02-02](E02-api-client/MOB-02-02-error-model.md) | Typed error model and retry policy | R1 | M |
| 7 | [MOB-02-03](E02-api-client/MOB-02-03-openapi-workflow.md) | OpenAPI types workflow | R1 | S |
| 8 | [MOB-01-03](E01-auth-session/MOB-01-03-anonymous-start.md) | Anonymous start | R1 | S |
| 9 | [MOB-01-06](E01-auth-session/MOB-01-06-session-loss.md) | 401 recovery and session loss | R1 | S |
| 10 | [MOB-03-01](E03-query-hooks/MOB-03-01-query-keys-cache-rules.md) | Query keys and cache rules | R1 | M |
| 11 | [MOB-03-02](E03-query-hooks/MOB-03-02-invalidation-map.md) | Invalidation map | R1 | S |
| 12 | [MOB-04-01](E04-offline-queue/MOB-04-01-outbox-storage.md) | Outbox: storage and allowed changes | R1 | M |
| 13 | [MOB-04-04](E04-offline-queue/MOB-04-04-pending-overlay.md) | Pending changes overlay | R1 | M |
| 14 | [MOB-04-02](E04-offline-queue/MOB-04-02-outbox-replay.md) | Outbox replay | R1 | M |
| 15 | [MOB-04-03](E04-offline-queue/MOB-04-03-conflict-rules.md) | Conflict rules | R1 | M |
| 16 | [MOB-06-01](E06-push/MOB-06-01-dev-build-notification-config.md) | Development build and notification setup | R1 | M |
| 17 | [MOB-06-02](E06-push/MOB-06-02-permission-timing.md) | When to ask for notification permission | R1 | S |
| 18 | [MOB-06-03](E06-push/MOB-06-03-device-registration.md) | Device and push token registration | R1 | M |
| 19 | [MOB-01-07](E01-auth-session/MOB-01-07-sign-out.md) | Sign-out and local wipe | R1 | M |
| 20 | [MOB-01-05](E01-auth-session/MOB-01-05-email-code-sign-in.md) | Sign-in with an e-mail code and anonymous upgrade | R1 | M |
| 21 | [MOB-01-08](E01-auth-session/MOB-01-08-account-deletion-export.md) | Account deletion and data export calls | R1 | S |
| 22 | [MOB-02-05](E02-api-client/MOB-02-05-entitlements-quota.md) | Entitlements and AI quota state | R1 | M |
| 23 | [MOB-07-01](E07-uploads/MOB-07-01-image-preparation.md) | Image preparation | R1–R3 | M |
| 24 | [MOB-07-02](E07-uploads/MOB-07-02-upload-transport.md) | Upload transport | R1 | M |
| 25 | [MOB-03-03](E03-query-hooks/MOB-03-03-account-household-hooks.md) | Account and household hooks | R1–R2 | M |
| 26 | [MOB-03-04](E03-query-hooks/MOB-03-04-profile-catalogue-hooks.md) | Kitchen profile and catalogue hooks | R1–R2 | M |
| 27 | [MOB-03-05](E03-query-hooks/MOB-03-05-fridge-barcode-intake-hooks.md) | Fridge, barcode and intake hooks | R1–R3 | L |
| 28 | [MOB-03-06](E03-query-hooks/MOB-03-06-library-plans-cooking-hooks.md) | Library, plans and cooking hooks | R1–R3 | L |
| 29 | [MOB-03-07](E03-query-hooks/MOB-03-07-shopping-list-hooks.md) | Shopping list hooks | R1–R2 | M |
| 30 | [MOB-03-09](E03-query-hooks/MOB-03-09-notifications-savings-hooks.md) | Notification preference and savings hooks | R1–R3 | S |
| 31 | [MOB-06-04](E06-push/MOB-06-04-receive-and-route.md) | Receiving notifications and opening the right place | R1 | M |
| 32 | [MOB-07-03](E07-uploads/MOB-07-03-flyer-inputs.md) | Flyer inputs: photos, PDF, link | R2 | S |
| 33 | [MOB-03-08](E03-query-hooks/MOB-03-08-offers-shops-price-matrix-hooks.md) | Shops, offers, flyers and price matrix hooks | R2–R3 | M |
| 34 | [MOB-05-01](E05-live-list/MOB-05-01-live-connection.md) | Live connection manager | R2 | M |
| 35 | [MOB-05-02](E05-live-list/MOB-05-02-live-list-updates.md) | Applying live list changes | R2 | M |
| 36 | [MOB-05-03](E05-live-list/MOB-05-03-live-household-notices.md) | Live stock notices | R2 | S |

Tasks spanning releases are built with their R1 part first; the R2/R3 rows in their hook tables follow the matching
gateway epics.

## Storage on the phone

| Where | What |
|---|---|
| SecureStore | Better Auth session cookie (`coolinary` prefix), `device.installation`, `health.profile` (existing, health data) |
| Memory only | JWT, auth state, intake proposals, list proposals, price matrix, live socket |
| AsyncStorage | `coolinary.query-cache` (allow-listed queries), `coolinary.outbox.v1`, `coolinary.auth-hint` (ids only), `coolinary.device-registration`, `coolinary.push-ask`, `coolinary.settings` |
| Cache directory (temporary) | prepared photos and picked PDFs, deleted after each upload |

## New packages (all in Expo Go except where noted)

`better-auth`, `@better-auth/expo`, `expo-network`, `expo-localization`, `expo-crypto`, `expo-file-system`,
`expo-image-manipulator`, `expo-notifications` (remote push needs a development build), `expo-document-picker` (R2).
Install with `bunx expo install`.

## Decisions

Settled on 1 November 2026 (full list: [../DECISIONS.md](../DECISIONS.md)).

| ID | Decision | Tasks |
|---|---|---|
| Q020 | Send `x-client: coolinary-mobile/<version> (<platform>)` on every request; 426 `app_update_required` shows an update screen | MOB-02-01, MOB-02-02 |
| Q021 | Offline outbox entries live at most 24 hours (the gateway's replay window), then are dropped with one summary message | MOB-04-01, MOB-04-02 |
| Q022 | Own-flyer PDFs up to 40 pages and 25 MB | MOB-07-03 |
| Q031 | E-mailed sign-in codes are valid 10 minutes; the countdown starts at 600 s | MOB-01-05 |
| Q032 | 60 s between code e-mails, as identity enforces | MOB-01-05 |
| Q035 | Only registered accounts can join a household | MOB-03-03 |
| Q036 | Signing in to an existing account moves the anonymous fridge into it; anonymous plans, list and profile are dropped | MOB-01-05 |
| Q043 | The anonymous account is created lazily, on the first request that needs one | MOB-01-03 |
| Q044 | Identity deletes anonymous accounts after 30 days without use | MOB-01-03 |
| Q045 | Members' daily kcal targets never go into the query cache or AsyncStorage | MOB-03-01 |
| Q076 | Household setting `contributeToOpenFoodFacts` (default on) | MOB-03-04 |
| Q111 | Free quota is a rolling 7 days; the app shows `nextFreeAt` | MOB-02-02, MOB-02-05 |
| Q114 | A long receipt is photographed in up to 3 parts | MOB-07-01 |
| Q115 | All pages of a receipt or flyer go up in one multipart request | MOB-07-02 |
| Q116 | The server accepts phone-made UUIDs for new stock items and list lines | MOB-03-05, MOB-03-07 |
| Q117 | `useSetPlanMeal()` puts a library dish into a day-plan meal (rule-checked, no AI) | MOB-03-06 |
| Q135 | Ask for notification permission only after Plus + dated fridge items, a category turned on, or a shared household | MOB-06-02 |
| Q136 | Sign-out always deletes `health.profile` | MOB-01-07, MOB-01-08 |
| Q137 | The owner of a shared household names a new owner before deleting the account | MOB-01-08 |
| Q142, Q143 | No export for anonymous accounts; the export is a direct JSON download | MOB-01-08 |
| Q181 | Invite links are `https://coolinary.app/i/<code>` | MOB-03-03 |
| Q187 | After Plus → Free, a shared household keeps its members; only invites and joining are blocked | MOB-03-03 |
| Q219 | Typed-sentence intake ships in R3 with voice | MOB-03-05 |

## Gateway routes

The hooks use the routes in the gateway's route catalogue (docs/EPICs/gateway/README.md "Route catalogue"), which
names the backend RPC behind each one; that table is the single list to keep in step. Points the communication
layer relies on:

- GW-01: error body with a stable `code` (`validation`, `unauthorized`, `plus_required`, `not_found`,
  `version_conflict`, `quota_reached`, `rate_limited`, `ai_unavailable`, `app_update_required` (426), ...);
  `Idempotency-Key` on every write, replayable for 24 hours.
- GW-02: Better Auth under `/api/auth/*` (anonymous sign-in, e-mail code, token); account deletion is two steps,
  `POST /v1/me/deletion-code` then `DELETE /v1/me` with the code (MOB-01-08).
- GW-03: quota headers `X-Quota-Limit`, `X-Quota-Remaining`, `X-Quota-Next-Free-At` on AI routes (rolling 7 days).
- GW-07: `PUT /v1/plans/{id}/meals/{meal}` puts a picked library dish into a day plan (Q117).
- GW-08: list writes are one line per call; ticks are `POST /v1/list/lines/{id}/tick|untick`; catch-up after a gap
  is `GET /v1/list/changes?since=`.
- GW-09: WebSocket `/v1/live`; control messages carry `type`, change notices carry `kind` (`list.changed`,
  `stock.changed`) with `household_id` and `version`; close codes 4401, 4403, 4429.
- GW-10: own flyers are asynchronous (`202 { flyerId }`, poll, then confirm the draft); PDFs up to 40 pages, 25 MB.
- GW-11: `DELETE /v1/devices/{installationId}` needs no body (notifier unregisters by installation id, Q130).
