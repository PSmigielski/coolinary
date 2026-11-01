// HTTP server for Better Auth (/api/auth/*). Not published: the gateway proxies to it.
import { Elysia } from "elysia";
import { auth } from "./auth";
import { config } from "./config";

export function startHttp() {
  return new Elysia().mount(auth.handler).listen(config.httpPort);
}
