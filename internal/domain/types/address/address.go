package address

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type AddressType string

var (
	AddressTypeBusiness = newAddressTypes("business")
	AddressTypeBilling  = newAddressTypes("billing")
	AddressTypeShipping = newAddressTypes("shipping")
)

var addressTypes = make(map[string]AddressType)

func newAddressTypes(v string) AddressType {
	at := AddressType(v)
	addressTypes[strings.ToLower(v)] = at
	return at
}

func ParseAddress(a string) (AddressType, error) {
	v, ok := addressTypes[strings.ToLower(a)]
	if !ok {
		return "", fmt.Errorf("invalid addresstype: %v", a)
	}
	return v, nil
}

type Address struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	Street     string
	City       string
	State      string
	PostalCode string
	Country    string
	IsDefault  bool
	Type       AddressType
}

func NewAddress(userID uuid.UUID, street, city, state, postalCode, country string, isDefault bool,
	addressType AddressType,
) *Address {
	return &Address{
		ID:         uuid.New(),
		UserID:     userID,
		Street:     street,
		City:       city,
		State:      state,
		PostalCode: postalCode,
		Country:    country,
		IsDefault:  isDefault,
		Type:       addressType,
	}
}
func (ad Address) GetUserID() uuid.UUID {
	return ad.UserID
}

func (ad Address) GetPostalCode() string {
	return ad.PostalCode
}

func (ad Address) GetStreet() string {
	return ad.Street
}

type Addresses []Address

func (u Addresses) AddAddresses(addr Address) Addresses {
	u = append(u, addr)
	return u
}

func (u Addresses) UpdateAddress(old Address, new Address) {
	for i, v := range u {
		if v.ID == old.ID {
			u[i] = new
		}
	}
}

func (u Addresses) GetDefaultAddress() Address {
	for _, v := range u {
		if v.IsDefault {
			return v
		}
	}
	return Address{}
}

func (u Addresses) GetDefaultShippigAddress() Address {
	for _, v := range u {
		if v.IsDefault {
			if v.Type == AddressTypeShipping {
				return v
			}
		}
	}
	return Address{}
}
