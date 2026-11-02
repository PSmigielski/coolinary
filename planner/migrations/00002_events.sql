-- +goose Up
-- Stream consumer (PLN-01-01): one row per applied event, written in the transaction of the handler that applied
-- it, so an event delivered again is acknowledged without running its handler twice. Rows older than 14 days are
-- deleted daily at 03:30 UTC; the streams keep about 7 days, so older ids can't come back. Same columns as kitchen's.
CREATE TABLE processed_events (
    event_id      uuid PRIMARY KEY,               -- the envelope id
    stream        text NOT NULL,
    type          text NOT NULL,
    processed_at  timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX processed_events_age ON processed_events (processed_at);

-- Entries the consumer gave up on: a handler failed 5 deliveries, the envelope's version is newer than its handler,
-- or the envelope is malformed. Kept for a person to look at. Keyed by stream entry (Q234), so an entry without valid
-- JSON or event id fits too.
CREATE TABLE dead_events (
    stream      text NOT NULL,
    entry_id    text NOT NULL,                    -- the stream entry id
    event_id    uuid,                             -- NULL when the envelope has no valid id
    type        text NOT NULL,                    -- the entry's type field
    envelope    text NOT NULL,                    -- the entry's envelope field exactly as received
    last_error  text NOT NULL,
    deliveries  int NOT NULL,
    created_at  timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (stream, entry_id)
);

-- +goose Down
DROP TABLE dead_events;
DROP TABLE processed_events;
