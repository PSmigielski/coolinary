# ADM-07 AI cost and system overview

**Goal.** One start page with health and waiting work, the AI cost per feature and plan against the business model's
limits (≤ €0.10 per Free household per month; ≤ 15% of a Plus household's net revenue), account and plan counts for
the phase gates, and the notifier delivery log and alerts for support.

**Flow.** admin → health checks, `GetQueueCounts` per admin service, Prometheus HTTP API, `AiAdminService.GetSpendStatus`,
`IdentityAdminService.GetAccountStats`, `NotifierAdminService` → one JSON per page.

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [ADM-07-01](ADM-07-01-system-overview.md) | System overview | R1 | M | ADM-01-05, KIT-01-07, PLN-01-06, OFF-01-06, NTF-06-02 |
| [ADM-07-03](ADM-07-03-account-stats.md) | Accounts and plans statistics | R1 | S | ADM-01-05, IDN-10-04 |
| [ADM-07-02](ADM-07-02-ai-cost-view.md) | AI cost and guardrails | R1 | M | ADM-07-03, AI-02-02, AI-02-04 |
| [ADM-07-04](ADM-07-04-delivery-log-and-alerts.md) | Delivery log, alerts and test sends | R1 | S | ADM-02-03, NTF-06 |
