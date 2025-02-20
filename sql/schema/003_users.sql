-- +goose Up
CREATE TABLE users (
    id UUID primary key,
    created_at timestamp not null,
    updated_at timestamp not null,
    email text unique not null,
    password text not null
);

-- +goose Down
Drop table users;