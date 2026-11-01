-- name: GetProductByEAN :one
SELECT ean, name, brand, pack_value, pack_unit
FROM products
WHERE ean = $1;
