import { pgTable, text, timestamp, uuid } from "drizzle-orm/pg-core";

// Delivery log: what was sent, when, and whether it worked. Recipient is stored
// as given for now; hash or drop it once retention rules are set.
export const deliveries = pgTable("deliveries", {
  id: uuid("id").primaryKey().defaultRandom(),
  channel: text("channel", { enum: ["email", "push"] }).notNull(),
  recipient: text("recipient").notNull(),
  subject: text("subject").notNull().default(""),
  status: text("status", { enum: ["sent", "failed"] }).notNull(),
  providerId: text("provider_id"),
  error: text("error"),
  createdAt: timestamp("created_at", { withTimezone: true }).notNull().defaultNow(),
});
