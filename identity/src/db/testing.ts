// Postgres for tests: the `identity_test` database that `make test` points TEST_DATABASE_URL at while `make dev` is
// up. Tests that need it skip when the variable is unset (`test.skipIf(!hasTestDb)`).
import { SQL } from "bun";
import { drizzle } from "drizzle-orm/bun-sql";
import { migrate } from "drizzle-orm/bun-sql/migrator";
import * as schema from "./schema";

export const hasTestDb = Boolean(process.env.TEST_DATABASE_URL);

/** Opens the test database and applies the migrations (under the same advisory lock as start-up). */
export async function openTestDb() {
  const url = process.env.TEST_DATABASE_URL;
  if (!url) throw new Error("TEST_DATABASE_URL is not set");
  const sql = new SQL(url, { max: 4 });
  const db = drizzle({ client: sql, schema });
  const conn = await sql.reserve();
  try {
    await conn`SELECT pg_advisory_lock(hashtext('identity:migrations'))`;
    await migrate(db, { migrationsFolder: `${import.meta.dir}/../../drizzle` });
  } finally {
    await conn`SELECT pg_advisory_unlock(hashtext('identity:migrations'))`;
    conn.release();
  }
  return { db, sql };
}
