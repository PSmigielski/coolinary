# Coolinary — EPICs and tasks

Every back-end feature from `docs/features.md` (releases 1–3) broken into atomic tasks, one file per task, grouped
by service and epic. Mobile covers only the communication layer for now; screens come later.

Read with: `docs/business-model.md`, `docs/features.md` (the product rules), `docs/infrastructure.md` and
`docs/service-conventions.md` (how services are built).

## How to read this folder

```
docs/EPICs/
  README.md                         this file: decisions, conventions, epic map, cross-service calls and events
  TEMPLATE.md                       the shape of every task file
  ROADMAP.md                        every task numbered in build order, grouped into 21 milestones
  DECISIONS.md                      every open small decision, with a recommended answer, by milestone
  <service>/README.md               the service's epics and tasks in build order, plus its "Decision needed" list
  <service>/E07-fridge-stock/
    README.md                       epic goal, flow, task list
    KIT-07-03-consume-stock.md      one task
```

- **Where to start:** [ROADMAP.md](ROADMAP.md) gives the order to build everything, milestone by milestone.
- **Task ID** `KIT-07-03` = service prefix, epic 07, task 03. **Epic ID** `KIT-07`.
- Prefixes: `GW` gateway, `IDN` identity, `NTF` notifier, `ADM` admin, `KIT` kitchen, `PLN` planner, `OFF` offers,
  `AI` ai, `MOB` mobile communication layer.
- **Release** `R1` launch, `R2` flyers and family, `R3` depth (`docs/features.md` section 7). Foundation work needed
  by R1 is tagged R1. An epic may hold tasks of several releases.
- **Size** `S` ≤ 1 day, `M` 2–3 days, `L` 4–5 days, for one developer. Nothing is bigger than L; bigger work is split.
- **Atomic** means: one service, one mergeable change, testable on its own once its dependencies are done.
- **Decided (Qnnn):** the 228 small choices the product docs left open were settled on 1 November 2026 (critical
  and high ones by you, the rest at their recommended answer). Each task states its rule as `Decided (Qnnn):`;
  [DECISIONS.md](DECISIONS.md) lists every answer. A marker still reading `DECISION NEEDED:` is genuinely open.

## Working through the tasks

Tasks are built in ROADMAP.md order, in parallel batches, by Claude Code using the project's `/next-task` command
(`.claude/skills/next-task/`, rules in `CLAUDE.md`).

1. Start the dev stack (`make dev`) and make sure the working tree is clean: `/next-task` refuses to start on
   changes that aren't from a task.
2. **Interactive:** run `claude` in the repo and type `/next-task`: it builds up to 4 ready tasks of the current
   milestone at once (`/next-task 6` for up to 6, `/next-task KIT-07-03 KIT-07-04` for those). Review, then repeat.
3. **Unattended:** `scripts/epics/run-tasks.sh` runs one `claude -p "/next-task <ID> ..."` per batch of up to
   `--jobs` tasks (default 4; `--jobs 1` builds one at a time), with an optional `--budget` cap in dollars per run.
   It stops at the end of each milestone, on a blocked task or a failed run. `--dry-run` shows what it would do.
   Run logs and costs: `.claude/state/runs/`.
4. Each milestone is built on its own branch (`m01-foundations`, …), one commit per task. Review and merge it, then
   continue.

How it fits together:

- `scripts/epics/tasks.py` holds the state: `next`, `ready [N]`, `show <ID>`, `start`, `done`,
  `block <ID> "<why>"`, `status`. Status lives in each task file's `Status` row; ROADMAP.md shows it after
  `scripts/epics/roadmap.sh`.
- In a batch, each task is built by a subagent in its own git worktree (`.claude/worktrees/<ID>`, branch
  `task/<ID>`, made by `scripts/epics/worktree.sh`), whose tests use their own `<svc>_test_<id>` databases. The
  main session then lands the tasks one at a time in roadmap order as one commit each. It resolves mechanical
  conflicts (regenerated code, lockfiles, a migration number taken twice, lines both sides only added). A conflict
  that needs a choice sends the task back to `todo` for the next batch. A worker never asks or decides: it stops,
  and the main session asks you (interactive) or blocks the task (unattended).
- A task is ready when every task in its "After" column is `done`. A `blocked` task is skipped; tasks that depend on
  it wait. Unblock one by fixing the cause and running `tasks.py set <ID> todo`.
