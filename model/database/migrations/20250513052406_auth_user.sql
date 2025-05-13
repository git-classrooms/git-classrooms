-- +goose Up
create table "public"."auth_users" (
    "id" uuid primary key default uuid_generate_v4(),
    "name" text not null,
    "username" text not null,
    "picture" text not null,
    "email" text not null,
    "email_verified" boolean not null default false
);

create table "public"."auth_accounts" (
    "id" uuid primary key default uuid_generate_v4(),
    "provider" text not null,
    "provider_id" text not null,
    "email" text not null,
    "email_verified" boolean not null default false,
    "access_token" text not null,
    "expiry" timestamp with time zone not null,
    "refresh_token" text not null,
    "refresh_expiry" timestamp with time zone not null,
    "token_type" text not null,
    "id_token" text not null,
    "name" text not null,
    "preferred_username" text not null,
    "nickname" text not null,
    "picture" text not null,
    "profile" text not null,
    "user_id" uuid not null,
    constraint "fk_auth_users_accounts" foreign key ("user_id") references "public"."auth_users"("id") on delete cascade
);

create unique index idx_unique_provider_provider_id on public.auth_accounts using btree (provider, provider_id);
create unique index idx_unique_provider_user_id on public.auth_accounts using btree (provider, user_id);

-- +goose Down
drop table "auth_accounts";
drop table "auth_users";
