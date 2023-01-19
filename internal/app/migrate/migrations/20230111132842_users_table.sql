-- +goose Up
-- +goose StatementBegin

-- Crypto extension
create extension if not exists "uuid-ossp";
create extension if not exists "pgcrypto";

-- Users table
create table if not exists users
(
    id         uuid                         default gen_random_uuid(),
    login      varchar(100) unique not null,
    password   varchar(100)        not null,
    created_at timestamp           not null default now(),
    updated_at timestamp           not null default now(),
    primary key (id),
    unique (login)
);
-- +goose StatementEnd


