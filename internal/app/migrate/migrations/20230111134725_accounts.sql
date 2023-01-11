-- +goose Up
-- +goose StatementBegin

-- Accounts table
create table if not exists accounts
(
    id         varchar(100) not null,
    balance    numeric not null,
    updated_at timestamp    not null default now(),
    primary key (id),
    unique (id)
);

-- +goose StatementEnd
