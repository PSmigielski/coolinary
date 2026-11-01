-- +goose Up
-- Recipe library entries the planner ranks. The full library (ingredients, tools,
-- cuisines, nutrition) arrives with the planner feature; this is the minimum to serve SuggestDishes.
CREATE TABLE dishes (
    id               text PRIMARY KEY,
    name             text NOT NULL,
    meal             smallint NOT NULL,         -- coolinary.planner.v1.Meal
    kcal_per_portion integer NOT NULL,
    created_at       timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX dishes_meal ON dishes (meal);

INSERT INTO dishes (id, name, meal, kcal_per_portion) VALUES
  ('owsianka',      'Porridge with apple and cinnamon',       1, 380),
  ('shakshuka',     'Shakshuka',                              1, 420),
  ('chicken-curry', 'Chicken curry with rice',                2, 520),
  ('pierogi',       'Pierogi ruskie',                         2, 610),
  ('salad-bowl',    'Grain bowl with roasted vegetables',     3, 450);

-- +goose Down
DROP TABLE dishes;
