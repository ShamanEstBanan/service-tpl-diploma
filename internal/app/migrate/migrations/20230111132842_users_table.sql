-- +goose Up
-- +goose StatementBegin

-- Users table
create table if not exists users
(
    id         varchar(100) not null,
    login      varchar(100) not null,
    password   varchar(100) not null,
    created_at timestamp    not null default now(),
    updated_at timestamp    not null default now(),
    primary key (id),
    unique (login)
);
-- +goose StatementEnd


