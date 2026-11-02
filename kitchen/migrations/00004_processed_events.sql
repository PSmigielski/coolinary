-- +goose Up
-- Stream consumer (KIT-01-04): one row per applied event, written in the transaction of the handler that applied
-- it, so an event delivered again is acknowledged without running its handler twice. Rows older than 14 days are
-- deleted daily; the streams keep about 7 days, so older ids can't come back.
CREATE TABLE processed_events (
    event_id      uuid PRIMARY KEY,               -- the envelope id
    stream        text NOT NULL,
    type          text NOT NULL,
    processed_at  timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX processed_events_age ON processed_events (processed_at);

-- +goose Down
DROP TABLE processed_events;
