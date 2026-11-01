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
  serviceName: env("SERVICE_NAME", "identity"),
  httpPort: Number(env("HTTP_PORT", "3000")),
  grpcPort: Number(env("GRPC_PORT", "50051")),
  adminPort: Number(env("ADMIN_PORT", "8081")),
  logLevel: env("LOG_LEVEL", "info"),
  grpcReflection: env("GRPC_REFLECTION", "false") === "true",
  databaseUrl: env("DATABASE_URL"),
  redisCacheUrl: env("REDIS_CACHE_URL", "redis://redis-cache:6379/0"),
  redisStateUrl: env("REDIS_STATE_URL", "redis://redis-state:6379/0"),
  // Public URL of the API (through the gateway), used by Better Auth for links and cookies.
  authUrl: env("BETTER_AUTH_URL"),
  authSecret: env("BETTER_AUTH_SECRET"),
  trustedOrigins: env("TRUSTED_ORIGINS", "coolinary://").split(",").map((s) => s.trim()).filter(Boolean),
  notifierAddr: env("NOTIFIER_GRPC_ADDR", "notifier:50051"),
} as const;
