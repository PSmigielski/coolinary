# KIT-11-06 GetMemberProfiles visibility and members in the planning snapshot

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-11 Per-person profiles](README.md) |
| Release | R2 |
| Size | M |
| Depends on | KIT-11-03 … KIT-11-05, KIT-10-08 |
| Unblocks | GW-05 (`GET /v1/profile/members`), PLN-12, PLN-13 |

## Goal

The app shows each person what they are allowed to see — never another adult's target — while planner gets every
person's target and limits to size the portions (features.md §2 "Privacy": "Other household members see portions
only, never another adult's weight or goal. Parents see their children's.").

## Business rules

1. `GetMemberProfiles(household_id, caller_member_id)` returns every member row of the household, with per field:

   | Field | Caller's own row | A child or a profile without account, caller adult | Another adult | Any row, caller a child |
   |---|---|---|---|---|
   | `daily_kcal_target` | ✓ (unless own row is a locked child) | ✓ | omitted | omitted |
   | `has_target` | ✓ | ✓ | ✓ | ✓ |
   | `target_kind`, `goal_set_by_parent` | ✓ | ✓ | ✓ | own only |
   | `calories_hidden` | ✓ | ✓ | omitted | own only |
   | diet section | ✓ | ✓ | ✓ | ✓ |

2. A locked child's own row never carries the target ("The child sees only their portion: no goal, no calories").
3. `is_me` marks the caller's row; rows are ordered owner first, then by creation.
4. **Planning snapshot** (`GetPlanningProfile.members`, KIT-10-08): every row, with `member_id`, `has_account`,
   `role`, `target_kind`, `daily_kcal_target` (or none), `calories_hidden`, diet section. On Free households the
   snapshot still lists every target; planner uses only the account member's (PLN-12-02 rule 6).
5. Neither response holds names, ages, weights, goals or consent details.
6. Export of a person's data (their own target and consent time) is KIT-13-07, not this RPC.

## Contract

`contracts/coolinary/kitchen/v1/profile.proto`

```proto
rpc GetMemberProfiles(GetMemberProfilesRequest) returns (GetMemberProfilesResponse);

message GetMemberProfilesRequest { string household_id = 1; string caller_member_id = 2; }
message GetMemberProfilesResponse { repeated MemberProfile members = 1; }   // MemberProfile: KIT-11-02

message PlanningMember {
  string member_id = 1;
  bool has_account = 2;
  string role = 3;                             // owner | member | profile
  TargetKind target_kind = 4;
  optional int32 daily_kcal_target = 5;
  bool calories_hidden = 6;
  DietSection diet = 7;
}
```

## Data

No schema change.

## Errors

| Case | gRPC status | Reason |
|---|---|---|
| Caller not in the household | `PERMISSION_DENIED` | `not_allowed` |

## Acceptance criteria

- [ ] Worked example (features.md): Adam calls → sees his own 2,000, Ola's 1,400 and Kuba's 3,200, and Ewa's row
      with `has_target = true` but no number.
- [ ] Kuba (locked child) calls → sees his own row without a target and no one else's targets.
- [ ] `GetPlanningProfile` for the same household lists all four targets (Verification 17/18 inputs).
- [ ] Verification 20: a second adult sees portions but not the other adult's goal — no field of Ewa's target reaches
      Adam's client (asserted on the gRPC response).
