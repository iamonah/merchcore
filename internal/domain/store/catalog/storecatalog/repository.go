package catalog

import (
	"context"

	"github.com/google/uuid"
)

type Filter struct {
	Limit  int
	Offset int
}

type Repository interface {
	CreateProduct(ctx context.Context, p Product) error
	GetProduct(ctx context.Context, tenantID uuid.UUID, productID ProductID) (*Product, error)
	ListProducts(ctx context.Context, tenantID uuid.UUID, filter Filter) ([]Product, error)
}
