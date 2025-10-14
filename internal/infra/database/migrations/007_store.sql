CREATE TABLE IF NOT EXISTS stores (
    id         BIGSERIAL PRIMARY KEY,
    store_id   UUID UNIQUE NOT NULL DEFAULT uuidv7(),
    tenant_id  UUID REFERENCES tenants(tenant_id),
    name       TEXT NOT NULL,
    domain     TEXT UNIQUE NOT NULL,
    logo_url   TEXT,
    created_at TIMESTAMPTZ DEFAULT now()
);

---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.

DROP TABLE IF EXISTS stores;