import { startAdmin } from "./admin";
import { runMigrations, sql } from "./db/client";
import { startGrpcServer } from "./grpc/server";
import { startHttp } from "./http";
import { log } from "./lib/log";
import { identityService } from "./service";

await runMigrations();
const { server } = await startGrpcServer([identityService]);
const http = startHttp();
const admin = startAdmin();
log.info("identity listening", { http: http.server?.port });

for (const signal of ["SIGINT", "SIGTERM"] as const) {
  process.on(signal, () => {
    log.info("shutting down", { signal });
    server.tryShutdown(async () => {
      await http.stop();
      admin.stop();
      await sql.close();
      process.exit(0);
    });
  });
}
