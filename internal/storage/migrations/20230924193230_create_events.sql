-- +goose Up
-- +goose StatementBegin
CREATE TABLE events (
    id serial primary key,
    num_transaction uuid NOT NULL,
    wallet_id int NOT NULL,
    status int not null ,
    CREATED_AT timestamp NOT NULL DEFAULT NOW()

);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS events;
-- +goose StatementEnd
