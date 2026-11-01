// Writes contracts/events/identity/<type>.v<version>.json from src/events/types.ts. Run it after changing an event
// definition (`bun run events:schemas`) and commit the files with the change.
import { mkdir } from "node:fs/promises";
import { eventSchemaDir, eventSchemaFile, eventSchemaText } from "../src/events/json-schema";
import { eventDefinitions, type IdentityEventType } from "../src/events/types";

await mkdir(eventSchemaDir, { recursive: true });
for (const type of Object.keys(eventDefinitions) as IdentityEventType[]) {
  await Bun.write(eventSchemaFile(type), eventSchemaText(type));
  console.log(`wrote contracts/events/identity/${eventSchemaFile(type).pathname.split("/").pop()}`);
}
