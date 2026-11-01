// Docker HEALTHCHECK for the distroless image (no shell, no curl).
export {};

const port = process.env.ADMIN_PORT ?? "8081";
try {
  const res = await fetch(`http://127.0.0.1:${port}/healthz`, { signal: AbortSignal.timeout(3000) });
  process.exit(res.ok ? 0 : 1);
} catch {
  process.exit(1);
}
