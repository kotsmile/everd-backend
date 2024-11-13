-- +goose Up
-- +goose StatementBegin
create table todos (
    id integer primary key,
    
    title varchar(100) not null,
    comment varchar(1000) not null default '',
    done boolean not null default false,

    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table todos;
-- +goose StatementEnd
