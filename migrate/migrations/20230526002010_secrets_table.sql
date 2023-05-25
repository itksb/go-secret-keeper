-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS secrets
(
    id      uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id uuid         NOT NULL,
    data    bytea        NOT NULL,
    meta    bytea        NOT NULL,
    type    varchar(128) NOT NULL
);
CREATE INDEX IF NOT EXISTS secrets_user_id_idx ON secrets (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS secrets;
-- +goose StatementEnd
