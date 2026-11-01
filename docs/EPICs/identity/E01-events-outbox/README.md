# IDN-01 Events: outbox and Redis Streams publisher (Bun)

Identity is the source of every household and account change other services must follow (members joining or
leaving, plan changes, deletions). This epic makes those changes reach `events.identity` on redis-state exactly as
committed: written in the same transaction as the change (outbox), published in order, at least once.

This is the Bun reference pattern. notifier and admin copy the consumer side (NTF-01); the Go services build their
own outbox in KIT-01, PLN-01 and OFF-01 with the same envelope.

## Flow

1. A business function changes rows and calls `emit(tx, ...)` in the same Drizzle transaction.
2. The publisher picks unpublished rows in id order and `XADD`s them to `events.identity`.
3. Consumers read with consumer groups and drop duplicates by envelope `id`.
4. A consumer that fell behind the 7-day stream window asks admin to republish snapshots (IDN-01-04).

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| IDN-01-01 | Outbox table and `emit()` helper | R1 | S | — |
| IDN-01-02 | Publisher to `events.identity` | R1 | M | IDN-01-01 |
| IDN-01-03 | Identity event types and schemas | R1 | S | IDN-01-01 |
| IDN-01-04 | Household snapshot republish | R1 | M | IDN-01-02, IDN-01-03, IDN-10-01 |
