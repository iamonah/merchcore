package money

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type Currency string

const (
	USD Currency = "USD"
	EUR Currency = "EUR"
)

type Money struct {
	Amount   decimal.Decimal
	Currency Currency
}

func New(amount decimal.Decimal, currency Currency) Money {
	return Money{
		Amount:   amount,
		Currency: currency,
	}
}

func (m Money) String() string {
	return fmt.Sprintf("%s %s", m.Amount.StringFixed(2), m.Currency)
}
