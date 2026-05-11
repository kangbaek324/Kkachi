-- name: CreateWallet :one
INSERT INTO wallets (user_id, wallet_number, nickname) VALUES($1, $2, $3) RETURNING *;

-- name: GetWalletByWalletNumber :one
SELECT user_id, wallet_number, nickname FROM wallets WHERE wallet_number = $1;

-- name: GetWallets :many
SELECT wallet_number, nickname FROM wallets WHERE user_id = $1;