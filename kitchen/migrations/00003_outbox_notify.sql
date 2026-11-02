-- +goose Up
-- Outbox relay (KIT-01-03): every statement that inserts outbox rows sends NOTIFY outbox when its transaction
-- commits, so the relay wakes at once instead of waiting for its 1-second poll. Notifications of one transaction
-- collapse into one.
-- +goose StatementBegin
CREATE FUNCTION outbox_notify() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN PERFORM pg_notify('outbox', ''); RETURN NULL; END $$;
-- +goose StatementEnd
CREATE TRIGGER outbox_notify AFTER INSERT ON outbox FOR EACH STATEMENT EXECUTE FUNCTION outbox_notify();

-- +goose Down
DROP TRIGGER outbox_notify ON outbox;
DROP FUNCTION outbox_notify();
