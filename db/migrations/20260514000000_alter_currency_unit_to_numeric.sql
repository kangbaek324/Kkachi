-- +goose Up
ALTER TABLE currencies ALTER COLUMN unit TYPE NUMERIC USING unit::numeric;

-- +goose Down
ALTER TABLE currencies ALTER COLUMN unit TYPE VARCHAR(10) USING unit::text;
