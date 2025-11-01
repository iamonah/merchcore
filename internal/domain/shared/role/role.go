package role

import "fmt"

var (
	SystemAdmin = newRole("system_admin")
	Admin       = newRole("admin")
	StoreOwner  = newRole("store_owner")
	Staff       = newRole("staff")
	Guest       = newRole("guest")
)

var roles = make(map[string]Role)

type Role struct {
	value string
}

func newRole(v string) Role {
	r := Role{v}
	roles[v] = r
	return r
}

func (r Role) String() string { return r.value }

func Parse(v string) (Role, error) {
	r, ok := roles[v]
	if !ok {
		return Role{}, fmt.Errorf("invalid role: %q", v)
	}
	return r, nil
}
