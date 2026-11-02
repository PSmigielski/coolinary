# Coolinary — roadmap

The order in which to build every task in this folder: all three releases, every service, the mobile
communication layer included. Work top to bottom. Each task comes after every task it depends on, so the
numbers are a valid build order for one developer. Tasks inside a milestone that don't depend on each other can
also be done in parallel.

- **Milestones** end with something you can try in the app or the back office ("Done when"). Release 1 ends at
  M11, release 2 at M17, release 3 at M21.
- **Estimate** uses S = 1, M = 2.5, L = 4.5 developer-days. Weeks assume 5 working days and one developer.
- **After** lists the task's dependencies (task numbers in this file).
- **Status** comes from each task file (`todo`, `in progress`, `done`, `blocked`). Work through it with
  `/next-task` in Claude Code or `scripts/epics/run-tasks.sh` (see README.md, "Working through the tasks").
- **Flags:** `decision` = the task still has an open `DECISION NEEDED:` (none since the 1 November 2026 decisions,
  see [DECISIONS.md](DECISIONS.md));
  `gate` = it must not go live until an outside condition is met (licence, retailer terms, contract).
- Built from the "Depends on" field of every task file. When a dependency names another service's epic, it points
  to the task in that epic that defines the RPC, event, route or table it names. Rebuild this file whenever
  dependencies change.

## Milestones

