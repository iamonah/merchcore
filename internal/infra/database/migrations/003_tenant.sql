CREATE TYPE status_type AS ENUM ('active', 'suspended', 'deleted');
CREATE TYPE plan_type AS ENUM ('free', 'pro');

CREATE TABLE tenant (
    tenant_id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(user_id),
    business_name VARCHAR(255),
    domain VARCHAR(255) UNIQUE NOT NULL,
    logo_url TEXT,
    plan status_type DEFAULT 'free',
    status status_type DEFAULT 'active', 
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    number_of_employees SERIAL,
    is_enabled      BOOLEAN,
    trial_end_at TIMESTAMPTZ
);
---- create above / drop below ----

