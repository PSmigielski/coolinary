// Bun -> Go gRPC smoke test: health, unary call, deadline, cancel and (optional) a
// call after an idle gap on the same channel. Usage: bun scripts/smoke.ts [idleSeconds]
import { credentials } from "@grpc/grpc-js";
import type { ResolveBarcodeResponse } from "../src/gen/coolinary/kitchen/v1/kitchen";
import { KitchenServiceClient } from "../src/gen/coolinary/kitchen/v1/kitchen";
import { channelOptions, checkHealth, clients } from "../src/grpc/clients";
import { codeName, unary } from "../src/grpc/unary";

const idleSeconds = Number(process.argv[2] ?? 0);
let failed = false;

async function step(name: string, fn: () => Promise<void>) {
  const t0 = performance.now();
  try {
    await fn();
    console.log(`ok   ${name.padEnd(28)} (${Math.round(performance.now() - t0)}ms)`);
  } catch (err) {
    failed = true;
    console.log(`FAIL ${name.padEnd(28)} ${String(err)} (${Math.round(performance.now() - t0)}ms)`);
  }
}

const resolve = (timeoutMs = 2000, signal?: AbortSignal) =>
  unary<ResolveBarcodeResponse>(
    "kitchen",
    "ResolveBarcode",
    (md, o, cb) => clients.kitchen.resolveBarcode({ ean: "9520000000017", householdId: "" }, md, o, cb),
    { timeoutMs, signal },
  );

await step("health check", async () => {
  const s = await checkHealth("kitchen");
  if (s !== "SERVING") throw new Error(s);
});
await step("unary ResolveBarcode", async () => {
  const res = await resolve();
  if (!res.found) throw new Error("example product not found");
});
await step("deadline fires (no answer)", async () => {
  // An unused address in the static range: the connection never completes, so only
  // the 300ms deadline can end the call.
  const blackhole = new KitchenServiceClient("10.77.0.99:50051", credentials.createInsecure(), channelOptions);
  const t0 = performance.now();
  try {
    await unary<ResolveBarcodeResponse>(
      "kitchen",
      "ResolveBarcode",
      (md, o, cb) => blackhole.resolveBarcode({ ean: "9520000000017", householdId: "" }, md, o, cb),
      { timeoutMs: 300 },
    );
    throw new Error("call succeeded against an unreachable address");
  } catch (err) {
    const code = (err as { code?: number }).code;
    if (code === undefined || codeName(code) !== "DEADLINE_EXCEEDED") throw err;
    if (performance.now() - t0 > 1500) throw new Error("deadline fired late");
  } finally {
    blackhole.close();
  }
});
await step("cancel is propagated", async () => {
  const ac = new AbortController();
  const p = resolve(5000, ac.signal);
  ac.abort();
  try {
    await p;
    throw new Error("call succeeded despite cancel");
  } catch (err) {
    const code = (err as { code?: number }).code;
    if (code === undefined || codeName(code) !== "CANCELLED") throw err;
  }
});
if (idleSeconds > 0) {
  console.log(`...  idle for ${idleSeconds}s`);
  await Bun.sleep(idleSeconds * 1000);
  await step("unary after idle", async () => {
    await resolve();
  });
}
process.exit(failed ? 1 : 0);
