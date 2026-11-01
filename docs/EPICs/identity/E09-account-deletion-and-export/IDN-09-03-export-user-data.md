# IDN-09-03 `ExportUserData`

| | |
|---|---|
| Service | identity |
| Epic | [IDN-09 Account deletion and data export](README.md) |
| Release | R1 |
| Size | S |
| Status | todo |
| Depends on | IDN-04-04 |
| Unblocks | GW-02 (assembles the full export, D9) |

## Goal

A person can get everything identity stores about them, as part of the app-wide data export.

## Business rules

1. `ExportUserData(actor_user_id)` returns the actor's own data only:
   - account: user id, e-mail (empty for anonymous), e-mail verified, anonymous flag, created at;
   - sessions: created at, last used, user agent, IP address (as stored by Better Auth);
   - membership: household id and name, role, own display name, joined at;
   - the household's plan, plan source, expiry, and plan history (only when the actor is the owner);
   - invites the actor created: created at, status, expiry (no codes).
2. Other members' data is not included (not even their names); the household name is included.
3. The gateway calls it, merges it with the other services' exports and returns one JSON file (D9).
4. Rate limit: at most 3 exports per account per day, enforced in the gateway (GW-02).

## Contract

```proto
rpc ExportUserData(ExportUserDataRequest) returns (ExportUserDataResponse);

message ExportUserDataRequest { string actor_user_id = 1; }
message ExportUserDataResponse {
  string json = 1;   // UTF-8 JSON document, schema "coolinary.identity.export.v1"
}
```

Decided (Q138): a JSON string field; each service owns its export shape and the gateway nests the exports under the
service name.

## Data

No schema change.

## Errors

| Case | gRPC status |
|---|---|
| Unknown user | `NOT_FOUND` |

## Acceptance criteria

- [ ] The export of a member of a 3-account household contains no other member's name or e-mail.
- [ ] An owner's export includes plan history; a member's does not.
