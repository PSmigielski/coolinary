# NTF-01 Events consumer (Bun)

**Goal.** notifier follows identity's events through Redis Streams (decision 4) to keep a mirror of who belongs to
which household and to erase or move a person's data when their account is deleted or upgraded. It is the first
Bun consumer; other Bun services copy its pattern.

**Flow.** identity outbox → `events.identity` stream → consumer group `notifier` → idempotent handler in one
Postgres transaction (`processed_events`) → `XACK`. Failures are retried; after 10 deliveries an entry goes to
`dead_events` and the founder is alerted.

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [NTF-01-01](NTF-01-01-streams-consumer.md) | Redis Streams consumer (Bun) | R1 | M | IDN-01 |
| [NTF-01-02](NTF-01-02-household-users-mirror.md) | Household users mirror | R1 | M | NTF-01-01, IDN-04, IDN-08/09 |
| [NTF-01-03](NTF-01-03-account-deleted.md) | Erase a deleted account's data | R1 | S | NTF-01-01, NTF-01-02, NTF-02-02, NTF-03-02, NTF-04-01, NTF-04-02, IDN-09 |
| [NTF-01-04](NTF-01-04-account-linked.md) | Move an anonymous account's devices to the signed-in account | R1 | S | NTF-01-01, NTF-01-02, NTF-03-02, NTF-04-01, NTF-04-02, IDN-03 |
