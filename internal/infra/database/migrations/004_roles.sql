CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT uuid7(),
    name TEXT UNIQUE NOT NULL,
    description TEXT
);

CREATE TABLE permissions (
    id UUID PRIMARY KEY DEFAULT uuid7(),
    code TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL
);

INSERT INTO permissions (code, name) VALUES
-- System-level
('system:tenant:create', 'Create a new tenant'),
('system:tenant:read',   'View tenants'),
('system:user:manage',   'Assign global roles'),

-- Store-level
('store:manage',          'Manage store (for store owner)'),
('product:write',         'Create/update/delete products'),
('product:read',          'View products'),
('order:write',           'Update/cancel orders'),
('order:read',            'View orders'),
('payment:write',         'Process payments'),
('payment:read',          'View payments'),
('report:read',           'View reports');

-- Insert Roles
INSERT INTO roles (name, description) VALUES
('system_admin', 'Full platform-level access'),
('admin',        'Limited platform-level powers'),
('store_owner',  'Full control of a tenant store'),
('store_admin',  'Delegated admin inside a store'),
('staff',        'Operational staff in a store'),
('customer',     'Read-only access for customers');

CREATE TABLE role_permissions (
    role_id UUID REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);

-- system_admin
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'system_admin'
  AND p.code IN (
    'system:tenant:create',
    'system:tenant:read',
    'system:user:manage'
  );

-- admin
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'admin'
  AND p.code IN (
  'system:tenant:read',
  'report:read'
);

-- store_owner
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'store_owner'
  AND p.code IN (
    'store:manage',
    'product:write',
    'product:read',
    'order:write',
    'order:read',
    'payment:write',
    'payment:read',
    'report:read'
  );

-- store_admin
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'store_admin'
  AND p.code IN (
    'product:write',
    'product:read',
    'order:write',
    'order:read',
    'payment:write',
    'payment:read',
    'report:read'
  );

-- staff
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'staff'
  AND p.code IN (
    'product:write',
    'product:read',
    'order:write',
    'order:read',
    'payment:read'
  );

-- customer
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'customer'
  AND p.code IN (
    'product:read',
    'order:read',
    'order:write',   -- includes cancel
    'payment:read'
  );


---- create above / drop below ----
