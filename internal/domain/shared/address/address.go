package address

import "github.com/google/uuid"

type Address struct {
	UserID     uuid.UUID
	ID         string
	Street     string
	City       string
	State      string
	PostalCode string
	Country    string
	IsDefault  bool
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

func (u Addresses) GetDefaultAdress() Address {
	for _, v := range u {
		if v.IsDefault {
			return v
		}
	}
	return Address{}
}

