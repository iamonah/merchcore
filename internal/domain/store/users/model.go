package customer

import "github.com/google/uuid"

type CustomerProfileCreate struct {
	FirstName string
	LastName  string
	OrgID     uuid.UUID

	Addresses []string
}

type CustomerProfileUpdate struct {
	FirstName *string
	LastName  *string

	Addresses []string
}
