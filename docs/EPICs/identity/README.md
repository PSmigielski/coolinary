# identity — epics and tasks

Bun + Elysia + Better Auth (`identity/`). Owns accounts, sign-in, sessions and JWTs, households and their members
(accounts and profiles), invites, the household plan and entitlements, account deletion, and the admin role. It is
the producer of `events.identity`, which kitchen, planner, offers and notifier follow.

What exists today: Better Auth with e-mail + password and the JWT plugin (15 min), a personal household created on
sign-up, JWT claims `householdId` and `plan`, and `IdentityService.GetHousehold`.

## Epics

| Epic | Title | Rel. | Tasks |
|---|---|---|---|
| [IDN-01](E01-events-outbox/README.md) | Events: outbox and Redis Streams publisher (Bun) | R1 | 4 |
| [IDN-02](E02-email-code-sign-in/README.md) | Sign-in with e-mail code; remove passwords | R1 | 5 |
| [IDN-03](E03-anonymous-accounts/README.md) | Anonymous accounts and upgrade | R1 | 4 |
| [IDN-04](E04-households-core/README.md) | Households and members core | R1 | 5 |
| [IDN-05](E05-plan-and-entitlements/README.md) | Plan and entitlements | R1–R2 | 5 |
| [IDN-06](E06-profiles-without-accounts/README.md) | Profiles without accounts | R2 | 2 |
| [IDN-07](E07-invites/README.md) | Invites (QR, code, link) | R2 | 5 |
| [IDN-08](E08-leave-and-remove/README.md) | Leaving and removing members | R2 | 4 |
| [IDN-09](E09-account-deletion-and-export/README.md) | Account deletion and data export | R1 | 3 |
| [IDN-10](E10-admin/README.md) | Admin RPCs and admin role | R1 | 4 |

41 tasks: 29 in R1 (one gated), 12 in R2.

## Build order

### Release 1

| # | Task | Size | Depends on |
|---|---|---|---|
| 1 | [IDN-01-01](E01-events-outbox/IDN-01-01-outbox-table-and-emit.md) Outbox table and `emit()` | S | — |
| 2 | [IDN-01-03](E01-events-outbox/IDN-01-03-event-types.md) Identity event types and schemas | S | IDN-01-01 |
| 3 | [IDN-02-01](E02-email-code-sign-in/IDN-02-01-notifier-client.md) Notifier gRPC client | S | — |
| 4 | [IDN-01-02](E01-events-outbox/IDN-01-02-publisher.md) Publisher to `events.identity` | M | IDN-01-01, NTF-05 |
| 5 | [IDN-04-01](E04-households-core/IDN-04-01-one-household-per-account.md) One household per account | S | — |
| 6 | [IDN-04-02](E04-households-core/IDN-04-02-household-created-events.md) Household creation events | S | IDN-01-01, IDN-04-01 |
| 7 | [IDN-04-03](E04-households-core/IDN-04-03-jwt-claims.md) JWT claims | S | IDN-04-01 |
| 8 | [IDN-04-04](E04-households-core/IDN-04-04-membership-checks.md) Membership checks, richer `GetHousehold` | S | IDN-04-01 |
| 9 | [IDN-04-05](E04-households-core/IDN-04-05-rename.md) Rename household and display names | S | IDN-04-04 |
| 10 | [IDN-02-02](E02-email-code-sign-in/IDN-02-02-email-code-plugin.md) E-mail code plugin and code e-mail | M | IDN-02-01, NTF-02 |
| 11 | [IDN-02-03](E02-email-code-sign-in/IDN-02-03-sign-in-rate-limits.md) Sign-in rate limits | M | IDN-02-02, GW-02 |
| 12 | [IDN-02-04](E02-email-code-sign-in/IDN-02-04-remove-password-sign-up.md) Remove password sign-up | S | IDN-02-02 |
| 13 | [IDN-02-05](E02-email-code-sign-in/IDN-02-05-session-lifetime.md) Session lifetime and JWT issuance | S | IDN-02-02 |
| 14 | [IDN-03-01](E03-anonymous-accounts/IDN-03-01-anonymous-sign-in.md) Anonymous sign-in | M | IDN-04-02, IDN-02-03 |
| 15 | [IDN-03-02](E03-anonymous-accounts/IDN-03-02-upgrade-new-account.md) Upgrade to a new e-mail account | M | IDN-03-01, IDN-02-02 |
| 16 | [IDN-03-03](E03-anonymous-accounts/IDN-03-03-upgrade-existing-account.md) Upgrade into an existing account | M | IDN-03-02 |
| 17 | [IDN-05-01](E05-plan-and-entitlements/IDN-05-01-plan-changes.md) Plan fields, history, `household.plan_changed` | M | IDN-01-01 |
| 18 | [IDN-05-02](E05-plan-and-entitlements/IDN-05-02-get-entitlements.md) `GetEntitlements` | M | IDN-05-01, IDN-04-04 |
| 19 | [IDN-05-03](E05-plan-and-entitlements/IDN-05-03-manual-plan-expiry.md) Manual plan expiry | S | IDN-05-01 |
| 20 | [IDN-10-01](E10-admin/IDN-10-01-admin-role.md) Admin role, `IsAdmin`, `IdentityAdminService` | M | IDN-04-01 |
| 21 | [IDN-10-02](E10-admin/IDN-10-02-find-account.md) Find an account | S | IDN-10-01, IDN-05-01 |
| 22 | [IDN-10-03](E10-admin/IDN-10-03-set-plan.md) Set a plan by hand | S | IDN-10-01, IDN-05-01 |
| 23 | [IDN-10-04](E10-admin/IDN-10-04-account-stats.md) Account numbers | S | IDN-10-01 |
| 24 | [IDN-01-04](E01-events-outbox/IDN-01-04-snapshot-republish.md) Household snapshot republish | M | IDN-01-02, IDN-01-03, IDN-10-01 |
| 25 | [IDN-09-01](E09-account-deletion-and-export/IDN-09-01-deletion-code.md) Deletion confirmation code | S | IDN-02-01, NTF-02 |
| 26 | [IDN-09-02](E09-account-deletion-and-export/IDN-09-02-delete-account.md) Delete an account | L | IDN-09-01, IDN-01-01, IDN-04-01 |
| 27 | [IDN-09-03](E09-account-deletion-and-export/IDN-09-03-export-user-data.md) `ExportUserData` | S | IDN-04-04 |
| 28 | [IDN-03-04](E03-anonymous-accounts/IDN-03-04-stale-anonymous-cleanup.md) Delete stale anonymous accounts | M | IDN-03-01, IDN-09-02 |
| — | [IDN-05-05](E05-plan-and-entitlements/IDN-05-05-store-subscription-link.md) Store subscription link | L | `GATE:` payments decision |

