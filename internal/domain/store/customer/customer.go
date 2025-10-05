package customer

import (
	"errors"
	"strings"

	"github.com/IamOnah/storefronthq/internal/domain/shared/address"

	"github.com/google/uuid"
)

type CustomerProfile struct {
	CustomerID uuid.UUID
	// UserID     uuid.UUID // global user identity
	OrgID     uuid.UUID // which store they belong to
	FirstName string
	LastName  string
	// Phone         PhoneNumber
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
