// Admin port (8081): /metrics for Prometheus, /healthz for Docker.
import { config } from "./config";
import { sql } from "./db/client";
import { registry } from "./lib/metrics";
import { redisCache, redisState } from "./lib/redis";

async function check(fn: () => Promise<unknown>): Promise<string> {
  try {
    await fn();
    return "ok";
  } catch (err) {
    return `down: ${(err as Error).message}`;
  }
}

export function startAdmin() {
  return Bun.serve({
    port: config.adminPort,
    async fetch(req) {
      const { pathname } = new URL(req.url);
      if (pathname === "/metrics") {
        return new Response(await registry.metrics(), { headers: { "Content-Type": registry.contentType } });
      }
      if (pathname === "/healthz") {
        const checks = {
          postgres: await check(() => sql`SELECT 1`),
          redis_cache: await check(() => redisCache.send("PING", [])),
          redis_state: await check(() => redisState.send("PING", [])),
        };
        const ok = Object.values(checks).every((v) => v === "ok");
        return Response.json({ ok, checks }, { status: ok ? 200 : 503 });
      }
      return new Response("not found", { status: 404 });
    },
  });
}
