-- +goose Up
-- Admin change log (PLN-01-06): every PlannerAdminService RPC that changes data writes one row in its own
-- transaction: who (the admin's user id), which RPC, the target ("dish:<id>", "import:<id>", ...) and a short JSON
-- summary of ids and counts, never free text from users. Rows older than 1 year are deleted by a monthly job (Q016).
CREATE TABLE admin_changes (
    id             bigserial PRIMARY KEY,
    admin_user_id  text NOT NULL,
    rpc            text NOT NULL,
    target         text NOT NULL,
    summary        jsonb NOT NULL DEFAULT '{}',
    created_at     timestamptz NOT NULL DEFAULT now()
);
-- One target's history, newest first, paged by id.
CREATE INDEX admin_changes_target ON admin_changes (target, id DESC);
-- The monthly cleanup.
CREATE INDEX admin_changes_age ON admin_changes (created_at);

-- +goose Down
DROP TABLE admin_changes;
