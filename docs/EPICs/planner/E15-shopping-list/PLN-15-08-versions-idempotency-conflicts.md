# PLN-15-08 Versions, idempotent offline edits, conflict rules

| | |
|---|---|
| Service | planner |
| Epic | [PLN-15 Shopping list](README.md) |
| Release | R1 |
| Size | M |
| Depends on | PLN-15-01 |
| Unblocks | PLN-15-05, PLN-15-06, PLN-15-07, PLN-15-09, MOB-04, GW-08 |

## Goal

Several people edit one list, sometimes offline in a shop with no signal. Every change must apply once, in a
predictable way, and the result must be the same on every phone.

## Business rules

1. Every list mutation RPC carries a `client_mutation_id` (UUID made on the phone). The first result is stored for
   7 days; a replay with the same id returns the stored result and changes nothing.
2. Mutations are applied in the order the server receives them, each in one transaction that takes a row lock on the
   household's list (`SELECT … FOR UPDATE` on `shopping_lists`).
3. Each committed mutation sets `list.version = version + 1`, and every line it touches gets that number as its
   `version`.
4. **Last writer wins per field:** amount, and ticked/unticked state. An edit with an older `base_version` is still
   applied; the response returns the line as it now is so the phone can show it.
5. **Remove beats edit:** an amount edit on a removed line returns `NOT_FOUND`; the phone drops the edit.
6. **Bought beats removed:** a tick on a line removed in the last 7 days restores it as ticked (PLN-15-07 rule 8).
7. **Same ingredient added twice** (two phones offline) → merged into one open line with both contributions.
8. **Adding an ingredient whose line is ticked** creates a new open line (the ticked one was already bought).
9. **Offline tick time:** the device's `ticked_at` is stored if it is within the last 7 days and not in the future;
   otherwise server time.
10. Mutations older than 7 days (device clock or queue) are rejected with `FAILED_PRECONDITION`
    (`mutation_too_old`) and the phone resyncs.
11. Committed changes are announced (PLN-15-09) after the transaction commits, never before.

## Contract

All list mutation requests (PLN-15-05/06/07) carry `string client_mutation_id`. Responses carry the new
`list_version`.

## Data

`planner/migrations/000NN_list_mutations.sql`

```sql
CREATE TABLE list_mutations (
    household_id       uuid NOT NULL,
    client_mutation_id uuid NOT NULL,
    list_version       bigint NOT NULL,
    response           bytea NOT NULL,           -- serialized response
    created_at         timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (household_id, client_mutation_id)
);
CREATE INDEX list_mutations_created ON list_mutations (created_at);
```

A daily job deletes rows older than 7 days and line tombstones older than 7 days.

## Errors

Reason codes (in backticks) travel as trailing gRPC metadata `x-error-code` (README decision 18); version
conflicts use `ABORTED`.

| Case | gRPC status |
|---|---|
| Missing `client_mutation_id` | `INVALID_ARGUMENT` |
| Mutation older than 7 days | `FAILED_PRECONDITION` `mutation_too_old` |

## Acceptance criteria

- [ ] Same `AddListLine` sent 3 times (flaky network): one line, one version bump.
- [ ] Phone A offline ticks milk; phone B removes milk; A comes online: milk is ticked and in the fridge.
- [ ] Phone A sets rice to 2 kg, phone B to 1 kg, B arrives last: 1 kg, both phones converge after resync.
- [ ] 50 concurrent mutations on one list: versions 1…50, no gaps, no duplicates (race test).
