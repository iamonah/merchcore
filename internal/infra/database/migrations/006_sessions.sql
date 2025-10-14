CREATE TABLE IF NOT EXISTS sessions (
    id              UUID PRIMARY KEY,
    user_id         UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    refresh_token   BYTEA NOT NULL,
    user_agent      TEXT NOT NULL,
    client_ip       TEXT NOT NULL,
    is_blocked      BOOLEAN NOT NULL DEFAULT false,
    expires_at      TIMESTAMPTZ NOT NULL,
    created_at      TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX IF NOT EXISTS sessions_refresh_token_idx ON sessions(refresh_token);

CREATE TABLE IF NOT EXISTS tokens(
    hash BYTEA PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    expiry TIMESTAMPTZ NOT NULL,
    scope TEXT NOT NULL
);

---- create above / drop below ----

DROP INDEX IF EXISTS sessions_refresh_token_idx;
DROP TABLE IF EXISTS tokens;
DROP TABLE IF EXISTS sessions;