- `make lint s=<svc>` and `make test s=<svc>` must pass for every service a task or batch touched. The Stop hook
  (`.claude/hooks/verify-task.sh`) checks this once the tasks are marked done, sends Claude back up to three times,
  and after the third failure blocks the tasks whose service still fails.
- Claude never makes a design decision on its own: an unclear task is blocked with a note saying what needs
  deciding (under "## Blocked" in the task file).

## Decisions taken (1 November 2026)

| # | Topic | Decision |
|---|---|---|
| 1 | Scope | All three releases, every task detailed. |
| 2 | Recipes | **planner** owns the recipe library (dishes, ingredients, steps, tags, calories per dish, review). Calories come from kitchen's nutrition calculation. `docs/infrastructure.md` still says kitchen; to be updated. |
| 3 | Kitchen profile | **kitchen** owns it: cuisines, tools, diets and allergens per person, dislikes, staples, household settings (portions, main meal, time per meal, cooking rhythm, meal shares), each person's daily kcal target and display flags. Planner reads it from kitchen. |
| 4 | Cross-service changes | **Redis Streams** on `redis-state`, written through a transactional outbox per service; consumer groups, at-least-once, idempotent handlers. No new container. |
| 5 | AI provider | **Claude API direct**, official Python SDK (`anthropic`). Message Batches (50% off) for every once-for-everyone job. Inference may run outside the EU; the privacy policy must say so. |
| 6 | Models | **Claude Haiku 5.5** (`claude-haiku-5-5`) reads and matches: text, photos, labels, receipts, flyers, catalogue matching, translation. **Claude Sonnet 5.5** (`claude-sonnet-5-5`) composes for Plus: day, family and week plans, swaps, rewrites, deals, top-up. Each job's model and effort are config values. Amended by decision 15. |
| 7 | Trip plan | Runs **on the phone**. Planner returns a price matrix (each list line × each candidate shop); home position never leaves the phone. The phone algorithm is a later frontend task. |
| 8 | Payments | **Decided later.** Tasks cover the plan flag, entitlements and a manual plan switch in admin; the store link is a placeholder. |
| 9 | Back office | New **admin** service (Bun + Elysia), reached at `admin.<domain>` through Caddy, only for accounts identity marks as admin. Each owning service exposes an `<Svc>AdminService` gRPC service that only admin calls. The admin UI comes with the frontend tasks. |
| 10 | Starter dishes | **Both**: a hand-written import file and an AI drafting batch job; both end in the same review queue. |
| 11 | First use | **Anonymous account** (Better Auth anonymous plugin) on first start; signing in later upgrades it and keeps the household. |
| 12 | Sign-in | **E-mail code only** (Better Auth e-mail OTP). Password sign-up is removed. |
| 13 | Name matching | **pg_trgm candidates + Haiku picks one or "none"**; each normalised name is matched once for everyone and stored. No embeddings, no pgvector. |
| 14 | Mobile now | Communication layer only: auth session, typed API client and query hooks, offline queue, live list socket, push registration, uploads. |
| 15 | Plan model by plan | A Sonnet day plan costs about $0.05, so 2 a week would put a Free household near $0.44/month against the €0.10 cap. **Free day plans compose on Haiku 5.5.** Plus compositions use Sonnet 5.5 only if it clearly beats Haiku on the golden set (AI-03); otherwise Haiku too. |
| 16 | Shared enums | `Cuisine`, `KitchenTool`, `Diet`, `Allergen`, `EffortTag` and the extra `Unit` values (tablespoon, teaspoon, pinch, to taste) live in `contracts/coolinary/common/v1/`. Every service imports one definition. |
| 17 | AI spend counters | `ai` gets `REDIS_STATE_URL` and keeps its daily and monthly spend counters in redis-state, so the hard spend limit (D11) works. `ai` still has no database. |
| 18 | Error reasons | Services set trailing gRPC metadata `x-error-code` with a stable snake_case reason (`plus_required`, `quota_reached`, `ai_refused`, ...). Version conflicts on shared data use `ABORTED`. The gateway maps both into its HTTP error body (GW-01). |

## Derived decisions (confirmed 1 November 2026 as Q001–Q013; D5 and D6 changed by Q004 and Q010)

