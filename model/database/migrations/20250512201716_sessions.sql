-- +goose Up
CREATE TABLE sessions (
  token TEXT PRIMARY KEY,
  data BYTEA NOT NULL,
  expiry TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_sessions_expiry ON public.sessions USING btree (expiry);

create table "public"."token_sessions" (
    "session_id" text primary key,
    "refresh_jti" text not null,
    "expiry" timestamp with time zone not null
);

CREATE INDEX idx_token_sessions_expiry ON public.token_sessions USING btree (expiry);
CREATE INDEX idx_token_sessions_refresh_jti ON public.token_sessions USING btree (refresh_jti);

-- +goose Down
DROP TABLE sessions;
DROP TABLE token_sessions;
