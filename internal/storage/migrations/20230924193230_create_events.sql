-- +goose Up
-- +goose StatementBegin
CREATE TABLE events (
    num_transaction text NOT NULL,
    wallet_id int NOT NULL,
    status int not null ,
    type_transaction int,
    data text,
    CREATED_AT timestamp NOT NULL DEFAULT NOW()

);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS events;
-- +goose StatementEnd
