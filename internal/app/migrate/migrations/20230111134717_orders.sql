-- +goose Up
-- +goose StatementBegin

-- Orders table
create table if not exists orders
(
    id          varchar(100)             not null,
    account_id  varchar(50)              not null,
    status      varchar(20)              not null default 'NEW',
    accrual     numeric                  not null default 0,
    uploaded_at timestamp with time zone not null default now(),
    updated_at  timestamp with time zone not null default now(),
    primary key (id),
    unique (id)
);

-- +goose StatementEnd
