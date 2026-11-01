-- +goose Up
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Shared product catalogue: every product a household has confirmed.
CREATE TABLE products (
    ean          text PRIMARY KEY,
    name         text NOT NULL,
    brand        text NOT NULL DEFAULT '',
    pack_value   double precision,
    pack_unit    smallint NOT NULL DEFAULT 0,   -- coolinary.common.v1.Unit
    confirmed    boolean NOT NULL DEFAULT false,
    created_at   timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX products_name_trgm ON products USING gin (name gin_trgm_ops);

-- Example row (prefix 952 is reserved by GS1 for examples).
INSERT INTO products (ean, name, brand, pack_value, pack_unit, confirmed)
VALUES ('9520000000017', 'Oat drink 1 l', 'Example', 1000, 2, true);

-- +goose Down
DROP TABLE products;
