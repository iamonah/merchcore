package customer

import (
	"errors"
	"strings"

	"github.com/iamonah/merchcore/internal/domain/types/address"

	"github.com/google/uuid"
)

type CustomerProfile struct {
	FirstName  string
	LastName   string
	CustomerID uuid.UUID

	OrgID uuid.UUID

	Addresses     []address.Address
	LoyaltyPoints int
}

func generateUserID() uuid.UUID {
	return uuid.New()
}

func NewCustomerProfile(cusProfile CustomerProfileCreate) (CustomerProfile, error) {
	customerID := generateUserID()
	cleanFirstName := strings.TrimSpace(cusProfile.FirstName)
	if cleanFirstName == "" {
		return CustomerProfile{}, errors.New("last-name cannot be emtpy")
	}
	cleanLastName := strings.TrimSpace(cusProfile.LastName)
	if cleanLastName == "" {
		return CustomerProfile{}, errors.New("first-name cannot be empty")
	}

	customerProfile := CustomerProfile{
		CustomerID: customerID,
		FirstName:  cleanFirstName,
		LastName:   cleanLastName,
		OrgID:      cusProfile.OrgID,
	}

	return customerProfile, nil
}
