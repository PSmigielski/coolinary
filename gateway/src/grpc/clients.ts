// gRPC clients for the services the gateway calls. Channels connect lazily, so
// the gateway starts even when a service is down.
import { Client, credentials, type ChannelOptions } from "@grpc/grpc-js";
import { BinaryReader, BinaryWriter } from "@bufbuild/protobuf/wire";
import { config, type ServiceName } from "../config";
import { IdentityServiceClient } from "../gen/coolinary/identity/v1/identity";
import { NotifierServiceClient } from "../gen/coolinary/notifier/v1/notifier";
import { KitchenServiceClient } from "../gen/coolinary/kitchen/v1/kitchen";
import { PlannerServiceClient } from "../gen/coolinary/planner/v1/planner";
import { OffersServiceClient } from "../gen/coolinary/offers/v1/offers";
import { codeName, unary } from "./unary";

// Keepalive pings detect dead connections after idle periods (see Bun HTTP/2 notes in docs/infrastructure.md).
export const channelOptions: ChannelOptions = {
  "grpc.keepalive_time_ms": 30_000,
  "grpc.keepalive_timeout_ms": 10_000,
  "grpc.keepalive_permit_without_calls": 1,
};

const insecure = credentials.createInsecure();

export const clients = {
  identity: new IdentityServiceClient(config.grpc.identity, insecure, channelOptions),
  notifier: new NotifierServiceClient(config.grpc.notifier, insecure, channelOptions),
  kitchen: new KitchenServiceClient(config.grpc.kitchen, insecure, channelOptions),
  planner: new PlannerServiceClient(config.grpc.planner, insecure, channelOptions),
  offers: new OffersServiceClient(config.grpc.offers, insecure, channelOptions),
};

// grpc.health.v1.Health/Check, encoded by hand so no extra proto is needed.
const healthClients = new Map<ServiceName, Client>();
const SERVING_STATUS = ["UNKNOWN", "SERVING", "NOT_SERVING", "SERVICE_UNKNOWN"];

export async function checkHealth(name: ServiceName, timeoutMs = 1000): Promise<string> {
  let client = healthClients.get(name);
  if (!client) {
    client = new Client(config.grpc[name], insecure, channelOptions);
    healthClients.set(name, client);
  }
  const c = client;
  try {
    const statusCode = await unary<number>("grpc.health.v1.Health", "Check", (md, opts, cb) =>
      c.makeUnaryRequest(
        "/grpc.health.v1.Health/Check",
        () => Buffer.from(new BinaryWriter().finish()), // empty HealthCheckRequest = whole server
        (buf: Buffer) => {
          const r = new BinaryReader(buf);
          let s = 0;
          while (r.pos < r.len) {
            const [field, wire] = r.tag();
            if (field === 1) s = r.int32();
            else r.skip(wire);
          }
          return s;
        },
        {},
        md,
        opts,
        cb,
      ), { timeoutMs });
    return SERVING_STATUS[statusCode] ?? "UNKNOWN";
  } catch (err) {
    const e = err as { code?: number };
    return `DOWN (${e.code !== undefined ? codeName(e.code) : String(err)})`;
  }
}
