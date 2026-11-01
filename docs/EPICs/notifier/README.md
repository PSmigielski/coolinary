# notifier — epics and tasks

Bun service: e-mail (SMTP), push (Expo Push API), devices and notification preferences, founder alerts, delivery
log. Calls no one. Today it serves only a raw `SendEmail`; these tasks make it template-based, private by default
and able to push. Conventions: `docs/EPICs/README.md`, `docs/service-conventions.md`.

## Epics

| Epic | Title | Rel. | Tasks |
|---|---|---|---|
| [NTF-01](E01-events-consumer/README.md) | Events consumer (Bun) | R1 | 4 |
| [NTF-02](E02-email-templates/README.md) | E-mail templates and languages | R1 | 5 |
| [NTF-03](E03-push-delivery/README.md) | Push delivery (Expo Push API) | R1 | 4 |
| [NTF-04](E04-devices-preferences/README.md) | Devices and notification preferences | R1 | 3 |
| [NTF-05](E05-founder-alerts/README.md) | Founder alerts | R1 | 2 |
| [NTF-06](E06-admin-rpcs/README.md) | Admin RPCs (delivery log) | R1 | 2 |

## Build order

Sign-in e-mail first (identity's e-mail code sign-in waits on it), then events, devices, push, cleanup, alerts, admin.

| # | Task | Title | Rel. | Size | Depends on |
|---|---|---|---|---|---|
| 1 | NTF-02-01 | Template registry and renderer | R1 | M | — |
| 2 | NTF-02-02 | Delivery log without personal data | R1 | M | — |
| 3 | NTF-02-03 | SendEmail with templates, languages and idempotency | R1 | M | NTF-02-01, NTF-02-02 |
| 4 | NTF-02-04 | Sign-in code e-mail | R1 | S | NTF-02-03 |
| 5 | NTF-05-02 | Founder alert e-mail and recipients | R1 | S | NTF-02-01 |
| 6 | NTF-05-01 | SendAlert with deduplication and throttling | R1 | M | NTF-05-02 |
| 7 | NTF-01-01 | Redis Streams consumer (Bun) | R1 | M | IDN-01 |
| 8 | NTF-01-02 | Household users mirror | R1 | M | NTF-01-01 |
| 9 | NTF-04-01 | Register and unregister devices | R1 | M | — |
| 10 | NTF-03-01 | Push templates and categories | R1 | M | NTF-02-01 |
| 11 | NTF-04-02 | Notification preferences and quiet hours | R1 | M | NTF-03-01 |
| 12 | NTF-03-02 | SendPush and the push queue | R1 | L | NTF-01-02, NTF-02-02, NTF-03-01, NTF-04-01, NTF-04-02 |
| 13 | NTF-03-03 | Expo delivery worker | R1 | M | NTF-03-02 |
| 14 | NTF-03-04 | Push receipts and cleanup | R1 | S | NTF-03-03 |
| 15 | NTF-01-03 | Erase a deleted account's data | R1 | S | NTF-01-02, NTF-03-02, NTF-04-02 |
| 16 | NTF-01-04 | Move an anonymous account's devices to the signed-in account | R1 | S | NTF-01-02, NTF-03-02, NTF-04-02 |
| 17 | NTF-04-03 | Export a user's notifier data | R1 | S | NTF-04-01, NTF-04-02 |
| 18 | NTF-02-05 | Account deletion confirmation e-mail | R1 | S | NTF-02-03 |
| 19 | NTF-06-01 | Delivery log search and statistics | R1 | M | NTF-02-02 |
| 20 | NTF-06-02 | Alerts list, acknowledge and test sends | R1 | S | NTF-06-01, NTF-05-01 |

20 tasks: 1 L, 11 M, 8 S.

## RPCs after these tasks

`NotifierService`: `SendEmail` (templated), `SendPush`, `SendAlert`, `RegisterDevice`, `UnregisterDevice`,
`GetPreferences`, `SetPreferences`, `ExportUserData`.
`NotifierAdminService`: `ListDeliveries`, `GetDeliveryStats`, `ListAlerts`, `AcknowledgeAlert`, `SendTestEmail`,
`SendTestPush`.

## Events consumed

From `events.identity`: `household.created`, `household.member_added`, `household.member_removed`,
`household.deleted` (mirror, NTF-01-02), `account.deleted` (NTF-01-03), `account.linked` (NTF-01-04).
`household.plan_changed` is not consumed: senders enforce Free/Plus themselves, reading the plan from identity's
`GetEntitlements` (Q004).

## Decisions

Settled with the founder on 1–2 November 2026 (full list: [DECISIONS.md](../DECISIONS.md)). Each task states its rule as `Decided (Qnnn)`; decisions owned by another service are listed with that service's task.

| Q | Decision | Tasks |
|---|---|---|
| Q004 | Ask identity each time — kitchen, planner and offers read the plan from identity's GetEntitlements on each request that needs it, and per household in background jobs; no local plan copy. |  (owner) |
| Q029 | 90 days. | NTF-02-02 |
| Q030 | Yes, keep the extra cap. | NTF-02-03 |
| Q037 | Yes, code in the subject. | NTF-02-04 |
| Q038 | "du". | NTF-02-04 |
| Q039 | Yes. | NTF-02-05 |
| Q040 | Critical every 1 hour, warning every 6 hours. | NTF-05-01 |
| Q042 | No: accept the risk and rely on the alert. | NTF-01-01 |
| Q130 | Yes; add installation_id to notifier's UnregisterDevice | NTF-04-01 |
| Q131 | 10; an 11th replaces the least recently seen. | NTF-04-01 |
| Q132 | 90 days. | NTF-04-01 |
| Q133 | 21:00–08:00. | NTF-04-02 |
| Q188 | Not needed: after a downgrade members stay (Q187), so no `sharing_ends` push. | NTF-03, NTF-04 |

## Gates

- NTF-02-01: Polish and German texts checked by a native speaker before release 1.

## New config and secrets

`DELIVERY_HASH_KEY` (secret `notifier_hash_key`), `EXPO_ACCESS_TOKEN` (secret `expo_access_token`),
`EXPO_PUSH_URL`, `SUPPORT_EMAIL`, `FOUNDER_ALERT_EMAILS`, `ADMIN_PUBLIC_URL`.