### Release 2

| # | Task | Size | Depends on |
|---|---|---|---|
| 29 | [IDN-08-03](E08-leave-and-remove/IDN-08-03-transfer-ownership.md) Transfer ownership | S | IDN-04-01, IDN-04-04 |
| 30 | [IDN-08-01](E08-leave-and-remove/IDN-08-01-leave-household.md) Move-out function and leaving | M | IDN-04-02, IDN-04-04 |
| 31 | [IDN-08-02](E08-leave-and-remove/IDN-08-02-remove-member.md) Owner removes a member | S | IDN-08-01 |
| 32 | [IDN-08-04](E08-leave-and-remove/IDN-08-04-revoke-live-tokens.md) Cut off removed members' live tokens | S | IDN-08-01, GW-02 |
| 33 | [IDN-06-01](E06-profiles-without-accounts/IDN-06-01-add-profile.md) Add a profile | S | IDN-04-05, IDN-05-01 |
| 34 | [IDN-06-02](E06-profiles-without-accounts/IDN-06-02-delete-profile.md) Delete a profile | S | IDN-06-01 |
| 35 | [IDN-07-01](E07-invites/IDN-07-01-invite-table-and-codes.md) Invite table and code generator | S | IDN-04-01 |
| 36 | [IDN-07-02](E07-invites/IDN-07-02-create-list-cancel.md) Create, list and cancel invites | M | IDN-07-01, IDN-05-01 |
| 37 | [IDN-07-03](E07-invites/IDN-07-03-invite-preview.md) Invite preview | S | IDN-07-02 |
| 38 | [IDN-07-04](E07-invites/IDN-07-04-accept-invite.md) Accept an invite | L | IDN-07-03, IDN-04-02 |
| 39 | [IDN-07-05](E07-invites/IDN-07-05-invite-expiry.md) Invite expiry and cleanup | S | IDN-07-01 |
| 40 | [IDN-05-04](E05-plan-and-entitlements/IDN-05-04-shared-household-downgrade.md) Downgrade of a shared household | M | IDN-05-01, IDN-08-01, NTF-03 |

## Cross-service needs

