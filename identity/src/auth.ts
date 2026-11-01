// Better Auth: e-mail + password accounts, sessions, and short-lived JWTs that the
// gateway verifies locally against /api/auth/jwks. The Expo plugin supports the app.
import { expo } from "@better-auth/expo";
import { betterAuth } from "better-auth";
import { drizzleAdapter } from "better-auth/adapters/drizzle";
import { jwt } from "better-auth/plugins";
import { config } from "./config";
import { db } from "./db/client";
import * as schema from "./db/schema";
import { createPersonalHousehold, householdClaims } from "./households";

export const auth = betterAuth({
  appName: "Coolinary",
  baseURL: config.authUrl,
  basePath: "/api/auth",
  secret: config.authSecret,
  database: drizzleAdapter(db, { provider: "pg", schema }),
  emailAndPassword: { enabled: true, minPasswordLength: 10 },
  trustedOrigins: config.trustedOrigins,
  databaseHooks: {
    user: {
      create: {
        after: async (user) => {
          await createPersonalHousehold(user);
        },
      },
    },
  },
  plugins: [
    expo(),
    jwt({
      jwt: {
        expirationTime: "15m",
        definePayload: async ({ user }) => ({ email: user.email, ...(await householdClaims(user.id)) }),
      },
    }),
  ],
});
