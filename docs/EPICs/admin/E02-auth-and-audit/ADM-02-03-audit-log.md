# ADM-02-03 Audit log

| | |
|---|---|
| Service | admin |
| Epic | [ADM-02 Admin sign-in, roles and audit log](README.md) |
| Release | R1 |
| Size | M |
| Depends on | ADM-01-03, ADM-02-02 |
| Unblocks | every mutating admin route |

## Goal

Every change made through the back office, and every look at personal data, is recorded: who, what, when, and
whether it worked (GDPR accountability; a record of how the catalogue and library were curated).

## Business rules

1. Every mutating route (`POST`, `PUT`, `PATCH`, `DELETE`) writes an audit row with status `pending` before the
   downstream call and sets `ok` or `error` (with the gRPC code) after it.
2. If the `pending` row cannot be written (Postgres down), the route answers 503 and makes no downstream call.
3. Reads of personal data are audited too, as action kind `read`: household lookup (ADM-05-01), delivery log search
   by e-mail (ADM-07-04), account statistics are not personal and not audited.
4. A row stores: `admin_user_id`, `roles`, `action` (`<service>.<object>.<verb>`, e.g. `kitchen.match.resolve`),
   `target_service`, `target_kind`, `target_id`, `request_id`, `status`, `grpc_code`, `details` (JSON with ids,
   counts, changed field names and their new values for catalogue, dish and offer edits), timestamps.
5. Never stored in `details`: e-mail addresses (a lookup by e-mail stores the HMAC used by identity/notifier, or
   just "by e-mail"), member display names, free text a user wrote. Admin's own notes (e.g. a review note) are
   allowed.
6. Rows are never edited after they are finished and never deleted by the application. Retention: rows older than
   2 years are deleted by a monthly job. Decided (Q047): 2 years.
7. `GET /admin/v1/audit` (role `admin`): filters `admin_user_id`, `action` prefix, `target_kind`, `target_id`,
   `from`, `to`; newest first; 50 per page with an opaque page token.

## Contract

| Method | Path | Response |
|---|---|---|
| GET | `/admin/v1/audit?admin_user_id=&action=&target_kind=&target_id=&from=&to=&page_token=` | `{ items: AuditEntry[], next_page_token }` |

## Data

`admin/src/db/schema.ts`

```sql
CREATE TABLE admin_audit (
  id             uuid PRIMARY KEY,          -- UUIDv7
  created_at     timestamptz NOT NULL DEFAULT now(),
  finished_at    timestamptz,
  admin_user_id  text NOT NULL,
  roles          text[] NOT NULL,
  kind           text NOT NULL CHECK (kind IN ('write', 'read', 'auth')),
  action         text NOT NULL,
  target_service text NOT NULL,
  target_kind    text NOT NULL DEFAULT '',
  target_id      text NOT NULL DEFAULT '',
  request_id     text NOT NULL,
  status         text NOT NULL CHECK (status IN ('pending', 'ok', 'error')),
  grpc_code      text,
  details        jsonb NOT NULL DEFAULT '{}'
);
CREATE INDEX admin_audit_created ON admin_audit (created_at DESC);
CREATE INDEX admin_audit_target ON admin_audit (target_kind, target_id);
CREATE INDEX admin_audit_admin ON admin_audit (admin_user_id, created_at DESC);
```

## Errors

| Case | HTTP |
|---|---|
| Audit write fails before a change | 503 `audit_unavailable` |

## Acceptance criteria

- [ ] Approving a match creates one row `kitchen.match.resolve` `ok` with the match id and chosen ingredient id.
- [ ] A downstream `FAILED_PRECONDITION` gives status `error`, `grpc_code = failed_precondition`.
- [ ] With admin's Postgres stopped, a write route returns 503 and kitchen receives no call.
- [ ] A household lookup by e-mail creates a `read` row without the address.
