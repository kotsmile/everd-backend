-- +goose Up
-- +goose StatementBegin
create table todolist (
    user_id integer not null,
    todo_id integer not null unique,

    constraint fk_todo foreign key (todo_id)
        references todos (id)
        on delete cascade
        on update cascade,

    constraint fk_user foreign key (user_id)
        references users (id)
        on delete cascade
        on update cascade,

    primary key (user_id, todo_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table todolist;
-- +goose StatementEnd
