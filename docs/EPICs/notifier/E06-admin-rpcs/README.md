# NTF-06 Admin RPCs (delivery log)

**Goal.** The admin service can search the delivery log (also by e-mail address, hashed on the fly), see delivery
statistics, list and acknowledge alerts, and send test messages to the admin's own address or phone.

**Flow.** admin (with `x-admin-user-id`) → `NotifierAdminService` → read-only queries, or a test send marked
`test = true`.

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [NTF-06-01](NTF-06-01-delivery-log-search.md) | Delivery log search and statistics | R1 | M | NTF-02-02 |
| [NTF-06-02](NTF-06-02-alerts-and-test-sends.md) | Alerts list, acknowledge and test sends | R1 | S | NTF-06-01, NTF-05-01, NTF-02-03, NTF-03-02 |
