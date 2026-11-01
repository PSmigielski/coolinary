# Coolinary — decisions

Every small choice the task files left open, settled on 1 November 2026. The task files state each answer as
`Decided (Qnnn):` and are written to match it, so this file is the index, not the source of the rules.

- **How they were settled:** Q001–Q028 were answered one by one; the remaining 200 were triaged. Critical items
  (child and allergen safety) and high items (legal and licences, data that can't be recovered, Free vs Plus, launch
  gates, visible household rules) were answered by you; medium and low items take their recommended answer.
- **Status:** `you chose` = your answer differs from the recommendation · `you confirmed` = you picked the
  recommendation · `default` = recommended answer, not asked. Any of them can still be changed: change the task's
  `Decided` rule and this line together.
- **Severity** of the items that were triaged: 5 critical, 42 high, 78 medium, 78 low; 28 were asked before triage.

231 decisions: 12 you chose, 63 you confirmed, 156 default.

## Where you chose differently from the recommendation

| ID | Decision | Your answer | Instead of | Tasks |
|---|---|---|---|---|
| Q004 | Should each service enforce its own Plus rules, using x-plan and a local plan copy kept from events? | **Ask identity each time — kitchen, planner and offers read the plan from identity's GetEntitlements on each request that needs it, and per household in background jobs; no local plan copy.** | Yes. | [README](README.md) |
| Q010 | Never store fridge, receipt, label or user-flyer photos; keep chain flyer pages only until offers expire + 14 days? | **Flyer pages kept 90 days — User photos are still never stored; chain flyer pages are kept until their offers expire plus 90 days.** | Yes. | [README](README.md) |
| Q016 | How long should planner keep its log of admin changes to dishes? | **1 year — Planner keeps its admin change log for 1 year.** | 2 years. | [PLN-01-06](planner/E01-platform/PLN-01-06-admin-service.md) |
| Q021 | How long may a phone stay offline and still have its queued edits replayed safely: 72 hours or 24 hours? | **24 hours — Idempotency keys replay for 24 hours; the app drops queued offline edits older than 24 hours.** | 72 hours, so a weekend offline loses nothing | [GW-01-03](gateway/E01-platform/GW-01-03-idempotency-keys.md) |
| Q022 | Should a flyer PDF a user uploads stay capped at 20 pages (25 MB), or go up to 40 pages? | **40 pages — Own-flyer PDFs may have up to 40 pages and 25 MB, from R2.** | Keep 20 pages for R2 and measure how often users hit it | [GW-01-05](gateway/E01-platform/GW-01-05-uploads.md), [OFF-10-01](offers/E10-user-flyers/OFF-10-01-upload-and-jobs.md), [MOB-07-03](mobile/E07-uploads/MOB-07-03-flyer-inputs.md) |
| Q076 | May a household opt out of having its confirmed product names sent to Open Food Facts? | **Household opt-out — A household setting, on by default, keeps its confirmed product names out of what is sent to Open Food Facts.** | No opt-out; only product facts confirmed by two households are sent. | [KIT-05-04](kitchen/E05-unknown-products/KIT-05-04-contribute-to-off.md) |
| Q103 | Should items past their use-by date count as "at home" when suggesting what to cook? | **Count with a warning — Items past their date count as at home, always marked "check it first" in suggestions, AI plans and cooking mode.** | No: expired items don't count. | [PLN-06-01](planner/E06-what-can-i-make-now/PLN-06-01-stock-coverage.md) |
| Q111 | When does Free's weekly AI allowance (2 day plans, 2 scans) reset: every Monday 00:00 UTC, or on a rolling 7 days? | **Rolling 7 days — Each use frees up again 7 days after it was made; the app shows when the next slot frees.** | Every Monday 00:00 UTC (01:00–02:00 in Poland and Germany) | [GW-03-01](gateway/E03-free-tier-quota/GW-03-01-quota-counters.md) |
| Q137 | When the owner of a shared household deletes their account, should ownership pass automatically to the longest-standing member? | **Owner chooses first — An owner of a shared household must name the new owner before deleting their account; with no other account member, the household is deleted with the account.** | Yes, automatically to the longest-standing member | [IDN-09-02](identity/E09-account-deletion-and-export/IDN-09-02-delete-account.md), [GW-02-05](gateway/E02-auth-and-account/GW-02-05-delete-account.md) |
| Q158 | If a chain's weekly spot-check scores below 97% two weeks in a row (or once below 90%), should it be taken offline automatically? | **Alert only — Bad weekly spot-checks only alert the founder; a chain is never taken offline automatically.** | Yes: moved back to syncing automatically, and the founder is alerted. | [OFF-08-03](offers/E08-accuracy-gate/OFF-08-03-weekly-spot-checks.md) |
| Q187 | When a shared household drops from Plus to Free, what happens to the other members? | **Members stay, no new invites — After a downgrade, existing members keep sharing; nobody new can be invited or join until Plus returns. This relaxes features.md's "sharing needs Plus".** | 14 days of grace, then non-owner members move to their own households. | [IDN-05-04](identity/E05-plan-and-entitlements/IDN-05-04-shared-household-downgrade.md) |
| Q195 | After a household drops from Plus to Free, should each person's allergens and diets still apply to plans? | **Allergens join household list — On a downgrade every person's allergens are added to the household allergen filter; per-person diets stop applying; editing per-person data waits for Plus.** | Yes, keep applying them; editing waits until the household is Plus again. | [KIT-11-04](kitchen/E11-per-person-profiles/KIT-11-04-per-person-restrictions.md) |

## All decisions by milestone

### M1 Foundations (R1)

| ID | Severity | Decision | Answer | Status | Tasks |
|---|---|---|---|---|---|
| Q001 (D10) | — | Should ai stay without a database, returning batch ids for callers to store and poll? | Yes (redis-cache and redis-state only). | you confirmed | [README](README.md) |
| Q002 (D11) | — | Should ai stop per-request AI calls above a global daily spend limit, so callers fall back to non-AI results? | Yes. | you confirmed | [README](README.md) |
| Q003 (D2) | — | Give admin IP 10.77.0.23, HTTP 3000, admin port 8081 and its own `admin` database for the audit log? | Yes, as written. | you confirmed | [README](README.md) |
| Q004 (D5) | — | Should each service enforce its own Plus rules, using x-plan and a local plan copy kept from events? | Ask identity each time — kitchen, planner and offers read the plan from identity's GetEntitlements on each request that needs it, and per household in background jobs; no local plan copy. | **you chose** | [README](README.md) |
| Q005 (D7) | — | Should voice input be the phone's own speech recognition feeding the text intake? | Yes; Claude takes no audio. | you confirmed | [README](README.md) |
| Q006 (D8) | — | Should savings counters live in kitchen, with counters and report in R2 and the share card in R3? | Yes. | you confirmed | [README](README.md) |
| Q007 (D1) | — | Add the calls offers → kitchen (name matching), planner → notifier (offer-ends reminders) and admin → every service? | Yes; the call graph stays free of cycles. | you confirmed | [README](README.md) |
| Q008 (D3) | — | Should user requests carry household_id as a message field that the gateway fills from the JWT and services trust? | Yes; never taken from client input. | you confirmed | [README](README.md) |
| Q009 (D4) | — | Add an `x-language` metadata key set by the gateway from Accept-Language? | Yes, BCP 47, forwarded like x-plan. | you confirmed | [README](README.md) |
| Q010 (D6) | — | Never store fridge, receipt, label or user-flyer photos; keep chain flyer pages only until offers expire + 14 days? | Flyer pages kept 90 days — User photos are still never stored; chain flyer pages are kept until their offers expire plus 90 days. | **you chose** | [README](README.md) |
| Q011 (D9) | — | Should the gateway assemble the GDPR export by calling each service, and deletion travel as events? | Yes. | you confirmed | [README](README.md) |
| Q012 (D12) | — | Release placement: swaps and wishes R1, tool rewrite R2, AI dishes into the library R2, offer-ends reminders R2? | Yes, as listed. | you confirmed | [README](README.md) |
| Q013 (D13) | — | Never compute distances on the server: the phone sorts shops and sends ids nearest first? | Yes. | you confirmed | [README](README.md) |
| Q014 | — | Where should the JSON Schemas of event payloads live? | contracts/events/<producer>/<type>.v1.json, next to the gRPC contracts. | you confirmed | [IDN-01-03](identity/E01-events-outbox/IDN-01-03-event-types.md), [README](README.md) |
| Q015 | — | Should each Redis Streams entry carry two fields, the event type and the full JSON envelope? | Two fields: type and envelope. | you confirmed | [KIT-01-02](kitchen/E01-platform/KIT-01-02-event-envelope-outbox.md), [IDN-01-02](identity/E01-events-outbox/IDN-01-02-publisher.md), [README](README.md) |
| Q016 | — | How long should planner keep its log of admin changes to dishes? | 1 year — Planner keeps its admin change log for 1 year. | **you chose** | [PLN-01-06](planner/E01-platform/PLN-01-06-admin-service.md) |
| Q017 | — | When Claude Sonnet wrongly declines a recipe request, should Anthropic automatically retry on another model, which may cost more? | On, with the answering model's real price recorded. | you confirmed | [AI-01-03](ai/E01-claude-provider/AI-01-03-claude-provider.md) |
| Q018 | — | What hard AI spend limits should stop AI calls: €20 a day overall, Free at most 40% of it, and €60 a month for batch jobs? | €20/day, Free ≤ 40% of it, €60/month for batches. | you confirmed | [AI-02-04](ai/E02-cost-and-spend-limit/AI-02-04-daily-spend-limit.md) |
| Q019 | — | Should a Free household's AI calls stop for the month once they cost €0.15? | Yes: hard guard at €0.15 per Free household per month. | you confirmed | [AI-02-05](ai/E02-cost-and-spend-limit/AI-02-05-household-cost.md) |
| Q020 | — | Should the app send its version and platform in an x-client header on every request? | Yes, for support and for refusing too-old app versions | you confirmed | [MOB-02-01](mobile/E02-api-client/MOB-02-01-request-middleware.md), [GW-01-01](gateway/E01-platform/GW-01-01-request-context.md) |
| Q021 | — | How long may a phone stay offline and still have its queued edits replayed safely: 72 hours or 24 hours? | 24 hours — Idempotency keys replay for 24 hours; the app drops queued offline edits older than 24 hours. | **you chose** | [GW-01-03](gateway/E01-platform/GW-01-03-idempotency-keys.md) |
| Q022 | — | Should a flyer PDF a user uploads stay capped at 20 pages (25 MB), or go up to 40 pages? | 40 pages — Own-flyer PDFs may have up to 40 pages and 25 MB, from R2. | **you chose** | [GW-01-05](gateway/E01-platform/GW-01-05-uploads.md), [OFF-10-01](offers/E10-user-flyers/OFF-10-01-upload-and-jobs.md), [MOB-07-03](mobile/E07-uploads/MOB-07-03-flyer-inputs.md) |

