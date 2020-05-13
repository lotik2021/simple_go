ALTER TABLE IF EXISTS maasapi.device_movista_search_history ADD COLUMN IF NOT EXISTS currency_code VARCHAR NOT NULL DEFAULT 'RUB';
ALTER TABLE IF EXISTS maasapi.device_movista_search_history ADD COLUMN IF NOT EXISTS culture_code VARCHAR NOT NULL DEFAULT 'ru';
ALTER TABLE IF EXISTS maasapi.device_movista_search_history ADD COLUMN IF NOT EXISTS comfort_type VARCHAR NOT NULL DEFAULT 'economy';
CREATE TABLE IF NOT EXISTS maasapi.customer (
  id serial PRIMARY KEY,
  age integer DEFAULT 36,
  seat_required boolean DEFAULT true,
  movista_search_id varchar NOT NULL
);

INSERT INTO maasapi.customer(movista_search_id)
  SELECT id from maasapi.device_movista_search_history;
