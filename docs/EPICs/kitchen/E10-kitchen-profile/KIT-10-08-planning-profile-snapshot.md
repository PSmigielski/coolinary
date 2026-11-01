# KIT-10-08 GetPlanningProfile snapshot for planner

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-10 Kitchen profile](README.md) |
| Release | R1 (members from R2, KIT-11-06) |
| Size | M |
| Status | todo |
| Depends on | KIT-10-02 … KIT-10-06, KIT-01-06 (plan from identity) |
| Unblocks | PLN-01 (profile snapshot cache), every planner feature |

## Goal

Planner gets everything it needs to plan for a household in one call: settings, cuisines, tools, diet, staples,
people and the plan, plus a version it can cache by (decision 3: planner reads the profile from kitchen).

## Business rules

1. `GetPlanningProfile(household_id)` is called by planner only (service-to-service); the gateway never exposes it.
2. Contents:
   1. household settings (KIT-10-01, KIT-10-06), including timezone, language, country and currency;
   2. cuisines per meal (KIT-10-02);
   3. tools (KIT-10-03);
   4. household-level diets, allergens and dislikes (KIT-10-04);
   5. staple ingredient ids, always including water (KIT-10-05);
   6. the household's plan, `free` or `plus`, asked from identity on every call (`PlanOf`, KIT-01-06, Q004) and
      never stored in the cached snapshot;
   7. people (R2): filled by KIT-11-06; in R1 the list is empty and planner treats the household as one profile.
      On Free, per-person diets and dislikes are left out; household allergens already include everyone's
      allergens (Q195, KIT-11-04).
3. `profile_version` is the household's `planning_version` (KIT-10-09). With `if_profile_version` equal to it, the
   response has `not_modified = true` and no body except `plan`, which is always current.
4. A household that never saved its profile gets the defaults of its stored country, or the "other" defaults
   when none is stored yet (KIT-10-01 rule 2), with `is_default = true`.
5. The snapshot never holds age, height, weight, sex, activity or goal type — kitchen does not have them. Only daily
   kcal targets (R2).
6. Kitchen caches the built snapshot in `redis-cache` under `kitchen:planprof:<household_id>:<planning_version>`,
   TTL 1 hour; any change bumps the version, so stale entries are never read. On a target deletion the household's
   keys are deleted at once (KIT-11-03).
7. Deadline set by planner; the call reads Postgres (or the cache) and asks identity for the plan. identity down →
   KIT-01-06 rule 5 (`free` unless `x-plan` says otherwise).

## Contract

`contracts/coolinary/kitchen/v1/profile.proto`

```proto
rpc GetPlanningProfile(GetPlanningProfileRequest) returns (GetPlanningProfileResponse);

message GetPlanningProfileRequest {
  string household_id = 1;
  int64 if_profile_version = 2;
}

message PlanningProfile {
  string household_id = 1;
  int64 profile_version = 2;
  bool is_default = 3;
  string plan = 4;                              // "free" | "plus"
  HouseholdSettings household = 5;
  CuisineSection cuisines = 6;
  repeated coolinary.common.v1.KitchenTool tools = 7;
  DietSection household_diet = 8;
  repeated string staple_ingredient_ids = 9;    // includes water
  repeated PlanningMember members = 10;         // KIT-11-06, R2
}

message GetPlanningProfileResponse {
  PlanningProfile profile = 1;
  bool not_modified = 2;
  string plan = 3;                              // current plan, also when not_modified
}
```

## Data

No schema change.

## Errors

| Case | gRPC status | Behaviour |
|---|---|---|
| Malformed household id | `INVALID_ARGUMENT` | |
| Database down | `UNAVAILABLE` | planner never guesses hard limits (PLN-01-03 rule 5) |

## Acceptance criteria

- [ ] A household with tools [microwave, kettle], diet [vegan] and staples [salt] returns those, vegetarian included,
      and staples [salt, water].
- [ ] Calling with the current `profile_version` returns `not_modified = true`.
- [ ] After any `UpdateProfile`, the version is higher and the next call returns the new values.
- [ ] A test scans the response type for fields named age, weight, height, sex or activity and finds none.
- [ ] With identity answering `free` and then `plus` and no profile change, two calls return those two plans; the
      second may be `not_modified` and still carries `plan = plus`.
- [ ] p95 under 20 ms with a warm cache (identity call included).
