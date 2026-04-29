-- +goose Up
ALTER TABLE wallets RENAME COLUMN account_number TO wallet_number;

-- +goose Down
ALTER TABLE wallets RENAME COLUMN wallet_number TO account_number;
