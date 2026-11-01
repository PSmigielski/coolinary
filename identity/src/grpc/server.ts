// gRPC server bootstrap: unary handler wrapper (metrics, logging, errors),
// grpc.health.v1 and, in dev, server reflection.
import { readFileSync } from "node:fs";
import {
  Server,
  ServerCredentials,
  status,
  type Metadata,
  type ServiceDefinition,
  type UntypedServiceImplementation,
  type handleUnaryCall,
} from "@grpc/grpc-js";
import { config } from "../config";
import { log } from "../lib/log";
import { rpcServerDuration, rpcServerRequests } from "../lib/metrics";
import { HealthStatus, healthService } from "./health";

/** Throw from a handler to return a specific gRPC status. */
export class RpcError extends Error {
  constructor(
    readonly code: status,
    message: string,
  ) {
    super(message);
  }
}

const codeName = (code: number) => ((status as unknown as Record<number, string>)[code] ?? "unknown").toLowerCase();

/** Wraps an async handler into a grpc-js unary handler. `path` is the generated method path. */
export function unary<Req, Res>(
  path: string,
  fn: (req: Req, metadata: Metadata) => Promise<Res>,
): handleUnaryCall<Req, Res> {
  const [, service = "", method = ""] = path.split("/");
  return (call, callback) => {
    const end = rpcServerDuration.startTimer({ rpc_service: service, rpc_method: method });
    const t0 = performance.now();
    fn(call.request, call.metadata).then(
      (res) => {
        end();
        rpcServerRequests.inc({ rpc_service: service, rpc_method: method, rpc_code: "ok" });
        log.debug("rpc", { method: path, code: "ok", duration_ms: Math.round(performance.now() - t0) });
        callback(null, res);
      },
      (err: unknown) => {
        end();
        const code = err instanceof RpcError ? err.code : status.INTERNAL;
        rpcServerRequests.inc({ rpc_service: service, rpc_method: method, rpc_code: codeName(code) });
        log.warn("rpc", { method: path, code: codeName(code), error: String(err) });
        callback({ code, details: err instanceof Error ? err.message : String(err) });
      },
    );
  };
}

export interface Registration {
  name: string;
  definition: ServiceDefinition;
  implementation: UntypedServiceImplementation;
}

export async function startGrpcServer(services: Registration[]) {
  const server = new Server({ "grpc.keepalive_permit_without_calls": 1, "grpc.http2.min_ping_interval_without_data_ms": 10_000 });
  const health = new Map<string, number>([["", HealthStatus.SERVING], ...services.map((s) => [s.name, HealthStatus.SERVING] as [string, number])]);
  for (const s of services) server.addService(s.definition, s.implementation);
  const h = healthService(health);
  server.addService(h.definition, h.implementation);
  if (config.grpcReflection) {
    // Dev only. Loaded lazily: these libraries read .proto files from disk, which the
    // bundled prod image doesn't have. The descriptor set is written by `make proto`.
    const { loadFileDescriptorSetFromBuffer } = await import("@grpc/proto-loader");
    const { ReflectionService } = await import("@grpc/reflection");
    const pkg = loadFileDescriptorSetFromBuffer(readFileSync("src/gen/descriptor.binpb"));
    new ReflectionService(pkg, { services: [...services.map((s) => s.name), "grpc.health.v1.Health"] }).addToServer(server);
  }
  await new Promise<void>((resolve, reject) =>
    server.bindAsync(`0.0.0.0:${config.grpcPort}`, ServerCredentials.createInsecure(), (err) => (err ? reject(err) : resolve())),
  );
  log.info("gRPC server listening", { port: config.grpcPort, services: services.map((s) => s.name), reflection: config.grpcReflection });
  return { server, health };
}
