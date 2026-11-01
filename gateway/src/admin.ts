// Admin port (8081): /metrics for Prometheus, /healthz for Docker.
import { config } from "./config";
import { registry } from "./lib/metrics";
import { ping, redisCache, redisState } from "./lib/redis";

export function startAdmin() {
  return Bun.serve({
    port: config.adminPort,
    async fetch(req) {
      const { pathname } = new URL(req.url);
      if (pathname === "/metrics") {
        return new Response(await registry.metrics(), { headers: { "Content-Type": registry.contentType } });
      }
      if (pathname === "/healthz") {
        const checks = { redis_cache: await ping(redisCache), redis_state: await ping(redisState) };
        const ok = Object.values(checks).every((v) => v === "ok");
        return Response.json({ ok, checks }, { status: ok ? 200 : 503 });
      }
      return new Response("not found", { status: 404 });
    },
  });
}
