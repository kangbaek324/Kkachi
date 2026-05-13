-- name: ListCurrencies :many
SELECT id, code FROM currencies;

-- name: ListCurrenciesWithRate :many
SELECT
    c.id,
    c.code,
    c.name,
    c.unit,
    er.rate,
    er.updated_at
FROM currencies c
LEFT JOIN exchange_rates er ON er.currency_id = c.id;

-- name: UpsertExchangeRate :exec
INSERT INTO exchange_rates (currency_id, rate, updated_at)
VALUES (@currency_id, @rate, NOW())
ON CONFLICT (currency_id)
DO UPDATE SET rate = EXCLUDED.rate, updated_at = NOW();