| Milestone | Release | Done when | Tasks | Done | Days | Weeks (cumulative) |
|---|---|---|---|---|---|---|
| [M1 Foundations](#m1-foundations) | R1 | Kitchen's event outbox and the Go consumers run, the gateway maps errors and idempotency keys, and ai answers through Claude with structured output, cost tracking and batches; the app has its typed API client. | 40 (#1–40) | 7/40 | 80.5 | 16 |
| [M2 Sign-in, households and plans](#m2-sign-in-households-and-plans) | R1 | Identity publishes household events; a phone starts anonymously, signs in with an e-mail code, keeps its household, and the gateway knows Free vs Plus; founder alerts arrive by e-mail. | 50 (#41–90) | 0/50 | 88 | 34 |
| [M3 Back office](#m3-back-office) | R1 | An admin signs in at admin.<domain>; every action is in the audit log. | 13 (#91–103) | 0/13 | 25 | 39 |
| [M4 Catalogue and kitchen profile](#m4-catalogue-and-kitchen-profile) | R1 | Ingredients are searchable in pl/de/en, names match to the catalogue, and onboarding saves cuisines, tools, diets and staples. | 38 (#104–141) | 0/38 | 77 | 54 |
| [M5 Fridge and barcodes](#m5-fridge-and-barcodes) | R1 | Typing, barcode scan and unpack mode fill the fridge, offline edits sync (Verification 14). | 35 (#142–176) | 0/35 | 82.5 | 71 |
| [M6 Nutrition and recipe library](#m6-nutrition-and-recipe-library) | R1 | Starter dishes are imported or drafted, reviewed, show calories, and 'what can I make now' works with AI off (Verification 1, 16). | 46 (#177–222) | 0/46 | 100.5 | 91 |
| [M7 AI intake](#m7-ai-intake) | R1 | Fridge photo, typed sentence and label photo become items to confirm; Free quota counts scans. | 21 (#223–243) | 0/21 | 49.5 | 101 |
| [M8 AI day plan](#m8-ai-day-plan) | R1 | A day plan comes from stock, passes the rule checker, falls back without AI, and swaps work (Verification 2-4). | 36 (#244–279) | 0/36 | 91 | 119 |
| [M9 Shopping list and cooking](#m9-shopping-list-and-cooking) | R1 | The list holds exactly the missing items, ticking moves them into the fridge, 'cooked it' reduces stock (Verification 12). | 23 (#280–302) | 0/23 | 45.5 | 128 |
| [M10 Notifications and reminders](#m10-notifications-and-reminders) | R1 | Push reaches a development build; Plus households get expiry reminders; preferences and quiet hours apply. | 23 (#303–325) | 0/23 | 43 | 136 |
| [M11 Account lifecycle and launch readiness](#m11-account-lifecycle-and-launch-readiness) | R1 | Deletion and export work across services, quotas and rate limits hold, costs stay in the caps (Verification 11). | 20 (#326–345) | 0/20 | 35.5 | 144 |
| [M12 Shops and flyer pipeline](#m12-shops-and-flyer-pipeline) | R2 | Shops from OpenStreetMap, postcode areas and editions exist; a flyer is fetched, read by batch AI, validated and held for review (Verification 8). | 50 (#346–395) | 0/50 | 113.5 | 166 |
| [M13 Launch chains go live](#m13-launch-chains-go-live) | R2 | Each launch chain passes its terms gate and the 97% sample, and syncs for two weeks (Verification 5, 7). | 5 (#396–400) | 0/5 | 12.5 | 169 |
| [M14 Offers in the app](#m14-offers-in-the-app) | R2 | Offers nearest first, my shops, own flyer scan, badges on the list, plan from deals and the price matrix (Verification 6, 9, 10, 13). | 37 (#401–437) | 0/37 | 76 | 184 |
| [M15 Shared household](#m15-shared-household) | R2 | Invites by QR, code and link, profiles without accounts, leaving and removing, live shared list (Verification 15). | 24 (#438–461) | 0/24 | 41 | 192 |
| [M16 Family goals and portions](#m16-family-goals-and-portions) | R2 | Per-person targets and one-dish portions, family day plans (Verification 17, 19, 20). | 34 (#462–495) | 0/34 | 76.5 | 208 |
| [M17 Library growth, rewrites and savings](#m17-library-growth-rewrites-and-savings) | R2 | Well-rated AI dishes join the library, tool rewrites, savings counters and monthly report. | 16 (#496–511) | 0/16 | 23.5 | 212 |
| [M18 Week plan](#m18-week-plan) | R3 | Week plan with fewest cooking sessions and leftovers (Verification 18). | 13 (#512–524) | 0/13 | 32.5 | 219 |
| [M19 Receipts and receipt prices](#m19-receipts-and-receipt-prices) | R3 | Receipt photos become items and feed chain prices into the price matrix. | 11 (#525–535) | 0/11 | 23.5 | 223 |
| [M20 Depth features](#m20-depth-features) | R3 | Cheapest top-up, restyle, share card and the remaining R3 items. | 12 (#536–547) | 0/12 | 25 | 228 |
| [M21 More chains and licensed feed](#m21-more-chains-and-licensed-feed) | R3 | Next chains and the licensed feed live behind their gates. | 11 (#548–558) | 0/11 | 32 | 235 |

Total: 558 tasks, about 1174 developer-days (235 weeks for one developer): R1 718 days, R2 343 days, R3 113 days. Frontend screens are not included yet.

## Release 1 — Launch

### M1 Foundations

**Done when:** Kitchen's event outbox and the Go consumers run, the gateway maps errors and idempotency keys, and ai answers through Claude with structured output, cost tracking and batches; the app has its typed API client.  
40 tasks, about 80.5 days.

| # | Task | Title | Service | Size | After | Flags | Status |
|---|---|---|---|---|---|---|---|
| 1 | [IDN-01-01](identity/E01-events-outbox/IDN-01-01-outbox-table-and-emit.md) | Outbox table and `emit()` helper | identity | S |  |  | done |
| 2 | [IDN-01-03](identity/E01-events-outbox/IDN-01-03-event-types.md) | Identity event types and schemas | identity | S | 1 |  | done |
| 3 | [KIT-01-01](kitchen/E01-platform/KIT-01-01-request-metadata.md) | Request metadata and household guard | kitchen | S |  |  | done |
| 4 | [KIT-01-02](kitchen/E01-platform/KIT-01-02-event-envelope-outbox.md) | Event envelope and outbox writer | kitchen | M |  |  | done |
| 5 | [KIT-01-03](kitchen/E01-platform/KIT-01-03-outbox-relay.md) | Outbox relay to Redis Streams | kitchen | M | 4 |  | done |
| 6 | [KIT-01-04](kitchen/E01-platform/KIT-01-04-stream-consumer.md) | Stream consumer with idempotent handlers | kitchen | L | 4 |  | done |
| 7 | [KIT-01-07](kitchen/E01-platform/KIT-01-07-admin-service-scaffold.md) | KitchenAdminService scaffold and change records | kitchen | M | 3 |  | done |
| 8 | [PLN-01-01](planner/E01-platform/PLN-01-01-event-consumer.md) | Event consumer (Redis Streams, Go) | planner | M | 3 |  | todo |
| 9 | [PLN-01-02](planner/E01-platform/PLN-01-02-service-clients.md) | Clients for kitchen, offers, ai, notifier, identity | planner | M |  |  | todo |
| 10 | [PLN-01-06](planner/E01-platform/PLN-01-06-admin-service.md) | PlannerAdminService scaffold and change log | planner | S |  |  | todo |
| 11 | [AI-01-01](ai/E01-claude-provider/AI-01-01-sdk-config-and-job-registry.md) | SDK dependency, settings and job registry | ai | M |  |  | todo |
| 12 | [AI-01-02](ai/E01-claude-provider/AI-01-02-provider-interface.md) | Provider interface for structured jobs | ai | M | 11 |  | todo |
| 13 | [AI-01-03](ai/E01-claude-provider/AI-01-03-claude-provider.md) | Claude provider | ai | L | 12 |  | todo |
| 14 | [AI-01-04](ai/E01-claude-provider/AI-01-04-error-mapping-and-circuit-breaker.md) | Error mapping, retries, deadlines and circuit breaker | ai | M | 13 |  | todo |
| 15 | [AI-01-05](ai/E01-claude-provider/AI-01-05-prompt-caching-layout.md) | Prompt caching and request layout | ai | M | 13 |  | todo |
| 16 | [AI-01-06](ai/E01-claude-provider/AI-01-06-output-validation-and-repair.md) | Output validation, calorie guard and one repair call | ai | M | 12 |  | todo |
| 17 | [AI-01-07](ai/E01-claude-provider/AI-01-07-untrusted-input.md) | Untrusted input handling (prompt injection) | ai | S | 12 |  | todo |
| 18 | [AI-01-08](ai/E01-claude-provider/AI-01-08-grpc-limits-and-context.md) | gRPC limits, request context and concurrency | ai | M | 11 |  | todo |
| 19 | [AI-01-09](ai/E01-claude-provider/AI-01-09-result-cache-rules.md) | Result cache rules per job | ai | S | 12 |  | todo |
| 20 | [AI-02-01](ai/E02-cost-and-spend-limit/AI-02-01-price-table-and-cost.md) | Price table and cost calculator | ai | S | 11 |  | todo |
| 21 | [AI-02-02](ai/E02-cost-and-spend-limit/AI-02-02-cost-metrics.md) | Cost and usage metrics, Grafana panels | ai | S | 20 |  | todo |
| 22 | [AI-02-03](ai/E02-cost-and-spend-limit/AI-02-03-langfuse-fields.md) | Langfuse fields for every job | ai | S | 20 |  | todo |
| 23 | [AI-02-04](ai/E02-cost-and-spend-limit/AI-02-04-daily-spend-limit.md) | Global daily spend limit | ai | M | 20 |  | todo |
| 24 | [AI-02-05](ai/E02-cost-and-spend-limit/AI-02-05-household-cost.md) | Per-household monthly cost and cost per active household | ai | M | 23 |  | todo |
| 25 | [AI-03-01](ai/E03-prompts-and-evals/AI-03-01-prompt-registry.md) | Prompt registry and versions | ai | S | 12 |  | todo |
| 26 | [AI-03-02](ai/E03-prompts-and-evals/AI-03-02-golden-set-format.md) | Golden set format, storage and privacy | ai | S | 25 |  | todo |
| 27 | [AI-03-03](ai/E03-prompts-and-evals/AI-03-03-eval-runner.md) | Eval runner CLI | ai | M | 20, 26 |  | todo |
| 28 | [AI-03-04](ai/E03-prompts-and-evals/AI-03-04-thresholds-and-policy.md) | Thresholds and release policy | ai | S | 27 |  | todo |
| 29 | [AI-03-05](ai/E03-prompts-and-evals/AI-03-05-composition-scorer.md) | Composition scorer | ai | M | 27 |  | todo |
| 30 | [AI-04-01](ai/E04-batch-jobs/AI-04-01-submit-batch.md) | SubmitBatch contract and request building | ai | L | 13, 15, 18, 23 |  | todo |
| 31 | [AI-04-02](ai/E04-batch-jobs/AI-04-02-get-and-cancel-batch.md) | GetBatch and CancelBatch | ai | S | 30 |  | todo |
| 32 | [AI-04-03](ai/E04-batch-jobs/AI-04-03-batch-results-stream.md) | GetBatchResults stream, validation and cost | ai | M | 16, 22, 31 |  | todo |
| 33 | [KIT-01-08](kitchen/E01-platform/KIT-01-08-ai-batch-tracker.md) | AI client and batch job tracker | kitchen | M | 3, 30, 31, 32 |  | todo |
| 34 | [AI-04-04](ai/E04-batch-jobs/AI-04-04-batch-failures-and-metrics.md) | Failed items, retry rules and batch metrics | ai | S | 32 |  | todo |
| 35 | [GW-01-01](gateway/E01-platform/GW-01-01-request-context.md) | Request context, metadata forwarding and deadline classes | gateway | S |  |  | todo |
| 36 | [GW-01-02](gateway/E01-platform/GW-01-02-error-model.md) | Problem JSON error model and gRPC → HTTP mapping | gateway | M | 35 |  | todo |
| 37 | [GW-01-03](gateway/E01-platform/GW-01-03-idempotency-keys.md) | Idempotency keys for mutations | gateway | M | 36 |  | todo |
| 38 | [GW-01-04](gateway/E01-platform/GW-01-04-openapi-rules.md) | OpenAPI rules, tags and export check | gateway | S | 36 |  | todo |
| 39 | [GW-01-05](gateway/E01-platform/GW-01-05-uploads.md) | Photo and PDF uploads | gateway | M | 36 |  | todo |
| 40 | [MOB-02-03](mobile/E02-api-client/MOB-02-03-openapi-workflow.md) | OpenAPI types workflow | mobile | S | 36 |  | todo |

### M2 Sign-in, households and plans

**Done when:** Identity publishes household events; a phone starts anonymously, signs in with an e-mail code, keeps its household, and the gateway knows Free vs Plus; founder alerts arrive by e-mail.  
50 tasks, about 88 days.

| # | Task | Title | Service | Size | After | Flags | Status |
|---|---|---|---|---|---|---|---|
| 41 | [IDN-02-01](identity/E02-email-code-sign-in/IDN-02-01-notifier-client.md) | Notifier gRPC client | identity | S |  |  | todo |
| 42 | [IDN-04-01](identity/E04-households-core/IDN-04-01-one-household-per-account.md) | One household per account, one owner per household | identity | S |  |  | todo |
| 43 | [IDN-04-02](identity/E04-households-core/IDN-04-02-household-created-events.md) | Household creation events and default names | identity | S | 1, 42 |  | todo |
| 44 | [IDN-04-03](identity/E04-households-core/IDN-04-03-jwt-claims.md) | JWT claims | identity | S | 42 |  | todo |
| 45 | [IDN-04-04](identity/E04-households-core/IDN-04-04-membership-checks.md) | Membership checks and richer `GetHousehold` | identity | S | 42 |  | todo |
| 46 | [IDN-04-05](identity/E04-households-core/IDN-04-05-rename.md) | Rename household and display names | identity | S | 45 |  | todo |
| 47 | [IDN-05-01](identity/E05-plan-and-entitlements/IDN-05-01-plan-changes.md) | Plan fields, plan history and `household.plan_changed` | identity | M | 1 |  | todo |
| 48 | [IDN-05-02](identity/E05-plan-and-entitlements/IDN-05-02-get-entitlements.md) | `GetEntitlements` | identity | M | 45, 47 |  | todo |
| 49 | [IDN-05-03](identity/E05-plan-and-entitlements/IDN-05-03-manual-plan-expiry.md) | Manual plan expiry | identity | S | 47 |  | todo |
| 50 | [IDN-05-05](identity/E05-plan-and-entitlements/IDN-05-05-store-subscription-link.md) | Store subscription link | identity | L | 47 | gate | todo |
| 51 | [IDN-10-01](identity/E10-admin/IDN-10-01-admin-role.md) | Admin role, `IsAdmin` and `IdentityAdminService` | identity | M | 42 |  | todo |
| 52 | [IDN-10-02](identity/E10-admin/IDN-10-02-find-account.md) | Find an account | identity | S | 47, 51 |  | todo |
| 53 | [IDN-10-03](identity/E10-admin/IDN-10-03-set-plan.md) | Set a plan by hand | identity | S | 47, 51 |  | todo |
| 54 | [IDN-10-04](identity/E10-admin/IDN-10-04-account-stats.md) | Account numbers | identity | S | 51 |  | todo |
| 55 | [NTF-02-01](notifier/E02-email-templates/NTF-02-01-template-registry.md) | Template registry and renderer | notifier | M |  | gate | todo |
| 56 | [NTF-02-02](notifier/E02-email-templates/NTF-02-02-delivery-log-privacy.md) | Delivery log without personal data | notifier | M |  |  | todo |
| 57 | [NTF-02-03](notifier/E02-email-templates/NTF-02-03-send-email-templated.md) | SendEmail with templates, languages and idempotency | notifier | M | 55, 56 |  | todo |
| 58 | [IDN-02-02](identity/E02-email-code-sign-in/IDN-02-02-email-code-plugin.md) | E-mail code plugin and code e-mail | identity | M | 41, 57 |  | todo |
| 59 | [IDN-02-03](identity/E02-email-code-sign-in/IDN-02-03-sign-in-rate-limits.md) | Sign-in rate limits and throttles | identity | M | 58 |  | todo |
| 60 | [IDN-02-04](identity/E02-email-code-sign-in/IDN-02-04-remove-password-sign-up.md) | Remove password sign-up | identity | S | 58 | gate | todo |
| 61 | [IDN-02-05](identity/E02-email-code-sign-in/IDN-02-05-session-lifetime.md) | Session lifetime and JWT issuance | identity | S | 58 |  | todo |
| 62 | [IDN-03-01](identity/E03-anonymous-accounts/IDN-03-01-anonymous-sign-in.md) | Anonymous sign-in | identity | M | 43, 59 |  | todo |
| 63 | [IDN-03-02](identity/E03-anonymous-accounts/IDN-03-02-upgrade-new-account.md) | Upgrade to a new e-mail account | identity | M | 1, 58, 62 |  | todo |
| 64 | [IDN-03-03](identity/E03-anonymous-accounts/IDN-03-03-upgrade-existing-account.md) | Upgrade into an existing account | identity | M | 63 |  | todo |
| 65 | [NTF-02-04](notifier/E02-email-templates/NTF-02-04-sign-in-code-email.md) | Sign-in code e-mail | notifier | S | 55, 57 |  | todo |
| 66 | [NTF-02-05](notifier/E02-email-templates/NTF-02-05-account-deleted-email.md) | Account deletion confirmation e-mail | notifier | S | 55, 57 |  | todo |
| 67 | [NTF-05-02](notifier/E05-founder-alerts/NTF-05-02-alert-email.md) | Founder alert e-mail and recipients | notifier | S | 55 |  | todo |
| 68 | [NTF-05-01](notifier/E05-founder-alerts/NTF-05-01-send-alert.md) | SendAlert with deduplication and throttling | notifier | M | 67 |  | todo |
| 69 | [IDN-01-02](identity/E01-events-outbox/IDN-01-02-publisher.md) | Publisher to `events.identity` | identity | M | 1, 68 |  | todo |
| 70 | [IDN-01-04](identity/E01-events-outbox/IDN-01-04-snapshot-republish.md) | Household snapshot republish | identity | M | 2, 51, 69 |  | todo |
| 71 | [NTF-01-01](notifier/E01-events-consumer/NTF-01-01-streams-consumer.md) | Redis Streams consumer (Bun) | notifier | M | 69 |  | todo |
| 72 | [NTF-01-02](notifier/E01-events-consumer/NTF-01-02-household-users-mirror.md) | Household users mirror | notifier | M | 2, 43, 71 |  | todo |
| 73 | [KIT-01-05](kitchen/E01-platform/KIT-01-05-dead-letters-lag.md) | Dead letters, lag metrics and alerts | kitchen | S | 5, 6, 68 |  | todo |
| 74 | [KIT-01-06](kitchen/E01-platform/KIT-01-06-household-plan-copy.md) | Household plan from identity | kitchen | S | 3, 48 |  | todo |
| 75 | [PLN-01-05](planner/E01-platform/PLN-01-05-household-plan-copy.md) | Household plan lookup from identity | planner | S | 9, 48 |  | todo |
| 76 | [GW-02-01](gateway/E02-auth-and-account/GW-02-01-jwt-claims.md) | JWT claims: member id, anonymous flag, registered-only routes | gateway | S | 36, 42, 62 |  | todo |
| 77 | [GW-02-02](gateway/E02-auth-and-account/GW-02-02-auth-proxy-allowlist.md) | Better Auth proxy allowlist for anonymous and e-mail code sign-in | gateway | M | 35, 41, 62 |  | todo |
| 78 | [GW-02-03](gateway/E02-auth-and-account/GW-02-03-me.md) | `GET /v1/me` | gateway | S | 76 |  | todo |
| 79 | [GW-04-01](gateway/E04-household-and-invites/GW-04-01-household-read-rename.md) | `GET` and `PATCH /v1/household` | gateway | S | 43, 46, 76 |  | todo |
| 80 | [MOB-01-01](mobile/E01-auth-session/MOB-01-01-auth-client-setup.md) | Better Auth client setup | mobile | S | 41, 62, 77 |  | todo |
| 81 | [MOB-01-02](mobile/E01-auth-session/MOB-01-02-auth-state.md) | Auth state and app start | mobile | M | 80 |  | todo |
| 82 | [MOB-01-03](mobile/E01-auth-session/MOB-01-03-anonymous-start.md) | Anonymous start | mobile | S | 62, 77, 81 |  | todo |
| 83 | [MOB-01-04](mobile/E01-auth-session/MOB-01-04-jwt-provider.md) | JWT provider and refresh | mobile | M | 77, 80, 81 |  | todo |
| 84 | [MOB-02-01](mobile/E02-api-client/MOB-02-01-request-middleware.md) | Request middleware | mobile | M | 35, 83 |  | todo |
| 85 | [MOB-01-06](mobile/E01-auth-session/MOB-01-06-session-loss.md) | 401 recovery and session loss | mobile | S | 83, 84 |  | todo |
| 86 | [MOB-02-04](mobile/E02-api-client/MOB-02-04-connectivity.md) | Connectivity and server reachability | mobile | S | 80 |  | todo |
| 87 | [MOB-02-02](mobile/E02-api-client/MOB-02-02-error-model.md) | Typed error model and retry policy | mobile | M | 35, 84, 86 |  | todo |
| 88 | [MOB-03-01](mobile/E03-query-hooks/MOB-03-01-query-keys-cache-rules.md) | Query keys and cache rules | mobile | M | 81, 87 |  | todo |
| 89 | [MOB-03-02](mobile/E03-query-hooks/MOB-03-02-invalidation-map.md) | Invalidation map | mobile | S | 88 |  | todo |
| 90 | [MOB-03-03](mobile/E03-query-hooks/MOB-03-03-account-household-hooks.md) | Account and household hooks | mobile | M | 76, 79, 88, 89 |  | todo |

### M3 Back office

**Done when:** An admin signs in at admin.<domain>; every action is in the audit log.  
13 tasks, about 25 days.

| # | Task | Title | Service | Size | After | Flags | Status |
|---|---|---|---|---|---|---|---|
| 91 | [ADM-01-01](admin/E01-service-skeleton/ADM-01-01-project-scaffold.md) | Project scaffold (Bun + Elysia) | admin | M |  |  | todo |
| 92 | [ADM-01-02](admin/E01-service-skeleton/ADM-01-02-compose-and-makefile.md) | Compose, Makefile and editor wiring | admin | M | 91 |  | todo |
| 93 | [ADM-01-03](admin/E01-service-skeleton/ADM-01-03-database-and-secrets.md) | Database, migrations and secrets | admin | S | 91 |  | todo |
| 94 | [ADM-01-04](admin/E01-service-skeleton/ADM-01-04-caddy-and-prometheus.md) | Caddy route and Prometheus scrape | admin | S | 92 |  | todo |
| 95 | [ADM-01-05](admin/E01-service-skeleton/ADM-01-05-grpc-clients-errors-openapi.md) | gRPC clients, error model and OpenAPI | admin | M | 35, 91 |  | todo |
| 96 | [ADM-02-01](admin/E02-auth-and-audit/ADM-02-01-sign-in-proxy.md) | Sign-in through the admin origin | admin | M | 41, 94, 95 |  | todo |
| 97 | [ADM-02-02](admin/E02-auth-and-audit/ADM-02-02-admin-check.md) | Admin check on every request | admin | M | 51, 95, 96 |  | todo |
| 98 | [ADM-02-03](admin/E02-auth-and-audit/ADM-02-03-audit-log.md) | Audit log | admin | M | 93, 97 |  | todo |
| 99 | [ADM-02-04](admin/E02-auth-and-audit/ADM-02-04-admin-accounts.md) | Admin accounts list | admin | S | 51, 97, 98 |  | todo |
| 100 | [ADM-05-01](admin/E05-households-support/ADM-05-01-household-lookup.md) | Household lookup | admin | M | 43, 52, 98 |  | todo |
| 101 | [ADM-05-02](admin/E05-households-support/ADM-05-02-manual-plan-switch.md) | Manual plan switch | admin | S | 47, 53, 100 |  | todo |
| 102 | [ADM-07-03](admin/E07-ai-cost-and-overview/ADM-07-03-account-stats.md) | Accounts and plans statistics | admin | S | 54, 95 |  | todo |
| 103 | [ADM-07-02](admin/E07-ai-cost-and-overview/ADM-07-02-ai-cost-view.md) | AI cost and guardrails | admin | M | 21, 23, 102 |  | todo |

### M4 Catalogue and kitchen profile

**Done when:** Ingredients are searchable in pl/de/en, names match to the catalogue, and onboarding saves cuisines, tools, diets and staples.  
38 tasks, about 77 days.

| # | Task | Title | Service | Size | After | Flags | Status |
|---|---|---|---|---|---|---|---|
| 104 | [KIT-02-01](kitchen/E02-ingredient-catalogue/KIT-02-01-ingredient-core-schema.md) | Ingredient core: categories, allergens, diet flags | kitchen | M |  |  | todo |
| 105 | [KIT-02-02](kitchen/E02-ingredient-catalogue/KIT-02-02-measures-and-units.md) | Typical weights, density and unit conversion | kitchen | M | 104 |  | todo |
| 106 | [KIT-02-03](kitchen/E02-ingredient-catalogue/KIT-02-03-storage-shelf-life-packs.md) | Storage place, shelf life and pack sizes | kitchen | S | 104 |  | todo |
| 107 | [KIT-03-01](kitchen/E03-name-matching/KIT-03-01-name-normaliser.md) | Name normaliser per language | kitchen | M |  |  | todo |
| 108 | [KIT-02-04](kitchen/E02-ingredient-catalogue/KIT-02-04-names-and-synonyms.md) | Names, synonyms and plurals per language | kitchen | M | 104, 107 |  | todo |
| 109 | [KIT-02-05](kitchen/E02-ingredient-catalogue/KIT-02-05-search-ingredients.md) | SearchIngredients (autocomplete) | kitchen | M | 108 |  | todo |
| 110 | [KIT-02-06](kitchen/E02-ingredient-catalogue/KIT-02-06-get-ingredients.md) | GetIngredients and catalogue version | kitchen | M | 105, 106, 108 |  | todo |
| 111 | [KIT-02-07](kitchen/E02-ingredient-catalogue/KIT-02-07-starter-catalogue-seed.md) | Starter catalogue seed and seed command | kitchen | M | 104, 105, 106, 108 |  | todo |
| 112 | [KIT-02-10](kitchen/E02-ingredient-catalogue/KIT-02-10-admin-catalogue-rpcs.md) | Admin catalogue editing RPCs | kitchen | M | 7, 104, 105, 106, 108 |  | todo |
| 113 | [KIT-03-02](kitchen/E03-name-matching/KIT-03-02-match-store.md) | Stored matches | kitchen | S | 104, 107 |  | todo |
| 114 | [KIT-02-09](kitchen/E02-ingredient-catalogue/KIT-02-09-merge-ingredients.md) | Merge duplicate ingredients | kitchen | M | 4, 7, 108, 113 |  | todo |
| 115 | [KIT-03-03](kitchen/E03-name-matching/KIT-03-03-match-names-deterministic.md) | MatchNames: stored, exact and trigram paths | kitchen | M | 108, 113 |  | todo |
| 116 | [KIT-03-06](kitchen/E03-name-matching/KIT-03-06-match-review-queue.md) | Match review queue (admin) | kitchen | M | 7, 113 |  | todo |
| 117 | [KIT-10-01](kitchen/E10-kitchen-profile/KIT-10-01-profile-record-and-defaults.md) | Profile record, country defaults, GetProfile and UpdateProfile | kitchen | M | 3 |  | todo |
| 118 | [KIT-10-02](kitchen/E10-kitchen-profile/KIT-10-02-cuisines-per-meal.md) | Cuisines per meal | kitchen | S | 117 |  | todo |
| 119 | [KIT-10-03](kitchen/E10-kitchen-profile/KIT-10-03-kitchen-tools-and-presets.md) | Kitchen tools and presets | kitchen | S | 117 |  | todo |
| 120 | [KIT-10-04](kitchen/E10-kitchen-profile/KIT-10-04-household-diet-allergens-dislikes.md) | Household diet, allergens and dislikes | kitchen | M | 104, 117 |  | todo |
| 121 | [KIT-10-05](kitchen/E10-kitchen-profile/KIT-10-05-staples.md) | Staples ("always at home") | kitchen | S | 104, 117 |  | todo |
| 122 | [KIT-10-06](kitchen/E10-kitchen-profile/KIT-10-06-household-rhythm-and-meal-shares.md) | Household rhythm and meal shares | kitchen | M | 117 |  | todo |
| 123 | [KIT-10-07](kitchen/E10-kitchen-profile/KIT-10-07-profile-options.md) | Profile options and onboarding defaults | kitchen | S | 118, 122 |  | todo |
| 124 | [KIT-10-08](kitchen/E10-kitchen-profile/KIT-10-08-planning-profile-snapshot.md) | GetPlanningProfile snapshot for planner | kitchen | M | 74, 118, 122 |  | todo |
| 125 | [KIT-10-09](kitchen/E10-kitchen-profile/KIT-10-09-profile-changed.md) | profile.changed event and planning version | kitchen | S | 3, 117 |  | todo |
| 126 | [PLN-01-03](planner/E01-platform/PLN-01-03-profile-snapshot.md) | Planning-profile snapshot cache | planner | M | 8, 9, 124, 125 |  | todo |
| 127 | [AI-08-01](ai/E08-catalogue-match/AI-08-01-pick-catalogue-matches.md) | PickCatalogueMatches (sync) | ai | M | 13, 16, 19 |  | todo |
| 128 | [KIT-03-04](kitchen/E03-name-matching/KIT-03-04-match-names-ai-pick.md) | MatchNames: Haiku pick and fallbacks | kitchen | M | 3, 115, 127 |  | todo |
| 129 | [KIT-03-05](kitchen/E03-name-matching/KIT-03-05-bulk-matching-batch.md) | Bulk matching through Message Batches | kitchen | M | 33, 128 |  | todo |
| 130 | [KIT-03-07](kitchen/E03-name-matching/KIT-03-07-rematch-on-catalogue-change.md) | Re-match after catalogue changes | kitchen | S | 112, 115, 129 |  | todo |
| 131 | [AI-08-02](ai/E08-catalogue-match/AI-08-02-batch-pick-matches.md) | Batch kind pick_catalogue_matches | ai | S | 32, 127 |  | todo |
| 132 | [AI-08-03](ai/E08-catalogue-match/AI-08-03-match-golden-set.md) | Golden set: 500 names | ai | M | 27, 127 |  | todo |
| 133 | [AI-09-01](ai/E09-catalogue-library-batch-jobs/AI-09-01-translate-ingredients.md) | translate_ingredients | ai | M | 32 |  | todo |
| 134 | [KIT-02-08](kitchen/E02-ingredient-catalogue/KIT-02-08-ai-translation-batch.md) | AI translation batch job | kitchen | M | 33, 108, 133 |  | todo |
| 135 | [AI-09-02](ai/E09-catalogue-library-batch-jobs/AI-09-02-match-nutrition-entries.md) | match_nutrition_entries | ai | M | 32 |  | todo |
| 136 | [ADM-03-01](admin/E03-catalogue-review/ADM-03-01-match-review-queue.md) | Name-match review queue | admin | M | 98, 116 |  | todo |
| 137 | [ADM-03-03](admin/E03-catalogue-review/ADM-03-03-ingredient-editor.md) | Ingredient editor and translation review | admin | M | 98, 109, 112 |  | todo |
| 138 | [ADM-03-04](admin/E03-catalogue-review/ADM-03-04-ingredient-merge.md) | Ingredient merge | admin | S | 114, 137 |  | todo |
| 139 | [ADM-03-06](admin/E03-catalogue-review/ADM-03-06-catalogue-jobs.md) | Catalogue batch jobs | admin | S | 23, 33, 98 |  | todo |
| 140 | [GW-05-01](gateway/E05-kitchen-profile/GW-05-01-household-profile.md) | `GET` and `PATCH /v1/profile` | gateway | M | 36, 37, 117 |  | todo |
| 141 | [GW-05-03](gateway/E05-kitchen-profile/GW-05-03-ingredient-search.md) | `GET /v1/catalogue/ingredients` (autocomplete) | gateway | S | 35, 109, 110 |  | todo |

### M5 Fridge and barcodes

**Done when:** Typing, barcode scan and unpack mode fill the fridge, offline edits sync (Verification 14).  
35 tasks, about 82.5 days.

| # | Task | Title | Service | Size | After | Flags | Status |
|---|---|---|---|---|---|---|---|
| 142 | [KIT-04-01](kitchen/E04-barcode-lookup/KIT-04-01-gtin-parsing.md) | GTIN parsing, check digits and prefix classes | kitchen | M |  |  | todo |
| 143 | [KIT-04-02](kitchen/E04-barcode-lookup/KIT-04-02-gs1-2d-codes.md) | GS1 2D codes and Digital Link | kitchen | M | 142 |  | todo |
| 144 | [KIT-04-03](kitchen/E04-barcode-lookup/KIT-04-03-internal-codes.md) | Shop-internal codes (prefixes 20–29) | kitchen | M | 117, 142 |  | todo |
| 145 | [KIT-04-04](kitchen/E04-barcode-lookup/KIT-04-04-household-products.md) | Household products (level 1) | kitchen | M | 142 |  | todo |
| 146 | [KIT-04-05](kitchen/E04-barcode-lookup/KIT-04-05-shared-catalogue.md) | Shared catalogue (level 2) | kitchen | M | 142 |  | todo |
| 147 | [KIT-04-06](kitchen/E04-barcode-lookup/KIT-04-06-open-food-facts-client.md) | Open Food Facts client (level 3) | kitchen | M | 142 | gate | todo |
| 148 | [KIT-04-07](kitchen/E04-barcode-lookup/KIT-04-07-open-food-facts-mapping.md) | Open Food Facts mapping and storage | kitchen | M | 107, 147 | gate | todo |
| 149 | [KIT-04-08](kitchen/E04-barcode-lookup/KIT-04-08-gs1-lookup.md) | GS1 register lookup (level 4) | kitchen | M | 142 | gate | todo |
| 150 | [KIT-04-09](kitchen/E04-barcode-lookup/KIT-04-09-resolve-chain.md) | ResolveBarcode: the lookup chain | kitchen | M | 143, 144, 145, 146, 147, 148 |  | todo |
| 151 | [KIT-04-10](kitchen/E04-barcode-lookup/KIT-04-10-ingredient-and-suggestions.md) | Product → ingredient, storage and use-by suggestion | kitchen | M | 106, 117, 128, 150 |  | todo |
| 152 | [KIT-04-11](kitchen/E04-barcode-lookup/KIT-04-11-multipacks.md) | Multipacks | kitchen | S | 107, 142, 146 |  | todo |
| 153 | [KIT-05-01](kitchen/E05-unknown-products/KIT-05-01-name-product.md) | Name an unknown product once | kitchen | M | 145, 146, 150 |  | todo |
| 154 | [KIT-05-02](kitchen/E05-unknown-products/KIT-05-02-submissions-and-confirmation.md) | Submissions, proposals and confirmation by a second household | kitchen | L | 153 |  | todo |
| 155 | [KIT-05-04](kitchen/E05-unknown-products/KIT-05-04-contribute-to-off.md) | Send confirmed products to Open Food Facts | kitchen | M | 117, 148, 154 | gate | todo |
| 156 | [KIT-05-05](kitchen/E05-unknown-products/KIT-05-05-admin-product-review.md) | Admin product review | kitchen | M | 7, 154 |  | todo |
| 157 | [KIT-07-01](kitchen/E07-fridge-stock/KIT-07-01-stock-schema.md) | Stock items and movements schema | kitchen | M | 3, 104, 145 |  | todo |
| 158 | [KIT-07-02](kitchen/E07-fridge-stock/KIT-07-02-storage-and-expiry-defaults.md) | Default storage place and use-by date | kitchen | M | 104, 117, 148, 157 |  | todo |
| 159 | [KIT-07-03](kitchen/E07-fridge-stock/KIT-07-03-add-stock-items.md) | AddStockItems: add, merge, idempotent adds | kitchen | M | 110, 142, 157, 158 |  | todo |
| 160 | [KIT-07-04](kitchen/E07-fridge-stock/KIT-07-04-update-and-remove.md) | Edit and remove items (used, thrown away, mistake) | kitchen | M | 157, 158 |  | todo |
| 161 | [KIT-07-05](kitchen/E07-fridge-stock/KIT-07-05-list-stock.md) | ListStock with expiry status and totals | kitchen | M | 104, 117, 157 |  | todo |
| 162 | [KIT-07-06](kitchen/E07-fridge-stock/KIT-07-06-proposals-and-unpack.md) | Stock proposals and unpack preview | kitchen | M | 104, 142, 153, 158 |  | todo |
| 163 | [KIT-07-07](kitchen/E07-fridge-stock/KIT-07-07-consume-stock.md) | ConsumeStock on "cooked it" | kitchen | L | 104, 121, 157 |  | todo |
| 164 | [KIT-07-08](kitchen/E07-fridge-stock/KIT-07-08-stock-changed.md) | stock.changed event and live notice | kitchen | S | 3, 159, 160, 163 |  | todo |
| 165 | [KIT-07-09](kitchen/E07-fridge-stock/KIT-07-09-link-unmatched-items.md) | Link items typed without a catalogue match | kitchen | M | 115, 159 |  | todo |
| 166 | [PLN-01-04](planner/E01-platform/PLN-01-04-stock-snapshot.md) | Stock snapshot cache | planner | S | 8, 9, 159, 161 |  | todo |
| 167 | [ADM-03-05](admin/E03-catalogue-review/ADM-03-05-product-review.md) | Unknown products review | admin | M | 98, 156 |  | todo |
| 168 | [ADM-03-07](admin/E03-catalogue-review/ADM-03-07-weighed-code-patterns.md) | Weighed-item barcode patterns | admin | S | 98, 144 |  | todo |
| 169 | [GW-06-01](gateway/E06-fridge-and-intake/GW-06-01-stock-routes.md) | Stock: list, add, edit, remove | gateway | M | 37, 158, 159, 160, 161, 162 |  | todo |
| 170 | [GW-06-02](gateway/E06-fridge-and-intake/GW-06-02-barcode-resolve.md) | Barcode and 2D code lookup | gateway | S | 36, 145, 148 |  | todo |
| 171 | [GW-06-03](gateway/E06-fridge-and-intake/GW-06-03-name-unknown-product.md) | Name an unknown product, confirm a shared one | gateway | S | 153, 154, 170 |  | todo |
| 172 | [MOB-04-01](mobile/E04-offline-queue/MOB-04-01-outbox-storage.md) | Outbox: storage and allowed changes | mobile | M | 81, 84 |  | todo |
| 173 | [MOB-03-04](mobile/E03-query-hooks/MOB-03-04-profile-catalogue-hooks.md) | Kitchen profile and catalogue hooks | mobile | M | 88, 89, 140, 172 |  | todo |
| 174 | [MOB-04-02](mobile/E04-offline-queue/MOB-04-02-outbox-replay.md) | Outbox replay | mobile | M | 82, 83, 86, 87, 89, 172 +1 |  | todo |
| 175 | [MOB-01-05](mobile/E01-auth-session/MOB-01-05-email-code-sign-in.md) | Sign-in with an e-mail code and anonymous upgrade | mobile | M | 41, 62, 76, 81, 83, 174 |  | todo |
| 176 | [MOB-04-04](mobile/E04-offline-queue/MOB-04-04-pending-overlay.md) | Pending changes overlay | mobile | M | 88, 172 |  | todo |

### M6 Nutrition and recipe library

**Done when:** Starter dishes are imported or drafted, reviewed, show calories, and 'what can I make now' works with AI off (Verification 1, 16).  
46 tasks, about 100.5 days.

| # | Task | Title | Service | Size | After | Flags | Status |
|---|---|---|---|---|---|---|---|
| 177 | [KIT-06-01](kitchen/E06-nutrition/KIT-06-01-nutrition-schema.md) | Nutrition entries, links and sources | kitchen | M | 104 |  | todo |
| 178 | [KIT-06-02](kitchen/E06-nutrition/KIT-06-02-import-ciqual.md) | Import Ciqual | kitchen | M | 177 | gate | todo |
| 179 | [KIT-06-03](kitchen/E06-nutrition/KIT-06-03-import-bls.md) | Import Bundeslebensmittelschlüssel 4.0 | kitchen | M | 177 | gate | todo |
| 180 | [KIT-06-04](kitchen/E06-nutrition/KIT-06-04-import-usda.md) | Import USDA FoodData Central | kitchen | M | 177 |  | todo |
| 181 | [KIT-06-05](kitchen/E06-nutrition/KIT-06-05-match-suggestions.md) | Ingredient → nutrition entry suggestions | kitchen | M | 33, 135, 177, 178 |  | todo |
| 182 | [KIT-06-07](kitchen/E06-nutrition/KIT-06-07-pack-nutrition.md) | Pack nutrition values of scanned products | kitchen | M | 148, 177 |  | todo |
| 183 | [KIT-06-08](kitchen/E06-nutrition/KIT-06-08-calculation-rules.md) | Calculation rules (library) | kitchen | M | 105, 177, 182 |  | todo |
| 184 | [KIT-06-09](kitchen/E06-nutrition/KIT-06-09-compute-nutrition-rpc.md) | ComputeNutrition RPC | kitchen | M | 183 |  | todo |
| 185 | [KIT-06-10](kitchen/E06-nutrition/KIT-06-10-nutrition-changed-event.md) | Nutrition version and `catalogue.nutrition_changed` | kitchen | S | 4, 105, 177 |  | todo |
| 186 | [KIT-06-06](kitchen/E06-nutrition/KIT-06-06-admin-nutrition-approval.md) | Admin approval of nutrition links | kitchen | M | 7, 181, 185 |  | todo |
| 187 | [PLN-02-01](planner/E02-recipe-library-model/PLN-02-01-dish-schema.md) | Dish schema (replaces skeleton `dishes`) | planner | M |  |  | todo |
| 188 | [PLN-02-02](planner/E02-recipe-library-model/PLN-02-02-dish-validation.md) | Dish validation rules | planner | M | 187 |  | todo |
| 189 | [PLN-02-03](planner/E02-recipe-library-model/PLN-02-03-derived-flags.md) | Derived diet and allergen flags | planner | M | 9, 110, 187 |  | todo |
| 190 | [PLN-02-04](planner/E02-recipe-library-model/PLN-02-04-dish-lifecycle.md) | Dish lifecycle and versions | planner | M | 187, 188 |  | todo |
| 191 | [PLN-02-05](planner/E02-recipe-library-model/PLN-02-05-portion-scaling.md) | Portion scaling rules | planner | S | 110, 187 |  | todo |
| 192 | [PLN-02-06](planner/E02-recipe-library-model/PLN-02-06-get-dish.md) | GetDish and BatchGetDishes | planner | M | 110, 187, 190, 191 |  | todo |
| 193 | [PLN-03-01](planner/E03-library-import-and-review/PLN-03-01-import-file-format.md) | Import file format and parser | planner | M | 188 |  | todo |
| 194 | [PLN-03-02](planner/E03-library-import-and-review/PLN-03-02-ingredient-name-resolution.md) | Ingredient name resolution | planner | M | 9, 115 |  | todo |
| 195 | [PLN-03-03](planner/E03-library-import-and-review/PLN-03-03-import-dishes-rpc.md) | ImportDishes admin RPC | planner | M | 10, 190, 193, 194 |  | todo |
| 196 | [PLN-04-01](planner/E04-dish-nutrition/PLN-04-01-compute-dish-nutrition.md) | Compute and store nutrition per dish version | planner | M | 9, 184, 187 |  | todo |
| 197 | [PLN-03-07](planner/E03-library-import-and-review/PLN-03-07-review-queue.md) | Review queue and decisions | planner | L | 10, 189, 190, 196 |  | todo |
| 198 | [PLN-03-08](planner/E03-library-import-and-review/PLN-03-08-dietitian-sign-off.md) | Dietitian sign-off | planner | S | 196, 197 |  | todo |
| 199 | [PLN-04-02](planner/E04-dish-nutrition/PLN-04-02-recompute-on-catalogue-events.md) | Recompute on catalogue events | planner | M | 8, 114, 186, 196 |  | todo |
| 200 | [PLN-04-03](planner/E04-dish-nutrition/PLN-04-03-nutrition-in-responses.md) | Nutrition in dish responses (kcal/kJ, rounding, incomplete) | planner | S | 192, 196 |  | todo |
| 201 | [PLN-04-04](planner/E04-dish-nutrition/PLN-04-04-scanned-product-override.md) | Pack values of the household's scanned products | planner | M | 166, 184, 196 |  | todo |
| 202 | [PLN-05-01](planner/E05-library-browse/PLN-05-01-search-dishes-rpc.md) | SearchDishes RPC (text search, scopes, sorting, paging) | planner | M | 192, 196 |  | todo |
| 203 | [PLN-05-02](planner/E05-library-browse/PLN-05-02-hard-limits.md) | Hard limits from the kitchen profile | planner | M | 126, 189, 202 |  | todo |
| 204 | [PLN-03-09](planner/E03-library-import-and-review/PLN-03-09-library-coverage-report.md) | Library coverage report | planner | M | 189, 203 |  | todo |
| 205 | [PLN-05-03](planner/E05-library-browse/PLN-05-03-filters-and-facets.md) | Filters and facet counts | planner | M | 196, 202 |  | todo |
| 206 | [PLN-05-04](planner/E05-library-browse/PLN-05-04-favourites.md) | Favourites | planner | S | 187 |  | todo |
| 207 | [PLN-05-05](planner/E05-library-browse/PLN-05-05-ratings.md) | Ratings | planner | S | 187 |  | todo |
| 208 | [PLN-05-06](planner/E05-library-browse/PLN-05-06-cooked-history.md) | Cooked history | planner | S | 187 |  | todo |
| 209 | [PLN-06-01](planner/E06-what-can-i-make-now/PLN-06-01-stock-coverage.md) | Stock coverage per dish | planner | M | 110, 126, 166, 191 |  | todo |
| 210 | [PLN-06-02](planner/E06-what-can-i-make-now/PLN-06-02-ranking.md) | Ranking: share on hand, expiry, tie-breaks | planner | M | 206, 208, 209 |  | todo |
| 211 | [PLN-06-03](planner/E06-what-can-i-make-now/PLN-06-03-suggest-dishes-rpc.md) | Extend SuggestDishes | planner | M | 203, 205, 210 |  | todo |
| 212 | [PLN-07-01](planner/E07-rule-checker/PLN-07-01-checker-core.md) | Checker core, candidate dish and verdict | planner | M | 126, 188 |  | todo |
| 213 | [PLN-07-02](planner/E07-rule-checker/PLN-07-02-tool-check.md) | Tool check | planner | S | 212 |  | todo |
| 214 | [PLN-07-03](planner/E07-rule-checker/PLN-07-03-allergen-diet-dislike-check.md) | Allergen, diet and dislike check per eater | planner | M | 110, 189, 212 |  | todo |
| 215 | [PLN-07-04](planner/E07-rule-checker/PLN-07-04-stock-check.md) | Stock and allowed-ingredient check | planner | M | 209, 212 |  | todo |
| 216 | [PLN-07-05](planner/E07-rule-checker/PLN-07-05-ai-output-checks.md) | AI output checks, time budget, metrics, 50-dish harness | planner | M | 188, 212, 215 |  | todo |
| 217 | [ADM-03-02](admin/E03-catalogue-review/ADM-03-02-nutrition-match-review.md) | Nutrition match approval | admin | M | 98, 186 |  | todo |
| 218 | [ADM-04-01](admin/E04-recipe-review/ADM-04-01-import-upload.md) | Recipe import upload | admin | M | 95, 98, 193, 195 |  | todo |
| 219 | [GW-05-02](gateway/E05-kitchen-profile/GW-05-02-profile-options.md) | `GET /v1/profile/options` (presets and reference lists) | gateway | S | 35, 123, 177 |  | todo |
| 220 | [GW-07-01](gateway/E07-recipes-and-plans/GW-07-01-library-browse.md) | Library browse and dish detail | gateway | M | 38, 196, 202 |  | todo |
| 221 | [GW-07-02](gateway/E07-recipes-and-plans/GW-07-02-favourites-ratings-history.md) | Favourites, ratings, history | gateway | S | 37, 206, 207, 208 |  | todo |
| 222 | [GW-07-03](gateway/E07-recipes-and-plans/GW-07-03-suggestions.md) | What can I make now | gateway | S | 36, 211 |  | todo |

### M7 AI intake

**Done when:** Fridge photo, typed sentence and label photo become items to confirm; Free quota counts scans.  
21 tasks, about 49.5 days.

| # | Task | Title | Service | Size | After | Flags | Status |
|---|---|---|---|---|---|---|---|
| 223 | [KIT-08-01](kitchen/E08-ai-intake/KIT-08-01-intake-pipeline.md) | Shared intake pipeline: match, propose, metrics | kitchen | M | 3, 115, 162 |  | todo |
| 224 | [AI-05-01](ai/E05-text-to-items/AI-05-01-text-to-items-claude.md) | ExtractItemsFromText on Claude | ai | M | 13, 16, 17, 19, 25 |  | todo |
| 225 | [AI-06-01](ai/E06-photo-to-items/AI-06-01-photo-to-items.md) | ExtractItemsFromPhoto | ai | L | 18, 224 |  | todo |
| 226 | [KIT-08-02](kitchen/E08-ai-intake/KIT-08-02-photo-intake.md) | Photo intake (fridge or shelf) | kitchen | M | 223, 225 |  | todo |
| 227 | [AI-06-02](ai/E06-photo-to-items/AI-06-02-photo-golden-set.md) | Golden set: 60 fridge and shelf photos | ai | M | 27, 225 |  | todo |
| 228 | [AI-07-01](ai/E07-label-reading/AI-07-01-read-product-label.md) | ReadProductLabel | ai | M | 225 |  | todo |
| 229 | [KIT-05-03](kitchen/E05-unknown-products/KIT-05-03-label-photo.md) | Read a product label photo | kitchen | M | 153, 228 |  | todo |
| 230 | [AI-07-02](ai/E07-label-reading/AI-07-02-label-golden-set.md) | Golden set: 50 labels | ai | S | 27, 228 |  | todo |
| 231 | [GW-03-01](gateway/E03-free-tier-quota/GW-03-01-quota-counters.md) | Quota counters with reserve / commit / release | gateway | M | 37, 48 |  | todo |
| 232 | [GW-03-02](gateway/E03-free-tier-quota/GW-03-02-quota-responses.md) | Quota in responses and the `quota_reached` error | gateway | S | 231 |  | todo |
| 233 | [GW-02-04](gateway/E02-auth-and-account/GW-02-04-entitlements.md) | `GET /v1/entitlements` | gateway | S | 48, 76, 232 |  | todo |
| 234 | [GW-03-03](gateway/E03-free-tier-quota/GW-03-03-plus-guard.md) | Plus-only route guard | gateway | S | 76 |  | todo |
| 235 | [GW-03-04](gateway/E03-free-tier-quota/GW-03-04-rate-limits.md) | Rate limits and Plus fair use | gateway | M | 35 |  | todo |
| 236 | [GW-03-05](gateway/E03-free-tier-quota/GW-03-05-ai-metrics.md) | AI usage metrics and dashboard panel | gateway | S | 231 |  | todo |
| 237 | [GW-06-04](gateway/E06-fridge-and-intake/GW-06-04-photo-intake.md) | Fridge or shelf photo → proposals | gateway | M | 39, 226, 231, 232 |  | todo |
| 238 | [GW-06-05](gateway/E06-fridge-and-intake/GW-06-05-label-photo.md) | Label photo for an unknown product | gateway | S | 39, 229, 231 |  | todo |
| 239 | [MOB-02-05](mobile/E02-api-client/MOB-02-05-entitlements-quota.md) | Entitlements and AI quota state | mobile | M | 47, 83, 87, 231, 233 |  | todo |
| 240 | [MOB-03-06](mobile/E03-query-hooks/MOB-03-06-library-plans-cooking-hooks.md) | Library, plans and cooking hooks | mobile | L | 88, 89, 172, 220, 239 |  | todo |
| 241 | [MOB-07-01](mobile/E07-uploads/MOB-07-01-image-preparation.md) | Image preparation | mobile | M |  |  | todo |
| 242 | [MOB-07-02](mobile/E07-uploads/MOB-07-02-upload-transport.md) | Upload transport | mobile | M | 84, 87, 169, 241 |  | todo |
| 243 | [MOB-03-05](mobile/E03-query-hooks/MOB-03-05-fridge-barcode-intake-hooks.md) | Fridge, barcode and intake hooks | mobile | L | 88, 89, 169, 172, 176, 242 |  | todo |

### M8 AI day plan

**Done when:** A day plan comes from stock, passes the rule checker, falls back without AI, and swaps work (Verification 2-4).  
36 tasks, about 91 days.

| # | Task | Title | Service | Size | After | Flags | Status |
|---|---|---|---|---|---|---|---|
| 244 | [PLN-08-01](planner/E08-ai-day-plan/PLN-08-01-plan-storage.md) | Plan storage, accept and read | planner | M | 187 |  | todo |
| 245 | [PLN-08-03](planner/E08-ai-day-plan/PLN-08-03-plan-allocation.md) | Ingredient allocation across meals | planner | M | 209 |  | todo |
| 246 | [PLN-08-05](planner/E08-ai-day-plan/PLN-08-05-library-fallback.md) | Library day plan (non-AI fallback) | planner | M | 210, 214, 245 |  | todo |
| 247 | [PLN-08-08](planner/E08-ai-day-plan/PLN-08-08-set-plan-meal.md) | SetPlanMeal: put a chosen dish into a plan meal | planner | M | 212, 213, 214, 215, 244, 245 +2 |  | todo |
| 248 | [AI-10-01](ai/E10-day-plan/AI-10-01-composition-contract.md) | Composition contract and shared dish types | ai | M | 20 |  | todo |
| 249 | [PLN-08-02](planner/E08-ai-day-plan/PLN-08-02-compose-request.md) | Compose request assembly | planner | M | 126, 166, 208, 210, 248 |  | todo |
| 250 | [PLN-08-04](planner/E08-ai-day-plan/PLN-08-04-validate-repair-store.md) | Validate, repair once, store AI dishes | planner | L | 190, 196, 216, 245, 248 |  | todo |
| 251 | [PLN-08-06](planner/E08-ai-day-plan/PLN-08-06-cuisine-choices.md) | Cuisine mismatch: three choices | planner | M | 244, 246, 250 |  | todo |
| 252 | [PLN-08-07](planner/E08-ai-day-plan/PLN-08-07-compose-day-plan-rpc.md) | ComposeDayPlan RPC, quota signal, idempotency | planner | M | 244, 251 |  | todo |
| 253 | [AI-09-03](ai/E09-catalogue-library-batch-jobs/AI-09-03-draft-dishes.md) | draft_dishes | ai | L | 32, 248 |  | todo |
| 254 | [PLN-03-04](planner/E03-library-import-and-review/PLN-03-04-ai-drafting-batch.md) | AI drafting batch (`draft_dishes`) | planner | L | 31, 32, 188, 189, 194, 253 +1 |  | todo |
| 255 | [AI-09-04](ai/E09-catalogue-library-batch-jobs/AI-09-04-translate-dishes.md) | translate_dishes | ai | M | 253 |  | todo |
| 256 | [PLN-03-05](planner/E03-library-import-and-review/PLN-03-05-dish-translation-batch.md) | Dish translation batch (`translate_dishes`) | planner | M | 30, 254, 255 |  | todo |
| 257 | [AI-09-05](ai/E09-catalogue-library-batch-jobs/AI-09-05-suggest-effort-tags.md) | suggest_effort_tags | ai | S | 32, 248 |  | todo |
| 258 | [PLN-03-06](planner/E03-library-import-and-review/PLN-03-06-effort-tag-suggestions.md) | Effort tag suggestions (`suggest_effort_tags`) | planner | S | 30, 254, 257 |  | todo |
| 259 | [AI-09-06](ai/E09-catalogue-library-batch-jobs/AI-09-06-batch-golden-sets.md) | Golden sets for catalogue and library jobs | ai | M | 27, 133, 257 |  | todo |
| 260 | [AI-10-02](ai/E10-day-plan/AI-10-02-compose-day-plan.md) | ComposeDayPlan prompt and call | ai | L | 13, 15, 19, 248 |  | todo |
| 261 | [AI-10-03](ai/E10-day-plan/AI-10-03-plan-validation.md) | Validation and repair of composed meals | ai | M | 16, 260 |  | todo |
| 262 | [AI-10-04](ai/E10-day-plan/AI-10-04-cuisine-choices.md) | Cuisine fit and the three choices | ai | M | 261 |  | todo |
| 263 | [AI-10-06](ai/E10-day-plan/AI-10-06-day-plan-golden-set.md) | Golden set and eval for day plans | ai | M | 29, 262 |  | todo |
| 264 | [AI-10-05](ai/E10-day-plan/AI-10-05-model-routing-by-plan.md) | Model routing by plan and the cost decision | ai | S | 24, 263 |  | todo |
| 265 | [AI-11-01](ai/E11-swaps-rewrites/AI-11-01-modify-meal-swap-expand.md) | ModifyMeal: contract, SWAP and EXPAND_ALTERNATIVE | ai | M | 261, 262 |  | todo |
| 266 | [PLN-09-01](planner/E09-swaps-and-rewrites/PLN-09-01-swap-meal.md) | Swap one meal | planner | M | 216, 245, 252, 265 |  | todo |
| 267 | [PLN-09-02](planner/E09-swaps-and-rewrites/PLN-09-02-replace-ingredient.md) | Replace a missing ingredient | planner | M | 209, 265, 266 |  | todo |
| 268 | [PLN-09-03](planner/E09-swaps-and-rewrites/PLN-09-03-free-text-wish.md) | Free-text wish for a meal | planner | M | 216, 265, 266 |  | todo |
| 269 | [AI-11-02](ai/E11-swaps-rewrites/AI-11-02-replace-ingredient.md) | ModifyMeal: REPLACE_INGREDIENT | ai | M | 265 |  | todo |
| 270 | [AI-11-03](ai/E11-swaps-rewrites/AI-11-03-free-text-wish.md) | ModifyMeal: WISH (free text) | ai | M | 17, 265 |  | todo |
| 271 | [AI-11-06](ai/E11-swaps-rewrites/AI-11-06-modify-golden-sets.md) | Golden sets for swaps, wishes and rewrites | ai | M | 29, 265 |  | todo |
| 272 | [ADM-04-02](admin/E04-recipe-review/ADM-04-02-ai-draft-batches.md) | AI dish draft batches | admin | M | 30, 98, 253, 254 |  | todo |
| 273 | [ADM-04-03](admin/E04-recipe-review/ADM-04-03-review-queue-and-editor.md) | Dish review queue and editor | admin | L | 10, 197, 218, 272 |  | todo |
| 274 | [ADM-04-04](admin/E04-recipe-review/ADM-04-04-publish-retire-translate.md) | Publish, retire and translate dishes | admin | M | 197, 254, 255, 256, 273 |  | todo |
| 275 | [ADM-04-05](admin/E04-recipe-review/ADM-04-05-dietitian-sign-off.md) | Dietitian sign-off | admin | S | 97, 197, 198, 273 |  | todo |
| 276 | [ADM-04-06](admin/E04-recipe-review/ADM-04-06-effort-tag-review.md) | Effort tag suggestions review | admin | S | 197, 254, 257, 258, 273 |  | todo |
| 277 | [GW-07-04](gateway/E07-recipes-and-plans/GW-07-04-ai-day-plan.md) | AI day plan and cuisine choices | gateway | M | 37, 231, 232, 249, 251 |  | todo |
| 278 | [GW-07-05](gateway/E07-recipes-and-plans/GW-07-05-plans-read-edit.md) | Read, accept, edit and delete plans | gateway | S | 244, 277 |  | todo |
| 279 | [GW-07-06](gateway/E07-recipes-and-plans/GW-07-06-swaps-wishes.md) | Swap a meal, replace an ingredient, free-text wish | gateway | M | 37, 234, 266, 267, 268 |  | todo |

### M9 Shopping list and cooking

**Done when:** The list holds exactly the missing items, ticking moves them into the fridge, 'cooked it' reduces stock (Verification 12).  
23 tasks, about 45.5 days.

| # | Task | Title | Service | Size | After | Flags | Status |
|---|---|---|---|---|---|---|---|
| 280 | [PLN-10-01](planner/E10-cooking-mode/PLN-10-01-cooking-steps.md) | Cooking mode payload | planner | S | 192, 201, 244 |  | todo |
| 281 | [PLN-10-02](planner/E10-cooking-mode/PLN-10-02-mark-cooked.md) | "Cooked it": consume stock, history, plan status | planner | M | 163, 191, 208, 245 |  | todo |
| 282 | [PLN-10-03](planner/E10-cooking-mode/PLN-10-03-rating-prompt.md) | Rating prompt after cooking | planner | S | 207, 281 |  | todo |
| 283 | [PLN-15-01](planner/E15-shopping-list/PLN-15-01-list-model.md) | List model and GetList | planner | M | 8, 110 |  | todo |
| 284 | [PLN-15-02](planner/E15-shopping-list/PLN-15-02-needs-from-plan.md) | Ingredient needs from the plan | planner | M | 187, 244, 283 |  | todo |
| 285 | [PLN-15-03](planner/E15-shopping-list/PLN-15-03-subtract-stock-and-staples.md) | Subtract stock and staples | planner | M | 124, 161, 284 |  | todo |
| 286 | [PLN-15-04](planner/E15-shopping-list/PLN-15-04-merge-and-round-to-packs.md) | Merge lines and round to packs | planner | M | 110, 285 |  | todo |
| 287 | [PLN-15-05](planner/E15-shopping-list/PLN-15-05-generate-review-accept.md) | Generate, review and accept proposed lines | planner | M | 286 |  | todo |
| 288 | [PLN-15-08](planner/E15-shopping-list/PLN-15-08-versions-idempotency-conflicts.md) | Versions, idempotent offline edits, conflict rules | planner | M | 283 |  | todo |
| 289 | [PLN-15-06](planner/E15-shopping-list/PLN-15-06-manual-lines-and-editing.md) | Manual lines, editing, removing, clearing | planner | M | 115, 283, 288 |  | todo |
| 290 | [PLN-15-07](planner/E15-shopping-list/PLN-15-07-tick-into-fridge.md) | Tick into the fridge (and undo) | planner | M | 159, 160, 283, 288 |  | todo |
| 291 | [PLN-15-09](planner/E15-shopping-list/PLN-15-09-live-notices-and-changes.md) | Live change notices and GetListChanges | planner | S | 288 |  | todo |
| 292 | [PLN-16-01](planner/E16-cuisine-base-pantry/PLN-16-01-base-pantry-lists.md) | Base pantry lists and admin edits | planner | S | 8, 110 |  | todo |
| 293 | [PLN-16-02](planner/E16-cuisine-base-pantry/PLN-16-02-base-pantry-checklist.md) | Base pantry checklist for a household | planner | S | 124, 161, 292 |  | todo |
| 294 | [PLN-16-03](planner/E16-cuisine-base-pantry/PLN-16-03-add-missing-base-items.md) | Add missing base items to the list | planner | S | 286, 288, 293 |  | todo |
| 295 | [ADM-04-08](admin/E04-recipe-review/ADM-04-08-library-settings.md) | Library settings: base pantries, no-cook extras, coverage | admin | M | 109, 204, 273, 292 |  | todo |
| 296 | [GW-07-07](gateway/E07-recipes-and-plans/GW-07-07-cooking-and-cooked.md) | Cooking mode and "cooked it" | gateway | M | 37, 280, 281, 282 |  | todo |
| 297 | [GW-08-01](gateway/E08-shopping-list/GW-08-01-list-lines.md) | Read the list; add, edit and delete lines | gateway | M | 37, 283, 289, 291 |  | todo |
| 298 | [GW-08-02](gateway/E08-shopping-list/GW-08-02-proposals.md) | Propose lines from the plan and accept them | gateway | M | 278, 284, 287, 297 |  | todo |
| 299 | [GW-08-03](gateway/E08-shopping-list/GW-08-03-tick.md) | Tick in the shop → into the fridge | gateway | S | 290, 297 |  | todo |
| 300 | [GW-08-04](gateway/E08-shopping-list/GW-08-04-base-pantry.md) | Cuisine base pantry | gateway | S | 293, 294, 297 |  | todo |
| 301 | [MOB-03-07](mobile/E03-query-hooks/MOB-03-07-shopping-list-hooks.md) | Shopping list hooks | mobile | M | 88, 89, 172, 176, 297 |  | todo |
| 302 | [MOB-04-03](mobile/E04-offline-queue/MOB-04-03-conflict-rules.md) | Conflict rules | mobile | M | 35, 169, 174, 297 |  | todo |

### M10 Notifications and reminders

**Done when:** Push reaches a development build; Plus households get expiry reminders; preferences and quiet hours apply.  
23 tasks, about 43 days.

| # | Task | Title | Service | Size | After | Flags | Status |
|---|---|---|---|---|---|---|---|
| 303 | [NTF-03-01](notifier/E03-push-delivery/NTF-03-01-push-templates.md) | Push templates and categories | notifier | M | 55 |  | todo |
| 304 | [NTF-04-01](notifier/E04-devices-preferences/NTF-04-01-devices.md) | Register and unregister devices | notifier | M |  |  | todo |
| 305 | [NTF-04-02](notifier/E04-devices-preferences/NTF-04-02-preferences.md) | Notification preferences and quiet hours | notifier | M | 303 |  | todo |
| 306 | [NTF-03-02](notifier/E03-push-delivery/NTF-03-02-send-push.md) | SendPush and the push queue | notifier | L | 56, 72, 303, 304, 305 |  | todo |
| 307 | [NTF-01-04](notifier/E01-events-consumer/NTF-01-04-account-linked.md) | Move an anonymous account's devices to the signed-in account | notifier | S | 63, 71, 72, 304, 305, 306 |  | todo |
| 308 | [NTF-03-03](notifier/E03-push-delivery/NTF-03-03-expo-delivery-worker.md) | Expo delivery worker | notifier | M | 306 |  | todo |
| 309 | [NTF-03-04](notifier/E03-push-delivery/NTF-03-04-push-receipts.md) | Push receipts and cleanup | notifier | S | 308 |  | todo |
| 310 | [NTF-04-03](notifier/E04-devices-preferences/NTF-04-03-export-user-data.md) | Export a user's notifier data | notifier | S | 56, 304, 305 |  | todo |
| 311 | [NTF-06-01](notifier/E06-admin-rpcs/NTF-06-01-delivery-log-search.md) | Delivery log search and statistics (admin) | notifier | M | 56 |  | todo |
| 312 | [NTF-06-02](notifier/E06-admin-rpcs/NTF-06-02-alerts-and-test-sends.md) | Alerts list, acknowledge and test sends (admin) | notifier | S | 57, 68, 306, 311 |  | todo |
| 313 | [KIT-09-01](kitchen/E09-expiry-reminders/KIT-09-01-reminder-selection.md) | Which items a reminder lists | kitchen | S | 74, 161 |  | todo |
| 314 | [KIT-09-02](kitchen/E09-expiry-reminders/KIT-09-02-reminder-scheduler.md) | Daily reminder scheduler per household timezone | kitchen | M | 3, 117, 313 |  | todo |
| 315 | [KIT-09-03](kitchen/E09-expiry-reminders/KIT-09-03-reminder-push.md) | Send the reminder and record it | kitchen | S | 303, 306, 314 |  | todo |
| 316 | [ADM-07-01](admin/E07-ai-cost-and-overview/ADM-07-01-system-overview.md) | System overview | admin | M | 7, 10, 95, 312 |  | todo |
| 317 | [ADM-07-04](admin/E07-ai-cost-and-overview/ADM-07-04-delivery-log-and-alerts.md) | Delivery log, alerts and test sends | admin | S | 98, 311 |  | todo |
| 318 | [GW-11-01](gateway/E11-devices-and-notifications/GW-11-01-devices.md) | Register and unregister a device | gateway | S | 76, 304 |  | todo |
| 319 | [GW-11-02](gateway/E11-devices-and-notifications/GW-11-02-notification-preferences.md) | Notification preferences and quiet hours | gateway | S | 76, 305 |  | todo |
| 320 | [MOB-03-09](mobile/E03-query-hooks/MOB-03-09-notifications-savings-hooks.md) | Notification preference and savings hooks | mobile | S | 88, 89, 318 |  | todo |
| 321 | [MOB-06-01](mobile/E06-push/MOB-06-01-dev-build-notification-config.md) | Development build and notification setup | mobile | M |  |  | todo |
| 322 | [MOB-06-02](mobile/E06-push/MOB-06-02-permission-timing.md) | When to ask for notification permission | mobile | S | 239, 321 |  | todo |
| 323 | [MOB-06-03](mobile/E06-push/MOB-06-03-device-registration.md) | Device and push token registration | mobile | M | 81, 304, 321, 322 |  | todo |
| 324 | [MOB-01-07](mobile/E01-auth-session/MOB-01-07-sign-out.md) | Sign-out and local wipe | mobile | M | 81, 88, 172, 318, 323 |  | todo |
| 325 | [MOB-06-04](mobile/E06-push/MOB-06-04-receive-and-route.md) | Receiving notifications and opening the right place | mobile | M | 89, 303, 323 |  | todo |

### M11 Account lifecycle and launch readiness

**Done when:** Deletion and export work across services, quotas and rate limits hold, costs stay in the caps (Verification 11).  
20 tasks, about 35.5 days.

| # | Task | Title | Service | Size | After | Flags | Status |
|---|---|---|---|---|---|---|---|
| 326 | [IDN-09-01](identity/E09-account-deletion-and-export/IDN-09-01-deletion-code.md) | Deletion confirmation code | identity | S | 41, 55 |  | todo |
| 327 | [IDN-09-02](identity/E09-account-deletion-and-export/IDN-09-02-delete-account.md) | Delete an account | identity | L | 1, 42, 326 |  | todo |
| 328 | [IDN-03-04](identity/E03-anonymous-accounts/IDN-03-04-stale-anonymous-cleanup.md) | Delete stale anonymous accounts | identity | M | 62, 327 |  | todo |
| 329 | [IDN-09-03](identity/E09-account-deletion-and-export/IDN-09-03-export-user-data.md) | `ExportUserData` | identity | S | 45 |  | todo |
| 330 | [NTF-01-03](notifier/E01-events-consumer/NTF-01-03-account-deleted.md) | Erase a deleted account's data | notifier | S | 71, 72, 304, 305, 306, 327 +1 |  | todo |
| 331 | [KIT-13-01](kitchen/E13-household-lifecycle/KIT-13-01-household-created.md) | household.created: default profile and stock state | kitchen | S | 6, 117 |  | todo |
| 332 | [KIT-13-02](kitchen/E13-household-lifecycle/KIT-13-02-member-joined-with-fridge.md) | Joining with a fridge: move stock, carry the person's data | kitchen | M | 2, 159, 331 |  | todo |
| 333 | [KIT-13-04](kitchen/E13-household-lifecycle/KIT-13-04-household-deleted.md) | household.deleted: erase the household | kitchen | M | 117, 142, 153, 157, 314 |  | todo |
| 334 | [KIT-13-05](kitchen/E13-household-lifecycle/KIT-13-05-account-deleted.md) | account.deleted: erase the person | kitchen | S |  |  | todo |
| 335 | [KIT-13-06](kitchen/E13-household-lifecycle/KIT-13-06-plan-changed.md) | Plan changes: planning version and effects | kitchen | S | 6, 74, 125 |  | todo |
| 336 | [KIT-13-07](kitchen/E13-household-lifecycle/KIT-13-07-export-household-data.md) | ExportHouseholdData | kitchen | M | 117, 157 |  | todo |
| 337 | [KIT-13-08](kitchen/E13-household-lifecycle/KIT-13-08-retention-job.md) | Retention clean-up job | kitchen | S | 157, 314 |  | todo |
| 338 | [KIT-13-09](kitchen/E13-household-lifecycle/KIT-13-09-ingredient-merge-household-tables.md) | Register household tables for ingredient merges | kitchen | S | 114, 120, 121, 157 |  | todo |
| 339 | [PLN-21-01](planner/E21-household-lifecycle/PLN-21-01-household-and-account-deletion.md) | Household and account deletion | planner | M | 8, 283 |  | todo |
| 340 | [PLN-21-04](planner/E21-household-lifecycle/PLN-21-04-ingredient-merge-for-lists.md) | Ingredient merge for lists, pantries and extras | planner | S | 8, 283 |  | todo |
| 341 | [PLN-21-06](planner/E21-household-lifecycle/PLN-21-06-export-household-data.md) | ExportHouseholdData | planner | M | 244, 283 |  | todo |
| 342 | [ADM-05-03](admin/E05-households-support/ADM-05-03-account-deletion-request.md) | Delete an account on request | admin | S | 100, 327 |  | todo |
| 343 | [GW-02-05](gateway/E02-auth-and-account/GW-02-05-delete-account.md) | Account deletion (`POST /v1/me/deletion-code`, `DELETE /v1/me`) | gateway | M | 37, 76, 326, 327 |  | todo |
| 344 | [GW-02-06](gateway/E02-auth-and-account/GW-02-06-data-export.md) | `POST /v1/me/export` (GDPR export) | gateway | M | 76, 310, 329, 331, 339 |  | todo |
| 345 | [MOB-01-08](mobile/E01-auth-session/MOB-01-08-account-deletion-export.md) | Account deletion and data export calls | mobile | S | 84, 90, 324, 326, 343, 344 +1 |  | todo |

## Release 2 — Flyers and family

### M12 Shops and flyer pipeline

**Done when:** Shops from OpenStreetMap, postcode areas and editions exist; a flyer is fetched, read by batch AI, validated and held for review (Verification 8).  
50 tasks, about 113.5 days.

| # | Task | Title | Service | Size | After | Flags | Status |
|---|---|---|---|---|---|---|---|
| 346 | [OFF-01-01](offers/E01-platform/OFF-01-01-event-outbox-publisher.md) | Event outbox and Redis Streams publisher | offers | M | 3 |  | todo |
| 347 | [OFF-01-02](offers/E01-platform/OFF-01-02-event-consumer.md) | Event consumer group and idempotent handlers | offers | M | 3, 346 |  | todo |
| 348 | [OFF-01-04](offers/E01-platform/OFF-01-04-service-clients.md) | Clients for identity, kitchen, ai and notifier | offers | S |  |  | todo |
| 349 | [OFF-01-03](offers/E01-platform/OFF-01-03-household-plan-copy.md) | Household plan lookup through identity | offers | S | 48, 348 |  | todo |
| 350 | [OFF-01-05](offers/E01-platform/OFF-01-05-scheduled-jobs.md) | Scheduled job runner | offers | M |  |  | todo |
| 351 | [OFF-01-06](offers/E01-platform/OFF-01-06-admin-service-shell.md) | OffersAdminService shell and change log | offers | S |  |  | todo |
| 352 | [OFF-02-01](offers/E02-chains-and-shops/OFF-02-01-chains-seed.md) | Chains table and seed | offers | S |  |  | todo |
| 353 | [OFF-02-02](offers/E02-chains-and-shops/OFF-02-02-shops-schema.md) | Shops schema rework | offers | S | 352 |  | todo |
| 354 | [OFF-02-03](offers/E02-chains-and-shops/OFF-02-03-osm-import.md) | OSM shop import job | offers | L | 350, 353 |  | todo |
| 355 | [OFF-02-04](offers/E02-chains-and-shops/OFF-02-04-shop-dedupe.md) | Merge duplicate shops | offers | S | 354 |  | todo |
| 356 | [OFF-02-05](offers/E02-chains-and-shops/OFF-02-05-user-adds-shop.md) | User adds a missing shop | offers | M | 353 |  | todo |
| 357 | [OFF-02-06](offers/E02-chains-and-shops/OFF-02-06-user-corrects-shop.md) | User proposes a shop correction | offers | S | 353 |  | todo |
| 358 | [OFF-02-07](offers/E02-chains-and-shops/OFF-02-07-shop-moderation-admin.md) | Shop moderation admin RPCs | offers | S | 351, 356, 357 |  | todo |
| 359 | [OFF-02-08](offers/E02-chains-and-shops/OFF-02-08-osm-attribution-gate.md) | OSM attribution and ODbL gate | offers | S | 354 | gate | todo |
| 360 | [OFF-03-01](offers/E03-areas-and-editions/OFF-03-01-postcode-dataset.md) | Postcode and town dataset import | offers | M | 350 |  | todo |
| 361 | [OFF-03-02](offers/E03-areas-and-editions/OFF-03-02-resolve-area.md) | Resolve a postcode or town to an area | offers | S | 360 |  | todo |
| 362 | [OFF-03-04](offers/E03-areas-and-editions/OFF-03-04-editions-model.md) | Flyer editions and shop-to-edition rules | offers | M | 352, 353 |  | todo |
| 363 | [OFF-03-03](offers/E03-areas-and-editions/OFF-03-03-shops-near-area.md) | Shops near an area | offers | M | 354, 355, 361, 362 |  | todo |
| 364 | [OFF-03-05](offers/E03-areas-and-editions/OFF-03-05-chain-store-numbers.md) | Link shops to chain store numbers | offers | M | 354, 362 |  | todo |
| 365 | [OFF-05-01](offers/E05-flyer-sync/OFF-05-01-sources-and-flyers-schema.md) | Flyer sources, flyers and versions | offers | M | 352, 362 | gate | todo |
| 366 | [OFF-05-02](offers/E05-flyer-sync/OFF-05-02-adapter-interface.md) | Adapter interface and fixture harness | offers | M | 365 |  | todo |
| 367 | [OFF-05-03](offers/E05-flyer-sync/OFF-05-03-polite-fetcher.md) | Polite fetcher | offers | M |  |  | todo |
| 368 | [OFF-05-04](offers/E05-flyer-sync/OFF-05-04-sync-and-detection.md) | Sync checks and change detection | offers | L | 350, 366, 367 |  | todo |
| 369 | [OFF-05-05](offers/E05-flyer-sync/OFF-05-05-page-storage.md) | Flyer page storage and retention | offers | M | 368 |  | todo |
| 370 | [OFF-05-06](offers/E05-flyer-sync/OFF-05-06-source-health-status.md) | Source health, alerts and flyer status for users | offers | M | 348, 368 |  | todo |
| 371 | [OFF-05-08](offers/E05-flyer-sync/OFF-05-08-sources-admin.md) | Sources and flyers admin RPCs | offers | M | 351, 365, 369 |  | todo |
| 372 | [OFF-07-01](offers/E07-offer-pipeline/OFF-07-01-offers-schema.md) | Offers schema rework | offers | M | 362, 365 |  | todo |
| 373 | [AI-13-01](ai/E13-flyer-offers/AI-13-01-offer-schema.md) | Offer schema and checks | ai | M | 16, 225 |  | todo |
| 374 | [AI-13-02](ai/E13-flyer-offers/AI-13-02-extract-flyer-offers-sync.md) | ExtractFlyerOffers for user flyers (sync, chunked) | ai | L | 18, 373 |  | todo |
| 375 | [OFF-07-02](offers/E07-offer-pipeline/OFF-07-02-batch-extraction.md) | Submit flyers for extraction and collect results | offers | L | 31, 32, 348, 369, 372, 374 +1 |  | todo |
| 376 | [OFF-07-03](offers/E07-offer-pipeline/OFF-07-03-normalise-offers.md) | Normalise prices, packs and conditions | offers | L | 375 |  | todo |
| 377 | [OFF-07-04](offers/E07-offer-pipeline/OFF-07-04-match-products.md) | Match offer products to catalogue ingredients | offers | M | 110, 115, 376 |  | todo |
| 378 | [OFF-07-08](offers/E07-offer-pipeline/OFF-07-08-duplicate-offers.md) | Duplicate and conflicting offers | offers | S | 376 |  | todo |
| 379 | [OFF-07-05](offers/E07-offer-pipeline/OFF-07-05-validation-rules.md) | Validation without AI | offers | M | 377, 378 |  | todo |
| 380 | [OFF-07-06](offers/E07-offer-pipeline/OFF-07-06-hold-back-queue.md) | Hold-back queue and admin review | offers | M | 351, 379 |  | todo |
| 381 | [OFF-07-07](offers/E07-offer-pipeline/OFF-07-07-publish-expire.md) | Publish, "next week" and expiry | offers | L | 346, 350, 379 |  | todo |
| 382 | [OFF-08-01](offers/E08-accuracy-gate/OFF-08-01-review-samples.md) | Review samples and verdicts | offers | M | 351, 369, 379 |  | todo |
| 383 | [OFF-08-02](offers/E08-accuracy-gate/OFF-08-02-go-live-gate.md) | Go-live gate per chain | offers | M | 381, 382 |  | todo |
| 384 | [OFF-08-03](offers/E08-accuracy-gate/OFF-08-03-weekly-spot-checks.md) | Weekly spot-checks | offers | S | 350, 382 |  | todo |
| 385 | [OFF-08-04](offers/E08-accuracy-gate/OFF-08-04-accuracy-record.md) | Accuracy record and metrics | offers | S | 382 |  | todo |
| 386 | [OFF-08-05](offers/E08-accuracy-gate/OFF-08-05-remove-chain.md) | Remove a chain within a day | offers | M | 369, 381 |  | todo |
| 387 | [AI-13-03](ai/E13-flyer-offers/AI-13-03-flyer-batch.md) | Batch kind extract_flyer_offers for chain flyers | ai | S | 32, 373 |  | todo |
| 388 | [AI-13-04](ai/E13-flyer-offers/AI-13-04-flyer-golden-set.md) | Golden set per launch chain | ai | M | 27, 374, 387 |  | todo |
| 389 | [ADM-06-01](admin/E06-flyer-review/ADM-06-01-source-status-board.md) | Flyer source status board | admin | M | 95, 371, 382, 385 |  | todo |
| 390 | [ADM-06-02](admin/E06-flyer-review/ADM-06-02-terms-gate-and-status.md) | Terms gate and chain status | admin | M | 371, 383, 386, 389 | gate | todo |
| 391 | [ADM-06-03](admin/E06-flyer-review/ADM-06-03-accuracy-sample.md) | Accuracy sample of 100 offers | admin | L | 371, 382, 390 |  | todo |
| 392 | [ADM-06-04](admin/E06-flyer-review/ADM-06-04-weekly-spot-check.md) | Weekly spot-check of 20 offers per chain | admin | S | 382, 384, 385, 391 |  | todo |
| 393 | [ADM-06-05](admin/E06-flyer-review/ADM-06-05-held-offers-queue.md) | Held-back offers queue | admin | M | 116, 380, 391 |  | todo |
| 394 | [ADM-06-06](admin/E06-flyer-review/ADM-06-06-shop-corrections.md) | Shop additions and corrections from users | admin | S | 98, 355, 358 |  | todo |
| 395 | [ADM-06-07](admin/E06-flyer-review/ADM-06-07-store-linking.md) | Link shops to chain store numbers | admin | S | 98, 364 |  | todo |

### M13 Launch chains go live

**Done when:** Each launch chain passes its terms gate and the 97% sample, and syncs for two weeks (Verification 5, 7).  
5 tasks, about 12.5 days.

| # | Task | Title | Service | Size | After | Flags | Status |
|---|---|---|---|---|---|---|---|
| 396 | [OFF-06-01](offers/E06-launch-chains/OFF-06-01-biedronka-pl.md) | Biedronka (PL) adapter | offers | M | 366, 367 | gate | todo |
| 397 | [OFF-06-02](offers/E06-launch-chains/OFF-06-02-lidl-pl.md) | Lidl PL adapter | offers | M | 366, 367 | gate | todo |
| 398 | [OFF-06-04](offers/E06-launch-chains/OFF-06-04-lidl-de.md) | Lidl DE adapter | offers | M | 366, 367, 397 | gate | todo |
| 399 | [OFF-06-05](offers/E06-launch-chains/OFF-06-05-aldi-sued-de.md) | Aldi Süd adapter | offers | M | 366, 367 | gate | todo |
| 400 | [OFF-06-06](offers/E06-launch-chains/OFF-06-06-aldi-nord-de.md) | Aldi Nord adapter (PDF) | offers | M | 366, 367 | gate | todo |

### M14 Offers in the app

**Done when:** Offers nearest first, my shops, own flyer scan, badges on the list, plan from deals and the price matrix (Verification 6, 9, 10, 13).  
37 tasks, about 76 days.

| # | Task | Title | Service | Size | After | Flags | Status |
|---|---|---|---|---|---|---|---|
| 401 | [OFF-04-01](offers/E04-my-shops/OFF-04-01-my-shops.md) | My shops: set and list | offers | M | 362, 363 |  | todo |
| 402 | [PLN-19-01](planner/E19-trip-price-matrix/PLN-19-01-candidate-shops.md) | Candidate shops for a trip | planner | S | 401 |  | todo |
| 403 | [OFF-04-02](offers/E04-my-shops/OFF-04-02-shopping-preferences.md) | Shopping preferences ("How you shop") | offers | S |  |  | todo |
| 404 | [OFF-04-03](offers/E04-my-shops/OFF-04-03-edition-demand.md) | Store-edition demand | offers | S | 362, 401 |  | todo |
| 405 | [OFF-04-04](offers/E04-my-shops/OFF-04-04-plan-change-rules.md) | My shops on plan changes | offers | S | 349, 350, 401 |  | todo |
| 406 | [OFF-04-05](offers/E04-my-shops/OFF-04-05-closed-merged-shops.md) | Closed and merged shops in my shops | offers | S | 355, 401 |  | todo |
| 407 | [OFF-05-07](offers/E05-flyer-sync/OFF-05-07-on-demand-editions.md) | Read store editions only on demand | offers | S | 368, 404 |  | todo |
| 408 | [OFF-06-03](offers/E06-launch-chains/OFF-06-03-kaufland-pl.md) | Kaufland PL adapter (per store) | offers | L | 364, 366, 367, 407 | gate | todo |
| 409 | [OFF-06-07](offers/E06-launch-chains/OFF-06-07-kaufland-de.md) | Kaufland DE adapter (per store) | offers | S | 408 | gate | todo |
| 410 | [OFF-09-01](offers/E09-offer-browse/OFF-09-01-browse-offers.md) | Browse offers of my shops | offers | L | 362, 381, 401 |  | todo |
| 411 | [OFF-09-02](offers/E09-offer-browse/OFF-09-02-search-offers.md) | Search offers by text | offers | M | 109, 410 |  | todo |
| 412 | [OFF-09-03](offers/E09-offer-browse/OFF-09-03-offer-detail.md) | Offer detail | offers | S | 410 |  | todo |
| 413 | [OFF-09-04](offers/E09-offer-browse/OFF-09-04-offer-categories.md) | Offer categories with counts | offers | S | 410 |  | todo |
| 414 | [OFF-09-05](offers/E09-offer-browse/OFF-09-05-find-offers-for-ingredients.md) | Find offers for ingredients (planner) | offers | M | 381, 401 |  | todo |
| 415 | [OFF-10-01](offers/E10-user-flyers/OFF-10-01-upload-and-jobs.md) | Upload photos or a PDF | offers | M | 372 |  | todo |
| 416 | [OFF-10-02](offers/E10-user-flyers/OFF-10-02-flyer-link.md) | Flyer from a link | offers | M | 415 |  | todo |
| 417 | [OFF-10-03](offers/E10-user-flyers/OFF-10-03-extract-draft.md) | Extract and check a draft | offers | M | 374, 376, 377, 379, 415 |  | todo |
| 418 | [OFF-10-04](offers/E10-user-flyers/OFF-10-04-confirm-private-offers.md) | Confirm and save private offers | offers | M | 381, 417 |  | todo |
| 419 | [OFF-10-05](offers/E10-user-flyers/OFF-10-05-shop-for-user-flyer.md) | Shop for a user flyer, in any country | offers | S | 356 |  | todo |
| 420 | [OFF-11-01](offers/E11-prices/OFF-11-01-price-observations.md) | Price observations from flyers | offers | M | 381 |  | todo |
| 421 | [OFF-11-03](offers/E11-prices/OFF-11-03-reference-prices.md) | Reference prices per country | offers | M | 350, 420 |  | todo |
| 422 | [OFF-11-04](offers/E11-prices/OFF-11-04-get-prices.md) | GetPrices for the trip price matrix | offers | M | 401, 420, 421 |  | todo |
| 423 | [PLN-17-01](planner/E17-offer-badges-and-reminders/PLN-17-01-offer-badges.md) | Offer badges on list lines | planner | M | 283, 401, 414, 422 |  | todo |
| 424 | [PLN-17-02](planner/E17-offer-badges-and-reminders/PLN-17-02-add-offer-to-list.md) | Add an offer to the list | planner | S | 288, 410, 423 |  | todo |
| 425 | [PLN-17-04](planner/E17-offer-badges-and-reminders/PLN-17-04-offer-ends-reminders.md) | "Offer ends tomorrow" reminders | planner | M | 306, 423 |  | todo |
| 426 | [PLN-19-02](planner/E19-trip-price-matrix/PLN-19-02-price-per-line-and-shop.md) | Price per line and shop | planner | M | 286, 402, 422 |  | todo |
| 427 | [PLN-19-03](planner/E19-trip-price-matrix/PLN-19-03-price-matrix-rpc.md) | Price matrix RPC and shop totals | planner | M | 283, 426 |  | todo |
| 428 | [AI-14-01](ai/E14-plan-from-deals/AI-14-01-compose-from-deals.md) | ComposeFromDeals | ai | L | 261 |  | todo |
| 429 | [AI-14-02](ai/E14-plan-from-deals/AI-14-02-deals-validation-and-golden-set.md) | Deals validation and golden set | ai | M | 29, 428 |  | todo |
| 430 | [ADM-06-08](admin/E06-flyer-review/ADM-06-08-reference-prices.md) | Reference prices | admin | S | 98, 421 |  | todo |
| 431 | [GW-10-01](gateway/E10-offers-and-shops/GW-10-01-shops-by-postcode.md) | Shops around a postcode; suggest a shop change | gateway | S | 36, 356, 357, 361, 363 |  | todo |
| 432 | [GW-10-02](gateway/E10-offers-and-shops/GW-10-02-my-shops.md) | My shops and shopping preferences | gateway | S | 401, 403, 431 |  | todo |
| 433 | [GW-10-03](gateway/E10-offers-and-shops/GW-10-03-offers-browse.md) | Offer browse, search and flyer status | gateway | M | 370, 410, 412, 413, 424, 432 |  | todo |
| 434 | [GW-10-04](gateway/E10-offers-and-shops/GW-10-04-own-flyer.md) | Scan a flyer that is not covered (async) | gateway | M | 39, 234, 415, 416, 418 |  | todo |
| 435 | [GW-10-05](gateway/E10-offers-and-shops/GW-10-05-price-matrix.md) | Trip price matrix | gateway | M | 234, 297, 402 |  | todo |
| 436 | [MOB-07-03](mobile/E07-uploads/MOB-07-03-flyer-inputs.md) | Flyer inputs: photos, PDF, link | mobile | S | 241, 242, 415, 433, 434 |  | todo |
| 437 | [MOB-03-08](mobile/E03-query-hooks/MOB-03-08-offers-shops-price-matrix-hooks.md) | Shops, offers, flyers and price matrix hooks | mobile | M | 88, 89, 431, 436 |  | todo |

### M15 Shared household

**Done when:** Invites by QR, code and link, profiles without accounts, leaving and removing, live shared list (Verification 15).  
24 tasks, about 41 days.

| # | Task | Title | Service | Size | After | Flags | Status |
|---|---|---|---|---|---|---|---|
| 438 | [IDN-06-01](identity/E06-profiles-without-accounts/IDN-06-01-add-profile.md) | Add a profile | identity | S | 1, 45, 46, 47 |  | todo |
| 439 | [IDN-06-02](identity/E06-profiles-without-accounts/IDN-06-02-delete-profile.md) | Delete a profile | identity | S | 438 |  | todo |
| 440 | [IDN-07-01](identity/E07-invites/IDN-07-01-invite-table-and-codes.md) | Invite table and code generator | identity | S | 42 |  | todo |
| 441 | [IDN-07-02](identity/E07-invites/IDN-07-02-create-list-cancel.md) | Create, list and cancel invites | identity | M | 46, 47, 440 |  | todo |
| 442 | [IDN-07-03](identity/E07-invites/IDN-07-03-invite-preview.md) | Invite preview | identity | S | 441 |  | todo |
| 443 | [IDN-07-04](identity/E07-invites/IDN-07-04-accept-invite.md) | Accept an invite | identity | L | 1, 43, 442 |  | todo |
| 444 | [IDN-05-04](identity/E05-plan-and-entitlements/IDN-05-04-shared-household-downgrade.md) | Downgrade of a shared household | identity | S | 47, 441, 443 |  | todo |
| 445 | [IDN-07-05](identity/E07-invites/IDN-07-05-invite-expiry.md) | Invite expiry and cleanup | identity | S | 440 |  | todo |
| 446 | [IDN-08-01](identity/E08-leave-and-remove/IDN-08-01-leave-household.md) | Move-out function and leaving | identity | M | 1, 43, 45 |  | todo |
| 447 | [IDN-08-02](identity/E08-leave-and-remove/IDN-08-02-remove-member.md) | Owner removes a member | identity | S | 446 |  | todo |
| 448 | [IDN-08-03](identity/E08-leave-and-remove/IDN-08-03-transfer-ownership.md) | Transfer ownership | identity | S | 42, 45 |  | todo |
| 449 | [IDN-08-04](identity/E08-leave-and-remove/IDN-08-04-revoke-live-tokens.md) | Cut off removed members' live tokens | identity | S | 76, 446 |  | todo |
| 450 | [OFF-14-01](offers/E14-household-lifecycle/OFF-14-01-erase-on-deletion.md) | Erase data on household or account deletion | offers | M | 327, 347, 446 |  | todo |
| 451 | [OFF-14-02](offers/E14-household-lifecycle/OFF-14-02-ingredient-merged.md) | Follow merged catalogue ingredients | offers | S | 114, 347 |  | todo |
| 452 | [OFF-14-03](offers/E14-household-lifecycle/OFF-14-03-export-household-data.md) | Export household data | offers | S | 401, 403, 418 |  | todo |
| 453 | [GW-04-02](gateway/E04-household-and-invites/GW-04-02-profiles-without-accounts.md) | Profiles without accounts | gateway | S | 46, 234, 438, 439 |  | todo |
| 454 | [GW-04-03](gateway/E04-household-and-invites/GW-04-03-invites-create-list-cancel.md) | Invites: create, list, cancel | gateway | M | 234, 441 |  | todo |
| 455 | [GW-04-04](gateway/E04-household-and-invites/GW-04-04-invite-preview-accept.md) | Invite preview and accept | gateway | M | 37, 235, 442, 443 |  | todo |
| 456 | [GW-04-05](gateway/E04-household-and-invites/GW-04-05-leave-remove.md) | Leave the household, remove a member | gateway | S | 37, 446, 447, 448 |  | todo |
| 457 | [GW-09-01](gateway/E09-live-list/GW-09-01-socket-auth.md) | `/v1/live` socket: auth, token refresh, limits, heartbeat | gateway | M | 76, 235 |  | todo |
| 458 | [GW-09-02](gateway/E09-live-list/GW-09-02-pubsub-relay.md) | Pub/sub relay and versioned change notices | gateway | M | 159, 291, 457 |  | todo |
| 459 | [MOB-05-01](mobile/E05-live-list/MOB-05-01-live-connection.md) | Live connection manager | mobile | M | 83, 86, 457 |  | todo |
| 460 | [MOB-05-02](mobile/E05-live-list/MOB-05-02-live-list-updates.md) | Applying live list changes | mobile | M | 291, 297, 301, 302, 457, 459 |  | todo |
| 461 | [MOB-05-03](mobile/E05-live-list/MOB-05-03-live-household-notices.md) | Live stock notices | mobile | S | 89, 457, 459 |  | todo |

### M16 Family goals and portions

**Done when:** Per-person targets and one-dish portions, family day plans (Verification 17, 19, 20).  
34 tasks, about 76.5 days.

| # | Task | Title | Service | Size | After | Flags | Status |
|---|---|---|---|---|---|---|---|
| 462 | [KIT-11-01](kitchen/E11-per-person-profiles/KIT-11-01-member-rows.md) | Member rows from identity events | kitchen | M | 2, 3, 117 |  | todo |
| 463 | [KIT-11-02](kitchen/E11-per-person-profiles/KIT-11-02-edit-rights.md) | UpdateMemberProfile and who may change what | kitchen | M | 462 |  | todo |
| 464 | [KIT-11-03](kitchen/E11-per-person-profiles/KIT-11-03-daily-kcal-target.md) | Daily kcal target: limits, consent, deletion | kitchen | M | 463 |  | todo |
| 465 | [KIT-11-04](kitchen/E11-per-person-profiles/KIT-11-04-per-person-restrictions.md) | Per-person diet, allergens and dislikes | kitchen | M | 74, 120, 335, 463 |  | todo |
| 466 | [KIT-11-05](kitchen/E11-per-person-profiles/KIT-11-05-display-and-family-flags.md) | Calories hidden and "eat lighter as a family" | kitchen | S | 122, 463 |  | todo |
| 467 | [KIT-11-06](kitchen/E11-per-person-profiles/KIT-11-06-member-visibility-and-snapshot.md) | GetMemberProfiles visibility and members in the planning snapshot | kitchen | M | 124, 464, 466 |  | todo |
| 468 | [KIT-13-03](kitchen/E13-household-lifecycle/KIT-13-03-member-removed.md) | member_removed: delete the person's data | kitchen | S | 462 |  | todo |
| 469 | [PLN-12-01](planner/E12-per-person-portions/PLN-12-01-meal-attendance.md) | Meal attendance (marking people away) | planner | M | 8, 244 |  | todo |
| 470 | [PLN-12-02](planner/E12-per-person-portions/PLN-12-02-meal-shares.md) | Meal shares per person | planner | S | 8, 464 |  | todo |
| 471 | [PLN-12-03](planner/E12-per-person-portions/PLN-12-03-portion-and-pot-calculator.md) | Portion and pot calculator | planner | M | 196, 470 |  | todo |
| 472 | [PLN-12-04](planner/E12-per-person-portions/PLN-12-04-no-cook-extras.md) | No-cook extras | planner | M | 110, 161, 184, 471 |  | todo |
| 473 | [PLN-12-05](planner/E12-per-person-portions/PLN-12-05-carry-over-and-day-balance.md) | Carry-over and day balance | planner | M | 471, 472 |  | todo |
| 474 | [PLN-12-06](planner/E12-per-person-portions/PLN-12-06-portion-allocation-for-planned-meals.md) | Portion allocation for planned meals | planner | M | 244, 469, 470, 471, 472, 473 |  | todo |
| 475 | [PLN-12-07](planner/E12-per-person-portions/PLN-12-07-calorie-visibility.md) | Who sees which numbers | planner | S | 474 |  | todo |
| 476 | [PLN-12-08](planner/E12-per-person-portions/PLN-12-08-scaled-ingredients.md) | Scaled ingredients for a planned meal | planner | S | 187, 474 |  | todo |
| 477 | [PLN-13-01](planner/E13-family-day-plan/PLN-13-01-shared-meal-limits.md) | Hard limits for a shared meal | planner | S | 212, 464 |  | todo |
| 478 | [PLN-13-02](planner/E13-family-day-plan/PLN-13-02-base-plus-toppings.md) | Base-plus-toppings plates | planner | M | 187, 476, 477 |  | todo |
| 479 | [PLN-13-03](planner/E13-family-day-plan/PLN-13-03-plan-comparator.md) | Plan comparator (priority order) | planner | M | 473 |  | todo |
| 480 | [PLN-13-04](planner/E13-family-day-plan/PLN-13-04-eat-lighter-as-a-family.md) | Eat lighter as a family | planner | S | 110, 472, 479 |  | todo |
| 481 | [PLN-13-05](planner/E13-family-day-plan/PLN-13-05-second-dish-when-forced.md) | Second dish only when an allergy or diet forces it | planner | S | 477 |  | todo |
| 482 | [PLN-13-06](planner/E13-family-day-plan/PLN-13-06-non-ai-family-day.md) | Family day without AI | planner | L | 474, 477, 478, 479, 480, 481 +1 |  | todo |
| 483 | [PLN-15-10](planner/E15-shopping-list/PLN-15-10-list-follows-plan.md) | List follows plan changes | planner | M | 287, 474 |  | todo |
| 484 | [PLN-18-01](planner/E18-plan-from-deals/PLN-18-01-deal-candidates.md) | Deal candidates for a household | planner | M | 401, 414, 477 |  | todo |
| 485 | [PLN-18-02](planner/E18-plan-from-deals/PLN-18-02-ai-plan-from-deals.md) | Plan from deals with AI | planner | L | 212, 244, 428, 474, 484 |  | todo |
| 486 | [PLN-18-03](planner/E18-plan-from-deals/PLN-18-03-non-ai-plan-from-deals.md) | Plan from deals without AI | planner | M | 482, 484 |  | todo |
| 487 | [PLN-18-04](planner/E18-plan-from-deals/PLN-18-04-deals-plan-to-list.md) | Deals plan to the shopping list | planner | S | 287, 423, 485 |  | todo |
| 488 | [AI-12-01](ai/E12-family-day-plan/AI-12-01-family-contract.md) | ComposeFamilyDay contract | ai | S | 248 |  | todo |
| 489 | [PLN-13-07](planner/E13-family-day-plan/PLN-13-07-ai-family-day.md) | Family day with AI | planner | L | 212, 244, 474, 475, 482, 488 |  | todo |
| 490 | [AI-12-02](ai/E12-family-day-plan/AI-12-02-family-compose.md) | Family composition prompt and call | ai | L | 260, 488 |  | todo |
| 491 | [AI-12-03](ai/E12-family-day-plan/AI-12-03-family-validation.md) | Validation for eaters, toppings and second dishes | ai | M | 261, 490 |  | todo |
| 492 | [AI-12-04](ai/E12-family-day-plan/AI-12-04-family-golden-set.md) | Golden set: families including the worked example | ai | M | 29, 491 |  | todo |
| 493 | [GW-05-04](gateway/E05-kitchen-profile/GW-05-04-member-profiles.md) | Per-person profiles and daily kcal targets | gateway | M | 140, 453, 463, 467 |  | todo |
| 494 | [GW-07-08](gateway/E07-recipes-and-plans/GW-07-08-portions-attendance.md) | Per-person portions and "away" marking | gateway | S | 278, 469, 474, 476, 477 |  | todo |
| 495 | [GW-07-10](gateway/E07-recipes-and-plans/GW-07-10-plan-from-deals.md) | Plan from deals | gateway | M | 234, 401, 432, 485 |  | todo |

### M17 Library growth, rewrites and savings

**Done when:** Well-rated AI dishes join the library, tool rewrites, savings counters and monthly report.  
16 tasks, about 23.5 days.

| # | Task | Title | Service | Size | After | Flags | Status |
|---|---|---|---|---|---|---|---|
| 496 | [KIT-12-01](kitchen/E12-savings/KIT-12-01-movement-valuation.md) | Grams and value of each stock movement | kitchen | S | 105, 157 |  | todo |
| 497 | [KIT-12-02](kitchen/E12-savings/KIT-12-02-ledger-rules.md) | Ledger rules and monthly totals | kitchen | M | 163, 496 |  | todo |
| 498 | [KIT-12-03](kitchen/E12-savings/KIT-12-03-get-savings.md) | GetSavings and ListSavingsReports | kitchen | S | 497 |  | todo |
| 499 | [KIT-12-04](kitchen/E12-savings/KIT-12-04-monthly-report-push.md) | Monthly savings push | kitchen | S | 303, 306, 314, 498 |  | todo |
| 500 | [PLN-11-01](planner/E11-library-growth/PLN-11-01-candidate-selection.md) | Promotion candidate job | planner | M | 75, 207, 208 |  | todo |
| 501 | [PLN-11-02](planner/E11-library-growth/PLN-11-02-anonymised-copy.md) | Anonymised copy into review | planner | M | 190, 197, 256, 258, 500 |  | todo |
| 502 | [PLN-11-03](planner/E11-library-growth/PLN-11-03-promotion-admin-and-metrics.md) | Promotion admin RPCs and savings metric | planner | S | 10, 500, 501 |  | todo |
| 503 | [PLN-21-02](planner/E21-household-lifecycle/PLN-21-02-member-added-and-removed.md) | Member added and removed | planner | S | 8, 469, 474 |  | todo |
| 504 | [PLN-21-03](planner/E21-household-lifecycle/PLN-21-03-plan-downgrade.md) | Plus → Free downgrade | planner | S | 8, 75, 126, 462, 474 |  | todo |
| 505 | [PLN-21-05](planner/E21-household-lifecycle/PLN-21-05-offers-events.md) | Offers events | planner | S | 8, 423 |  | todo |
| 506 | [PLN-17-03](planner/E17-offer-badges-and-reminders/PLN-17-03-badge-refresh-and-expiry.md) | Badge refresh and expiry | planner | S | 423, 505 |  | todo |
| 507 | [AI-11-04](ai/E11-swaps-rewrites/AI-11-04-rewrite-for-tools.md) | RewriteDish: TOOLS | ai | M | 265 |  | todo |
| 508 | [PLN-09-04](planner/E09-swaps-and-rewrites/PLN-09-04-tool-rewrite.md) | Rewrite a recipe for my tools | planner | M | 190, 196, 216, 507 |  | todo |
| 509 | [ADM-04-07](admin/E04-recipe-review/ADM-04-07-ai-dish-promotion.md) | Promote well-rated AI dishes to the library | admin | S | 273, 501, 502 |  | todo |
| 510 | [GW-07-09](gateway/E07-recipes-and-plans/GW-07-09-tool-rewrite.md) | Rewrite a dish for the tools at hand | gateway | S | 234, 508 |  | todo |
| 511 | [GW-12-01](gateway/E12-savings/GW-12-01-savings-counters-report.md) | Savings counters and monthly reports | gateway | S | 234, 498 |  | todo |

## Release 3 — Depth

### M18 Week plan

**Done when:** Week plan with fewest cooking sessions and leftovers (Verification 18).  
13 tasks, about 32.5 days.

| # | Task | Title | Service | Size | After | Flags | Status |
|---|---|---|---|---|---|---|---|
| 512 | [PLN-14-01](planner/E14-week-plan/PLN-14-01-week-plan-model.md) | Week plan model, read and edit | planner | M | 244, 474 |  | todo |
| 513 | [PLN-14-02](planner/E14-week-plan/PLN-14-02-cooking-days.md) | Cooking days from the rhythm | planner | S | 124, 512 |  | todo |
| 514 | [PLN-14-03](planner/E14-week-plan/PLN-14-03-cook-once-eat-twice.md) | Cook once, eat twice | planner | M | 187, 476, 513 |  | todo |
| 515 | [PLN-14-04](planner/E14-week-plan/PLN-14-04-effort-summary.md) | Effort summary (sessions, stove time) | planner | S | 478, 514 |  | todo |
| 516 | [PLN-14-05](planner/E14-week-plan/PLN-14-05-week-balance.md) | Week balance | planner | M | 473, 514 |  | todo |
| 517 | [PLN-14-06](planner/E14-week-plan/PLN-14-06-non-ai-week.md) | Week plan without AI | planner | L | 482, 513, 514, 515, 516 |  | todo |
| 518 | [AI-15-01](ai/E15-week-plan/AI-15-01-week-plan-contract.md) | ComposeWeekPlan contract | ai | S | 488 |  | todo |
| 519 | [PLN-14-07](planner/E14-week-plan/PLN-14-07-ai-week.md) | Week plan with AI | planner | L | 212, 512, 517, 518 |  | todo |
| 520 | [AI-15-02](ai/E15-week-plan/AI-15-02-week-plan-compose.md) | Week composition prompt and call | ai | L | 518 |  | todo |
| 521 | [AI-15-03](ai/E15-week-plan/AI-15-03-week-plan-validation.md) | Validation of sessions and leftovers | ai | M | 491, 520 |  | todo |
| 522 | [AI-15-04](ai/E15-week-plan/AI-15-04-week-plan-golden-set.md) | Golden set: family weeks | ai | M | 29, 521 |  | todo |
| 523 | [GW-07-11](gateway/E07-recipes-and-plans/GW-07-11-week-plan.md) | Week plan | gateway | M | 37, 234, 512, 513 |  | todo |
| 524 | [GW-07-12](gateway/E07-recipes-and-plans/GW-07-12-restyle.md) | Restyle a dish | gateway | S | 234, 508 |  | todo |

### M19 Receipts and receipt prices

**Done when:** Receipt photos become items and feed chain prices into the price matrix.  
11 tasks, about 23.5 days.

| # | Task | Title | Service | Size | After | Flags | Status |
|---|---|---|---|---|---|---|---|
| 525 | [KIT-08-03](kitchen/E08-ai-intake/KIT-08-03-text-intake.md) | Text and voice intake with a non-AI fallback | kitchen | M | 109, 223, 224 |  | todo |
| 526 | [KIT-08-05](kitchen/E08-ai-intake/KIT-08-05-receipt-name-memory.md) | Receipt-name memory per household and chain | kitchen | M | 159 |  | todo |
| 527 | [PLN-19-04](planner/E19-trip-price-matrix/PLN-19-04-receipt-prices.md) | Receipt prices in the matrix | planner | S | 422, 426 |  | todo |
| 528 | [AI-16-01](ai/E16-receipt-reading/AI-16-01-extract-receipt.md) | ExtractReceipt | ai | M | 225 |  | todo |
| 529 | [KIT-08-04](kitchen/E08-ai-intake/KIT-08-04-receipt-intake.md) | Receipt intake and line matching | kitchen | L | 115, 142, 223, 526, 528 |  | todo |
| 530 | [KIT-08-06](kitchen/E08-ai-intake/KIT-08-06-receipt-scanned-event.md) | receipt.scanned event for shop prices | kitchen | M | 3, 529 |  | todo |
| 531 | [OFF-11-02](offers/E11-prices/OFF-11-02-receipt-prices.md) | Receipt prices from scanned receipts | offers | M | 347, 420, 530 |  | todo |
| 532 | [AI-16-02](ai/E16-receipt-reading/AI-16-02-receipt-privacy-and-totals.md) | Receipt privacy filter and totals check | ai | S | 528 |  | todo |
| 533 | [AI-16-03](ai/E16-receipt-reading/AI-16-03-receipt-golden-set.md) | Golden set: 40 receipts | ai | M | 27, 532 |  | todo |
| 534 | [GW-06-06](gateway/E06-fridge-and-intake/GW-06-06-text-intake.md) | Typed or spoken sentence → proposals | gateway | S | 37, 231, 525 |  | todo |
| 535 | [GW-06-07](gateway/E06-fridge-and-intake/GW-06-07-receipt-intake.md) | Receipt photo → proposals | gateway | S | 39, 231, 529 |  | todo |

### M20 Depth features

**Done when:** Cheapest top-up, restyle, share card and the remaining R3 items.  
12 tasks, about 25 days.

| # | Task | Title | Service | Size | After | Flags | Status |
|---|---|---|---|---|---|---|---|
| 536 | [KIT-12-05](kitchen/E12-savings/KIT-12-05-share-card.md) | Shareable savings card data | kitchen | S | 498 |  | todo |
| 537 | [PLN-09-05](planner/E09-swaps-and-rewrites/PLN-09-05-restyle.md) | Restyle a dish ("make this Italian-style") | planner | M | 508 |  | todo |
| 538 | [PLN-20-01](planner/E20-cheapest-top-up/PLN-20-01-needs-and-allowed-swaps.md) | Top-up needs and allowed swaps | planner | M | 212, 284, 285, 484 |  | todo |
| 539 | [PLN-20-03](planner/E20-cheapest-top-up/PLN-20-03-non-ai-top-up.md) | Top-up without AI | planner | S | 426, 538 |  | todo |
| 540 | [AI-05-02](ai/E05-text-to-items/AI-05-02-text-golden-set.md) | Golden set: typed and spoken item lists | ai | S | 27, 224 |  | todo |
| 541 | [AI-11-05](ai/E11-swaps-rewrites/AI-11-05-restyle.md) | RewriteDish: RESTYLE | ai | S | 507 |  | todo |
| 542 | [AI-17-01](ai/E17-cheapest-top-up/AI-17-01-cheapest-top-up.md) | CheapestTopUp | ai | L | 428 |  | todo |
| 543 | [PLN-20-02](planner/E20-cheapest-top-up/PLN-20-02-ai-cheapest-top-up.md) | Cheapest top-up with AI | planner | L | 426, 538, 542 |  | todo |
| 544 | [PLN-20-04](planner/E20-cheapest-top-up/PLN-20-04-apply-top-up.md) | Applying a top-up | planner | M | 474, 483, 512, 543 |  | todo |
| 545 | [AI-17-02](ai/E17-cheapest-top-up/AI-17-02-top-up-validation-and-golden-set.md) | Top-up validation and golden set | ai | M | 29, 542 |  | todo |
| 546 | [GW-10-06](gateway/E10-offers-and-shops/GW-10-06-cheapest-top-up.md) | Cheapest top-up | gateway | S | 435, 523, 543, 544 |  | todo |
| 547 | [GW-12-02](gateway/E12-savings/GW-12-02-savings-card.md) | Shareable savings card data | gateway | S | 511, 536 |  | todo |

### M21 More chains and licensed feed

**Done when:** Next chains and the licensed feed live behind their gates.  
11 tasks, about 32 days.

| # | Task | Title | Service | Size | After | Flags | Status |
|---|---|---|---|---|---|---|---|
| 548 | [OFF-12-01](offers/E12-next-chains/OFF-12-01-dino-pl.md) | Dino (PL) adapter | offers | M | 366, 367 | gate | todo |
| 549 | [OFF-12-02](offers/E12-next-chains/OFF-12-02-netto-pl.md) | Netto PL adapter | offers | M | 366, 367 | gate | todo |
| 550 | [OFF-12-03](offers/E12-next-chains/OFF-12-03-carrefour-pl.md) | Carrefour PL adapter (per store format) | offers | L | 362, 366, 367 | gate | todo |
| 551 | [OFF-12-04](offers/E12-next-chains/OFF-12-04-auchan-pl.md) | Auchan PL adapter | offers | M | 366, 367 | gate | todo |
| 552 | [OFF-12-05](offers/E12-next-chains/OFF-12-05-aldi-pl.md) | Aldi PL adapter | offers | M | 366, 367 | gate | todo |
| 553 | [OFF-12-06](offers/E12-next-chains/OFF-12-06-penny-de.md) | Penny DE adapter | offers | M | 366, 367 | gate | todo |
| 554 | [OFF-13-01](offers/E13-licensed-feed/OFF-13-01-feed-client-mapping.md) | Feed client and retailer/store mapping | offers | L | 352, 364 | gate | todo |
| 555 | [OFF-13-02](offers/E13-licensed-feed/OFF-13-02-feed-ingest.md) | Ingest feed offers | offers | L | 376, 377, 379, 381, 554 |  | todo |
| 556 | [OFF-13-03](offers/E13-licensed-feed/OFF-13-03-per-store-chains-via-feed.md) | Per-store chains through the feed (Rewe, Edeka, Netto DE) | offers | M | 407, 555 |  | todo |
| 557 | [OFF-13-04](offers/E13-licensed-feed/OFF-13-04-feed-reporting.md) | View and click reporting to the network | offers | M | 555 | gate | todo |
| 558 | [ADM-06-09](admin/E06-flyer-review/ADM-06-09-feed-mapping.md) | Licensed feed: retailer and store mapping | admin | S | 395, 554 | gate | todo |

## Later releases that extend earlier tasks

These earlier tasks name work from a later release (usually a table or hook added later). Build the earlier task
in its own release; when you build the later task, also extend the earlier one as described.

| Earlier task (#) | Release | Extended by (#) | Release |
|---|---|---|---|
| [ADM-04-08](admin/E04-recipe-review/ADM-04-08-library-settings.md) (295) | R1 | [PLN-12-04](planner/E12-per-person-portions/PLN-12-04-no-cook-extras.md) (472) | R2 |
| [ADM-07-01](admin/E07-ai-cost-and-overview/ADM-07-01-system-overview.md) (316) | R1 | [OFF-01-06](offers/E01-platform/OFF-01-06-admin-service-shell.md) (351) | R2 |
| [AI-11-06](ai/E11-swaps-rewrites/AI-11-06-modify-golden-sets.md) (271) | R1 | [AI-11-05](ai/E11-swaps-rewrites/AI-11-05-restyle.md) (541) | R3 |
| [GW-02-06](gateway/E02-auth-and-account/GW-02-06-data-export.md) (344) | R1 | [OFF-14-03](offers/E14-household-lifecycle/OFF-14-03-export-household-data.md) (452) | R2 |
| [GW-07-04](gateway/E07-recipes-and-plans/GW-07-04-ai-day-plan.md) (277) | R1 | [PLN-13-06](planner/E13-family-day-plan/PLN-13-06-non-ai-family-day.md) (482) | R2 |
| [KIT-13-02](kitchen/E13-household-lifecycle/KIT-13-02-member-joined-with-fridge.md) (332) | R1 | [KIT-11-01](kitchen/E11-per-person-profiles/KIT-11-01-member-rows.md) (462) | R2 |
| [KIT-13-04](kitchen/E13-household-lifecycle/KIT-13-04-household-deleted.md) (333) | R1 | [KIT-08-05](kitchen/E08-ai-intake/KIT-08-05-receipt-name-memory.md) (526) | R3 |
| [KIT-13-04](kitchen/E13-household-lifecycle/KIT-13-04-household-deleted.md) (333) | R1 | [KIT-11-01](kitchen/E11-per-person-profiles/KIT-11-01-member-rows.md) (462) | R2 |
| [KIT-13-04](kitchen/E13-household-lifecycle/KIT-13-04-household-deleted.md) (333) | R1 | [KIT-12-01](kitchen/E12-savings/KIT-12-01-movement-valuation.md) (496) | R2 |
| [KIT-13-05](kitchen/E13-household-lifecycle/KIT-13-05-account-deleted.md) (334) | R1 | [KIT-11-01](kitchen/E11-per-person-profiles/KIT-11-01-member-rows.md) (462) | R2 |
| [KIT-13-06](kitchen/E13-household-lifecycle/KIT-13-06-plan-changed.md) (335) | R1 | [KIT-11-04](kitchen/E11-per-person-profiles/KIT-11-04-per-person-restrictions.md) (465) | R2 |
| [KIT-13-07](kitchen/E13-household-lifecycle/KIT-13-07-export-household-data.md) (336) | R1 | [KIT-08-05](kitchen/E08-ai-intake/KIT-08-05-receipt-name-memory.md) (526) | R3 |
| [KIT-13-07](kitchen/E13-household-lifecycle/KIT-13-07-export-household-data.md) (336) | R1 | [KIT-11-06](kitchen/E11-per-person-profiles/KIT-11-06-member-visibility-and-snapshot.md) (467) | R2 |
| [KIT-13-07](kitchen/E13-household-lifecycle/KIT-13-07-export-household-data.md) (336) | R1 | [KIT-12-01](kitchen/E12-savings/KIT-12-01-movement-valuation.md) (496) | R2 |
| [KIT-13-08](kitchen/E13-household-lifecycle/KIT-13-08-retention-job.md) (337) | R1 | [KIT-08-05](kitchen/E08-ai-intake/KIT-08-05-receipt-name-memory.md) (526) | R3 |
| [KIT-13-08](kitchen/E13-household-lifecycle/KIT-13-08-retention-job.md) (337) | R1 | [KIT-08-06](kitchen/E08-ai-intake/KIT-08-06-receipt-scanned-event.md) (530) | R3 |
| [KIT-13-09](kitchen/E13-household-lifecycle/KIT-13-09-ingredient-merge-household-tables.md) (338) | R1 | [KIT-08-05](kitchen/E08-ai-intake/KIT-08-05-receipt-name-memory.md) (526) | R3 |
| [MOB-03-09](mobile/E03-query-hooks/MOB-03-09-notifications-savings-hooks.md) (320) | R1 | [GW-12-01](gateway/E12-savings/GW-12-01-savings-counters-report.md) (511) | R2 |

## Dependency loops resolved

Two task files depend on each other in a loop. The roadmap breaks each loop at the link shown last, which is
safe to build in either order; fix the "Depends on" fields when you start these tasks.

- GW-02-01 -> IDN-03-01 -> IDN-02-03 -> GW-02-01
- KIT-08-04 -> KIT-08-05 -> KIT-08-04

