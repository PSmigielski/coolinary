// Identity's events: one definition per type, with its version and a strict payload schema. emit() checks every
// payload against it before writing the outbox row. IDN-01-03 adds the remaining events and their JSON Schemas.
//
// Payloads hold ids, enums, booleans and timestamps only: no e-mail address, name or other user text. JSON field
// names are snake_case. Strict objects reject any field not listed here.
import { z } from "zod";

export interface EventDefinition {
  /** Starts at 1. Raised only for a breaking payload change, never set by the caller. */
  version: number;
  payload: z.ZodType;
}

export const eventDefinitions = {
  "household.deleted": {
    version: 1,
    payload: z.strictObject({
      household_id: z.uuid(),
      reason: z.enum(["owner_deleted", "joined_other", "linked_merge", "anonymous_cleanup"]),
    }),
  },
} as const satisfies Record<string, EventDefinition>;

export type IdentityEventType = keyof typeof eventDefinitions;

export type IdentityEventPayload<T extends IdentityEventType> = z.infer<(typeof eventDefinitions)[T]["payload"]>;
