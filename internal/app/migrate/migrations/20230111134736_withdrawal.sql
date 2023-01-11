-- +goose Up
-- +goose StatementBegin

-- Withdrawals table
create table if not exists withdrawals
(
    id         integer,
    account_id varchar(100) not null,
    order_id   varchar(100) not null,
    points     numeric not null ,
    updated_at timestamp    not null default now(),
    primary key (id),
    unique (order_id)
);

-- +goose StatementEnd