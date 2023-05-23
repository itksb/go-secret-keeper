-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS auth_accounts
(
    id       uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    login    varchar(100) NOT NULL UNIQUE,
    password varchar(255) NOT NULL
);

CREATE INDEX IF NOT EXISTS auth_accounts_login_pwd ON auth_accounts (login, password);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS auth_accounts;
-- +goose StatementEnd
