CREATE TABLE stores (
    id         BIGSERIAL PRIMARY KEY,
    store_id   UUID UNIQUE NOT NULL,
    tenant_id  UUID REFERENCES tenants(id),
    name       TEXT NOT NULL,
    domain     TEXT UNIQUE NOT NULL,
    logo_url   TEXT,
    created_at TIMESTAMPTZ DEFAULT now()
);
---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
