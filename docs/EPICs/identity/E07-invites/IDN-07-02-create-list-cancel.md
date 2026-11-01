# IDN-07-02 Create, list and cancel invites

| | |
|---|---|
| Service | identity |
| Epic | [IDN-07 Invites](README.md) |
| Release | R2 |
| Size | M |
| Depends on | IDN-07-01, IDN-05-01, IDN-04-05 |
| Unblocks | GW-04, IDN-07-03 |

## Goal

A Plus household can invite up to four more people, see who hasn't joined yet and take an invite back.

## Business rules

1. `CreateInvite(household_id, actor_user_id)`:
   1. The actor is an account member and not anonymous; anonymous → `FAILED_PRECONDITION` `sign_in_required`.
   2. The household is Plus (IDN-05-04); otherwise `FAILED_PRECONDITION`
      `plus_required`. On Free the app opens the Plus offer instead (features.md).
   3. Accounts plus pending invites must stay ≤ 5: with 3 accounts at most 2 pending invites. Otherwise
      `FAILED_PRECONDITION` `household_full`.
      Decided (Q182): pending invites count toward the 5, so the limit holds even when all accept.
   4. The inviter's display name must not be empty, so the preview can say who invited; otherwise
      `FAILED_PRECONDITION` `name_required`.
   5. Creates a `pending` invite valid for 48 hours (`expires_at = now + 48 h`).
   6. Returns `invite_id`, `code`, `display_code` (`ABC-DEF`), `link`, `expires_at`.
2. Who may invite: any account member.
   Decided (Q183): any account member.
3. `ListInvites(household_id, actor_user_id)` returns pending, unexpired invites: id, display code, link,
   creator's display name, created and expiry time, newest first. Any account member sees all of them.
4. `CancelInvite(household_id, actor_user_id, invite_id)`: the creator or the owner. Sets `cancelled` and
   `ended_at`. Cancelling an invite that is already cancelled, expired or accepted is a no-op that returns the
   invite as it is.
5. Invites are not announced by e-mail; the app shares the link through the phone's share sheet.
6. No events: other services don't track invites.

## Contract

```proto
rpc CreateInvite(CreateInviteRequest) returns (Invite);
rpc ListInvites(ListInvitesRequest) returns (ListInvitesResponse);
rpc CancelInvite(CancelInviteRequest) returns (Invite);

message CreateInviteRequest { string household_id = 1; string actor_user_id = 2; }
message ListInvitesRequest { string household_id = 1; string actor_user_id = 2; }
message CancelInviteRequest { string household_id = 1; string actor_user_id = 2; string invite_id = 3; }
message ListInvitesResponse { repeated Invite invites = 1; }
message Invite {
  string id = 1;
  string code = 2;              // "KX7M4P"
  string display_code = 3;      // "KX7-M4P"
  string link = 4;
  string status = 5;
  string created_by_display_name = 6;
  string created_at = 7;
  string expires_at = 8;
}
```

## Data

No schema change (IDN-07-01).

## Errors

Reason codes (in backticks) travel as trailing gRPC metadata `x-error-code` (README decision 18); version
conflicts use `ABORTED`.

| Case | gRPC status / reason |
|---|---|
| Anonymous actor | `FAILED_PRECONDITION` `sign_in_required` |
| Free | `FAILED_PRECONDITION` `plus_required` |
| Accounts + pending = 5 | `FAILED_PRECONDITION` `household_full` |
| Inviter has no display name | `FAILED_PRECONDITION` `name_required` |
| Cancel by someone else than creator or owner | `PERMISSION_DENIED` |

## Acceptance criteria

- [ ] Free household → `plus_required`.
- [ ] Plus household with 3 accounts: two invites succeed, the third returns `household_full`; cancelling one makes
      room again.
- [ ] Listed invites never include expired ones, even before the expiry job runs.
- [ ] Verification 15 (features.md): the same invite works by QR, by code and by link (with IDN-07-04).
