-- +goose Up
ALTER TABLE exchange_rates ADD CONSTRAINT exchange_rates_currency_id_unique UNIQUE (currency_id);

-- +goose Down
ALTER TABLE exchange_rates DROP CONSTRAINT exchange_rates_currency_id_unique;
