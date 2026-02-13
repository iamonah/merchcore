package tenantdb

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/iamonah/merchcore/internal/domain/tenant"
	"github.com/iamonah/merchcore/internal/infra/database"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"
)

type tenantStore struct {
	conn database.DBTX
}

var _ tenant.TenantRepository = (*tenantStore)(nil)

func NewTenantStore(conn *pgxpool.Pool) *tenantStore {
	return &tenantStore{
		conn: conn,
	}
}

func (t *tenantStore) CreateTenant(ctx context.Context, te *tenant.TenantProfile) error {
	conn := database.GetTXFromContext(ctx, t.conn)

	query := `
		INSERT INTO tenants (
		id, user_id, business_name, domain, subdomain, logo_url, plan, status, business_mode, 
		number_of_employees, trial_end_at, trial_start_at, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`
	_, err := conn.Exec(
		ctx,
		query,
		te.ID,
		te.UserID,
		te.BusinessName,
		te.Domain,
		te.Subdomain,
		te.LogoURL,
		te.Plan,
		te.Status,
		te.BusinessMode,
		te.NumberOfEmployees,
		te.TrialEndAt,
		te.TrialStartAt,
		te.CreatedAt,
		te.UpdatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.ConstraintName {
			case "domain_uq":
				return tenant.ErrDomain
			case "subdomain_uq":
				return tenant.ErrSubDomain
			}

			switch pgErr.Code {
			case "22P02":
				return fmt.Errorf("%w:%w", tenant.ErrDatabase, tenant.ErrInvalidEnumValue) // invalid enum input
			}
		}
		return fmt.Errorf("%w: %w", tenant.ErrDatabase, err)
	}
	return nil
}

func (t *tenantStore) CreateTenantSchema(ctx context.Context, userID uuid.UUID) error {
	conn := database.GetTXFromContext(ctx, t.conn)

	schemaName := fmt.Sprintf("tenant_%s", userID)
	quotedSchema := pq.QuoteIdentifier(schemaName)

	createSchemaSQL := fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s;", quotedSchema)
	if _, err := conn.Exec(ctx, createSchemaSQL); err != nil {
		return fmt.Errorf("create schema: %w", err)
	}

	const placeholderCount = 53

	args := make([]any, placeholderCount)
	for i := range args {
		args[i] = quotedSchema
	}

	ddl := fmt.Sprintf(tenantInitSQL, args...)

	if _, err := conn.Exec(ctx, ddl); err != nil {
		return fmt.Errorf("initialize schema objects for tenant %s: %w", schemaName, err)
	}

	return nil
}

func (t *tenantStore) CheckSubdomainAvailability(ctx context.Context, subdomain string) (bool, error) {
	query := `
		SELECT EXISTS (SELECT 1 FROM tenants WHERE subdomain = $1);
	`
	var exist bool
	err := t.conn.QueryRow(ctx, query, subdomain).Scan(&exist)
	if err != nil {
		return false, fmt.Errorf("%w: %w", tenant.ErrDatabase, err)
	}
	return exist, nil
}
