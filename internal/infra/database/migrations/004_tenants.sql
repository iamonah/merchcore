CREATE TYPE IF NOT EXISTS status_type AS ENUM ('maintenance','active', 'suspended', 'archived');
CREATE TYPE IF NOT EXISTS plan_type AS ENUM ('free','core', 'pro', 'enterprise');
CREATE TYPE IF NOT EXISTS business_mode AS ENUM ('online', 'physical', 'hybrid');

CREATE TABLE IF NOT EXISTS tenants (
    id                  UUID PRIMARY KEY,
    user_id             UUID NOT NULL REFERENCES users(id),
    business_name       VARCHAR(255) NOT NULL,
    domain              VARCHAR(255) NOT NULL,
    subdomain           VARCHAR(255) NOT NULL,
    logo_url            TEXT,
    plan                plan_type DEFAULT 'free',
    status              status_type DEFAULT 'active',
    business_mode       business_mode DEFAULT 'online',
    number_of_employees INTEGER,
    -- is_enabled          BOOLEAN DEFAULT true,
    trial_start_at      TIMESTAMPTZ DEFAULT NULL,
    trial_end_at        TIMESTAMPTZ DEFAULT NULL, --paid user
    created_at          TIMESTAMPTZ DEFAULT now(),
    updated_at          TIMESTAMPTZ DEFAULT now(),
    deleted_at          TIMESTAMPTZ DEFAULT NULL,

    CONSTRAINT domain_uq UNIQUE (domain),
    CONSTRAINT subdomain_uq UNIQUE (subdomain)
);


---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.

DROP TABLE IF EXISTS tenants;
DROP TYPE IF EXISTS status_type;
DROP TYPE IF EXISTS plan_type;
DROP TYPE IF EXISTS business_mode;
