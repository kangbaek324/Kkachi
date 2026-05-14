CREATE TYPE user_role AS ENUM ('user', 'admin');

CREATE TABLE users (
    id          BIGSERIAL    PRIMARY KEY,
    username    VARCHAR(50)  NOT NULL UNIQUE,
    password    VARCHAR(255) NOT NULL,
    role        VARCHAR(20)  NOT NULL DEFAULT 'user',
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE currencies (
    id      BIGSERIAL    PRIMARY KEY,
    code    VARCHAR(10)  NOT NULL UNIQUE,
    name    VARCHAR(50)  NOT NULL,
    unit    NUMERIC      NOT NULL
);

CREATE TABLE exchange_rates (
    id          BIGSERIAL      PRIMARY KEY,
    currency_id BIGINT         NOT NULL REFERENCES currencies(id),
    rate        NUMERIC(18, 6) NOT NULL,
    updated_at  TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE TABLE wallets (
    id             BIGSERIAL    PRIMARY KEY,
    user_id        BIGINT       NOT NULL REFERENCES users(id),
    wallet_number  VARCHAR(20)  NOT NULL UNIQUE,
    nickname       VARCHAR(50),
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE balances (
    id          BIGSERIAL      PRIMARY KEY,
    wallet_id   BIGINT         NOT NULL REFERENCES wallets(id),
    currency_id BIGINT         NOT NULL REFERENCES currencies(id),
    amount      NUMERIC(18, 6) NOT NULL DEFAULT 0,
    UNIQUE (wallet_id, currency_id)
);

CREATE TABLE refresh_tokens (
    id         BIGSERIAL    PRIMARY KEY,
    user_id    BIGINT       NOT NULL REFERENCES users(id),
    token      VARCHAR(512) NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ  NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);