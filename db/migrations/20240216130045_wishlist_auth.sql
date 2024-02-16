-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS wish_users(
    uid SERIAL PRIMARY KEY,
    lname VARCHAR NOT NULL,
    fname VARCHAR NOT NULL,
    email VARCHAR NOT NULL UNIQUE,
    pwd_hash VARCHAR NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE wish_users;
-- +goose StatementEnd
