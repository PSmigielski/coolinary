# ADM-05 Households and plans support

**Goal.** Support tools on identity's data: find a household by e-mail or id (no fridge, plans or health data),
switch a household's plan by hand until payments are decided (decision 8), and delete an account on request (GDPR).

**Flow.** admin → `IdentityAdminService.FindAccount` (+ `IdentityService.GetHousehold`) / `SetPlan` / `AdminDeleteAccount` → identity publishes
`household.plan_changed` or `account.deleted` → the other services follow.

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [ADM-05-01](ADM-05-01-household-lookup.md) | Household lookup | R1 | M | ADM-02-03, IDN-10-02, IDN-04 |
| [ADM-05-02](ADM-05-02-manual-plan-switch.md) | Manual plan switch | R1 | S | ADM-05-01, IDN-05, IDN-10-03 |
| [ADM-05-03](ADM-05-03-account-deletion-request.md) | Delete an account on request | R1 | S | ADM-05-01, IDN-09-02 |
