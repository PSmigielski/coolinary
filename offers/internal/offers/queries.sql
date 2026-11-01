-- name: ListCurrentOffers :many
SELECT o.id, o.shop_id, o.product_name, o.price_minor, o.currency, o.conditions,
       o.valid_from, o.valid_to, o.source_url
FROM offers o
JOIN shops s ON s.id = o.shop_id
WHERE s.country_code = $1
  AND o.valid_to >= CURRENT_DATE
  AND o.id > $2
ORDER BY o.id
LIMIT $3;
