CREATE TYPE provider_type AS ENUM ('local', 'google');

CREATE TABLE IF NOT EXISTS users (
    user_id         UUID NOT NULL,
    password_hash   BYTEA,
    email           citext NOT NULL,
    first_name      VARCHAR(256) NOT NULL,
    last_name       VARCHAR(256) NOT NULL,
    provider_id     TEXT,
    phone_number    VARCHAR(15) NOT NULL,
    provider        provider_type NOT NULL DEFAULT 'local',   
    created_at      TIMESTAMPTZ DEFAULT now(),
    updated_at      TIMESTAMPTZ DEFAULT now(),
    is_verified     BOOLEAN DEFAULT false,
    country         VARCHAR(255)
    deleted_at      TIMESTAMPTZ,
    role_id         REFERENCES role_permissions(role_id)

    CONSTRAINT users_email_uq UNIQUE (email),
    CONSTRAINT users_provider_uq UNIQUE (provider_id),
    CONSTRAINT users_user_uq UNIQUE (user_id),
    CONSTRAINT users_phone_number_uq UNIQUE (phone_number)

    CONSTRAINT provider_fields_chk CHECK (
        (provider = 'local' AND password_hash IS NOT NULL)
        OR
        (provider = 'google' AND provider_id IS NOT NULL)
    )
);

CREATE TABLE IF NOT EXISTS addresses (
    id          BIGSERIAL PRIMARY KEY,
    user_id     UUID REFERENCES users(id),
    street      TEXT NOT NULL,
    city        TEXT NOT NULL,
    state       TEXT NOT NULL,
    postal_code TEXT NOT NULL,
    country     TEXT NOT NULL,
    is_default  BOOLEAN NOT NULL
);

CREATE INDEX IF NOT EXISTS users_id_idx ON users(id);
CREATE INDEX IF NOT EXISTS users_user_id_idx ON users(user_id);
CREATE INDEX IF NOT EXISTS users_email_idx ON users(email);
CREATE INDEX IF NOT EXISTS users_phone_number_idx ON users(phone_number);
CREATE INDEX IF NOT EXISTS users_provider_id_idx ON users(provider_id);

---- create above / drop below ----

DROP INDEX IF EXISTS users_id_idx;
DROP INDEX IF EXISTS users_user_id_idx;
DROP INDEX IF EXISTS users_email_idx;
DROP INDEX IF EXISTS users_phone_number_idx;
DROP INDEX IF EXISTS users_provider_id_idx;

DROP TABLE IF EXISTS users
DROP TABLE IF EXISTS addresses


