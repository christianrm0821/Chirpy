-- +goose Up
alter table users
add is_chirpy_red boolean not null
default false;

-- +goose Down
alter table users
drop column is_chirpy_red;