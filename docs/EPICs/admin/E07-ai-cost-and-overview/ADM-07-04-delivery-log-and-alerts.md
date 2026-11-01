# ADM-07-04 Delivery log, alerts and test sends

| | |
|---|---|
| Service | admin |
| Epic | [ADM-07 AI cost and system overview](README.md) |
| Release | R1 |
| Size | S |
| Status | todo |
| Depends on | ADM-02-03; NTF-06 (`ListDeliveries`, `GetDeliveryStats`, `ListAlerts`, `AcknowledgeAlert`, `SendTestEmail`, `SendTestPush`) |
| Unblocks | — |

## Goal

Support can answer "I never got my sign-in code", the founder can work through founder alerts, and templates can be
checked by sending them to oneself.

## Business rules

1. Delivery log search: filters as NTF-06-01 (channel, status, template, category, user id, e-mail, date range ≤ 31
   days). A search by e-mail is audited as a `read` without the address (ADM-02-03 rule 5).
2. Delivery stats per day for up to 90 days.
3. Alerts: list open (default) or all of the last 30 days; acknowledge one (audited `notifier.alert.ack`).
4. Test e-mail: admin sends `x-admin-email` = the e-mail in the admin's JWT and `to` = the same address; template and
   language chosen from notifier's list. Test push goes to the admin's own devices. At most 10 test sends per admin
   per hour (notifier enforces). Audited `notifier.test.send`.
5. Role `admin` only.

## Contract

| Method | Path | Body / response |
|---|---|---|
| GET | `/admin/v1/notifier/deliveries?channel=&status=&template_id=&category=&user_id=&email=&from=&to=&page_token=` | `{ items, next_page_token }` |
| GET | `/admin/v1/notifier/delivery-stats?from=&to=` | `{ stats }` |
| GET | `/admin/v1/notifier/alerts?state=open\|all` | `{ items }` |
| POST | `/admin/v1/notifier/alerts/{id}/acknowledge` | alert |
| POST | `/admin/v1/notifier/test-email` | `{ template_id, language }` → `{ delivery_id }` |
| POST | `/admin/v1/notifier/test-push` | `{ template_id, language }` → `{ queued }` |

## Data

None in admin.

## Errors

| Case | HTTP |
|---|---|
| Range too long | 400 |
| No device for test push | 404 |
| Too many tests | 429 |

## Acceptance criteria

- [ ] Searching the delivery log by an e-mail shows its sign-in code deliveries and writes a `read` audit row
      without the address.
- [ ] A test e-mail arrives only at the admin's own address.
