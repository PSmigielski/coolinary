// Writes the gateway's OpenAPI document (used by mobile/ to generate its typed client).
// Usage: bun scripts/export-openapi.ts ../contracts/openapi/gateway.json
// No service needs to run: the app is created in memory and asked for /openapi/json.
export {};

process.env.AUTH_ISSUER ??= "https://api.example.com"; // required by config, unused here

const { createApp } = await import("../src/http/app");
const out = process.argv[2] ?? "openapi.json";
const res = await createApp().handle(new Request("http://localhost/openapi/json"));
if (!res.ok) throw new Error(`openapi export failed: ${res.status}`);
await Bun.write(out, `${JSON.stringify(await res.json(), null, 2)}\n`);
console.log(`wrote ${out}`);
process.exit(0);
