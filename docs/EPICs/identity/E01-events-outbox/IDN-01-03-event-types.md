# IDN-01-03 Identity event types and schemas

| | |
|---|---|
| Service | identity |
| Epic | [IDN-01 Events: outbox and publisher](README.md) |
| Release | R1 |
| Size | S |
| Status | done |
| Depends on | IDN-01-01 |
| Unblocks | All identity tasks that emit; consumers KIT-13, PLN-21, OFF-14, NTF-01 |

## Goal

Every identity event has one written definition that the producer checks against and that consumers in Go, Bun
and Python can read, so a payload field never means two things.

## Business rules

1. `src/events/types.ts` defines each event as a TypeScript type plus a runtime validator. JSON field names are
   snake_case. Ids are UUID strings except user ids, which are Better Auth text ids.
2. The events and their v1 payloads:

   | Type | Payload |
   |---|---|
   | `household.created` | `household_id`, `owner_user_id`, `owner_member_id` |
   | `household.member_added` | `household_id`, `member_id`, `user_id` (null for profiles), `has_account`, `role` (`owner`/`member`/`profile`), `brought_fridge_from` (household id or null), `previous_member_id` (member id or null) |
   | `household.member_removed` | `household_id`, `member_id`, `user_id` or null, `reason` (`left`/`removed`/`profile_deleted`/`plan_ended`), `new_household_id` or null |
   | `household.plan_changed` | `household_id`, `plan` (`free`/`plus`), `previous_plan`, `source` (`default`/`manual`/`store`) |
   | `household.deleted` | `household_id`, `reason` (`owner_deleted`/`joined_other`/`linked_merge`/`anonymous_cleanup`) |
   | `account.linked` | `anonymous_user_id`, `user_id`, `household_id`, `anonymous_household_id`, `merged` |
   | `account.deleted` | `user_id`, `member_ids` |
   | `household.snapshot` | `household_id`, `plan`, `members`: list of `member_id`, `user_id` or null, `has_account`, `role` |

3. Consumer rules, written in the same file's header comment and in `identity/EVENTS.md`:
   - Drop duplicates by envelope `id`.
   - `household.member_added` and `household.snapshot` are upserts by `member_id` / `household_id`: applying the
     same one twice, or a member that already exists, is not an error.
   - Within `events.identity`, events about one household arrive in the order identity committed them. Example:
     `member_added` with `brought_fridge_from = X` always comes before `household.deleted` for X.
   - Unknown event types are acknowledged and ignored. Unknown payload fields are ignored.
4. Versioning: adding an optional field keeps the version. Removing, renaming or changing the meaning of a field
   creates version 2; identity then emits v2 only after every consumer accepts both versions (one release).
5. A JSON Schema per event is written next to the types and exported to the shared contracts folder so non-Bun
   consumers can validate in tests.

## Contract

Decided (Q014): event schemas live in `contracts/events/identity/<type>.v1.json` (JSON Schema 2020-12), next to
the gRPC contracts. `make proto` copies nothing; tests read the schemas directly.

## Data

No schema change.

## Errors

None at runtime beyond IDN-01-01's validation.

## Acceptance criteria

- [ ] Every `emit` call site type-checks against these definitions.
- [ ] `identity/EVENTS.md` lists all eight events with an example payload each.
- [ ] A test validates one sample of every event against its JSON Schema.
- [ ] No payload type contains a field outside the list above (test enumerates keys).

## Notes

Differences from docs/EPICs/README.md, to be folded back: `household.created` drops `country` (identity doesn't
know it; the country lives in the kitchen profile, decision 3) and adds `owner_member_id`;
`household.member_added` adds `role` and `previous_member_id` (lets kitchen carry a person's diet and target into a
new household, IDN-07-04, IDN-08-01); `household.member_removed` keeps reason `plan_ended` reserved but identity never emits it (Q187: members stay after a downgrade, IDN-05-04);
`household.deleted` adds `reason`; `account.linked` adds `anonymous_household_id` and `merged` (IDN-03-03); new
event `household.snapshot` (IDN-01-04).
