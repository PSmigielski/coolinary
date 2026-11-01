// coolinary.identity.v1.IdentityService
import { status } from "@grpc/grpc-js";
import { asc, eq } from "drizzle-orm";
import { db } from "./db/client";
import { householdMembers, households } from "./db/schema";
import type { GetHouseholdRequest, GetHouseholdResponse } from "./gen/coolinary/identity/v1/identity";
import { type IdentityServiceServer, IdentityServiceService, Plan, Role } from "./gen/coolinary/identity/v1/identity";
import { type Registration, RpcError, unary } from "./grpc/server";

const UUID = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;
const roles = { owner: Role.ROLE_OWNER, member: Role.ROLE_MEMBER, profile: Role.ROLE_PROFILE } as const;

const implementation: IdentityServiceServer = {
  getHousehold: unary<GetHouseholdRequest, GetHouseholdResponse>(
    IdentityServiceService.getHousehold.path,
    async (req) => {
      if (!UUID.test(req.householdId)) throw new RpcError(status.INVALID_ARGUMENT, "household_id must be a UUID");
      const [h] = await db.select().from(households).where(eq(households.id, req.householdId)).limit(1);
      if (!h) throw new RpcError(status.NOT_FOUND, "household not found");
      const members = await db
        .select()
        .from(householdMembers)
        .where(eq(householdMembers.householdId, h.id))
        .orderBy(asc(householdMembers.createdAt));
      return {
        household: {
          id: h.id,
          name: h.name,
          ownerUserId: h.ownerUserId,
          plan: h.plan === "plus" ? Plan.PLAN_PLUS : Plan.PLAN_FREE,
          members: members.map((m) => ({
            id: m.id,
            displayName: m.displayName,
            role: roles[m.role],
            hasAccount: m.userId !== null,
          })),
        },
      };
    },
  ),
};

export const identityService: Registration = {
  name: "coolinary.identity.v1.IdentityService",
  definition: IdentityServiceService,
  implementation,
};
