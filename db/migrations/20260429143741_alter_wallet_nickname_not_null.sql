-- +goose Up
ALTER TABLE wallets ALTER COLUMN nickname SET NOT NULL;

-- +goose Down
ALTER TABLE wallets ALTER COLUMN nickname DROP NOT NULL;
