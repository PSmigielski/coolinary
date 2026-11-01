# ADM-05-02 Manual plan switch

| | |
|---|---|
| Service | admin |
| Epic | [ADM-05 Households and plans support](README.md) |
| Release | R1 |
| Size | S |
| Status | todo |
| Depends on | ADM-05-01; IDN-05 (plan source and expiry), IDN-10-03 (`IdentityAdminService.SetPlan`) |
| Unblocks | — |

## Goal

Until payments are decided (decision 8), the founder can give a household Plus by hand (founding pre-sale buyers,
testers, support cases) and take it back.

## Business rules

1. Input: household id, plan (`free` or `plus`), reason (`founding`, `tester`, `support`, `refund`, `other`), note
   (≤ 200 characters, required for `other`), optional expiry date (`plus` only; must be in the future, at most 3
   years ahead).
2. Admin calls `SetPlan(household_id, plan, expires_at, reason)` (IDN-10-03). identity's `reason` is one text of
   3–200 characters, so admin sends `"<reason code>: <note>"` (e.g. `"tester: beta group 2"`). identity sets plan
   source `manual` for Plus and `default` for Free, stores the admin user id, and publishes `household.plan_changed`;
   at expiry identity returns the household to Free (IDN-05). Setting Free on a shared household starts identity's
   sharing grace period (IDN-10-03 rule 6); admin shows that.
3. Setting the plan the household already has with the same expiry → 200, no change, no event.
4. A household whose plan source is `store` (once payments exist) cannot be switched here; identity returns
   `FAILED_PRECONDITION`.
5. The app sees the change after its next JWT refresh, at most 15 minutes (JWT lifetime); identity's answer says so
   with `token_stale: true`, and admin states it in the response.
6. Audited (`identity.plan.set`, household id, old and new plan, reason, expiry).
7. Role `admin` only.

## Contract

| Method | Path | Body / response |
|---|---|---|
| POST | `/admin/v1/households/{id}/plan` | `{ plan, reason, note?, expires_on? }` → `{ plan, plan_source, expires_on, effective_within_minutes: 15 }` |

## Data

None in admin.

## Errors

| Case | HTTP |
|---|---|
| Bad plan, reason, date, missing note | 400 |
| Plan comes from a store | 409 `store_plan` |
| Household not found | 404 |

## Acceptance criteria

- [ ] Switching a Free household to Plus emits one `household.plan_changed`; offers allows more than 2 shops on its
      next request for it, because it reads the plan from identity (Q004).
- [ ] A Plus with expiry tomorrow becomes Free after identity's expiry job runs.
