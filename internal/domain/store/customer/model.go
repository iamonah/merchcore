package customer

import "github.com/google/uuid"

type CustomerProfileCreate struct {
	FirstName string
	LastName  string
	OrgID     uuid.UUID
	// Phone         PhoneNumber
	Addresses []string
}

type CustomerProfileUpdate struct {
	FirstName *string
	LastName  *string
	// Phone         PhoneNumber
	Addresses []string
}
