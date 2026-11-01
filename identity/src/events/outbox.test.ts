import { afterAll, beforeAll, describe, expect, test } from "bun:test";
import { eq, inArray } from "drizzle-orm";
import type { BunSQLDatabase } from "drizzle-orm/bun-sql";
import type * as schema from "../db/schema";
import { households, outboxEvents, user } from "../db/schema";
import { hasTestDb, openTestDb } from "../db/testing";
import { checkPayload, emit, InvalidEventPayloadError, type Tx } from "./outbox";
import type { IdentityEventPayload } from "./types";

const uuidV7 = /^[0-9a-f]{8}-[0-9a-f]{4}-7[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/;

function deletedPayload(householdId: string): IdentityEventPayload<"household.deleted"> {
  return { household_id: householdId, reason: "owner_deleted" };
}

describe("checkPayload", () => {
  test("accepts a payload with exactly the listed fields", () => {
    const payload = deletedPayload(crypto.randomUUID());
    expect(checkPayload("household.deleted", payload)).toEqual(payload);
  });

  test("rejects a payload with an extra field such as email, without echoing its value", () => {
    const payload = { ...deletedPayload(crypto.randomUUID()), email: "someone@example.test" };
    expect(() => checkPayload("household.deleted", payload)).toThrow(InvalidEventPayloadError);
    try {
      checkPayload("household.deleted", payload);
    } catch (err) {
      expect((err as Error).message).toContain("email");
      expect((err as Error).message).not.toContain("someone@example.test");
    }
  });

  test("rejects a missing field, a wrong enum value and a non-uuid id", () => {
    expect(() => checkPayload("household.deleted", { household_id: crypto.randomUUID() })).toThrow(
      InvalidEventPayloadError,
    );
    expect(() => checkPayload("household.deleted", { household_id: crypto.randomUUID(), reason: "bored" })).toThrow(
      InvalidEventPayloadError,
    );
    expect(() => checkPayload("household.deleted", { household_id: "42", reason: "owner_deleted" })).toThrow(
      InvalidEventPayloadError,
    );
  });
});

// Checked by tsc (`make lint`), never run: emit() needs a transaction and a known event type.
declare const notATransaction: BunSQLDatabase<typeof schema>;
async function compileTimeChecks(tx: Tx) {
  const id = crypto.randomUUID();
  // @ts-expect-error the database is not a transaction
  await emit(notATransaction, "household.deleted", id, deletedPayload(id));
  // @ts-expect-error unknown event types don't compile
  await emit(tx, "household.renamed", id, { household_id: id });
  // @ts-expect-error fields outside the definition don't compile
  await emit(tx, "household.deleted", id, { ...deletedPayload(id), email: "someone@example.test" });
}

test("emit only compiles with a transaction, a known event type and its payload", () => {
  expect(compileTimeChecks).toBeFunction();
});

describe.skipIf(!hasTestDb)("emit (Postgres)", () => {
  let t: Awaited<ReturnType<typeof openTestDb>>;
  const eventIds: string[] = [];
  const userIds: string[] = [];

  beforeAll(async () => {
    t = await openTestDb();
  });

  afterAll(async () => {
    if (eventIds.length) await t.db.delete(outboxEvents).where(inArray(outboxEvents.id, eventIds));
    if (userIds.length) await t.db.delete(user).where(inArray(user.id, userIds)); // cascades to households
    await t.sql.close();
  });

  async function insertHousehold(tx: Tx | typeof t.db) {
    const userId = crypto.randomUUID();
    userIds.push(userId);
    const now = new Date();
    await tx
      .insert(user)
      .values({ id: userId, name: "Test", email: `${userId}@example.test`, createdAt: now, updatedAt: now });
    const [h] = await tx.insert(households).values({ name: "Test home", ownerUserId: userId }).returning();
    if (!h) throw new Error("household insert returned no row");
    return h.id;
  }

  async function outboxRows(ids: string[]) {
    return t.db.select().from(outboxEvents).where(inArray(outboxEvents.id, ids)).orderBy(outboxEvents.id);
  }

  test("a transaction that emits and then throws leaves no row", async () => {
    let id = "";
    const run = t.db.transaction(async (tx) => {
      const householdId = crypto.randomUUID();
      id = await emit(tx, "household.deleted", householdId, deletedPayload(householdId));
      throw new Error("change failed after emit");
    });
    await expect(run).rejects.toThrow("change failed after emit");
    expect(id).toMatch(uuidV7);
    expect(await outboxRows([id])).toHaveLength(0);
  });

  test("an invalid payload throws and rolls back the whole transaction", async () => {
    let householdId = "";
    const run = t.db.transaction(async (tx) => {
      householdId = await insertHousehold(tx);
      const withEmail = { ...deletedPayload(householdId), email: "someone@example.test" };
      await emit(tx, "household.deleted", householdId, withEmail as IdentityEventPayload<"household.deleted">);
    });
    await expect(run).rejects.toThrow(InvalidEventPayloadError);
    expect(await t.db.select().from(households).where(eq(households.id, householdId))).toHaveLength(0);
    expect(await t.db.select().from(outboxEvents).where(eq(outboxEvents.householdId, householdId))).toHaveLength(0);
  });

  test("two emits in a committed transaction leave two rows in call order, sharing the transaction time", async () => {
    const a = crypto.randomUUID();
    const b = crypto.randomUUID();
    const ids = await t.db.transaction(async (tx) => {
      const first = await emit(tx, "household.deleted", a, deletedPayload(a));
      const second = await emit(tx, "household.deleted", b, { household_id: b, reason: "joined_other" });
      return [first, second];
    });
    eventIds.push(...ids);

    const rows = await outboxRows(ids);
    expect(rows.map((r) => r.id)).toEqual(ids);
    expect(rows.map((r) => r.householdId)).toEqual([a, b]);
    expect(rows.map((r) => r.payload)).toEqual([deletedPayload(a), { household_id: b, reason: "joined_other" }]);
    for (const r of rows) {
      expect(r.id).toMatch(uuidV7);
      expect(r.type).toBe("household.deleted");
      expect(r.version).toBe(1); // from the definition, not the caller
      expect(r.publishedAt).toBeNull();
      expect(r.attempts).toBe(0);
      expect(r.lastError).toBeNull();
    }
    expect(rows[0]?.occurredAt).toEqual(rows[1]?.occurredAt as Date);
  });

  test("ids generated in the same millisecond still sort in call order", async () => {
    const householdId = crypto.randomUUID();
    const ids = await t.db.transaction(async (tx) => {
      const out: string[] = [];
      for (let i = 0; i < 50; i++)
        out.push(await emit(tx, "household.deleted", householdId, deletedPayload(householdId)));
      return out;
    });
    eventIds.push(...ids);
    expect((await outboxRows(ids)).map((r) => r.id)).toEqual(ids);
  });

  test("household_id has no foreign key: household.deleted is emitted in the transaction that deletes it", async () => {
    const householdId = await insertHousehold(t.db);
    const id = await t.db.transaction(async (tx) => {
      await tx.delete(households).where(eq(households.id, householdId));
      return emit(tx, "household.deleted", householdId, deletedPayload(householdId));
    });
    eventIds.push(id);

    expect(await t.db.select().from(households).where(eq(households.id, householdId))).toHaveLength(0);
    const [row] = await outboxRows([id]);
    expect(row?.householdId).toBe(householdId);
    const [fks] = await t.sql`
      SELECT count(*)::int AS n FROM pg_constraint WHERE conrelid = 'outbox_events'::regclass AND contype = 'f'`;
    expect(fks.n).toBe(0);
  });

  test("household_id may be null, for account-level events about an account without a household", async () => {
    const householdId = crypto.randomUUID();
    const id = await t.db.transaction((tx) => emit(tx, "household.deleted", null, deletedPayload(householdId)));
    eventIds.push(id);
    const [row] = await outboxRows([id]);
    expect(row?.householdId).toBeNull();
  });
});
