// Prometheus metrics. Names are shared by all services (docs/service-conventions.md).
import client from "prom-client";

export const registry = new client.Registry();
client.collectDefaultMetrics({ register: registry });

export const rpcServerRequests = new client.Counter({
  name: "rpc_server_requests_total",
  help: "RPCs handled, by service, method and status code.",
  labelNames: ["rpc_service", "rpc_method", "rpc_code"] as const,
  registers: [registry],
});

export const rpcServerDuration = new client.Histogram({
  name: "rpc_server_duration_seconds",
  help: "RPC handling time.",
  labelNames: ["rpc_service", "rpc_method"] as const,
  buckets: [0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10, 30],
  registers: [registry],
});
