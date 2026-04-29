-- name: CreateWallet :one
INSERT INTO wallets (user_id, wallet_number, nickname) VALUES($1, $2, $3) RETURNING *;