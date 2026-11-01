import createClient from "openapi-fetch";
import createQueryClient from "openapi-react-query";

import type { operations, paths } from "./schema";

/**
 * Gateway base URL, inlined at bundle time from `.env` (EXPO_PUBLIC_API_URL).
 * On a phone this must be the dev machine's LAN IP, not localhost.
 */
export const API_URL = process.env.EXPO_PUBLIC_API_URL ?? "";

/** Plain typed fetch client — use for one-off calls outside React. */
export const client = createClient<paths>({ baseUrl: API_URL });

/** Typed TanStack Query hooks: `$api.useQuery("get", "/v1/system/status")`. */
export const $api = createQueryClient(client);

export type SystemStatus =
  operations["getV1SystemStatus"]["responses"][200]["content"]["application/json"];
