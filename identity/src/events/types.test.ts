import { describe, expect, test } from "bun:test";
import Ajv2020 from "ajv/dist/2020";
import { eventSchemaFile, eventSchemaText } from "./json-schema";
import { checkPayload, emit, InvalidEventPayloadError, type Tx } from "./outbox";
import { eventDefinitions, type IdentityEventPayload, type IdentityEventType } from "./types";

// IDN-01-03 rule 2: the eight events and every field of their v1 payloads. A field is a household or member id
// (UUID), a Better Auth user id (text), a boolean, one of a list of values, a list of member ids or (snapshot only)
// a list of members. "?" marks the fields that may be null.
type Field = "uuid" | "uuid?" | "user_id" | "user_id?" | "boolean" | "uuid[]" | "members" | readonly string[];

const plans = ["free", "plus"] as const;
const roles = ["owner", "member", "profile"] as const;

const fields: { [T in IdentityEventType]: { [K in keyof IdentityEventPayload<T>]-?: Field } } = {
  "household.created": { household_id: "uuid", owner_user_id: "user_id", owner_member_id: "uuid" },
  "household.member_added": {
    household_id: "uuid",
    member_id: "uuid",
    user_id: "user_id?",
    has_account: "boolean",
    role: roles,
    brought_fridge_from: "uuid?",
    previous_member_id: "uuid?",
  },
  "household.member_removed": {
    household_id: "uuid",
    member_id: "uuid",
    user_id: "user_id?",
    reason: ["left", "removed", "profile_deleted", "plan_ended"],
    new_household_id: "uuid?",
  },
  "household.plan_changed": {
    household_id: "uuid",
    plan: plans,
    previous_plan: plans,
    source: ["default", "manual", "store"],
  },
  "household.deleted": {
    household_id: "uuid",
    reason: ["owner_deleted", "joined_other", "linked_merge", "anonymous_cleanup"],
  },
  "account.linked": {
    anonymous_user_id: "user_id",
    user_id: "user_id",
    household_id: "uuid",
    anonymous_household_id: "uuid",
    merged: "boolean",
  },
  "account.deleted": { user_id: "user_id", member_ids: "uuid[]" },
  "household.snapshot": { household_id: "uuid", plan: plans, members: "members" },
};

type SnapshotMember = IdentityEventPayload<"household.snapshot">["members"][number];
const snapshotMemberFields: { [K in keyof SnapshotMember]-?: Field } = {
  member_id: "uuid",
  user_id: "user_id?",
  has_account: "boolean",
  role: roles,
};

const eventTypes = Object.keys(fields) as IdentityEventType[];

// Made-up ids. Better Auth user ids are 32 random letters and digits.
const household = "69220c46-f782-4cf9-b1e9-007cc895e0f9";
const otherHousehold = "944f0589-901c-46f2-99cd-d27b41985853";
const ownerMember = "07cdc834-8deb-4937-bafa-b4394c952d0b";
const member = "a4bba725-063d-41ec-9220-96582fd724e0";
const previousMember = "2c7e3f56-5129-428a-b5ff-a58c0707c551";
const profileMember = "cc53d33c-2e9c-416d-a9de-666ebbce5764";
const user = "Xq7LmP2vT9sKcR4wYb8NdF3hJ6gA1eZu";
const otherUser = "Hn4sW8kPz2QmV6tR1cY9bL3xJ7fD5aGe";

// One valid payload per event. Typed against the definitions, so tsc proves each is a valid emit() payload.
const samples: { [T in IdentityEventType]: IdentityEventPayload<T> } = {
  "household.created": { household_id: household, owner_user_id: user, owner_member_id: ownerMember },
  "household.member_added": {
    household_id: household,
    member_id: member,
    user_id: otherUser,
    has_account: true,
    role: "member",
    brought_fridge_from: otherHousehold,
    previous_member_id: previousMember,
  },
  "household.member_removed": {
    household_id: household,
    member_id: member,
    user_id: otherUser,
    reason: "left",
    new_household_id: otherHousehold,
  },
  "household.plan_changed": { household_id: household, plan: "plus", previous_plan: "free", source: "manual" },
  "household.deleted": { household_id: household, reason: "owner_deleted" },
  "account.linked": {
    anonymous_user_id: otherUser,
    user_id: user,
    household_id: household,
    anonymous_household_id: otherHousehold,
    merged: true,
  },
  "account.deleted": { user_id: user, member_ids: [ownerMember] },
  "household.snapshot": {
    household_id: household,
    plan: "free",
    members: [
      { member_id: ownerMember, user_id: user, has_account: true, role: "owner" },
      { member_id: profileMember, user_id: null, has_account: false, role: "profile" },
    ],
  },
};

const snapshotMember = samples["household.snapshot"].members[0] as SnapshotMember;