| # | Decision | Why |
|---|---|---|
| D1 | New call-graph edges: `offers → kitchen` (name matching), `planner → notifier` (offer-ends reminders), `admin → identity, notifier, kitchen, planner, offers, ai`, and `kitchen, planner, offers → identity` for the plan (Q004). Still no cycles. | Matching lives in kitchen; reminders need the list, which lives in planner. |
| D2 | `admin` gets IP `10.77.0.23`, HTTP `3000`, admin port `8081`, its own database `admin` (audit log only). | Next free Bun address; conventions. |
| D3 | User-scoped requests carry `household_id` (and `member_id` where needed) as message fields. The gateway always fills them from the JWT, never from client input. Services trust them (internal network) and service-to-service calls pass them on. | Matches the existing contracts. |
| D4 | New metadata key `x-language` (BCP 47), set by the gateway from `Accept-Language`, forwarded like the others. | Catalogue, recipes and AI output are per language. |
| D5 | Each service enforces its own Plus rules, reading the household's plan from identity (`GetEntitlements`) on each request that needs it and per household per run in background jobs. No local plan copy (Q004). The gateway still forwards `x-plan`. The rules are in "Free and Plus" below. | One source of truth for the plan; changed from a local copy by Q004. |
| D6 | Photos (fridge, shelf, receipt, label, user flyer) are processed in memory and never stored. Chain flyer pages are stored until their offers expire plus 90 days, for review and disputes (Q010). | "Collect the minimum" (business model). |
| D7 | Voice input = the phone's speech recognition, then the text intake. | Claude takes no audio. |
| D8 | Savings counters live in kitchen (stock knows what was used before expiry and what it cost); release R2 for counters and report, R3 for the share card. | Not placed in features.md section 7. |
| D9 | GDPR data export is assembled by the gateway, calling each service's export RPC. Deletion goes through the `account.deleted` / `household.deleted` events. | Export is read-only; deletion must be durable. |
| D10 | `ai` has no database (redis-cache for results, redis-state for spend counters, decision 17). For batch jobs it returns the batch id; the calling service stores it and polls. | Keeps `ai` free of a database. |
| D11 | `ai` has a global daily spend limit; above it, per-request calls return `RESOURCE_EXHAUSTED` and callers fall back to non-AI results. | Cost guardrails in the business model. |
| D12 | Release placement where features.md section 7 is silent: swaps and free-text wishes R1 (part of the day plan), tool rewrite R2, adding well-rated AI dishes to the library R2, offer-ends reminders R2. | Closest release by dependency. |
| D13 | Distances are never computed on the server. Shop lists come from the postcode area with coordinates; the phone sorts, filters by range, and sends shop ids nearest first when order matters. | features.md privacy rule. |

## Conventions for every task

- **Contracts.** One public gRPC service per back-end service (`KitchenService`, ...) and one admin service
  (`KitchenAdminService` in `contracts/coolinary/kitchen/v1/admin.proto`). Messages may be split into per-domain files
  in the same package (`catalogue.proto`, `stock.proto`). Sketches in tasks use proto3 and the shared types in
  `coolinary/common/v1/common.proto` (`Money`, `Quantity`, `Unit`, paging). New shared types go there.
- **Migrations.** Sketches name the file `000NN_<topic>.sql` (goose) or describe the Drizzle schema change; the number is
  given at implementation.
- **Errors.** gRPC codes per `docs/service-conventions.md`; add `FAILED_PRECONDITION` for a broken business rule
  (e.g. "Plus required"), `RESOURCE_EXHAUSTED` for quotas and the AI spend limit, `ALREADY_EXISTS`,
  `PERMISSION_DENIED` (wrong household or role), `ABORTED` for version conflicts. The reason travels in trailing
  metadata `x-error-code` (decision 18). The gateway maps both to HTTP (GW-01).
- **Shared types.** Enums used by more than one service live in `coolinary/common/v1` (decision 16).
- **AI output is a proposal.** Every AI result is checked without AI before it is shown or stored, and every AI
  feature has a non-AI fallback when `ai` is down, over budget or refuses. AI never states a calorie number.
- **Privacy.** Age, height, weight and goals never leave the phone; only each person's daily kcal target does. No
  exact location on the server, only a postcode. No personal data or user free text in logs or traces.
- **Languages.** Launch languages are Polish, German and English (wave 1). Every user-facing text the server produces
  (e-mails, push, dish names, steps) exists per language.
- **Admin RPCs.** Called only by admin, which sets `x-admin-user-id` metadata; the owning service writes it into its
  own change records where the task says so.

## Free and Plus (from features.md section 6)

