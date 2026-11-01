# IDN-07-05 Invite expiry and cleanup

| | |
|---|---|
| Service | identity |
| Epic | [IDN-07 Invites](README.md) |
| Release | R2 |
| Size | S |
| Depends on | IDN-07-01 |
| Unblocks | — |

## Goal

Expired invites disappear from the waiting list and old rows don't keep codes reserved forever.

## Business rules

1. Every 15 minutes (advisory lock): `pending` invites with `expires_at <= now` become `expired` with `ended_at`.
2. Preview, accept and list treat `expires_at <= now` as expired whether or not the job has run.
3. Daily at 04:30 UTC: delete invites whose `ended_at` is more than 30 days ago (this frees their codes).
4. An invite is also ended when its household is deleted (cascade) or its creator leaves (the invite stays valid;
   `created_by_member_id` becomes null and the preview shows the owner as inviter).

## Contract

No RPC.

## Data

Index `invites (expires_at) WHERE status = 'pending'`.

## Errors

None.

## Acceptance criteria

- [ ] An invite created 48 h 1 min ago shows `invite_expired` in preview and is marked expired by the next run.
- [ ] A row ended 31 days ago is deleted; one ended 29 days ago is kept.