// Values each kind of field must accept and reject, beyond the sample's own value.
function casesFor(field: Field): { accepts: unknown[]; rejects: unknown[] } {
  if (Array.isArray(field)) return { accepts: [...field], rejects: ["unknown", "", null] };
  switch (field) {
    case "uuid":
      return { accepts: [], rejects: ["42", "", user, null] };
    case "uuid?":
      return { accepts: [null], rejects: ["42", "", user] };
    case "user_id":
      return { accepts: [user], rejects: ["someone@example.test", "Test User", "", null] };
    case "user_id?":
      return { accepts: [null], rejects: ["someone@example.test", "Test User", ""] };
    case "boolean":
      return { accepts: [true, false], rejects: ["true", 1, null] };
    case "uuid[]":
      return { accepts: [[], [ownerMember, member]], rejects: [["42"], [user], ownerMember, null] };
    case "members":
      return { accepts: [[]], rejects: [snapshotMember, null] };
  }
  throw new Error(`unknown field kind ${String(field)}`);
}

/** One row per field and value: the values a table's fields must accept, or must reject. */
function valueCases(table: Record<string, Field>, which: "accepts" | "rejects") {
  return Object.entries(table).flatMap(([name, field]) => casesFor(field)[which].map((value) => ({ name, value })));
}

function withField(payload: object, name: string, value: unknown): unknown {
  return { ...payload, [name]: value };
}

function withoutField(payload: object, name: string): unknown {
  const { [name as keyof typeof payload]: _dropped, ...rest } = payload;
  return rest;
}

describe("identity event definitions", () => {
  test("define exactly the eight identity events, each at version 1", () => {
    const versions = Object.fromEntries(Object.entries(eventDefinitions).map(([type, d]) => [type, d.version]));
    expect(versions).toEqual(Object.fromEntries(eventTypes.map((type) => [type, 1])));
  });

  describe.each(eventTypes)("%s", (type) => {
    test("accepts a payload with every listed field", () => {
      expect(checkPayload(type, samples[type])).toEqual(samples[type]);
    });

    test.each(Object.keys(fields[type]))("rejects a payload without %s", (name) => {
      expect(() => checkPayload(type, withoutField(samples[type], name))).toThrow(InvalidEventPayloadError);
    });

    test("rejects a field outside the list, naming it but not its value", () => {
      const payload = withField(samples[type], "email", "someone@example.test");
      expect(() => checkPayload(type, payload)).toThrow(/email/);
      expect(() => checkPayload(type, payload)).not.toThrow(/someone@example\.test/);
    });

    test.each(valueCases(fields[type], "accepts"))("accepts $name = $value", ({ name, value }) => {
      const payload = withField(samples[type], name, value);
      expect(checkPayload(type, payload)).toEqual(payload as IdentityEventPayload<typeof type>);
    });

    test.each(valueCases(fields[type], "rejects"))("rejects $name = $value", ({ name, value }) => {
      expect(() => checkPayload(type, withField(samples[type], name, value))).toThrow(InvalidEventPayloadError);
    });
  });

  describe("household.snapshot members", () => {
    const snapshot = samples["household.snapshot"];
    const withMember = (m: unknown) => ({ ...snapshot, members: [m] });

    test.each(Object.keys(snapshotMemberFields))("rejects a member without %s", (name) => {
      const payload = withMember(withoutField(snapshotMember, name));
      expect(() => checkPayload("household.snapshot", payload)).toThrow(InvalidEventPayloadError);
    });

    test("rejects a member field outside the list, such as a display name", () => {
      const payload = withMember({ ...snapshotMember, display_name: "Test User" });
      expect(() => checkPayload("household.snapshot", payload)).toThrow(/display_name/);
    });

    test.each(valueCases(snapshotMemberFields, "accepts"))("accepts member $name = $value", ({ name, value }) => {
      const payload = withMember(withField(snapshotMember, name, value));
      expect(checkPayload("household.snapshot", payload)).toEqual(payload as typeof snapshot);
    });

    test.each(valueCases(snapshotMemberFields, "rejects"))("rejects member $name = $value", ({ name, value }) => {
      const payload = withMember(withField(snapshotMember, name, value));
      expect(() => checkPayload("household.snapshot", payload)).toThrow(InvalidEventPayloadError);
    });
  });
});

// The JSON Schemas consumers in other languages validate their test fixtures against (Q014).
type JsonSchema = {
  $schema?: string;
  properties?: Record<string, JsonSchema>;
  required?: string[];
  additionalProperties?: boolean;
  items?: JsonSchema;
};

const schemaDir = new URL("../../../contracts/events/identity/", import.meta.url);

async function committedSchema(type: IdentityEventType): Promise<JsonSchema | null> {
  const file = Bun.file(new URL(`${type}.v${eventDefinitions[type].version}.json`, schemaDir));
  return (await file.exists()) ? ((await file.json()) as JsonSchema) : null;
}

