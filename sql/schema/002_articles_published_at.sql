-- +goose Up
alter table articles
add column published_at timestamp not null;

-- +goose Down
drop table articles;