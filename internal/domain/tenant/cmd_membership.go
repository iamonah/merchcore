package tenant

import (
	"time"

	"github.com/IamOnah/storefronthq/internal/domain/shared/role"

	"github.com/google/uuid"
)

type Membership struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	OrganizationID uuid.UUID
	Role           role.Role // enum: Owner, Admin, Manager, Staff, Customer
	CreatedAt      time.Time
}
