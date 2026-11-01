// Each identity event's payload as JSON Schema 2020-12, generated from the definitions in types.ts. Consumers in other
// languages validate their test fixtures against these (Q014): never live events, which may carry fields added later.
// `bun run events:schemas` writes them to contracts/events/identity/; types.test.ts fails while a file is out of date.
import { z } from "zod";
import { eventDefinitions, type IdentityEventType } from "./types";

export const eventSchemaDir = new URL("../../../contracts/events/identity/", import.meta.url);

/** contracts/events/identity/<type>.v<version>.json */
export function eventSchemaFile(type: IdentityEventType): URL {
  return new URL(`${type}.v${eventDefinitions[type].version}.json`, eventSchemaDir);
}

export function eventJsonSchema(type: IdentityEventType): Record<string, unknown> {
  const { version, payload } = eventDefinitions[type];
  const { $schema, description, ...rest } = z.toJSONSchema(payload);
  return { $schema, title: `${type} v${version}`, description, ...rest };
}

/** The file content: two-space JSON with a final newline. */
export function eventSchemaText(type: IdentityEventType): string {
  return `${JSON.stringify(eventJsonSchema(type), null, 2)}\n`;
}
