CREATE TYPE status_type AS ENUM ('maintainance','active', 'suspended', 'archived');
CREATE TYPE plan_type AS ENUM ('free','core', 'pro', 'enterprise');
CREATE TYPE business_mode AS ENUM ('online', 'physical', 'hybrid');

CREATE TABLE IF NOT EXISTS tenants (
    id                  UUID PRIMARY KEY,
    user_id             UUID NOT NULL REFERENCES users(user_id),
    business_name       VARCHAR(255) NOT NULL,
    domain              VARCHAR(255) UNIQUE NOT NULL,
    subdomain           VARCHAR(255) UNIQUE NOT NULL,
    logo_url            TEXT,
    plan                plan_type DEFAULT 'free',
    status              status_type DEFAULT 'active',
    business_mode       business_mode DEFAULT 'online',
    number_of_employees INTEGER,
    is_enabled          BOOLEAN DEFAULT true,
    trial_end_at        TIMESTAMPTZ,
    created_at          TIMESTAMPTZ DEFAULT now(),
    updated_at          TIMESTAMPTZ DEFAULT now(),
    deleted_at          TIMESTAMPTZ
);


---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.

DROP TYPE IF EXISTS status_type;
DROP TYPE IF EXISTS plan_type;
DROP TYPE IF EXISTS business_mode;
DROP TABLE IF EXISTS tenants;