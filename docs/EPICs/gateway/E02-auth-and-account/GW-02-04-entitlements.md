# GW-02-04 `GET /v1/entitlements`

| | |
|---|---|
| Service | gateway |
| Epic | [GW-02 Auth and account routes](README.md) |
| Release | R1 |
| Size | S |
| Depends on | GW-02-01, GW-03-02, IDN-05 (`GetEntitlements`) |
| Unblocks | MOB-03 (paywall and limits shown in the app) |

## Goal

One call tells the app what the household's plan allows and how much Free AI quota is left this week, so the app
can show limits and the paywall before the user hits a 403 or 429.

## Business rules

1. Calls identity `GetEntitlements(household_id, actor_user_id)` (IDN-05), which returns the plan, where it came from and the
   feature map. The gateway adds the live quota counters from GW-03-02.
2. The plan in the response is identity's current value, which can be newer than the JWT's `plan` claim (up to 15
   minutes older). If they differ, the response sets `tokenRefreshRequired: true` so the app fetches a new JWT.
3. Feature map (values from identity; the table in docs/EPICs/README.md "Free and Plus" is the source):
   `aiDayPlansPerWeek` (2 or `null` = unlimited), `photoScansPerWeek` (2 or `null`), `otherAi`, `expiryReminders`,
   `calorieGoals` (`"self"` or `"everyone"`), `perPersonDiets`, `shopsWithOffers` (2 or `null`), `tripPlan`,
   `sharedHousehold`, `maxAccounts` (1 or 5), `savingsReport`, `nearShopReminders`.
4. Cached per household in `redis-cache` for 60 s (`gw:ent:<householdId>`); the quota part is never cached.
   Identity down → serve the cached value if any, else derive the map from the JWT plan with `stale: true`.

## Contract

`GET /v1/entitlements` — `auth: true`, tag `account`.

```json
200 {
  "plan": "free", "planSource": "none", "tokenRefreshRequired": false, "stale": false,
  "features": {
    "aiDayPlansPerWeek": 2, "photoScansPerWeek": 2, "otherAi": false, "expiryReminders": false,
    "calorieGoals": "self", "perPersonDiets": false, "shopsWithOffers": 2, "tripPlan": false,
    "sharedHousehold": false, "maxAccounts": 1, "savingsReport": false, "nearShopReminders": false
  },
  "quota": {
    "dayPlan": { "limit": 2, "used": 1, "remaining": 1, "nextFreeAt": "2026-11-05T18:42:10Z" },
    "scan":    { "limit": 2, "used": 0, "remaining": 2, "nextFreeAt": null }
  }
}
```

`planSource`: `none`, `manual` (set in admin), `store` (later, payments not decided).

## Data (redis-cache)

| Key | Value | TTL |
|---|---|---|
| `gw:ent:<householdId>` | identity's response JSON | 60 s |

## Errors

| Case | HTTP | `code` |
|---|---|---|
| No household | 403 | `no_household` |
| Identity down, nothing cached | 200 with `stale: true` | — |

## Acceptance criteria

- [ ] Free household with one AI day plan this week → `dayPlan.remaining: 1`.
- [ ] Plus household → `limit: null` and `remaining: null` for both kinds.
- [ ] After identity switches a household to Plus, the next call (cache expired) returns `plan: "plus"` and
      `tokenRefreshRequired: true` while the JWT still says `free`.
- [ ] With identity stopped and an empty cache, the route returns 200, `stale: true`, features from the JWT plan.
