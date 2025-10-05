package role


type Permission struct {
	ID   int
	Code string
	Name string
}

type Role struct {
	ID          int
	Code        string
	Name        string
	Permissions []Permission
}

// RolePermissionsSeed maps roles to their assigned permissions
var RolePermissionsSeed = map[string][]string{
    "system_admin": {
        "system:tenant:create",
        "system:tenant:read",
        "system:user:manage",
    },

    "admin": {
        "system:tenant:read",
    },

    "store_owner": {
        "store:manage",
        "product:write",
        "product:read",
        "order:write",
        "order:read",
        "payment:write",
        "payment:read",
        "report:read",
    },

    "store_admin": {
        "product:write",
        "product:read",
        "order:write",
        "order:read",
        "payment:write",
        "payment:read",
        "report:read",
    },

    "staff": {
        "product:write",
        "product:read",
        "order:write",
        "order:read",
        "payment:read",
    },

    "customer": {
        "product:read",
        "order:read",
        "order:write",   // cancel included
        "payment:read",
    },
}


// system_admin = root operator of your SaaS (rare, trusted accounts).

// admin = platform admin with visibility and support powers, but not root.

// store_owner = tenant root (for their own org).

// store_admin/staff/customer = scoped to their store.