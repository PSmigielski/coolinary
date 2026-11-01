// Public HTTP API for the mobile app. Each route maps to gRPC calls on the
// services that own the data. OpenAPI is served at /openapi/json.

import { openapi } from "@elysiajs/openapi";
import { Elysia, t } from "elysia";
import { config, type ServiceName } from "../config";
import { unitToJSON } from "../gen/coolinary/common/v1/common";
import { type GetHouseholdResponse, planToJSON, roleToJSON } from "../gen/coolinary/identity/v1/identity";
import { type ResolveBarcodeResponse, resolveSourceToJSON } from "../gen/coolinary/kitchen/v1/kitchen";
import { checkHealth, clients } from "../grpc/clients";
import { type CallContext, codeName, unary } from "../grpc/unary";
import { verifyBearer } from "../lib/auth";
import { log } from "../lib/log";
import { httpDuration, httpRequests } from "../lib/metrics";
import { ping, redisCache, redisState } from "../lib/redis";

const started = new WeakMap<Request, number>();
const backendServices: ServiceName[] = ["identity", "notifier", "kitchen", "planner", "offers", "ai"];

function callContext(request: Request, user?: { id: string; householdId?: string; plan: string }): CallContext {
  return {
    signal: request.signal,
    metadata: {
      "x-request-id": request.headers.get("x-request-id") ?? crypto.randomUUID(),
      traceparent: request.headers.get("traceparent") ?? undefined,
      "x-user-id": user?.id,
      "x-household-id": user?.householdId,
      "x-plan": user?.plan,
    },
  };
}

const User = t.Object({
  id: t.String(),
  email: t.Optional(t.String()),
  householdId: t.Optional(t.String()),
  plan: t.Union([t.Literal("free"), t.Literal("plus")]),
});

const Product = t.Object({
  ean: t.String(),
  name: t.String(),
  brand: t.String(),
  packSize: t.Optional(t.Object({ value: t.Number(), unit: t.String() })),
});

export function createApp() {
  return (
    new Elysia()
      .use(
        openapi({
          path: "/openapi",
          documentation: {
            info: { title: "Coolinary API", version: "0.1.0", description: "Public API used by the Coolinary app." },
          },
          exclude: { paths: ["/api/auth/*"] },
        }),
      )
      .onRequest(({ request }) => {
        started.set(request, performance.now());
      })
      .onAfterResponse(({ request, route, set }) => {
        const t0 = started.get(request);
        const labels = { method: request.method, route: route || "unknown" };
        httpRequests.inc({ ...labels, status: String(set.status ?? 200) });
        if (t0 !== undefined) httpDuration.observe(labels, (performance.now() - t0) / 1000);
      })
      .onError(({ code, error, request }) => {
        if (code !== "NOT_FOUND" && code !== "VALIDATION")
          log.error("request failed", { path: new URL(request.url).pathname, code, error: String(error) });
      })
      .macro({
        auth: {
          async resolve({ headers, status }) {
            const user = await verifyBearer(headers.authorization);
            if (!user) return status(401, { error: "unauthorized" });
            return { user };
          },
        },
      })
      .get("/healthz", () => ({ ok: true }), { detail: { hide: true } })
      .get(
        "/v1/system/status",
        async () => {
          const entries = await Promise.all(backendServices.map(async (s) => [s, await checkHealth(s)] as const));
          const services = Object.fromEntries(entries);
          const stores = { redisCache: await ping(redisCache), redisState: await ping(redisState) };
          const ok = entries.every(([, v]) => v === "SERVING") && Object.values(stores).every((v) => v === "ok");
          return { ok, gateway: "SERVING", services, stores, checkedAt: new Date().toISOString() };
        },
        {
          detail: { summary: "Health of every back-end service, checked over gRPC", tags: ["system"] },
          response: t.Object({
            ok: t.Boolean(),
            gateway: t.String(),
            services: t.Record(t.String(), t.String()),
            stores: t.Object({ redisCache: t.String(), redisState: t.String() }),
            checkedAt: t.String(),
          }),
        },
      )
      .get(
        "/v1/products/:ean",
        async ({ params, request, status }) => {
          let res: ResolveBarcodeResponse;
          try {
            res = await unary(
              "kitchen",
              "ResolveBarcode",
              (md, o, cb) => clients.kitchen.resolveBarcode({ ean: params.ean, householdId: "" }, md, o, cb),
              callContext(request),
            );
          } catch (err) {
            const e = err as { code: number; details?: string };
            return status(e.code === 3 ? 400 : 502, { error: e.details ?? codeName(e.code) });
          }
          if (!res.found || !res.product) return status(404, { error: "unknown product" });
          const p = res.product;
          return {
            ean: p.ean,
            name: p.name,
            brand: p.brand,
            source: resolveSourceToJSON(res.source),
            packSize: p.packSize ? { value: p.packSize.value, unit: unitToJSON(p.packSize.unit) } : undefined,
          };
        },
        {
          params: t.Object({ ean: t.String({ pattern: "^[0-9]{8,14}$" }) }),
          detail: { summary: "Look up a scanned barcode", tags: ["kitchen"] },
          response: {
            200: t.Composite([Product, t.Object({ source: t.String() })]),
            400: t.Object({ error: t.String() }),
            404: t.Object({ error: t.String() }),
            502: t.Object({ error: t.String() }),
          },
        },
      )
      .get("/v1/me", ({ user }) => user, {
        auth: true,
        detail: { summary: "The signed-in user, from the JWT", tags: ["identity"] },
        response: { 200: User, 401: t.Object({ error: t.String() }) },
      })
      .get(
        "/v1/household",
        async ({ user, request, status }) => {
          const householdId = user.householdId;
          if (!householdId) return status(404, { error: "no household" });
          try {
            const res = await unary<GetHouseholdResponse>(
              "identity",
              "GetHousehold",
              (md, o, cb) => clients.identity.getHousehold({ householdId }, md, o, cb),
              callContext(request, user),
            );
            const h = res.household;
            if (!h) return status(502, { error: "identity returned no household" });
            return {
              id: h.id,
              name: h.name,
              plan: planToJSON(h.plan),
              members: h.members.map((m) => ({
                id: m.id,
                displayName: m.displayName,
                role: roleToJSON(m.role),
                hasAccount: m.hasAccount,
              })),
            };
          } catch (err) {
            const e = err as { code: number; details?: string };
            return status(e.code === 5 ? 404 : 502, { error: e.details ?? codeName(e.code) });
          }
        },
        {
          auth: true,
          detail: { summary: "The signed-in user's household (identity over gRPC)", tags: ["identity"] },
          response: {
            200: t.Object({
              id: t.String(),
              name: t.String(),
              plan: t.String(),
              members: t.Array(
                t.Object({ id: t.String(), displayName: t.String(), role: t.String(), hasAccount: t.Boolean() }),
              ),
            }),
            401: t.Object({ error: t.String() }),
            404: t.Object({ error: t.String() }),
            502: t.Object({ error: t.String() }),
          },
        },
      )
      // Better Auth runs in identity and works over HTTP; pass these through unchanged.
      .all(
        "/api/auth/*",
        ({ request }) => {
          const url = new URL(request.url);
          return fetch(config.identityHttpUrl + url.pathname + url.search, {
            method: request.method,
            headers: request.headers,
            body: request.method === "GET" || request.method === "HEAD" ? undefined : request.body,
            redirect: "manual",
          });
        },
        { parse: "none", detail: { hide: true } },
      )
  );
}

export type App = ReturnType<typeof createApp>;
