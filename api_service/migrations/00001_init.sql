-- +goose Up
-- +goose StatementBegin
CREATE TYPE user_role AS ENUM ('admin', 'username');
CREATE TABLE users (
    id SERIAL primary key,
    name varchar unique not null,
    password varchar unique not null,
    created_dt timestamp with time zone default now() not null,
    role user_role not null
);

CREATE TABLE services (
    id SERIAL primary key,
    name varchar unique not null,
    owner_id integer references users(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE services;
DROP TABLE users;
DROP TYPE user_role;
-- +goose StatementEnd