| Rule | Free | Plus | Enforced in |
|---|---|---|---|
| AI day plans | 2 per week | unlimited | gateway (GW-03) |
| AI photo scans (fridge, shelf, label, receipt) | 2 per week | unlimited | gateway (GW-03) |
| Other AI (swaps, rewrites, week plan, family plan, deals, top-up, own flyer) | not available | unlimited | gateway + owning service |
| Expiry reminders | no | yes | kitchen (KIT-09) |
| Calorie goals | one personal goal | every member, children included | kitchen (KIT-11) |
| Per-person diet profiles | no (one household profile); on a downgrade every person's allergens join the household list (Q195) | yes | kitchen (KIT-11) |
| Shops with offers | 2 shops | all in range | offers (OFF-04) |
| Trip plan across shops (price matrix) | no | yes | planner (PLN-19) |
| Shared household (up to 5 accounts) | no new sharing; after a downgrade existing members stay, but nobody new can be invited or join (Q187) | yes, held by the owner | identity (IDN-07) |
| Near-shop reminders, savings report | no | yes | phone / kitchen (KIT-12) |

Anonymous accounts count as Free households.

## Epic map

### gateway (Bun + Elysia) — `gateway/`

| Epic | Title | Rel. |
|---|---|---|
| GW-01 | Platform: request context, error model, idempotency keys, OpenAPI rules | R1 |
| GW-02 | Auth and account routes (anonymous, e-mail code, me, entitlements, deletion, export) | R1 |
| GW-03 | Free-tier AI quota and rate limits | R1 |
| GW-04 | Household, profiles and invite routes | R1–R2 |
| GW-05 | Kitchen profile and catalogue routes | R1–R2 |
| GW-06 | Fridge, barcode and intake routes | R1–R3 |
| GW-07 | Recipe library, plans and cooking routes | R1–R3 |
| GW-08 | Shopping list routes | R1 |
| GW-09 | Live shopping list (WebSocket) | R2 |
| GW-10 | Shops, offers, flyers and price matrix routes | R2–R3 |
| GW-11 | Devices and notification preference routes | R1 |
| GW-12 | Savings routes | R2–R3 |

### identity (Bun + Elysia + Better Auth) — `identity/`

| Epic | Title | Rel. |
|---|---|---|
| IDN-01 | Events: outbox and Redis Streams publisher (Bun) | R1 |
| IDN-02 | Sign-in with e-mail code; remove passwords | R1 |
| IDN-03 | Anonymous accounts and upgrade | R1 |
| IDN-04 | Households and members core | R1 |
| IDN-05 | Plan and entitlements | R1 |
| IDN-06 | Profiles without accounts | R2 |
| IDN-07 | Invites (QR, code, link) | R2 |
| IDN-08 | Leaving and removing members | R2 |
| IDN-09 | Account deletion and data export | R1 |
| IDN-10 | Admin RPCs and admin role | R1 |

### notifier (Bun) — `notifier/`

| Epic | Title | Rel. |
|---|---|---|
| NTF-01 | Events consumer (Bun) | R1 |
| NTF-02 | E-mail templates and languages | R1 |
| NTF-03 | Push delivery (Expo Push API) | R1 |
| NTF-04 | Devices and notification preferences | R1 |
| NTF-05 | Founder alerts | R1 |
| NTF-06 | Admin RPCs (delivery log) | R1 |

### admin (Bun + Elysia, new) — `admin/`

| Epic | Title | Rel. |
|---|---|---|
| ADM-01 | Service skeleton (compose, Caddy, database, monitoring) | R1 |
| ADM-02 | Admin sign-in, roles and audit log | R1 |
| ADM-03 | Catalogue, matching and nutrition review | R1 |
| ADM-04 | Recipe import and dish review | R1 |
| ADM-05 | Households and plans support | R1 |
| ADM-06 | Flyer sources, offers and shops review | R2 |
| ADM-07 | AI cost and system overview | R1 |

### kitchen (Go) — `kitchen/`

