package tenant

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/iamonah/merchcore/internal/domain/shared/address"
)

type CreateTenant struct {
	UserID           uuid.UUID
	BusinessName     string
	Description      string
	Subdomain        *string
	Domain           *string
	LogoURL          *string
	BusinessMode     string
	BusinessCategory string
	Plan             string
	BusinessAddress  *AddressInput
	BillingAddress   *AddressInput
}

// Primitive form used only for JSON input (e.g. during tenant creation)
type AddressInput struct {
	UserID     uuid.UUID
	Street     string
	City       string
	State      string
	PostalCode string
	Country    string
	IsDefault  bool
	Type       string
}

func NewAddress(p AddressInput) (*address.Address, error) {
	addtype, err := address.ParseAddress(p.Type)
	if err != nil {
		return nil, err
	}
	if addtype == address.AddressTypeBilling || addtype == address.AddressTypeBusiness {
		return nil, fmt.Errorf("invalid store type (business or billing)")
	}
	return &address.Address{
		ID:         uuid.New(),
		UserID:     p.UserID,
		Street:     p.Street,
		City:       p.City,
		State:      p.State,
		PostalCode: p.PostalCode,
		Country:    p.Country,
		IsDefault:  p.IsDefault,
		Type:       addtype,
	}, nil
}
