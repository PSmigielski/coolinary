// Prometheus metrics. Names are shared by all services (docs/service-conventions.md).
import client from "prom-client";

export const registry = new client.Registry();
client.collectDefaultMetrics({ register: registry });

export const httpRequests = new client.Counter({
  name: "http_server_requests_total",
  help: "HTTP requests handled, by method, route and status.",
  labelNames: ["method", "route", "status"] as const,
  registers: [registry],
});

export const httpDuration = new client.Histogram({
  name: "http_server_duration_seconds",
  help: "HTTP request handling time.",
  labelNames: ["method", "route"] as const,
  buckets: [0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10],
  registers: [registry],
});

export const rpcClientRequests = new client.Counter({
  name: "rpc_client_requests_total",
  help: "Outgoing RPCs, by target service, method and status code.",
  labelNames: ["rpc_service", "rpc_method", "rpc_code"] as const,
  registers: [registry],
});

export const rpcClientDuration = new client.Histogram({
  name: "rpc_client_duration_seconds",
  help: "Outgoing RPC time.",
  labelNames: ["rpc_service", "rpc_method"] as const,
  buckets: [0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10],
  registers: [registry],
});
