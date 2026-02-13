package catalog

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/iamonah/merchcore/internal/domain/types/money"
)

type ProductID uuid.UUID

type Product struct {
	ID          ProductID
	TenantID    uuid.UUID
	Name        string
	Description string
	Price       money.Money
	Active      bool
	Images      []string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewProduct(tenantID uuid.UUID, name, description string, price money.Money) (*Product, error) {
	if name == "" {
		return nil, errors.New("product name cannot be empty")
	}
	if price.Amount.IsNegative() {
		return nil, errors.New("product price cannot be negative")
	}

	return &Product{
		ID:          ProductID(uuid.New()),
		TenantID:    tenantID,
		Name:        name,
		Description: description,
		Price:       price,
		Active:      true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Images:      make([]string, 0),
	}, nil
}