| Epic | Title | Rel. |
|---|---|---|
| KIT-01 | Platform: events outbox, publisher and consumer (Go), plan lookup, admin service | R1 |
| KIT-02 | Ingredient catalogue and translations | R1 |
| KIT-03 | Name matching (trigram + Haiku) | R1 |
| KIT-04 | Barcode lookup chain | R1 |
| KIT-05 | Unknown products and shared-catalogue confirmation | R1 |
| KIT-06 | Nutrition tables and calculation | R1 |
| KIT-07 | Fridge stock | R1 |
| KIT-08 | AI intake: photo, text, receipt | R1–R3 |
| KIT-09 | Expiry reminders | R1 |
| KIT-10 | Kitchen profile and household settings | R1 |
| KIT-11 | Per-person profiles and calorie targets | R2 |
| KIT-12 | Savings ledger and monthly report | R2–R3 |
| KIT-13 | Household lifecycle: event consumers and data export | R1–R2 |

### planner (Go) — `planner/`

| Epic | Title | Rel. |
|---|---|---|
| PLN-01 | Platform: events (Go), identity/kitchen/offers/ai clients, profile and stock snapshots, plan lookup | R1 |
| PLN-02 | Recipe library model and storage | R1 |
| PLN-03 | Library import, AI drafting and review | R1 |
| PLN-04 | Dish nutrition and calories | R1 |
| PLN-05 | Library browse, favourites, ratings, history | R1 |
| PLN-06 | What can I make now | R1 |
| PLN-07 | Rule checker for AI and library dishes | R1 |
| PLN-08 | AI day plan and cuisine choice | R1 |
| PLN-09 | Swaps, wishes, tool rewrite, restyle | R1–R3 |
| PLN-10 | Cooking mode and "cooked it" | R1 |
| PLN-11 | Library growth from well-rated AI dishes | R2 |
| PLN-12 | Per-person portions (one dish, one pot) | R2 |
| PLN-13 | Family day plan | R2 |
| PLN-14 | Week plan with leftovers and fewest cooking sessions | R3 |
| PLN-15 | Shopping list | R1 |
| PLN-16 | Cuisine base pantry | R1 |
| PLN-17 | Offer badges and offer-ends reminders | R2 |
| PLN-18 | Plan from deals | R2 |
| PLN-19 | Trip price matrix | R2–R3 |
| PLN-20 | Cheapest top-up | R3 |
| PLN-21 | Household lifecycle: event consumers and data export | R1–R2 |

### offers (Go) — `offers/`

| Epic | Title | Rel. |
|---|---|---|
| OFF-01 | Platform: events (Go), identity and kitchen clients, plan lookup, admin service | R2 |
| OFF-02 | Chains and shops from OpenStreetMap | R2 |
| OFF-03 | Postcode areas and flyer editions | R2 |
| OFF-04 | My shops and shopping preferences | R2 |
| OFF-05 | Flyer source framework and sync scheduler | R2 |
| OFF-06 | Launch chain adapters (Biedronka, Lidl PL, Kaufland PL, Lidl DE, Aldi Süd, Aldi Nord, Kaufland DE) | R2 |
| OFF-07 | Offer extraction, validation and publishing | R2 |
| OFF-08 | Accuracy gate and weekly spot-checks | R2 |
| OFF-09 | Offer browse, search and lookup by ingredient | R2 |
| OFF-10 | User-supplied flyers | R2 |
| OFF-11 | Prices: flyer, receipt and reference | R2–R3 |
| OFF-12 | Next chain adapters (Dino, Netto PL, Carrefour PL, Auchan PL, Aldi PL, Penny DE) | R3 |
| OFF-13 | Licensed offers feed (ShopFully) | R3 |
| OFF-14 | Household lifecycle: event consumers and data export | R2 |

### ai (Python + FastAPI) — `ai/`

| Epic | Title | Rel. |
|---|---|---|
| AI-01 | Claude provider foundation (SDK, structured output, caching, refusals, errors) | R1 |
| AI-02 | Cost tracking and spend limit | R1 |
| AI-03 | Prompts, golden sets and evals | R1 |
| AI-04 | Batch jobs (Message Batches) | R1 |
| AI-05 | Text to items | R1 |
| AI-06 | Photo to items | R1 |
| AI-07 | Product label reading | R1 |
| AI-08 | Catalogue match pick | R1 |
| AI-09 | Catalogue and library batch jobs (translation, nutrition matches, effort tags, dish drafts) | R1 |
| AI-10 | Day plan composition and cuisine choices | R1 |
| AI-11 | Swaps, wishes, tool rewrite, restyle | R1–R3 |
| AI-12 | Family day plan | R2 |
| AI-13 | Flyer offer extraction | R2 |
| AI-14 | Plan from deals | R2 |
| AI-15 | Week plan | R3 |
| AI-16 | Receipt reading | R3 |
| AI-17 | Cheapest top-up | R3 |

