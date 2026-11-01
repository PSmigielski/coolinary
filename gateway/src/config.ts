// Settings from environment variables. Any X can be given as X_FILE (Docker secrets).
import { readFileSync } from "node:fs";

function env(key: string, fallback?: string): string {
  const file = process.env[`${key}_FILE`];
  if (file) return readFileSync(file, "utf8").trim();
  const value = process.env[key] ?? fallback;
  if (value === undefined) throw new Error(`${key} is required`);
  return value;
}

export const config = {
  serviceName: env("SERVICE_NAME", "gateway"),
  httpPort: Number(env("HTTP_PORT", "3000")),
  adminPort: Number(env("ADMIN_PORT", "8081")),
  logLevel: env("LOG_LEVEL", "info"),
  redisCacheUrl: env("REDIS_CACHE_URL", "redis://redis-cache:6379/0"),
  redisStateUrl: env("REDIS_STATE_URL", "redis://redis-state:6379/0"),
  identityHttpUrl: env("IDENTITY_HTTP_URL", "http://identity:3000"),
  // JWTs from identity must carry this issuer and audience (the public API URL).
  authIssuer: env("AUTH_ISSUER"),
  grpc: {
    identity: env("IDENTITY_GRPC_ADDR", "identity:50051"),
    notifier: env("NOTIFIER_GRPC_ADDR", "notifier:50051"),
    kitchen: env("KITCHEN_GRPC_ADDR", "kitchen:50051"),
    planner: env("PLANNER_GRPC_ADDR", "planner:50051"),
    offers: env("OFFERS_GRPC_ADDR", "offers:50051"),
    ai: env("AI_GRPC_ADDR", "ai:50051"),
  },
} as const;

export type ServiceName = keyof typeof config.grpc;
