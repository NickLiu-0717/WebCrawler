-- +goose Up
CREATE TABLE articles (
    id UUID primary key,
    url text unique not null,
    title text not null,
    content text not null,
    catagory text not null,
    image_url text,
    created_at timestamp not null
);

-- +goose Down
Drop table articles;