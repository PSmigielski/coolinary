-- name: GetProductByEAN :one
SELECT ean, name, brand, pack_value, pack_unit
FROM products
WHERE ean = $1;

-- name: InsertOutboxEvent :exec
-- Only events.Add calls this, always with a transaction's Queries (q.WithTx(tx)).
INSERT INTO outbox (id, type, version, household_id, payload, occurred_at)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: ListenOutbox :exec
-- The relay's own connection subscribes to the NOTIFY the outbox trigger sends on commit.
LISTEN outbox;

-- name: LockUnpublishedOutbox :many
-- One relay pass: the oldest unpublished rows in id order, locked until the pass's transaction ends. Rows another
-- relay holds are skipped.
SELECT id, type, version, household_id, payload, occurred_at, published_at, attempts, last_error
FROM outbox
WHERE published_at IS NULL
ORDER BY id
LIMIT $1
FOR UPDATE SKIP LOCKED;

-- name: MarkOutboxPublished :exec
UPDATE outbox SET published_at = now() WHERE id = ANY(sqlc.arg(ids)::uuid[]);

-- name: MarkOutboxFailed :exec
UPDATE outbox SET attempts = attempts + 1, last_error = $2 WHERE id = $1;

-- name: OutboxBacklog :one
SELECT count(*) AS pending, min(occurred_at)::timestamptz AS oldest
FROM outbox
WHERE published_at IS NULL;

-- name: DeletePublishedOutbox :execrows
-- At most 10,000 published rows older than the cutoff per run; unpublished rows are never deleted.
DELETE FROM outbox
WHERE id IN (SELECT old.id FROM outbox AS old WHERE old.published_at < sqlc.arg(cutoff)::timestamptz LIMIT 10000);
