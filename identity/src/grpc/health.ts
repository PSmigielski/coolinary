// grpc.health.v1.Health/Check, with the two tiny messages encoded by hand, so no .proto
// file has to be read at runtime (the bundled prod image has none).
import { status, type ServiceDefinition, type handleUnaryCall } from "@grpc/grpc-js";
import { BinaryReader, BinaryWriter, WireType } from "@bufbuild/protobuf/wire";

export const HealthStatus = { UNKNOWN: 0, SERVING: 1, NOT_SERVING: 2 } as const;

interface CheckRequest {
  service: string;
}
interface CheckResponse {
  status: number;
}

function readField<T>(buf: Uint8Array, field: number, read: (r: BinaryReader) => T, fallback: T): T {
  const r = new BinaryReader(buf);
  let value = fallback;
  while (r.pos < r.len) {
    const [no, wire] = r.tag();
    if (no === field) value = read(r);
    else r.skip(wire);
  }
  return value;
}

const definition: ServiceDefinition = {
  check: {
    path: "/grpc.health.v1.Health/Check",
    requestStream: false,
    responseStream: false,
    requestSerialize: (m: CheckRequest) => Buffer.from(new BinaryWriter().tag(1, WireType.LengthDelimited).string(m.service).finish()),
    requestDeserialize: (b: Buffer): CheckRequest => ({ service: readField(b, 1, (r) => r.string(), "") }),
    responseSerialize: (m: CheckResponse) => Buffer.from(new BinaryWriter().tag(1, WireType.Varint).int32(m.status).finish()),
    responseDeserialize: (b: Buffer): CheckResponse => ({ status: readField(b, 1, (r) => r.int32(), 0) }),
  },
};

/** statuses: service name ("" = whole server) -> HealthStatus. Mutate it to change what Check reports. */
export function healthService(statuses: Map<string, number>) {
  const check: handleUnaryCall<CheckRequest, CheckResponse> = (call, callback) => {
    const s = statuses.get(call.request.service);
    if (s === undefined) callback({ code: status.NOT_FOUND, details: "unknown service" });
    else callback(null, { status: s });
  };
  return { definition, implementation: { check } };
}
