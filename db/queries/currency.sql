-- name: ListCurrencies :many
SELECT id, code FROM currencies;

-- name: UpsertExchangeRate :exec
INSERT INTO exchange_rates (currency_id, rate, updated_at)
VALUES (@currency_id, @rate, NOW())
ON CONFLICT (currency_id)
DO UPDATE SET rate = EXCLUDED.rate, updated_at = NOW();
