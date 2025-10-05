package tenant

import (
	"fmt"
	"time"

	"github.com/IamOnah/storefronthq/internal/domain/shared/address"

	"github.com/google/uuid"
)

type TenantProfile struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	Domain       string // custom URL like victorelectronics.shop
	BusinessName string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	LogoURL      *string
	Address      []address.Address
}

// client pays for store and a store is created
func NewTenant(tenantData CreateTenant) (TenantProfile, error) {
	tenantID := uuid.New()

	userID, err := uuid.Parse(tenantData.UserID)
	if err != nil {
		return TenantProfile{}, fmt.Errorf("invalid tenant id: %w", err)
	}

	return TenantProfile{
		ID:           tenantID,
		UserID:       userID,
		Domain:       tenantData.Domain,
		BusinessName: tenantData.BusinessName,
		UpdatedAt:    time.Now(),
		CreatedAt:    time.Now(),
	}, nil
}

func (tp *TenantProfile) GetUserID() uuid.UUID {
	return tp.UserID
}

func (tp *TenantProfile) GetBusinessName() string {
	return tp.BusinessName
}

func (tp *TenantProfile) GetTenantDomain() string {
	return tp.Domain
}
