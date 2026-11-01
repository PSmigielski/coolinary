import { startAdmin } from "./admin";
import { runMigrations, sql } from "./db/client";
import { startGrpcServer } from "./grpc/server";
import { log } from "./lib/log";
import { notifierService } from "./service";

await runMigrations();
const { server } = await startGrpcServer([notifierService]);
const admin = startAdmin();

for (const signal of ["SIGINT", "SIGTERM"] as const) {
  process.on(signal, () => {
    log.info("shutting down", { signal });
    server.tryShutdown(async () => {
      admin.stop();
      await sql.close();
      process.exit(0);
    });
  });
}
