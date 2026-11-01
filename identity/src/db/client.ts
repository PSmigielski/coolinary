// Postgres via Bun's built-in SQL client + Drizzle. Migrations (drizzle/) run at
// start-up under an advisory lock so concurrent starts never race.
import { SQL } from "bun";
import { drizzle } from "drizzle-orm/bun-sql";
import { migrate } from "drizzle-orm/bun-sql/migrator";
import { config } from "../config";
import { log } from "../lib/log";
import * as schema from "./schema";

export const sql = new SQL(config.databaseUrl, { max: 10 });
export const db = drizzle({ client: sql, schema });

export async function runMigrations() {
  const lockKey = `${config.serviceName}:migrations`;
  const conn = await sql.reserve();
  try {
    await conn`SELECT pg_advisory_lock(hashtext(${lockKey}))`;
    await migrate(db, { migrationsFolder: "drizzle" });
    log.info("migrations applied");
  } finally {
    await conn`SELECT pg_advisory_unlock(hashtext(${lockKey}))`;
    conn.release();
  }
}
