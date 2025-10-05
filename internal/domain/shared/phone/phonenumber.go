package phone

import (
	"errors"
	"fmt"

	"github.com/nyaruka/phonenumbers"
)

type PhoneNumber struct {
	Number  string
	Country string
}

func NewPhoneNumber(num, country string) PhoneNumber {
	return PhoneNumber{
		Number:  num,
		Country: country,
	}
}
func (ph *PhoneNumber) ValidateNumber() error {
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

func (ph *PhoneNumber) GetPhoneNumber() string {
	return ph.Number
}

func (ph *PhoneNumber) GetCountry() string {
	return ph.Country
}
