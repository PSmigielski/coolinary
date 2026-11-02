-- +goose Up
-- Admin change records (KIT-01-07): every KitchenAdminService RPC that changes data writes one row in its own
-- transaction: who (the admin's user id), what (the RPC name), which entity, and the entity's catalogue fields
-- before and after. Never household data. Rows older than 3 years are deleted by a monthly job.
CREATE TABLE admin_changes (
    id             bigserial PRIMARY KEY,
    admin_user_id  text NOT NULL,
    action         text NOT NULL,
    entity_type    text NOT NULL,
    entity_id      text NOT NULL,
    before         jsonb,                          -- NULL when the change created the entity
    after          jsonb,                          -- NULL when the change deleted it
    created_at     timestamptz NOT NULL DEFAULT now()
);
-- One entity's history, newest first, paged by id.
CREATE INDEX admin_changes_entity ON admin_changes (entity_type, entity_id, id DESC);
-- The monthly cleanup.
CREATE INDEX admin_changes_age ON admin_changes (created_at);

-- +goose Down
DROP TABLE admin_changes;
