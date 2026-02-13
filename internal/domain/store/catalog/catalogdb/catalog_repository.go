package catalogdb

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	catalog "github.com/iamonah/merchcore/internal/domain/store/catalog/storecatalog"
	"github.com/iamonah/merchcore/internal/domain/types/money"
	"github.com/iamonah/merchcore/internal/infra/database"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type CatalogDB struct {
	conn database.DBTX
}

func NewCatalogDB(conn *pgxpool.Pool) *CatalogDB {
	return &CatalogDB{
		conn: conn,
	}
}

func (db *CatalogDB) CreateProduct(ctx context.Context, p catalog.Product) error {
	conn := database.GetTXFromContext(ctx, db.conn)

	query := `
		INSERT INTO products (
			id, tenant_id, name, description, 
			price_amount, price_currency, active, images, 
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := conn.Exec(ctx, query,
		p.ID,
		p.TenantID,
		p.Name,
		p.Description,
		p.Price.Amount,
		p.Price.Currency,
		p.Active,
		p.Images,
		p.CreatedAt,
		p.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to insert product: %w", err)
	}

	return nil
}

func (db *CatalogDB) GetProduct(ctx context.Context, tenantID, productID catalog.ProductID) (*catalog.Product, error) {
	conn := database.GetTXFromContext(ctx, db.conn)

	query := `
		SELECT 
			id, tenant_id, name, description, 
			price_amount, price_currency, active, images, 
			created_at, updated_at
		FROM products
		WHERE id = $1 AND tenant_id = $2
	`

	var p catalog.Product
	var priceAmount decimal.Decimal
	var priceCurrency string
	var id uuid.UUID

	err := conn.QueryRow(ctx, query, uuid.UUID(productID), tenantID).Scan(
		&id,
		&p.TenantID,
		&p.Name,
		&p.Description,
		&priceAmount,
		&priceCurrency,
		&p.Active,
		&p.Images,
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("product not found")
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	p.ID = catalog.ProductID(id)
	p.Price = money.New(priceAmount, money.Currency(priceCurrency))

	return &p, nil
}

func (db *CatalogDB) ListProducts(ctx context.Context, tenantID uuid.UUID, filter catalog.Filter) ([]catalog.Product, error) {
	conn := database.GetTXFromContext(ctx, db.conn)

	query := `
		SELECT 
			id, tenant_id, name, description, 
			price_amount, price_currency, active, images, 
			created_at, updated_at
		FROM products
		WHERE tenant_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	limit := filter.Limit
	if limit == 0 {
		limit = 10
	}

	rows, err := conn.Query(ctx, query, tenantID, limit, filter.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}
	defer rows.Close()

	var products []catalog.Product
	for rows.Next() {
		var p catalog.Product
		var priceAmount decimal.Decimal
		var priceCurrency string
		var id uuid.UUID

		if err := rows.Scan(
			&id,
			&p.TenantID,
			&p.Name,
			&p.Description,
			&priceAmount,
			&priceCurrency,
			&p.Active,
			&p.Images,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}

		p.ID = catalog.ProductID(id)
		p.Price = money.New(priceAmount, money.Currency(priceCurrency))
		products = append(products, p)
	}

	return products, nil
}
