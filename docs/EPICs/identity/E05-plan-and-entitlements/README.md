# IDN-05 Plan and entitlements

The plan (Free or Plus) belongs to the household and is held by its owner (features.md §6, business-model §3).
Identity is the only writer of the plan. Payments are decided later (decision 8), so for now a plan changes only by
admin (manual, optionally until a date); the store link is a gated placeholder.

## Flow

1. Admin sets Plus for a household (IDN-10-03) → `changePlan` writes the household, a history row and
   `household.plan_changed`.
2. Services keep their own copy of the plan from the event (derived decision D5); the next JWT carries the new plan.
3. The app reads `GetEntitlements` (through the gateway) to show what Free and Plus include and how many AI uses
   are left.

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| IDN-05-01 | Plan fields, plan history and `household.plan_changed` | R1 | M | IDN-01-01 |
| IDN-05-02 | `GetEntitlements` | R1 | M | IDN-05-01, IDN-04-04 |
| IDN-05-03 | Manual plan expiry | R1 | S | IDN-05-01 |
| IDN-05-04 | Downgrade of a shared household | R2 | M | IDN-05-01, IDN-08-01, NTF-03 (SendPush) |
| IDN-05-05 | Store subscription link | R1 | L | IDN-05-01 — `GATE:` payments decision |
