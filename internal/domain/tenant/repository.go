package tenant

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrDomain           = errors.New("domain already exists")
	ErrSubDomain        = errors.New("subdomain already exists")
	ErrInvalidUser      = errors.New("invalid or non-existent user")
	ErrInvalidEnumValue = errors.New("invalid value for enum field")
	ErrDatabase         = errors.New("database error")
)

type TenantRepository interface {
	CreateTenant(ctx context.Context, tenant *TenantProfile) error
	CreateTenantSchema(ctx context.Context, userID uuid.UUID) error
	CheckSubdomainAvailability(ctx context.Context, subdomain string) (bool, error)
}
