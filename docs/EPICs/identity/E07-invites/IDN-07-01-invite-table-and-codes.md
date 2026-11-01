# IDN-07-01 Invite table and code generator

| | |
|---|---|
| Service | identity |
| Epic | [IDN-07 Invites](README.md) |
| Release | R2 |
| Size | S |
| Status | todo |
| Depends on | IDN-04-01 |
| Unblocks | IDN-07-02 … IDN-07-05 |

## Goal

Invite codes are short enough to read aloud, can't be confused with each other, and can't be guessed.

## Business rules

1. Alphabet: `ABCDEFGHJKMNPQRSTUVWXYZ23456789` (31 characters; no `0 O 1 I L`). Length 6 → about 887 million codes.
2. Codes come from a cryptographic random source. On a collision with a code still in the table, generate again,
   up to 5 times; then `INTERNAL`.
3. Input normalisation (preview and accept): upper-case, remove spaces and dashes. Anything outside the alphabet
   after that → `INVALID_ARGUMENT` `invalid_code`. The app displays codes as `ABC-DEF`.
4. Link: `<INVITE_LINK_BASE>/<CODE>` (e.g. `https://coolinary.app/i/KX7M4P`); the QR code encodes the link. The app
   scheme `coolinary://invite/<CODE>` is the fallback when universal links aren't set up.
   Decided (Q181): `INVITE_LINK_BASE = https://coolinary.app/i` once the domain is bought; that domain serves
   `apple-app-site-association` and `assetlinks.json`.
5. Statuses: `pending`, `accepted`, `cancelled`, `expired`. Only `pending` with `expires_at > now` is usable.
6. A code stays reserved while its row exists, so an old link never opens someone else's household; rows are
   deleted 30 days after they end (IDN-07-05), which frees the code.

## Contract

No RPC (used by IDN-07-02 … 04). New setting `INVITE_LINK_BASE` in `config.ts` and the three env files.

## Data

```ts
export const invites = pgTable("invites", {
  id: uuid("id").primaryKey().defaultRandom(),
  householdId: uuid("household_id").notNull().references(() => households.id, { onDelete: "cascade" }),
  code: text("code").notNull(),
  createdByMemberId: uuid("created_by_member_id").references(() => householdMembers.id, { onDelete: "set null" }),
  status: text("status", { enum: ["pending", "accepted", "cancelled", "expired"] }).notNull().default("pending"),
  createdAt: timestamp("created_at", { withTimezone: true }).notNull().defaultNow(),
  expiresAt: timestamp("expires_at", { withTimezone: true }).notNull(),
  endedAt: timestamp("ended_at", { withTimezone: true }),
  acceptedByUserId: text("accepted_by_user_id"),
}, (t) => [
  uniqueIndex("invites_code_unique").on(t.code),
  index("invites_household_pending_idx").on(t.householdId).where(sql`status = 'pending'`),
]);
```

## Errors

| Case | gRPC status |
|---|---|
| Code with characters outside the alphabet | `INVALID_ARGUMENT` `invalid_code` |

## Acceptance criteria

- [ ] 100,000 generated codes contain only alphabet characters and no duplicates.
- [ ] `" kx7-m4p "` normalises to `KX7M4P`; `KX7-M0P` is refused (`0` is not in the alphabet).
