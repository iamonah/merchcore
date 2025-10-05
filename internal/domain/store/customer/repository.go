package customer

import (
	"errors"
)

var (
	ErrCustomerNotFound    = errors.New("customer not found")
	ErrFailedToAddCustomer = errors.New("customer not created")
	ErrUpdateCustomer      = errors.New("customer profile not updated")
)

// CustomerRepository is a interface that defines the rules around what a customer repository
// Has to be able to perform
type CustomerRepository interface {
	// Get(uuid.UUID) (aggregate.Customer, error)
	// Add(aggregate.Customer) error
	// Update(aggregate.Customer) error
}
