-- +goose Up
-- +goose StatementBegin
create table users (
    id integer primary key,

    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    deleted_at timestamp default null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table users;
-- +goose StatementEnd
