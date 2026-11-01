# IDN-05-01 Plan fields, plan history and `household.plan_changed`

| | |
|---|---|
| Service | identity |
| Epic | [IDN-05 Plan and entitlements](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | IDN-01-01 |
| Unblocks | IDN-05-02 … IDN-05-05, IDN-10-03, kitchen/planner/offers plan copies (D5) |

## Goal

A household's plan changes in exactly one place, every change is recorded, and every service learns about it.

## Business rules

1. `changePlan(tx, { householdId, to, source, expiresAt, reason, changedBy })` in `src/plans.ts` is the only code
   that writes `households.plan`.
2. `source`: `default` (Free with nothing bought), `manual` (set by an admin), `store` (reserved for IDN-05-05;
   refused until that task ships).
3. `to = plus` with `source = manual` may carry `expiresAt` (a future time, at most 3 years ahead) or none (until
   changed). `to = free` always clears `expiresAt` and sets `source = default`.
4. A call that changes nothing (same plan, source and expiry) writes nothing and emits nothing.
5. Every real change writes a `plan_changes` row (from, to, source, expiry, reason, changed_by, time) and emits
   `household.plan_changed` (`household_id`, `plan`, `previous_plan`, `source`) in the same transaction. A change of
   only the expiry date writes the history row but emits no event (the plan didn't change).
6. `changed_by` is the admin's user id for manual changes, or `system` for jobs (IDN-05-03, IDN-05-04).
7. `reason` is required for manual changes: 3–200 characters, free text for support notes; stored, never logged.
8. The plan applies to the whole household: every member, account or not, gets Plus features while the household
   is Plus.
9. The JWT `plan` claim follows on the next token (IDN-04-03).

## Contract

No public RPC (admin RPC in IDN-10-03).

## Data

```ts
// households: new columns
planSource: text("plan_source", { enum: ["default", "manual", "store"] }).notNull().default("default"),
planExpiresAt: timestamp("plan_expires_at", { withTimezone: true }),
planChangedAt: timestamp("plan_changed_at", { withTimezone: true }).notNull().defaultNow(),

export const planChanges = pgTable("plan_changes", {
  id: uuid("id").primaryKey().defaultRandom(),
  householdId: uuid("household_id").notNull().references(() => households.id, { onDelete: "cascade" }),
  fromPlan: text("from_plan", { enum: ["free", "plus"] }).notNull(),
  toPlan: text("to_plan", { enum: ["free", "plus"] }).notNull(),
  source: text("source", { enum: ["default", "manual", "store"] }).notNull(),
  expiresAt: timestamp("expires_at", { withTimezone: true }),
  reason: text("reason").notNull().default(""),
  changedBy: text("changed_by").notNull(),          // admin user id or "system"
  createdAt: timestamp("created_at", { withTimezone: true }).notNull().defaultNow(),
}, (t) => [index("plan_changes_household_idx").on(t.householdId, t.createdAt)]);
```

## Events

Publishes `household.plan_changed`.

## Errors

| Case | Result |
|---|---|
| `source = store` before IDN-05-05 | `FAILED_PRECONDITION` `store_not_enabled` |
| `expiresAt` in the past or > 3 years ahead | `INVALID_ARGUMENT` |
| Manual change without a reason | `INVALID_ARGUMENT` |

## Acceptance criteria

- [ ] Free → Plus writes one history row and one event with `previous_plan: "free"`.
- [ ] Repeating the same change writes nothing.
- [ ] Changing only the expiry of a Plus plan writes a history row and no event.
- [ ] Deleting a household deletes its history rows.
