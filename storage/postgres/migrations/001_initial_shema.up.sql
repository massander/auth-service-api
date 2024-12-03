create extension if not exists "uuid-ossp";

create table if not exists
    refresh_tokens (
        id uuid primary key,
        "token" text not null,
        user_id uuid not null,
        revoked bool not null default false,
        client_ip varchar(64) not null,
        created_at timestamptz null,
        updated_at timestamptz null
    );

create table if not exists
    access_tokens (
        id uuid not null primary key,
        parent_id uuid not null,
        user_id uuid not null,
        revoked bool not null default false,
        client_ip varchar(64) not null,
        created_at timestamptz null default now(),
        updated_at timestamptz null default now()
    );