function sortedKeys(schema: JsonSchema | null | undefined) {
  return { properties: Object.keys(schema?.properties ?? {}).sort(), required: [...(schema?.required ?? [])].sort() };
}

describe("JSON Schemas in contracts/events/identity", () => {
  // Formats are annotations here; the patterns zod writes are what check the ids.
  const ajv = new Ajv2020({ strict: true, allErrors: true, validateFormats: false });

  describe.each(eventTypes)("%s", (type) => {
    test("is committed as <type>.v1.json in JSON Schema 2020-12", async () => {
      expect((await committedSchema(type))?.$schema).toBe("https://json-schema.org/draft/2020-12/schema");
    });

    test("lists exactly the payload's fields, all required, and no others", async () => {
      const schema = await committedSchema(type);
      const listed = Object.keys(fields[type]).sort();
      expect(sortedKeys(schema)).toEqual({ properties: listed, required: listed });
      expect(schema?.additionalProperties).toBe(false);
    });

    test("accepts the sample payload, and rejects it with a field missing or added", async () => {
      const validate = ajv.compile((await committedSchema(type)) ?? {});
      const [firstField] = Object.keys(fields[type]);
      expect(validate(samples[type])).toBe(true);
      expect(validate(withoutField(samples[type], firstField as string))).toBe(false);
      expect(validate(withField(samples[type], "email", "someone@example.test"))).toBe(false);
    });
  });

  test.each(eventTypes)("%s: the committed file is what `bun run events:schemas` writes", async (type) => {
    const file = Bun.file(eventSchemaFile(type));
    expect((await file.exists()) ? await file.text() : "").toBe(eventSchemaText(type));
  });

  test("household.snapshot lists exactly the member fields, all required, and no others", async () => {
    const members = (await committedSchema("household.snapshot"))?.properties?.members?.items;
    const listed = Object.keys(snapshotMemberFields).sort();
    expect(sortedKeys(members)).toEqual({ properties: listed, required: listed });
    expect(members?.additionalProperties).toBe(false);
  });
});

// identity/EVENTS.md: one "### `<type>`" section per event, each with a ```json example payload.
async function documentedExamples(): Promise<Map<string, unknown>> {
  const file = Bun.file(new URL("../../EVENTS.md", import.meta.url));
  const doc = (await file.exists()) ? await file.text() : "";
  const examples = new Map<string, unknown>();
  for (const section of doc.split(/^### /m).slice(1)) {
    const type = section.match(/^`([^`]+)`/)?.[1];
    const json = section.match(/```json\n([\s\S]*?)```/)?.[1];
    if (type) examples.set(type, json === undefined ? undefined : JSON.parse(json));
  }
  return examples;
}

describe("identity/EVENTS.md", () => {
  test("has one section per event, and no others", async () => {
    expect([...(await documentedExamples()).keys()].sort()).toEqual([...eventTypes].sort());
  });

  test.each(eventTypes)("%s: the example payload is one emit() accepts", async (type) => {
    const example = (await documentedExamples()).get(type);
    expect(checkPayload(type, example)).toEqual(example as IdentityEventPayload<typeof type>);
  });
});

// Checked by tsc (`make lint`), never run: every event's payload type-checks at an emit() call site, and payloads
// that break the definitions don't compile.
async function compileTimeChecks(tx: Tx) {
  await emit(tx, "household.created", household, samples["household.created"]);
  await emit(tx, "household.member_added", household, samples["household.member_added"]);
  await emit(tx, "household.member_removed", household, samples["household.member_removed"]);
  await emit(tx, "household.plan_changed", household, samples["household.plan_changed"]);
  await emit(tx, "household.deleted", household, samples["household.deleted"]);
  await emit(tx, "account.linked", household, samples["account.linked"]);
  await emit(tx, "account.deleted", null, samples["account.deleted"]);
  await emit(tx, "household.snapshot", household, samples["household.snapshot"]);
  // @ts-expect-error owner_user_id is never null
  await emit(tx, "household.created", household, { ...samples["household.created"], owner_user_id: null });
  // @ts-expect-error role is owner, member or profile
  await emit(tx, "household.member_added", household, { ...samples["household.member_added"], role: "admin" });
  // @ts-expect-error has_account is a boolean
  await emit(tx, "household.member_added", household, { ...samples["household.member_added"], has_account: "yes" });
  // @ts-expect-error member_ids is a list
  await emit(tx, "account.deleted", null, { user_id: user, member_ids: ownerMember });
  await emit(tx, "household.snapshot", household, {
    ...samples["household.snapshot"],
    // @ts-expect-error snapshot members carry no display name
    members: [{ ...snapshotMember, display_name: "Test User" }],
  });
}

test("every event's payload compiles at an emit() call site", () => {
  expect(compileTimeChecks).toBeFunction();
});
