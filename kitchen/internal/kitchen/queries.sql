-- name: GetProductByEAN :one
SELECT ean, name, brand, pack_value, pack_unit
FROM products
WHERE ean = $1;

-- name: InsertOutboxEvent :exec
-- Only events.Add calls this, always with a transaction's Queries (q.WithTx(tx)).
INSERT INTO outbox (id, type, version, household_id, payload, occurred_at)
VALUES ($1, $2, $3, $4, $5, $6);
