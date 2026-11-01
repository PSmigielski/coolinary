-- +goose Up
-- Shops and offers. Flyer sync, shop import from OpenStreetMap and editions by area
-- arrive with the offers feature; this is the minimum to serve ListOffers.
CREATE TABLE shops (
    id           text PRIMARY KEY,
    chain        text NOT NULL,
    country_code text NOT NULL,
    name         text NOT NULL
);

CREATE TABLE offers (
    id            text PRIMARY KEY,
    shop_id       text NOT NULL REFERENCES shops (id),
    product_name  text NOT NULL,
    price_minor   bigint NOT NULL CHECK (price_minor > 0),
    currency      text NOT NULL,
    conditions    text NOT NULL DEFAULT '',
    valid_from    date NOT NULL,
    valid_to      date NOT NULL CHECK (valid_to >= valid_from),
    source_url    text NOT NULL
);
CREATE INDEX offers_valid ON offers (valid_to);

INSERT INTO shops (id, chain, country_code, name) VALUES
  ('example-pl-1', 'Example Market', 'PL', 'Example Market Warszawa');
INSERT INTO offers (id, shop_id, product_name, price_minor, currency, conditions, valid_from, valid_to, source_url) VALUES
  ('example-1', 'example-pl-1', 'Chicken thighs 1 kg', 1599, 'PLN', '', CURRENT_DATE, CURRENT_DATE + 6, 'https://example.com/flyer'),
  ('example-2', 'example-pl-1', 'Basmati rice 1 kg',    699, 'PLN', '2+1', CURRENT_DATE, CURRENT_DATE + 6, 'https://example.com/flyer');

-- +goose Down
DROP TABLE offers;
DROP TABLE shops;
