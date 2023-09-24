-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    wallet_id SERIAL PRIMARY KEY,
    currency_code int,
    actual_balance double precision,
    frozen_balance double precision
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
