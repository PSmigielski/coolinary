import { startAdmin } from "./admin";
import { config } from "./config";
import { createApp } from "./http/app";
import { log } from "./lib/log";

const app = createApp().listen(config.httpPort);
const admin = startAdmin();
log.info("gateway listening", { http: config.httpPort, admin: config.adminPort });

for (const signal of ["SIGINT", "SIGTERM"] as const) {
  process.on(signal, async () => {
    log.info("shutting down", { signal });
    await app.stop();
    admin.stop();
    process.exit(0);
  });
}
