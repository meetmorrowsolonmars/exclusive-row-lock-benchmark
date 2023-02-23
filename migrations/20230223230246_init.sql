-- +goose Up
-- +goose StatementBegin
create table mutex
(
    id     serial primary key,
    amount numeric not null
);

create table optimistic_lock
(
    id     serial primary key,
    amount numeric not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table mutex;
drop table optimistic_lock;
-- +goose StatementEnd
