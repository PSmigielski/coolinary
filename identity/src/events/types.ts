// Identity's events: one definition per type, with its version and a strict payload schema. emit() checks every
// payload against it before writing the outbox row, and json-schema.ts exports each one as JSON Schema to
// contracts/events/identity/ for consumers in other languages (`bun run events:schemas`). identity/EVENTS.md
// describes every event with an example.
//
// Payloads hold ids, enums and booleans only: no e-mail address, name or other user text. JSON field names are
// snake_case. Ids are UUIDs, except user ids, which are Better Auth text ids. A field that may be null is always
// present. Strict objects reject any field not listed here.
//
// Rules for consumers (kitchen, planner, offers, notifier):
// - Drop duplicates by envelope `id`: delivery is at least once.
// - `household.member_added` and `household.snapshot` are upserts by `member_id` / `household_id`: applying the same
//   one twice, or a member that already exists, is not an error.
// - Within `events.identity`, events about one household arrive in the order identity committed them. Example:
//   `member_added` with `brought_fridge_from = X` always comes before `household.deleted` for X.
// - Acknowledge and ignore unknown event types. Ignore unknown payload fields.
//
// Versioning: adding an optional field keeps the version. Removing, renaming or changing the meaning of a field
// creates version 2; identity emits v2 only after every consumer accepts both versions (one release).
import { z } from "zod";

export interface EventDefinition {
  /** Starts at 1. Raised only for a breaking payload change, never set by the caller. */
  version: number;
  payload: z.ZodType;
}

/** A Better Auth user id: random letters and digits, never an e-mail address or a name. */
const userId = z.string().regex(/^[A-Za-z0-9_-]{1,64}$/);
const householdId = z.uuid().describe("The household the event is about.");
const memberId = z.uuid().describe("The member: one person in one household, with or without an account.");
const plan = z.enum(["free", "plus"]);
const role = z
  .enum(["owner", "member", "profile"])
  .describe("owner and member have an account; a profile is a person without one (a child, someone without a phone).");
const hasAccount = z.boolean().describe("True for owner and member, false for a profile.");
const memberUserId = userId.nullable().describe("The member's account, or null for a profile.");

export const eventDefinitions = {
  "household.created": {
    version: 1,
    payload: z
      .strictObject({
        household_id: householdId,
        owner_user_id: userId.describe("The account that owns the new household."),
        owner_member_id: memberId.describe("The owner's member id in the new household."),
      })
      .describe("A household was created. Its owner's household.member_added follows in the same transaction."),
  },
  "household.member_added": {
    version: 1,
    payload: z
      .strictObject({
        household_id: householdId,
        member_id: memberId,
        user_id: memberUserId,
        has_account: hasAccount,
        role,
        brought_fridge_from: z
          .uuid()
          .nullable()
          .describe("The household whose stock moves into this one, or null. Its household.deleted follows."),
        previous_member_id: z
          .uuid()
          .nullable()
          .describe(
            "The same person's member id in the household they came from, to carry over their settings, or null.",
          ),
      })
      .describe("A person became a member of a household. An upsert by member_id."),
  },
  "household.member_removed": {
    version: 1,
    payload: z
      .strictObject({
        household_id: householdId,
        member_id: memberId,
        user_id: memberUserId,
        reason: z
          .enum(["left", "removed", "profile_deleted", "plan_ended"])
          .describe("plan_ended is reserved: identity doesn't emit it (members stay after a downgrade, Q187)."),
        new_household_id: z
          .uuid()
          .nullable()
          .describe("The household the person moves to, or null if they move nowhere (a deleted profile)."),
      })
      .describe("A person stopped being a member of a household."),
  },
  "household.plan_changed": {
    version: 1,
    payload: z
      .strictObject({
        household_id: householdId,
        plan: plan.describe("The plan from now on."),
        previous_plan: plan.describe("The plan before this change."),
        source: z
          .enum(["default", "manual", "store"])
          .describe(
            "default: Free with nothing bought; manual: set by an admin; store: reserved for store purchases (IDN-05-05).",
          ),
      })
      .describe("A household's plan changed. Read the current plan from identity (GetEntitlements), Q004."),
  },
  "household.deleted": {
    version: 1,
    payload: z
      .strictObject({
        household_id: householdId,
        reason: z.enum(["owner_deleted", "joined_other", "linked_merge", "anonymous_cleanup"]),
      })
      .describe("A household and everything in it is gone. Delete its data."),
  },
  "account.linked": {
    version: 1,
    payload: z
      .strictObject({
        anonymous_user_id: userId.describe("The anonymous account that signed in."),
        user_id: userId.describe("The account it became or was merged into."),
        household_id: householdId.describe("The account's household after linking."),
        anonymous_household_id: z.uuid().describe("The anonymous account's household before linking."),
        merged: z
          .boolean()
          .describe(
            "False: the anonymous account was upgraded in place (household_id = anonymous_household_id). True: it signed in to an existing account, its household's stock moved into that account's household and the anonymous household was deleted.",
          ),
      })
      .describe("An anonymous account signed in with an e-mail code."),
  },
  "account.deleted": {
    version: 1,
    payload: z
      .strictObject({
        user_id: userId,
        member_ids: z.array(z.uuid()).describe("The account's member ids; empty if it had no household."),
      })
      .describe("An account was deleted. Delete what belongs to the account and these members."),
  },
  "household.snapshot": {
    version: 1,
    payload: z
      .strictObject({
        household_id: householdId,
        plan,
        members: z
          .array(z.strictObject({ member_id: memberId, user_id: memberUserId, has_account: hasAccount, role }))
          .describe("Every member: accounts and profiles."),
      })
      .describe("The household's current state, republished on request (IDN-01-04). An upsert by household_id."),
  },
} as const satisfies Record<string, EventDefinition>;

export type IdentityEventType = keyof typeof eventDefinitions;

export type IdentityEventPayload<T extends IdentityEventType> = z.infer<(typeof eventDefinitions)[T]["payload"]>;
