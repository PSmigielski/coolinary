# KIT-11-02 UpdateMemberProfile and who may change what

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-11 Per-person profiles](README.md) |
| Release | R2 |
| Size | M |
| Status | todo |
| Depends on | KIT-11-01 |
| Unblocks | KIT-11-03, KIT-11-04, KIT-11-05, GW-05 (`PUT /v1/profile/members/{memberId}`) |

## Goal

One RPC changes a person's profile, and it enforces who may change what: adults manage their own data and their
children's, nobody sees or edits another adult's goal, and a teenager cannot change a goal a parent set
(features.md §2 "Children and teenagers": "Only a parent sets the goal").

## Business rules

1. `UpdateMemberProfile(household_id, caller_member_id, member_id, version, fields…)`. The caller row must belong to
   the household and have an account whose `user_id` equals `x-user-id`; else `PERMISSION_DENIED`.
2. Each field is optional; a field not sent is unchanged. Fields: `diet` section (KIT-11-04), `daily_kcal_target`
   or `clear_target` and `health_consent` (KIT-11-03), `target_kind`, `goal_set_by_parent`, `calories_hidden`
   (KIT-11-05).
3. `version` must equal the row's version, else `ABORTED` with the current version (gateway → 409). Success adds 1.
4. **An adult** is a caller whose own row has `target_kind = adult`. The server cannot verify age; the kind is what
   the household set.
5. **Rights matrix** (✓ allowed, ✗ `PERMISSION_DENIED` reason `not_allowed`, 🔒 `PERMISSION_DENIED` reason
   `goal_locked`):

   | Field | Own row, adult | Own row, child with `goal_set_by_parent` | Own row, child without | Another adult account | A child or a profile without account |
   |---|---|---|---|---|---|
   | target, consent | ✓ | 🔒 | ✓ | ✗ | ✓ (adult caller) |
   | `target_kind`, `goal_set_by_parent` | ✓ (`adult` only) | 🔒 | ✓ (`adult` only) | ✗, except the owner (rule 6) | ✓ (adult caller) |
   | `calories_hidden` | ✓ | 🔒 | ✓ | ✗ | ✓ (adult caller) |
   | diets, allergens | ✓ | 🔒 | ✓ | ✗ | ✓ (adult caller) |
   | dislikes | ✓ | ✓ | ✓ | ✗ | ✓ (adult caller) |

6. Marking an **account member** as a child (`target_kind = child`, usually with `goal_set_by_parent = true`) is
   allowed only to the household owner, and only the owner can undo it; for profiles without an account any adult
   may set it. Decided (Q192): owner-only for accounts, so one adult cannot lock another.
7. A child caller (`target_kind = child`) has no rights over anyone else's row.
8. Free households (plan from identity, KIT-01-06): only the caller's own row may be changed, and only target, consent and
   `calories_hidden`; anything else → `FAILED_PRECONDITION` reason `plus_required` (features.md §6: one personal
   goal on Free; per-person profiles are Plus).
9. Every successful update bumps `planning_version` and writes `profile.changed` with `members` (KIT-10-09).
10. Values are never logged; the log line holds the RPC name, result code and the household hash.

## Contract

`contracts/coolinary/kitchen/v1/profile.proto`

```proto
rpc UpdateMemberProfile(UpdateMemberProfileRequest) returns (MemberProfile);

enum TargetKind { TARGET_KIND_UNSPECIFIED = 0; TARGET_KIND_ADULT = 1; TARGET_KIND_CHILD = 2; }

message HealthConsent {
  string version = 1;        // consent text version, e.g. "2026-11"
  string given_at = 2;       // RFC 3339
}

message UpdateMemberProfileRequest {
  string household_id = 1;
  string caller_member_id = 2;
  string member_id = 3;
  int32 version = 4;
  optional DietSection diet = 5;              // KIT-11-04
  optional int32 daily_kcal_target = 6;       // KIT-11-03
  bool clear_target = 7;
  HealthConsent health_consent = 8;
  optional TargetKind target_kind = 9;
  optional bool goal_set_by_parent = 10;
  optional bool calories_hidden = 11;         // KIT-11-05
}

message MemberProfile {
  string member_id = 1;
  bool is_me = 2;
  int32 version = 3;
  bool has_account = 4;
  TargetKind target_kind = 5;
  bool goal_set_by_parent = 6;
  bool has_target = 7;
  optional int32 daily_kcal_target = 8;       // only where KIT-11-06 allows
  bool calories_hidden = 9;
  DietSection diet = 10;
}
```

## Data

No schema change (`member_profiles` from KIT-11-01).

## Errors

Reason codes (in backticks) travel as trailing gRPC metadata `x-error-code` (README decision 18); version
conflicts use `ABORTED`.

| Case | gRPC status | Reason |
|---|---|---|
| Caller not in the household / not the JWT user | `PERMISSION_DENIED` | `not_allowed` |
| Field not allowed for this caller | `PERMISSION_DENIED` | `not_allowed` |
| Child changing a parent-set field | `PERMISSION_DENIED` | `goal_locked` (gateway → 403) |
| Free household, not own row or not a Free field | `FAILED_PRECONDITION` | `plus_required` (gateway → 403) |
| Version mismatch | `ABORTED` | `version_conflict` |
| Unknown member | `NOT_FOUND` | |

## Acceptance criteria

- [ ] Ewa (adult) sets Ola's (profile, child) target and hides her calories: allowed.
- [ ] Ewa tries to read or set Adam's (adult account) target: `not_allowed`.
- [ ] Kuba (15, account, child, `goal_set_by_parent`) changes his own target: `goal_locked`
      (features.md Verification 19, "a teenager's own account cannot change their goal"); adding a dislike works.
- [ ] Adam (adult, not owner) marks Ewa's account as a child: `not_allowed`; the owner can.
- [ ] On a Free household, the owner sets their own target (allowed) and a per-person diet (`plus_required`).
