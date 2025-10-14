CREATE TABLE IF NOT EXISTS tenant_customers (
    id         BIGSERIAL PRIMARY KEY,
    tenant_id  UUID NOT NULL REFERENCES tenants(tenant_id) ON DELETE CASCADE,
    user_id    UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    joined_at  TIMESTAMPTZ DEFAULT now(),
    UNIQUE (tenant_id, user_id)
);

CREATE INDEX IF NOT EXISTS tenant_customers_tenant_id_idx ON tenant_customers(tenant_id);
CREATE INDEX IF NOT EXISTS tenant_customers_user_id_idx ON tenant_customers(user_id);

---- create above / drop below ----

DROP INDEX IF EXISTS tenant_customers_tenant_id_idx;
DROP INDEX IF EXISTS tenant_customers_user_id_idx;
DROP TABLE IF EXISTS tenant_customers;