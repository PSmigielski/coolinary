// Verifies Bearer JWTs issued by identity (Better Auth JWT plugin) against its JWKS.
// Keys are cached by jose; no call to identity per request.
import { createRemoteJWKSet, jwtVerify, type JWTPayload } from "jose";
import { config } from "../config";

const jwks = createRemoteJWKSet(new URL("/api/auth/jwks", config.identityHttpUrl), { cooldownDuration: 30_000 });

export interface AuthUser {
  id: string;
  email?: string;
  householdId?: string;
  plan: "free" | "plus";
}

export async function verifyBearer(header: string | undefined | null): Promise<AuthUser | null> {
  if (!header?.startsWith("Bearer ")) return null;
  try {
    const { payload } = await jwtVerify(header.slice(7), jwks, { issuer: config.authIssuer, audience: config.authIssuer });
    return toUser(payload);
  } catch {
    return null;
  }
}

function toUser(p: JWTPayload): AuthUser | null {
  if (!p.sub) return null;
  return {
    id: p.sub,
    email: typeof p.email === "string" ? p.email : undefined,
    householdId: typeof p.householdId === "string" ? p.householdId : undefined,
    plan: p.plan === "plus" ? "plus" : "free",
  };
}
