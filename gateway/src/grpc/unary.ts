// Promise wrapper around grpc-js unary calls: deadline, metadata forwarding,
// cancellation via AbortSignal, and client metrics.
import { Metadata, status, type CallOptions, type ClientUnaryCall, type ServiceError } from "@grpc/grpc-js";
import { rpcClientDuration, rpcClientRequests } from "../lib/metrics";

export type Invoke<Res> = (
  metadata: Metadata,
  options: Partial<CallOptions>,
  callback: (err: ServiceError | null, res?: Res) => void,
) => ClientUnaryCall;

export interface CallContext {
  /** Deadline in milliseconds (every call has one). Default 2000. */
  timeoutMs?: number;
  /** Forwarded as gRPC metadata: x-user-id, x-household-id, x-plan, x-request-id, traceparent. */
  metadata?: Record<string, string | undefined>;
  signal?: AbortSignal;
}

export function codeName(code: number): string {
  return (status as unknown as Record<number, string>)[code] ?? String(code);
}

export function unary<Res>(service: string, method: string, invoke: Invoke<Res>, ctx: CallContext = {}): Promise<Res> {
  const md = new Metadata();
  for (const [k, v] of Object.entries(ctx.metadata ?? {})) if (v) md.set(k, v);
  const deadline = new Date(Date.now() + (ctx.timeoutMs ?? 2000));
  const end = rpcClientDuration.startTimer({ rpc_service: service, rpc_method: method });

  return new Promise<Res>((resolve, reject) => {
    const call = invoke(md, { deadline }, (err, res) => {
      end();
      rpcClientRequests.inc({ rpc_service: service, rpc_method: method, rpc_code: err ? codeName(err.code).toLowerCase() : "ok" });
      if (err) reject(err);
      else resolve(res as Res);
    });
    ctx.signal?.addEventListener("abort", () => call.cancel(), { once: true });
  });
}
