import { eq } from "drizzle-orm";
import { db } from "./db/client";
import { householdMembers, households } from "./db/schema";

/** Every new account starts with its own household and an empty fridge. */
export async function createPersonalHousehold(user: { id: string; name: string }) {
  await db.transaction(async (tx) => {
    const [h] = await tx
      .insert(households)
      .values({ name: user.name ? `${user.name}'s home` : "My home", ownerUserId: user.id })
      .returning({ id: households.id });
    if (!h) throw new Error("household insert returned no row");
    await tx
      .insert(householdMembers)
      .values({ householdId: h.id, userId: user.id, displayName: user.name || "Me", role: "owner" });
  });
}

/** Claims added to the JWT the gateway verifies: which household, which plan. */
export async function householdClaims(userId: string): Promise<{ householdId?: string; plan: "free" | "plus" }> {
  const [row] = await db
    .select({ householdId: households.id, plan: households.plan })
    .from(householdMembers)
    .innerJoin(households, eq(households.id, householdMembers.householdId))
    .where(eq(householdMembers.userId, userId))
    .limit(1);
  return row ? { householdId: row.householdId, plan: row.plan } : { plan: "free" };
}
