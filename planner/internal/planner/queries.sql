-- name: ListDishesForMeal :many
SELECT id, name, kcal_per_portion
FROM dishes
WHERE meal = $1 OR $1 = 0
ORDER BY name
LIMIT $2;
