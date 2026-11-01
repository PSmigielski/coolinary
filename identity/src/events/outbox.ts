// The only way to write an event: in the caller's transaction, so a rolled-back change never produces one and a
// committed change never loses one. No network I/O here; the publisher (IDN-01-02) sends the rows to Redis.
import type { db } from "../db/client";
import { outboxEvents } from "../db/schema";
import { eventDefinitions, type IdentityEventPayload, type IdentityEventType } from "./types";

/** A Drizzle transaction. `db` itself is not one, so emit() can't be called outside a transaction. */
export type Tx = Parameters<Parameters<typeof db.transaction>[0]>[0];

/** A payload that doesn't match its event definition. Always a bug: the transaction rolls back (RPC: INTERNAL). */
export class InvalidEventPayloadError extends Error {
  constructor(
    readonly type: string,
    readonly problems: string[],
  ) {
    super(`invalid ${type} payload: ${problems.join("; ")}`);
    this.name = "InvalidEventPayloadError";
  }
}

/**
 * Checks a payload against its event definition and returns it. Throws InvalidEventPayloadError naming the fields
 * that fail, never their values (a rejected payload may hold exactly the personal data that must not leak).
 */
export function checkPayload<T extends IdentityEventType>(type: T, payload: unknown): IdentityEventPayload<T> {
  const result = eventDefinitions[type].payload.safeParse(payload);
  if (result.success) return result.data as IdentityEventPayload<T>;
  throw new InvalidEventPayloadError(
    type,
    result.error.issues.map((issue) => {
      const at = issue.path.length ? issue.path.join(".") : "payload";
      return issue.code === "unrecognized_keys"
        ? `${at}: unknown fields ${issue.keys.join(", ")}`
        : `${at}: ${issue.code}`;
    }),
  );
}

/**
 * Writes one event to the outbox in `tx` and returns its envelope id (UUIDv7, so id order is creation order).
 * `householdId` is the household the event is about, or null for an account without one. `occurred_at` is the
 * transaction time, shared by every event of the transaction.
 */
export async function emit<T extends IdentityEventType>(
  tx: Tx,
  type: T,
  householdId: string | null,
  payload: IdentityEventPayload<T>,
): Promise<string> {
  const checked = checkPayload(type, payload);
  const id = Bun.randomUUIDv7();
  await tx
    .insert(outboxEvents)
    .values({ id, type, version: eventDefinitions[type].version, householdId, payload: checked });
  return id;
}
