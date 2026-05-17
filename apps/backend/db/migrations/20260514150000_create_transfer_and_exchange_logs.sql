-- +goose Up
CREATE TABLE transfer_logs (
    id               BIGSERIAL      PRIMARY KEY,
    from_wallet_id   BIGINT         NOT NULL REFERENCES wallets(id),
    to_wallet_id     BIGINT         NOT NULL REFERENCES wallets(id),
    currency_id      BIGINT         NOT NULL REFERENCES currencies(id),
    amount           NUMERIC(18, 6) NOT NULL,
    transferred_at   TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_transfer_logs_from_wallet ON transfer_logs(from_wallet_id, transferred_at DESC);
CREATE INDEX idx_transfer_logs_to_wallet   ON transfer_logs(to_wallet_id, transferred_at DESC);

CREATE TABLE exchange_logs (
    id              BIGSERIAL      PRIMARY KEY,
    wallet_id       BIGINT         NOT NULL REFERENCES wallets(id),
    from_currency_id BIGINT        NOT NULL REFERENCES currencies(id),
    to_currency_id   BIGINT        NOT NULL REFERENCES currencies(id),
    from_amount     NUMERIC(18, 6) NOT NULL,
    to_amount       NUMERIC(18, 6) NOT NULL,
    from_rate       NUMERIC(18, 6) NOT NULL,
    from_unit       NUMERIC(18, 6) NOT NULL,
    to_rate         NUMERIC(18, 6) NOT NULL,
    to_unit         NUMERIC(18, 6) NOT NULL,
    krw_amount      NUMERIC(18, 6) NOT NULL,
    exchanged_at    TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_exchange_logs_wallet ON exchange_logs(wallet_id, exchanged_at DESC);

-- +goose Down
DROP TABLE exchange_logs;
DROP TABLE transfer_logs;
