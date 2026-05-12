-- name: CreateWallet :one
INSERT INTO wallets (user_id, wallet_number, nickname) VALUES($1, $2, $3) RETURNING *;

-- name: GetWallet :one
SELECT user_id, wallet_number, nickname FROM wallets WHERE wallet_number = $1;

-- name: GetWallets :many
SELECT wallet_number, nickname FROM wallets WHERE user_id = $1;

-- name: GetWalletBalances :many
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

-- name: GetWalletBalance :one
SELECT
    w.user_id,
    COALESCE(b.amount, 0) AS amount
FROM wallets w
JOIN currencies c ON c.code = $2
LEFT JOIN balances b
    ON b.currency_id = c.id
    AND b.wallet_id = w.id
WHERE w.wallet_number = $1;

-- name: UpsertBalance :exec
INSERT INTO balances (wallet_id, currency_id, amount)
VALUES ($1, $2, $3)
ON CONFLICT (wallet_id, currency_id)
DO UPDATE SET amount = balances.amount + EXCLUDED.amount;

-- name: GetWalletBalanceLock :one
SELECT
    w.id AS wallet_id,
    w.user_id,
    c.id AS currency_id,
    COALESCE(b.amount, 0) AS amount
FROM wallets w
JOIN currencies c ON c.code = $2
LEFT JOIN balances b
    ON b.currency_id = c.id
    AND b.wallet_id = w.id
WHERE w.wallet_number = $1
FOR UPDATE OF w;
