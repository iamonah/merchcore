package tenant

import "github.com/IamOnah/storefronthq/internal/domain/shared/address"

type CreateTenant struct {
	UserID       string
	BusinessName string
	Domain       string
	Address      []address.Address
}
