# IDN-01-01 Outbox table and `emit()` helper

| | |
|---|---|
| Service | identity |
| Epic | [IDN-01 Events: outbox and publisher](README.md) |
| Release | R1 |
| Size | S |
| Depends on | — |
| Unblocks | IDN-01-02, IDN-01-03, IDN-04-02, every task that emits an event |

## Goal

Any identity change that other services must learn about writes its event in the same database transaction as the
change. A rolled-back change never produces an event, and a committed change never loses one.

## Business rules

1. Events are written only through `emit(tx, type, householdId, payload)` in `src/events/outbox.ts`. Its first
   parameter is the Drizzle transaction type, so calling it without a transaction does not compile.
2. `emit` generates the envelope id as a UUIDv7 in code. Ordering rows by `id` therefore equals creation order;
   the publisher (IDN-01-02) relies on this.
3. `occurred_at` is the database transaction time (`now()`), so all events of one transaction share it.
4. `household_id` is the household the event is about. Account-level events (`account.deleted`,
   `account.linked`) carry the account's household at the time of the event, or null if it has none.
5. `version` starts at 1 per event type and comes from the event type definition (IDN-01-03), never from the caller.
6. Before insert, `payload` is checked against the event's type definition. A payload that fails is a bug: `emit`
   throws and the whole transaction rolls back (the RPC returns `INTERNAL`).
7. Payloads hold ids, enums, booleans and timestamps only. No e-mail address, display name, household name or any
   user text (privacy rule in docs/EPICs/README.md). IDN-01-03 lists the allowed fields.
8. `emit` does no network I/O. Publishing is the publisher's job.

## Contract

No RPC. Internal TypeScript API:

```ts
export async function emit<T extends IdentityEventType>(
  tx: Tx, type: T, householdId: string | null, payload: IdentityEventPayload<T>,
): Promise<string /* envelope id */>;
```

## Data

`identity/src/db/schema.ts` → `bun run db:generate`:

```ts
export const outboxEvents = pgTable("outbox_events", {
  id: uuid("id").primaryKey(),                       // UUIDv7, set in code
  type: text("type").notNull(),                      // e.g. "household.member_added"
  version: integer("version").notNull(),
  householdId: uuid("household_id"),                 // no FK: the household may be deleted by the same tx
  payload: jsonb("payload").notNull(),
  occurredAt: timestamp("occurred_at", { withTimezone: true }).notNull().defaultNow(),
  publishedAt: timestamp("published_at", { withTimezone: true }),
  attempts: integer("attempts").notNull().default(0),
  lastError: text("last_error"),
}, (t) => [index("outbox_unpublished_idx").on(t.id).where(sql`published_at IS NULL`)]);
```

## Errors

| Case | Result |
|---|---|
| Payload doesn't match the type | throws; transaction rolls back; RPC returns `INTERNAL` |
| Unknown event type | compile error (typed union) |

## Acceptance criteria

- [ ] A test transaction that calls `emit` and then throws leaves no row in `outbox_events`.
- [ ] A committed transaction with two `emit` calls leaves two rows whose ids sort in call order.
- [ ] A payload containing an extra field (e.g. `email`) is rejected.
- [ ] `household_id` has no foreign key, so `household.deleted` can be emitted in the transaction that deletes it.

## Notes

UUIDv7: Bun has `Bun.randomUUIDv7()`. The partial index keeps the publisher's query cheap as published rows pile up
until cleanup (IDN-01-02).
