package contact

import (
	"errors"
	"fmt"

	"github.com/nyaruka/phonenumbers"
)

type Contact struct {
	Number  string
	Country string
}

func NewContact(num, country string) Contact {
	return Contact{
		Number:  num,
		Country: country,
	}
}

func (ph *Contact) ValidateContact() error {
	num, err := phonenumbers.Parse(ph.Number, ph.Country)
	if err != nil {
		return fmt.Errorf("unable to parse %s for region %s: %w", ph.Number, ph.Country, err)
	}

	if !phonenumbers.IsPossibleNumber(num) {
		return errors.New("not a valid phone-number")
	}

	if !phonenumbers.IsValidNumberForRegion(num, ph.Country) {
		return fmt.Errorf("not a valid phone-number for this country: %s", ph.Country)

	}
	//E.164 format:
	ph.Number = phonenumbers.Format(num, phonenumbers.E164)
	return nil
}

func (ph *Contact) GetPhoneNumber() string {
	return ph.Number
}

func (ph *Contact) GetCountry() string {
	return ph.Country
}