### mobile communication layer (Expo) — `mobile/`

| Epic | Title | Rel. |
|---|---|---|
| MOB-01 | Auth session (anonymous start, e-mail code, JWT refresh) | R1 |
| MOB-02 | API client and error model | R1 |
| MOB-03 | Query keys, hooks and cache rules | R1–R3 |
| MOB-04 | Offline mutation queue | R1 |
| MOB-05 | Live shopping list client | R2 |
| MOB-06 | Push notifications | R1 |
| MOB-07 | Uploads (photos, PDFs) | R1 |

## Cross-service calls

Call graph (no cycles): `gateway → identity, notifier, kitchen, planner, offers`; `identity → notifier`;
`kitchen → identity, ai, notifier`; `planner → identity, kitchen, offers, ai, notifier`; `offers → identity, kitchen, ai, notifier`;
`admin → identity, notifier, kitchen, planner, offers, ai`. `ai` and `notifier` call no one.

RPCs that other services depend on (names are binding; message fields are designed in the owning task). Each
service README lists all of its RPCs with the task that defines them:

| Owner | RPC | Epic | Called by |
|---|---|---|---|
| identity | `GetHousehold` (exists), `GetEntitlements` | IDN-04, IDN-05 | gateway |
| identity | member, profile, invite, leave/remove RPCs | IDN-04, 06–08 | gateway |
| identity | `IdentityAdminService.IsAdmin` | IDN-10 | admin |
| notifier | `SendEmail` (extended with template + language) | NTF-02 | identity |
| notifier | `SendPush` (to users or a household, with a category) | NTF-03 | kitchen, planner |
| notifier | `SendAlert` (founder alert with a dedup key) | NTF-05 | identity, kitchen, planner, offers |
| notifier | `RegisterDevice`, `UnregisterDevice`, `GetPreferences`, `SetPreferences` | NTF-04 | gateway |
| kitchen | `ResolveBarcode` (exists, extended) | KIT-04 | gateway |
| kitchen | `SearchIngredients`, `GetIngredients`, `GetCatalogueVersion` | KIT-02 | gateway, planner, offers |
| kitchen | `MatchNames` | KIT-03 | offers, planner |
| kitchen | `ComputeNutrition` | KIT-06 | planner |
| kitchen | `ListStock`, `AddStockItems`, `UpdateStockItem`, `RemoveStockItems`, `ConsumeStock` | KIT-07 | gateway, planner |
| kitchen | `GetPlanningProfile` (one snapshot of everything planning needs) | KIT-10, KIT-11 | planner |
| planner | `SuggestDishes` (exists), library, plan, list RPCs | PLN-* | gateway |
| offers | `ListMyShops` | OFF-04 | planner |
| offers | `FindOffersForIngredients`, `GetOffers` | OFF-09 | planner |
| offers | `GetPrices` | OFF-11 | planner |
| ai | `ExtractItemsFromText` (exists), `ExtractItemsFromPhoto`, `ReadProductLabel`, `ExtractReceipt` | AI-05, 06, 07, 16 | kitchen |
| ai | `PickCatalogueMatches` | AI-08 | kitchen |
| ai | `SubmitBatch`, `GetBatch`, `GetBatchResults` (server-streaming), `CancelBatch` (job kinds below) | AI-04 | kitchen, planner, offers |
| ai | `ComposeDayPlan`, `ModifyMeal`, `RewriteDish`, `ComposeFamilyDay`, `ComposeFromDeals`, `ComposeWeekPlan`, `CheapestTopUp` | AI-10–17 | planner |
| ai | `ExtractFlyerOffers` (per request) | AI-13 | offers |
| every service | `<Svc>AdminService` RPCs | each epic's admin task | admin |
| kitchen, planner, offers | `ExportHouseholdData` | KIT-13, PLN-21, OFF-14 | gateway |
| identity, notifier | `ExportUserData` | IDN-09, NTF-04 | gateway |
| offers | `RecordFeedUsage` (R3) | OFF-13 | gateway |

Batch job kinds (AI-04, AI-09, AI-13): `translate_ingredients`, `match_nutrition_entries`, `pick_catalogue_matches`,
`draft_dishes`, `translate_dishes`, `suggest_effort_tags`, `extract_flyer_offers`.

## Events (Redis Streams on redis-state)

