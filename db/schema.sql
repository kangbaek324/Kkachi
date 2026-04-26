CREATE TYPE user_role AS ENUM ('user', 'admin');

CREATE TABLE users (
    id         bigserial PRIMARY KEY,
    username   varchar(15)  NOT NULL,
    password   varchar(60)  NOT NULL,
    role       user_role    NOT NULL DEFAULT 'user',
    created_at timestamptz  NOT NULL DEFAULT now(),
    CONSTRAINT users_username_unique UNIQUE (username)
);

CREATE TABLE currencies (
    id         bigserial PRIMARY KEY,
    code       varchar(15)  NOT NULL,
    name       varchar(15)  NOT NULL,
    unit       int          NOT NULL DEFAULT 1,
    created_at timestamptz  NOT NULL DEFAULT now(),
    CONSTRAINT currencies_code_unique UNIQUE (code)
);

CREATE TABLE exchange_rates (
    id           bigserial PRIMARY KEY,
    currency_id  bigint       NOT NULL REFERENCES currencies(id),
    rate         decimal(10,4) NOT NULL,
    last_updated timestamptz  NOT NULL DEFAULT now()
);

CREATE TABLE wallets (
    id             bigserial PRIMARY KEY,
    account_number varchar(10)  NOT NULL,
    user_id        bigint       NOT NULL REFERENCES users(id),
    nickname       varchar(15)  NULL,
    created_at     timestamptz  NOT NULL DEFAULT now(),
    CONSTRAINT wallets_account_number_unique UNIQUE (account_number)
);

CREATE TABLE balances (
    account_id  bigint        NOT NULL REFERENCES wallets(id),
    currency_id bigint        NOT NULL REFERENCES currencies(id),
    balance     decimal(19,4) NOT NULL DEFAULT 0,
    PRIMARY KEY (account_id, currency_id)
);