### M2 Sign-in, households and plans (R1)

| ID | Severity | Decision | Answer | Status | Tasks |
|---|---|---|---|---|---|
| Q023 | — | Should a new household's name be stored empty and shown as a translated default in the app? | Store empty; app shows "Mój dom" / "Mein Zuhause" / "My home". | you confirmed | [IDN-04-02](identity/E04-households-core/IDN-04-02-household-created-events.md) |
| Q024 | — | Who may rename the household? | The owner only. | you confirmed | [IDN-04-05](identity/E04-households-core/IDN-04-05-rename.md) |
| Q025 | — | Should the gateway read Free-tier limits (2 plans, 2 scans a week) from identity instead of hard-coding them? | Yes: from GetEntitlements, cached 5 minutes. | you confirmed | [IDN-05-02](identity/E05-plan-and-entitlements/IDN-05-02-get-entitlements.md) |
| Q026 | — | Must someone sign in with an e-mail code before buying Plus? | Yes, sign in first. | you confirmed | [IDN-05-05](identity/E05-plan-and-entitlements/IDN-05-05-store-subscription-link.md) |
| Q027 | — | Should admin rights be stored in identity's own table and granted only from a command line? | Own admin_users table plus a CLI command, not Better Auth's admin plugin. | you confirmed | [IDN-10-01](identity/E10-admin/IDN-10-01-admin-role.md) |
| Q028 | — | Should the back office have a separate dietitian role that can only read dishes and record checks? | Two roles: admin (everything) and dietitian. | you confirmed | [ADM-02-02](admin/E02-auth-and-audit/ADM-02-02-admin-check.md), [IDN-10-01](identity/E10-admin/IDN-10-01-admin-role.md) |
| Q029 | medium | How long should the e-mail and push delivery log be kept? | 90 days. | default | [NTF-02-02](notifier/E02-email-templates/NTF-02-02-delivery-log-privacy.md) |
| Q030 | low | Should notifier cap e-mails per address (5 per 15 minutes, 20 per day) on top of identity's limits? | Yes, keep the extra cap. | default | [NTF-02-03](notifier/E02-email-templates/NTF-02-03-send-email-templated.md) |
| Q031 | low | How long should an e-mailed sign-in code stay valid? | 10 minutes, to allow for slow e-mail; the app countdown follows | default | [IDN-02-02](identity/E02-email-code-sign-in/IDN-02-02-email-code-plugin.md), [MOB-01-05](mobile/E01-auth-session/MOB-01-05-email-code-sign-in.md) |
| Q032 | low | Should the app make users wait 60 seconds before requesting another sign-in code? | Yes, 60 s, matching identity's limit of one code per minute | default | [IDN-02-03](identity/E02-email-code-sign-in/IDN-02-03-sign-in-rate-limits.md), [MOB-01-05](mobile/E01-auth-session/MOB-01-05-email-code-sign-in.md) |
| Q033 | medium | How long should a phone stay signed in without being used? | 60 days, extended on every daily use. | default | [IDN-02-05](identity/E02-email-code-sign-in/IDN-02-05-session-lifetime.md) |
| Q034 | medium | How many anonymous accounts may one IP address create per hour, given each gets its own Free AI allowance? | 5 per hour (identity's draft), with the gateway using the same number | default | [IDN-03-01](identity/E03-anonymous-accounts/IDN-03-01-anonymous-sign-in.md), [GW-03-04](gateway/E03-free-tier-quota/GW-03-04-rate-limits.md) |
| Q035 | high | Must someone have a registered account (with e-mail) before joining a shared household? | Yes, only registered accounts can join | you confirmed | [GW-04-04](gateway/E04-household-and-invites/GW-04-04-invite-preview-accept.md), [IDN-03-01](identity/E03-anonymous-accounts/IDN-03-01-anonymous-sign-in.md) |
| Q036 | high | When someone using the app anonymously signs in to an account they already have, what happens to the anonymous fridge? | Move the fridge items into the account; drop plans, list and profile. | you confirmed | [IDN-03-03](identity/E03-anonymous-accounts/IDN-03-03-upgrade-existing-account.md) |
| Q037 | medium | Should the sign-in code appear in the e-mail subject, so it shows in the phone notification? | Yes, code in the subject. | default | [NTF-02-04](notifier/E02-email-templates/NTF-02-04-sign-in-code-email.md) |
| Q038 | medium | Should German texts address users as "du" or "Sie"? | "du". | default | [NTF-02-04](notifier/E02-email-templates/NTF-02-04-sign-in-code-email.md) |
| Q039 | low | Should identity send an e-mail confirming that an account was deleted? | Yes. | default | [NTF-02-05](notifier/E02-email-templates/NTF-02-05-account-deleted-email.md) |
| Q040 | low | How often should a repeating founder alert be e-mailed again? | Critical every 1 hour, warning every 6 hours. | default | [NTF-05-01](notifier/E05-founder-alerts/NTF-05-01-send-alert.md) |
| Q041 | low | Should every service use the same six metric names for publishing and consuming events? | Yes: six shared names, written into service-conventions.md. | default | [IDN-01-02](identity/E01-events-outbox/IDN-01-02-publisher.md) |
| Q042 | medium | If notifier misses events because it was down longer than 7 days, should it rebuild membership automatically? | No: accept the risk and rely on the alert. | default | [NTF-01-01](notifier/E01-events-consumer/NTF-01-01-streams-consumer.md), [KIT-01-05](kitchen/E01-platform/KIT-01-05-dead-letters-lag.md) |
| Q043 | medium | Should the anonymous account be created only when the app first needs the server, rather than at first start? | Only on the first request that needs an account | default | [MOB-01-03](mobile/E01-auth-session/MOB-01-03-anonymous-start.md) |
| Q044 | high | How long should an unused anonymous account and its fridge be kept before it is deleted? | 30 days without use (identity's draft) | you confirmed | [IDN-03-04](identity/E03-anonymous-accounts/IDN-03-04-stale-anonymous-cleanup.md), [MOB-01-03](mobile/E01-auth-session/MOB-01-03-anonymous-start.md) |
| Q045 | high | Should members' daily calorie targets be kept out of the phone's unencrypted offline cache? | Yes; read them from secure storage or refetch | you confirmed | [MOB-03-01](mobile/E03-query-hooks/MOB-03-01-query-keys-cache-rules.md) |
| Q229 | high | If kitchen, planner or offers can't reach identity for a household's plan, what do they do? | Fall back to the plan in the signed-in user's token (x-plan, ≤ 15 min old), Free if missing; jobs skip the household until their next run. | you confirmed | [KIT-01-06](kitchen/E01-platform/KIT-01-06-household-plan-copy.md), [PLN-01-05](planner/E01-platform/PLN-01-05-household-plan-copy.md), [OFF-01-03](offers/E01-platform/OFF-01-03-household-plan-copy.md), [OFF-04-04](offers/E04-my-shops/OFF-04-04-plan-change-rules.md) |

### M3 Back office (R1)

| ID | Severity | Decision | Answer | Status | Tasks |
|---|---|---|---|---|---|
| Q046 | medium | Should the back office sign in through its own proxy rather than calling the public API with CORS? | Proxy through the admin service. | default | [ADM-02-01](admin/E02-auth-and-audit/ADM-02-01-sign-in-proxy.md) |
| Q047 | medium | How long should the back-office audit log be kept? | 2 years. | default | [ADM-02-03](admin/E02-auth-and-audit/ADM-02-03-audit-log.md) |
| Q048 | medium | Should admins be able to grant admin rights from the back office? | No; only identity's command line grants them. | default | [ADM-02-04](admin/E02-auth-and-audit/ADM-02-04-admin-accounts.md) |
| Q049 | medium | Until payments exist, what monthly net revenue per Plus household should the 15% AI cost check use? | €1.64 per month, as a config value. | default | [ADM-07-02](admin/E07-ai-cost-and-overview/ADM-07-02-ai-cost-view.md) |

### M4 Catalogue and kitchen profile (R1)

| ID | Severity | Decision | Answer | Status | Tasks |
|---|---|---|---|---|---|
| Q050 | medium | When an ingredient measured in millilitres has no known density, should we assume 1 g per ml? | No default; the dish shows "calories incomplete" until a density is added. | default | [KIT-02-02](kitchen/E02-ingredient-catalogue/KIT-02-02-measures-and-units.md) |
| Q051 | low | How should "a pinch" and "to taste" count in calorie and stock maths? | Pinch = 0.5 g; "to taste" = 0 g, shown as "not counted". | default | [KIT-02-02](kitchen/E02-ingredient-catalogue/KIT-02-02-measures-and-units.md) |
| Q052 | low | Should name matching split German compound words such as "Hähnchenschenkel" before searching? | No splitting; Haiku resolves compounds when text matching is unsure. | default | [KIT-03-01](kitchen/E03-name-matching/KIT-03-01-name-normaliser.md) |
| Q053 | high | How big should the starter ingredient catalogue be, and who writes and checks it? | About 700 ingredients, AI-drafted, checked by the founder and a native speaker per language. | you confirmed | [KIT-02-07](kitchen/E02-ingredient-catalogue/KIT-02-07-starter-catalogue-seed.md) |
| Q054 | high | Should kitchen announce allergen or diet flag changes on an ingredient so planner re-checks every dish that uses it? | Yes: new event catalogue.ingredient_flags_changed. | you confirmed | [KIT-02-10](kitchen/E02-ingredient-catalogue/KIT-02-10-admin-catalogue-rpcs.md), [PLN-02-03](planner/E02-recipe-library-model/PLN-02-03-derived-flags.md), [ADM-03-03](admin/E03-catalogue-review/ADM-03-03-ingredient-editor.md), [README](README.md) |
| Q055 | low | Can an admin retire an ingredient that stock, products or recipes still use? | Yes: hidden from search and new matches; existing references keep working. | default | [KIT-02-10](kitchen/E02-ingredient-catalogue/KIT-02-10-admin-catalogue-rpcs.md) |
| Q056 | medium | Must an admin be able to undo an ingredient merge? | No undo; create a new ingredient if a merge was wrong. | default | [KIT-02-09](kitchen/E02-ingredient-catalogue/KIT-02-09-merge-ingredients.md) |
| Q057 | low | How similar must a text match be to accept it without AI or review? | Best score ≥ 0.85 and ahead of the runner-up by ≥ 0.15. | default | [KIT-03-03](kitchen/E03-name-matching/KIT-03-03-match-names-deterministic.md) |
| Q058 | medium | When an admin corrects a name match, should kitchen announce it so offers and planner fix stored links? | Yes: new event catalogue.match_changed. | default | [KIT-03-06](kitchen/E03-name-matching/KIT-03-06-match-review-queue.md), [README](README.md) |
| Q059 | medium | What should onboarding default to for users outside Poland and Germany? | Austria like Germany; others English, EUR, Brussels time, no home cuisine. | default | [KIT-10-01](kitchen/E10-kitchen-profile/KIT-10-01-profile-record-and-defaults.md) |
| Q060 | medium | Which tools should the "basic kitchen" preset tick? | Hob, oven, kettle, toaster, hand blender, baking tins. | default | [KIT-10-03](kitchen/E10-kitchen-profile/KIT-10-03-kitchen-tools-and-presets.md), [PLN-03-09](planner/E03-library-import-and-review/PLN-03-09-library-coverage-report.md) |
| Q061 | low | Should "kettle" be a kitchen tool that dishes can require? | Yes, the kettle is a tool (the "student room" preset uses it). | default | [KIT-10-03](kitchen/E10-kitchen-profile/KIT-10-03-kitchen-tools-and-presets.md) |
| Q062 | low | Within what limits may a household change the default 25/40/25/10 meal shares? | Breakfast 10–40%, main 25–60%, other 10–40%, snacks 0–20%. | default | [KIT-10-06](kitchen/E10-kitchen-profile/KIT-10-06-household-rhythm-and-meal-shares.md) |
| Q063 | medium | Should the default meal split differ by country (Polish midday main meal, German cold supper)? | One default, 25/40/25/10, for every country for now. | default | [KIT-10-06](kitchen/E10-kitchen-profile/KIT-10-06-household-rhythm-and-meal-shares.md) |
| Q064 | low | Which AI confidence levels accept a match, send it to review, or reject it? | ≥ 0.8 accept; 0.5–0.8 accept and review; below 0.5 none, review. | default | [KIT-03-04](kitchen/E03-name-matching/KIT-03-04-match-names-ai-pick.md) |
| Q065 | low | How should offers learn that bulk name matches are ready: by asking again or through an event? | Ask again: offers re-calls MatchNames in bulk mode every 30 minutes. | default | [KIT-03-05](kitchen/E03-name-matching/KIT-03-05-bulk-matching-batch.md) |
| Q066 | high | Should AI-translated ingredient names go live before a native speaker reviews them? | Yes, live at once; the native-speaker check per language follows. | you confirmed | [KIT-02-08](kitchen/E02-ingredient-catalogue/KIT-02-08-ai-translation-batch.md) |

### M5 Fridge and barcodes (R1)

| ID | Severity | Decision | Answer | Status | Tasks |
|---|---|---|---|---|---|
| Q067 | low | Should every barcode be stored in its 14-digit, zero-padded form as the single key? | Yes: 14-digit GTIN as the key; display 8, 12 or 13 digits. | default | [KIT-04-01](kitchen/E04-barcode-lookup/KIT-04-01-gtin-parsing.md) |
| Q068 | low | For weighed shop codes (prefix 20–29) from a chain we have no pattern for, which digits identify the item? | The first 7 digits (prefix plus 5 item digits), the common European layout. | default | [KIT-04-03](kitchen/E04-barcode-lookup/KIT-04-03-internal-codes.md) |
| Q069 | low | Should names given to weighed shop codes stay with the household (and chain) instead of being shared with everyone? | Per household, keyed by country, item digits and an optional chain. | default | [KIT-04-03](kitchen/E04-barcode-lookup/KIT-04-03-internal-codes.md) |
| Q070 | low | How many scanned products should a household's own product memory keep? | Up to 5,000; a nightly job drops the least recently scanned. | default | [KIT-04-04](kitchen/E04-barcode-lookup/KIT-04-04-household-products.md) |
| Q071 | high | Should Open Food Facts data be kept in its own table, apart from Coolinary's own catalogue? | Yes: a separate `off_products` table, never mixed with our catalogue. | you confirmed | [KIT-04-07](kitchen/E04-barcode-lookup/KIT-04-07-open-food-facts-mapping.md) |
| Q072 | medium | When a multipack is scanned (e.g. 6 × 0.5 l water), is it one fridge line or six? | One line with the total (3,000 ml, or 6 pieces). | default | [KIT-04-11](kitchen/E04-barcode-lookup/KIT-04-11-multipacks.md) |
| Q073 | low | How many unknown products may one household name per day? | 50 per day, then refused until the next day. | default | [KIT-05-01](kitchen/E05-unknown-products/KIT-05-01-name-product.md) |
| Q074 | low | When should a shared product go to admin review because households disagree about it? | 3 conflicting names within 30 days, or 3 households reporting it wrong. | default | [KIT-05-02](kitchen/E05-unknown-products/KIT-05-02-submissions-and-confirmation.md) |
| Q075 | medium | Should product names proposed to other households be checked against a swear-word list per language? | Yes: a word list per language (pl, de, en) before names are shared. | default | [KIT-05-02](kitchen/E05-unknown-products/KIT-05-02-submissions-and-confirmation.md) |
| Q076 | high | May a household opt out of having its confirmed product names sent to Open Food Facts? | Household opt-out — A household setting, on by default, keeps its confirmed product names out of what is sent to Open Food Facts. | **you chose** | [KIT-05-04](kitchen/E05-unknown-products/KIT-05-04-contribute-to-off.md) |
| Q077 | low | Should a household's fridge be capped at 1,000 items? | Yes, 1,000 items, as an abuse guard. | default | [KIT-07-01](kitchen/E07-fridge-stock/KIT-07-01-stock-schema.md), [KIT-13-02](kitchen/E13-household-lifecycle/KIT-13-02-member-joined-with-fridge.md) |
| Q078 | high | When food moves from the freezer to the fridge, how many days may its new use-by date be at most? | At most 2 days, or the fridge shelf life if shorter. | you confirmed | [KIT-07-02](kitchen/E07-fridge-stock/KIT-07-02-storage-and-expiry-defaults.md) |
| Q079 | medium | May users add fridge items that don't match any catalogue ingredient? | Yes: stored as "not matched" and matched later. | default | [KIT-07-03](kitchen/E07-fridge-stock/KIT-07-03-add-stock-items.md) |
| Q080 | low | When "cooked it" reduces stock, should frozen items be used only after fridge and pantry items? | Yes: freezer last; otherwise oldest date first. | default | [KIT-07-07](kitchen/E07-fridge-stock/KIT-07-07-consume-stock.md) |
| Q081 | low | Below what leftover amount should a partly used fridge item disappear automatically? | Below 2% of the original, or under 5 g, 5 ml or 0.1 piece. | default | [KIT-07-07](kitchen/E07-fridge-stock/KIT-07-07-consume-stock.md) |
| Q231 | high | When a household turns off "contribute to Open Food Facts", what happens to names it already sent? | Only future sends stop; names already sent stay in Open Food Facts. | you confirmed | [KIT-05-04](kitchen/E05-unknown-products/KIT-05-04-contribute-to-off.md) |

### M6 Nutrition and recipe library (R1)

| ID | Severity | Decision | Answer | Status | Tasks |
|---|---|---|---|---|---|
| Q082 | low | How should nutrition-table values written as "traces" or "< x" be counted? | "Traces" = 0; "< x" = half of x; both flagged. | default | [KIT-06-01](kitchen/E06-nutrition/KIT-06-01-nutrition-schema.md) |
| Q083 | low | Which checks must pack nutrition values from Open Food Facts pass before we use them? | Each 0–100 g, macros ≤ 101 g, kcal within ±20% of computed. | default | [KIT-06-07](kitchen/E06-nutrition/KIT-06-07-pack-nutrition.md) |
| Q084 | medium | For deep-fried dishes, how much frying oil counts toward calories when the recipe doesn't state it? | Recipe must state absorbed oil; otherwise "calories incomplete" and not publishable. | default | [KIT-06-08](kitchen/E06-nutrition/KIT-06-08-calculation-rules.md), [PLN-04-01](planner/E04-dish-nutrition/PLN-04-01-compute-dish-nutrition.md) |
| Q085 | low | Should nutrient values shown to users be rounded the way EU food labels round them? | Yes: EU label rounding; kJ to the nearest 10. | default | [KIT-06-08](kitchen/E06-nutrition/KIT-06-08-calculation-rules.md) |
| Q086 | medium | Can a dish be published in the library while its calories are incomplete? | Not for the 200 starter dishes; later yes, if an admin confirms. | default | [PLN-02-02](planner/E02-recipe-library-model/PLN-02-02-dish-validation.md) |
| Q087 | medium | Should a dish needing only a microwave (or kettle, toaster, blender) count as "no cooking"? | Yes: microwave counts as no cooking. | default | [PLN-02-02](planner/E02-recipe-library-model/PLN-02-02-dish-validation.md) |
| Q088 | low | When a recipe is scaled above twice its base servings, should it warn that cooking may take longer? | Yes: show "cooking times may be longer" above 2× base servings. | default | [PLN-02-05](planner/E02-recipe-library-model/PLN-02-05-portion-scaling.md) |
| Q089 | low | Should kitchen mark ingredients like eggs as whole-only, so a scaled recipe never asks for 1.5 eggs? | Yes: add whole_units_only; round half up, never below 1. | default | [PLN-02-05](planner/E02-recipe-library-model/PLN-02-05-portion-scaling.md) |
| Q090 | low | Should kitchen's name matching accept "recipe" as a context, next to product, receipt and offer? | Yes: add context recipe to MatchNames. | default | [PLN-03-02](planner/E03-library-import-and-review/PLN-03-02-ingredient-name-resolution.md) |
| Q091 | low | Which match scores should auto-accept or reject ingredient names in imported or drafted recipes? | Accept ≥ 0.90, review 0.60–0.90, reject < 0.60; tune after 200 dishes. | default | [PLN-03-02](planner/E03-library-import-and-review/PLN-03-02-ingredient-name-resolution.md) |
| Q092 | medium | For base-plus-toppings dishes, which toppings should the headline calories assume? | Base plus the first option of each topping group, plus a min–max range. | default | [PLN-04-01](planner/E04-dish-nutrition/PLN-04-01-compute-dish-nutrition.md) |
| Q093 | high | Must a dietitian check each starter dish before it can be published? | No: publish, but launch waits until every starter dish is checked. | you confirmed | [PLN-03-08](planner/E03-library-import-and-review/PLN-03-08-dietitian-sign-off.md) |
| Q094 | high | What share of starter dishes must be within 10% of the dietitian's calories before launch? | 95%, with every starter dish checked. | you confirmed | [PLN-03-08](planner/E03-library-import-and-review/PLN-03-08-dietitian-sign-off.md) |
| Q095 | low | When nutrition data changes, should calories in past history and old plans be recalculated? | No: history keeps the values it showed. | default | [PLN-04-02](planner/E04-dish-nutrition/PLN-04-02-recompute-on-catalogue-events.md) |
| Q096 | high | Should the server remove calories for people who hide them, or leave hiding to the app? | The app hides them, using the profile flag. | you confirmed | [PLN-04-03](planner/E04-dish-nutrition/PLN-04-03-nutrition-in-responses.md) |
| Q097 | low | Should kitchen's nutrition calculation accept a scanned product (EAN) per ingredient line? | Yes: optional EAN per line in ComputeNutrition. | default | [PLN-04-04](planner/E04-dish-nutrition/PLN-04-04-scanned-product-override.md) |
| Q098 | low | Where should a household's own scanned pack values replace nutrition-table values? | Only in dish detail and cooking mode, not browse or plans. | default | [PLN-04-04](planner/E04-dish-nutrition/PLN-04-04-scanned-product-override.md) |
| Q099 | medium | When browsing the library, whose diets and allergens filter dishes by default? | Everyone in the household, with a "for me" switch. | default | [PLN-05-02](planner/E05-library-browse/PLN-05-02-hard-limits.md) |
| Q100 | high | How many starter dishes should exist per cuisine and meal? | ≥ 8 per cuisine×meal (basic kitchen), ≥ 4 student room, ≥ 3 per diet. | you confirmed | [PLN-03-09](planner/E03-library-import-and-review/PLN-03-09-library-coverage-report.md) |
| Q101 | medium | Should favourites be one shared household list or kept per person? | One shared list per household. | default | [PLN-05-04](planner/E05-library-browse/PLN-05-04-favourites.md) |
| Q102 | medium | How long should the "cooked it" history be kept? | 2 years, then deleted by a monthly job. | default | [PLN-05-06](planner/E05-library-browse/PLN-05-06-cooked-history.md) |
| Q103 | high | Should items past their use-by date count as "at home" when suggesting what to cook? | Count with a warning — Items past their date count as at home, always marked "check it first" in suggestions, AI plans and cooking mode. | **you chose** | [PLN-06-01](planner/E06-what-can-i-make-now/PLN-06-01-stock-coverage.md) |
| Q104 | medium | How should "what can I make now" weigh ingredients at home against items about to expire? | 70% share at home, 30% expiry; +0.05 favourite, −0.20/−0.10 recently cooked; hide < 30%. | default | [PLN-06-02](planner/E06-what-can-i-make-now/PLN-06-02-ranking.md) |
| Q105 | medium | Should one kitchen tool ever stand in for another (e.g. a hob for a wok) when checking dishes? | No substitutions at launch. | default | [PLN-07-02](planner/E07-rule-checker/PLN-07-02-tool-check.md) |
| Q106 | critical | Should Coolinary track "may contain" (cross-contact) allergen warnings from product packs? | No: state in the app and terms that cross-contact isn't checked. | you confirmed | [PLN-07-03](planner/E07-rule-checker/PLN-07-03-allergen-diet-dislike-check.md) |
| Q107 | low | How far may a plan exceed the amount at home before a meal is flagged "not enough at home"? | 10%: above it the meal is flagged; 0–10% gives a soft warning. | default | [PLN-07-04](planner/E07-rule-checker/PLN-07-04-stock-check.md) |
| Q108 | low | How many minutes over the household's time per meal may a dish run before it's rejected? | 10 minutes; within that it's marked "tight". | default | [PLN-07-05](planner/E07-rule-checker/PLN-07-05-ai-output-checks.md) |

### M7 AI intake (R1)

| ID | Severity | Decision | Answer | Status | Tasks |
|---|---|---|---|---|---|
| Q109 | medium | After a photo or text scan, which proposed items should be ticked ready to add? | Matched items with confidence ≥ 0.5 are pre-ticked. | default | [KIT-08-01](kitchen/E08-ai-intake/KIT-08-01-intake-pipeline.md) |
| Q110 | medium | Should AI read the nutrition table from a label photo in release 1? | No: read name, brand and size only in R1. | default | [KIT-05-03](kitchen/E05-unknown-products/KIT-05-03-label-photo.md), [AI-07-01](ai/E07-label-reading/AI-07-01-read-product-label.md) |
| Q111 | high | When does Free's weekly AI allowance (2 day plans, 2 scans) reset: every Monday 00:00 UTC, or on a rolling 7 days? | Rolling 7 days — Each use frees up again 7 days after it was made; the app shows when the next slot frees. | **you chose** | [GW-03-01](gateway/E03-free-tier-quota/GW-03-01-quota-counters.md) |
| Q112 | high | Do label photos of unknown products and typed or spoken item lists count toward Free's 2 scans per week? | Yes, every AI reading request counts as one scan | you confirmed | [GW-03-01](gateway/E03-free-tier-quota/GW-03-01-quota-counters.md), [GW-06-05](gateway/E06-fridge-and-intake/GW-06-05-label-photo.md), [GW-06-06](gateway/E06-fridge-and-intake/GW-06-06-text-intake.md), [KIT-08-03](kitchen/E08-ai-intake/KIT-08-03-text-intake.md) |
| Q113 | high | Should Plus have a fair-use cap of 60 AI responses per household per day, stated in the terms? | Yes, 60 per household per day | you confirmed | [GW-03-04](gateway/E03-free-tier-quota/GW-03-04-rate-limits.md) |
| Q114 | low | Should a long receipt be photographed in up to 3 parts rather than as one tall photo? | Yes, up to 3 parts so small print stays readable | default | [MOB-07-01](mobile/E07-uploads/MOB-07-01-image-preparation.md) |
| Q115 | medium | Should several flyer pages or receipt parts go up in one request (one AI call, one quota count) instead of one request per page? | One request with all pages | default | [MOB-07-02](mobile/E07-uploads/MOB-07-02-upload-transport.md) |
| Q116 | low | Should the server accept ids made by the phone for new fridge items and list lines, so offline edits can refer to them? | Yes, accept phone-made UUIDs for stock items and list lines | default | [MOB-03-05](mobile/E03-query-hooks/MOB-03-05-fridge-barcode-intake-hooks.md), [MOB-03-07](mobile/E03-query-hooks/MOB-03-07-shopping-list-hooks.md) |

### M8 AI day plan (R1)

| ID | Severity | Decision | Answer | Status | Tasks |
|---|---|---|---|---|---|
| Q117 | high | In R1, can users put a library dish they picked into a day-plan meal themselves? | Yes; add a rule-checked SetPlanMeal (no AI, no quota) | you confirmed | [GW-07-05](gateway/E07-recipes-and-plans/GW-07-05-plans-read-edit.md), [PLN-08-01](planner/E08-ai-day-plan/PLN-08-01-plan-storage.md) |
| Q118 | medium | How long should day plans be kept? | Unaccepted plans 7 days; accepted plans 1 year. | default | [PLN-08-01](planner/E08-ai-day-plan/PLN-08-01-plan-storage.md) |
| Q119 | medium | Should breakfast and the non-main meal prefer no-cooking dishes from launch, not only in family plans? | Yes, as a soft preference from R1; households can turn it off. | default | [PLN-08-02](planner/E08-ai-day-plan/PLN-08-02-compose-request.md) |
| Q120 | medium | If every AI meal fails the checks and the plan falls back to library dishes, does it use one of Free's 2 weekly plans? | No: count only when at least one AI meal survives. | default | [PLN-08-04](planner/E08-ai-day-plan/PLN-08-04-validate-repair-store.md) |
| Q121 | high | Should AI draft new dishes in English and translate them, or write each cuisine's dishes in its home language first? | Draft in English, keep the original dish name, then translate to pl and de. | you confirmed | [AI-09-03](ai/E09-catalogue-library-batch-jobs/AI-09-03-draft-dishes.md) |
| Q122 | high | Can dishes go live with machine-translated names and steps? | Yes, after a native speaker spot-checks 10% of dishes per language. | you confirmed | [PLN-03-05](planner/E03-library-import-and-review/PLN-03-05-dish-translation-batch.md) |
| Q123 | medium | How much better must Sonnet score than Haiku on the day-plan test set before Plus plans use Sonnet (which costs about 20× more)? | 10 percentage points more plans passing first time; Sonnet repairs Plus only. | default | [AI-10-05](ai/E10-day-plan/AI-10-05-model-routing-by-plan.md) |
| Q124 | medium | Should Free users get a simple non-AI "replace this ingredient" (same food category)? | No at launch: Free gets "add to shopping list" instead. | default | [PLN-09-02](planner/E09-swaps-and-rewrites/PLN-09-02-replace-ingredient.md) |
| Q125 | low | From R2, should every Plus household's day plan go through the family planner, even single-person households? | Yes, one code path; a single person is just one plate | default | [GW-07-04](gateway/E07-recipes-and-plans/GW-07-04-ai-day-plan.md) |

### M9 Shopping list and cooking (R1)

| ID | Severity | Decision | Answer | Status | Tasks |
|---|---|---|---|---|---|
| Q126 | low | Should kitchen's ConsumeStock accept an idempotency key and preferred stock items? | Yes, both. | default | [PLN-10-02](planner/E10-cooking-mode/PLN-10-02-mark-cooked.md) |
| Q127 | medium | Should stock that expires before a planned meal still be subtracted from the shopping list? | No: it doesn't count, so the item goes on the list. | default | [PLN-15-03](planner/E15-shopping-list/PLN-15-03-subtract-stock-and-staples.md) |
| Q128 | low | How should shopping-list amounts round up to packs? | ≤ 2 pack sizes, ≤ 20 packs; else round to 10 g (≤ 100) or 50 g. | default | [PLN-15-04](planner/E15-shopping-list/PLN-15-04-merge-and-round-to-packs.md) |
| Q129 | low | When a typed list line not matched to the catalogue is ticked, should it go into the fridge? | No: it's ticked on the list only. | default | [PLN-15-07](planner/E15-shopping-list/PLN-15-07-tick-into-fridge.md) |

### M10 Notifications and reminders (R1)

| ID | Severity | Decision | Answer | Status | Tasks |
|---|---|---|---|---|---|
| Q130 | low | Should the app unregister push notifications by its installation id instead of having to keep the push token until sign-out? | Yes; add installation_id to notifier's UnregisterDevice | default | [GW-11-01](gateway/E11-devices-and-notifications/GW-11-01-devices.md), [NTF-04-01](notifier/E04-devices-preferences/NTF-04-01-devices.md) |
| Q131 | low | How many phones may receive push for one account? | 10; an 11th replaces the least recently seen. | default | [NTF-04-01](notifier/E04-devices-preferences/NTF-04-01-devices.md) |
| Q132 | low | After how many days unseen should a phone's push registration be deleted? | 90 days. | default | [NTF-04-01](notifier/E04-devices-preferences/NTF-04-01-devices.md) |
| Q133 | medium | What should the default quiet hours be, when no push is sent? | 21:00–08:00. | default | [NTF-04-02](notifier/E04-devices-preferences/NTF-04-02-preferences.md) |
| Q134 | medium | At what local time should the daily expiry reminder go out? | 09:00 local time, in time to plan the day's meals. | default | [KIT-09-02](kitchen/E09-expiry-reminders/KIT-09-02-reminder-scheduler.md) |
| Q135 | medium | When should the app first ask permission to send notifications? | When a Plus household saves dated items, turns on a category, or starts sharing | default | [MOB-06-02](mobile/E06-push/MOB-06-02-permission-timing.md) |
| Q136 | high | Should signing out delete the age, height and weight stored on the phone? | Yes, so the next person on the phone doesn't inherit them | you confirmed | [MOB-01-07](mobile/E01-auth-session/MOB-01-07-sign-out.md) |

### M11 Account lifecycle and launch readiness (R1)

| ID | Severity | Decision | Answer | Status | Tasks |
|---|---|---|---|---|---|
| Q137 | high | When the owner of a shared household deletes their account, should ownership pass automatically to the longest-standing member? | Owner chooses first — An owner of a shared household must name the new owner before deleting their account; with no other account member, the household is deleted with the account. | **you chose** | [IDN-09-02](identity/E09-account-deletion-and-export/IDN-09-02-delete-account.md), [GW-02-05](gateway/E02-auth-and-account/GW-02-05-delete-account.md) |
| Q138 | low | In the GDPR export, should each service return its data as its own JSON, nested by service? | Yes: one JSON string per service, nested by the gateway. | default | [IDN-09-03](identity/E09-account-deletion-and-export/IDN-09-03-export-user-data.md), [GW-02-06](gateway/E02-auth-and-account/GW-02-06-data-export.md) |
| Q139 | medium | When someone moves to another household, should their old household's diet settings travel with them as their personal diet? | Yes, if they were that household's only account member. | default | [KIT-13-02](kitchen/E13-household-lifecycle/KIT-13-02-member-joined-with-fridge.md) |
| Q140 | low | When someone joins with their fridge and the total would pass 1,000 items, what happens? | Move the soonest-expiring items first, up to 1,000; the rest stay behind. | default | [KIT-13-02](kitchen/E13-household-lifecycle/KIT-13-02-member-joined-with-fridge.md) |
| Q141 | medium | How long should individual fridge movements (used, cooked, thrown away) be kept? | 24 months; monthly savings totals are kept longer. | default | [KIT-13-08](kitchen/E13-household-lifecycle/KIT-13-08-retention-job.md) |
| Q142 | high | Should anonymous accounts (no e-mail) be able to download a GDPR data export? | No; they can delete their account instead | you confirmed | [GW-02-06](gateway/E02-auth-and-account/GW-02-06-data-export.md) |
| Q143 | low | Should the data export be an immediate JSON download, or an e-mailed link to a file built in the background? | Immediate JSON download; one household's data is small | default | [GW-02-06](gateway/E02-auth-and-account/GW-02-06-data-export.md) |

### M12 Shops and flyer pipeline (R2)

| ID | Severity | Decision | Answer | Status | Tasks |
|---|---|---|---|---|---|
| Q144 | low | Should shop locations come from the Overpass API (one query per chain and country, monthly) or from downloaded country map files? | Overpass API, one query per chain per country, monthly. | default | [OFF-02-03](offers/E02-chains-and-shops/OFF-02-03-osm-import.md) |
| Q145 | low | When the founder rejects a shop a user added, should it disappear for everyone, including the household that added it? | Hidden for everyone and removed from the adder's my shops. | default | [OFF-02-05](offers/E02-chains-and-shops/OFF-02-05-user-adds-shop.md) |
| Q146 | medium | Should approved shop corrections from users be uploaded to OpenStreetMap automatically? | No; the founder edits OpenStreetMap by hand when worthwhile. | default | [OFF-02-06](offers/E02-chains-and-shops/OFF-02-06-user-corrects-shop.md) |
| Q147 | high | Which postcode dataset should map a typed postcode to a place: GeoNames (free, credit required) or another source? | GeoNames postal codes, CC BY 4.0, credit "GeoNames", refreshed yearly. | you confirmed | [OFF-03-01](offers/E03-areas-and-editions/OFF-03-01-postcode-dataset.md) |
| Q148 | medium | Should Aldi Nord's regional extra offers be skipped in release 2 and only its national flyer read? | Skip regional extras in R2; read the national flyer only. | default | [OFF-03-04](offers/E03-areas-and-editions/OFF-03-04-editions-model.md), [OFF-06-06](offers/E06-launch-chains/OFF-06-06-aldi-nord-de.md) |
| Q149 | low | How many shops at most should one "shops near this postcode" answer return? | 1,000 shops (central Berlin has about 400). | default | [OFF-03-03](offers/E03-areas-and-editions/OFF-03-03-shops-near-area.md) |
| Q150 | high | Should Kaufland shops be linked to their store numbers using the chain's public store finder, under the same terms check as its flyers? | Yes: read the public store finder, covered by the same retailer-terms gate. | you confirmed | [OFF-03-05](offers/E03-areas-and-editions/OFF-03-05-chain-store-numbers.md) |
| Q151 | low | How long should the flyer reader wait between two requests to the same retailer website? | 3 seconds, or the site's own Crawl-delay if larger. | default | [OFF-05-03](offers/E05-flyer-sync/OFF-05-03-polite-fetcher.md) |
| Q152 | medium | When a known flyer disappears from the chain's flyer page, when should its offers be removed? | Missing twice before it starts: withdrawn. Once started: kept to its end date. | default | [OFF-05-04](offers/E05-flyer-sync/OFF-05-04-sync-and-detection.md) |
| Q153 | low | Where should downloaded flyer pages be stored until review: a Docker volume or the SeaweedFS store already used by Langfuse? | Named Docker volume offers-flyers at /data/flyers. | default | [OFF-05-05](offers/E05-flyer-sync/OFF-05-05-page-storage.md) |
| Q154 | medium | When a flyer is found less than 36 hours before it starts, should it be read at full AI price instead of the half-price batch? | Yes: read at once at full price when < 36 h before start. | default | [OFF-07-02](offers/E07-offer-pipeline/OFF-07-02-batch-extraction.md) |
| Q155 | medium | Should non-food offers (electronics, clothes, garden) be thrown away rather than stored? | Drop them; never stored. | default | [OFF-07-03](offers/E07-offer-pipeline/OFF-07-03-normalise-offers.md) |
| Q156 | low | Should offers be grouped by the app's own ingredient categories or by each chain's flyer sections? | The catalogue's ingredient categories, one set across the app. | default | [OFF-07-04](offers/E07-offer-pipeline/OFF-07-04-match-products.md) |
| Q157 | medium | Besides the checks in features.md, should offers be held for review when AI confidence is below 0.6 or the price per kg is outside 0.2×–5× the usual price? | Yes: hold below 0.6 confidence or outside 0.2×–5× reference price. | default | [OFF-07-05](offers/E07-offer-pipeline/OFF-07-05-validation-rules.md) |
| Q158 | high | If a chain's weekly spot-check scores below 97% two weeks in a row (or once below 90%), should it be taken offline automatically? | Alert only — Bad weekly spot-checks only alert the founder; a chain is never taken offline automatically. | **you chose** | [OFF-08-03](offers/E08-accuracy-gate/OFF-08-03-weekly-spot-checks.md) |
| Q159 | low | Should a shop correction be accepted automatically when 3 households propose the same change? | No; review every correction by hand until volumes are known. | default | [ADM-06-06](admin/E06-flyer-review/ADM-06-06-shop-corrections.md) |

### M13 Launch chains go live (R2)

| ID | Severity | Decision | Answer | Status | Tasks |
|---|---|---|---|---|---|
| Q160 | high | When a Biedronka flyer shows only a start date, should the app assume it lasts 7 days? | End = start + 6 days; dates printed on offers may shorten it. | you confirmed | [OFF-06-01](offers/E06-launch-chains/OFF-06-01-biedronka-pl.md) |

### M14 Offers in the app (R2)

| ID | Severity | Decision | Answer | Status | Tasks |
|---|---|---|---|---|---|
| Q161 | low | How many shops at most may a Plus household save as "my shops"? | 40 shops (Free stays at 2). | default | [OFF-04-01](offers/E04-my-shops/OFF-04-01-my-shops.md) |
| Q162 | medium | How many nearby shops may the phone compare in one trip plan? | Up to 10 shops (175 shop combinations with 3 per trip) | default | [PLN-19-01](planner/E19-trip-price-matrix/PLN-19-01-candidate-shops.md), [GW-10-05](gateway/E10-offers-and-shops/GW-10-05-price-matrix.md) |
| Q163 | medium | How many shops should one trip include by default: 1, 2 or 3? | 2 shops per trip. | default | [OFF-04-02](offers/E04-my-shops/OFF-04-02-shopping-preferences.md) |
| Q164 | medium | What travel range should a new household start with: 1 km walking, 3 km short trip, or 10 km by car? | 3 km (short trip). | default | [OFF-04-02](offers/E04-my-shops/OFF-04-02-shopping-preferences.md) |
| Q165 | medium | How far should a user be allowed to walk between two shops on one trip, and what is the default? | 0–1,500 m in 100 m steps; default 500 m. | default | [OFF-04-02](offers/E04-my-shops/OFF-04-02-shopping-preferences.md) |
| Q166 | medium | features.md sets €2 as the minimum saving for adding an extra shop; what should it be in Poland? | PLN 8.00 in Poland; €2.00 in euro countries. | default | [OFF-04-02](offers/E04-my-shops/OFF-04-02-shopping-preferences.md) |
| Q167 | low | When no household wants a Kaufland store edition any more, for how many days should it still be read? | 30 more days. | default | [OFF-04-03](offers/E04-my-shops/OFF-04-03-edition-demand.md) |
| Q168 | low | After a Plus household drops to Free, how long should its extra shops (beyond 2) be kept in case it upgrades again? | 90 days, then deleted. | default | [OFF-04-04](offers/E04-my-shops/OFF-04-04-plan-change-rules.md) |
| Q169 | low | How many Kaufland store editions per chain should be checked each day at most? | 300 per day, highest demand first, with a daily alert above it. | default | [OFF-05-07](offers/E05-flyer-sync/OFF-05-07-on-demand-editions.md) |
| Q170 | medium | Should offers show the product name as printed in the flyer, with the matched ingredient name in the user's language underneath? | As printed, plus the ingredient name in the user's language. | default | [OFF-09-01](offers/E09-offer-browse/OFF-09-01-browse-offers.md) |
| Q171 | medium | Should a Plus household be limited to 10 own-flyer uploads per day, even though Plus is "unlimited"? | Yes: 10 uploads per household per day as an abuse guard. | default | [OFF-10-01](offers/E10-user-flyers/OFF-10-01-upload-and-jobs.md) |
| Q172 | low | For own-flyer links, should only direct PDF or image links be accepted in release 2? | Only direct PDF or image links in R2. | default | [OFF-10-02](offers/E10-user-flyers/OFF-10-02-flyer-link.md) |
| Q173 | high | Should offers read from a flyer a user uploads be visible only to that household? | Yes, private to the uploading household; never used as public prices | you confirmed | [OFF-10-04](offers/E10-user-flyers/OFF-10-04-confirm-private-offers.md), [GW-10-04](gateway/E10-offers-and-shops/GW-10-04-own-flyer.md) |
| Q174 | high | Where should the "usual price" of each ingredient come from when no shop has an offer? | Weekly median of 90 days of observed prices (≥ 3 from ≥ 2 chains); admin CSV overrides. | you confirmed | [OFF-11-03](offers/E11-prices/OFF-11-03-reference-prices.md) |
| Q175 | medium | Should prices that need a shop's app or loyalty card count in the trip price comparison? | Include them, flagged as "with app or card". | default | [OFF-11-04](offers/E11-prices/OFF-11-04-get-prices.md) |
| Q176 | low | What is the smallest discount shown on an offer badge? | 5%, rounded down to a whole percent; else the price only. | default | [PLN-17-01](planner/E17-offer-badges-and-reminders/PLN-17-01-offer-badges.md) |
| Q177 | low | Should list lines show offers starting within 3 days when no current offer exists? | Yes, as "upcoming" badges, 3-day window. | default | [PLN-17-01](planner/E17-offer-badges-and-reminders/PLN-17-01-offer-badges.md) |
| Q178 | medium | Do Free households get "offer ends tomorrow" reminders? | Yes: Free and Plus both. | default | [PLN-17-04](planner/E17-offer-badges-and-reminders/PLN-17-04-offer-ends-reminders.md) |
| Q179 | low | At what local time should "offer ends tomorrow" reminders be sent? | 18:00 local time. | default | [PLN-17-04](planner/E17-offer-badges-and-reminders/PLN-17-04-offer-ends-reminders.md) |

### M15 Shared household (R2)

| ID | Severity | Decision | Answer | Status | Tasks |
|---|---|---|---|---|---|
| Q180 | low | How many profiles without an account (children, people without a phone) may one household have? | At most 10. | default | [IDN-06-01](identity/E06-profiles-without-accounts/IDN-06-01-add-profile.md) |
| Q181 | high | Which web domain should invite links use (it needs universal-link and app-link files)? | https://coolinary.app/i/<code>, once the domain is bought | you confirmed | [IDN-07-01](identity/E07-invites/IDN-07-01-invite-table-and-codes.md), [GW-04-03](gateway/E04-household-and-invites/GW-04-03-invites-create-list-cancel.md) |
| Q182 | low | Should waiting invites count toward the 5-account limit? | Yes: accounts plus waiting invites never exceed 5. | default | [IDN-07-02](identity/E07-invites/IDN-07-02-create-list-cancel.md) |
| Q183 | medium | Who may invite people into the household? | Any member with an account. | default | [IDN-07-02](identity/E07-invites/IDN-07-02-create-list-cancel.md) |
| Q184 | low | Should the owner get a push when someone joins the household? | Yes: "Ewa joined your household". | default | [IDN-07-04](identity/E07-invites/IDN-07-04-accept-invite.md) |
| Q185 | high | When the owner leaves, does Plus stay with the household or go with the person? | Stays with the household until the payment route is chosen. | you confirmed | [IDN-08-01](identity/E08-leave-and-remove/IDN-08-01-leave-household.md) |
| Q186 | medium | When the owner leaves a shared household, must they name the new owner first? | Yes, the owner picks the new owner before leaving | default | [IDN-08-01](identity/E08-leave-and-remove/IDN-08-01-leave-household.md), [GW-04-05](gateway/E04-household-and-invites/GW-04-05-leave-remove.md) |
| Q187 | high | When a shared household drops from Plus to Free, what happens to the other members? | Members stay, no new invites — After a downgrade, existing members keep sharing; nobody new can be invited or join until Plus returns. This relaxes features.md's "sharing needs Plus". | **you chose** | [IDN-05-04](identity/E05-plan-and-entitlements/IDN-05-04-shared-household-downgrade.md) |
| Q188 | low | Should members get a push warning before shared access ends? | Yes: push in category household, template sharing_ends. | default | [IDN-05-04](identity/E05-plan-and-entitlements/IDN-05-04-shared-household-downgrade.md), NTF-03, NTF-04 |
| Q189 | medium | Should a removed member get a push telling them? | No push; the app says "You are no longer in <household>" on next use. | default | [IDN-08-02](identity/E08-leave-and-remove/IDN-08-02-remove-member.md) |
| Q190 | low | Should identity announce an ownership change as an event to other services? | No event; no other service stores who the owner is. | default | [IDN-08-03](identity/E08-leave-and-remove/IDN-08-03-transfer-ownership.md) |
| Q191 | high | When someone is removed from a household, should their access end at once, or within up to 15 minutes when their token expires? | At once: a revocation key checked on every request (one Redis read) | you confirmed | [IDN-08-04](identity/E08-leave-and-remove/IDN-08-04-revoke-live-tokens.md), [GW-04-05](gateway/E04-household-and-invites/GW-04-05-leave-remove.md) |

### M16 Family goals and portions (R2)

| ID | Severity | Decision | Answer | Status | Tasks |
|---|---|---|---|---|---|
| Q192 | critical | Who may mark another account in the household as a child, which locks that person's goal? | Only the household owner; any adult for profiles without an account. | you confirmed | [KIT-11-02](kitchen/E11-per-person-profiles/KIT-11-02-edit-rights.md) |
| Q193 | medium | What is the highest daily calorie target an adult may set? | 5,000 kcal. | default | [KIT-11-03](kitchen/E11-per-person-profiles/KIT-11-03-daily-kcal-target.md) |
| Q194 | critical | Which absolute range should the server accept for a child's daily target, pending a paediatric dietitian's review? | 500–4,000 kcal on the server; age-based limits checked on the phone. | you confirmed | [KIT-11-03](kitchen/E11-per-person-profiles/KIT-11-03-daily-kcal-target.md) |
| Q195 | critical | After a household drops from Plus to Free, should each person's allergens and diets still apply to plans? | Allergens join household list — On a downgrade every person's allergens are added to the household allergen filter; per-person diets stop applying; editing per-person data waits for Plus. | **you chose** | [KIT-11-04](kitchen/E11-per-person-profiles/KIT-11-04-per-person-restrictions.md) |
| Q196 | medium | Besides one-off "away" marks, should people get weekly rules like "Ola away for dinner Mon–Fri"? | Yes: weekly rules, at most 20 per household. | default | [PLN-12-01](planner/E12-per-person-portions/PLN-12-01-meal-attendance.md) |
| Q197 | medium | If a household cooks for more people than it has profiles, how are the extra eaters counted? | As unnamed eaters with 1 standard portion each. | default | [PLN-12-02](planner/E12-per-person-portions/PLN-12-02-meal-shares.md) |
| Q198 | medium | When a dish's calories are incomplete, how are portions split in a household with goals? | Everyone gets 1 standard portion, flagged "calories incomplete". | default | [PLN-12-03](planner/E12-per-person-portions/PLN-12-03-portion-and-pot-calculator.md) |
| Q199 | low | In which order are no-cook extras offered to someone short on calories? | Bread, fruit, yoghurt, milk, cheese, nuts, raw veg; lighter order for "eat lighter". | default | [PLN-12-04](planner/E12-per-person-portions/PLN-12-04-no-cook-extras.md) |
| Q200 | critical | Should nuts be excluded as a no-cook extra for every child profile? | Yes: never offer nuts to child profiles. | you confirmed | [PLN-12-04](planner/E12-per-person-portions/PLN-12-04-no-cook-extras.md) |
| Q201 | low | How much may one meal's calorie gap shift the next meal? | At most ±25% of the next meal's target. | default | [PLN-12-05](planner/E12-per-person-portions/PLN-12-05-carry-over-and-day-balance.md) |
| Q202 | medium | Should a person's planned day be checked against their target minus the 10% snack share? | Yes: compare meals with 90% of the target. | default | [PLN-12-05](planner/E12-per-person-portions/PLN-12-05-carry-over-and-day-balance.md) |
| Q203 | high | Who may see a child's calories and target when they are not hidden? | Every adult account member of the household. | you confirmed | [PLN-12-07](planner/E12-per-person-portions/PLN-12-07-calorie-visibility.md) |
| Q204 | low | When does one person's diet "narrow the choice too much", so planner prefers base-plus-toppings dishes? | When fewer than 5 dishes fit everyone, or under 30% of those without them. | default | [PLN-13-01](planner/E13-family-day-plan/PLN-13-01-shared-meal-limits.md) |
| Q205 | medium | Where does "eat lighter as a family" rank when planner compares candidate plans? | After each person's day target, before using up stock. | default | [PLN-13-03](planner/E13-family-day-plan/PLN-13-03-plan-comparator.md) |
| Q206 | low | How should "use the oldest stock first" be scored when comparing family plans? | Share used × 1/(1 + days to expiry); undated items count 0.05. | default | [PLN-13-03](planner/E13-family-day-plan/PLN-13-03-plan-comparator.md) |
| Q207 | medium | When a plan change raises amounts or adds ingredients, should the list update itself or ask first? | Ask: increases go into a "plan changed, review" proposal. | default | [PLN-15-10](planner/E15-shopping-list/PLN-15-10-list-follows-plan.md) |
| Q208 | medium | How strict is "fridge first" when planning from this week's deals? | Each meal uses ≥ 1 stock item; ≥ 50% by weight from stock; ≤ 3 non-offer buys a day. | default | [PLN-18-02](planner/E18-plan-from-deals/PLN-18-02-ai-plan-from-deals.md) |
| Q230 | high | On a Plus → Free downgrade, do the health-related diets gluten-free, lactose-free and nut-free join the household list with the allergens? | Yes; vegetarian, vegan and dislikes stop applying. | you confirmed | [KIT-11-04](kitchen/E11-per-person-profiles/KIT-11-04-per-person-restrictions.md), [KIT-13-06](kitchen/E13-household-lifecycle/KIT-13-06-plan-changed.md), [PLN-21-03](planner/E21-household-lifecycle/PLN-21-03-plan-downgrade.md) |

### M17 Library growth, rewrites and savings (R2)

| ID | Severity | Decision | Answer | Status | Tasks |
|---|---|---|---|---|---|
| Q209 | medium | How should euro savings handle fridge items whose purchase price is unknown? | Count kilograms only; euros only where the price paid is known. | default | [KIT-12-01](kitchen/E12-savings/KIT-12-01-movement-valuation.md) |
| Q210 | medium | Should savings counters run for Free and anonymous households too? | Yes, for everyone, so a household that upgrades sees its history. | default | [KIT-12-02](kitchen/E12-savings/KIT-12-02-ledger-rules.md) |
| Q211 | high | How close to its use-by date must food be used to count as "saved"? | Used within the 3 days before its date (the app's "soon"). | you confirmed | [KIT-12-02](kitchen/E12-savings/KIT-12-02-ledger-rules.md) |
| Q212 | high | What savings does a Free household see? | This month's counters on Free; monthly report and past months on Plus | you confirmed | [GW-12-01](gateway/E12-savings/GW-12-01-savings-counters-report.md), [KIT-12-03](kitchen/E12-savings/KIT-12-03-get-savings.md) |
| Q213 | low | When should the monthly savings report be pushed? | The 1st of the month at 10:00 local time. | default | [KIT-12-04](kitchen/E12-savings/KIT-12-04-monthly-report-push.md) |
| Q214 | medium | When does a household's AI dish become a candidate for the shared library? | Rated ≥ 4.0, cooked ≥ 2 times, at least 14 days old. | default | [PLN-11-01](planner/E11-library-growth/PLN-11-01-candidate-selection.md) |
| Q215 | high | How are users told their AI dishes may join the shared library in anonymised form? | A sentence in the terms and on the AI-dish card; no opt-in. | you confirmed | [PLN-11-02](planner/E11-library-growth/PLN-11-02-anonymised-copy.md) |
| Q216 | medium | Should planner skip the AI call when library dishes already cover ≥ 95% of the stock for every meal? | Measure first; decide with launch data. | default | [PLN-11-03](planner/E11-library-growth/PLN-11-03-promotion-admin-and-metrics.md) |

### M18 Week plan (R3)

| ID | Severity | Decision | Answer | Status | Tasks |
|---|---|---|---|---|---|
| Q217 | medium | For "cook twice a week", which days are cooking days, and may one pot cover three meals? | Days 1 and 4, double amounts only; other main meals no-cook. | default | [PLN-14-02](planner/E14-week-plan/PLN-14-02-cooking-days.md) |
| Q218 | low | How should the week's "time at the stove" estimate treat double amounts and toppings? | Double amount ×1.25 time; +10 minutes per topping needing its own pan. | default | [PLN-14-04](planner/E14-week-plan/PLN-14-04-effort-summary.md) |

### M19 Receipts and receipt prices (R3)

| ID | Severity | Decision | Answer | Status | Tasks |
|---|---|---|---|---|---|
| Q219 | high | Should typed-sentence intake ("bought eggs, two peppers") ship in release 1, or with voice in release 3? | Release 3, together with voice, as features.md groups them. | you confirmed | [KIT-08-03](kitchen/E08-ai-intake/KIT-08-03-text-intake.md) |
| Q220 | low | How old may a receipt price be before the trip plan ignores it? | 60 days. | default | [PLN-19-04](planner/E19-trip-price-matrix/PLN-19-04-receipt-prices.md) |
| Q221 | low | Between reading a receipt and the user confirming it, should the server keep the extracted lines for 2 hours? | Yes: a 2-hour copy in the cache; no image kept. | default | [KIT-08-04](kitchen/E08-ai-intake/KIT-08-04-receipt-intake.md) |
| Q222 | medium | Should the receipt-prices event sent to offers leave out the household id? | Yes: no household, user or receipt id in the event. | default | [KIT-08-06](kitchen/E08-ai-intake/KIT-08-06-receipt-scanned-event.md) |

### M20 Depth features (R3)

| ID | Severity | Decision | Answer | Status | Tasks |
|---|---|---|---|---|---|
| Q223 | low | Which ingredient categories may "restyle this dish" change? | Seasonings, herbs, sauces and fats by kitchen category codes; mains stay. | default | [PLN-09-05](planner/E09-swaps-and-rewrites/PLN-09-05-restyle.md) |
| Q224 | medium | What limits apply when "cheapest top-up" swaps an ingredient for a cheaper one on offer? | Calories change ≤ 15%, saves ≥ 10% or €0.50, ≤ 3 swaps a week. | default | [PLN-20-01](planner/E20-cheapest-top-up/PLN-20-01-needs-and-allowed-swaps.md) |
| Q225 | low | Should the shareable savings card image be drawn on the phone or rendered as a PNG by the server? | Drawn on the phone; the server sends data only | default | [GW-12-02](gateway/E12-savings/GW-12-02-savings-card.md) |

### M21 More chains and licensed feed (R3)

| ID | Severity | Decision | Answer | Status | Tasks |
|---|---|---|---|---|---|
| Q226 | medium | Should Auchan Poland be covered with national flyers only, skipping its per-store flyers? | National flyers only in R3. | default | [OFF-12-04](offers/E12-next-chains/OFF-12-04-auchan-pl.md) |
| Q227 | low | Should Penny Germany be covered with national offers only, skipping its local per-store offers? | National offers only. | default | [OFF-12-06](offers/E12-next-chains/OFF-12-06-penny-de.md) |
| Q228 | medium | Should shops be created from the licensed offers feed where OpenStreetMap has none? | Only where OSM has none, marked as coming from the feed. | default | [OFF-13-01](offers/E13-licensed-feed/OFF-13-01-feed-client-mapping.md) |

