-- name: ListDishesForMeal :many
SELECT id, name, kcal_per_portion
FROM dishes
WHERE meal = $1 OR $1 = 0
ORDER BY name
LIMIT $2;

-- name: InsertProcessedEvent :execrows
-- The stream consumer records an event in its handler's transaction. 0 rows: the event was already applied (an
-- uncommitted row of another transaction makes this wait for that transaction's end).
INSERT INTO processed_events (event_id, stream, type)
VALUES ($1, $2, $3)
ON CONFLICT (event_id) DO NOTHING;

-- name: DeleteProcessedEvents :execrows
DELETE FROM processed_events WHERE processed_at < sqlc.arg(cutoff)::timestamptz;

-- name: InsertDeadEvent :exec
-- The stream consumer parks an entry it gives up on, then acknowledges it. An entry already parked (its XACK failed
-- before) is kept as it is.
INSERT INTO dead_events (stream, entry_id, event_id, type, envelope, last_error, deliveries)
VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (stream, entry_id) DO NOTHING;
