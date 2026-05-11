-- name: CreateWallet :one
INSERT INTO wallets (user_id, wallet_number, nickname) VALUES($1, $2, $3) RETURNING *;

-- name: GetWalletByWalletNumber :one
SELECT user_id, wallet_number, nickname FROM wallets WHERE wallet_number = $1;

-- name: GetWallets :many
SELECT wallet_number, nickname FROM wallets WHERE user_id = $1;

-- name: GetWalletBalance :many
SELECT 
    w.user_id,
    c.code,
    c.name,
    COALESCE(b.amount, 0) AS amount
FROM wallets w
JOIN currencies c ON true
LEFT JOIN balances b 
    ON b.currency_id = c.id
    AND b.wallet_id = w.id
WHERE w.wallet_number = $1
ORDER BY COALESCE(b.amount, 0) DESC;