Streams: `events.identity`, `events.kitchen`, `events.offers` (planner publishes no events). Each stream entry has
two fields: `type` (the event type, for metrics and skipping) and `envelope` (JSON) (Q015). Envelope: `id` (UUIDv7, the idempotency key), `type`, `version`,
`occurred_at` (RFC 3339 UTC), `producer`, `household_id` (empty for catalogue-wide events), `payload`. Field names
are snake_case. No personal data beyond ids. Streams are trimmed to about 7 days (`MINID ~`); a consumer that falls
further behind asks the owner to republish snapshots (`household.snapshot`, IDN-01-04). Adding an optional field
keeps the version; anything else raises it. Event JSON Schemas live in
`contracts/events/<producer>/<type>.v1.json` (Q014).

| Event | Producer (epic) | Payload | Consumers |
|---|---|---|---|
| `household.created` | identity (IDN-01-03) | `household_id`, `owner_user_id`, `owner_member_id` | kitchen, planner, notifier |
| `household.member_added` | identity | `household_id`, `member_id`, `user_id` or null, `has_account`, `role` (`owner`/`member`/`profile`), `brought_fridge_from` or null, `previous_member_id` or null | kitchen, planner, notifier |
| `household.member_removed` | identity | `household_id`, `member_id`, `user_id` or null, `reason` (`left`/`removed`/`profile_deleted`; `plan_ended` is reserved and not emitted, Q187), `new_household_id` or null | kitchen, planner, notifier |
| `household.plan_changed` | identity | `household_id`, `plan`, `previous_plan`, `source` (`default`/`manual`/`store`) | kitchen, planner, offers (to react to up- and downgrades; the plan itself is read from identity, Q004) |
| `household.deleted` | identity | `household_id`, `reason` (`owner_deleted`/`joined_other`/`linked_merge`/`anonymous_cleanup`) | kitchen, planner, offers, notifier |
| `household.snapshot` | identity (IDN-01-04) | `household_id`, `plan`, `members` (`member_id`, `user_id` or null, `has_account`, `role`) | kitchen, planner, offers, notifier |
| `account.linked` | identity (IDN-03) | `anonymous_user_id`, `user_id`, `household_id`, `anonymous_household_id`, `merged` | kitchen, planner, notifier |
| `account.deleted` | identity (IDN-09) | `user_id`, `member_ids` | kitchen, planner, offers, notifier |
| `catalogue.ingredient_merged` | kitchen (KIT-02-09) | `from_ingredient_id`, `into_ingredient_id` | planner, offers |
| `catalogue.ingredient_flags_changed` | kitchen (KIT-02-10, Q054) | `ingredient_ids` | planner |
| `catalogue.match_changed` | kitchen (KIT-03-06, Q058) | normalised name, language, old and new ingredient id | offers, planner |
| `catalogue.nutrition_changed` | kitchen (KIT-06) | `ingredient_ids` (≤ 500 per event) | planner |
| `profile.changed` | kitchen (KIT-10/11) | `household_id`, `sections`, `planning_version` | planner |
| `stock.changed` | kitchen (KIT-07-08) | `household_id`, `stock_version` | planner (cache) |
| `receipt.scanned` | kitchen (KIT-08-06, R3) | no household id; `chain_key`, `shop_name`, `country`, `currency`, `purchased_on`, lines with `ingredient_id`, `ean`, `unit_price`, `quantity`, `item_quantity` | offers |
| `offers.published` | offers (OFF-07) | `offer_ids`, `chain_id`, `valid_from`, `valid_to`, optional `ingredient_ids` | planner |
| `offers.expired` | offers (OFF-07) | `offer_ids` | planner |

**Live notices** (not streams): Redis pub/sub on redis-state, channel `live.household.<id>`, JSON
`{"kind": "...", "household_id": "...", "version": n}`. Planner publishes `list.changed` with the list version
(PLN-15); kitchen publishes `stock.changed` with the stock version (KIT-07-08). The gateway relays them over the
WebSocket (GW-09); clients refetch when a version is skipped.

## Existing docs and code to update

- Done: `docs/infrastructure.md` (recipe library in planner, admin service, call graph, events, Claude, no pgvector)
  and `docs/service-conventions.md` (metadata, admin services, shared types, error codes and `x-error-code`, events).
- With ADM-01: `platform/postgres/initdb/01-databases.sh`, the compose files, the Caddyfile and Prometheus.
- With the first contract task of each service: `contracts/coolinary/common/v1/` gets the shared enums (decision 16).
