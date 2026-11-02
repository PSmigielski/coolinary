-- +goose Up
-- Transactional outbox (KIT-01-02): events.Add writes one row in the transaction of the change it describes;
-- the relay (KIT-01-03) publishes unpublished rows in id order to the events.kitchen stream.
-- household_id has no foreign key: an event may be written in the transaction that deletes its household.
CREATE TABLE outbox (
    id            uuid PRIMARY KEY,               -- UUIDv7, also the envelope id
    type          text NOT NULL,
    version       integer NOT NULL DEFAULT 1,
    household_id  uuid,                           -- NULL for catalogue-wide events
    payload       jsonb NOT NULL,
    occurred_at   timestamptz NOT NULL,
    published_at  timestamptz,
    attempts      integer NOT NULL DEFAULT 0,
    last_error    text
);
CREATE INDEX outbox_unpublished ON outbox (id) WHERE published_at IS NULL;

-- +goose Down
DROP TABLE outbox;
