CREATE TABLE IF NOT EXISTS roles (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    name TEXT UNIQUE NOT NULL,
    description TEXT
);

-- Insert the fixed roles
INSERT INTO roles (name, description) VALUES
('system_admin', 'Full platform-level access'),
('admin',        'Limited platform-level powers'),
('store_owner',  'Owner of a store with full control'),
('store_admin',  'Store admin with delegated powers'),
('guest',        'Unregistered or new user');

---- create above / drop below ----
DROP TABLE IF EXISTS roles;
