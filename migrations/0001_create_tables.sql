-- +goose Up
CREATE TABLE brand (
  id int GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  name TEXT NOT NULL,
  UNIQUE(id, name)
);

CREATE TABLE price (
  id int GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY, -- renamed from price_list
  brand_id INTEGER NOT NULL,
  start_date TIMESTAMP WITHOUT TIME ZONE NOT NULL, -- use UTC for times
  end_date TIMESTAMP WITHOUT TIME ZONE NOT NULL, -- use UTC for times
  product_id INTEGER NOT NULL,
  priority INTEGER NOT NULL,
  price INTEGER NOT NULL, -- lowest unit, eg: cents in USD, yen in JPY
  curr TEXT NOT NULL,
  CHECK (start_date <= end_date), -- prevent start dates being after end dates, could be managed in package storage
  CONSTRAINT fk_brand_id
    FOREIGN KEY(brand_id)
	    REFERENCES brand(id)
);

CREATE INDEX brand_ix_id ON brand (id);
CREATE INDEX price_ix_product_id ON price (product_id);

select column_name, data_type, character_maximum_length, column_default, is_nullable from INFORMATION_SCHEMA.COLUMNS where table_name = 'brand';

select column_name, data_type, character_maximum_length, column_default, is_nullable from INFORMATION_SCHEMA.COLUMNS where table_name = 'price';

-- +goose Down
DROP TABLE IF EXISTS price;
DROP TABLE IF EXISTS brand;