- **gateway (GW-02, GW-04):** reads the new JWT claims (`memberId`, `anon`; `email` only for e-mail accounts);
  forwards the client IP in `x-forwarded-for` on the `/api/auth/*` proxy; checks the redis-state key
  `auth:revoked:<user_id>:<household_id>` (IDN-08-04); fetches a new token when an RPC answers `token_stale: true`;
  passes `actor_user_id` on every identity RPC; limits data exports to 3 per day.
- **notifier (NTF-02, NTF-03, NTF-05):** `SendEmail` templates `sign_in_code` and `account_deletion_code` in pl, de,
  en; `SendPush` category `household` (owner notified on join, IDN-07-04); `SendAlert` key `identity.outbox_stuck`.
- **kitchen (KIT-13):** on `household.member_added` with `previous_member_id`, copy that person's diet and target to
  the new member; with `brought_fridge_from`, move stock (also used for the anonymous merge, IDN-03-03).
- **all consumers:** treat `household.member_added` and `household.snapshot` as upserts.
- **kitchen, planner, offers (Q004):** call `GetEntitlements` whenever they need a household's plan (requests and
  jobs, `actor_user_id` empty for jobs); no local plan copies.

## Decisions

Settled with the founder on 1–2 November 2026 (full list: [DECISIONS.md](../DECISIONS.md)). Each task states its rule as `Decided (Qnnn)`.

| Q | Decision | Tasks |
|---|---|---|
| Q004 | Ask identity each time — kitchen, planner and offers read the plan from identity's GetEntitlements on each request that needs it, and per household in background jobs; no local plan copy. | IDN-05-02 |
| Q014 | contracts/events/<producer>/<type>.v1.json, next to the gRPC contracts. | IDN-01-03 |
| Q015 | Two fields: type and envelope. | IDN-01-02 |
| Q023 | Store empty; app shows "Mój dom" / "Mein Zuhause" / "My home". | IDN-04-02 |
| Q024 | The owner only. | IDN-04-05 |
| Q025 | Yes: from GetEntitlements, cached 5 minutes. | IDN-05-02 |
| Q026 | Yes, sign in first. | IDN-05-05 |
| Q027 | Own admin_users table plus a CLI command, not Better Auth's admin plugin. | IDN-10-01 |
| Q028 | Two roles: admin (everything) and dietitian. | IDN-10-01 |
| Q031 | 10 minutes, to allow for slow e-mail; the app countdown follows | IDN-02-02 |
| Q032 | Yes, 60 s, matching identity's limit of one code per minute | IDN-02-03 |
| Q033 | 60 days, extended on every daily use. | IDN-02-05 |
| Q034 | 5 per hour (identity's draft), with the gateway using the same number | IDN-03-01 |
| Q035 | Yes, only registered accounts can join | IDN-03-01 |
| Q036 | Move the fridge items into the account; drop plans, list and profile. | IDN-03-03 |
| Q041 | Yes: six shared names, written into service-conventions.md. | IDN-01-02 |
| Q044 | 30 days without use (identity's draft) | IDN-03-04 |
| Q137 | Owner chooses first — An owner of a shared household must name the new owner before deleting their account; with no other account member, the household is deleted with the account. | IDN-09-02 |
| Q138 | Yes: one JSON string per service, nested by the gateway. | IDN-09-03 |
| Q180 | At most 10. | IDN-06-01 |
| Q181 | https://coolinary.app/i/<code>, once the domain is bought | IDN-07-01 |
| Q182 | Yes: accounts plus waiting invites never exceed 5. | IDN-07-02 |
| Q183 | Any member with an account. | IDN-07-02 |
| Q184 | Yes: "Ewa joined your household". | IDN-07-04 |
| Q185 | Stays with the household until the payment route is chosen. | IDN-08-01 |
| Q186 | Yes, the owner picks the new owner before leaving | IDN-08-01 |
| Q187 | Members stay, no new invites — After a downgrade, existing members keep sharing; nobody new can be invited or join until Plus returns. This relaxes features.md's "sharing needs Plus". | IDN-05-04 |
| Q188 | Not needed: after a downgrade members stay (Q187), so no `sharing_ends` push. | IDN-05-04 |
| Q189 | No push; the app says "You are no longer in <household>" on next use. | IDN-08-02 |
| Q190 | No event; no other service stores who the owner is. | IDN-08-03 |
| Q191 | At once: a revocation key checked on every request (one Redis read) | IDN-08-04 |
